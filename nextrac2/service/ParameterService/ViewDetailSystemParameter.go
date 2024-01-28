package ParameterService

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
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input parameterService) ViewParameter(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ParameterRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewParameter(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input parameterService) doViewParameter(inputStruct in.ParameterRequest, contextModel *applicationModel.ContextModel) (result out.ViewDetailParameterDTOOut, err errorModel.ErrorModel) {
	funcName := "doViewParameter"
	Parameter := repository.ParameterModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	}

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*contextModel)
	if isOnlyHaveOwnAccess {
		Parameter.CreatedBy.Int64 = userID
	}

	Parameter, err = dao.ParameterDAO.ViewParameter(serverconfig.ServerAttribute.DBConnection, Parameter)
	if err.Error != nil {
		return
	}

	if Parameter.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	result = reformatDAOtoDTO(Parameter)
	return
}

func reformatDAOtoDTO(ownerModel repository.ParameterModel) out.ViewDetailParameterDTOOut {
	return out.ViewDetailParameterDTOOut{
		ID:          ownerModel.ID.Int64,
		Permission:  ownerModel.Permission.String,
		Name:        ownerModel.Name.String,
		Value:       ownerModel.Value.String,
		Description: ownerModel.Description.String,
		CreatedBy:   ownerModel.CreatedBy.Int64,
		UpdatedAt:   ownerModel.UpdatedAt.Time,
	}
}

func (input parameterService) validateView(inputStruct *in.ParameterRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewParameter()
}
