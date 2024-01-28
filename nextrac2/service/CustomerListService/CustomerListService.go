package CustomerListService

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"strconv"
)

type customerListService struct {
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var CustomerListService = customerListService{}.New()

func (input customerListService) New() (output customerListService) {
	output.FileName = "CustomerListService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidSearchBy = []string{"company_id", "branch_id", "company_name"}
	output.ValidOrderBy = []string{"id"}
	return
}

func (input customerListService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(customerListRequest *in.CustomerListRequest) errorModel.ErrorModel) (inputStruct in.CustomerListRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])

	for _, branchDataElm := range inputStruct.BranchData {
		if branchDataElm.ID == 0 {
			branchDataElm.ID = int64(id)
		}
	}

	err = validation(&inputStruct)
	return
}

func (input customerListService) readBodyAndValidateForView(request *http.Request, contextModel *applicationModel.ContextModel, validation func(customerListRequest *in.CustomerListRequest) errorModel.ErrorModel) (inputStruct in.CustomerListRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	inputStruct.ID = int64(id)

	err = validation(&inputStruct)
	return
}

func (input customerListService) GetImportValidator() (output []service.ImportValidatorData) {
	output = append(output, service.ImportValidatorData{
		FieldName:  constanta.CompanyID,
		HeaderName: "Cust ID",
		MinLength:  1,
		MaxLength:  20,
		Validator:  util.ValidateStringWithMinMaxMandatory,
	})

	output = append(output, service.ImportValidatorData{
		FieldName:  constanta.BranchID,
		HeaderName: "Branch ID",
		MinLength:  1,
		MaxLength:  20,
		Validator:  util.ValidateStringWithMinMaxMandatory,
	})

	output = append(output, service.ImportValidatorData{
		FieldName:  constanta.CompanyName,
		HeaderName: "Company Name",
		MinLength:  1,
		MaxLength:  50,
		Validator:  util.ValidateStringWithMinMaxMandatory,
	})

	output = append(output, service.ImportValidatorData{
		FieldName:  constanta.City,
		HeaderName: "City",
		MinLength:  1,
		MaxLength:  50,
		Validator:  util.ValidateStringWithMinMaxOptional,
	})

	output = append(output, service.ImportValidatorData{
		FieldName:  constanta.Implementer,
		HeaderName: "Implementer",
		MinLength:  1,
		MaxLength:  50,
		Validator:  util.ValidateStringWithMinMaxOptional,
	})

	output = append(output, service.ImportValidatorData{
		FieldName:  constanta.ImplementationAt,
		HeaderName: "Implementasi",
		Validator:  util.ValidateDateTimeOptional,
	})

	output = append(output, service.ImportValidatorData{
		FieldName:  constanta.Product,
		HeaderName: "Product",
		MinLength:  1,
		MaxLength:  50,
		Validator:  util.ValidateStringWithMinMaxMandatory,
	})

	output = append(output, service.ImportValidatorData{
		FieldName:  constanta.Version,
		HeaderName: "Version",
		MinLength:  1,
		MaxLength:  50,
		Validator:  util.ValidateStringWithMinMaxOptional,
	})

	output = append(output, service.ImportValidatorData{
		FieldName:  constanta.LicenseType,
		HeaderName: "License Type",
		MinLength:  1,
		MaxLength:  50,
		Validator:  util.ValidateStringWithMinMaxOptional,
	})

	output = append(output, service.ImportValidatorData{
		HeaderName: "User",
		Validator:  util.ValidateParseInteger,
	})

	output = append(output, service.ImportValidatorData{
		FieldName:  constanta.ExpDate,
		HeaderName: "Expired",
		Validator:  util.ValidateDateWithFindMandatory,
	})

	return
}

func (input customerListService) ConvertImportDataToDTO(importData []string, _ *map[int]map[string]int64, _ *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var customerStruct in.CustomerListImportRequest

	//-------- Preparing the data
	companyID := importData[0]
	branchID := importData[1]
	companyName := importData[2]
	city := importData[3]
	implementer := importData[4]
	implementation := util.DateConvert(importData[5])
	product := importData[6]
	version := importData[7]
	licenseType := importData[8]
	userAmount, _ := strconv.Atoi(importData[9])

	var indexFillColumn int
	for i := 17; i >= 10; i-- {
		if importData[i] != "" {
			indexFillColumn = i
			break
		}
	}

	expDate := util.DateConvert(importData[indexFillColumn])

	//-------- Copy to struct
	customerStruct = in.CustomerListImportRequest{
		CompanyID:      companyID,
		BranchID:       branchID,
		CompanyName:    companyName,
		City:           city,
		Implementer:    implementer,
		Implementation: implementation,
		Product:        product,
		Version:        version,
		LicenseType:    licenseType,
		UserAmount:     userAmount,
		ExpDate:        expDate,
	}

	result = customerStruct
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListService) TruncateTableCustomerList(tx *sql.Tx, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	funcName := "truncateTableCustomerList"

	//------- Step 1 truncate table customer
	err = dao.CustomerListDAO.TruncateTableCustomer(tx)
	if err.Error != nil {
		return
	}

	//------- Step 2 count row
	var counts int64
	counts, err = dao.CustomerListDAO.CountRowCustomerList(tx)
	if err.Error != nil {
		return
	}

	//------- Step 3 reset sequence ID
	if counts == 0 {
		err = dao.CustomerListDAO.SetValSequenceCustomerListTable(tx)
		if err.Error != nil {
			return
		}
	} else {
		newMessageError := GenerateI18NMessage("COUNT_IS_NOT_ZERO_MESSAGE", contextModel.AuthAccessTokenModel.Locale)
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errors.New(newMessageError))
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
