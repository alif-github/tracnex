package master_data_dao

import (
	"encoding/json"
	"errors"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"strconv"
)

func GetListPersonTitleFromMasterData(inputStruct in.PersonTitleRequest, contextModel *applicationModel.ContextModel) (listPersonTitle []out.PersonTitleResponse, err errorModel.ErrorModel) {
	fileName := "PersonTitleDAO.go"
	funcName := "GetListPersonTitleFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	listPersonTitlePath := masterData.Host + masterData.PathRedirect.PersonTitle.GetList

	headerRequest := make(map[string][]string)
	authorize := resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, constanta.DefaultApplicationsLanguage)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	statusCode, _, bodyResult, errorS := common.HitAPI(listPersonTitlePath, headerRequest, util2.StructToJSON(inputStruct), "POST", *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	var listPersonTitleResponse out.GetListPersonTitle
	_ = json.Unmarshal([]byte(bodyResult), &listPersonTitleResponse)

	if statusCode == 200 {
		listPersonTitle = listPersonTitleResponse.Nexsoft.Payload.Data.Content
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(listPersonTitleResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(fileName, funcName, statusCode, listPersonTitleResponse.Nexsoft.Payload.Status.Code, causedBy)
	}

	return
}

func ViewDetailPersonTitleFromMasterData(id int, contextModel *applicationModel.ContextModel) (dataPersonTitle out.PersonTitleResponse, err errorModel.ErrorModel) {
	fileName := "PersonTitleDAO.go"
	funcName := "ViewDetailPersonTitleFromMasterData"
	masterData := config.ApplicationConfiguration.GetMasterData()
	viewPersonTitlePath := masterData.Host + masterData.PathRedirect.PersonTitle.View + "/" + strconv.Itoa(id)

	headerRequest := make(map[string][]string)
	authorize := resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, constanta.DefaultApplicationsLanguage)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{authorize}

	statusCode, _, bodyResult, errorS := common.HitAPI(viewPersonTitlePath, headerRequest, "", "GET", *contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	var viewPersonTitleResponse out.ViewPersonTitle
	_ = json.Unmarshal([]byte(bodyResult), &viewPersonTitleResponse)

	if statusCode == 200 {
		dataPersonTitle = viewPersonTitleResponse.Nexsoft.Payload.Data.Content
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(viewPersonTitleResponse.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(fileName, funcName, statusCode, viewPersonTitleResponse.Nexsoft.Payload.Status.Code, causedBy)
	}

	return
}