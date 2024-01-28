package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"strings"
)

type FileInfo struct {
	FileName	string	`json:"file_name"`
	SizeFile	int64	`json:"size_file"`
	TypeFile	string	`json:"type_file"`
}

type ImportRequest struct {
	TypeData	string	`json:"type_data"`
}

type ImportConfirmRequest struct {
	TypeData	string	`json:"type_data"`
	Filename	string	`json:"filename"`
	Truncate	bool	`json:"truncate"`
	Confirm		bool	`json:"confirm"`
}

func (input *ImportRequest) ValidateImport() (err errorModel.ErrorModel) {
	fileName := "ProcessFileCustomerListDTO.go"
	funcName := "ValidateImport"

	if util.IsStringEmpty(input.TypeData) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.GetTypeImportFile)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *ImportConfirmRequest) ValidateConfirm() (err errorModel.ErrorModel) {
	fileName := "ProcessFileCustomerListDTO.go"
	funcName := "ValidateConfirm"

	if util.IsStringEmpty(input.TypeData) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Type)
	}

	if util.IsStringEmpty(input.Filename) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Filename)
	}

	splitFilename := strings.Split(input.Filename, "_")
	if len(splitFilename) != 2 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Filename)
	}

	if input.TypeData != splitFilename[0] {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Type)
	}

	return errorModel.GenerateNonErrorModel()
}