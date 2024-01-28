package EmployeeGradeService

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

func (input employeeGradeGetListService) GetEmployeeGradeMatrix(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct         in.GetListDataDTO
		searchByParam       []in.SearchByParam
		createdBy           int64
		isOnlyHaveOwnAccess bool
		grades              []interface{}
		db                  = serverconfig.ServerAttribute.DBConnection
		validSearchBy       = []string{"id", "level_id"}
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, validSearchBy, input.ValidOrderBy, applicationModel.GetListEmployeeGradeValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	createdBy, isOnlyHaveOwnAccess = service.CheckIsOnlyHaveOwnPermission(*context)
	if !isOnlyHaveOwnAccess {
		createdBy = 0
	}

	grades, err = dao.EmployeeGradeDAO.GetEmployeeGradeMatrix(db, inputStruct, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.convertRepoToDTO(grades)
	output.Status = service.GetResponseMessages("SUCCESS_GET_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeGradeGetListService) InitiateEmployeeGradeMatrix(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam       []in.SearchByParam
		db                  = serverconfig.ServerAttribute.DBConnection
		isOnlyHaveOwnAccess bool
		createdBy           int64
		countData           int
		validSearchBy       = []string{"id", "level_id"}
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, validSearchBy, applicationModel.GetListEmployeeGradeValidOperator)
	if err.Error != nil {
		return
	}

	createdBy, isOnlyHaveOwnAccess = service.CheckIsOnlyHaveOwnPermission(*context)
	if !isOnlyHaveOwnAccess {
		createdBy = 0
	}

	countData, err = dao.EmployeeGradeDAO.GetCountGradeMatrix(db, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: validSearchBy,
		ValidLimit:    input.ValidLimit,
		CountData:     countData,
		ValidOperator: applicationModel.GetListEmployeeGradeValidOperator,
	}

	output.Status = service.GetResponseMessages("SUCCESS_INITIATE_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}
