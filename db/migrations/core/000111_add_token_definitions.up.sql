-- Ensures that changes to a contract are mirrored in the token_definitions table. A unique constraint on (chain, address) already exists to guarantee uniqueness within the table itself.
create unique index contracts_id_chain_address_idx on contracts(id, chain, address);

create table if not exists token_definitions (
  id character varying(255) primary key,
  created_at timestamp with time zone not null default current_timestamp,
  last_updated timestamp with time zone not null default current_timestamp,
  deleted boolean not null default false,
  name character varying,
  description character varying,
  token_type character varying,
  token_id character varying,
  external_url character varying,
  chain integer,
  metadata jsonb,
  fallback_media jsonb,
  contract_address character varying(255) not null,
  contract_id character varying(255),
  token_media_id character varying(255) references token_medias(id),
  foreign key(contract_id, chain, contract_address) references contracts(id, chain, address) on update cascade
);
create unique index if not exists token_definitions_chain_contract_id_token_idx on token_definitions(chain, contract_id, token_id) where not deleted;
create unique index if not exists token_definitions_chain_contract_address_token_idx on token_definitions(chain, contract_address, token_id) where not deleted;
create index token_definitions_contract_id_idx on token_definitions(contract_id) where not deleted;

-- Make nullable in a later migration after back filling
alter table tokens add column if not exists token_definition_id character varying(255) references token_definitions(id);

create unique index if not exists tokens_owner_token_definition_idx on tokens(owner_user_id, token_definition_id) where not deleted;
alter table tokens rename column contract to contract__deprecated;
alter table tokens rename column chain to chain__deprecated;
alter table tokens rename column token_id to token_id__deprecated;
alter table tokens rename column name to name__deprecated;
alter table tokens rename column description to description__deprecated;
alter table tokens rename column token_type to token_type__deprecated;
alter table tokens rename column ownership_history to ownership_history__deprecated;
alter table tokens rename column external_url to external_url__deprecated;
alter table tokens rename column is_provider_marked_spam to is_provider_marked_spam__deprecated;
alter table tokens rename column token_uri to token_uri__deprecated;
alter table tokens rename column fallback_media to fallback_media__deprecated;
alter table tokens rename column token_media_id to token_media_id__deprecated;

alter table token_medias rename column name to name__deprecated;
alter table token_medias rename column description to description__deprecated;
alter table token_medias rename column metadata to metadata__deprecated;
alter table token_medias rename column contract_id to contract_id__deprecated;
alter table token_medias rename column token_id to token_id__deprecated;
alter table token_medias rename column chain to chain__deprecated;
