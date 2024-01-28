package ValidationLicenseService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/cryptoModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
)

func GenerateI18NMessage(messageID string, language string) (output string) {
	return util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ValidationLicenseBundle, messageID, language, nil)
}

func (input validationLicenseService) GenerateLicenseEncrypt(inputStruct cryptoModel.EncryptLicenseRequestModel, productKey string) (result cryptoModel.EncryptLicenseResponseModel, err errorModel.ErrorModel) {
	var (
		funcName                             = "GenerateLicenseEncrypt"
		errorS                               error
		argumentMap                          map[string]string
		validResponse                        in.GenerateDataValidationResponse
		dataLicense, dataProductReqGenerator []byte
		dataResponseEncrypt                  []byte
	)

	result = cryptoModel.EncryptLicenseResponseModel{
		Notification: "OK",
	}

	var productComponent []in.GenerateDataComponent
	for i := 0; i < len(inputStruct.LicenseConfigData.ProductComponent); i++ {
		productComponent = append(productComponent, in.GenerateDataComponent{
			Name:  inputStruct.LicenseConfigData.ProductComponent[i].ComponentName,
			Value: inputStruct.LicenseConfigData.ProductComponent[i].ComponentValue,
		})
	}

	licenseConfigReqGenerator := in.GenerateDataLicenseConfiguration{
		InstallationId:     inputStruct.LicenseConfigData.InstallationID,
		ClientId:           inputStruct.LicenseConfigData.ClientID,
		ProductId:          inputStruct.LicenseConfigData.ProductID,
		LicenseVariantName: inputStruct.LicenseConfigData.LicenseVariantName,
		LicenseTypeName:    inputStruct.LicenseConfigData.LicenseTypeName,
		DeploymentMethod:   inputStruct.LicenseConfigData.DeploymentMethod,
		NoOfUser:           inputStruct.LicenseConfigData.NumberOfUser,
		UniqueId1:          inputStruct.LicenseConfigData.UniqueID1,
		UniqueId2:          inputStruct.LicenseConfigData.UniqueID2,
		ProductValidFrom:   inputStruct.LicenseConfigData.ProductValidFrom.String(),
		ProductValidThru:   inputStruct.LicenseConfigData.ProductValidThru.String(),
		LicenseStatus:      inputStruct.LicenseConfigData.LicenseStatus,
		ModuleName1:        inputStruct.LicenseConfigData.ModuleName1,
		ModuleName2:        inputStruct.LicenseConfigData.ModuleName2,
		ModuleName3:        inputStruct.LicenseConfigData.ModuleName3,
		ModuleName4:        inputStruct.LicenseConfigData.ModuleName4,
		ModuleName5:        inputStruct.LicenseConfigData.ModuleName5,
		ModuleName6:        inputStruct.LicenseConfigData.ModuleName6,
		ModuleName7:        inputStruct.LicenseConfigData.ModuleName7,
		ModuleName8:        inputStruct.LicenseConfigData.ModuleName8,
		ModuleName9:        inputStruct.LicenseConfigData.ModuleName9,
		ModuleName10:       inputStruct.LicenseConfigData.ModuleName10,
		MaxOfflineDays:     inputStruct.LicenseConfigData.MaxOfflineDays,
		IsConcurrentUser:   inputStruct.LicenseConfigData.IsConcurrentUser,
		Component:          productComponent,
	}

	dataLicense, errorS = json.Marshal(licenseConfigReqGenerator)
	if errorS != nil {
		service.LogMessage(errorS.Error(), http.StatusInternalServerError)
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	productReqGenerator := in.GenerateDataProductConfiguration{
		ProductKey:   productKey,
		SignatureKey: inputStruct.SignatureKey,
		ClientSecret: inputStruct.ClientSecret,
		EncryptKey:   util.HashingPassword(inputStruct.LicenseConfigData.ClientID, inputStruct.ClientSecret),
		HardwareId:   inputStruct.Hwid,
	}

	dataProductReqGenerator, errorS = json.Marshal(productReqGenerator)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	argumentMap = make(map[string]string)
	argumentMap["args1"] = string(dataLicense)
	argumentMap["args2"] = string(dataProductReqGenerator)

	dataResponseEncrypt, err = service.GeneratorLicense("ProductEncrypt", argumentMap)
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

	result.ProductKey = validResponse.ProductKey
	result.ProductEncrypt = validResponse.ProductEncrypt
	result.ProductSignature = validResponse.ProductSignature

	err = errorModel.GenerateNonErrorModel()
	return
}
