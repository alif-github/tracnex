package token

import (
	"github.com/dgrijalva/jwt-go"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type JWTToken struct {
}

func (input JWTToken) generateJWT(Payload jwt.Claims, key string) (string, errorModel.ErrorModel) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, Payload)
	token, err := jwtToken.SignedString([]byte(key))
	if err != nil {
		return "", errorModel.GenerateUnknownError("JWT.go", "generateJWT", err)
	}
	return token, errorModel.GenerateNonErrorModel()
}

func (input JWTToken) GenerateToken(Payload jwt.Claims, key string) (string, errorModel.ErrorModel) {
	return input.generateJWT(Payload, key)
}
