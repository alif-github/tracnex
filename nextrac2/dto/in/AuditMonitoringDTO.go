package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type AuditMonitoringRequest struct {
	AbstractDTO
	ID        int64 `json:"id"`
}

func (input *AuditMonitoringRequest) ValidateViewAuditMonitoring() (err errorModel.ErrorModel) {
	fileName := "AuditMonitoringDTO.go"
	funcName := "ValidateViewAuditMonitoring"

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	return errorModel.GenerateNonErrorModel()
}
