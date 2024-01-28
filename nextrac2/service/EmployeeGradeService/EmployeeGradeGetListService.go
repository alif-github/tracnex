package EmployeeGradeService

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

type employeeGradeGetListService struct {
	service.GetListData
	FileName string
}

var EmployeeGradeGetListService = employeeGradeGetListService{}.New()

func (input employeeGradeGetListService) New() (output employeeGradeGetListService) {
	output.FileName = "EmployeeGradeGetListService.go"
	output.ValidSearchBy = []string{"id"}
	output.ValidOrderBy = []string{"id"}
	output.ValidLimit = service.DefaultLimit
	return
}

func (input employeeGradeGetListService) GetEmployeeGrade(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
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

	levels, err := dao.EmployeeGradeDAO.GetEmployeeGrade(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.convertRepoToDTO(levels)
	output.Status = service.GetResponseMessages("SUCCESS_GET_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeGradeGetListService) InitiateEmployeeGrade(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	_, _, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListEmployeeValidOperator)
	if err.Error != nil {
		return
	}

	countData, err := dao.EmployeeGradeDAO.GetCountEmployeeGrade(serverconfig.ServerAttribute.DBConnection)
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

func (input employeeGradeGetListService) convertRepoToDTO(data []interface{}) (grades []out.EmployeeGradeResponse) {
	for _, item := range data {
		grade := item.(repository.EmployeeGradeModel)
		grades = append(grades, out.EmployeeGradeResponse{
			ID:                grade.ID.Int64,
			Grade:             grade.Grade.String,
			Description:       grade.Description.String,
			UpdatedAt:         grade.UpdatedAt.Time,
		})
	}
	return
}