package extern_services

import (
	"fmt"
	"context"
	"encoding/json"
	"strings"
	log "github.com/sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
	"github.com/mitchellh/mapstructure"
	gfcore "github.com/gloflow/gloflow/go/gf_core"
	// "github.com/davecgh/go-spew/spew"
)

//-------------------------------------------------------------
type GLRYopenSeaAsset struct {
	IDstr               string `bson:"id"        mapstructure:"id"`
	TokenIDstr          string `bson:"token_id"  mapstructure:"token_id"`
	NumberOfSalesInt    int    `bson:"num_sales" mapstructure:"num_sales"`
	
	// IMAGES
	ImageURLstr         string                   `bson:"image_url"           mapstructure:"image_url"`
	ImagePreviewURLstr  string                   `bson:"image_preview_url"   mapstructure:"image_preview_url"`
	ImageThumbURLstr    string                   `bson:"image_thumbnail_url" mapstructure:"image_thumbnail_url"`
	ImageOriginalURLstr string                   `bson:"image_original_url"  mapstructure:"image_original_url"`
	AnimationURLstr     string                   `bson:"animation_url"       mapstructure:"animation_url"`

	NameStr             string                   `bson:"name"           mapstructure:"name"`
	DescriptionStr      string                   `bson:"description"    mapstructure:"description"`
	ExternLinkStr       string                   `bson:"external_link"  mapstructure:"external_link"`
	AssetContract       GLRYopenSeaAssetContract `bson:"asset_contract" mapstructure:"asset_contract"`
	Owner               GLRYopenSeaOwner         `bson:"owner"          mapstructure:"owner"`
	PermaLinkStr        string                   `bson:"permalink"      mapstructure:"permalink"`

	// IMPORTANT!! - OpenSea (unlike Gallery) only allows an Asset to be in a single collection 
	Collection GLRYopenSeaCollection `bson:"collection" mapstructure:"collection"`
	
	Creator        GLRYopenSeaCreator `bson:"creator"      mapstructure:"creator"`
	LastSaleStr    string             `bson:"last_sale"    mapstructure:"last_sale"`
	ListingDateStr string             `bson:"listing_date" mapstructure:"listing_date"`
}

type GLRYopenSeaAssetContract struct {
	AddressStr     string `bson:"address"      mapstructure:"address"`
	CreatedDateStr string `bson:"created_date" mapstructure:"created_date"`
	NameStr        string `bson:"name"         mapstructure:"name"`
	OwnerInt       int    `bson:"owner"        mapstructure:"owner"`
	SymbolStr      string `bson:"symbol"       mapstructure:"symbol"`
	DescriptionStr string `bson:"description"  mapstructure:"description"`
}

type GLRYopenSeaOwner struct {
	UserStr            string `bson:"user"            mapstructure:"user"`
	ProfileImageURLstr string `bson:"profile_img_url" mapstructure:"profile_img_url"`
	AddressStr         string `bson:"address"         mapstructure:"address"`
}

type GLRYopenSeaCollection struct {
	CreatedDateStr       string `bson:"created_date"       mapstructure:"created_date"`
	DescriptionStr       string `bson:"description"        mapstructure:"description"`
	ExternalURLstr       string `bson:"external_url"       mapstructure:"external_url"`
	ImageURLstr          string `bson:"image_url"          mapstructure:"image_url"`
	NameStr              string `bson:"name"               mapstructure:"name"`
	PayoutAddressStr     string `bson:"payout_address"     mapstructure:"payout_address"`
	TwitterUsernameStr   string `bson:"twitter_username"   mapstructure:"twitter_username"`
	InstagramUsernameStr string `bson:"instagram_username" mapstructure:"instagram_username"`
}

type GLRYopenSeaCreator struct {
	UserStr            string `bson:"user"            mapstructure:"user"`
	ProfileImageURLstr string `bson:"profile_img_url" mapstructure:"profile_img_url"`
	AddressStr         string `bson:"address"         mapstructure:"address"`
}

//-------------------------------------------------------------
func OpenSeaFetchAssetsForAcc(pOwnerWalletAddressStr string,
	pCtx context.Context,
	pRuntimeSys *gfcore.Runtime_sys) ([]*GLRYopenSeaAsset, *gfcore.Gf_error) {

	/*{
	*	"id": 21976544,
	*	"token_id": "1137",
	*	"num_sales": 0,
		"background_color": null,
	*	"image_url": "https://lh3.googleusercontent.com/8S2uhc_74_JijJwYnNOEQvlnHs6dI4lU86k8Zj2WcelVG9Gp4hx62UDzf2B_R4cTMdd_03SLOV_rFZFF8_5vwFSEz76OX61of4ZPaA",
	*	"image_preview_url": "https://lh3.googleusercontent.com/8S2uhc_74_JijJwYnNOEQvlnHs6dI4lU86k8Zj2WcelVG9Gp4hx62UDzf2B_R4cTMdd_03SLOV_rFZFF8_5vwFSEz76OX61of4ZPaA=s250",
	*	"image_thumbnail_url": "https://lh3.googleusercontent.com/8S2uhc_74_JijJwYnNOEQvlnHs6dI4lU86k8Zj2WcelVG9Gp4hx62UDzf2B_R4cTMdd_03SLOV_rFZFF8_5vwFSEz76OX61of4ZPaA=s128",
	*	"image_original_url": "https://coin-nfts.s3.us-east-2.amazonaws.com/coin-500px.gif",
	*	"animation_url": null,
		"animation_original_url": null,
	*	"name": "April 14 2021",
	*	"description": "A special thank you for all of the hard work you put in to usher in an open financial system and bring Coinbase public.",
	*	"external_link": "https://www.coinbase.com/",
		"asset_contract": {
	*		"address": "0x6966ac85200cadd8c66d14d6c1a5431353edc8c9",
			"asset_contract_type": "non-fungible",
	*		"created_date": "2021-04-14T19:18:41.926804",
	*		"name": "Coinbase Direct Public Offering",
			"nft_version": "3.0",
			"opensea_version": null,
	*		"owner": 833403,
			"schema_name": "ERC721",
	*		"symbol": "$COIN",
			"total_supply": "0",
	*		"description": "We did it! This commemorative NFT is a special thank you to all of the people who worked tirelessly to bring Coinbase public.",
			"external_link": "https://www.coinbase.com/",
			"image_url": "https://lh3.googleusercontent.com/cudWSCgwfsRLdmHrZ7wBx74pk5xBLDOcAvxYMgQicyZ1wG3VeASwL5WXoJbR1P70CjGTCE-wc6mPpvV-AMGVhVZ9QTPzzcHGpal1=s120",
			"default_to_fiat": false,
			"dev_buyer_fee_basis_points": 0,
			"dev_seller_fee_basis_points": 0,
			"only_proxied_transfers": false,
			"opensea_buyer_fee_basis_points": 0,
			"opensea_seller_fee_basis_points": 250,
			"buyer_fee_basis_points": 0,
			"seller_fee_basis_points": 250,
			"payout_address": null
		},
		"owner": {
			"user": null,
	*		"profile_img_url": "https://storage.googleapis.com/opensea-static/opensea-profile/25.png",
	*		"address": "0x70d04384b5c3a466ec4d8cfb8213efc31c6a9d15",
			"config": "",
			"discord_id": ""
		},
		"permalink": "https://opensea.io/assets/0x6966ac85200cadd8c66d14d6c1a5431353edc8c9/1137",
		"collection": {
			"banner_image_url": null,
			"chat_url": null,
	*		"created_date": "2021-04-14T19:22:29.035968",
			"default_to_fiat": false,
	*		"description": "We did it! This commemorative NFT is a special thank you to all of the people who worked tirelessly to bring Coinbase public.",
			"dev_buyer_fee_basis_points": "0",
			"dev_seller_fee_basis_points": "0",
			"discord_url": null,
			"display_data": {
				"card_display_style": "contain",
				"images": []
			},
	*		"external_url": "https://www.coinbase.com/",
			"featured": false,
			"featured_image_url": null,
			"hidden": true,
			"safelist_request_status": "not_requested",
	*		"image_url": "https://lh3.googleusercontent.com/cudWSCgwfsRLdmHrZ7wBx74pk5xBLDOcAvxYMgQicyZ1wG3VeASwL5WXoJbR1P70CjGTCE-wc6mPpvV-AMGVhVZ9QTPzzcHGpal1=s120",
			"is_subject_to_whitelist": false,
			"large_image_url": "https://lh3.googleusercontent.com/cudWSCgwfsRLdmHrZ7wBx74pk5xBLDOcAvxYMgQicyZ1wG3VeASwL5WXoJbR1P70CjGTCE-wc6mPpvV-AMGVhVZ9QTPzzcHGpal1",
			"medium_username": null,
	*		"name": "Coinbase Direct Public Offering",
			"only_proxied_transfers": false,
			"opensea_buyer_fee_basis_points": "0",
			"opensea_seller_fee_basis_points": "250",
	*		"payout_address": null,
			"require_email": false,
			"short_description": null,
			"slug": "coinbase-direct-public-offering",
			"telegram_url": null,
	*		"twitter_username": null,
	*		"instagram_username": null,
			"wiki_url": null
		},
		"decimals": 0,
		"sell_orders": null,
		"creator": {
	*		"user": null,
	*		"profile_img_url": "https://storage.googleapis.com/opensea-static/opensea-profile/25.png",
	*		"address": "0x70d04384b5c3a466ec4d8cfb8213efc31c6a9d15",
			"config": "",
			"discord_id": ""
		},
		"traits": [],
	*	"last_sale": null,
		"top_bid": null,
	*	"listing_date": null,
		"is_presale": false,
		"transfer_fee_payment_token": null,
		"transfer_fee": null
	},*/


	offsetInt := 0
	limitInt := 50
	qsArgsMap := map[string]string{
		"owner":           pOwnerWalletAddressStr,
		"order_direction": "desc",
		"offset":          fmt.Sprintf("%d", offsetInt),
		"limit":           fmt.Sprintf("%d", limitInt),

	}




	qsLst := []string{}
	for k, v := range qsArgsMap {
		qsLst = append(qsLst, fmt.Sprintf("%s=%s", k, v))
	}
	qsStr := strings.Join(qsLst, "&")
	urlStr := fmt.Sprintf("https://api.opensea.io/api/v1/assets?%s", qsStr)



	log.WithFields(log.Fields{
			"url":                  urlStr,
			"owner_wallet_address": pOwnerWalletAddressStr,
		}).Info("making HTTP request to OpenSea API")
	_, respBytes, errs := gorequest.New().Get(urlStr).EndBytes()
	if len(errs) > 0 {
		

	}

	var response map[string]interface{}
	err := json.Unmarshal(respBytes, &response)
	if err != nil {
		gErr := gfcore.Error__create(fmt.Sprintf("failed to parse json response from OpenSea API"), 
			"json_decode_error",
			map[string]interface{}{"url": urlStr,},
			err, "glry_extern_services", pRuntimeSys)
		return nil, gErr
	}




	assetsLst := response["assets"].([]interface{})



	

	assetsForAccLst := []*GLRYopenSeaAsset{}

	for _, aMap := range assetsLst {

		var asset GLRYopenSeaAsset
		err := mapstructure.Decode(aMap, &asset)
		if err != nil {
			
			gErr := gfcore.Error__create("failed to load OpenSea asset map into a GLRYopenSeaAsset struct",
				"mapstruct__decode",
				map[string]interface{}{
					"url":                  urlStr,
					"owner_wallet_address": pOwnerWalletAddressStr,
				},
				err, "glry_extern_services", pRuntimeSys)
			
			return nil, gErr
		}

		assetsForAccLst = append(assetsForAccLst, &asset)
	}
	
	// spew.Dump(assetsForAccLst)

	return assetsForAccLst, nil
}