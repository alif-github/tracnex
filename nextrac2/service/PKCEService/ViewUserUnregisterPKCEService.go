package PKCEService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input pkceService) ViewUserPKCEForUnregister(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	inputStruct, err := input.readUrlPathUnregisPKCE(request, input.validateViewUserPKCEForUnregister)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewUserPKCEForUnregister(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse {
		Code: 		util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message:	GenerateI18NMessage("SUCCESS_VIEW_CUSTOM_PKCE_UNREGIS_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) doViewUserPKCEForUnregister(inputStruct in.PKCERequest, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	fileName := "ViewUserUnregisterPKCEService.go"
	funcName := "doViewUserPKCEForUnregister"

	userParam := repository.PKCEClientMappingModel {
		Username: 	sql.NullString{String: inputStruct.Username},
		CreatedBy:	sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	}

	resultData, err := dao.PKCEClientMappingDAO.ViewForUnregisterUserPKCE(serverconfig.ServerAttribute.DBConnection, userParam)
	if err.Error != nil {
		return
	}

	if resultData.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.User)
		return
	}

	output = input.convertToDTOForViewUnregister(resultData)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) convertToDTOForViewUnregister(dataModel repository.PKCEClientMappingModel) out.ViewPKCEResponse {
	return out.ViewPKCEResponse {
		UserID: 			dataModel.ID.Int64,
		ParentClientID:		dataModel.ParentClientID.String,
		ClientID: 			dataModel.ClientID.String,
		Username:			dataModel.Username.String,
		CreatedBy: 			dataModel.CreatedBy.Int64,
		UpdatedAt: 			dataModel.UpdatedAt.Time,
	}
}

func (input pkceService) validateViewUserPKCEForUnregister(inputStruct *in.PKCERequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewForUnregisterPKCE()
}