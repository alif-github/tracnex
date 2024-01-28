package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"strconv"
	"strings"
	"time"
)

type UserRequest struct {
	AbstractDTO
	ID                    int64  `json:"id"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	ConfirmPassword       string `json:"confirm_password"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	Email                 string `json:"email"`
	ClientID              string `json:"client_id"`
	Role                  string `json:"role"`
	DataGroupID           string `json:"data_group_id"`
	Group                 string `json:"group"`
	CountryCode           string `json:"country_code"`
	Phone                 string `json:"phone"`
	Device                string `json:"device"`
	Locale                string `json:"locale"`
	IPWhitelist           string `json:"ip_whitelist"`
	IsAdmin               bool   `json:"is_admin"`
	Status                string `json:"status"`
	AuthUserID            int64  `json:"auth_user_id"`
	ResourceID            string `json:"resource_id"`
	PlatformDevice        string `json:"platform_device"`
	Currency              string `json:"currency"`
	Scope                 string
	ClientSecret          string
	SignatureKey          string
	UpdatedAt             time.Time
	UpdatedAtStr          string                        `json:"updated_at"`
	AdditionalInformation []model.AdditionalInformation `json:"additional_information"`
}

func (input UserRequest) AdditionalInformationString() string {
	if input.AdditionalInformation != nil {
		return util.StructToJSON(input.AdditionalInformation)
	} else {
		return ""
	}
}

func (input UserRequest) ValidateAddResourceToUser() (err errorModel.ErrorModel) {
	fileName := "UserDTOIn.go"
	funcName := "ValidateAddResourceToUser"

	if util.IsStringEmpty(input.Username) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Username)
	}

	if !util.IsEmailAddress(input.Email) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "Email")
	}

	err = input.ValidateMinMaxString(input.Email, constanta.Email, 6, 256)
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}

func (input UserRequest) ValidationCheckUserAuth() (err errorModel.ErrorModel) {
	var (
		fileName         = "UserDTOIn.go"
		funcName         = "ValidationCheckUserAuth"
		validationResult bool
	)

	if util.IsStringEmpty(input.Email) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Email)
	}

	if !util.IsEmailAddress(input.Email) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
	}

	err = input.ValidateMinMaxString(input.Email, constanta.Email, 6, 256)
	if err.Error != nil {
		return err
	}

	if util.IsStringEmpty(input.Phone) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Phone)
	}

	validationResult = util.IsPhoneNumberWithCountryCode(input.Phone)
	if !validationResult {
		return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Phone)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input UserRequest) ValidationCheckUsernameAuth() (err errorModel.ErrorModel) {
	fileName := "UserDTOIn.go"
	funcName := "ValidationCheckUsernameAuth"

	if util.IsStringEmpty(input.Username) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Username)
	}

	validationResult, errField, additionalInfo := util.IsNexsoftUsernameStandardValid(input.Username)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.Username, additionalInfo)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *UserRequest) ValidateInternalInsertUser() errorModel.ErrorModel {
	err := input.validateMandatoryUser()
	if err.Error != nil {
		return err
	}

	if input.Role == "" {
		input.Role = constanta.DefaultRoleUser
	}

	if input.Group == "" {
		input.Group = constanta.TanggerangDataGoup
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *UserRequest) ValidateViewUserAndResendOTP() errorModel.ErrorModel {
	fileName := "UserDTOIn.go"
	funcName := "ValidateViewUserAndResendOTP"
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.User)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *UserRequest) ValidateInsertUser() errorModel.ErrorModel {
	var (
		fileName = "UserDTOIn.go"
		funcName = "ValidateInsertUser"
		err      errorModel.ErrorModel
	)

	err = input.validateMandatoryUser()
	if err.Error != nil {
		return err
	}

	if util.IsStringEmpty(input.Role) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "Role")
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *UserRequest) ValidateDeleteUser() (err errorModel.ErrorModel) {
	fileName := "UserDTO.go"
	funcName := "ValidateDeleteUser"

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

func (input *UserRequest) ValidateUpdateUser() (err errorModel.ErrorModel) {
	fileName := "UserDTO.go"
	funcName := "ValidateUpdateUser"

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

	if (input.Status != constanta.StatusActive) && (input.Status != constanta.StatusNonActive) && (input.Status != "") {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.StatusUserRegex, constanta.Status, "")
		return
	}

	err = input.validateMandatoryUpdateUser()
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *UserRequest) ValidateUpdateUserProfile() (err errorModel.ErrorModel) {
	fileName := "UserDTO.go"
	funcName := "ValidateUpdateUserProfile"

	if util.IsStringEmpty(input.UpdatedAtStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	err = input.validateMandatoryUpdateUser()
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *UserRequest) validateMandatoryUpdateUser() errorModel.ErrorModel {
	var (
		fileName         = "UserDTO.go"
		funcName         = "validateMandatoryUser"
		validationResult bool
		errField         string
		additionalInfo   string
		err              errorModel.ErrorModel
	)

	validationResult, errField, additionalInfo = util.IsNexsoftNameStandardValid(input.FirstName)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, "FirstName", additionalInfo)
	}

	if input.LastName != "" {
		err = input.ValidateMinMaxString(input.LastName, "LastName", 1, 50)
		if err.Error != nil {
			return err
		}

		validationResult, errField, additionalInfo = util.IsNexsoftNameStandardValid(input.LastName)
		if !validationResult {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, "LastName", additionalInfo)
		}
	}

	if !util.IsEmailAddress(input.Email) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "Email")
	}

	err = input.ValidateMinMaxString(input.Email, constanta.Email, 6, 256)
	if err.Error != nil {
		return err
	}

	if !util.IsStringEmpty(input.Email) {
		input.Email = strings.ToLower(input.Email)
	}

	if !util.IsCountryCode(input.CountryCode) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "Country Code")
	}

	validationResult = util.IsPhoneNumberWithCountryCode(input.CountryCode + "-" + input.Phone)
	if !validationResult {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "Phone")
	}

	//tutup sementara validasi platform device
	//if !util.IsStringEmpty(input.PlatformDevice) {
	//	return errorModel.GenerateFormatFieldError(fileName, funcName, "Platform Device")
	//}

	return errorModel.GenerateNonErrorModel()
}

func (input *UserRequest) validateMandatoryUser() errorModel.ErrorModel {
	var (
		fileName         = "UserDTOIn.go"
		funcName         = "validateMandatoryUser"
		validationResult bool
		errField         string
		additionalInfo   string
		err              errorModel.ErrorModel
	)

	if util.IsStringEmpty(input.FirstName) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "Firstname")
	}

	err = input.ValidateMinMaxString(input.FirstName, "Firstname", 2, 50)
	if err.Error != nil {
		return err
	}

	validationResult, errField, additionalInfo = util.IsNexsoftNameStandardValid(input.FirstName)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, "Firstname", additionalInfo)
	}

	if input.LastName != "" {
		err = input.ValidateMinMaxString(input.LastName, "Lastname", 2, 50)
		if err.Error != nil {
			return err
		}

		validationResult, errField, additionalInfo = util.IsNexsoftNameStandardValid(input.LastName)
		if !validationResult {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, "Lastname", additionalInfo)
		}
	}

	if util.IsStringEmpty(input.Username) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "Username")
	}

	validationResult, errField, additionalInfo = util.IsNexsoftUsernameStandardValid(input.Username)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, "Username", additionalInfo)
	}

	if input.AuthUserID < 1 {
		if util.IsStringEmpty(input.Password) {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, "Password")
		}

		if util.IsStringEmpty(input.ConfirmPassword) {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, "Confirm Password")
		}

		if input.Password != input.ConfirmPassword {
			return errorModel.GenerateInvalidDifferentCompareData(fileName, funcName, "Password", "Confirm Password")
		}

		validationResult, errField, additionalInfo = util.IsNexsoftPasswordStandardValid(input.Password)
		if !validationResult {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, "Password", additionalInfo)
		}
	}

	if util.IsStringEmpty(input.Email) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "Email")
	}

	if !util.IsEmailAddress(input.Email) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "Email")
	}

	err = input.ValidateMinMaxString(input.Email, constanta.Email, 6, 256)
	if err.Error != nil {
		return err
	}

	if util.IsStringEmpty(input.Email) {
		input.Email = strings.ToLower(input.Email)
	}

	if util.IsStringEmpty(input.CountryCode) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "Country Code")
	}

	if !util.IsCountryCode(input.CountryCode) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "Country Code")
	}

	if util.IsStringEmpty(input.Phone) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "Phone")
	}

	_, validationResult, errField = input.IsPhoneNumberValid(input.Phone)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, "Phone", "")
	}

	validationResult = util.IsPhoneNumberWithCountryCode(input.CountryCode + "-" + input.Phone)
	if !validationResult {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "Phone")
	}

	if input.AdditionalInformation != nil {
		for i := 0; i < len(input.AdditionalInformation); i++ {
			validationResult, errField = util.IsNexsoftAdditionalInformationKeyStandardValid(input.AdditionalInformation[i].Key)
			if !validationResult {
				return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, "Additional Information ["+strconv.Itoa(i+1)+"]", "")
			}

			for j := 1 + i; j < len(input.AdditionalInformation); j++ {
				if input.AdditionalInformation[j].Key == input.AdditionalInformation[i].Key {
					if input.AdditionalInformation[j].Value == input.AdditionalInformation[i].Value {
						input.AdditionalInformation = append(input.AdditionalInformation[:i], input.AdditionalInformation[i+1:]...)
						j--
					} else {
						return errorModel.GenerateFormatFieldError(fileName, funcName, "Additional Information")
					}
				}
			}
		}
	}

	input.Username = strings.Trim(strings.ToLower(input.Username), " ")

	if input.Locale == "" {
		input.Locale = constanta.DefaultApplicationsLanguage
	}

	input.Email = strings.ToLower(input.Email)

	//tutup sementara validasi platform device
	//if util.IsStringEmpty(input.PlatformDevice) {
	//	return errorModel.GenerateEmptyFieldError(fileName, funcName, "Platform Device")
	//}

	return errorModel.GenerateNonErrorModel()
}
