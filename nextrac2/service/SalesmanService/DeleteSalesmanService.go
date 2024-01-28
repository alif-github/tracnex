package SalesmanService

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
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/DataScopeService"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"time"
)

func (input salesmanService) DeleteSalesman(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteSalesman"
		inputStruct in.SalesmanRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDelete)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteSalesman, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- Additional Function
	})

	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_DELETE_SALESMAN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanService) doDeleteSalesman(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName      = "DeleteSalesmanService.go"
		funcName      = "doDeleteSalesman"
		inputStruct   = inputStructInterface.(in.SalesmanRequest)
		salesmanModel repository.SalesmanModel
		salesmanOnDB  repository.SalesmanModel
		tempDataAudit []repository.AuditSystemModel
		scopeLimit    map[string]interface{}
	)

	scopeLimit, err = input.validateDataScopeSalesman(contextModel)
	if err.Error != nil {
		return
	}

	salesmanModel = repository.SalesmanModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	salesmanModel.CreatedBy.Int64 = 0
	salesmanOnDB, err = dao.SalesmanDAO.GetSalesmanForUpdateDelete(serverconfig.ServerAttribute.DBConnection, salesmanModel, scopeLimit, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if salesmanOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.SalesmanID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, salesmanOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if salesmanOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(fileName, funcName, "Salesman")
		return
	}

	if salesmanOnDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.UpdatedAt)
		return
	}

	// ----------- Update for delete
	encodedStr, errorS := service.RandToken(8)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	salesmanModel.Nik.String = salesmanOnDB.Nik.String + encodedStr
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.SalesmanDAO.TableName, salesmanOnDB.ID.Int64, 0)...)
	err = dao.SalesmanDAO.DeleteSalesman(tx, salesmanModel)
	if err.Error != nil {
		return
	}

	scope := repository.DataScopeModel{Scope: sql.NullString{String: constanta.SalesmanDataScope + ":" + strconv.Itoa(int(salesmanOnDB.ID.Int64))}}
	_, tempDataAudit, err = DataScopeService.DataScopeService.DoDeleteDataScope(tx, scope, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, tempDataAudit...)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanService) validateDelete(inputStruct *in.SalesmanRequest) errorModel.ErrorModel {
	return inputStruct.ValidationDeleteSalesman()
}
