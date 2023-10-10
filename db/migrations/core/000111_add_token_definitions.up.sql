create table if not exists token_definitions (
  id character varying(255) primary key,
  created_at timestamp with time zone not null default current_timestamp,
  last_updated timestamp with time zone not null default current_timestamp,
  deleted boolean not null default false,
  name character varying,
  description character varying,
  token_type character varying,
  token_id character varying,
  token_uri character varying,
  external_url character varying,
  chain integer,
  is_provider_marked_spam boolean,
  fallback_media jsonb,
  contract_id character varying(255) references contracts(id),
  token_media_id character varying(255) references token_media(id)
);

create unique index if not exists token_definitions_chain_contract_token_idx on token_definitions(chain, contract_id, token_id) where not deleted;
create index token_definitions_chain_contract_token_idx on token_definitions(chain, contract_id, token_id) where not deleted;

alter table tokens add column if not exists token_definition_id character varying(255) references token_definitions(id);
alter table tokens rename column contract to contract_id;
alter table tokens rename column name to name__deprecated;
alter table tokens rename column description to description__deprecated;
alter table tokens rename column token_type to token_type__deprecated;
alter table tokens rename column ownership_history to ownership_history__deprecated;
alter table tokens rename column external_url to external_url__deprecated;
alter table tokens rename column block_number to block_number__deprecated;
alter table tokens rename column is_provider_marked_spam to is_provider_marked_spam__deprecated;
alter table tokens rename column token_uri to token_uri__deprecated;
alter table tokens rename column fallback_media to fallback_media__deprecated;
alter table tokens rename column token_media_id to token_media_id__deprecated;

alter table token_medias rename column name to name__deprecated;
alter table token_medias rename column description to description__deprecated;
-- XXX: alter table token_medias rename column contract_id to contract_id__deprecated;
-- XXX: alter table token_medias rename column token_id to token_id__deprecated;
-- XXX: alter table token_medias rename column chain to chain__deprecated;
