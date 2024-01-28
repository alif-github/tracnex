package ParameterService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
)

type employeeParameterViewService struct {
	service.GetListData
	FileName string
}

var EmployeeParameterViewService = employeeParameterViewService{FileName: "EmployeeParameterViewService.go"}

func (input employeeParameterViewService) ViewParameterEmployee(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	levels, err := dao.ParameterDAO.GetParameterForEmployee(serverconfig.ServerAttribute.DBConnection)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.convertRepoToDTO(levels)
	output.Status = service.GetResponseMessages("SUCCESS_GET_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeParameterViewService) convertRepoToDTO(data []repository.ParameterModel) (parameters []out.ParameterForView) {
	for _, parameter := range data {
		parameters = append(parameters, out.ParameterForView{
			ID:                parameter.ID.Int64,
			Name:              parameter.Name.String,
			Value:             parameter.Value.String,
			UpdatedAt:         parameter.UpdatedAt.Time,
		})
	}
	return
}
