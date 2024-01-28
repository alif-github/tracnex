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

func GetListPositionFromMasterData(inputStruct master_data_request.PositionGetListRequest, contextModel *applicationModel.ContextModel) (listPosition []master_data_response.PositionResponse, err errorModel.ErrorModel){
	fileName := "PositionDAO.go"
	funcName := "GetListPositionFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	listPositionPath := masterData.Host + masterData.PathRedirect.Position.GetList

	strData, err := hitAPIToMasterData(fileName, funcName, listPositionPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(strData), &listPosition)

	err = errorModel.GenerateNonErrorModel()
	return
}

func GetViewPositionFromMasterData(inputStruct master_data_request.PositionGetListRequest, contextModel *applicationModel.ContextModel) (result master_data_response.PositionResponse, err errorModel.ErrorModel){
	fileName := "PositionDAO.go"
	funcName := "GetViewPositionFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	urlPath := masterData.Host + masterData.PathRedirect.Position.View + "/" + strconv.Itoa(int(inputStruct.ID))

	strData, err := hitAPIToMasterData(fileName, funcName, urlPath, http.MethodGet, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(strData), &result)

	err = errorModel.GenerateNonErrorModel()
	return
}