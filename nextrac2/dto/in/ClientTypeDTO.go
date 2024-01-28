package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type ClientTypeRequest struct {
	AbstractDTO
	ID                 int64   `json:"id"`
	ClientType         string  `json:"client_type"`
	Description        string  `json:"description"`
	ParentClientTypeID int64   `json:"parent_client_type_id"`
	Remarks            []int64 `json:"remarks"`
	UpdatedAtStr       string  `json:"updated_at"`
	UpdatedAt          time.Time
}

func (input *ClientTypeRequest) ValidateInsert() (err errorModel.ErrorModel) {
	var (
		funcName = "ValidateInsert"
		fileName = "ClientTypeDTO.go"
	)

	err = input.mandatoryFieldValidation(fileName, funcName)
	if err.Error != nil {
		return
	}

	if !util.IsStringEmpty(input.Description) {
		err = input.ValidateMinMaxString(input.Description, constanta.Description, 1, 200)
		if err.Error != nil {
			return
		}

		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.Description, input.Description)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *ClientTypeRequest) mandatoryFieldValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	if util.IsStringEmpty(input.ClientType) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientType)
	}

	err = input.ValidateMinMaxString(input.ClientType, constanta.ClientType, 1, 22)
	if err.Error != nil {
		return
	}

	err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.ClientType, input.ClientType)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *ClientTypeRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	var (
		funcName = "ValidateUpdate"
		fileName = "ClientTypeDTO.go"
	)

	err = input.validationForUpdateAndDelete(fileName, funcName)
	if err.Error != nil {
		return
	}

	err = input.mandatoryFieldValidation(fileName, funcName)
	if err.Error != nil {
		return
	}

	if !util.IsStringEmpty(input.Description) {
		err = input.ValidateMinMaxString(input.Description, constanta.Description, 1, 200)
		if err.Error != nil {
			return
		}

		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.Description, input.Description)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *ClientTypeRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
	if input.ID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ID)
		return
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

func (input *ClientTypeRequest) ValidateDelete() (err errorModel.ErrorModel) {
	var (
		fileName = "ClientTypeDTO.go"
		funcName = "ValidateDelete"
	)

	return input.validationForUpdateAndDelete(fileName, funcName)
}

func (input *ClientTypeRequest) ValidateView() (err errorModel.ErrorModel) {
	var (
		funcName = "ValidateView"
		fileName = "ClientTypeDTO.go"
	)

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
