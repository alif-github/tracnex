package ActivationLicenseService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Azure/go-autorest/autorest/date"
	"math"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
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
	"strings"
	"sync"
	"time"
)

func (input activationLicenseService) ActivateLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName       = "ActivateLicense"
		inputStruct    in.ActivationLicenseRequest
		additionalInfo interface{}
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateActivateLicense)
	if err.Error != nil {
		return
	}

	additionalInfo, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doActivateLicense, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- Function Additional
	})
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

func (input activationLicenseService) doActivateLicense(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName            = "doActivateLicense"
		inputStruct         = inputStructInterface.(in.ActivationLicenseRequest)
		isError             bool
		errorDetail         []out.ActivationLicenseErrorDetail
		licenseConfigs      []repository.LicenseConfigModel
		licenseDataActivate []in.ActivateLicenseDataRequest
		clientOnDB          repository.ClientCredentialModel
	)

	isError, errorDetail = input.validateRequestDetail(inputStruct, contextModel)
	if isError && len(errorDetail) > 0 {
		err = errorModel.GenerateFormatFieldError(input.FileName, funcName, constanta.ApplicationDetail)
		return
	}

	if contextModel.AuthAccessTokenModel.ClientID != inputStruct.ClientID {
		err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
		return
	}

	clientOnDB, err = dao.ClientCredentialDAO.GetClientCredentialForActivationLicense(serverconfig.ServerAttribute.DBConnection, repository.ClientCredentialModel{
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

	//--- Create JSON File
	licenseConfigs, errorDetail, isError = input.doGetLicenseConfig(inputStruct, contextModel)
	if isError {
		err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseGetLicenseConfigError, nil)
		output = errorDetail
		return
	}

	// Validate HWID
	for _, licenseItem := range licenseConfigs {
		if licenseItem.DeploymentMethod.String == "O" || licenseItem.DeploymentMethod.String == "M" {
			if util.IsStringEmpty(strings.TrimSpace(inputStruct.Hwid)) {
				err = errorModel.GenerateEmptyFieldError(input.FileName, funcName, constanta.HWID)
				errorDetail = append(errorDetail, out.ActivationLicenseErrorDetail{
					UniqueID1:       licenseItem.UniqueID1.String,
					UniqueID2:       licenseItem.UniqueID2.String,
					LicenseConfigID: licenseItem.ID.Int64,
					Message:         service.GetErrorMessage(err, *contextModel),
				})
				return
			}
		}
	}

	//--- Generate Product Key, Product Encrypt, dan Product Signature
	licenseDataActivate, errorDetail = input.doGenerateLicenseEncrypt(clientOnDB, inputStruct.Hwid, licenseConfigs, contextModel)
	if len(errorDetail) > 0 {
		err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseEncryptLicenseError, nil)
		output = errorDetail
		return
	}

	//--- Active License
	if len(licenseDataActivate) > 0 {
		//--- Insert Product License and User License
		output, dataAudit, err = input.DoActivateLicenseOnDB(tx, licenseDataActivate, contextModel, timeNow)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input activationLicenseService) DoActivateLicenseOnDB(tx *sql.Tx, inputStruct []in.ActivateLicenseDataRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName                              = "DoActivateLicenseOnDB"
		productLicenses                       []repository.ProductLicenseModel
		userLicenses                          []repository.UserLicenseModel
		response                              []out.LicenseResponse
		licenseConfigOnDB                     repository.LicenseConfigModel
		insertedIDsProductLicense             []int64
		insertedUserLicenseID, totalDataMoved int64
	)

	for _, item := range inputStruct {
		productLicenses = append(productLicenses, repository.ProductLicenseModel{
			LicenseConfigId:  sql.NullInt64{Int64: item.LicenseConfigID},
			ProductKey:       sql.NullString{String: item.ProductKey},
			ProductEncrypt:   sql.NullString{String: item.ProductEncrypt},
			ProductSignature: sql.NullString{String: item.ProductSignature},
			ClientId:         sql.NullString{String: item.ClientID},
			ClientSecret:     sql.NullString{String: item.ClientSecret},
			HWID:             sql.NullString{String: item.Hwid},
			ActivationDate:   sql.NullTime{Time: timeNow},
			LicenseStatus:    sql.NullInt32{Int32: 1},
			CreatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			CreatedAt:        sql.NullTime{Time: timeNow},
			CreatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		})

		response = append(response, out.LicenseResponse{
			ProductKey:       item.ProductKey,
			ProductEncrypt:   item.ProductEncrypt,
			ProductSignature: item.ProductSignature,
			UniqueID1:        item.UniqueID1,
			UniqueID2:        item.UniqueID2,
			ClientTypeID:     item.ClientTypeID,
		})

		if item.IsUserConcurrent == "N" {
			userLicenses = append(userLicenses, repository.UserLicenseModel{
				ParentCustomerId: sql.NullInt64{Int64: item.ParentCustomerID},
				CustomerId:       sql.NullInt64{Int64: item.CustomerID},
				SiteId:           sql.NullInt64{Int64: item.SiteID},
				InstallationId:   sql.NullInt64{Int64: item.InstallationID},
				ClientID:         sql.NullString{String: item.ClientID},
				UniqueId1:        sql.NullString{String: item.UniqueID1},
				UniqueId2:        sql.NullString{String: item.UniqueID2},
				ProductValidFrom: sql.NullTime{Time: item.ProductValidFrom.Time},
				ProductValidThru: sql.NullTime{Time: item.ProductValidThru.Time},
				TotalLicense:     sql.NullInt64{Int64: item.NumberOfUser},
				LicenseConfigId:  sql.NullInt64{Int64: item.LicenseConfigID},
				ProductKey:       sql.NullString{String: item.ProductKey},
				ProductSignature: sql.NullString{String: item.ProductSignature},
				ProductEncrypt:   sql.NullString{String: item.ProductEncrypt},
				CreatedAt:        sql.NullTime{Time: timeNow},
				CreatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
				CreatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
				UpdatedAt:        sql.NullTime{Time: timeNow},
				UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
				UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			})
		}
	}

	if len(productLicenses) > 0 {
		insertedIDsProductLicense, err = dao.ProductLicenseDAO.InsertBulkProductLicense(tx, productLicenses)
		if err.Error != nil {
			err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseCreateProductLicenseError, err.CausedBy)
			return
		}

		for _, id := range insertedIDsProductLicense {
			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.ProductLicenseDAO.TableName, id, 0)...)
		}
	}

	if len(userLicenses) > 0 {
		/* old flow activation
		insertedIDs, err = dao.UserLicenseDAO.InsertBulkUserLicenses(tx, userLicenses)
		if err.Error != nil {
			err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseCreateUserLicenseError, err.CausedBy)
			return
		}

		for _, id := range insertedIDs {
			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.UserLicenseDAO.TableName, id, 0)...)
		}
		*/

		for _, item := range userLicenses {
			// insert user license
			item.ProductLicenseID.Int64, err = dao.ProductLicenseDAO.GetIDByLicenseConfigID(tx, repository.ProductLicenseModel{LicenseConfigId: item.LicenseConfigId})
			if err.Error != nil {
				return
			}

			insertedUserLicenseID, err = dao.UserLicenseDAO.InsertUserLicense(tx, item)
			if err.Error != nil {
				return
			}

			// check is license config has user reg detail
			licenseConfigOnDB, err = dao.LicenseConfigDAO.CheckPreviousLicenseConfig(serverconfig.ServerAttribute.DBConnection, repository.LicenseConfigModel{
				ID: sql.NullInt64{Int64: item.LicenseConfigId.Int64},
			})
			if err.Error != nil {
				return
			}

			// move user reg
			if licenseConfigOnDB.IsHasPrevLicenseConfig.Bool {
				totalDataMoved, err = dao.UserRegistrationDetailDAO.MoveUserRegistrationDetail(tx, repository.UserRegistrationDetailModel{
					UserLicenseID:   sql.NullInt64{Int64: insertedUserLicenseID},
					LicenseConfigID: sql.NullInt64{Int64: item.LicenseConfigId.Int64},
				})
				if err.Error != nil {
					return
				}
			}

			// update total activated user license
			err = dao.UserLicenseDAO.UpdateTotalActivatedMovedUserLicense(tx, repository.UserLicenseModel{
				TotalActivated: sql.NullInt64{Int64: totalDataMoved},
				UpdatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
				UpdatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
				UpdatedAt:      sql.NullTime{Time: timeNow},
				ID:             sql.NullInt64{Int64: insertedUserLicenseID},
			})
			if err.Error != nil {
				return
			}

			// update total activated previous user license
			err = dao.UserLicenseDAO.UpdateTotalActivatedPrevUserLicense(tx, repository.UserLicenseModel{
				LicenseConfigId: licenseConfigOnDB.PrevLicenseConfigID,
				UpdatedBy:       sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
				UpdatedClient:   sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
				UpdatedAt:       sql.NullTime{Time: timeNow},
			})
			if err.Error != nil {
				return
			}
		}
	}

	output = response
	return
}

func (input activationLicenseService) doGenerateLicenseEncrypt(clientRequested repository.ClientCredentialModel, hwid string, licenseConfigs []repository.LicenseConfigModel, contextModel *applicationModel.ContextModel) (result []in.ActivateLicenseDataRequest, errorDetail []out.ActivationLicenseErrorDetail) {
	var (
		funcName     = "doGenerateLicenseEncrypt"
		isJsonError  = false
		encryptModel cryptoModel.EncryptLicenseRequestModel
		err          errorModel.ErrorModel
	)

	for _, config := range licenseConfigs {
		var (
			encryptResponse     cryptoModel.EncryptLicenseResponseModel
			activateLicenseData in.ActivateLicenseDataRequest
		)

		// Get JSON License
		encryptModel = cryptoModel.EncryptLicenseRequestModel{
			SignatureKey: clientRequested.SignatureKey.String,
			ClientSecret: clientRequested.ClientSecret.String,
			Hwid:         hwid,
		}

		encryptModel.LicenseConfigData, err = input.generateJSONLicense(config)

		if err.Error != nil {
			err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseEncryptLicenseError, nil)
			errorDetail = append(errorDetail, out.ActivationLicenseErrorDetail{
				UniqueID1:       config.UniqueID1.String,
				UniqueID2:       config.UniqueID2.String,
				LicenseConfigID: config.ID.Int64,
				Message:         service.GetErrorMessage(err, *contextModel),
			})
			isJsonError = true
			continue
		}

		if isJsonError {
			continue
		}

		// Encrypt License
		encryptResponse, err = input.GenerateLicenseEncrypt(encryptModel)
		if err.Error != nil {
			errorDetail = append(errorDetail, out.ActivationLicenseErrorDetail{
				UniqueID1:       config.UniqueID1.String,
				UniqueID2:       config.UniqueID2.String,
				LicenseConfigID: config.ID.Int64,
				Message:         service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		activateLicenseData = input.GetActivateLicenseData(config, encryptResponse)
		activateLicenseData.ClientSecret = clientRequested.ClientSecret.String
		activateLicenseData.Hwid = hwid

		result = append(result, activateLicenseData)
	}

	if len(licenseConfigs) == len(errorDetail) {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, nil)
		return
	}

	return
}

func (input activationLicenseService) GetActivateLicenseData(licenseConfig repository.LicenseConfigModel, encryptResponse cryptoModel.EncryptLicenseResponseModel) (result in.ActivateLicenseDataRequest) {
	return in.ActivateLicenseDataRequest{
		LicenseConfigID:  licenseConfig.ID.Int64,
		ProductKey:       encryptResponse.ProductKey,
		ProductEncrypt:   encryptResponse.ProductEncrypt,
		ProductSignature: encryptResponse.ProductSignature,
		IsUserConcurrent: licenseConfig.IsUserConcurrent.String,
		ClientID:         licenseConfig.ClientID.String,
		UniqueID1:        licenseConfig.UniqueID1.String,
		UniqueID2:        licenseConfig.UniqueID2.String,
		ParentCustomerID: licenseConfig.ParentCustomerID.Int64,
		NumberOfUser:     licenseConfig.NoOfUser.Int64,
		CustomerID:       licenseConfig.CustomerID.Int64,
		SiteID:           licenseConfig.SiteID.Int64,
		InstallationID:   licenseConfig.InstallationID.Int64,
		ProductValidFrom: date.Date{licenseConfig.ProductValidFrom.Time},
		ProductValidThru: date.Date{licenseConfig.ProductValidThru.Time},
		ClientTypeID:     licenseConfig.ClientTypeID.Int64,
	}
}

func (input activationLicenseService) GenerateLicenseEncrypt(inputStruct cryptoModel.EncryptLicenseRequestModel) (result cryptoModel.EncryptLicenseResponseModel, err errorModel.ErrorModel) {
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
		SignatureKey: inputStruct.SignatureKey,
		ClientSecret: inputStruct.ClientSecret,
		EncryptKey:   util2.HashingPassword(inputStruct.LicenseConfigData.ClientID, inputStruct.ClientSecret),
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

func (input activationLicenseService) doGetLicenseConfig(inputStruct in.ActivationLicenseRequest, contextModel *applicationModel.ContextModel) (payload []repository.LicenseConfigModel, errorDetail []out.ActivationLicenseErrorDetail, isError bool) {
	var wg sync.WaitGroup
	isError = false

	totalPage := math.Ceil(float64(len(inputStruct.DetailClient)) / float64(constanta.TotalDataProductLicensePerChannel))
	resultLicense := make(chan []repository.LicenseConfigModel, len(inputStruct.DetailClient))
	resultError := make(chan []out.ActivationLicenseErrorDetail, len(inputStruct.DetailClient))

	for i := 1; i <= int(totalPage); i++ {
		wg.Add(1)
		var uniqueIDs []in.UniqueIDClient

		offset := dao.CountOffset(i, constanta.TotalDataProductLicensePerChannel)
		until := offset + constanta.TotalDataProductLicensePerChannel

		if i == int(totalPage) {
			uniqueIDs = append(uniqueIDs, inputStruct.DetailClient[offset:]...)
		} else {
			uniqueIDs = append(uniqueIDs, inputStruct.DetailClient[offset:until]...)
		}

		go input.getJSONFromLicense(resultLicense, resultError, uniqueIDs, inputStruct.ClientID, contextModel, &wg)
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

	if len(errorDetail) > 0 {
		isError = true
	}

	return
}

func (input activationLicenseService) getJSONFromLicense(resultLicense chan []repository.LicenseConfigModel, resultError chan []out.ActivationLicenseErrorDetail, applicationDetail []in.UniqueIDClient, clientID string, contextModel *applicationModel.ContextModel, wg *sync.WaitGroup) {

	funcName := "getJSONFromLicense"
	var err errorModel.ErrorModel
	var tempResult []repository.LicenseConfigModel
	var tempErrorDetail []out.ActivationLicenseErrorDetail

	for _, license := range applicationDetail {
		var licenseOnDB []repository.LicenseConfigModel
		licenseOnDB, err = dao.LicenseConfigDAO.GetLicenseForJSONFile(serverconfig.ServerAttribute.DBConnection, repository.LicenseConfigModel{
			ClientID:  sql.NullString{String: clientID},
			UniqueID1: sql.NullString{String: license.UniqueID1},
			UniqueID2: sql.NullString{String: license.UniqueID2},
		})

		if err.Error != nil {
			service.LogMessage(fmt.Sprintf("Error get license on DB for %s - %s", license.UniqueID1, license.UniqueID2), http.StatusBadRequest)
			tempErrorDetail = append(tempErrorDetail, out.ActivationLicenseErrorDetail{
				UniqueID1: license.UniqueID1,
				UniqueID2: license.UniqueID2,
				Message:   service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		if len(licenseOnDB) < 1 {
			service.LogMessage(fmt.Sprintf("No license on DB for %s - %s", license.UniqueID1, license.UniqueID2), http.StatusBadRequest)
			err = errorModel.GenerateActivationLicenseError(input.FileName, funcName, constanta.ActivationLicenseGetLicenseConfigError, nil)
			tempErrorDetail = append(tempErrorDetail, out.ActivationLicenseErrorDetail{
				UniqueID1: license.UniqueID1,
				UniqueID2: license.UniqueID2,
				Message:   service.GetErrorMessage(err, *contextModel),
			})
			continue
		}

		tempResult = append(tempResult, licenseOnDB...)

	}

	defer func() {
		resultLicense <- tempResult
		resultError <- tempErrorDetail
		wg.Done()
	}()
}

func (input activationLicenseService) generateJSONLicense(licenseConfig repository.LicenseConfigModel) (result cryptoModel.JSONFileActivationLicenseModel, err errorModel.ErrorModel) {
	funcName := "generateJSONLicense"
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
		LicenseStatus:      1,
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
			service.LogMessage(service.GetErrorMessage(err, applicationModel.ContextModel{}), http.StatusBadRequest)
			return
		}
		result.ProductComponent = components
	}

	return
}

func (input activationLicenseService) validateRequestDetail(inputStruct in.ActivationLicenseRequest, contextModel *applicationModel.ContextModel) (isError bool, errorDetail []out.ActivationLicenseErrorDetail) {
	var err errorModel.ErrorModel
	for _, license := range inputStruct.DetailClient {
		err = inputStruct.ValidateMinMaxString(license.UniqueID1, constanta.UniqueID1, 1, 20)
		if err.Error != nil {
			isError = true
			errorDetail = append(errorDetail, out.ActivationLicenseErrorDetail{
				UniqueID1: license.UniqueID1,
				UniqueID2: license.UniqueID2,
				Message:   service.GetErrorMessage(err, *contextModel),
			})
		}

		if !util.IsStringEmpty(license.UniqueID2) {
			err = inputStruct.ValidateMinMaxString(license.UniqueID2, constanta.UniqueID2, 1, 20)
			if err.Error != nil {
				isError = true
				errorDetail = append(errorDetail, out.ActivationLicenseErrorDetail{
					UniqueID1: license.UniqueID1,
					UniqueID2: license.UniqueID2,
					Message:   service.GetErrorMessage(err, *contextModel),
				})
			}
		}
	}

	return
}

func (input activationLicenseService) validateActivateLicense(inputStruct *in.ActivationLicenseRequest) errorModel.ErrorModel {
	return inputStruct.ValidateActivateLicense()
}
