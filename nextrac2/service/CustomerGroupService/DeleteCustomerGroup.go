package CustomerGroupService

import (
	"database/sql"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/DataScopeService"
	"time"
)

func (input customerGroupService) DeleteCustomerGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteCustomerGroup"
		inputStruct in.CustomerGroupRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDeleteCustomerGroup)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteCustomerGroup, func(interface{}, applicationModel.ContextModel) {
		// Additional function
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input customerGroupService) doDeleteCustomerGroup(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName           = "DeleteCustomerGroup.go"
		funcName           = "doUpdateCustomerGroup"
		inputStruct        = inputStructInterface.(in.CustomerGroupRequest)
		customerGroupModel repository.CustomerGroupModel
		customerGroupOnDB  repository.CustomerGroupModel
		scope              map[string]interface{}
		tempDataAudit      []repository.AuditSystemModel
	)

	customerGroupModel = repository.CustomerGroupModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	//--- Get scope
	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	//--- Validate ID to DB
	customerGroupOnDB, err = dao.CustomerGroupDAO.GetCustomerGroupForDelete(tx, repository.CustomerGroupModel{
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

	err = input.CheckUserLimitedByOwnAccess(contextModel, customerGroupOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if customerGroupOnDB.UpdatedAt.Time.Sub(inputStruct.UpdatedAt) != time.Duration(0) {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.CustomerGroup)
		return
	}

	//--- Update for delete
	encodedStr, errorS := service.RandToken(constanta.RandTokenForDeleteLength)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	customerGroupModel.CustomerGroupID.String = customerGroupOnDB.CustomerGroupID.String + encodedStr
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.CustomerGroupDAO.TableName, customerGroupModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	err = dao.CustomerGroupDAO.DeleteCustomerGroup(tx, customerGroupModel)
	if err.Error != nil {
		return
	}

	//--- Delete data scope
	_, tempDataAudit, err = DataScopeService.DataScopeService.DoDeleteDataScope(tx, repository.DataScopeModel{
		Scope: sql.NullString{String: fmt.Sprintf("%s:%d", constanta.CustomerGroupDataScope, customerGroupOnDB.ID.Int64)},
	}, contextModel, timeNow)

	dataAudit = append(dataAudit, tempDataAudit...)
	return

}

func (input customerGroupService) validateDeleteCustomerGroup(inputStruct *in.CustomerGroupRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDeleteCustomerGroup()
}
