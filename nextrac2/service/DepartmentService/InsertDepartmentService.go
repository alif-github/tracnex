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

func (input departmentService) InsertDepartment(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertDepartment"
		inputStruct in.DepartmentRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertDepartment, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input departmentService) validateInsert(inputStruct *in.DepartmentRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsert()
}

func (input departmentService) doInsertDepartment(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct = inputStructInterface.(in.DepartmentRequest)
		inputModel  = input.convertDTOToModel(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
	)

	insertedID, err := dao.DepartmentDAO.InsertDepartment(tx, inputModel)
	if err.Error != nil {
		err = input.checkDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.DepartmentDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: insertedID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input departmentService) convertDTOToModel(inputStruct in.DepartmentRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.DepartmentModel {
	return repository.DepartmentModel{
		Name:          sql.NullString{String: inputStruct.DepartmentName},
		Description:   sql.NullString{String: inputStruct.Description},
		CreatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		CreatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input departmentService) checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_name_department") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.DepartmentName)
		}

		if service.CheckDBError(err, "uq_department_department_name") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.DepartmentName)
		}
	}

	return err
}
