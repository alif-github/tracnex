package ClientService

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/ClientRegistrationLogService"
	"nexsoft.co.id/nextrac2/service/HitAddExternalResource"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
	"time"
)

type clientService struct {
	service.RegistrationPrepared
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var ClientService = clientService{}.New()

func (input clientService) New() (output clientService) {
	output.FileName = "ClientService.go"
	output.ValidLimit = []int{10, 20, 50, 100}
	output.ValidOrderBy = []string{"id"}
	output.IdResourceAllowed = []int64{1}
	return
}

func (input clientService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(clientRequest *in.ClientRequest) errorModel.ErrorModel) (inputStruct in.ClientRequest, err errorModel.ErrorModel) {
	var (
		funcName            = "readBodyAndValidate"
		stringBody          string
		attributeRequestStr *in.ClientRequestAttributeRequest
		isAllowed           bool
	)

	//---------- Read Body Request
	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	//---------- Unmarshal String Body to Main struct
	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		input.LogError(err, *contextModel)
		return
	}

	//---------- Init, Important request must be exists
	if inputStruct.ClientTypeID == 0 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(input.FileName, funcName, constanta.ClientTypeID)
		return
	}

	//---------- Check to DB, client type exist on table ?
	preparedError := errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientTypeID)
	err = input.CheckIsClientTypeExist(inputStruct.ClientTypeID, preparedError)
	if err.Error != nil {
		return
	}

	//---------- Unmarshal String Body to Attribute request
	errorS = json.Unmarshal([]byte(stringBody), &attributeRequestStr)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		input.LogError(err, *contextModel)
		return
	}

	//---------- Check client type allowing
	for _, idResourceItem := range input.IdResourceAllowed {
		if idResourceItem == inputStruct.ClientTypeID {
			isAllowed = true
			break
		}
	}

	//---------- Is not allowed, then forbidden to access
	if !isAllowed {
		err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
		return
	}

	//---------- Adding to main struct -> attribute request
	b, errorS := json.Marshal(attributeRequestStr)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		input.LogError(err, *contextModel)
		return
	}

	inputStruct.Body = string(b)

	//---------- Main validation
	err = validation(&inputStruct)

	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
	logModel.Message = "--->> Sukses Read Body Validate!"
	logModel.Status = 200
	util.LogInfo(logModel.ToLoggerObject())
	return
}

func (input clientService) addClientToAuthenticationServer(inputStruct in.ClientRequest, contextModel *applicationModel.ContextModel) (registerClientContent authentication_response.RegisterClientContent, err errorModel.ErrorModel) {
	funcName := "addClientToAuthenticationServer"

	var data authentication_response.RegisterClientAuthenticationResponse
	var additionalInfo []model.AdditionalInformation

	if inputStruct.SocketID != "" {
		additionalInfo = append(additionalInfo, model.AdditionalInformation{
			Key:   "socket_id",
			Value: inputStruct.SocketID,
		})
	}

	clientStruct := authentication_request.ClientAuthentication{
		ResourceID:           config.ApplicationConfiguration.GetServerResourceID(),
		Scope:                constanta.ScopeClient,
		GrantTypes:           constanta.GrantTypes,
		AccessTokenValidity:  constanta.AccessTokenValidity,
		RefreshTokenValidity: constanta.RefreshTokenValidity,
		MaxAuthFail:          constanta.MaxAuthFail,
		Locale:               constanta.IndonesianLanguage,
	}

	if inputStruct.SocketID != "" {
		clientStruct.ClientInformation = additionalInfo
	}

	internalToken := resource_common_service.GenerateInternalToken("auth", 0, contextModel.AuthAccessTokenModel.ClientID, config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	registerClientUrl := authenticationServer.Host + authenticationServer.PathRedirect.InternalClient.CrudClient

	statusCode, bodyResult, errorS := common.HitRegisterClientAuthenticationServer(internalToken, registerClientUrl, clientStruct, contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	readError := json.Unmarshal([]byte(bodyResult), &data)
	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, readError)
		return
	}
	registerClientContent.ClientID = data.Nexsoft.Payload.Data.Content.ClientID
	registerClientContent.SignatureKey = data.Nexsoft.Payload.Data.Content.SignatureKey
	registerClientContent.ClientSecret = data.Nexsoft.Payload.Data.Content.ClientSecret

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(data.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, data.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}

func (input clientService) addResourceClient(clientCredentialResult authentication_response.RegisterClientContent, contextModel *applicationModel.ContextModel) (addResourceResponse authentication_response.AddResourceClientContent, err errorModel.ErrorModel) {
	funcName := "addClientToAuthenticationServer"

	var data authentication_response.AddResourceClientAuthenticationResponse

	clientStruct := authentication_request.AddResourceClient{
		ClientID:     clientCredentialResult.ClientID,
		ClientSecret: clientCredentialResult.ClientSecret,
		ResourceID:   config.ApplicationConfiguration.GetServerResourceID(),
	}

	internalToken := resource_common_service.GenerateInternalToken("auth", 0, contextModel.AuthAccessTokenModel.ClientID, config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	registerClientUrl := authenticationServer.Host + authenticationServer.PathRedirect.InternalClient.AddResourceAdmin

	statusCode, _, errorS := common.HitAddResourceClientAuthenticationServer(internalToken, registerClientUrl, clientStruct, contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(data.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, data.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}

func (input clientService) addResourceExternal(inputStructInterface interface{}, clientCredential authentication_response.RegisterClientContent, contextModel *applicationModel.ContextModel) (successResource []string, failedResource []string, err errorModel.ErrorModel) {
	funcName := "addResourceExternal"
	inputStruct := inputStructInterface.(in.ClientRequest)
	var errorS []errorModel.ErrorModel
	var errorTemp errorModel.ErrorModel

	err = HitAddExternalResource.HitAddExternalResource.AddResourceNexcloud(in.AddResourceNexcloud{
		FirstName: inputStruct.ClientName,
		LastName:  constanta.Nexdistribution,
		ClientID:  clientCredential.ClientID,
	}, contextModel)

	if err.Error != nil {
		_, errorTemp = service.NewErrorAddResource(serverconfig.ServerAttribute.ClientBundle, contextModel, "DETAIL_ERROR_FAILED_RESOURCE_NEXCLOUD_MESSAGE", input.FileName, funcName)
		errorS = append(errorS, errorTemp)
		failedResource = append(failedResource, constanta.NexCloudResourceID)
	} else {
		successResource = append(successResource, constanta.NexCloudResourceID)
	}

	//err = HitAddExternalResource.HitAddExternalResource.AddResourceNexdrive(in.AddResourceNexdrive{
	//	FirstName: 	inputStruct.ClientName,
	//	LastName: 	constanta.ND6,
	//	ClientID: 	clientCredential.ClientID,
	//}, contextModel)
	//
	//if err.Error != nil {
	//	_, errorTemp = service.NewErrorAddResource(serverconfig.ServerAttribute.ClientBundle, contextModel, "DETAIL_ERROR_FAILED_RESOURCE_NEXDRIVE_MESSAGE", input.FileName, funcName)
	//	errorS = append(errorS, errorTemp)
	//	failedResource = append(failedResource, constanta.NexdriveResourceID)
	//} else {
	//	successResource = append(successResource, constanta.NexdriveResourceID)
	//}

	if len(errorS) == 2 {
		_, err = service.NewErrorAddResource(serverconfig.ServerAttribute.ClientBundle, contextModel, "DETAIL_ERROR_FAILED_RESOURCE_MESSAGE", input.FileName, funcName)
		return
	} else if len(errorS) == 1 {
		for _, errElm := range errorS {
			err = errElm
		}
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) DoCheckClientMappingForRegister(tx *sql.Tx, inputStructInterface interface{}) (output interface{}, clientCredential repository.ClientCredentialModel, err errorModel.ErrorModel) {
	var (
		inputStruct        = inputStructInterface.(in.ClientRequest)
		modelClientCdr     []repository.ClientCredentialModel
		clientCdr          []repository.ClientCredentialModel
		modelClientMapping []repository.ClientMappingModel
		result             []repository.ClientMappingModel
	)

	//---------- Input to client mapping model
	modelClientMapping = input.inputToClientMappingModel(inputStruct, false)

	//---------- Check to client mapping and get client ID
	result, err = dao.ClientMappingDAO.CheckClientMapping(tx, modelClientMapping, false)
	if err.Error != nil {
		return
	}

	//---------- Empty result would be return old struct
	if len(result) < 1 {
		return inputStruct, repository.ClientCredentialModel{}, errorModel.GenerateNonErrorModel()
	}

	//---------- Set client id result from client mapping for get credential
	for _, resultValue := range result {
		modelClientCdr = append(modelClientCdr, repository.ClientCredentialModel{ClientID: sql.NullString{String: resultValue.ClientID.String}})
	}

	//---------- Get client credential
	clientCdr, err = dao.ClientCredentialDAO.GetClientCredential(serverconfig.ServerAttribute.DBConnection, modelClientCdr)
	if err.Error != nil {
		return
	}

	//---------- Check the client ID if different must return error
	for index, clientCdrValue := range clientCdr {
		if index == 0 {
			clientCredential = repository.ClientCredentialModel{
				ClientID:     sql.NullString{String: clientCdrValue.ClientID.String},
				ClientSecret: sql.NullString{String: clientCdrValue.ClientSecret.String},
				SignatureKey: sql.NullString{String: clientCdrValue.SignatureKey.String},
			}
			continue
		}
		if clientCdrValue.ClientID.String != clientCredential.ClientID.String {
			err = input.prepareErrorRegisteredClient(result, true)
			return
		}
	}

	output = input.RemoveDataRegistered(inputStructInterface.(in.ClientRequest), result)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) DoCheckClientMapping(tx *sql.Tx, inputStructInterface interface{}, isCheckClientID bool) (output interface{}, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.ClientRequest)
	var modelClientMapping []repository.ClientMappingModel
	var newInputStruct in.ClientRequest

	//---------- Input to client mapping model
	modelClientMapping = input.inputToClientMappingModel(inputStruct, isCheckClientID)

	//---------- Check registered data in client mapping
	result, err := dao.ClientMappingDAO.CheckClientMapping(tx, modelClientMapping, isCheckClientID)
	if err.Error != nil {
		return
	}

	//---------- Decision for check client id or not
	if !isCheckClientID {
		if len(result) != 0 {
			err = input.prepareErrorRegisteredClient(result, false)
			return
		} else {
			output = inputStruct
		}
	} else {
		if len(result) != 0 {
			var newCompanyData []in.CompanyData
			var newBranchData []in.BranchData

			for _, resultElm := range result {
				newBranchData = append(newBranchData, in.BranchData{
					ClientID:    resultElm.ClientID.String,
					BranchID:    resultElm.BranchID.String,
					ClientAlias: resultElm.ClientAlias.String,
				})
				newCompanyData = append(newCompanyData, in.CompanyData{
					CompanyID:  resultElm.CompanyID.String,
					BranchData: newBranchData,
				})
			}
			newInputStruct = in.ClientRequest{
				ClientTypeID: inputStruct.ClientTypeID,
				CompanyData:  newCompanyData,
			}

			output = newInputStruct
		} else {
			output = in.ClientRequest{}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) RemoveDataRegistered(inputStruct in.ClientRequest, result []repository.ClientMappingModel) (output interface{}) {

	var (
		clientStruct    in.ClientRequest //this for template when data one by one remove
		urlStruct       in.URLRequest
		bodyReqStruct   in.BodyRequest
		clientModelTemp []repository.ClientMappingForRemoveModel
	)

	for _, companyDataElm := range inputStruct.CompanyData {
		for _, branchDataElm := range companyDataElm.BranchData {
			clientModelTemp = append(clientModelTemp, repository.ClientMappingForRemoveModel{
				CompanyID:   sql.NullString{String: companyDataElm.CompanyID},
				BranchID:    sql.NullString{String: branchDataElm.BranchID},
				BranchName:  sql.NullString{String: branchDataElm.BranchName},
				ClientAlias: sql.NullString{String: branchDataElm.ClientAlias},
			})
		}
	}

	if clientModelTemp == nil {
		output = in.ClientRequest{}
		return
	}

	//-----------------------------------------Remove process-----------------------------------------

	for i := 0; i < len(clientModelTemp); i++ {
		for j := 0; j < len(result); j++ {
			if (clientModelTemp[i].CompanyID.String == result[j].CompanyID.String) &&
				(clientModelTemp[i].BranchID.String == result[j].BranchID.String) {

				copy(clientModelTemp[i:], clientModelTemp[i+1:])
				clientModelTemp[len(clientModelTemp)-1] = repository.ClientMappingForRemoveModel{}
				clientModelTemp = clientModelTemp[:len(clientModelTemp)-1]

				i = -1
				break
			}
		}
	}

	//---------- Insert to new struct
	var companyData []in.CompanyData
	for _, clientModelTempElm := range clientModelTemp {

		var branchData []in.BranchData

		branchData = append(branchData, in.BranchData{
			BranchID:    clientModelTempElm.BranchID.String,
			BranchName:  clientModelTempElm.BranchName.String,
			ClientAlias: clientModelTempElm.ClientAlias.String,
		})

		companyData = append(companyData, in.CompanyData{
			CompanyID:  clientModelTempElm.CompanyID.String,
			BranchData: branchData,
		})
	}

	urlStruct = in.URLRequest{
		UrlAPI: inputStruct.UrlAPI,
	}

	bodyReqStruct = in.BodyRequest{
		Body: inputStruct.Body,
	}

	clientStruct = in.ClientRequest{
		URLRequest:   urlStruct,
		BodyRequest:  bodyReqStruct,
		ClientTypeID: inputStruct.ClientTypeID,
		ClientName:   inputStruct.ClientName,
		SocketID:     inputStruct.SocketID,
		CompanyData:  companyData,
	}

	output = clientStruct
	return
}

func (input clientService) DoCheckCustomerExist(tx *sql.Tx, inputStructInterface interface{}, errors errorModel.ErrorModel,
	_ *applicationModel.ContextModel, timeNow time.Time, isCheckExp bool) (output interface{}, err errorModel.ErrorModel) {
	var (
		funcName      = "DoCheckCustomerExist"
		inputStruct   = inputStructInterface.(in.ClientRequest)
		modelCustomer []repository.CustomerListModel
	)

	//------- Copy to model
	for _, companyDataElm := range inputStruct.CompanyData {
		for _, branchDataElm := range companyDataElm.BranchData {
			modelCustomer = append(modelCustomer, repository.CustomerListModel{
				CompanyID:  sql.NullString{String: companyDataElm.CompanyID},
				BranchID:   sql.NullString{String: branchDataElm.BranchID},
				BranchName: sql.NullString{String: branchDataElm.BranchName},
			})
		}
	}

	var (
		result               []repository.CustomerListModel
		errorTemp            = errorModel.GenerateErrorCustomActivationCode(input.FileName, funcName, fmt.Sprintf(`Company ID atau branch ID tidak valid`))
		errorBulkData        []out.CompanyBranchErrorBulk
		validCounting        = 0
		dataNotFoundCounting = 0
	)

	for _, modelCustomerItem := range modelCustomer {
		var resultTemp []repository.CustomerListModel

		//------- Check customer to DB
		resultTemp, err = dao.CustomerListDAO.CheckCustomerByProductName(tx, modelCustomerItem)
		if err.Error != nil {
			return
		}

		//------- If empty then registration because may input wrong company or branch id
		if resultTemp == nil {
			var custInstallationOnDB []repository.CustomerInstallationForConfig
			custInstallationOnDB, err = dao.CustomerInstallationDAO.GetCustomerInstallationByUniqueID(serverconfig.ServerAttribute.DBConnection, []repository.CustomerInstallationDetail{
				{
					UniqueID1: modelCustomerItem.CompanyID,
					UniqueID2: modelCustomerItem.BranchID,
				},
			}, false)
			if err.Error != nil {
				return
			}

			if len(custInstallationOnDB) == 0 {
				dataNotFoundCounting++
				errorBulkData = append(errorBulkData, out.CompanyBranchErrorBulk{
					CompanyID:    modelCustomerItem.CompanyID.String,
					BranchID:     modelCustomerItem.BranchID.String,
					ErrorMessage: errorTemp.ErrorParameter[0].ErrorParameterValue,
				})
			} else {
				if len(custInstallationOnDB) > 0 {
					validCounting++

					//--- Report ND6 Old WAR (24/05/2023)
					if modelCustomerItem.BranchName.String == "" {
						modelCustomerItem.BranchName.String = modelCustomerItem.BranchID.String
					}

					result = append(result, repository.CustomerListModel{
						CompanyID:   sql.NullString{String: modelCustomerItem.CompanyID.String},
						BranchID:    sql.NullString{String: modelCustomerItem.BranchID.String},
						CompanyName: sql.NullString{String: modelCustomerItem.BranchName.String},
					})
				}
			}
		}

		//------- Append to result
		if len(resultTemp) > 0 && isCheckExp {
			for index, resultTempItem := range resultTemp {

				//------- Last, direct to append
				if (len(resultTemp)-index == 1) && (resultTempItem.ExpDate.Time.After(timeNow) || resultTempItem.ExpDate.Time.Equal(timeNow)) {
					validCounting++
					result = append(result, resultTempItem)
					break
				} else if len(resultTemp)-index == 1 {
					result = append(result, resultTempItem)
					break
				}

				//------- Take it the first data after timeNow
				if resultTempItem.ExpDate.Time.After(timeNow) || resultTempItem.ExpDate.Time.Equal(timeNow) {
					validCounting++
					result = append(result, resultTempItem)
					break
				}
			}
		} else if len(resultTemp) > 0 && !isCheckExp {
			for _, resultItem := range resultTemp {
				//------- Append it the first data
				validCounting++
				result = append(result, resultItem)
				break
			}
		}
	}

	if dataNotFoundCounting > 0 {
		//detailTemp1 := util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ClientBundle, "AMOUNT_VALID_DATA_MESSAGE", contextModel.AuthAccessTokenModel.Locale, nil)
		//detailTemp2 := util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ClientBundle, "DETAIL_BRANCH_COMPANY_NOT_FOUND_MESSAGE", contextModel.AuthAccessTokenModel.Locale, nil)
		//detail := detailTemp1 + strconv.Itoa(len(result)) + detailTemp2
		err = errorModel.GenerateInvalidAddBranch(input.FileName, funcName, errorBulkData)
		return
	}

	//------- If empty then invalid registration
	if result == nil {
		err = errors
		return
	} else if validCounting < 1 {
		err = errors
		return
	}

	//------- Get company name for client alias
	for indexCompany, companyDataElm2 := range inputStruct.CompanyData {
		for indexBranch, branchDataElm2 := range companyDataElm2.BranchData {
			for _, resultElm := range result {
				if companyDataElm2.CompanyID == resultElm.CompanyID.String && branchDataElm2.BranchID == resultElm.BranchID.String {
					inputStruct.CompanyData[indexCompany].BranchData[indexBranch].ClientAlias = resultElm.CompanyName.String
				}
			}
		}
	}

	output = inputStruct
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) isND6Registration(inputStructInterface interface{}, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	funcName := "isND6Registration"
	inputStruct := inputStructInterface.(in.ClientRequest)

	result, err := dao.ClientTypeDAO.CheckClientType(serverconfig.ServerAttribute.DBConnection, &repository.ClientTypeModel{
		ClientType: sql.NullString{String: constanta.ND6},
	})
	if err.Error != nil {
		return
	}
	if result.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientTypeID)
		return
	}

	if result.ID.Int64 != inputStruct.ClientTypeID {
		detail := util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ClientBundle, "DETAIL_INVALID_CLIENTTYPE_ND6_MESSAGE", contextModel.AuthAccessTokenModel.Locale, nil)
		err = errorModel.GenerateInvalidRegistrationClient(input.FileName, funcName, []string{detail})
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) clientErrorHandle(inputStruct in.ClientRequest, contextModel *applicationModel.ContextModel,
	resourceIDList []out.ResourceList, successResourceID []string, failedResourceID []string, registerClientContent authentication_response.RegisterClientContent,
	timeNow time.Time, errorS errorModel.ErrorModel) (result interface{}, err errorModel.ErrorModel) {

	var message string
	var resourceList []string

	for _, resourceIDElm := range successResourceID {
		resourceIDList = append(resourceIDList, out.ResourceList{
			ResourceID: resourceIDElm,
			Status:     "OK",
		})
	}

	for _, resourceIDElm := range failedResourceID {
		resourceIDList = append(resourceIDList, out.ResourceList{
			ResourceID: resourceIDElm,
			Status:     "FAIL",
		})
	}

	newRegisterClientContent := input.addResourceInfo(registerClientContent, resourceIDList)
	result, message, err = service.CustomFailedResponsePayload(newRegisterClientContent, errorS, contextModel)

	for _, resourceElm := range resourceIDList {
		if resourceElm.Status != "FAIL" {
			resourceList = append(resourceList, resourceElm.ResourceID)
		}
	}

	resourceListStr := strings.Join(resourceList, " ")

	_, err = ClientRegistrationLogService.ClientRegistrationLogService.InsertClientRegistrationLog(in.ClientRegisterLogRequest{
		ClientID:          registerClientContent.ClientID,
		ClientTypeID:      inputStruct.ClientTypeID,
		AttributeRequest:  inputStruct.Body,
		SuccessStatusAuth: true,
		Resource:          resourceListStr,
		MessageAuth:       message,
		Details:           errorS.AdditionalInformation[0],
		Code:              errorS.Error.Error(),
		RequestTimestamp:  timeNow,
		RequestCount:      1,
	}, contextModel)

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) customResponseAndAddResourceInfo(registerClientContent authentication_response.RegisterClientContent,
	contextModel *applicationModel.ContextModel, resourceIDList []out.ResourceList, messageID string) (successMessage string, code string, result interface{}) {

	newRegisterClientContent := input.addResourceInfo(registerClientContent, resourceIDList)
	successMessage = GenerateI18NMessage(messageID, contextModel.AuthAccessTokenModel.Locale)
	code, result = service.CustomSuccessResponsePayload(newRegisterClientContent, successMessage, contextModel)
	return
}

func (input clientService) clientSuccessHandle(registerClientContent authentication_response.RegisterClientContent,
	contextModel *applicationModel.ContextModel, inputStruct in.ClientRequest,
	resourceIDList []out.ResourceList, messageID string, timeNow time.Time) (result interface{}, err errorModel.ErrorModel) {

	var code string
	var resourceList []string
	var successMessage string

	successMessage, code, result = input.customResponseAndAddResourceInfo(registerClientContent, contextModel, resourceIDList, messageID)

	for _, resourceElm := range resourceIDList {
		if resourceElm.Status != "FAIL" {
			resourceList = append(resourceList, resourceElm.ResourceID)
		}
	}

	resourceListStr := strings.Join(resourceList, " ")

	_, err = ClientRegistrationLogService.ClientRegistrationLogService.InsertClientRegistrationLog(in.ClientRegisterLogRequest{
		ClientID:              registerClientContent.ClientID,
		ClientTypeID:          inputStruct.ClientTypeID,
		AttributeRequest:      inputStruct.Body,
		SuccessStatusAuth:     true,
		SuccessStatusNexcloud: true,
		Resource:              resourceListStr,
		MessageAuth:           successMessage,
		Code:                  code,
		RequestTimestamp:      timeNow,
		RequestCount:          1,
	}, contextModel)

	return
}

func (input clientService) prepareErrorRegisteredClient(clientModel []repository.ClientMappingModel, isRegister bool) (err errorModel.ErrorModel) {
	funcName := "prepareErrorRegisteredClient"
	var attributeError in.AttributeRequestErrorRegisteredClient
	var detail []string

	for _, resultItem := range clientModel {
		attributeError = in.AttributeRequestErrorRegisteredClient{
			CompanyID: resultItem.CompanyID.String,
			BranchID:  resultItem.BranchID.String,
		}
		b, _ := json.Marshal(attributeError)
		detail = append(detail, string(b))
	}

	if isRegister {
		err = errorModel.GenerateDataUsedRegisterClientDiffClientIDError(input.FileName, funcName, detail)
	} else {
		err = errorModel.GenerateDataUsedRegisterClientError(input.FileName, funcName, detail)
	}
	return
}

func (input clientService) inputToClientMappingModel(inputStruct in.ClientRequest, isCheckClientID bool) (modelClientMapping []repository.ClientMappingModel) {
	for _, companyDataElm := range inputStruct.CompanyData {
		for _, branchDataElm := range companyDataElm.BranchData {
			if isCheckClientID {
				modelClientMapping = append(modelClientMapping, repository.ClientMappingModel{
					CompanyID: sql.NullString{String: companyDataElm.CompanyID},
					BranchID:  sql.NullString{String: branchDataElm.BranchID},
					ClientID:  sql.NullString{String: branchDataElm.ClientID},
				})
			} else {
				modelClientMapping = append(modelClientMapping, repository.ClientMappingModel{
					CompanyID: sql.NullString{String: companyDataElm.CompanyID},
					BranchID:  sql.NullString{String: branchDataElm.BranchID},
				})
			}
		}
	}
	return
}

func (input clientService) DoCheckClientMappingSpecialInsertNewBranch(tx *sql.Tx, inputStructInterface interface{}, isCheckClientID bool, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		inputStruct        = inputStructInterface.(in.ClientRequest)
		modelClientMapping []repository.ClientMappingModel
		newInputStruct     in.ClientRequest
		result             []repository.ClientMappingModel
		diffClientID       bool
	)

	//---------- Input to client mapping model
	modelClientMapping = input.inputToClientMappingModel(inputStruct, isCheckClientID)

	//---------- Check registered data in client mapping
	result, err = dao.ClientMappingDAO.CheckClientMapping(tx, modelClientMapping, isCheckClientID)
	if err.Error != nil {
		return
	}

	//---------- Decision for check client id or not
	if !isCheckClientID {
		if len(modelClientMapping) != len(result) {

			for _, resultValueA := range result {
				if resultValueA.ClientID.String != contextModel.AuthAccessTokenModel.ClientID {
					err = errorModel.GenerateDataUsedRegisterClientError(input.FileName, "DoCheckClientMappingSpecialInsertNewBranch", nil)
					output = in.ClientRequest{}
					diffClientID = true
					break
				}
			}

			if diffClientID {
				return
			}

			resultFinal := input.RemoveDataRegistered(inputStruct, result)
			output = resultFinal.(in.ClientRequest)
			err = errorModel.GenerateNonErrorModel()
			return
		} else {
			for _, resultValueB := range result {
				if resultValueB.ClientID.String != contextModel.AuthAccessTokenModel.ClientID {
					err = errorModel.GenerateDataUsedRegisterClientError(input.FileName, "DoCheckClientMappingSpecialInsertNewBranch", nil)
					output = in.ClientRequest{}
					diffClientID = true
					break
				}
			}

			if diffClientID {
				return
			}

			output = in.ClientRequest{}
			err = errorModel.GenerateNonErrorModel()
			return
		}
	} else {
		if len(result) != 0 {
			var newCompanyData []in.CompanyData
			var newBranchData []in.BranchData

			for _, resultElm := range result {
				newBranchData = append(newBranchData, in.BranchData{
					ClientID:    resultElm.ClientID.String,
					BranchID:    resultElm.BranchID.String,
					ClientAlias: resultElm.ClientAlias.String,
				})
				newCompanyData = append(newCompanyData, in.CompanyData{
					CompanyID:  resultElm.CompanyID.String,
					BranchData: newBranchData,
				})
			}
			newInputStruct = in.ClientRequest{
				ClientTypeID: inputStruct.ClientTypeID,
				CompanyData:  newCompanyData,
			}

			output = newInputStruct
		} else {
			output = in.ClientRequest{}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) DoCheckCustomerExistForRegistrationClientID(tx *sql.Tx, inputStructInterface interface{}, errors errorModel.ErrorModel,
	timeNow time.Time, isCheckExp bool) (output interface{}, err errorModel.ErrorModel) {
	var (
		funcName      = "DoCheckCustomerExistForRegistrationClientID"
		inputStruct   = inputStructInterface.(in.ClientRequest)
		modelCustomer []repository.CustomerListModel
	)

	//------- Copy to model
	for _, companyDataElm := range inputStruct.CompanyData {
		for _, branchDataElm := range companyDataElm.BranchData {
			modelCustomer = append(modelCustomer, repository.CustomerListModel{
				CompanyID:  sql.NullString{String: companyDataElm.CompanyID},
				BranchID:   sql.NullString{String: branchDataElm.BranchID},
				BranchName: sql.NullString{String: branchDataElm.BranchName},
			})
		}
	}

	var (
		result          []repository.CustomerListModel
		invalidItemList []repository.CustomerListModel
		validCounting   = 0
	)

	for _, modelCustomerItem := range modelCustomer {
		var resultTemp []repository.CustomerListModel

		//------- Check customer to DB
		resultTemp, err = dao.CustomerListDAO.CheckCustomerByProductName(tx, modelCustomerItem)
		if err.Error != nil {
			return
		}

		////------- If empty then registration because may input wrong company or branch id
		//if resultTemp == nil {
		//	invalidItemList = append(invalidItemList, modelCustomerItem)
		//	continue
		//}

		if resultTemp == nil {
			var custInstallationOnDB []repository.CustomerInstallationForConfig
			custInstallationOnDB, err = dao.CustomerInstallationDAO.GetCustomerInstallationByUniqueID(serverconfig.ServerAttribute.DBConnection, []repository.CustomerInstallationDetail{
				{
					UniqueID1: modelCustomerItem.CompanyID,
					UniqueID2: modelCustomerItem.BranchID,
				},
			}, false)
			if err.Error != nil {
				return
			}

			if len(custInstallationOnDB) == 0 {
				invalidItemList = append(invalidItemList, modelCustomerItem)
				continue
			} else {
				if len(custInstallationOnDB) > 0 {
					validCounting++

					//--- Report ND6 WAR 24/05/2023
					if modelCustomerItem.BranchName.String == "" {
						modelCustomerItem.BranchName.String = modelCustomerItem.BranchID.String
					}

					result = append(result, repository.CustomerListModel{
						CompanyID:   sql.NullString{String: modelCustomerItem.CompanyID.String},
						BranchID:    sql.NullString{String: modelCustomerItem.BranchID.String},
						CompanyName: sql.NullString{String: modelCustomerItem.BranchName.String},
					})
				}
			}
		}

		//------- Append to result
		if len(resultTemp) > 0 && isCheckExp {
			for index, resultTempItem := range resultTemp {

				//------- Last, direct to append
				if (len(resultTemp)-index == 1) && (resultTempItem.ExpDate.Time.After(timeNow) || resultTempItem.ExpDate.Time.Equal(timeNow)) {
					validCounting++
					result = append(result, resultTempItem)
					break
				} else if len(resultTemp)-index == 1 {
					result = append(result, resultTempItem)
					break
				}

				//------- Take it the first data after timeNow
				if resultTempItem.ExpDate.Time.After(timeNow) || resultTempItem.ExpDate.Time.Equal(timeNow) {
					validCounting++
					result = append(result, resultTempItem)
					break
				} else {
					invalidItemList = append(invalidItemList, resultTempItem)
				}
			}
		} else if len(resultTemp) > 0 && !isCheckExp {
			for _, resultItem := range resultTemp {

				//------- Append it the first data
				validCounting++
				result = append(result, resultItem)
				break
			}
		}
	}

	/*
		Kemungkinan :
		1. result nya nol dan valid counting nya 0, maka dia error
		2. result nya ada dan valid counting nya 0, maka dia error
		3. result nya ada dan valid counting nya > 0, maka dia sukses
	*/

	if len(invalidItemList) > 0 {
		var invalidItemListResponse []interface{}
		for _, invalidItem := range invalidItemList {
			invalidItemResponse := out.DetailErrorRegistrationClientID{
				CompanyId:    invalidItem.CompanyID.String,
				BranchID:     invalidItem.BranchID.String,
				ErrorMessage: "ID perusahaan atau cabang tidak valid",
			}

			invalidItemListResponse = append(invalidItemListResponse, invalidItemResponse)
		}

		err = errorModel.GenerateInvalidRegistrationClientWithDetailData(input.FileName, funcName, invalidItemListResponse)
		return
	}

	if (len(result) < 1 && validCounting < 1) || (len(result) > 0 && validCounting < 1) {
		var invalidItemListResponse []interface{}
		for _, invalidItem := range invalidItemList {
			invalidItemResponse := out.DetailErrorRegistrationClientID{
				CompanyId:    invalidItem.CompanyID.String,
				BranchID:     invalidItem.BranchID.String,
				ErrorMessage: "ID perusahaan atau cabang tidak valid",
			}

			invalidItemListResponse = append(invalidItemListResponse, invalidItemResponse)
		}

		err = errorModel.GenerateInvalidRegistrationClientWithDetailData(input.FileName, funcName, invalidItemListResponse)
		return
	} else if len(result) > 0 && validCounting > 0 {

		//------- Get company name for client alias
		for indexCompany, companyDataElm2 := range inputStruct.CompanyData {
			for indexBranch, branchDataElm2 := range companyDataElm2.BranchData {
				for _, resultElm := range result {
					if companyDataElm2.CompanyID == resultElm.CompanyID.String && branchDataElm2.BranchID == resultElm.BranchID.String {

						inputStruct.CompanyData[indexCompany].BranchData[indexBranch].ClientAlias = resultElm.CompanyName.String
					}
				}
			}
		}

	}

	output = inputStruct
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) checkDataInCustomerInstallation(inputStruct in.ClientRequest) (output []repository.CustomerInstallationForConfig, err errorModel.ErrorModel) {
	var (
		installation []repository.CustomerInstallationDetail
		db           = serverconfig.ServerAttribute.DBConnection
	)

	for _, companyDataElm := range inputStruct.CompanyData {
		for _, branchDataElm := range companyDataElm.BranchData {
			installation = append(installation, repository.CustomerInstallationDetail{
				UniqueID1: sql.NullString{String: companyDataElm.CompanyID},
				UniqueID2: sql.NullString{String: branchDataElm.BranchID},
			})
		}
	}

	output, err = dao.CustomerInstallationDAO.GetCustomerInstallationByUniqueID(db, installation, true)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) CheckDataMustUpdateInCustomerInstallation(tx *sql.Tx, inputStruct []repository.CustomerInstallationForConfig, clientID string, ctx *applicationModel.ContextModel, timeNow time.Time) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		updatedBy     int64
		updatedClient string
	)

	if ctx.AuthAccessTokenModel.ResourceUserID == 0 {
		updatedBy = 1
		updatedClient = constanta.SystemClient
	} else {
		updatedBy = ctx.AuthAccessTokenModel.ResourceUserID
		updatedClient = ctx.AuthAccessTokenModel.ClientID
	}

	for _, item := range inputStruct {
		var ClientMappingOnDB repository.ClientMappingModel
		ClientMappingOnDB, err = dao.ClientMappingDAO.GetClientMappingByUniqueID12(tx, repository.ClientMappingModel{
			ClientID:  sql.NullString{String: clientID},
			CompanyID: item.UniqueID1,
			BranchID:  item.UniqueID2,
		})

		if err.Error != nil {
			return
		}

		if ClientMappingOnDB.ID.Int64 != 0 {
			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *ctx, timeNow, dao.CustomerInstallationDAO.TableName, item.ID.Int64, 0)...)
			err = dao.CustomerInstallationDAO.UpdateCustomerInstallationClientMappingID(tx, repository.CustomerInstallationModel{
				ID:              item.ID,
				ClientMappingID: ClientMappingOnDB.ID,
				UpdatedAt:       sql.NullTime{Time: timeNow},
				UpdatedBy:       sql.NullInt64{Int64: updatedBy},
				UpdatedClient:   sql.NullString{String: updatedClient},
			})

			if err.Error != nil {
				return
			}
		}
	}

	return
}
