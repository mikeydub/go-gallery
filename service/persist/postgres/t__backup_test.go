package postgres

import (
	"context"
	"testing"

	"github.com/mikeydub/go-gallery/service/memstore/redis"
	"github.com/mikeydub/go-gallery/service/persist"
)

func TestBackupRestore_Success(t *testing.T) {
	a, db := setupTest(t)

	galleryRepo := NewGalleryRepository(db, redis.NewCache(0))
	collectionRepo := NewCollectionRepository(db, galleryRepo)
	nftRepo := NewNFTRepository(db, galleryRepo)
	userRepo := NewUserRepository(db)
	backupRepo := NewBackupRepository(db)

	user := persist.User{

		Username:           "username",
		UsernameIdempotent: "username-idempotent",
		Addresses: []persist.Address{
			"0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
		},
	}

	userID, err := userRepo.Create(context.Background(), user)
	a.NoError(err)

	nfts := []persist.NFT{
		{
			OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
			Name:         "name",
			OpenseaID:    1,
		},
		{
			OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
			Name:         "blah blah",
			OpenseaID:    10,
		},
	}

	ids, err := nftRepo.CreateBulk(context.Background(), nfts)
	a.NoError(err)

	collection := persist.CollectionDB{
		Name:        "name",
		OwnerUserID: userID,
		NFTs:        ids,
	}

	collID, err := collectionRepo.Create(context.Background(), collection)
	a.NoError(err)

	gallery := persist.GalleryDB{
		OwnerUserID: userID,
		Collections: []persist.DBID{collID},
	}

	id, err := galleryRepo.Create(context.Background(), gallery)
	a.NoError(err)

	g, err := galleryRepo.GetByID(context.Background(), id)
	a.NoError(err)

	err = backupRepo.Insert(context.Background(), g)
	a.NoError(err)

	backups, err := backupRepo.Get(context.Background(), userID)
	a.NoError(err)
	a.Len(backups, 1)
	a.Len(backups[0].Gallery.Collections, 1)

	collection2 := persist.CollectionDB{
		Name:        "name2",
		OwnerUserID: userID,
		NFTs:        ids,
	}

	collID2, err := collectionRepo.Create(context.Background(), collection2)
	a.NoError(err)

	err = galleryRepo.Update(context.Background(), id, userID, persist.GalleryUpdateInput{
		Collections: []persist.DBID{collID, collID2},
	})
	a.NoError(err)

	err = nftRepo.UpdateByID(context.Background(), ids[0], userID, persist.NFTUpdateOwnerAddressInput{
		OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d3",
	})
	a.NoError(err)

	err = backupRepo.Restore(context.Background(), userID, backups[0].ID)
	a.NoError(err)

	galleries, err := galleryRepo.GetByUserID(context.Background(), userID)
	a.NoError(err)

	a.Len(galleries, 1)

	a.Len(galleries[0].Collections, 1)
}

func TestOneBackupADay_Success(t *testing.T) {
	a, db := setupTest(t)

	galleryRepo := NewGalleryRepository(db, redis.NewCache(0))
	collectionRepo := NewCollectionRepository(db, galleryRepo)
	nftRepo := NewNFTRepository(db, galleryRepo)
	userRepo := NewUserRepository(db)
	backupRepo := NewBackupRepository(db)

	user := persist.User{

		Username:           "username",
		UsernameIdempotent: "username-idempotent",
		Addresses: []persist.Address{
			"0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
		},
	}

	userID, err := userRepo.Create(context.Background(), user)
	a.NoError(err)

	nfts := []persist.NFT{
		{
			OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
			Name:         "name",
			OpenseaID:    1,
		},
		{
			OwnerAddress: "0x8914496dc01efcc49a2fa340331fb90969b6f1d2",
			Name:         "blah blah",
			OpenseaID:    10,
		},
	}

	ids, err := nftRepo.CreateBulk(context.Background(), nfts)
	a.NoError(err)

	collection := persist.CollectionDB{
		Name:        "name",
		OwnerUserID: userID,
		NFTs:        ids,
	}

	collID, err := collectionRepo.Create(context.Background(), collection)
	a.NoError(err)

	gallery := persist.GalleryDB{
		OwnerUserID: userID,
		Collections: []persist.DBID{collID},
	}

	id, err := galleryRepo.Create(context.Background(), gallery)
	a.NoError(err)

	g, err := galleryRepo.GetByID(context.Background(), id)
	a.NoError(err)

	err = backupRepo.Insert(context.Background(), g)
	a.NoError(err)

	backups, err := backupRepo.Get(context.Background(), userID)
	a.NoError(err)
	a.Len(backups, 1)
	a.Len(backups[0].Gallery.Collections, 1)

	collection2 := persist.CollectionDB{
		Name:        "name2",
		OwnerUserID: userID,
		NFTs:        ids,
	}

	collID2, err := collectionRepo.Create(context.Background(), collection2)
	a.NoError(err)

	err = galleryRepo.Update(context.Background(), id, userID, persist.GalleryUpdateInput{
		Collections: []persist.DBID{collID, collID2},
	})
	a.NoError(err)

	err = backupRepo.Insert(context.Background(), g)
	a.NoError(err)

	backups, err = backupRepo.Get(context.Background(), userID)
	a.NoError(err)
	a.Len(backups, 1)
	a.Len(backups[0].Gallery.Collections, 1)
}
