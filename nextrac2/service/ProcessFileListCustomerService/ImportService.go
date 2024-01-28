package ProcessFileListCustomerService

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/CustomerListService"
)

type importService struct {
	service.AbstractService
	service.ImportDataService
	AvailableTypeImport map[string]service.ImportDataType
}

var ImportService = importService{}.New()

func (input importService) New() (output importService) {
	output.FileName = "ImportService.go"
	output.RowStartAt = 1
	output.AvailableTypeImport = make(map[string]service.ImportDataType)
	output.AvailableTypeImport["customer"] = service.ImportDataType{
		Delimiter:       constanta.PipaDelimiter,
		IsActive:        true,
		Message:         input.messageImportFile(),
		ImportValidator: CustomerListService.CustomerListService.GetImportValidator(),
		ConvertToDTO:    CustomerListService.CustomerListService.ConvertImportDataToDTO,
		DoInsert:        CustomerListService.CustomerListService.DoInsertCustomerByImport,
		Truncate:        CustomerListService.CustomerListService.TruncateTableCustomerList,
	}
	return
}

func (input *importService) GetImportValidator(inputStruct in.ImportRequest) (output []service.ImportValidatorData, err errorModel.ErrorModel) {
	funcName := "GetImportValidator"

	if input.AvailableTypeImport[inputStruct.TypeData].ImportValidator != nil {
		output = input.AvailableTypeImport[inputStruct.TypeData].ImportValidator
	} else {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "Type")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input importService) messageImportFile() string {
	return "COMPANY_ID\t:\tMandatory (Maximum 20 Character)\n" +
		"BRANCH_ID\t:\tMandatory (Maximum 20 Character)\n" +
		"COMPANY_NAME\t:\tMandatory (Maximum 50 Character)\n" +
		"CITY\t:\tOptional (Maximum 50 Character)\n" +
		"IMPLEMENTER\t:\tOptional (Maximum 50 Character)\n" +
		"IMPLEMENTATION\t:\tOptional (YYYY-MM-DD)\n" +
		"PRODUCT\t:\tMandatory (Maximum 50 Character)\n" +
		"VERSION\t:\tOptional (Maximum 50 Character)\n" +
		"LICENSE_TYPE\t:\tOptional (Maximum 50 Character)\n" +
		"USER_AMOUNT\t:\tOptional (Number Only)\n" +
		"EXP_DATE\t:\tMandatory (YYYY-MM-DD)\n"
}
