package resource_common_service

import (
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	token2 "nexsoft.co.id/nextrac2/resource_common_service/token"
	"nexsoft.co.id/nextrac2/token"
	"strings"
	"time"
)

func MandatoryValidateJWTToken(jwtTokenStr string, apiAllowedScope string, resourceID string, key string) (jwtToken token.PayloadJWTToken, err model.ResourceCommonErrorModel) {
	jwtToken, err = token2.ValidateJWT(jwtTokenStr, key)
	if err.Error != nil {
		return
	}

	if !common.CheckIsResourceIDExist(jwtToken.Resource, resourceID) {
		err = common.GenerateForbiddenByResourceID()
		return
	}

	scopeSplit := strings.Split(apiAllowedScope, " ")
	for i := 0; i < len(scopeSplit); i++ {
		if !common.CheckIsScopeExist(jwtToken.Scope, scopeSplit[i]) {
			err = common.GenerateForbiddenByScope()
			return
		}
	}

	return
}

func ValidateTokenWithoutCheckSignature(jwtTokenStr string, apiAllowedScope string, resourceID string) (jwtToken token.PayloadJWTToken, err model.ResourceCommonErrorModel) {
	jwtToken, err = ConvertJWTToPayload(jwtTokenStr)
	if err.Error != nil {
		return
	}

	if time.Now().Unix() > jwtToken.ExpiresAt {
		err = common.GenerateExpiredToken()
		return
	}

	if !common.CheckIsResourceIDExist(jwtToken.Resource, resourceID) {
		fmt.Println("jwtToken >>", jwtToken.Resource)
		err = common.GenerateForbiddenByResourceID()
		return
	}

	scopeSplit := strings.Split(apiAllowedScope, " ")
	for i := 0; i < len(scopeSplit); i++ {
		if !common.CheckIsScopeExist(jwtToken.Scope, scopeSplit[i]) {
			err = common.GenerateForbiddenByScope()
			return
		}
	}

	return
}

func ConvertJWTToPayload(jwtTokenStr string) (jwtToken token.PayloadJWTToken, err model.ResourceCommonErrorModel) {
	splitJWT := strings.Split(jwtTokenStr, ".")
	if len(splitJWT) == 3 {
		payload := splitJWT[1]

		byteData, errs := util.Base64decoder(payload)
		if errs != nil {
			err = common.GenerateUnknownError(errs)
			return
		}
		_ = json.Unmarshal(byteData, &jwtToken)
	} else {
		err = common.GenerateUnauthorizedTokenError()
	}

	return
}
