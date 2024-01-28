package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type DetailClientMappingDTO struct {
	AbstractDTO
	ClientTypeID             int64      `json:"client_type_id"`
	CompanyID                string     `json:"company_id"`
	BranchID                 string     `json:"branch_id"`
}

func (input *DetailClientMappingDTO) ValidateDetailClientMapping() errorModel.ErrorModel {
	err := input.validateMandatoryField()
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *DetailClientMappingDTO) validateMandatoryField() errorModel.ErrorModel {
	fileName := "DetailClientMappingDTO.go"
	funcName := "validateMandatoryField"
	var validationResult bool
	var err errorModel.ErrorModel

	if input.ClientTypeID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientTypeID)
	}

	validationResult = util.IsStringEmpty(input.CompanyID)
	if validationResult {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CompanyID)
	}

	err = input.ValidateMinMaxString(input.CompanyID, constanta.CompanyID, 1, 20)
	if err.Error != nil {
		return err
	}

	validationResult = util.IsStringEmpty(input.BranchID)
	if validationResult {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.BranchID)
	}

	err = input.ValidateMinMaxString(input.BranchID, constanta.BranchID, 1, 20)
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}