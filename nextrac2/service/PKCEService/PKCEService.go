package PKCEService

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
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
	"nexsoft.co.id/nextrac2/resource_common_service/dto/nexcloud_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/nexcloud_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/ClientRegistrationLogService"
	"nexsoft.co.id/nextrac2/service/ClientService"
	"nexsoft.co.id/nextrac2/service/UserService/CRUDUserService"
	"nexsoft.co.id/nextrac2/util"
	"strings"
	"time"
)

type pkceService struct {
	service.RegistrationPrepared
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var PkceService = pkceService{}.New()

func (input pkceService) New() (output pkceService) {
	output.FileName = "PKCEService.go"
	output.ValidLimit = []int{10, 20, 50, 100}
	output.ValidSearchBy = []string{"username", "parent_client_id"}
	output.ValidOrderBy = []string{"id"}
	output.IdResourceAllowed = []int64{2}
	return
}

func (input pkceService) readUrlPathUnregisPKCE(request *http.Request, validation func(pkceRequest *in.PKCERequest) errorModel.ErrorModel) (inputStruct in.PKCERequest, err errorModel.ErrorModel) {
	funcName := "readUrlPathUnregisPKCE"

	strUsername, isExit := mux.Vars(request)["USERNAME"]

	if !isExit {
		err = errorModel.GenerateUnsupportedRequestParam(input.FileName, funcName)
		return
	}

	inputStruct.Username = strUsername

	err = validation(&inputStruct)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) readBodyAndValidateRegisUnregisPKCE(request *http.Request, contextModel *applicationModel.ContextModel, isForRegister bool,
	validation func(pkceRequest *in.PKCERequest) errorModel.ErrorModel) (inputStruct in.PKCERequest, err errorModel.ErrorModel) {

	funcName := "readBodyAndValidateRegisUnregisPKCE"
	var stringBody string
	var attributeRequestStr *in.AttributeRequestRegistPKCE
	var isAllowed bool

	//---------- Read Body Request
	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	//---------- Unmarshal String Body to Main struct
	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	//---------- Init, Important request must be exist
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

	if isForRegister {
		//---------- Unmarshal String Body to Attribute request
		errorS = json.Unmarshal([]byte(stringBody), &attributeRequestStr)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
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

	//---------- Login must be same with parent client ID
	isNexmile := inputStruct.ClientTypeID == constanta.ResourceNexmileID
	if isNexmile {
		if inputStruct.ParentClientID == "" {
			err = errorModel.GenerateEmptyFieldError(input.FileName, funcName, constanta.ParentClientID)
			return
		}

		if contextModel.AuthAccessTokenModel.ClientID != inputStruct.ParentClientID {
			err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
			return
		}
	}

	if isForRegister {
		//---------- Adding to main struct -> attribute request
		b, _ := json.Marshal(attributeRequestStr)
		inputStruct.Body = string(b)
	}

	//---------- Main validation
	err = validation(&inputStruct)
	return
}

func (input pkceService) checkClientMappingValid(tx *sql.Tx, inputStructInterface interface{}) (output interface{}, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.PKCERequest)
	var clientData interface{}
	var branchData []in.BranchData
	var companyData []in.CompanyData

	branchData = append(branchData, in.BranchData{
		ClientID: inputStruct.ParentClientID,
		BranchID: inputStruct.BranchID,
	})
	companyData = append(companyData, in.CompanyData{
		CompanyID:  inputStruct.CompanyID,
		BranchData: branchData,
	})

	clientData = in.ClientRequest{
		ClientTypeID: inputStruct.ClientTypeID,
		CompanyData:  companyData,
	}

	// Memisahkan data yang telah ada di client mapping
	result, err := ClientService.ClientService.DoCheckClientMapping(tx, clientData, true)

	if err.Error != nil {
		return
	}

	newStructClient := result.(in.ClientRequest)
	var newStructPKCE in.PKCERequest

	for _, companyDataElm := range newStructClient.CompanyData {
		for _, branchDataElm := range companyDataElm.BranchData {
			newStructPKCE = in.PKCERequest{
				BodyRequest:    in.BodyRequest{Body: inputStruct.Body},
				ParentClientID: branchDataElm.ClientID,
				ClientTypeID:   newStructClient.ClientTypeID,
				ClientAlias:    branchDataElm.ClientAlias,
				CompanyID:      companyDataElm.CompanyID,
				BranchID:       branchDataElm.BranchID,
				Username:       inputStruct.Username,
				Password:       inputStruct.Password,
				FirstName:      inputStruct.FirstName,
				LastName:       inputStruct.LastName,
				Email:          inputStruct.Email,
				Phone:          inputStruct.Phone,
			}
		}
	}

	output = newStructPKCE
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) HitUserRegistrationToAuthServer(inputStructInterface interface{}, contextModel *applicationModel.ContextModel) (pkceUserResponse out.PKCEResponse, err errorModel.ErrorModel) {
	funcName := "HitUserRegistrationToAuthServer"
	inputStruct := inputStructInterface.(in.PKCERequest)
	userRequest := in.UserRequest{
		Username:    inputStruct.Username,
		Password:    inputStruct.Password,
		FirstName:   inputStruct.FirstName,
		LastName:    inputStruct.LastName,
		Email:       inputStruct.Email,
		CountryCode: constanta.IndonesianCodeNumber,
		Phone:       inputStruct.Phone,
		Locale:      constanta.IndonesianLanguage,
	}

	registerUserResponse, err := CRUDUserService.UserService.AddUserToAuthenticationServer(userRequest, contextModel, true)

	if err.Error != nil {
		if err.CausedBy.Error() == "User already exist" {
			detail := util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.PKCEUserBundle, "DETAIL_ERROR_INVALID_USERNAME_MESSAGE", contextModel.AuthAccessTokenModel.Locale, nil)
			err = errorModel.GenerateInvalidRegistrationPKCE(input.FileName, funcName, []string{detail})
			return
		} else {
			return
		}
	}

	pkceUserResponse = out.PKCEResponse{
		UserID:   registerUserResponse.Nexsoft.Payload.Data.Content.UserID,
		ClientID: registerUserResponse.Nexsoft.Payload.Data.Content.ClientID,
		Username: inputStruct.Username,
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) addResourceNextracToAuthServer(inputStruct out.PKCEResponse, contextModel *applicationModel.ContextModel) (addUserResourceResult authentication_response.AddClientAuthenticationResponse, err errorModel.ErrorModel) {
	funcName := "addResourceNextracToAuthServer"
	internalToken := resource_common_service.GenerateInternalToken("auth", 0, contextModel.AuthAccessTokenModel.ClientID, config.ApplicationConfiguration.GetServerResourceID(), contextModel.AuthAccessTokenModel.Locale)
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	addResourceUserUrl := authenticationServer.Host + authenticationServer.PathRedirect.AddResourceClient

	statusCode, bodyResult, errorS := common.HitAddClientResource(internalToken, addResourceUserUrl, inputStruct.ClientID, config.ApplicationConfiguration.GetServerResourceID(), contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	if statusCode == 200 {
		_ = json.Unmarshal([]byte(bodyResult), &addUserResourceResult)
	} else {
		err = common.ReadAuthServerError(funcName, statusCode, bodyResult, contextModel)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) addResourceUserNexcloudToNexcloudServer(pkceResponseAuth out.PKCEResponse, pkceRequest in.PKCERequest, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	addResourceNexcloud := in.AddResourceNexcloud{
		FirstName: pkceRequest.FirstName,
		LastName:  pkceRequest.LastName,
		ClientID:  pkceResponseAuth.ClientID,
	}
	err = service.AddResourceNexcloudToNexcloud(addResourceNexcloud, contextModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) checkUserToAuthServer(pkceResponseAuth out.PKCEResponse, contextModel *applicationModel.ContextModel) (resourceID []out.ResourceList, err errorModel.ErrorModel) {
	checkUserStruct := authentication_request.CheckClientOrUser{
		ClientID: pkceResponseAuth.ClientID,
	}
	checkClientUserResp, err := service.CheckClientOrUserInAuth(checkUserStruct, contextModel)
	if err.Error != nil {
		return
	}

	isExist := checkClientUserResp.Nexsoft.Payload.Data.Content.IsExist

	if isExist {
		resourceIDResp := checkClientUserResp.Nexsoft.Payload.Data.Content.AdditionalInformation.ResourceID
		resourceIDRespArray := strings.Split(resourceIDResp, " ")
		for _, resourceIDRespElm := range resourceIDRespArray {
			resourceID = append(resourceID, out.ResourceList{
				ResourceID: resourceIDRespElm,
			})
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) hitUnregisterNexcloudToNexcloud(pkceModel repository.PKCEClientMappingModel, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	funcName := "addResourceNextracToAuthServer"
	var viewUserByClientIDResponse nexcloud_response.UserNexcloudResponse

	internalToken := resource_common_service.GenerateInternalToken(constanta.NexCloudResourceID, 0, contextModel.AuthAccessTokenModel.ClientID, config.ApplicationConfiguration.GetServerResourceID(), contextModel.AuthAccessTokenModel.Locale)
	nexcloudServer := config.ApplicationConfiguration.GetNexcloudAPI()
	deleteResourceClient := nexcloudServer.Host + nexcloudServer.PathRedirect.CrudClient + "/" + pkceModel.ClientID.String

	viewUserByClientIDResponse, err = input.internalGetUserNexcloud(internalToken, deleteResourceClient, contextModel)
	if err.Error != nil {
		return
	}

	unregisterClientDTO := nexcloud_request.UnregisterClient{
		UpdatedAt: viewUserByClientIDResponse.Nexsoft.Payload.Data.Content.UpdatedAt,
	}

	statusCode, bodyResult, errorS := common.HitUnregisterClientResourceNexcloud(internalToken, deleteResourceClient, unregisterClientDTO, contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		err = common.ReadAuthServerError(funcName, statusCode, bodyResult, contextModel)
		return
	}

	return
}

func (input pkceService) internalGetUserNexcloud(token string, path string, contextModel *applicationModel.ContextModel) (viewUserByClientIDResponse nexcloud_response.UserNexcloudResponse, err errorModel.ErrorModel) {
	funcName := "internalGetUserNexcloud"
	var statusCode int
	var bodyResult string
	var errs error

	header := make(map[string][]string)
	header[constanta.TokenHeaderNameConstanta] = []string{token}
	statusCode, _, bodyResult, errs = common.HitAPI(path, header, "", "GET", *contextModel)
	if errs != nil {
		err = errorModel.GenerateUnknownError("PKCEService.go", funcName, errs)
		return
	}

	if statusCode != 200 {
		err = common.ReadAuthServerError(funcName, statusCode, bodyResult, contextModel)
		return
	} else {
		_ = json.Unmarshal([]byte(bodyResult), &viewUserByClientIDResponse)
		return
	}
}

func (input pkceService) getSuccessResponse(inputStruct in.PKCERequest, contextModel *applicationModel.ContextModel) string {
	param := make(map[string]interface{})
	param["EMAIL"] = inputStruct.Email

	return util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.PKCEUserBundle, "SUCCESS_PKCE_REGIS_MESSAGE", contextModel.AuthAccessTokenModel.Locale, param)
}

func (input pkceService) isUnregisterBefore(inputStruct in.PKCERequest) (modelPKCEClientMapping repository.CheckPKCEClientMappingModel, err errorModel.ErrorModel) {

	modelPKCEClientMapping = repository.CheckPKCEClientMappingModel{
		Username:  sql.NullString{String: inputStruct.Username},
		Email:     sql.NullString{String: inputStruct.Email},
		Phone:     sql.NullString{String: constanta.IndonesianCodeNumber + "-" + inputStruct.Phone},
		Firstname: sql.NullString{String: inputStruct.FirstName},
		Lastname:  sql.NullString{String: inputStruct.LastName},
	}

	modelPKCEClientMapping, err = dao.PKCEClientMappingDAO.CheckUserPKCEUnregisterBefore(serverconfig.ServerAttribute.DBConnection, modelPKCEClientMapping)
	if err.Error != nil {
		return
	}

	if modelPKCEClientMapping.ID.Int64 < 1 {
		modelPKCEClientMapping.IsRegisteredBefore.Bool = false
	} else {
		modelPKCEClientMapping.IsRegisteredBefore.Bool = true
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) addResourceInfo(inputStruct out.PKCEResponse, resourceDataList []out.ResourceList) (result out.PKCEResponse) {
	inputStruct.ResourceList = resourceDataList
	result = inputStruct
	return
}

func (input pkceService) regisPKCEErrorHandle(resourceIDList []out.ResourceList, registerClientContent out.PKCEResponse, inputStruct in.PKCERequest,
	contextModel *applicationModel.ContextModel, errorS errorModel.ErrorModel, timeNow time.Time) (output interface{}, err errorModel.ErrorModel) {

	var resourceList []string
	var result interface{}
	var message string

	result = input.addResourceInfo(registerClientContent, resourceIDList)
	output, message, err = service.CustomFailedResponsePayload(result, errorS, contextModel)

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

	return
}

func (input pkceService) regisPKCESuccessHandle(responseAuth out.PKCEResponse, contextModel *applicationModel.ContextModel, inputStruct in.PKCERequest,
	resourceIDList []out.ResourceList, timeNow time.Time) (result interface{}, err errorModel.ErrorModel) {

	var code string
	var resourceList []string

	newRegisterClientContent := input.addResourceInfo(responseAuth, resourceIDList)
	successMessage := input.getSuccessResponse(inputStruct, contextModel)
	code, result = service.CustomSuccessResponsePayload(newRegisterClientContent, successMessage, contextModel)

	for _, resourceElm := range resourceIDList {
		if resourceElm.Status != "FAIL" {
			resourceList = append(resourceList, resourceElm.ResourceID)
		}
	}

	resourceListStr := strings.Join(resourceList, " ")

	_, err = ClientRegistrationLogService.ClientRegistrationLogService.InsertClientRegistrationLog(in.ClientRegisterLogRequest{
		ClientID:              responseAuth.ClientID,
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