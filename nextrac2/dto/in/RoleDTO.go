package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"regexp"
	"time"
)

type RoleRequest struct {
	AbstractDTO
	ID           int64    `json:"id"`
	RoleID       string   `json:"role_id"`
	Description  string   `json:"description"`
	Permission   []string `json:"permission"`
	Level        int      `json:"level"`
	UpdatedAtStr string   `json:"updated_at"`
	UpdatedAt    time.Time
}

func (input *RoleRequest) ValidateInsertRole() (err errorModel.ErrorModel) {
	funcName := "ValidateInsertRole"

	return input.mandatoryValidation(RoleDTOFilename, funcName)
}

func (input *RoleRequest) mandatoryValidation(fileName string, funcName string) (err errorModel.ErrorModel) {

	//---------- Check is string empty for role ID
	if util.IsStringEmpty(input.RoleID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Role)
	}

	//---------- Role minimum 5 and maximum 50
	err = input.ValidateMinMaxString(input.RoleID, constanta.Role, 5, 50)
	if err.Error != nil {
		return
	}

	//---------- Role format regex
	validationResult, errField := input.isNexsoftRoleIDStandarValid(input.RoleID)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.Role, "")
	}

	//---------- Check is string empty description
	if util.IsStringEmpty(input.Description) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Description)
	}

	//---------- Description minimum 1 and maximum 256
	err = input.ValidateMinMaxString(input.Description, constanta.Description, 1, 256)
	if err.Error != nil {
		return
	}

	err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.Description, input.Description)
	if err.Error != nil {
		return
	}

	//---------- Len permission must more than 0
	if len(input.Permission) == 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Permission)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *RoleRequest) ValidationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
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

func (input *RoleRequest) ValidateDeleteRole() (err errorModel.ErrorModel) {
	funcName := "ValidateDeleteRole"

	return input.ValidationForUpdateAndDelete(RoleDTOFilename, funcName)
}

func (input *RoleRequest) ValidateUpdateRole() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdateRole"

	err = input.ValidationForUpdateAndDelete(RoleDTOFilename, funcName)
	if err.Error != nil {
		return
	}

	return input.mandatoryValidation(RoleDTOFilename, funcName)
}

func (input *RoleRequest) ValidateViewRole() (err errorModel.ErrorModel) {
	fileName := "RoleDTO.go"
	funcName := "ValidateViewRole"

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *RoleRequest) isNexsoftRoleIDStandarValid(profileName string) (bool, string) {
	NameOrTitle := regexp.MustCompile("^[A-Z0-9](?:|(?:[a-z0-9]+|(?:[a-z0-9]|[a-z0-9])(?:([_-]|)[a-z0-9])+)|[ ]([A-Z0-9](?:|(?:[a-z0-9]+|(?:[a-z0-9]|[a-z0-9])(?:([_-]|)[a-z0-9])+))+)+)+$")
	return NameOrTitle.MatchString(profileName), "PROFILE_NAME_REGEX_MESSAGE"
}
