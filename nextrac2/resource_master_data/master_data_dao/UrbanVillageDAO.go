package master_data_dao

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"strconv"
)

func GetViewUrbanVillageFromMasterData(inputStruct master_data_request.UrbanVillageRequest, contextModel *applicationModel.ContextModel) (result master_data_response.UrbanVillageResponse, err errorModel.ErrorModel) {
	var (
		fileName         = "UrbanVillage.go"
		funcName         = "GetViewUrbanVillageFromMasterData"
		masterData       = config.ApplicationConfiguration.GetMasterData()
		listPositionPath = masterData.Host + masterData.PathRedirect.UrbanVillage.View + "/" + strconv.Itoa(int(inputStruct.ID))
		strData          string
	)

	strData, err = hitAPIToMasterData(fileName, funcName, listPositionPath, http.MethodGet, contextModel, "")
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(strData), &result)
	err = errorModel.GenerateNonErrorModel()
	return
}

func CountAllUrbanVillageFromMasterData(resultOnDB repository.UrbanVillageModel, contextModel *applicationModel.ContextModel) (result int64, err errorModel.ErrorModel) {
	var (
		inputStruct master_data_request.UrbanVillageRequest
		resultTemp  string
		fileName    = "UrbanVillageDAO.go"
		funcName    = "CountAllUrbanVillageFromMasterData"
		masterData  = config.ApplicationConfiguration.GetMasterData()
		countPath   = masterData.Host + masterData.PathRedirect.UrbanVillage.Count + "/all"
	)

	inputStruct = master_data_request.UrbanVillageRequest{
		UpdatedAtStart: resultOnDB.LastSync.Time,
		AbstractDTO: in.AbstractDTO{
			Page:  -99,
			Limit: -99,
		},
	}

	resultTemp, err = hitAPIToMasterData(fileName, funcName, countPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorS := json.Unmarshal([]byte(resultTemp), &result)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(countryDAOName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func GetListForSyncUrbanVillageFromMDB(inputStruct master_data_request.UrbanVillageRequest, contextModel *applicationModel.ContextModel) (result []master_data_response.UrbanVillageResponse, err errorModel.ErrorModel) {
	var (
		fileName    = "UrbanVillageDAO.go"
		funcName    = "GetListForSyncUrbanVillageFromMDB"
		masterData  = config.ApplicationConfiguration.GetMasterData()
		getListPath = masterData.Host + masterData.PathRedirect.UrbanVillage.GetList + "/all"
		resultTemp  string
	)

	resultTemp, err = hitAPIToMasterData(fileName, funcName, getListPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorS := json.Unmarshal([]byte(resultTemp), &result)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
