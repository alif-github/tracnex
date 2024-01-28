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

func (input employeeService) InitiateGetEmployeeReimbursementTypes(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		db = serverconfig.ServerAttribute.DBConnection
		searchByParam []in.SearchByParam
		validOrderBy = []string{"id"}
	)

	_, searchByParam, errModel = input.ReadAndValidateGetCountData(request, []string{}, applicationModel.GetListEmployeeReimbursementTypesValidOperator)
	if errModel.Error != nil {
		return
	}

	employee, errModel := dao.EmployeeDAO.GetByUserId(db, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	result, errModel := dao.BenefitDAO.InitiateListByBenefitType(db, searchByParam, repository.Benefit{
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
		ValidOperator: applicationModel.GetListEmployeeReimbursementTypesValidOperator,
		CountData:     result,
	}
	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	return
}

func (input employeeService) GetEmployeeReimbursementTypes(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		db = serverconfig.ServerAttribute.DBConnection
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validOrderBy = []string{"id"}
	)

	inputStruct, searchByParam, errModel = input.ReadAndValidateGetListData(request, []string{}, validOrderBy, applicationModel.GetListEmployeeReimbursementTypesValidOperator, service.DefaultLimit)
	if errModel.Error != nil {
		return
	}

	employee, errModel := dao.EmployeeDAO.GetByUserId(db, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	results, errModel := dao.BenefitDAO.GetListByEmployeeLevelIdAndEmployeeGradeId(db, inputStruct, searchByParam, repository.Benefit{
		EmployeeGradeId: sql.NullInt64{Int64: employee.GradeID.Int64},
		EmployeeLevelId: sql.NullInt64{Int64: employee.LevelID.Int64},
	})
	if errModel.Error != nil {
		return
	}

	output.Data.Content = input.toEmployeeReimbursementTypesResponse(results)
	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) toEmployeeReimbursementTypesResponse(data []interface{}) (result []out.EmployeeReimbursementType) {
	for _, item := range data {
		benefit, _ := item.(repository.Benefit)

		result = append(result, out.EmployeeReimbursementType{
			Id:          benefit.ID.Int64,
			BenefitName: benefit.BenefitName.String,
			BenefitType: benefit.BenefitType.String,
			Description: benefit.Description.String,
			CreatedAt:   benefit.CreatedAt.Time,
			UpdatedAt:   benefit.UpdatedAt.Time,
		})
	}

	return
}
