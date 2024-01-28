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

func GetlistCompanyTitleFromMasterData(inputStruct master_data_request.CompanyTitleRequest, contextModel *applicationModel.ContextModel) (result []master_data_response.CompanyTitleResponse, err errorModel.ErrorModel){
	fileName := "PositionDAO.go"
	funcName := "GetlistCompanyTitleFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	listPositionPath := masterData.Host + masterData.PathRedirect.CompanyTitle.GetList

	strData, err := hitAPIToMasterData(fileName, funcName, listPositionPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(strData), &result)

	err = errorModel.GenerateNonErrorModel()
	return
}

func GetViewCompanyTitleFromMasterData(inputStruct master_data_request.CompanyTitleRequest, contextModel *applicationModel.ContextModel) (result master_data_response.CompanyTitleResponse, err errorModel.ErrorModel){
	fileName := "PositionDAO.go"
	funcName := "GetViewCompanyTitleFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	listPositionPath := masterData.Host + masterData.PathRedirect.CompanyTitle.View + "/" + strconv.Itoa(int(inputStruct.ID))

	strData, err := hitAPIToMasterData(fileName, funcName, listPositionPath, http.MethodGet, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(strData), &result)

	err = errorModel.GenerateNonErrorModel()
	return
}