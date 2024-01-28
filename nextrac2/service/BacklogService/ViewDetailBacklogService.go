package BacklogService

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
	"strconv"
)

func (input backlogService) ViewDetailBacklog(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.BacklogRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewDetailBacklog)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewBacklog(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input backlogService) validateViewDetailBacklog(inputStruct *in.BacklogRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}

func (input backlogService) doViewBacklog(inputStruct in.BacklogRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		funcName    = "doViewBacklog"
		db          = serverconfig.ServerAttribute.DBConnection
		backlogOnDB repository.BacklogModel
	)

	backlogOnDB, err = dao.BacklogDAO.ViewDetailBacklog(db, repository.BacklogModel{ID: sql.NullInt64{Int64: inputStruct.ID}})
	if err.Error != nil {
		return
	}

	if backlogOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.BacklogConstanta)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, backlogOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.convertModelToResponseDetail(backlogOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input backlogService) convertModelToResponseDetail(backlogOnDB repository.BacklogModel) out.ViewDetailBacklogResponse {
	return out.ViewDetailBacklogResponse{
		ID:              backlogOnDB.ID.Int64,
		Layer1:          backlogOnDB.Layer1.String,
		Layer2:          backlogOnDB.Layer2.String,
		Layer3:          backlogOnDB.Layer3.String,
		Layer4:          backlogOnDB.Layer4.String,
		Layer5:          backlogOnDB.Layer5.String,
		Subject:         backlogOnDB.Subject.String,
		Tracker:         backlogOnDB.Tracker.String,
		Feature:         backlogOnDB.Feature.Int64,
		RedmineNumber:   strconv.FormatInt(backlogOnDB.RedmineNumber.Int64, 10),
		Sprint:          backlogOnDB.Sprint.String,
		SprintName:      backlogOnDB.SprintName.String,
		ReferenceTicket: backlogOnDB.ReferenceTicket.Int64,
		Description:     backlogOnDB.Description.String,
		PicId:           backlogOnDB.EmployeeId.Int64,
		Pic:             backlogOnDB.EmployeeName.String,
		Status:          backlogOnDB.Status.String,
		Mandays:         backlogOnDB.Mandays.Float64,
		MandaysDone:     backlogOnDB.EstimateTime.Float64,
		FlowChanged:     backlogOnDB.FlowChanged.String,
		AdditionalData:  backlogOnDB.AdditionalData.String,
		Note:            backlogOnDB.Note.String,
		Url:             backlogOnDB.Url.String,
		Page:            backlogOnDB.Page.String,
		DepartmentId:    backlogOnDB.DepartmentId.Int64,
		DepartmentName:  backlogOnDB.DepartmentName.String,
		UrlFile:         backlogOnDB.FileUploadData.Host.String + backlogOnDB.FileUploadData.Path.String + backlogOnDB.FileUploadData.FileName.String,
		UpdatedAt:       backlogOnDB.UpdatedAt.Time,
		CreatedName:     backlogOnDB.CreatedName.String,
		UpdatedName:     backlogOnDB.UpdatedName.String,
		CreatedBy:       backlogOnDB.CreatedBy.Int64,
		CreatedAt:       backlogOnDB.CreatedAt.Time,
	}
}
