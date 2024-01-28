package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type ClientMappingForUIRequest struct {
	AbstractDTO
	ID              int64     `json:"id"`
	ClientTypeId    int64     `json:"client_type_id"`
	ClientId        string    `json:"client_id"`
	SocketID        string    `json:"socket_id"`
	SocketPassword  string    `json:"-"`
	CompanyID       string    `json:"company_id"`
	BranchID        string    `json:"branch_id"`
	ClientAlias     string    `json:"client_alias"`
	UpdatedAtString string    `json:"updated_at"`
	UpdatedAt       time.Time `json:"-"`
}

func (input *ClientMappingForUIRequest) ValidateUpdateClientID() (err errorModel.ErrorModel) {
	var (
		fileName         = "ClientMappingForUIDTO.go"
		funcName         = "ValidateUpdateClientID"
		validationResult bool
		errField         string
		additionalInfo   string
	)

	if input.ClientTypeId < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, "client_type_id")
		return
	}

	validationResult = util.IsStringEmpty(input.SocketID)
	if !validationResult {
		err = input.ValidateMinMaxString(input.SocketID, "socket_id", 1, 50)
		if err.Error != nil {
			return
		}

		validationResult, errField, additionalInfo = IsOnlyWordCharacterValid(input.SocketID)
		if !validationResult {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, "socket_id", additionalInfo)
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input ClientMappingForUIRequest) ValidateViewCLientMapping() (err errorModel.ErrorModel) {
	var (
		fileName = "ClientMappingForUIDTO.go"
		funcName = "ValidateViewCLientMapping"
	)

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *ClientMappingForUIRequest) ValidateUpdateClientName() (err errorModel.ErrorModel) {
	var (
		fileName = "ClientMappingForUIDTO.go"
		funcName = "ValidateUpdateClientName"
	)

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	if util.IsStringEmpty(input.UpdatedAtString) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtString, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	if !util.IsStringEmpty(input.ClientAlias) {
		err = input.ValidateMinMaxString(input.ClientAlias, constanta.ClientAlias, 1, 50)
		if err.Error != nil {
			return
		}

		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.ClientAlias, input.ClientAlias)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
