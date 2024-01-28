package ClientRegistrationNonOnPremiseService

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type clientRegistrationNonOnPremiseService struct {
	service.AbstractService
	service.GetListData
}

var ClientRegistrationNonOnPremiseService = clientRegistrationNonOnPremiseService{}.New()

func (input clientRegistrationNonOnPremiseService) New() (output clientRegistrationNonOnPremiseService) {
	output.FileName = "ClientRegistrationNonOnPremiseService.go"
	return
}

func (input clientRegistrationNonOnPremiseService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ClientRegistrationNonOnPremiseRequest) errorModel.ErrorModel) (inputStruct in.ClientRegistrationNonOnPremiseRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	err = validation(&inputStruct)
	return
}

func (input clientRegistrationNonOnPremiseService) createModelClientRegistrationNonOnPremise(inputStruct in.ClientRegistrationNonOnPremiseRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (clientRegistModel repository.ClientRegistNonOnPremiseModel) {
	var modelDetailUniqueID []repository.DetailUniqueID

	for _, itemUnique := range inputStruct.DetailClient {
		var detailUnique repository.DetailUniqueID

		detailUnique.UniqueID1.String = itemUnique.UniqueID1
		if itemUnique.UniqueID2 != "" {
			detailUnique.UniqueID2.String = itemUnique.UniqueID2
		}

		modelDetailUniqueID = append(modelDetailUniqueID, detailUnique)
	}

	clientRegistModel = repository.ClientRegistNonOnPremiseModel{
		ClientID:      sql.NullString{String: inputStruct.ClientID},
		ClientTypeID:  sql.NullInt64{Int64: inputStruct.ClientTypeID},
		DetailUnique:  modelDetailUniqueID,
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	return
}

func (input clientRegistrationNonOnPremiseService) getClientForValidateToAuthenticationServer(clientModel repository.ClientRegistNonOnPremiseModel, contextModel *applicationModel.ContextModel) (resultAuth repository.ClientRegistNonOnPremiseModel, err errorModel.ErrorModel) {
	var (
		funcName             = "getClientForValidateToAuthenticationServer"
		data                 authentication_response.CheckClientByViewResponse
		internalToken        string
		getClientInfoUrl     string
		bodyResult           string
		statusCode           int
		errorS               error
		authenticationServer = config.ApplicationConfiguration.GetAuthenticationServer()
	)

	internalToken = resource_common_service.GenerateInternalToken(constanta.AuthDestination, 0, "", config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
	getClientInfoUrl = authenticationServer.Host + authenticationServer.PathRedirect.InternalClient.CrudClient + "/" + clientModel.ClientID.String

	statusCode, bodyResult, errorS = common.HitGetClientAuthenticationServer(internalToken, getClientInfoUrl, contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	readError := json.Unmarshal([]byte(bodyResult), &data)
	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, readError)
		return
	}

	resultAuth.ClientID.String = data.Nexsoft.Payload.Data.Content.ClientID
	resultAuth.ClientSecret.String = data.Nexsoft.Payload.Data.Content.ClientSecret
	resultAuth.SignatureKey.String = data.Nexsoft.Payload.Data.Content.SignatureKey
	resultAuth.AliasName.String = data.Nexsoft.Payload.Data.Content.AliasName

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(data.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, data.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}

func (input clientRegistrationNonOnPremiseService) compareClientCredential(clientModel *repository.ClientRegistNonOnPremiseModel, clientAuth repository.ClientRegistNonOnPremiseModel) (err errorModel.ErrorModel) {
	funcName := "compareClientCredential"

	if clientModel.ClientID.String != clientAuth.ClientID.String {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientID)
		return
	}

	if clientModel.ClientSecret.String != clientAuth.ClientSecret.String {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientID)
		return
	}

	if clientModel.SignatureKey.String != clientAuth.SignatureKey.String {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientID)
		return
	}

	clientModel.AliasName.String = clientAuth.AliasName.String
	clientModel.FirstName.String = clientAuth.AliasName.String
	err = errorModel.GenerateNonErrorModel()
	return
}

func checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_user_clientid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ClientID)
		} else if service.CheckDBError(err, "uq_user_authuserid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.AuthUserID)
		} else if service.CheckDBError(err, "uq_clientrolescope_clientid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ClientID)
		}
	}
	return err
}
