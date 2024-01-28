package EmployeeAllowanceService

import (
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

type employeeAllowanceGetListService struct {
	service.GetListData
	FileName string
}

var EmployeeAllowanceGetListService = employeeAllowanceGetListService{}.New()

func (input employeeAllowanceGetListService) New() (output employeeAllowanceGetListService) {
	output.FileName = "EmployeeAllowanceGetListService.go"
	output.ValidSearchBy = []string{"id"}
	output.ValidOrderBy = []string{"id"}
	output.ValidLimit = service.DefaultLimit
	return
}

func (input employeeAllowanceGetListService) GetEmployeeAllowance(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam
	var createdBy int64 = 0

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListEmployeeValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	createdBy, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*context)
	if !isOnlyHaveOwnAccess {
		createdBy = 0
	}

	allowances, err := dao.EmployeeAllowanceDAO.GetEmployeeAllowance(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.convertRepoToDTO(allowances)
	output.Status = service.GetResponseMessages("SUCCESS_GET_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeAllowanceGetListService) InitiateEmployeeAllowance(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	_, _, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListEmployeeValidOperator)
	if err.Error != nil {
		return
	}

	countData, err := dao.EmployeeAllowanceDAO.GetCountEmployeeAllowance(serverconfig.ServerAttribute.DBConnection)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		CountData:     int(countData),
		ValidOperator: applicationModel.GetListEmployeeValidOperator,
	}
	output.Status = service.GetResponseMessages("SUCCESS_INITIATE_MESSAGE", context)

	return
}

func (input employeeAllowanceGetListService) convertRepoToDTO(data []interface{}) (allowances []out.EmployeeAllowanceResponse) {
	for _, item := range data {
		allowance := item.(repository.EmpAllowanceModel)
		allowances = append(allowances, out.EmployeeAllowanceResponse{
			ID:                allowance.ID.Int64,
			AllowanceName:     allowance.AllowanceName.String,
			AllowanceType:     allowance.Type.String,
			UpdatedAt:         allowance.UpdatedAt.Time,
		})
	}
	return
}
