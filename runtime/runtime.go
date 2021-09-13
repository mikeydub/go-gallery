package runtime

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/gin-gonic/gin"
)

const (
	// GalleryDBName represents the database name for gallery related information
	GalleryDBName = "gallery"
	// InfraDBName represents the database name for eth infrastructure related information
	InfraDBName = "infra"
)

// Runtime represents the runtime of the application and its services
type Runtime struct {
	Config *Config
	DB     *DB
	Router *gin.Engine
}

// DB is an abstract represenation of a MongoDB database and Client to interact with it
type DB struct {
	MongoClient *mongo.Client
}

// GetRuntime sets up the runtime to be used at the start of the application
func GetRuntime(pConfig *Config) (*Runtime, error) {

	//------------------
	// LOGS
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	// log.SetLevel(log.WarnLevel)
	log.SetLevel(log.DebugLevel)

	//------------------
	// DB

	mongoURLstr := pConfig.MongoURL

	db, err := dbInit(mongoURLstr, pConfig)

	if err != nil {
		return nil, err
	}

	err = setupMongoIndexes(db.MongoClient.Database(GalleryDBName))
	if err != nil {
		return nil, err
	}

	// RUNTIME
	runtime := &Runtime{
		Config: pConfig,
		DB:     db,
	}

	// TEST REDIS CONNECTION
	client := redis.NewClient(&redis.Options{
		Addr:     runtime.Config.RedisURL,
		Password: runtime.Config.RedisPassword,
		DB:       0,
	})
	if err = client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %s\n connecting with URL %s", err, runtime.Config.RedisURL)
	}
	log.Info("redis connected! ✅")

	return runtime, nil
}

func dbInit(pMongoURLstr string,
	pConfig *Config) (*DB, error) {

	log.WithFields(log.Fields{}).Info("connecting to mongo...")

	var tlsConf *tls.Config
	if pConfig.MongoUseTLS {
		tlsCerts, err := accessSecret(context.Background(), viper.GetString(mongoTLSSecretName))
		if err != nil {
			return nil, err
		}
		tlsConf, err = dbGetCustomTLSConfig(tlsCerts)
		if err != nil {
			return nil, err
		}
	}
	mongoClient, err := connectMongo(pMongoURLstr, tlsConf)
	if err != nil {
		return nil, err
	}
	log.Info("mongo connected! ✅")

	db := &DB{
		MongoClient: mongoClient,
	}

	return db, nil
}

func dbGetCustomTLSConfig(pCerts []byte) (*tls.Config, error) {

	tlsConfig := new(tls.Config)
	tlsConfig.RootCAs = x509.NewCertPool()

	ok := tlsConfig.RootCAs.AppendCertsFromPEM(pCerts)
	if !ok {
		return nil, fmt.Errorf("unable to append certs from pem")
	}

	return tlsConfig, nil
}

func setupMongoIndexes(db *mongo.Database) error {
	b := true
	db.Collection("users").Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.M{"username_idempotent": 1},
		Options: &options.IndexOptions{
			Unique: &b,
			Sparse: &b,
		},
	})
	return nil
}

func connectMongo(pMongoURL string,
	pTLS *tls.Config,
) (*mongo.Client, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	mOpts := options.Client().ApplyURI(pMongoURL)

	// TLS
	if pTLS != nil {
		mOpts.SetTLSConfig(pTLS)
	}

	mClient, err := mongo.Connect(ctx, mOpts)
	if err != nil {
		return nil, err
	}

	err = mClient.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return mClient, nil
}
