package CustomerGroupService

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

func (input customerGroupService) InsertCustomerGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "InsertCustomerGroup"
	var inputStruct in.CustomerGroupRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.insertCustomerGroup, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)

	return
}

func (input customerGroupService) insertCustomerGroup(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	//funcName := "insertCustomerGroup"
	inputStruct := inputStructInterface.(in.CustomerGroupRequest)

	customerGroupModel := input.convertSturctToModelInsert(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	idCustomerGroup, err := dao.CustomerGroupDAO.InsertCustomerGroup(tx, customerGroupModel)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.CustomerGroupDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idCustomerGroup},
	})

	var dataAuditTemp repository.AuditSystemModel
	dataAuditTemp, err =input.GenerateDataScope(tx, idCustomerGroup, dao.CustomerGroupDAO.TableName, constanta.CustomerGroupDataScope, contextModel.AuthAccessTokenModel.ResourceUserID, contextModel.AuthAccessTokenModel.ClientID, timeNow)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, dataAuditTemp)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerGroupService) convertSturctToModelInsert(inputStruct in.CustomerGroupRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.CustomerGroupModel {
	return repository.CustomerGroupModel{
		CustomerGroupID:   sql.NullString{String: inputStruct.CustomerGroupID},
		CustomerGroupName: sql.NullString{String: inputStruct.CustomerGroupName},
		CreatedBy:         sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedAt:         sql.NullTime{Time: timeNow},
		CreatedClient:     sql.NullString{String: authAccessToken.ClientID},
		UpdatedBy:         sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:         sql.NullTime{Time: timeNow},
		UpdatedClient:     sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input customerGroupService) validateInsert(inputStruct *in.CustomerGroupRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertCustomerGroup()
}
