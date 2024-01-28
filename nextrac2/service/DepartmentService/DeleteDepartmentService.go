package DepartmentService

import (
	"database/sql"
	"github.com/google/uuid"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input departmentService) DeleteDepartment(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteDepartment"
		inputStruct in.DepartmentRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDeleteDepartment)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteDepartment, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input departmentService) validateDeleteDepartment(inputStruct *in.DepartmentRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}

func (input departmentService) doDeleteDepartment(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct    = inputStructInterface.(in.DepartmentRequest)
		inputModel     = input.convertDTOToModelForDelete(inputStruct, contextModel, timeNow)
		departmentOnDB repository.DepartmentModel
	)

	// Lock Sementara
	if inputModel.ID.Int64 == 1 || inputModel.ID.Int64 == 2 {
		err = errorModel.GenerateSimpleErrorModel(400, "Data ini tidak diizinkan untuk dihapus")
		return
	}

	// Validate Module Check DB
	departmentOnDB, err = input.validateDepartmentOnDB(tx, inputStruct, inputModel, contextModel)
	if err.Error != nil {
		return
	}

	// Create Random Module DB
	input.randTokenGenerator(&inputModel, departmentOnDB)

	// Delete Module
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.DepartmentDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	err = dao.DepartmentDAO.DeleteDepartment(tx, inputModel)

	return
}

func (input departmentService) convertDTOToModelForDelete(inputStruct in.DepartmentRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) repository.DepartmentModel {
	return repository.DepartmentModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}
}

func (input departmentService) validateDepartmentOnDB(tx *sql.Tx, inputStruct in.DepartmentRequest, inputModel repository.DepartmentModel, contextModel *applicationModel.ContextModel) (departmentOnDB repository.DepartmentModel, err errorModel.ErrorModel) {
	var (
		funcName = "validateDepartmentOnDB"
	)

	departmentOnDB, err = dao.DepartmentDAO.GetDepartmentForUpdate(tx, repository.DepartmentModel{ID: inputModel.ID})
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

	return
}

func (input departmentService) randTokenGenerator(inputModel *repository.DepartmentModel, modelOnDB repository.DepartmentModel) {
	encodedStr := service.RandTimeToken(constanta.RandTokenForDeleteLength, uuid.New().String())
	inputModel.Name.String = modelOnDB.Name.String + encodedStr
}
