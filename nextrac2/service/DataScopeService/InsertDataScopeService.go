package DataScopeService

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
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input dataScopeService) InsertDataScope(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertDataScope"
		inputStruct in.ScopeRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateInsertScope)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertDataScope, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input dataScopeService) doInsertDataScope(tx *sql.Tx, inputStructInterface interface{}, _ *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName    = "InsertDataScopeService.go"
		funcName    = "doInsertDataScope"
		inputStruct = inputStructInterface.(in.ScopeRequest)
		dataScopeID int64
	)

	if inputStruct.ScopeType != constanta.CustomerGroupDataScope &&
		inputStruct.ScopeType != constanta.CustomerCategoryDataScope &&
		inputStruct.ScopeType != constanta.ProvinceDataScope &&
		inputStruct.ScopeType != constanta.DistrictDataScope &&
		inputStruct.ScopeType != constanta.ProductGroupDataScope &&
		inputStruct.ScopeType != constanta.ClientTypeDataScope &&
		inputStruct.ScopeType != constanta.SalesmanDataScope &&
		inputStruct.ScopeType != constanta.EmployeeDataScope {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.DataScope)
		return
	}

	if inputStruct.ScopeID > 0 {
		var dataAuditTemp repository.AuditSystemModel
		dataAuditTemp, err = input.GenerateDataScope(tx, inputStruct.ScopeID, "-", inputStruct.ScopeType, constanta.SystemID, constanta.SystemClient, timeNow)
		if err.Error != nil {
			return
		}

		dataAudit = append(dataAudit, dataAuditTemp)
		err = errorModel.GenerateNonErrorModel()
		return
	}

	dataScope := repository.DataScopeModel{
		Scope:         sql.NullString{String: inputStruct.ScopeType + ":all"},
		Description:   sql.NullString{String: "Data Scope For Nextrac2 With ID for table : - "},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	dataScopeID, err = dao.DataScopeDAO.InsertDataScope(tx, dataScope)
	if err.Error != nil {
		if err.CausedBy != nil {
			if service.CheckDBError(err, "uq_scope") {
				err = errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.Scope)
				return
			} else if service.CheckDBError(err, "uq_datascope_scope") {
				err = errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.Scope)
				return
			}
		}
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.DataScopeDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: dataScopeID},
	})

	return
}

func CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_product_product_name") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ProductName)
		} else if service.CheckDBError(err, "uq_product_productid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ProductID)
		}
	}
	return err
}

func (input dataScopeService) ValidateInsertScope(inputStruct *in.ScopeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsert()
}
