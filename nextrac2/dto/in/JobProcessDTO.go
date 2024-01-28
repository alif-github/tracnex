package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type JobProcessRequest struct {
	JobID	string	`json:"job_id"`
}

func (input *JobProcessRequest) ViewDetailJobProcess() errorModel.ErrorModel {
	fileName := "JobProcessDTO.go"
	funcName := "ViewDetailJobProcess"

	if util.IsStringEmpty(input.JobID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.JobID)
	}

	return errorModel.GenerateNonErrorModel()
}