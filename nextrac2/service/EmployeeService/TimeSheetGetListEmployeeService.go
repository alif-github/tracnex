package EmployeeService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input employeeService) InitiateEmployeeTimeSheet(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		countData     interface{}
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListEmployeeValidOperator)
	if err.Error != nil {
		return
	}

	searchByTs := input.ValidSearchBy
	for i := 0; i < len(searchByTs); i++ {
		if searchByTs[i] == "is_timesheet" {
			searchByTs = append(searchByTs[:i], searchByTs[i+1:]...)
			i = -1
		} else if searchByTs[i] == "is_redmine_check" {
			searchByTs = append(searchByTs[:i], searchByTs[i+1:]...)
			i = -1
		}
	}

	countData, err = input.doInitiateEmployeeTimeSheet(searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: searchByTs,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListEmployeeValidOperator,
		CountData:     countData.(int),
	}

	return
}

func (input employeeService) doInitiateEmployeeTimeSheet(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		createdBy = contextModel.LimitedByCreatedBy
		db        = serverconfig.ServerAttribute.DBConnection
		scope     map[string]interface{}
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	output, err = dao.EmployeeDAO.GetCountEmployeeTimeSheet(db, searchByParam, createdBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) GetListEmployeeTimeSheet(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListEmployeeValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListEmployeeTimeSheet(inputStruct, searchByParam, contextModel, false)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) doGetListEmployeeTimeSheet(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel, isAdmin bool) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult  []interface{}
		scope     map[string]interface{}
		createdBy int64
		db        = serverconfig.ServerAttribute.DBConnection
	)

	if isAdmin {
		scope = make(map[string]interface{})
		scope[constanta.EmployeeDataScope] = []interface{}{"all"}
	} else {
		//--- Add userID when have own permission
		createdBy = contextModel.LimitedByCreatedBy

		//--- Get scope
		scope, err = input.validateDataScope(contextModel)
		if err.Error != nil {
			return
		}
	}

	dbResult, err = dao.EmployeeDAO.GetListEmployeeTimeSheet(db, inputStruct, searchByParam, createdBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseGetListTimeSheet(dbResult)
	return
}

func (input employeeService) convertModelToResponseGetListTimeSheet(dbResult []interface{}) (result interface{}) {
	var resultGetList []out.GetListEmployeeTimeSheetResponse
	for _, dbResultItem := range dbResult {
		var (
			item       = dbResultItem.(repository.EmployeeModel)
			trackerDev in.TrackerDeveloper
			trackerQA  in.TrackerQA
		)

		if item.DepartmentId.Int64 == constanta.QAQCDepartmentID {
			_ = json.Unmarshal([]byte(item.MandaysRate.String), &trackerQA)
		} else {
			_ = json.Unmarshal([]byte(item.MandaysRate.String), &trackerDev)
		}

		//--- Get List
		resultGetList = append(resultGetList, out.GetListEmployeeTimeSheetResponse{
			ID:                    item.ID.Int64,
			Name:                  item.Name.String,
			Nik:                   item.NIK.Int64,
			RedmineId:             item.RedmineId.Int64,
			DepartmentName:        item.DepartmentName.String,
			MandaysRate:           trackerDev.Task,
			MandaysRateAutomation: trackerQA.Automation,
			MandaysRateManual:     trackerQA.Manual,
			CreatedAt:             item.CreatedAt.Time,
			UpdatedAt:             item.UpdatedAt.Time,
			UpdatedBy:             item.UpdatedBy.Int64,
			UpdatedName:           item.UpdatedName.String,
		})
	}

	//--- Return DDL
	return resultGetList
}
