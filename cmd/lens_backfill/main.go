package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/mikeydub/go-gallery/db/gen/coredb"
	"github.com/mikeydub/go-gallery/env"
	"github.com/mikeydub/go-gallery/service/lens"
	"github.com/mikeydub/go-gallery/service/persist"
	"github.com/mikeydub/go-gallery/service/persist/postgres"
	"github.com/mikeydub/go-gallery/util"
	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/conc/pool"
	"github.com/spf13/viper"
)

func main() {

	setDefaults()

	pg := postgres.NewPgxClient()

	l := lens.NewAPI(&http.Client{Timeout: 10 * time.Second})

	queries := coredb.New(pg)

	ctx := context.Background()

	// get every wallet with their owner user ID
	rows, err := pg.Query(ctx, `select u.id, w.address from users u join wallets w on w.id = any(u.wallets) where u.deleted = false and w.chain = 0 and w.deleted = false and u.universal = false order by u.created_at desc;`)
	if err != nil {
		panic(err)
	}

	p := pool.New().WithMaxGoroutines(10).WithErrors()

	for rows.Next() {
		var userID persist.DBID
		var walletAddress persist.Address

		err := rows.Scan(&userID, &walletAddress)
		if err != nil {
			panic(err)
		}

		p.Go(func() error {
			logrus.Infof("getting user %s (%s)", userID, walletAddress)
			u, err := l.DefaultProfileByAddress(ctx, walletAddress)
			if err != nil {
				logrus.Error(err)
				return nil
			}
			logrus.Infof("got user %s %s %s %s", u.Name, u.Handle, u.Picture.Optimized.URL, u.Bio)
			return queries.AddSocialToUser(ctx, coredb.AddSocialToUserParams{
				UserID: userID,
				Socials: persist.Socials{
					persist.SocialProviderLens: persist.SocialUserIdentifiers{
						Provider: persist.SocialProviderLens,
						ID:       u.ID,
						Display:  true,
						Metadata: map[string]interface{}{
							"username":          u.Handle,
							"name":              util.FirstNonEmptyString(u.Name, u.Handle),
							"profile_image_url": util.FirstNonEmptyString(u.Picture.Optimized.URL, u.Picture.URI),
							"bio":               u.Bio,
						},
					},
				},
			})

		})

	}

	err = p.Wait()
	if err != nil {
		panic(err)
	}

}

func setDefaults() {
	viper.SetDefault("ENV", "local")
	viper.SetDefault("POSTGRES_HOST", "0.0.0.0")
	viper.SetDefault("POSTGRES_PORT", 5432)
	viper.SetDefault("POSTGRES_USER", "gallery_backend")
	viper.SetDefault("POSTGRES_PASSWORD", "")
	viper.SetDefault("POSTGRES_DB", "postgres")

	viper.AutomaticEnv()

	if env.GetString("ENV") != "local" {
		logrus.Info("running in non-local environment, skipping environment configuration")
	} else {
		fi := "local"
		if len(os.Args) > 1 {
			fi = os.Args[1]
		}
		envFile := util.ResolveEnvFile("backend", fi)
		util.LoadEncryptedEnvFile(envFile)
	}

}
