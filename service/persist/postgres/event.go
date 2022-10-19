package postgres

import (
	"context"
	"time"

	db "github.com/mikeydub/go-gallery/db/gen/coredb"
	"github.com/mikeydub/go-gallery/service/persist"
	"github.com/spf13/viper"
)

type EventRepository struct {
	Queries *db.Queries
}

func (r *EventRepository) Get(ctx context.Context, eventID persist.DBID) (db.Event, error) {
	return r.Queries.GetEvent(ctx, eventID)
}

func (r *EventRepository) Add(ctx context.Context, event db.Event) (*db.Event, error) {
	switch event.ResourceTypeID {
	case persist.ResourceTypeUser:
		return r.AddUserEvent(ctx, event)
	case persist.ResourceTypeToken:
		return r.AddTokenEvent(ctx, event)
	case persist.ResourceTypeCollection:
		return r.AddCollectionEvent(ctx, event)
	default:
		return nil, persist.ErrUnknownResourceType{ResourceType: event.ResourceTypeID}
	}
}

func (r *EventRepository) AddUserEvent(ctx context.Context, event db.Event) (*db.Event, error) {
	event, err := r.Queries.CreateUserEvent(ctx, db.CreateUserEventParams{
		ID:             persist.GenerateID(),
		ActorID:        event.ActorID,
		Action:         event.Action,
		ResourceTypeID: event.ResourceTypeID,
		UserID:         event.SubjectID,
		Data:           event.Data,
	})
	return &event, err
}

func (r *EventRepository) AddTokenEvent(ctx context.Context, event db.Event) (*db.Event, error) {
	event, err := r.Queries.CreateTokenEvent(ctx, db.CreateTokenEventParams{
		ID:             persist.GenerateID(),
		ActorID:        event.ActorID,
		Action:         event.Action,
		ResourceTypeID: event.ResourceTypeID,
		TokenID:        event.SubjectID,
		Data:           event.Data,
	})
	return &event, err
}

func (r *EventRepository) AddCollectionEvent(ctx context.Context, event db.Event) (*db.Event, error) {
	event, err := r.Queries.CreateCollectionEvent(ctx, db.CreateCollectionEventParams{
		ID:             persist.GenerateID(),
		ActorID:        event.ActorID,
		Action:         event.Action,
		ResourceTypeID: event.ResourceTypeID,
		CollectionID:   event.SubjectID,
		Data:           event.Data,
		Caption:        event.Caption,
	})
	return &event, err
}

// WindowActive checks if there are more recent events with an action that matches the provided event.
func (r *EventRepository) WindowActive(ctx context.Context, event db.Event) (bool, error) {
	return r.Queries.IsWindowActive(ctx, db.IsWindowActiveParams{
		ActorID:     event.ActorID,
		Action:      event.Action,
		WindowStart: event.CreatedAt,
		WindowEnd:   event.CreatedAt.Add(time.Duration(viper.GetInt("FEED_WINDOW_SIZE")) * time.Second),
	})
}

// WindowActiveForSubject checks if there are more recent events with an action on a specific resource such as
// as a collection or a token.
func (r *EventRepository) WindowActiveForSubject(ctx context.Context, event db.Event) (bool, error) {
	return r.Queries.IsWindowActiveWithSubject(ctx, db.IsWindowActiveWithSubjectParams{
		ActorID:     event.ActorID,
		Action:      event.Action,
		SubjectID:   event.SubjectID,
		WindowStart: event.CreatedAt,
		WindowEnd:   event.CreatedAt.Add(time.Duration(viper.GetInt("FEED_WINDOW_SIZE")) * time.Second),
	})
}

// EventsInWindow returns events belonging to the same window of activity as the given eventID.
func (r *EventRepository) EventsInWindow(ctx context.Context, eventID persist.DBID, windowSeconds int) ([]db.Event, error) {
	return r.Queries.GetEventsInWindow(ctx, db.GetEventsInWindowParams{
		ID:   eventID,
		Secs: float64(windowSeconds),
	})
}
