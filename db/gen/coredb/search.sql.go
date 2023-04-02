// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: search.sql

package coredb

import (
	"context"
)

const searchContracts = `-- name: SearchContracts :many
with min_content_score as (
    select score from contract_relevance where id is null
),
poap_weight as (
    -- Using a CTE as a workaround because sqlc has trouble with this as an inline value
    -- in the ts_rank_cd statement below. We want non-POAP addresses to get crazy high weighting,
    -- but ts_rank weights have to be in the [0, 1] range, so we divide the POAP weight by 1000000000
    -- to offset the fact that we're going to multiply all addresses by 1000000000.
    select $5::float4 / 1000000000 as weight
)
select contracts.id, contracts.deleted, contracts.version, contracts.created_at, contracts.last_updated, contracts.name, contracts.symbol, contracts.address, contracts.creator_address, contracts.chain, contracts.profile_banner_url, contracts.profile_image_url, contracts.badge_url, contracts.description, contracts.parent_id from contracts left join contract_relevance on contract_relevance.id = contracts.id,
     to_tsquery('simple', websearch_to_tsquery('simple', $1)::text || ':*') simple_partial_query,
     websearch_to_tsquery('simple', $1) simple_full_query,
     websearch_to_tsquery('english', $1) english_full_query,
     min_content_score,
     poap_weight,
     greatest (
        ts_rank_cd(concat('{', $2::float4, ', 1, 1, 1}')::float4[], fts_name, simple_partial_query, 1),
        ts_rank_cd(concat('{', $3::float4, ', 1, 1, 1}')::float4[], fts_description_english, english_full_query, 1),
        ts_rank_cd(concat('{', poap_weight.weight::float4, ', 1, 1, 1}')::float4[], fts_address, simple_full_query, 1) * 1000000000
        ) as match_score,
     coalesce(contract_relevance.score, min_content_score.score) as content_score
where (
    simple_full_query @@ fts_address or
    simple_partial_query @@ fts_name or
    english_full_query @@ fts_description_english
    )
    and contracts.deleted = false
order by content_score * match_score desc, content_score desc, match_score desc
limit $4
`

type SearchContractsParams struct {
	Query             string
	NameWeight        float32
	DescriptionWeight float32
	Limit             int32
	PoapAddressWeight float32
}

func (q *Queries) SearchContracts(ctx context.Context, arg SearchContractsParams) ([]Contract, error) {
	rows, err := q.db.Query(ctx, searchContracts,
		arg.Query,
		arg.NameWeight,
		arg.DescriptionWeight,
		arg.Limit,
		arg.PoapAddressWeight,
	)
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
			&i.Description,
			&i.ParentID,
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

const searchGalleries = `-- name: SearchGalleries :many
with min_content_score as (
    select score from gallery_relevance where id is null
)
select galleries.id, galleries.deleted, galleries.last_updated, galleries.created_at, galleries.version, galleries.owner_user_id, galleries.collections, galleries.name, galleries.description, galleries.hidden, galleries.position from galleries left join gallery_relevance on gallery_relevance.id = galleries.id,
    to_tsquery('simple', websearch_to_tsquery('simple', $1)::text || ':*') simple_partial_query,
    websearch_to_tsquery('english', $1) english_full_query,
    min_content_score,
    greatest(
        ts_rank_cd(concat('{', $2::float4, ', 1, 1, 1}')::float4[], fts_name, simple_partial_query, 1),
        ts_rank_cd(concat('{', $3::float4, ', 1, 1, 1}')::float4[], fts_description_english, english_full_query, 1)
        ) as match_score,
    coalesce(gallery_relevance.score, min_content_score.score) as content_score
where (
    simple_partial_query @@ fts_name or
    english_full_query @@ fts_description_english
    )
    and deleted = false and hidden = false
order by content_score * match_score desc, content_score desc, match_score desc
limit $4
`

type SearchGalleriesParams struct {
	Query             string
	NameWeight        float32
	DescriptionWeight float32
	Limit             int32
}

func (q *Queries) SearchGalleries(ctx context.Context, arg SearchGalleriesParams) ([]Gallery, error) {
	rows, err := q.db.Query(ctx, searchGalleries,
		arg.Query,
		arg.NameWeight,
		arg.DescriptionWeight,
		arg.Limit,
	)
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
			&i.Name,
			&i.Description,
			&i.Hidden,
			&i.Position,
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

const searchUsers = `-- name: SearchUsers :many
with min_content_score as (
    select score from user_relevance where id is null
)
select u.id, u.deleted, u.version, u.last_updated, u.created_at, u.username, u.username_idempotent, u.wallets, u.bio, u.traits, u.universal, u.notification_settings, u.email_verified, u.email_unsubscriptions, u.featured_gallery, u.primary_wallet_id, u.user_experiences from users u left join user_relevance on u.id = user_relevance.id,
    -- Adding the search condition to the wallet join statement is a very helpful optimization, but we can't use
    -- "simple_full_query" at this point in the statement, so we're repeating the "websearch_to_tsquery..." part here
    unnest(u.wallets) as wallet_id left join wallets w on w.id = wallet_id and w.deleted = false and websearch_to_tsquery('simple', $1) @@ w.fts_address,
    to_tsquery('simple', websearch_to_tsquery('simple', $1)::text || ':*') simple_partial_query,
    websearch_to_tsquery('simple', $1) simple_full_query,
    websearch_to_tsquery('english', $1) english_full_query,
    min_content_score,
    greatest(
        ts_rank_cd(concat('{', $2::float4, ', 1, 1, 1}')::float4[], u.fts_username, simple_partial_query, 1),
        ts_rank_cd(concat('{', $3::float4, ', 1, 1, 1}')::float4[], u.fts_bio_english, english_full_query, 1),
        ts_rank_cd('{1, 1, 1, 1}', w.fts_address, simple_full_query) * 1000000000
        ) as match_score,
    coalesce(user_relevance.score, min_content_score.score) as content_score
where (
    simple_partial_query @@ u.fts_username or
    english_full_query @@ u.fts_bio_english or
    simple_full_query @@ w.fts_address
    )
    and u.universal = false and u.deleted = false
group by (u.id, content_score * match_score, content_score, match_score)
order by content_score * match_score desc, content_score desc, match_score desc, length(u.username_idempotent) asc
limit $4
`

type SearchUsersParams struct {
	Query          string
	UsernameWeight float32
	BioWeight      float32
	Limit          int32
}

func (q *Queries) SearchUsers(ctx context.Context, arg SearchUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, searchUsers,
		arg.Query,
		arg.UsernameWeight,
		arg.BioWeight,
		arg.Limit,
	)
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
			&i.Universal,
			&i.NotificationSettings,
			&i.EmailVerified,
			&i.EmailUnsubscriptions,
			&i.FeaturedGallery,
			&i.PrimaryWalletID,
			&i.UserExperiences,
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
