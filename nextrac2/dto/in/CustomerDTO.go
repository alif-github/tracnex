package in

import (
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	common2 "nexsoft.co.id/nextrac2/resource_common_service/common"
	util2 "nexsoft.co.id/nextrac2/util"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type CustomerRequest struct {
	AbstractDTO
	ID                          int64                    `json:"id"`
	IsPrincipal                 bool                     `json:"is_principal"`
	IsParent                    bool                     `json:"is_parent"`
	ParentCustomerID            int64                    `json:"parent_customer_id"`
	MDBParentCustomerID         int64                    `json:"mdb_parent_customer_id"`
	MDBCompanyProfileID         int64                    `json:"mdb_company_profile_id"`
	Npwp                        string                   `json:"npwp"`
	MDBCompanyTitleID           int64                    `json:"mdb_company_title_id"`
	CompanyTitle                string                   `json:"company_title"`
	CustomerName                string                   `json:"customer_name"`
	Address                     string                   `json:"address"`
	Address2                    string                   `json:"address_2"`
	Address3                    string                   `json:"address_3"`
	Hamlet                      string                   `json:"hamlet"`
	Neighbourhood               string                   `json:"neighbourhood"`
	CountryID                   int64                    `json:"country_id"`
	ProvinceID                  int64                    `json:"province_id"`
	DistrictID                  int64                    `json:"district_id"`
	SubDistrictID               int64                    `json:"sub_district_id"`
	UrbanVillageID              int64                    `json:"urban_village_id"`
	PostalCodeID                int64                    `json:"postal_code_id"`
	Longitude                   float64                  `json:"long"`
	Latitude                    float64                  `json:"lat"`
	PhoneCountryCode            string                   `json:"phone_country_code"`
	Phone                       string                   `json:"phone"`
	AlternativePhoneCountryCode string                   `json:"alternative_phone_country_code"`
	AlternativePhone            string                   `json:"alternative_phone"`
	Fax                         string                   `json:"fax"`
	CompanyEmail                string                   `json:"company_email"`
	AlternativeCompanyEmail     string                   `json:"alternative_company_email"`
	CustomerSource              string                   `json:"customer_source"`
	TaxName                     string                   `json:"tax_name"`
	TaxAddress                  string                   `json:"tax_address"`
	SalesmanID                  int64                    `json:"salesman_id"`
	RefCustomerID               int64                    `json:"ref_customer_id"`
	DistributorOF               string                   `json:"distributor_of"`
	CustomerGroupID             int64                    `json:"customer_group_id"`
	CustomerCategoryID          int64                    `json:"customer_category_id"`
	Status                      string                   `json:"status"`
	CustomerContact             []CustomerContactRequest `json:"customer_contact"`
	UpdatedAtStr                string                   `json:"updated_at"`
	UpdatedAt                   time.Time
	IsSuccess                   bool
	MDBProvinceID               int64
	MDBDistrictID               int64
	MDBSubDistrictID            int64
	MDBUrbanVillageID           int64
	MDBPostalCodeID             int64
}

func (input *CustomerRequest) ValidateDelete() (err errorModel.ErrorModel) {
	funcName := "ValidateDelete"

	return input.validationForUpdateAndDelete(CustomerDTOFileName, funcName)
}

func (input *CustomerRequest) ValidateUpdate(isUpdateAll bool) (err errorModel.ErrorModel) {
	funcName := "ValidateUpdate"

	err = input.mandatoryField(isUpdateAll, CustomerDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	err = input.validationForUpdateAndDelete(CustomerDTOFileName, funcName)
	if err.Error != nil {
		return
	}

	//if input.MDBCompanyProfileID < 1 {
	//	err = errorModel.GenerateEmptyFieldError(CustomerDTOFileName, funcName, constanta.MDBCompanyProfileID)
	//	return
	//}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input CustomerRequest) ValidateView() (err errorModel.ErrorModel) {
	if input.ID < 1 {
		err = errorModel.GenerateEmptyFieldError(CustomerDTOFileName, "ValidateView", constanta.ID)
		return
	}
	return
}

func (input *CustomerRequest) ValidateInsert() (err errorModel.ErrorModel) {
	funcName := "ValidateInsert"
	return input.mandatoryField(true, CustomerDTOFileName, funcName)
}

func (input *CustomerRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
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

func (input *CustomerRequest) mandatoryField(isUpdateAll bool, fileName string, funcName string) (err errorModel.ErrorModel) {
	var validationResult bool

	if isUpdateAll {
		//--- NPWP Is Empty
		if util.IsStringEmpty(input.Npwp) {
			err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.NPWP)
			return
		}

		//--- NPWP Validation
		if err = input.ValidateMinMaxString(input.Npwp, constanta.NPWP, 1, 38); err.Error != nil {
			return
		}

		//--- NPWP Must Digit
		//if !util.IsNPWPValid(input.Npwp) {
		//	err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "MUST_DIGIT_NPWP_RULE", constanta.NPWP, "")
		//	return
		//}

		//--- NPWP Must Digit
		npwpRegex := regexp.MustCompile(fmt.Sprintf(`^[0-9]{15,16}$`))
		if !npwpRegex.MatchString(input.Npwp) {
			err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "MUST_DIGIT_NPWP_RULE", constanta.NPWP, "")
			return
		}

		//--- NPWP Birth Date
		d, _ := strconv.Atoi(string(input.Npwp[0]))
		if d != 0 && len(input.Npwp) == 16 {
			birthNPWPRegex := regexp.MustCompile(fmt.Sprintf(`(0[1-9]|[12]\d|[4-6]\d|[37][01])`))
			monthNPWPRegex := regexp.MustCompile(fmt.Sprintf(`(0[1-9]|1[0-2])`))
			if !birthNPWPRegex.MatchString(input.Npwp[6:8]) {
				err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "DATE_BIRTH_NPWP_RULE", constanta.NPWP, "")
				return
			}

			if !monthNPWPRegex.MatchString(input.Npwp[8:10]) {
				err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "MONTH_BIRTH_NPWP_RULE", constanta.NPWP, "")
				return
			}
		}
	}

	// Status Company Validation
	if input.IsPrincipal && !input.IsParent {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "CUSTOMER_PRINCIPAL_VALIDATION", constanta.Customer, "")
		return
	}

	if !input.IsPrincipal && !input.IsParent && input.ParentCustomerID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ParentCustomerID)
		return
	}

	//if input.IsParent && input.ParentCustomerID > 0 {
	//	err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.ParentCustomerID)
	//	return
	//}

	// Status Validation
	if util.IsStringEmpty(input.Status) {
		input.Status = "N"
	}

	if input.Status != "N" && input.Status != "A" && input.Status != "P" {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Status)
		return
	}

	// MDB Company Title
	if input.MDBCompanyTitleID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CompanyTitle)
		return
	}

	//Customer Name Validation
	if err = input.ValidateMinMaxString(input.CustomerName, constanta.CustomerName, 1, 50); err.Error != nil {
		return
	}

	//if err = input.ValidateMinMaxString(input.CustomerName, constanta.CustomerName, 1, 100); err.Error != nil {
	//	return
	//}

	if err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.CustomerName, input.CustomerName); err.Error != nil {
		return
	}

	// Address Validation
	if err = input.ValidateMinMaxString(input.Address, constanta.Address, 1, 100); err.Error != nil {
		return
	}

	//if err = input.ValidateMinMaxString(input.Address, constanta.Address, 1, 150); err.Error != nil {
	//	return
	//}

	// Province Validation
	if input.ProvinceID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ProvinceID)
		return
	}

	// District Validation
	if input.DistrictID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.DistrictID)
		return
	}

	// Phone Validation
	if err = input.validateContactNumber(fileName, funcName); err.Error != nil {
		return
	}

	// Email Validation
	if err = input.ValidateMinMaxString(input.CompanyEmail, constanta.Email, 1, 100); err.Error != nil {
		return
	}

	if !util.IsEmailAddress(input.CompanyEmail) {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
		return
	}

	if input.CompanyEmail != "" {
		input.CompanyEmail = strings.ToLower(input.CompanyEmail)
	}

	validationResult = !util.IsStringEmpty(input.AlternativeCompanyEmail)
	if err = input.ValidateMinMaxString(input.AlternativeCompanyEmail, constanta.AlternativeEmail, 1, 100); err.Error != nil && validationResult {
		return
	}

	if !util.IsEmailAddress(input.AlternativeCompanyEmail) && !util.IsStringEmpty(input.AlternativeCompanyEmail) {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.AlternativeEmail)
		return
	}

	if input.AlternativeCompanyEmail != "" {
		input.AlternativeCompanyEmail = strings.ToLower(input.AlternativeCompanyEmail)
	}

	if err = input.validateOptionalField(fileName, funcName); err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *CustomerRequest) validateOptionalField(fileName, funcName string) (err errorModel.ErrorModel) {
	var validationResult bool
	var customerSourceArr = []string{"Nexsoft", "Reseller", "Refree"}

	validationResult = !util.IsStringEmpty(input.DistributorOF)
	if err = input.ValidateMinMaxString(input.DistributorOF, constanta.DistributorOf, 1, 255); err.Error != nil && validationResult {
		return
	}

	if err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.DistributorOf, input.DistributorOF); err.Error != nil && validationResult {
		return
	}

	// Validate Customer Source
	if !util.IsStringEmpty(input.CustomerSource) && !common2.ValidateStringContainInStringArray(customerSourceArr, input.CustomerSource) {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.CustomerSource)
		return
	}

	validationResult = !util.IsStringEmpty(input.Hamlet)
	if validationResult {
		if err = input.ValidateMinMaxString(input.Hamlet, constanta.Hamlet, 1, 5); err.Error != nil {
			return
		}

		AutoAddNumberRTRW(&input.Hamlet, 3)
		if err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.Hamlet, input.Hamlet); err.Error != nil {
			return
		}
	}

	validationResult = !util.IsStringEmpty(input.Neighbourhood)
	if validationResult {
		if err = input.ValidateMinMaxString(input.Neighbourhood, constanta.Neighbourhood, 1, 5); err.Error != nil {
			return
		}

		AutoAddNumberRTRW(&input.Neighbourhood, 3)
		if err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.Neighbourhood, input.Neighbourhood); err.Error != nil {
			return
		}
	}

	validationResult = !util.IsStringEmpty(input.TaxName)
	if err = input.ValidateMinMaxString(input.TaxName, constanta.TaxName, 1, 100); err.Error != nil && validationResult {
		return
	}

	validationResult = !util.IsStringEmpty(input.TaxAddress)
	if err = input.ValidateMinMaxString(input.TaxAddress, constanta.TaxAddress, 1, 255); err.Error != nil && validationResult {
		return
	}

	if !util.IsStringEmpty(input.Fax) {
		isValid := IsPhoneNumberWithCountryCodeMDB(input.Fax)
		if !isValid {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.PhoneRegex, constanta.Fax, "")
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *CustomerRequest) validateContactNumber(fileName, funcName string) (err errorModel.ErrorModel) {
	var validationResult bool

	if err = input.ValidateMinMaxString(input.Phone, constanta.Phone, 1, 18); err.Error != nil {
		return
	}

	if input.Phone, err = input.validateCountryCodeAndPhone(fileName, funcName, input.Phone, constanta.Phone, false); err.Error != nil {
		return
	}

	validationResult = util.IsStringEmpty(input.AlternativePhone)
	if err = input.ValidateMinMaxString(input.AlternativePhone, constanta.AlternativePhone, 1, 18); err.Error != nil && !validationResult {
		return
	}

	if input.AlternativePhone, err = input.validateCountryCodeAndPhone(fileName, funcName, input.AlternativePhone, constanta.AlternativePhone, false); err.Error != nil && !validationResult {
		return
	}

	validationResult = !util.IsStringEmpty(input.Fax)
	if err = input.ValidateMinMaxString(input.Fax, constanta.Fax, 1, 18); err.Error != nil && validationResult {
		return
	}

	if input.Fax, err = input.validateCountryCodeAndPhone(fileName, funcName, input.Fax, constanta.Fax, true); err.Error != nil && validationResult {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *CustomerRequest) validateCountryCodeAndPhone(fileName, funcName, phoneStr, fieldName string, isFax bool) (result string, err errorModel.ErrorModel) {
	var splitPhone []string
	var phone int
	var validationResult bool

	if !IsPhoneNumberWithCountryCodeMDB(input.Phone) {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.PhoneRegex, constanta.Phone, "")
		return
	}

	splitPhone = strings.Split(phoneStr, "-")
	if !util.IsCountryCode(splitPhone[0]) {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
		return
	}

	if !util.IsFacsimileValid(splitPhone[1]) && isFax {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
		return
	}

	if phone, validationResult = util.IsPhoneNumber(splitPhone[1]); !validationResult && !isFax {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
		return
	}

	result = splitPhone[0] + "-" + strconv.Itoa(phone)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input CustomerRequest) validateCustomerContact(inputStruct *CustomerContactRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsert(false)
}
