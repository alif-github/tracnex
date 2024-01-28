package ComponentService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"time"
)

func (input componentService) InsertComponent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertComponent"
		inputStruct in.ComponentRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsertComponent)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertComponent, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input componentService) doInsertComponent(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct = inputStructInterface.(in.ComponentRequest)
		inputModel  = input.convertDTOToModelInsert(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
		insertedID  int64
	)

	insertedID, err = dao.ComponentDAO.InsertComponent(tx, inputModel)
	if err.Error != nil {
		err = input.checkDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.ComponentDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: insertedID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input componentService) convertDTOToModelInsert(inputStruct in.ComponentRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.ComponentModel {
	return repository.ComponentModel{
		ComponentName: sql.NullString{String: inputStruct.ComponentName},
		CreatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		CreatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input componentService) validateInsertComponent(inputStruct *in.ComponentRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsert()
}
