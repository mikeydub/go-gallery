-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 AND deleted = false;

-- name: GetUserWithPIIByID :one
select * from pii.user_view where id = @user_id and deleted = false;

-- name: GetUserByIdBatch :batchone
SELECT * FROM users WHERE id = $1 AND deleted = false;

-- name: GetUsersByIDs :many
SELECT * FROM users WHERE id = ANY(@user_ids) AND deleted = false
    AND (created_at, id) < (@cur_before_time, @cur_before_id)
    AND (created_at, id) > (@cur_after_time, @cur_after_id)
    ORDER BY CASE WHEN @paging_forward::bool THEN (created_at, id) END ASC,
             CASE WHEN NOT @paging_forward::bool THEN (created_at, id) END DESC
    LIMIT $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username_idempotent = lower(sqlc.arg('username')) AND deleted = false;

-- name: GetUserByUsernameBatch :batchone
SELECT * FROM users WHERE username_idempotent = lower($1) AND deleted = false;

-- name: GetUserByVerifiedEmailAddress :one
select u.* from users u join pii.for_users p on u.id = p.user_id
where p.pii_email_address = lower($1)
  and u.email_verified != 0
  and p.deleted = false
  and u.deleted = false;

-- name: GetUserByAddressBatch :batchone
select users.*
from users, wallets
where wallets.address = sqlc.arg('address')
	and wallets.chain = sqlc.arg('chain')::int
	and array[wallets.id] <@ users.wallets
	and wallets.deleted = false
	and users.deleted = false;

-- name: GetUsersWithTrait :many
SELECT * FROM users WHERE (traits->$1::string) IS NOT NULL AND deleted = false;

-- name: GetUsersWithTraitBatch :batchmany
SELECT * FROM users WHERE (traits->$1::string) IS NOT NULL AND deleted = false;

-- name: GetGalleryById :one
SELECT * FROM galleries WHERE id = $1 AND deleted = false;

-- name: GetGalleryByIdBatch :batchone
SELECT * FROM galleries WHERE id = $1 AND deleted = false;

-- name: GetGalleryByCollectionId :one
SELECT g.* FROM galleries g, collections c WHERE c.id = $1 AND c.deleted = false AND $1 = ANY(g.collections) AND g.deleted = false;

-- name: GetGalleryByCollectionIdBatch :batchone
SELECT g.* FROM galleries g, collections c WHERE c.id = $1 AND c.deleted = false AND $1 = ANY(g.collections) AND g.deleted = false;

-- name: GetGalleriesByUserId :many
SELECT * FROM galleries WHERE owner_user_id = $1 AND deleted = false order by position;

-- name: GetGalleriesByUserIdBatch :batchmany
SELECT * FROM galleries WHERE owner_user_id = $1 AND deleted = false order by position;

-- name: GetCollectionById :one
SELECT * FROM collections WHERE id = $1 AND deleted = false;

-- name: GetCollectionByIdBatch :batchone
SELECT * FROM collections WHERE id = $1 AND deleted = false;

-- name: GetCollectionsByGalleryId :many
SELECT c.* FROM galleries g, unnest(g.collections)
    WITH ORDINALITY AS x(coll_id, coll_ord)
    INNER JOIN collections c ON c.id = x.coll_id
    WHERE g.id = $1 AND g.deleted = false AND c.deleted = false ORDER BY x.coll_ord;

-- name: GetCollectionsByGalleryIdBatch :batchmany
SELECT c.* FROM galleries g, unnest(g.collections)
    WITH ORDINALITY AS x(coll_id, coll_ord)
    INNER JOIN collections c ON c.id = x.coll_id
    WHERE g.id = $1 AND g.deleted = false AND c.deleted = false ORDER BY x.coll_ord;

-- name: GetTokenById :one
SELECT * FROM tokens WHERE id = $1 AND deleted = false;

-- name: GetTokenByIdBatch :batchone
SELECT * FROM tokens WHERE id = $1 AND deleted = false;

-- name: GetTokensByCollectionId :many
SELECT t.* FROM users u, collections c, unnest(c.nfts)
    WITH ORDINALITY AS x(nft_id, nft_ord)
    INNER JOIN tokens t ON t.id = x.nft_id
    WHERE u.id = t.owner_user_id AND t.owned_by_wallets && u.wallets
    AND c.id = sqlc.arg('collection_id') AND u.deleted = false AND c.deleted = false AND t.deleted = false ORDER BY x.nft_ord LIMIT sqlc.narg('limit');

-- name: GetTokensByCollectionIdBatch :batchmany
SELECT t.* FROM users u, collections c, unnest(c.nfts)
    WITH ORDINALITY AS x(nft_id, nft_ord)
    INNER JOIN tokens t ON t.id = x.nft_id
    WHERE u.id = t.owner_user_id AND t.owned_by_wallets && u.wallets
    AND c.id = sqlc.arg('collection_id') AND u.deleted = false AND c.deleted = false AND t.deleted = false ORDER BY x.nft_ord LIMIT sqlc.narg('limit');

-- name: GetNewTokensByFeedEventIdBatch :batchmany
WITH new_tokens AS (
    SELECT added.id, row_number() OVER () added_order
    FROM (SELECT jsonb_array_elements_text(data -> 'collection_new_token_ids') id FROM feed_events f WHERE f.id = $1 AND f.deleted = false) added
)
SELECT t.* FROM new_tokens a JOIN tokens t ON a.id = t.id AND t.deleted = false ORDER BY a.added_order;

-- name: GetMembershipByMembershipId :one
SELECT * FROM membership WHERE id = $1 AND deleted = false;

-- name: GetMembershipByMembershipIdBatch :batchone
SELECT * FROM membership WHERE id = $1 AND deleted = false;

-- name: GetWalletByID :one
SELECT * FROM wallets WHERE id = $1 AND deleted = false;

-- name: GetWalletByIDBatch :batchone
SELECT * FROM wallets WHERE id = $1 AND deleted = false;

-- name: GetWalletByChainAddress :one
SELECT wallets.* FROM wallets WHERE address = $1 AND chain = $2 AND deleted = false;

-- name: GetWalletByChainAddressBatch :batchone
SELECT wallets.* FROM wallets WHERE address = $1 AND chain = $2 AND deleted = false;

-- name: GetWalletsByUserID :many
SELECT w.* FROM users u, unnest(u.wallets) WITH ORDINALITY AS a(wallet_id, wallet_ord)INNER JOIN wallets w on w.id = a.wallet_id WHERE u.id = $1 AND u.deleted = false AND w.deleted = false ORDER BY a.wallet_ord;

-- name: GetWalletsByUserIDBatch :batchmany
SELECT w.* FROM users u, unnest(u.wallets) WITH ORDINALITY AS a(wallet_id, wallet_ord)INNER JOIN wallets w on w.id = a.wallet_id WHERE u.id = $1 AND u.deleted = false AND w.deleted = false ORDER BY a.wallet_ord;

-- name: GetContractByID :one
select * FROM contracts WHERE id = $1 AND deleted = false;

-- name: GetContractsByIDs :many
SELECT * from contracts WHERE id = ANY(@contract_ids) AND deleted = false;

-- name: GetContractByChainAddress :one
select * FROM contracts WHERE address = $1 AND chain = $2 AND deleted = false;

-- name: GetContractByChainAddressBatch :batchone
select * FROM contracts WHERE address = $1 AND chain = $2 AND deleted = false;

-- name: GetContractsByUserID :many
SELECT DISTINCT ON (contracts.id) contracts.* FROM contracts, tokens
    WHERE tokens.owner_user_id = $1 AND tokens.contract = contracts.id
    AND tokens.deleted = false AND contracts.deleted = false;

-- name: GetContractsByUserIDBatch :batchmany
SELECT DISTINCT ON (contracts.id) contracts.* FROM contracts, tokens
    WHERE tokens.owner_user_id = $1 AND tokens.contract = contracts.id
    AND tokens.deleted = false AND contracts.deleted = false;

-- name: GetContractsDisplayedByUserIDBatch :batchmany
with last_refreshed as (
  select last_updated from owned_contracts limit 1
),
displayed as (
  select contract_id
  from owned_contracts
  where owned_contracts.user_id = $1 and displayed = true
  union
  select contracts.id
  from last_refreshed, galleries, contracts, tokens
  join collections on tokens.id = any(collections.nfts) and collections.deleted = false
  where tokens.owner_user_id = $1
    and tokens.contract = contracts.id
    and collections.owner_user_id = tokens.owner_user_id
    and galleries.owner_user_id = tokens.owner_user_id
    and tokens.deleted = false
    and galleries.deleted = false
    and contracts.deleted = false
    and galleries.last_updated > last_refreshed.last_updated
    and collections.last_updated > last_refreshed.last_updated
)
select contracts.* from contracts, displayed
where contracts.id = displayed.contract_id and contracts.deleted = false;

-- name: GetFollowersByUserIdBatch :batchmany
SELECT u.* FROM follows f
    INNER JOIN users u ON f.follower = u.id
    WHERE f.followee = $1 AND f.deleted = false
    ORDER BY f.last_updated DESC;

-- name: GetFollowingByUserIdBatch :batchmany
SELECT u.* FROM follows f
    INNER JOIN users u ON f.followee = u.id
    WHERE f.follower = $1 AND f.deleted = false
    ORDER BY f.last_updated DESC;

-- name: GetTokensByWalletIds :many
SELECT * FROM tokens WHERE owned_by_wallets && $1 AND deleted = false
    ORDER BY tokens.created_at DESC, tokens.name DESC, tokens.id DESC;

-- name: GetTokensByWalletIdsBatch :batchmany
SELECT * FROM tokens WHERE owned_by_wallets && $1 AND deleted = false
    ORDER BY tokens.created_at DESC, tokens.name DESC, tokens.id DESC;

-- name: GetTokensByContractId :many
SELECT * FROM tokens WHERE contract = $1 AND deleted = false
    ORDER BY tokens.created_at DESC, tokens.name DESC, tokens.id DESC;

-- name: GetTokensByContractIdBatch :batchmany
SELECT * FROM tokens WHERE contract = $1 AND deleted = false
    ORDER BY tokens.created_at DESC, tokens.name DESC, tokens.id DESC;

-- name: GetTokensByContractIdPaginate :many
SELECT t.* FROM tokens t
    JOIN users u ON u.id = t.owner_user_id
    WHERE t.contract = $1 AND t.deleted = false
    AND (NOT @gallery_users_only::bool OR u.universal = false)
    AND (u.universal,t.created_at,t.id) < (@cur_before_universal, @cur_before_time::timestamptz, @cur_before_id)
    AND (u.universal,t.created_at,t.id) > (@cur_after_universal, @cur_after_time::timestamptz, @cur_after_id)
    ORDER BY CASE WHEN @paging_forward::bool THEN (u.universal,t.created_at,t.id) END ASC,
             CASE WHEN NOT @paging_forward::bool THEN (u.universal,t.created_at,t.id) END DESC
    LIMIT $2;

-- name: GetTokensByContractIdBatchPaginate :batchmany
SELECT t.* FROM tokens t
    JOIN users u ON u.id = t.owner_user_id
    WHERE t.contract = sqlc.arg('contract') AND t.deleted = false
    AND (NOT @gallery_users_only::bool OR u.universal = false)
    AND (u.universal,t.created_at,t.id) < (@cur_before_universal, @cur_before_time::timestamptz, @cur_before_id)
    AND (u.universal,t.created_at,t.id) > (@cur_after_universal, @cur_after_time::timestamptz, @cur_after_id)
    ORDER BY CASE WHEN @paging_forward::bool THEN (u.universal,t.created_at,t.id) END ASC,
             CASE WHEN NOT @paging_forward::bool THEN (u.universal,t.created_at,t.id) END DESC
    LIMIT sqlc.arg('limit');

-- name: CountTokensByContractId :one
SELECT count(*) FROM tokens JOIN users ON users.id = tokens.owner_user_id WHERE contract = $1 AND (NOT @gallery_users_only::bool OR users.universal = false) AND tokens.deleted = false;

-- name: GetOwnersByContractIdBatchPaginate :batchmany
-- Note: sqlc has trouble recognizing that the output of the "select distinct" subquery below will
--       return complete rows from the users table. As a workaround, aliasing the subquery to
--       "users" seems to fix the issue (along with aliasing the users table inside the subquery
--       to "u" to avoid confusion -- otherwise, sqlc creates a custom row type that includes
--       all users.* fields twice).
select users.* from (
    select distinct on (u.id) u.* from users u, tokens t
        where t.contract = sqlc.arg('contract') and t.owner_user_id = u.id
        and (not @gallery_users_only::bool or u.universal = false)
        and t.deleted = false and u.deleted = false
    ) as users
    where (users.universal,users.created_at,users.id) < (@cur_before_universal, @cur_before_time::timestamptz, @cur_before_id)
    and (users.universal,users.created_at,users.id) > (@cur_after_universal, @cur_after_time::timestamptz, @cur_after_id)
    order by case when @paging_forward::bool then (users.universal,users.created_at,users.id) end asc,
         case when not @paging_forward::bool then (users.universal,users.created_at,users.id) end desc limit sqlc.narg('limit');


-- name: CountOwnersByContractId :one
SELECT count(DISTINCT users.id) FROM users, tokens
    WHERE tokens.contract = $1 AND tokens.owner_user_id = users.id
    AND (NOT @gallery_users_only::bool OR users.universal = false)
    AND tokens.deleted = false AND users.deleted = false;

-- name: GetTokenOwnerByID :one
SELECT u.* FROM tokens t
    JOIN users u ON u.id = t.owner_user_id
    WHERE t.id = $1 AND t.deleted = false AND u.deleted = false;

-- name: GetTokenOwnerByIDBatch :batchone
SELECT u.* FROM tokens t
    JOIN users u ON u.id = t.owner_user_id
    WHERE t.id = $1 AND t.deleted = false AND u.deleted = false;

-- name: GetPreviewURLsByContractIdAndUserId :many
SELECT (MEDIA->>'thumbnail_url')::varchar as thumbnail_url FROM tokens WHERE CONTRACT = $1 AND DELETED = false AND OWNER_USER_ID = $2 AND LENGTH(MEDIA->>'thumbnail_url'::varchar) > 0 ORDER BY ID LIMIT 3;

-- name: GetTokensByUserId :many
SELECT tokens.* FROM tokens, users
    WHERE tokens.owner_user_id = $1 AND users.id = $1
      AND tokens.owned_by_wallets && users.wallets
      AND tokens.deleted = false AND users.deleted = false
    ORDER BY tokens.created_at DESC, tokens.name DESC, tokens.id DESC;

-- name: GetTokensByUserIdBatch :batchmany
SELECT tokens.* FROM tokens, users
    WHERE tokens.owner_user_id = $1 AND users.id = $1
      AND tokens.owned_by_wallets && users.wallets
      AND tokens.deleted = false AND users.deleted = false
    ORDER BY tokens.created_at DESC, tokens.name DESC, tokens.id DESC;

-- name: GetTokensByUserIdAndContractID :many
SELECT tokens.* FROM tokens, users
    WHERE tokens.owner_user_id = $1 AND users.id = $1
      AND tokens.owned_by_wallets && users.wallets
      AND tokens.contract = $2
      AND tokens.deleted = false AND users.deleted = false
    ORDER BY tokens.created_at DESC, tokens.name DESC, tokens.id DESC;

-- name: GetTokensByUserIdAndContractIDBatch :batchmany
SELECT tokens.* FROM tokens, users
    WHERE tokens.owner_user_id = $1 AND users.id = $1
      AND tokens.owned_by_wallets && users.wallets
      AND tokens.contract = $2
      AND tokens.deleted = false AND users.deleted = false
    ORDER BY tokens.created_at DESC, tokens.name DESC, tokens.id DESC;

-- name: GetTokensByUserIdAndChainBatch :batchmany
SELECT tokens.* FROM tokens, users
WHERE tokens.owner_user_id = $1 AND users.id = $1
  AND tokens.owned_by_wallets && users.wallets
  AND tokens.deleted = false AND users.deleted = false
  AND tokens.chain = $2
ORDER BY tokens.created_at DESC, tokens.name DESC, tokens.id DESC;

-- name: CreateUserEvent :one
INSERT INTO events (id, actor_id, action, resource_type_id, user_id, subject_id, data, group_id, caption) VALUES ($1, $2, $3, $4, $5, $5, $6, $7, $8) RETURNING *;

-- name: CreateTokenEvent :one
INSERT INTO events (id, actor_id, action, resource_type_id, token_id, subject_id, data, group_id, caption, gallery_id, collection_id) VALUES ($1, $2, $3, $4, $5, $5, $6, $7, $8, $9, $10) RETURNING *;

-- name: CreateCollectionEvent :one
INSERT INTO events (id, actor_id, action, resource_type_id, collection_id, subject_id, data, caption, group_id, gallery_id) VALUES ($1, $2, $3, $4, $5, $5, $6, $7, $8, $9) RETURNING *;

-- name: CreateGalleryEvent :one
INSERT INTO events (id, actor_id, action, resource_type_id, gallery_id, subject_id, data, external_id, group_id, caption) VALUES ($1, $2, $3, $4, $5, $5, $6, $7, $8, $9) RETURNING *;

-- name: CreateAdmireEvent :one
INSERT INTO events (id, actor_id, action, resource_type_id, admire_id, feed_event_id, subject_id, data, group_id, caption) VALUES ($1, $2, $3, $4, $5, $6, $5, $7, $8, $9) RETURNING *;

-- name: CreateCommentEvent :one
INSERT INTO events (id, actor_id, action, resource_type_id, comment_id, feed_event_id, subject_id, data, group_id, caption) VALUES ($1, $2, $3, $4, $5, $6, $5, $7, $8, $9) RETURNING *;

-- name: GetEvent :one
SELECT * FROM events WHERE id = $1 AND deleted = false;

-- name: GetEventsInWindow :many
with recursive activity as (
    select * from events where events.id = $1 and deleted = false
    union
    select e.* from events e, activity a
    where e.actor_id = a.actor_id
        and e.action = any(@actions)
        and e.created_at < a.created_at
        and e.created_at >= a.created_at - make_interval(secs => $2)
        and e.deleted = false
        and e.caption is null
        and (not @include_subject::bool or e.subject_id = a.subject_id)
)
select * from events where id = any(select id from activity) order by (created_at, id) asc;

-- name: GetGalleryEventsInWindow :many
with recursive activity as (
    select * from events where events.id = $1 and deleted = false
    union
    select e.* from events e, activity a
    where e.actor_id = a.actor_id
        and e.action = any(@actions)
        and e.gallery_id = @gallery_id
        and e.created_at < a.created_at
        and e.created_at >= a.created_at - make_interval(secs => $2)
        and e.deleted = false
        and e.caption is null
        and (not @include_subject::bool or e.subject_id = a.subject_id)
)
select * from events where id = any(select id from activity) order by (created_at, id) asc;

-- name: GetEventsInGroup :many
select * from events where group_id = @group_id and deleted = false order by(created_at, id) asc;

-- name: GetActorForGroup :one
select actor_id from events where group_id = @group_id and deleted = false order by(created_at, id) asc limit 1;

-- name: HasLaterGroupedEvent :one
select exists(
  select 1 from events where deleted = false
  and group_id = @group_id
  and id > @event_id
);

-- name: IsActorActionActive :one
select exists(
  select 1 from events where deleted = false
  and actor_id = $1
  and action = any(@actions)
  and created_at > @window_start and created_at <= @window_end
);

-- name: IsActorSubjectActive :one
select exists(
  select 1 from events where deleted = false
  and actor_id = $1
  and subject_id = $2
  and created_at > @window_start and created_at <= @window_end
);

-- name: IsActorGalleryActive :one
select exists(
  select 1 from events where deleted = false
  and actor_id = $1
  and gallery_id = $2
  and created_at > @window_start and created_at <= @window_end
);


-- name: IsActorSubjectActionActive :one
select exists(
  select 1 from events where deleted = false
  and actor_id = $1
  and subject_id = $2
  and action = any(@actions)
  and created_at > @window_start and created_at <= @window_end
);

-- name: PaginateGlobalFeed :batchmany
SELECT * FROM feed_events WHERE deleted = false
    AND (event_time, id) < (sqlc.arg('cur_before_time'), sqlc.arg('cur_before_id'))
    AND (event_time, id) > (sqlc.arg('cur_after_time'), sqlc.arg('cur_after_id'))
    ORDER BY CASE WHEN sqlc.arg('paging_forward')::bool THEN (event_time, id) END ASC,
            CASE WHEN NOT sqlc.arg('paging_forward')::bool THEN (event_time, id) END DESC
    LIMIT sqlc.arg('limit');

-- name: PaginatePersonalFeedByUserID :batchmany
SELECT fe.* FROM feed_events fe, follows fl WHERE fe.deleted = false AND fl.deleted = false
    AND fe.owner_id = fl.followee AND fl.follower = sqlc.arg('follower')
    AND (fe.event_time, fe.id) < (sqlc.arg('cur_before_time'), sqlc.arg('cur_before_id'))
    AND (fe.event_time, fe.id) > (sqlc.arg('cur_after_time'), sqlc.arg('cur_after_id'))
    ORDER BY CASE WHEN sqlc.arg('paging_forward')::bool THEN (fe.event_time, fe.id) END ASC,
            CASE WHEN NOT sqlc.arg('paging_forward')::bool THEN (fe.event_time, fe.id) END DESC
    LIMIT sqlc.arg('limit');

-- name: PaginateUserFeedByUserID :batchmany
SELECT * FROM feed_events WHERE owner_id = sqlc.arg('owner_id') AND deleted = false
    AND (event_time, id) < (sqlc.arg('cur_before_time'), sqlc.arg('cur_before_id'))
    AND (event_time, id) > (sqlc.arg('cur_after_time'), sqlc.arg('cur_after_id'))
    ORDER BY CASE WHEN sqlc.arg('paging_forward')::bool THEN (event_time, id) END ASC,
            CASE WHEN NOT sqlc.arg('paging_forward')::bool THEN (event_time, id) END DESC
    LIMIT sqlc.arg('limit');

-- name: PaginateTrendingFeed :many
select f.* from feed_events f join unnest(@feed_event_ids::text[]) with ordinality t(id, pos) using(id) where f.deleted = false
  and t.pos > @cur_before_pos::int
  and t.pos < @cur_after_pos::int
  order by case when @paging_forward::bool then t.pos end desc,
          case when not @paging_forward::bool then t.pos end asc
  limit sqlc.arg('limit');

-- name: GetEventByIdBatch :batchone
SELECT * FROM feed_events WHERE id = $1 AND deleted = false;

-- name: CreateFeedEvent :one
INSERT INTO feed_events (id, owner_id, action, data, event_time, event_ids, group_id, caption) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *;

-- name: IsFeedEventExistsForGroup :one
SELECT exists(
  SELECT 1 FROM feed_events WHERE deleted = false
  AND group_id = $1
);

-- name: UpdateFeedEventCaptionByGroup :one
UPDATE feed_events SET caption = (select caption from events where events.group_id = $1) WHERE group_id = $1 AND deleted = false returning *;

-- name: GetLastFeedEventForUser :one
select * from feed_events where deleted = false
    and owner_id = $1
    and action = any(@actions)
    and event_time < $2
    order by event_time desc
    limit 1;

-- name: GetLastFeedEventForToken :one
select * from feed_events where deleted = false
    and owner_id = $1
    and action = any(@actions)
    and data ->> 'token_id' = @token_id::varchar
    and event_time < $2
    order by event_time desc
    limit 1;

-- name: GetLastFeedEventForCollection :one
select * from feed_events where deleted = false
    and owner_id = $1
    and action = any(@actions)
    and data ->> 'collection_id' = @collection_id
    and event_time < $2
    order by event_time desc
    limit 1;

-- name: IsFeedUserActionBlocked :one
SELECT EXISTS(SELECT 1 FROM feed_blocklist WHERE user_id = $1 AND (action = $2 or action = '') AND deleted = false);

-- name: BlockUserFromFeed :exec
INSERT INTO feed_blocklist (id, user_id, action) VALUES ($1, $2, $3);

-- name: UnblockUserFromFeed :exec
UPDATE feed_blocklist SET deleted = true WHERE user_id = $1;

-- name: GetAdmireByAdmireID :one
SELECT * FROM admires WHERE id = $1 AND deleted = false;

-- name: GetAdmiresByAdmireIDs :many
SELECT * from admires WHERE id = ANY(@admire_ids) AND deleted = false;

-- name: GetAdmireByAdmireIDBatch :batchone
SELECT * FROM admires WHERE id = $1 AND deleted = false;

-- name: GetAdmiresByActorID :many
SELECT * FROM admires WHERE actor_id = $1 AND deleted = false ORDER BY created_at DESC;

-- name: GetAdmiresByActorIDBatch :batchmany
SELECT * FROM admires WHERE actor_id = $1 AND deleted = false ORDER BY created_at DESC;

-- name: PaginateAdmiresByFeedEventIDBatch :batchmany
SELECT * FROM admires WHERE feed_event_id = sqlc.arg('feed_event_id') AND deleted = false
    AND (created_at, id) < (sqlc.arg('cur_before_time'), sqlc.arg('cur_before_id')) AND (created_at, id) > (sqlc.arg('cur_after_time'), sqlc.arg('cur_after_id'))
    ORDER BY CASE WHEN sqlc.arg('paging_forward')::bool THEN (created_at, id) END ASC,
             CASE WHEN NOT sqlc.arg('paging_forward')::bool THEN (created_at, id) END DESC
    LIMIT sqlc.arg('limit');

-- name: CountAdmiresByFeedEventIDBatch :batchone
SELECT count(*) FROM admires WHERE feed_event_id = $1 AND deleted = false;

-- name: GetCommentByCommentID :one
SELECT * FROM comments WHERE id = $1 AND deleted = false;

-- name: GetCommentsByCommentIDs :many
SELECT * from comments WHERE id = ANY(@comment_ids) AND deleted = false;

-- name: GetCommentByCommentIDBatch :batchone
SELECT * FROM comments WHERE id = $1 AND deleted = false;

-- name: PaginateCommentsByFeedEventIDBatch :batchmany
SELECT * FROM comments WHERE feed_event_id = sqlc.arg('feed_event_id') AND deleted = false
    AND (created_at, id) < (sqlc.arg('cur_before_time'), sqlc.arg('cur_before_id'))
    AND (created_at, id) > (sqlc.arg('cur_after_time'), sqlc.arg('cur_after_id'))
    ORDER BY CASE WHEN sqlc.arg('paging_forward')::bool THEN (created_at, id) END ASC,
             CASE WHEN NOT sqlc.arg('paging_forward')::bool THEN (created_at, id) END DESC
    LIMIT sqlc.arg('limit');

-- name: CountCommentsByFeedEventIDBatch :batchone
SELECT count(*) FROM comments WHERE feed_event_id = $1 AND deleted = false;

-- name: GetCommentsByActorID :many
SELECT * FROM comments WHERE actor_id = $1 AND deleted = false ORDER BY created_at DESC;

-- name: GetCommentsByActorIDBatch :batchmany
SELECT * FROM comments WHERE actor_id = $1 AND deleted = false ORDER BY created_at DESC;

-- name: GetUserNotifications :many
SELECT * FROM notifications WHERE owner_id = $1 AND deleted = false
    AND (created_at, id) < (@cur_before_time, @cur_before_id)
    AND (created_at, id) > (@cur_after_time, @cur_after_id)
    ORDER BY CASE WHEN @paging_forward::bool THEN (created_at, id) END ASC,
             CASE WHEN NOT @paging_forward::bool THEN (created_at, id) END DESC
    LIMIT $2;

-- name: GetUserUnseenNotifications :many
SELECT * FROM notifications WHERE owner_id = $1 AND deleted = false AND seen = false
    AND (created_at, id) < (@cur_before_time, @cur_before_id)
    AND (created_at, id) > (@cur_after_time, @cur_after_id)
    ORDER BY CASE WHEN @paging_forward::bool THEN (created_at, id) END ASC,
             CASE WHEN NOT @paging_forward::bool THEN (created_at, id) END DESC
    LIMIT $2;

-- name: GetRecentUnseenNotifications :many
SELECT * FROM notifications WHERE owner_id = @owner_id AND deleted = false AND seen = false and created_at > @created_after order by created_at desc limit @lim;

-- name: GetUserNotificationsBatch :batchmany
SELECT * FROM notifications WHERE owner_id = sqlc.arg('owner_id') AND deleted = false
    AND (created_at, id) < (sqlc.arg('cur_before_time'), sqlc.arg('cur_before_id'))
    AND (created_at, id) > (sqlc.arg('cur_after_time'), sqlc.arg('cur_after_id'))
    ORDER BY CASE WHEN sqlc.arg('paging_forward')::bool THEN (created_at, id) END ASC,
             CASE WHEN NOT sqlc.arg('paging_forward')::bool THEN (created_at, id) END DESC
    LIMIT sqlc.arg('limit');

-- name: CountUserNotifications :one
SELECT count(*) FROM notifications WHERE owner_id = $1 AND deleted = false;

-- name: CountUserUnseenNotifications :one
SELECT count(*) FROM notifications WHERE owner_id = $1 AND deleted = false AND seen = false;

-- name: GetNotificationByID :one
SELECT * FROM notifications WHERE id = $1 AND deleted = false;

-- name: GetNotificationByIDBatch :batchone
SELECT * FROM notifications WHERE id = $1 AND deleted = false;

-- name: GetMostRecentNotificationByOwnerIDForAction :one
select * from notifications
    where owner_id = $1
    and action = $2
    and deleted = false
    and (not @only_for_feed_event::bool or feed_event_id = $3)
    order by created_at desc
    limit 1;

-- name: GetNotificationsByOwnerIDForActionAfter :many
SELECT * FROM notifications
    WHERE owner_id = $1 AND action = $2 AND deleted = false AND created_at > @created_after
    ORDER BY created_at DESC;

-- name: CreateAdmireNotification :one
INSERT INTO notifications (id, owner_id, action, data, event_ids, feed_event_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: CreateCommentNotification :one
INSERT INTO notifications (id, owner_id, action, data, event_ids, feed_event_id, comment_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: CreateFollowNotification :one
INSERT INTO notifications (id, owner_id, action, data, event_ids) VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: CreateViewGalleryNotification :one
INSERT INTO notifications (id, owner_id, action, data, event_ids, gallery_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: UpdateNotification :exec
UPDATE notifications SET data = $2, event_ids = event_ids || $3, amount = $4, last_updated = now(), seen = false WHERE id = $1 AND deleted = false AND NOT amount = $4;

-- name: UpdateNotificationSettingsByID :exec
UPDATE users SET notification_settings = $2 WHERE id = $1;

-- name: ClearNotificationsForUser :many
UPDATE notifications SET seen = true WHERE owner_id = $1 AND seen = false RETURNING *;

-- name: PaginateInteractionsByFeedEventIDBatch :batchmany
SELECT interactions.created_At, interactions.id, interactions.tag FROM (
    SELECT t.created_at, t.id, sqlc.arg('admire_tag')::int as tag FROM admires t WHERE sqlc.arg('admire_tag') != 0 AND t.feed_event_id = sqlc.arg('feed_event_id') AND t.deleted = false
        AND (sqlc.arg('admire_tag'), t.created_at, t.id) < (sqlc.arg('cur_before_tag')::int, sqlc.arg('cur_before_time'), sqlc.arg('cur_before_id')) AND (sqlc.arg('admire_tag'), t.created_at, t.id) > (sqlc.arg('cur_after_tag')::int, sqlc.arg('cur_after_time'), sqlc.arg('cur_after_id'))
                                                                    UNION
    SELECT t.created_at, t.id, sqlc.arg('comment_tag')::int as tag FROM comments t WHERE sqlc.arg('comment_tag') != 0 AND t.feed_event_id = sqlc.arg('feed_event_id') AND t.deleted = false
        AND (sqlc.arg('comment_tag'), t.created_at, t.id) < (sqlc.arg('cur_before_tag')::int, sqlc.arg('cur_before_time'), sqlc.arg('cur_before_id')) AND (sqlc.arg('comment_tag'), t.created_at, t.id) > (sqlc.arg('cur_after_tag')::int, sqlc.arg('cur_after_time'), sqlc.arg('cur_after_id'))
) as interactions

ORDER BY CASE WHEN sqlc.arg('paging_forward')::bool THEN (tag, created_at, id) END ASC,
         CASE WHEN NOT sqlc.arg('paging_forward')::bool THEN (tag, created_at, id) END DESC
LIMIT sqlc.arg('limit');

-- name: CountInteractionsByFeedEventIDBatch :batchmany
SELECT count(*), sqlc.arg('admire_tag')::int as tag FROM admires t WHERE sqlc.arg('admire_tag') != 0 AND t.feed_event_id = sqlc.arg('feed_event_id') AND t.deleted = false
                                                        UNION
SELECT count(*), sqlc.arg('comment_tag')::int as tag FROM comments t WHERE sqlc.arg('comment_tag') != 0 AND t.feed_event_id = sqlc.arg('feed_event_id') AND t.deleted = false;

-- name: GetAdmireByActorIDAndFeedEventID :batchone
SELECT * FROM admires WHERE actor_id = $1 AND feed_event_id = $2 AND deleted = false;


-- for some reason this query will not allow me to use @tags for $1
-- name: GetUsersWithEmailNotificationsOnForEmailType :many
select * from pii.user_view
    where (email_unsubscriptions->>'all' = 'false' or email_unsubscriptions->>'all' is null)
    and (email_unsubscriptions->>sqlc.arg(email_unsubscription)::varchar = 'false' or email_unsubscriptions->>sqlc.arg(email_unsubscription)::varchar is null)
    and deleted = false and pii_email_address is not null and email_verified = $1
    and (created_at, id) < (@cur_before_time, @cur_before_id)
    and (created_at, id) > (@cur_after_time, @cur_after_id)
    order by case when @paging_forward::bool then (created_at, id) end asc,
             case when not @paging_forward::bool then (created_at, id) end desc
    limit $2;

-- name: GetUsersWithEmailNotificationsOn :many
-- TODO: Does not appear to be used
select * from pii.user_view
    where (email_unsubscriptions->>'all' = 'false' or email_unsubscriptions->>'all' is null)
    and deleted = false and pii_email_address is not null and email_verified = $1
    and (created_at, id) < (@cur_before_time, @cur_before_id)
    and (created_at, id) > (@cur_after_time, @cur_after_id)
    order by case when @paging_forward::bool then (created_at, id) end asc,
             case when not @paging_forward::bool then (created_at, id) end desc
    limit $2;

-- name: GetUsersWithRolePaginate :many
select u.* from users u, user_roles ur where u.deleted = false and ur.deleted = false
    and u.id = ur.user_id and ur.role = @role
    and (u.username_idempotent, u.id) < (@cur_before_key::varchar, @cur_before_id)
    and (u.username_idempotent, u.id) > (@cur_after_key::varchar, @cur_after_id)
    order by case when @paging_forward::bool then (u.username_idempotent, u.id) end asc,
             case when not @paging_forward::bool then (u.username_idempotent, u.id) end desc
    limit $1;

-- name: GetUsersByPositionPaginate :many
select u.* from users u join unnest(@user_ids::text[]) with ordinality t(id, pos) using(id) where u.deleted = false
  and t.pos > @cur_before_pos::int
  and t.pos < @cur_after_pos::int
  order by case when @paging_forward::bool then t.pos end desc,
          case when not @paging_forward::bool then t.pos end asc
  limit sqlc.arg('limit');

-- name: UpdateUserVerificationStatus :exec
UPDATE users SET email_verified = $2 WHERE id = $1;

-- name: UpdateUserEmail :exec
with upsert_pii as (
    insert into pii.for_users (user_id, pii_email_address) values (@user_id, @email_address)
        on conflict (user_id) do update set pii_email_address = excluded.pii_email_address
),

upsert_metadata as (
    insert into dev_metadata_users (user_id, has_email_address) values (@user_id, (@email_address is not null))
        on conflict (user_id) do update set has_email_address = excluded.has_email_address
)

update users set email_verified = 0 where users.id = @user_id;

-- name: UpdateUserEmailUnsubscriptions :exec
UPDATE users SET email_unsubscriptions = $2 WHERE id = $1;

-- name: UpdateUserPrimaryWallet :exec
update users set primary_wallet_id = @wallet_id from wallets
    where users.id = @user_id and wallets.id = @wallet_id
    and wallets.id = any(users.wallets) and wallets.deleted = false;

-- name: GetUsersByChainAddresses :many
select users.*,wallets.address from users, wallets where wallets.address = ANY(@addresses::varchar[]) AND wallets.chain = @chain::int AND ARRAY[wallets.id] <@ users.wallets AND users.deleted = false AND wallets.deleted = false;

-- name: GetFeedEventByID :one
SELECT * FROM feed_events WHERE id = $1 AND deleted = false;

-- name: AddUserRoles :exec
insert into user_roles (id, user_id, role, created_at, last_updated)
select unnest(@ids::varchar[]), $1, unnest(@roles::varchar[]), now(), now()
on conflict (user_id, role) do update set deleted = false, last_updated = now();

-- name: DeleteUserRoles :exec
update user_roles set deleted = true, last_updated = now() where user_id = $1 and role = any(@roles);

-- name: GetUserRolesByUserId :many
with membership_roles(role) as (
    select (case when exists(
        select 1
        from tokens
        where owner_user_id = @user_id
            and token_id = any(@membership_token_ids::varchar[])
            and contract = (select id from contracts where address = @membership_address and contracts.chain = @chain and contracts.deleted = false)
            and exists(select 1 from users where id = @user_id and email_verified = 1 and deleted = false)
            and deleted = false
    ) then @granted_membership_role else null end)::varchar
)
select role from user_roles where user_id = @user_id and deleted = false
union
select role from membership_roles where role is not null;

-- name: RedeemMerch :one
update merch set redeemed = true, token_id = @token_hex, last_updated = now() where id = (select m.id from merch m where m.object_type = @object_type and m.token_id is null and m.redeemed = false and m.deleted = false order by m.id limit 1) and token_id is null and redeemed = false returning discount_code;

-- name: GetMerchDiscountCodeByTokenID :one
select discount_code from merch where token_id = @token_hex and redeemed = true and deleted = false;

-- name: GetUserOwnsTokenByIdentifiers :one
select exists(select 1 from tokens where owner_user_id = @user_id and token_id = @token_hex and contract = @contract and chain = @chain and deleted = false) as owns_token;

-- name: UpdateGalleryHidden :one
update galleries set hidden = @hidden, last_updated = now() where id = @id and deleted = false returning *;

-- name: UpdateGalleryPositions :exec
with updates as (
    select unnest(@gallery_ids::text[]) as id, unnest(@positions::text[]) as position
)
update galleries g set position = updates.position, last_updated = now() from updates where g.id = updates.id and deleted = false and g.owner_user_id = @owner_user_id;

-- name: UserHasDuplicateGalleryPositions :one
select exists(select position,count(*) from galleries where owner_user_id = $1 and deleted = false group by position having count(*) > 1);

-- name: UpdateGalleryInfo :exec
update galleries set name = case when @name_set::bool then @name else name end, description = case when @description_set::bool then @description else description end, last_updated = now() where id = @id and deleted = false;

-- name: UpdateGalleryCollections :exec
update galleries set collections = @collections, last_updated = now() where galleries.id = @gallery_id and galleries.deleted = false and (select count(*) from collections c where c.id = any(@collections) and c.gallery_id = @gallery_id and c.deleted = false) = cardinality(@collections);

-- name: UpdateUserFeaturedGallery :exec
update users set featured_gallery = @gallery_id, last_updated = now() from galleries where users.id = @user_id and galleries.id = @gallery_id and galleries.owner_user_id = @user_id and galleries.deleted = false;

-- name: GetGalleryTokenMediasByGalleryID :many
select m.media from tokens t, collections c, galleries g, token_medias m where g.id = $1 and c.id = any(g.collections) and t.id = any(c.nfts) and t.deleted = false and g.deleted = false and c.deleted = false and (length(t.media->>'thumbnail_url'::varchar) > 0 or length(t.media->>'media_url'::varchar) > 0) and t.token_media_id = m.id and m.deleted = false and m.active order by array_position(g.collections, c.id),array_position(c.nfts, t.id) limit $2;

-- name: GetTokenByTokenIdentifiers :one
select * from tokens where tokens.token_id = @token_hex and contract = (select contracts.id from contracts where contracts.address = @contract_address) and tokens.chain = @chain and tokens.deleted = false;

-- name: GetTokensByIDs :many
select * from tokens join unnest(@token_ids::varchar[]) with ordinality t(id, pos) using (id) where deleted = false order by t.pos asc;

-- name: DeleteCollections :exec
update collections set deleted = true, last_updated = now() where id = any(@ids::varchar[]);

-- name: UpdateCollectionsInfo :exec
with updates as (
    select unnest(@ids::varchar[]) as id, unnest(@names::varchar[]) as name, unnest(@collectors_notes::varchar[]) as collectors_note, unnest(@layouts::jsonb[]) as layout, unnest(@token_settings::jsonb[]) as token_settings, unnest(@hidden::bool[]) as hidden
)
update collections c set collectors_note = updates.collectors_note, layout = updates.layout, token_settings = updates.token_settings, hidden = updates.hidden, name = updates.name, last_updated = now(), version = 1 from updates where c.id = updates.id and c.deleted = false;

-- name: GetCollectionTokensByCollectionID :one
select nfts from collections where id = $1 and deleted = false;

-- name: UpdateCollectionTokens :exec
update collections set nfts = @nfts, last_updated = now() where id = @id and deleted = false;

-- name: CreateCollection :one
insert into collections (id, version, name, collectors_note, owner_user_id, gallery_id, layout, nfts, hidden, token_settings, created_at, last_updated) values (@id, 1, @name, @collectors_note, @owner_user_id, @gallery_id, @layout, @nfts, @hidden, @token_settings, now(), now()) returning id;

-- name: GetGalleryIDByCollectionID :one
select gallery_id from collections where id = $1 and deleted = false;

-- name: GetAllTimeTrendingUserIDs :many
select users.id
from events, galleries, users
left join legacy_views on users.id = legacy_views.user_id and legacy_views.deleted = false
where action = 'ViewedGallery'
  and events.gallery_id = galleries.id
  and users.id = galleries.owner_user_id
  and galleries.deleted = false
  and users.deleted = false
group by users.id
order by row_number() over(order by count(events.id) + coalesce(max(legacy_views.view_count), 0) desc, max(users.created_at) desc) asc
limit $1;

-- name: GetWindowedTrendingUserIDs :many
with viewers as (
  select gallery_id, count(distinct coalesce(actor_id, external_id)) viewer_count
  from events
  where action = 'ViewedGallery' and events.created_at >= @window_end
  group by gallery_id
),
edit_events as (
  select actor_id
  from events
  where action in (
    'CollectionCreated',
    'CollectorsNoteAddedToCollection',
    'CollectorsNoteAddedToToken',
    'TokensAddedToCollection',
    'GalleryInfoUpdated'
  ) and created_at >= @window_end
  group by actor_id
)
select users.id
from viewers, galleries, users, edit_events
where viewers.gallery_id = galleries.id
	and galleries.owner_user_id = users.id
	and users.deleted = false
	and galleries.deleted = false
  and users.id = edit_events.actor_id
group by users.id
order by row_number() over(order by sum(viewers.viewer_count) desc, max(users.created_at) desc) asc
limit $1;

-- name: GetUserExperiencesByUserID :one
select user_experiences from users where id = $1;

-- name: UpdateUserExperience :exec
update users set user_experiences = user_experiences || @experience where id = @user_id;

-- name: GetTrendingUsersByIDs :many
select users.* from users join unnest(@user_ids::varchar[]) with ordinality t(id, pos) using (id) where deleted = false order by t.pos asc;

-- name: GetTrendingFeedEventIDs :many
select feed_events.id, feed_events.created_at, count(*)
from events as interactions, feed_events
where interactions.action IN ('CommentedOnFeedEvent', 'AdmiredFeedEvent') and interactions.created_at >= @window_end and interactions.feed_event_id is not null and interactions.feed_event_id = feed_events.id
group by feed_events.id, feed_events.created_at;

-- name: UpdateCollectionGallery :exec
update collections set gallery_id = @gallery_id, last_updated = now() where id = @id and deleted = false;

-- name: AddCollectionToGallery :exec
update galleries set collections = array_append(collections, @collection_id), last_updated = now() where id = @gallery_id and deleted = false;

-- name: RemoveCollectionFromGallery :exec
update galleries set collections = array_remove(collections, @collection_id), last_updated = now() where id = @gallery_id and deleted = false;

-- name: UserOwnsGallery :one
select exists(select 1 from galleries where id = $1 and owner_user_id = $2 and deleted = false);

-- name: UserOwnsCollection :one
select exists(select 1 from collections where id = $1 and owner_user_id = $2 and deleted = false);

-- name: GetSocialAuthByUserID :one
select * from pii.socials_auth where user_id = $1 and provider = $2 and deleted = false;

-- name: UpsertSocialOAuth :exec
insert into pii.socials_auth (id, user_id, provider, access_token, refresh_token) values (@id, @user_id, @provider, @access_token, @refresh_token) on conflict (user_id, provider) where deleted = false do update set access_token = @access_token, refresh_token = @refresh_token, last_updated = now();

-- name: AddSocialToUser :exec
insert into pii.for_users (user_id, pii_socials) values (@user_id, @socials) on conflict (user_id) where deleted = false do update set pii_socials = for_users.pii_socials || @socials;

-- name: RemoveSocialFromUser :exec
update pii.for_users set pii_socials = pii_socials - @social::varchar where user_id = @user_id;

-- name: GetSocialsByUserID :one
select pii_socials from pii.user_view where id = $1;

-- name: UpdateUserSocials :exec
update pii.for_users set pii_socials = @socials where user_id = @user_id;

-- name: UpdateEventCaptionByGroup :exec
update events set caption = @caption where group_id = @group_id and deleted = false;

-- this query will take in enoug info to create a sort of fake table of social accounts matching them up to users in gallery with twitter connected.
-- it will also go and search for whether the specified user follows any of the users returned
-- name: GetSocialConnectionsPaginate :many
select s.*, user_view.id as user_id, user_view.created_at as user_created_at, (f.id is not null)::bool as already_following
from (select unnest(@social_ids::varchar[]) as social_id, unnest(@social_usernames::varchar[]) as social_username, unnest(@social_displaynames::varchar[]) as social_displayname, unnest(@social_profile_images::varchar[]) as social_profile_image) as s
    inner join pii.user_view on user_view.pii_socials->sqlc.arg('social')::text->>'id'::varchar = s.social_id and user_view.deleted = false
    left outer join follows f on f.follower = @user_id and f.followee = user_view.id and f.deleted = false
where case when @only_unfollowing::bool then f.id is null else true end
    and (f.id is not null,user_view.created_at,user_view.id) < (@cur_before_following::bool, @cur_before_time::timestamptz, @cur_before_id)
    and (f.id is not null,user_view.created_at,user_view.id) > (@cur_after_following::bool, @cur_after_time::timestamptz, @cur_after_id)
order by case when @paging_forward::bool then (f.id is not null,user_view.created_at,user_view.id) end asc,
    case when not @paging_forward::bool then (f.id is not null,user_view.created_at,user_view.id) end desc
limit $1;

-- name: GetSocialConnections :many
select s.*, user_view.id as user_id, user_view.created_at as user_created_at, (f.id is not null)::bool as already_following
from (select unnest(@social_ids::varchar[]) as social_id, unnest(@social_usernames::varchar[]) as social_username, unnest(@social_displaynames::varchar[]) as social_displayname, unnest(@social_profile_images::varchar[]) as social_profile_image) as s
    inner join pii.user_view on user_view.pii_socials->sqlc.arg('social')::text->>'id'::varchar = s.social_id and user_view.deleted = false
    left outer join follows f on f.follower = @user_id and f.followee = user_view.id and f.deleted = false
where case when @only_unfollowing::bool then f.id is null else true end
order by (f.id is not null,user_view.created_at,user_view.id);

-- name: CountSocialConnections :one
select count(*)
from (select unnest(@social_ids::varchar[]) as social_id) as s
    inner join pii.user_view on user_view.pii_socials->sqlc.arg('social')::text->>'id'::varchar = s.social_id and user_view.deleted = false
    left outer join follows f on f.follower = @user_id and f.followee = user_view.id and f.deleted = false
where case when @only_unfollowing::bool then f.id is null else true end;

-- name: AddManyFollows :exec
insert into follows (id, follower, followee, deleted) select unnest(@ids::varchar[]), @follower, unnest(@followees::varchar[]), false on conflict (follower, followee) where deleted = false do update set deleted = false, last_updated = now() returning last_updated > created_at;

-- name: GetSharedFollowersBatchPaginate :batchmany
select users.*, a.created_at followed_on
from users, follows a, follows b
where a.follower = @follower
	and a.followee = b.follower
	and b.followee = @followee
	and users.id = b.follower
	and a.deleted = false
	and b.deleted = false
	and users.deleted = false
  and (a.created_at, users.id) > (sqlc.arg('cur_before_time'), sqlc.arg('cur_before_id'))
  and (a.created_at, users.id) < (sqlc.arg('cur_after_time'), sqlc.arg('cur_after_id'))
order by case when sqlc.arg('paging_forward')::bool then (a.created_at, users.id) end desc,
        case when not sqlc.arg('paging_forward')::bool then (a.created_at, users.id) end asc
limit sqlc.arg('limit');

-- name: CountSharedFollows :one
select count(*)
from users, follows a, follows b
where a.follower = @follower
	and a.followee = b.follower
	and b.followee = @followee
	and users.id = b.follower
	and a.deleted = false
	and b.deleted = false
	and users.deleted = false;

-- name: GetSharedContractsBatchPaginate :batchmany
select contracts.*, a.displayed as displayed_by_user_a, b.displayed as displayed_by_user_b, a.owned_count
from owned_contracts a, owned_contracts b, contracts
left join marketplace_contracts on contracts.id = marketplace_contracts.contract_id
where a.user_id = @user_a_id
  and b.user_id = @user_b_id
  and a.contract_id = b.contract_id
  and a.contract_id = contracts.id
  and marketplace_contracts.contract_id is null
  and contracts.name is not null
  and contracts.name != ''
  and contracts.name != 'Unidentified contract'
  and (
    a.displayed,
    b.displayed,
    a.owned_count,
    contracts.id
  ) > (
    sqlc.arg('cur_before_displayed_by_user_a'),
    sqlc.arg('cur_before_displayed_by_user_b'),
    sqlc.arg('cur_before_owned_count')::int,
    sqlc.arg('cur_before_contract_id')
  )
  and (
    a.displayed,
    b.displayed,
    a.owned_count,
    contracts.id
  ) < (
    sqlc.arg('cur_after_displayed_by_user_a'),
    sqlc.arg('cur_after_displayed_by_user_b'),
    sqlc.arg('cur_after_owned_count')::int,
    sqlc.arg('cur_after_contract_id')
  )
order by case when sqlc.arg('paging_forward')::bool then (a.displayed, b.displayed, a.owned_count, contracts.id) end desc,
        case when not sqlc.arg('paging_forward')::bool then (a.displayed, b.displayed, a.owned_count, contracts.id) end asc
limit sqlc.arg('limit');

-- name: CountSharedContracts :one
select count(*)
from owned_contracts a, owned_contracts b, contracts
left join marketplace_contracts on contracts.id = marketplace_contracts.contract_id
where a.user_id = @user_a_id
  and b.user_id = @user_b_id
  and a.contract_id = b.contract_id
  and a.contract_id = contracts.id
  and marketplace_contracts.contract_id is null
  and contracts.name is not null
  and contracts.name != ''
  and contracts.name != 'Unidentified contract';

-- name: AddPiiAccountCreationInfo :exec
insert into pii.account_creation_info (user_id, ip_address, created_at) values (@user_id, @ip_address, now())
  on conflict do nothing;

-- name: GetUserByWalletID :one
select * from users where array[@wallet::varchar]::varchar[] <@ wallets and deleted = false;

-- name: DeleteUserByID :exec
update users set deleted = true where id = $1;

-- name: InsertWallet :exec
insert into wallets (id, address, chain, wallet_type) values ($1, $2, $3, $4);

-- name: DeleteWalletByID :exec
update wallets set deleted = true, last_updated = now() where id = $1;

-- name: InsertUser :exec
insert into users (id, username, username_idempotent, bio, wallets, universal, email_unsubscriptions, primary_wallet_id) values ($1, $2, $3, $4, $5, $6, $7, $8) returning id;

-- name: AddWalletToUserByID :exec
update users set wallets = array_append(wallets, @wallet_id::varchar) where id = @user_id;

-- name: IsExistsActiveTokenMediaByTokenIdentifers :one
select exists(select 1 from token_medias where token_medias.contract_id = $1 and token_medias.token_id = $2 and token_medias.chain = $3 and active = true and deleted = false);

-- name: InsertTokenPipelineResults :exec
with insert_job(id) as (
    insert into token_processing_jobs (id, token_properties, pipeline_metadata, processing_cause, processor_version)
    values (@processing_job_id, @token_properties, @pipeline_metadata, @processing_cause, @processor_version)
    returning id
),
-- Optionally create an inactive record of the existing active record if the new media is also active
insert_media_move_active_record(last_updated) as (
    insert into token_medias (id, contract_id, token_id, chain, metadata, media, name, description, processing_job_id, active, created_at, last_updated)
    (
        select @copy_media_id, contract_id, token_id, chain, metadata, media, name, description, processing_job_id, false, created_at, now()
        from token_medias
        where contract_id = @contract_id
            and token_id = @token_id
            and chain = @chain
            and active
            and not deleted
            and @active = true
        limit 1
    )
    returning last_updated
),
-- Update the existing active record with the new media data
insert_media_add_record(insert_id, active, is_new) as (
    insert into token_medias (id, contract_id, token_id, chain, metadata, media, name, description, processing_job_id, active, created_at, last_updated)
    values (@new_media_id, @contract_id, @token_id, @chain, @metadata, @media, @name, @description, (select id from insert_job), @active,
        -- Using timestamps generated from insert_media_move_active_record ensures that the new record is only inserted after the current media is moved
        (select coalesce((select last_updated from insert_media_move_active_record), now())),
        (select coalesce((select last_updated from insert_media_move_active_record), now()))
    )
    on conflict (contract_id, token_id, chain) where active and not deleted do update
        set metadata = excluded.metadata,
            media = excluded.media,
            name = excluded.name,
            description = excluded.description,
            processing_job_id = excluded.processing_job_id,
            last_updated = now()
    returning id as insert_id, active, id = @new_media_id is_new
),
-- This will return the existing active record if it exists. If the incoming record is active,
-- this will still return the active record before the update, and not the new record.
existing_active(id) as (
    select id
    from token_medias
    where chain = @chain and contract_id = @contract_id and token_id = @token_id and active and not deleted
    limit 1
)
update tokens
set token_media_id = (
    case
        -- The pipeline didn't produce active media, but one already exists so use that one
        when not insert_medias.active and (select id from existing_active) is not null
        then (select id from existing_active)

        -- The pipeline produced active media, or didn't produce active media but no active media existed before
        else insert_medias.insert_id
    end
)
from insert_media_add_record insert_medias
where
    tokens.chain = @chain
    and tokens.contract = @contract_id
    and tokens.token_id = @token_id
    and not tokens.deleted
    and (
        -- The case statement below handles which token instances get updated:
        case
            -- If the active media already existed, update tokens that have no media (new tokens that haven't been processed before) or tokens that don't use this media yet
            when insert_medias.active and not insert_medias.is_new
            then (tokens.token_media_id is null or tokens.token_media_id != insert_medias.insert_id)

            -- Brand new active media, update all tokens in the filter to use this media
            when insert_medias.active and insert_medias.is_new
            then 1 = 1

            -- The pipeline run produced inactive media, only update the token instance (since it may have not been processed before)
            -- Since there is no db constraint on inactive media, all inactive media is new
            when not insert_medias.active
            then tokens.id = @token_dbid

            else 1 = 1
        end
    );

-- name: InsertSpamContracts :exec
with insert_spam_contracts as (
    insert into alchemy_spam_contracts (id, chain, address, created_at, is_spam) (
        select unnest(@id::varchar[])
        , unnest(@chain::int[])
        , unnest(@address::varchar[])
        , unnest(@created_at::timestamptz[])
        , unnest(@is_spam::bool[])
    ) on conflict(chain, address) do update set created_at = excluded.created_at, is_spam = excluded.is_spam
    returning created_at
)
delete from alchemy_spam_contracts where created_at < (select created_at from insert_spam_contracts limit 1);

-- name: GetPushTokenByPushToken :one
select * from push_notification_tokens where push_token = @push_token and deleted = false;

-- name: CreatePushTokenForUser :one
insert into push_notification_tokens (id, user_id, push_token, created_at, deleted) values (@id, @user_id, @push_token, now(), false) returning *;

-- name: DeletePushTokensByIDs :exec
update push_notification_tokens set deleted = true where id = any(@ids) and deleted = false;

-- name: GetPushTokensByUserID :many
select * from push_notification_tokens where user_id = @user_id and deleted = false;

-- name: GetPushTokensByIDs :many
select t.* from unnest(@ids::text[]) ids join push_notification_tokens t on t.id = ids and t.deleted = false;

-- name: CreatePushTickets :exec
insert into push_notification_tickets (id, push_token_id, ticket_id, created_at, check_after, num_check_attempts, status, deleted) values
  (
   unnest(@ids::text[]),
   unnest(@push_token_ids::text[]),
   unnest(@ticket_ids::text[]),
   now(),
   now() + interval '15 minutes',
   0,
   'pending',
   false
  );

-- name: UpdatePushTickets :exec
with updates as (
    select unnest(@ids::text[]) as id, unnest(@check_after::timestamptz[]) as check_after, unnest(@num_check_attempts::int[]) as num_check_attempts, unnest(@status::text[]) as status, unnest(@deleted::bool[]) as deleted
)
update push_notification_tickets t set check_after = updates.check_after, num_check_attempts = updates.num_check_attempts, status = updates.status, deleted = updates.deleted from updates where t.id = updates.id and t.deleted = false;

-- name: GetCheckablePushTickets :many
select * from push_notification_tickets where check_after <= now() and deleted = false limit sqlc.arg('limit');

-- name: GetAllTokensWithContractsByIDs :many
SELECT
    tokens.*,
    contracts.*,
    (
        SELECT wallets.address
        FROM wallets
        WHERE wallets.id = ANY(tokens.owned_by_wallets) and wallets.deleted = false
        LIMIT 1
    ) AS wallet_address
FROM tokens
JOIN contracts ON contracts.id = tokens.contract
LEFT JOIN token_medias on token_medias.id = tokens.token_media_id
WHERE tokens.deleted = false
AND (tokens.token_media_id IS NULL or token_medias.active = false)
AND tokens.id >= @start_id AND tokens.id < @end_id
ORDER BY tokens.id;

-- name: GetMissingThumbnailTokensByIDRange :many
SELECT
    tokens.*,
    contracts.*,
    (
        SELECT wallets.address
        FROM wallets
        WHERE wallets.id = ANY(tokens.owned_by_wallets) and wallets.deleted = false
        LIMIT 1
    ) AS wallet_address
FROM tokens
JOIN contracts ON contracts.id = tokens.contract
left join token_medias on tokens.token_media_id = token_medias.id where tokens.deleted = false and token_medias.active = true and token_medias.media->>'media_type' = 'html' and (token_medias.media->>'thumbnail_url' is null or token_medias.media->>'thumbnail_url' = '')
AND tokens.id >= @start_id AND tokens.id < @end_id
ORDER BY tokens.id;

-- name: GetSVGTokensWithContractsByIDs :many
SELECT
    tokens.*,
    contracts.*,
    (
        SELECT wallets.address
        FROM wallets
        WHERE wallets.id = ANY(tokens.owned_by_wallets) and wallets.deleted = false
        LIMIT 1
    ) AS wallet_address
FROM tokens
JOIN contracts ON contracts.id = tokens.contract
LEFT JOIN token_medias on token_medias.id = tokens.token_media_id
WHERE tokens.deleted = false
AND token_medias.active = true
AND token_medias.media->>'media_type' = 'svg'
AND tokens.id >= @start_id AND tokens.id < @end_id
ORDER BY tokens.id;

-- name: GetReprocessJobRangeByID :one
select * from reprocess_jobs where id = $1;

-- name: GetMediaByTokenID :batchone
select m.*
from token_medias m
where m.id = (select token_media_id from tokens where tokens.id = $1) and m.active and not m.deleted;

-- name: UpsertSession :one
insert into sessions (id, user_id,
                      created_at, created_with_user_agent, created_with_platform, created_with_os,
                      last_refreshed, last_user_agent, last_platform, last_os, current_refresh_id, active_until, invalidated, last_updated, deleted)
    values (@id, @user_id, now(), @user_agent, @platform, @os, now(), @user_agent, @platform, @os, @current_refresh_id, @active_until, false, now(), false)
    on conflict (id) where deleted = false do update set
        last_refreshed = case when sessions.invalidated then sessions.last_refreshed else excluded.last_refreshed end,
        last_user_agent = case when sessions.invalidated then sessions.last_user_agent else excluded.last_user_agent end,
        last_platform = case when sessions.invalidated then sessions.last_platform else excluded.last_platform end,
        last_os = case when sessions.invalidated then sessions.last_os else excluded.last_os end,
        current_refresh_id = case when sessions.invalidated then sessions.current_refresh_id else excluded.current_refresh_id end,
        last_updated = case when sessions.invalidated then sessions.last_updated else excluded.last_updated end,
        active_until = case when sessions.invalidated then sessions.active_until else greatest(sessions.active_until, excluded.active_until) end
    returning *;

-- name: InvalidateSession :exec
update sessions set invalidated = true, active_until = least(active_until, now()), last_updated = now() where id = @id and deleted = false and invalidated = false;

-- name: UpdateTokenMetadataFieldsByTokenIdentifiers :exec
update tokens set name = @name, description = @description, last_updated = now() where token_id = @token_id and contract = (select id from contracts where address = @contract_address) and deleted = false;

-- name: GetTopCollectionsForCommunity :many
with contract_tokens as (
	select t.id, t.owner_user_id
	from tokens t
	join contracts c on t.contract = c.id
	where not t.deleted and not c.deleted and t.contract = c.id and c.chain = $1 and c.address = $2
),
ranking as (
	select col.id, rank() over (order by count(col.id) desc, col.created_at desc) score
	from collections col
	join contract_tokens on col.owner_user_id = contract_tokens.owner_user_id and contract_tokens.id = any(col.nfts)
	join users on col.owner_user_id = users.id
	where not col.deleted and not col.hidden and not users.deleted
	group by col.id
)
select collections.id from collections join ranking using(id) where score <= 100 order by score asc;

-- name: GetVisibleCollectionsByIDsPaginate :many
select collections.*
from collections, unnest(@collection_ids::varchar[]) with ordinality as t(id, pos)
where collections.id = t.id and not deleted and not hidden and t.pos < @cur_before_pos::int and t.pos > @cur_after_pos::int
order by case when @paging_forward::bool then t.pos end asc, case when not @paging_forward::bool then t.pos end desc
limit $1;
