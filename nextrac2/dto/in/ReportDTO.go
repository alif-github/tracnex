package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type ErrorBundleReport struct {
	RedmineNumber int64  `json:"redmine_number"`
	Message       string `json:"message"`
}

type ReportRequest struct {
	AbstractDTO
	order    string
	filter   string
	employee string
}

type ReportHistory struct {
	ID            int64
	SuccessTicket []int64
	DataActual    string
	DepartmentID  int64
}

func (input *ReportRequest) ValidateUpdatePayment(validLimit []int, validOrder []string) {
	input.ValidateInputPageLimitAndOrderBy(validLimit, validOrder)
}

func (input *ReportHistory) ValidateView() errorModel.ErrorModel {
	var (
		fileName = "ReportDTO.go"
		funcName = "ValidateView"
	)

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	return errorModel.GenerateNonErrorModel()
}
