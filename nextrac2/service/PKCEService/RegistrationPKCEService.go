package PKCEService

import (
	"database/sql"
	"net/http"
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
	"nexsoft.co.id/nextrac2/util"
	"strings"
	"time"
)

func (input pkceService) RegistrationUserPKCE(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "RegistrationUserPKCE"
	var registerUser interface{}

	inputStruct, err := input.readBodyAndValidateRegisUnregisPKCE(request, contextModel, true, input.validateRegistrationPKCE)
	if err.Error != nil {
		return
	}

	registerUser, err = input.InsertServiceWithAuditCustom(funcName, inputStruct, contextModel, input.doRegistrationUserPKCE, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output = registerUser.(out.Payload)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) doRegistrationUserPKCE(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, isReRegistrationPKCE bool, err errorModel.ErrorModel) {
	fileName := "RegistrationPKCEService.go"
	funcName := "doRegistrationPKCEService"
	var detail string
	var resourceIDDataList []out.ResourceList
	var result interface{}
	var checkUnregisteredModel repository.CheckPKCEClientMappingModel
	inputStruct := inputStructInterface.(in.PKCERequest)
	isNexmile := inputStruct.ClientTypeID == constanta.ResourceNexmileID

	//---------- Check validity of client mapping
	var newInputStructInterface interface{}
	if isNexmile {
		newInputStructInterface, err = input.checkClientMappingValid(tx, inputStructInterface)
		if err.Error != nil {
			return
		}

		inputStruct = newInputStructInterface.(in.PKCERequest)

		if (inputStruct == in.PKCERequest{}) {
			detail = util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.PKCEUserBundle, "DETAIL_ERROR_INVALID_ID_MESSAGE", contextModel.AuthAccessTokenModel.Locale, nil)
			err = errorModel.GenerateInvalidRegistrationPKCE(fileName, funcName, []string{detail})
			return
		}
	}

	//---------- Check if user re registration
	if isNexmile {
		checkUnregisteredModel, err = input.isUnregisterBefore(inputStruct)
		if err.Error != nil {
			return
		}
	}

	//---------- Hit authentication user to auth server
	result, err = input.HitUserRegistrationToAuthServer(inputStruct, contextModel)
	if err.Error != nil {

		//---------- Re-registering user pkce nexmile
		if err.Error.Error() == "E-4-AUT-SRV-002" && checkUnregisteredModel.IsRegisteredBefore.Bool && isNexmile {
			//isReRegistrationPKCE = true //---------- Status for re-registration pkce
			//
			//copyPKCERequest := in.PKCEReRequest {
			//	ID: 					checkUnregisteredModel.ID.Int64,
			//	ClientID: 				checkUnregisteredModel.ClientID.String,
			//	AuthUserID: 			checkUnregisteredModel.AuthUserID.Int64,
			//	PKCEClientMappingID:	checkUnregisteredModel.PKCEClientMappingID.Int64,
			//	PKCERequest: 			inputStruct,
			//}

			//---------- Do re-registering user pkce nexmile
			//output, dataAudit, err = input.doReRegistrationUserPKCE(tx, copyPKCERequest, contextModel, timeNow)
			return
		}

		return
	}

	responseAuth := result.(out.PKCEResponse)

	//---------- Add resource nextrac
	_, err = input.addResourceNextracToAuthServer(responseAuth, contextModel)
	if err.Error != nil {
		return
	} else {
		resourceIDDataList = append(resourceIDDataList, out.ResourceList {
			ResourceID: 	config.ApplicationConfiguration.GetServerResourceID(),
			Status: 		"OK",
		})
	}

	//---------- Insert new registered user nexmile to user
	idUser, idClientRoleScope, err := input.doInsertToUser(tx, inputStruct, responseAuth, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	//---------- Insert new registered user to client mapping
	idPKCEClientMapping, err := input.doInsertPKCEClientMapping(tx, inputStruct, responseAuth, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName: 	sql.NullString{String: dao.UserDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idUser},
	}, repository.AuditSystemModel{
		TableName: 	sql.NullString{String: dao.ClientRoleScopeDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idClientRoleScope},
	}, repository.AuditSystemModel{
		TableName: 	sql.NullString{String: dao.PKCEClientMappingDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idPKCEClientMapping},
	})

	//---------- Add resource nexcloud
	err = input.addResourceUserNexcloudToNexcloudServer(responseAuth, inputStruct, contextModel)
	if err.Error != nil {
		_, err = service.NewErrorAddResource(serverconfig.ServerAttribute.PKCEUserBundle, contextModel, "DETAIL_ERROR_FAILED_RESOURCE_NEXCLOUD_MESSAGE", fileName, funcName)
		resourceIDDataList = append(resourceIDDataList, out.ResourceList {
			ResourceID: 	constanta.NexCloudResourceID,
			Status: 		"FAIL",
		})
		output, err = input.regisPKCEErrorHandle(resourceIDDataList, responseAuth, inputStruct, contextModel, err, timeNow)
		return
	} else {
		resourceIDDataList = append(resourceIDDataList, out.ResourceList {
			ResourceID: 	constanta.NexCloudResourceID,
			Status: 		"OK",
		})
	}

	output, err = input.regisPKCESuccessHandle(responseAuth, contextModel, inputStruct, resourceIDDataList, timeNow)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) checkResourceAfterError(fileName string, funcName string, pkceResponse out.PKCEResponse, messageID string, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var detail string

	detail = util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.PKCEUserBundle, messageID, contextModel.AuthAccessTokenModel.Locale, nil)
	err = errorModel.GenerateAuthenticationServerAddResourceError(fileName, funcName, []string{detail})
	resourceIDList, errS := input.checkUserToAuthServer(pkceResponse, contextModel)
	if errS.Error != nil {
		err = errS
		return
	}

	pkceResponse.ResourceList = resourceIDList
	result = pkceResponse
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) doInsertToUser(tx *sql.Tx, pkceRequest in.PKCERequest, pkceResponse out.PKCEResponse, contextModel *applicationModel.ContextModel, timeNow time.Time) (idUser int64, idClientRoleScope int64, err errorModel.ErrorModel)  {
	var roleID int64
	idUser, err = dao.UserDAO.InsertUser(tx, repository.UserModel{
		ClientID: 		sql.NullString{String: pkceResponse.ClientID},
		AuthUserID:		sql.NullInt64{Int64: pkceResponse.UserID},
		Locale: 		sql.NullString{String: constanta.IndonesianLanguage},
		FirstName: 		sql.NullString{String: pkceRequest.FirstName},
		LastName: 		sql.NullString{String: pkceRequest.LastName},
		Email: 			sql.NullString{String: pkceRequest.Email},
		Username: 		sql.NullString{String: pkceRequest.Username},
		Phone: 			sql.NullString{String: constanta.IndonesianCodeNumber+"-"+pkceRequest.Phone},
		Status: 		sql.NullString{String: constanta.PendingOnApproval},
		CreatedBy:		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt: 		sql.NullTime{Time: timeNow},
		UpdatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: 	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt: 		sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		return
	}

	if pkceRequest.ClientTypeID == constanta.ResourceNexmileID {
		roleID = constanta.RoleUserNexMile
	}

	idClientRoleScope, err = dao.ClientRoleScopeDAO.InsertClientRoleScope(tx, repository.ClientRoleScopeModel{
		ClientID: 		sql.NullString{String: pkceResponse.ClientID},
		RoleID: 		sql.NullInt64{Int64: roleID},
		CreatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: 	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt: 		sql.NullTime{Time: timeNow},
		UpdatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: 	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt: 		sql.NullTime{Time: timeNow},
	})
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) doInsertPKCEClientMapping(tx *sql.Tx, pkceRequest in.PKCERequest, pkceResponse out.PKCEResponse, contextModel *applicationModel.ContextModel, timeNow time.Time) (idClientMapping int64, err errorModel.ErrorModel) {
	var clientAlias string
	var isClientDependant bool

	//------- Client alias for user
	if pkceRequest.ClientTypeID == constanta.ResourceNexmileID {
		clientAlias = pkceRequest.FirstName + " SFA " + pkceRequest.ClientAlias
		isClientDependant = true
	} else {
		clientAlias = pkceRequest.FirstName + pkceRequest.LastName
		isClientDependant = false
	}

	idClientMapping, err = dao.PKCEClientMappingDAO.InsertPKCEClientMapping(tx, &repository.PKCEClientMappingModel {
		ParentClientID: sql.NullString{String: pkceRequest.ParentClientID},
		ClientID: 		sql.NullString{String: pkceResponse.ClientID},
		ClientAlias: 	sql.NullString{String: clientAlias},
		AuthUserID: 	sql.NullInt64{Int64: pkceResponse.UserID},
		Username: 		sql.NullString{String: pkceRequest.Username},
		CompanyID: 		sql.NullString{String: pkceRequest.CompanyID},
		BranchID: 		sql.NullString{String: pkceRequest.BranchID},
		ClientTypeID: 	sql.NullInt64{Int64: pkceRequest.ClientTypeID},
		CreatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: 	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt: 		sql.NullTime{Time: timeNow},
		UpdatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: 	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt: 		sql.NullTime{Time: timeNow},
	}, isClientDependant)

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) splitResourceFromAuth(resourceID string) (resourceList []out.ResourceList) {
	splitResourceID := strings.Split(resourceID, " ")
	for _, resourceIDElm := range splitResourceID {
		resourceList = append(resourceList, out.ResourceList{
			ResourceID: resourceIDElm,
		})
	}
	return
}

func (input pkceService) validateRegistrationPKCE(inputStruct *in.PKCERequest) errorModel.ErrorModel {
	return inputStruct.ValidateRegistrationPKCE()
}
