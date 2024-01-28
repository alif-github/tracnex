package OTPGeneratorService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type otpGeneratorService struct {
	service.AbstractService
	service.GetListData
}

var OTPGeneratorService = otpGeneratorService{}.New()

func (input otpGeneratorService) New() (output otpGeneratorService) {
	output.FileName = "OTPGeneratorService.go"
	output.ServiceName = "OTP Generator"
	return
}

func (input otpGeneratorService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.OTPGenerator) errorModel.ErrorModel) (inputStruct in.OTPGenerator, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidate"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	err = validation(&inputStruct)
	return
}

func (input otpGeneratorService) GenerateOTP(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.OTPGenerator
		timeNow     = time.Now()
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateGenerateOTP)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGenerateOTP(inputStruct, timeNow)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input otpGeneratorService) doGenerateOTP(inputStruct in.OTPGenerator, timeNow time.Time) (result interface{}, err errorModel.ErrorModel) {
	var (
		//funcName   = "doGenerateOTP"
		timeNowStr = timeNow.Format(constanta.DefaultTimeFormat)
	)

	timeNow, _ = time.Parse(constanta.DefaultTimeFormat, timeNowStr)

	//--- Todo Get WhiteList On DB (Grochat)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input otpGeneratorService) validateGenerateOTP(inputStruct *in.OTPGenerator) errorModel.ErrorModel {
	return inputStruct.ValidateOTPGenerator()
}
