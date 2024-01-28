package EmployeeContractService

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
	"time"
)

func (input employeeContractService) InsertEmployeeContract(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.EmployeeContractRequest
		funcName    = "InsertEmployeeContract"
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateInsertEmployeeContract)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertEmployeeContract, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input employeeContractService) doInsertEmployeeContract(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName        = "InsertEmployeeContractService.go"
		funcName        = "doInsertEmployeeContract"
		inputStruct     = inputStructInterface.(in.EmployeeContractRequest)
		db              = serverconfig.ServerAttribute.DBConnection
		repoModel       repository.EmployeeContractModel
		isExistEmployee bool
		insertedID      int64
	)

	repoModel = input.convertDTOToModel(inputStruct, contextModel, timeNow, false)
	isExistEmployee, _, err = dao.EmployeeDAO.CheckEmployeeIDByID(db, repoModel.EmployeeID.Int64)
	if err.Error != nil {
		return
	}

	if !isExistEmployee {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.EmployeeID)
		return
	}

	//--- Insert Employee Contract
	insertedID, err = dao.EmployeeContractDAO.InsertEmployeeContract(tx, repoModel)
	if err.Error != nil {
		err = CheckDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.EmployeeContractDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: insertedID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeContractService) convertDTOToModel(inputStruct in.EmployeeContractRequest, contextModel *applicationModel.ContextModel, timeNow time.Time, isUpdate bool) (repo repository.EmployeeContractModel) {
	repo = repository.EmployeeContractModel{
		ContractNo:    sql.NullString{String: inputStruct.Contract},
		Information:   sql.NullString{String: inputStruct.Information},
		EmployeeID:    sql.NullInt64{Int64: inputStruct.EmployeeID},
		FromDate:      sql.NullTime{Time: inputStruct.FromDate},
		ThruDate:      sql.NullTime{Time: inputStruct.ThruDate},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	if isUpdate {
		repo.ID = sql.NullInt64{Int64: inputStruct.ID}
		repo.CreatedBy = sql.NullInt64{}
		repo.CreatedClient = sql.NullString{}
		repo.CreatedAt = sql.NullTime{}
	}

	return
}

func (input employeeContractService) ValidateInsertEmployeeContract(inputStruct *in.EmployeeContractRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsert()
}
