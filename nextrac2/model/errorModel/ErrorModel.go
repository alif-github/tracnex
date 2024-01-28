package errorModel

import "errors"

type ErrorModel struct {
	Code                  int
	Error                 error
	FileName              string
	FuncName              string
	CausedBy              error
	ErrorParameter        []ErrorParameter
	AdditionalInformation []string
	OtherData             interface{}
}

type ErrorLogModel struct {
	Code         int
	ErrorMessage string
	FileName     string
	FuncName     string
}

type ErrorParameter struct {
	ErrorParameterKey   string
	ErrorParameterValue string
}

func GenerateErrorModel(code int, err string, fileName string, funcName string, causedBy error) ErrorModel {
	var errModel ErrorModel
	errModel.Code = code
	errModel.Error = errors.New(err)
	errModel.FileName = fileName
	errModel.FuncName = funcName
	errModel.CausedBy = causedBy
	return errModel
}

func GenerateErrorModelWithoutCaused(code int, err string, fileName string, funcName string) ErrorModel {
	var errModel ErrorModel
	errModel.Code = code
	errModel.Error = errors.New(err)
	errModel.FileName = fileName
	errModel.FuncName = funcName
	return errModel
}

func GenerateErrorModelWithErrorParam(code int, err string, fileName string, funcName string, errorParam []ErrorParameter) ErrorModel {
	var errModel ErrorModel
	errModel.Code = code
	errModel.Error = errors.New(err)
	errModel.FileName = fileName
	errModel.FuncName = funcName
	errModel.ErrorParameter = errorParam
	return errModel
}

func GenerateSimpleErrorModel(code int, err string) ErrorModel {
	var errModel ErrorModel
	errModel.Code = code
	errModel.Error = errors.New(err)
	return errModel
}

func GenerateNonErrorModel() ErrorModel {
	var errModel ErrorModel
	errModel.Code = 200
	errModel.Error = nil
	return errModel
}

func GenerateErrorModelWithAdditionalInformation(code int, err string, fileName string, funcName string, additionalInformation []string) ErrorModel {
	var errModel ErrorModel
	errModel.Code = code
	errModel.Error = errors.New(err)
	errModel.FileName = fileName
	errModel.FuncName = funcName
	errModel.AdditionalInformation = additionalInformation
	return errModel
}

func GenerateErrorModelWithAdditionalInformationWithData(code int, err string, fileName string, funcName string, other interface{}) ErrorModel {
	var errModel ErrorModel
	errModel.Code = code
	errModel.Error = errors.New(err)
	errModel.FileName = fileName
	errModel.FuncName = funcName
	errModel.OtherData = other
	return errModel
}

func GenerateErrorModelWithAdditionalInformationAndErrorParam(code int, err string, fileName string, funcName string, additionalInformation []string, errorParam []ErrorParameter) ErrorModel {
	var errModel ErrorModel
	errModel.Code = code
	errModel.Error = errors.New(err)
	errModel.FileName = fileName
	errModel.FuncName = funcName
	errModel.ErrorParameter = errorParam
	errModel.AdditionalInformation = additionalInformation
	return errModel
}

func GenerateErrorModelWithErrorParamAndCaused(code int, err string, fileName string, funcName string, errorParam []ErrorParameter, causedBy error) ErrorModel {
	var errModel ErrorModel
	errModel.Code = code
	errModel.Error = errors.New(err)
	errModel.FileName = fileName
	errModel.FuncName = funcName
	errModel.CausedBy = causedBy
	errModel.ErrorParameter = errorParam
	return errModel
}
