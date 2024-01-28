package master_data_dao

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"strconv"
)

var personProfileDAOFileName = "PersonProfileDAO.go"

func GetListPersonProfileFromMasterData(inputStruct master_data_request.PersonProfileGetListRequest, contextModel *applicationModel.ContextModel) (result []master_data_response.GetListPersonProfileResponse, err errorModel.ErrorModel) {
	funcName := "GetListPersonProfileFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	masterDataPath := masterData.Host + masterData.PathRedirect.PersonProfile.GetList

	strData, err := hitAPIToMasterData(personProfileDAOFileName, funcName, masterDataPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorUnmarshal := json.Unmarshal([]byte(strData), &result)
	if errorUnmarshal != nil {
		err = errorModel.GenerateUnknownError(personProfileDAOFileName, funcName, errorUnmarshal)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func GetViewPersonProfileFromMasterData(inputStruct master_data_request.PersonProfileGetListRequest, contextModel *applicationModel.ContextModel) (result master_data_response.ViewPersonProfileResponse, err errorModel.ErrorModel) {
	funcName := "GetViewPersonProfileFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	masterDataPath := masterData.Host + masterData.PathRedirect.PersonProfile.View + "/" + strconv.Itoa(int(inputStruct.ID))

	strData, err := hitAPIToMasterData(personProfileDAOFileName, funcName, masterDataPath, http.MethodGet, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorUnmarshal := json.Unmarshal([]byte(strData), &result)
	if errorUnmarshal != nil {
		err = errorModel.GenerateUnknownError(personProfileDAOFileName, funcName, errorUnmarshal)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func InsertPersonProfileToMasterData(inputStruct master_data_request.PersonProfileWriteRequest, contextModel *applicationModel.ContextModel) (insertedID int64, err errorModel.ErrorModel) {
	var (
		funcName       = "InsertPersonProfileToMasterData"
		masterData     = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterData.Host + masterData.PathRedirect.PersonProfile.View
		authorize      = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, contextModel.AuthAccessTokenModel.Locale)
		multipartData  = make(map[string]common.MultiPartData)
		apiResponse    out.APIResponse
		resultRequest  master_data_response.MasterDataInsertedIDResponse
	)

	headerRequest := make(map[string][]string)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	multipartData["content"] = common.MultiPartData{Data: util2.StructToJSON(inputStruct)}

	statusCode, _, bodyResult, errorS := common.HitMultipartFormDataRequest(masterDataPath, http.MethodPost, headerRequest, multipartData, *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(personProfileDAOFileName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &apiResponse)

	if statusCode == 200 {
		tempData := apiResponse.Nexsoft.Payload.Data.Content
		result := util2.StructToJSON(tempData)
		errorS = json.Unmarshal([]byte(result), &resultRequest)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(companyProfileDAOName, funcName, errorS)
			return
		}
		insertedID = resultRequest.ID
	} else {
		causedBy := errors.New(apiResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(companyProfileDAOName, funcName, statusCode, apiResponse.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func ValidateBulkPersonProfileToMasterData(inputStruct []master_data_request.PersonProfileWriteRequest, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		funcName       = "ValidateBulkPersonProfileToMasterData"
		masterData     = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterData.Host + masterData.PathRedirect.PersonProfile.Validate
		authorize      = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, contextModel.AuthAccessTokenModel.Locale)
		multipartData  = make(map[string]common.MultiPartData)
		bodyResult     string
		statusCode     int
		errorS         error
		apiResponse    out.APIResponse
	)

	// Add token to header
	headerRequest := make(map[string][]string)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	// Validate request
	for _, request := range inputStruct {
		// Add body request
		multipartData["content"] = common.MultiPartData{Data: util2.StructToJSON(request)}

		//do validate
		statusCode, _, bodyResult, errorS = common.HitMultipartFormDataRequest(masterDataPath, http.MethodPost, headerRequest, multipartData, *contextModel)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(personProfileDAOFileName, funcName, errorS)
			return
		}

		_ = json.Unmarshal([]byte(bodyResult), &apiResponse)

		if statusCode != 200 {
			causedBy := errors.New(apiResponse.Nexsoft.Payload.Status.Message)
			err = errorModel.GenerateAuthenticationServerError(companyProfileDAOName, funcName, statusCode, apiResponse.Nexsoft.Payload.Status.Code, causedBy)
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func ValidatePersonProfileToMasterData(inputStruct master_data_request.PersonProfileWriteRequest, contextModel *applicationModel.ContextModel, isUpdate bool) (err errorModel.ErrorModel) {
	var (
		funcName       = "ValidatePersonProfileToMasterData"
		masterData     = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterData.Host + masterData.PathRedirect.PersonProfile.Validate
		authorize      = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, contextModel.AuthAccessTokenModel.Locale)
		multipartData  = make(map[string]common.MultiPartData)
		bodyResult     string
		statusCode     int
		errorS         error
		apiResponse    out.APIResponse
		methode        = http.MethodPost
	)

	if isUpdate {
		methode = http.MethodPut
	}

	// Add token to header
	headerRequest := make(map[string][]string)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	// Validate request section
	// Add body request
	multipartData["content"] = common.MultiPartData{Data: util2.StructToJSON(inputStruct)}

	//do validate
	statusCode, _, bodyResult, errorS = common.HitMultipartFormDataRequest(masterDataPath, methode, headerRequest, multipartData, *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(personProfileDAOFileName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &apiResponse)

	if statusCode != 200 {
		causedBy := errors.New(apiResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(companyProfileDAOName, funcName, statusCode, apiResponse.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func UpdatePersonProfileToMasterData(inputStruct master_data_request.PersonProfileWriteRequest, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		funcName       = "UpdatePersonProfileToMasterData"
		masterData     = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterData.Host + masterData.PathRedirect.PersonProfile.View + "/" + strconv.Itoa(int(inputStruct.ID))
		authorize      = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, contextModel.AuthAccessTokenModel.Locale)
		multipartData  = make(map[string]common.MultiPartData)
		apiResponse    out.APIResponse
	)

	fmt.Println("[Update Person Profile To Master Data Req] => ", util2.StructToJSON(inputStruct))
	headerRequest := make(map[string][]string)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	multipartData["content"] = common.MultiPartData{Data: util2.StructToJSON(inputStruct)}

	statusCode, _, bodyResult, errorS := common.HitMultipartFormDataRequest(masterDataPath, http.MethodPut, headerRequest, multipartData, *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(personProfileDAOFileName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &apiResponse)

	if statusCode != 200 {
		causedBy := errors.New(apiResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(companyProfileDAOName, funcName, statusCode, apiResponse.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
