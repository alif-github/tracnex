package ProductLicenseService

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/cryptoModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"strings"
)

type productLicenseService struct {
	service.AbstractService
	service.GetListData
}

var ProductLicenseService = productLicenseService{}.New()

func (input productLicenseService) New() (output productLicenseService) {
	output.FileName = "ProductLicenseService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"customer_name",
		"unique_id_1",
		"unique_id_2",
		"installation_id",
		"product_name",
		"license_variant_name",
		"license_type_name",
		"product_valid_from",
		"product_valid_thru",
		"license_status",
	}
	output.ValidSearchBy = []string{
		"customer_name",
		"customer_id",
	}
	output.ServiceName = constanta.ProductLicense

	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.CustomerCategoryDataScope] = applicationModel.MappingScopeDB{
		View:  "c.customer_category_id",
		Count: "c.customer_category_id",
	}
	output.MappingScopeDB[constanta.CustomerGroupDataScope] = applicationModel.MappingScopeDB{
		View:  "c.customer_group_id",
		Count: "c.customer_group_id",
	}
	output.MappingScopeDB[constanta.SalesmanDataScope] = applicationModel.MappingScopeDB{
		View:  "c.salesman_id",
		Count: "c.salesman_id",
	}
	output.MappingScopeDB[constanta.ProvinceDataScope] = applicationModel.MappingScopeDB{
		View:  "c.province_id",
		Count: "c.province_id",
	}
	output.MappingScopeDB[constanta.DistrictDataScope] = applicationModel.MappingScopeDB{
		View:  "c.district_id",
		Count: "c.district_id",
	}
	output.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "pr.client_type_id",
		Count: "pr.client_type_id",
	}

	output.ListScope = output.SetListScope()

	return
}

func (input productLicenseService) checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_product_license_license_config_id") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.LicenseConfig)
		}
	}

	return err
}

func (input productLicenseService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ProductLicenseRequest) errorModel.ErrorModel) (inputStruct in.ProductLicenseRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	id, _ := strconv.Atoi(mux.Vars(request)["id"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func (input productLicenseService) validateDataScope(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	output = service.ValidateScope(contextModel, input.ListScope)
	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, "validateDataScope")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) readBodyAndValidateHWID(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.UpdateLicenseHWIDRequest) errorModel.ErrorModel) (inputStruct in.UpdateLicenseHWIDRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidateHWID"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	err = validation(&inputStruct)
	return
}

func (input productLicenseService) decryptLicense(decryptStruct cryptoModel.DecryptLicenseRequestModel, signatureKey string, productEncrypt string) (jsonDecrypt cryptoModel.DecryptLicenseResponseModel, err errorModel.ErrorModel) {
	var (
		funcName                = "decryptLicense"
		argumentMap             = make(map[string]string)
		tempResult              cryptoModel.DecryptLicenseResponseModel
		dataByteDecryptLicense  []byte
		dataByteResponseDecrypt []byte
		errorS                  error
	)

	//--- Create Model
	decryptStruct.EncryptKey = util2.HashingPassword(decryptStruct.ClientId, decryptStruct.ClientSecret)
	decryptStruct.SignatureKey = signatureKey

	//--- Start Decrypt License
	dataByteDecryptLicense, errorS = json.Marshal(decryptStruct)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	//--- Set Argument Generator License
	argumentMap["args1"] = productEncrypt
	argumentMap["args2"] = string(dataByteDecryptLicense)

	dataByteResponseDecrypt, err = service.GeneratorLicense(constanta.ProductDecryptAction, argumentMap)
	if err.Error != nil {
		if err.CausedBy != nil {
			errorS = errors.New(err.CausedBy.Error())
		}
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	//--- Convert Response Generator
	errorS = json.Unmarshal(dataByteResponseDecrypt, &tempResult)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	if strings.ToLower(tempResult.MessageCode) != constanta.StatusMessage {
		err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, tempResult.Notification, errors.New(tempResult.Message))
		return
	}

	jsonDecrypt = tempResult
	err = errorModel.GenerateNonErrorModel()
	return
}
