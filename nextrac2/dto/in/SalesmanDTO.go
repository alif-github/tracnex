package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
	"time"
)

type SalesmanRequest struct {
	AbstractDTO
	ID            int64  `json:"id"`
	PersonTitleID int64  `json:"person_title_id"`
	Sex           string `json:"sex"`
	Nik           string `json:"nik"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Address       string `json:"address"`
	Hamlet        string `json:"hamlet"`
	Neighbourhood string `json:"neighbourhood"`
	ProvinceID    int64  `json:"province_id"`
	DistrictID    int64  `json:"district_id"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	Status        string `json:"status"`
	CreatedAtStr  string `json:"created_at"`
	UpdatedAtStr  string `json:"updated_at"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (input *SalesmanRequest) ValidationInsertSalesman() (err errorModel.ErrorModel) {
	var (
		fileName = input.fileNameFuncNameSalesman()
		funcName = "ValidationInsertSalesman"
	)

	if util.IsStringEmpty(input.Nik) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.NIK)
	}

	if !util.IsNIKValid(input.Nik) {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.NIKRegex, constanta.NIK, "")
	}

	return input.mandatoryValidation()
}

func (input SalesmanRequest) ValidateViewSalesman() (err errorModel.ErrorModel) {
	fileName := input.fileNameFuncNameSalesman()
	funcName := "ValidateViewSalesman"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ID)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *SalesmanRequest) ValidationDeleteSalesman() (err errorModel.ErrorModel) {
	fileName := input.fileNameFuncNameSalesman()
	funcName := "ValidationDeleteSalesman"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ID)
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

func (input *SalesmanRequest) ValidationUpdateSalesman() (err errorModel.ErrorModel) {
	fileName := input.fileNameFuncNameSalesman()
	funcName := "ValidationUpdateSalesman"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ID)
	}

	if util.IsStringEmpty(input.UpdatedAtStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return input.mandatoryValidation()
}

func (input *SalesmanRequest) mandatoryValidation() (err errorModel.ErrorModel) {
	var (
		fileName = input.fileNameFuncNameSalesman()
		funcName = "mandatoryValidation"
		isValid  bool
	)

	if input.PersonTitleID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.PersonTitleID)
	}

	if util.IsStringEmpty(input.Sex) {
		input.Sex = "N"
	}

	if (input.Sex != "P") && (input.Sex != "L") && (input.Sex != "N") {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.SexRegex, constanta.Sex, "")
	}

	if util.IsStringEmpty(input.FirstName) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.FirstName)
	}

	err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.FirstName, input.FirstName)
	if err.Error != nil {
		return
	}

	err = input.ValidateMinMaxString(input.FirstName, constanta.FirstName, 1, 20)
	if err.Error != nil {
		return
	}

	//isValid, errField = util.IsNexsoftProfileNameStandardValid(input.FirstName)
	//if !isValid {
	//	return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.FirstName, "")
	//}

	if !util.IsStringEmpty(input.LastName) {
		err = input.ValidateMinMaxString(input.LastName, constanta.LastName, 1, 35)
		if err.Error != nil {
			return
		}

		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.LastName, input.LastName)
		if err.Error != nil {
			return
		}

		//isValid, errField = util.IsNexsoftProfileNameStandardValid(input.LastName)
		//if !isValid {
		//	return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.LastName, "")
		//}
	}

	if util.IsStringEmpty(input.Address) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Address)
	}

	err = input.ValidateMinMaxString(input.Address, constanta.Address, 10, 256)
	if err.Error != nil {
		return
	}

	if !util.IsStringEmpty(input.Hamlet) {
		err = input.ValidateMinMaxString(input.Hamlet, constanta.Hamlet, 1, 5)
		if err.Error != nil {
			return
		}

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

		err = util2.ValidateSpecialCharacter(fileName, funcName, constanta.Neighbourhood, input.Neighbourhood)
		if err.Error != nil {
			return
		}
	}

	if input.ProvinceID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ProvinceID)
	}

	if input.DistrictID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.DistrictID)
	}

	if util.IsStringEmpty(input.Phone) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Phone)
	}

	isValid = IsPhoneNumberWithCountryCodeMDB(input.Phone)
	if !isValid {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.PhoneRegex, constanta.Phone, "")
	}

	if util.IsStringEmpty(input.Email) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Email)
	}

	if !IsEmailAddressMDB(input.Email) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
	}

	err = input.ValidateMinMaxString(input.Email, constanta.Email, 5, 100)
	if err.Error != nil {
		return err
	}

	if input.Email != "" {
		input.Email = strings.ToLower(input.Email)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input SalesmanRequest) fileNameFuncNameSalesman() (fileName string) {
	return "SalesmanDTO.go"
}
