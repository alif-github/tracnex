package EmployeeService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input employeeService) GetMedicalRemainingBalance(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	benefit, errModel := dao.EmployeeBenefitsDAO.GetByUserId(serverconfig.ServerAttribute.DBConnection, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	output.Data.Content = out.EmployeeReimbursementBenefit{
		CurrentMedicalValue: benefit.CurrentMedicalValue.Float64,
	}
	return
}