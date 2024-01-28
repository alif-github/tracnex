package EmployeeFacilitiesActiveService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

func getEmpMatrixBody(request *http.Request, fileName string) (inputStruct in.EmployeeMatrixRequest, bodySize int, err errorModel.ErrorModel) {
	funcName := "getEmpMatrixBody"
	jsonString, bodySize, readError := util.ReadBody(request)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	readError = json.Unmarshal([]byte(jsonString), &inputStruct)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

