package ComponentService

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
)

func (input componentService) ViewComponent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ComponentRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewComponent)
	if err.Error != nil {
		return
	}

	output.Data.Content, err =input.doViewComponent(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status =input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input componentService) doViewComponent(inputStruct in.ComponentRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	funcName := "doViewComponent"

	componentOnDB, err := dao.ComponentDAO.ViewComponent(serverconfig.ServerAttribute.DBConnection, repository.ComponentModel{
		ID:        sql.NullInt64{Int64: inputStruct.ID},
	})

	if componentOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Component)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, componentOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.convertModelToResponseDetail(componentOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input componentService) convertModelToResponseDetail(inputModel repository.ComponentModel) out.ComponentResponse {
	return out.ComponentResponse{
		ID:            inputModel.ID.Int64,
		ComponentName: inputModel.ComponentName.String,
		CreatedAt:     inputModel.CreatedAt.Time,
		UpdatedAt:     inputModel.UpdatedAt.Time,
		UpdatedBy:     inputModel.UpdatedBy.Int64,
		UpdatedName:   inputModel.UpdatedName.String,
	}
}

func (input componentService) validateViewComponent(inputStruct *in.ComponentRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}