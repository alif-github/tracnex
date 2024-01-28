package ClientRegistrationLogService

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"time"
)

func (input clientRegistrationLogService) InsertClientRegistrationLog(inputStruct in.ClientRegisterLogRequest, contextModel *applicationModel.ContextModel) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertClientRegistrationLog"
	var result interface{}

	result, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertClientRegistrationLog, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	id = result.(int64)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientRegistrationLogService) doInsertClientRegistrationLog(tx *sql.Tx, inputStructInterface interface{}, _ *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.ClientRegisterLogRequest)
	var clientRegistrationLogModel repository.ClientRegistrationLogModel
	var id int64

	clientRegistrationLogModel = repository.ClientRegistrationLogModel{
		ClientID:              sql.NullString{String: inputStruct.ClientID},
		ClientTypeID:          sql.NullInt64{Int64: inputStruct.ClientTypeID},
		AttributeRequest:      sql.NullString{String: inputStruct.AttributeRequest},
		SuccessStatusAuth:     sql.NullBool{Bool: inputStruct.SuccessStatusAuth},
		SuccessStatusNexcloud: sql.NullBool{Bool: inputStruct.SuccessStatusNexcloud},
		SuccessStatusNexdrive: sql.NullBool{Bool: inputStruct.SuccessStatusNexdrive},
		Resource:              sql.NullString{String: inputStruct.Resource},
		MessageAuth:           sql.NullString{String: inputStruct.MessageAuth},
		MessageNexcloud:       sql.NullString{String: inputStruct.MessageNexcloud},
		MessageNexdrive:       sql.NullString{String: inputStruct.MessageNexdrive},
		Details:               sql.NullString{String: inputStruct.Details},
		Code:                  sql.NullString{String: inputStruct.Code},
		RequestTimeStamp:      sql.NullTime{Time: timeNow},
		RequestCount: 		   sql.NullInt64{Int64: inputStruct.RequestCount},
		CreatedBy:             sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:         sql.NullString{String: constanta.SystemClient},
		CreatedAt:             sql.NullTime{Time: timeNow},
		UpdatedBy:             sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:         sql.NullString{String: constanta.SystemClient},
		UpdatedAt:             sql.NullTime{Time: timeNow},
	}

	id, err = dao.ClientRegistrationLogDAO.InsertClientRegistrationLog(tx, &clientRegistrationLogModel)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName: 		sql.NullString{String: dao.ClientRegistrationLogDAO.TableName},
		PrimaryKey: 	sql.NullInt64{Int64: id},
	})

	output = id

	err = errorModel.GenerateNonErrorModel()
	return
}
