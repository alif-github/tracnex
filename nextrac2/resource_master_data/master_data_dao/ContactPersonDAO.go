package master_data_dao

import (
	"encoding/json"
	"fmt"
	"net/http"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"strconv"
)

var contactPersonDAOName = "ContactPersonDAO.go"

func GetListContactPersonFromMasterData(inputStruct master_data_request.ContactPersonGetListRequest, contextModel *applicationModel.ContextModel) (result []master_data_response.ContactPersonResponse, err errorModel.ErrorModel) {
	var (
		funcName       = "GetListContactPersonFromMasterData"
		masterDataConf = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterDataConf.Host + masterDataConf.PathRedirect.ContactPerson.GetList
	)

	strData, err := hitAPIToMasterData(contactPersonDAOName, funcName, masterDataPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorUnmarshal := json.Unmarshal([]byte(strData), &result)
	if errorUnmarshal != nil {
		err = errorModel.GenerateUnknownError(contactPersonDAOName, funcName, errorUnmarshal)
		return
	}

	return
}

func GetViewContactPerson(inputStruct master_data_request.ContactPersonGetListRequest, contextModel *applicationModel.ContextModel) (result master_data_response.ViewContactPersonResponse, err errorModel.ErrorModel) {
	var (
		funcName       = "GetViewContactPerson"
		masterDataConf = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterDataConf.Host + masterDataConf.PathRedirect.ContactPerson.BaseUrl + "/" + strconv.Itoa(int(inputStruct.ID))
	)

	strData, err := hitAPIToMasterData(contactPersonDAOName, funcName, masterDataPath, http.MethodGet, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorUnmarshal := json.Unmarshal([]byte(strData), &result)
	if errorUnmarshal != nil {
		err = errorModel.GenerateUnknownError(contactPersonDAOName, funcName, errorUnmarshal)
		return
	}

	err = errorModel.GenerateNonErrorModel()

	return
}

func InsertContactPerson(inputStruct master_data_request.ContactPersonWriteRequest, contextModel *applicationModel.ContextModel) (insertedID int64, err errorModel.ErrorModel) {
	var (
		funcName       = "InsertContactPerson"
		masterDataConf = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterDataConf.Host + masterDataConf.PathRedirect.ContactPerson.BaseUrl
		resultRequest  master_data_response.MasterDataInsertedIDResponse
	)

	strData, err := hitAPIToMasterData(contactPersonDAOName, funcName, masterDataPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	errorS := json.Unmarshal([]byte(strData), &resultRequest)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(contactPersonDAOName, funcName, errorS)
		return
	}
	insertedID = resultRequest.ID

	err = errorModel.GenerateNonErrorModel()
	return
}

func UpdateContactPerson(inputStruct master_data_request.ContactPersonWriteRequest, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		funcName           = "UpdateContactPerson"
		masterDataConf     = config.ApplicationConfiguration.GetMasterData()
		masterDataPath     = masterDataConf.Host + masterDataConf.PathRedirect.ContactPerson.BaseUrl + "/" + strconv.Itoa(int(inputStruct.ID))
		contactPersonOnMDB master_data_response.ViewContactPersonResponse
	)

	//--- Get updated_at contact_person
	fmt.Println("[Request to master data for view] => ", util2.StructToJSON(inputStruct))
	contactPersonOnMDB, err = GetViewContactPerson(master_data_request.ContactPersonGetListRequest{
		ID: inputStruct.ID,
	}, contextModel)

	if err.Error != nil {
		return
	}

	inputStruct.UpdatedAt = contactPersonOnMDB.UpdatedAt

	//--- Update contact_person
	fmt.Println("[Request to master data for update] => ", util2.StructToJSON(inputStruct))
	_, err = hitAPIToMasterData(contactPersonDAOName, funcName, masterDataPath, http.MethodPut, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	fmt.Println("[End of update to master data]")
	return
}

func ValidateContactPersonOnMDB(inputStruct master_data_request.ContactPersonGetListRequest, contextModel *applicationModel.ContextModel) (contactPersonID int64, isValid bool, err errorModel.ErrorModel) {
	var (
		funcName       = "GetListContactPersonFromMasterData"
		masterDataConf = config.ApplicationConfiguration.GetMasterData()
		masterDataPath = masterDataConf.Host + masterDataConf.PathRedirect.ContactPerson.GetList
		result         []master_data_response.ContactPersonResponse
	)

	inputStruct.Page = 1

	fmt.Println("Request to master data: ", util2.StructToJSON(inputStruct))
	strData, err := hitAPIToMasterData(contactPersonDAOName, funcName, masterDataPath, http.MethodPost, contextModel, inputStruct)
	if err.Error != nil {
		return
	}

	fmt.Println("Response from master data: ", strData)
	errorUnmarshal := json.Unmarshal([]byte(strData), &result)
	if errorUnmarshal != nil {
		err = errorModel.GenerateUnknownError(contactPersonDAOName, funcName, errorUnmarshal)
		return
	}

	fmt.Println("Response succeed convert to result: ", len(result))
	if len(result) > 0 {
		isValid = true
		contactPersonID = result[0].ID
	}

	return
}
