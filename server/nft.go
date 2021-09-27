package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikeydub/go-gallery/copy"
	"github.com/mikeydub/go-gallery/persist"
	"github.com/mikeydub/go-gallery/runtime"
	"github.com/mikeydub/go-gallery/util"
)

type getNftsByIDInput struct {
	NftID persist.DBID `json:"id" form:"id" binding:"required"`
}

type getNftsByUserIDInput struct {
	UserID persist.DBID `json:"user_id" form:"user_id" binding:"required"`
}

type getUnassignedNFTByUserIDInput struct {
	SkipCache bool `json:"skip_cache" form:"skip_cache"`
}

type getOpenseaNftsInput struct {
	// Comma separated list of wallet addresses
	WalletAddresses string `json:"addresses" form:"addresses"`
	SkipCache       bool   `json:"skip_cache" form:"skip_cache"`
}

type getNftsOutput struct {
	Nfts []*persist.Nft `json:"nfts"`
}

type getNftByIDOutput struct {
	Nft *persist.Nft `json:"nft"`
}

type getUnassignedNftsOutput struct {
	Nfts []*persist.CollectionNft `json:"nfts"`
}

type updateNftByIDInput struct {
	ID             persist.DBID `form:"id" binding:"required"`
	CollectorsNote string       `form:"collectors_note"`
}

type getOwnershipHistoryInput struct {
	NftID     persist.DBID `json:"id" form:"id" binding:"required"`
	SkipCache bool         `json:"skip_cache" form:"skip_cache"`
}

type getOwnershipHistoryOutput struct {
	OwnershipHistory *persist.OwnershipHistory `json:"ownership_history"`
}

func getNftByID(pRuntime *runtime.Runtime) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := &getNftsByIDInput{}

		if err := c.ShouldBindQuery(input); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse{
				Error: copy.NftIDQueryNotProvided,
			})
			return
		}

		nfts, err := persist.NftGetByID(c, input.NftID, pRuntime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
			return
		}
		if len(nfts) == 0 {
			c.JSON(http.StatusNotFound, errorResponse{
				Error: fmt.Sprintf("no nfts found with id: %s", input.NftID),
			})
			return
		}

		if len(nfts) > 1 {
			nfts = nfts[:1]
			// TODO log that this should not be happening
		}
		c.JSON(http.StatusOK, getNftByIDOutput{Nft: nfts[0]})
	}
}

// Must specify nft id in json input
func updateNftByID(pRuntime *runtime.Runtime) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := &updateNftByIDInput{}
		if err := c.ShouldBindJSON(input); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse{
				Error: err.Error(),
			})
			return
		}

		userID, ok := getUserIDfromCtx(c)
		if !ok {
			c.JSON(http.StatusBadRequest, errorResponse{Error: "user id not found in context"})
			return
		}

		update := &persist.UpdateNFTInfoInput{CollectorsNote: input.CollectorsNote}

		err := persist.NftUpdateByID(c, input.ID, userID, update, pRuntime)
		if err != nil {
			if err.Error() == copy.CouldNotFindDocument {
				c.JSON(http.StatusNotFound, errorResponse{Error: err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, successOutput{Success: true})
	}
}

func getNftsForUser(pRuntime *runtime.Runtime) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := &getNftsByUserIDInput{}
		if err := c.ShouldBindQuery(input); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}
		nfts, err := persist.NftGetByUserID(c, input.UserID, pRuntime)
		if len(nfts) == 0 || err != nil {
			nfts = []*persist.Nft{}
		}

		c.JSON(http.StatusOK, getNftsOutput{Nfts: nfts})
	}
}

func getUnassignedNftsForUser(pRuntime *runtime.Runtime) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := &getUnassignedNFTByUserIDInput{}
		if err := c.ShouldBindQuery(input); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}

		userID, ok := getUserIDfromCtx(c)
		if !ok {
			c.JSON(http.StatusBadRequest, errorResponse{Error: "user id not found in context"})
			return
		}
		coll, err := persist.CollGetUnassigned(c, userID, input.SkipCache, pRuntime)
		if coll == nil || err != nil {
			coll = &persist.Collection{Nfts: []*persist.CollectionNft{}}
		}

		c.JSON(http.StatusOK, getUnassignedNftsOutput{Nfts: coll.Nfts})
	}
}

func getNftsFromOpensea(pRuntime *runtime.Runtime) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := &getOpenseaNftsInput{}
		if err := c.ShouldBindQuery(input); err != nil {
			c.JSON(http.StatusBadRequest, errorResponse{Error: err.Error()})
			return
		}

		userID, ok := getUserIDfromCtx(c)
		if !ok {
			c.JSON(http.StatusBadRequest, errorResponse{Error: "user id not found in context"})
			return
		}

		addresses := strings.Split(input.WalletAddresses, ",")
		if len(addresses) > 0 {
			ownsWallet, err := doesUserOwnWallets(c, userID, addresses, pRuntime)
			if err != nil {
				c.JSON(http.StatusInternalServerError, errorResponse{Error: err.Error()})
				return
			}
			if !ownsWallet {
				c.JSON(http.StatusBadRequest, errorResponse{Error: "user does not own wallet"})
				return
			}
		}

		nfts, err := openSeaPipelineAssetsForAcc(c, userID, addresses, input.SkipCache, pRuntime)
		if len(nfts) == 0 || err != nil {
			nfts = []*persist.Nft{}
		}

		c.JSON(http.StatusOK, getNftsOutput{Nfts: nfts})
	}
}

func doesUserOwnWallets(pCtx context.Context, userID persist.DBID, walletAddresses []string, pRuntime *runtime.Runtime) (bool, error) {
	user, err := persist.UserGetByID(pCtx, userID, pRuntime)
	if err != nil {
		return false, err
	}
	for _, walletAddress := range walletAddresses {
		if !util.Contains(user.Addresses, walletAddress) {
			return false, nil
		}
	}
	return true, nil
}
