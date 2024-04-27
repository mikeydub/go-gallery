// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package mirrordb

import (
	"time"

	"github.com/jackc/pgtype"
	"github.com/mikeydub/go-gallery/service/persist"
)

type BaseCollection struct {
	ID                    string     `db:"id" json:"id"`
	SimplehashLookupNftID string     `db:"simplehash_lookup_nft_id" json:"simplehash_lookup_nft_id"`
	LastSimplehashSync    *time.Time `db:"last_simplehash_sync" json:"last_simplehash_sync"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
	LastUpdated           time.Time  `db:"last_updated" json:"last_updated"`
}

type BaseContract struct {
	Address               string     `db:"address" json:"address"`
	SimplehashLookupNftID string     `db:"simplehash_lookup_nft_id" json:"simplehash_lookup_nft_id"`
	LastSimplehashSync    *time.Time `db:"last_simplehash_sync" json:"last_simplehash_sync"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
	LastUpdated           time.Time  `db:"last_updated" json:"last_updated"`
}

type BaseOwner struct {
	SimplehashKafkaKey       string           `db:"simplehash_kafka_key" json:"simplehash_kafka_key"`
	SimplehashNftID          *string          `db:"simplehash_nft_id" json:"simplehash_nft_id"`
	ContractAddress          *persist.Address `db:"contract_address" json:"contract_address"`
	TokenID                  pgtype.Numeric   `db:"token_id" json:"token_id"`
	OwnerAddress             *persist.Address `db:"owner_address" json:"owner_address"`
	Quantity                 pgtype.Numeric   `db:"quantity" json:"quantity"`
	CollectionID             *string          `db:"collection_id" json:"collection_id"`
	FirstAcquiredDate        *time.Time       `db:"first_acquired_date" json:"first_acquired_date"`
	LastAcquiredDate         *time.Time       `db:"last_acquired_date" json:"last_acquired_date"`
	FirstAcquiredTransaction *string          `db:"first_acquired_transaction" json:"first_acquired_transaction"`
	LastAcquiredTransaction  *string          `db:"last_acquired_transaction" json:"last_acquired_transaction"`
	MintedToThisWallet       *bool            `db:"minted_to_this_wallet" json:"minted_to_this_wallet"`
	AirdroppedToThisWallet   *bool            `db:"airdropped_to_this_wallet" json:"airdropped_to_this_wallet"`
	SoldToThisWallet         *bool            `db:"sold_to_this_wallet" json:"sold_to_this_wallet"`
	CreatedAt                time.Time        `db:"created_at" json:"created_at"`
	LastUpdated              time.Time        `db:"last_updated" json:"last_updated"`
	KafkaOffset              *int64           `db:"kafka_offset" json:"kafka_offset"`
	KafkaPartition           *int32           `db:"kafka_partition" json:"kafka_partition"`
	KafkaTimestamp           *time.Time       `db:"kafka_timestamp" json:"kafka_timestamp"`
}

type BaseToken struct {
	SimplehashKafkaKey string           `db:"simplehash_kafka_key" json:"simplehash_kafka_key"`
	SimplehashNftID    *string          `db:"simplehash_nft_id" json:"simplehash_nft_id"`
	ContractAddress    *persist.Address `db:"contract_address" json:"contract_address"`
	TokenID            pgtype.Numeric   `db:"token_id" json:"token_id"`
	Name               *string          `db:"name" json:"name"`
	Description        *string          `db:"description" json:"description"`
	Previews           pgtype.JSONB     `db:"previews" json:"previews"`
	ImageUrl           *string          `db:"image_url" json:"image_url"`
	VideoUrl           *string          `db:"video_url" json:"video_url"`
	AudioUrl           *string          `db:"audio_url" json:"audio_url"`
	ModelUrl           *string          `db:"model_url" json:"model_url"`
	OtherUrl           *string          `db:"other_url" json:"other_url"`
	BackgroundColor    *string          `db:"background_color" json:"background_color"`
	ExternalUrl        *string          `db:"external_url" json:"external_url"`
	OnChainCreatedDate *time.Time       `db:"on_chain_created_date" json:"on_chain_created_date"`
	Status             *string          `db:"status" json:"status"`
	TokenCount         pgtype.Numeric   `db:"token_count" json:"token_count"`
	OwnerCount         pgtype.Numeric   `db:"owner_count" json:"owner_count"`
	Contract           pgtype.JSONB     `db:"contract" json:"contract"`
	CollectionID       *string          `db:"collection_id" json:"collection_id"`
	LastSale           pgtype.JSONB     `db:"last_sale" json:"last_sale"`
	FirstCreated       pgtype.JSONB     `db:"first_created" json:"first_created"`
	Rarity             pgtype.JSONB     `db:"rarity" json:"rarity"`
	ExtraMetadata      *string          `db:"extra_metadata" json:"extra_metadata"`
	ImageProperties    pgtype.JSONB     `db:"image_properties" json:"image_properties"`
	VideoProperties    pgtype.JSONB     `db:"video_properties" json:"video_properties"`
	AudioProperties    pgtype.JSONB     `db:"audio_properties" json:"audio_properties"`
	ModelProperties    pgtype.JSONB     `db:"model_properties" json:"model_properties"`
	OtherProperties    pgtype.JSONB     `db:"other_properties" json:"other_properties"`
	CreatedAt          time.Time        `db:"created_at" json:"created_at"`
	LastUpdated        time.Time        `db:"last_updated" json:"last_updated"`
	KafkaOffset        *int64           `db:"kafka_offset" json:"kafka_offset"`
	KafkaPartition     *int32           `db:"kafka_partition" json:"kafka_partition"`
	KafkaTimestamp     *time.Time       `db:"kafka_timestamp" json:"kafka_timestamp"`
}

type EthereumCollection struct {
	ID                    string     `db:"id" json:"id"`
	SimplehashLookupNftID string     `db:"simplehash_lookup_nft_id" json:"simplehash_lookup_nft_id"`
	LastSimplehashSync    *time.Time `db:"last_simplehash_sync" json:"last_simplehash_sync"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
	LastUpdated           time.Time  `db:"last_updated" json:"last_updated"`
}

type EthereumContract struct {
	Address               string     `db:"address" json:"address"`
	SimplehashLookupNftID string     `db:"simplehash_lookup_nft_id" json:"simplehash_lookup_nft_id"`
	LastSimplehashSync    *time.Time `db:"last_simplehash_sync" json:"last_simplehash_sync"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
	LastUpdated           time.Time  `db:"last_updated" json:"last_updated"`
}

type EthereumOwner struct {
	SimplehashKafkaKey       string           `db:"simplehash_kafka_key" json:"simplehash_kafka_key"`
	SimplehashNftID          *string          `db:"simplehash_nft_id" json:"simplehash_nft_id"`
	ContractAddress          *persist.Address `db:"contract_address" json:"contract_address"`
	TokenID                  pgtype.Numeric   `db:"token_id" json:"token_id"`
	OwnerAddress             *persist.Address `db:"owner_address" json:"owner_address"`
	Quantity                 pgtype.Numeric   `db:"quantity" json:"quantity"`
	CollectionID             *string          `db:"collection_id" json:"collection_id"`
	FirstAcquiredDate        *time.Time       `db:"first_acquired_date" json:"first_acquired_date"`
	LastAcquiredDate         *time.Time       `db:"last_acquired_date" json:"last_acquired_date"`
	FirstAcquiredTransaction *string          `db:"first_acquired_transaction" json:"first_acquired_transaction"`
	LastAcquiredTransaction  *string          `db:"last_acquired_transaction" json:"last_acquired_transaction"`
	MintedToThisWallet       *bool            `db:"minted_to_this_wallet" json:"minted_to_this_wallet"`
	AirdroppedToThisWallet   *bool            `db:"airdropped_to_this_wallet" json:"airdropped_to_this_wallet"`
	SoldToThisWallet         *bool            `db:"sold_to_this_wallet" json:"sold_to_this_wallet"`
	CreatedAt                time.Time        `db:"created_at" json:"created_at"`
	LastUpdated              time.Time        `db:"last_updated" json:"last_updated"`
	KafkaOffset              *int64           `db:"kafka_offset" json:"kafka_offset"`
	KafkaPartition           *int32           `db:"kafka_partition" json:"kafka_partition"`
	KafkaTimestamp           *time.Time       `db:"kafka_timestamp" json:"kafka_timestamp"`
}

type EthereumToken struct {
	SimplehashKafkaKey string           `db:"simplehash_kafka_key" json:"simplehash_kafka_key"`
	SimplehashNftID    *string          `db:"simplehash_nft_id" json:"simplehash_nft_id"`
	ContractAddress    *persist.Address `db:"contract_address" json:"contract_address"`
	TokenID            pgtype.Numeric   `db:"token_id" json:"token_id"`
	Name               *string          `db:"name" json:"name"`
	Description        *string          `db:"description" json:"description"`
	Previews           pgtype.JSONB     `db:"previews" json:"previews"`
	ImageUrl           *string          `db:"image_url" json:"image_url"`
	VideoUrl           *string          `db:"video_url" json:"video_url"`
	AudioUrl           *string          `db:"audio_url" json:"audio_url"`
	ModelUrl           *string          `db:"model_url" json:"model_url"`
	OtherUrl           *string          `db:"other_url" json:"other_url"`
	BackgroundColor    *string          `db:"background_color" json:"background_color"`
	ExternalUrl        *string          `db:"external_url" json:"external_url"`
	OnChainCreatedDate *time.Time       `db:"on_chain_created_date" json:"on_chain_created_date"`
	Status             *string          `db:"status" json:"status"`
	TokenCount         pgtype.Numeric   `db:"token_count" json:"token_count"`
	OwnerCount         pgtype.Numeric   `db:"owner_count" json:"owner_count"`
	Contract           pgtype.JSONB     `db:"contract" json:"contract"`
	CollectionID       *string          `db:"collection_id" json:"collection_id"`
	LastSale           pgtype.JSONB     `db:"last_sale" json:"last_sale"`
	FirstCreated       pgtype.JSONB     `db:"first_created" json:"first_created"`
	Rarity             pgtype.JSONB     `db:"rarity" json:"rarity"`
	ExtraMetadata      *string          `db:"extra_metadata" json:"extra_metadata"`
	ImageProperties    pgtype.JSONB     `db:"image_properties" json:"image_properties"`
	VideoProperties    pgtype.JSONB     `db:"video_properties" json:"video_properties"`
	AudioProperties    pgtype.JSONB     `db:"audio_properties" json:"audio_properties"`
	ModelProperties    pgtype.JSONB     `db:"model_properties" json:"model_properties"`
	OtherProperties    pgtype.JSONB     `db:"other_properties" json:"other_properties"`
	CreatedAt          time.Time        `db:"created_at" json:"created_at"`
	LastUpdated        time.Time        `db:"last_updated" json:"last_updated"`
	KafkaOffset        *int64           `db:"kafka_offset" json:"kafka_offset"`
	KafkaPartition     *int32           `db:"kafka_partition" json:"kafka_partition"`
	KafkaTimestamp     *time.Time       `db:"kafka_timestamp" json:"kafka_timestamp"`
}

type ZoraCollection struct {
	ID                    string     `db:"id" json:"id"`
	SimplehashLookupNftID string     `db:"simplehash_lookup_nft_id" json:"simplehash_lookup_nft_id"`
	LastSimplehashSync    *time.Time `db:"last_simplehash_sync" json:"last_simplehash_sync"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
	LastUpdated           time.Time  `db:"last_updated" json:"last_updated"`
}

type ZoraContract struct {
	Address               string     `db:"address" json:"address"`
	SimplehashLookupNftID string     `db:"simplehash_lookup_nft_id" json:"simplehash_lookup_nft_id"`
	LastSimplehashSync    *time.Time `db:"last_simplehash_sync" json:"last_simplehash_sync"`
	CreatedAt             time.Time  `db:"created_at" json:"created_at"`
	LastUpdated           time.Time  `db:"last_updated" json:"last_updated"`
}

type ZoraOwner struct {
	SimplehashKafkaKey       string           `db:"simplehash_kafka_key" json:"simplehash_kafka_key"`
	SimplehashNftID          *string          `db:"simplehash_nft_id" json:"simplehash_nft_id"`
	ContractAddress          *persist.Address `db:"contract_address" json:"contract_address"`
	TokenID                  pgtype.Numeric   `db:"token_id" json:"token_id"`
	OwnerAddress             *persist.Address `db:"owner_address" json:"owner_address"`
	Quantity                 pgtype.Numeric   `db:"quantity" json:"quantity"`
	CollectionID             *string          `db:"collection_id" json:"collection_id"`
	FirstAcquiredDate        *time.Time       `db:"first_acquired_date" json:"first_acquired_date"`
	LastAcquiredDate         *time.Time       `db:"last_acquired_date" json:"last_acquired_date"`
	FirstAcquiredTransaction *string          `db:"first_acquired_transaction" json:"first_acquired_transaction"`
	LastAcquiredTransaction  *string          `db:"last_acquired_transaction" json:"last_acquired_transaction"`
	MintedToThisWallet       *bool            `db:"minted_to_this_wallet" json:"minted_to_this_wallet"`
	AirdroppedToThisWallet   *bool            `db:"airdropped_to_this_wallet" json:"airdropped_to_this_wallet"`
	SoldToThisWallet         *bool            `db:"sold_to_this_wallet" json:"sold_to_this_wallet"`
	CreatedAt                time.Time        `db:"created_at" json:"created_at"`
	LastUpdated              time.Time        `db:"last_updated" json:"last_updated"`
	KafkaOffset              *int64           `db:"kafka_offset" json:"kafka_offset"`
	KafkaPartition           *int32           `db:"kafka_partition" json:"kafka_partition"`
	KafkaTimestamp           *time.Time       `db:"kafka_timestamp" json:"kafka_timestamp"`
}

type ZoraToken struct {
	SimplehashKafkaKey string           `db:"simplehash_kafka_key" json:"simplehash_kafka_key"`
	SimplehashNftID    *string          `db:"simplehash_nft_id" json:"simplehash_nft_id"`
	ContractAddress    *persist.Address `db:"contract_address" json:"contract_address"`
	TokenID            pgtype.Numeric   `db:"token_id" json:"token_id"`
	Name               *string          `db:"name" json:"name"`
	Description        *string          `db:"description" json:"description"`
	Previews           pgtype.JSONB     `db:"previews" json:"previews"`
	ImageUrl           *string          `db:"image_url" json:"image_url"`
	VideoUrl           *string          `db:"video_url" json:"video_url"`
	AudioUrl           *string          `db:"audio_url" json:"audio_url"`
	ModelUrl           *string          `db:"model_url" json:"model_url"`
	OtherUrl           *string          `db:"other_url" json:"other_url"`
	BackgroundColor    *string          `db:"background_color" json:"background_color"`
	ExternalUrl        *string          `db:"external_url" json:"external_url"`
	OnChainCreatedDate *time.Time       `db:"on_chain_created_date" json:"on_chain_created_date"`
	Status             *string          `db:"status" json:"status"`
	TokenCount         pgtype.Numeric   `db:"token_count" json:"token_count"`
	OwnerCount         pgtype.Numeric   `db:"owner_count" json:"owner_count"`
	Contract           pgtype.JSONB     `db:"contract" json:"contract"`
	CollectionID       *string          `db:"collection_id" json:"collection_id"`
	LastSale           pgtype.JSONB     `db:"last_sale" json:"last_sale"`
	FirstCreated       pgtype.JSONB     `db:"first_created" json:"first_created"`
	Rarity             pgtype.JSONB     `db:"rarity" json:"rarity"`
	ExtraMetadata      *string          `db:"extra_metadata" json:"extra_metadata"`
	ImageProperties    pgtype.JSONB     `db:"image_properties" json:"image_properties"`
	VideoProperties    pgtype.JSONB     `db:"video_properties" json:"video_properties"`
	AudioProperties    pgtype.JSONB     `db:"audio_properties" json:"audio_properties"`
	ModelProperties    pgtype.JSONB     `db:"model_properties" json:"model_properties"`
	OtherProperties    pgtype.JSONB     `db:"other_properties" json:"other_properties"`
	CreatedAt          time.Time        `db:"created_at" json:"created_at"`
	LastUpdated        time.Time        `db:"last_updated" json:"last_updated"`
	KafkaOffset        *int64           `db:"kafka_offset" json:"kafka_offset"`
	KafkaPartition     *int32           `db:"kafka_partition" json:"kafka_partition"`
	KafkaTimestamp     *time.Time       `db:"kafka_timestamp" json:"kafka_timestamp"`
}
