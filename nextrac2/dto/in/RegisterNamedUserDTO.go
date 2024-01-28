package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
	"time"
)

type RegisterNamedUserRequest struct {
	ParentClientID   string `json:"parent_client_id"`
	ClientID         string `json:"client_id"`
	ClientTypeID     int64  `json:"client_type_id"`
	ClientType       string `json:"clienr_type"`
	AuthUserID       int64  `json:"auth_user_id"`
	Firstname        string `json:"firstname"`
	Lastname         string `json:"lastname"`
	Username         string `json:"username"`
	UserID           string `json:"user_id"`
	Password         string `json:"password"`
	ClientAliases    string `json:"client_aliases"`
	SalesmanID       string `json:"salesman_id"`
	AndroidID        string `json:"android_id"`
	RegDateStr       string `json:"reg_date"`
	Email            string `json:"email"`
	UniqueID1        string `json:"unique_id_1"`
	UniqueID2        string `json:"unique_id_2"`
	NoTelp           string `json:"no_telp"`
	SalesmanCategory string `json:"salesman_category"`
	RegDate          time.Time
	UserType         int64  `json:"user_type"`
	Channel          string `json:"channel"`
	CompanyName      string `json:"company_name"`
	BranchName       string `json:"branch_name"`
	UserAdmin        string `json:"user_admin"`
	PasswordAdmin    string `json:"password_admin"`
	CountryCode      string `json:"country_code"`
}

type CheckNamedUserBeforeInsertRequest struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func (input CheckNamedUserBeforeInsertRequest) ValidateCheckNamedUserBeforeInsert() (err errorModel.ErrorModel) {
	var (
		funcName = "ValidateCheckNamedUserBeforeInsert"
		fileName = "RegisterNamedUserDTO.go"
	)

	if util.IsStringEmpty(input.Email) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Email)
	}

	if util.IsStringEmpty(input.Phone) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Phone)
	}

	if !IsEmailAddressMDB(input.Email) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
	}

	isPhoneValid := IsPhoneNumberWithCountryCodeMDB(input.Phone)
	if !isPhoneValid {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.PhoneRegex, constanta.Phone, "")
	}

	return
}

func (input *RegisterNamedUserRequest) ValidateRegisterOrRenewLicenseNamedUser() (err errorModel.ErrorModel) {
	var (
		funcName = "ValidateRegisterOrRenewLicenseNamedUser"
		fileName = input.fileNameGenerate()
	)

	if util.IsStringEmpty(input.ParentClientID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ParentClientID)
	}

	if input.ClientTypeID < 1 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientTypeID)
	}

	if util.IsStringEmpty(input.UniqueID1) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UniqueID1)
	}

	if util.IsStringEmpty(input.UserAdmin) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UserAdmin)
	}

	if !util.IsStringEmpty(input.UserAdmin) {
		err = util2.ValidateMinMaxString(input.UserAdmin, constanta.UserAdmin, 1, 100)
		if err.Error != nil {
			return err
		}
	}

	if util.IsStringEmpty(input.PasswordAdmin) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.PasswordAdmin)
	}

	if !util.IsStringEmpty(input.PasswordAdmin) {
		err = util2.ValidateMinMaxString(input.PasswordAdmin, constanta.PasswordAdmin, 1, 100)
		if err.Error != nil {
			return err
		}
	}

	if util.IsStringEmpty(input.Firstname) {
		err = util2.ValidateMinMaxString(input.Firstname, constanta.FirstName, 1, 50)
		if err.Error != nil {
			return err
		}

		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.FirstName, input.Firstname)
		if err.Error != nil {
			return err
		}
	}

	if util.IsStringEmpty(input.Password) {
		err = util2.ValidateMinMaxString(input.Password, constanta.ClientMappingPassword, 1, 100)
		if err.Error != nil {
			return err
		}
	}

	// Optional Validation

	if !util.IsStringEmpty(input.UserID) {
		err = util2.ValidateMinMaxString(input.UserID, constanta.UserID, 3, 20)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.Lastname) {
		err = util2.ValidateMinMaxString(input.Lastname, constanta.LastName, 1, 50)
		if err.Error != nil {
			return err
		}

		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.LastName, input.Lastname)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.AndroidID) {
		err = util2.ValidateMinMaxString(input.AndroidID, constanta.ClientMappingAndroidID, 1, 100)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.SalesmanID) {
		err = util2.ValidateMinMaxString(input.SalesmanID, constanta.SalesmanID, 1, 30)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.Email) {
		if !IsEmailAddressMDB(input.Email) {
			return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
		}
	}

	if !util.IsStringEmpty(input.Email) {
		input.Email = strings.ToLower(input.Email)
	}

	if !util.IsStringEmpty(input.NoTelp) || !util.IsStringEmpty(input.CountryCode) {
		if util.IsStringEmpty(input.CountryCode) {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CountryCode)
		}

		isValid := IsPhoneNumberWithCountryCodeMDB(input.CountryCode + "-" + input.NoTelp)
		if !isValid {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.PhoneRegex, constanta.Phone, "")
		}
	}

	if !util.IsStringEmpty(input.UniqueID2) {
		err = util2.ValidateMinMaxString(input.UniqueID2, constanta.UniqueID2, 1, 20)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.SalesmanCategory) {
		err = util2.ValidateMinMaxString(input.SalesmanCategory, constanta.SalesmanCategory, 1, 10)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.ClientAliases) {
		err = util2.ValidateMinMaxString(input.ClientAliases, constanta.ClientMappingClientAlias, 1, 50)
		if err.Error != nil {
			return
		}

		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.ClientMappingClientAlias, input.ClientAliases)
		if err.Error != nil {
			return err
		}
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

	return
}

func (input *RegisterNamedUserRequest) ValidateRegisterNamedUserRequest() (err errorModel.ErrorModel) {
	var (
		fileName         = input.fileNameGenerate()
		funcName         = "ValidateRegisterNamedUser"
		validationResult bool
		errField         string
	)

	if util.IsStringEmpty(input.ParentClientID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientMappingParentClientID)
	}

	if util.IsStringEmpty(input.ClientID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientMappingClientID)
	}

	validationResult, errField, _ = util2.IsClientIDValid(input.ClientID)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.ClientMappingClientID, "")
	}

	if input.ClientTypeID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.NewClientType)
	}

	if input.AuthUserID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ClientMappingAuthUserID)
	}

	if !util.IsStringEmpty(input.UserID) {
		err = util2.ValidateMinMaxString(input.UserID, constanta.ClientMappingUserID, 1, 30)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.Firstname) {
		err = util2.ValidateMinMaxString(input.Firstname, constanta.FirstName, 1, 50)
		if err.Error != nil {
			return err
		}

		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.FirstName, input.Firstname)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.Lastname) {
		err = util2.ValidateMinMaxString(input.Lastname, constanta.LastName, 1, 50)
		if err.Error != nil {
			return err
		}

		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.LastName, input.Lastname)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.ClientAliases) {
		err = util2.ValidateMinMaxString(input.ClientAliases, constanta.ClientMappingClientAlias, 1, 50)
		if err.Error != nil {
			return
		}

		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.ClientMappingClientAlias, input.ClientAliases)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.Password) {
		err = util2.ValidateMinMaxString(input.Password, constanta.ClientMappingPassword, 1, 100)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.SalesmanID) {
		err = util2.ValidateMinMaxString(input.SalesmanID, constanta.SalesmanID, 1, 30)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.AndroidID) {
		err = util2.ValidateMinMaxString(input.AndroidID, constanta.ClientMappingAndroidID, 1, 100)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.RegDateStr) {
		input.RegDate, err = TimeStrToTimeWithTimeFormat(input.RegDateStr, constanta.ClientMappingRegDate, constanta.DefaultInstallationTimeFormat)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.Email) {
		if !IsEmailAddressMDB(input.Email) {
			return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
		}
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

	if !util.IsStringEmpty(input.NoTelp) {
		isValid := IsPhoneNumberWithCountryCodeMDB(input.NoTelp)
		if !isValid {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.PhoneRegex, constanta.Phone, "")
		}
	}

	if !util.IsStringEmpty(input.SalesmanCategory) {
		err = util2.ValidateMinMaxString(input.SalesmanCategory, constanta.SalesmanCategory, 1, 10)
		if err.Error != nil {
			return err
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *RegisterNamedUserRequest) ValidateRegisterNamedUserClientMappingRequest() (err errorModel.ErrorModel) {
	fileName := input.fileNameGenerate()
	funcName := "ValidateRegisterNamedUserClientMappingRequest"
	var validationResult bool
	var errField string

	if util.IsStringEmpty(input.ParentClientID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientMappingParentClientID)
	}

	if util.IsStringEmpty(input.ClientID) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientMappingClientID)
	}

	validationResult, errField, _ = util2.IsClientIDValid(input.ClientID)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.ClientMappingClientID, "")
	}

	if input.ClientTypeID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.NewClientType)
	}

	if input.AuthUserID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ClientMappingAuthUserID)
	}

	if !util.IsStringEmpty(input.UserID) {
		err = util2.ValidateMinMaxString(input.UserID, constanta.ClientMappingUserID, 1, 30)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.ClientAliases) {
		err = util2.ValidateMinMaxString(input.ClientAliases, constanta.ClientMappingClientAlias, 1, 50)
		if err.Error != nil {
			return
		}
	}

	if !util.IsStringEmpty(input.Password) {
		err = util2.ValidateMinMaxString(input.Password, constanta.ClientMappingPassword, 1, 100)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.SalesmanID) {
		err = util2.ValidateMinMaxString(input.SalesmanID, constanta.SalesmanID, 1, 30)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.AndroidID) {
		err = util2.ValidateMinMaxString(input.AndroidID, constanta.ClientMappingAndroidID, 1, 100)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.Email) {
		if !util.IsEmailAddress(input.Email) {
			return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
		}
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

	if !util.IsStringEmpty(input.NoTelp) {
		var isValid bool

		isValid = strings.Contains(input.NoTelp, "-")
		if !isValid {
			return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Phone)
		}

		phoneNew := strings.Split(input.NoTelp, "-")
		if len(phoneNew) != 2 {
			return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Phone)
		}

		isValid = util.IsCountryCode(phoneNew[0])
		if !isValid {
			return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Phone)
		}

		_, isValid = util.IsPhoneNumber(phoneNew[1])
		if !isValid {
			return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Phone)
		}
	}

	if !util.IsStringEmpty(input.SalesmanCategory) {
		err = util2.ValidateMinMaxString(input.SalesmanCategory, constanta.SalesmanCategory, 1, 10)
		if err.Error != nil {
			return err
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input RegisterNamedUserRequest) fileNameGenerate() (fileName string) {
	return "RegisterNamedUserDTO.go"
}
