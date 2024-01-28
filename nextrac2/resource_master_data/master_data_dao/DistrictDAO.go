package master_data_dao

import (
	"encoding/json"
	"errors"
	"net/http"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"strconv"
)

var (
	districtDAOFileName = "DisrictDAO.go"
)

func GetListDistrictFromMasterData(inputStruct in.DistrictRequest, contextModel *applicationModel.ContextModel) (listDistrict []out.DistrictResponse, err errorModel.ErrorModel) {
	fileName := "DistrictDAO.go"
	funcName := "GetListDistrictFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	listDistrictPath := masterData.Host + masterData.PathRedirect.District.GetList

	headerRequest := make(map[string][]string)
	authorize := resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, constanta.DefaultApplicationsLanguage)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	statusCode, _, bodyResult, errorS := common.HitAPI(listDistrictPath, headerRequest, util2.StructToJSON(inputStruct), "POST", *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	var listDistrictResponse out.GetListDistrict
	_ = json.Unmarshal([]byte(bodyResult), &listDistrictResponse)

	if statusCode == 200 {
		listDistrict = listDistrictResponse.Nexsoft.Payload.Data.Content
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(listDistrictResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(fileName, funcName, statusCode, listDistrictResponse.Nexsoft.Payload.Status.Code, causedBy)
	}

	return
}

func ViewDetailDistrictFromMasterData(id int, contextModel *applicationModel.ContextModel, internalToken string) (dataDistrict out.DistrictResponse, err errorModel.ErrorModel) {
	var (
		fileName         = "DistrictDAO.go"
		funcName         = "ViewDetailDistrictFromMasterData"
		masterData       = config.ApplicationConfiguration.GetMasterData()
		viewDistrictPath = masterData.Host + masterData.PathRedirect.District.View + "/" + strconv.Itoa(id)
	)

	headerRequest := make(map[string][]string)
	authorize := resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, constanta.DefaultApplicationsLanguage)

	if internalToken != "" {
		authorize = internalToken
	}

	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	statusCode, _, bodyResult, errorS := common.HitAPI(viewDistrictPath, headerRequest, "", "GET", *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	var viewDistrictResponse out.ViewDistrict
	_ = json.Unmarshal([]byte(bodyResult), &viewDistrictResponse)

	if statusCode == 200 {
		dataDistrict = viewDistrictResponse.Nexsoft.Payload.Data.Content
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(viewDistrictResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(fileName, funcName, statusCode, viewDistrictResponse.Nexsoft.Payload.Status.Code, causedBy)
	}

	return
}

func CountAllDistrictFromMasterData(inputStruct master_data_request.ProvinceRequest, contextModel *applicationModel.ContextModel) (result int64, err errorModel.ErrorModel) {
	var (
		funcName       = "CountAllDistrictFromMasterData"
		masterData     = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterData.Host + masterData.PathRedirect.District.View + "/count/all"
	)

	inputStruct.Page = -99

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

func GetListForSyncDistrictFromMasterData(inputStruct master_data_request.DistrictRequest, contextModel *applicationModel.ContextModel) (result []master_data_response.DistrictResponse, err errorModel.ErrorModel) {
	funcName := "GetListForSyncDistrictFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	masterDataPath := masterData.Host + masterData.PathRedirect.District.GetList

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

func GetListAllDistrictFromMasterData(inputStruct master_data_request.DistrictRequest, contextModel *applicationModel.ContextModel) (result []master_data_response.DistrictResponse, err errorModel.ErrorModel) {
	funcName := "GetListAllDistrictFromMasterData"
	fileName := districtDAOFileName
	masterData := config.ApplicationConfiguration.GetMasterData()
	masterDataPath := masterData.Host + masterData.PathRedirect.District.GetList + "/all"

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
