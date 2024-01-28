package EmployeeLevelService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
)

func (input employeeLevelGetListService) GetEmployeeLevelMatrix(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct         in.GetListDataDTO
		searchByParam       []in.SearchByParam
		isOnlyHaveOwnAccess bool
		createdBy           int64
		levels              []interface{}
		db                  = serverconfig.ServerAttribute.DBConnection
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListEmployeeLevelValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	createdBy, isOnlyHaveOwnAccess = service.CheckIsOnlyHaveOwnPermission(*context)
	if !isOnlyHaveOwnAccess {
		createdBy = 0
	}

	levels, err = dao.EmployeeLevelDAO.GetEmployeeLevelMatrix(db, inputStruct, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.convertRepoToDTO(levels)
	output.Status = service.GetResponseMessages("SUCCESS_GET_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeLevelGetListService) InitiateEmployeeLevelMatrix(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam       []in.SearchByParam
		db                  = serverconfig.ServerAttribute.DBConnection
		isOnlyHaveOwnAccess bool
		createdBy           int64
		countData           int
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListEmployeeLevelValidOperator)
	if err.Error != nil {
		return
	}

	createdBy, isOnlyHaveOwnAccess = service.CheckIsOnlyHaveOwnPermission(*context)
	if !isOnlyHaveOwnAccess {
		createdBy = 0
	}

	countData, err = dao.EmployeeLevelDAO.GetCountLevelMatrix(db, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		CountData:     countData,
		ValidOperator: applicationModel.GetListEmployeeLevelValidOperator,
	}

	output.Status = service.GetResponseMessages("SUCCESS_INITIATE_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}
