package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type CustomerGroupRequest struct {
	AbstractDTO
	ID                int64  `json:"id"`
	CustomerGroupID   string `json:"customer_group_id"`
	CustomerGroupName string `json:"customer_group_name"`
	UpdatedAtStr      string `json:"updated_at"`
	UpdatedAt         time.Time
}

func (input *CustomerGroupRequest) ValidateViewCustomerGroup() (err errorModel.ErrorModel) {
	funcName := "ValidateViewCustomerGroup"
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(CustomerGroupDTOFileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *CustomerGroupRequest) ValidateInsertCustomerGroup() (err errorModel.ErrorModel) {
	funcName := "ValidateInsertCustomerGroup"
	return input.mandatoryValidation(CustomerGroupDTOFileName, funcName)
}

func (input *CustomerGroupRequest) ValidateUpdateCustomerGroup() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdateCustomerGroup"

	//---------- Check is string empty for Customer Group Name
	if util.IsStringEmpty(input.CustomerGroupName) {
		err = errorModel.GenerateEmptyFieldError(CustomerGroupDTOFileName, funcName, constanta.CustomerGroupName)
		return
	}
	err = input.ValidateMinMaxString(input.CustomerGroupName, constanta.CustomerGroupName, 1, 50)
	if err.Error != nil {
		return
	}

	if err = util2.ValidateSpecialCharacter(CustomerGroupDTOFileName, funcName, constanta.CustomerGroupName, input.CustomerGroupName); err.Error != nil {
		return
	}

	err = input.validationForUpdateAndDelete(CustomerGroupDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	return
}

func (input *CustomerGroupRequest) ValidateDeleteCustomerGroup() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdateCustomerGroup"

	err = input.validationForUpdateAndDelete(CustomerGroupDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	return
}

func (input *CustomerGroupRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
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

func (input *CustomerGroupRequest) mandatoryValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	//var isValid bool
	//var errField string
	//---------- Check is string empty for Customer Group ID
	if util.IsStringEmpty(input.CustomerGroupID) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CustomerGroupID)
		return
	}
	err = input.ValidateMinMaxString(input.CustomerGroupID, constanta.CustomerGroupID, 1, 22)
	if err.Error != nil {
		return
	}

	if err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.CustomerGroupID, input.CustomerGroupID); err.Error != nil {
		return
	}

	//---------- Check is string empty for Customer Group Name
	if util.IsStringEmpty(input.CustomerGroupName) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CustomerGroupName)
		return
	}
	err = input.ValidateMinMaxString(input.CustomerGroupName, constanta.CustomerGroupName, 1, 50)
	if err.Error != nil {
		return
	}

	if err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.CustomerGroupName, input.CustomerGroupName); err.Error != nil {
		return
	}

	return
}
