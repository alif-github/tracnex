package master_data_dao

import (
	"encoding/json"
	"errors"
	"net/http"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"strconv"
)

var provinceDAOName = "ProvinceDAO.go"

func GetListProvinceFromMasterData(inputStruct in.ProvinceRequest, contextModel *applicationModel.ContextModel) (listProvince []master_data_response.ProvinceResponse, err errorModel.ErrorModel) {
	funcName := "GetListProvinceFromMasterData"
	fileName := provinceDAOName
	masterData := config.ApplicationConfiguration.GetMasterData()
	listProvincePath := masterData.Host + masterData.PathRedirect.Province.GetList

	headerRequest := make(map[string][]string)
	authorize := resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, constanta.DefaultApplicationsLanguage)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	statusCode, _, bodyResult, errorS := common.HitAPI(listProvincePath, headerRequest, util2.StructToJSON(inputStruct), "POST", *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	var listProvinceResponse master_data_response.GetListProvince
	_ = json.Unmarshal([]byte(bodyResult), &listProvinceResponse)

	switch statusCode {
	case 200:
		listProvince = listProvinceResponse.Nexsoft.Payload.Data.Content
		err = errorModel.GenerateNonErrorModel()
	case 404:
		err = errorModel.GenerateUnknownError(fileName, funcName, errors.New(bodyResult))
	default:
		causedBy := errors.New(listProvinceResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(fileName, funcName, statusCode, listProvinceResponse.Nexsoft.Payload.Status.Code, causedBy)
	}

	return
}

func ViewDetailProvinceFromMasterData(id int, contextModel *applicationModel.ContextModel, internalToken string) (dataProvince master_data_response.ProvinceResponse, err errorModel.ErrorModel) {
	var (
		funcName             = "ViewDetailProvinceFromMasterData"
		masterData           = config.ApplicationConfiguration.GetMasterData()
		viewProvincePath     = masterData.Host + masterData.PathRedirect.Province.View + "/" + strconv.Itoa(id)
		headerRequest        = make(map[string][]string)
		authorize            = internalToken
		viewProvinceResponse master_data_response.ViewProvince
	)

	if internalToken == "" {
		authorize = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, constanta.DefaultApplicationsLanguage)
	}
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	statusCode, _, bodyResult, errorS := common.HitAPI(viewProvincePath, headerRequest, "", "GET", *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(provinceDAOName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &viewProvinceResponse)
	if statusCode == 200 {
		dataProvince = viewProvinceResponse.Nexsoft.Payload.Data.Content
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(viewProvinceResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(provinceDAOName, funcName, statusCode, viewProvinceResponse.Nexsoft.Payload.Status.Code, causedBy)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func CountAllProvinceFromMasterData(inputStruct master_data_request.ProvinceRequest, contextModel *applicationModel.ContextModel) (result int64, err errorModel.ErrorModel) {
	funcName := "CountAllProvinceFromMasterData"
	fileName := provinceDAOName
	masterData := config.ApplicationConfiguration.GetMasterData()
	masterDataPath := masterData.Host + masterData.PathRedirect.Province.View + "/count/all"

	inputStruct.Page = -99

	strData, err := hitAPIToMasterData(fileName, funcName, masterDataPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorUnmarshal := json.Unmarshal([]byte(strData), &result)
	if errorUnmarshal != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorUnmarshal)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func CountProvinceFromMasterData(inputStruct master_data_request.ProvinceRequest, contextModel *applicationModel.ContextModel) (result int64, err errorModel.ErrorModel) {
	funcName := "CountAllProvinceFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	masterDataPath := masterData.Host + masterData.PathRedirect.Province.View + "/count"

	strData, err := hitAPIToMasterData(provinceDAOName, funcName, masterDataPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorUnmarshal := json.Unmarshal([]byte(strData), &result)
	if errorUnmarshal != nil {
		err = errorModel.GenerateUnknownError(provinceDAOName, funcName, errorUnmarshal)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func GetListAllProvinceFromMasterData(inputStruct master_data_request.ProvinceRequest, contextModel *applicationModel.ContextModel) (result []master_data_response.ProvinceResponse, err errorModel.ErrorModel) {
	funcName := "GetListCompanyProfileFromMasterData"
	fileName := provinceDAOName
	masterData := config.ApplicationConfiguration.GetMasterData()
	masterDataPath := masterData.Host + masterData.PathRedirect.Province.GetList + "/all"

	inputStruct.Page = -99

	strData, err := hitAPIToMasterData(fileName, funcName, masterDataPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorUnmarshal := json.Unmarshal([]byte(strData), &result)
	if errorUnmarshal != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorUnmarshal)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func GetListForSyncProvinceFromMasterData(inputStruct master_data_request.ProvinceRequest, contextModel *applicationModel.ContextModel) (result []master_data_response.ProvinceResponse, err errorModel.ErrorModel) {
	funcName := "GetListCompanyProfileFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	masterDataPath := masterData.Host + masterData.PathRedirect.Province.GetList

	strData, err := hitAPIToMasterData(countryDAOName, funcName, masterDataPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorUnmarshal := json.Unmarshal([]byte(strData), &result)
	if errorUnmarshal != nil {
		err = errorModel.GenerateUnknownError(countryDAOName, funcName, errorUnmarshal)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
