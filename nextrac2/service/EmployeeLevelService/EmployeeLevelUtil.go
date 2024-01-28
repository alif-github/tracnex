package EmployeeLevelService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

func getLevelBody(request *http.Request, fileName string) (levelBody in.EmployeeLevelRequest, bodySize int, err errorModel.ErrorModel) {
	funcName := "getLevelBody"
	jsonString, bodySize, readError := util.ReadBody(request)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	readError = json.Unmarshal([]byte(jsonString), &levelBody)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
