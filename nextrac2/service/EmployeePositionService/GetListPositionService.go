package EmployeePositionService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input employeePositionService) InitiatePosition(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		countData     interface{}
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListEmployeePositionValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiatePosition(searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListEmployeePositionValidOperator,
		CountData:     countData.(int),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeePositionService) doInitiatePosition(searchByParam []in.SearchByParam, _ *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		createdBy int64
		db        = serverconfig.ServerAttribute.DBConnection
	)

	output = 0
	output, err = dao.EmployeePositionDAO.GetCountPosition(db, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeePositionService) GetListPosition(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListEmployeePositionValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListPosition(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeePositionService) doGetListPosition(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
		db       = serverconfig.ServerAttribute.DBConnection
	)

	dbResult, err = dao.EmployeePositionDAO.GetListPosition(db, inputStruct, searchByParam)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseGetList(dbResult)
	return
}

func (input employeePositionService) convertModelToResponseGetList(dbResult []interface{}) (result []out.ListEmployeePosition) {
	for _, dbResultItem := range dbResult {
		item := dbResultItem.(repository.EmployeePositionModel)
		result = append(result, out.ListEmployeePosition{
			ID:          item.ID.Int64,
			Name:        item.Name.String,
			Description: item.Description.String,
			UpdatedAt:   item.UpdatedAt.Time,
		})
	}

	return result
}
