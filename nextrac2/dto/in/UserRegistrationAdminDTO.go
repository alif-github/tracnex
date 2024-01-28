package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
)

type UserRegistrationAdminRequest struct {
	ID            int64  `json:"id"`
	UniqueID1     string `json:"unique_id_1"`
	UniqueID2     string `json:"unique_id_2"`
	UserAdmin     string `json:"user_admin"`
	PasswordAdmin string `json:"password_admin"`
	CompanyName   string `json:"company_name"`
	BranchName    string `json:"branch_name"`
	ClientID      string `json:"client_id"`
	ClientTypeID  int64  `json:"client_type_id"`
}

func (input UserRegistrationAdminRequest) ValidateViewUserRegistrationAdmin() (err errorModel.ErrorModel) {
	fileName := "UserRegistrationAdminDTO.go"
	funcName := "ValidateViewUserRegistrationAdmin"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ID)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input UserRegistrationAdminRequest) ValidateInsertUserRegistrationAdmin() (err errorModel.ErrorModel) {
	fileName := "UserRegistrationAdminDTO.go"
	funcName := "ValidateInsertUserRegistrationAdmin"

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

	if util.IsStringEmpty(input.UserAdmin) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UserAdmin)
	}

	err = util2.ValidateMinMaxString(input.UserAdmin, constanta.UserAdmin, 1, 100)
	if err.Error != nil {
		return err
	}

	if util.IsStringEmpty(input.PasswordAdmin) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.PasswordAdmin)
	}

	err = util2.ValidateMinMaxString(input.PasswordAdmin, constanta.PasswordAdmin, 1, 100)
	if err.Error != nil {
		return err
	}

	if !util.IsStringEmpty(input.CompanyName) {
		err = util2.ValidateMinMaxString(input.CompanyName, constanta.CompanyName, 1, 100)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.BranchName) {
		err = util2.ValidateMinMaxString(input.BranchName, constanta.BranchName, 1, 100)
		if err.Error != nil {
			return err
		}
	}

	if util.IsStringEmpty(input.ClientID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientMappingClientID)
	}

	if input.ClientTypeID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.NewClientType)
	}

	return errorModel.GenerateNonErrorModel()
}
