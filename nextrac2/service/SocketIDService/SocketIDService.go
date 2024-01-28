package SocketIDService

import (
	"database/sql"
	"encoding/json"
	"errors"
	utilCommon "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func DoUpdateSocketID(tx *sql.Tx, body interface{}, contextModel *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	fileName := "SocketIDService.go"
	funcName := "DoUpdateSocketID"
	clientMappingBody := body.(in.ClientMappingForUIRequest)
	clientMappingRepository := getClientmappingRepository(clientMappingBody, contextModel.AuthAccessTokenModel, now)

	clientMappingOnDB, err := dao.ClientMappingDAO.GetClientMappingsForChangeSocketID(tx, clientMappingRepository)
	if err.Error != nil {
		return
	}

	if len(clientMappingOnDB) < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientID)
		return
	}

	err = validate(clientMappingBody, contextModel)
	if err.Error != nil {
		return
	}

	for _, item := range clientMappingOnDB {
		auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, now, dao.ClientMappingDAO.TableName, item.ID.Int64, 0)...)
	}

	err = updateClientToAuthenticationServer(clientMappingBody, contextModel)
	if err.Error != nil {
		return
	}

	err = dao.ClientMappingDAO.UpdateSocketID(tx, clientMappingRepository)
	return
}

func validate(clientMappingBody in.ClientMappingForUIRequest, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	fileName := "SocketIDService.go"
	funcName := "validate"
	err = clientMappingBody.ValidateUpdateClientID()
	if err.Error != nil {
		return
	}

	//---------- Check is client type ND6 or not
	clientType, err := dao.ClientTypeDAO.ValidationClientType(serverconfig.ServerAttribute.DBConnection, repository.ClientTypeModel{
		ID: sql.NullInt64{Int64: clientMappingBody.ClientTypeId},
		ClientType: sql.NullString{String: constanta.ND6},
	})
	if err.Error != nil {
		return
	}
	if clientType.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, "client_type_id")
	}

	//---------- Check client id
	//_, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*contextModel)
	//if isOnlyHaveOwnAccess {
	//	if clientMappingBody.ClientId != contextModel.LoggerModel.ClientID {
	//		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
	//		return
	//	}
	//}

	if clientMappingBody.ClientId != contextModel.LoggerModel.ClientID {
		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func getClientmappingRepository(clientMapingBody in.ClientMappingForUIRequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.ClientMappingModel {
	return repository.ClientMappingModel{
		ID:             sql.NullInt64{Int64: clientMapingBody.ID},
		SocketID:       sql.NullString{String: clientMapingBody.SocketID},
		ClientTypeID:   sql.NullInt64{Int64: clientMapingBody.ClientTypeId},
		ClientID:       sql.NullString{String: clientMapingBody.ClientId},
		UpdatedAt:      sql.NullTime{Time: now},
		UpdatedBy:      sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedClient:  sql.NullString{String: authAccessToken.ClientID},
	}
}

func updateClientToAuthenticationServer(clientBody in.ClientMappingForUIRequest, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	fileName := "SocketIDService.go"
	funcName := "UpdateClientToAuthenticationServer"
	var responseApi out.APIResponse

	//---------- Todo get credential token
	authServer := config.ApplicationConfiguration.GetAuthenticationServer()
	internalToken := resource_common_service.GenerateInternalToken(constanta.AuthDestination, 0, contextModel.AuthAccessTokenModel.ClientID, config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
	if err.Error != nil {
		return
	}

	//---------- Todo get detail client
	getDetailClientURL := authServer.Host + authServer.PathRedirect.InternalClient.CrudClient + "/" + clientBody.ClientId
	dataTobeUpdate, err := getClientOnAuthentication(internalToken, getDetailClientURL, contextModel, clientBody)
	if err.Error != nil {
		return
	}

	//---------- Todo update client
	updateClientUrl := authServer.Host + authServer.PathRedirect.InternalClient.CrudClient + "/" + clientBody.ClientId
	statusCode, bodyResult, errorHit := common.HitUpdateClientAuthenticationServer(internalToken, updateClientUrl, dataTobeUpdate, contextModel)
	if errorHit != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorHit)
		return
	}

	readError := json.Unmarshal([]byte(bodyResult), &responseApi)
	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}


	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	}else {
		causedBy := errors.New(responseApi.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(fileName, funcName, statusCode, responseApi.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func getClientOnAuthentication(internalToken string, getDetailClientURL string, contextModel *applicationModel.ContextModel,
	clientBody in.ClientMappingForUIRequest) (dataToUpdate authentication_request.ClientUpdateRequest, err errorModel.ErrorModel) {

	fileName := "SocketIDService.go"
	funcName := "getClientOnAuthentication"
	var responseApi out.APIResponse
	var clientOnAuth authentication_response.DetailClientResponse
	var additonalInfoToUpdate []model2.AdditionalInformation

	additonalInfoToUpdate = append(additonalInfoToUpdate, model2.AdditionalInformation{
		Key:   "socket_id",
		Value: clientBody.SocketID,
	})

	statusCode, bodyResult, errorHit := common.HitGetDetailClientAuthenticationServer(internalToken, getDetailClientURL, contextModel)
	if errorHit != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorHit)
		return
	}

	readError := json.Unmarshal([]byte(bodyResult), &responseApi)
	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	if statusCode == 200 {
		readError = json.Unmarshal([]byte(utilCommon.StructToJSON(responseApi.Nexsoft.Payload.Data.Content)), &clientOnAuth)
		if readError != nil {
			err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
			return
		}
		err = errorModel.GenerateNonErrorModel()
	}else {
		causedBy := errors.New(responseApi.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(fileName, funcName, statusCode, responseApi.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	dataToUpdate = authentication_request.ClientUpdateRequest{
		Scope:                clientOnAuth.Scope,
		RedirectUri:          clientOnAuth.RedirectUri,
		IPWhitelist:          clientOnAuth.IPWhitelist,
		AccessTokenValidity:  clientOnAuth.AccessTokenValidity,
		RefreshTokenValidity: clientOnAuth.RefreshTokenValidity,
		MultipleLogin:        clientOnAuth.MultipleLogin,
		MaxAuthFail:          clientOnAuth.MaxAuthFail,
		Locale:               clientOnAuth.Locale,
		UpdatedAtString:      clientOnAuth.UpdatedAtString,
		ClientInformation:    clientOnAuth.ClientInformation,
	}

	for i, updateItem := range dataToUpdate.ClientInformation {
		for j, itemToUpdate := range additonalInfoToUpdate {
			if updateItem.Key == itemToUpdate.Key {
				dataToUpdate.ClientInformation[i] = itemToUpdate
				additonalInfoToUpdate = append(additonalInfoToUpdate[:j], additonalInfoToUpdate[j+1:]...)
			}
		}
	}

	if len(additonalInfoToUpdate) > 0 {
		dataToUpdate.ClientInformation = append(dataToUpdate.ClientInformation, additonalInfoToUpdate...)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}