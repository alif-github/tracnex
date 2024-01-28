package GenerateInternalToken

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/token"
	"strconv"
	"time"
)

func (input generateInternalTokenService) StartGenerateInToken(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GenerateInternalTokenRequestDTO

	inputStruct, err = input.readBodyAndValidate(request, contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.generateInternalTokenCustom(inputStruct)

	output.Status = out.StatusResponse{
		Code: 		"OK",
		Message: 	"Sukses Generate Token",
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input generateInternalTokenService) generateInternalTokenCustom(inputStructInterface interface{}) string {
	inputStruct := inputStructInterface.(in.GenerateInternalTokenRequestDTO)

	userClientID := config.ApplicationConfiguration.GetClientCredentialsClientID()
	if inputStruct.ClientID != "" {
		userClientID = inputStruct.ClientID
	}

	if inputStruct.ClientID == "" {
		inputStruct.ClientID = userClientID
	}

	usedUserID := config.ApplicationConfiguration.GetClientCredentialsAuthUserID()
	if inputStruct.AuthUserID > 0 {
		usedUserID = inputStruct.AuthUserID
	}

	tokenCode := token.PayloadJWTInternal{
		Locale:     inputStruct.Locale,
		ClientID:   inputStruct.ClientID,
		UserClient: userClientID,
		Resource:   inputStruct.Destination,
		Version:    config.ApplicationConfiguration.GetServerVersion(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * 3 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    inputStruct.Issuer,
			Subject:   strconv.Itoa(int(usedUserID)),
		},
	}

	jwtToken, _ := token.JWTToken{}.GenerateToken(tokenCode, config.ApplicationConfiguration.GetJWTToken().Internal)

	return jwtToken
}