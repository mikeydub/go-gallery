package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/lib/pq"
	"github.com/mikeydub/go-gallery/service/logger"
	"time"

	"github.com/mikeydub/go-gallery/service/memstore"
	"github.com/mikeydub/go-gallery/service/persist"
)

const staleCommunityTime = time.Hour * 2

// CommunityRepository represents a repository for interacting with persisted communities
type CommunityRepository struct {
	cache memstore.Cache
	db    *sql.DB

	getInfoStmt          *sql.Stmt
	getUserByAddressStmt *sql.Stmt
	getPreviewNFTsStmt   *sql.Stmt
}

// NewCommunityRepository returns a new CommunityRepository
func NewCommunityRepository(db *sql.DB, cache memstore.Cache) *CommunityRepository {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	getInfoStmt, err := db.PrepareContext(ctx,
		`SELECT n.OWNER_ADDRESS,n.CONTRACT,n.TOKEN_COLLECTION_NAME,n.CREATOR_ADDRESS,n.IMAGE_PREVIEW_URL
	FROM galleries g, unnest(g.COLLECTIONS) WITH ORDINALITY AS u(coll, coll_ord)
	LEFT JOIN collections c ON c.ID = coll
	LEFT JOIN LATERAL (SELECT n.*,nft,nft_ord FROM nfts n, unnest(c.NFTS) WITH ORDINALITY AS x(nft, nft_ord)) n ON n.ID = n.nft
	WHERE n.CONTRACT ->> 'address' = $1 AND g.DELETED = false AND c.DELETED = false AND n.DELETED = false ORDER BY coll_ord,n.nft_ord;`,
	)
	checkNoErr(err)

	getUserByAddressStmt, err := db.PrepareContext(ctx, `SELECT ID,USERNAME FROM users WHERE ADDRESSES @> ARRAY[$1]:: varchar[] AND DELETED = false`)
	checkNoErr(err)

	getPreviewNFTsStmt, err := db.PrepareContext(ctx, `SELECT IMAGE_THUMBNAIL_URL FROM nfts WHERE CONTRACT->>'address' = $1 AND DELETED = false AND OWNER_ADDRESS = ANY($2) AND LENGTH(IMAGE_THUMBNAIL_URL) > 0 ORDER BY ID LIMIT 3`)
	checkNoErr(err)

	return &CommunityRepository{
		cache:                cache,
		db:                   db,
		getInfoStmt:          getInfoStmt,
		getUserByAddressStmt: getUserByAddressStmt,
		getPreviewNFTsStmt:   getPreviewNFTsStmt,
	}
}

// GetByAddress returns a community by its address
func (c *CommunityRepository) GetByAddress(ctx context.Context, pAddress persist.Address, forceRefresh bool) (persist.Community, error) {
	var community persist.Community

	if !forceRefresh {
		bs, err := c.cache.Get(ctx, pAddress.String())
		if err == nil && len(bs) > 0 {
			err = json.Unmarshal(bs, &community)
			if err != nil {
				return persist.Community{}, err
			}
			return community, nil
		}
	}

	community = persist.Community{ContractAddress: pAddress}

	addresses := make([]persist.Address, 0, 20)
	contract := persist.NFTContract{}

	rows, err := c.getInfoStmt.QueryContext(ctx, pAddress)
	if err != nil {
		return persist.Community{}, fmt.Errorf("error getting community info: %w", err)
	}
	defer rows.Close()

	seenAddress := map[persist.Address]bool{}
	for rows.Next() {
		var address persist.Address
		err = rows.Scan(&address, &contract, &community.Name, &community.CreatorAddress, &community.PreviewImage)
		if err != nil {
			return persist.Community{}, fmt.Errorf("error scanning community info: %w", err)
		}

		if !seenAddress[address] {
			addresses = append(addresses, address)
		}
		seenAddress[address] = true
	}

	if err = rows.Err(); err != nil {
		return persist.Community{}, fmt.Errorf("error getting community info: %w", err)
	}

	if len(seenAddress) == 0 {
		return persist.Community{}, persist.ErrCommunityNotFound{CommunityAddress: pAddress}
	}

	community.Description = contract.ContractDescription

	if contract.ContractImage.String() != "" {
		community.PreviewImage = contract.ContractImage
	}
	if community.Name.String() == "" {
		community.Name = contract.ContractName
	}

	tokenHolders := map[persist.DBID]*persist.TokenHolder{}
	for _, address := range addresses {
		var username persist.NullString
		var userID persist.DBID
		err := c.getUserByAddressStmt.QueryRowContext(ctx, address).Scan(&userID, &username)
		if err != nil {
			logger.For(ctx).Warnf("error getting member of community '%s' by address '%s': %s", pAddress, address, err)
			continue
		}

		if username.String() == "" {
			continue
		}

		if tokenHolder, ok := tokenHolders[userID]; ok {
			tokenHolder.Addresses = append(tokenHolder.Addresses, address)
		} else {
			tokenHolders[userID] = &persist.TokenHolder{
				UserID:      userID,
				Addresses:   []persist.Address{address},
				Username:    username,
				PreviewNFTs: nil,
			}
		}
	}

	community.Owners = make([]persist.TokenHolder, 0, len(tokenHolders))

	for _, tokenHolder := range tokenHolders {
		previewNFTs := make([]persist.NullString, 0, 3)

		rows, err = c.getPreviewNFTsStmt.QueryContext(ctx, pAddress, pq.Array(tokenHolder.Addresses))
		defer rows.Close()

		if err != nil {
			logger.For(ctx).Warnf("error getting preview NFTs of community '%s' by addresses '%s': %s", pAddress, tokenHolder.Addresses, err)
		} else {
			for rows.Next() {
				var imageURL persist.NullString
				err = rows.Scan(&imageURL)
				if err != nil {
					logger.For(ctx).Warnf("error scanning preview NFT of community '%s' by addresses '%s': %s", pAddress, tokenHolder.Addresses, err)
					continue
				}
				previewNFTs = append(previewNFTs, imageURL)
			}
		}

		tokenHolder.PreviewNFTs = previewNFTs
		community.Owners = append(community.Owners, *tokenHolder)
	}

	community.LastUpdated = persist.LastUpdatedTime(time.Now())

	bs, err := json.Marshal(community)
	if err != nil {
		return persist.Community{}, err
	}
	err = c.cache.Set(ctx, pAddress.String(), bs, staleCommunityTime)
	if err != nil {
		return persist.Community{}, err
	}

	return community, nil

}
