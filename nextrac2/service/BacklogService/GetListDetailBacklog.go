package BacklogService

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
)

func (input backlogService) GetListDetailBacklog(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validOrderBy  = []string{
			"id",
			"layer_1",
			"layer_2",
			"layer_3",
			"redmine_number",
			"sprint",
			"pic",
			"status",
			"mandays",
		}
		validSearchBy = []string{
			"pic",
			"redmine_number",
			"sprint",
		}
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, applicationModel.GetListDetailBacklogValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListDetailBacklog(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input backlogService) doGetListDetailBacklog(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		funcName         = "doGetListDetailBacklog"
		dbResult         []interface{}
		isFilterBySprint bool
		scope            map[string]interface{}
	)

	// validate filter by sprint
	for i := 0; i < len(searchByParam); i++ {
		if searchByParam[i].SearchKey == "sprint" {
			isFilterBySprint = true
		}
	}

	if !isFilterBySprint {
		err = errorModel.GenerateSprintFilterError(input.FileName, funcName)
		return
	}

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	dbResult, err = dao.BacklogDAO.GetListDetailBacklog(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseGetListDetail(dbResult)
	return
}

func (input backlogService) convertModelToResponseGetListDetail(dbResult []interface{}) (result []out.DetailBacklogResponse) {
	for _, dbResultItem := range dbResult {
		var (
			item           = dbResultItem.(repository.BacklogModel)
			departmentCode string
		)

		if item.DepartmentId.Int64 == 1 {
			departmentCode = constanta.DepartmentDeveloper
		} else if item.DepartmentId.Int64 == 2 {
			departmentCode = constanta.DepartmentQAQC
		}

		result = append(result, out.DetailBacklogResponse{
			ID:             item.ID.Int64,
			Layer1:         item.Layer1.String,
			Layer2:         item.Layer2.String,
			Layer3:         item.Layer3.String,
			Redmine:        item.RedmineNumber.Int64,
			Sprint:         item.Sprint.String,
			Pic:            item.EmployeeName.String,
			Status:         item.Status.String,
			Mandays:        item.Mandays.Float64,
			DepartmentCode: departmentCode,
			UpdatedAt:      item.UpdatedAt.Time,
		})
	}

	return result
}

func (input backlogService) InitiateGetListDetailParentBacklog(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		countData     interface{}
		validSearchBy = []string{
			"pic",
			"redmine_number",
			"sprint",
		}
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, validSearchBy, applicationModel.GetListDetailBacklogValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateDetailBacklog(searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	if countData == nil {
		countData = 0
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListDetailBacklogValidOperator,
		CountData:     countData.(int),
	}
	return
}

func (input backlogService) doInitiateDetailBacklog(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		funcName         = "doInitiateDetailBacklog"
		isFilterBySprint bool
		scope            map[string]interface{}
	)

	// validate filter by sprint
	for i := 0; i < len(searchByParam); i++ {
		if searchByParam[i].SearchKey == "sprint" {
			isFilterBySprint = true
		}
	}

	if !isFilterBySprint {
		err = errorModel.GenerateSprintFilterError(input.FileName, funcName)
		return
	}

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	output, err = dao.BacklogDAO.GetCountDetailBacklog(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
