// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package coredb

import (
	"database/sql"
	"time"

	"github.com/jackc/pgtype"
	"github.com/mikeydub/go-gallery/service/persist"
)

type Admire struct {
	ID          persist.DBID
	Version     int32
	FeedEventID persist.DBID
	ActorID     persist.DBID
	Deleted     bool
	CreatedAt   time.Time
	LastUpdated time.Time
}

type Collection struct {
	ID             persist.DBID
	Deleted        bool
	OwnerUserID    persist.DBID
	Nfts           persist.DBIDList
	Version        sql.NullInt32
	LastUpdated    time.Time
	CreatedAt      time.Time
	Hidden         bool
	CollectorsNote sql.NullString
	Name           sql.NullString
	Layout         persist.TokenLayout
	TokenSettings  map[persist.DBID]persist.CollectionTokenSettings
	GalleryID      persist.DBID
}

type Comment struct {
	ID          persist.DBID
	Version     int32
	FeedEventID persist.DBID
	ActorID     persist.DBID
	ReplyTo     persist.DBID
	Comment     string
	Deleted     bool
	CreatedAt   time.Time
	LastUpdated time.Time
}

type Contract struct {
	ID               persist.DBID
	Deleted          bool
	Version          sql.NullInt32
	CreatedAt        time.Time
	LastUpdated      time.Time
	Name             sql.NullString
	Symbol           sql.NullString
	Address          persist.Address
	CreatorAddress   persist.Address
	Chain            persist.Chain
	ProfileBannerUrl sql.NullString
	ProfileImageUrl  sql.NullString
	BadgeUrl         sql.NullString
	Description      sql.NullString
	OwnerAddress     persist.Address
}

type ContractRelevance struct {
	ID    persist.DBID
	Score int32
}

type DevMetadataUser struct {
	UserID          persist.DBID
	HasEmailAddress persist.Email
	Deleted         bool
}

type EarlyAccess struct {
	Address persist.Address
}

type Event struct {
	ID             persist.DBID
	Version        int32
	ActorID        sql.NullString
	ResourceTypeID persist.ResourceType
	SubjectID      persist.DBID
	UserID         persist.DBID
	TokenID        persist.DBID
	CollectionID   persist.DBID
	Action         persist.Action
	Data           persist.EventData
	Deleted        bool
	LastUpdated    time.Time
	CreatedAt      time.Time
	GalleryID      persist.DBID
	CommentID      persist.DBID
	AdmireID       persist.DBID
	FeedEventID    persist.DBID
	ExternalID     sql.NullString
	Caption        sql.NullString
	GroupID        sql.NullString
}

type FeedBlocklist struct {
	ID          persist.DBID
	UserID      persist.DBID
	Action      persist.Action
	LastUpdated time.Time
	CreatedAt   time.Time
	Deleted     bool
}

type FeedEvent struct {
	ID          persist.DBID
	Version     int32
	OwnerID     persist.DBID
	Action      persist.Action
	Data        persist.FeedEventData
	EventTime   time.Time
	EventIds    persist.DBIDList
	Deleted     bool
	LastUpdated time.Time
	CreatedAt   time.Time
	Caption     sql.NullString
	GroupID     sql.NullString
}

type Follow struct {
	ID          persist.DBID
	Follower    persist.DBID
	Followee    persist.DBID
	Deleted     bool
	CreatedAt   time.Time
	LastUpdated time.Time
}

type Gallery struct {
	ID          persist.DBID
	Deleted     bool
	LastUpdated time.Time
	CreatedAt   time.Time
	Version     sql.NullInt32
	OwnerUserID persist.DBID
	Collections persist.DBIDList
	Name        string
	Description string
	Hidden      bool
	Position    string
}

type GalleryRelevance struct {
	ID    persist.DBID
	Score int32
}

type LegacyView struct {
	UserID      persist.DBID
	ViewCount   sql.NullInt32
	LastUpdated time.Time
	CreatedAt   time.Time
	Deleted     sql.NullBool
}

type MarketplaceContract struct {
	ContractID persist.DBID
}

type Membership struct {
	ID          persist.DBID
	Deleted     bool
	Version     sql.NullInt32
	CreatedAt   time.Time
	LastUpdated time.Time
	TokenID     persist.DBID
	Name        sql.NullString
	AssetUrl    sql.NullString
	Owners      persist.TokenHolderList
}

type Merch struct {
	ID           persist.DBID
	Deleted      bool
	Version      sql.NullInt32
	CreatedAt    time.Time
	LastUpdated  time.Time
	TokenID      persist.TokenID
	ObjectType   int32
	DiscountCode sql.NullString
	Redeemed     bool
}

type Nonce struct {
	ID          persist.DBID
	Deleted     bool
	Version     sql.NullInt32
	LastUpdated time.Time
	CreatedAt   time.Time
	UserID      persist.DBID
	Address     persist.Address
	Value       sql.NullString
	Chain       persist.Chain
}

type Notification struct {
	ID          persist.DBID
	Deleted     bool
	OwnerID     persist.DBID
	Version     sql.NullInt32
	LastUpdated time.Time
	CreatedAt   time.Time
	Action      persist.Action
	Data        persist.NotificationData
	EventIds    persist.DBIDList
	FeedEventID persist.DBID
	CommentID   persist.DBID
	GalleryID   persist.DBID
	Seen        bool
	Amount      int32
}

type OwnedContract struct {
	UserID         persist.DBID
	UserCreatedAt  time.Time
	ContractID     persist.DBID
	OwnedCount     int64
	DisplayedCount int64
	Displayed      bool
	LastUpdated    time.Time
}

type PiiAccountCreationInfo struct {
	UserID    persist.DBID
	IpAddress string
	CreatedAt time.Time
}

type PiiForUser struct {
	UserID          persist.DBID
	PiiEmailAddress persist.Email
	Deleted         bool
	PiiSocials      persist.Socials
}

type PiiSocialsAuth struct {
	ID           persist.DBID
	Deleted      bool
	Version      sql.NullInt32
	CreatedAt    time.Time
	LastUpdated  time.Time
	UserID       persist.DBID
	Provider     persist.SocialProvider
	AccessToken  sql.NullString
	RefreshToken sql.NullString
}

type PiiUserView struct {
	ID                   persist.DBID
	Deleted              bool
	Version              sql.NullInt32
	LastUpdated          time.Time
	CreatedAt            time.Time
	Username             sql.NullString
	UsernameIdempotent   sql.NullString
	Wallets              persist.WalletList
	Bio                  sql.NullString
	Traits               pgtype.JSONB
	Universal            bool
	NotificationSettings persist.UserNotificationSettings
	EmailVerified        persist.EmailVerificationStatus
	EmailUnsubscriptions persist.EmailUnsubscriptions
	FeaturedGallery      *persist.DBID
	PrimaryWalletID      persist.DBID
	UserExperiences      pgtype.JSONB
	PiiEmailAddress      persist.Email
	PiiSocials           persist.Socials
}

type RecommendationResult struct {
	ID                persist.DBID
	Version           sql.NullInt32
	UserID            persist.DBID
	RecommendedUserID persist.DBID
	RecommendedCount  sql.NullInt32
	CreatedAt         time.Time
	LastUpdated       time.Time
	Deleted           bool
}

type ScrubbedPiiAccountCreationInfo struct {
	UserID    persist.DBID
	IpAddress persist.Address
	CreatedAt time.Time
}

type ScrubbedPiiForUser struct {
	UserID          persist.DBID
	PiiEmailAddress persist.Email
	Deleted         bool
	PiiSocials      persist.Socials
}

type SpamUserScore struct {
	UserID        persist.DBID
	Score         int32
	DecidedIsSpam sql.NullBool
	DecidedAt     sql.NullTime
	Deleted       bool
	CreatedAt     time.Time
}

type Token struct {
	ID                   persist.DBID
	Deleted              bool
	Version              sql.NullInt32
	CreatedAt            time.Time
	LastUpdated          time.Time
	Name                 sql.NullString
	Description          sql.NullString
	CollectorsNote       sql.NullString
	Media                persist.Media
	TokenUri             sql.NullString
	TokenType            sql.NullString
	TokenID              persist.TokenID
	Quantity             sql.NullString
	OwnershipHistory     persist.AddressAtBlockList
	TokenMetadata        persist.TokenMetadata
	ExternalUrl          sql.NullString
	BlockNumber          sql.NullInt64
	OwnerUserID          persist.DBID
	OwnedByWallets       persist.DBIDList
	Chain                persist.Chain
	Contract             persist.DBID
	IsUserMarkedSpam     sql.NullBool
	IsProviderMarkedSpam sql.NullBool
	LastSynced           time.Time
}

type TopRecommendedUser struct {
	RecommendedUserID persist.DBID
	Frequency         int64
	LastUpdated       interface{}
}

type User struct {
	ID                   persist.DBID
	Deleted              bool
	Version              sql.NullInt32
	LastUpdated          time.Time
	CreatedAt            time.Time
	Username             sql.NullString
	UsernameIdempotent   sql.NullString
	Wallets              persist.WalletList
	Bio                  sql.NullString
	Traits               pgtype.JSONB
	Universal            bool
	NotificationSettings persist.UserNotificationSettings
	EmailVerified        persist.EmailVerificationStatus
	EmailUnsubscriptions persist.EmailUnsubscriptions
	FeaturedGallery      *persist.DBID
	PrimaryWalletID      persist.DBID
	UserExperiences      pgtype.JSONB
}

type UserRelevance struct {
	ID    persist.DBID
	Score int32
}

type UserRole struct {
	ID          persist.DBID
	UserID      persist.DBID
	Role        persist.Role
	Version     int32
	Deleted     bool
	CreatedAt   time.Time
	LastUpdated time.Time
}

type Wallet struct {
	ID          persist.DBID
	CreatedAt   time.Time
	LastUpdated time.Time
	Deleted     bool
	Version     sql.NullInt32
	Address     persist.Address
	WalletType  persist.WalletType
	Chain       persist.Chain
}
