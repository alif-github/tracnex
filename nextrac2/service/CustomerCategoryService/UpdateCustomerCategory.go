package CustomerCategoryService

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

func (input customerCategoryService) UpdateCustomerCategory(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateCustomerCategory"
		inputStruct in.CustomerCategoryRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdate)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateCustomerCategory, func(interface{}, applicationModel.ContextModel) {
		//additional function
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input customerCategoryService) doUpdateCustomerCategory(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName         = "doUpdateCustomerCategory"
		inputStruct      = inputStructInterface.(in.CustomerCategoryRequest)
		customerCatModel repository.CustomerCategoryModel
		scope            map[string]interface{}
		customerCatOnDB  repository.CustomerCategoryModel
	)

	customerCatModel = input.convertStructToModelUpdate(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	//--- Get scope
	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	//--- Validate ID to DB
	customerCatOnDB, err = dao.CustomerCategoryDAO.GetCustomerCategoryForUpdate(tx, repository.CustomerCategoryModel{
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

	if customerCatOnDB.UpdatedAt.Time.Sub(inputStruct.UpdatedAt) != time.Duration(0) {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.CustomerCategory)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.CustomerCategoryDAO.TableName, customerCatModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	err = dao.CustomerCategoryDAO.UpdateCustomerCategory(tx, customerCatModel)
	return
}

func (input customerCategoryService) convertStructToModelUpdate(inputStruct in.CustomerCategoryRequest, authAccessModel model2.AuthAccessTokenModel, timeNow time.Time) repository.CustomerCategoryModel {
	return repository.CustomerCategoryModel{
		ID:                   sql.NullInt64{Int64: inputStruct.ID},
		CustomerCategoryName: sql.NullString{String: inputStruct.CustomerCategoryName},
		UpdatedBy:            sql.NullInt64{Int64: authAccessModel.ResourceUserID},
		UpdatedAt:            sql.NullTime{Time: timeNow},
		UpdatedClient:        sql.NullString{String: authAccessModel.ClientID},
	}
}

func (input customerCategoryService) validateUpdate(inputStruct *in.CustomerCategoryRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdate()
}
