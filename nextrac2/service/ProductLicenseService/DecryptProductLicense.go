package ProductLicenseService

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
)

func (input productLicenseService) DecryptProductLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.ProductLicenseRequest
	)

	inputStruct, err = input.readBodyAndValidateForViewProductLicense(request, input.validateDecryptProductLicense)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doDecryptProductLicenseByIDProductLicense(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_PRODUCT_LICENSE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) validateDecryptProductLicense(inputStruct *in.ProductLicenseRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewProductLicense()
}

func (input productLicenseService) doDecryptProductLicenseByIDProductLicense(inputStruct in.ProductLicenseRequest, contextModel *applicationModel.ContextModel) (output in.GenerateDataLicenseConfiguration, err errorModel.ErrorModel) {
	var (
		funcName                = "doDecryptProductLicenseByIDProductLicense"
		productLicenseOnDB      repository.ProductLicenseModel
		dataDecryptLicense      in.GenerateDataProductConfiguration
		validResponse           in.GenerateDataValidationResponse
		argumentMap             map[string]string
		dataByteResponseDecrypt []byte
		dataByteDecryptLicense  []byte
		errorS                  error
		scope                   map[string]interface{}
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	productLicenseOnDB, err = dao.ProductLicenseDAO.GetDataForDecryptProductLicense(serverconfig.ServerAttribute.DBConnection, repository.ProductLicenseModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	}, scope, input.MappingScopeDB, input.ListScope)
	if err.Error != nil {
		return
	}

	if productLicenseOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateDataNotFound(input.FileName, funcName)
		return
	}

	// convert to request argument generator exe
	dataDecryptLicense = in.GenerateDataProductConfiguration{
		SignatureKey:     productLicenseOnDB.SignatureKey.String,
		ProductSignature: productLicenseOnDB.ProductSignature.String,
		ClientId:         productLicenseOnDB.ClientId.String,
		ClientSecret:     productLicenseOnDB.ClientSecret.String,
		EncryptKey:       util2.HashingPassword(productLicenseOnDB.ClientId.String, productLicenseOnDB.ClientSecret.String),
		HardwareId:       productLicenseOnDB.HWID.String,
		ProductKey:       productLicenseOnDB.ProductKey.String,
		ProductId:        productLicenseOnDB.ProductId.String,
	}

	dataByteDecryptLicense, errorS = json.Marshal(dataDecryptLicense)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	// set argument generator license
	argumentMap = make(map[string]string)
	argumentMap["args1"] = productLicenseOnDB.ProductEncrypt.String
	argumentMap["args2"] = string(dataByteDecryptLicense)

	dataByteResponseDecrypt, err = service.GeneratorLicense(constanta.ProductDecryptAction, argumentMap)
	if err.Error != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	// convert response generator
	errorS = json.Unmarshal(dataByteResponseDecrypt, &validResponse)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	if strings.ToLower(validResponse.Message) != constanta.StatusMessage {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	output = validResponse.Configuration
	err = errorModel.GenerateNonErrorModel()
	return
}
