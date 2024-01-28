package StandarManhourService

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
	"nexsoft.co.id/nextrac2/serverconfig"
	"time"
)

func (input standarManhourService) InsertStandarManhour(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertStandarManhour"
		inputStruct in.StandarManhourRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateInsertStandarManhour)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertStandarManhour, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input standarManhourService) doInsertStandarManhour(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName    = "InsertStandarManhourService.go"
		funcName    = "doInsertStandarManhour"
		inputStruct = inputStructInterface.(in.StandarManhourRequest)
		inputModel  = input.convertDTOToModel(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
		db          = serverconfig.ServerAttribute.DBConnection
		idInserted  int64
		isExist     bool
	)

	//-- Check ID Department
	isExist, err = input.DepartmentDAO.CheckIDDepartment(db, repository.DepartmentModel{ID: sql.NullInt64{Int64: inputModel.DepartmentID.Int64}})
	if err.Error != nil {
		return
	}

	if !isExist {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.DepartmentId)
		return
	}

	//-- Insert To DB
	idInserted, err = input.StandarManhourDAO.InsertStandarManhour(tx, inputModel)
	if err.Error != nil {
		return
	}

	//-- Insert To Data Audit
	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.StandarManhourDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idInserted},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input standarManhourService) convertDTOToModel(inputStruct in.StandarManhourRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.StandarManhourModel {
	return repository.StandarManhourModel{
		Case:          sql.NullString{String: inputStruct.Case},
		DepartmentID:  sql.NullInt64{Int64: inputStruct.DepartmentID},
		Manhour:       sql.NullFloat64{Float64: inputStruct.Manhour},
		CreatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		CreatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input standarManhourService) ValidateInsertStandarManhour(inputStruct *in.StandarManhourRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsert()
}
