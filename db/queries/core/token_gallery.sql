-- name: UpsertTokens :many
with token_definitions_insert as (
  insert into token_definitions
  ( 
    id
    , created_at
    , last_updated
    , deleted
    , name
    , description
    , token_type
    , token_id
    , external_url
    , chain
    , fallback_media
    , contract_address
    , contract_id
    , metadata
  ) (
    select unnest(@definition_dbid::varchar[]) as id
      , now()
      , now()
      , false
      , unnest(@definition_name::varchar[]) as name
      , unnest(@definition_description::varchar[]) as description
      , unnest(@definition_token_type::varchar[]) as token_type
      , unnest(@definition_token_id::varchar[]) as token_id
      , unnest(@definition_external_url::varchar[]) as external_url
      , unnest(@definition_chain::int[]) as chain
      , unnest(@definition_fallback_media::jsonb[]) as fallback_media
      , unnest(@definition_contract_address::varchar[]) as contract_address
      , unnest(@definition_contract_id::varchar[]) as contract_id
      , unnest(@definition_metadata::jsonb[]) as metadata
  )
  on conflict (chain, contract_id, token_id) where deleted = false
  do update set
    last_updated = excluded.last_updated
    , name = coalesce(nullif(excluded.name, ''), nullif(token_definitions.name, ''))
    , description = coalesce(nullif(excluded.description, ''), nullif(token_definitions.description, ''))
    , token_type = excluded.token_type
    , external_url = coalesce(nullif(excluded.external_url, ''), nullif(token_definitions.external_url, ''))
    , fallback_media = excluded.fallback_media
    , contract_address = excluded.contract_address
    , metadata = excluded.metadata
  returning *
)
, tokens_insert as (
  insert into tokens
  (
    id
    , deleted
    , version
    , created_at
    , last_updated
    , collectors_note
    , quantity
    , block_number
    , owner_user_id
    , owned_by_wallets
    , is_creator_token
    , last_synced
    , token_definition_id
    , contract_id
  ) (
    select
      bulk_upsert.id
      , false
      , bulk_upsert.version
      , now()
      , now()
      , bulk_upsert.collectors_note
      , bulk_upsert.quantity
      , bulk_upsert.block_number
      , bulk_upsert.owner_user_id
      , case when @set_holder_fields::bool then bulk_upsert.owned_by_wallets[bulk_upsert.owned_by_wallets_start_idx::int:bulk_upsert.owned_by_wallets_end_idx::int] else '{}' end
      , case when @set_creator_fields::bool then bulk_upsert.is_creator_token else false end
      , now()
      , token_definitions_insert.id
      , bulk_upsert.contract_id
    from (
      select unnest(@token_dbid::varchar[]) as id
        , unnest(@token_version::int[]) as version
        , unnest(@token_collectors_note::varchar[]) as collectors_note
        , unnest(@token_quantity::varchar[]) as quantity
        , unnest(@token_block_number::bigint[]) as block_number
        , unnest(@token_owner_user_id::varchar[]) as owner_user_id
        , @token_owned_by_wallets::varchar[] as owned_by_wallets
        , unnest(@token_owned_by_wallets_start_idx::int[]) as owned_by_wallets_start_idx
        , unnest(@token_owned_by_wallets_end_idx::int[]) as owned_by_wallets_end_idx
        , unnest(@token_is_creator_token::bool[]) as is_creator_token
        , unnest(@token_token_id::varchar[]) as token_id
        , unnest(@token_contract_address::varchar[]) as contract_address
        , unnest(@token_chain::int[]) as chain
        , unnest(@token_contract_id::varchar[]) as contract_id
    ) bulk_upsert
    join token_definitions_insert on (bulk_upsert.chain, bulk_upsert.contract_address, bulk_upsert.token_id) = (token_definitions_insert.chain, token_definitions_insert.contract_address, token_definitions_insert.token_id)
  )
  on conflict (owner_user_id, token_definition_id) where deleted = false
  do update set
    quantity = excluded.quantity
    , owned_by_wallets = case when @set_holder_fields then excluded.owned_by_wallets else tokens.owned_by_wallets end
    , is_creator_token = case when @set_creator_fields then excluded.is_creator_token else tokens.is_creator_token end
    , block_number = excluded.block_number
    , version = excluded.version
    , last_updated = excluded.last_updated
    , last_synced = greatest(excluded.last_synced,tokens.last_synced)
    , contract_id = excluded.contract_id
  returning *
)
select sqlc.embed(tokens), sqlc.embed(token_definitions), sqlc.embed(contracts)
from tokens_insert tokens
join token_definitions_insert token_definitions on tokens.token_definition_id = token_definitions.id
join contracts on token_definitions.contract_id = contracts.id;

-- name: DeleteTokensBeforeTimestamp :execrows
update tokens t
set owned_by_wallets = case when @remove_holder_status::bool then '{}' else owned_by_wallets end,
    is_creator_token = case when @remove_creator_status::bool then false else is_creator_token end,
    last_updated = now()
from token_definitions td
where
  -- Guard against only_from_user_id and only_from_contract_ids both being null/empty, as this would
  -- result in deleting more tokens than intended.
  (sqlc.narg('only_from_user_id')::text is not null or cardinality(@only_from_contract_ids::text[]) > 0)
  and (sqlc.narg('only_from_user_id') is null or t.owner_user_id = @only_from_user_id)
  and (cardinality(@only_from_contract_ids) = 0 or td.contract_id = any(@only_from_contract_ids))
  and (cardinality(@only_from_chains::int[]) = 0 or td.chain = any(@only_from_chains))
  and t.deleted = false
  and ((@remove_holder_status and t.is_holder_token) or (@remove_creator_status and t.is_creator_token))
  and t.last_synced < @timestamp;
