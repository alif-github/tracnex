package master_data_dao

import (
	"encoding/json"
	"errors"
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

var companyProfileDAOName = "CompanyProfileDAO.go"

func GetListCompanyProfileFromMasterData(inputStruct master_data_request.CompanyProfileGetListRequest, contextModel *applicationModel.ContextModel) (result []master_data_response.CompanyProfileResponse, err errorModel.ErrorModel) {
	var (
		funcName       = "GetListCompanyProfileFromMasterData"
		masterData     = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterData.Host + masterData.PathRedirect.CompanyProfile.GetList
		strData        string
	)

	strData, err = hitAPIToMasterData(companyProfileDAOName, funcName, masterDataPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorUnmarshal := json.Unmarshal([]byte(strData), &result)
	if errorUnmarshal != nil {
		err = errorModel.GenerateUnknownError(companyProfileDAOName, funcName, errorUnmarshal)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func GetViewCompanyProfileFromMasterData(inputStruct master_data_request.CompanyProfileGetListRequest, contextModel *applicationModel.ContextModel) (result master_data_response.ViewCompanyProfileResponse, err errorModel.ErrorModel) {
	var (
		funcName       = "GetViewCompanyProfileFromMasterData"
		masterData     = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterData.Host + masterData.PathRedirect.CompanyProfile.View + "/" + strconv.Itoa(int(inputStruct.ID))
		strData        string
	)

	strData, err = hitAPIToMasterData(companyProfileDAOName, funcName, masterDataPath, http.MethodGet, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorUnmarshal := json.Unmarshal([]byte(strData), &result)
	if errorUnmarshal != nil {
		err = errorModel.GenerateUnknownError(companyProfileDAOName, funcName, errorUnmarshal)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func InsertCompanyProfile(inputStruct master_data_request.CompanyProfileWriteRequest, contextModel *applicationModel.ContextModel) (insertedID int64, err errorModel.ErrorModel) {
	var (
		funcName       = "InsertCompanyProfile"
		masterData     = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterData.Host + masterData.PathRedirect.CompanyProfile.View
		authorize      = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, contextModel.AuthAccessTokenModel.Locale)
		multipartData  = make(map[string]common.MultiPartData)
		resultRequest  master_data_response.MasterDataInsertedIDResponse
	)

	headerRequest := make(map[string][]string)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	multipartData["content"] = common.MultiPartData{Data: util2.StructToJSON(inputStruct)}

	statusCode, _, bodyResult, errorS := common.HitMultipartFormDataRequest(masterDataPath, http.MethodPost, headerRequest, multipartData, *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(companyProfileDAOName, funcName, errorS)
		return
	}

	var apiResponse out.APIResponse
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

func UpdateCompanyProfile(inputStruct master_data_request.CompanyProfileWriteRequest, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		funcName       = "UpdateCompanyProfile"
		masterData     = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterData.Host + masterData.PathRedirect.CompanyProfile.View + "/" + strconv.Itoa(int(inputStruct.ID))
		authorize      = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, contextModel.AuthAccessTokenModel.Locale)
		multipartData  = make(map[string]common.MultiPartData)
		companyProfile master_data_response.ViewCompanyProfileResponse
	)

	// Get updated at company profile
	companyProfile, err = GetViewCompanyProfileFromMasterData(master_data_request.CompanyProfileGetListRequest{ID: inputStruct.ID}, contextModel)
	if err.Error != nil {
		return
	}

	inputStruct.UpdatedAt = companyProfile.UpdatedAt

	headerRequest := make(map[string][]string)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	multipartData["content"] = common.MultiPartData{Data: util2.StructToJSON(inputStruct)}

	statusCode, _, bodyResult, errorS := common.HitMultipartFormDataRequest(masterDataPath, http.MethodPut, headerRequest, multipartData, *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(companyProfileDAOName, funcName, errorS)
		return
	}

	var apiResponse out.APIResponse
	_ = json.Unmarshal([]byte(bodyResult), &apiResponse)

	if statusCode != 200 {
		causedBy := errors.New(apiResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(companyProfileDAOName, funcName, statusCode, apiResponse.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
