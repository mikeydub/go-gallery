package notifications

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/mikeydub/go-gallery/service/limiters"
	"github.com/mikeydub/go-gallery/service/redis"
	sentryutil "github.com/mikeydub/go-gallery/service/sentry"
	"github.com/mikeydub/go-gallery/service/task"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"cloud.google.com/go/pubsub"
	"github.com/bsm/redislock"
	"github.com/gin-gonic/gin"
	"github.com/googleapis/gax-go/v2/apierror"
	db "github.com/mikeydub/go-gallery/db/gen/coredb"
	"github.com/mikeydub/go-gallery/env"
	"github.com/mikeydub/go-gallery/service/logger"
	"github.com/mikeydub/go-gallery/service/persist"
	"github.com/mikeydub/go-gallery/util"
	"google.golang.org/grpc/codes"
)

type lockKey struct {
	ownerID persist.DBID
	action  persist.Action
}

const viewWindow = 24 * time.Hour
const groupedWindow = 10 * time.Minute
const notificationTimeout = 10 * time.Second
const maxLockTimeout = 2 * time.Minute
const NotificationHandlerContextKey = "notification.notificationHandlers"

type NotificationHandlers struct {
	Notifications            *notificationDispatcher
	UserNewNotifications     map[persist.DBID]chan db.Notification
	UserUpdatedNotifications map[persist.DBID]chan db.Notification
	pubSub                   *pubsub.Client
}

type pushLimiter struct {
	comments *limiters.KeyRateLimiter
	admires  *limiters.KeyRateLimiter
	follows  *limiters.KeyRateLimiter
	tokens   *limiters.KeyRateLimiter
}

func newPushLimiter() *pushLimiter {
	cache := redis.NewCache(redis.PushNotificationRateLimitersCache)
	ctx := context.Background()
	return &pushLimiter{
		comments: limiters.NewKeyRateLimiter(ctx, cache, "comments", 5, time.Minute),
		admires:  limiters.NewKeyRateLimiter(ctx, cache, "admires", 1, time.Minute*10),
		follows:  limiters.NewKeyRateLimiter(ctx, cache, "follows", 1, time.Minute*10),
		tokens:   limiters.NewKeyRateLimiter(ctx, cache, "tokens", 1, time.Minute*10),
	}
}

func (p *pushLimiter) tryComment(ctx context.Context, sendingUserID persist.DBID, receivingUserID persist.DBID, feedEventID persist.DBID) error {
	key := fmt.Sprintf("%s:%s:%s", sendingUserID.String(), receivingUserID.String(), feedEventID.String())
	if p.isActionAllowed(ctx, p.comments, key) {
		return nil
	}

	return errRateLimited{
		limiterName: p.comments.Name(),
		senderID:    sendingUserID,
		receiverID:  receivingUserID,
		feedEventID: feedEventID,
	}
}

func (p *pushLimiter) tryAdmire(ctx context.Context, sendingUserID persist.DBID, receivingUserID persist.DBID, feedEventID persist.DBID) error {
	key := fmt.Sprintf("%s:%s:%s", sendingUserID.String(), receivingUserID.String(), feedEventID.String())
	if p.isActionAllowed(ctx, p.admires, key) {
		return nil
	}

	return errRateLimited{
		limiterName: p.admires.Name(),
		senderID:    sendingUserID,
		receiverID:  receivingUserID,
		feedEventID: feedEventID,
	}
}

func (p *pushLimiter) tryFollow(ctx context.Context, sendingUserID persist.DBID, receivingUserID persist.DBID) error {
	key := fmt.Sprintf("%s:%s", sendingUserID, receivingUserID)
	if p.isActionAllowed(ctx, p.follows, key) {
		return nil
	}

	return errRateLimited{
		limiterName: p.follows.Name(),
		senderID:    sendingUserID,
		receiverID:  receivingUserID,
	}
}

func (p *pushLimiter) tryTokens(ctx context.Context, sendingUserID persist.DBID, tokenID persist.DBID) error {
	key := fmt.Sprintf("%s:%s", sendingUserID, tokenID)
	if p.isActionAllowed(ctx, p.tokens, key) {
		return nil
	}

	return errRateLimited{
		limiterName: p.follows.Name(),
		senderID:    sendingUserID,
	}
}

func (p *pushLimiter) isActionAllowed(ctx context.Context, limiter *limiters.KeyRateLimiter, key string) bool {
	canContinue, _, err := limiter.ForKey(ctx, key)
	if err != nil {
		logger.For(ctx).Warnf("error getting rate limit for key %s: %s", key, err.Error())
		return false
	}

	return canContinue
}

// New registers specific notification handlers
func New(queries *db.Queries, pub *pubsub.Client, taskClient *cloudtasks.Client, lock *redislock.Client) *NotificationHandlers {
	notifDispatcher := notificationDispatcher{handlers: map[persist.Action]notificationHandler{}, lock: lock}
	limiter := newPushLimiter()

	def := defaultNotificationHandler{queries: queries, pubSub: pub, taskClient: taskClient, limiter: limiter}
	group := groupedNotificationHandler{queries: queries, pubSub: pub, taskClient: taskClient, limiter: limiter}
	tokenIDGrouped := tokenIDGroupedNotificationHandler{queries: queries, pubSub: pub, taskClient: taskClient, limiter: limiter}
	view := viewedNotificationHandler{queries: queries, pubSub: pub, taskClient: taskClient, limiter: limiter}

	// grouped notification actions
	notifDispatcher.AddHandler(persist.ActionUserFollowedUsers, group)
	notifDispatcher.AddHandler(persist.ActionAdmiredFeedEvent, group)
	notifDispatcher.AddHandler(persist.ActionAdmiredPost, group)
	notifDispatcher.AddHandler(persist.ActionAdmiredToken, group)

	// single notification actions (default)
	notifDispatcher.AddHandler(persist.ActionCommentedOnFeedEvent, def)
	notifDispatcher.AddHandler(persist.ActionCommentedOnPost, def)

	// notification actions that are grouped by token id
	notifDispatcher.AddHandler(persist.ActionNewTokensReceived, tokenIDGrouped)

	// viewed notifications are handled separately
	notifDispatcher.AddHandler(persist.ActionViewedGallery, view)

	new := map[persist.DBID]chan db.Notification{}
	updated := map[persist.DBID]chan db.Notification{}

	notificationHandlers := &NotificationHandlers{Notifications: &notifDispatcher, UserNewNotifications: new, UserUpdatedNotifications: updated, pubSub: pub}
	if pub != nil {
		go notificationHandlers.receiveNewNotificationsFromPubSub()
		go notificationHandlers.receiveUpdatedNotificationsFromPubSub()
	} else {
		logger.For(nil).Warn("pubsub not configured, notifications will not be received")
	}
	return notificationHandlers
}

// Register specific notification handlers
func AddTo(ctx *gin.Context, notificationHandlers *NotificationHandlers) {
	ctx.Set(NotificationHandlerContextKey, notificationHandlers)
}

func For(ctx context.Context) *NotificationHandlers {
	gc := util.MustGetGinContext(ctx)
	return gc.Value(NotificationHandlerContextKey).(*NotificationHandlers)
}

func (n *NotificationHandlers) GetNewNotificationsForUser(userID persist.DBID) chan db.Notification {
	if sub, ok := n.UserNewNotifications[userID]; ok && sub != nil {
		logger.For(context.Background()).Infof("returning existing new notification channel for user: %s", userID)
		return sub
	}
	sub := make(chan db.Notification)
	n.UserNewNotifications[userID] = sub
	logger.For(context.Background()).Infof("created new new notification channel for user: %s", userID)
	return sub
}

func (n *NotificationHandlers) GetUpdatedNotificationsForUser(userID persist.DBID) chan db.Notification {
	if sub, ok := n.UserUpdatedNotifications[userID]; ok && sub != nil {
		logger.For(context.Background()).Infof("returning existing updated notification channel for user: %s", userID)
		return sub
	}
	sub := make(chan db.Notification)
	n.UserUpdatedNotifications[userID] = sub
	logger.For(context.Background()).Infof("created new updated notification channel for user: %s", userID)
	return sub
}

func (n *NotificationHandlers) UnsubscribeNewNotificationsForUser(userID persist.DBID) {
	logger.For(context.Background()).Infof("unsubscribing new notifications for user: %s", userID)
	delete(n.UserNewNotifications, userID)
}

func (n *NotificationHandlers) UnsubscribeUpdatedNotificationsForUser(userID persist.DBID) {
	logger.For(context.Background()).Infof("unsubscribing updated notifications for user: %s", userID)
	delete(n.UserUpdatedNotifications, userID)
}

type notificationHandler interface {
	Handle(context.Context, db.Notification) error
}

type notificationDispatcher struct {
	handlers map[persist.Action]notificationHandler
	lock     *redislock.Client
}

func (d *notificationDispatcher) AddHandler(action persist.Action, handler notificationHandler) {
	d.handlers[action] = handler
}

func (d *notificationDispatcher) Dispatch(ctx context.Context, notif db.Notification) error {
	if handler, ok := d.handlers[notif.Action]; ok {
		l, _ := d.lock.Obtain(ctx, lockKey{ownerID: notif.OwnerID, action: notif.Action}.String(), maxLockTimeout, &redislock.Options{RetryStrategy: redislock.LinearBackoff(5 * time.Second)})
		if l != nil {
			defer l.Release(ctx)
		}
		return handler.Handle(ctx, notif)
	}
	logger.For(ctx).Warnf("no handler registered for action: %s", notif.Action)
	return nil
}

type defaultNotificationHandler struct {
	queries    *db.Queries
	pubSub     *pubsub.Client
	taskClient *cloudtasks.Client
	limiter    *pushLimiter
}

func (h defaultNotificationHandler) Handle(ctx context.Context, notif db.Notification) error {
	return insertAndPublishNotif(ctx, notif, h.queries, h.pubSub, h.taskClient, h.limiter)
}

type groupedNotificationHandler struct {
	queries    *db.Queries
	pubSub     *pubsub.Client
	taskClient *cloudtasks.Client
	limiter    *pushLimiter
}

func (h groupedNotificationHandler) Handle(ctx context.Context, notif db.Notification) error {
	var curNotif db.Notification

	// Bucket notifications on the feed event if it has one
	onlyForFeed := notif.FeedEventID != ""
	onlyForPost := notif.PostID != ""
	curNotif, _ = h.queries.GetMostRecentNotificationByOwnerIDForAction(ctx, db.GetMostRecentNotificationByOwnerIDForActionParams{
		OwnerID:          notif.OwnerID,
		Action:           notif.Action,
		OnlyForFeedEvent: onlyForFeed,
		FeedEventID:      notif.FeedEventID,
		PostID:           notif.PostID,
		OnlyForPost:      onlyForPost,
	})

	if time.Since(curNotif.CreatedAt) < groupedWindow {
		logger.For(ctx).Infof("grouping notification %s: %s-%s", curNotif.ID, notif.Action, notif.OwnerID)
		return updateAndPublishNotif(ctx, notif, curNotif, h.queries, h.pubSub, h.taskClient, h.limiter)
	}
	logger.For(ctx).Infof("not grouping notification: %s-%s", notif.Action, notif.OwnerID)
	return insertAndPublishNotif(ctx, notif, h.queries, h.pubSub, h.taskClient, h.limiter)

}

type tokenIDGroupedNotificationHandler struct {
	queries    *db.Queries
	pubSub     *pubsub.Client
	taskClient *cloudtasks.Client
	limiter    *pushLimiter
}

func (h tokenIDGroupedNotificationHandler) Handle(ctx context.Context, notif db.Notification) error {
	var curNotif db.Notification

	// Bucket notifications on the feed event if it has one
	onlyForFeed := notif.FeedEventID != ""
	onlyForPost := notif.PostID != ""
	curNotif, _ = h.queries.GetMostRecentNotificationByOwnerIDTokenIDForAction(ctx, db.GetMostRecentNotificationByOwnerIDTokenIDForActionParams{
		OwnerID:          notif.OwnerID,
		Action:           notif.Action,
		TokenID:          notif.TokenID,
		OnlyForFeedEvent: onlyForFeed,
		FeedEventID:      notif.FeedEventID,
		PostID:           notif.PostID,
		OnlyForPost:      onlyForPost,
	})

	if time.Since(curNotif.CreatedAt) < groupedWindow {
		logger.For(ctx).Infof("grouping notification %s: %s-%s", curNotif.ID, notif.Action, notif.OwnerID)
		return updateAndPublishNotif(ctx, notif, curNotif, h.queries, h.pubSub, h.taskClient, h.limiter)
	}
	logger.For(ctx).Infof("not grouping notification: %s-%s", notif.Action, notif.OwnerID)
	return insertAndPublishNotif(ctx, notif, h.queries, h.pubSub, h.taskClient, h.limiter)

}

type viewedNotificationHandler struct {
	queries    *db.Queries
	pubSub     *pubsub.Client
	taskClient *cloudtasks.Client
	limiter    *pushLimiter
}

// will return the beginning of the week (sunday) in PST
func beginningOfWeek(t time.Time) time.Time {

	pst, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}

	y, m, d := t.In(pst).Date()

	newD := d - int(t.Weekday())

	return time.Date(y, m, newD, 0, 0, 0, 0, pst)
}

// this handler will still group notifications in the usual window, but it will also ensure that each viewer does
// does not show up mutliple times in a week
func (h viewedNotificationHandler) Handle(ctx context.Context, notif db.Notification) error {
	// all of this user's view notifications in the current week
	notifs, _ := h.queries.GetNotificationsByOwnerIDForActionAfter(ctx, db.GetNotificationsByOwnerIDForActionAfterParams{
		OwnerID:      notif.OwnerID,
		Action:       notif.Action,
		CreatedAfter: beginningOfWeek(time.Now()),
	})
	if len(notifs) == 0 {
		// if there are no notifications this week, then we definitely are going to insert this one
		logger.For(ctx).Debugf("no notifications this week, inserting: %s-%s", notif.Action, notif.OwnerID)
		return insertAndPublishNotif(ctx, notif, h.queries, h.pubSub, h.taskClient, h.limiter)
	}

	mostRecentNotif := notifs[0]

	if notif.Data.UnauthedViewerIDs != nil && len(notif.Data.UnauthedViewerIDs) > 0 {

		resultIDs := []string{}
		// add each of the unauthed viewer ids in the passed in notif to the map unless it is already in one of the notifications this week
		for _, id := range notif.Data.UnauthedViewerIDs {
			add := true
		firstInner:
			for _, n := range notifs {
				if util.ContainsString(n.Data.UnauthedViewerIDs, id) {
					add = false
					break firstInner
				}
			}
			if add {
				resultIDs = append(resultIDs, id)
			}
		}

		notif.Data.UnauthedViewerIDs = resultIDs
	}

	if notif.Data.AuthedViewerIDs != nil && len(notif.Data.AuthedViewerIDs) > 0 {
		// go through each of the authed viewer ids in the passed in notif and add them to the map unless they are already in one of the notifications this week
		resultIDs := []persist.DBID{}
		for _, id := range notif.Data.AuthedViewerIDs {
			add := true
		secondInner:
			for _, n := range notifs {
				if persist.ContainsDBID(n.Data.AuthedViewerIDs, id) {
					add = false
					break secondInner
				}
			}
			if add {
				resultIDs = append(resultIDs, id)
			}
		}

		notif.Data.AuthedViewerIDs = resultIDs
	}

	// if the most recent notification in the last week is within the grouping window then we will update it, if not, insert it
	if time.Since(mostRecentNotif.CreatedAt) < viewWindow {
		logger.For(ctx).Debugf("grouping notification %s: %s-%s", mostRecentNotif.ID, notif.Action, notif.OwnerID)
		return updateAndPublishNotif(ctx, notif, mostRecentNotif, h.queries, h.pubSub, h.taskClient, h.limiter)
	}
	logger.For(ctx).Debugf("not grouping notification: %s-%s", notif.Action, notif.OwnerID)
	return insertAndPublishNotif(ctx, notif, h.queries, h.pubSub, h.taskClient, h.limiter)
}

// subscribe returns a subscription to the given topic
func (n *NotificationHandlers) subscribe(ctx context.Context, topic, name string) (*pubsub.Subscription, error) {
	sub, err := createSubscription(ctx, n.pubSub, topic, name)
	if err == nil {
		return sub, nil
	}

	if errTopicMissing(err) {
		if _, err := n.pubSub.CreateTopic(ctx, topic); err != nil {
			return nil, err
		}
	}

	return createSubscription(ctx, n.pubSub, topic, name)
}

func (n *NotificationHandlers) receiveNewNotificationsFromPubSub() {
	sub, err := n.subscribe(context.Background(), env.GetString("PUBSUB_TOPIC_NEW_NOTIFICATIONS"), fmt.Sprintf("new-notifications-%s", persist.GenerateID()))
	if err != nil {
		logger.For(nil).Errorf("error creating new notifications subscription: %s", err)
		panic(err)
	}

	logger.For(nil).Info("subscribing to new notifications pubsub topic")

	err = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {

		logger.For(ctx).Debugf("received new notification from pubsub: %s", string(msg.Data))

		defer msg.Ack()
		notif := db.Notification{}
		err := json.Unmarshal(msg.Data, &notif)
		if err != nil {
			logger.For(ctx).Warnf("failed to unmarshal pubsub message: %s", err)
			return
		}

		logger.For(ctx).Infof("received new notification from pubsub: %s", notif.OwnerID)

		if sub, ok := n.UserNewNotifications[notif.OwnerID]; ok && sub != nil {
			select {
			case sub <- notif:
				logger.For(ctx).Debugf("sent new notification to user: %s", notif.OwnerID)
			case <-time.After(notificationTimeout):
				logger.For(ctx).Debugf("notification create channel not open for user: %s", notif.OwnerID)
				n.UnsubscribeNewNotificationsForUser(notif.OwnerID)
			}
		} else {
			logger.For(ctx).Debugf("no notification create channel open for user: %s", notif.OwnerID)
		}
	})
	if err != nil {
		logger.For(nil).Errorf("error receiving new notifications from pubsub: %s", err)
		panic(err)
	}
}

func (n *NotificationHandlers) receiveUpdatedNotificationsFromPubSub() {
	sub, err := n.subscribe(context.Background(), env.GetString("PUBSUB_TOPIC_UPDATED_NOTIFICATIONS"), fmt.Sprintf("updated-notifications-%s", persist.GenerateID()))
	if err != nil {
		logger.For(nil).Errorf("error creating updated notifications subscription: %s", err)
		panic(err)
	}

	logger.For(nil).Infof("subscribed to updated notifications pubsub")

	err = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {

		logger.For(ctx).Debugf("received updated notification from pubsub: %s", string(msg.Data))

		defer msg.Ack()
		notif := db.Notification{}
		err := json.Unmarshal(msg.Data, &notif)
		if err != nil {
			logger.For(ctx).Warnf("failed to unmarshal pubsub message: %s", err)
			return
		}

		logger.For(ctx).Infof("received updated notification from pubsub: %s", notif.OwnerID)

		if sub, ok := n.UserUpdatedNotifications[notif.OwnerID]; ok && sub != nil {
			select {
			case sub <- notif:
				logger.For(ctx).Debugf("sent updated notification to user: %s", notif.OwnerID)
			case <-time.After(notificationTimeout):
				logger.For(ctx).Debugf("notification update channel not open for user: %s", notif.OwnerID)
				n.UnsubscribeUpdatedNotificationsForUser(notif.OwnerID)
			}
		} else {
			logger.For(ctx).Debugf("no notification update channel open for user: %s", notif.OwnerID)
		}
	})
	if err != nil {
		logger.For(nil).Errorf("error receiving new notifications from pubsub: %s", err)
		panic(err)
	}
}

func createPushMessage(ctx context.Context, notif db.Notification, queries *db.Queries, limiter *pushLimiter) (task.PushNotificationMessage, error) {
	badgeCount, badgeErr := queries.CountUserUnseenNotifications(ctx, notif.OwnerID)
	if badgeErr != nil {
		return task.PushNotificationMessage{}, badgeErr
	}

	message := task.PushNotificationMessage{
		Title: "Gallery",
		Sound: true,
		Badge: int(badgeCount),
		Data: map[string]any{
			"action":          notif.Action,
			"notification_id": notif.ID,
		},
	}

	if notif.Action == persist.ActionAdmiredFeedEvent || notif.Action == persist.ActionAdmiredPost {
		admirer, err := queries.GetUserById(ctx, notif.Data.AdmirerIDs[0])
		if err != nil {
			return task.PushNotificationMessage{}, err
		}

		if err = limiter.tryAdmire(ctx, admirer.ID, notif.OwnerID, notif.FeedEventID); err != nil {
			return task.PushNotificationMessage{}, err
		}

		if !admirer.Username.Valid {
			return task.PushNotificationMessage{}, fmt.Errorf("user with ID=%s has no username", admirer.ID)
		}

		ac := "activity"
		if notif.Action == persist.ActionAdmiredPost {
			ac = "post"
		}

		message.Body = fmt.Sprintf("%s admired your %s", admirer.Username.String, ac)
		return message, nil
	}

	if notif.Action == persist.ActionCommentedOnFeedEvent || notif.Action == persist.ActionCommentedOnPost {
		comment, err := queries.GetCommentByCommentID(ctx, notif.CommentID)
		if err != nil {
			return task.PushNotificationMessage{}, err
		}

		commenter, err := queries.GetUserById(ctx, comment.ActorID)
		if err != nil {
			return task.PushNotificationMessage{}, err
		}

		if err = limiter.tryComment(ctx, commenter.ID, notif.OwnerID, notif.FeedEventID); err != nil {
			return task.PushNotificationMessage{}, err
		}

		if !commenter.Username.Valid {
			return task.PushNotificationMessage{}, fmt.Errorf("user with ID=%s has no username", commenter.ID)
		}

		ac := "activity"
		if notif.Action == persist.ActionCommentedOnPost {
			ac = "post"
		}

		message.Body = fmt.Sprintf("%s commented on your %s", commenter.Username.String, ac)
		return message, nil
	}

	if notif.Action == persist.ActionUserFollowedUsers {
		follower, err := queries.GetUserById(ctx, notif.Data.FollowerIDs[0])
		if err != nil {
			return task.PushNotificationMessage{}, err
		}

		if err = limiter.tryFollow(ctx, follower.ID, notif.OwnerID); err != nil {
			return task.PushNotificationMessage{}, err
		}

		if !follower.Username.Valid {
			return task.PushNotificationMessage{}, fmt.Errorf("user with ID=%s has no username", follower.ID)
		}

		var body string
		if notif.Data.FollowedBack {
			body = fmt.Sprintf("%s followed you back", follower.Username.String)
		} else {
			body = fmt.Sprintf("%s followed you", follower.Username.String)
		}

		message.Body = body
		return message, nil
	}
	if notif.Action == persist.ActionNewTokensReceived {

		tm, err := queries.GetTokenMediaByTokenId(ctx, notif.Data.NewTokenID)
		if err != nil {
			return task.PushNotificationMessage{}, err
		}
		name := tm.Name
		if name == "" {
			to, err := queries.GetTokenById(ctx, notif.Data.NewTokenID)
			if err != nil {
				return task.PushNotificationMessage{}, err
			}
			name = to.Name.String
		}

		name = util.TruncateWithEllipsis(name, 20)

		if err := limiter.tryTokens(ctx, notif.OwnerID, notif.Data.NewTokenID); err != nil {
			return task.PushNotificationMessage{}, err
		}
		amount := notif.Data.NewTokenQuantity
		i := amount.BigInt().Uint64()
		if i > 1 {
			message.Body = fmt.Sprintf("You received %d new %s tokens", i, name)
		} else {
			message.Body = fmt.Sprintf("You received a new %s token", name)
		}
	}

	return task.PushNotificationMessage{}, fmt.Errorf("unsupported notification action: %s", notif.Action)
}

func actionSupportsPushNotifications(action persist.Action) bool {
	switch action {
	case persist.ActionAdmiredFeedEvent:
		return true
	case persist.ActionCommentedOnFeedEvent:
		return true
	case persist.ActionUserFollowedUsers:
		return true
	case persist.ActionNewTokensReceived:
		return true
	case persist.ActionCommentedOnPost:
		return true
	case persist.ActionAdmiredPost:
		return true
	default:
		return false
	}
}

func sendPushNotifications(ctx context.Context, notif db.Notification, queries *db.Queries, taskClient *cloudtasks.Client, limiter *pushLimiter) error {
	if !actionSupportsPushNotifications(notif.Action) {
		return nil
	}

	pushTokens, err := queries.GetPushTokensByUserID(ctx, notif.OwnerID)
	if err != nil {
		return err
	}

	// Don't try to send anything if the user doesn't have any registered push tokens
	if len(pushTokens) == 0 {
		return nil
	}

	message, err := createPushMessage(ctx, notif, queries, limiter)
	if err != nil {
		if _, ok := err.(errRateLimited); ok {
			// Rate limiting is expected and shouldn't be propagated upward as an error
			logger.For(ctx).Infof("couldn't create push message: %s", err)
			return nil
		}

		return err
	}

	for _, token := range pushTokens {
		toSend := message
		toSend.PushTokenID = token.ID
		err = task.CreateTaskForPushNotification(ctx, toSend, taskClient)
		if err != nil {
			sentryutil.ReportError(ctx, err)
			logger.For(ctx).Errorf("failed to create task for push notification: %s", err)
		}
	}

	return nil
}

func insertAndPublishNotif(ctx context.Context, notif db.Notification, queries *db.Queries, ps *pubsub.Client, taskClient *cloudtasks.Client, limiter *pushLimiter) error {
	newNotif, err := addNotification(ctx, notif, queries)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	err = sendPushNotifications(ctx, notif, queries, taskClient, limiter)
	if err != nil {
		logger.For(ctx).Errorf("failed to send push notifications for notification with DBID=%s, error: %s", notif.ID, err)
	}

	marshalled, err := json.Marshal(newNotif)
	if err != nil {
		return err
	}
	t := ps.Topic(env.GetString("PUBSUB_TOPIC_NEW_NOTIFICATIONS"))
	result := t.Publish(ctx, &pubsub.Message{
		Data: marshalled,
	})

	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to publish new notification: %w", err)
	}

	logger.For(ctx).Infof("pushed new notification to pubsub: %s", notif.OwnerID)

	return nil
}

func updateAndPublishNotif(ctx context.Context, notif db.Notification, mostRecentNotif db.Notification, queries *db.Queries, ps *pubsub.Client, taskClient *cloudtasks.Client, limiter *pushLimiter) error {
	var amount = notif.Amount
	resultData := mostRecentNotif.Data.Concat(notif.Data)
	switch notif.Action {
	case persist.ActionAdmiredFeedEvent, persist.ActionAdmiredPost:
		amount = int32(len(resultData.AdmirerIDs))
	case persist.ActionViewedGallery:
		amount = int32(len(resultData.AuthedViewerIDs) + len(resultData.UnauthedViewerIDs))
	case persist.ActionUserFollowedUsers:
		amount = int32(len(resultData.FollowerIDs))
	case persist.ActionNewTokensReceived:
		amount = int32(resultData.NewTokenQuantity.BigInt().Uint64())
	default:
		amount = mostRecentNotif.Amount + notif.Amount
	}
	err := queries.UpdateNotification(ctx, db.UpdateNotificationParams{
		ID: mostRecentNotif.ID,
		// this concat will put the notif.Data values at the beginning of the array, sorted from most recently added to oldest added
		Data:   resultData,
		Amount: amount,
		// this will also get concatenated at the DB level
		EventIds: notif.EventIds,
	})
	if err != nil {
		return fmt.Errorf("error updating notification: %w", err)
	}

	err = sendPushNotifications(ctx, notif, queries, taskClient, limiter)
	if err != nil {
		logger.For(ctx).Errorf("failed to send push notifications for notification with DBID=%s, error: %s", notif.ID, err)
	}

	updatedNotif, err := queries.GetNotificationByID(ctx, mostRecentNotif.ID)
	if err != nil {
		return fmt.Errorf("error getting updated notification by %s: %w", mostRecentNotif.ID, err)
	}
	marshalled, err := json.Marshal(updatedNotif)
	if err != nil {
		return err
	}
	t := ps.Topic(env.GetString("PUBSUB_TOPIC_UPDATED_NOTIFICATIONS"))
	result := t.Publish(ctx, &pubsub.Message{
		Data: marshalled,
	})
	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("error publishing updated notification: %w", err)
	}

	logger.For(ctx).Infof("pushed updated notification to pubsub: %s", updatedNotif.OwnerID)
	return nil
}

func addNotification(ctx context.Context, notif db.Notification, queries *db.Queries) (db.Notification, error) {
	id := persist.GenerateID()
	switch notif.Action {
	case persist.ActionAdmiredFeedEvent:
		return queries.CreateAdmireNotification(ctx, db.CreateAdmireNotificationParams{
			ID:        id,
			OwnerID:   notif.OwnerID,
			Action:    notif.Action,
			Data:      notif.Data,
			EventIds:  notif.EventIds,
			FeedEvent: sql.NullString{String: notif.FeedEventID.String(), Valid: true},
		})
	case persist.ActionCommentedOnFeedEvent:
		return queries.CreateCommentNotification(ctx, db.CreateCommentNotificationParams{
			ID:        id,
			OwnerID:   notif.OwnerID,
			Action:    notif.Action,
			Data:      notif.Data,
			EventIds:  notif.EventIds,
			FeedEvent: sql.NullString{String: notif.FeedEventID.String(), Valid: true},
			CommentID: notif.CommentID,
		})
	case persist.ActionAdmiredPost:
		return queries.CreateAdmireNotification(ctx, db.CreateAdmireNotificationParams{
			ID:       id,
			OwnerID:  notif.OwnerID,
			Action:   notif.Action,
			Data:     notif.Data,
			EventIds: notif.EventIds,
			Post:     sql.NullString{String: notif.PostID.String(), Valid: true},
		})
	case persist.ActionCommentedOnPost:
		return queries.CreateCommentNotification(ctx, db.CreateCommentNotificationParams{
			ID:        id,
			OwnerID:   notif.OwnerID,
			Action:    notif.Action,
			Data:      notif.Data,
			EventIds:  notif.EventIds,
			Post:      sql.NullString{String: notif.PostID.String(), Valid: true},
			CommentID: notif.CommentID,
		})
	case persist.ActionUserFollowedUsers:
		return queries.CreateSimpleNotification(ctx, db.CreateSimpleNotificationParams{
			ID:       id,
			OwnerID:  notif.OwnerID,
			Action:   notif.Action,
			Data:     notif.Data,
			EventIds: notif.EventIds,
		})
	case persist.ActionViewedGallery:
		return queries.CreateViewGalleryNotification(ctx, db.CreateViewGalleryNotificationParams{
			ID:        id,
			OwnerID:   notif.OwnerID,
			Action:    notif.Action,
			Data:      notif.Data,
			EventIds:  notif.EventIds,
			GalleryID: notif.GalleryID,
		})
	case persist.ActionNewTokensReceived:
		amount := notif.Data.NewTokenQuantity.BigInt().Int64()
		return queries.CreateTokenNotification(ctx, db.CreateTokenNotificationParams{
			ID:       id,
			OwnerID:  notif.OwnerID,
			Action:   notif.Action,
			Data:     notif.Data,
			EventIds: notif.EventIds,
			TokenID:  notif.TokenID,
			Amount:   int32(amount),
		})
	default:
		return db.Notification{}, fmt.Errorf("unknown notification action: %s", notif.Action)
	}
}

func (l lockKey) String() string {
	return fmt.Sprintf("%s:%s", l.ownerID, l.action)
}

func errTopicMissing(err error) bool {
	var aErr *apierror.APIError
	if ok := errors.As(err, &aErr); ok && aErr.GRPCStatus().Code() == codes.NotFound {
		return true
	}
	return false
}

func createSubscription(ctx context.Context, client *pubsub.Client, topic, name string) (*pubsub.Subscription, error) {
	return client.CreateSubscription(ctx, name, pubsub.SubscriptionConfig{
		Topic:            client.Topic(topic),
		ExpirationPolicy: time.Hour * 24 * 3,
	})
}

type errRateLimited struct {
	limiterName string
	senderID    persist.DBID
	receiverID  persist.DBID
	feedEventID persist.DBID
}

func (e errRateLimited) Error() string {
	str := fmt.Sprintf("rate limit exceeded for limiter=%s, sender=%s, receiver=%s", e.limiterName, e.senderID, e.receiverID)
	if e.feedEventID != "" {
		str += fmt.Sprintf(", feedEvent=%s", e.feedEventID)
	}
	return str
}
