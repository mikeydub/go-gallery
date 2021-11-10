package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mikeydub/go-gallery/middleware"
	"github.com/mikeydub/go-gallery/persist"
	"github.com/mikeydub/go-gallery/util"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestConfig struct {
	server    *httptest.Server
	serverURL string
	repos     *repositories
	mgoClient *mongo.Client
	user1     *TestUser
	user2     *TestUser
}

var tc *TestConfig

type TestUser struct {
	id       persist.DBID
	address  persist.Address
	jwt      string
	username string
}

func generateTestUser(repos *repositories, username string) *TestUser {
	ctx := context.Background()

	address := persist.Address(strings.ToLower(fmt.Sprintf("0x%s", util.RandStringBytes(40))))
	user := &persist.User{
		UserName:           username,
		UserNameIdempotent: strings.ToLower(username),
		Addresses:          []persist.Address{address},
	}
	id, err := repos.userRepository.Create(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	jwt, err := middleware.JWTGenerate(ctx, id)
	if err != nil {
		log.Fatal(err)
	}
	authNonceRotateDb(ctx, address, id, repos.nonceRepository)
	log.Info(id, username)
	return &TestUser{id, address, jwt, username}
}

// Should be called at the beginning of every integration test
// Initializes the runtime, connects to mongodb, and starts a test server
func initializeTestEnv() *TestConfig {
	gin.SetMode(gin.ReleaseMode) // Prevent excessive logs
	ts := httptest.NewServer(CoreInit())

	mclient := newMongoClient()
	repos := newRepos()
	log.Info("test server connected! ✅")
	return &TestConfig{
		server:    ts,
		serverURL: fmt.Sprintf("%s/glry/v1", ts.URL),
		repos:     repos,
		mgoClient: mclient,
		user1:     generateTestUser(repos, "bob"),
		user2:     generateTestUser(repos, "john"),
	}
}

// Should be called at the end of every integration test
func teardown() {
	log.Info("tearing down test suite...")
	tc.server.Close()
	clearDB()
}

func clearDB() {
	tc.mgoClient.Database("gallery").Drop(context.Background())
}

func assertValidResponse(assert *assert.Assertions, resp *http.Response) {
	assert.Equal(http.StatusOK, resp.StatusCode, "Status should be 200")
}

func assertValidJSONResponse(assert *assert.Assertions, resp *http.Response) {
	assertValidResponse(assert, resp)
	val, ok := resp.Header["Content-Type"]
	assert.True(ok, "Content-Type header should be set")
	assert.Equal("application/json; charset=utf-8", val[0], "Response should be in JSON")
}

func assertErrorResponse(assert *assert.Assertions, resp *http.Response) {
	assert.NotEqual(http.StatusOK, resp.StatusCode, "Status should not be 200")
	val, ok := resp.Header["Content-Type"]
	assert.True(ok, "Content-Type header should be set")
	assert.Equal("application/json; charset=utf-8", val[0], "Response should be in JSON")
}

func setupTest(t *testing.T) *assert.Assertions {
	tc = initializeTestEnv()
	t.Cleanup(teardown)
	return assert.New(t)
}
