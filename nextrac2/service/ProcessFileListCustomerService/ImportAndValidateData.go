package ProcessFileListCustomerService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/util"
)

func (input importService) ImportAndValidateData(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	fileName := "ImportAndValidateData.go"
	funcName := "ImportAndValidateData"

	inputStruct, buffer, extension, err := input.ReadRequestMultipartForm(request, contextModel, input.validateImport)
	if err.Error != nil {
		return
	}

	importValidator, err := input.GetImportValidator(inputStruct)
	if err.Error != nil {
		return
	}

	fileDataName, result, totalData, err, multipleError := input.ValidateImportData(fileName, funcName, buffer, extension,
		constanta.PipaDelimiter, importValidator, inputStruct.TypeData, contextModel.AuthAccessTokenModel.Locale)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse {
		Code: 		util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
	}

	if len(multipleError) != 0 {
		err = errorModel.GenerateMultipleErrorAcquired(fileName, funcName)
		output.Data.Content = multipleError
		output.Status.Message = GenerateI18NMessage("FAILED_IMPORT_VALIDATE_MESSAGE", contextModel.AuthAccessTokenModel.Locale)
	} else {
		output.Data.Content = out.ImportDataResponse {
			Filename: 	fileDataName,
			Data: 		result,
			TotalData: 	totalData,
		}
		output.Status.Message = GenerateI18NMessage("SUCCESS_IMPORT_VALIDATE_MESSAGE", contextModel.AuthAccessTokenModel.Locale)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input importService) validateImport(inputStruct *in.ImportRequest) errorModel.ErrorModel {
	return inputStruct.ValidateImport()
}