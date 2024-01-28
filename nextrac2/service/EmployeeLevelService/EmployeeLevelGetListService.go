package EmployeeLevelService

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

type employeeLevelGetListService struct {
	service.GetListData
	FileName string
}

var EmployeeLevelGetListService = employeeLevelGetListService{}.New()

func (input employeeLevelGetListService) New() (output employeeLevelGetListService) {
	output.FileName = "EmployeeLevelGetListService.go"
	output.ValidSearchBy = []string{"id"}
	output.ValidOrderBy = []string{"id"}
	output.ValidLimit = service.DefaultLimit
	return
}

func (input employeeLevelGetListService) GetEmployeeLevel(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct         in.GetListDataDTO
		searchByParam       []in.SearchByParam
		isOnlyHaveOwnAccess bool
		createdBy           int64
		levels              []interface{}
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListEmployeeValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	createdBy, isOnlyHaveOwnAccess = service.CheckIsOnlyHaveOwnPermission(*context)
	if !isOnlyHaveOwnAccess {
		createdBy = 0
	}

	levels, err = dao.EmployeeLevelDAO.GetEmployeeLevel(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	levels, err = dao.EmployeeLevelDAO.GetEmployeeLevelMatrix(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.convertRepoToDTO(levels)
	output.Status = service.GetResponseMessages("SUCCESS_GET_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeLevelGetListService) InitiateEmployeeLevel(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	_, _, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListEmployeeValidOperator)
	if err.Error != nil {
		return
	}

	countData, err := dao.EmployeeLevelDAO.GetCountEmployeeLevel(serverconfig.ServerAttribute.DBConnection)
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

func (input employeeLevelGetListService) convertRepoToDTO(data []interface{}) (levels []out.EmployeeLevelResponse) {
	for _, item := range data {
		level := item.(repository.EmployeeLevelModel)
		levels = append(levels, out.EmployeeLevelResponse{
			ID:          level.ID.Int64,
			Level:       level.Level.String,
			Description: level.Description.String,
			UpdatedAt:   level.UpdatedAt.Time,
		})
	}
	return
}
