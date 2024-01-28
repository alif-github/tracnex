package ParameterService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
)

func GenerateI18NMessage(messageID string, language string) (output string) {
	//return util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ParameterBundle, messageID, language, nil)
	return util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.AuditMonitoringBundle, messageID, language, nil)
}

func getParameterBody(request *http.Request, fileName string) (parameterBody []in.ParameterRequest, bodySize int, err errorModel.ErrorModel) {
	funcName := "getParameterBody"
	jsonString, bodySize, readError := util.ReadBody(request)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	readError = json.Unmarshal([]byte(jsonString), &parameterBody)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

