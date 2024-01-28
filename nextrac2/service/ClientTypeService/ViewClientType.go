package ClientTypeService

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
	"nexsoft.co.id/nextrac2/util"
)

func (input clientTypeService) ViewClientType(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ClientTypeRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewClientType)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewClientType(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_CLIENT_TYPE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input clientTypeService) doViewClientType(inputStruct in.ClientTypeRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		funcName = "doViewClientType"
	)

	clientTypeOnDB, err := dao.ClientTypeDAO.ViewClientType(serverconfig.ServerAttribute.DBConnection, repository.ClientTypeModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	})

	if clientTypeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientType)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, clientTypeOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.convertDTOToModelForView(clientTypeOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeService) validateViewClientType(inputStruct *in.ClientTypeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}
