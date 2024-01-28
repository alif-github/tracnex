package PKCEService

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input pkceService) doReRegistrationUserPKCE(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	fileName := "Re-registrationPKCEService.go"
	funcName := "doRegistrationUserPKCE"

	inputStruct := inputStructInterface.(in.PKCEReRequest)
	var resourceIDDataList []out.ResourceList
	var result out.PKCEResponse

	//---------- STEP 1. Add resource nextrac to auth server
	structForAddResource := out.PKCEResponse {
		ClientID: inputStruct.ClientID,
	}

	//---------- STEP 2. Add resource nextrac to auth server
	_, err = input.addResourceNextracToAuthServer(structForAddResource, contextModel)
	if err.Error != nil {
		return
	} else {
		resourceIDDataList = append(resourceIDDataList, out.ResourceList {
			ResourceID: 	config.ApplicationConfiguration.GetServerResourceID(),
			Status: 		"OK",
		})
	}

	//---------- STEP 1. Update user to be deleted == false
	userStruct := in.UserRequest {
		ID: 		inputStruct.ID,
		ClientID:	inputStruct.ClientID,
		AuthUserID:	inputStruct.AuthUserID,
	}

	//---------- STEP 2. Update user to be deleted == false
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel,
		timeNow, dao.UserDAO.TableName, userStruct.ID, 0)...)

	//---------- STEP 3. Update user to be deleted == false
	err = input.doUpdateToUser(tx, userStruct, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	//---------- STEP 1. Update pkce client mapping
	pkceMappingStruct := in.PKCERequest {
		ParentClientID: 		inputStruct.PKCERequest.ParentClientID,
		PKCEClientMappingID: 	inputStruct.PKCEClientMappingID,
		FirstName: 				inputStruct.PKCERequest.FirstName,
		CompanyID: 				inputStruct.PKCERequest.CompanyID,
		BranchID:				inputStruct.PKCERequest.BranchID,
		ClientAlias: 			inputStruct.PKCERequest.ClientAlias,
	}

	//---------- STEP 2. Update pkce client mapping
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel,
		timeNow, dao.PKCEClientMappingDAO.TableName, pkceMappingStruct.PKCEClientMappingID, 0)...)

	//---------- STEP 3. Update pkce client mapping
	err = input.doUpdatePKCEClientMapping(tx, pkceMappingStruct, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	//---------- STEP 1. Add resource nexcloud
	addResourceNexcloud := in.AddResourceNexcloud {
		FirstName: 	inputStruct.PKCERequest.FirstName,
		LastName: 	inputStruct.PKCERequest.LastName,
		ClientID: 	inputStruct.ClientID,
	}

	//---------- Preparing for response body
	result = out.PKCEResponse{
		UserID: 	userStruct.AuthUserID,
		ClientID: 	userStruct.ClientID,
		Username: 	inputStruct.PKCERequest.Username,
	}

	//---------- STEP 2. Add resource nexcloud
	err = input.addResourceUserNexcloud(addResourceNexcloud, contextModel)
	if err.Error != nil {
		_, err = service.NewErrorAddResource(serverconfig.ServerAttribute.PKCEUserBundle, contextModel, "DETAIL_ERROR_FAILED_RESOURCE_NEXCLOUD_MESSAGE", fileName, funcName)
		resourceIDDataList = append(resourceIDDataList, out.ResourceList {
			ResourceID: 	constanta.NexCloudResourceID,
			Status: 		"FAIL",
		})
		result = input.addResourceInfoReRegistration(result, resourceIDDataList)
		output, _, err = service.CustomFailedResponsePayload(result, err, contextModel)
		return
	} else {
		resourceIDDataList = append(resourceIDDataList, out.ResourceList {
			ResourceID: 	constanta.NexCloudResourceID,
			Status: 		"OK",
		})
	}

	//---------- Success preparation
	result = input.addResourceInfoReRegistration(result, resourceIDDataList)
	successMessage := "sukses"
	_, output = service.CustomSuccessResponsePayload(result, successMessage, contextModel)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) doUpdateToUser(tx *sql.Tx, userStruct in.UserRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	err = dao.UserDAO.ReActiveUser(tx, repository.UserModel {
		ID: 			sql.NullInt64{Int64: userStruct.ID},
		CreatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedAt: 		sql.NullTime{Time: timeNow},
		CreatedClient: 	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:		sql.NullTime{Time: timeNow},
		UpdatedClient:	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	})

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) doUpdatePKCEClientMapping(tx *sql.Tx, pkceMappingStruct in.PKCERequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {

	//------- Client alias for user
	clientAlias := pkceMappingStruct.FirstName + " SFA " + pkceMappingStruct.ClientAlias

	err = dao.PKCEClientMappingDAO.UpdatePKCEClientMappingWithClientID(tx, repository.PKCEClientMappingModel{
		ID: 			sql.NullInt64{Int64: pkceMappingStruct.PKCEClientMappingID},
		ParentClientID: sql.NullString{String: pkceMappingStruct.ParentClientID},
		CompanyID: 		sql.NullString{String: pkceMappingStruct.CompanyID},
		BranchID: 		sql.NullString{String: pkceMappingStruct.BranchID},
		ClientAlias: 	sql.NullString{String: clientAlias},
		CreatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedAt: 		sql.NullTime{Time: timeNow},
		CreatedClient: 	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:		sql.NullTime{Time: timeNow},
		UpdatedClient:	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	})

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) addResourceUserNexcloud(addResourceNexcloud in.AddResourceNexcloud, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {

	//------- Add resource to nexcloud
	err = service.AddResourceNexcloudToNexcloud(addResourceNexcloud, contextModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) addResourceInfoReRegistration(inputStruct out.PKCEResponse, resourceDataList []out.ResourceList) (result out.PKCEResponse) {
	inputStruct.ResourceList = resourceDataList
	result = inputStruct
	return
}