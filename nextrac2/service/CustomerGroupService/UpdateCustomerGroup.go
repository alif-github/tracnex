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
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input customerGroupService) UpdateCustomerGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateCustomerGroup"
		inputStruct in.CustomerGroupRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateCustomerGroup)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateCustomerGroup, func(interface{}, applicationModel.ContextModel) {
		// Additional function
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input customerGroupService) doUpdateCustomerGroup(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName           = "doUpdateCustomerGroup"
		inputStruct        = inputStructInterface.(in.CustomerGroupRequest)
		customerGroupModel repository.CustomerGroupModel
		customerGroupOnDB  repository.CustomerGroupModel
		scope              map[string]interface{}
	)

	customerGroupModel = input.convertStructToModelUpdate(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	//--- Get scope
	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	//--- Validate ID to DB
	customerGroupOnDB, err = dao.CustomerGroupDAO.GetCustomerGroupForUpdate(tx, repository.CustomerGroupModel{
		ID: customerGroupModel.ID,
	}, scope, input.MappingScopeDB)

	if err.Error != nil {
		return
	}

	if customerGroupOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, customerGroupOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if customerGroupOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.CustomerGroup)
		return
	}

	if customerGroupOnDB.UpdatedAt.Time.Sub(inputStruct.UpdatedAt) != time.Duration(0) {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.CustomerGroup)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.CustomerGroupDAO.TableName, customerGroupModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.CustomerGroupDAO.UpdateCustomerGroup(tx, customerGroupModel)
	return
}

func (input customerGroupService) convertStructToModelUpdate(inputStruct in.CustomerGroupRequest, authAccessModel model2.AuthAccessTokenModel, timeNow time.Time) repository.CustomerGroupModel {
	return repository.CustomerGroupModel{
		ID:                sql.NullInt64{Int64: inputStruct.ID},
		CustomerGroupName: sql.NullString{String: inputStruct.CustomerGroupName},
		UpdatedBy:         sql.NullInt64{Int64: authAccessModel.ResourceUserID},
		UpdatedAt:         sql.NullTime{Time: timeNow},
		UpdatedClient:     sql.NullString{String: authAccessModel.ClientID},
	}
}

func (input customerGroupService) validateUpdateCustomerGroup(inputStruct *in.CustomerGroupRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateCustomerGroup()
}
