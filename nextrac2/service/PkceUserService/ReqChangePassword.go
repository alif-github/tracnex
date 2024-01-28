package PkceUserService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type reqChangePasswordService struct {
	service.RegistrationPrepared
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var ReqChangePasswordService = reqChangePasswordService{}.New()

func (input reqChangePasswordService) New() (output reqChangePasswordService) {
	output.FileName = "ChangePassword.go"
	output.ValidLimit = []int{10, 20, 50, 100}
	output.ValidOrderBy = []string{"id"}
	output.IdResourceAllowed = []int64{2}
	return
}

func (input reqChangePasswordService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(changePassword *in.ChangePassword) errorModel.ErrorModel) (inputStruct in.ChangePassword, err errorModel.ErrorModel) {
	funcName := "RequestChangePassword"
	var stringBody string
	var isAllowed bool

	//---------- Read Body Request
	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	//---------- Unmarshal String Body to Main struct
	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	//---------- Init, Important request must be exist
	if inputStruct.ClientTypeID == 0 {
		err = errorModel.GenerateEmptyFieldError(input.FileName, funcName, constanta.ClientTypeID)
		return
	}

	//---------- Check to DB, client type exist on table ?
	preparedError := errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientTypeID)
	err = input.CheckIsClientTypeExist(inputStruct.ClientTypeID, preparedError)
	if err.Error != nil {
		return
	}

	//---------- Check client type allowing
	for _, idResourceItem := range input.IdResourceAllowed {
		if idResourceItem == inputStruct.ClientTypeID {
			isAllowed = true
			break
		}
	}

	//---------- Is not allowed, then forbidden to access
	if !isAllowed {
		err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
		return
	}

	//---------- Login must be same with parent client ID
	isNexmile := inputStruct.ClientTypeID == constanta.ResourceNexmileID
	if isNexmile {
		if inputStruct.ParentClientID == "" {
			err = errorModel.GenerateEmptyFieldError(input.FileName, funcName, constanta.ParentClientID)
			return
		}

		if contextModel.AuthAccessTokenModel.ClientID != inputStruct.ParentClientID {
			err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
			return
		}
	}

	//---------- Main validation
	err = validation(&inputStruct)
	return
}