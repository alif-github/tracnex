package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type AddResourceNexcloud struct {
	FirstName	string	`json:"first_name"`
	LastName	string	`json:"last_name"`
	ClientID	string	`json:"client_id"`
	ClientAlias	string	`json:"client_alias"`
}

type AddResourceNexdrive struct {
	FirstName	string	`json:"first_name"`
	LastName	string	`json:"last_name"`
	ClientID	string	`json:"client_id"`
	ClientAlias	string	`json:"client_alias"`
}

type AddResourceExternalRequest struct {
	AbstractDTO
	BodyRequest
	ID					int64		`json:"id"`
	ClientTypeID		int64		`json:"client_type_id"`
	ClientID			string		`json:"client_id"`
	OldResource			string		`json:"old_resource"`
	UpdatedAtStr		string		`json:"updated_at"`
	UpdatedAt			time.Time
}

func (input *AddResourceExternalRequest) ValidateAddResourceExternal() errorModel.ErrorModel {
	return input.validateMandatoryAddResourceExternal()
}

func (input *AddResourceExternalRequest) validateMandatoryAddResourceExternal() errorModel.ErrorModel {
	fileName := "AddResourceNexcloudDTOIn.go"
	funcName := "validateMandatoryAddResourceExternal"
	var validationResult bool
	var err errorModel.ErrorModel

	validationResult = util.IsStringEmpty(input.ClientID)
	if validationResult {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientID)
	}
	err = input.ValidateMinMaxString(input.ClientID, constanta.ClientID, 1, 256)
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *AddResourceExternalRequest) ValidateViewLogForAddResource() errorModel.ErrorModel {
	fileName := "AddResourceExternalDTO.go"
	funcName := "ValidateViewLogForAddResource"

	if util.IsStringEmpty(input.ClientID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientID)
	}

	return errorModel.GenerateNonErrorModel()
}