package membership

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gammazero/workerpool"
	"github.com/mikeydub/go-gallery/service/eth"
	"github.com/mikeydub/go-gallery/service/opensea"
	"github.com/mikeydub/go-gallery/service/persist"
	"github.com/mikeydub/go-gallery/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// MembershipTierIDs is a list of all membership tiers
var MembershipTierIDs = []persist.TokenID{"3", "4", "5", "6", "8"}

// UpdateMembershipTiers fetches all membership cards for all token IDs
func UpdateMembershipTiers(membershipRepository persist.MembershipRepository, userRepository persist.UserRepository, nftRepository persist.NFTRepository, ethClient *eth.Client) ([]persist.MembershipTier, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	membershipTiers := make([]persist.MembershipTier, len(MembershipTierIDs))
	tierChan := make(chan persist.MembershipTier)
	for _, v := range MembershipTierIDs {
		events, err := OpenseaFetchMembershipCards(persist.Address(viper.GetString("CONTRACT_ADDRESS")), persist.TokenID(v), 0, 0)
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch membership cards for token: %s, %w", v, err)
		}
		if len(events) == 0 {
			continue
		}
		time.Sleep(time.Second)
		go func(id persist.TokenID, events []opensea.Event) {
			tierChan <- processEvents(ctx, id, events, ethClient, userRepository, nftRepository, membershipRepository)
		}(v, events)
	}

	for i := 0; i < len(MembershipTierIDs); i++ {
		membershipTiers[i] = <-tierChan
	}
	return membershipTiers, nil
}

// UpdateMembershipTier fetches all membership cards for a token ID
func UpdateMembershipTier(pTokenID persist.TokenID, membershipRepository persist.MembershipRepository, userRepository persist.UserRepository, nftRepository persist.NFTRepository, ethClient *eth.Client) (persist.MembershipTier, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	events, err := OpenseaFetchMembershipCards(persist.Address(viper.GetString("CONTRACT_ADDRESS")), pTokenID, 0, 0)
	if err != nil {
		return persist.MembershipTier{}, fmt.Errorf("Failed to fetch membership cards for token: %s, %w", pTokenID, err)
	}
	if len(events) == 0 {
		return persist.MembershipTier{}, fmt.Errorf("No membership cards found for token: %s", pTokenID)
	}

	return processEvents(ctx, pTokenID, events, ethClient, userRepository, nftRepository, membershipRepository), nil
}

// UpdateMembershipTiersToken fetches all membership cards for a token ID
func UpdateMembershipTiersToken(membershipRepository persist.MembershipRepository, userRepository persist.UserRepository, nftRepository persist.TokenRepository, ethClient *eth.Client) ([]persist.MembershipTier, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	membershipTiers := make([]persist.MembershipTier, len(MembershipTierIDs))
	tierChan := make(chan persist.MembershipTier)
	for _, v := range MembershipTierIDs {
		go func(id persist.TokenID) {
			tierChan <- processEventsToken(ctx, id, ethClient, userRepository, nftRepository, membershipRepository)
		}(v)
	}

	for i := 0; i < len(MembershipTierIDs); i++ {
		membershipTiers[i] = <-tierChan
	}
	return membershipTiers, nil
}

// UpdateMembershipTierToken fetches all membership cards for a token ID
func UpdateMembershipTierToken(pTokenID persist.TokenID, membershipRepository persist.MembershipRepository, userRepository persist.UserRepository, nftRepository persist.TokenRepository, ethClient *eth.Client) (persist.MembershipTier, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	events, err := OpenseaFetchMembershipCards(persist.Address(viper.GetString("CONTRACT_ADDRESS")), pTokenID, 0, 0)
	if err != nil {
		return persist.MembershipTier{}, fmt.Errorf("Failed to fetch membership cards for token: %s, %w", pTokenID, err)
	}
	if len(events) == 0 {
		return persist.MembershipTier{}, fmt.Errorf("No membership cards found for token: %s", pTokenID)
	}

	return processEventsToken(ctx, pTokenID, ethClient, userRepository, nftRepository, membershipRepository), nil
}

// OpenseaFetchMembershipCards recursively fetches all membership cards for a token ID
func OpenseaFetchMembershipCards(contractAddress persist.Address, tokenID persist.TokenID, pOffset int, pRetry int) ([]opensea.Event, error) {

	client := &http.Client{
		Timeout: time.Minute,
	}

	result := []opensea.Event{}

	urlStr := fmt.Sprintf("https://api.opensea.io/api/v1/events?asset_contract_address=%s&token_id=%s&only_opensea=false&offset=%d&limit=50", contractAddress, tokenID.Base10String(), pOffset)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-KEY", viper.GetString("OPENSEA_API_KEY"))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		if resp.StatusCode == 429 {
			if pRetry > 3 {
				return nil, fmt.Errorf("timed out fetching membership cards %d at url: %s", tokenID.Base10Int(), urlStr)
			}

			logrus.Warnf("Opensea API rate limit exceeded, retrying in 5 seconds")
			time.Sleep(time.Second * 2)
			return OpenseaFetchMembershipCards(contractAddress, tokenID, pOffset, pRetry+1)
		}
		return nil, fmt.Errorf("unexpected status code: %d - url: %s", resp.StatusCode, urlStr)
	}

	response := &opensea.Events{}
	err = util.UnmarshallBody(response, resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}
	result = append(result, response.Events...)
	if len(response.Events) == 50 {
		next, err := OpenseaFetchMembershipCards(contractAddress, tokenID, pOffset+50, pRetry)
		if err != nil {
			return nil, err
		}
		result = append(result, next...)
	}

	return result, nil
}

func processEvents(ctx context.Context, id persist.TokenID, events []opensea.Event, ethClient *eth.Client, userRepository persist.UserRepository, nftRepository persist.NFTRepository, membershipRepository persist.MembershipRepository) persist.MembershipTier {
	tier := persist.MembershipTier{
		TokenID:     id,
		LastUpdated: persist.LastUpdatedTime(time.Now()),
	}
	logrus.Infof("Fetching membership tier: %s", id)

	asset := events[0].Asset
	tier.Name = persist.NullString(asset.Name)
	tier.AssetURL = persist.NullString(asset.ImageURL)

	logrus.Infof("Fetched membership cards for token %s with name %s and asset URL %s ", id, tier.Name, tier.AssetURL)
	tier.Owners = make([]persist.MembershipOwner, 0, len(events))

	ownersChan := make(chan persist.MembershipOwner)
	wp := workerpool.New(10)
	for _, event := range events {
		f := func() {
			membershipOwner := persist.MembershipOwner{Address: event.ToAccount.Address}

			if event.ToAccount.Address.Address() != common.BigToAddress(big.NewInt(0)) {
				// does to have the NFT?
				hasNFT, _ := ethClient.HasNFT(ctx, id, event.ToAccount.Address)
				if hasNFT {
					if glryUser, err := userRepository.GetByAddress(ctx, event.ToAccount.Address); err == nil && glryUser.Username != "" {
						membershipOwner.Username = glryUser.Username
						membershipOwner.UserID = glryUser.ID

						nfts, err := nftRepository.GetByUserID(ctx, glryUser.ID)
						if err == nil && len(nfts) > 0 {
							nftURLs := make([]persist.NullString, 0, 3)
							for i, nft := range nfts {
								if i == 3 {
									break
								}
								if nft.ImagePreviewURL != "" {
									nftURLs = append(nftURLs, nft.ImagePreviewURL)
								} else {
									i--
									continue
								}
							}
							membershipOwner.PreviewNFTs = nftURLs
						}
					}

				}
			}
			ownersChan <- membershipOwner
		}
		wp.Submit(f)
	}
	receivedOwners := map[persist.Address]bool{}
	for i := 0; i < len(events); i++ {
		owner := <-ownersChan
		if receivedOwners[owner.Address] {
			continue
		}
		tier.Owners = append(tier.Owners, owner)
		receivedOwners[owner.Address] = true
	}
	wp.StopWait()
	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		err := membershipRepository.UpsertByTokenID(ctx, id, tier)
		if err != nil {
			logrus.Errorf("Error upserting membership tier: %s", err)
		}
	}()
	return tier
}

func processEventsToken(ctx context.Context, id persist.TokenID, ethClient *eth.Client, userRepository persist.UserRepository, nftRepository persist.TokenRepository, membershipRepository persist.MembershipRepository) persist.MembershipTier {
	tier := persist.MembershipTier{
		TokenID:     id,
		LastUpdated: persist.LastUpdatedTime(time.Now()),
	}
	logrus.Infof("Fetching membership tier: %s", id)

	tokens, err := nftRepository.GetByTokenIdentifiers(ctx, persist.TokenID(id), persist.Address(viper.GetString("CONTRACT_ADDRESS")), -1, 0)
	if err != nil || len(tokens) == 0 {
		logrus.WithError(err).Errorf("Failed to fetch membership cards for token: %s", id)
		return tier
	}
	initialToken := tokens[0]

	tier.Name = persist.NullString(initialToken.Name)
	tier.AssetURL = persist.NullString(initialToken.Media.MediaURL)

	logrus.Infof("Fetched membership cards for token %s with name %s and asset URL %s ", id, tier.Name, tier.AssetURL)

	tier.Owners = make([]persist.MembershipOwner, 0, len(tokens))

	ownersChan := make(chan persist.MembershipOwner)
	wp := workerpool.New(10)
	for _, token := range tokens {
		f := func() {
			membershipOwner := persist.MembershipOwner{Address: token.OwnerAddress}
			if glryUser, err := userRepository.GetByAddress(ctx, token.OwnerAddress); err == nil && glryUser.Username != "" {
				membershipOwner.Username = glryUser.Username
				membershipOwner.UserID = glryUser.ID

				nfts, err := nftRepository.GetByUserID(ctx, glryUser.ID, -1, 0)
				if err == nil && len(nfts) > 0 {
					nftURLs := make([]persist.NullString, 0, 3)
					for i, nft := range nfts {
						if i == 3 {
							break
						}
						if nft.Media.PreviewURL != "" {
							nftURLs = append(nftURLs, nft.Media.PreviewURL)
						} else if nft.Media.MediaURL != "" {
							nftURLs = append(nftURLs, nft.Media.MediaURL)
						} else {
							i--
							continue
						}
					}
					membershipOwner.PreviewNFTs = nftURLs
				}
			}
			ownersChan <- membershipOwner
		}
		wp.Submit(f)
	}
	for i := 0; i < len(tokens); i++ {
		tier.Owners = append(tier.Owners, <-ownersChan)
	}
	wp.StopWait()
	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		err := membershipRepository.UpsertByTokenID(ctx, id, tier)
		if err != nil {
			logrus.Errorf("Error upserting membership tier: %s", err)
		}
	}()
	return tier
}
