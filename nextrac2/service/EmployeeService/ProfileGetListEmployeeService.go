package EmployeeService

import (
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
	"strconv"
)

func (input employeeService) InitiateEmployee(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		countData     interface{}
		validSearchBy []string
		validOrderBy  []string
		validOperator map[string]applicationModel.DefaultOperator
	)

	validSearchBy, validOrderBy, validOperator = input.getValidEmployee()
	_, searchByParam, err = input.ReadAndValidateGetCountData(request, validSearchBy, validOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateEmployee(searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  validOrderBy,
		ValidSearchBy: validSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: validOperator,
		CountData:     countData.(int),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) doInitiateEmployee(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		createdBy = contextModel.LimitedByCreatedBy
		db        = serverconfig.ServerAttribute.DBConnection
		scope     map[string]interface{}
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	output, err = dao.EmployeeDAO.GetCountEmployee(db, searchByParam, createdBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) GetListEmployee(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validSearchBy []string
		validOrderBy  []string
		validOperator map[string]applicationModel.DefaultOperator
	)

	validSearchBy, validOrderBy, validOperator = input.getValidEmployee()
	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, validOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListEmployee(inputStruct, searchByParam, contextModel, false, false)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) GetListEmployeeDDL(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		operator      = make(map[string]applicationModel.DefaultOperator)
		validLimit    = service.DefaultLimit
		validSearchBy = []string{
			"nik",
			"redmine_id",
			"name",
			"department",
			"department_id",
			"is_timesheet",
			"is_redmine_check",
		}
		validOrderBy = []string{
			"nik",
			"redmine_id",
			"name",
			"department",
			"updated_at",
		}
	)

	operator["nik"] = applicationModel.DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	operator["redmine_id"] = applicationModel.DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	operator["name"] = applicationModel.DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	operator["department"] = applicationModel.DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	operator["department_id"] = applicationModel.DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	operator["id"] = applicationModel.DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	operator["id_card"] = applicationModel.DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	operator["first_name"] = applicationModel.DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	operator["is_timesheet"] = applicationModel.DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	operator["is_redmine_check"] = applicationModel.DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	contextModel.LimitedByCreatedBy = 0
	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, operator, validLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListEmployee(inputStruct, searchByParam, contextModel, false, true)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) GetListEmployeeByAdmin(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListEmployeeValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListEmployee(inputStruct, searchByParam, contextModel, true, true)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) doGetListEmployee(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel, isAdmin, isDDL bool) (output interface{}, err errorModel.ErrorModel) {
	var (
		fileName                      = "ProfileGetListEmployeeService.go"
		funcName                      = "doGetListEmployee"
		dbResult                      []interface{}
		scope                         map[string]interface{}
		createdBy                     int64
		isTimeSheet, isRedmineCheck   bool
		idxTimeSheet, idxRedmineCheck int
		db                            = serverconfig.ServerAttribute.DBConnection
	)

	if isAdmin {
		scope = make(map[string]interface{})
		scope[constanta.EmployeeDataScope] = []interface{}{"all"}
	} else {
		createdBy = contextModel.LimitedByCreatedBy        //--- Add userID when have own permission
		scope, err = input.validateDataScope(contextModel) //--- Get scope
		if err.Error != nil {
			return
		}
	}

	for idx, itemSearchByParam := range searchByParam {
		if itemSearchByParam.SearchKey == "is_timesheet" {
			isTimeSheetTemp, errS := strconv.ParseBool(itemSearchByParam.SearchValue)
			if errS != nil {
				err = errorModel.GenerateUnknownError(fileName, funcName, errS)
				return
			}

			if isTimeSheetTemp {
				isTimeSheet = true
				idxTimeSheet = idx
			}
		} else if itemSearchByParam.SearchKey == "is_redmine_check" {
			isRedmineCheckTemp, errS := strconv.ParseBool(itemSearchByParam.SearchValue)
			if errS != nil {
				err = errorModel.GenerateUnknownError(fileName, funcName, errS)
				return
			}

			if isRedmineCheckTemp {
				isRedmineCheck = true
				idxRedmineCheck = idx
			}
		}
	}

	if len(searchByParam) > 0 && idxTimeSheet > 0 {
		searchByParam = append(searchByParam[:idxTimeSheet], searchByParam[idxTimeSheet+1:]...)
	}

	if len(searchByParam) > 0 && idxRedmineCheck > 0 {
		searchByParam = append(searchByParam[:idxRedmineCheck], searchByParam[idxRedmineCheck+1:]...)
	}

	if isDDL && isTimeSheet {
		dbResult, err = dao.EmployeeDAO.GetListEmployee(db, inputStruct, searchByParam, createdBy, scope, input.MappingScopeDB, true)
		if err.Error != nil {
			return
		}
	} else if isDDL && isRedmineCheck {
		dbResult, err = dao.EmployeeDAO.GetListEmployeeReport(db, inputStruct, searchByParam, createdBy, scope, input.MappingScopeDB)
		if err.Error != nil {
			return
		}
	} else {
		dbResult, err = dao.EmployeeDAO.GetListEmployee(db, inputStruct, searchByParam, createdBy, scope, input.MappingScopeDB, false)
		if err.Error != nil {
			return
		}
	}

	output = input.convertModelToResponseGetList(dbResult, isDDL)
	return
}

func (input employeeService) convertModelToResponseGetList(dbResult []interface{}, isDDL bool) (result interface{}) {
	var (
		resultGetList []out.GetListEmployeeResponse
		resultDDL     []out.GetListEmployeeForDDLResponse
	)

	for _, dbResultItem := range dbResult {
		var item = dbResultItem.(repository.EmployeeModel)

		//--- Get List
		resultGetList = append(resultGetList, out.GetListEmployeeResponse{
			ID:           item.ID.Int64,
			Name:         item.Name.String,
			NIK:          item.IDCard.String,
			DepartmentID: item.DepartmentId.Int64,
			Department:   item.DepartmentName.String,
			CreatedAt:    item.CreatedAt.Time,
			UpdatedAt:    item.UpdatedAt.Time,
			UpdatedBy:    item.UpdatedBy.Int64,
			UpdatedName:  item.UpdatedName.String,
		})

		//--- Get DDL
		resultDDL = append(resultDDL, out.GetListEmployeeForDDLResponse{
			ID:             item.ID.Int64,
			Name:           item.Name.String,
			Nik:            item.IDCard.String,
			RedmineId:      item.RedmineId.Int64,
			DepartmentName: item.DepartmentName.String,
			CreatedAt:      item.CreatedAt.Time,
			UpdatedAt:      item.UpdatedAt.Time,
			UpdatedBy:      item.UpdatedBy.Int64,
			UpdatedName:    item.UpdatedName.String,
		})
	}

	//--- Is DDL
	if isDDL {
		result = resultDDL
		return
	}

	//--- Return DDL
	return resultGetList
}

func (input employeeService) getValidEmployee() (searchBy, orderBy []string, operator map[string]applicationModel.DefaultOperator) {
	searchBy = []string{
		"nik",
		"name",
		"department",
	}

	orderBy = []string{
		"nik",
		"name",
		"department",
	}

	operator = make(map[string]applicationModel.DefaultOperator)
	operator["nik"] = applicationModel.DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	operator["name"] = applicationModel.DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	operator["department"] = applicationModel.DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	return
}
