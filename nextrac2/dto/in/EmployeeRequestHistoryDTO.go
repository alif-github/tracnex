package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

type EmployeeRequestHistory struct {
	ID                 int64     `json:"id"`
	Status             string    `json:"-"`
	RequestType        string    `json:"request_type"`
	CancellationReason string    `json:"cancellation_reason"`
	StrUpdatedAt       string    `json:"updated_at"`
	UpdatedAt          time.Time `json:"-"`
}

func (input *EmployeeRequestHistory) ValidateCancellation() errorModel.ErrorModel {
	fileName := "EmployeeRequestHistoryDTO.go"
	funcName := "ValidateCancellation"

	if errModel := util.ValidateMinMax(input.CancellationReason, "cancellation_reason", 1, 256); errModel.Error != nil {
		return errModel
	}

	if !input.isRequestTypeValid(input.RequestType) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "request_type")
	}

	updatedAt, errModel := TimeStrToTime(input.StrUpdatedAt, "updated_at")
	if errModel.Error != nil {
		return errModel
	}

	input.UpdatedAt = updatedAt

	return errorModel.GenerateNonErrorModel()
}

func (input *EmployeeRequestHistory) isRequestTypeValid(requestType string) bool {
	return requestType == constanta.LeaveType ||
		requestType == constanta.PermitType ||
		requestType == constanta.SickLeaveType ||
		requestType == constanta.ReimbursementType
}