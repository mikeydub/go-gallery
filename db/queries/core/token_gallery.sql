-- name: UpsertTokens :many
insert into tokens
(
  id
  , deleted
  , version
  , created_at
  , last_updated
  , name
  , description
  , collectors_note
  , token_type
  , token_id
  , quantity
  , ownership_history
  , media
  , fallback_media
  , token_metadata
  , external_url
  , block_number
  , owner_user_id
  , owned_by_wallets
  , is_creator_token
  , chain
  , contract
  , is_provider_marked_spam
  , last_synced
  , token_uri
  , token_media_id
) (
  select
    id
    , false
    , version
    , now()
    , now()
    , name
    , description
    , collectors_note
    , token_type
    , token_id
    , quantity
    , case when @set_holder_fields::bool then ownership_history[ownership_history_start_idx::int:ownership_history_end_idx::int] else '{}' end
    , media
    , fallback_media
    , token_metadata
    , external_url
    , block_number
    , owner_user_id
    , case when @set_holder_fields then owned_by_wallets[owned_by_wallets_start_idx::int:owned_by_wallets_end_idx::int] else '{}' end
    , case when @set_creator_fields::bool then is_creator_token else false end
    , chain
    , contract
    , is_provider_marked_spam
    , now()
    , token_uri
    , (select tm.id
       from token_medias tm
       where tm.token_id = bulk_upsert.token_id
         and tm.contract_id = bulk_upsert.contract
         and tm.chain = bulk_upsert.chain
         and tm.active = true
         and tm.deleted = false
        limit 1
      ) as token_media_id
  from (
    select unnest(@id::varchar[]) as id
      , unnest(@version::int[]) as version
      , unnest(@name::varchar[]) as name
      , unnest(@description::varchar[]) as description
      , unnest(@collectors_note::varchar[]) as collectors_note
      , unnest(@token_type::varchar[]) as token_type
      , unnest(@quantity::varchar[]) as quantity
      , @ownership_history::jsonb[] as ownership_history
      , unnest(@ownership_history_start_idx::int[]) as ownership_history_start_idx
      , unnest(@ownership_history_end_idx::int[]) as ownership_history_end_idx
      , unnest(@media::jsonb[]) as media
      , unnest(@fallback_media::jsonb[]) as fallback_media
      , unnest(@token_metadata::jsonb[]) as token_metadata
      , unnest(@external_url::varchar[]) as external_url
      , unnest(@block_number::bigint[]) as block_number
      , unnest(@owner_user_id::varchar[]) as owner_user_id
      , @owned_by_wallets::varchar[] as owned_by_wallets
      , unnest(@owned_by_wallets_start_idx::int[]) as owned_by_wallets_start_idx
      , unnest(@owned_by_wallets_end_idx::int[]) as owned_by_wallets_end_idx
      , unnest(@is_creator_token::bool[]) as is_creator_token
      , unnest(@is_provider_marked_spam::bool[]) as is_provider_marked_spam
      , unnest(@token_uri::varchar[]) as token_uri
      , unnest(@token_id::varchar[]) as token_id
      , unnest(@contract::varchar[]) as contract
      , unnest(@chain::int[]) as chain
  ) bulk_upsert
)
on conflict (token_id, contract, chain, owner_user_id) where deleted = false
do update set
  token_type = excluded.token_type
  , name = excluded.name
  , description = excluded.description
  , token_uri = excluded.token_uri
  , quantity = excluded.quantity
  , owned_by_wallets = case when @set_holder_fields then excluded.owned_by_wallets else tokens.owned_by_wallets end
  , ownership_history = case when @set_holder_fields then tokens.ownership_history || excluded.ownership_history else tokens.ownership_history end
  , is_creator_token = case when @set_creator_fields then excluded.is_creator_token else tokens.is_creator_token end
  , fallback_media = excluded.fallback_media
  , token_metadata = excluded.token_metadata
  , external_url = excluded.external_url
  , block_number = excluded.block_number
  , version = excluded.version
  , last_updated = excluded.last_updated
  , is_provider_marked_spam = excluded.is_provider_marked_spam
  , last_synced = greatest(excluded.last_synced,tokens.last_synced)
returning *;

-- name: DeleteTokensBeforeTimestamp :execrows
update tokens t
set owned_by_wallets = case when @remove_holder_status::bool then '{}' else owned_by_wallets end,
    is_creator_token = case when @remove_creator_status::bool then false else is_creator_token end,
    last_updated = now()
where
  -- Guard against only_from_user_id and only_from_contract_ids both being null/empty, as this would
  -- result in deleting more tokens than intended.
  (sqlc.narg('only_from_user_id')::text is not null or cardinality(@only_from_contract_ids::text[]) > 0)
  and (sqlc.narg('only_from_user_id') is null or owner_user_id = @only_from_user_id)
  and (cardinality(@only_from_contract_ids) = 0 or contract = any(@only_from_contract_ids))
  and (cardinality(@only_from_chains::int[]) = 0 or chain = any(@only_from_chains))
  and deleted = false
  and ((@remove_holder_status and is_holder_token) or (@remove_creator_status and is_creator_token))
  and last_synced < @timestamp;