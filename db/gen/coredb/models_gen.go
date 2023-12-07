// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package coredb

import (
	"database/sql"
	"time"

	"github.com/jackc/pgtype"
	"github.com/mikeydub/go-gallery/service/persist"
)

type Admire struct {
	ID          persist.DBID `db:"id" json:"id"`
	Version     int32        `db:"version" json:"version"`
	FeedEventID persist.DBID `db:"feed_event_id" json:"feed_event_id"`
	ActorID     persist.DBID `db:"actor_id" json:"actor_id"`
	Deleted     bool         `db:"deleted" json:"deleted"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	LastUpdated time.Time    `db:"last_updated" json:"last_updated"`
	PostID      persist.DBID `db:"post_id" json:"post_id"`
	TokenID     persist.DBID `db:"token_id" json:"token_id"`
}

type AlchemySpamContract struct {
	ID        persist.DBID    `db:"id" json:"id"`
	Chain     persist.Chain   `db:"chain" json:"chain"`
	Address   persist.Address `db:"address" json:"address"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
	IsSpam    bool            `db:"is_spam" json:"is_spam"`
}

type Collection struct {
	ID             persist.DBID                                     `db:"id" json:"id"`
	Deleted        bool                                             `db:"deleted" json:"deleted"`
	OwnerUserID    persist.DBID                                     `db:"owner_user_id" json:"owner_user_id"`
	Nfts           persist.DBIDList                                 `db:"nfts" json:"nfts"`
	Version        sql.NullInt32                                    `db:"version" json:"version"`
	LastUpdated    time.Time                                        `db:"last_updated" json:"last_updated"`
	CreatedAt      time.Time                                        `db:"created_at" json:"created_at"`
	Hidden         bool                                             `db:"hidden" json:"hidden"`
	CollectorsNote sql.NullString                                   `db:"collectors_note" json:"collectors_note"`
	Name           sql.NullString                                   `db:"name" json:"name"`
	Layout         persist.TokenLayout                              `db:"layout" json:"layout"`
	TokenSettings  map[persist.DBID]persist.CollectionTokenSettings `db:"token_settings" json:"token_settings"`
	GalleryID      persist.DBID                                     `db:"gallery_id" json:"gallery_id"`
}

type Comment struct {
	ID                persist.DBID `db:"id" json:"id"`
	Version           int32        `db:"version" json:"version"`
	FeedEventID       persist.DBID `db:"feed_event_id" json:"feed_event_id"`
	ActorID           persist.DBID `db:"actor_id" json:"actor_id"`
	ReplyTo           persist.DBID `db:"reply_to" json:"reply_to"`
	Comment           string       `db:"comment" json:"comment"`
	Deleted           bool         `db:"deleted" json:"deleted"`
	CreatedAt         time.Time    `db:"created_at" json:"created_at"`
	LastUpdated       time.Time    `db:"last_updated" json:"last_updated"`
	PostID            persist.DBID `db:"post_id" json:"post_id"`
	Removed           bool         `db:"removed" json:"removed"`
	TopLevelCommentID persist.DBID `db:"top_level_comment_id" json:"top_level_comment_id"`
}

type Community struct {
	ID                      persist.DBID          `db:"id" json:"id"`
	Version                 int32                 `db:"version" json:"version"`
	CommunityType           persist.CommunityType `db:"community_type" json:"community_type"`
	Key1                    string                `db:"key1" json:"key1"`
	Key2                    string                `db:"key2" json:"key2"`
	Key3                    string                `db:"key3" json:"key3"`
	Key4                    string                `db:"key4" json:"key4"`
	Name                    string                `db:"name" json:"name"`
	OverrideName            sql.NullString        `db:"override_name" json:"override_name"`
	Description             string                `db:"description" json:"description"`
	OverrideDescription     sql.NullString        `db:"override_description" json:"override_description"`
	ProfileImageUrl         sql.NullString        `db:"profile_image_url" json:"profile_image_url"`
	OverrideProfileImageUrl sql.NullString        `db:"override_profile_image_url" json:"override_profile_image_url"`
	BadgeImageUrl           sql.NullString        `db:"badge_image_url" json:"badge_image_url"`
	OverrideBadgeImageUrl   sql.NullString        `db:"override_badge_image_url" json:"override_badge_image_url"`
	ContractID              persist.DBID          `db:"contract_id" json:"contract_id"`
	CreatedAt               time.Time             `db:"created_at" json:"created_at"`
	LastUpdated             time.Time             `db:"last_updated" json:"last_updated"`
	Deleted                 bool                  `db:"deleted" json:"deleted"`
}

type CommunityCreator struct {
	ID                    persist.DBID                 `db:"id" json:"id"`
	Version               int32                        `db:"version" json:"version"`
	CreatorType           persist.CommunityCreatorType `db:"creator_type" json:"creator_type"`
	CommunityID           persist.DBID                 `db:"community_id" json:"community_id"`
	CreatorUserID         persist.DBID                 `db:"creator_user_id" json:"creator_user_id"`
	CreatorAddress        persist.Address              `db:"creator_address" json:"creator_address"`
	CreatorAddressL1Chain persist.L1Chain              `db:"creator_address_l1_chain" json:"creator_address_l1_chain"`
	CreatorAddressChain   persist.Chain                `db:"creator_address_chain" json:"creator_address_chain"`
	CreatedAt             time.Time                    `db:"created_at" json:"created_at"`
	LastUpdated           time.Time                    `db:"last_updated" json:"last_updated"`
	Deleted               bool                         `db:"deleted" json:"deleted"`
}

type Contract struct {
	ID                    persist.DBID    `db:"id" json:"id"`
	Deleted               bool            `db:"deleted" json:"deleted"`
	Version               sql.NullInt32   `db:"version" json:"version"`
	CreatedAt             time.Time       `db:"created_at" json:"created_at"`
	LastUpdated           time.Time       `db:"last_updated" json:"last_updated"`
	Name                  sql.NullString  `db:"name" json:"name"`
	Symbol                sql.NullString  `db:"symbol" json:"symbol"`
	Address               persist.Address `db:"address" json:"address"`
	CreatorAddress        persist.Address `db:"creator_address" json:"creator_address"`
	Chain                 persist.Chain   `db:"chain" json:"chain"`
	ProfileBannerUrl      sql.NullString  `db:"profile_banner_url" json:"profile_banner_url"`
	ProfileImageUrl       sql.NullString  `db:"profile_image_url" json:"profile_image_url"`
	BadgeUrl              sql.NullString  `db:"badge_url" json:"badge_url"`
	Description           sql.NullString  `db:"description" json:"description"`
	OwnerAddress          persist.Address `db:"owner_address" json:"owner_address"`
	IsProviderMarkedSpam  bool            `db:"is_provider_marked_spam" json:"is_provider_marked_spam"`
	ParentID              persist.DBID    `db:"parent_id" json:"parent_id"`
	OverrideCreatorUserID persist.DBID    `db:"override_creator_user_id" json:"override_creator_user_id"`
	L1Chain               persist.L1Chain `db:"l1_chain" json:"l1_chain"`
}

type ContractCommunityMembership struct {
	ID          persist.DBID `db:"id" json:"id"`
	Version     int32        `db:"version" json:"version"`
	ContractID  persist.DBID `db:"contract_id" json:"contract_id"`
	CommunityID persist.DBID `db:"community_id" json:"community_id"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	LastUpdated time.Time    `db:"last_updated" json:"last_updated"`
	Deleted     bool         `db:"deleted" json:"deleted"`
}

type ContractCommunityType struct {
	ID            persist.DBID          `db:"id" json:"id"`
	Version       int32                 `db:"version" json:"version"`
	ContractID    persist.DBID          `db:"contract_id" json:"contract_id"`
	CommunityType persist.CommunityType `db:"community_type" json:"community_type"`
	IsValidType   bool                  `db:"is_valid_type" json:"is_valid_type"`
	CreatedAt     time.Time             `db:"created_at" json:"created_at"`
	LastUpdated   time.Time             `db:"last_updated" json:"last_updated"`
	Deleted       bool                  `db:"deleted" json:"deleted"`
}

type ContractCreator struct {
	ContractID     persist.DBID    `db:"contract_id" json:"contract_id"`
	CreatorUserID  persist.DBID    `db:"creator_user_id" json:"creator_user_id"`
	Chain          persist.Chain   `db:"chain" json:"chain"`
	CreatorAddress persist.Address `db:"creator_address" json:"creator_address"`
}

type ContractRelevance struct {
	ID    persist.DBID `db:"id" json:"id"`
	Score int32        `db:"score" json:"score"`
}

type DevMetadataUser struct {
	UserID          persist.DBID  `db:"user_id" json:"user_id"`
	HasEmailAddress persist.Email `db:"has_email_address" json:"has_email_address"`
	Deleted         bool          `db:"deleted" json:"deleted"`
}

type EarlyAccess struct {
	Address persist.Address `db:"address" json:"address"`
}

type Event struct {
	ID             persist.DBID         `db:"id" json:"id"`
	Version        int32                `db:"version" json:"version"`
	ActorID        sql.NullString       `db:"actor_id" json:"actor_id"`
	ResourceTypeID persist.ResourceType `db:"resource_type_id" json:"resource_type_id"`
	SubjectID      persist.DBID         `db:"subject_id" json:"subject_id"`
	UserID         persist.DBID         `db:"user_id" json:"user_id"`
	TokenID        persist.DBID         `db:"token_id" json:"token_id"`
	CollectionID   persist.DBID         `db:"collection_id" json:"collection_id"`
	Action         persist.Action       `db:"action" json:"action"`
	Data           persist.EventData    `db:"data" json:"data"`
	Deleted        bool                 `db:"deleted" json:"deleted"`
	LastUpdated    time.Time            `db:"last_updated" json:"last_updated"`
	CreatedAt      time.Time            `db:"created_at" json:"created_at"`
	GalleryID      persist.DBID         `db:"gallery_id" json:"gallery_id"`
	CommentID      persist.DBID         `db:"comment_id" json:"comment_id"`
	AdmireID       persist.DBID         `db:"admire_id" json:"admire_id"`
	FeedEventID    persist.DBID         `db:"feed_event_id" json:"feed_event_id"`
	ExternalID     sql.NullString       `db:"external_id" json:"external_id"`
	Caption        sql.NullString       `db:"caption" json:"caption"`
	GroupID        sql.NullString       `db:"group_id" json:"group_id"`
	PostID         persist.DBID         `db:"post_id" json:"post_id"`
	ContractID     persist.DBID         `db:"contract_id" json:"contract_id"`
	MentionID      persist.DBID         `db:"mention_id" json:"mention_id"`
}

type ExternalSocialConnection struct {
	ID                persist.DBID `db:"id" json:"id"`
	Version           int32        `db:"version" json:"version"`
	SocialAccountType string       `db:"social_account_type" json:"social_account_type"`
	FollowerID        persist.DBID `db:"follower_id" json:"follower_id"`
	FolloweeID        persist.DBID `db:"followee_id" json:"followee_id"`
	CreatedAt         time.Time    `db:"created_at" json:"created_at"`
	LastUpdated       time.Time    `db:"last_updated" json:"last_updated"`
	Deleted           bool         `db:"deleted" json:"deleted"`
}

type FeedBlocklist struct {
	ID          persist.DBID   `db:"id" json:"id"`
	UserID      persist.DBID   `db:"user_id" json:"user_id"`
	Action      persist.Action `db:"action" json:"action"`
	LastUpdated time.Time      `db:"last_updated" json:"last_updated"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	Deleted     bool           `db:"deleted" json:"deleted"`
}

type FeedEntity struct {
	ID             persist.DBID `db:"id" json:"id"`
	FeedEntityType int32        `db:"feed_entity_type" json:"feed_entity_type"`
	CreatedAt      time.Time    `db:"created_at" json:"created_at"`
	ActorID        persist.DBID `db:"actor_id" json:"actor_id"`
}

type FeedEntityScore struct {
	ID             persist.DBID     `db:"id" json:"id"`
	CreatedAt      time.Time        `db:"created_at" json:"created_at"`
	ActorID        persist.DBID     `db:"actor_id" json:"actor_id"`
	Action         persist.Action   `db:"action" json:"action"`
	ContractIds    persist.DBIDList `db:"contract_ids" json:"contract_ids"`
	Interactions   int32            `db:"interactions" json:"interactions"`
	FeedEntityType int32            `db:"feed_entity_type" json:"feed_entity_type"`
	LastUpdated    time.Time        `db:"last_updated" json:"last_updated"`
}

type FeedEntityScoreView struct {
	ID             persist.DBID     `db:"id" json:"id"`
	CreatedAt      time.Time        `db:"created_at" json:"created_at"`
	ActorID        persist.DBID     `db:"actor_id" json:"actor_id"`
	Action         persist.Action   `db:"action" json:"action"`
	ContractIds    persist.DBIDList `db:"contract_ids" json:"contract_ids"`
	Interactions   int32            `db:"interactions" json:"interactions"`
	FeedEntityType int32            `db:"feed_entity_type" json:"feed_entity_type"`
	LastUpdated    time.Time        `db:"last_updated" json:"last_updated"`
}

type FeedEvent struct {
	ID          persist.DBID          `db:"id" json:"id"`
	Version     int32                 `db:"version" json:"version"`
	OwnerID     persist.DBID          `db:"owner_id" json:"owner_id"`
	Action      persist.Action        `db:"action" json:"action"`
	Data        persist.FeedEventData `db:"data" json:"data"`
	EventTime   time.Time             `db:"event_time" json:"event_time"`
	EventIds    persist.DBIDList      `db:"event_ids" json:"event_ids"`
	Deleted     bool                  `db:"deleted" json:"deleted"`
	LastUpdated time.Time             `db:"last_updated" json:"last_updated"`
	CreatedAt   time.Time             `db:"created_at" json:"created_at"`
	Caption     sql.NullString        `db:"caption" json:"caption"`
	GroupID     sql.NullString        `db:"group_id" json:"group_id"`
}

type Follow struct {
	ID          persist.DBID `db:"id" json:"id"`
	Follower    persist.DBID `db:"follower" json:"follower"`
	Followee    persist.DBID `db:"followee" json:"followee"`
	Deleted     bool         `db:"deleted" json:"deleted"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	LastUpdated time.Time    `db:"last_updated" json:"last_updated"`
}

type Gallery struct {
	ID          persist.DBID     `db:"id" json:"id"`
	Deleted     bool             `db:"deleted" json:"deleted"`
	LastUpdated time.Time        `db:"last_updated" json:"last_updated"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
	Version     sql.NullInt32    `db:"version" json:"version"`
	OwnerUserID persist.DBID     `db:"owner_user_id" json:"owner_user_id"`
	Collections persist.DBIDList `db:"collections" json:"collections"`
	Name        string           `db:"name" json:"name"`
	Description string           `db:"description" json:"description"`
	Hidden      bool             `db:"hidden" json:"hidden"`
	Position    string           `db:"position" json:"position"`
}

type GalleryRelevance struct {
	ID    persist.DBID `db:"id" json:"id"`
	Score int32        `db:"score" json:"score"`
}

type LegacyView struct {
	UserID      persist.DBID  `db:"user_id" json:"user_id"`
	ViewCount   sql.NullInt32 `db:"view_count" json:"view_count"`
	LastUpdated time.Time     `db:"last_updated" json:"last_updated"`
	CreatedAt   time.Time     `db:"created_at" json:"created_at"`
	Deleted     sql.NullBool  `db:"deleted" json:"deleted"`
}

type MarketplaceContract struct {
	ContractID persist.DBID `db:"contract_id" json:"contract_id"`
}

type MediaValidationRule struct {
	ID        persist.DBID `db:"id" json:"id"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	MediaType string       `db:"media_type" json:"media_type"`
	Property  string       `db:"property" json:"property"`
	Required  bool         `db:"required" json:"required"`
}

type Membership struct {
	ID          persist.DBID            `db:"id" json:"id"`
	Deleted     bool                    `db:"deleted" json:"deleted"`
	Version     sql.NullInt32           `db:"version" json:"version"`
	CreatedAt   time.Time               `db:"created_at" json:"created_at"`
	LastUpdated time.Time               `db:"last_updated" json:"last_updated"`
	TokenID     persist.DBID            `db:"token_id" json:"token_id"`
	Name        sql.NullString          `db:"name" json:"name"`
	AssetUrl    sql.NullString          `db:"asset_url" json:"asset_url"`
	Owners      persist.TokenHolderList `db:"owners" json:"owners"`
}

type Mention struct {
	ID         persist.DBID  `db:"id" json:"id"`
	PostID     persist.DBID  `db:"post_id" json:"post_id"`
	CommentID  persist.DBID  `db:"comment_id" json:"comment_id"`
	UserID     persist.DBID  `db:"user_id" json:"user_id"`
	ContractID persist.DBID  `db:"contract_id" json:"contract_id"`
	Start      sql.NullInt32 `db:"start" json:"start"`
	Length     sql.NullInt32 `db:"length" json:"length"`
	CreatedAt  time.Time     `db:"created_at" json:"created_at"`
	Deleted    bool          `db:"deleted" json:"deleted"`
}

type Merch struct {
	ID           persist.DBID    `db:"id" json:"id"`
	Deleted      bool            `db:"deleted" json:"deleted"`
	Version      sql.NullInt32   `db:"version" json:"version"`
	CreatedAt    time.Time       `db:"created_at" json:"created_at"`
	LastUpdated  time.Time       `db:"last_updated" json:"last_updated"`
	TokenID      persist.TokenID `db:"token_id" json:"token_id"`
	ObjectType   int32           `db:"object_type" json:"object_type"`
	DiscountCode sql.NullString  `db:"discount_code" json:"discount_code"`
	Redeemed     bool            `db:"redeemed" json:"redeemed"`
}

type MigrationValidation struct {
	ID                       persist.DBID   `db:"id" json:"id"`
	MediaID                  persist.DBID   `db:"media_id" json:"media_id"`
	ProcessingJobID          persist.DBID   `db:"processing_job_id" json:"processing_job_id"`
	Chain                    persist.Chain  `db:"chain" json:"chain"`
	Contract                 sql.NullString `db:"contract" json:"contract"`
	TokenID                  persist.DBID   `db:"token_id" json:"token_id"`
	MediaType                interface{}    `db:"media_type" json:"media_type"`
	RemappedTo               interface{}    `db:"remapped_to" json:"remapped_to"`
	OldMedia                 pgtype.JSONB   `db:"old_media" json:"old_media"`
	NewMedia                 pgtype.JSONB   `db:"new_media" json:"new_media"`
	MediaTypeValidation      string         `db:"media_type_validation" json:"media_type_validation"`
	DimensionsValidation     string         `db:"dimensions_validation" json:"dimensions_validation"`
	MediaUrlValidation       string         `db:"media_url_validation" json:"media_url_validation"`
	ThumbnailUrlValidation   string         `db:"thumbnail_url_validation" json:"thumbnail_url_validation"`
	LivePreviewUrlValidation string         `db:"live_preview_url_validation" json:"live_preview_url_validation"`
	LastRefreshed            interface{}    `db:"last_refreshed" json:"last_refreshed"`
}

type Nonce struct {
	ID          persist.DBID    `db:"id" json:"id"`
	Deleted     bool            `db:"deleted" json:"deleted"`
	Version     sql.NullInt32   `db:"version" json:"version"`
	LastUpdated time.Time       `db:"last_updated" json:"last_updated"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UserID      persist.DBID    `db:"user_id" json:"user_id"`
	Address     persist.Address `db:"address" json:"address"`
	Value       sql.NullString  `db:"value" json:"value"`
	Chain       persist.Chain   `db:"chain" json:"chain"`
	L1Chain     persist.L1Chain `db:"l1_chain" json:"l1_chain"`
}

type Notification struct {
	ID          persist.DBID             `db:"id" json:"id"`
	Deleted     bool                     `db:"deleted" json:"deleted"`
	OwnerID     persist.DBID             `db:"owner_id" json:"owner_id"`
	Version     sql.NullInt32            `db:"version" json:"version"`
	LastUpdated time.Time                `db:"last_updated" json:"last_updated"`
	CreatedAt   time.Time                `db:"created_at" json:"created_at"`
	Action      persist.Action           `db:"action" json:"action"`
	Data        persist.NotificationData `db:"data" json:"data"`
	EventIds    persist.DBIDList         `db:"event_ids" json:"event_ids"`
	FeedEventID persist.DBID             `db:"feed_event_id" json:"feed_event_id"`
	CommentID   persist.DBID             `db:"comment_id" json:"comment_id"`
	GalleryID   persist.DBID             `db:"gallery_id" json:"gallery_id"`
	Seen        bool                     `db:"seen" json:"seen"`
	Amount      int32                    `db:"amount" json:"amount"`
	PostID      persist.DBID             `db:"post_id" json:"post_id"`
	TokenID     persist.DBID             `db:"token_id" json:"token_id"`
	ContractID  persist.DBID             `db:"contract_id" json:"contract_id"`
	MentionID   persist.DBID             `db:"mention_id" json:"mention_id"`
}

type OwnedContract struct {
	UserID         persist.DBID `db:"user_id" json:"user_id"`
	UserCreatedAt  time.Time    `db:"user_created_at" json:"user_created_at"`
	ContractID     persist.DBID `db:"contract_id" json:"contract_id"`
	OwnedCount     int64        `db:"owned_count" json:"owned_count"`
	DisplayedCount int64        `db:"displayed_count" json:"displayed_count"`
	Displayed      bool         `db:"displayed" json:"displayed"`
	LastUpdated    time.Time    `db:"last_updated" json:"last_updated"`
}

type PiiAccountCreationInfo struct {
	UserID    persist.DBID `db:"user_id" json:"user_id"`
	IpAddress string       `db:"ip_address" json:"ip_address"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
}

type PiiForUser struct {
	UserID          persist.DBID    `db:"user_id" json:"user_id"`
	PiiEmailAddress persist.Email   `db:"pii_email_address" json:"pii_email_address"`
	Deleted         bool            `db:"deleted" json:"deleted"`
	PiiSocials      persist.Socials `db:"pii_socials" json:"pii_socials"`
}

type PiiSocialsAuth struct {
	ID           persist.DBID           `db:"id" json:"id"`
	Deleted      bool                   `db:"deleted" json:"deleted"`
	Version      sql.NullInt32          `db:"version" json:"version"`
	CreatedAt    time.Time              `db:"created_at" json:"created_at"`
	LastUpdated  time.Time              `db:"last_updated" json:"last_updated"`
	UserID       persist.DBID           `db:"user_id" json:"user_id"`
	Provider     persist.SocialProvider `db:"provider" json:"provider"`
	AccessToken  sql.NullString         `db:"access_token" json:"access_token"`
	RefreshToken sql.NullString         `db:"refresh_token" json:"refresh_token"`
}

type PiiUserView struct {
	ID                   persist.DBID                     `db:"id" json:"id"`
	Deleted              bool                             `db:"deleted" json:"deleted"`
	Version              sql.NullInt32                    `db:"version" json:"version"`
	LastUpdated          time.Time                        `db:"last_updated" json:"last_updated"`
	CreatedAt            time.Time                        `db:"created_at" json:"created_at"`
	Username             sql.NullString                   `db:"username" json:"username"`
	UsernameIdempotent   sql.NullString                   `db:"username_idempotent" json:"username_idempotent"`
	Wallets              persist.WalletList               `db:"wallets" json:"wallets"`
	Bio                  sql.NullString                   `db:"bio" json:"bio"`
	Traits               pgtype.JSONB                     `db:"traits" json:"traits"`
	Universal            bool                             `db:"universal" json:"universal"`
	NotificationSettings persist.UserNotificationSettings `db:"notification_settings" json:"notification_settings"`
	EmailVerified        persist.EmailVerificationStatus  `db:"email_verified" json:"email_verified"`
	EmailUnsubscriptions persist.EmailUnsubscriptions     `db:"email_unsubscriptions" json:"email_unsubscriptions"`
	FeaturedGallery      *persist.DBID                    `db:"featured_gallery" json:"featured_gallery"`
	PrimaryWalletID      persist.DBID                     `db:"primary_wallet_id" json:"primary_wallet_id"`
	UserExperiences      pgtype.JSONB                     `db:"user_experiences" json:"user_experiences"`
	PiiEmailAddress      persist.Email                    `db:"pii_email_address" json:"pii_email_address"`
	PiiSocials           persist.Socials                  `db:"pii_socials" json:"pii_socials"`
}

type Post struct {
	ID          persist.DBID     `db:"id" json:"id"`
	Version     int32            `db:"version" json:"version"`
	TokenIds    persist.DBIDList `db:"token_ids" json:"token_ids"`
	ContractIds persist.DBIDList `db:"contract_ids" json:"contract_ids"`
	ActorID     persist.DBID     `db:"actor_id" json:"actor_id"`
	Caption     sql.NullString   `db:"caption" json:"caption"`
	CreatedAt   time.Time        `db:"created_at" json:"created_at"`
	LastUpdated time.Time        `db:"last_updated" json:"last_updated"`
	Deleted     bool             `db:"deleted" json:"deleted"`
}

type ProfileImage struct {
	ID           persist.DBID               `db:"id" json:"id"`
	UserID       persist.DBID               `db:"user_id" json:"user_id"`
	TokenID      persist.DBID               `db:"token_id" json:"token_id"`
	SourceType   persist.ProfileImageSource `db:"source_type" json:"source_type"`
	Deleted      bool                       `db:"deleted" json:"deleted"`
	CreatedAt    time.Time                  `db:"created_at" json:"created_at"`
	LastUpdated  time.Time                  `db:"last_updated" json:"last_updated"`
	WalletID     persist.DBID               `db:"wallet_id" json:"wallet_id"`
	EnsAvatarUri sql.NullString             `db:"ens_avatar_uri" json:"ens_avatar_uri"`
	EnsDomain    sql.NullString             `db:"ens_domain" json:"ens_domain"`
}

type PushNotificationTicket struct {
	ID               persist.DBID `db:"id" json:"id"`
	PushTokenID      persist.DBID `db:"push_token_id" json:"push_token_id"`
	TicketID         string       `db:"ticket_id" json:"ticket_id"`
	CreatedAt        time.Time    `db:"created_at" json:"created_at"`
	CheckAfter       time.Time    `db:"check_after" json:"check_after"`
	NumCheckAttempts int32        `db:"num_check_attempts" json:"num_check_attempts"`
	Deleted          bool         `db:"deleted" json:"deleted"`
	Status           string       `db:"status" json:"status"`
}

type PushNotificationToken struct {
	ID        persist.DBID `db:"id" json:"id"`
	UserID    persist.DBID `db:"user_id" json:"user_id"`
	PushToken string       `db:"push_token" json:"push_token"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	Deleted   bool         `db:"deleted" json:"deleted"`
}

type RecommendationResult struct {
	ID                persist.DBID  `db:"id" json:"id"`
	Version           sql.NullInt32 `db:"version" json:"version"`
	UserID            persist.DBID  `db:"user_id" json:"user_id"`
	RecommendedUserID persist.DBID  `db:"recommended_user_id" json:"recommended_user_id"`
	RecommendedCount  sql.NullInt32 `db:"recommended_count" json:"recommended_count"`
	CreatedAt         time.Time     `db:"created_at" json:"created_at"`
	LastUpdated       time.Time     `db:"last_updated" json:"last_updated"`
	Deleted           bool          `db:"deleted" json:"deleted"`
}

type ReprocessJob struct {
	ID           int          `db:"id" json:"id"`
	TokenStartID persist.DBID `db:"token_start_id" json:"token_start_id"`
	TokenEndID   persist.DBID `db:"token_end_id" json:"token_end_id"`
}

type ScrubbedPiiAccountCreationInfo struct {
	UserID    persist.DBID    `db:"user_id" json:"user_id"`
	IpAddress persist.Address `db:"ip_address" json:"ip_address"`
	CreatedAt time.Time       `db:"created_at" json:"created_at"`
}

type ScrubbedPiiForUser struct {
	UserID          persist.DBID    `db:"user_id" json:"user_id"`
	PiiEmailAddress persist.Email   `db:"pii_email_address" json:"pii_email_address"`
	Deleted         bool            `db:"deleted" json:"deleted"`
	PiiSocials      persist.Socials `db:"pii_socials" json:"pii_socials"`
}

type Session struct {
	ID                   persist.DBID `db:"id" json:"id"`
	UserID               persist.DBID `db:"user_id" json:"user_id"`
	CreatedAt            time.Time    `db:"created_at" json:"created_at"`
	CreatedWithUserAgent string       `db:"created_with_user_agent" json:"created_with_user_agent"`
	CreatedWithPlatform  string       `db:"created_with_platform" json:"created_with_platform"`
	CreatedWithOs        string       `db:"created_with_os" json:"created_with_os"`
	LastRefreshed        time.Time    `db:"last_refreshed" json:"last_refreshed"`
	LastUserAgent        string       `db:"last_user_agent" json:"last_user_agent"`
	LastPlatform         string       `db:"last_platform" json:"last_platform"`
	LastOs               string       `db:"last_os" json:"last_os"`
	CurrentRefreshID     string       `db:"current_refresh_id" json:"current_refresh_id"`
	ActiveUntil          time.Time    `db:"active_until" json:"active_until"`
	Invalidated          bool         `db:"invalidated" json:"invalidated"`
	LastUpdated          time.Time    `db:"last_updated" json:"last_updated"`
	Deleted              bool         `db:"deleted" json:"deleted"`
}

type SpamUserScore struct {
	UserID        persist.DBID `db:"user_id" json:"user_id"`
	Score         int32        `db:"score" json:"score"`
	DecidedIsSpam sql.NullBool `db:"decided_is_spam" json:"decided_is_spam"`
	DecidedAt     sql.NullTime `db:"decided_at" json:"decided_at"`
	Deleted       bool         `db:"deleted" json:"deleted"`
	CreatedAt     time.Time    `db:"created_at" json:"created_at"`
}

type Token struct {
	ID                persist.DBID      `db:"id" json:"id"`
	Deleted           bool              `db:"deleted" json:"deleted"`
	Version           sql.NullInt32     `db:"version" json:"version"`
	CreatedAt         time.Time         `db:"created_at" json:"created_at"`
	LastUpdated       time.Time         `db:"last_updated" json:"last_updated"`
	CollectorsNote    sql.NullString    `db:"collectors_note" json:"collectors_note"`
	Quantity          persist.HexString `db:"quantity" json:"quantity"`
	BlockNumber       sql.NullInt64     `db:"block_number" json:"block_number"`
	OwnerUserID       persist.DBID      `db:"owner_user_id" json:"owner_user_id"`
	OwnedByWallets    persist.DBIDList  `db:"owned_by_wallets" json:"owned_by_wallets"`
	ContractID        persist.DBID      `db:"contract_id" json:"contract_id"`
	IsUserMarkedSpam  sql.NullBool      `db:"is_user_marked_spam" json:"is_user_marked_spam"`
	LastSynced        time.Time         `db:"last_synced" json:"last_synced"`
	IsCreatorToken    bool              `db:"is_creator_token" json:"is_creator_token"`
	TokenDefinitionID persist.DBID      `db:"token_definition_id" json:"token_definition_id"`
	IsHolderToken     bool              `db:"is_holder_token" json:"is_holder_token"`
	Displayable       bool              `db:"displayable" json:"displayable"`
}

type TokenCommunityMembership struct {
	ID                persist.DBID `db:"id" json:"id"`
	Version           int32        `db:"version" json:"version"`
	TokenDefinitionID persist.DBID `db:"token_definition_id" json:"token_definition_id"`
	CommunityID       persist.DBID `db:"community_id" json:"community_id"`
	CreatedAt         time.Time    `db:"created_at" json:"created_at"`
	LastUpdated       time.Time    `db:"last_updated" json:"last_updated"`
	Deleted           bool         `db:"deleted" json:"deleted"`
}

type TokenDefinition struct {
	ID              persist.DBID          `db:"id" json:"id"`
	CreatedAt       time.Time             `db:"created_at" json:"created_at"`
	LastUpdated     time.Time             `db:"last_updated" json:"last_updated"`
	Deleted         bool                  `db:"deleted" json:"deleted"`
	Name            sql.NullString        `db:"name" json:"name"`
	Description     sql.NullString        `db:"description" json:"description"`
	TokenType       persist.TokenType     `db:"token_type" json:"token_type"`
	TokenID         persist.TokenID       `db:"token_id" json:"token_id"`
	ExternalUrl     sql.NullString        `db:"external_url" json:"external_url"`
	Chain           persist.Chain         `db:"chain" json:"chain"`
	Metadata        persist.TokenMetadata `db:"metadata" json:"metadata"`
	FallbackMedia   persist.FallbackMedia `db:"fallback_media" json:"fallback_media"`
	ContractAddress persist.Address       `db:"contract_address" json:"contract_address"`
	ContractID      persist.DBID          `db:"contract_id" json:"contract_id"`
	TokenMediaID    persist.DBID          `db:"token_media_id" json:"token_media_id"`
}

type TokenMedia struct {
	ID              persist.DBID  `db:"id" json:"id"`
	CreatedAt       time.Time     `db:"created_at" json:"created_at"`
	LastUpdated     time.Time     `db:"last_updated" json:"last_updated"`
	Version         int32         `db:"version" json:"version"`
	Active          bool          `db:"active" json:"active"`
	Media           persist.Media `db:"media" json:"media"`
	ProcessingJobID persist.DBID  `db:"processing_job_id" json:"processing_job_id"`
	Deleted         bool          `db:"deleted" json:"deleted"`
}

type TokenMediasActive struct {
	ID               persist.DBID `db:"id" json:"id"`
	LastUpdated      time.Time    `db:"last_updated" json:"last_updated"`
	MediaType        interface{}  `db:"media_type" json:"media_type"`
	JobID            persist.DBID `db:"job_id" json:"job_id"`
	TokenProperties  pgtype.JSONB `db:"token_properties" json:"token_properties"`
	PipelineMetadata pgtype.JSONB `db:"pipeline_metadata" json:"pipeline_metadata"`
}

type TokenMediasMissingProperty struct {
	ID          persist.DBID `db:"id" json:"id"`
	MediaType   interface{}  `db:"media_type" json:"media_type"`
	LastUpdated time.Time    `db:"last_updated" json:"last_updated"`
	IsValid     bool         `db:"is_valid" json:"is_valid"`
	Reason      []byte       `db:"reason" json:"reason"`
}

type TokenMediasNoValidationRule struct {
	ID          persist.DBID `db:"id" json:"id"`
	MediaType   interface{}  `db:"media_type" json:"media_type"`
	LastUpdated time.Time    `db:"last_updated" json:"last_updated"`
	IsValid     bool         `db:"is_valid" json:"is_valid"`
	Reason      string       `db:"reason" json:"reason"`
}

type TokenProcessingJob struct {
	ID               persist.DBID             `db:"id" json:"id"`
	CreatedAt        time.Time                `db:"created_at" json:"created_at"`
	TokenProperties  persist.TokenProperties  `db:"token_properties" json:"token_properties"`
	PipelineMetadata persist.PipelineMetadata `db:"pipeline_metadata" json:"pipeline_metadata"`
	ProcessingCause  persist.ProcessingCause  `db:"processing_cause" json:"processing_cause"`
	ProcessorVersion string                   `db:"processor_version" json:"processor_version"`
	Deleted          bool                     `db:"deleted" json:"deleted"`
}

type TopRecommendedUser struct {
	RecommendedUserID persist.DBID `db:"recommended_user_id" json:"recommended_user_id"`
	Frequency         int64        `db:"frequency" json:"frequency"`
	LastUpdated       interface{}  `db:"last_updated" json:"last_updated"`
}

type User struct {
	ID                   persist.DBID                     `db:"id" json:"id"`
	Deleted              bool                             `db:"deleted" json:"deleted"`
	Version              sql.NullInt32                    `db:"version" json:"version"`
	LastUpdated          time.Time                        `db:"last_updated" json:"last_updated"`
	CreatedAt            time.Time                        `db:"created_at" json:"created_at"`
	Username             sql.NullString                   `db:"username" json:"username"`
	UsernameIdempotent   sql.NullString                   `db:"username_idempotent" json:"username_idempotent"`
	Wallets              persist.WalletList               `db:"wallets" json:"wallets"`
	Bio                  sql.NullString                   `db:"bio" json:"bio"`
	Traits               pgtype.JSONB                     `db:"traits" json:"traits"`
	Universal            bool                             `db:"universal" json:"universal"`
	NotificationSettings persist.UserNotificationSettings `db:"notification_settings" json:"notification_settings"`
	EmailVerified        persist.EmailVerificationStatus  `db:"email_verified" json:"email_verified"`
	EmailUnsubscriptions persist.EmailUnsubscriptions     `db:"email_unsubscriptions" json:"email_unsubscriptions"`
	FeaturedGallery      *persist.DBID                    `db:"featured_gallery" json:"featured_gallery"`
	PrimaryWalletID      persist.DBID                     `db:"primary_wallet_id" json:"primary_wallet_id"`
	UserExperiences      pgtype.JSONB                     `db:"user_experiences" json:"user_experiences"`
	ProfileImageID       persist.DBID                     `db:"profile_image_id" json:"profile_image_id"`
}

type UserRelevance struct {
	ID    persist.DBID `db:"id" json:"id"`
	Score int32        `db:"score" json:"score"`
}

type UserRole struct {
	ID          persist.DBID `db:"id" json:"id"`
	UserID      persist.DBID `db:"user_id" json:"user_id"`
	Role        persist.Role `db:"role" json:"role"`
	Version     int32        `db:"version" json:"version"`
	Deleted     bool         `db:"deleted" json:"deleted"`
	CreatedAt   time.Time    `db:"created_at" json:"created_at"`
	LastUpdated time.Time    `db:"last_updated" json:"last_updated"`
}

type Wallet struct {
	ID          persist.DBID       `db:"id" json:"id"`
	CreatedAt   time.Time          `db:"created_at" json:"created_at"`
	LastUpdated time.Time          `db:"last_updated" json:"last_updated"`
	Deleted     bool               `db:"deleted" json:"deleted"`
	Version     sql.NullInt32      `db:"version" json:"version"`
	Address     persist.Address    `db:"address" json:"address"`
	WalletType  persist.WalletType `db:"wallet_type" json:"wallet_type"`
	Chain       persist.Chain      `db:"chain" json:"chain"`
	L1Chain     persist.L1Chain    `db:"l1_chain" json:"l1_chain"`
}
