package ReportService

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"strconv"
)

func (input reportService) ViewReportHistory(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ReportHistory
	inputStruct, err = input.readBodyAndValidateReportHistory(request, contextModel, input.validateViewReportHistory)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewReportHistory(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input reportService) doViewReportHistory(inputStruct in.ReportHistory, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		funcName          = "doViewReportHistory"
		db                = serverconfig.ServerAttribute.DBConnection
		reportHistoryOnDB repository.ReportHistoryModel
	)

	reportHistoryOnDB, err = dao.ReportHistoryDAO.ViewReportHistory(db, repository.ReportHistoryModel{ID: sql.NullInt64{Int64: inputStruct.ID}})
	if err.Error != nil {
		return
	}

	if reportHistoryOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ReportConstanta)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, reportHistoryOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.convertModelToResponseDetail(reportHistoryOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) convertModelToResponseDetail(inputModel repository.ReportHistoryModel) out.ViewReportHistoryResponse {
	var data out.ResultsReportResponse
	_ = json.Unmarshal([]byte(inputModel.Data.String), &data)
	return out.ViewReportHistoryResponse{
		ID:                inputModel.ID.Int64,
		Department:        inputModel.DepartmentName.String,
		SuccessTicket:     inputModel.SuccessTicket.String,
		PaymentDate:       inputModel.CreatedAt.Time,
		PersonResponsible: inputModel.CreatedName.String,
		Data:              data,
	}
}

func (input reportService) validateViewReportHistory(inputStruct *in.ReportHistory) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}

func (input reportService) readBodyAndValidateReportHistory(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ReportHistory) errorModel.ErrorModel) (inputStruct in.ReportHistory, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidateReportHistory"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}
