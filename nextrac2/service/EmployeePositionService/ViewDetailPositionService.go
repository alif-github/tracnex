package EmployeePositionService

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

func (input employeePositionService) ViewDetailEmployeePosition(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.EmployeePosition
	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewEmployeePosition)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewDetailEmployeePosition(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input employeePositionService) doViewDetailEmployeePosition(inputStruct in.EmployeePosition, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		funcName     = "doViewModule"
		db           = serverconfig.ServerAttribute.DBConnection
		positionOnDB repository.EmployeePositionModel
	)

	//--- View Position
	positionOnDB, err = dao.EmployeePositionDAO.ViewEmployeePosition(db, repository.EmployeePositionModel{ID: sql.NullInt64{Int64: inputStruct.ID}})
	if err.Error != nil {
		return
	}

	if positionOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Position)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, positionOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.convertModelToResponseDetail(positionOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeePositionService) convertModelToResponseDetail(inputModel repository.EmployeePositionModel) out.DetailEmployeePosition {
	return out.DetailEmployeePosition{
		ID:          inputModel.ID.Int64,
		Name:        inputModel.Name.String,
		Description: inputModel.Description.String,
		CompanyName: inputModel.CompanyName.String,
		CreatedAt:   inputModel.CreatedAt.Time,
		UpdatedAt:   inputModel.UpdatedAt.Time,
		CreatedBy:   inputModel.CreatedBy.Int64,
		CreatedName: inputModel.CreatedName.String,
		UpdatedBy:   inputModel.UpdatedBy.Int64,
		UpdatedName: inputModel.UpdatedName.String,
	}
}

func (input employeePositionService) validateViewEmployeePosition(inputStruct *in.EmployeePosition) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}
