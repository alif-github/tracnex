package Login

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/session"
	"nexsoft.co.id/nextrac2/token"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type nexmileLoginService struct {
	service.AbstractService
}

var NexmileLoginService = nexmileLoginService{}.New()

func (input nexmileLoginService) New() (output nexmileLoginService) {
	output.FileName = "NexmileLoginService.go"
	return
}

func (input nexmileLoginService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.VerifyDTOIn) errorModel.ErrorModel) (inputStruct in.VerifyDTOIn, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidate"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	err = validation(&inputStruct)
	return
}

func (input nexmileLoginService) LoginNexmileService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct    in.VerifyDTOIn
		tokenStruct    in.TokenDTOIn
		codeChallenger string
		codeVerifier   = util2.GenerateCryptoRandom()
	)

	//--- Generate Verifier Code
	codeChallenger = util2.GenerateSHA256(codeVerifier)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateLoginNexmile)
	if err.Error != nil {
		return
	}

	//--- Authorize
	err = input.authorizeNexmileProcess(&inputStruct, codeChallenger, contextModel)
	if err.Error != nil {
		return
	}

	//--- Verify
	tokenStruct, err = input.verifyNexmileProcess(inputStruct, codeVerifier, contextModel)
	if err.Error != nil {
		return
	} else {
		if inputStruct.CheckPassword {
			output.Status = out.StatusResponse{
				Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
				Message: session.GenerateLoginI18NMessage("TOKEN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
			}

			err = errorModel.GenerateNonErrorModel()
			return
		}
	}

	//--- Token
	header, err = input.tokenNexmileProcess(tokenStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: session.GenerateLoginI18NMessage("TOKEN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexmileLoginService) authorizeNexmileProcess(inputStruct *in.VerifyDTOIn, codeChallenger string, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		headerToken map[string]string
		funcName    = "authorizeNexmileProcess"
		tokenName   = constanta.TokenHeaderNameConstanta
	)

	headerToken, err = AuthorizeService.HitAuthorizeAuthenticationServer(in.AuthorizeDTOIn{CodeChallenger: codeChallenger}, contextModel)
	if err.Error != nil {
		return
	}

	inputStruct.Authorize = headerToken[tokenName]
	if util.IsStringEmpty(inputStruct.Authorize) {
		err = errorModel.GenerateUnauthorizedClientError(input.FileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexmileLoginService) verifyNexmileProcess(inputStruct in.VerifyDTOIn, codeVerifier string, contextModel *applicationModel.ContextModel) (tokenStruct in.TokenDTOIn, err errorModel.ErrorModel) {
	var resultVerify out.VerifyDTOOut
	resultVerify, err = VerifyService.HitVerifyAuthenticationServer(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	tokenStruct = in.TokenDTOIn{
		AuthorizationCode: resultVerify.AuthorizationCode,
		CodeVerifier:      codeVerifier,
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexmileLoginService) tokenNexmileProcess(tokenStruct in.TokenDTOIn, contextModel *applicationModel.ContextModel) (headerResult map[string]string, err errorModel.ErrorModel) {
	_, _, headerResult, err = TokenService.HitTokenAuthenticationServer(tokenStruct, contextModel, input.roleMapping)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexmileLoginService) validateLoginNexmile(inputStruct *in.VerifyDTOIn) errorModel.ErrorModel {
	return inputStruct.ValidateNexmileLoginDTO()
}

func (input nexmileLoginService) roleMapping(clientID string, token string, payload token.PayloadJWTToken) (authAccessTokenModel model2.AuthAccessTokenModel, authenticationRoleModel model2.AuthenticationRoleModel, authenticationDataModel model2.AuthenticationDataModel, errs error) {
	var (
		roleModel        repository.RoleMappingPersonProfileModel
		err              errorModel.ErrorModel
		tx               *sql.Tx
		userModel        repository.UserModel
		clientTokenModel repository.ClientTokenModel
		db               = serverconfig.ServerAttribute.DBConnection
		temp             = make(map[string][]string)
		tempDataScope    = make(map[string]interface{})
	)

	userModel = repository.UserModel{ClientID: sql.NullString{String: clientID}}
	roleModel, err = dao.UserDAO.RoleMappingUser(db, userModel)
	if err.Error != nil {
		errs = err.CausedBy
		return
	}

	authAccessTokenModel.ResourceUserID = roleModel.PersonProfileID.Int64
	authAccessTokenModel.RedisAuthAccessTokenModel = model2.RedisAuthAccessTokenModel{
		ResourceUserID: roleModel.PersonProfileID.Int64,
		IPWhiteList:    roleModel.IPWhitelist.String,
		SignatureKey:   roleModel.SignatureKey.String,
	}

	_ = json.Unmarshal([]byte(roleModel.Permissions.String), &temp)
	authenticationRoleModel.Role = roleModel.RoleName.String
	authenticationRoleModel.Permission = temp

	_ = json.Unmarshal([]byte(roleModel.Scope.String), &tempDataScope)
	authenticationDataModel.Group = roleModel.GroupName.String
	authenticationDataModel.Scope = tempDataScope

	tx, errs = db.Begin()
	if errs != nil {
		return
	}

	defer func() {
		if errs != nil && err.Error != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
		return
	}()

	clientTokenModel = repository.ClientTokenModel{
		ClientID:      sql.NullString{String: clientID},
		AuthUserID:    roleModel.AuthUserID,
		Token:         sql.NullString{String: token},
		ExpiredAt:     sql.NullTime{Time: time.Unix(payload.ExpiresAt, 0)},
		CreatedBy:     roleModel.AuthUserID,
		CreatedClient: sql.NullString{String: clientID},
	}

	err = dao.ClientTokenDAO.InsertClientToken(tx, clientTokenModel)
	if err.Error != nil {
		errs = err.CausedBy
		return
	}

	err = dao.UserDAO.UpdateLastTokenUser(tx, userModel)
	if err.Error != nil {
		errs = err.CausedBy
		return
	}

	return
}
