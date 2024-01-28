package master_data_dao

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"strconv"
)

var subDistrictDAOName = "SubDistrict.go"

func GetViewSubDistrictFromMasterData(inputStruct int64, contextModel *applicationModel.ContextModel, internalToken string) (result master_data_response.SubDistrictResponse, err errorModel.ErrorModel) {
	fileName := subDistrictDAOName
	funcName := "GetViewSubDistrictFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	listPositionPath := masterData.Host + masterData.PathRedirect.SubDistrict.View + "/" + strconv.Itoa(int(inputStruct))

	strData, err := hitAPIToMasterData(fileName, funcName, listPositionPath, http.MethodGet, contextModel, "")
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(strData), &result)

	err = errorModel.GenerateNonErrorModel()
	return
}

func CountAllSubDistrictFromMasterData(inputStruct master_data_request.SubDistrictRequest, contextModel *applicationModel.ContextModel, internalToken string) (result int, err errorModel.ErrorModel) {
	fileName := subDistrictDAOName
	funcName := "CountAllSubDistrictFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	listPositionPath := masterData.Host + masterData.PathRedirect.SubDistrict.Count + "/all"
	inputStruct.Page = -99
	inputStruct.Limit = -99

	strData, err := hitAPIToMasterDataWithInternalToken(fileName, funcName, listPositionPath, http.MethodPost, internalToken, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(strData), &result)

	err = errorModel.GenerateNonErrorModel()
	return
}

func GetListSubDistrictFromMasterData(inputStruct master_data_request.SubDistrictRequest, contextModel *applicationModel.ContextModel, internalToken string) (result []master_data_response.SubDistrictResponse, err errorModel.ErrorModel) {
	fileName := subDistrictDAOName
	funcName := "GetListSubDistrictFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	listPositionPath := masterData.Host + masterData.PathRedirect.SubDistrict.GetList

	strData, err := hitAPIToMasterDataWithInternalToken(fileName, funcName, listPositionPath, http.MethodPost, internalToken, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(strData), &result)

	err = errorModel.GenerateNonErrorModel()
	return
}

func GetListAllSubDistrictFromMasterData(inputStruct master_data_request.SubDistrictRequest, contextModel *applicationModel.ContextModel, internalToken string) (result []master_data_response.SubDistrictResponse, err errorModel.ErrorModel) {
	fileName := subDistrictDAOName
	funcName := "GetListAllSubDistrictFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	listPositionPath := masterData.Host + masterData.PathRedirect.SubDistrict.GetList + "/all"
	inputStruct.Page = -99

	strData, err := hitAPIToMasterDataWithInternalToken(fileName, funcName, listPositionPath, http.MethodPost, internalToken, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(strData), &result)

	err = errorModel.GenerateNonErrorModel()
	return
}
