package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikeydub/go-gallery/eth"
	"github.com/mikeydub/go-gallery/persist"
	"github.com/mikeydub/go-gallery/util"
)

const (
	analyticsKeyUserUpdateWithBio = "updated_bio"
)

var bannedUsernames = map[string]bool{
	"password":      true,
	"auth":          true,
	"welcome":       true,
	"edit":          true,
	"404":           true,
	"nuke":          true,
	"account":       true,
	"settings":      true,
	"artists":       true,
	"artist":        true,
	"collections":   true,
	"collection":    true,
	"nft":           true,
	"members":       true,
	"nfts":          true,
	"bookmarks":     true,
	"messages":      true,
	"guestbook":     true,
	"notifications": true,
	"explore":       true,
	"analytics":     true,
	"gallery":       true,
	"investors":     true,
	"team":          true,
	"faq":           true,
	"info":          true,
	"about":         true,
	"contact":       true,
	"terms":         true,
	"privacy":       true,
	"help":          true,
	"support":       true,
	"feed":          true,
	"feeds":         true,
	"membership":    true,
}

type userUpdateInput struct {
	UserName string `json:"username" binding:"username"`
	BioStr   string `json:"bio"`
}

type userGetInput struct {
	UserID   persist.DBID    `json:"user_id" form:"user_id"`
	Address  persist.Address `json:"address" form:"address" binding:"eth_addr"` // len=42"` // standard ETH "0x"-prefixed address
	Username string          `json:"username" form:"username"`
}

type userGetOutput struct {
	UserID      persist.DBID         `json:"id"`
	UserNameStr string               `json:"username"`
	BioStr      string               `json:"bio"`
	Addresses   []persist.Address    `json:"addresses"`
	CreatedAt   persist.CreationTime `json:"created_at"`
}

// INPUT - USER_CREATE - initial user creation is just an empty user, to store it in the DB.
//         this is to allow for users interupting the onboarding flow, and to be able to come back to it later
//         and the system recognize that their user already exists.
//         the users entering details on the user as they onboard are all user-update operations.
type userAddAddressInput struct {

	// needed because this is a new user that cant be logged into, and the client creating
	// the user still needs to prove ownership of their address.
	Signature string          `json:"signature" binding:"required,signature"`
	Address   persist.Address `json:"address"   binding:"required,eth_addr"` // len=42"` // standard ETH "0x"-prefixed address
}

type userRemoveAddressesInput struct {
	Addresses []persist.Address `json:"addresses"   binding:"required"`
}

type userCreateOutput struct {
	SignatureValid bool         `json:"signature_valid"`
	JWTtoken       string       `json:"jwt_token"` // JWT token is sent back to user to use to continue onboarding
	UserID         persist.DBID `json:"user_id"`
	GalleryID      persist.DBID `json:"gallery_id"`
}

type userAddAddressOutput struct {
	SignatureValid bool `json:"signature_valid"`
}

type errUserNotFound struct {
	userID   persist.DBID
	address  persist.Address
	username string
}

type errNonceNotFound struct {
	address persist.Address
}
type errUserExistsWithAddress struct {
	address persist.Address
}

var errUserIDNotInCtx = errors.New("expected user ID to be in request context")
var errMustResolveENS = errors.New("one of user's addresses must resolve to ENS to set ENS as username")
var errUserCannotRemoveAllAddresses = errors.New("user does not have enough addresses to remove")

func updateUserInfo(userRepository persist.UserRepository, ethClient *eth.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		input := &userUpdateInput{}

		if err := c.ShouldBindJSON(input); err != nil {
			util.ErrResponse(c, http.StatusBadRequest, err)
			return
		}

		userID := getUserIDfromCtx(c)
		if userID == "" {
			util.ErrResponse(c, http.StatusBadRequest, errUserIDNotInCtx)
			return
		}

		if strings.HasSuffix(strings.ToLower(input.UserName), ".eth") {
			user, err := userRepository.GetByID(c, userID)
			if err != nil {
				util.ErrResponse(c, http.StatusInternalServerError, err)
				return
			}
			can := false
			for _, addr := range user.Addresses {
				if resolves, _ := ethClient.ResolvesENS(c, input.UserName, addr); resolves {
					can = true
					break
				}
			}
			if !can {
				c.JSON(http.StatusBadRequest, util.ErrorResponse{Error: errMustResolveENS.Error()})
				return
			}
		}

		err := userRepository.UpdateByID(
			c,
			userID,
			&persist.UserUpdateInfoInput{
				UserNameIdempotent: strings.ToLower(input.UserName),
				UserName:           input.UserName,
				Bio:                sanitizationPolicy.Sanitize(input.BioStr),
			},
		)
		if err != nil {
			util.ErrResponse(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, util.SuccessResponse{Success: true})

	}
}

func getUser(userRepository persist.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		input := &userGetInput{}

		if err := c.ShouldBindQuery(input); err != nil {
			util.ErrResponse(c, http.StatusBadRequest, err)
			return
		}

		output, err := userGetDb(
			c,
			input,
			userRepository,
		)
		if err != nil {
			status := http.StatusInternalServerError
			switch err.(type) {
			case persist.ErrUserNotFoundByAddress, persist.ErrUserNotFoundByID, persist.ErrUserNotFoundByUsername:
				status = http.StatusNotFound
			}
			util.ErrResponse(c, status, err)
			return
		}

		c.JSON(http.StatusOK, output)

	}
}

func createUserToken(userRepository persist.UserRepository, nonceRepository persist.NonceRepository, galleryRepository persist.GalleryTokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		input := &userAddAddressInput{}

		if err := c.ShouldBindJSON(input); err != nil {
			util.ErrResponse(c, http.StatusBadRequest, err)
			return
		}

		output, err := userCreateDbToken(c, input, userRepository, nonceRepository, galleryRepository)
		if err != nil {
			util.ErrResponse(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, output)

	}
}
func addUserAddress(userRepository persist.UserRepository, nonceRepository persist.NonceRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		input := &userAddAddressInput{}

		if err := c.ShouldBindJSON(input); err != nil {
			util.ErrResponse(c, http.StatusBadRequest, err)
			return
		}

		userID := getUserIDfromCtx(c)
		if userID == "" {
			util.ErrResponse(c, http.StatusBadRequest, errUserIDNotInCtx)
			return
		}

		output, err := addAddressToUserDB(c, userID, input, userRepository, nonceRepository)
		if err != nil {
			util.ErrResponse(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, output)

	}
}

func removeAddressesToken(userRepository persist.UserRepository, collRepo persist.CollectionTokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {

		input := &userRemoveAddressesInput{}

		if err := c.ShouldBindJSON(input); err != nil {
			util.ErrResponse(c, http.StatusBadRequest, err)
			return
		}

		userID := getUserIDfromCtx(c)
		if userID == "" {
			util.ErrResponse(c, http.StatusBadRequest, errUserIDNotInCtx)
			return
		}

		err := removeAddressesFromUserDBToken(c, userID, input, userRepository, collRepo)
		if err != nil {
			util.ErrResponse(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, util.SuccessResponse{Success: true})

	}
}

func userCreateDbToken(pCtx context.Context, pInput *userAddAddressInput,
	userRepo persist.UserRepository, nonceRepo persist.NonceRepository, galleryRepo persist.GalleryTokenRepository) (*userCreateOutput, error) {

	output := &userCreateOutput{}

	nonceValueStr, id, _ := getUserWithNonce(pCtx, pInput.Address, userRepo, nonceRepo)
	if nonceValueStr == "" {
		return nil, errNonceNotFound{address: pInput.Address}
	}
	if id != "" {
		return nil, errUserExistsWithAddress{address: pInput.Address}
	}

	sigValidBool, err := authVerifySignatureAllMethods(pInput.Signature,
		nonceValueStr,
		pInput.Address)
	if err != nil {
		return nil, err
	}

	output.SignatureValid = sigValidBool
	if !sigValidBool {
		return output, nil
	}

	user := &persist.User{
		Addresses: []persist.Address{pInput.Address},
	}

	userID, err := userRepo.Create(pCtx, user)
	if err != nil {
		return nil, err
	}

	output.UserID = userID

	jwtTokenStr, err := jwtGeneratePipeline(pCtx, userID)
	if err != nil {
		return nil, err
	}

	output.JWTtoken = jwtTokenStr

	err = authNonceRotateDb(pCtx, pInput.Address, userID, nonceRepo)
	if err != nil {
		return nil, err
	}

	galleryInsert := &persist.GalleryTokenDB{OwnerUserID: userID, Collections: []persist.DBID{}}

	galleryID, err := galleryRepo.Create(pCtx, galleryInsert)
	if err != nil {
		return nil, err
	}

	output.GalleryID = galleryID

	return output, nil
}

func addAddressToUserDB(pCtx context.Context, pUserID persist.DBID, pInput *userAddAddressInput,
	userRepo persist.UserRepository, nonceRepo persist.NonceRepository) (*userAddAddressOutput, error) {

	output := &userAddAddressOutput{}

	nonceValueStr, userID, _ := getUserWithNonce(pCtx, pInput.Address, userRepo, nonceRepo)
	if nonceValueStr == "" {
		return nil, errNonceNotFound{pInput.Address}
	}
	if userID != "" {
		return nil, errUserExistsWithAddress{pInput.Address}
	}

	dataStr := nonceValueStr
	sigValidBool, err := authVerifySignatureAllMethods(pInput.Signature,
		dataStr,
		pInput.Address)
	if err != nil {
		return nil, err
	}

	output.SignatureValid = sigValidBool
	if !sigValidBool {
		return output, nil
	}

	if err = userRepo.AddAddresses(pCtx, pUserID, []persist.Address{pInput.Address}); err != nil {
		return nil, err
	}

	err = authNonceRotateDb(pCtx, pInput.Address, pUserID, nonceRepo)
	if err != nil {
		return nil, err
	}

	return output, nil
}
func removeAddressesFromUserDBToken(pCtx context.Context, pUserID persist.DBID, pInput *userRemoveAddressesInput,
	userRepo persist.UserRepository, collRepo persist.CollectionTokenRepository) error {

	user, err := userRepo.GetByID(pCtx, pUserID)
	if err != nil {
		return err
	}

	if len(user.Addresses) < len(pInput.Addresses) {
		return errUserCannotRemoveAllAddresses
	}

	err = userRepo.RemoveAddresses(pCtx, pUserID, pInput.Addresses)
	if err != nil {
		return err
	}
	return collRepo.RemoveNFTsOfAddresses(pCtx, pUserID, pInput.Addresses)
}

func userGetDb(pCtx context.Context, pInput *userGetInput,
	userRepo persist.UserRepository) (*userGetOutput, error) {

	//------------------

	var user *persist.User
	var err error
	switch {
	case pInput.UserID != "":
		user, err = userRepo.GetByID(pCtx, pInput.UserID)
		if err != nil {
			return nil, err
		}
		break
	case pInput.Username != "":
		user, err = userRepo.GetByUsername(pCtx, pInput.Username)
		if err != nil {
			return nil, err
		}
		break
	case pInput.Address != "":
		user, err = userRepo.GetByAddress(pCtx, pInput.Address)
		if err != nil {
			return nil, err
		}
		break
	}

	if user == nil {
		return nil, errUserNotFound{pInput.UserID, pInput.Address, pInput.Username}
	}

	output := &userGetOutput{
		UserID:      user.ID,
		UserNameStr: user.UserName,
		BioStr:      user.Bio,
		CreatedAt:   user.CreationTime,
		Addresses:   user.Addresses,
	}

	return output, nil
}

// returns nonce value string, user id
// will return empty strings and error if no nonce found
// will return empty string if no user found
func getUserWithNonce(pCtx context.Context, pAddress persist.Address,
	userRepo persist.UserRepository, nonceRepo persist.NonceRepository) (nonceValue string, userID persist.DBID, err error) {

	nonce, err := nonceRepo.Get(pCtx, pAddress)
	if err != nil {
		return nonceValue, userID, err
	}
	if nonce != nil {
		nonceValue = nonce.Value
	} else {
		return nonceValue, userID, errNonceNotFound{pAddress}
	}

	user, err := userRepo.GetByAddress(pCtx, pAddress)
	if err != nil {
		return nonceValue, userID, err
	}
	if user != nil {
		userID = user.ID
	} else {
		return nonceValue, userID, errUserNotFound{address: pAddress}
	}

	return nonceValue, userID, nil
}

func (e errUserNotFound) Error() string {
	return fmt.Sprintf("user not found: address: %s, ID: %s, Username: %s", e.address, e.userID, e.username)
}

func (e errNonceNotFound) Error() string {
	return fmt.Sprintf("nonce not found for address: %s", e.address)
}

func (e errUserExistsWithAddress) Error() string {
	return fmt.Sprintf("user already exists with address: %s", e.address)
}
