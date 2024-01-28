package ProductLicenseService

import (
	"database/sql"
	"encoding/json"
	"github.com/Azure/go-autorest/autorest/date"
	"net/http"
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
	"time"
)

func (input productLicenseService) UpdateHWIDProductLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName       = "ActivateLicense"
		inputStruct    in.UpdateLicenseHWIDRequest
		additionalInfo interface{}
	)

	inputStruct, err = input.readBodyAndValidateHWID(request, contextModel, input.validateUpdateHWID)
	if err.Error != nil {
		return
	}

	additionalInfo, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateLicenseHWID, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- Function Additional
	})

	if err.Error != nil {
		return
	}

	if additionalInfo != nil {
		output.Data.Content = additionalInfo
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) doUpdateLicenseHWID(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName             = "doUpdateLicenseHWID"
		inputStruct          = inputStructInterface.(in.UpdateLicenseHWIDRequest)
		clientOnDB           repository.ClientCredentialModel
		productLicenses      []repository.ProductLicenseModel
		decryptedLicense     []in.ValidationLicenseJSONFile
		encryptedNewLicenses []in.ActivateLicenseDataRequest
		db                   = serverconfig.ServerAttribute.DBConnection
	)

	//--- Validate Credential Token
	clientOnDB, err = dao.ClientCredentialDAO.GetClientCredentialByClientID(db, repository.ClientCredentialModel{
		ClientID: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	})

	if err.Error != nil {
		return
	}

	if clientOnDB.ClientID.String == "" {
		err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
		return
	}

	//--- Get Product Licenses
	productLicenses, err = input.getProductLicensesForUpdateHWID(inputStruct, clientOnDB)
	if err.Error != nil {
		return
	}

	//--- Decrypt License
	decryptedLicense, err = input.decryptLicenseAndGetNewConfig(productLicenses, clientOnDB)
	if err.Error != nil {
		return
	}

	//--- Re Encrypt Licenses
	encryptedNewLicenses, err = input.generateNewLicense(decryptedLicense, clientOnDB, inputStruct)
	if err.Error != nil {
		return
	}

	//--- Update Product License
	output, dataAudit, err = input.doUpdateProductLicenseHWID(tx, encryptedNewLicenses, contextModel, timeNow)
	return
}

func (input productLicenseService) doUpdateProductLicenseHWID(tx *sql.Tx, encryptedNewLicenses []in.ActivateLicenseDataRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var response []out.LicenseResponse

	for _, license := range encryptedNewLicenses {
		var updatedID int64
		productLicenseModel := repository.ProductLicenseModel{
			LicenseConfigId:  sql.NullInt64{Int64: license.LicenseConfigID},
			ProductKey:       sql.NullString{String: license.ProductKey},
			ProductEncrypt:   sql.NullString{String: license.ProductEncrypt},
			ProductSignature: sql.NullString{String: license.ProductSignature},
			UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			HWID:             sql.NullString{String: license.Hwid},
		}

		updatedID, err = dao.ProductLicenseDAO.UpdateProductLicenseForHWID(tx, productLicenseModel)
		if err.Error != nil {
			return
		}

		dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ProductLicenseDAO.TableName, updatedID, 0)...)
		response = append(response, out.LicenseResponse{
			ProductKey:       license.ProductKey,
			ProductEncrypt:   license.ProductEncrypt,
			ProductSignature: license.ProductSignature,
			ClientTypeID:     license.ClientTypeID,
			UniqueID1:        license.UniqueID1,
			UniqueID2:        license.UniqueID2,
		})
	}

	output = response
	return
}

func (input productLicenseService) generateNewLicense(decryptedLicense []in.ValidationLicenseJSONFile, clientOnDB repository.ClientCredentialModel, inputStruct in.UpdateLicenseHWIDRequest) (result []in.ActivateLicenseDataRequest, err errorModel.ErrorModel) {
	var (
		funcName = "generateNewLicense"
		errorS   error
		dataLicense, dataProductReqGenerator,
		dataResponseEncrypt []byte
	)

	for i := 0; i < len(decryptedLicense); i++ {
		var validResponse cryptoModel.EncryptLicenseResponseModel

		dataLicense, errorS = json.Marshal(cryptoModel.JSONFileActivationLicenseModel{
			InstallationID:      decryptedLicense[i].InstallationID,
			ClientID:            decryptedLicense[i].ClientID,
			ProductID:           decryptedLicense[i].ProductID,
			LicenseVariantName:  decryptedLicense[i].LicenseVariantName,
			LicenseTypeName:     decryptedLicense[i].LicenseTypeName,
			DeploymentMethod:    decryptedLicense[i].DeploymentMethod,
			NumberOfUser:        decryptedLicense[i].NumberOfUser,
			UniqueID1:           decryptedLicense[i].UniqueID1,
			UniqueID2:           decryptedLicense[i].UniqueID2,
			ProductValidFromStr: decryptedLicense[i].ProductValidFromStr,
			ProductValidThruStr: decryptedLicense[i].ProductValidThruStr,
			LicenseStatus:       decryptedLicense[i].LicenseStatus,
			ModuleName1:         decryptedLicense[i].ModuleName1,
			ModuleName2:         decryptedLicense[i].ModuleName2,
			ModuleName3:         decryptedLicense[i].ModuleName3,
			ModuleName4:         decryptedLicense[i].ModuleName4,
			ModuleName5:         decryptedLicense[i].ModuleName5,
			ModuleName6:         decryptedLicense[i].ModuleName6,
			ModuleName7:         decryptedLicense[i].ModuleName7,
			ModuleName8:         decryptedLicense[i].ModuleName8,
			ModuleName9:         decryptedLicense[i].ModuleName9,
			ModuleName10:        decryptedLicense[i].ModuleName10,
			MaxOfflineDays:      decryptedLicense[i].MaxOfflineDays,
			IsConcurrentUser:    decryptedLicense[i].IsConcurrentUser,
			ProductComponent:    decryptedLicense[i].ProductComponent,
		})

		if errorS != nil {
			service.LogMessage(errorS.Error(), http.StatusInternalServerError)
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
			return
		}

		dataProductReqGenerator, errorS = json.Marshal(cryptoModel.EncryptLicenseRequest{
			SignatureKey: clientOnDB.SignatureKey.String,
			ClientSecret: clientOnDB.ClientSecret.String,
			EncryptKey:   util2.HashingPassword(clientOnDB.ClientID.String, clientOnDB.ClientSecret.String),
			ProductKey:   decryptedLicense[i].ProductKey,
			HardwareId:   inputStruct.Hwid,
		})
		if errorS != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
			return
		}

		argumentMap := make(map[string]string)
		argumentMap["args1"] = string(dataLicense)
		argumentMap["args2"] = string(dataProductReqGenerator)

		dataResponseEncrypt, err = service.GeneratorLicense(constanta.ProductEncryptAction, argumentMap)
		if err.Error != nil {
			service.LogMessage(err.CausedBy.Error(), 500)
			err = errorModel.GenerateUnknownError(input.FileName, funcName, err.CausedBy)
			return
		}

		errorS = json.Unmarshal(dataResponseEncrypt, &validResponse)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
			return
		}

		result = append(result, in.ActivateLicenseDataRequest{
			LicenseConfigID:  decryptedLicense[i].LicenseConfigID,
			ProductKey:       validResponse.ProductKey,
			ProductEncrypt:   validResponse.ProductEncrypt,
			ProductSignature: validResponse.ProductSignature,
			IsUserConcurrent: decryptedLicense[i].IsConcurrentUser,
			ClientID:         decryptedLicense[i].ClientID,
			ClientSecret:     clientOnDB.ClientSecret.String,
			Hwid:             inputStruct.Hwid,
			UniqueID1:        decryptedLicense[i].UniqueID1,
			UniqueID2:        decryptedLicense[i].UniqueID2,
			ClientTypeID:     decryptedLicense[i].ClientTypeID,
		})
	}

	return
}

func (input productLicenseService) decryptLicenseAndGetNewConfig(productLicenses []repository.ProductLicenseModel, clientOnDB repository.ClientCredentialModel) (decryptedLicenses []in.ValidationLicenseJSONFile, err errorModel.ErrorModel) {
	for i := 0; i < len(productLicenses); i++ {
		var (
			licenseConfigJSON in.ValidationLicenseJSONFile
			signatureKey      = clientOnDB.SignatureKey.String
			productEncrypt    = productLicenses[i].ProductEncrypt.String
		)

		//--- Start Decrypt License
		_, err = input.decryptLicense(cryptoModel.DecryptLicenseRequestModel{
			ProductSignature: productLicenses[i].ProductSignature.String,
			ClientId:         productLicenses[i].ClientId.String,
			ClientSecret:     productLicenses[i].ClientSecret.String,
			HardwareId:       productLicenses[i].HWID.String,
			ProductKey:       productLicenses[i].ProductKey.String,
			ProductId:        productLicenses[i].ProductId.String,
		}, signatureKey, productEncrypt)
		if err.Error != nil {
			return
		}
		//--- End Decrypt License

		//--- Start Get License Config
		licenseConfigJSON, err = input.getJSONLicense(productLicenses[i])
		if err.Error != nil {
			return
		}

		decryptedLicenses = append(decryptedLicenses, licenseConfigJSON)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) getJSONLicense(productLicense repository.ProductLicenseModel) (result in.ValidationLicenseJSONFile, err errorModel.ErrorModel) {
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

	return
}

func (input productLicenseService) getProductLicensesForUpdateHWID(inputStruct in.UpdateLicenseHWIDRequest, clientOnDB repository.ClientCredentialModel) (result []repository.ProductLicenseModel, err errorModel.ErrorModel) {
	var funcName = "getProductLicensesForUpdateHWID"
	for _, license := range inputStruct.Licenses {
		var productLicenseOnDB repository.ProductLicenseModel
		productLicenseOnDB, err = dao.ProductLicenseDAO.GetProductLicenseForUpdateHWID(serverconfig.ServerAttribute.DBConnection, repository.ProductLicenseModel{
			ProductKey:     sql.NullString{String: license.ProductKey},
			ProductEncrypt: sql.NullString{String: license.ProductEncrypt},
			ClientId:       clientOnDB.ClientID,
			ClientSecret:   clientOnDB.ClientSecret,
		})

		if err.Error != nil {
			return
		}

		if productLicenseOnDB.ID.Int64 < 1 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ProductLicense)
			return
		}

		result = append(result, productLicenseOnDB)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) validateUpdateHWID(inputStruct *in.UpdateLicenseHWIDRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateHWID()
}
