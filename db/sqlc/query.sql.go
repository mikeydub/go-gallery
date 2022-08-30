// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: query.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"

	"github.com/mikeydub/go-gallery/service/persist"
)

const createCollectionEvent = `-- name: CreateCollectionEvent :one
INSERT INTO events (id, actor_id, action, resource_type_id, collection_id, subject_id, data) VALUES ($1, $2, $3, $4, $5, $5, $6) RETURNING id, version, actor_id, resource_type_id, subject_id, user_id, token_id, collection_id, action, data, deleted, last_updated, created_at
`

type CreateCollectionEventParams struct {
	ID             persist.DBID
	ActorID        persist.DBID
	Action         persist.Action
	ResourceTypeID persist.ResourceType
	CollectionID   persist.DBID
	Data           persist.EventData
}

func (q *Queries) CreateCollectionEvent(ctx context.Context, arg CreateCollectionEventParams) (Event, error) {
	row := q.db.QueryRow(ctx, createCollectionEvent,
		arg.ID,
		arg.ActorID,
		arg.Action,
		arg.ResourceTypeID,
		arg.CollectionID,
		arg.Data,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Version,
		&i.ActorID,
		&i.ResourceTypeID,
		&i.SubjectID,
		&i.UserID,
		&i.TokenID,
		&i.CollectionID,
		&i.Action,
		&i.Data,
		&i.Deleted,
		&i.LastUpdated,
		&i.CreatedAt,
	)
	return i, err
}

const createFeedEvent = `-- name: CreateFeedEvent :one
INSERT INTO feed_events (id, owner_id, action, data, event_time, event_ids) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, version, owner_id, action, data, event_time, event_ids, deleted, last_updated, created_at
`

type CreateFeedEventParams struct {
	ID        persist.DBID
	OwnerID   persist.DBID
	Action    persist.Action
	Data      persist.FeedEventData
	EventTime time.Time
	EventIds  persist.DBIDList
}

func (q *Queries) CreateFeedEvent(ctx context.Context, arg CreateFeedEventParams) (FeedEvent, error) {
	row := q.db.QueryRow(ctx, createFeedEvent,
		arg.ID,
		arg.OwnerID,
		arg.Action,
		arg.Data,
		arg.EventTime,
		arg.EventIds,
	)
	var i FeedEvent
	err := row.Scan(
		&i.ID,
		&i.Version,
		&i.OwnerID,
		&i.Action,
		&i.Data,
		&i.EventTime,
		&i.EventIds,
		&i.Deleted,
		&i.LastUpdated,
		&i.CreatedAt,
	)
	return i, err
}

const createTokenEvent = `-- name: CreateTokenEvent :one
INSERT INTO events (id, actor_id, action, resource_type_id, token_id, subject_id, data) VALUES ($1, $2, $3, $4, $5, $5, $6) RETURNING id, version, actor_id, resource_type_id, subject_id, user_id, token_id, collection_id, action, data, deleted, last_updated, created_at
`

type CreateTokenEventParams struct {
	ID             persist.DBID
	ActorID        persist.DBID
	Action         persist.Action
	ResourceTypeID persist.ResourceType
	TokenID        persist.DBID
	Data           persist.EventData
}

func (q *Queries) CreateTokenEvent(ctx context.Context, arg CreateTokenEventParams) (Event, error) {
	row := q.db.QueryRow(ctx, createTokenEvent,
		arg.ID,
		arg.ActorID,
		arg.Action,
		arg.ResourceTypeID,
		arg.TokenID,
		arg.Data,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Version,
		&i.ActorID,
		&i.ResourceTypeID,
		&i.SubjectID,
		&i.UserID,
		&i.TokenID,
		&i.CollectionID,
		&i.Action,
		&i.Data,
		&i.Deleted,
		&i.LastUpdated,
		&i.CreatedAt,
	)
	return i, err
}

const createUserEvent = `-- name: CreateUserEvent :one
INSERT INTO events (id, actor_id, action, resource_type_id, user_id, subject_id, data) VALUES ($1, $2, $3, $4, $5, $5, $6) RETURNING id, version, actor_id, resource_type_id, subject_id, user_id, token_id, collection_id, action, data, deleted, last_updated, created_at
`

type CreateUserEventParams struct {
	ID             persist.DBID
	ActorID        persist.DBID
	Action         persist.Action
	ResourceTypeID persist.ResourceType
	UserID         persist.DBID
	Data           persist.EventData
}

func (q *Queries) CreateUserEvent(ctx context.Context, arg CreateUserEventParams) (Event, error) {
	row := q.db.QueryRow(ctx, createUserEvent,
		arg.ID,
		arg.ActorID,
		arg.Action,
		arg.ResourceTypeID,
		arg.UserID,
		arg.Data,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Version,
		&i.ActorID,
		&i.ResourceTypeID,
		&i.SubjectID,
		&i.UserID,
		&i.TokenID,
		&i.CollectionID,
		&i.Action,
		&i.Data,
		&i.Deleted,
		&i.LastUpdated,
		&i.CreatedAt,
	)
	return i, err
}

const getCollectionById = `-- name: GetCollectionById :one
SELECT id, deleted, owner_user_id, nfts, version, last_updated, created_at, hidden, collectors_note, name, layout, token_settings FROM collections WHERE id = $1 AND deleted = false
`

func (q *Queries) GetCollectionById(ctx context.Context, id persist.DBID) (Collection, error) {
	row := q.db.QueryRow(ctx, getCollectionById, id)
	var i Collection
	err := row.Scan(
		&i.ID,
		&i.Deleted,
		&i.OwnerUserID,
		&i.Nfts,
		&i.Version,
		&i.LastUpdated,
		&i.CreatedAt,
		&i.Hidden,
		&i.CollectorsNote,
		&i.Name,
		&i.Layout,
		&i.TokenSettings,
	)
	return i, err
}

const getCollectionsByGalleryId = `-- name: GetCollectionsByGalleryId :many
SELECT c.id, c.deleted, c.owner_user_id, c.nfts, c.version, c.last_updated, c.created_at, c.hidden, c.collectors_note, c.name, c.layout, c.token_settings FROM galleries g, unnest(g.collections)
    WITH ORDINALITY AS x(coll_id, coll_ord)
    INNER JOIN collections c ON c.id = x.coll_id
    WHERE g.id = $1 AND g.deleted = false AND c.deleted = false ORDER BY x.coll_ord
`

func (q *Queries) GetCollectionsByGalleryId(ctx context.Context, id persist.DBID) ([]Collection, error) {
	rows, err := q.db.Query(ctx, getCollectionsByGalleryId, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Collection
	for rows.Next() {
		var i Collection
		if err := rows.Scan(
			&i.ID,
			&i.Deleted,
			&i.OwnerUserID,
			&i.Nfts,
			&i.Version,
			&i.LastUpdated,
			&i.CreatedAt,
			&i.Hidden,
			&i.CollectorsNote,
			&i.Name,
			&i.Layout,
			&i.TokenSettings,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getContractByChainAddress = `-- name: GetContractByChainAddress :one
select id, deleted, version, created_at, last_updated, name, symbol, address, creator_address, chain, profile_banner_url, profile_image_url, badge_url FROM contracts WHERE address = $1 AND chain = $2 AND deleted = false
`

type GetContractByChainAddressParams struct {
	Address persist.Address
	Chain   sql.NullInt32
}

func (q *Queries) GetContractByChainAddress(ctx context.Context, arg GetContractByChainAddressParams) (Contract, error) {
	row := q.db.QueryRow(ctx, getContractByChainAddress, arg.Address, arg.Chain)
	var i Contract
	err := row.Scan(
		&i.ID,
		&i.Deleted,
		&i.Version,
		&i.CreatedAt,
		&i.LastUpdated,
		&i.Name,
		&i.Symbol,
		&i.Address,
		&i.CreatorAddress,
		&i.Chain,
		&i.ProfileBannerUrl,
		&i.ProfileImageUrl,
		&i.BadgeUrl,
	)
	return i, err
}

const getContractByID = `-- name: GetContractByID :one
select id, deleted, version, created_at, last_updated, name, symbol, address, creator_address, chain, profile_banner_url, profile_image_url, badge_url FROM contracts WHERE id = $1 AND deleted = false
`

func (q *Queries) GetContractByID(ctx context.Context, id persist.DBID) (Contract, error) {
	row := q.db.QueryRow(ctx, getContractByID, id)
	var i Contract
	err := row.Scan(
		&i.ID,
		&i.Deleted,
		&i.Version,
		&i.CreatedAt,
		&i.LastUpdated,
		&i.Name,
		&i.Symbol,
		&i.Address,
		&i.CreatorAddress,
		&i.Chain,
		&i.ProfileBannerUrl,
		&i.ProfileImageUrl,
		&i.BadgeUrl,
	)
	return i, err
}

const getContractsByUserID = `-- name: GetContractsByUserID :many
SELECT DISTINCT ON (contracts.id) contracts.id, contracts.deleted, contracts.version, contracts.created_at, contracts.last_updated, contracts.name, contracts.symbol, contracts.address, contracts.creator_address, contracts.chain, contracts.profile_banner_url, contracts.profile_image_url, contracts.badge_url FROM contracts, tokens
    WHERE tokens.owner_user_id = $1 AND tokens.contract = contracts.id
    AND tokens.deleted = false AND contracts.deleted = false
`

func (q *Queries) GetContractsByUserID(ctx context.Context, ownerUserID persist.DBID) ([]Contract, error) {
	rows, err := q.db.Query(ctx, getContractsByUserID, ownerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Contract
	for rows.Next() {
		var i Contract
		if err := rows.Scan(
			&i.ID,
			&i.Deleted,
			&i.Version,
			&i.CreatedAt,
			&i.LastUpdated,
			&i.Name,
			&i.Symbol,
			&i.Address,
			&i.CreatorAddress,
			&i.Chain,
			&i.ProfileBannerUrl,
			&i.ProfileImageUrl,
			&i.BadgeUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEvent = `-- name: GetEvent :one
SELECT id, version, actor_id, resource_type_id, subject_id, user_id, token_id, collection_id, action, data, deleted, last_updated, created_at FROM events WHERE id = $1 AND deleted = false
`

func (q *Queries) GetEvent(ctx context.Context, id persist.DBID) (Event, error) {
	row := q.db.QueryRow(ctx, getEvent, id)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Version,
		&i.ActorID,
		&i.ResourceTypeID,
		&i.SubjectID,
		&i.UserID,
		&i.TokenID,
		&i.CollectionID,
		&i.Action,
		&i.Data,
		&i.Deleted,
		&i.LastUpdated,
		&i.CreatedAt,
	)
	return i, err
}

const getEventsInWindow = `-- name: GetEventsInWindow :many
WITH RECURSIVE activity AS (
    SELECT id, version, actor_id, resource_type_id, subject_id, user_id, token_id, collection_id, action, data, deleted, last_updated, created_at FROM events WHERE events.id = $1 AND deleted = false
    UNION
    SELECT e.id, e.version, e.actor_id, e.resource_type_id, e.subject_id, e.user_id, e.token_id, e.collection_id, e.action, e.data, e.deleted, e.last_updated, e.created_at FROM events e, activity a
    WHERE e.actor_id = a.actor_id
        AND e.action = a.action
        AND e.created_at < a.created_at
        AND e.created_at >= a.created_at - make_interval(secs => $2)
        AND e.deleted = false
)
SELECT id, version, actor_id, resource_type_id, subject_id, user_id, token_id, collection_id, action, data, deleted, last_updated, created_at FROM events WHERE id = ANY(SELECT id FROM activity) ORDER BY created_at DESC
`

type GetEventsInWindowParams struct {
	ID   persist.DBID
	Secs float64
}

func (q *Queries) GetEventsInWindow(ctx context.Context, arg GetEventsInWindowParams) ([]Event, error) {
	rows, err := q.db.Query(ctx, getEventsInWindow, arg.ID, arg.Secs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Version,
			&i.ActorID,
			&i.ResourceTypeID,
			&i.SubjectID,
			&i.UserID,
			&i.TokenID,
			&i.CollectionID,
			&i.Action,
			&i.Data,
			&i.Deleted,
			&i.LastUpdated,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGalleriesByUserId = `-- name: GetGalleriesByUserId :many
SELECT id, deleted, last_updated, created_at, version, owner_user_id, collections FROM galleries WHERE owner_user_id = $1 AND deleted = false
`

func (q *Queries) GetGalleriesByUserId(ctx context.Context, ownerUserID persist.DBID) ([]Gallery, error) {
	rows, err := q.db.Query(ctx, getGalleriesByUserId, ownerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Gallery
	for rows.Next() {
		var i Gallery
		if err := rows.Scan(
			&i.ID,
			&i.Deleted,
			&i.LastUpdated,
			&i.CreatedAt,
			&i.Version,
			&i.OwnerUserID,
			&i.Collections,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getGalleryByCollectionId = `-- name: GetGalleryByCollectionId :one
SELECT g.id, g.deleted, g.last_updated, g.created_at, g.version, g.owner_user_id, g.collections FROM galleries g, collections c WHERE c.id = $1 AND c.deleted = false AND $1 = ANY(g.collections) AND g.deleted = false
`

func (q *Queries) GetGalleryByCollectionId(ctx context.Context, id persist.DBID) (Gallery, error) {
	row := q.db.QueryRow(ctx, getGalleryByCollectionId, id)
	var i Gallery
	err := row.Scan(
		&i.ID,
		&i.Deleted,
		&i.LastUpdated,
		&i.CreatedAt,
		&i.Version,
		&i.OwnerUserID,
		&i.Collections,
	)
	return i, err
}

const getGalleryById = `-- name: GetGalleryById :one
SELECT id, deleted, last_updated, created_at, version, owner_user_id, collections FROM galleries WHERE id = $1 AND deleted = false
`

func (q *Queries) GetGalleryById(ctx context.Context, id persist.DBID) (Gallery, error) {
	row := q.db.QueryRow(ctx, getGalleryById, id)
	var i Gallery
	err := row.Scan(
		&i.ID,
		&i.Deleted,
		&i.LastUpdated,
		&i.CreatedAt,
		&i.Version,
		&i.OwnerUserID,
		&i.Collections,
	)
	return i, err
}

const getLastFeedEvent = `-- name: GetLastFeedEvent :one
SELECT id, version, owner_id, action, data, event_time, event_ids, deleted, last_updated, created_at FROM feed_events
    WHERE owner_id = $1 AND action = $2 AND event_time < $3 AND deleted = false
    ORDER BY event_time DESC
    LIMIT 1
`

type GetLastFeedEventParams struct {
	OwnerID   persist.DBID
	Action    persist.Action
	EventTime time.Time
}

func (q *Queries) GetLastFeedEvent(ctx context.Context, arg GetLastFeedEventParams) (FeedEvent, error) {
	row := q.db.QueryRow(ctx, getLastFeedEvent, arg.OwnerID, arg.Action, arg.EventTime)
	var i FeedEvent
	err := row.Scan(
		&i.ID,
		&i.Version,
		&i.OwnerID,
		&i.Action,
		&i.Data,
		&i.EventTime,
		&i.EventIds,
		&i.Deleted,
		&i.LastUpdated,
		&i.CreatedAt,
	)
	return i, err
}

const getLastFeedEventForCollection = `-- name: GetLastFeedEventForCollection :one
SELECT id, version, owner_id, action, data, event_time, event_ids, deleted, last_updated, created_at FROM feed_events
    WHERE owner_id = $1 and action = $2 AND data ->> 'collection_id' = $4::varchar AND event_time < $3 AND deleted = false
    ORDER BY event_time DESC
    LIMIT 1
`

type GetLastFeedEventForCollectionParams struct {
	OwnerID      persist.DBID
	Action       persist.Action
	EventTime    time.Time
	CollectionID string
}

func (q *Queries) GetLastFeedEventForCollection(ctx context.Context, arg GetLastFeedEventForCollectionParams) (FeedEvent, error) {
	row := q.db.QueryRow(ctx, getLastFeedEventForCollection,
		arg.OwnerID,
		arg.Action,
		arg.EventTime,
		arg.CollectionID,
	)
	var i FeedEvent
	err := row.Scan(
		&i.ID,
		&i.Version,
		&i.OwnerID,
		&i.Action,
		&i.Data,
		&i.EventTime,
		&i.EventIds,
		&i.Deleted,
		&i.LastUpdated,
		&i.CreatedAt,
	)
	return i, err
}

const getLastFeedEventForToken = `-- name: GetLastFeedEventForToken :one
SELECT id, version, owner_id, action, data, event_time, event_ids, deleted, last_updated, created_at FROM feed_events
    WHERE owner_id = $1 and action = $2 AND data ->> 'token_id' = $4::varchar AND event_time < $3 AND deleted = false
    ORDER BY event_time DESC
    LIMIT 1
`

type GetLastFeedEventForTokenParams struct {
	OwnerID   persist.DBID
	Action    persist.Action
	EventTime time.Time
	TokenID   string
}

func (q *Queries) GetLastFeedEventForToken(ctx context.Context, arg GetLastFeedEventForTokenParams) (FeedEvent, error) {
	row := q.db.QueryRow(ctx, getLastFeedEventForToken,
		arg.OwnerID,
		arg.Action,
		arg.EventTime,
		arg.TokenID,
	)
	var i FeedEvent
	err := row.Scan(
		&i.ID,
		&i.Version,
		&i.OwnerID,
		&i.Action,
		&i.Data,
		&i.EventTime,
		&i.EventIds,
		&i.Deleted,
		&i.LastUpdated,
		&i.CreatedAt,
	)
	return i, err
}

const getMembershipByMembershipId = `-- name: GetMembershipByMembershipId :one
SELECT id, deleted, version, created_at, last_updated, token_id, name, asset_url, owners FROM membership WHERE id = $1 AND deleted = false
`

func (q *Queries) GetMembershipByMembershipId(ctx context.Context, id persist.DBID) (Membership, error) {
	row := q.db.QueryRow(ctx, getMembershipByMembershipId, id)
	var i Membership
	err := row.Scan(
		&i.ID,
		&i.Deleted,
		&i.Version,
		&i.CreatedAt,
		&i.LastUpdated,
		&i.TokenID,
		&i.Name,
		&i.AssetUrl,
		&i.Owners,
	)
	return i, err
}

const getTokenById = `-- name: GetTokenById :one
SELECT id, deleted, version, created_at, last_updated, name, description, collectors_note, media, token_uri, token_type, token_id, quantity, ownership_history, token_metadata, external_url, block_number, owner_user_id, owned_by_wallets, chain, contract, is_user_marked_spam, is_provider_marked_spam FROM tokens WHERE id = $1 AND deleted = false
`

func (q *Queries) GetTokenById(ctx context.Context, id persist.DBID) (Token, error) {
	row := q.db.QueryRow(ctx, getTokenById, id)
	var i Token
	err := row.Scan(
		&i.ID,
		&i.Deleted,
		&i.Version,
		&i.CreatedAt,
		&i.LastUpdated,
		&i.Name,
		&i.Description,
		&i.CollectorsNote,
		&i.Media,
		&i.TokenUri,
		&i.TokenType,
		&i.TokenID,
		&i.Quantity,
		&i.OwnershipHistory,
		&i.TokenMetadata,
		&i.ExternalUrl,
		&i.BlockNumber,
		&i.OwnerUserID,
		&i.OwnedByWallets,
		&i.Chain,
		&i.Contract,
		&i.IsUserMarkedSpam,
		&i.IsProviderMarkedSpam,
	)
	return i, err
}

const getTokensByCollectionId = `-- name: GetTokensByCollectionId :many
SELECT t.id, t.deleted, t.version, t.created_at, t.last_updated, t.name, t.description, t.collectors_note, t.media, t.token_uri, t.token_type, t.token_id, t.quantity, t.ownership_history, t.token_metadata, t.external_url, t.block_number, t.owner_user_id, t.owned_by_wallets, t.chain, t.contract, t.is_user_marked_spam, t.is_provider_marked_spam FROM users u, collections c, unnest(c.nfts)
    WITH ORDINALITY AS x(nft_id, nft_ord)
    INNER JOIN tokens t ON t.id = x.nft_id
    WHERE u.id = t.owner_user_id AND t.owned_by_wallets && u.wallets
    AND c.id = $1 AND u.deleted = false AND c.deleted = false AND t.deleted = false ORDER BY x.nft_ord
`

func (q *Queries) GetTokensByCollectionId(ctx context.Context, id persist.DBID) ([]Token, error) {
	rows, err := q.db.Query(ctx, getTokensByCollectionId, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Token
	for rows.Next() {
		var i Token
		if err := rows.Scan(
			&i.ID,
			&i.Deleted,
			&i.Version,
			&i.CreatedAt,
			&i.LastUpdated,
			&i.Name,
			&i.Description,
			&i.CollectorsNote,
			&i.Media,
			&i.TokenUri,
			&i.TokenType,
			&i.TokenID,
			&i.Quantity,
			&i.OwnershipHistory,
			&i.TokenMetadata,
			&i.ExternalUrl,
			&i.BlockNumber,
			&i.OwnerUserID,
			&i.OwnedByWallets,
			&i.Chain,
			&i.Contract,
			&i.IsUserMarkedSpam,
			&i.IsProviderMarkedSpam,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTokensByUserId = `-- name: GetTokensByUserId :many
SELECT tokens.id, tokens.deleted, tokens.version, tokens.created_at, tokens.last_updated, tokens.name, tokens.description, tokens.collectors_note, tokens.media, tokens.token_uri, tokens.token_type, tokens.token_id, tokens.quantity, tokens.ownership_history, tokens.token_metadata, tokens.external_url, tokens.block_number, tokens.owner_user_id, tokens.owned_by_wallets, tokens.chain, tokens.contract, tokens.is_user_marked_spam, tokens.is_provider_marked_spam FROM tokens, users
    WHERE tokens.owner_user_id = $1 AND users.id = $1
      AND tokens.owned_by_wallets && users.wallets
      AND tokens.deleted = false AND users.deleted = false
    ORDER BY tokens.created_at DESC, tokens.name DESC, tokens.id DESC
`

func (q *Queries) GetTokensByUserId(ctx context.Context, ownerUserID persist.DBID) ([]Token, error) {
	rows, err := q.db.Query(ctx, getTokensByUserId, ownerUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Token
	for rows.Next() {
		var i Token
		if err := rows.Scan(
			&i.ID,
			&i.Deleted,
			&i.Version,
			&i.CreatedAt,
			&i.LastUpdated,
			&i.Name,
			&i.Description,
			&i.CollectorsNote,
			&i.Media,
			&i.TokenUri,
			&i.TokenType,
			&i.TokenID,
			&i.Quantity,
			&i.OwnershipHistory,
			&i.TokenMetadata,
			&i.ExternalUrl,
			&i.BlockNumber,
			&i.OwnerUserID,
			&i.OwnedByWallets,
			&i.Chain,
			&i.Contract,
			&i.IsUserMarkedSpam,
			&i.IsProviderMarkedSpam,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTokensByWalletIds = `-- name: GetTokensByWalletIds :many
SELECT id, deleted, version, created_at, last_updated, name, description, collectors_note, media, token_uri, token_type, token_id, quantity, ownership_history, token_metadata, external_url, block_number, owner_user_id, owned_by_wallets, chain, contract, is_user_marked_spam, is_provider_marked_spam FROM tokens WHERE owned_by_wallets && $1 AND deleted = false
    ORDER BY tokens.created_at DESC, tokens.name DESC, tokens.id DESC
`

func (q *Queries) GetTokensByWalletIds(ctx context.Context, ownedByWallets persist.DBIDList) ([]Token, error) {
	rows, err := q.db.Query(ctx, getTokensByWalletIds, ownedByWallets)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Token
	for rows.Next() {
		var i Token
		if err := rows.Scan(
			&i.ID,
			&i.Deleted,
			&i.Version,
			&i.CreatedAt,
			&i.LastUpdated,
			&i.Name,
			&i.Description,
			&i.CollectorsNote,
			&i.Media,
			&i.TokenUri,
			&i.TokenType,
			&i.TokenID,
			&i.Quantity,
			&i.OwnershipHistory,
			&i.TokenMetadata,
			&i.ExternalUrl,
			&i.BlockNumber,
			&i.OwnerUserID,
			&i.OwnedByWallets,
			&i.Chain,
			&i.Contract,
			&i.IsUserMarkedSpam,
			&i.IsProviderMarkedSpam,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserById = `-- name: GetUserById :one
SELECT id, deleted, version, last_updated, created_at, username, username_idempotent, wallets, bio, traits FROM users WHERE id = $1 AND deleted = false
`

func (q *Queries) GetUserById(ctx context.Context, id persist.DBID) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Deleted,
		&i.Version,
		&i.LastUpdated,
		&i.CreatedAt,
		&i.Username,
		&i.UsernameIdempotent,
		&i.Wallets,
		&i.Bio,
		&i.Traits,
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, deleted, version, last_updated, created_at, username, username_idempotent, wallets, bio, traits FROM users WHERE username_idempotent = lower($1) AND deleted = false
`

func (q *Queries) GetUserByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Deleted,
		&i.Version,
		&i.LastUpdated,
		&i.CreatedAt,
		&i.Username,
		&i.UsernameIdempotent,
		&i.Wallets,
		&i.Bio,
		&i.Traits,
	)
	return i, err
}

const getUsersWithTrait = `-- name: GetUsersWithTrait :many
SELECT id, deleted, version, last_updated, created_at, username, username_idempotent, wallets, bio, traits FROM users WHERE (traits->$1::string) IS NOT NULL AND deleted = false
`

func (q *Queries) GetUsersWithTrait(ctx context.Context, dollar_1 string) ([]User, error) {
	rows, err := q.db.Query(ctx, getUsersWithTrait, dollar_1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Deleted,
			&i.Version,
			&i.LastUpdated,
			&i.CreatedAt,
			&i.Username,
			&i.UsernameIdempotent,
			&i.Wallets,
			&i.Bio,
			&i.Traits,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getWalletByChainAddress = `-- name: GetWalletByChainAddress :one
SELECT wallets.id, wallets.created_at, wallets.last_updated, wallets.deleted, wallets.version, wallets.address, wallets.wallet_type, wallets.chain FROM wallets WHERE address = $1 AND chain = $2 AND deleted = false
`

type GetWalletByChainAddressParams struct {
	Address persist.Address
	Chain   sql.NullInt32
}

func (q *Queries) GetWalletByChainAddress(ctx context.Context, arg GetWalletByChainAddressParams) (Wallet, error) {
	row := q.db.QueryRow(ctx, getWalletByChainAddress, arg.Address, arg.Chain)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.LastUpdated,
		&i.Deleted,
		&i.Version,
		&i.Address,
		&i.WalletType,
		&i.Chain,
	)
	return i, err
}

const getWalletByID = `-- name: GetWalletByID :one
SELECT id, created_at, last_updated, deleted, version, address, wallet_type, chain FROM wallets WHERE id = $1 AND deleted = false
`

func (q *Queries) GetWalletByID(ctx context.Context, id persist.DBID) (Wallet, error) {
	row := q.db.QueryRow(ctx, getWalletByID, id)
	var i Wallet
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.LastUpdated,
		&i.Deleted,
		&i.Version,
		&i.Address,
		&i.WalletType,
		&i.Chain,
	)
	return i, err
}

const getWalletsByUserID = `-- name: GetWalletsByUserID :many
SELECT w.id, w.created_at, w.last_updated, w.deleted, w.version, w.address, w.wallet_type, w.chain FROM users u, unnest(u.wallets) WITH ORDINALITY AS a(wallet_id, wallet_ord)INNER JOIN wallets w on w.id = a.wallet_id WHERE u.id = $1 AND u.deleted = false AND w.deleted = false ORDER BY a.wallet_ord
`

func (q *Queries) GetWalletsByUserID(ctx context.Context, id persist.DBID) ([]Wallet, error) {
	rows, err := q.db.Query(ctx, getWalletsByUserID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Wallet
	for rows.Next() {
		var i Wallet
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.LastUpdated,
			&i.Deleted,
			&i.Version,
			&i.Address,
			&i.WalletType,
			&i.Chain,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const globalFeedHasMoreEvents = `-- name: GlobalFeedHasMoreEvents :one
SELECT
    CASE WHEN $2::bool
    THEN EXISTS(
        SELECT 1
        FROM feed_events
        WHERE event_time > (SELECT event_time FROM feed_events f WHERE f.id = $1)
        AND deleted = false
        LIMIT 1
    )
    ELSE EXISTS(
        SELECT 1
        FROM feed_events
        WHERE event_time < (SELECT event_time FROM feed_events f WHERE f.id = $1)
        AND deleted = false
        LIMIT 1)
    END::bool
`

type GlobalFeedHasMoreEventsParams struct {
	ID        persist.DBID
	FromFirst bool
}

func (q *Queries) GlobalFeedHasMoreEvents(ctx context.Context, arg GlobalFeedHasMoreEventsParams) (bool, error) {
	row := q.db.QueryRow(ctx, globalFeedHasMoreEvents, arg.ID, arg.FromFirst)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}

const isFeedUserActionBlocked = `-- name: IsFeedUserActionBlocked :one
SELECT EXISTS(SELECT 1 FROM feed_blocklist WHERE user_id = $1 AND action = $2 AND deleted = false)
`

type IsFeedUserActionBlockedParams struct {
	UserID persist.DBID
	Action persist.Action
}

func (q *Queries) IsFeedUserActionBlocked(ctx context.Context, arg IsFeedUserActionBlockedParams) (bool, error) {
	row := q.db.QueryRow(ctx, isFeedUserActionBlocked, arg.UserID, arg.Action)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const isWindowActive = `-- name: IsWindowActive :one
SELECT EXISTS(
    SELECT 1 FROM events
    WHERE actor_id = $1 AND action = $2 AND deleted = false
    AND created_at > $3 AND created_at <= $4
    LIMIT 1
)
`

type IsWindowActiveParams struct {
	ActorID     persist.DBID
	Action      persist.Action
	WindowStart time.Time
	WindowEnd   time.Time
}

func (q *Queries) IsWindowActive(ctx context.Context, arg IsWindowActiveParams) (bool, error) {
	row := q.db.QueryRow(ctx, isWindowActive,
		arg.ActorID,
		arg.Action,
		arg.WindowStart,
		arg.WindowEnd,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const isWindowActiveWithSubject = `-- name: IsWindowActiveWithSubject :one
SELECT EXISTS(
    SELECT 1 FROM events
    WHERE actor_id = $1 AND action = $2 AND subject_id = $3 AND deleted = false
    AND created_at > $4 AND created_at <= $5
    LIMIT 1
)
`

type IsWindowActiveWithSubjectParams struct {
	ActorID     persist.DBID
	Action      persist.Action
	SubjectID   persist.DBID
	WindowStart time.Time
	WindowEnd   time.Time
}

func (q *Queries) IsWindowActiveWithSubject(ctx context.Context, arg IsWindowActiveWithSubjectParams) (bool, error) {
	row := q.db.QueryRow(ctx, isWindowActiveWithSubject,
		arg.ActorID,
		arg.Action,
		arg.SubjectID,
		arg.WindowStart,
		arg.WindowEnd,
	)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const userFeedHasMoreEvents = `-- name: UserFeedHasMoreEvents :one
SELECT
    CASE WHEN $3::bool
    THEN EXISTS(
        SELECT 1
        FROM feed_events fe
        INNER JOIN follows fl ON fe.owner_id = fl.followee AND fl.follower = $1
        WHERE event_time > (SELECT event_time FROM feed_events f WHERE f.id = $2)
        AND fe.deleted = false AND fl.deleted = false
        LIMIT 1)
    ELSE EXISTS(
        SELECT 1
        FROM feed_events fe
        INNER JOIN follows fl ON fe.owner_id = fl.followee AND fl.follower = $1
        WHERE event_time < (SELECT event_time FROM feed_events f WHERE f.id = $2)
        AND fe.deleted = false AND fl.deleted = false
        LIMIT 1
    )
    END::bool
`

type UserFeedHasMoreEventsParams struct {
	Follower  persist.DBID
	ID        persist.DBID
	FromFirst bool
}

func (q *Queries) UserFeedHasMoreEvents(ctx context.Context, arg UserFeedHasMoreEventsParams) (bool, error) {
	row := q.db.QueryRow(ctx, userFeedHasMoreEvents, arg.Follower, arg.ID, arg.FromFirst)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}
