package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/mikeydub/go-gallery/persist"
	"github.com/mikeydub/go-gallery/util"
	"github.com/stretchr/testify/assert"
)

func TestUpdateCollectionNameByID_Success(t *testing.T) {
	assert := setupTest(t, 2)

	// seed DB with collection
	collID, err := tc.repos.collectionTokenRepository.Create(context.Background(), &persist.CollectionTokenDB{
		Name:        "very cool collection",
		OwnerUserID: tc.user1.id,
	})
	assert.Nil(err)

	// build update request body
	update := collectionUpdateInfoByIDInput{Name: "new coll name", ID: collID}
	resp := updateCollectionInfoRequest(assert, update, tc.user1.jwt)
	assertValidResponse(assert, resp)

	// retrieve updated nft
	resp, err = http.Get(fmt.Sprintf("%s/collections/get?id=%s", tc.serverURL, collID))
	assert.Nil(err)
	assertValidJSONResponse(assert, resp)

	type CollectionGetResponse struct {
		Collection *persist.Collection `json:"collection"`
		Error      string              `json:"error"`
	}
	// ensure nft was updated
	body := CollectionGetResponse{}
	util.UnmarshallBody(&body, resp.Body)
	assert.NotNil(body.Collection)
	assert.Empty(body.Error)
	assert.Equal(update.Name, body.Collection.Name)
}

func TestCreateCollection_Success_Token(t *testing.T) {
	assert := setupTest(t, 2)

	nftIDs := seedTokens(assert)

	gid, err := tc.repos.galleryTokenRepository.Create(context.Background(), &persist.GalleryTokenDB{OwnerUserID: tc.user1.id})
	assert.Nil(err)

	input := collectionCreateInputToken{GalleryID: gid, Nfts: nftIDs}
	resp := createCollectionRequestToken(assert, input, tc.user1.jwt)
	assertValidResponse(assert, resp)

	type CreateResp struct {
		ID    persist.DBID `json:"collection_id"`
		Error string       `json:"error"`
	}

	createResp := &CreateResp{}
	err = util.UnmarshallBody(createResp, resp.Body)
	assert.Nil(err)
	assert.Empty(createResp.Error)

	// retrieve updated nft
	resp, err = http.Get(fmt.Sprintf("%s/collections/get?id=%s", tc.serverURL, createResp.ID))
	assert.Nil(err)
	assertValidJSONResponse(assert, resp)

	type CollectionGetResponse struct {
		Collection *persist.Collection `json:"collection"`
		Error      string              `json:"error"`
	}
	// ensure nft was updated
	body := CollectionGetResponse{}
	util.UnmarshallBody(&body, resp.Body)
	assert.NotNil(body.Collection)
	assert.Len(body.Collection.Nfts, 3)
	assert.Empty(body.Error)

	gallery, err := tc.repos.galleryTokenRepository.GetByID(context.Background(), gid, true)
	assert.Nil(err)
	assert.Len(gallery.Collections, 1)
}

func TestGetUnassignedCollection_Success_Token(t *testing.T) {
	assert := setupTest(t, 2)

	nftIDs := seedTokens(assert)
	// seed DB with collection
	_, err := tc.repos.collectionTokenRepository.Create(context.Background(), &persist.CollectionTokenDB{
		Name:        "very cool collection",
		OwnerUserID: tc.user1.id,
		Nfts:        nftIDs[:2],
	})
	assert.Nil(err)

	resp := getUnassignedNFTsRequestToken(assert, tc.user1.id, tc.user1.jwt)
	assertValidResponse(assert, resp)

	type NftsResponse struct {
		Nfts  []*persist.Token `json:"nfts"`
		Error string           `json:"error"`
	}
	// ensure nft was updated
	body := NftsResponse{}
	util.UnmarshallBody(&body, resp.Body)
	assert.Len(body.Nfts, 1)
	assert.Empty(body.Error)
}

func TestDeleteCollection_Success_Token(t *testing.T) {
	assert := setupTest(t, 2)

	collID := createCollectionInDbForUserIDToken(assert, "COLLECTION NAME", tc.user1.id)
	verifyCollectionExistsInDbForIDToken(assert, collID)

	resp := sendCollDeleteRequestToken(assert, collectionDeleteInputToken{ID: collID}, tc.user1)

	assertValidResponse(assert, resp)

	// Assert that the collection was deleted
	_, err := tc.repos.collectionTokenRepository.GetByID(context.Background(), collID, false)
	assert.NotNil(err)

}

func TestDeleteCollection_Failure_Unauthenticated_Token(t *testing.T) {
	assert := setupTest(t, 2)

	collID := createCollectionInDbForUserIDToken(assert, "COLLECTION NAME", tc.user1.id)
	verifyCollectionExistsInDbForIDToken(assert, collID)

	resp := sendCollDeleteRequestToken(assert, collectionDeleteInputToken{ID: collID}, nil)

	assert.Equal(401, resp.StatusCode)
}

func TestDeleteCollection_Failure_DifferentUsersCollection_Token(t *testing.T) {
	assert := setupTest(t, 2)

	collID := createCollectionInDbForUserIDToken(assert, "COLLECTION NAME", tc.user1.id)
	verifyCollectionExistsInDbForIDToken(assert, collID)

	resp := sendCollDeleteRequestToken(assert, collectionDeleteInputToken{ID: collID}, tc.user2)
	assert.Equal(404, resp.StatusCode)
}

func TestGetHiddenCollections_Success_Token(t *testing.T) {
	assert := setupTest(t, 2)

	nftIDs := seedTokens(assert)
	_, err := tc.repos.collectionTokenRepository.Create(context.Background(), &persist.CollectionTokenDB{
		Name:        "very cool collection",
		OwnerUserID: tc.user1.id,
		Nfts:        nftIDs,
		Hidden:      true,
	})
	assert.Nil(err)

	resp := sendCollUserGetRequestToken(assert, string(tc.user1.id), tc.user1)

	type CollectionsResponse struct {
		Collections []*persist.Collection `json:"collections"`
		Error       string                `json:"error"`
	}

	body := CollectionsResponse{}
	util.UnmarshallBody(&body, resp.Body)
	assert.Len(body.Collections, 1)
	assert.Empty(body.Error)
}

func TestGetNoHiddenCollections_Success_Token(t *testing.T) {
	assert := setupTest(t, 2)

	nftIDs := seedTokens(assert)
	_, err := tc.repos.collectionTokenRepository.Create(context.Background(), &persist.CollectionTokenDB{
		Name:        "very cool collection",
		OwnerUserID: tc.user1.id,
		Nfts:        nftIDs[0:1],
		Hidden:      false,
	})
	_, err = tc.repos.collectionTokenRepository.Create(context.Background(), &persist.CollectionTokenDB{
		Name:        "very cool collection",
		OwnerUserID: tc.user1.id,
		Nfts:        nftIDs[1:],
		Hidden:      true,
	})
	assert.Nil(err)

	resp := sendCollUserGetRequestToken(assert, string(tc.user1.id), tc.user2)

	type CollectionsResponse struct {
		Collections []*persist.Collection `json:"collections"`
		Error       string                `json:"error"`
	}

	body := CollectionsResponse{}
	util.UnmarshallBody(&body, resp.Body)
	assert.Len(body.Collections, 1)
	assert.Empty(body.Error)
}

func TestUpdateCollectionNftsOrder_Success_Token(t *testing.T) {
	assert := setupTest(t, 2)

	nftIDs := seedTokens(assert)

	collID, err := tc.repos.collectionTokenRepository.Create(context.Background(), &persist.CollectionTokenDB{
		Name:        "very cool collection",
		OwnerUserID: tc.user1.id,
		Nfts:        nftIDs,
	})
	assert.Nil(err)

	updatedIDs := make([]persist.DBID, 3)
	updatedIDs[0] = nftIDs[0]
	updatedIDs[1] = nftIDs[2]
	updatedIDs[2] = nftIDs[1]

	update := collectionUpdateNftsByIDInputToken{ID: collID, Nfts: updatedIDs}
	resp := updateCollectionNftsRequestToken(assert, update, tc.user1.jwt)
	assertValidResponse(assert, resp)

	errResp := util.ErrorResponse{}
	util.UnmarshallBody(&errResp, resp.Body)
	assert.Empty(errResp.Error)

	// retrieve updated nft
	resp, err = http.Get(fmt.Sprintf("%s/collections/get?id=%s", tc.serverURL, collID))
	assert.Nil(err)
	assertValidJSONResponse(assert, resp)

	type CollectionGetResponse struct {
		Collection persist.Collection `json:"collection"`
		Error      string             `json:"error"`
	}
	// ensure nft was updated
	body := CollectionGetResponse{}
	util.UnmarshallBody(&body, resp.Body)

	assert.Empty(body.Error)
	assert.NotEqual(updatedIDs[1], body.Collection.Nfts[2].ID)
	assert.Equal(updatedIDs[1], body.Collection.Nfts[1].ID)
}
func TestUpdateCollectionNfts_Success_Token(t *testing.T) {
	assert := setupTest(t, 2)

	nftIDs := seedTokens(assert)

	collID, err := tc.repos.collectionTokenRepository.Create(context.Background(), &persist.CollectionTokenDB{
		Name:        "very cool collection",
		OwnerUserID: tc.user1.id,
		Nfts:        nftIDs,
	})
	assert.Nil(err)

	update := collectionUpdateNftsByIDInputToken{ID: collID, Nfts: nftIDs[:2]}
	resp := updateCollectionNftsRequestToken(assert, update, tc.user1.jwt)
	assertValidResponse(assert, resp)

	errResp := util.ErrorResponse{}
	util.UnmarshallBody(&errResp, resp.Body)
	assert.Empty(errResp.Error)

	// retrieve updated nft
	resp, err = http.Get(fmt.Sprintf("%s/collections/get?id=%s", tc.serverURL, collID))
	assert.Nil(err)
	assertValidJSONResponse(assert, resp)

	type CollectionGetResponse struct {
		Collection persist.Collection `json:"collection"`
		Error      string             `json:"error"`
	}
	// ensure nft was updated
	body := CollectionGetResponse{}
	util.UnmarshallBody(&body, resp.Body)

	assert.Empty(body.Error)
	assert.Len(body.Collection.Nfts, 2)
}

func verifyCollectionExistsInDbForIDToken(assert *assert.Assertions, collID persist.DBID) {
	collectionsBeforeDelete, err := tc.repos.collectionTokenRepository.GetByID(context.Background(), collID, false)
	assert.Nil(err)
	assert.Equal(collectionsBeforeDelete.ID, collID)
}

func sendCollDeleteRequestToken(assert *assert.Assertions, requestBody interface{}, authenticatedUser *TestUser) *http.Response {
	data, err := json.Marshal(requestBody)
	assert.Nil(err)

	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/collections/delete", tc.serverURL),
		bytes.NewBuffer(data))
	assert.Nil(err)

	if authenticatedUser != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authenticatedUser.jwt))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(err)

	return resp
}

func getUnassignedNFTsRequestToken(assert *assert.Assertions, userID persist.DBID, jwt string) *http.Response {
	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/nfts/unassigned/get?user_id=%s&skip_cache=false", tc.serverURL, userID),
		nil)
	assert.Nil(err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(err)
	return resp
}

func sendCollUserGetRequestToken(assert *assert.Assertions, forUserID string, authenticatedUser *TestUser) *http.Response {

	req, err := http.NewRequest("GET",
		fmt.Sprintf("%s/collections/user_get?user_id=%s", tc.serverURL, forUserID),
		nil)
	assert.Nil(err)

	if authenticatedUser != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authenticatedUser.jwt))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(err)
	assertValidResponse(assert, resp)

	return resp
}

func createCollectionRequestToken(assert *assert.Assertions, input collectionCreateInputToken, jwt string) *http.Response {
	data, err := json.Marshal(input)
	assert.Nil(err)

	// send update request
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/collections/create", tc.serverURL),
		bytes.NewBuffer(data))
	assert.Nil(err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(err)
	return resp
}

func updateCollectionInfoRequestToken(assert *assert.Assertions, input collectionUpdateInfoByIDInputToken, jwt string) *http.Response {
	data, err := json.Marshal(input)
	assert.Nil(err)

	// send update request
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/collections/update/info", tc.serverURL),
		bytes.NewBuffer(data))
	assert.Nil(err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(err)
	return resp
}
func updateCollectionNftsRequestToken(assert *assert.Assertions, input collectionUpdateNftsByIDInputToken, jwt string) *http.Response {
	data, err := json.Marshal(input)
	assert.Nil(err)

	// send update request
	req, err := http.NewRequest("POST",
		fmt.Sprintf("%s/collections/update/nfts", tc.serverURL),
		bytes.NewBuffer(data))
	assert.Nil(err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt))
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(err)
	return resp
}

func createCollectionInDbForUserIDToken(assert *assert.Assertions, collectionName string, userID persist.DBID) persist.DBID {
	collID, err := tc.repos.collectionTokenRepository.Create(context.Background(), &persist.CollectionTokenDB{
		Name:        collectionName,
		OwnerUserID: userID,
	})
	assert.Nil(err)

	return collID
}

func seedTokens(assert *assert.Assertions) []persist.DBID {
	nfts := []*persist.Token{
		{CollectorsNote: "asd", OwnerAddress: tc.user1.address},
		{CollectorsNote: "bbb", OwnerAddress: tc.user1.address},
		{CollectorsNote: "wowowowow", OwnerAddress: tc.user1.address},
	}
	nftIDs, err := tc.repos.tokenRepository.CreateBulk(context.Background(), nfts)
	assert.Nil(err)
	return nftIDs
}
