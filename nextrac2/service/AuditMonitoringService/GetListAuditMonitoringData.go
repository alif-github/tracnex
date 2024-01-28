package AuditMonitoringService

import (
	"database/sql"
	"net/http"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
)

func (input auditMonitoringService) GetListAuditMonitoring(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListAuditMonitoringValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListAuditMonitoringData(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)

	return
}

func (input auditMonitoringService) doGetListAuditMonitoringData(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output []out.AuditMonitoringResponse, err errorModel.ErrorModel) {
	var dbResult []interface{}
	var listScope map[string]interface{}
	funcName := "doGetListAuditMonitoringData"

	if err = input.validateSearchBy(searchByParam); err.Error != nil {
		return
	}

	// Validate and get table name with menu_code
	if menuCodeIdx := input.getIndexSearchByParam(searchByParam, constanta.MenuCodeDDLName); menuCodeIdx >= 0 {
		if input.getIndexSearchByParam(searchByParam, constanta.TableNameDDLName) >= 0 {
			err = errorModel.GenerateFormatFieldError("GetListAuditMonitoringData.go", funcName, constanta.Filter)
			return
		} else {
			if err = input.getTableNameByMenuCode(searchByParam, menuCodeIdx, contextModel.IsAdmin); err.Error != nil {
				return
			}
		}
	}

	if contextModel.IsAdmin {
		listScope = make(map[string]interface{})
		listScope[constanta.NexsoftDataScope] = []interface{}{"all"}
	} else {
		listScope, err = input.validateDataScope(contextModel)
		if err.Error != nil {
			return
		}
	}

	dbResult, err = dao.AuditSystemDAO.GetListAuditData(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, listScope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertToListDTOOut(dbResult)
	return
}

func (input auditMonitoringService) getTableNameByMenuCode(searchByParam []in.SearchByParam, menuCodeIdx int, isAdmin bool) (err errorModel.ErrorModel) {
	var tableName string
	funcName := "getTableNameByMenuCode"

	tableName, err = dao.MenuDAO.GetTableNameWithMenuCode(serverconfig.ServerAttribute.DBConnection, repository.MenuModel{
		MenuCode: sql.NullString{String: searchByParam[menuCodeIdx].SearchValue},
	}, isAdmin)

	if err.Error != nil {
		return
	}

	if util2.IsStringEmpty(tableName) {
		tableName = getMenuCodeConstanta(searchByParam[menuCodeIdx].SearchValue)
		if util2.IsStringEmpty(tableName) {
			err = errorModel.GenerateUnknownDataError("GetListAuditMonitoringData.go", funcName, constanta.MenuCodeDDLName)
		}
	}

	searchByParam[menuCodeIdx] = in.SearchByParam{
		SearchKey:      constanta.TableNameDDLName,
		DataType:       "char",
		SearchOperator: "eq",
		SearchValue:    tableName,
		SearchType:     searchByParam[menuCodeIdx].SearchType,
	}

	return
}

func (input auditMonitoringService) getIndexSearchByParam(searchByParam []in.SearchByParam, searchBy string) int {
	for i, param := range searchByParam {
		if param.SearchKey == searchBy {
			return i
		}
	}
	return -1
}

func (input auditMonitoringService) validateSearchBy(searchByParam []in.SearchByParam) (err errorModel.ErrorModel) {
	var isPrimaryKey bool
	var isTableName bool
	for i := 0; i < len(searchByParam); i++ {
		if searchByParam[i].SearchKey == "primary_key" {
			isPrimaryKey = true
		}
		if searchByParam[i].SearchKey == "table_name" || searchByParam[i].SearchKey == "menu_code" {
			isTableName = true
		}
	}
	if isPrimaryKey {
		if !isTableName {
			err = errorModel.GenerateEmptyFieldError(input.FileName, "doGetListAuditMonitoringData", constanta.TableName)
			return
		}
	}

	return
}

func (input auditMonitoringService) convertToListDTOOut(dbResult []interface{}) (result []out.AuditMonitoringResponse) {
	for i := 0; i < len(dbResult); i++ {
		repo := dbResult[i].(repository.AuditSystemModel)
		result = append(result, out.AuditMonitoringResponse{
			ID:            repo.ID.Int64,
			TableName:     repo.TableName.String,
			PrimaryKey:    repo.PrimaryKey.Int64,
			Action:        repo.Action.Int32,
			CreatedName:   repo.CreatedName.String,
			CreatedBy:     repo.CreatedBy.Int64,
			CreatedClient: repo.CreatedClient.String,
			CreatedAt:     repo.CreatedAt.Time,
		})
	}
	return result
}

func (input auditMonitoringService) InitiateGetListAuditMonitoringData(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData int

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListAuditMonitoringValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateListMonitoringData(searchByParam, *contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_LIST_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListAuditMonitoringValidOperator,
		CountData:     countData,
	}

	return
}

func (input auditMonitoringService) doInitiateListMonitoringData(searchByParam []in.SearchByParam, contextModel applicationModel.ContextModel) (output int, err errorModel.ErrorModel) {
	output, err = dao.AuditSystemDAO.CountAuditData(serverconfig.ServerAttribute.DBConnection, searchByParam, false, contextModel.LimitedByCreatedBy)
	return
}
