package token

import (
	"github.com/dgrijalva/jwt-go"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/token"
	"strings"
)

func ValidateJWT(jwtTokenStr string, key string) (payload token.PayloadJWTToken, errors model.ResourceCommonErrorModel) {
	claims := &token.PayloadJWTToken{}

	jwtToken, err := jwt.ParseWithClaims(jwtTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		if strings.Contains(err.Error(), "expired") {
			errors = common.GenerateExpiredToken()
		} else {
			errors = common.GenerateUnauthorizedTokenError()
		}
		return
	}

	if jwtToken.Header["alg"] != "HS512" && jwtToken.Header["alg"] != "HS256" {
		return payload, common.GenerateInvalidMethode()
	}

	payload = *jwtToken.Claims.(*token.PayloadJWTToken)
	return
}

func ValidateJWTInternal(jwtTokenStr string, key string) (payload token.PayloadJWTInternal, err model.ResourceCommonErrorModel) {
	claims := &token.PayloadJWTInternal{}
	jwtToken, errors := jwt.ParseWithClaims(jwtTokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if errors != nil {
		if strings.Contains(errors.Error(), "expired") {
			err = common.GenerateExpiredToken()
		} else {
			err = common.GenerateUnauthorizedTokenError()
		}
		return
	}

	if jwtToken.Header["alg"] != "HS512" && jwtToken.Header["alg"] != "HS256" {
		return payload, common.GenerateInvalidMethode()
	}

	payload = *jwtToken.Claims.(*token.PayloadJWTInternal)
	return
}
