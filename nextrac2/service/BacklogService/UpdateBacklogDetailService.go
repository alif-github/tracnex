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
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input backlogService) UpdateDetailBacklog(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateDetailBacklog"
		inputStruct in.BacklogRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateDetailBacklog)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateDetailBacklog, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input backlogService) validateUpdateDetailBacklog(inputStruct *in.BacklogRequest) (err errorModel.ErrorModel) {
	return inputStruct.ValidateUpdate()
}

func (input backlogService) doUpdateDetailBacklog(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct          = inputStructInterface.(in.BacklogRequest)
		model                = input.convertToDTOToInputModelForUpdate(inputStruct, *contextModel, timeNow)
		newUpdatedFileUpload int64
		scope                map[string]interface{}
	)

	mappingScopeDB := make(map[string]applicationModel.MappingScopeDB)
	mappingScopeDB[constanta.EmployeeDataScope] = applicationModel.MappingScopeDB{
		View:  "e.id",
		Count: "e.id",
	}

	createdBy := contextModel.LimitedByCreatedBy       //--- Add userID when have own permission
	scope, err = input.validateDataScope(contextModel) //--- Get scope
	if err.Error != nil {
		return
	}

	//--- Check and Validate PIC On DB
	err = input.checkPICOnDB(inputStruct, createdBy, scope, mappingScopeDB)
	if err.Error != nil {
		return
	}

	//--- Check and Validate Backlog On DB
	err = input.checkAndLocDetailBacklogOnDB(inputStruct, model, contextModel)
	if err.Error != nil {
		return
	}

	// handle Form Perubahan
	listBacklogUpdate := []*in.BacklogRequest{
		&inputStruct,
	}

	newUpdatedFileUpload, err = input.HandleFile(tx, listBacklogUpdate, contextModel)
	if err.Error != nil {
		return
	}

	// update file perubahan jika ada
	model.FileUploadId.Int64 = newUpdatedFileUpload

	//--- Update On DB
	dataAudit, err = input.updateDetailBacklogOnDB(tx, contextModel, timeNow, model)
	return
}

func (input backlogService) convertToDTOToInputModelForUpdate(inputStruct in.BacklogRequest, contextModel applicationModel.ContextModel, timeNow time.Time) repository.BacklogModel {
	return repository.BacklogModel{
		ID:              sql.NullInt64{Int64: inputStruct.ID},
		Layer1:          sql.NullString{String: inputStruct.Layer1},
		Layer2:          sql.NullString{String: inputStruct.Layer2},
		Layer3:          sql.NullString{String: inputStruct.Layer3},
		Layer4:          sql.NullString{String: inputStruct.Layer4},
		Layer5:          sql.NullString{String: inputStruct.Layer5},
		Feature:         sql.NullInt64{Int64: inputStruct.Feature},
		Subject:         sql.NullString{String: inputStruct.Subject},
		ReferenceTicket: sql.NullInt64{Int64: inputStruct.ReferenceTicket},
		RedmineNumber:   sql.NullInt64{Int64: inputStruct.RedmineNumber},
		Sprint:          sql.NullString{String: inputStruct.Sprint},
		SprintName:      sql.NullString{String: inputStruct.SprintName},
		EmployeeId:      sql.NullInt64{Int64: inputStruct.PicId},
		Status:          sql.NullString{String: inputStruct.Status},
		Mandays:         sql.NullFloat64{Float64: inputStruct.Mandays},
		EstimateTime:    sql.NullFloat64{Float64: inputStruct.MandaysDone},
		FlowChanged:     sql.NullString{String: inputStruct.FlowChanged},
		AdditionalData:  sql.NullString{String: inputStruct.AdditionalData},
		Note:            sql.NullString{String: inputStruct.Note},
		Url:             sql.NullString{String: inputStruct.Url},
		Page:            sql.NullString{String: inputStruct.Page},
		Tracker:         sql.NullString{String: inputStruct.Tracker},
		UpdatedBy:       sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:       sql.NullTime{Time: timeNow},
		UpdatedClient:   sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}
}

func (input backlogService) checkAndLocDetailBacklogOnDB(inputStruct in.BacklogRequest, model repository.BacklogModel, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		funcName    = "checkAndLocDetailBacklogOnDB"
		db          = serverconfig.ServerAttribute.DBConnection
		backlogOnDB repository.BacklogModel
	)

	backlogOnDB, err = dao.BacklogDAO.GetDetailBacklogForUpdateOrDelete(db, model)
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

	if backlogOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.BacklogConstanta)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input backlogService) checkPICOnDB(inputStruct in.BacklogRequest, createdBy int64, scopeLimit map[string]interface{}, scopeDB map[string]applicationModel.MappingScopeDB) (err errorModel.ErrorModel) {
	var (
		funcName     = "checkPICOnDB"
		db           = serverconfig.ServerAttribute.DBConnection
		employeeOnDB repository.EmployeeModel
	)

	employeeOnDB, err = dao.EmployeeDAO.ViewEmployee(db, repository.EmployeeModel{ID: sql.NullInt64{Int64: inputStruct.PicId}}, createdBy, scopeLimit, scopeDB)
	if err.Error != nil {
		return
	}

	if employeeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.EmployeeID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input backlogService) updateDetailBacklogOnDB(tx *sql.Tx, contextModel *applicationModel.ContextModel, timeNow time.Time, model repository.BacklogModel) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.BacklogDAO.TableName, model.ID.Int64, contextModel.LimitedByCreatedBy)...)
	err = dao.BacklogDAO.UpdateBacklog(tx, model)
	if err.Error != nil {
		err = input.checkDuplicateError(err)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
