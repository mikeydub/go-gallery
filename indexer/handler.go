package indexer

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"cloud.google.com/go/storage"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/mikeydub/go-gallery/service/persist"
	"github.com/mikeydub/go-gallery/service/persist/mongodb"
	"github.com/mikeydub/go-gallery/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"google.golang.org/api/option"
)

// Init initializes the indexer
func Init() {
	router, i := coreInit()
	logrus.Info("Starting indexer...")
	go i.Start()
	http.Handle("/", router)
}

func coreInit() (*gin.Engine, *Indexer) {

	setDefaults()

	events := []eventHash{transferBatchEventHash, transferEventHash, transferSingleEventHash}

	tokenRepo, contractRepo, userRepo := newRepos()
	var s *storage.Client
	var err error
	if viper.GetString("ENV") != "local" {
		s, err = storage.NewClient(context.Background())
	} else {
		s, err = storage.NewClient(context.Background(), option.WithCredentialsFile("./decrypted/service-key.json"))
	}
	if err != nil {
		panic(err)
	}
	i := NewIndexer(newEthClient(), newIPFSShell(), s, tokenRepo, contractRepo, userRepo, persist.Chain(viper.GetString("CHAIN")), events, "stats.json")

	router := gin.Default()

	if viper.GetString("ENV") == "local" {
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Info("Registering handlers...")
	return handlersInit(router, i, tokenRepo), i
}

func handlersInit(router *gin.Engine, i *Indexer, tokenRepository persist.TokenRepository) *gin.Engine {
	router.GET("/status", getStatus(i, tokenRepository))

	return router
}

func getStatus(i *Indexer, tokenRepository persist.TokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		total, _ := tokenRepository.Count(ctx, persist.CountTypeTotal)
		mostRecent, _ := tokenRepository.MostRecentBlock(ctx)
		noMetadata, _ := tokenRepository.Count(ctx, persist.CountTypeNoMetadata)
		erc721, _ := tokenRepository.Count(ctx, persist.CountTypeERC721)
		erc1155, _ := tokenRepository.Count(ctx, persist.CountTypeERC1155)

		c.JSON(http.StatusOK, gin.H{
			"total_tokens": total,
			"recent_block": i.mostRecentBlock,
			"most_recent":  mostRecent,
			"bad_uris":     i.badURIs,
			"no_metadata":  noMetadata,
			"erc721":       erc721,
			"erc1155":      erc1155,
		})
	}
}

func setDefaults() {
	viper.SetDefault("RPC_URL", "wss://eth-mainnet.alchemyapi.io/v2/Lxc2B4z57qtwik_KfOS0I476UUUmXT86")
	viper.SetDefault("IPFS_URL", "https://ipfs.io")
	viper.SetDefault("CHAIN", "ETH")
	viper.SetDefault("MONGO_URL", "mongodb://localhost:27017/")
	viper.SetDefault("ENV", "local")
	viper.SetDefault("GCLOUD_TOKEN_LOGS_BUCKET", "eth-token-logs")
	viper.SetDefault("GCLOUD_TOKEN_CONTENT_BUCKET", "token-content")

	viper.AutomaticEnv()
}

func newEthClient() *ethclient.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dialer := *websocket.DefaultDialer
	dialer.ReadBufferSize = 1024 * 20
	rpcClient, err := rpc.DialWebsocketWithDialer(ctx, viper.GetString("RPC_URL"), "", dialer)
	if err != nil {
		panic(err)
	}

	return ethclient.NewClient(rpcClient)

}

func newIPFSShell() *shell.Shell {
	sh := shell.NewShell(viper.GetString("IPFS_URL"))
	sh.SetTimeout(time.Second * 15)
	return sh
}

func newRepos() (persist.TokenRepository, persist.ContractRepository, persist.UserRepository) {
	mgoClient := newMongoClient()
	return mongodb.NewTokenMongoRepository(mgoClient, nil), mongodb.NewContractMongoRepository(mgoClient), mongodb.NewUserMongoRepository(mgoClient)
}

func newMongoClient() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	mgoURL := viper.GetString("MONGO_URL")
	if viper.GetString("ENV") != "local" {
		mongoSecretName := viper.GetString("MONGO_SECRET_NAME")
		secret, err := util.AccessSecret(context.Background(), mongoSecretName)
		if err != nil {
			panic(err)
		}
		mgoURL = string(secret)
	}

	logrus.Infof("Connecting to mongo at %s", mgoURL)

	mOpts := options.Client().ApplyURI(mgoURL)
	mOpts.SetRegistry(mongodb.CustomRegistry)
	mOpts.SetWriteConcern(writeconcern.New(writeconcern.J(true), writeconcern.W(1)))
	mOpts.SetRetryWrites(true)
	mOpts.SetRetryReads(true)

	return mongodb.NewMongoClient(ctx, mOpts)
}

func redirectStderr(f *os.File) {
	err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	if err != nil {
		log.Fatalf("Failed to redirect stderr to file: %v", err)
	}
}
