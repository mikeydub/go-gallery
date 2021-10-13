package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikeydub/go-gallery/copy"
	"github.com/mikeydub/go-gallery/eth"
	"github.com/mikeydub/go-gallery/persist"
	"github.com/mikeydub/go-gallery/util"
	"github.com/spf13/viper"
)

const (
	userIDcontextKey = "user_id"
	authContextKey   = "authenticated"
)

var rateLimiter = NewIPRateLimiter(1, 5)

var errInvalidJWT = errors.New("invalid JWT")

var errRateLimited = errors.New("rate limited")

type errUserDoesNotHaveRequiredNFT struct {
	userID persist.DBID
}

func jwtRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse{Error: copy.InvalidAuthHeader})
			return
		}
		authHeaders := strings.Split(header, " ")
		if len(authHeaders) == 2 {
			if authHeaders[0] != "Bearer" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse{Error: copy.InvalidAuthHeader})
				return
			}
			// get string after "Bearer"
			jwt := authHeaders[1]
			// use an env variable as jwt secret as upposed to using a stateful secret stored in
			// database that is unique to every user and session
			valid, userID, err := authJwtParse(jwt, viper.GetString("JWT_SECRET"))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse{Error: err.Error()})
				return
			}

			if !valid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse{Error: errInvalidJWT.Error()})
				return
			}

			c.Set(userIDcontextKey, userID)
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, util.ErrorResponse{Error: copy.InvalidAuthHeader})
			return
		}
		c.Next()
	}
}

func jwtOptional() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header != "" {
			authHeaders := strings.Split(header, " ")
			if len(authHeaders) == 2 {
				// get string after "Bearer"
				jwt := authHeaders[1]
				valid, userID, _ := authJwtParse(jwt, viper.GetString("JWT_SECRET"))
				c.Set(authContextKey, valid)
				c.Set(userIDcontextKey, userID)
			} else {
				c.Set(authContextKey, false)
				c.Set(userIDcontextKey, persist.DBID(""))
			}
		} else {
			c.Set(authContextKey, false)
			c.Set(userIDcontextKey, persist.DBID(""))
		}
		c.Next()
	}
}

func rateLimited() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := rateLimiter.GetLimiter(c.ClientIP())
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusBadRequest, util.ErrorResponse{Error: errRateLimited.Error()})
			return
		}
		c.Next()
	}
}

func requireNFT(userRepository persist.UserRepository, ethClient *eth.Client, tokenIDs []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserIDfromCtx(c)
		if userID != "" {
			user, err := userRepository.GetByID(c, userID)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, util.ErrorResponse{Error: err.Error()})
				return
			}
			has := false
			for _, addr := range user.Addresses {
				if res, _ := ethClient.HasNFTs(c, tokenIDs, addr); res {
					has = true
					break
				}
			}
			if !has {
				c.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse{Error: errUserDoesNotHaveRequiredNFT{userID}.Error()})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse{Error: errUserIDNotInCtx.Error()})
			return
		}
		c.Next()
	}
}

func handleCORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestOrigin := c.Request.Header.Get("Origin")
		allowedOrigins := strings.Split(viper.GetString("ALLOWED_ORIGINS"), ",")

		if util.Contains(allowedOrigins, requestOrigin) || (strings.ToLower(viper.GetString("ENV")) == "development" && strings.HasPrefix(requestOrigin, "https://gallery-git-") && strings.HasSuffix(requestOrigin, "-gallery-so.vercel.app")) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", requestOrigin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func getUserIDfromCtx(c *gin.Context) persist.DBID {
	return c.MustGet(userIDcontextKey).(persist.DBID)
}

func (e errUserDoesNotHaveRequiredNFT) Error() string {
	return fmt.Sprintf("user %s does not have required NFT", e.userID)
}
