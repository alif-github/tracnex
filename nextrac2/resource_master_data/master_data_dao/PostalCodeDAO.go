package master_data_dao

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"strconv"
)

func GetViewPostalCodeFromMasterData(inputStruct master_data_request.PostalCodeRequest, contextModel *applicationModel.ContextModel) (result master_data_response.UrbanVillageResponse, err errorModel.ErrorModel) {
	fileName := "PostalCode.go"
	funcName := "GetViewPostalCodeFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	listPositionPath := masterData.Host + masterData.PathRedirect.PostalCode.View + "/" + strconv.Itoa(int(inputStruct.ID))

	strData, err := hitAPIToMasterData(fileName, funcName, listPositionPath, http.MethodGet, contextModel, "")
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(strData), &result)

	result.UpdatedAt, err = in.TimeStrToTime(result.UpdatedAtStr, constanta.UpdatedAt)
	return
}

func CountAllPostalCodeFromMasterData(resultOnDB repository.PostalCodeModel, contextModel *applicationModel.ContextModel) (result int64, err errorModel.ErrorModel) {
	var (
		inputStruct master_data_request.PostalCodeRequest
		resultTemp  string
		fileName    = "PostalCodeDAO.go"
		funcName    = "CountAllPostalCodeFromMasterData"
		masterData  = config.ApplicationConfiguration.GetMasterData()
		countPath   = masterData.Host + masterData.PathRedirect.PostalCode.Count + "/all"
	)

	inputStruct = master_data_request.PostalCodeRequest{
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

func GetListForSyncPostalCodeFromMDB(inputStruct master_data_request.PostalCodeRequest, contextModel *applicationModel.ContextModel) (result []master_data_response.PostalCodeResponse, err errorModel.ErrorModel) {
	var (
		fileName    = "PostalCodeDAO.go"
		funcName    = "GetListForSyncPostalCodeFromMDB"
		masterData  = config.ApplicationConfiguration.GetMasterData()
		getListPath = masterData.Host + masterData.PathRedirect.PostalCode.GetList + "/all"
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
