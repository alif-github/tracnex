package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type EmployeeApprovalRequest struct {
	Id           int64     `json:"-"`
	RequestType  string    `json:"request_type"`
	Status       string    `json:"status"`
	StrUpdatedAt string    `json:"updated_at"`
	UpdatedAt    time.Time `json:"-"`
}

func (input *EmployeeApprovalRequest) ValidateApprovalRequest() (errModel errorModel.ErrorModel) {
	fileName := "EmployeeApprovalDTO.go"
	funcName := "ValidateApprovalRequest"

	input.UpdatedAt, errModel = TimeStrToTime(input.StrUpdatedAt, "updated_at")
	if errModel.Error != nil {
		return
	}

	if !input.isStatusValid() {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "status")
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *EmployeeApprovalRequest) isStatusValid() bool {
	return input.Status == constanta.ApprovedRequestStatus ||
		input.Status == constanta.RejectedRequestStatus
}

//func (input *EmployeeApprovalRequest) isRequestTypeValid() bool {
//	return input.RequestType == constanta.LeaveAllowanceType ||
//		input.RequestType == constanta.AnnualLeaveAllowanceType ||
//		input.RequestType == constanta.SickLeaveAllowanceType ||
//		input.RequestType == constanta.PermitAllowanceType
//}
