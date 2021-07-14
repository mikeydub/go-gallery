package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mikeydub/go-gallery/persist"
	"github.com/mikeydub/go-gallery/runtime"
)

type collectionCreateInput struct {
	OwnerUserIdStr    string `json:"user_id" validate:"required"`
	NameStr           string `json:"name"        validate:"required,min=4,max=50"`
	CollectorsNoteStr string `json:"collectors_note" validate:"required,min=0,max=500"`
}

type collectionDeleteInput struct {
	IDstr string `json:"id"`
}

//-------------------------------------------------------------
// HANDLERS

func getAllCollectionsForUser(pRuntime *runtime.Runtime) gin.HandlerFunc {
	return func(c *gin.Context) {
		//------------------
		// INPUT

		userIDstr := c.Query("userid")

		//------------------
		// CREATE

		colls, err := persist.CollGetByUserID(persist.DbId(userIDstr), c, pRuntime)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"collections": colls})
	}
}

func createCollection(pRuntime *runtime.Runtime) gin.HandlerFunc {
	return func(c *gin.Context) {

		// TODO sanatize input
		input := &collectionCreateInput{}
		if err := c.ShouldBindJSON(input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//------------------
		// CREATE

		_, err := collectionCreateDb(input, input.OwnerUserIdStr, c, pRuntime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.Status(http.StatusOK)
	}
}

func deleteCollection(pRuntime *runtime.Runtime) gin.HandlerFunc {
	return func(c *gin.Context) {
		input := collectionCreateInput{}
		if err := c.ShouldBindJSON(input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// TODO make a db func for delete

		c.Status(http.StatusOK)
	}
}

//-------------------------------------------------------------
// CREATE
func collectionCreateDb(pInput *collectionCreateInput,
	pUserIDstr string,
	pCtx context.Context,
	pRuntime *runtime.Runtime) (persist.DbId, error) {

	err := runtime.Validate(pInput, pRuntime)
	if err != nil {
		return "", err
	}

	//------------------

	nameStr := pInput.NameStr
	ownerUserIDstr := pUserIDstr

	coll := &persist.CollectionDb{
		NameStr:           nameStr,
		CollectorsNoteStr: pInput.CollectorsNoteStr,
		OwnerUserIDstr:    ownerUserIDstr,
		DeletedBool:       false,
		NFTsLst:           []persist.DbId{},
	}

	return persist.CollCreate(coll, pCtx, pRuntime)

}

//-------------------------------------------------------------
