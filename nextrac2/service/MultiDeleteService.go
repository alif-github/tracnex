package service

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type MultiDeleteData struct {
}

func (input MultiDeleteData) ReadMultiDeleteData(request *http.Request, contextModel *applicationModel.ContextModel) (result in.MultiDeleteRequest, err errorModel.ErrorModel) {
	fileName := "MultiDeleteService.go"
	funcName := "readGetListData"

	var stringBody string
	var errorS error

	stringBody, contextModel.LoggerModel.ByteIn, errorS = util.ReadBody(request)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, errorS)
		return
	}

	if contextModel.IsSignatureCheck {
		digest := util.GenerateMessageDigest(stringBody)
		if !ValidateSignature(digest, contextModel.AuthAccessTokenModel.SignatureKey, request) {
			err = errorModel.GenerateInvalidSignatureError(fileName, funcName)
			return
		}
	}

	_ = json.Unmarshal([]byte(stringBody), &result)
	return
}
