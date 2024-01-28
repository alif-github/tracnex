package HitAddExternalResource

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
)

func (input hitAddExternalResourceService) doAddClientMappingResourceNexcloudService(inputStructInterface interface{}, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	fileName := "AddResourceNexcloudChannel.go"
	funcName := "doAddClientMappingResourceNexcloudService"

	inputStruct := inputStructInterface.(in.AddResourceExternalRequest)
	var firstName string
	var errTemp errorModel.ErrorModel

	//---------- forbidden when client id request different with client id logged in
	if inputStruct.ClientID != contextModel.AuthAccessTokenModel.ClientID {
		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		return
	}

	//---- Check client type
	//err = input.checkClientTypeByID(&inputStruct)
	//if err.Error != nil {
	//	return
	//}

	//---------- Check resource exist
	if strings.Contains(inputStruct.OldResource, constanta.NexCloudResourceID) {
		detail := util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.AddResourceExternalBundle, "FAILED_ADD_RESOURCE_EXISTING", contextModel.AuthAccessTokenModel.Locale, nil)
		err = errorModel.GenerateInvalidAddResourceNexcloud(fileName, funcName, []string{detail})
		return
	}

	//---------- Get FirstName from user
	firstName, err = input.doGetFirstNameUser(inputStruct)
	if err.Error != nil {
		return
	}

	//---------- Add resource nexcloud
	err = input.AddResourceNexcloud(in.AddResourceNexcloud {
		FirstName: 	firstName,
		LastName: 	constanta.Nexdistribution,
		ClientID: 	inputStruct.ClientID,
	}, contextModel)

	if err.Error != nil {
		if err.Error.Error() == constanta.ErrorClientIDExistNexcloud {
			return
		}

		//---------- Error must be update to registration log
		errTemp = input.updateRegistrationLogWithAudit(inputStruct, constanta.NexCloudResourceID, contextModel, err)
		if errTemp.Error != nil {
			return
		}

		return
	}

	//---------- Success must be update to registration log
	err = input.updateRegistrationLogWithAudit(inputStruct, constanta.NexCloudResourceID, contextModel, err)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hitAddExternalResourceService) doAddPKCEClientMappingResourceNexcloudService(inputStructInterface interface{}, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	fileName := "AddResourceNexcloudChannel.go"
	funcName := "doAddPKCEClientMappingResourceNexcloudService"

	inputStruct := inputStructInterface.(in.AddResourceExternalRequest)
	var errTemp errorModel.ErrorModel


	//---------- Check parent Client ID
	pkceClientMappingOnDB, err := dao.PKCEClientMappingDAO.GetPKCEClientMappingForAddResource(serverconfig.ServerAttribute.DBConnection, repository.PKCEClientMappingModel{
		ClientID:       sql.NullString{String: inputStruct.ClientID},
		ParentClientID: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	})

	if err.Error != nil {
		return
	}

	if pkceClientMappingOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		return
	}

	//---- Check client type
	//err = input.checkClientTypeByID(&inputStruct)
	//if err.Error != nil {
	//	return
	//}

	//---------- Check resource exist
	if strings.Contains(inputStruct.OldResource, constanta.NexCloudResourceID) {
		detail := util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.AddResourceExternalBundle, "FAILED_ADD_RESOURCE_EXISTING", contextModel.AuthAccessTokenModel.Locale, nil)
		err = errorModel.GenerateInvalidAddResourceNexcloud(fileName, funcName, []string{detail})
		return
	}

	//---------- Get FirstName from user
	userOnDB, err := dao.UserDAO.GetUserByClientID(serverconfig.ServerAttribute.DBConnection, repository.UserModel{ClientID: sql.NullString{String: inputStruct.ClientID}})
	if err.Error != nil {
		return
	}

	if userOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.User)
		return
	}

	//---------- Add resource nexcloud
	err = input.AddResourceNexcloud(in.AddResourceNexcloud {
		FirstName: 	userOnDB.FirstName.String,
		LastName: 	userOnDB.LastName.String,
		ClientID: 	inputStruct.ClientID,
	}, contextModel)

	if err.Error != nil {
		if err.Error.Error() == constanta.ErrorClientIDExistNexcloud {
			return
		}

		//---------- Error must be update to registration log
		errTemp = input.updateRegistrationLogWithAudit(inputStruct, constanta.NexCloudResourceID, contextModel, err)
		if errTemp.Error != nil {
			return
		}

		return
	}

	//---------- Success must be update to registration log
	err = input.updateRegistrationLogWithAudit(inputStruct, constanta.NexCloudResourceID, contextModel, err)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
