package ProductLicenseService

import (
	"database/sql"
	"encoding/json"
	"errors"
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
)

func (input productLicenseService) UpdateHWIDByPassProductLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		result      interface{}
		inputStruct in.UpdateLicenseHWIDByPassRequest
	)

	//--- Read Body and Validate
	inputStruct, err = input.readBodyAndValidateByPassProductLicense(request, contextModel, input.validateUpdateHWIDByPass)
	if err.Error != nil {
		return
	}

	//--- Main Function
	result, err = input.doUpdateHWIDByPass(inputStruct)
	if err.Error != nil {
		return
	}

	//--- Data Content
	if result != nil {
		output.Data.Content = result.([]out.LicenseResponse)
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) doUpdateHWIDByPass(inputStruct in.UpdateLicenseHWIDByPassRequest) (result interface{}, err errorModel.ErrorModel) {
	var (
		fileName         = "UpdateHWIDByPassService.go"
		funcName         = "doUpdateHWIDByPass"
		db               = serverconfig.ServerAttribute.DBConnection
		decryptedLicense []in.ValidationLicenseJSONFile
		clientCredential repository.ClientCredentialModel
		resultGenerator  []in.ActivateLicenseDataRequest
		configurationArr []cryptoModel.DecryptLicenseResponseModel
		isExist          bool
	)

	//--- Validation Credential
	if err = input.validationCredential(db, inputStruct); err.Error != nil {
		return
	}

	//--- Validation HWID to Database
	isExist, err = dao.WhiteListDevice.GetValidateWhiteListDevice(db, repository.WhiteListDeviceModel{Device: sql.NullString{String: inputStruct.HWIDInternal}})
	if err.Error != nil {
		return
	}

	if !isExist {
		err = errorModel.GenerateDataNotFoundWithParam(fileName, funcName, constanta.HardwareID)
		return
	}

	for i := 0; i < len(inputStruct.License); i++ {
		//--- Decrypt License
		var configuration cryptoModel.DecryptLicenseResponseModel
		configuration, err = input.decryptLicense(cryptoModel.DecryptLicenseRequestModel{
			ProductSignature: inputStruct.License[i].ProductSignature,
			ClientId:         inputStruct.ClientID,
			ClientSecret:     inputStruct.ClientSecret,
			HardwareId:       inputStruct.HWID, //--- Use old HWID
			ProductKey:       inputStruct.License[i].ProductKey,
		}, inputStruct.SignatureKey, inputStruct.License[i].ProductEncrypt)
		if err.Error != nil {
			return
		}

		configurationArr = append(configurationArr, configuration)
	}

	//--- Generate New License
	decryptedLicense, clientCredential = input.preparingNewLicenseModel(configurationArr, inputStruct)
	resultGenerator, err = input.generateNewLicense(decryptedLicense, clientCredential, in.UpdateLicenseHWIDRequest{Hwid: inputStruct.HWIDInternal})
	if err.Error != nil {
		return
	}

	//--- Return Result Generator
	if len(resultGenerator) < 1 {
		err = errorModel.GenerateUnknownError(fileName, funcName, errors.New("result generator empty"))
		return
	}

	//--- Output
	var response []out.LicenseResponse
	for _, itemResultGenerator := range resultGenerator {
		response = append(response, out.LicenseResponse{
			ProductKey:       itemResultGenerator.ProductKey,
			ProductEncrypt:   itemResultGenerator.ProductEncrypt,
			ProductSignature: itemResultGenerator.ProductSignature,
			ClientTypeID:     itemResultGenerator.ClientTypeID,
			UniqueID1:        itemResultGenerator.UniqueID1,
			UniqueID2:        itemResultGenerator.UniqueID2,
		})
	}

	result = response
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) preparingNewLicenseModel(configuration []cryptoModel.DecryptLicenseResponseModel, inputStruct in.UpdateLicenseHWIDByPassRequest) (decryptedLicense []in.ValidationLicenseJSONFile, clientCredential repository.ClientCredentialModel) {

	//--- JSON File Model
	for i := 0; i < len(configuration); i++ {
		jsonFileModel := cryptoModel.JSONFileActivationLicenseModel{
			InstallationID:      configuration[i].Configuration.InstallationID,
			ClientID:            inputStruct.ClientIDInternal,
			ProductID:           configuration[i].Configuration.ProductID,
			LicenseVariantName:  configuration[i].Configuration.LicenseVariantName,
			LicenseTypeName:     configuration[i].Configuration.LicenseTypeName,
			DeploymentMethod:    configuration[i].Configuration.DeploymentMethod,
			NumberOfUser:        configuration[i].Configuration.NumberOfUser,
			UniqueID1:           configuration[i].Configuration.UniqueID1,
			UniqueID2:           configuration[i].Configuration.UniqueID2,
			ProductValidFromStr: configuration[i].Configuration.ProductValidFromStr,
			ProductValidThruStr: configuration[i].Configuration.ProductValidThruStr,
			LicenseStatus:       configuration[i].Configuration.LicenseStatus,
			ModuleName1:         configuration[i].Configuration.ModuleName1,
			ModuleName2:         configuration[i].Configuration.ModuleName2,
			ModuleName3:         configuration[i].Configuration.ModuleName3,
			ModuleName4:         configuration[i].Configuration.ModuleName4,
			ModuleName5:         configuration[i].Configuration.ModuleName5,
			ModuleName6:         configuration[i].Configuration.ModuleName6,
			ModuleName7:         configuration[i].Configuration.ModuleName7,
			ModuleName8:         configuration[i].Configuration.ModuleName8,
			ModuleName9:         configuration[i].Configuration.ModuleName9,
			ModuleName10:        configuration[i].Configuration.ModuleName10,
			MaxOfflineDays:      configuration[i].Configuration.MaxOfflineDays,
			IsConcurrentUser:    configuration[i].Configuration.IsConcurrentUser,
			ProductComponent:    configuration[i].Configuration.ProductComponent,
		}

		//--- Decrypt License
		decryptedLicense = append(decryptedLicense, in.ValidationLicenseJSONFile{
			JSONFileActivationLicenseModel: jsonFileModel,
			ProductKey:                     configuration[i].ProductKey,
			ClientTypeID:                   inputStruct.ClientTypeID,
		})
	}

	//--- Client Credential
	clientCredential = repository.ClientCredentialModel{
		ClientID:     sql.NullString{String: inputStruct.ClientIDInternal},
		ClientSecret: sql.NullString{String: inputStruct.ClientSecretInternal},
		SignatureKey: sql.NullString{String: inputStruct.SignatureKeyInternal},
	}

	return
}

func (input productLicenseService) readBodyAndValidateByPassProductLicense(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.UpdateLicenseHWIDByPassRequest) errorModel.ErrorModel) (inputStruct in.UpdateLicenseHWIDByPassRequest, err errorModel.ErrorModel) {
	var (
		fileName   = "UpdateHWIDByPassService.go"
		funcName   = "readBodyAndValidateByPassProductLicense"
		stringBody string
	)

	//--- Read Body
	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	//--- Unmarshal Body
	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	//--- General Validation
	if err = validation(&inputStruct); err.Error != nil {
		return
	}

	//--- Client ID Validation
	if inputStruct.ClientIDInternal != contextModel.AuthAccessTokenModel.ClientID {
		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		return
	}

	//if inputStruct.ClientTypeID != constanta.ResourceND6ID {
	//	err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
	//	return
	//}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) validationCredential(db *sql.DB, inputStruct in.UpdateLicenseHWIDByPassRequest) (err errorModel.ErrorModel) {
	var (
		fileName         = "UpdateHWIDByPassService.go"
		funcName         = "validationCredential"
		credentialResult repository.ClientCredentialModel
	)

	//--- Get Client Credential
	credentialResult, err = dao.ClientCredentialDAO.GetClientCredentialByClientID(db, repository.ClientCredentialModel{ClientID: sql.NullString{String: inputStruct.ClientIDInternal}})
	if err.Error != nil {
		return
	}

	if credentialResult.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientID)
		return
	}

	if credentialResult.ClientSecret.String != inputStruct.ClientSecretInternal {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientSecret)
		return
	}

	if credentialResult.SignatureKey.String != inputStruct.SignatureKeyInternal {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.SignatureKey)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) validateUpdateHWIDByPass(inputStruct *in.UpdateLicenseHWIDByPassRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateHWIDByPass()
}
