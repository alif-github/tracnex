package EmployeeService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input employeeService) GetRemainingLeave(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	benefit, errModel := dao.EmployeeBenefitsDAO.GetByUserId(serverconfig.ServerAttribute.DBConnection, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	totalLeave := benefit.CurrentAnnualLeave.Int64 + benefit.LastAnnualLeave.Int64
	currentAnnualLeave := benefit.CurrentAnnualLeave.Int64
	negativeLeave := int64(0)

	if totalLeave < 0 {
		totalLeave = 0
	}

	if currentAnnualLeave < 0 {
		negativeLeave = currentAnnualLeave
		currentAnnualLeave = 0
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	output.Data.Content = out.EmployeeAnnualLeaveBenefit{
		TotalLeave:         totalLeave,
		CurrentAnnualLeave: currentAnnualLeave,
		LastAnnualLeave:    benefit.LastAnnualLeave.Int64,
		NegativeLeave:      negativeLeave,
	}
	return
}