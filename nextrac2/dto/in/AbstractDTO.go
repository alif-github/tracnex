package in

import (
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	common2 "nexsoft.co.id/nextrac2/resource_common_service/common"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type AbstractDTO struct {
	Page           int    `json:"page"`
	Limit          int    `json:"limit"`
	OrderBy        string `json:"order_by"`
	UpdatedAtStart string `json:"updated_at_start"`
	UpdatedAtEnd   string `json:"updated_at_end"`
}

func TimeStrToTime(timeStr string, fieldName string) (output time.Time, error errorModel.ErrorModel) {
	validationResult := util.IsStringEmpty(timeStr)
	if validationResult {
		error = errorModel.GenerateEmptyFieldError("AbstractDTO.go", "TimeStrToTime", constanta.UpdatedAt)
		return
	}

	output, err := time.Parse(constanta.DefaultTimeFormat, timeStr)
	if err != nil {
		error = errorModel.GenerateFormatFieldError("AbstractDTO.go", "TimeStrToTime", fieldName)
		return
	}

	return output, errorModel.GenerateNonErrorModel()
}

func TimeStrToTimeWithTimeFormat(timeStr string, fieldName string, format string) (output time.Time, error errorModel.ErrorModel) {
	output, err := time.Parse(format, timeStr)
	if err != nil {
		error = errorModel.GenerateFormatFieldError("AbstractDTO.go", "TimeStrToTime", fieldName)
		return
	}

	return output, errorModel.GenerateNonErrorModel()
}

func TimeDBStrToTime(timeStr string, fieldName string) (output time.Time, error errorModel.ErrorModel) {
	output, err := time.Parse(constanta.DefaultDBTimeFormat, timeStr)
	if err != nil {
		error = errorModel.GenerateFormatFieldError("AbstractDTO.go", "TimeStrToTime", fieldName)
		return
	}

	return output, errorModel.GenerateNonErrorModel()
}

func (input AbstractDTO) ValidateMinMaxString(inputStr string, fieldName string, min int, max int) errorModel.ErrorModel {
	lenStr := len(inputStr)
	if lenStr < min {
		if min == 1 {
			return errorModel.GenerateEmptyFieldError("AbstractDTO.go", "ValidateMinMaxString", fieldName)
		} else {
			return errorModel.GenerateFieldFormatWithRuleError("AbstractDTO.go", "ValidateMinMaxString", "NEED_MORE_THAN", fieldName, strconv.Itoa(min))
		}
	}
	if lenStr > max {
		return errorModel.GenerateFieldFormatWithRuleError("AbstractDTO.go", "ValidateMinMaxString", "NEED_LESS_THAN", fieldName, strconv.Itoa(max))
	}

	return errorModel.GenerateNonErrorModel()
}

func (input AbstractDTO) ValidateStatus(status string) errorModel.ErrorModel {
	if status != "A" && status != "N" {
		return errorModel.GenerateUnknownDataError("AbstractDTO.go", "ValidateStatus", constanta.Status)
	}
	return errorModel.GenerateNonErrorModel()
}

func (input AbstractDTO) ValidateBooleanValue(userInput string, fieldName string) errorModel.ErrorModel {
	if userInput != "Y" && userInput != "N" {
		return errorModel.GenerateUnknownDataError("AbstractDTO.go", "ValidateBooleanValue", fieldName)
	}
	return errorModel.GenerateNonErrorModel()
}

func (input *AbstractDTO) ValidateInputPageLimitAndOrderBy(validLimit []int, validOrderBy []string) (err errorModel.ErrorModel) {
	funcName := "ValidateInputPageAndLimit"
	if input.Page < 1 && input.Page != -99 {
		return errorModel.GenerateFieldFormatWithRuleError(AbstractDTOFileName, funcName, constanta.NeedMoreThan, constanta.Page, "0")
	}

	if input.Limit != -99 {
		input.Limit = checkLimit(validLimit, input.Limit)
	}

	if input.Limit < 1 && input.Limit != -99 && util.IsStringEmpty(input.UpdatedAtStart) {
		return errorModel.GenerateFieldFormatWithRuleError(AbstractDTOFileName, funcName, constanta.NeedMoreThan, constanta.Limit, "0")
	}

	// Validate Order
	arrOrder := strings.Split(input.OrderBy, ", ")
	var resultOrder string
	for i := 0; i < len(arrOrder); i++ {
		orderItem := strings.Trim(arrOrder[i], " ")
		if orderItem == "" {
			if len(arrOrder) != 1 && i > 0 {
				continue
			}
			resultOrder = validOrderBy[0]
		} else {
			orderItem, err = ValidateOrderBy(validOrderBy, orderItem)
			if err.Error != nil {
				return
			}
			resultOrder += orderItem
			if i != len(arrOrder)-1 {
				resultOrder += ", "
			}
		}
	}
	input.OrderBy = resultOrder

	return errorModel.GenerateNonErrorModel()
}

func (input *AbstractDTO) ValidateInputPageLimitAndOrderByDistributor(validLimit []int, validOrderBy []string) (err errorModel.ErrorModel) {
	funcName := "ValidateInputPageLimitAndOrderByDistributor"
	if input.Page < 1 {
		return errorModel.GenerateFieldFormatWithRuleError(AbstractDTOFileName, funcName, constanta.NeedMoreThan, constanta.Page, "0")
	}

	if input.Limit < 1 {
		return errorModel.GenerateFieldFormatWithRuleError(AbstractDTOFileName, funcName, constanta.NeedMoreThan, constanta.Limit, "0")
	}

	// Validate Order
	arrOrder := strings.Split(input.OrderBy, ", ")
	var resultOrder string
	for i := 0; i < len(arrOrder); i++ {
		orderItem := strings.Trim(arrOrder[i], " ")
		if orderItem == "" {
			if len(arrOrder) != 1 && i > 0 {
				continue
			}
			resultOrder = validOrderBy[0]
		} else {
			orderItem, err = ValidateOrderBy(validOrderBy, orderItem)
			if err.Error != nil {
				return
			}
			resultOrder += orderItem
			if i != len(arrOrder)-1 {
				resultOrder += ", "
			}
		}
	}
	input.OrderBy = resultOrder

	return errorModel.GenerateNonErrorModel()
}

func (input AbstractDTO) ValidateUOM(userInput string, fieldName string) errorModel.ErrorModel {
	err := input.ValidateMinMaxString(userInput, fieldName, 1, 4)
	if err.Error != nil {
		return err
	}
	return errorModel.GenerateNonErrorModel()
}

func (input AbstractDTO) ValidateReportColumnHeader(userInput string, fieldName string) errorModel.ErrorModel {
	if !util.IsStringEmpty(userInput) {
		err := input.ValidateMinMaxString(userInput, fieldName, 1, 10)
		if err.Error != nil {
			return err
		}
	}
	return errorModel.GenerateNonErrorModel()
}

func (input AbstractDTO) ValidateIsContainSpaceString(fileName string, funcName string, fieldName string, userInput string) errorModel.ErrorModel {
	result, errorS := regexp.Compile("^[a-zA-Z0-9]+$")

	if errorS != nil {
		return errorModel.GenerateUnknownError(fileName, funcName, errorS)
	}
	if !result.MatchString(userInput) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input AbstractDTO) ValidateSpecialCharacter(fileName string, funcName string, fieldName string, userInput string) errorModel.ErrorModel {
	result, errorS := regexp.Compile("^[a-zA-Z0-9 ]+[a-zA-Z0-9-._! ]*$")

	if errorS != nil {
		return errorModel.GenerateUnknownError(fileName, funcName, errorS)
	}
	if !result.MatchString(userInput) {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.ErrorSpecialCharacter, fieldName, "")
	}

	return errorModel.GenerateNonErrorModel()
}

func IsOnlyWordCharacterValid(inputText string) (bool, string, string) {
	onlyCharacterRegex := regexp.MustCompile("^[\\w]*$")
	return onlyCharacterRegex.MatchString(inputText), "IS_ONLY_WORD_CHARACTER_REGEX_MESSAGE", ""
}

func IsOnlyDigitValid(inputText string) (bool, string, string) {
	onlyCharacterRegex := regexp.MustCompile("^[\\d]*$")
	return onlyCharacterRegex.MatchString(inputText), "IS_ONLY_DIGIT_REGEX_MESSAGE", ""
}

func IsOnlyAlfaNumerikWithOrWithoutUnderScoredValid(inputText string) (bool, string, string) {
	onlyCharacterRegex := regexp.MustCompile("^[a-zA-Z0-9]+[a-zA-Z0-9_ ]{0,}$")
	return onlyCharacterRegex.MatchString(inputText), "IS_ONLY_ALFANUMERIK_REGEX_MESSAGE", ""
}

func IsOnlyAlfaNumerikValid(inputText string) (bool, string, string) {
	onlyCharacterRegex := regexp.MustCompile("^[a-zA-Z0-9]+[a-zA-Z0-9 ]{0,}$")
	return onlyCharacterRegex.MatchString(inputText), "IS_ONLY_ALFANUMERIK_REGEX_MESSAGE2", ""
}

func (input AbstractDTO) UniqueStrArray(fileName string, funcName string, stringSlice []string) errorModel.ErrorModel {
	keys := make(map[string]bool)
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
		} else {
			return errorModel.GenerateDataDuplicateInDTOError(fileName, funcName)
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input AbstractDTO) IsPhoneNumberValid(phone string) (number int, isValid bool, regex string) {
	number, isValid = util.IsNumeric(phone)
	return number, len(phone) <= 13 && isValid, "PHONE_FORMAT_REGEX_MESSAGE"
}

type regexValue struct {
	regex    string
	ruleName string
}

type defaultField struct {
	enum       map[string][]string
	dateFormat map[string]string
	regex      map[string]regexValue
	autoFix    map[string]func(value reflect.Value)
}

var (
	defaultFields = initiate()
	intType       = []string{reflect.Int64.String(), reflect.Int.String(), reflect.Int32.String()}
	floatType     = []string{reflect.Float32.String(), reflect.Float64.String()}
)

func initiate() defaultField {
	enum := make(map[string][]string)
	enum["record_status"] = []string{"A", "N", "P"}
	enum["sharing_permission"] = []string{"edit", "view"}

	dateFormat := make(map[string]string)
	dateFormat["default"] = constanta.DefaultTimeFormat
	//dateFormat["date_only"] = constanta.DateOnlyTimeFormat

	regex := make(map[string]regexValue)
	regex["profile_name"] = regexValue{
		regex:    util.ProfileNameStandardRegex,
		ruleName: "PROFILE_NAME_REGEX_MESSAGE",
	}
	regex["directory_name"] = regexValue{
		regex:    util.DirectoryNameStandardRegex,
		ruleName: "DIRECTORY_NAME_REGEX_MESSAGE",
	}

	regex["code"] = regexValue{
		regex:    util.AlphanumericRegex,
		ruleName: "ALPHANUMERIC_REGEX",
	}
	regex["parameter_permission"] = regexValue{
		regex:    "^(nexsoft[.])(([a-z]+)([_\\.])([a-z]+))+$",
		ruleName: "NAME_PARAMETER_FORMAT",
	}
	regex["parameter_name"] = regexValue{
		regex:    "^([a-z_]+[a-z])$",
		ruleName: "PERMISSION_PARAMETER_FORMAT",
	}

	//regex["name"] = regexValue{
	//	regex:    util.NameStandardRegex,
	//	ruleName: "NAME_REGEX_MESSAGE",
	//}

	regex["description"] = regexValue{
		regex:    "^[A-Za-z ]+$",
		ruleName: "DESCRIPTION_REGEX_MESSAGE",
	}

	autoFix := make(map[string]func(value reflect.Value))
	autoFix["filename"] = autoFixFilename

	return defaultField{
		enum:       enum,
		dateFormat: dateFormat,
		regex:      regex,
		autoFix:    autoFix,
	}
}

func (input AbstractDTO) basicValidatorByTag(dto interface{}, menu string) (err errorModel.ErrorModel) {
	funcName := "basicValidatorByTag"
	fileName := "AbstractDTO.go"

	reflectType := reflect.TypeOf(dto).Elem()
	reflectValue := reflect.ValueOf(dto).Elem()

	max := 0
	min := 0
	isMinFound := false
	isMaxFound := false
	for i := 0; i < reflectType.NumField(); i++ {
		currentField := reflectType.Field(i)
		currentValue := reflectValue.FieldByName(currentField.Name)

		if currentField.Name == "AbstractDTO" {
			continue
		}

		if currentField.Type.Kind() == reflect.Struct {
			newDTO := currentValue.Addr().Interface()
			err = input.basicValidatorByTag(newDTO, menu)
			if err.Error != nil {
				return
			}
		}

		required := currentField.Tag.Get("required")
		requiredArray := strings.Split(required, ",")
		reservedValues := currentField.Tag.Get("reserved")
		if reservedValues != "" {
			reservedValue := strings.Split(reservedValues, ",")
			if common2.ValidateStringContainInStringArray(reservedValue, currentValue.String()) {
				err = errorModel.GenerateReservedValueString(fileName, funcName, currentField.Name)
				return
			}
		}
		if common.ValidateStringContainInStringArray(requiredArray, menu) {
			defaultValue := currentField.Tag.Get("default")
			min, isMinFound, max, isMaxFound = getMinMaxValue(currentField)
			if common.ValidateStringContainInStringArray(intType, currentField.Type.String()) {
				if currentValue.IsZero() {
					valueIn, _ := strconv.Atoi(defaultValue)
					currentValue.SetInt(int64(valueIn))
				}

				value := currentValue.Int()
				if isMinFound {
					if min != 0 && int(value) == 0 {
						err = errorModel.GenerateEmptyFieldError(fileName, funcName, currentField.Name)
						return
					}
					if int(value) < min {
						err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_MORE_THAN", currentField.Name, strconv.Itoa(min))
						return
					}
				}
				if isMaxFound {
					if int(value) > max {
						err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_LESS_THAN", currentField.Name, strconv.Itoa(max))
						return
					}
				}
			} else if common.ValidateStringContainInStringArray(floatType, currentField.Type.String()) {
				if currentValue.IsZero() {
					valueIn, _ := strconv.ParseFloat(defaultValue, 64)
					currentValue.SetFloat(valueIn)
				}
				value := currentValue.Float()
				if isMinFound {
					if value < float64(min) {
						err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_MORE_THAN", currentField.Name, strconv.Itoa(min))
						return
					}
				}
				if isMaxFound {
					if value > float64(max) {
						err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_LESS_THAN", currentField.Name, strconv.Itoa(max))
						return
					}
				}
			} else if reflect.String.String() == currentField.Type.String() {
				currentValue.SetString(strings.Trim(currentValue.String(), " "))
				if currentValue.IsZero() {
					currentValue.SetString(defaultValue)
				}

				value := currentValue.String()
				err = input.ValidateMinMaxString(value, currentField.Name, min, max)
				if err.Error != nil {
					return
				}

				enumField := currentField.Tag.Get("enum")
				if enumField != "" {
					if !common.ValidateStringContainInStringArray(defaultFields.enum[enumField], currentValue.String()) {
						err = errorModel.GenerateUnknownDataError(fileName, funcName, currentField.Name)
						return
					}
				}

				autoFix := currentField.Tag.Get("auto_fix")
				if autoFix != "" {
					if defaultFields.autoFix[autoFix] != nil {
						defaultFields.autoFix[autoFix](currentValue)
					}
				}

				regexField := currentField.Tag.Get("regex")
				if regexField != "" {
					if defaultFields.regex[regexField].regex != "" {
						if len(value) > 0 || min != 0 {
							if !regexp.MustCompile(defaultFields.regex[regexField].regex).MatchString(currentValue.String()) {
								err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, defaultFields.regex[regexField].ruleName, currentField.Name, "")
								return
							}
						}

					}
				}

			} else if "time.Time" == currentField.Type.String() {
				var timeObject time.Time
				dateFormatTag := currentField.Tag.Get("dateFormat")
				strField := currentField.Name + "Str"

				timeFormatUsed := defaultFields.dateFormat["default"]
				if defaultFields.dateFormat[dateFormatTag] != "" {
					timeFormatUsed = defaultFields.dateFormat[dateFormatTag]
				}

				timeObject, err = TimeStrToTimeWithTimeFormat(reflectValue.FieldByName(strField).String(), currentField.Name, timeFormatUsed)
				if err.Error != nil {
					return
				}
				currentValue.Set(reflect.ValueOf(timeObject))
			} else if currentValue.Kind() == reflect.Slice {
				if isMinFound {
					if currentValue.Len() == 0 {
						err = errorModel.GenerateEmptyFieldError(fileName, funcName, currentField.Name)
						return
					}
					if currentValue.Len() < min {
						err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_MORE_THAN", currentField.Name, strconv.Itoa(min))
						return
					}
				}
				if isMaxFound {
					if currentValue.Len() > max {
						err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_LESS_THAN", currentField.Name, strconv.Itoa(max))
						return
					}
				}
				for i := 0; i < currentValue.Len(); i++ {
					temp := currentValue.Index(i)
					if temp.Type().String() == reflect.Struct.String() || temp.Type().Kind().String() == reflect.Struct.String() {
						newDTO := currentValue.Index(i).Addr().Interface()
						err = input.basicValidatorByTag(newDTO, menu)
						if err.Error != nil {
							return
						}
					}
				}
			}
		}
	}
	return
}

func getMinMaxValue(field reflect.StructField) (min int, isMinFound bool, max int, isMaxFound bool) {
	maxStr, isMaxFound := field.Tag.Lookup("max")
	minStr, isMinFound := field.Tag.Lookup("min")

	min, _ = strconv.Atoi(minStr)
	max, _ = strconv.Atoi(maxStr)

	return
}

func autoFixFilename(input reflect.Value) {
	if input.Type().Kind() != reflect.String {
		return
	}

	temp := input.String()

	temp = strings.ReplaceAll(temp, "/", "_")
	temp = strings.ReplaceAll(temp, "\\", "_")

	input.SetString(temp)
}

func IsEmailAddressMDB(input string) (output bool) {
	str := fmt.Sprintf(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`)
	emailRegexp := regexp.MustCompile(str)
	return emailRegexp.MatchString(input)
}

func IsPhoneNumberWithCountryCodeMDB(input string) bool {
	phoneNumberRegexp := regexp.MustCompile("[+][0-9]{1,3}[-][1-9][0-9]{8,12}$")
	return phoneNumberRegexp.MatchString(input)
}

func ValidateOrderBy(validOrderBy []string, orderBy string) (result string, err errorModel.ErrorModel) {
	orderBySplit := strings.Split(orderBy, " ")
	funcName := "ValidateOrderBy"
	var isAscending bool

	if !(len(orderBySplit) >= 1 && len(orderBySplit) <= 2) {
		err = errorModel.GenerateFormatFieldError("GetListDataDTO.go", funcName, constanta.OrderBy)
		return
	}

	if len(orderBySplit) == 1 {
		isAscending = true
	} else {
		if strings.ToUpper(orderBySplit[1]) == "ASC" {
			isAscending = true
		} else if strings.ToUpper(orderBySplit[1]) == "DESC" {
			isAscending = false
		} else {
			err = errorModel.GenerateFormatFieldError("GetListDataDTO.go", funcName, constanta.OrderBy)
			return
		}
	}

	if !validateOrderBy(orderBySplit[0], validOrderBy) {
		err = errorModel.GenerateFormatFieldError("GetListDataDTO.go", funcName, constanta.OrderBy)
		return
	}

	result = orderBySplit[0] + " "
	if isAscending {
		result += "ASC"
	} else {
		result += "DESC"
	}

	return
}

func AutoAddNumberRTRW(input *string, numberMin int) {
	for {
		if len(*input) >= numberMin {
			break
		}

		//--- Add 0
		*input = "0" + *input
	}
}
