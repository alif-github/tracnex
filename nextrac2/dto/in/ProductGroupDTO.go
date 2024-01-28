package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type ProductGroupRequest struct {
	AbstractDTO
	ID               int64  `json:"id"`
	ProductGroupName string `json:"product_group_name"`
	UpdatedAtStr     string `json:"updated_at"`
	UpdatedAt        time.Time
}

func (input *ProductGroupRequest) ValidateView() (err errorModel.ErrorModel) {
	funcName := "ValidateView"
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(ProductGroupDTOFileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *ProductGroupRequest) ValidateDelete() (err errorModel.ErrorModel) {
	funcName := "ValidateDelete"
	return input.validationForUpdateAndDelete(ProductGroupDTOFileName, funcName)
}

func (input *ProductGroupRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdate"

	err = input.mandatoryFieldValidation(ProductGroupDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	err = input.validationForUpdateAndDelete(ProductGroupDTOFileName, funcName)
	if err.Error != nil {
		return
	}
	return
}

func (input *ProductGroupRequest) ValidateInsert() (err errorModel.ErrorModel) {
	funcName := "ValidateInsert"
	return input.mandatoryFieldValidation(ProductGroupDTOFileName, funcName)
}

func (input *ProductGroupRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	if util.IsStringEmpty(input.UpdatedAtStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *ProductGroupRequest) mandatoryFieldValidation(fileName string, funcName string) (err errorModel.ErrorModel) {

	err = input.ValidateMinMaxString(input.ProductGroupName, constanta.ProductGroupName, 1, 22)
	if err.Error != nil {
		return
	}

	err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.ProductGroupName, input.ProductGroupName)
	if err.Error != nil {
		return
	}

	return
}
