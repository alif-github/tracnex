package EmployeeService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strings"
)

type employeeLeaveDetailService struct {
	FileName string
	service.AbstractService
}

var EmployeeLeaveDetailService = employeeLeaveDetailService{FileName: "EmployeeLeaveDetailService.go"}

func (input employeeLeaveDetailService) DetailEmployeeLeave(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "DetailEmployeeLeave"

	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	leave, err := dao.EmployeeLeaveDAO.GetDetailEmployeeLeave(serverconfig.ServerAttribute.DBConnection, id)
	if err.Error != nil {
		return
	}

	if leave.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employee_leave")
		return
	}

	output.Data.Content = input.getResponse(leave)
	output.Status = service.GetResponseMessages("SUCCESS_GET_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeLeaveDetailService) getResponse(model repository.EmployeeLeaveModel) out.EmployeeLeave {

	var date []string
	_ = json.Unmarshal([]byte(model.StrDateList.String), &date)

	separatedPaths := strings.Split(model.Path.String, "/")
	fileName := separatedPaths[len(separatedPaths) - 1]

	return out.EmployeeLeave{
		ID:            model.ID.Int64,
		IDCard:        model.IDCard.String,
		FirstName:     model.Firstname.String,
		LastName:      model.Lastname.String,
		FullName:      model.Firstname.String + " " + model.Lastname.String,
		Department:    model.Department.String,
		AllowanceName: model.AllowanceName.String,
		Value:         model.Value.Int64,
		Type :         model.Type.String,
		Status: 	   model.Status.String,
		StartDate:     model.StartDate.Time,
		EndDate:       model.EndDate.Time,
		LeaveTime:     model.LeaveTime.Time,
		CancellationReason: model.CancellationReason.String,
        CurrentAnnualLeave: model.CurrentAnnualLeave.Int64,
        LastAnnualLeave:    model.LastAnnualLeave.Int64,
        Date:         date,
        Filename:     fileName,
        Attachment:   model.Host.String + model.Path.String,
        Description:  model.Description.String,
	}
}