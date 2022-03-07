package publicapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-playground/validator/v10"
	"github.com/mikeydub/go-gallery/graphql/dataloader"
	"github.com/mikeydub/go-gallery/service/persist"
	"github.com/mikeydub/go-gallery/service/pubsub"
	"github.com/mikeydub/go-gallery/validate"
)

const maxNftsPerCollection = 1000

// TODO: Convert this to a validation error, and enforce in all potential contexts here
var errTooManyNFTsInCollection = errors.New(fmt.Sprintf("maximum of %d NFTs in a collection", maxNftsPerCollection))

type PublicCollectionAPI interface {
	CreateCollection(ctx context.Context, galleryID persist.DBID, name string, collectorsNote string, nfts []persist.DBID, layout persist.TokenLayout) (*persist.Collection, error)
	DeleteCollection(ctx context.Context, collectionID persist.DBID) error
	UpdateCollectionNfts(ctx context.Context, collectionID persist.DBID, nfts []persist.DBID, layout persist.TokenLayout) error
	UpdateCollection(ctx context.Context, collectionID persist.DBID, name string, collectorsNote string) error
}

type CollectionAPI struct {
	repos     *persist.Repositories
	loaders   *dataloader.Loaders
	validator *validator.Validate
	ethClient *ethclient.Client
	pubsub    pubsub.PubSub
}

func NewCollectionAPI(repos *persist.Repositories, loaders *dataloader.Loaders, ethClient *ethclient.Client, pubsub pubsub.PubSub) *CollectionAPI {
	return &CollectionAPI{
		repos:     repos,
		loaders:   loaders,
		ethClient: ethClient,
		pubsub:    pubsub,
	}
}

func (api CollectionAPI) CreateCollection(ctx context.Context, galleryID persist.DBID, name string, collectorsNote string, nfts []persist.DBID, layout persist.TokenLayout) (*persist.Collection, error) {
	// Validate
	if err := validateFields(api.validator, validationMap{
		"galleryID": {galleryID, "required"},
		// TODO: Validate name length
		"collectorsNote": {collectorsNote, "medium"},
		"nfts":           {nfts, "required,unique"},
	}); err != nil {
		return nil, err
	}

	layout, err := persist.ValidateLayout(layout, nfts)
	if err != nil {
		return nil, err
	}

	// Sanitize
	name = validate.SanitizationPolicy.Sanitize(name)
	collectorsNote = validate.SanitizationPolicy.Sanitize(collectorsNote)

	userID, err := getAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	collection := persist.CollectionDB{
		OwnerUserID:    userID,
		NFTs:           nfts,
		Layout:         layout,
		Name:           persist.NullString(name),
		CollectorsNote: persist.NullString(collectorsNote),
	}

	collectionID, err := api.repos.CollectionRepository.Create(ctx, collection)
	if err != nil {
		return nil, err
	}

	err = api.repos.GalleryRepository.AddCollections(ctx, galleryID, userID, []persist.DBID{collectionID})
	if err != nil {
		return nil, err
	}

	// TODO: Get a shallow collection instead of a fully unnested one. Can we roll these into a single struct with
	// multiple fields (nftIds, nfts) and assume it's not hydrated if nfts is null? And then maybe include a parameter
	// for whether to hydrate the hierarchy or not?
	createdCollection, err := dataloader.For(ctx).CollectionByCollectionId.Load(collectionID)
	if err != nil {
		return nil, err
	}

	return &createdCollection, nil
}

func (api CollectionAPI) DeleteCollection(ctx context.Context, collectionID persist.DBID) error {
	// Validate
	if err := validateFields(api.validator, validationMap{
		"collectionID": {collectionID, "required"},
	}); err != nil {
		return err
	}

	userID, err := getAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	return api.repos.CollectionRepository.Delete(ctx, collectionID, userID)
}

func (api CollectionAPI) UpdateCollection(ctx context.Context, collectionID persist.DBID, name string, collectorsNote string) error {
	// Validate
	if err := validateFields(api.validator, validationMap{
		"collectionID":   {collectionID, "required"},
		"name":           {name, "required"}, // TODO: Validate length
		"collectorsNote": {collectorsNote, "required,medium"},
	}); err != nil {
		return err
	}

	// Sanitize
	name = validate.SanitizationPolicy.Sanitize(name)
	collectorsNote = validate.SanitizationPolicy.Sanitize(collectorsNote)

	userID, err := getAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	update := persist.CollectionUpdateInfoInput{
		Name:           persist.NullString(name),
		CollectorsNote: persist.NullString(collectorsNote),
	}

	return api.repos.CollectionRepository.Update(ctx, collectionID, userID, update)
}

func (api CollectionAPI) UpdateCollectionNfts(ctx context.Context, collectionID persist.DBID, nfts []persist.DBID, layout persist.TokenLayout) error {
	// Validate
	if err := validateFields(api.validator, validationMap{
		"collectionID": {collectionID, "required"},
		"nfts":         {nfts, "required,unique"},
	}); err != nil {
		return err
	}

	layout, err := persist.ValidateLayout(layout, nfts)
	if err != nil {
		return err
	}

	if len(nfts) > maxNftsPerCollection {
		return errTooManyNFTsInCollection
	}

	userID, err := getAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	update := persist.CollectionUpdateNftsInput{NFTs: nfts, Layout: layout}

	err = api.repos.CollectionRepository.UpdateNFTs(ctx, collectionID, userID, update)
	if err != nil {
		return err
	}

	backupGalleriesForUser(ctx, userID, api.repos)

	return nil
}
