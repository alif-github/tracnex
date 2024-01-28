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
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input productGroupService) UpdateProductGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "UpdateProductGroup"
	var inputStruct in.ProductGroupRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdate)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateProductGroup, func(interface{}, applicationModel.ContextModel) {
		// additional Function
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input productGroupService) doUpdateProductGroup(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doUpdateProductGroup"
	inputStruct := inputStructInterface.(in.ProductGroupRequest)

	productGroupModel := input.convertStructToModelUpdate(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	// Get scope
	scope, err := input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	// Validate ID to DB
	productGroupOnDB, err := dao.ProductGroupDAO.GetProductGroupForDelete(tx, repository.ProductGroupModel{
		ID: productGroupModel.ID,
	}, scope, input.MappingScopeDB)

	if err.Error != nil {
		return
	}

	if productGroupOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	if productGroupOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.ProductGroup)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, productGroupOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if productGroupOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.ProductGroup)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ProductGroupDAO.TableName, productGroupModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.ProductGroupDAO.UpdateProductGroup(tx, productGroupModel)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}
	return
}

func (input productGroupService) convertStructToModelUpdate(inputStruct in.ProductGroupRequest, authAccessModel model2.AuthAccessTokenModel, timeNow time.Time) repository.ProductGroupModel {
	return repository.ProductGroupModel{
		ID:               sql.NullInt64{Int64: inputStruct.ID},
		ProductGroupName: sql.NullString{String: inputStruct.ProductGroupName},
		UpdatedBy:        sql.NullInt64{Int64: authAccessModel.ResourceUserID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		UpdatedClient:    sql.NullString{String: authAccessModel.ClientID},
	}
}

func (input productGroupService) validateUpdate(inputStruct *in.ProductGroupRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdate()
}
