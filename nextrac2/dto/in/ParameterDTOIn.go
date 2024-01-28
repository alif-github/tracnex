package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ParameterRequest struct {
	ID               int64  `json:"id"`
	ParameterGroupID int64  `json:"parameter_group_id"`
	Permission       string `json:"permission"`
	Name             string `json:"name"`
	NameEn           string `json:"name_en"`
	Value            string `json:"value"`
	Code             string `json:"code"`
	Sequence         int32  `json:"sequence"`
	Level            int32  `json:"level"`
	Type             string `json:"type"`
	Length           int32  `json:"length"`
	MinVal           int32  `json:"min_val"`
	MaxVal           int32  `json:"max_val"`
	DefaultVal       string `json:"default_val"`
	SelectListValue  string `json:"select_list_value"`
	SelectListName   string `json:"select_list_name"`
	Description      string `json:"description"`
	UpdateAtStr      string `json:"updated_at"`
	UpdatedAt         time.Time
}

func (input *ParameterRequest) ValidateViewParameter() (err errorModel.ErrorModel) {
	fileName := "ParameterDTO.go"
	funcName := "ValidateViewParameter"
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *ParameterRequest) ValidateInsertParameter() (err errorModel.ErrorModel) {
	//fileName := "ParameterDTO.go"
	//funcName := "ValidateInsertParameter"
	//

	return errorModel.GenerateNonErrorModel()
}

func (input *ParameterRequest) ValidateUpdateParameter() (err errorModel.ErrorModel) {
	//fileName := "ParameterDTO.go"
	//funcName := "ValidateInsertParameter"
	//

	return errorModel.GenerateNonErrorModel()
}

func (input *ParameterRequest) ValidateDeleteParameter() (err errorModel.ErrorModel) {
	//fileName := "ParameterDTO.go"
	//funcName := "ValidateInsertParameter"
	//

	return errorModel.GenerateNonErrorModel()
}

func (input *ParameterRequest) ValidateUpdateParameterEmployee() (err errorModel.ErrorModel) {
	fileName := "ParameterDTO.go"
	funcName := "ValidateUpdateParameterEmployee"
	if input.Value == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "value id "+strconv.Itoa(int(input.ID)))
	}

	validValue := regexp.MustCompile(`\d{2}-\d{2}`)
	isValidValue := validValue.MatchString(input.Value)
	if !isValidValue && input.Name == "cutOffAnualLeave" {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "format yang digunakan dd-mm", "value", "")
	}

	if input.Name == "cutOffAnualLeave" {
		dateArr := strings.Split(input.Value, "-")
		dd, _ := strconv.ParseInt(dateArr[0], 10, 64)
		mm, _ := strconv.ParseInt(dateArr[1], 10, 64)
		if dd == 0 || mm == 0 || dd > 31 || mm > 12{
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "format yang digunakan dd-mm", "value", "")
		}
	}

	if util.IsStringEmpty(input.UpdateAtStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	if (input.Value != "1" && input.Value != "0") && input.Name == "anualLeaveAfterProbation" {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "value hanya diisi dengan 1 (ya) atau 0 (tidak)", "value anualLeaveAfterProbation", "")
	}

	if input.Name == "expiredMedicalClaim"{
		i, e:= strconv.ParseInt(input.Value, 10, 64)
		if e != nil {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "harus menggunakan angka", "value", "")
		}

		if i >= 1000{
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "maksimal 999", "value", "")
		}
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdateAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}