package ProductGroupService

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

func (input productGroupService) DeleteProductGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "DeleteProductGroup"
	var inputStruct in.ProductGroupRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDelete)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteProductGroup, func(interface{}, applicationModel.ContextModel) {
		// additional Function
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input productGroupService) doDeleteProductGroup(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	fileName := "DeleteProductGroup.go"
	funcName := "doDeleteProductGroup"
	inputStruct := inputStructInterface.(in.ProductGroupRequest)

	productGroupModel := repository.ProductGroupModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

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

	err = input.CheckUserLimitedByOwnAccess(contextModel, productGroupOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if productGroupOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.ProductGroup)
		return
	}

	if productGroupOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.ProductGroup)
		return
	}

	// ----------- Update for delete
	encodedStr, errorS := service.RandToken(10)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	productGroupModel.ProductGroupName.String = productGroupOnDB.ProductGroupName.String + encodedStr

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.ProductGroupDAO.TableName, productGroupModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.ProductGroupDAO.DeleteProductGroup(tx, productGroupModel)
	if err.Error != nil {
		return
	}

	// delete data scope
	_, tempDataAudit, err := DataScopeService.DataScopeService.DoDeleteDataScope(tx, repository.DataScopeModel{
		Scope: sql.NullString{String: fmt.Sprintf("%s:%d", constanta.ProductGroupDataScope, productGroupOnDB.ID.Int64)},
	}, contextModel, timeNow)

	dataAudit = append(dataAudit, tempDataAudit...)
	return

}

func (input productGroupService) validateDelete(inputStruct *in.ProductGroupRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}
