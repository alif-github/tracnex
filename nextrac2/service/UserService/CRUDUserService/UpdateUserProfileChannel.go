package CRUDUserService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	"nexsoft.co.id/nextrac2/util"
)

func (input userService) UpdateAdminProfile(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	return input.updateProfileService(request, contextModel, true)
}

func (input userService) UpdateUserProfile(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	return input.updateProfileService(request, contextModel, false)
}

func (input userService) updateProfileService(request *http.Request, contextModel *applicationModel.ContextModel, isAdmin bool) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "updateProfileService"
		inputStruct in.UserRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateUserProfile)
	if err.Error != nil {
		return
	}

	inputStruct.IsAdmin = isAdmin
	inputStruct.ID = contextModel.AuthAccessTokenModel.ResourceUserID
	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateProfile, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_UPDATE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
