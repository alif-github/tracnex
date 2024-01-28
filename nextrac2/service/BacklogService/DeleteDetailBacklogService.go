package BacklogService

import (
	"database/sql"
	"math/rand"
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
	"strconv"
	"time"
)

func (input backlogService) DeleteDetailBacklog(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteDetailBacklog"
		inputStruct in.BacklogRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDelete)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteBacklog, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input backlogService) validateDelete(inputStruct *in.BacklogRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}

func (input backlogService) doDeleteBacklog(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct = inputStructInterface.(in.BacklogRequest)
	)

	inputModel := repository.BacklogModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	backlogOnDB, err := input.validationBacklogOnDB(tx, inputStruct, inputModel, contextModel)
	if err.Error != nil {
		return
	}

	// assign redmine number with random string
	inputModel.RedmineNumber = backlogOnDB.RedmineNumber

	// delete backlog
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.BacklogDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.BacklogDAO.DeleteBacklog(tx, inputModel)
	if err.Error != nil {
		return
	}

	return
}

func randomNumber2Digit() int64 {
	// Inisialisasi generator angka acak dengan seed yang berubah
	rand.Seed(time.Now().UnixNano())

	// Menghasilkan nilai acak antara 10 dan 99
	return int64(rand.Intn(90) + 10)
}

func (input backlogService) validationBacklogOnDB(tx *sql.Tx, inputStruct in.BacklogRequest, inputModel repository.BacklogModel, contextModel *applicationModel.ContextModel) (backlogOnDB repository.BacklogModel, err errorModel.ErrorModel) {
	var (
		funcName = "validationBacklogOnDB"
	)

	backlogOnDB, err = dao.BacklogDAO.GetDetailBacklogForUpdateOrDelete(serverconfig.ServerAttribute.DBConnection, repository.BacklogModel{ID: inputModel.ID})
	if err.Error != nil {
		return
	}

	if backlogOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
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

	// generate random number
	redmineNumberString := strconv.FormatInt(backlogOnDB.RedmineNumber.Int64, 10)
	randomString2Digit := strconv.FormatInt(randomNumber2Digit(), 10)

	redmineNumberNewString := redmineNumberString + randomString2Digit
	redmineNumberNewInt, errs := strconv.ParseInt(redmineNumberNewString, 10, 36)
	if errs != nil {
		err = errorModel.GenerateSimpleErrorModel(400, "Error parsing unik 2 digit")
		return
	}

	backlogOnDB.RedmineNumber.Int64 = redmineNumberNewInt

	err = errorModel.GenerateNonErrorModel()
	return
}
