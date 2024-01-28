package CustomerCategoryService

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

func (input customerCategoryService) DeleteCustomerCategory(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteCustomerCategory"
		inputStruct in.CustomerCategoryRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDelete)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteCustomerCategory, func(interface{}, applicationModel.ContextModel) {
		//additional function
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input customerCategoryService) doDeleteCustomerCategory(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName         = "DeleteCustomerCategory.go"
		funcName         = "doDeleteCustomerCategory"
		inputStruct      = inputStructInterface.(in.CustomerCategoryRequest)
		customerCatModel repository.CustomerCategoryModel
		customerCatOnDB  repository.CustomerCategoryModel
		tempDataAudit    []repository.AuditSystemModel
		scope            map[string]interface{}
	)

	customerCatModel = repository.CustomerCategoryModel{
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

	// Validate ID to DB
	customerCatOnDB, err = dao.CustomerCategoryDAO.GetCustomerCategoryForDelete(tx, repository.CustomerCategoryModel{
		ID: customerCatModel.ID,
	}, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if customerCatOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, customerCatOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if customerCatOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.CustomerCategory)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, customerCatOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if customerCatOnDB.UpdatedAt.Time.Sub(inputStruct.UpdatedAt) != time.Duration(0) {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.CustomerCategory)
		return
	}

	//--- Update for delete
	encodedStr, errorS := service.RandToken(constanta.RandTokenForDeleteLength)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	customerCatModel.CustomerCategoryID.String = customerCatOnDB.CustomerCategoryID.String + encodedStr

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.CustomerCategoryDAO.TableName, customerCatModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	err = dao.CustomerCategoryDAO.DeleteCustomerCategory(tx, customerCatModel)
	if err.Error != nil {
		return
	}

	//--- delete data scope
	_, tempDataAudit, err = DataScopeService.DataScopeService.DoDeleteDataScope(tx, repository.DataScopeModel{
		Scope: sql.NullString{String: fmt.Sprintf("%s:%d", constanta.CustomerCategoryDataScope, customerCatOnDB.ID.Int64)},
	}, contextModel, timeNow)

	dataAudit = append(dataAudit, tempDataAudit...)
	return
}

func (input customerCategoryService) validateDelete(inputStruct *in.CustomerCategoryRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}
