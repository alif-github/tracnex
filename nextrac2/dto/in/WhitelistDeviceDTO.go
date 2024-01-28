package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type WhiteListDeviceRequest struct {
	AbstractDTO
	ID           int64  `json:"id"`
	Device       string `json:"device"`
	Description  string `json:"description"`
	UpdatedAtStr string `json:"updated_at"`
	UpdatedAt    time.Time
}

const whiteListDTOFileName = "WhiteListDTO.go"

func (input *WhiteListDeviceRequest) ValidateForInsert() (err errorModel.ErrorModel) {
	funcName := "ValidateForInsert"
	return input.mandatoryFieldValidation(whiteListDTOFileName, funcName)
}

func (input *WhiteListDeviceRequest) mandatoryFieldValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	//--- Field Device
	if util.IsStringEmpty(input.Device) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Device)
	}

	err = input.ValidateMinMaxString(input.Device, constanta.Device, 1, 200)
	if err.Error != nil {
		return
	}

	err = util2.ValidateWhiteListSpecialCharacter(fileName, funcName, constanta.Device, input.Device)
	if err.Error != nil {
		return
	}

	//--- Field Description, open validation empty string
	//if util.IsStringEmpty(input.Description) {
	//	return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Description)
	//}

	err = input.ValidateMinMaxString(input.Description, constanta.Description, 0, 200)
	if err.Error != nil {
		return
	}

	//err = util2.ValidateWhiteListSpecialCharacter(fileName, funcName, constanta.Description, input.Description)
	//if err.Error != nil {
	//	return
	//}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *WhiteListDeviceRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdate"

	err = input.mandatoryFieldValidation(whiteListDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	err = input.validationForUpdateAndDelete(ProductGroupDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *WhiteListDeviceRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
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

func (input *WhiteListDeviceRequest) ValidateDelete() (err errorModel.ErrorModel) {
	var (
		funcName = "ValidateDelete"
	)

	return input.validationForUpdateAndDelete(whiteListDTOFileName, funcName)
}

func (input *WhiteListDeviceRequest) ValidateView() (err errorModel.ErrorModel) {
	funcName := "ValidateView"
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(whiteListDTOFileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
