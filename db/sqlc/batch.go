// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.13.0
// source: batch.go
package sqlc

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v4"
	"github.com/mikeydub/go-gallery/service/persist"
)

const getCollectionByIdBatch = `-- name: GetCollectionByIdBatch :batchone
SELECT id, deleted, owner_user_id, nfts, version, last_updated, created_at, hidden, collectors_note, name, layout FROM collections WHERE id = $1 AND deleted = false
`

type GetCollectionByIdBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetCollectionByIdBatch(ctx context.Context, id []persist.DBID) *GetCollectionByIdBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range id {
		vals := []interface{}{
			a,
		}
		batch.Queue(getCollectionByIdBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetCollectionByIdBatchBatchResults{br, 0}
}

func (b *GetCollectionByIdBatchBatchResults) QueryRow(f func(int, Collection, error)) {
	for {
		row := b.br.QueryRow()
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
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, i, err)
		}
		b.ind++
	}
}

func (b *GetCollectionByIdBatchBatchResults) Close() error {
	return b.br.Close()
}

const getCollectionsByGalleryIdBatch = `-- name: GetCollectionsByGalleryIdBatch :batchmany
SELECT c.id, c.deleted, c.owner_user_id, c.nfts, c.version, c.last_updated, c.created_at, c.hidden, c.collectors_note, c.name, c.layout FROM galleries g, unnest(g.collections)
    WITH ORDINALITY AS x(coll_id, coll_ord)
    INNER JOIN collections c ON c.id = x.coll_id
    WHERE g.id = $1 AND g.deleted = false AND c.deleted = false ORDER BY x.coll_ord
`

type GetCollectionsByGalleryIdBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetCollectionsByGalleryIdBatch(ctx context.Context, id []persist.DBID) *GetCollectionsByGalleryIdBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range id {
		vals := []interface{}{
			a,
		}
		batch.Queue(getCollectionsByGalleryIdBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetCollectionsByGalleryIdBatchBatchResults{br, 0}
}

func (b *GetCollectionsByGalleryIdBatchBatchResults) Query(f func(int, []Collection, error)) {
	for {
		rows, err := b.br.Query()
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
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
				break
			}
			items = append(items, i)
		}

		if f != nil {
			f(b.ind, items, rows.Err())
		}
		b.ind++
	}
}

func (b *GetCollectionsByGalleryIdBatchBatchResults) Close() error {
	return b.br.Close()
}

const getContractByChainAddressBatch = `-- name: GetContractByChainAddressBatch :batchone
select id, deleted, version, created_at, last_updated, name, symbol, address, latest_block, creator_address, chain FROM contracts WHERE address = $1 AND chain = $2 AND deleted = false
`

type GetContractByChainAddressBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

type GetContractByChainAddressBatchParams struct {
	Address persist.Address
	Chain   sql.NullInt32
}

func (q *Queries) GetContractByChainAddressBatch(ctx context.Context, arg []GetContractByChainAddressBatchParams) *GetContractByChainAddressBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.Address,
			a.Chain,
		}
		batch.Queue(getContractByChainAddressBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetContractByChainAddressBatchBatchResults{br, 0}
}

func (b *GetContractByChainAddressBatchBatchResults) QueryRow(f func(int, Contract, error)) {
	for {
		row := b.br.QueryRow()
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
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, i, err)
		}
		b.ind++
	}
}

func (b *GetContractByChainAddressBatchBatchResults) Close() error {
	return b.br.Close()
}

const getContractByIDBatch = `-- name: GetContractByIDBatch :batchone
select id, deleted, version, created_at, last_updated, name, symbol, address, latest_block, creator_address, chain FROM contracts WHERE id = $1 AND deleted = false
`

type GetContractByIDBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetContractByIDBatch(ctx context.Context, id []persist.DBID) *GetContractByIDBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range id {
		vals := []interface{}{
			a,
		}
		batch.Queue(getContractByIDBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetContractByIDBatchBatchResults{br, 0}
}

func (b *GetContractByIDBatchBatchResults) QueryRow(f func(int, Contract, error)) {
	for {
		row := b.br.QueryRow()
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
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, i, err)
		}
		b.ind++
	}
}

func (b *GetContractByIDBatchBatchResults) Close() error {
	return b.br.Close()
}

const getFollowersByUserIdBatch = `-- name: GetFollowersByUserIdBatch :batchmany
SELECT u.id, u.deleted, u.version, u.last_updated, u.created_at, u.username, u.username_idempotent, u.wallets, u.bio FROM follows f
    INNER JOIN users u ON f.follower = u.id
    WHERE f.followee = $1 AND f.deleted = false
    ORDER BY f.last_updated DESC
`

type GetFollowersByUserIdBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetFollowersByUserIdBatch(ctx context.Context, followee []persist.DBID) *GetFollowersByUserIdBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range followee {
		vals := []interface{}{
			a,
		}
		batch.Queue(getFollowersByUserIdBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetFollowersByUserIdBatchBatchResults{br, 0}
}

func (b *GetFollowersByUserIdBatchBatchResults) Query(f func(int, []User, error)) {
	for {
		rows, err := b.br.Query()
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
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
			); err != nil {
				break
			}
			items = append(items, i)
		}

		if f != nil {
			f(b.ind, items, rows.Err())
		}
		b.ind++
	}
}

func (b *GetFollowersByUserIdBatchBatchResults) Close() error {
	return b.br.Close()
}

const getFollowingByUserIdBatch = `-- name: GetFollowingByUserIdBatch :batchmany
SELECT u.id, u.deleted, u.version, u.last_updated, u.created_at, u.username, u.username_idempotent, u.wallets, u.bio FROM follows f
    INNER JOIN users u ON f.followee = u.id
    WHERE f.follower = $1 AND f.deleted = false
    ORDER BY f.last_updated DESC
`

type GetFollowingByUserIdBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetFollowingByUserIdBatch(ctx context.Context, follower []persist.DBID) *GetFollowingByUserIdBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range follower {
		vals := []interface{}{
			a,
		}
		batch.Queue(getFollowingByUserIdBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetFollowingByUserIdBatchBatchResults{br, 0}
}

func (b *GetFollowingByUserIdBatchBatchResults) Query(f func(int, []User, error)) {
	for {
		rows, err := b.br.Query()
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
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
			); err != nil {
				break
			}
			items = append(items, i)
		}

		if f != nil {
			f(b.ind, items, rows.Err())
		}
		b.ind++
	}
}

func (b *GetFollowingByUserIdBatchBatchResults) Close() error {
	return b.br.Close()
}

const getGalleriesByUserIdBatch = `-- name: GetGalleriesByUserIdBatch :batchmany
SELECT id, deleted, last_updated, created_at, version, owner_user_id, collections FROM galleries WHERE owner_user_id = $1 AND deleted = false
`

type GetGalleriesByUserIdBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetGalleriesByUserIdBatch(ctx context.Context, ownerUserID []persist.DBID) *GetGalleriesByUserIdBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range ownerUserID {
		vals := []interface{}{
			a,
		}
		batch.Queue(getGalleriesByUserIdBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetGalleriesByUserIdBatchBatchResults{br, 0}
}

func (b *GetGalleriesByUserIdBatchBatchResults) Query(f func(int, []Gallery, error)) {
	for {
		rows, err := b.br.Query()
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
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
				break
			}
			items = append(items, i)
		}

		if f != nil {
			f(b.ind, items, rows.Err())
		}
		b.ind++
	}
}

func (b *GetGalleriesByUserIdBatchBatchResults) Close() error {
	return b.br.Close()
}

const getGalleryByCollectionIdBatch = `-- name: GetGalleryByCollectionIdBatch :batchone
SELECT g.id, g.deleted, g.last_updated, g.created_at, g.version, g.owner_user_id, g.collections FROM galleries g, collections c WHERE c.id = $1 AND c.deleted = false AND $1 = ANY(g.collections) AND g.deleted = false
`

type GetGalleryByCollectionIdBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetGalleryByCollectionIdBatch(ctx context.Context, id []persist.DBID) *GetGalleryByCollectionIdBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range id {
		vals := []interface{}{
			a,
		}
		batch.Queue(getGalleryByCollectionIdBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetGalleryByCollectionIdBatchBatchResults{br, 0}
}

func (b *GetGalleryByCollectionIdBatchBatchResults) QueryRow(f func(int, Gallery, error)) {
	for {
		row := b.br.QueryRow()
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
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, i, err)
		}
		b.ind++
	}
}

func (b *GetGalleryByCollectionIdBatchBatchResults) Close() error {
	return b.br.Close()
}

const getGalleryByIdBatch = `-- name: GetGalleryByIdBatch :batchone
SELECT id, deleted, last_updated, created_at, version, owner_user_id, collections FROM galleries WHERE id = $1 AND deleted = false
`

type GetGalleryByIdBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetGalleryByIdBatch(ctx context.Context, id []persist.DBID) *GetGalleryByIdBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range id {
		vals := []interface{}{
			a,
		}
		batch.Queue(getGalleryByIdBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetGalleryByIdBatchBatchResults{br, 0}
}

func (b *GetGalleryByIdBatchBatchResults) QueryRow(f func(int, Gallery, error)) {
	for {
		row := b.br.QueryRow()
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
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, i, err)
		}
		b.ind++
	}
}

func (b *GetGalleryByIdBatchBatchResults) Close() error {
	return b.br.Close()
}

const getMembershipByMembershipIdBatch = `-- name: GetMembershipByMembershipIdBatch :batchone
SELECT id, deleted, version, created_at, last_updated, token_id, name, asset_url, owners FROM membership WHERE id = $1 AND deleted = false
`

type GetMembershipByMembershipIdBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetMembershipByMembershipIdBatch(ctx context.Context, id []persist.DBID) *GetMembershipByMembershipIdBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range id {
		vals := []interface{}{
			a,
		}
		batch.Queue(getMembershipByMembershipIdBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetMembershipByMembershipIdBatchBatchResults{br, 0}
}

func (b *GetMembershipByMembershipIdBatchBatchResults) QueryRow(f func(int, Membership, error)) {
	for {
		row := b.br.QueryRow()
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
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, i, err)
		}
		b.ind++
	}
}

func (b *GetMembershipByMembershipIdBatchBatchResults) Close() error {
	return b.br.Close()
}

const getTokenByIdBatch = `-- name: GetTokenByIdBatch :batchone
SELECT id, deleted, version, created_at, last_updated, name, description, collectors_note, media, token_uri, token_type, token_id, quantity, ownership_history, token_metadata, external_url, block_number, owner_user_id, owned_by_wallets, chain, contract FROM tokens WHERE id = $1 AND deleted = false
`

type GetTokenByIdBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetTokenByIdBatch(ctx context.Context, id []persist.DBID) *GetTokenByIdBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range id {
		vals := []interface{}{
			a,
		}
		batch.Queue(getTokenByIdBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetTokenByIdBatchBatchResults{br, 0}
}

func (b *GetTokenByIdBatchBatchResults) QueryRow(f func(int, Token, error)) {
	for {
		row := b.br.QueryRow()
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
		)
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, i, err)
		}
		b.ind++
	}
}

func (b *GetTokenByIdBatchBatchResults) Close() error {
	return b.br.Close()
}

const getTokensByCollectionIdBatch = `-- name: GetTokensByCollectionIdBatch :batchmany
SELECT t.id, t.deleted, t.version, t.created_at, t.last_updated, t.name, t.description, t.collectors_note, t.media, t.token_uri, t.token_type, t.token_id, t.quantity, t.ownership_history, t.token_metadata, t.external_url, t.block_number, t.owner_user_id, t.owned_by_wallets, t.chain, t.contract FROM collections c, unnest(c.nfts)
    WITH ORDINALITY AS x(nft_id, nft_ord)
    INNER JOIN tokens t ON t.id = x.nft_id
    WHERE c.id = $1 AND c.deleted = false AND t.deleted = false ORDER BY x.nft_ord
`

type GetTokensByCollectionIdBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetTokensByCollectionIdBatch(ctx context.Context, id []persist.DBID) *GetTokensByCollectionIdBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range id {
		vals := []interface{}{
			a,
		}
		batch.Queue(getTokensByCollectionIdBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetTokensByCollectionIdBatchBatchResults{br, 0}
}

func (b *GetTokensByCollectionIdBatchBatchResults) Query(f func(int, []Token, error)) {
	for {
		rows, err := b.br.Query()
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
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
			); err != nil {
				break
			}
			items = append(items, i)
		}

		if f != nil {
			f(b.ind, items, rows.Err())
		}
		b.ind++
	}
}

func (b *GetTokensByCollectionIdBatchBatchResults) Close() error {
	return b.br.Close()
}

const getTokensByUserIdBatch = `-- name: GetTokensByUserIdBatch :batchmany
SELECT tokens.id, tokens.deleted, tokens.version, tokens.created_at, tokens.last_updated, tokens.name, tokens.description, tokens.collectors_note, tokens.media, tokens.token_uri, tokens.token_type, tokens.token_id, tokens.quantity, tokens.ownership_history, tokens.token_metadata, tokens.external_url, tokens.block_number, tokens.owner_user_id, tokens.owned_by_wallets, tokens.chain, tokens.contract FROM tokens, users
    WHERE tokens.owner_user_id = $1 AND users.id = $1
      AND tokens.owned_by_wallets && users.wallets
      AND tokens.deleted = false AND users.deleted = false
`

type GetTokensByUserIdBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetTokensByUserIdBatch(ctx context.Context, ownerUserID []persist.DBID) *GetTokensByUserIdBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range ownerUserID {
		vals := []interface{}{
			a,
		}
		batch.Queue(getTokensByUserIdBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetTokensByUserIdBatchBatchResults{br, 0}
}

func (b *GetTokensByUserIdBatchBatchResults) Query(f func(int, []Token, error)) {
	for {
		rows, err := b.br.Query()
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
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
			); err != nil {
				break
			}
			items = append(items, i)
		}

		if f != nil {
			f(b.ind, items, rows.Err())
		}
		b.ind++
	}
}

func (b *GetTokensByUserIdBatchBatchResults) Close() error {
	return b.br.Close()
}

const getTokensByWalletIdsBatch = `-- name: GetTokensByWalletIdsBatch :batchmany
SELECT id, deleted, version, created_at, last_updated, name, description, collectors_note, media, token_uri, token_type, token_id, quantity, ownership_history, token_metadata, external_url, block_number, owner_user_id, owned_by_wallets, chain, contract FROM tokens WHERE owned_by_wallets && $1 AND deleted = false
`

type GetTokensByWalletIdsBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetTokensByWalletIdsBatch(ctx context.Context, ownedByWallets []persist.DBIDList) *GetTokensByWalletIdsBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range ownedByWallets {
		vals := []interface{}{
			a,
		}
		batch.Queue(getTokensByWalletIdsBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetTokensByWalletIdsBatchBatchResults{br, 0}
}

func (b *GetTokensByWalletIdsBatchBatchResults) Query(f func(int, []Token, error)) {
	for {
		rows, err := b.br.Query()
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
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
			); err != nil {
				break
			}
			items = append(items, i)
		}

		if f != nil {
			f(b.ind, items, rows.Err())
		}
		b.ind++
	}
}

func (b *GetTokensByWalletIdsBatchBatchResults) Close() error {
	return b.br.Close()
}

const getUserByIdBatch = `-- name: GetUserByIdBatch :batchone
SELECT id, deleted, version, last_updated, created_at, username, username_idempotent, wallets, bio FROM users WHERE id = $1 AND deleted = false
`

type GetUserByIdBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetUserByIdBatch(ctx context.Context, id []persist.DBID) *GetUserByIdBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range id {
		vals := []interface{}{
			a,
		}
		batch.Queue(getUserByIdBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetUserByIdBatchBatchResults{br, 0}
}

func (b *GetUserByIdBatchBatchResults) QueryRow(f func(int, User, error)) {
	for {
		row := b.br.QueryRow()
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
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, i, err)
		}
		b.ind++
	}
}

func (b *GetUserByIdBatchBatchResults) Close() error {
	return b.br.Close()
}

const getUserByUsernameBatch = `-- name: GetUserByUsernameBatch :batchone
SELECT id, deleted, version, last_updated, created_at, username, username_idempotent, wallets, bio FROM users WHERE username_idempotent = lower($1) AND deleted = false
`

type GetUserByUsernameBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetUserByUsernameBatch(ctx context.Context, lower []string) *GetUserByUsernameBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range lower {
		vals := []interface{}{
			a,
		}
		batch.Queue(getUserByUsernameBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetUserByUsernameBatchBatchResults{br, 0}
}

func (b *GetUserByUsernameBatchBatchResults) QueryRow(f func(int, User, error)) {
	for {
		row := b.br.QueryRow()
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
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, i, err)
		}
		b.ind++
	}
}

func (b *GetUserByUsernameBatchBatchResults) Close() error {
	return b.br.Close()
}

const getWalletByChainAddressBatch = `-- name: GetWalletByChainAddressBatch :batchone
SELECT wallets.id, wallets.created_at, wallets.last_updated, wallets.deleted, wallets.version, wallets.address, wallets.wallet_type, wallets.chain FROM wallets WHERE address = $1 AND chain = $2 AND deleted = false
`

type GetWalletByChainAddressBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

type GetWalletByChainAddressBatchParams struct {
	Address persist.Address
	Chain   sql.NullInt32
}

func (q *Queries) GetWalletByChainAddressBatch(ctx context.Context, arg []GetWalletByChainAddressBatchParams) *GetWalletByChainAddressBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range arg {
		vals := []interface{}{
			a.Address,
			a.Chain,
		}
		batch.Queue(getWalletByChainAddressBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetWalletByChainAddressBatchBatchResults{br, 0}
}

func (b *GetWalletByChainAddressBatchBatchResults) QueryRow(f func(int, Wallet, error)) {
	for {
		row := b.br.QueryRow()
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
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, i, err)
		}
		b.ind++
	}
}

func (b *GetWalletByChainAddressBatchBatchResults) Close() error {
	return b.br.Close()
}

const getWalletByIDBatch = `-- name: GetWalletByIDBatch :batchone
SELECT id, created_at, last_updated, deleted, version, address, wallet_type, chain FROM wallets WHERE id = $1 AND deleted = false
`

type GetWalletByIDBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetWalletByIDBatch(ctx context.Context, id []persist.DBID) *GetWalletByIDBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range id {
		vals := []interface{}{
			a,
		}
		batch.Queue(getWalletByIDBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetWalletByIDBatchBatchResults{br, 0}
}

func (b *GetWalletByIDBatchBatchResults) QueryRow(f func(int, Wallet, error)) {
	for {
		row := b.br.QueryRow()
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
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
		}
		if f != nil {
			f(b.ind, i, err)
		}
		b.ind++
	}
}

func (b *GetWalletByIDBatchBatchResults) Close() error {
	return b.br.Close()
}

const getWalletsByUserIDBatch = `-- name: GetWalletsByUserIDBatch :batchmany
SELECT w.id, w.created_at, w.last_updated, w.deleted, w.version, w.address, w.wallet_type, w.chain FROM users u, unnest(u.wallets) WITH ORDINALITY AS a(wallet_id, wallet_ord)INNER JOIN wallets w on w.id = a.wallet_id WHERE u.id = $1 AND u.deleted = false AND w.deleted = false ORDER BY a.wallet_ord
`

type GetWalletsByUserIDBatchBatchResults struct {
	br  pgx.BatchResults
	ind int
}

func (q *Queries) GetWalletsByUserIDBatch(ctx context.Context, id []persist.DBID) *GetWalletsByUserIDBatchBatchResults {
	batch := &pgx.Batch{}
	for _, a := range id {
		vals := []interface{}{
			a,
		}
		batch.Queue(getWalletsByUserIDBatch, vals...)
	}
	br := q.db.SendBatch(ctx, batch)
	return &GetWalletsByUserIDBatchBatchResults{br, 0}
}

func (b *GetWalletsByUserIDBatchBatchResults) Query(f func(int, []Wallet, error)) {
	for {
		rows, err := b.br.Query()
		if err != nil && (err.Error() == "no result" || err.Error() == "batch already closed") {
			break
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
				break
			}
			items = append(items, i)
		}

		if f != nil {
			f(b.ind, items, rows.Err())
		}
		b.ind++
	}
}

func (b *GetWalletsByUserIDBatchBatchResults) Close() error {
	return b.br.Close()
}
