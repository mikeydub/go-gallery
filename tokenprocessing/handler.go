package tokenprocessing

import (
	"context"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/gin-gonic/gin"

	db "github.com/mikeydub/go-gallery/db/gen/coredb"
	"github.com/mikeydub/go-gallery/service/multichain"
	"github.com/mikeydub/go-gallery/service/persist"
	"github.com/mikeydub/go-gallery/service/persist/postgres"
	sentryutil "github.com/mikeydub/go-gallery/service/sentry"
	"github.com/mikeydub/go-gallery/service/throttle"
	"github.com/mikeydub/go-gallery/service/tokenmanage"
)

const defaultSyncMaxRetries = 4

var contractSpecificRetries = map[persist.ContractIdentifiers]int{
	persist.NewContractIdentifiers("0x47a91457a3a1f700097199fd63c039c4784384ab", persist.ChainArbitrum): 24, // Prohibition
}

func handlersInitServer(ctx context.Context, router *gin.Engine, tp *tokenProcessor, mc *multichain.Provider, repos *postgres.Repositories, throttler *throttle.Locker, taskClient *cloudtasks.Client) *gin.Engine {
	// Retry tokens that failed during syncs, but don't retry tokens that failed during manual refreshes
	refreshManager := tokenmanage.New(ctx, taskClient)
	syncManager := tokenmanage.NewWithRetries(ctx, taskClient, syncMaxRetriesF(ctx, mc.Queries))

	mediaGroup := router.Group("/media")
	mediaGroup.POST("/process", func(c *gin.Context) {
		if hub := sentryutil.SentryHubFromContext(c); hub != nil {
			hub.Scope().AddEventProcessor(sentryutil.SpanFilterEventProcessor(c, 1000, 1*time.Millisecond, 8, true))
		}
		processBatch(tp, mc.Queries, syncManager)(c)
	})
	mediaGroup.POST("/process/token", processMediaForTokenIdentifiers(tp, mc.Queries, refreshManager))
	mediaGroup.POST("/tokenmanage/process/token", processMediaForTokenManaged(tp, mc.Queries, syncManager))
	mediaGroup.POST("/process/post-preflight", processPostPreflight(tp, syncManager, mc, repos.UserRepository))
	ownersGroup := router.Group("/owners")
	ownersGroup.POST("/process/contract", processOwnersForContractTokens(mc, throttler))
	ownersGroup.POST("/process/user", processOwnersForUserTokens(mc, mc.Queries))
	ownersGroup.POST("/process/alchemy", processOwnersForAlchemyTokens(mc, mc.Queries))
	ownersGroup.POST("/process/wallet-removal", processWalletRemoval(mc.Queries))
	contractsGroup := router.Group("/contracts")
	contractsGroup.POST("/detect-spam", detectSpamContracts(mc.Queries))

	return router
}

func syncMaxRetriesF(ctx context.Context, q *db.Queries) func(id persist.DBID) int {
	return func(id persist.DBID) int {
		td, _ := q.GetTokenDefinitionById(ctx, id)
		cID := persist.NewContractIdentifiers(td.ContractAddress, td.Chain)
		if retries, ok := contractSpecificRetries[cID]; ok {
			return retries
		}
		return defaultSyncMaxRetries
	}
}
