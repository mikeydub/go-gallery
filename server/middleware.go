package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mikeydub/go-gallery/copy"
	"github.com/mikeydub/go-gallery/persist"
	"github.com/mikeydub/go-gallery/runtime"
	"github.com/mikeydub/go-gallery/util"
)

const (
	userIDcontextKey = "user_id"
	authContextKey   = "authenticated"
)

var rateLimiter = NewIPRateLimiter(1, 5)

func jwtRequired(runtime *runtime.Runtime) gin.HandlerFunc {
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
			valid, userID, err := authJwtParse(jwt, os.Getenv("JWT_SECRET"), runtime)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, util.ErrorResponse{Error: err.Error()})
				return
			}

			if !valid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid jwt"})
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

func jwtOptional(runtime *runtime.Runtime) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header != "" {
			authHeaders := strings.Split(c.GetHeader("Authorization"), " ")
			if len(authHeaders) == 2 {
				// get string after "Bearer"
				jwt := authHeaders[1]
				valid, userID, _ := authJwtParse(jwt, os.Getenv("JWT_SECRET"), runtime)
				c.Set(authContextKey, valid)
				c.Set(userIDcontextKey, userID)
			}
		}
		c.Next()
	}
}

func rateLimited(runtime *runtime.Runtime) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := rateLimiter.GetLimiter(c.ClientIP())
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusBadRequest, errorResponse{Error: "rate limited"})
			return
		}
		c.Next()
	}
}

func handleCORS(runtimeConfig *runtime.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestOrigin := c.Request.Header.Get("Origin")
		allowedOrigins := strings.Split(runtimeConfig.AllowedOrigins, ",")

		if util.Contains(allowedOrigins, requestOrigin) {
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

func getUserIDfromCtx(c *gin.Context) (persist.DBID, bool) {
	val, ok := c.Get(userIDcontextKey)
	if !ok {
		return "", false
	}
	userID, ok := val.(persist.DBID)
	if !ok {
		return "", false
	}
	return userID, true
}
