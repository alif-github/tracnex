package master_data_dao

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
)

var countryDAOName = "CountryDAO.go"

func GetListCountryCountryFromMasterData(inputStruct master_data_request.CountryRequest, contextModel *applicationModel.ContextModel) (result []master_data_response.CountryResponse, err errorModel.ErrorModel) {
	funcName := "GetListCompanyProfileFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	masterDataPath := masterData.Host + masterData.PathRedirect.Country.GetList

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

func GetListAllCountryCountryFromMasterData(contextModel *applicationModel.ContextModel) (result []master_data_response.CountryResponse, err errorModel.ErrorModel) {
	funcName := "GetListCompanyProfileFromMasterData"
	var inputStruct master_data_request.CountryRequest
	masterData := config.ApplicationConfiguration.GetMasterData()
	masterDataPath := masterData.Host + masterData.PathRedirect.Country.GetList

	inputStruct.Page = -99
	inputStruct.Limit = -99

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

func CountAllCountryCountryFromMasterData(contextModel *applicationModel.ContextModel) (result int64, err errorModel.ErrorModel) {
	funcName := "GetListCompanyProfileFromMasterData"
	var inputStruct interface{}
	masterData := config.ApplicationConfiguration.GetMasterData()
	masterDataPath := masterData.Host + masterData.PathRedirect.Country.Count

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
