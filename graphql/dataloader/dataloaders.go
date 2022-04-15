//go:generate go run github.com/vektah/dataloaden UserLoaderByID github.com/mikeydub/go-gallery/service/persist.DBID github.com/mikeydub/go-gallery/db/sqlc.User
//go:generate go run github.com/vektah/dataloaden UserLoaderByAddress github.com/mikeydub/go-gallery/service/persist.Address github.com/mikeydub/go-gallery/db/sqlc.User
//go:generate go run github.com/vektah/dataloaden UserLoaderByString string github.com/mikeydub/go-gallery/db/sqlc.User
//go:generate go run github.com/vektah/dataloaden GalleryLoaderByID github.com/mikeydub/go-gallery/service/persist.DBID github.com/mikeydub/go-gallery/db/sqlc.Gallery
//go:generate go run github.com/vektah/dataloaden GalleriesLoaderByID github.com/mikeydub/go-gallery/service/persist.DBID []github.com/mikeydub/go-gallery/db/sqlc.Gallery
//go:generate go run github.com/vektah/dataloaden CollectionLoaderByID github.com/mikeydub/go-gallery/service/persist.DBID github.com/mikeydub/go-gallery/db/sqlc.Collection
//go:generate go run github.com/vektah/dataloaden CollectionsLoaderByID github.com/mikeydub/go-gallery/service/persist.DBID []github.com/mikeydub/go-gallery/db/sqlc.Collection
//go:generate go run github.com/vektah/dataloaden NftLoaderByID github.com/mikeydub/go-gallery/service/persist.DBID github.com/mikeydub/go-gallery/db/sqlc.Nft
//go:generate go run github.com/vektah/dataloaden NftsLoaderByID github.com/mikeydub/go-gallery/service/persist.DBID []github.com/mikeydub/go-gallery/db/sqlc.Nft
//go:generate go run github.com/vektah/dataloaden NftsLoaderByAddress github.com/mikeydub/go-gallery/service/persist.Address []github.com/mikeydub/go-gallery/db/sqlc.Nft
//go:generate go run github.com/vektah/dataloaden MembershipLoaderById github.com/mikeydub/go-gallery/service/persist.DBID github.com/mikeydub/go-gallery/db/sqlc.Membership

package dataloader

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/mikeydub/go-gallery/db/sqlc"
	"github.com/mikeydub/go-gallery/service/persist"
)

const defaultMaxBatchOne = 100 // Default for queries that return a single result
const defaultMaxBatchMany = 10 // Default for queries that return many results
const defaultWaitTime = 2 * time.Millisecond

// Loaders will cache and batch lookups. They are short-lived and should never persist beyond
// a single request, nor should they be shared between requests (since the data returned is
// relative to the current request context, including the user and their auth status).
type Loaders struct {

	// Every entry here must have a corresponding entry in the Clear___Caches methods below

	UserByUserId             UserLoaderByID
	UserByUsername           UserLoaderByString
	UserByAddress            UserLoaderByAddress
	GalleryByGalleryId       GalleryLoaderByID
	GalleryByCollectionId    GalleryLoaderByID
	GalleriesByUserId        GalleriesLoaderByID
	CollectionByCollectionId CollectionLoaderByID
	CollectionsByGalleryId   CollectionsLoaderByID
	NftByNftId               NftLoaderByID
	NftsByOwnerAddress       NftsLoaderByAddress
	NftsByCollectionId       NftsLoaderByID
	MembershipByMembershipId MembershipLoaderById
}

func NewLoaders(ctx context.Context, q *sqlc.Queries) *Loaders {
	loaders := &Loaders{}

	loaders.UserByUserId = UserLoaderByID{
		maxBatch: defaultMaxBatchOne,
		wait:     defaultWaitTime,
		fetch:    loadUserByUserId(ctx, loaders, q),
	}

	loaders.UserByUsername = UserLoaderByString{
		maxBatch: defaultMaxBatchOne,
		wait:     defaultWaitTime,
		fetch:    loadUserByUsername(ctx, loaders, q),
	}

	loaders.UserByAddress = UserLoaderByAddress{
		maxBatch: defaultMaxBatchOne,
		wait:     defaultWaitTime,
		fetch:    loadUserByAddress(ctx, loaders, q),
	}

	loaders.GalleryByGalleryId = GalleryLoaderByID{
		maxBatch: defaultMaxBatchOne,
		wait:     defaultWaitTime,
		fetch:    loadGalleryByGalleryId(ctx, loaders, q),
	}

	loaders.GalleryByCollectionId = GalleryLoaderByID{
		maxBatch: defaultMaxBatchOne,
		wait:     defaultWaitTime,
		fetch:    loadGalleryByCollectionId(ctx, loaders, q),
	}

	loaders.GalleriesByUserId = GalleriesLoaderByID{
		maxBatch: defaultMaxBatchMany,
		wait:     defaultWaitTime,
		fetch:    loadGalleriesByUserId(ctx, loaders, q),
	}

	loaders.CollectionByCollectionId = CollectionLoaderByID{
		maxBatch: defaultMaxBatchOne,
		wait:     defaultWaitTime,
		fetch:    loadCollectionByCollectionId(ctx, loaders, q),
	}

	loaders.CollectionsByGalleryId = CollectionsLoaderByID{
		maxBatch: defaultMaxBatchMany,
		wait:     defaultWaitTime,
		fetch:    loadCollectionsByGalleryId(ctx, loaders, q),
	}

	loaders.NftByNftId = NftLoaderByID{
		maxBatch: defaultMaxBatchOne,
		wait:     defaultWaitTime,
		fetch:    loadNftByNftId(ctx, loaders, q),
	}

	loaders.NftsByOwnerAddress = NftsLoaderByAddress{
		maxBatch: defaultMaxBatchMany,
		wait:     defaultWaitTime,
		fetch:    loadNftsByOwnerAddress(ctx, loaders, q),
	}

	loaders.NftsByCollectionId = NftsLoaderByID{
		maxBatch: defaultMaxBatchMany,
		wait:     defaultWaitTime,
		fetch:    loadNftsByCollectionId(ctx, loaders, q),
	}

	loaders.MembershipByMembershipId = MembershipLoaderById{
		maxBatch: defaultMaxBatchOne,
		wait:     defaultWaitTime,
		fetch:    loadMembershipByMembershipId(ctx, loaders, q),
	}

	return loaders
}

// These are pretty verbose and repetitive; hopefully generics make this cleaner in the future

func (l *Loaders) ClearAllCaches() {
	l.ClearUserCaches()
	l.ClearGalleryCaches()
	l.ClearCollectionCaches()
	l.ClearNftCaches()
	l.ClearMembershipCaches()
}

func (l *Loaders) ClearUserCaches() {
	l.UserByUserId.mu.Lock()
	l.UserByUserId.cache = nil
	l.UserByUserId.mu.Unlock()

	l.UserByAddress.mu.Lock()
	l.UserByAddress.cache = nil
	l.UserByAddress.mu.Unlock()

	l.UserByUsername.mu.Lock()
	l.UserByUsername.cache = nil
	l.UserByUsername.mu.Unlock()
}

func (l *Loaders) ClearGalleryCaches() {
	l.GalleryByGalleryId.mu.Lock()
	l.GalleryByGalleryId.cache = nil
	l.GalleryByGalleryId.mu.Unlock()

	l.GalleryByCollectionId.mu.Lock()
	l.GalleryByCollectionId.cache = nil
	l.GalleryByCollectionId.mu.Unlock()

	l.GalleriesByUserId.mu.Lock()
	l.GalleriesByUserId.cache = nil
	l.GalleriesByUserId.mu.Unlock()
}

func (l *Loaders) ClearCollectionCaches() {
	l.CollectionByCollectionId.mu.Lock()
	l.CollectionByCollectionId.cache = nil
	l.CollectionByCollectionId.mu.Unlock()

	l.CollectionsByGalleryId.mu.Lock()
	l.CollectionsByGalleryId.cache = nil
	l.CollectionsByGalleryId.mu.Unlock()
}

func (l *Loaders) ClearNftCaches() {
	l.NftByNftId.mu.Lock()
	l.NftByNftId.cache = nil
	l.NftByNftId.mu.Unlock()

	l.NftsByOwnerAddress.mu.Lock()
	l.NftsByOwnerAddress.cache = nil
	l.NftsByOwnerAddress.mu.Unlock()

	l.NftsByCollectionId.mu.Lock()
	l.NftsByCollectionId.cache = nil
	l.NftsByCollectionId.mu.Unlock()
}

func (l *Loaders) ClearMembershipCaches() {
	l.MembershipByMembershipId.mu.Lock()
	l.MembershipByMembershipId.cache = nil
	l.MembershipByMembershipId.mu.Unlock()
}

func loadUserByUserId(ctx context.Context, loaders *Loaders, q *sqlc.Queries) func([]persist.DBID) ([]sqlc.User, []error) {
	return func(userIds []persist.DBID) ([]sqlc.User, []error) {
		users := make([]sqlc.User, len(userIds))
		errors := make([]error, len(userIds))

		b := q.GetUserByIdBatch(ctx, userIds)
		defer b.Close()

		b.QueryRow(func(i int, user sqlc.User, err error) {
			if err == pgx.ErrNoRows {
				err = persist.ErrUserNotFound{UserID: userIds[i]}
			}

			// Add results to other loaders' caches
			if err == nil {
				loaders.UserByUsername.Prime(user.Username.String, user)
				for _, address := range user.Addresses {
					loaders.UserByAddress.Prime(address, user)
				}
			}

			users[i], errors[i] = user, err
		})

		return users, errors
	}
}

func loadUserByUsername(ctx context.Context, loaders *Loaders, q *sqlc.Queries) func([]string) ([]sqlc.User, []error) {
	return func(usernames []string) ([]sqlc.User, []error) {
		users := make([]sqlc.User, len(usernames))
		errors := make([]error, len(usernames))

		b := q.GetUserByUsernameBatch(ctx, usernames)
		defer b.Close()

		b.QueryRow(func(i int, user sqlc.User, err error) {
			if err == pgx.ErrNoRows {
				err = persist.ErrUserNotFound{Username: usernames[i]}
			}

			// Add results to other loaders' caches
			if err == nil {
				loaders.UserByUserId.Prime(user.ID, user)
				for _, address := range user.Addresses {
					loaders.UserByAddress.Prime(address, user)
				}
			}

			users[i], errors[i] = user, err
		})

		return users, errors
	}
}

func loadUserByAddress(ctx context.Context, loaders *Loaders, q *sqlc.Queries) func([]persist.EthereumAddress) ([]sqlc.User, []error) {
	return func(addresses []persist.EthereumAddress) ([]sqlc.User, []error) {
		users := make([]sqlc.User, len(addresses))
		errors := make([]error, len(addresses))

		addressStrings := make([]string, len(addresses))
		for i, address := range addresses {
			addressStrings[i] = address.String()
		}

		b := q.GetUserByAddressBatch(ctx, addressStrings)
		defer b.Close()

		b.QueryRow(func(i int, user sqlc.User, err error) {
			if err == pgx.ErrNoRows {
				err = persist.ErrUserNotFound{Address: addresses[i]}
			}

			// Add results to other loaders' caches
			if err == nil {
				loaders.UserByUserId.Prime(user.ID, user)
				loaders.UserByUsername.Prime(user.Username.String, user)
			}

			users[i], errors[i] = user, err
		})

		return users, errors
	}
}

func loadGalleryByGalleryId(ctx context.Context, loaders *Loaders, q *sqlc.Queries) func([]persist.DBID) ([]sqlc.Gallery, []error) {
	return func(galleryIds []persist.DBID) ([]sqlc.Gallery, []error) {
		galleries := make([]sqlc.Gallery, len(galleryIds))
		errors := make([]error, len(galleryIds))

		b := q.GetGalleryByIdBatch(ctx, galleryIds)
		defer b.Close()

		b.QueryRow(func(i int, g sqlc.Gallery, err error) {
			galleries[i] = g
			errors[i] = err

			if errors[i] == pgx.ErrNoRows {
				errors[i] = persist.ErrGalleryNotFoundByID{ID: galleryIds[i]}
			}

			// Add results to other loaders' caches
			if errors[i] == nil {
				for _, collection := range galleries[i].Collections {
					loaders.GalleryByCollectionId.Prime(collection, galleries[i])
				}
			}
		})

		return galleries, errors
	}
}

func loadGalleryByCollectionId(ctx context.Context, loaders *Loaders, q *sqlc.Queries) func([]persist.DBID) ([]sqlc.Gallery, []error) {
	return func(collectionIds []persist.DBID) ([]sqlc.Gallery, []error) {
		galleries := make([]sqlc.Gallery, len(collectionIds))
		errors := make([]error, len(collectionIds))

		b := q.GetGalleryByCollectionIdBatch(ctx, collectionIds)
		defer b.Close()

		b.QueryRow(func(i int, g sqlc.Gallery, err error) {
			galleries[i] = g
			errors[i] = err

			if errors[i] == pgx.ErrNoRows {
				errors[i] = persist.ErrGalleryNotFoundByCollectionID{ID: collectionIds[i]}
			}

			// Add results to other loaders' caches
			if errors[i] == nil {
				loaders.GalleryByGalleryId.Prime(galleries[i].ID, galleries[i])
			}
		})

		return galleries, errors
	}
}

func loadGalleriesByUserId(ctx context.Context, loaders *Loaders, q *sqlc.Queries) func([]persist.DBID) ([][]sqlc.Gallery, []error) {
	return func(userIds []persist.DBID) ([][]sqlc.Gallery, []error) {
		galleries := make([][]sqlc.Gallery, len(userIds))
		errors := make([]error, len(userIds))

		b := q.GetGalleriesByUserIdBatch(ctx, userIds)
		defer b.Close()

		b.Query(func(i int, g []sqlc.Gallery, err error) {
			galleries[i] = g
			errors[i] = err

			// Add results to other loaders' caches
			if errors[i] == nil {
				for _, gallery := range galleries[i] {
					loaders.GalleryByGalleryId.Prime(gallery.ID, gallery)
					for _, collection := range gallery.Collections {
						loaders.GalleryByCollectionId.Prime(collection, gallery)
					}
				}
			}
		})

		return galleries, errors
	}
}

func loadCollectionByCollectionId(ctx context.Context, loaders *Loaders, q *sqlc.Queries) func([]persist.DBID) ([]sqlc.Collection, []error) {
	return func(collectionIds []persist.DBID) ([]sqlc.Collection, []error) {
		collections := make([]sqlc.Collection, len(collectionIds))
		errors := make([]error, len(collectionIds))

		b := q.GetCollectionByIdBatch(ctx, collectionIds)
		defer b.Close()

		b.QueryRow(func(i int, c sqlc.Collection, err error) {
			collections[i] = c
			errors[i] = err

			if errors[i] == pgx.ErrNoRows {
				errors[i] = persist.ErrCollectionNotFoundByID{ID: collectionIds[i]}
			}
		})

		return collections, errors
	}
}

func loadCollectionsByGalleryId(ctx context.Context, loaders *Loaders, q *sqlc.Queries) func([]persist.DBID) ([][]sqlc.Collection, []error) {
	return func(galleryIds []persist.DBID) ([][]sqlc.Collection, []error) {
		collections := make([][]sqlc.Collection, len(galleryIds))
		errors := make([]error, len(galleryIds))

		b := q.GetCollectionsByGalleryIdBatch(ctx, galleryIds)
		defer b.Close()

		b.Query(func(i int, c []sqlc.Collection, err error) {
			collections[i] = c
			errors[i] = err

			// Add results to the CollectionByCollectionId loader's cache
			if errors[i] == nil {
				for _, collection := range collections[i] {
					loaders.CollectionByCollectionId.Prime(collection.ID, collection)
				}
			}
		})

		return collections, errors
	}
}

func loadNftByNftId(ctx context.Context, loaders *Loaders, q *sqlc.Queries) func([]persist.DBID) ([]sqlc.Nft, []error) {
	return func(nftIds []persist.DBID) ([]sqlc.Nft, []error) {
		nfts := make([]sqlc.Nft, len(nftIds))
		errors := make([]error, len(nftIds))

		b := q.GetNftByIdBatch(ctx, nftIds)
		defer b.Close()

		b.QueryRow(func(i int, n sqlc.Nft, err error) {
			nfts[i] = n
			errors[i] = err

			if errors[i] == pgx.ErrNoRows {
				errors[i] = persist.ErrNFTNotFoundByID{ID: nftIds[i]}
			}
		})

		return nfts, errors
	}
}

func loadNftsByOwnerAddress(ctx context.Context, loaders *Loaders, q *sqlc.Queries) func([]persist.EthereumAddress) ([][]sqlc.Nft, []error) {
	return func(addresses []persist.EthereumAddress) ([][]sqlc.Nft, []error) {
		nfts := make([][]sqlc.Nft, len(addresses))
		errors := make([]error, len(addresses))

		b := q.GetNftsByOwnerAddressBatch(ctx, addresses)
		defer b.Close()

		b.Query(func(i int, n []sqlc.Nft, err error) {
			nfts[i] = n
			errors[i] = err

			// Add results to the NftByNftId loader's cache
			if errors[i] == nil {
				for _, nft := range nfts[i] {
					loaders.NftByNftId.Prime(nft.ID, nft)
				}
			}
		})

		return nfts, errors
	}
}

func loadNftsByCollectionId(ctx context.Context, loaders *Loaders, q *sqlc.Queries) func([]persist.DBID) ([][]sqlc.Nft, []error) {
	return func(collectionIds []persist.DBID) ([][]sqlc.Nft, []error) {
		nfts := make([][]sqlc.Nft, len(collectionIds))
		errors := make([]error, len(collectionIds))

		b := q.GetNftsByCollectionIdBatch(ctx, collectionIds)
		defer b.Close()

		b.Query(func(i int, n []sqlc.Nft, err error) {
			nfts[i] = n
			errors[i] = err

			// Add results to the NftByNftId loader's cache
			if errors[i] == nil {
				for _, nft := range nfts[i] {
					loaders.NftByNftId.Prime(nft.ID, nft)
				}
			}
		})

		return nfts, errors
	}
}

func loadMembershipByMembershipId(ctx context.Context, loaders *Loaders, q *sqlc.Queries) func([]persist.DBID) ([]sqlc.Membership, []error) {
	return func(membershipIds []persist.DBID) ([]sqlc.Membership, []error) {
		memberships := make([]sqlc.Membership, len(membershipIds))
		errors := make([]error, len(membershipIds))

		b := q.GetMembershipByMembershipIdBatch(ctx, membershipIds)
		defer b.Close()

		b.QueryRow(func(i int, m sqlc.Membership, err error) {
			memberships[i] = m
			errors[i] = err

			if errors[i] == pgx.ErrNoRows {
				errors[i] = persist.ErrMembershipNotFoundByID{ID: membershipIds[i]}
			}
		})

		return memberships, errors
	}
}
