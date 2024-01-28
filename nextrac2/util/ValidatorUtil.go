package util

import (
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

func ValidateStringWithMinMaxMandatory(fileName string, funcName string, input string, fieldName string, minLength int, maxLength int) (_ interface{}, output errorModel.ErrorModel) {
	var err errorModel.ErrorModel

	validationResult := util.IsStringEmpty(input)
	if validationResult {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, fieldName)
		return nil, err
	}

	err = ValidateMinMaxString(input, fieldName, minLength, maxLength)
	if err.Error != nil {
		return nil, err
	}

	return nil, errorModel.GenerateNonErrorModel()
}

func ValidateStringWithMinMaxOptional(fileName string, funcName string, input string, fieldName string, minLength int, maxLength int) (_ interface{}, output errorModel.ErrorModel) {
	var err errorModel.ErrorModel

	if input != "" {
		validationResult := util.IsStringEmpty(input)
		if validationResult {
			err = errorModel.GenerateEmptyFieldError(fileName, funcName, fieldName)
			return nil, err
		}

		err = ValidateMinMaxString(input, fieldName, minLength, maxLength)
		if err.Error != nil {
			return nil, err
		}
	}

	return nil, errorModel.GenerateNonErrorModel()
}

func ValidateDateTimeOptional(fileName string, funcName string, input string, fieldName string, _ int, _ int) (result interface{}, output errorModel.ErrorModel) {
	var err errorModel.ErrorModel
	var timeResult time.Time

	if input != "" {
		timeResult = DateConvert(input)
		validationResult := timeResult.IsZero()

		if validationResult {
			err = errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
			return timeResult, err
		}
	}

	return timeResult, errorModel.GenerateNonErrorModel()
}

func ValidateDateWithFindMandatory(fileName string, funcName string, input string, fieldName string, _ int, _ int) (result interface{}, output errorModel.ErrorModel) {
	var err errorModel.ErrorModel
	var timeResult time.Time

	timeResult = DateConvert(input)
	validationResult := timeResult.IsZero()

	if validationResult {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
		return timeResult, err
	}

	return timeResult, errorModel.GenerateNonErrorModel()
}

func ValidateParseInteger(_ string, _ string, input string, _ string, _ int, _ int) (result interface{}, output errorModel.ErrorModel) {

	parseIntegerStr, _ := strconv.Atoi(input)

	return parseIntegerStr, errorModel.GenerateNonErrorModel()
}

func ValidateMinMaxString(inputStr string, fieldName string, min int, max int) errorModel.ErrorModel {
	fileName := "ValidatorUtil.go"
	funcName := "ValidateMinMaxString"

	if len(inputStr) < min {
		if min == 1 {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, fieldName)
		} else {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_MORE_THAN", fieldName, strconv.Itoa(min))
		}
	}
	if len(inputStr) > max {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_LESS_THAN", fieldName, strconv.Itoa(max))
	}

	return errorModel.GenerateNonErrorModel()
}

func ValidateMinMaxInteger(inputData int64, fieldName string, min int, max int) errorModel.ErrorModel {
	fileName := "ValidatorUtil.go"
	funcName := "ValidateMinMaxInteger"

	if len(strconv.Itoa(int(inputData))) > max {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_LESS_THAN", fieldName, strconv.Itoa(max))
	} else if len(strconv.Itoa(int(inputData))) < min {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_MORE_THAN", fieldName, strconv.Itoa(min))
	}

	return errorModel.GenerateNonErrorModel()
}

func ValidateMinMaxFloat(inputData float64, fieldName string, min int, max int) errorModel.ErrorModel {
	fileName := "ValidatorUtil.go"
	funcName := "ValidateMinMaxInteger"

	if len(strconv.Itoa(int(inputData))) > max {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_LESS_THAN", fieldName, strconv.Itoa(max))
	} else if len(strconv.Itoa(int(inputData))) < min {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_MORE_THAN", fieldName, strconv.Itoa(min))
	}

	return errorModel.GenerateNonErrorModel()
}

func ValidateMinMax(inputData interface{}, fieldName string, min int, max int) errorModel.ErrorModel {
	fileName := "ValidatorUtil.go"
	funcName := "ValidateMinMaxString"

	switch inputData.(type) {
	case string:
		inputStr := inputData.(string)
		if len(inputStr) < min {
			if min == 1 {
				return errorModel.GenerateEmptyFieldError(fileName, funcName, fieldName)
			} else {
				return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_MORE_THAN", fieldName, strconv.Itoa(min))
			}
		}
		if len(inputStr) > max {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_LESS_THAN", fieldName, strconv.Itoa(max))
		}
		break
	case int, int64, int32, int16, int8:
		valueData := reflect.ValueOf(inputData)
		if int(valueData.Int()) < min {
			if min == 1 {
				return errorModel.GenerateEmptyFieldError(fileName, funcName, fieldName)
			} else {
				return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_MORE_THAN", fieldName, strconv.Itoa(min))
			}
		}
		if int(valueData.Int()) > max {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_LESS_THAN", fieldName, strconv.Itoa(max))
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func IsDataEmpty(inputData interface{}) bool {
	valueData := reflect.ValueOf(inputData)
	return valueData.IsZero()
}

func DateConvert(date string) (output time.Time) {
	const (
		layoutISO = "2006-01-02"
		layoutUS  = "2006-01-02T15:04:05Z"
	)

	t, _ := time.Parse(layoutISO, date)
	timeString := t.Format(layoutUS)
	output, _ = time.Parse(layoutUS, timeString)
	return
}

func IsFieldNumericEmpty(input int64) bool {
	return input <= 0
}

func ValidateSpecialCharacter(fileName string, funcName string, fieldName string, userInput string) errorModel.ErrorModel {
	//result, errorS := regexp.Compile("[a-zA-Z0-9]+[^:=,.\\-/&)\\\\(;@?%!#~\\|$^*\\`/]+$")
	//result, errorS := regexp.Compile("^[a-zA-Z0-9 ]+[a-zA-Z0-9-._@ ]*$")
	result, errorS := regexp.Compile("^[a-zA-Z0-9]+[a-zA-Z0-9-._@ ]*$")

	if errorS != nil {
		return errorModel.GenerateUnknownError(fileName, funcName, errorS)
	}

	if !result.MatchString(userInput) {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.ErrorSpecialCharacter, fieldName, "")
	}

	return errorModel.GenerateNonErrorModel()
}

func ValidateWhiteListSpecialCharacter(fileName string, funcName string, fieldName string, userInput string) errorModel.ErrorModel {
	result, errorS := regexp.Compile("^[a-zA-Z0-9-]*$")

	if errorS != nil {
		return errorModel.GenerateUnknownError(fileName, funcName, errorS)
	}

	if !result.MatchString(userInput) {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.ErrorWhitelistSpecialCharacter, fieldName, "")
	}

	return errorModel.GenerateNonErrorModel()
}

func ValidateEmptyField(fileName string, funcName string, fieldName string, userInput string) errorModel.ErrorModel {
	result, errorS := regexp.Compile(fmt.Sprintf(`^\s*$`))

	if errorS != nil {
		return errorModel.GenerateUnknownError(fileName, funcName, errorS)
	}

	if result.MatchString(userInput) {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.ErrorEmptyString, fieldName, "")
	}

	return errorModel.GenerateNonErrorModel()
}

func ValidateSpecialCharacterAlphabet(fileName string, funcName string, fieldName string, userInput string) errorModel.ErrorModel {
	result, errorS := regexp.Compile("^[a-zA-Z]+[a-zA-Z ]*$")

	if errorS != nil {
		return errorModel.GenerateUnknownError(fileName, funcName, errorS)
	}

	if !result.MatchString(userInput) {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.ErrorSpecialCharacterAlphabet, fieldName, "")
	}

	return errorModel.GenerateNonErrorModel()
}

func ValidateSpecialCharacterTrimSpace(fileName string, funcName string, fieldName string, userInput string) errorModel.ErrorModel {
	//result, errorS := regexp.Compile("[a-zA-Z0-9]+[^:=,.\\-/&)\\\\(;@?%!#~\\|$^*\\`/]+$")
	result, errorS := regexp.Compile("^[a-zA-Z0-9]+[a-zA-Z0-9-._@]*$")

	if errorS != nil {
		return errorModel.GenerateUnknownError(fileName, funcName, errorS)
	}

	if !result.MatchString(userInput) {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.ErrorSpecialCharacter, fieldName, "")
	}

	return errorModel.GenerateNonErrorModel()
}

// ---------------------------------------------------------------------------------------------------------------------
// ---------------------------------------------------------------------------------------------------------------------

//func ValidateProfileNameWithMinMaxRegex(fileName string, funcName string, input string, fieldName string, minLength int, maxLength int) (output errorModel.ErrorModel) {
//	if !util.IsStringEmpty(input) || minLength != 0 {
//		isValid, msg := util.IsNexsoftProfileNameStandardValid(input)
//		if !isValid {
//			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, msg, fieldName, "")
//		}
//	}
//
//	return ValidateStringWithMinMaxMandatory(fileName, funcName, input, fieldName, minLength, maxLength)
//}
//
//func ValidateCountryPhoneCodeRegex(fileName string, funcName string, input string, _ string, _ int, _ int) (output errorModel.ErrorModel) {
//	isValid := util.IsCountryCode(input)
//	if !isValid {
//		return errorModel.GenerateFormatFieldError(fileName, funcName, "PHONE")
//	}
//	return errorModel.GenerateNonErrorModel()
//}
//
//func ValidateIsNumeric(fileName string, funcName string, input string, fieldName string, minLength int, _ int) (err errorModel.ErrorModel) {
//	if !util.IsStringEmpty(input) || minLength != 0 {
//		_, isNumeric := util.IsNumeric(input)
//		if !isNumeric {
//			return errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
//		}
//	}
//	return errorModel.GenerateNonErrorModel()
//}
//
//func ValidateIsNumericWithMinMax(fileName string, funcName string, input string, fieldName string, minLength int, maxLength int) (err errorModel.ErrorModel) {
//	if !util.IsStringEmpty(input) {
//		err = ValidateIsNumeric(fileName, funcName, input, fieldName, minLength, maxLength)
//		if err.Error != nil {
//			return
//		}
//
//		err = ValidateStringWithMinMax(fileName, funcName, input, fieldName, minLength, maxLength)
//		if err.Error != nil {
//			return
//		}
//	}
//
//	return errorModel.GenerateNonErrorModel()
//}
//
//func ValidatePhone(fileName string, funcName string, input string, fieldName string, minLength int, _ int) (err errorModel.ErrorModel) {
//	if !util.IsStringEmpty(input) || minLength != 0 {
//		phoneSplit := strings.Split(input, "-")
//		if len(phoneSplit) != 2 {
//			return errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
//		}
//
//		if !util.IsCountryCode(phoneSplit[0]) {
//			return errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
//		}
//
//		_, validationResult := util.IsPhoneNumber(phoneSplit[1])
//		if !validationResult {
//			return errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
//		}
//	}
//
//	return errorModel.GenerateNonErrorModel()
//}
//
//func ValidateEmailWithMinMax(fileName string, funcName string, input string, fieldName string, minLength int, maxLength int) (err errorModel.ErrorModel) {
//	if !util.IsStringEmpty(input) || minLength != 0 {
//		if !util.IsEmailAddress(input) {
//			return errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
//		}
//
//		err = ValidateStringWithMinMax(fileName, funcName, input, fieldName, minLength, maxLength)
//		if err.Error != nil {
//			return
//		}
//	}
//	return errorModel.GenerateNonErrorModel()
//}
//
//func ValidateNPWP(fileName string, funcName string, input string, _ string, minLength int, _ int) (err errorModel.ErrorModel) {
//	if !util.IsStringEmpty(input) || minLength != 0 {
//		if !util.IsNPWPValid(input) {
//			return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.NPWP)
//		}
//	}
//	return errorModel.GenerateNonErrorModel()
//}
//
//func ValidateText(_ string, _ string, _ string, _ string, _ int, _ int) (err errorModel.ErrorModel) {
//	return errorModel.GenerateNonErrorModel()
//}
//
//func ValidateAdditionalInfo(fileName string, funcName string, input string, fieldName string, _ int, _ int) (err errorModel.ErrorModel) {
//	var additionalInfo []model.AdditionalInformation
//	if !util.IsStringEmpty(input) {
//		_ = json.Unmarshal([]byte(input), &additionalInfo)
//		if len(additionalInfo) < 1 {
//			err = errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
//			return
//		}
//		for i := 0; i < len(additionalInfo); i++ {
//			validationResult, errField := util.IsNexsoftAdditionalInformationKeyStandardValid(additionalInfo[i].Key)
//			if !validationResult {
//				return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, fieldName+" ["+strconv.Itoa(i+1)+"]", "")
//			}
//
//			for j := 1 + i; j < len(additionalInfo); j++ {
//				if additionalInfo[j].Key == additionalInfo[i].Key {
//					if additionalInfo[j].Value == additionalInfo[i].Value {
//						additionalInfo = append(additionalInfo[:i], additionalInfo[i+1:]...)
//						j--
//					} else {
//						return errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
//					}
//				}
//			}
//		}
//	}
//	return errorModel.GenerateNonErrorModel()
//}
//
//func ValidateFloat(fileName string, funcName string, input string, fieldName string, _ int, _ int) (err errorModel.ErrorModel) {
//	if !util.IsStringEmpty(input) {
//		_, errS := strconv.ParseFloat(input, 64)
//		if errS != nil {
//			return errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
//		}
//	}
//
//	return errorModel.GenerateNonErrorModel()
//}
//
//func ValidateFax(fileName string, funcName string, input string, fieldName string, minLength int, _ int) (err errorModel.ErrorModel) {
//	if !util.IsStringEmpty(input) || minLength != 0 {
//		faxSplit := strings.Split(input, "-")
//		if len(faxSplit) != 2 {
//			return errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
//		}
//
//		if !util.IsCountryCode(faxSplit[0]) {
//			return errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
//		}
//
//		if !util.IsFacsimileValid(faxSplit[1]) {
//			return errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
//		}
//	}
//
//	return errorModel.GenerateNonErrorModel()
//}

//func ValidateBoolean(_ string, _ string, input string, fieldName string, _ int, _ int) (err errorModel.ErrorModel) {
//	if !util.IsStringEmpty(input) {
//		err = in.ValidateBooleanValue(input, fieldName)
//		if err.Error != nil {
//			return
//		}
//	}
//
//	return errorModel.GenerateNonErrorModel()
//}

//func ValidateDateWithCurrentDate(fileName string, funcName string, input string, fieldName string, minLength int, _ int) (err errorModel.ErrorModel) {
//	var timeDate time.Time
//	var validationResult bool
//
//	if !util.IsStringEmpty(input) || minLength != 0 {
//		timeDate, err = in.ValidateShortDate(fileName, funcName, fieldName, input)
//		if err.Error != nil {
//			return
//		}
//
//		validationResult, err = in.ValidateDate(fileName, funcName, fieldName, timeDate)
//		if !validationResult {
//			return
//		}
//	}
//
//	return
//}

//func ValidateStringMustEmpty(fileName string, funcName string, input string, fieldName string, _ int, _ int) (err errorModel.ErrorModel) {
//	if !util.IsStringEmpty(input) {
//		return errorModel.GenerateFieldMustEmpty(fileName, funcName, fieldName)
//	}
//	return errorModel.GenerateNonErrorModel()
//}
