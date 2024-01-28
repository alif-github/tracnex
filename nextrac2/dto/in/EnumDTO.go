package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type EnumRequest struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

func (input *EnumRequest) ValidateView() (err errorModel.ErrorModel) {
	var (
		funcName = "ValidateView"
	)

	if util.IsStringEmpty(input.Type) {
		return errorModel.GenerateUnknownDataError("EnumDTO.go", funcName, constanta.Type)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
