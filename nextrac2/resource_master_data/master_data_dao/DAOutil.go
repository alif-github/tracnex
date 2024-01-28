package master_data_dao

import (
	"encoding/json"
	"errors"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
)

func hitAPIToMasterData(fileName, funcName, masterDataPath, method string, contextModel *applicationModel.ContextModel, inputStruct interface{}) (result string, err errorModel.ErrorModel) {
	headerRequest := make(map[string][]string)
	authorize := resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, contextModel.AuthAccessTokenModel.Locale)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	statusCode, _, bodyResult, errorS := common.HitAPI(masterDataPath, headerRequest, util2.StructToJSON(inputStruct), method, *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	var apiResponse out.APIResponse
	_ = json.Unmarshal([]byte(bodyResult), &apiResponse)

	switch statusCode {
	case 200:
		tempData := apiResponse.Nexsoft.Payload.Data.Content
		result = util2.StructToJSON(tempData)
		err = errorModel.GenerateNonErrorModel()
	case 404:
		err = errorModel.GenerateUnknownError(fileName, funcName, errors.New(bodyResult))
	default:
		causedBy := errors.New(apiResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(fileName, funcName, statusCode, apiResponse.Nexsoft.Payload.Status.Code, causedBy)
	}

	return
}

func hitAPIToMasterDataWithInternalToken(fileName, funcName, masterDataPath, method, token string, contextModel *applicationModel.ContextModel, inputStruct interface{}) (result string, err errorModel.ErrorModel) {
	headerRequest := make(map[string][]string)
	if token == "" {
		token = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, contextModel.AuthAccessTokenModel.Locale)
	}
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{token}

	statusCode, _, bodyResult, errorS := common.HitAPI(masterDataPath, headerRequest, util2.StructToJSON(inputStruct), method, *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	var apiResponse out.APIResponse
	_ = json.Unmarshal([]byte(bodyResult), &apiResponse)

	switch statusCode {
	case 200:
		tempData := apiResponse.Nexsoft.Payload.Data.Content
		result = util2.StructToJSON(tempData)
		err = errorModel.GenerateNonErrorModel()
	case 404:
		err = errorModel.GenerateUnknownError(fileName, funcName, errors.New(bodyResult))
	default:
		causedBy := errors.New(apiResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(fileName, funcName, statusCode, apiResponse.Nexsoft.Payload.Status.Code, causedBy)
	}

	return
}
