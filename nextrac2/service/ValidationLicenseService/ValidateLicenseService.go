package ValidationLicenseService

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/cryptoModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input validationLicenseService) ValidateLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "ValidateLicense"
		inputStruct in.ValidationLicenseRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateRequestValidationLicense)
	if err.Error != nil {
		return
	}

	additionalInfo, err := input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doValidateLicense, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		if additionalInfo != nil {
			output.Other = additionalInfo
		}
		return
	}

	if additionalInfo != nil {
		output.Data.Content = additionalInfo
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_ACTIVATION_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input validationLicenseService) doValidateLicense(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName        = "doValidateLicense"
		inputStruct     = inputStructInterface.(in.ValidationLicenseRequest)
		errorDetail     []out.ValidationLicenseErrorDetail
		clientOnDB      repository.ClientCredentialModel
		productLicenses []repository.ProductLicenseModel
		jsonLicenses    []in.ValidationLicenseJSONFile
		activationData  []in.ActivateLicenseDataRequest
	)

	if len(inputStruct.ProductDetail) > 0 {
		byt, _ := json.Marshal(inputStruct.ProductDetail)
		fmt.Println(fmt.Sprintf(`Data Product Detail -> %s`, string(byt)))
	}

	//--- Check Client
	clientOnDB, err = input.clientValidation(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	//--- Check Product License
	productLicenses, errorDetail = input.checkLicenseOnProductLicense(inputStruct, clientOnDB, contextModel)
	if input.checkError(errorDetail) {
		err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseGetProductLicenseError, nil)
		output = errorDetail
		return
	}

	//--- Decrypt Data License
	var dataAuditTemp []repository.AuditSystemModel
	jsonLicenses, dataAuditTemp, errorDetail = input.decryptProductLicense(tx, timeNow, productLicenses, clientOnDB, contextModel)
	if input.checkError(errorDetail) {
		err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseDecryptLicenseError, nil)
		output = errorDetail
		return
	}

	//--- Data Audit
	dataAudit = append(dataAudit, dataAuditTemp...)

	//--- Generate JSON License and Encrypt Data
	activationData, errorDetail = input.generateJSONLicenseAndEncrypt(jsonLicenses, clientOnDB, inputStruct, contextModel)
	if input.checkError(errorDetail) {
		err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseGetJSONFileError, nil)
		output = errorDetail
		return
	}

	output, dataAudit, err = input.updateProductLicense(tx, activationData, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	return
}

func (input validationLicenseService) clientValidation(inputStruct in.ValidationLicenseRequest, contextModel *applicationModel.ContextModel) (clientOnDB repository.ClientCredentialModel, err errorModel.ErrorModel) {
	var (
		funcName = "clientValidation"
		db       = serverconfig.ServerAttribute.DBConnection
	)

	if contextModel.AuthAccessTokenModel.ClientID != inputStruct.ClientID {
		err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
		return
	}

	clientOnDB, err = dao.ClientCredentialDAO.GetClientCredentialForActivationLicense(db, repository.ClientCredentialModel{
		ClientID:     sql.NullString{String: inputStruct.ClientID},
		ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
		SignatureKey: sql.NullString{String: inputStruct.SignatureKey},
	})

	if err.Error != nil {
		return
	}

	if clientOnDB.ClientID.String == "" {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientMappingClientID)
		return
	}

	return
}

func (input validationLicenseService) updateProductLicense(tx *sql.Tx, inputStruct []in.ActivateLicenseDataRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName        = "updateProductLicense"
		productLicenses repository.ProductLicenseModel
		response        []out.LicenseResponseWithSalesman
	)

	for i := 0; i < len(inputStruct); i++ {
		productLicenses = repository.ProductLicenseModel{
			LicenseConfigId:  sql.NullInt64{Int64: inputStruct[i].LicenseConfigID},
			ProductKey:       sql.NullString{String: inputStruct[i].ProductKey},
			ProductEncrypt:   sql.NullString{String: inputStruct[i].ProductEncrypt},
			ProductSignature: sql.NullString{String: inputStruct[i].ProductSignature},
			UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		}

		err = dao.ProductLicenseDAO.UpdateProductLicenseForValidationLicense(tx, productLicenses)
		if err.Error != nil {
			err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseCreateProductLicenseError, err.CausedBy)
			return
		}

		var salesmanListOut []out.SalesmanLicenseListOut
		for _, itemSalesman := range inputStruct[i].SalesmanList {
			salesmanList := out.SalesmanLicenseListOut{
				ID:         itemSalesman.ID,
				AuthUserID: itemSalesman.AuthUserID,
				UserID:     itemSalesman.UserID,
				Status:     itemSalesman.Status,
			}
			salesmanListOut = append(salesmanListOut, salesmanList)
		}

		response = append(response, out.LicenseResponseWithSalesman{
			LicenseResponse: out.LicenseResponse{
				ProductKey:       inputStruct[i].ProductKey,
				ProductEncrypt:   inputStruct[i].ProductEncrypt,
				ProductSignature: inputStruct[i].ProductSignature,
				ClientTypeID:     inputStruct[i].ClientTypeID,
				UniqueID1:        inputStruct[i].UniqueID1,
				UniqueID2:        inputStruct[i].UniqueID2,
			},
			SalesmanList: salesmanListOut,
		})
	}

	output = response
	return
}

func (input validationLicenseService) generateJSONLicenseAndEncrypt(inputStruct []in.ValidationLicenseJSONFile, clientRequested repository.ClientCredentialModel, userStruct in.ValidationLicenseRequest, contextModel *applicationModel.ContextModel) (result []in.ActivateLicenseDataRequest, errorDetail []out.ValidationLicenseErrorDetail) {
	var (
		err      errorModel.ErrorModel
		funcName = "generateNewLicense"
		errorS   error
		dataLicense, dataProductReqGenerator,
		dataResponseEncrypt []byte
	)

	for i := 0; i < len(inputStruct); i++ {
		var validResponse cryptoModel.EncryptLicenseResponseModel

		dataLicense, errorS = json.Marshal(cryptoModel.JSONFileActivationLicenseModel{
			InstallationID:      inputStruct[i].InstallationID,
			ClientID:            inputStruct[i].ClientID,
			ProductID:           inputStruct[i].ProductID,
			LicenseVariantName:  inputStruct[i].LicenseVariantName,
			LicenseTypeName:     inputStruct[i].LicenseTypeName,
			DeploymentMethod:    inputStruct[i].DeploymentMethod,
			NumberOfUser:        inputStruct[i].NumberOfUser,
			UniqueID1:           inputStruct[i].UniqueID1,
			UniqueID2:           inputStruct[i].UniqueID2,
			ProductValidFromStr: inputStruct[i].ProductValidFromStr,
			ProductValidThruStr: inputStruct[i].ProductValidThruStr,
			LicenseStatus:       inputStruct[i].LicenseStatus,
			ModuleName1:         inputStruct[i].ModuleName1,
			ModuleName2:         inputStruct[i].ModuleName2,
			ModuleName3:         inputStruct[i].ModuleName3,
			ModuleName4:         inputStruct[i].ModuleName4,
			ModuleName5:         inputStruct[i].ModuleName5,
			ModuleName6:         inputStruct[i].ModuleName6,
			ModuleName7:         inputStruct[i].ModuleName7,
			ModuleName8:         inputStruct[i].ModuleName8,
			ModuleName9:         inputStruct[i].ModuleName9,
			ModuleName10:        inputStruct[i].ModuleName10,
			MaxOfflineDays:      inputStruct[i].MaxOfflineDays,
			IsConcurrentUser:    inputStruct[i].IsConcurrentUser,
			ProductComponent:    inputStruct[i].ProductComponent,
		})

		if errorS != nil {
			service.LogMessage(errorS.Error(), http.StatusInternalServerError)
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
			errorDetail = append(errorDetail, out.ValidationLicenseErrorDetail{
				LicenseConfigID: inputStruct[i].LicenseConfigID,
				Message:         service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		dataProductReqGenerator, errorS = json.Marshal(cryptoModel.EncryptLicenseRequest{
			SignatureKey: clientRequested.SignatureKey.String,
			ClientSecret: clientRequested.ClientSecret.String,
			EncryptKey:   util2.HashingPassword(clientRequested.ClientID.String, clientRequested.ClientSecret.String),
			ProductKey:   inputStruct[i].ProductKey,
			HardwareId:   inputStruct[i].HWID,
		})
		if errorS != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
			errorDetail = append(errorDetail, out.ValidationLicenseErrorDetail{
				LicenseConfigID: inputStruct[i].LicenseConfigID,
				Message:         service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		argumentMap := make(map[string]string)
		argumentMap["args1"] = string(dataLicense)
		argumentMap["args2"] = string(dataProductReqGenerator)

		dataResponseEncrypt, err = service.GeneratorLicense(constanta.ProductEncryptAction, argumentMap)
		if err.Error != nil {
			service.LogMessage(err.CausedBy.Error(), 500)
			err = errorModel.GenerateUnknownError(input.FileName, funcName, err.CausedBy)
			errorDetail = append(errorDetail, out.ValidationLicenseErrorDetail{
				LicenseConfigID: inputStruct[i].LicenseConfigID,
				Message:         service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		errorS = json.Unmarshal(dataResponseEncrypt, &validResponse)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
			errorDetail = append(errorDetail, out.ValidationLicenseErrorDetail{
				LicenseConfigID: inputStruct[i].LicenseConfigID,
				Message:         service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		result = append(result, in.ActivateLicenseDataRequest{
			LicenseConfigID:  inputStruct[i].LicenseConfigID,
			ProductKey:       validResponse.ProductKey,
			ProductEncrypt:   validResponse.ProductEncrypt,
			ProductSignature: validResponse.ProductSignature,
			IsUserConcurrent: inputStruct[i].IsConcurrentUser,
			ClientID:         inputStruct[i].ClientID,
			ClientSecret:     clientRequested.ClientSecret.String,
			Hwid:             inputStruct[i].HWID,
			UniqueID1:        inputStruct[i].UniqueID1,
			UniqueID2:        inputStruct[i].UniqueID2,
			ClientTypeID:     inputStruct[i].ClientTypeID,
			SalesmanList:     inputStruct[i].SalesmanLicenseList,
		})
	}

	return
}

func (input validationLicenseService) generateJSONLicense(inputStruct in.ValidationLicenseJSONFile, clientRequested repository.ClientCredentialModel) (result cryptoModel.EncryptLicenseRequestModel) {
	result = cryptoModel.EncryptLicenseRequestModel{
		SignatureKey: clientRequested.SignatureKey.String,
		ClientSecret: clientRequested.ClientSecret.String,
		Hwid:         inputStruct.HWID,
	}

	result.LicenseConfigData = cryptoModel.JSONFileActivationLicenseModel{
		InstallationID:      inputStruct.InstallationID,
		ClientID:            inputStruct.ClientID,
		ProductID:           inputStruct.ProductID,
		LicenseVariantName:  inputStruct.LicenseVariantName,
		LicenseTypeName:     inputStruct.LicenseTypeName,
		DeploymentMethod:    inputStruct.DeploymentMethod,
		NumberOfUser:        inputStruct.NumberOfUser,
		UniqueID1:           inputStruct.UniqueID1,
		UniqueID2:           inputStruct.UniqueID2,
		ProductValidFrom:    date.Date{inputStruct.ProductValidFrom.Time},
		ProductValidFromStr: inputStruct.ProductValidFrom.String(),
		ProductValidThru:    date.Date{inputStruct.ProductValidThru.Time},
		ProductValidThruStr: inputStruct.ProductValidThru.String(),
		LicenseStatus:       inputStruct.LicenseStatus,
		ModuleName1:         inputStruct.ModuleName1,
		ModuleName2:         inputStruct.ModuleName2,
		ModuleName3:         inputStruct.ModuleName3,
		ModuleName4:         inputStruct.ModuleName4,
		ModuleName5:         inputStruct.ModuleName5,
		ModuleName6:         inputStruct.ModuleName6,
		ModuleName7:         inputStruct.ModuleName7,
		ModuleName8:         inputStruct.ModuleName8,
		ModuleName9:         inputStruct.ModuleName9,
		ModuleName10:        inputStruct.ModuleName10,
		MaxOfflineDays:      inputStruct.MaxOfflineDays,
		IsConcurrentUser:    inputStruct.IsConcurrentUser,
		ProductComponent:    inputStruct.ProductComponent,
	}

	return
}

func (input validationLicenseService) getLicenseConfig(jsonLicenses []in.ValidationLicenseJSONFile, contextModel *applicationModel.ContextModel) (errorDetail []out.ValidationLicenseErrorDetail) {
	//var wg sync.WaitGroup
	funcName := "getLicenseConfig"

	for i := 0; i < len(jsonLicenses); i++ {
		var components []cryptoModel.ProductComponents
		jsonLicenseOnDB, err := dao.LicenseConfigDAO.GetLicenseConfigForValidation(serverconfig.ServerAttribute.DBConnection, repository.LicenseConfigModel{
			ID: sql.NullInt64{Int64: jsonLicenses[i].LicenseConfigID},
		})

		if err.Error != nil {
			//err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseGetLicenseConfigError, nil)
			errorDetail = append(errorDetail, out.ValidationLicenseErrorDetail{
				ProductKey:     jsonLicenses[i].ProductKey,
				ProductEncrypt: jsonLicenses[i].ProductEncrypt,
				Message:        service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		if jsonLicenseOnDB.ID.Int64 < 1 {
			err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseGetLicenseConfigError, nil)
			errorDetail = append(errorDetail, out.ValidationLicenseErrorDetail{
				LicenseConfigID: jsonLicenseOnDB.ID.Int64,
				Message:         service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		jsonLicenses[i].InstallationID = jsonLicenseOnDB.InstallationID.Int64
		jsonLicenses[i].ClientID = jsonLicenseOnDB.ClientID.String
		jsonLicenses[i].ProductID = jsonLicenseOnDB.ProductCode.String
		jsonLicenses[i].LicenseVariantName = jsonLicenseOnDB.LicenseVariantName.String
		jsonLicenses[i].LicenseTypeName = jsonLicenseOnDB.LicenseTypeName.String
		jsonLicenses[i].DeploymentMethod = jsonLicenseOnDB.DeploymentMethod.String
		jsonLicenses[i].NumberOfUser = jsonLicenseOnDB.NoOfUser.Int64
		jsonLicenses[i].UniqueID1 = jsonLicenseOnDB.UniqueID1.String
		jsonLicenses[i].UniqueID2 = jsonLicenseOnDB.UniqueID2.String
		jsonLicenses[i].ProductValidFrom = date.Date{Time: jsonLicenseOnDB.ProductValidFrom.Time}
		jsonLicenses[i].ProductValidThru = date.Date{Time: jsonLicenseOnDB.ProductValidThru.Time}
		jsonLicenses[i].LicenseStatus = jsonLicenseOnDB.ProductLicenseStatus.Int64
		jsonLicenses[i].ModuleName1 = jsonLicenseOnDB.ModuleIDName1.String
		jsonLicenses[i].ModuleName2 = jsonLicenseOnDB.ModuleIDName2.String
		jsonLicenses[i].ModuleName3 = jsonLicenseOnDB.ModuleIDName3.String
		jsonLicenses[i].ModuleName4 = jsonLicenseOnDB.ModuleIDName4.String
		jsonLicenses[i].ModuleName5 = jsonLicenseOnDB.ModuleIDName5.String
		jsonLicenses[i].ModuleName6 = jsonLicenseOnDB.ModuleIDName6.String
		jsonLicenses[i].ModuleName7 = jsonLicenseOnDB.ModuleIDName7.String
		jsonLicenses[i].ModuleName8 = jsonLicenseOnDB.ModuleIDName8.String
		jsonLicenses[i].ModuleName9 = jsonLicenseOnDB.ModuleIDName9.String
		jsonLicenses[i].ModuleName10 = jsonLicenseOnDB.ModuleIDName10.String
		jsonLicenses[i].MaxOfflineDays = jsonLicenseOnDB.MaxOfflineDays.Int64
		jsonLicenses[i].IsConcurrentUser = jsonLicenseOnDB.IsUserConcurrent.String
		jsonLicenses[i].LicenseConfigID = jsonLicenseOnDB.ID.Int64
		jsonLicenses[i].ClientTypeID = jsonLicenseOnDB.ClientTypeID.Int64

		if jsonLicenseOnDB.ComponentSting.String != "" {
			errorUnmarshal := json.Unmarshal([]byte(jsonLicenseOnDB.ComponentSting.String), &components)
			if errorUnmarshal != nil {
				err = errorModel.GenerateUnknownError(input.FileName, funcName, errorUnmarshal)
				return
			}
			jsonLicenses[i].ProductComponent = components
		}
	}

	return
}

func (input validationLicenseService) checkLicenseConfig(resultLicense chan []repository.LicenseConfigModel, resultError chan []out.ValidationLicenseErrorDetail, licenseDetail []in.ValidationLicenseJSONFile, contextModel *applicationModel.ContextModel, wg *sync.WaitGroup) {
	var (
		funcName                 = "checkLicenseConfig"
		err                      errorModel.ErrorModel
		tempResult, licenseModel []repository.LicenseConfigModel
		tempErrorDetail          []out.ValidationLicenseErrorDetail
	)

	for i := 0; i < len(licenseDetail); i++ {
		licenseModel = append(licenseModel, repository.LicenseConfigModel{
			ID: sql.NullInt64{Int64: licenseDetail[i].LicenseConfigID},
		})
	}

	tempResult, err = dao.LicenseConfigDAO.GetLicenseConfigForValidationLicense(serverconfig.ServerAttribute.DBConnection, licenseModel)
	if err.Error != nil || len(tempResult) < 1 {
		for i := 0; i < len(licenseDetail); i++ {
			err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseGetLicenseConfigError, nil)
			tempErrorDetail = append(tempErrorDetail, out.ValidationLicenseErrorDetail{
				ProductKey:     licenseDetail[i].ProductKey,
				ProductEncrypt: licenseDetail[i].ProductEncrypt,
				Message:        service.GetErrorMessage(err, *contextModel),
			})
		}
		return
	}

	if len(tempResult) != len(licenseDetail) {
		differentLicenses := input.checkDifferentLicense(licenseModel, tempResult)
		for i := 0; i < len(differentLicenses); i++ {
			err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseGetLicenseConfigError, nil)
			tempErrorDetail = append(tempErrorDetail, out.ValidationLicenseErrorDetail{
				LicenseConfigID: differentLicenses[i].ID.Int64,
				Message:         service.GetErrorMessage(err, *contextModel),
			})
		}
		return
	}

	defer func() {
		resultLicense <- tempResult
		resultError <- tempErrorDetail
		wg.Done()
	}()
}

func (input validationLicenseService) decryptProductLicense(tx *sql.Tx, timeNow time.Time, productLicenses []repository.ProductLicenseModel, clientRequested repository.ClientCredentialModel, contextModel *applicationModel.ContextModel) (result []in.ValidationLicenseJSONFile, dataAuditTemp []repository.AuditSystemModel, errorDetail []out.ValidationLicenseErrorDetail) {
	var (
		funcName = "decryptLicenseAndGetNewConfig"
		err      errorModel.ErrorModel
	)

	for i := 0; i < len(productLicenses); i++ {
		var (
			tempResult              cryptoModel.DecryptLicenseResponseModel
			licenseConfigJSON       in.ValidationLicenseJSONFile
			dataByteResponseDecrypt []byte
		)

		//--- Data Audit Validate License
		dataAuditTemp = append(dataAuditTemp, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ProductLicenseDAO.TableName, productLicenses[i].ID.Int64, 0)...)

		// ### Start Decrypt License
		dataByteDecryptLicense, errorS := json.Marshal(cryptoModel.DecryptLicenseRequestModel{
			SignatureKey:     clientRequested.SignatureKey.String,
			ProductSignature: productLicenses[i].ProductSignature.String,
			ClientId:         productLicenses[i].ClientId.String,
			ClientSecret:     productLicenses[i].ClientSecret.String,
			EncryptKey:       util2.HashingPassword(productLicenses[i].ClientId.String, productLicenses[i].ClientSecret.String),
			HardwareId:       productLicenses[i].HWID.String,
			ProductKey:       productLicenses[i].ProductKey.String,
			ProductId:        productLicenses[i].ProductId.String,
		})
		if errorS != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
			return
		}

		//--- Set Argument Generator License
		argumentMap := make(map[string]string)
		argumentMap["args1"] = productLicenses[i].ProductEncrypt.String
		argumentMap["args2"] = string(dataByteDecryptLicense)

		dataByteResponseDecrypt, err = service.GeneratorLicense(constanta.ProductDecryptAction, argumentMap)
		if err.Error != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
			errorDetail = append(errorDetail, out.ValidationLicenseErrorDetail{
				ProductKey:     productLicenses[i].ProductKey.String,
				ProductEncrypt: productLicenses[i].ProductEncrypt.String,
				Message:        service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		//--- Convert Response Generator
		errorS = json.Unmarshal(dataByteResponseDecrypt, &tempResult)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
			errorDetail = append(errorDetail, out.ValidationLicenseErrorDetail{
				ProductKey:     productLicenses[i].ProductKey.String,
				ProductEncrypt: productLicenses[i].ProductEncrypt.String,
				Message:        service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		if strings.ToLower(tempResult.MessageCode) != constanta.StatusMessage {
			err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, tempResult.Notification, errors.New(tempResult.Message))
			errorDetail = append(errorDetail, out.ValidationLicenseErrorDetail{
				ProductKey:     productLicenses[i].ProductKey.String,
				ProductEncrypt: productLicenses[i].ProductEncrypt.String,
				Message:        service.GetErrorMessage(err, *contextModel),
			})
			continue
		}
		// ### End Decrypt License

		// ### Start Get License Config
		var salesman []cryptoModel.SalesmanList
		licenseConfigJSON, salesman, err = input.getJSONLicense(productLicenses[i])
		if err.Error != nil {
			errorDetail = append(errorDetail, out.ValidationLicenseErrorDetail{
				ProductKey:     productLicenses[i].ProductKey.String,
				ProductEncrypt: productLicenses[i].ProductEncrypt.String,
				Message:        service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		// ### End Get License Config
		var salesmanList []in.SalesmanLicenseList
		for _, itemSalesman := range salesman {
			salesmanList = append(salesmanList, in.SalesmanLicenseList{
				ID:         itemSalesman.ID,
				AuthUserID: itemSalesman.AuthUserID,
				UserID:     itemSalesman.UserID,
				Status:     itemSalesman.Status,
			})
		}

		licenseConfigJSON.SalesmanLicenseList = salesmanList
		result = append(result, licenseConfigJSON)
	}

	return
}

func (input validationLicenseService) getJSONLicense(productLicense repository.ProductLicenseModel) (result in.ValidationLicenseJSONFile, salesman []cryptoModel.SalesmanList, err errorModel.ErrorModel) {
	var (
		licenseConfigOnDB repository.LicenseConfigModel
		components        []cryptoModel.ProductComponents
		funcName          = "getJSONLicense"
	)

	licenseConfigOnDB, err = dao.LicenseConfigDAO.GetLicenseConfigForValidation(serverconfig.ServerAttribute.DBConnection, repository.LicenseConfigModel{
		ID: productLicense.LicenseConfigId,
	})

	if err.Error != nil {
		return
	}

	if licenseConfigOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseGetLicenseConfigError, nil)
		return
	}

	tempJSONFile := cryptoModel.JSONFileActivationLicenseModel{
		InstallationID:     licenseConfigOnDB.InstallationID.Int64,
		ClientID:           licenseConfigOnDB.ClientID.String,
		ProductID:          licenseConfigOnDB.ProductCode.String,
		LicenseVariantName: licenseConfigOnDB.LicenseVariantName.String,
		LicenseTypeName:    licenseConfigOnDB.LicenseTypeName.String,
		DeploymentMethod:   licenseConfigOnDB.DeploymentMethod.String,
		NumberOfUser:       licenseConfigOnDB.NoOfUser.Int64,
		UniqueID1:          licenseConfigOnDB.UniqueID1.String,
		UniqueID2:          licenseConfigOnDB.UniqueID2.String,
		ProductValidFrom:   date.Date{Time: licenseConfigOnDB.ProductValidFrom.Time},
		ProductValidThru:   date.Date{Time: licenseConfigOnDB.ProductValidThru.Time},
		LicenseStatus:      licenseConfigOnDB.ProductLicenseStatus.Int64,
		ModuleName1:        licenseConfigOnDB.ModuleIDName1.String,
		ModuleName2:        licenseConfigOnDB.ModuleIDName2.String,
		ModuleName3:        licenseConfigOnDB.ModuleIDName3.String,
		ModuleName4:        licenseConfigOnDB.ModuleIDName4.String,
		ModuleName5:        licenseConfigOnDB.ModuleIDName5.String,
		ModuleName6:        licenseConfigOnDB.ModuleIDName6.String,
		ModuleName7:        licenseConfigOnDB.ModuleIDName7.String,
		ModuleName8:        licenseConfigOnDB.ModuleIDName8.String,
		ModuleName9:        licenseConfigOnDB.ModuleIDName9.String,
		ModuleName10:       licenseConfigOnDB.ModuleIDName10.String,
		MaxOfflineDays:     licenseConfigOnDB.MaxOfflineDays.Int64,
		IsConcurrentUser:   licenseConfigOnDB.IsUserConcurrent.String,
	}

	result = in.ValidationLicenseJSONFile{
		JSONFileActivationLicenseModel: tempJSONFile,
		LicenseConfigID:                licenseConfigOnDB.ID.Int64,
		ClientTypeID:                   licenseConfigOnDB.ClientTypeID.Int64,
		ProductKey:                     productLicense.ProductKey.String,
		ProductEncrypt:                 productLicense.ProductEncrypt.String,
		HWID:                           productLicense.HWID.String,
	}

	result.ProductValidFromStr = result.ProductValidFrom.String()
	result.ProductValidThruStr = result.ProductValidThru.String()

	if licenseConfigOnDB.ComponentSting.String != "" {
		errorUnmarshal := json.Unmarshal([]byte(licenseConfigOnDB.ComponentSting.String), &components)
		if errorUnmarshal != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorUnmarshal)
			return
		}
		result.ProductComponent = components
	}

	if licenseConfigOnDB.SalesmanString.String != "" {
		errorUnmarshal := json.Unmarshal([]byte(licenseConfigOnDB.SalesmanString.String), &salesman)
		if errorUnmarshal != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorUnmarshal)
			return
		}
	}

	return
}

func (input validationLicenseService) checkLicenseOnProductLicense(inputStruct in.ValidationLicenseRequest, clientModel repository.ClientCredentialModel, contextModel *applicationModel.ContextModel) (payload []repository.ProductLicenseModel, errorDetail []out.ValidationLicenseErrorDetail) {
	var wg sync.WaitGroup

	totalPage := math.Ceil(float64(len(inputStruct.ProductDetail)) / float64(constanta.TotalDataProductLicensePerChannel))
	resultLicense := make(chan []repository.ProductLicenseModel, len(inputStruct.ProductDetail))
	resultError := make(chan []out.ValidationLicenseErrorDetail, len(inputStruct.ProductDetail))

	for i := 1; i <= int(totalPage); i++ {
		wg.Add(1)
		var licenses []in.ProductEncryptDetail

		offset := dao.CountOffset(i, constanta.TotalDataProductLicensePerChannel)
		until := offset + constanta.TotalDataProductLicensePerChannel

		if i == int(totalPage) {
			licenses = append(licenses, inputStruct.ProductDetail[offset:]...)
		} else {
			licenses = append(licenses, inputStruct.ProductDetail[offset:until]...)
		}

		go input.doGetProductLicenses(inputStruct, resultLicense, resultError, licenses, clientModel, contextModel, &wg)
	}

	for i := 0; i < int(totalPage); i++ {
		tempResult := <-resultLicense
		tempErrorDetail := <-resultError
		payload = append(payload, tempResult...)
		errorDetail = append(errorDetail, tempErrorDetail...)
	}

	wg.Wait()
	close(resultLicense)
	close(resultError)
	return
}

// Action
func (input validationLicenseService) doGenerateJSONLicense(licenseConfig repository.LicenseConfigModel) (result cryptoModel.JSONFileActivationLicenseModel, err errorModel.ErrorModel) {
	funcName := "doGenerateJSONLicense"
	var components []cryptoModel.ProductComponents
	result = cryptoModel.JSONFileActivationLicenseModel{
		InstallationID:     licenseConfig.InstallationID.Int64,
		ClientID:           licenseConfig.ClientID.String,
		ProductID:          licenseConfig.ProductCode.String,
		LicenseVariantName: licenseConfig.LicenseVariantName.String,
		LicenseTypeName:    licenseConfig.LicenseTypeName.String,
		DeploymentMethod:   licenseConfig.DeploymentMethod.String,
		NumberOfUser:       licenseConfig.NoOfUser.Int64,
		UniqueID1:          licenseConfig.UniqueID1.String,
		UniqueID2:          licenseConfig.UniqueID2.String,
		ProductValidFrom:   date.Date{licenseConfig.ProductValidFrom.Time},
		ProductValidThru:   date.Date{licenseConfig.ProductValidThru.Time},
		LicenseStatus:      licenseConfig.ProductLicenseStatus.Int64,
		ModuleName1:        licenseConfig.ModuleIDName1.String,
		ModuleName2:        licenseConfig.ModuleIDName2.String,
		ModuleName3:        licenseConfig.ModuleIDName3.String,
		ModuleName4:        licenseConfig.ModuleIDName4.String,
		ModuleName5:        licenseConfig.ModuleIDName5.String,
		ModuleName6:        licenseConfig.ModuleIDName6.String,
		ModuleName7:        licenseConfig.ModuleIDName7.String,
		ModuleName8:        licenseConfig.ModuleIDName8.String,
		ModuleName9:        licenseConfig.ModuleIDName9.String,
		ModuleName10:       licenseConfig.ModuleIDName10.String,
		MaxOfflineDays:     licenseConfig.MaxOfflineDays.Int64,
		IsConcurrentUser:   licenseConfig.IsUserConcurrent.String,
	}

	if licenseConfig.ComponentSting.String != "" {
		errorUnmarshal := json.Unmarshal([]byte(licenseConfig.ComponentSting.String), &components)
		if errorUnmarshal != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorUnmarshal)
			return
		}
		result.ProductComponent = components
	}

	return
}

func (input validationLicenseService) doDecryptLicense(inputStruct in.GenerateDataProductConfiguration, productLicenseOnDB repository.ProductLicenseModel) (result cryptoModel.JSONFileActivationLicenseModel, err errorModel.ErrorModel) {
	funcName := "doDecryptLicense"
	var validResponse in.GenerateDataValidationResponse

	// todo decrypt
	dataByteDecryptLicense, errorS := json.Marshal(inputStruct)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	// set argument generator license
	argumentMap := make(map[string]string)
	argumentMap["args1"] = productLicenseOnDB.ProductEncrypt.String
	argumentMap["args2"] = string(dataByteDecryptLicense)

	dataByteResponseDecrypt, err := service.GeneratorLicense(constanta.ProductDecryptAction, argumentMap)
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

	if strings.ToLower(validResponse.MessageCode) != constanta.StatusMessage {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errors.New(validResponse.Message))
		return
	}

	productValidThru, _ := in.TimeStrToTimeWithTimeFormat(validResponse.Configuration.ProductValidThru, constanta.ProductValidThru, constanta.DefaultInstallationTimeFormat)
	productValidFrom, _ := in.TimeStrToTimeWithTimeFormat(validResponse.Configuration.ProductValidFrom, constanta.ProductValidFrom, constanta.DefaultInstallationTimeFormat)

	var productComponents []cryptoModel.ProductComponents

	for _, component := range validResponse.Configuration.Component {
		productComponents = append(productComponents, cryptoModel.ProductComponents{
			ComponentName:  component.Name,
			ComponentValue: component.Value,
		})
	}

	result = cryptoModel.JSONFileActivationLicenseModel{
		InstallationID:     validResponse.Configuration.InstallationId,
		ClientID:           validResponse.Configuration.ClientId,
		ProductID:          validResponse.Configuration.ProductId,
		LicenseVariantName: validResponse.Configuration.LicenseVariantName,
		LicenseTypeName:    validResponse.Configuration.LicenseTypeName,
		DeploymentMethod:   validResponse.Configuration.DeploymentMethod,
		NumberOfUser:       validResponse.Configuration.NoOfUser,
		UniqueID1:          validResponse.Configuration.UniqueId1,
		UniqueID2:          validResponse.Configuration.UniqueId2,
		ProductValidFrom:   date.Date{Time: productValidFrom},
		ProductValidThru:   date.Date{Time: productValidThru},
		LicenseStatus:      validResponse.Configuration.LicenseStatus,
		ModuleName1:        validResponse.Configuration.ModuleName1,
		ModuleName2:        validResponse.Configuration.ModuleName2,
		ModuleName3:        validResponse.Configuration.ModuleName3,
		ModuleName4:        validResponse.Configuration.ModuleName4,
		ModuleName5:        validResponse.Configuration.ModuleName5,
		ModuleName6:        validResponse.Configuration.ModuleName6,
		ModuleName7:        validResponse.Configuration.ModuleName7,
		ModuleName8:        validResponse.Configuration.ModuleName8,
		ModuleName9:        validResponse.Configuration.ModuleName9,
		ModuleName10:       validResponse.Configuration.ModuleName10,
		MaxOfflineDays:     validResponse.Configuration.MaxOfflineDays,
		IsConcurrentUser:   validResponse.Configuration.IsConcurrentUser,
		ProductComponent:   productComponents,
	}

	return
}

func (input validationLicenseService) doGetProductLicenses(inputStruct in.ValidationLicenseRequest, resultLicenses chan []repository.ProductLicenseModel, resultError chan []out.ValidationLicenseErrorDetail, productLicenses []in.ProductEncryptDetail, clientOnDB repository.ClientCredentialModel, contextModel *applicationModel.ContextModel, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		funcName        = "doGetProductLicenses"
		err             errorModel.ErrorModel
		tempResult      []repository.ProductLicenseModel
		tempErrorDetail []out.ValidationLicenseErrorDetail
	)

	for _, licenses := range productLicenses {
		var productLicenseOnDB repository.ProductLicenseModel
		productLicenseOnDB, err = dao.ProductLicenseDAO.GetProductLicenseForValidation(serverconfig.ServerAttribute.DBConnection, repository.ProductLicenseModel{
			ProductKey:     sql.NullString{String: licenses.ProductKey},
			ProductEncrypt: sql.NullString{String: licenses.ProductEncrypt},
			ClientId:       clientOnDB.ClientID,
			ClientSecret:   clientOnDB.ClientSecret,
			HWID:           sql.NullString{String: inputStruct.HwID},
		})

		if err.Error != nil {
			err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseGetProductLicenseError, err.CausedBy)
			tempErrorDetail = append(tempErrorDetail, out.ValidationLicenseErrorDetail{
				ProductKey:     licenses.ProductKey,
				ProductEncrypt: licenses.ProductEncrypt,
				Message:        service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		if productLicenseOnDB.ID.Int64 < 1 {
			err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseGetProductLicenseError, nil)
			tempErrorDetail = append(tempErrorDetail, out.ValidationLicenseErrorDetail{
				ProductKey:     licenses.ProductKey,
				ProductEncrypt: licenses.ProductEncrypt,
				Message:        service.GetErrorMessage(err, *contextModel),
			})
			continue
		}
		tempResult = append(tempResult, productLicenseOnDB)
	}

	resultLicenses <- tempResult
	resultError <- tempErrorDetail
}

// Util
func (input validationLicenseService) validateRequestValidationLicense(inputStruct *in.ValidationLicenseRequest) errorModel.ErrorModel {
	return inputStruct.ValidateRequestValidationLicense()
}

func (input validationLicenseService) checkDifferentLicense(licensesID, existingLicenses []repository.LicenseConfigModel) (result []repository.LicenseConfigModel) {
	tempMap := make(map[int64]bool)

	for _, item := range licensesID {
		tempMap[item.ID.Int64] = true
	}

	for _, item := range existingLicenses {
		if _, ok := tempMap[item.ID.Int64]; !ok {
			result = append(result, item)
		}
	}
	return
}

func (input validationLicenseService) checkError(errorDetail []out.ValidationLicenseErrorDetail) bool {
	if len(errorDetail) > 0 {
		return true
	}

	return false
}
