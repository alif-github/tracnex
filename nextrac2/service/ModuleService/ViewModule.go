package ModuleService

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

func (input moduleService) ViewModule(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ModuleRequest
	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewModule)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewModule(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input moduleService) doViewModule(inputStruct in.ModuleRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		funcName   = "doViewModule"
		moduleOnDB repository.ModuleModel
	)

	moduleOnDB, err = dao.ModuleDAO.ViewModule(serverconfig.ServerAttribute.DBConnection, repository.ModuleModel{ID: sql.NullInt64{Int64: inputStruct.ID}})
	if err.Error != nil {
		return
	}

	if moduleOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Module)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, moduleOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.convertModelToResponseDetail(moduleOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input moduleService) convertModelToResponseDetail(inputModel repository.ModuleModel) out.ModuleResponse {
	return out.ModuleResponse{
		ID:          inputModel.ID.Int64,
		ModuleName:  inputModel.ModuleName.String,
		CreatedAt:   inputModel.CreatedAt.Time,
		UpdatedAt:   inputModel.UpdatedAt.Time,
		UpdatedBy:   inputModel.UpdatedBy.Int64,
		UpdatedName: inputModel.UpdatedName.String,
	}
}

func (input moduleService) validateViewModule(inputStruct *in.ModuleRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}
