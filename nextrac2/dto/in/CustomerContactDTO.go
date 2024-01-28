package in

import (
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"regexp"
	"strings"
	"time"
)

type CustomerContactRequest struct {
	AbstractDTO
	ID                  int64  `json:"id"`
	CustomerID          int64  `json:"customer_id"`
	MdbCompanyProfileID int64  `json:"-"`
	MdbPersonProfileID  int64  `json:"mdb_person_profile_id"`
	Nik                 string `json:"nik"`
	MdbPersonTitleID    int64  `json:"mdb_person_title_id"`
	PersonTitle         string `json:"person_title"`
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	Sex                 string `json:"sex"`
	Address             string `json:"address"`
	Address2            string `json:"address_2"`
	Address3            string `json:"address_3"`
	Hamlet              string `json:"hamlet"`
	Neighbourhood       string `json:"neighbourhood"`
	ProvinceID          int64  `json:"province_id"`
	DistrictID          int64  `json:"district_id"`
	PhoneCountryCode    string `json:"phone_country_code"`
	Phone               string `json:"phone"`
	Email               string `json:"email"`
	MdbPositionID       int64  `json:"mdb_position_id"`
	PositionName        string `json:"position_name"`
	Status              string `json:"status"`
	MdbContactPersonID  int64  `json:"mdb_contact_person_id"`
	Action              int64  `json:"action"`
	UpdatedAtStr        string `json:"updated_at"`
	UpdatedAt           time.Time
	IsSuccess           bool
	MDBProvinceID       int64
	MDBDistrictID       int64
}

func (input *CustomerContactRequest) ValidateUpdateForUpdateCustomer() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdateForUpdateCustomer"

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(CustomerContactDTOFileName, funcName, constanta.ID)
	}

	err = input.mandatoryFieldValidation(CustomerContactDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	return
}

func (input *CustomerContactRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdate"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldError(CustomerContactDTOFileName, funcName, constanta.ID)
	}

	if input.MdbPersonProfileID < 1 {
		return errorModel.GenerateEmptyFieldError(CustomerContactDTOFileName, funcName, constanta.MDBPersonProfileID)
	}

	err = input.mandatoryFieldValidation(CustomerContactDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	return
}

func (input *CustomerContactRequest) ValidateDelete() (err errorModel.ErrorModel) {
	funcName := "ValidateDelete"

	err = input.validateForUpdateAndDelete(CustomerContactDTOFileName, funcName)
	return
}

func (input *CustomerContactRequest) ValidateInsert(isCheckCustomerID bool) (err errorModel.ErrorModel) {
	funcName := "ValidateInsert"

	// Validate CustomerID
	if isCheckCustomerID {
		if input.CustomerID < 1 {
			err = errorModel.GenerateEmptyFieldError(CustomerContactDTOFileName, funcName, constanta.CustomerID)
			return
		}
	}

	return input.mandatoryFieldValidation(CustomerContactDTOFileName, funcName)
}

func (input *CustomerContactRequest) validateForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
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

func (input *CustomerContactRequest) mandatoryFieldValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	var (
		validationResult bool
		errField         string
		additionalInfo   string
	)

	//--- NIK Min 16 and Max 20
	if !util.IsStringEmpty(input.Nik) {
		err = util2.ValidateMinMaxString(input.Nik, constanta.NIK, 16, 20)
		if err.Error != nil {
			return
		}

		//--- NPWP Must Digit
		nikRegex := regexp.MustCompile(fmt.Sprintf(`^[0-9]{16,20}$`))
		if !nikRegex.MatchString(input.Nik) {
			err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "MUST_DIGIT_NIK_RULE", constanta.NIK, "")
			return
		}

		//--- NPWP Birth Date
		birthNPWPRegex := regexp.MustCompile(fmt.Sprintf(`(0[1-9]|[12]\d|[4-6]\d|[37][01])`))
		monthNPWPRegex := regexp.MustCompile(fmt.Sprintf(`(0[1-9]|1[0-2])`))
		if !birthNPWPRegex.MatchString(input.Nik[6:8]) {
			err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "DATE_BIRTH_NPWP_RULE", constanta.NIK, "")
			return
		}

		if !monthNPWPRegex.MatchString(input.Nik[8:10]) {
			err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "MONTH_BIRTH_NPWP_RULE", constanta.NIK, "")
			return
		}
	}

	// Validate MDB Person Title ID
	//if input.MdbPersonTitleID < 1 {
	//	err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.MDBPersonTitleID)
	//	return
	//}

	if input.MdbPositionID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.MDBPositionID)
		return
	}

	err = util2.ValidateMinMaxString(input.FirstName, constanta.FirstName, 3, 20)
	if err.Error != nil {
		return
	}

	validationResult, errField, additionalInfo = util.IsNexsoftNameStandardValid(input.FirstName)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.FirstName, additionalInfo)
	}

	if !util.IsStringEmpty(input.LastName) {
		err = util2.ValidateMinMaxString(input.LastName, constanta.LastName, 3, 50)
		if err.Error != nil {
			return
		}

		validationResult, errField, additionalInfo = util.IsNexsoftNameStandardValid(input.LastName)
		if !validationResult {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.LastName, additionalInfo)
		}
	}

	//--- Validate Sex
	if !util.IsStringEmpty(input.Sex) {
		if input.Sex != "N" && input.Sex != "L" && input.Sex != "P" {
			err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Sex)
			return
		}
	}

	err = util2.ValidateMinMaxString(input.Address, constanta.Address, 1, 100)
	if err.Error != nil {
		return
	}

	if !util.IsStringEmpty(input.Hamlet) {
		err = input.ValidateMinMaxString(input.Hamlet, constanta.Hamlet, 1, 5)
		if err.Error != nil {
			return
		}

		//--- Auto Add Number RT
		AutoAddNumberRTRW(&input.Hamlet, 3)
		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.Hamlet, input.Hamlet)
		if err.Error != nil {
			return
		}
	}

	if !util.IsStringEmpty(input.Neighbourhood) {
		err = input.ValidateMinMaxString(input.Neighbourhood, constanta.Neighbourhood, 1, 5)
		if err.Error != nil {
			return
		}

		//--- Auto Add Number RW
		AutoAddNumberRTRW(&input.Neighbourhood, 3)
		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.Neighbourhood, input.Neighbourhood)
		if err.Error != nil {
			return
		}
	}

	//--- Validate Phone
	err = input.ValidateMinMaxString(input.Phone, constanta.Phone, 1, 18)
	if err.Error != nil {
		return
	}

	if !IsPhoneNumberWithCountryCodeMDB(input.Phone) {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.PhoneRegex, constanta.Phone, "")
		return
	}

	//--- Email Validation
	err = util2.ValidateMinMaxString(input.Email, constanta.Email, 1, 100)
	if err.Error != nil {
		return
	}

	if !util.IsEmailAddress(input.Email) {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
		return
	}

	if input.Email != "" {
		input.Email = strings.ToLower(input.Email)
	}

	//--- Status Validation
	if util.IsStringEmpty(input.Status) {
		input.Status = "N"
	}

	if input.Status != "N" && input.Status != "A" && input.Status != "P" {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Status)
		return
	}

	return
}
