package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"strings"
	"time"
)

type GetListDataDTO struct {
	AbstractDTO
	ID                   int64  `json:"id"`
	Filter               string `json:"filter"`
	Search               string `json:"search"`
	UpdatedAtStartString string
	UpdatedAtEndString   string
	UpdatedAtStart       time.Time
	UpdatedAtEnd         time.Time
}

type SearchByParam struct {
	SearchKey      string
	DataType       string
	SearchOperator string
	SearchValue    string
	SearchType     string
}

func (input *GetListDataDTO) ValidateGetListData(validSearchKey []string, validOrderBy []string, validOperator map[string]applicationModel.DefaultOperator, validLimit []int) (searchBy []SearchByParam, err errorModel.ErrorModel) {
	err = input.ValidateInputPageLimitAndOrderBy(validLimit, validOrderBy)
	if err.Error != nil {
		return
	}

	return input.validateFilter(validSearchKey, validOperator)
}

func (input *GetListDataDTO) ValidateDeleteListData(validSearchKey []string, validOperator map[string]applicationModel.DefaultOperator, validLimit []int) (searchBy []SearchByParam, err errorModel.ErrorModel) {
	return input.validateFilter(validSearchKey, validOperator)
}

func (input *GetListDataDTO) ValidateGetCountData(validSearchKey []string, validOperator map[string]applicationModel.DefaultOperator) (searchBy []SearchByParam, err errorModel.ErrorModel) {
	return input.validateFilter(validSearchKey, validOperator)
}

func (input *GetListDataDTO) ValidateGetListDataWithID(id int64, validSearchKey []string, validOrderBy []string, validOperator map[string]applicationModel.DefaultOperator, validLimit []int) (searchBy []SearchByParam, err errorModel.ErrorModel) {
	if id < 1 {
		err = errorModel.GenerateEmptyFieldError("GetListDataDTO.go", "ValidateGetListDataWithID", constanta.ID)
		return
	}

	return input.ValidateGetListData(validSearchKey, validOrderBy, validOperator, validLimit)
}

func (input GetListDataDTO) validateFilter(validSearchKey []string, validOperator map[string]applicationModel.DefaultOperator) (searchBy []SearchByParam, err errorModel.ErrorModel) {
	filter := input.Filter
	search := input.Search
	if filter != "" {
		filterSplitComma := strings.Split(filter, ", ")
		for i := 0; i < len(filterSplitComma); i++ {
			filterIndex := filterSplitComma[i]
			filterIndexSplitSpace := strings.Split(filterIndex, " ")
			if len(filterIndexSplitSpace) > 2 {
				searchKey := strings.Trim(filterIndexSplitSpace[0], " ")
				operator := strings.Trim(filterIndexSplitSpace[1], " ")
				searchValue := ""
				for j := 2; j < len(filterIndexSplitSpace); j++ {
					searchValue += filterIndexSplitSpace[j] + " "
				}

				searchValue = strings.Trim(searchValue, " ")

				if !isOperatorValid(searchKey, searchValue, operator, validOperator) {
					err = errorModel.GenerateFormatFieldError("GetListDataDTO.go", "validateFilter", constanta.Filter)
					return
				}

				err = input.validateSearchValue("GetListDataDTO.go", "validateFilter", constanta.Filter, searchValue)
				if err.Error != nil {
					return
				}

				searchBy = append(searchBy, SearchByParam{
					DataType:       validOperator[searchKey].DataType,
					SearchKey:      searchKey,
					SearchOperator: operator,
					SearchValue:    searchValue,
					SearchType:     constanta.Filter,
				})

				if !common.ValidateStringContainInStringArray(validSearchKey, searchKey) {
					err = errorModel.GenerateFormatFieldError("GetListDataDTO.go", "validateFilter", constanta.Filter)
					return
				}
			} else {
				err = errorModel.GenerateFormatFieldError("GetListDataDTO.go", "validateFilter", constanta.Filter)
				return
			}
		}
	}
	if search != "" {
		searchSplitComma := strings.Split(search, ", ")
		for i := 0; i < len(searchSplitComma); i++ {
			searchIndex := searchSplitComma[i]
			searchIndexSplitSpace := strings.Split(searchIndex, " ")
			if len(searchIndexSplitSpace) >= 3 {
				searchKey := strings.Trim(searchIndexSplitSpace[0], " ")
				operator := strings.Trim(searchIndexSplitSpace[1], " ")
				searchValue := strings.Trim(searchIndexSplitSpace[2], " ")

				err = input.validateSearchValue("GetListDataDTO.go", "validateFilter", constanta.Filter, searchValue)
				if err.Error != nil {
					return
				}
				searchBy = append(searchBy, SearchByParam{
					DataType:       validOperator[searchKey].DataType,
					SearchKey:      searchKey,
					SearchOperator: operator,
					SearchValue:    searchValue,
					SearchType:     constanta.Search,
				})

				if !common.ValidateStringContainInStringArray(validSearchKey, searchKey) {
					err = errorModel.GenerateFormatFieldError("GetListDataDTO.go", "validateFilter", constanta.Search)
					return
				}
				if !isOperatorValid(searchKey, searchValue, operator, validOperator) {
					err = errorModel.GenerateFormatFieldError("GetListDataDTO.go", "validateFilter", constanta.Search)
					return
				}
			} else {
				err = errorModel.GenerateFormatFieldError("GetListDataDTO.go", "validateFilter", constanta.Search)
				return
			}
		}
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func validateOrderBy(orderBy string, validOrderBy []string) bool {
	return common.ValidateStringContainInStringArray(validOrderBy, orderBy)
}

func (input GetListDataDTO) validateSearchValue(fileName, funcName, fieldName, searchValue string) (err errorModel.ErrorModel) {
	if util.IsStringEmpty(searchValue) {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
		return
	}

	err = util2.ValidateSpecialCharacter(fileName, funcName, fieldName, searchValue)
	if err.Error != nil {
		return
	}

	return
}

func isOperatorValid(key string, value string, operator string, validOperator map[string]applicationModel.DefaultOperator) bool {
	if validOperator[key].Operator == nil {
		return false
	} else {
		if validOperator[key].DataType == "number" {
			_, err := strconv.Atoi(value)
			if err != nil {
				return false
			}
		}
		return common.ValidateStringContainInStringArray(validOperator[key].Operator, operator)
	}
}

func checkLimit(validLimit []int, limit int) (result int) {
	if len(validLimit) < 0 {
		return 0
	}
	if limit < validLimit[0] {
		return validLimit[0]
	}
	oldValidLimit := validLimit[0]
	for i := 0; i < len(validLimit); i++ {
		if validLimit[i] == limit {
			return limit
		} else {
			if i > 0 {
				if oldValidLimit < limit && validLimit[i] > limit {
					if limit-oldValidLimit > validLimit[i]-limit {
						return validLimit[i]
					} else {
						return oldValidLimit
					}
				}
			}
		}
		oldValidLimit = validLimit[i]
		if i == len(validLimit)-1 {
			return validLimit[i]
		}
	}
	return 0
}

func (input *GetListDataDTO) ValidateGetCountDataWithID(id int64, validSearchKey []string, validOperator map[string]applicationModel.DefaultOperator) (searchBy []SearchByParam, err errorModel.ErrorModel) {
	if id < 1 {
		err = errorModel.GenerateEmptyFieldError("GetListDataDTO.go", "ValidateGetCountDataWithID", constanta.ID)
		return
	}
	return input.ValidateGetCountData(validSearchKey, validOperator)

}

func (input *GetListDataDTO) ValidateUpdatedAtRange() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdatedAtRange"
	var errS error

	if !util.IsStringEmpty(input.UpdatedAtStartString) {
		input.UpdatedAtStart, errS = time.Parse(constanta.DefaultTimeFormat, input.UpdatedAtStartString)
		if errS != nil {
			return errorModel.GenerateFormatFieldError(GetListDataDTOFileName, funcName, constanta.UpdatedAtStart)
		}
	} else if input.Page == -99 || input.Limit == -99 {
		err = errorModel.GenerateEmptyFieldError(GetListDataDTOFileName, funcName, constanta.UpdatedAtStart)
		return
	}

	if !util.IsStringEmpty(input.UpdatedAtEndString) {
		input.UpdatedAtEnd, errS = time.Parse(constanta.DefaultTimeFormat, input.UpdatedAtEndString)
		if errS != nil {
			return errorModel.GenerateFormatFieldError(GetListDataDTOFileName, funcName, constanta.UpdatedAtEnd)
		}
	}

	if !util.IsStringEmpty(input.UpdatedAtStartString) && util.IsStringEmpty(input.UpdatedAtEndString) {
		input.UpdatedAtEnd, _ = time.Parse(constanta.DefaultTimeFormat, time.Now().Format(constanta.DefaultTimeFormat))
	}

	if !util.IsStringEmpty(input.UpdatedAtEndString) && util.IsStringEmpty(input.UpdatedAtStartString) {
		return errorModel.GenerateEmptyFieldError(GetListDataDTOFileName, funcName, constanta.UpdatedAtStart)
	}

	if input.UpdatedAtEnd.Before(input.UpdatedAtStart) {
		err = errorModel.GenerateFieldFormatWithRuleError(GetListDataDTOFileName, funcName, constanta.NeedMoreThan, constanta.UpdatedAtEnd, "")
	}

	return errorModel.GenerateNonErrorModel()
}
