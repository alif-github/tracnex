package ForgetPassword

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

type forgetPasswordService struct {
	FileName string
}

var ForgetPasswordService = forgetPasswordService{}.New()

func (input forgetPasswordService) New() (output forgetPasswordService) {
	output.FileName = "ForgetPassword.go"
	return
}

func (input forgetPasswordService) readBodyForResetPassword(request *http.Request, validation func(input *in.ResetPasswordRequest) errorModel.ErrorModel) (inputStruct in.ResetPasswordRequest, err errorModel.ErrorModel) {
	var (
		funcName = "readBodyForResetPassword"
		fileName = input.FileName
	)

	stringJSON, _, errs := util.ReadBody(request)
	if errs != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, errs)
		return
	}

	errs = json.Unmarshal([]byte(stringJSON), &inputStruct)
	if errs != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, errs)
		return
	}

	if err = validation(&inputStruct); err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input forgetPasswordService) readBodyForChangePassword(request *http.Request, validation func(input *in.ChangePasswordByEmailRequest) errorModel.ErrorModel) (inputStruct in.ChangePasswordByEmailRequest, err errorModel.ErrorModel) {
	var (
		funcName = "readBodyForChangePassword"
		fileName = input.FileName
	)

	stringJSON, _, errs := util.ReadBody(request)
	if errs != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, errs)
		return
	}

	errs = json.Unmarshal([]byte(stringJSON), &inputStruct)
	if errs != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, errs)
		return
	}

	if err = validation(&inputStruct); err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
