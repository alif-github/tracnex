package EmployeeService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
)

func (input employeeService) InitiateGetEmployeeLeaveTypes(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		db = serverconfig.ServerAttribute.DBConnection
		searchByParam []in.SearchByParam
		validOrderBy = []string{"id"}
	)

	_, searchByParam, errModel = input.ReadAndValidateGetCountData(request, []string{}, applicationModel.GetListEmployeeLeaveTypesValidOperator)
	if errModel.Error != nil {
		return
	}

	employee, errModel := dao.EmployeeDAO.GetByUserId(db, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	result, errModel := dao.AllowanceDAO.InitiateListByEmployeeLevelIdAndEmployeeGradeId(db, searchByParam, repository.Allowance{
		EmployeeGradeId: sql.NullInt64{Int64: employee.GradeID.Int64},
		EmployeeLevelId: sql.NullInt64{Int64: employee.LevelID.Int64},
	})
	if errModel.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  validOrderBy,
		ValidSearchBy: nil,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListEmployeeLeaveTypesValidOperator,
		CountData:     result,
	}
	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	return
}

func (input employeeService) GetEmployeeLeaveTypes(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		db = serverconfig.ServerAttribute.DBConnection
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validOrderBy = []string{"id"}
	)

	inputStruct, searchByParam, errModel = input.ReadAndValidateGetListData(request, []string{}, validOrderBy, applicationModel.GetListEmployeeLeaveTypesValidOperator, service.DefaultLimit)
	if errModel.Error != nil {
		return
	}

	employee, errModel := dao.EmployeeDAO.GetByUserId(db, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	results, errModel := dao.AllowanceDAO.GetListByEmployeeLevelIdAndEmployeeGradeId(db, inputStruct, searchByParam, repository.Allowance{
		EmployeeGradeId: sql.NullInt64{Int64: employee.GradeID.Int64},
		EmployeeLevelId: sql.NullInt64{Int64: employee.LevelID.Int64},
	})
	if errModel.Error != nil {
		return
	}

	output.Data.Content = input.toEmployeeLeaveAllowanceResponse(results)
	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) toEmployeeLeaveAllowanceResponse(data []interface{}) (result []out.EmployeeLeaveAllowance) {
	for _, item := range data {
		allowance, _ := item.(repository.Allowance)

		result = append(result, out.EmployeeLeaveAllowance{
			Id:            allowance.ID.Int64,
			AllowanceName: allowance.AllowanceName.String,
			AllowanceType: allowance.AllowanceType.String,
			Value:         allowance.Value.String,
			CreatedAt:     allowance.CreatedAt.Time,
			UpdatedAt: 	   allowance.UpdatedAt.Time,
		})
	}

	return
}
