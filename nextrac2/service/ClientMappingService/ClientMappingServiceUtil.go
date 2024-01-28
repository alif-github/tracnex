package ClientMappingService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	utilCommon "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
	"strconv"
)

func GenerateI18Message(messageID string, language string) (output string) {
	return util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ClientMappingBundle, messageID, language, nil)
}

func GetClientMappingBodies(request *http.Request, fileName string) (clients in.ClientMappingRequest , bodySize int, err errorModel.ErrorModel) {
	funcName := "GetClientMappingBodies"
	jsonString, bodySize, readError := utilCommon.ReadBody(request)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	readError = json.Unmarshal([]byte(jsonString), &clients)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func GetClientMappingBody(request *http.Request, fileName string) (clientMappingBody in.ClientMappingRequest , bodySize int, err errorModel.ErrorModel) {
	funcName := "GetClientMappingBody"
	jsonString, bodySize, readError := utilCommon.ReadBody(request)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	readError = json.Unmarshal([]byte(jsonString), &clientMappingBody)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func readPathParam(request *http.Request) (id int64, err errorModel.ErrorModel) {
	funcName := "readPathParam"

	strId, ok := mux.Vars(request)["ID"]
	idParam, errConvert := strconv.Atoi(strId)
	id = int64(idParam)

	if !ok || errConvert != nil {
		err = errorModel.GenerateUnsupportedRequestParam("ComplaintSubServiceUtil.go", funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}