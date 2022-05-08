// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/mikeydub/go-gallery/service/persist"
)

type AddUserAddressPayloadOrError interface {
	IsAddUserAddressPayloadOrError()
}

type AuthorizationError interface {
	IsAuthorizationError()
}

type CollectionByIDOrError interface {
	IsCollectionByIDOrError()
}

type CollectionNftByIDOrError interface {
	IsCollectionNftByIDOrError()
}

type CommunityByAddressOrError interface {
	IsCommunityByAddressOrError()
}

type CreateCollectionPayloadOrError interface {
	IsCreateCollectionPayloadOrError()
}

type CreateUserPayloadOrError interface {
	IsCreateUserPayloadOrError()
}

type DeleteCollectionPayloadOrError interface {
	IsDeleteCollectionPayloadOrError()
}

type Error interface {
	IsError()
}

type GalleryUserOrWallet interface {
	IsGalleryUserOrWallet()
}

type GetAuthNoncePayloadOrError interface {
	IsGetAuthNoncePayloadOrError()
}

type LoginPayloadOrError interface {
	IsLoginPayloadOrError()
}

type Media interface {
	IsMedia()
}

type MediaSubtype interface {
	IsMediaSubtype()
}

type NftByIDOrError interface {
	IsNftByIDOrError()
}

type Node interface {
	IsNode()
}

type RefreshOpenSeaNftsPayloadOrError interface {
	IsRefreshOpenSeaNftsPayloadOrError()
}

type RemoveUserAddressesPayloadOrError interface {
	IsRemoveUserAddressesPayloadOrError()
}

type UpdateCollectionHiddenPayloadOrError interface {
	IsUpdateCollectionHiddenPayloadOrError()
}

type UpdateCollectionInfoPayloadOrError interface {
	IsUpdateCollectionInfoPayloadOrError()
}

type UpdateCollectionNftsPayloadOrError interface {
	IsUpdateCollectionNftsPayloadOrError()
}

type UpdateGalleryCollectionsPayloadOrError interface {
	IsUpdateGalleryCollectionsPayloadOrError()
}

type UpdateNftInfoPayloadOrError interface {
	IsUpdateNftInfoPayloadOrError()
}

type UpdateUserInfoPayloadOrError interface {
	IsUpdateUserInfoPayloadOrError()
}

type UserByUsernameOrError interface {
	IsUserByUsernameOrError()
}

type ViewerOrError interface {
	IsViewerOrError()
}

type AddUserAddressPayload struct {
	Viewer *Viewer `json:"viewer"`
}

func (AddUserAddressPayload) IsAddUserAddressPayloadOrError() {}

type AudioMedia struct {
	PreviewURLs      *PreviewURLSet `json:"previewURLs"`
	MediaURL         *string        `json:"mediaURL"`
	MediaType        *string        `json:"mediaType"`
	ContentRenderURL *string        `json:"contentRenderURL"`
}

func (AudioMedia) IsMediaSubtype() {}
func (AudioMedia) IsMedia()        {}

type AuthMechanism struct {
	DebugAuth   *DebugAuth       `json:"debugAuth"`
	EthereumEoa *EthereumEoaAuth `json:"ethereumEoa"`
	GnosisSafe  *GnosisSafeAuth  `json:"gnosisSafe"`
}

type AuthNonce struct {
	Nonce      *string `json:"nonce"`
	UserExists *bool   `json:"userExists"`
}

func (AuthNonce) IsGetAuthNoncePayloadOrError() {}

type Collection struct {
	Dbid           persist.DBID      `json:"dbid"`
	Version        *int              `json:"version"`
	Name           *string           `json:"name"`
	CollectorsNote *string           `json:"collectorsNote"`
	Gallery        *Gallery          `json:"gallery"`
	Layout         *CollectionLayout `json:"layout"`
	Hidden         *bool             `json:"hidden"`
	Nfts           []*CollectionNft  `json:"nfts"`
}

func (Collection) IsNode()                  {}
func (Collection) IsCollectionByIDOrError() {}

type CollectionLayout struct {
	Columns    *int   `json:"columns"`
	Whitespace []*int `json:"whitespace"`
}

type CollectionLayoutInput struct {
	Columns    int   `json:"columns"`
	Whitespace []int `json:"whitespace"`
}

type CollectionNft struct {
	HelperCollectionNftData
	Nft        *Nft        `json:"nft"`
	Collection *Collection `json:"collection"`
}

func (CollectionNft) IsNode()                     {}
func (CollectionNft) IsCollectionNftByIDOrError() {}

type Community struct {
	LastUpdated     *time.Time        `json:"lastUpdated"`
	ContractAddress *persist.Address  `json:"contractAddress"`
	CreatorAddress  *persist.Address  `json:"creatorAddress"`
	Name            *string           `json:"name"`
	Description     *string           `json:"description"`
	PreviewImage    *string           `json:"previewImage"`
	Owners          []*CommunityOwner `json:"owners"`
}

func (Community) IsNode()                      {}
func (Community) IsCommunityByAddressOrError() {}

type CommunityOwner struct {
	Address  *persist.Address `json:"address"`
	Username *string          `json:"username"`
}

type CreateCollectionInput struct {
	GalleryID      persist.DBID           `json:"galleryId"`
	Name           string                 `json:"name"`
	CollectorsNote string                 `json:"collectorsNote"`
	Nfts           []persist.DBID         `json:"nfts"`
	Layout         *CollectionLayoutInput `json:"layout"`
}

type CreateCollectionPayload struct {
	Collection *Collection `json:"collection"`
}

func (CreateCollectionPayload) IsCreateCollectionPayloadOrError() {}

type CreateUserPayload struct {
	UserID    *persist.DBID `json:"userId"`
	GalleryID *persist.DBID `json:"galleryId"`
	Viewer    *Viewer       `json:"viewer"`
}

func (CreateUserPayload) IsCreateUserPayloadOrError() {}

type DebugAuth struct {
	UserID    *persist.DBID     `json:"userId"`
	Addresses []persist.Address `json:"addresses"`
}

type DeleteCollectionPayload struct {
	Gallery *Gallery `json:"gallery"`
}

func (DeleteCollectionPayload) IsDeleteCollectionPayloadOrError() {}

type ErrAuthenticationFailed struct {
	Message string `json:"message"`
}

func (ErrAuthenticationFailed) IsAddUserAddressPayloadOrError() {}
func (ErrAuthenticationFailed) IsError()                        {}
func (ErrAuthenticationFailed) IsLoginPayloadOrError()          {}
func (ErrAuthenticationFailed) IsCreateUserPayloadOrError()     {}

type ErrCollectionNotFound struct {
	Message string `json:"message"`
}

func (ErrCollectionNotFound) IsError()                          {}
func (ErrCollectionNotFound) IsCollectionByIDOrError()          {}
func (ErrCollectionNotFound) IsCollectionNftByIDOrError()       {}
func (ErrCollectionNotFound) IsDeleteCollectionPayloadOrError() {}

type ErrCommunityNotFound struct {
	Message string `json:"message"`
}

func (ErrCommunityNotFound) IsCommunityByAddressOrError() {}
func (ErrCommunityNotFound) IsError()                     {}

type ErrDoesNotOwnRequiredNft struct {
	Message string `json:"message"`
}

func (ErrDoesNotOwnRequiredNft) IsGetAuthNoncePayloadOrError() {}
func (ErrDoesNotOwnRequiredNft) IsAuthorizationError()         {}
func (ErrDoesNotOwnRequiredNft) IsError()                      {}
func (ErrDoesNotOwnRequiredNft) IsLoginPayloadOrError()        {}
func (ErrDoesNotOwnRequiredNft) IsCreateUserPayloadOrError()   {}

type ErrInvalidInput struct {
	Message    string   `json:"message"`
	Parameters []string `json:"parameters"`
	Reasons    []string `json:"reasons"`
}

func (ErrInvalidInput) IsUserByUsernameOrError()                  {}
func (ErrInvalidInput) IsCreateCollectionPayloadOrError()         {}
func (ErrInvalidInput) IsDeleteCollectionPayloadOrError()         {}
func (ErrInvalidInput) IsUpdateCollectionInfoPayloadOrError()     {}
func (ErrInvalidInput) IsUpdateCollectionNftsPayloadOrError()     {}
func (ErrInvalidInput) IsUpdateCollectionHiddenPayloadOrError()   {}
func (ErrInvalidInput) IsUpdateGalleryCollectionsPayloadOrError() {}
func (ErrInvalidInput) IsUpdateNftInfoPayloadOrError()            {}
func (ErrInvalidInput) IsAddUserAddressPayloadOrError()           {}
func (ErrInvalidInput) IsRemoveUserAddressesPayloadOrError()      {}
func (ErrInvalidInput) IsUpdateUserInfoPayloadOrError()           {}
func (ErrInvalidInput) IsError()                                  {}

type ErrInvalidToken struct {
	Message string `json:"message"`
}

func (ErrInvalidToken) IsAuthorizationError() {}
func (ErrInvalidToken) IsError()              {}

type ErrNftNotFound struct {
	Message string `json:"message"`
}

func (ErrNftNotFound) IsNftByIDOrError()           {}
func (ErrNftNotFound) IsError()                    {}
func (ErrNftNotFound) IsCollectionNftByIDOrError() {}

type ErrNoCookie struct {
	Message string `json:"message"`
}

func (ErrNoCookie) IsAuthorizationError() {}
func (ErrNoCookie) IsError()              {}

type ErrNotAuthorized struct {
	Message string             `json:"message"`
	Cause   AuthorizationError `json:"cause"`
}

func (ErrNotAuthorized) IsViewerOrError()                          {}
func (ErrNotAuthorized) IsCreateCollectionPayloadOrError()         {}
func (ErrNotAuthorized) IsDeleteCollectionPayloadOrError()         {}
func (ErrNotAuthorized) IsUpdateCollectionInfoPayloadOrError()     {}
func (ErrNotAuthorized) IsUpdateCollectionNftsPayloadOrError()     {}
func (ErrNotAuthorized) IsUpdateCollectionHiddenPayloadOrError()   {}
func (ErrNotAuthorized) IsUpdateGalleryCollectionsPayloadOrError() {}
func (ErrNotAuthorized) IsUpdateNftInfoPayloadOrError()            {}
func (ErrNotAuthorized) IsAddUserAddressPayloadOrError()           {}
func (ErrNotAuthorized) IsRemoveUserAddressesPayloadOrError()      {}
func (ErrNotAuthorized) IsUpdateUserInfoPayloadOrError()           {}
func (ErrNotAuthorized) IsRefreshOpenSeaNftsPayloadOrError()       {}
func (ErrNotAuthorized) IsError()                                  {}

type ErrOpenSeaRefreshFailed struct {
	Message string `json:"message"`
}

func (ErrOpenSeaRefreshFailed) IsRefreshOpenSeaNftsPayloadOrError() {}
func (ErrOpenSeaRefreshFailed) IsError()                            {}

type ErrUserAlreadyExists struct {
	Message string `json:"message"`
}

func (ErrUserAlreadyExists) IsUpdateUserInfoPayloadOrError() {}
func (ErrUserAlreadyExists) IsError()                        {}
func (ErrUserAlreadyExists) IsCreateUserPayloadOrError()     {}

type ErrUserNotFound struct {
	Message string `json:"message"`
}

func (ErrUserNotFound) IsUserByUsernameOrError() {}
func (ErrUserNotFound) IsError()                 {}
func (ErrUserNotFound) IsLoginPayloadOrError()   {}

type EthereumEoaAuth struct {
	Address   persist.Address `json:"address"`
	Nonce     string          `json:"nonce"`
	Signature string          `json:"signature"`
}

type Gallery struct {
	Dbid        persist.DBID  `json:"dbid"`
	Owner       *GalleryUser  `json:"owner"`
	Collections []*Collection `json:"collections"`
}

func (Gallery) IsNode() {}

type GalleryUser struct {
	Dbid                persist.DBID   `json:"dbid"`
	Username            *string        `json:"username"`
	Bio                 *string        `json:"bio"`
	Wallets             []*Wallet      `json:"wallets"`
	Galleries           []*Gallery     `json:"galleries"`
	IsAuthenticatedUser *bool          `json:"isAuthenticatedUser"`
	Followers           []*GalleryUser `json:"followers"`
	Following           []*GalleryUser `json:"following"`
}

func (GalleryUser) IsNode()                  {}
func (GalleryUser) IsGalleryUserOrWallet()   {}
func (GalleryUser) IsUserByUsernameOrError() {}

type GltfMedia struct {
	PreviewURLs      *PreviewURLSet `json:"previewURLs"`
	MediaURL         *string        `json:"mediaURL"`
	MediaType        *string        `json:"mediaType"`
	ContentRenderURL *string        `json:"contentRenderURL"`
}

func (GltfMedia) IsMediaSubtype() {}
func (GltfMedia) IsMedia()        {}

type GnosisSafeAuth struct {
	Address persist.Address `json:"address"`
	Nonce   string          `json:"nonce"`
}

type HTMLMedia struct {
	PreviewURLs      *PreviewURLSet `json:"previewURLs"`
	MediaURL         *string        `json:"mediaURL"`
	MediaType        *string        `json:"mediaType"`
	ContentRenderURL *string        `json:"contentRenderURL"`
}

func (HTMLMedia) IsMediaSubtype() {}
func (HTMLMedia) IsMedia()        {}

type ImageMedia struct {
	PreviewURLs       *PreviewURLSet `json:"previewURLs"`
	MediaURL          *string        `json:"mediaURL"`
	MediaType         *string        `json:"mediaType"`
	ContentRenderURLs *ImageURLSet   `json:"contentRenderURLs"`
}

func (ImageMedia) IsMediaSubtype() {}
func (ImageMedia) IsMedia()        {}

type ImageURLSet struct {
	Raw    *string `json:"raw"`
	Small  *string `json:"small"`
	Medium *string `json:"medium"`
	Large  *string `json:"large"`
}

type InvalidMedia struct {
	PreviewURLs      *PreviewURLSet `json:"previewURLs"`
	MediaURL         *string        `json:"mediaURL"`
	MediaType        *string        `json:"mediaType"`
	ContentRenderURL *string        `json:"contentRenderURL"`
}

func (InvalidMedia) IsMediaSubtype() {}
func (InvalidMedia) IsMedia()        {}

type JSONMedia struct {
	PreviewURLs      *PreviewURLSet `json:"previewURLs"`
	MediaURL         *string        `json:"mediaURL"`
	MediaType        *string        `json:"mediaType"`
	ContentRenderURL *string        `json:"contentRenderURL"`
}

func (JSONMedia) IsMediaSubtype() {}
func (JSONMedia) IsMedia()        {}

type LoginPayload struct {
	UserID *persist.DBID `json:"userId"`
	Viewer *Viewer       `json:"viewer"`
}

func (LoginPayload) IsLoginPayloadOrError() {}

type LogoutPayload struct {
	Viewer *Viewer `json:"viewer"`
}

type MembershipOwner struct {
	Dbid        persist.DBID     `json:"dbid"`
	Address     *persist.Address `json:"address"`
	User        *GalleryUser     `json:"user"`
	PreviewNfts []*string        `json:"previewNfts"`
}

type MembershipTier struct {
	Dbid     persist.DBID       `json:"dbid"`
	Name     *string            `json:"name"`
	AssetURL *string            `json:"assetUrl"`
	TokenID  *string            `json:"tokenId"`
	Owners   []*MembershipOwner `json:"owners"`
}

func (MembershipTier) IsNode() {}

type Nft struct {
	Dbid                  persist.DBID        `json:"dbid"`
	CreationTime          *time.Time          `json:"creationTime"`
	LastUpdated           *time.Time          `json:"lastUpdated"`
	CollectorsNote        *string             `json:"collectorsNote"`
	Media                 MediaSubtype        `json:"media"`
	TokenType             *TokenType          `json:"tokenType"`
	Chain                 *Chain              `json:"chain"`
	Name                  *string             `json:"name"`
	Description           *string             `json:"description"`
	TokenURI              *string             `json:"tokenUri"`
	TokenID               *string             `json:"tokenId"`
	Quantity              *string             `json:"quantity"`
	Owner                 GalleryUserOrWallet `json:"owner"`
	OwnershipHistory      []*OwnerAtBlock     `json:"ownershipHistory"`
	TokenMetadata         *string             `json:"tokenMetadata"`
	ContractAddress       *persist.Address    `json:"contractAddress"`
	ExternalURL           *string             `json:"externalUrl"`
	BlockNumber           *string             `json:"blockNumber"`
	CreatorAddress        *persist.Address    `json:"creatorAddress"`
	OpenseaCollectionName *string             `json:"openseaCollectionName"`
	OpenseaID             *int                `json:"openseaId"`
}

func (Nft) IsNode()           {}
func (Nft) IsNftByIDOrError() {}

type OwnerAtBlock struct {
	Owner       GalleryUserOrWallet `json:"owner"`
	BlockNumber *string             `json:"blockNumber"`
}

type PreviewURLSet struct {
	Raw    *string `json:"raw"`
	Small  *string `json:"small"`
	Medium *string `json:"medium"`
	Large  *string `json:"large"`
}

type RefreshOpenSeaNftsPayload struct {
	Viewer *Viewer `json:"viewer"`
}

func (RefreshOpenSeaNftsPayload) IsRefreshOpenSeaNftsPayloadOrError() {}

type RemoveUserAddressesPayload struct {
	Viewer *Viewer `json:"viewer"`
}

func (RemoveUserAddressesPayload) IsRemoveUserAddressesPayloadOrError() {}

type TextMedia struct {
	PreviewURLs      *PreviewURLSet `json:"previewURLs"`
	MediaURL         *string        `json:"mediaURL"`
	MediaType        *string        `json:"mediaType"`
	ContentRenderURL *string        `json:"contentRenderURL"`
}

func (TextMedia) IsMediaSubtype() {}
func (TextMedia) IsMedia()        {}

type UnknownMedia struct {
	PreviewURLs      *PreviewURLSet `json:"previewURLs"`
	MediaURL         *string        `json:"mediaURL"`
	MediaType        *string        `json:"mediaType"`
	ContentRenderURL *string        `json:"contentRenderURL"`
}

func (UnknownMedia) IsMediaSubtype() {}
func (UnknownMedia) IsMedia()        {}

type UpdateCollectionHiddenInput struct {
	CollectionID persist.DBID `json:"collectionId"`
	Hidden       bool         `json:"hidden"`
}

type UpdateCollectionHiddenPayload struct {
	Collection *Collection `json:"collection"`
}

func (UpdateCollectionHiddenPayload) IsUpdateCollectionHiddenPayloadOrError() {}

type UpdateCollectionInfoInput struct {
	CollectionID   persist.DBID `json:"collectionId"`
	Name           string       `json:"name"`
	CollectorsNote string       `json:"collectorsNote"`
}

type UpdateCollectionInfoPayload struct {
	Collection *Collection `json:"collection"`
}

func (UpdateCollectionInfoPayload) IsUpdateCollectionInfoPayloadOrError() {}

type UpdateCollectionNftsInput struct {
	CollectionID persist.DBID           `json:"collectionId"`
	Nfts         []persist.DBID         `json:"nfts"`
	Layout       *CollectionLayoutInput `json:"layout"`
}

type UpdateCollectionNftsPayload struct {
	Collection *Collection `json:"collection"`
}

func (UpdateCollectionNftsPayload) IsUpdateCollectionNftsPayloadOrError() {}

type UpdateGalleryCollectionsInput struct {
	GalleryID   persist.DBID   `json:"galleryId"`
	Collections []persist.DBID `json:"collections"`
}

type UpdateGalleryCollectionsPayload struct {
	Gallery *Gallery `json:"gallery"`
}

func (UpdateGalleryCollectionsPayload) IsUpdateGalleryCollectionsPayloadOrError() {}

type UpdateNftInfoInput struct {
	NftID          persist.DBID  `json:"nftId"`
	CollectorsNote string        `json:"collectorsNote"`
	CollectionID   *persist.DBID `json:"collectionId"`
}

type UpdateNftInfoPayload struct {
	Nft *Nft `json:"nft"`
}

func (UpdateNftInfoPayload) IsUpdateNftInfoPayloadOrError() {}

type UpdateUserInfoInput struct {
	Username string `json:"username"`
	Bio      string `json:"bio"`
}

type UpdateUserInfoPayload struct {
	Viewer *Viewer `json:"viewer"`
}

func (UpdateUserInfoPayload) IsUpdateUserInfoPayloadOrError() {}

type VideoMedia struct {
	PreviewURLs       *PreviewURLSet `json:"previewURLs"`
	MediaURL          *string        `json:"mediaURL"`
	MediaType         *string        `json:"mediaType"`
	ContentRenderURLs *VideoURLSet   `json:"contentRenderURLs"`
}

func (VideoMedia) IsMediaSubtype() {}
func (VideoMedia) IsMedia()        {}

type VideoURLSet struct {
	Raw    *string `json:"raw"`
	Small  *string `json:"small"`
	Medium *string `json:"medium"`
	Large  *string `json:"large"`
}

type Viewer struct {
	User            *GalleryUser     `json:"user"`
	ViewerGalleries []*ViewerGallery `json:"viewerGalleries"`
}

func (Viewer) IsViewerOrError() {}

type ViewerGallery struct {
	Gallery *Gallery `json:"gallery"`
}

type Wallet struct {
	Address *persist.Address `json:"address"`
	Nfts    []*Nft           `json:"nfts"`
}

func (Wallet) IsNode()                {}
func (Wallet) IsGalleryUserOrWallet() {}

type Chain string

const (
	ChainEthereum Chain = "Ethereum"
	ChainArbitrum Chain = "Arbitrum"
	ChainPolygon  Chain = "Polygon"
	ChainOptimism Chain = "Optimism"
)

var AllChain = []Chain{
	ChainEthereum,
	ChainArbitrum,
	ChainPolygon,
	ChainOptimism,
}

func (e Chain) IsValid() bool {
	switch e {
	case ChainEthereum, ChainArbitrum, ChainPolygon, ChainOptimism:
		return true
	}
	return false
}

func (e Chain) String() string {
	return string(e)
}

func (e *Chain) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Chain(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Chain", str)
	}
	return nil
}

func (e Chain) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type TokenType string

const (
	TokenTypeErc721  TokenType = "ERC721"
	TokenTypeErc1155 TokenType = "ERC1155"
	TokenTypeErc20   TokenType = "ERC20"
)

var AllTokenType = []TokenType{
	TokenTypeErc721,
	TokenTypeErc1155,
	TokenTypeErc20,
}

func (e TokenType) IsValid() bool {
	switch e {
	case TokenTypeErc721, TokenTypeErc1155, TokenTypeErc20:
		return true
	}
	return false
}

func (e TokenType) String() string {
	return string(e)
}

func (e *TokenType) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = TokenType(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid TokenType", str)
	}
	return nil
}

func (e TokenType) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
