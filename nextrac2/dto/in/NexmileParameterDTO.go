package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
)

type NexmileParameterRequest struct {
	ClientID      string           `json:"client_id"`
	ClientTypeID  int64            `json:"client_type_id"`
	UniqueID1     string           `json:"unique_id_1"`
	UniqueID2     string           `json:"unique_id_2"`
	ParameterData []ParameterValue `json:"parameter_data"`
}

type ParameterValue struct {
	ParameterID string `json:"parameter_id"`
	Value       string `json:"value"`
}

type NexmileParameterRequestForView struct {
	ClientID     string `json:"client_id"`
	AuthUserID   int64  `json:"auth_user_id"`
	UserId       string `json:"user_id"`
	AndroidId    string `json:"android_id"`
	ClientTypeId int64  `json:"client_type_id"`
	Password     string `json:"password"`
}

func (input NexmileParameterRequest) ValidateInsertNexmileParameter() errorModel.ErrorModel {
	var (
		fileName = "NexmileParameterDTO.go"
		funcName = "ValidateInsertNexmileParameter"
		err      errorModel.ErrorModel
	)

	if util.IsStringEmpty(input.ClientID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientMappingClientID)
	}

	err = util2.ValidateMinMaxString(input.ClientID, constanta.ClientMappingClientID, 1, 256)
	if err.Error != nil {
		return err
	}

	if input.ClientTypeID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ClientMappingAuthUserID)
	}

	if util.IsStringEmpty(input.UniqueID1) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UniqueID1)
	}

	err = util2.ValidateMinMaxString(input.UniqueID1, constanta.UniqueID1, 1, 20)
	if err.Error != nil {
		return err
	}

	if !util.IsStringEmpty(input.UniqueID2) {
		err = util2.ValidateMinMaxString(input.UniqueID2, constanta.UniqueID2, 1, 20)
		if err.Error != nil {
			return err
		}
	}

	for _, itemParameter := range input.ParameterData {
		if util.IsStringEmpty(itemParameter.ParameterID) {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ParameterID)
		}

		err = util2.ValidateMinMaxString(itemParameter.ParameterID, constanta.ParameterID, 1, 250)
		if err.Error != nil {
			return err
		}

		if !util.IsStringEmpty(itemParameter.Value) {
			err = util2.ValidateMinMaxString(itemParameter.ParameterID, constanta.ParameterID, 1, 250)
			if err.Error != nil {
				return err
			}
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input NexmileParameterRequestForView) ValidateViewNexmileParameter() errorModel.ErrorModel {
	fileName := "NexmileParameterDTO.go"
	funcName := "ValidateViewNexmileParameter"

	if util.IsStringEmpty(input.UserId) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UserID)
	}

	if util.IsStringEmpty(input.Password) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Password)
	}

	if util.IsStringEmpty(input.AndroidId) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.AndroidID)
	}

	if input.ClientTypeId < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ClientTypeID)
	}

	return errorModel.GenerateNonErrorModel()
}
