package HitAddExternalResource

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input hitAddExternalResourceService) InsertAddResourceNexcloudService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	fileName := "InsertAddResourceNexcloudService.go"
	funcName := "InsertAddResourceNexcloudService"

	inputStruct, err := input.readBodyAndValidate(request, contextModel, input.validateInsertAddResource)
	if err.Error != nil {
		return
	}

	//--- Check client type
	//clientTypeOnDB, err := dao.ClientTypeDAO.ValidateClientTypeByID(serverconfig.ServerAttribute.DBConnection, repository.ClientTypeModel{ID: sql.NullInt64{Int64: inputStruct.ClientTypeID}})
	//if err.Error != nil {
	//	return
	//}
	//
	//if clientTypeOnDB.ID.Int64 < 1 {
	//	err = errorModel.GenerateEmptyFieldError(input.FileName, funcName, constanta.ClientTypeID)
	//	return
	//}

	if inputStruct.ClientTypeID == constanta.ResourceND6ID{
		err = input.doAddClientMappingResourceNexcloudService(inputStruct, contextModel)
	}else if inputStruct.ClientTypeID == constanta.ResourceNexmileID{
		err = input.doAddPKCEClientMappingResourceNexcloudService(inputStruct, contextModel)
	}else {
		err = errorModel.GenerateUnsupportedResponseTypeError(fileName, funcName)
	}

	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse {
		Code: 		util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: 	GenerateI18NMessage("SUCCESS_ADD_RESOURCE_NEXCLOUD_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}