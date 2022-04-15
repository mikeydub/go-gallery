package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"testing"

	"github.com/mikeydub/go-gallery/service/auth"
	"github.com/mikeydub/go-gallery/service/persist"
	"github.com/mikeydub/go-gallery/service/user"
	"github.com/mikeydub/go-gallery/util"
	"github.com/stretchr/testify/assert"
)

func TestUserRoutes(t *testing.T) {

	setupDoubles(t)

	t.Run("get user by ID", func(t *testing.T) {
		setupTest(t, 1)

		assert := assert.New(t)

		resp, err := http.Get(fmt.Sprintf("%s/users/get?user_id=%s", tc.serverURL, tc.user1.id))
		assert.Nil(err)
		assertValidJSONResponse(assert, resp)

		body := persist.User{}
		util.UnmarshallBody(&body, resp.Body)
		assert.Equal(tc.user1.username, body.Username.String())
	})

	t.Run("get user by address", func(t *testing.T) {
		assert := setupTest(t, 1)

		resp, err := http.Get(fmt.Sprintf("%s/users/get?address=%s", tc.serverURL, tc.user2.address))
		assert.Nil(err)
		assertValidJSONResponse(assert, resp)

		body := persist.User{}
		util.UnmarshallBody(&body, resp.Body)
		assert.Equal(tc.user2.username, body.Username.String())
	})

	t.Run("get user by username", func(t *testing.T) {
		assert := setupTest(t, 1)

		resp, err := http.Get(fmt.Sprintf("%s/users/get?username=%s", tc.serverURL, tc.user1.username))
		assert.Nil(err)
		assertValidJSONResponse(assert, resp)

		body := persist.User{}
		util.UnmarshallBody(&body, resp.Body)
		assert.Equal(tc.user1.username, body.Username.String())
	})

	t.Run("get authenticated user includes address", func(t *testing.T) {
		assert := setupTest(t, 1)

		userID := tc.user1.id
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/users/get?user_id=%s", tc.serverURL, userID), nil)
		assert.Nil(err)

		resp, err := tc.user1.client.Do(req)
		assert.Nil(err)
		assertValidJSONResponse(assert, resp)

		body := persist.User{}
		util.UnmarshallBody(&body, resp.Body)
		assert.Equal(userID, body.ID)
		assert.NotEmpty(body.Addresses)
	})

	t.Run("can update user if authenticated", func(t *testing.T) {
		assert := setupTest(t, 1)

		update := user.UpdateUserInput{
			UserName: "kaito",
		}
		resp := updateUserInfoRequest(assert, update, tc.user1)
		assertValidJSONResponse(assert, resp)

		user, err := tc.repos.UserRepository.GetByID(context.Background(), tc.user1.id)
		assert.Nil(err)
		assert.Equal(update.UserName, user.Username.String())
	})

	// Updating the username to itself should not trigger an error, despite the DB
	// having a user entity with that username already
	t.Run("can update user with no change", func(t *testing.T) {
		assert := setupTest(t, 1)

		update := user.UpdateUserInput{
			UserName: "bob",
		}
		resp := updateUserInfoRequest(assert, update, tc.user1)
		assertValidJSONResponse(assert, resp)

		user, err := tc.repos.UserRepository.GetByID(context.Background(), tc.user1.id)
		assert.Nil(err)
		assert.Equal(update.UserName, user.Username.String())
	})

	t.Run("can't update user if not authenticated", func(t *testing.T) {
		assert := setupTest(t, 1)

		update := user.UpdateUserInput{
			UserName: "kaito",
		}
		resp := updateUserInfoNoAuthRequest(assert, update)
		assertErrorResponse(assert, resp)
	})

	t.Run("can't update username to already existing username", func(t *testing.T) {
		assert := setupTest(t, 1)

		user2, err := tc.repos.UserRepository.GetByID(context.Background(), tc.user2.id)
		assert.Nil(err)
		log.Println(user2.Username)

		update := user.UpdateUserInput{
			UserName: tc.user2.username,
		}
		resp := updateUserInfoRequest(assert, update, tc.user1)
		assertErrorResponse(assert, resp)

		user, err := tc.repos.UserRepository.GetByID(context.Background(), tc.user1.id)
		assert.Nil(err)
		assert.NotEqual(update.UserName, user.Username)
	})

	t.Run("can't update username if username is invalid", func(t *testing.T) {
		assert := setupTest(t, 1)

		update := user.UpdateUserInput{
			UserName: "92ks&$m__",
		}
		resp := updateUserInfoRequest(assert, update, tc.user1)
		assertErrorResponse(assert, resp)

		user, err := tc.repos.UserRepository.GetByID(context.Background(), tc.user1.id)
		assert.Nil(err)
		assert.NotEqual(update.UserName, user.Username)
	})

	t.Run("can add addresses", func(t *testing.T) {
		assert := setupTest(t, 1)

		nonce := persist.UserNonce{
			Value:   "TestNonce",
			Address: persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5")),
		}
		err := tc.repos.NonceRepository.Create(context.Background(), nonce)
		assert.Nil(err)

		update := user.AddUserAddressesInput{
			Address:   persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5")),
			Signature: "0x7d3b810c5ae6efa6e5457f5ed85fe048f623b0f1127a7825f119a86714b72fec444d3fa301c05887ba1b94b77e5d68c8567171404cff43b7790e8f4d928b752a1b",
		}
		resp := userAddAddressesRequest(assert, update, tc.user1)
		assertValidJSONResponse(assert, resp)

		errResp := &util.ErrorResponse{}
		err = util.UnmarshallBody(errResp, resp.Body)
		assert.Nil(err)
		assert.Empty(errResp.Error)

		updatedUser, err := tc.repos.UserRepository.GetByID(context.Background(), tc.user1.id)
		assert.Nil(err)
		assert.Equal(update.Address, updatedUser.Addresses[1])
	})

	t.Run("can't add address if wrong nonce", func(t *testing.T) {
		assert := setupTest(t, 1)

		nonce := persist.UserNonce{
			Value:   "Wrong Nonce",
			Address: persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5")),
		}
		err := tc.repos.NonceRepository.Create(context.Background(), nonce)
		assert.Nil(err)

		update := user.AddUserAddressesInput{
			Address:   persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5")),
			Signature: "0x7d3b810c5ae6efa6e5457f5ed85fe048f623b0f1127a7825f119a86714b72fec444d3fa301c05887ba1b94b77e5d68c8567171404cff43b7790e8f4d928b752a1b",
		}
		resp := userAddAddressesRequest(assert, update, tc.user1)
		assertErrorResponse(assert, resp)
	})

	t.Run("can't add address if other user owns address", func(t *testing.T) {
		assert := setupTest(t, 1)

		u := persist.User{
			Addresses: []persist.EthereumAddress{persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5"))},
		}
		_, err := tc.repos.UserRepository.Create(context.Background(), u)

		nonce := persist.UserNonce{
			Value:   "TestNonce",
			Address: persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5")),
		}
		err = tc.repos.NonceRepository.Create(context.Background(), nonce)
		assert.Nil(err)

		update := user.AddUserAddressesInput{
			Address:   persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5")),
			Signature: "0x7d3b810c5ae6efa6e5457f5ed85fe048f623b0f1127a7825f119a86714b72fec444d3fa301c05887ba1b94b77e5d68c8567171404cff43b7790e8f4d928b752a1b",
		}
		resp := userAddAddressesRequest(assert, update, tc.user1)
		assertErrorResponse(assert, resp)
	})

	t.Run("can remove address", func(t *testing.T) {
		assert := setupTest(t, 1)

		u := persist.User{
			Addresses:          []persist.EthereumAddress{"0xcb1b78568d0Ef81585f074b0Dfd6B743959070D9", persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5"))},
			Username:           "TestUser",
			UsernameIdempotent: "testuser",
		}
		userID, err := tc.repos.UserRepository.Create(context.Background(), u)
		assert.Nil(err)

		nft := persist.NFT{
			OwnerAddress: persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5")),
			Name:         "test",
		}
		nftID, err := tc.repos.NftRepository.Create(context.Background(), nft)

		nft2 := persist.NFT{
			OwnerAddress: persist.EthereumAddress(strings.ToLower("0xcb1b78568d0Ef81585f074b0Dfd6B743959070D9")),
			Name:         "blah",
		}
		nftID2, err := tc.repos.NftRepository.Create(context.Background(), nft2)

		coll := persist.CollectionDB{
			NFTs:        []persist.DBID{nftID, nftID2},
			Name:        "test-coll",
			OwnerUserID: userID,
		}
		collID, err := tc.repos.CollectionRepository.Create(context.Background(), coll)

		jwt, err := auth.JWTGeneratePipeline(context.Background(), userID)
		assert.Nil(err)

		update := user.RemoveUserAddressesInput{
			Addresses: []persist.EthereumAddress{persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5"))},
		}
		j, err := cookiejar.New(nil)
		client := &http.Client{Jar: j}
		tu := &TestUser{
			id:     userID,
			jwt:    jwt,
			client: client,
		}
		getFakeCookie(assert, jwt, client)

		resp := userRemoveAddressesRequest(assert, update, tu)
		assertValidJSONResponse(assert, resp)

		errResp := &util.ErrorResponse{}
		util.UnmarshallBody(errResp, resp.Body)
		assert.Empty(errResp.Error)

		nfts, err := tc.repos.NftRepository.GetByUserID(context.Background(), userID)
		assert.Nil(err)
		assert.Len(nfts, 1)

		res, err := tc.repos.CollectionRepository.GetByID(context.Background(), collID, true)
		assert.Nil(err)
		assert.Len(res.NFTs, 2)

		user, err := tc.repos.UserRepository.GetByID(context.Background(), userID)
		assert.Nil(err)
		assert.Len(user.Addresses, 1)
	})

	t.Run("removing address moves NFTS", func(t *testing.T) {
		assert := setupTest(t, 1)

		u := persist.User{
			Addresses:          []persist.EthereumAddress{"0xcb1b78568d0Ef81585f074b0Dfd6B743959070D9", persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5"))},
			Username:           "TestUser",
			UsernameIdempotent: "testuser",
		}
		userID, err := tc.repos.UserRepository.Create(context.Background(), u)
		assert.NoError(err)

		nft := persist.NFT{
			OwnerAddress: persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5")),
			Name:         "test",
			OpenseaID:    2,
		}
		nftID, err := tc.repos.NftRepository.Create(context.Background(), nft)

		nft2 := persist.NFT{
			OwnerAddress: persist.EthereumAddress(strings.ToLower("0xcb1b78568d0Ef81585f074b0Dfd6B743959070D9")),
			Name:         "blah",
			OpenseaID:    1,
		}
		nftID2, err := tc.repos.NftRepository.Create(context.Background(), nft2)

		coll := persist.CollectionDB{
			NFTs:        []persist.DBID{nftID, nftID2},
			Name:        "test-coll",
			OwnerUserID: userID,
		}
		collID, err := tc.repos.CollectionRepository.Create(context.Background(), coll)

		jwt, err := auth.JWTGeneratePipeline(context.Background(), userID)
		assert.NoError(err)

		update := user.RemoveUserAddressesInput{
			Addresses: []persist.EthereumAddress{persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5"))},
		}
		j, err := cookiejar.New(nil)
		client := &http.Client{Jar: j}
		tu := &TestUser{
			id:     userID,
			jwt:    jwt,
			client: client,
		}
		getFakeCookie(assert, jwt, client)

		resp := userRemoveAddressesRequest(assert, update, tu)
		assertValidJSONResponse(assert, resp)

		errResp := &util.ErrorResponse{}
		util.UnmarshallBody(errResp, resp.Body)
		assert.Empty(errResp.Error)

		nfts, err := tc.repos.NftRepository.GetByUserID(context.Background(), userID)
		assert.NoError(err)
		assert.Len(nfts, 1)

		res, err := tc.repos.CollectionRepository.GetByID(context.Background(), collID, true)
		assert.NoError(err)
		assert.Len(res.NFTs, 2)

		user, err := tc.repos.UserRepository.GetByID(context.Background(), userID)
		assert.NoError(err)
		assert.Len(user.Addresses, 1)

		err = tc.repos.NftRepository.UpdateByIDUnsafe(context.Background(), nftID2, persist.NFTUpdateOwnerAddressInput{
			OwnerAddress: persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5")),
		})
		assert.NoError(err)

		newNFT2, err := tc.repos.NftRepository.GetByID(context.Background(), nftID2)
		assert.NoError(err)
		assert.Equal(newNFT2.OwnerAddress, persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5")))

		// This would happen in an opensea refresh
		err = tc.repos.CollectionRepository.RemoveNFTsOfOldAddresses(context.Background(), userID)
		assert.NoError(err)

		res, err = tc.repos.CollectionRepository.GetByID(context.Background(), collID, true)
		assert.NoError(err)
		assert.Len(res.NFTs, 0)
	})

	t.Run("can't remove address if not owned by user", func(t *testing.T) {
		assert := setupTest(t, 1)

		u := persist.User{
			Addresses: []persist.EthereumAddress{persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5")), "0xcb1b78568d0Ef81585f074b0Dfd6B743959070D9"},
		}
		userID, err := tc.repos.UserRepository.Create(context.Background(), u)
		assert.Nil(err)

		jwt, err := auth.JWTGeneratePipeline(context.Background(), userID)
		assert.Nil(err)

		update := user.RemoveUserAddressesInput{
			Addresses: []persist.EthereumAddress{tc.user1.address},
		}

		j, err := cookiejar.New(nil)
		client := &http.Client{Jar: j}
		tu := &TestUser{
			id:     userID,
			jwt:    jwt,
			client: client,
		}
		getFakeCookie(assert, jwt, client)

		resp := userRemoveAddressesRequest(assert, update, tu)
		assertErrorResponse(assert, resp)
	})

	t.Run("can't remove all addresses", func(t *testing.T) {
		assert := setupTest(t, 1)

		u := persist.User{
			Addresses: []persist.EthereumAddress{persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5")), "0xcb1b78568d0Ef81585f074b0Dfd6B743959070D9"},
		}
		userID, err := tc.repos.UserRepository.Create(context.Background(), u)
		assert.Nil(err)

		jwt, err := auth.JWTGeneratePipeline(context.Background(), userID)
		assert.Nil(err)

		update := user.RemoveUserAddressesInput{
			Addresses: u.Addresses,
		}

		j, err := cookiejar.New(nil)
		client := &http.Client{Jar: j}
		tu := &TestUser{
			id:     userID,
			jwt:    jwt,
			client: client,
		}
		getFakeCookie(assert, jwt, client)

		resp := userRemoveAddressesRequest(assert, update, tu)
		assertErrorResponse(assert, resp)
	})

	t.Run("get user previews", func(t *testing.T) {
		assert := setupTest(t, 1)

		nfts := []persist.NFT{
			{
				OwnerAddress:      tc.user1.address,
				Name:              "test",
				OpenseaTokenID:    "10",
				OpenseaID:         2,
				ImageURL:          "https://example.com/test.png",
				ImagePreviewURL:   "https://example.com/test.png",
				ImageThumbnailURL: "https://example.com/test.png",
			},
			{
				OwnerAddress:      tc.user1.address,
				Name:              "test",
				OpenseaTokenID:    "9",
				OpenseaID:         1,
				ImageURL:          "https://example.com/test.png",
				ImagePreviewURL:   "https://example.com/test.png",
				ImageThumbnailURL: "https://example.com/test.png",
			},
			{
				OwnerAddress:      tc.user1.address,
				Name:              "test",
				OpenseaTokenID:    "8",
				OpenseaID:         10,
				ImageURL:          "https://example.com/test.png",
				ImagePreviewURL:   "https://example.com/test.png",
				ImageThumbnailURL: "https://example.com/test.png",
			},
		}
		nftIDs, err := tc.repos.NftRepository.CreateBulk(context.Background(), nfts)
		assert.Nil(err)

		coll := persist.CollectionDB{
			Name:           "test",
			CollectorsNote: "test",
			NFTs:           nftIDs,
			OwnerUserID:    tc.user1.id,
		}
		collID, err := tc.repos.CollectionRepository.Create(context.Background(), coll)
		assert.Nil(err)

		gallery := persist.GalleryDB{
			OwnerUserID: tc.user1.id,
			Collections: []persist.DBID{collID},
		}
		_, err = tc.repos.GalleryRepository.Create(context.Background(), gallery)
		assert.Nil(err)

		resp, err := http.Get(fmt.Sprintf("%s/users/previews?user_id=%s", tc.serverURL, tc.user1.id))
		assert.Nil(err)
		assertValidJSONResponse(assert, resp)

		type previews struct {
			Previews []string `json:"previews"`
		}
		var body previews
		util.UnmarshallBody(&body, resp.Body)
		assert.Len(body.Previews, 3)
	})

	t.Run("can merge users", func(t *testing.T) {
		assert := setupTest(t, 1)

		// Create user
		u := persist.User{
			Addresses: []persist.EthereumAddress{"0xcb1b78568d0Ef81585f074b0Dfd6B743959070D9"},
		}
		userID, err := tc.repos.UserRepository.Create(context.Background(), u)
		assert.Nil(err)

		nft := persist.NFT{
			OwnerAddress:   persist.EthereumAddress(strings.ToLower("0xcb1b78568d0Ef81585f074b0Dfd6B743959070D9")),
			OpenseaTokenID: "1",
		}
		galleryID := fillGallery(assert, nft, userID)

		// Set up cookie to make an authenticated request
		jwt, err := auth.JWTGeneratePipeline(context.Background(), userID)
		assert.Nil(err)
		tu := createCookieAndSetOnUser(assert, userID, jwt)

		u2 := persist.User{
			Addresses: []persist.EthereumAddress{persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5"))},
		}
		userID2, err := tc.repos.UserRepository.Create(context.Background(), u2)
		assert.Nil(err)

		nft2 := persist.NFT{
			OwnerAddress:   u2.Addresses[0],
			OpenseaTokenID: "2",
		}
		deletedGallery := fillGallery(assert, nft2, userID2)

		nonce := persist.UserNonce{
			Value:   "TestNonce",
			Address: u2.Addresses[0],
		}
		err = tc.repos.NonceRepository.Create(context.Background(), nonce)
		assert.Nil(err)

		mergeInput := user.MergeUsersInput{
			SecondUserID: userID2,
			Signature:    "0x7d3b810c5ae6efa6e5457f5ed85fe048f623b0f1127a7825f119a86714b72fec444d3fa301c05887ba1b94b77e5d68c8567171404cff43b7790e8f4d928b752a1b",
			Nonce:        "TestNonce",
			Address:      u2.Addresses[0],
			WalletType:   auth.WalletTypeEOA,
		}

		resp := mergeUserRequest(assert, mergeInput, tu)
		defer resp.Body.Close()
		assertValidResponse(assert, resp)

		errResp := util.ErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		assert.Nil(err)
		assert.Empty(errResp.Error)

		newInitialUser, err := tc.repos.UserRepository.GetByID(context.Background(), userID)
		assert.Nil(err)
		assert.Len(newInitialUser.Addresses, 2)

		newGallery, err := tc.repos.GalleryRepository.GetByID(context.Background(), galleryID)
		assert.Nil(err)
		assert.Len(newGallery.Collections, 2)

		_, err = tc.repos.GalleryRepository.GetByID(context.Background(), deletedGallery)
		assert.NotNil(err)

		_, err = tc.repos.UserRepository.GetByID(context.Background(), userID2)
		assert.NotNil(err)
	})

	t.Run("can merge users with no galleries", func(t *testing.T) {
		assert := setupTest(t, 1)

		// Create user
		u := persist.User{
			Addresses: []persist.EthereumAddress{"0xcb1b78568d0Ef81585f074b0Dfd6B743959070D9"},
		}
		userID, err := tc.repos.UserRepository.Create(context.Background(), u)
		assert.Nil(err)

		// Set up cookie to make an authenticated request
		jwt, err := auth.JWTGeneratePipeline(context.Background(), userID)
		assert.Nil(err)
		tu := createCookieAndSetOnUser(assert, userID, jwt)

		u2 := persist.User{
			Addresses: []persist.EthereumAddress{persist.EthereumAddress(strings.ToLower("0x9a3f9764B21adAF3C6fDf6f947e6D3340a3F8AC5"))},
		}
		userID2, err := tc.repos.UserRepository.Create(context.Background(), u2)
		assert.Nil(err)

		nonce := persist.UserNonce{
			Value:   "TestNonce",
			Address: u2.Addresses[0],
		}
		err = tc.repos.NonceRepository.Create(context.Background(), nonce)
		assert.Nil(err)

		mergeInput := user.MergeUsersInput{
			SecondUserID: userID2,
			Signature:    "0x7d3b810c5ae6efa6e5457f5ed85fe048f623b0f1127a7825f119a86714b72fec444d3fa301c05887ba1b94b77e5d68c8567171404cff43b7790e8f4d928b752a1b",
			Nonce:        "TestNonce",
			Address:      u2.Addresses[0],
			WalletType:   auth.WalletTypeEOA,
		}

		resp := mergeUserRequest(assert, mergeInput, tu)
		defer resp.Body.Close()
		assertValidResponse(assert, resp)

		newInitialUser, err := tc.repos.UserRepository.GetByID(context.Background(), userID)
		assert.Nil(err)
		assert.Len(newInitialUser.Addresses, 2)

		_, err = tc.repos.UserRepository.GetByID(context.Background(), userID2)
		assert.NotNil(err)
	})
}

func fillGallery(assert *assert.Assertions, nft persist.NFT, userID persist.DBID) persist.DBID {
	nftID, err := tc.repos.NftRepository.Create(context.Background(), nft)
	assert.Nil(err)

	coll := persist.CollectionDB{
		OwnerUserID: userID,
		NFTs:        []persist.DBID{nftID},
		Name:        "test-coll",
	}
	collID, err := tc.repos.CollectionRepository.Create(context.Background(), coll)
	assert.Nil(err)

	gallery := persist.GalleryDB{
		OwnerUserID: userID,
		Collections: []persist.DBID{collID},
	}
	galleryID, err := tc.repos.GalleryRepository.Create(context.Background(), gallery)
	assert.Nil(err)
	return galleryID
}

func updateUserInfoRequest(assert *assert.Assertions, input user.UpdateUserInput, tu *TestUser) *http.Response {
	data, err := json.Marshal(input)
	assert.Nil(err)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/update/info", tc.serverURL), bytes.NewBuffer(data))
	assert.Nil(err)

	resp, err := tu.client.Do(req)
	assert.Nil(err)
	return resp
}
func userAddAddressesRequest(assert *assert.Assertions, input user.AddUserAddressesInput, tu *TestUser) *http.Response {
	data, err := json.Marshal(input)
	assert.Nil(err)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/update/addresses/add", tc.serverURL), bytes.NewBuffer(data))
	assert.Nil(err)

	resp, err := tu.client.Do(req)
	assert.Nil(err)
	return resp
}

func userRemoveAddressesRequest(assert *assert.Assertions, input user.RemoveUserAddressesInput, tu *TestUser) *http.Response {
	data, err := json.Marshal(input)
	assert.Nil(err)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/update/addresses/remove", tc.serverURL), bytes.NewBuffer(data))
	assert.Nil(err)

	resp, err := tu.client.Do(req)
	assert.Nil(err)
	return resp
}

func updateUserInfoNoAuthRequest(assert *assert.Assertions, input user.UpdateUserInput) *http.Response {
	data, err := json.Marshal(input)
	assert.Nil(err)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/update/info", tc.serverURL), bytes.NewBuffer(data))
	assert.Nil(err)
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(err)
	return resp
}

func mergeUserRequest(assert *assert.Assertions, input user.MergeUsersInput, tu *TestUser) *http.Response {
	data, err := json.Marshal(input)
	assert.Nil(err)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/users/merge", tc.serverURL), bytes.NewBuffer(data))
	assert.Nil(err)

	resp, err := tu.client.Do(req)
	assert.Nil(err)
	return resp
}
