// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: query.sql

package sqlc

import (
	"context"
	"database/sql"

	"github.com/mikeydub/go-gallery/service/persist"
)

const getCollectionById = `-- name: GetCollectionById :one
SELECT id, deleted, owner_user_id, nfts, version, last_updated, created_at, hidden, collectors_note, name, layout FROM collections WHERE id = $1 AND deleted = false
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
	)
	return i, err
}

const getCollectionsByGalleryId = `-- name: GetCollectionsByGalleryId :many
SELECT c.id, c.deleted, c.owner_user_id, c.nfts, c.version, c.last_updated, c.created_at, c.hidden, c.collectors_note, c.name, c.layout FROM galleries g, unnest(g.collections)
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
select id, deleted, version, created_at, last_updated, name, symbol, address, latest_block, creator_address, chain FROM contracts WHERE address = $1 AND chain = $2 AND deleted = false
`

type GetContractByChainAddressParams struct {
	Address sql.NullString
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
		&i.LatestBlock,
		&i.CreatorAddress,
		&i.Chain,
	)
	return i, err
}

const getContractByID = `-- name: GetContractByID :one
select id, deleted, version, created_at, last_updated, name, symbol, address, latest_block, creator_address, chain FROM contracts WHERE id = $1 AND deleted = false
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
		&i.LatestBlock,
		&i.CreatorAddress,
		&i.Chain,
	)
	return i, err
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
SELECT id, deleted, version, created_at, last_updated, name, description, contract_address, collectors_note, media, token_uri, token_type, token_id, quantity, ownership_history, token_metadata, external_url, block_number, owner_user_id, owned_by_wallets, chain FROM tokens WHERE id = $1 AND deleted = false
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
		&i.ContractAddress,
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
	)
	return i, err
}

const getTokensByCollectionId = `-- name: GetTokensByCollectionId :many
SELECT t.id, t.deleted, t.version, t.created_at, t.last_updated, t.name, t.description, t.contract_address, t.collectors_note, t.media, t.token_uri, t.token_type, t.token_id, t.quantity, t.ownership_history, t.token_metadata, t.external_url, t.block_number, t.owner_user_id, t.owned_by_wallets, t.chain FROM collections c, unnest(c.nfts)
    WITH ORDINALITY AS x(nft_id, nft_ord)
    INNER JOIN tokens t ON t.id = x.nft_id
    WHERE c.id = $1 AND c.deleted = false AND t.deleted = false ORDER BY x.nft_ord
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
			&i.ContractAddress,
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
SELECT id, deleted, version, created_at, last_updated, name, description, contract_address, collectors_note, media, token_uri, token_type, token_id, quantity, ownership_history, token_metadata, external_url, block_number, owner_user_id, owned_by_wallets, chain FROM tokens WHERE owner_user_id = $1 AND deleted = false
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
			&i.ContractAddress,
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
SELECT id, deleted, version, created_at, last_updated, name, description, contract_address, collectors_note, media, token_uri, token_type, token_id, quantity, ownership_history, token_metadata, external_url, block_number, owner_user_id, owned_by_wallets, chain FROM tokens WHERE owned_by_wallets && $1 AND deleted = false
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
			&i.ContractAddress,
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
SELECT id, deleted, version, last_updated, created_at, username, username_idempotent, wallets, bio FROM users WHERE id = $1 AND deleted = false
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
	)
	return i, err
}

const getUserByUsername = `-- name: GetUserByUsername :one
SELECT id, deleted, version, last_updated, created_at, username, username_idempotent, wallets, bio FROM users WHERE username_idempotent = lower($1) AND deleted = false
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
	)
	return i, err
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
