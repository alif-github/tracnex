package DepartmentService

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
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input departmentService) UpdateDepartment(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateDepartment"
		inputStruct in.DepartmentRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateDepartment)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateDepartment, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input departmentService) validateUpdateDepartment(inputStruct *in.DepartmentRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdate()
}

func (input departmentService) doUpdateDepartment(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName    = "doUpdateDepartment"
		inputStruct = inputStructInterface.(in.DepartmentRequest)
		inputModel  = input.convertDTOToModelForUpdate(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
	)

	departmentOnDB, err := dao.DepartmentDAO.GetDepartmentForUpdate(tx, repository.DepartmentModel{
		ID: inputModel.ID,
	})
	if err.Error != nil {
		return
	}

	if departmentOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Department)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, departmentOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if departmentOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.Department)
		return
	}

	if departmentOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Department)
		return
	}

	err = dao.DepartmentDAO.UpdateDepartmentByID(tx, inputModel)
	if err.Error != nil {
		err = input.checkDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.DepartmentDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	return
}

func (input departmentService) convertDTOToModelForUpdate(inputStruct in.DepartmentRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.DepartmentModel {
	return repository.DepartmentModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		Name:          sql.NullString{String: inputStruct.DepartmentName},
		Description:   sql.NullString{String: inputStruct.Description},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}
