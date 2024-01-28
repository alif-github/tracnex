package ProductGroupService

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
	"time"
)

func (input productGroupService) InsertProductGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "InsertProductGroup"
	var inputStruct in.ProductGroupRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertProductGroup, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input productGroupService) doInsertProductGroup(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.ProductGroupRequest)
	productGroupModel := input.convertStructToModel(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	insertedID, err := dao.ProductGroupDAO.InsertProductGroup(tx, productGroupModel)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.ProductGroupDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: insertedID},
	})

	var dataAuditTemp repository.AuditSystemModel
	dataAuditTemp, err =input.GenerateDataScope(tx, insertedID, dao.ProductGroupDAO.TableName, constanta.ProductGroupDataScope, contextModel.AuthAccessTokenModel.ResourceUserID, contextModel.AuthAccessTokenModel.ClientID, timeNow)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, dataAuditTemp)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productGroupService) convertStructToModel(inputStruct in.ProductGroupRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.ProductGroupModel {
	return repository.ProductGroupModel{
		ProductGroupName: sql.NullString{String: inputStruct.ProductGroupName},
		CreatedBy:        sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedAt:        sql.NullTime{Time: timeNow},
		CreatedClient:    sql.NullString{String: authAccessToken.ClientID},
		UpdatedBy:        sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		UpdatedClient:    sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input productGroupService) validateInsert(inputStruct *in.ProductGroupRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsert()
}
