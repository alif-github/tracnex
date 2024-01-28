package HitAddExternalResource

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type hitAddExternalResourceService struct {
	service.RegistrationPrepared
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var HitAddExternalResource = hitAddExternalResourceService{}.New()

func (input hitAddExternalResourceService) New() (output hitAddExternalResourceService) {
	output.FileName = "AddResourceExternalService.go"
	output.ValidLimit = []int{10, 20, 50, 100}
	output.ValidOrderBy = []string{"id"}
	output.IdResourceAllowed = []int64{1,2}
	return
}

func (input hitAddExternalResourceService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(addResourceRequest *in.AddResourceExternalRequest) errorModel.ErrorModel) (inputStruct in.AddResourceExternalRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string
	var attributeRequest *in.AttributeRequestAddResource
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

	//---------- Init, Important request must be exist
	if util.IsStringEmpty(inputStruct.ClientID) {
		err = errorModel.GenerateEmptyFieldError(input.FileName, funcName, constanta.ClientID)
		return
	}

	//---------- Unmarshal String Body to Attribute request
	errorS = json.Unmarshal([]byte(stringBody), &attributeRequest)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
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
	b, _ := json.Marshal(attributeRequest)
	inputStruct.Body = string(b)

	//---------- Get Log
	result, err := dao.ClientRegistrationLogDAO.GetDataClientRegistrationLogForUpdate(serverconfig.ServerAttribute.DBConnection, repository.ClientRegistrationLogModel{ClientID: sql.NullString{String: inputStruct.ClientID}})
	if err.Error != nil {
		return
	}

	//---------- If different client type id, then unknown data error
	if inputStruct.ClientTypeID != result.ClientTypeID.Int64 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientTypeID + " & " + constanta.ClientID)
		return
	}

	if result.ID.Int64 > 0 {
		inputStruct.ID = result.ID.Int64
		inputStruct.OldResource = result.Resource.String
		inputStruct.UpdatedAt = result.UpdatedAt.Time
	} else {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientTypeID + " & " + constanta.ClientID)
		return
	}

	//---------- Main validation
	err = validation(&inputStruct)
	return
}

func (input hitAddExternalResourceService) readBodyAndValidateForView(request *http.Request, _ *applicationModel.ContextModel, validation func(addResourceRequest *in.AddResourceExternalRequest) errorModel.ErrorModel) (inputStruct in.AddResourceExternalRequest, err errorModel.ErrorModel) {

	clientID, _ := mux.Vars(request)["CLIENT"]

	if inputStruct.ClientID == "" {
		inputStruct.ClientID = clientID
	}

	err = validation(&inputStruct)
	return
}

func (input hitAddExternalResourceService) AddResourceNexcloud(inputStruct in.AddResourceNexcloud, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	funcName := "AddResourceNexcloud"
	var payloadMessage authentication_response.AuthenticationErrorResponse

	internalToken := resource_common_service.GenerateInternalToken(constanta.NexCloudResourceID, 0, contextModel.AuthAccessTokenModel.ClientID, config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
	nexcloudAPIServer := config.ApplicationConfiguration.GetNexcloudAPI()
	addResourceNexcloudUrl := nexcloudAPIServer.Host + nexcloudAPIServer.PathRedirect.AddResourceClient

	header := make(map[string][]string)
	header[common.AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, errorS := common.HitAPI(addResourceNexcloudUrl, header, util.StructToJSON(inputStruct), "POST", *contextModel)

	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &payloadMessage)

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(payloadMessage.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, payloadMessage.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}

func (input hitAddExternalResourceService) AddResourceNexdrive(inputStruct in.AddResourceNexdrive, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	funcName := "AddResourceNexdrive"
	var payloadMessage authentication_response.AuthenticationErrorResponse

	internalToken := resource_common_service.GenerateInternalToken("drive", 0, contextModel.AuthAccessTokenModel.ClientID, config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
	nexdriveConfig := config.ApplicationConfiguration.GetNexdrive()
	addResourceNexdriveUrl := nexdriveConfig.Host + nexdriveConfig.PathRedirect.AddResourceClient

	header := make(map[string][]string)
	header[common.AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, errorS := common.HitAPI(addResourceNexdriveUrl, header, util.StructToJSON(inputStruct), "POST", *contextModel)

	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &payloadMessage)

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(payloadMessage.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, payloadMessage.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}

func (input hitAddExternalResourceService) doGetFirstNameUser(inputStruct in.AddResourceExternalRequest) (firstName string, err errorModel.ErrorModel) {
	fileName := "AddResourceExternalService.go"
	funcName := "doGetFirstNameUser"

	userModel, err := dao.UserDAO.GetIdAndFirstNameUser(serverconfig.ServerAttribute.DBConnection, repository.UserModel {
		ClientID: 		sql.NullString{String: inputStruct.ClientID},
	})

	firstName = userModel.FirstName.String

	if err.Error != nil {
		return
	}

	if firstName == "" {
		err = errorModel.GenerateDataNotFound(fileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hitAddExternalResourceService) updateRegistrationLogWithAudit(inputStruct in.AddResourceExternalRequest, doForResource string,
	contextModel *applicationModel.ContextModel, errors errorModel.ErrorModel) (err errorModel.ErrorModel) {

	funcName := "updateRegistrationLogWithAudit"

	structForUpdateLog := input.fillModelForRegistrationLog(inputStruct, contextModel, errors, doForResource)
	_, err = input.ServiceWithDataAuditPreparedByService(funcName, structForUpdateLog, contextModel, input.doUpdateRegistrationLog, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hitAddExternalResourceService) doUpdateRegistrationLog(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel,
	timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {

	fileName := "AddResourceExternalService.go"
	funcName := "doUpdateRegistrationLog"

	inputStruct := inputStructInterface.(in.PreparedRepositoryClientRegisterLog)

	clientRegistrationLogModel := repository.ClientRegistrationLogModel {
		ClientID:              sql.NullString{String: inputStruct.ClientID},
		AttributeRequest:      sql.NullString{String: inputStruct.AttributeRequest},
		SuccessStatus: 		   sql.NullBool{Bool: inputStruct.Status},
		Message:       		   sql.NullString{String: inputStruct.Message},
		Details:               sql.NullString{String: inputStruct.Detail},
		Code:                  sql.NullString{String: inputStruct.Code},
		Resource:              sql.NullString{String: inputStruct.Resource},
		RequestTimeStamp:      sql.NullTime{Time: timeNow},
		UpdatedBy:             sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:         sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:             sql.NullTime{Time: timeNow},
	}

	result, err := dao.ClientRegistrationLogDAO.GetDataClientRegistrationLogForUpdate(serverconfig.ServerAttribute.DBConnection, repository.ClientRegistrationLogModel{ClientID: sql.NullString{String: inputStruct.ClientID}})
	if err.Error != nil {
		return
	}

	if inputStruct.UpdatedAt != result.UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(fileName, funcName, "Client ID: "+ inputStruct.ClientID)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ClientRegistrationLogDAO.TableName, result.ID.Int64, 0)...)

	err = dao.ClientRegistrationLogDAO.UpdateRegistrationLogForAddResource(tx, clientRegistrationLogModel, inputStruct.ProcessForResource)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hitAddExternalResourceService) fillModelForRegistrationLog(inputStruct in.AddResourceExternalRequest, contextModel *applicationModel.ContextModel,
	errorS errorModel.ErrorModel, doForResource string) (preparedLogStruct in.PreparedRepositoryClientRegisterLog) {

	var detail, addResource, messageID string

	switch doForResource {
	case constanta.NexCloudResourceID :
		addResource = constanta.NexCloudResourceID
		messageID = "SUCCESS_ADD_RESOURCE_NEXCLOUD_MESSAGE"
		break
	case constanta.NexdriveResourceID :
		addResource = constanta.NexdriveResourceID
		messageID = "SUCCESS_ADD_RESOURCE_NEXDRIVE_MESSAGE"
	}

	if errorS.Error != nil {
		if len(errorS.AdditionalInformation) > 0 {
			detail = errorS.AdditionalInformation[0]
		}
		preparedLogStruct = in.PreparedRepositoryClientRegisterLog{
			Status: 			false,
			Code: 				errorS.Error.Error(),
			Detail: 			detail,
			Resource: 			inputStruct.OldResource,
			Message: 			util2.GenerateI18NErrorMessage(errorS, contextModel.AuthAccessTokenModel.Locale),
		}
	} else {
		preparedLogStruct = in.PreparedRepositoryClientRegisterLog{
			Status: 			true,
			Code: 				util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
			Resource: 			inputStruct.OldResource + " " + addResource,
			Message: 			GenerateI18NMessage(messageID, contextModel.AuthAccessTokenModel.Locale),
		}
	}

	preparedLogStruct.AttributeRequest = inputStruct.Body
	preparedLogStruct.ClientID = inputStruct.ClientID
	preparedLogStruct.ProcessForResource = doForResource
	preparedLogStruct.UpdatedAt = inputStruct.UpdatedAt

	return
}

func (input hitAddExternalResourceService) validateInsertAddResource(inputStruct *in.AddResourceExternalRequest) errorModel.ErrorModel {
	return inputStruct.ValidateAddResourceExternal()
}

func (input hitAddExternalResourceService) checkClientTypeByID(inputStruct *in.AddResourceExternalRequest) (err errorModel.ErrorModel) {
	funcName := "checkClientTypeByID"

	var result repository.ClientTypeModel

	result, err = dao.ClientTypeDAO.CheckClientTypeByID(serverconfig.ServerAttribute.DBConnection, &repository.ClientTypeModel{
		ID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
	})

	if err.Error != nil {
		return
	}

	if result.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientTypeID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}