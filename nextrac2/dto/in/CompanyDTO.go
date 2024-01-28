package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"regexp"
	"time"
)

type CompanyRequest struct {
	AbstractDTO
	ID                   int64     `json:"id"`
	CompanyTitle         string    `json:"company_title"`
	CompanyName          string    `json:"company_name"`
	Address              string    `json:"address"`
	Address2             string    `json:"address2"`
	Neighbourhood        string    `json:"neighbourhood"`
	Hamlet               string    `json:"hamlet"`
	ProvinceId           int64     `json:"province_id"`
	ProvinceName         string    `json:"province_name"`
	DistrictId           int64     `json:"district_id"`
	DistrictName         string    `json:"district_name"`
	SubDistrictId        int64     `json:"sub_district_id"`
	SubDistrictName      string    `json:"sub_district_name"`
	VillageId            int64     `json:"urban_village_id"`
	Village              string    `json:"village"`
	PostalCodeId         int64     `json:"postal_code_id"`
	PostalCode           string    `json:"postal_code"`
	Longitude            string    `json:"longitude"`
	Latitude             string    `json:"latitude"`
	Telephone            string    `json:"telephone"`
	TelephoneAlternate   string    `json:"alternate_telephone"`
	Fax                  string   `json:"fax"`
	Email                string     `json:"email"`
	AlternateEmail       string     `json:"alternate_email"`
	Npwp                 string     `json:"npwp"`
	TaxName              string      `json:"tax_name"`
	TaxAddress           string     `json:"tax_address"`
	UpdatedAtStr         string    `json:"updated_at"`
	UpdatedAt            time.Time
}

func (input *CompanyRequest) ValidateCompany(isUpdate bool) (err errorModel.ErrorModel) {
	funcName := "ValidateCompany"
	fileName := "CompanyDTO.go"

	if input.CompanyTitle == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "company_title")
	}

	err = input.ValidateMinMaxString(input.CompanyTitle, "company_title", 2, 10)
	if err.Error != nil {
		return
	}

	if input.CompanyName == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "company_name")
	}

	if len(input.CompanyName) > 100{
		return errorModel.GenerateFieldHaveMaxLimitError(fileName, funcName, "company_name", 100)
	}

	if input.ProvinceId == 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "province_id")
	}

	if input.DistrictId == 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "district_id")
	}

	if input.SubDistrictId == 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "sub_district_id")
	}

	if input.VillageId == 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "urban_village_id")
	}

	if input.Address == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "address")
	}

	if len(input.Address) > 256{
		return errorModel.GenerateFieldHaveMaxLimitError(fileName, funcName, "address", 256)
	}

	if len(input.Neighbourhood) > 5 && input.Neighbourhood != ""{
		return errorModel.GenerateFieldHaveMaxLimitError(fileName, funcName, "neighbourhood", 5)
	}

	if len(input.Hamlet) > 5 && input.Hamlet != ""{
		return errorModel.GenerateFieldHaveMaxLimitError(fileName, funcName, "hamlet", 5)
	}

	if input.Telephone == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "telephone")
	}

	if input.Telephone != ""{
		if !IsPhoneValid(input.Telephone){
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "telephone harus diawali dengan +62-", "telephone", "")
		}
	}

	if input.Npwp == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "npwp")
	}

	if len(input.Npwp) != 15 && len(input.Npwp) != 16{
        return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "npwp hanya bisa diisi dengan 15-16 karakter", "npwp", "")
	}

	if input.Email != "" && !IsEmailAddress(input.Email){
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "format email tidak sesuai", "email", "")

	}

	if input.AlternateEmail != "" && !IsEmailAddress(input.Email){
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "format email tidak sesuai", "alternate_email", "")
	}

	if len(input.TaxName) > 100 && input.TaxName != ""{
		return errorModel.GenerateFieldHaveMaxLimitError(fileName, funcName, "tax_name", 100)
	}

	if len(input.TaxAddress) > 256 && input.TaxAddress != ""{
		return errorModel.GenerateFieldHaveMaxLimitError(fileName, funcName, "tax_address", 256)
	}

	if isUpdate {
		if util.IsStringEmpty(input.UpdatedAtStr) {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
		}

		input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
		if err.Error != nil {
			return
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func IsEmailAddress(input string) (output bool) {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegexp.MatchString(input)
}

func IsPhoneValid(input string) (output bool) {
	phoneRegexp := regexp.MustCompile(`\+62-[0-9]{9,}$`)
	return phoneRegexp.MatchString(input)
}