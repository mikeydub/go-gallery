package postgres

import (
	"context"
	"testing"

	"github.com/mikeydub/go-gallery/service/memstore/redis"
	"github.com/mikeydub/go-gallery/service/persist"
)

func TestCollectionGetByUserID_Success(t *testing.T) {
	a, db := setupTest(t)

	collectionRepo := NewCollectionRepository(db)
	nftRepo := NewNFTRepository(db, redis.NewCache(0), redis.NewCache(1))
	userRepo := NewUserRepository(db)

	user := persist.User{

		Username:           "username",
		UsernameIdempotent: "username-idempotent",
		Addresses: []persist.Address{
			"0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
		},
	}

	id, err := userRepo.Create(context.Background(), user)
	a.NoError(err)

	nfts := []persist.NFTDB{
		{
			OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
			Name:         "name",
		},
		{
			OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d1",
			Name:         "blah blah",
		},
	}

	ids, err := nftRepo.CreateBulk(context.Background(), nfts)
	a.NoError(err)

	collection := persist.CollectionDB{
		Name:        "name",
		OwnerUserID: id,
		NFTs:        ids,
	}

	_, err = collectionRepo.Create(context.Background(), collection)
	a.NoError(err)

	collections, err := collectionRepo.GetByUserID(context.Background(), id, true)
	a.NoError(err)

	a.Equal(1, len(collections))

	a.Greater(len(collections[0].NFTs), 0)
}

func TestCollectionGetByID_Success(t *testing.T) {
	a, db := setupTest(t)

	collectionRepo := NewCollectionRepository(db)
	nftRepo := NewNFTRepository(db, redis.NewCache(0), redis.NewCache(1))
	userRepo := NewUserRepository(db)

	user := persist.User{

		Username:           "username",
		UsernameIdempotent: "username-idempotent",
		Addresses: []persist.Address{
			"0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
		},
	}

	id, err := userRepo.Create(context.Background(), user)
	a.NoError(err)
	a.NotEmpty(id)

	nfts := []persist.NFTDB{
		{
			OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d1",
			Name:         "name",
		},
		{
			OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
			Name:         "blah blah",
		},
	}

	ids, err := nftRepo.CreateBulk(context.Background(), nfts)
	a.NoError(err)
	a.NotEmpty(ids)

	collection := persist.CollectionDB{
		Name:        "name",
		OwnerUserID: id,
		NFTs:        ids,
	}

	collID, err := collectionRepo.Create(context.Background(), collection)
	a.NoError(err)
	a.NotEmpty(collID)

	coll, err := collectionRepo.GetByID(context.Background(), collID, true)
	a.NoError(err)

	a.Equal(collection.Name, coll.Name)

	a.Greater(len(coll.NFTs), 0)

}

func TestCollectionUpdate_Success(t *testing.T) {
	a, db := setupTest(t)

	collectionRepo := NewCollectionRepository(db)
	nftRepo := NewNFTRepository(db, redis.NewCache(0), redis.NewCache(1))
	userRepo := NewUserRepository(db)

	user := persist.User{

		Username:           "username",
		UsernameIdempotent: "username-idempotent",
		Addresses: []persist.Address{
			"0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
		},
	}

	userID, err := userRepo.Create(context.Background(), user)
	a.NoError(err)
	a.NotEmpty(userID)

	nfts := []persist.NFTDB{
		{
			OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d1",
			Name:         "name",
		},
		{
			OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
			Name:         "blah blah",
		},
	}

	ids, err := nftRepo.CreateBulk(context.Background(), nfts)
	a.NoError(err)
	a.NotEmpty(ids)

	collection := persist.CollectionDB{
		Name:        "name",
		OwnerUserID: userID,
		NFTs:        ids,
	}

	collID, err := collectionRepo.Create(context.Background(), collection)
	a.NoError(err)
	a.NotEmpty(collID)

	update := persist.CollectionUpdateInfoInput{Name: "new name"}

	err = collectionRepo.Update(context.Background(), collID, userID, update)
	a.NoError(err)

	coll, err := collectionRepo.GetByID(context.Background(), collID, true)
	a.NoError(err)

	a.Equal(update.Name, coll.Name)
}

func TestCollectionUpdateNFTOrder_Success(t *testing.T) {
	a, db := setupTest(t)

	collectionRepo := NewCollectionRepository(db)
	nftRepo := NewNFTRepository(db, redis.NewCache(0), redis.NewCache(1))
	userRepo := NewUserRepository(db)

	user := persist.User{

		Username:           "username",
		UsernameIdempotent: "username-idempotent",
		Addresses: []persist.Address{
			"0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
		},
	}

	userID, err := userRepo.Create(context.Background(), user)
	a.NoError(err)
	a.NotEmpty(userID)

	nfts := []persist.NFTDB{
		{
			OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d1",
			Name:         "name",
		},
		{
			OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
			Name:         "blah blah",
		},
	}

	ids, err := nftRepo.CreateBulk(context.Background(), nfts)
	a.NoError(err)
	a.NotEmpty(ids)

	collection := persist.CollectionDB{
		Name:        "name",
		OwnerUserID: userID,
		NFTs:        ids,
	}

	collID, err := collectionRepo.Create(context.Background(), collection)
	a.NoError(err)
	a.NotEmpty(collID)

	temp := ids[1]
	ids[1] = ids[0]
	ids[0] = temp

	update := persist.CollectionUpdateNftsInput{NFTs: ids}

	err = collectionRepo.UpdateNFTs(context.Background(), collID, userID, update)
	a.NoError(err)

	coll, err := collectionRepo.GetByID(context.Background(), collID, true)
	a.NoError(err)

	idsResult := make([]persist.DBID, len(coll.NFTs))
	for i, resNFT := range coll.NFTs {
		idsResult[i] = resNFT.ID
	}

	a.Equal(ids, idsResult)
}
