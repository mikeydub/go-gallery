package glry_lib

import (
	// "fmt"
	// "time"
	"net/http"
	"context"
	log "github.com/sirupsen/logrus"
	gfcore "github.com/gloflow/gloflow/go/gf_core"
	gfrpclib "github.com/gloflow/gloflow/go/gf_rpc_lib"
	"github.com/mikeydub/go-gallery/glry_core"
	"github.com/mikeydub/go-gallery/glry_db"
	"github.com/mikeydub/go-gallery/glry_extern_services"
)

//-------------------------------------------------------------
func HandlersInit(pRuntime *glry_core.Runtime) {



	log.WithFields(log.Fields{}).Debug("initializing HTTP handlers")

	//-------------------------------------------------------------
	// HEALTH

	gfrpclib.Create_handler__http("/glry/v1/health",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gfcore.Gf_error) {

			//------------------
			// INPUT


			//------------------
			

			log.WithFields(log.Fields{}).Debug("/health")

			//------------------
			// OUTPUT
			dataMap := map[string]interface{}{
				"msg": "gallery operational",
			}

			//------------------

			return dataMap, nil
		},
		pRuntime.RuntimeSys)
	
	//-------------------------------------------------------------
	// ME

	gfrpclib.Create_handler__http("/glry/v1/auth/me",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gfcore.Gf_error) {

			//------------------
			// INPUT


			//------------------
			
			//------------------
			// OUTPUT
			dataMap := map[string]interface{}{
				"msg": "gallery operational",
			}

			//------------------

			return dataMap, nil
		},
		pRuntime.RuntimeSys)

	//-------------------------------------------------------------
	// AUTH_USER_LOGIN

	gfrpclib.Create_handler__http("/glry/v1/auth/login",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gfcore.Gf_error) {

			//------------------
			// INPUT

			var input GLRYauthUserLoginInput
			inputParsed, gErr := gfrpclib.Get_http_input_to_struct(input, pResp, pReq, pRuntime.RuntimeSys)
			if gErr != nil {
				return nil, gErr
			}

			//------------------
			
			// USER_LOGIN__PIPELINE
			validBool, userJWTtokenStr, gErr := AuthUserLoginAndMemorizeAttemptPipeline(inputParsed.(*GLRYauthUserLoginInput),
				pReq,
				pCtx,
				pRuntime)
			if gErr != nil {
				return nil, gErr
			}

			/*
			// ADD!! - going forward we should follow this approach, after v1
			// SET_JWT_COOKIE
			expirationTime := time.Now().Add(time.Duration(pRuntime.Config.JWTtokenTTLsecInt/60) * time.Minute)
			http.SetCookie(pResp, &http.Cookie{
				Name:    "glry_token",
				Value:   userJWTtokenStr,
				Expires: expirationTime,
			})*/

			//------------------
			// OUTPUT
			dataMap := map[string]interface{}{
				"valid":     validBool,
				"jwt_token": userJWTtokenStr,
			}

			//------------------

			return dataMap, nil
		},
		pRuntime.RuntimeSys)

	//-------------------------------------------------------------
	// AUTH_GET_PREFLIGHT

	gfrpclib.Create_handler__http("/glry/v1/auth/get_preflight",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gfcore.Gf_error) {

			//------------------
			// INPUT

			qMap        := pReq.URL.Query()
			userAddrStr := qMap["addr"][0]

			input := &GLRYauthUserGetPublicInfoInput{
				AddressStr: glry_db.GLRYuserAddress(userAddrStr),
			}

			//------------------
			
			// GET_PUBLIC_INFO
			output, gErr := AuthUserGetPreflightPipeline(input, pCtx, pRuntime)
			if gErr != nil {
				return nil, gErr
			}

			//------------------
			// OUTPUT
			dataMap := map[string]interface{}{
				"nonce": output.NonceStr,
			}

			//------------------

			return dataMap, nil
		},
		pRuntime.RuntimeSys)

	//-------------------------------------------------------------
	// AUTH_SIGNUP

	gfrpclib.Create_handler__http("/glry/v1/auth/signup",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gfcore.Gf_error) {
			
			if pReq.Method == "POST" {
				//------------------
				// INPUT

				var input GLRYauthUserCreateInput
				inputParsed, gErr := gfrpclib.Get_http_input_to_struct(input, pResp, pReq, pRuntime.RuntimeSys)
				if gErr != nil {
					return nil, gErr
				}

				//------------------
				// GET_PUBLIC_INFO
				user, gErr := AuthUserCreatePipeline(inputParsed.(*GLRYauthUserCreateInput), pCtx, pRuntime)
				if gErr != nil {
					return nil, gErr
				}

				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					"id": user.IDstr,
					// "nonce": user.NonceInt,
				}

				//------------------

				return dataMap, nil
			}

			return nil, nil
		},
		pRuntime.RuntimeSys)

	//-------------------------------------------------------------
	// COLLECTION_CREATE

	gfrpclib.Create_handler__http("/glry/v1/collections/create",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gfcore.Gf_error) {

			//------------------
			// INPUT

			var input GLRYcollCreateInput
			inputParsed, gErr := gfrpclib.Get_http_input_to_struct(input, pResp, pReq, pRuntime.RuntimeSys)
			if gErr != nil {
				return nil, gErr
			}

			// FINISH!! - get user_id mechanism
			userIDstr := ""

			//------------------


			coll, gErr := CollCreatePipeline(inputParsed.(*GLRYcollCreateInput), userIDstr, pCtx, pRuntime)
			if gErr != nil {
				return nil, gErr
			}
			
			//------------------
			// OUTPUT
			dataMap := map[string]interface{}{
				"coll_id": coll.IDstr,
			}

			//------------------

			return dataMap, nil
		},
		pRuntime.RuntimeSys)

	//-------------------------------------------------------------
	// COLLECTION_DELETE
	gfrpclib.Create_handler__http("/glry/v1/collections/delete",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gfcore.Gf_error) {

			//------------------
			// INPUT

			var input GLRYcollDeleteInput
			inputParsed, gErr := gfrpclib.Get_http_input_to_struct(input, pResp, pReq, pRuntime.RuntimeSys)
			if gErr != nil {
				return nil, gErr
			}

			//------------------

			_, gErr = CollDeletePipeline(inputParsed.(*GLRYcollDeleteInput), pCtx, pRuntime)
			if gErr != nil {
				return nil, gErr
			}

			//------------------
			// OUTPUT
			dataMap := map[string]interface{}{
		
			}

			//------------------

			return dataMap, nil
		},
		pRuntime.RuntimeSys)
		
	//-------------------------------------------------------------
	// NFTS_FOR_USER__GET
	gfrpclib.Create_handler__http("/glry/v1/nfts/user_get",
		func(pCtx context.Context, pResp http.ResponseWriter, pReq *http.Request) (map[string]interface{}, *gfcore.Gf_error) {

			if pReq.Method == "GET" {

				//------------------
				// INPUT

				//------------------

				userIDstr := "7bfaafcc-722e-4dce-986f-fe0d9bee2047"
				nfts, gErr := glry_db.NFTgetByUserID(userIDstr, pCtx, pRuntime)
				if gErr != nil {
					return nil, gErr
				}

				//------------------
				// OUTPUT
				dataMap := map[string]interface{}{
					"nfts": nfts,
				}

				//------------------

				return dataMap, nil
			}

			return nil, nil
		},
		pRuntime.RuntimeSys)

	//-------------------------------------------------------------
	// NFTS_FROM_OPENSEA__GET
	gfrpclib.Create_handler__http("/glry/v1/nfts/opensea_get",
		func(pCtx context.Context, pResp http.ResponseWriter,
			pReq *http.Request) (map[string]interface{}, *gfcore.Gf_error) {

			//------------------
			// INPUT

			//------------------

			ownerWalletAddressStr := "0x70d04384b5c3a466ec4d8cfb8213efc31c6a9d15"
			_, gErr := glry_extern_services.OpenSeaPipelineAssetsForAcc(ownerWalletAddressStr, pCtx, pRuntime)
			if gErr != nil {
				return nil, gErr
			}


			//------------------
			// OUTPUT
			dataMap := map[string]interface{}{
	
			}

			//------------------

			return dataMap, nil
		},
		pRuntime.RuntimeSys)

	//-------------------------------------------------------------
}