package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type ScopeRequest struct {
	ScopeType string `json:"scope_type"`
	ScopeID   int64  `json:"scope_id"`
}

func (input ScopeRequest) ValidateInsert() errorModel.ErrorModel {
	var (
		fileName = "ScopeDTO.go"
		funcName = "ValidateInsert"
	)

	//--- Scope Type
	if util.IsStringEmpty(input.ScopeType) {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Scope)
	}

	return errorModel.GenerateNonErrorModel()
}
