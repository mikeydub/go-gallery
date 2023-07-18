// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: recommend.sql

package coredb

import (
	"context"

	"github.com/mikeydub/go-gallery/service/persist"
)

const getDisplayedContracts = `-- name: GetDisplayedContracts :many
select user_id, contract_id, displayed
from owned_contracts
where contract_id not in (
	select id from contracts where chain || ':' || address = any($1::varchar[])
) and displayed
`

type GetDisplayedContractsRow struct {
	UserID     persist.DBID `json:"user_id"`
	ContractID persist.DBID `json:"contract_id"`
	Displayed  bool         `json:"displayed"`
}

func (q *Queries) GetDisplayedContracts(ctx context.Context, excludedContracts []string) ([]GetDisplayedContractsRow, error) {
	rows, err := q.db.Query(ctx, getDisplayedContracts, excludedContracts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetDisplayedContractsRow
	for rows.Next() {
		var i GetDisplayedContractsRow
		if err := rows.Scan(&i.UserID, &i.ContractID, &i.Displayed); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFollowEdgesByUserID = `-- name: GetFollowEdgesByUserID :many
select id, follower, followee, deleted, created_at, last_updated from follows f where f.follower = $1 and f.deleted = false
`

func (q *Queries) GetFollowEdgesByUserID(ctx context.Context, follower persist.DBID) ([]Follow, error) {
	rows, err := q.db.Query(ctx, getFollowEdgesByUserID, follower)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Follow
	for rows.Next() {
		var i Follow
		if err := rows.Scan(
			&i.ID,
			&i.Follower,
			&i.Followee,
			&i.Deleted,
			&i.CreatedAt,
			&i.LastUpdated,
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

const getFollowGraphSource = `-- name: GetFollowGraphSource :many
select
  follows.follower,
  follows.followee
from
  follows,
  users as followers,
  users as followees,
  (
    select owner_user_id
    from collections
    where cardinality(nfts) > 0 and deleted = false
    group by owner_user_id
  ) displaying
where
  follows.follower = followers.id
  and follows.followee = displaying.owner_user_id
  and followers.deleted is false
  and follows.followee = followees.id
  and followees.deleted is false
  and follows.deleted = false
`

type GetFollowGraphSourceRow struct {
	Follower persist.DBID `json:"follower"`
	Followee persist.DBID `json:"followee"`
}

func (q *Queries) GetFollowGraphSource(ctx context.Context) ([]GetFollowGraphSourceRow, error) {
	rows, err := q.db.Query(ctx, getFollowGraphSource)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFollowGraphSourceRow
	for rows.Next() {
		var i GetFollowGraphSourceRow
		if err := rows.Scan(&i.Follower, &i.Followee); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopRecommendedUserIDs = `-- name: GetTopRecommendedUserIDs :many
select recommended_user_id from top_recommended_users
`

func (q *Queries) GetTopRecommendedUserIDs(ctx context.Context) ([]persist.DBID, error) {
	rows, err := q.db.Query(ctx, getTopRecommendedUserIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []persist.DBID
	for rows.Next() {
		var recommended_user_id persist.DBID
		if err := rows.Scan(&recommended_user_id); err != nil {
			return nil, err
		}
		items = append(items, recommended_user_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserLabels = `-- name: GetUserLabels :many
select follower id from follows where not deleted group by 1
union
select followee id from follows where not deleted group by 1
union
select user_id id from owned_contracts where displayed group by 1
`

func (q *Queries) GetUserLabels(ctx context.Context) ([]persist.DBID, error) {
	rows, err := q.db.Query(ctx, getUserLabels)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []persist.DBID
	for rows.Next() {
		var id persist.DBID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatedRecommendationResults = `-- name: UpdatedRecommendationResults :exec
insert into recommendation_results
(
  id
  , user_id
  , recommended_user_id
  , recommended_count
) (
  select
    unnest($1::varchar[])
    , unnest($2::varchar[])
    , unnest($3::varchar[])
    , unnest($4::int[])
)
on conflict (user_id, recommended_user_id, version) where deleted = false
do update set
  recommended_count = recommendation_results.recommended_count + excluded.recommended_count,
  last_updated = now()
`

type UpdatedRecommendationResultsParams struct {
	ID                []string `json:"id"`
	UserID            []string `json:"user_id"`
	RecommendedUserID []string `json:"recommended_user_id"`
	RecommendedCount  []int32  `json:"recommended_count"`
}

func (q *Queries) UpdatedRecommendationResults(ctx context.Context, arg UpdatedRecommendationResultsParams) error {
	_, err := q.db.Exec(ctx, updatedRecommendationResults,
		arg.ID,
		arg.UserID,
		arg.RecommendedUserID,
		arg.RecommendedCount,
	)
	return err
}
