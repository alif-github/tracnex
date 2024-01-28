package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type DataGroupRequest struct {
	AbstractDTO
	ID           int64     `json:"id"`
	GroupID      string    `json:"group_id"`
	Description  string    `json:"description"`
	Scope        []string  `json:"scope"`
	UpdatedAt    time.Time `json:"-"`
	UpdatedAtStr string    `json:"updated_at"`
}

func (input *DataGroupRequest) ValidateInsertDataGroup() (err errorModel.ErrorModel) {
	funcName := "ValidateInsertDataGroup"

	return input.mandatoryValidation(DataGroupDTOFileName, funcName)
}

func (input *DataGroupRequest) ValidateUpdateDataGroup() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdateDataGroup"

	err = input.validationForUpdateAndDelete(DataGroupDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	err = input.ValidateMinMaxString(input.Description, constanta.Description, 1, 200)
	if err.Error != nil {
		return
	}

	return
}

func (input *DataGroupRequest) ValidateDeleteDataGroup() (err errorModel.ErrorModel) {
	funcName := "ValidateDeleteDataGroup"

	return input.validationForUpdateAndDelete(DataGroupDTOFileName, funcName)
}

func (input *DataGroupRequest) ValidateViewDataGroup() (err errorModel.ErrorModel) {
	funcName := "ValidateViewDataGroup"

	if input.ID < 0 {
		return errorModel.GenerateUnknownDataError(DataGroupDTOFileName, funcName, constanta.ID)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *DataGroupRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
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

func (input *DataGroupRequest) mandatoryValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	err = input.ValidateMinMaxString(input.GroupID, constanta.Group, 1, 23)
	if err.Error != nil {
		return
	}

	validationResult, errField := util.IsNexsoftProfileNameStandardValid(input.GroupID)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.Group, "")
	}

	//err = input.ValidateIsContainSpaceString(fileName, funcName, constanta.Group, input.GroupID)
	//if err.Error != nil {
	//	return
	//}

	err = input.ValidateMinMaxString(input.Description, constanta.Description, 1, 200)
	if err.Error != nil {
		return
	}

	err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.Description, input.Description)
	if err.Error != nil {
		return
	}

	if len(input.Scope) == 0 {
		return errorModel.GenerateEmptyFieldError(DataGroupDTOFileName, funcName, constanta.DataScope)
	}

	return errorModel.GenerateNonErrorModel()
}
