package HitAddExternalResource

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
)

func (input hitAddExternalResourceService) InsertAddResourceNexdriveService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	fileName := "InsertAddResourceNexdriveService.go"
	funcName := "InsertAddResourceNexdriveService"

	inputStruct, err := input.readBodyAndValidate(request, contextModel, input.validateInsertAddResource)
	if err.Error != nil {
		return
	}

	//forbidden when client id request different with client id logged in
	if inputStruct.ClientID != contextModel.AuthAccessTokenModel.ClientID {
		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		return
	}

	err = input.doInsertAddResourceNexDriveService(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse {
		Code: 		util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: 	GenerateI18NMessage("SUCCESS_ADD_RESOURCE_NEXDRIVE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hitAddExternalResourceService) doInsertAddResourceNexDriveService(inputStructInterface interface{}, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	fileName := "InsertAddResourceNexdriveService.go"
	funcName := "doInsertAddResourceNexdriveService"

	inputStruct := inputStructInterface.(in.AddResourceExternalRequest)
	var firstName string
	var errTemp errorModel.ErrorModel

	//---- Check resource exist
	if strings.Contains(inputStruct.OldResource, constanta.NexdriveResourceID) {
		detail := util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.AddResourceExternalBundle, "FAILED_ADD_RESOURCE_EXISTING", contextModel.AuthAccessTokenModel.Locale, nil)
		err = errorModel.GenerateInvalidAddResourceNexdrive(fileName, funcName, []string{detail})
		return
	}

	//---- Get FirstName from user
	firstName, err = input.doGetFirstNameUser(inputStruct)
	if err.Error != nil {
		return
	}

	//---- Add resource nexdrive
	err = input.AddResourceNexdrive(in.AddResourceNexdrive {
		FirstName: 	firstName,
		LastName: 	constanta.ND6,
		ClientID: 	inputStruct.ClientID,
	}, contextModel)

	if err.Error != nil {
		//todo ganti sesuai dengan status error nexdrive
		if err.Error.Error() == constanta.ErrorClientIDExistNexcloud {
			return
		}

		//---- Error must be update to registration log
		errTemp = input.updateRegistrationLogWithAudit(inputStruct, constanta.NexdriveResourceID, contextModel, err)
		if errTemp.Error != nil {
			return
		}

		return
	}

	//---- Success must be update to registration log
	err = input.updateRegistrationLogWithAudit(inputStruct, constanta.NexdriveResourceID, contextModel, err)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}