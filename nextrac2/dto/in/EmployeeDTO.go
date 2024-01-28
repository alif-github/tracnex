package in

import (
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type EmployeeRequest struct {
	AbstractDTO
	ID                    int64    `json:"id"`
	NIK                   int64    `json:"nik"`
	RedmineId             int64    `json:"redmine_id"`
	Name                  string   `json:"name"`
	IDCard                string   `json:"id_card"`
	FirstName             string   `json:"first_name"`
	LastName              string   `json:"last_name"`
	Gender                string   `json:"gender"`
	Email                 string   `json:"email"`
	Phone                 string   `json:"phone"`
	PlaceOfBirth          string   `json:"place_of_birth"`
	DateOfBirthStr        string   `json:"date_of_birth"`
	Type                  string   `json:"type"`
	MothersMaiden         string   `json:"mothers_maiden"`
	TaxMethod             string   `json:"tax_method"`
	AddressResidence      string   `json:"address_residence"`
	AddressTax            string   `json:"address_tax"`
	NPWP                  string   `json:"npwp"`
	Religion              string   `json:"religion"`
	DateJoinStr           string   `json:"date_join"`
	DateOutStr            string   `json:"date_out"`
	ReasonResignation     string   `json:"reason_resignation"`
	Status                string   `json:"status"`
	Position              int64    `json:"position_id"`
	MaritalStatus         string   `json:"marital_status"`
	NumberOfDependents    int64    `json:"number_of_dependents"`
	Nationality           string   `json:"nationality"`
	DepartmentId          int64    `json:"department_id"`
	NoBpjs                string   `json:"bpjs_no"`
	NoBpjsTk              string   `json:"bpjs_tk_no"`
	Level                 int64    `json:"level_id"`
	Grade                 int64    `json:"grade_id"`
	IsHaveMember          bool     `json:"is_have_member"`
	MemberID              []string `json:"member"`
	MandaysRate           float64  `json:"mandays_rate"`
	MandaysRateAutomation float64  `json:"mandays_rate_automation"`
	MandaysRateManual     float64  `json:"mandays_rate_manual"`
	UpdatedAtStr          string   `json:"updated_at"`
	Education             string   `json:"education"`
	Active                bool     `json:"active"`
	MemberIDStr           string
	Files                 []MultipartFileDTO
	UpdatedAt             time.Time
	DateJoin              time.Time
	DateOut               time.Time
	DateOfBirth           time.Time
}

type EmployeeJSONDB struct {
	IDCard             string `json:"id_card"`
	NPWP               string `json:"npwp"`
	DepartmentID       int64  `json:"department_id"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	Gender             string `json:"gender"`
	PlaceOfBirth       string `json:"place_of_birth"`
	DateOfBirth        string `json:"date_of_birth"`
	AddressResidence   string `json:"address_residence"`
	AddressTax         string `json:"address_tax"`
	DateJoin           string `json:"date_join"`
	DateOut            string `json:"date_out"`
	Religion           string `json:"religion"`
	Type               string `json:"type"`
	Status             string `json:"status"`
	EmployeePositionID int64  `json:"employee_position_id"`
	MaritalStatus      string `json:"marital_status"`
	Education          string `json:"education"`
	MothersMaiden      string `json:"mothers_maiden"`
	NumberOfDependents int64  `json:"number_of_dependents"`
	Nationality        string `json:"nationality"`
	TaxMethod          string `json:"tax_method"`
	ReasonResignation  string `json:"reason_resignation"`
	IsHaveMember       bool   `json:"is_have_member"`
	Member             string `json:"member"`
	Active             bool   `json:"active"`
	Photo              string `json:"photo"`
	NoBpjs             string `json:"bpjs_no"`
	NoBpjsTk           string `json:"bpjs_tk_no"`
	Level              int64  `json:"employee_level_id"`
	Grade              int64  `json:"employee_grade_id"`
	FileUploadID       int64  `json:"file_upload_id"`
}

type TrackerDeveloper struct {
	Task float64 `json:"task"`
}

type TrackerQA struct {
	Automation float64 `json:"automation"`
	Manual     float64 `json:"manual"`
}

type MemberList struct {
	MemberID []string `json:"member_id"`
}

func (input *EmployeeRequest) ValidateInsert() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeeDTO"
		funcName = "ValidateInsert"
	)

	return input.mandatoryFieldValidation(fileName, funcName)
}

func (input *EmployeeRequest) ValidateUpdateEmployeeTimeSheet() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeeDTO"
		funcName = "ValidateUpdateEmployeeTimeSheet"
	)

	if err = input.validationForUpdateAndDelete(fileName, funcName); err.Error != nil {
		return
	}

	return input.mandatoryFieldValidationTimeSheet(fileName, funcName)
}

func (input *EmployeeRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeeDTO"
		funcName = "ValidateUpdate"
	)

	if err = input.validationForUpdateAndDelete(fileName, funcName); err.Error != nil {
		return
	}

	return input.mandatoryFieldValidation(fileName, funcName)
}

func (input *EmployeeRequest) ValidateViewEmployee() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeeDTO"
		funcName = "ValidateViewEmployee"
	)

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *EmployeeRequest) mandatoryFieldValidation(fileName string, funcName string) (err errorModel.ErrorModel) {
	var (
		isValid bool
		errMsg  string
		addInfo string
	)

	//--- ID Card : Check Empty
	if util.IsStringEmpty(input.IDCard) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.IDCard)
		return
	}

	//--- ID Card : Check Min Max
	err = input.ValidateMinMaxString(input.IDCard, constanta.IDCard, 1, 50)
	if err.Error != nil {
		return
	}

	//--- ID Card : Check Is Only Alfa Numeric
	isValid, errMsg, addInfo = IsOnlyAlfaNumerikValid(input.IDCard)
	if !isValid {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errMsg, constanta.IDCard, addInfo)
		return
	}

	//--- Department : Department ID Check
	if input.DepartmentId < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.DepartmentId)
		return
	}

	//--- FirstName : Check Empty
	if util.IsStringEmpty(input.FirstName) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.FirstName)
	}

	//--- FirstName : Check Min Max
	err = input.ValidateMinMaxString(input.FirstName, constanta.FirstName, 1, 50)
	if err.Error != nil {
		return
	}

	//--- FirstName : Is Nexsoft Name Standard Valid ?
	isValid = false
	errMsg = ""
	addInfo = ""
	isValid, errMsg, addInfo = util.IsNexsoftNameStandardValid(input.FirstName)
	if !isValid {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errMsg, constanta.FirstName, addInfo)
		return
	}

	//--- LastName : Check Empty
	if !util.IsStringEmpty(input.LastName) {
		//--- LastName : Check Min Max
		err = input.ValidateMinMaxString(input.LastName, constanta.LastName, 1, 100)
		if err.Error != nil {
			return
		}

		//--- LastName : Is Nexsoft Name Standard Valid ?
		isValid = false
		errMsg = ""
		addInfo = ""
		isValid, errMsg, addInfo = util.IsNexsoftNameStandardValid(input.LastName)
		if !isValid {
			err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errMsg, constanta.LastName, addInfo)
			return
		}
	}

	//--- Gender : Check Enum
	if util.IsStringEmpty(input.Gender) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Sex)
		return
	}

	//--- Gender : Check Option
	if input.Gender != "Male" && input.Gender != "Female" && input.Gender != "Neutral" {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "GENDER_RULE", constanta.Sex, "")
		return
	}

	//--- NPWP : Check Empty
	if util.IsStringEmpty(input.NPWP) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.NPWP)
		return
	}

	//--- NPWP : Check Min Max
	if err = input.ValidateMinMaxString(input.NPWP, constanta.NPWP, 1, 38); err.Error != nil {
		return
	}

	//--- NPWP : Check Regex
	npwpRegex := regexp.MustCompile(fmt.Sprintf(`^[0-9]{15,16}$`))
	if !npwpRegex.MatchString(input.NPWP) {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "MUST_DIGIT_NPWP_RULE", constanta.NPWP, "")
		return
	}

	//--- Religion Check Empty
	if util.IsStringEmpty(input.Religion) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Religion)
		return
	}

	//--- Religion Check Option
	if input.Religion != "Islam" && input.Religion != "Protestant" && input.Religion != "Catholic" && input.Religion != "Hindu" && input.Religion != "Buddha" && input.Religion != "Confucius" && input.Religion != "Others" {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "RELIGION_RULE", constanta.Religion, "")
		return
	}

	//--- Date Join Check Empty
	if util.IsStringEmpty(input.DateJoinStr) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.DateJoin)
		return
	}

	//--- Date Join Check Format Time
	input.DateJoin, err = TimeStrToTimeWithTimeFormat(input.DateJoinStr, constanta.DateJoin, constanta.DefaultInstallationTimeFormat)
	if err.Error != nil {
		return
	}

	//--- Date Out Check Empty
	if !util.IsStringEmpty(input.DateOutStr) {
		//--- Date Out Check Format Time
		input.DateOut, err = TimeStrToTimeWithTimeFormat(input.DateOutStr, constanta.DateOut, constanta.DefaultInstallationTimeFormat)
		if err.Error != nil {
			return
		}

		//--- Date Join And Date Out Compare
		if input.DateOut.Before(input.DateJoin) {
			err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "DATE_JOIN_OUT_RULE", constanta.DateOut, "")
			return
		}
	}

	//--- Status Check Empty
	if util.IsStringEmpty(input.Status) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Status)
		return
	}

	//--- Status Check Option
	if input.Status != "Probation" && input.Status != "PKWTT" && input.Status != "PKWT" {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "STATUS_RULE", constanta.Status, "")
		return
	}

	//--- Marital Status Check Empty
	if util.IsStringEmpty(input.MaritalStatus) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Marital)
		return
	}

	//--- Marital Status Check Option
	if input.MaritalStatus != "Single" && input.MaritalStatus != "Married" {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "MARITAL_STATUS_RULE", constanta.Marital, "")
		return
	}

	//--- Nationality Check Empty
	if util.IsStringEmpty(input.Nationality) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Nationality)
		return
	}

	//--- Nationality Check Option
	if input.Nationality != "WNI" && input.Nationality != "WNA" {
		err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NATIONALITY_RULE", constanta.Nationality, "")
		return
	}

	//--- BPJS Check Empty
	if !util.IsStringEmpty(input.NoBpjs) {
		//--- BPJS Set Default
		isValid = false
		errMsg = ""
		addInfo = ""

		//--- BPJS Check Min Max
		if err = input.ValidateMinMaxString(input.NoBpjs, constanta.Bpjs, 1, 25); err.Error != nil {
			return
		}

		//--- BPJS Check Digit Valid
		isValid, errMsg, addInfo = IsOnlyDigitValid(input.NoBpjs)
		if !isValid {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errMsg, constanta.Bpjs, addInfo)
		}
	}

	//--- BPJS TK Check Empty
	if !util.IsStringEmpty(input.NoBpjsTk) {
		//--- BPJS TK Set Default
		isValid = false
		errMsg = ""
		addInfo = ""

		//--- BPJS TK Check Min Max
		if err = input.ValidateMinMaxString(input.NoBpjsTk, constanta.BpjsTk, 1, 25); err.Error != nil {
			return
		}

		//--- BPJS TK Check Digit Valid
		isValid, errMsg, addInfo = IsOnlyDigitValid(input.NoBpjsTk)
		if !isValid {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errMsg, constanta.BpjsTk, addInfo)
		}
	}

	//--- Level Check Empty
	if input.Level < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.Level)
		return
	}

	//--- Grade Check Empty
	if input.Grade < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.Grade)
		return
	}

	var member MemberList
	if input.IsHaveMember {
		//--- Member List Check
		for i, itemMember := range input.MemberID {
			if itemMember == "all" {
				if i > 0 {
					err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "ALL_MEMBER_RULE", constanta.Member, "")
					return
				}
				member = MemberList{MemberID: []string{itemMember}}
				break
			} else {
				num, errs := strconv.Atoi(itemMember)
				if errs != nil {
					err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "ID_MEMBER_RULE", constanta.Member, "")
					return
				}
				if num == 0 {
					err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "MEMBER_NOT_NULL_RULE", constanta.Member, "")
					return
				}
				member.MemberID = append(member.MemberID, itemMember)
			}
		}

		//--- JSON Marshal Member ID
		if len(input.MemberID) > 0 {
			strByte, _ := json.Marshal(member)
			input.MemberIDStr = string(strByte)
		}
	}

	//--- Education Check Empty
	if !util.IsStringEmpty(input.Education) {
		//--- Education : Min Max String
		err = input.ValidateMinMaxString(input.Education, constanta.Education, 1, 10)
		if err.Error != nil {
			return
		}

		//--- Education : Alfa Numeric Valid
		isValid = false
		errMsg = ""
		addInfo = ""
		isValid, errMsg, addInfo = IsOnlyAlfaNumerikValid(input.Education)
		if !isValid {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errMsg, constanta.Education, addInfo)
		}
	}

	//--- Email
	if !util.IsStringEmpty(input.Email) {
		if err = input.ValidateMinMaxString(input.Email, constanta.Email, 1, 100); err.Error != nil {
			return
		}

		if !util.IsEmailAddress(input.Email) {
			err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Email)
			return
		}
	}

	//--- Phone
	if !util.IsStringEmpty(input.Phone) {
		if err = input.ValidateMinMaxString(input.Phone, constanta.Phone, 1, 18); err.Error != nil {
			return
		}

		if !util.IsPhoneNumberWithCountryCode(input.Phone) {
			err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.PhoneRegex, constanta.Phone, "")
			return
		}

		splitPhone := strings.Split(input.Phone, "-")
		if !util.IsCountryCode(splitPhone[0]) {
			err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.PhoneCode)
			return
		}

		phone, isValid := util.IsPhoneNumber(splitPhone[1])
		if !isValid {
			err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.Phone)
			return
		}

		input.Phone = splitPhone[0] + "-" + strconv.Itoa(phone)
	}

	//--- Place Of Birth
	if !util.IsStringEmpty(input.PlaceOfBirth) {
		if err = input.ValidateMinMaxString(input.PlaceOfBirth, constanta.BirthPlace, 1, 256); err.Error != nil {
			return
		}
	}

	//--- Date Of Birth Check Empty
	if !util.IsStringEmpty(input.DateOfBirthStr) {
		//--- Date Of Birth Check Format Time
		input.DateOfBirth, err = TimeStrToTimeWithTimeFormat(input.DateOfBirthStr, constanta.BirthDate, constanta.DefaultInstallationTimeFormat)
		if err.Error != nil {
			return
		}
	}

	//--- Type Check Empty
	if !util.IsStringEmpty(input.Type) {
		if err = input.ValidateMinMaxString(input.Type, constanta.Type, 1, 50); err.Error != nil {
			return
		}
	}

	//--- Mothers Maiden Check Empty
	if !util.IsStringEmpty(input.MothersMaiden) {
		if err = input.ValidateMinMaxString(input.MothersMaiden, constanta.MothersMaiden, 1, 50); err.Error != nil {
			return
		}
	}

	//--- Dependants if is not 0
	if input.NumberOfDependents > 0 {
		if input.NumberOfDependents > 99 {
			err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NUMBERS_DEPENDANTS_RULE", constanta.NumbersDependants, "")
			return
		}
	}

	//--- Tax Method Check Empty
	if !util.IsStringEmpty(input.TaxMethod) {
		if err = input.ValidateMinMaxString(input.TaxMethod, constanta.TaxMethod, 1, 50); err.Error != nil {
			return
		}
	}

	//--- Address Check
	if !util.IsStringEmpty(input.AddressResidence) {
		if err = input.ValidateMinMaxString(input.AddressResidence, constanta.Address, 1, 256); err.Error != nil {
			return
		}
	}

	//--- Address Tax Check
	if !util.IsStringEmpty(input.AddressTax) {
		if err = input.ValidateMinMaxString(input.AddressTax, constanta.AddressTax, 1, 256); err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *EmployeeRequest) mandatoryFieldValidationTimeSheet(fileName string, funcName string) (err errorModel.ErrorModel) {
	//--- Redmine ID Check Is Empty
	if input.RedmineId < 1 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.RedmineId)
	}

	//--- Redmine ID Validate Min Max Integer
	err = util2.ValidateMinMaxInteger(input.RedmineId, constanta.RedmineId, 1, 5)
	if err.Error != nil {
		return
	}

	//--- ID Card Check Empty
	if util.IsStringEmpty(input.IDCard) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.IDCard)
		return
	}

	//--- ID Card Check Min Max
	err = input.ValidateMinMaxString(input.IDCard, constanta.IDCard, 1, 50)
	if err.Error != nil {
		return
	}

	////--- Name Check Empty
	//if util.IsStringEmpty(input.Name) {
	//	return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Name)
	//}

	//--- Department Check
	if input.DepartmentId < 1 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.DepartmentId)
	}

	//--- Mandays Rate Check
	if input.DepartmentId == constanta.QAQCDepartmentID {
		//open validation mandays rate with minimum value 0
		if input.MandaysRateAutomation < 0 || input.MandaysRateManual < 0 {
			err = errorModel.GenerateSimpleErrorModel(400, "Data Mandays Rate Automation atau Mandays Rate Manual tidak boleh bernilai kurang dari 0")
			return
		}

		err = util2.ValidateMinMaxFloat(input.MandaysRateAutomation, constanta.MandaysRateAutomation, 0, 7)
		if err.Error != nil {
			return
		}

		err = util2.ValidateMinMaxFloat(input.MandaysRateManual, constanta.MandaysRateManual, 0, 7)
		if err.Error != nil {
			return
		}
	} else {
		//open validation mandays rate with minimum value 0
		if input.MandaysRate < 0 {
			err = errorModel.GenerateSimpleErrorModel(400, "Data Mandays Rate tidak boleh kurang dari 0")
			return
		}

		err = util2.ValidateMinMaxFloat(input.MandaysRate, constanta.MandaysRate, 0, 7)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *EmployeeRequest) validationForUpdateAndDelete(fileName string, funcName string) (err errorModel.ErrorModel) {
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

func (input *EmployeeRequest) ValidateDelete() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeeDTO.go"
		funcName = "ValidateDelete"
	)

	return input.validationForUpdateAndDelete(fileName, funcName)
}

func (input *EmployeeRequest) ValidateGetEmployeeTimeSheetRedmineByNIK() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeeDTO"
		funcName = "ValidateGetEmployeeTimeSheetRedmineByNIK"
	)

	//--- ID Card
	if util.IsStringEmpty(input.IDCard) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.IDCard)
	}

	//--- Department ID
	if input.DepartmentId < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.DepartmentId)
	}

	return
}
