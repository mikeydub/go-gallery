package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mikeydub/go-gallery/service/persist"
	"github.com/mikeydub/go-gallery/util"
)

var errGetNFTsInput = errors.New("address or user_id must be provided")

type getNFTsInput struct {
	Address persist.Address `form:"address"`
	UserID  persist.DBID    `form:"user_id"`
}

func getNFTs(nftRepo persist.NFTRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input getNFTsInput
		if err := c.ShouldBindQuery(&input); err != nil {
			util.ErrResponse(c, http.StatusBadRequest, err)
			return
		}

		if input.Address == "" && input.UserID == "" {
			util.ErrResponse(c, http.StatusBadRequest, errGetNFTsInput)
			return
		}

		var nfts []persist.NFT
		var err error

		if input.Address == "" {
			nfts, err = nftRepo.GetByUserID(c, input.UserID)
		} else {
			nfts, err = nftRepo.GetByAddresses(c, []persist.Address{input.Address})
		}
		if err != nil {
			util.ErrResponse(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, nfts)
	}
}
