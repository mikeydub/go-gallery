CREATE TABLE users (
    ID varchar(32) PRIMARY KEY,
    DELETED boolean NOT NULL DEFAULT false,
    VERSION int,
    LAST_UPDATED timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CREATED_AT timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    USERNAME varchar(255),
    USERNAME_IDEMPOTENT varchar(255),
    ADDRESSES varchar(255) []
);

CREATE UNIQUE INDEX users_username_idempotent ON users (USERNAME_IDEMPOTENT);

CREATE TABLE galleries (
    ID varchar(32) PRIMARY KEY,
    DELETED boolean NOT NULL DEFAULT false,
    OWNER_USER_ID varchar(32),
    COLLECTIONS varchar(32) [],
    VERSION int
);

CREATE TABLE nfts (
    ID varchar(32) PRIMARY KEY,
    DELETED boolean NOT NULL DEFAULT false,
    VERSION int,
    LAST_UPDATED timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    NAME varchar,
    DESCRIPTION varchar,
    EXTERNAL_URL varchar,
    CREATOR_ADDRESS varchar(255),
    CREATOR_NAME varchar,
    OWNER_ADDRESS char(42),
    MULTIPLE_OWNERS boolean,
    CONTRACT json,
    OPENSEA_ID bigint,
    OPENSEA_TOKEN_ID varchar(255),
    IMAGE_URL varchar,
    IMAGE_THUMBNAIL_URL varchar,
    IMAGE_PREVIEW_URL varchar,
    IMAGE_ORIGINAL_URL varchar,
    ANIMATION_URL varchar,
    ANIMATION_ORIGINAL_URL varchar
);

CREATE TABLE collections (
    ID varchar(32) PRIMARY KEY,
    DELETED boolean NOT NULL DEFAULT false,
    OWNER_USER_ID varchar(32),
    NFTS varchar(32) [],
    VERSION int,
    LAST_UPDATED timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CREATED_AT timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    HIDDEN boolean NOT NULL DEFAULT false,
    COLLECTORS_NOTE varchar,
    NAME varchar(255),
    LAYOUT json
);

CREATE TABLE nonces (
    ID varchar(32) PRIMARY KEY,
    DELETED boolean NOT NULL DEFAULT false,
    VERSION int,
    LAST_UPDATED timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CREATED_AT timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    USER_ID varchar(32),
    ADDRESS char(42),
    VALUE varchar(255)
);

CREATE TABLE tokens (
    ID varchar(32) PRIMARY KEY,
    DELETED boolean NOT NULL DEFAULT false,
    VERSION int,
    CREATED_AT timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    NAME varchar,
    DESCRIPTION varchar,
    CONTRACT_ADDRESS char(42),
    COLLECTORS_NOTE varchar,
    MEDIA json,
    OWNER_ADDRESS char(42),
    TOKEN_URI varchar,
    TOKEN_TYPE varchar,
    TOKEN_ID varchar,
    QUANTITY varchar,
    OWNERSHIP_HISTORY json [],
    TOKEN_METADATA json,
    EXTERNAL_URL varchar,
    BLOCK_NUMBER bigint
);

CREATE UNIQUE INDEX token_id_contract_address_idx ON tokens (TOKEN_ID, CONTRACT_ADDRESS);

CREATE UNIQUE INDEX token_id_contract_address_owner_address_idx ON tokens (TOKEN_ID, CONTRACT_ADDRESS, OWNER_ADDRESS);

CREATE INDEX block_number_idx ON tokens (BLOCK_NUMBER);

CREATE TABLE contracts (
    ID varchar(32) PRIMARY KEY,
    DELETED boolean NOT NULL DEFAULT false,
    VERSION int,
    CREATED_AT timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    LAST_UPDATED timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    NAME varchar,
    SYMBOL varchar,
    ADDRESS char(42),
    LATEST_BLOCK bigint
);

CREATE UNIQUE INDEX address_idx ON contracts (ADDRESS);

CREATE TABLE login_attempts (
    ID varchar(32) PRIMARY KEY,
    DELETED boolean NOT NULL DEFAULT false,
    VERSION int,
    CREATED_AT timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    LAST_UPDATED timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ADDRESS char(42),
    REQUEST_HOST_ADDRESS varchar(255),
    USER_EXISTS boolean,
    SIGNATURE varchar(255),
    SIGNATURE_VALID boolean,
    REQUEST_HEADERS json,
    NONCE_VALUE varchar
);

CREATE TABLE features (
    ID varchar(32) PRIMARY KEY,
    DELETED boolean NOT NULL DEFAULT false,
    VERSION int,
    LAST_UPDATED timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CREATED_AT timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    REQUIRED_TOKEN varchar,
    REQUIRED_AMOUNT bigint,
    TOKEN_TYPE varchar,
    NAME varchar,
    IS_ENABLED boolean,
    ADMIN_ONLY boolean,
    FORCE_ENABLED_USER_IDS varchar(32) []
);

CREATE UNIQUE INDEX feature_name_idx ON features (NAME);

CREATE UNIQUE INDEX feature_required_token_idx ON features (REQUIRED_TOKEN);

CREATE TABLE backups (
    ID varchar(32) PRIMARY KEY,
    DELETED boolean NOT NULL DEFAULT false,
    VERSION int,
    CREATED_AT timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    LAST_UPDATED timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    GALLERY_ID varchar(32),
    GALLERY json
);

CREATE TABLE membership (
    ID varchar(32) PRIMARY KEY,
    DELETED boolean NOT NULL DEFAULT false,
    VERSION int,
    CREATED_AT timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    LAST_UPDATED timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    NAME varchar,
    ASSET_URL varchar,
    OWNERS json []
);

CREATE TABLE access (
    ID varchar(32) PRIMARY KEY,
    DELETED boolean NOT NULL DEFAULT false,
    VERSION int,
    CREATED_AT timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    LAST_UPDATED timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    USER_ID varchar(32),
    MOST_RECENT_BLOCK bigint,
    REQUIRED_TOKENS_OWNED json,
    IS_ADMIN boolean
);