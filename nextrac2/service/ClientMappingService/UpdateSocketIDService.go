package ClientMappingService

import (
	"encoding/json"
	"net/http"
	utilCommon "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/SocketIDService"
	"nexsoft.co.id/nextrac2/util"
)

type updateSocketIDService struct {
	service.AbstractService
}

var UpdateSocketIDService = updateSocketIDService{}.New()

func (input updateSocketIDService) New() (output updateSocketIDService) {
	output.FileName = "UpdateSocketIDService.go"
	return
}

func (input updateSocketIDService) UpdateSocketID(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	clientMappingBody, err := input.readParamAndBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService("UpdateSocketID", clientMappingBody, contextModel, SocketIDService.DoUpdateSocketID, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code: 		util.GenerateConstantaI18n("OK", contextModel.AuthAccessTokenModel.Locale, nil),
		Message:	GenerateI18Message("SUCCESS_UPDATE_SOCKET_ID_ND6_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	return
}

func (input updateSocketIDService) readParamAndBody(request *http.Request, contextModel *applicationModel.ContextModel) (clientMappingBody in.ClientMappingForUIRequest, err errorModel.ErrorModel) {
	funcName := "readParamAndBody"

	jsonString, bodySize, readError := utilCommon.ReadBody(request)
	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, readError)
		return
	}
	contextModel.LoggerModel.ByteIn = bodySize

	readError = json.Unmarshal([]byte(jsonString), &clientMappingBody)
	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, readError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}