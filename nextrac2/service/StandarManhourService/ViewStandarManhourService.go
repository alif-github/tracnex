package StandarManhourService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"strconv"
)

func (input standarManhourService) ViewStandarManhour(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.StandarManhourRequest
	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewStandarManhour)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewStandarManhour(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input standarManhourService) doViewStandarManhour(inputStruct in.StandarManhourRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		funcName   = "doViewStandarManhour"
		resultOnDB repository.StandarManhourModel
		db         = serverconfig.ServerAttribute.DBConnection
	)

	resultOnDB, err = input.StandarManhourDAO.ViewStandarManhour(db, repository.StandarManhourModel{ID: sql.NullInt64{Int64: inputStruct.ID}})
	if err.Error != nil {
		return
	}

	if resultOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.StandarManhour)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, resultOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.convertModelToResponseDetail(resultOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input standarManhourService) convertModelToResponseDetail(inputModel repository.StandarManhourModel) out.DetailStandarManhourResponse {
	return out.DetailStandarManhourResponse{
		ID:           inputModel.ID.Int64,
		Case:         inputModel.Case.String,
		DepartmentID: inputModel.DepartmentID.Int64,
		Department:   inputModel.Department.String,
		Manhour:      strconv.FormatFloat(inputModel.Manhour.Float64, 'f', -1, 64) + " jam",
		UpdatedAt:    inputModel.UpdatedAt.Time,
	}
}

func (input standarManhourService) validateViewStandarManhour(inputStruct *in.StandarManhourRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}
