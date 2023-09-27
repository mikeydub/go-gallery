// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: recommend.sql

package coredb

import (
	"context"
	"time"

	"github.com/mikeydub/go-gallery/service/persist"
)

const getContractLabels = `-- name: GetContractLabels :many
select user_id, contract_id, displayed
from owned_contracts
where contract_id not in (
  select id from contracts where chain || ':' || address = any($1::varchar[])
) and displayed
`

type GetContractLabelsRow struct {
	UserID     persist.DBID `json:"user_id"`
	ContractID persist.DBID `json:"contract_id"`
	Displayed  bool         `json:"displayed"`
}

func (q *Queries) GetContractLabels(ctx context.Context, excludedContracts []string) ([]GetContractLabelsRow, error) {
	rows, err := q.db.Query(ctx, getContractLabels, excludedContracts)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetContractLabelsRow
	for rows.Next() {
		var i GetContractLabelsRow
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

const getExternalFollowGraphSource = `-- name: GetExternalFollowGraphSource :many
select
  external_social_connections.follower_id,
  external_social_connections.followee_id
from
  external_social_connections,
  users as followers,
  users as followees
where
  external_social_connections.follower_id = followers.id
  and followers.deleted is false
  and external_social_connections.followee_id = followees.id
  and followees.deleted is false
  and external_social_connections.deleted = false
`

type GetExternalFollowGraphSourceRow struct {
	FollowerID persist.DBID `json:"follower_id"`
	FolloweeID persist.DBID `json:"followee_id"`
}

func (q *Queries) GetExternalFollowGraphSource(ctx context.Context) ([]GetExternalFollowGraphSourceRow, error) {
	rows, err := q.db.Query(ctx, getExternalFollowGraphSource)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetExternalFollowGraphSourceRow
	for rows.Next() {
		var i GetExternalFollowGraphSourceRow
		if err := rows.Scan(&i.FollowerID, &i.FolloweeID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeedEntityScores = `-- name: GetFeedEntityScores :many
with refreshed as (
  select greatest((select last_updated from feed_entity_scores limit 1), $1::timestamptz) last_updated
)
select feed_entity_scores.id, feed_entity_scores.created_at, feed_entity_scores.actor_id, feed_entity_scores.action, feed_entity_scores.contract_ids, feed_entity_scores.interactions, feed_entity_scores.feed_entity_type, feed_entity_scores.last_updated, posts.id, posts.version, posts.token_ids, posts.contract_ids, posts.actor_id, posts.caption, posts.created_at, posts.last_updated, posts.deleted
from feed_entity_scores
join posts on feed_entity_scores.id = posts.id
where feed_entity_scores.created_at > $1::timestamptz
  and ($2::bool or feed_entity_scores.actor_id != $3)
  and feed_entity_scores.feed_entity_type == $4
  and not posts.deleted
union
select feed_entity_scores.id, feed_entity_scores.created_at, feed_entity_scores.actor_id, feed_entity_scores.action, feed_entity_scores.contract_ids, feed_entity_scores.interactions, feed_entity_scores.feed_entity_type, feed_entity_scores.last_updated, posts.id, posts.version, posts.token_ids, posts.contract_ids, posts.actor_id, posts.caption, posts.created_at, posts.last_updated, posts.deleted
from feed_entity_score_view feed_entity_scores
join posts on feed_entity_score_view.id = posts.id
where created_at > (select last_updated from refreshed limit 1)
and ($2::bool or feed_entity_score_view.actor_id != $3)
and feed_entity_score_view.feed_entity_type == $4
`

type GetFeedEntityScoresParams struct {
	WindowEnd      time.Time    `json:"window_end"`
	IncludeViewer  bool         `json:"include_viewer"`
	ViewerID       persist.DBID `json:"viewer_id"`
	PostEntityType int32        `json:"post_entity_type"`
}

type GetFeedEntityScoresRow struct {
	FeedEntityScore FeedEntityScore `json:"feedentityscore"`
	Post            Post            `json:"post"`
}

func (q *Queries) GetFeedEntityScores(ctx context.Context, arg GetFeedEntityScoresParams) ([]GetFeedEntityScoresRow, error) {
	rows, err := q.db.Query(ctx, getFeedEntityScores,
		arg.WindowEnd,
		arg.IncludeViewer,
		arg.ViewerID,
		arg.PostEntityType,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedEntityScoresRow
	for rows.Next() {
		var i GetFeedEntityScoresRow
		if err := rows.Scan(
			&i.FeedEntityScore.ID,
			&i.FeedEntityScore.CreatedAt,
			&i.FeedEntityScore.ActorID,
			&i.FeedEntityScore.Action,
			&i.FeedEntityScore.ContractIds,
			&i.FeedEntityScore.Interactions,
			&i.FeedEntityScore.FeedEntityType,
			&i.FeedEntityScore.LastUpdated,
			&i.Post.ID,
			&i.Post.Version,
			&i.Post.TokenIds,
			&i.Post.ContractIds,
			&i.Post.ActorID,
			&i.Post.Caption,
			&i.Post.CreatedAt,
			&i.Post.LastUpdated,
			&i.Post.Deleted,
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
  users as followees
where
  follows.follower = followers.id
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
select follower_id id from external_social_connections where not deleted group by 1
union
select followee_id id from external_social_connections where not deleted group by 1
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
