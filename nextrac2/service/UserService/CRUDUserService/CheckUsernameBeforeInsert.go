package CRUDUserService

import (
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input userService) CheckUsernameAuth(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.UserRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateCheckUsernameAuth)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doCheckUsernameAuth(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_CHECK_USERNAME_AUTH", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) validateCheckUsernameAuth(inputStruct *in.UserRequest) errorModel.ErrorModel {
	return inputStruct.ValidationCheckUsernameAuth()
}

func (input userService) doCheckUsernameAuth(inputStructInterface interface{}, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		funcName             = "doCheckUsernameAuth"
		checkUsernameRequest = input.setRequestForCheckUsernameAuth(inputStructInterface)
	)

	userAuthByUsername, err := HitAuthenticateServerForGetDetailUserAuth(checkUsernameRequest, contextModel)
	if err.Error != nil && err.CausedBy.Error() != constanta.AuthenticationDataNotFound {
		return
	}

	if !util.IsStringEmpty(userAuthByUsername.Nexsoft.Payload.Data.Content.Username) {
		err = errorModel.GenerateUsernameAlreadyUsed(input.FileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) setRequestForCheckUsernameAuth(inputStructInterface interface{}) (userRequest in.UserRequest) {
	return in.UserRequest{
		Username: inputStructInterface.(in.UserRequest).Username,
	}
}
