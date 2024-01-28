package RegistrationNamedUserService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/service/UserService/CRUDUserService"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
)

func (input registrationNamedUserService) CheckEmailAndPhoneBeforeRegisterNamedUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.CheckNamedUserBeforeInsertRequest
		response    out.CheckEmailAndPhoneBeforeInsertResponse
	)

	inputStruct, err = input.readBodyAndValidateBeforeRegisterNamedUser(request, contextModel, input.validationCheckEmailAndPhoneBeforeRegisterNamedUser)
	if err.Error != nil {
		return
	}

	response, err = input.doCheckEmailAndPhoneBeforeRegisterNamedUser(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content = response.Content
	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n(response.Status.Code, contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage(response.Status.Message, contextModel.AuthAccessTokenModel.Locale),
	}

	if output.Status.Message == GenerateI18NMessage("FAILED_FOUND_DIFFERENT_USER_MESSAGE", contextModel.AuthAccessTokenModel.Locale) {
		output.Status.Code = "E-9FTR-TRAC-SRV-0013"
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) validationCheckEmailAndPhoneBeforeRegisterNamedUser(inputStruct *in.CheckNamedUserBeforeInsertRequest) errorModel.ErrorModel {
	return inputStruct.ValidateCheckNamedUserBeforeInsert()
}

func (input registrationNamedUserService) doCheckEmailAndPhoneBeforeRegisterNamedUser(inputStructInterface interface{}, contextModel *applicationModel.ContextModel) (output out.CheckEmailAndPhoneBeforeInsertResponse, err errorModel.ErrorModel) {
	var (
		inputStruct             = inputStructInterface.(in.CheckNamedUserBeforeInsertRequest)
		userAuthByEmailAndPhone authentication_response.UserAuthenticationResponse
		dataUserByEmailAndPhone authentication_response.UserContent
		userAuthByEmail         authentication_response.UserAuthenticationResponse
		userAuthByPhone         authentication_response.UserAuthenticationResponse
		dataUserAuthByPhone     authentication_response.UserContent
		dataUserAuthByEmail     authentication_response.UserContent
	)

	userAuthByEmailAndPhone, err = CRUDUserService.HitAuthenticateServerForGetDetailUserAuth(in.UserRequest{Email: inputStruct.Email, Phone: inputStruct.Phone}, contextModel)
	dataUserByEmailAndPhone = userAuthByEmailAndPhone.Nexsoft.Payload.Data.Content
	if err.Error != nil && err.CausedBy.Error() != constanta.AuthenticationDataNotFound {
		return
	}

	if dataUserByEmailAndPhone.UserID > 0 {
		output = input.convertResponseSuccessFoundDataByEmailAndPhone(dataUserByEmailAndPhone, contextModel)
		return
	} else {
		userAuthByEmail, err = CRUDUserService.HitAuthenticateServerForGetDetailUserAuth(in.UserRequest{Email: inputStruct.Email}, contextModel)
		dataUserAuthByEmail = userAuthByEmail.Nexsoft.Payload.Data.Content
		if err.Error != nil && err.CausedBy.Error() != constanta.AuthenticationDataNotFound {
			return
		}

		userAuthByPhone, err = CRUDUserService.HitAuthenticateServerForGetDetailUserAuth(in.UserRequest{Phone: inputStruct.Phone}, contextModel)
		dataUserAuthByPhone = userAuthByPhone.Nexsoft.Payload.Data.Content
		if err.Error != nil && err.CausedBy.Error() != constanta.AuthenticationDataNotFound {
			return
		}

		if dataUserAuthByEmail.UserID > 0 && dataUserAuthByPhone.UserID > 0 {
			output = input.convertResponseFailedDifferentDataByEmailAndPhone(dataUserAuthByEmail, dataUserAuthByPhone, contextModel)

			return
		}

		if dataUserAuthByEmail.UserID > 0 {
			output = input.convertResponseSuccessFoundDataByEmail(dataUserAuthByEmail, contextModel)
			err = errorModel.GenerateNonErrorModel()
			return
		} else if dataUserAuthByPhone.UserID > 0 {
			output = input.convertResponseSuccessFoundDataByPhone(dataUserAuthByPhone, contextModel)
			err = errorModel.GenerateNonErrorModel()
			return
		}
	}

	output = input.convertResponseSuccessNotFound()

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input registrationNamedUserService) convertResponseSuccessFoundDataByEmailAndPhone(data authentication_response.UserContent, contextModel *applicationModel.ContextModel) out.CheckEmailAndPhoneBeforeInsertResponse {
	return out.CheckEmailAndPhoneBeforeInsertResponse{
		Status: out.Status{
			Code:    strings.ToUpper(constanta.StatusMessage),
			Message: "SUCCESS_FOUND_DATA_BY_EMAIL_PHONE_MESSAGE",
		},
		Content: []out.CheckEmailAndPhoneResponse{
			{
				Email:      data.Email,
				Phone:      data.Phone,
				AuthUserId: data.UserID,
				Message:    GenerateI18NMessage("SUCCESS_FOUND_DATA_BY_EMAIL_PHONE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
			},
		},
	}
}

func (input registrationNamedUserService) convertResponseSuccessFoundDataByEmail(data authentication_response.UserContent, contextModel *applicationModel.ContextModel) out.CheckEmailAndPhoneBeforeInsertResponse {
	return out.CheckEmailAndPhoneBeforeInsertResponse{
		Status: out.Status{
			Code:    strings.ToUpper(constanta.StatusMessage),
			Message: "SUCCESS_FOUND_DATA_BY_EMAIL_PHONE_MESSAGE",
		},
		Content: []out.CheckEmailAndPhoneResponse{
			{
				Phone:      data.Phone,
				Email:      data.Email,
				AuthUserId: data.UserID,
				Message:    GenerateI18NMessage("SUCCESS_FOUND_DATA_BY_EMAIL_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
			},
		},
	}
}

func (input registrationNamedUserService) convertResponseSuccessFoundDataByPhone(data authentication_response.UserContent, contextModel *applicationModel.ContextModel) out.CheckEmailAndPhoneBeforeInsertResponse {
	return out.CheckEmailAndPhoneBeforeInsertResponse{
		Status: out.Status{
			Code:    strings.ToUpper(constanta.StatusMessage),
			Message: "SUCCESS_FOUND_DATA_BY_EMAIL_PHONE_MESSAGE",
		},
		Content: []out.CheckEmailAndPhoneResponse{
			{
				Email:      data.Email,
				Phone:      data.Phone,
				AuthUserId: data.UserID,
				Message:    GenerateI18NMessage("SUCCESS_FOUND_DATA_BY_PHONE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
			},
		},
	}
}

func (input registrationNamedUserService) convertResponseFailedDifferentDataByEmailAndPhone(dataAuthByEmail authentication_response.UserContent, dataAuthByPhone authentication_response.UserContent, contextModel *applicationModel.ContextModel) out.CheckEmailAndPhoneBeforeInsertResponse {
	return out.CheckEmailAndPhoneBeforeInsertResponse{
		Status: out.Status{
			Code:    strings.ToUpper(constanta.StatusMessage),
			Message: "FAILED_FOUND_DIFFERENT_USER_MESSAGE",
		},
		Content: []out.CheckEmailAndPhoneResponse{
			{
				Email:      dataAuthByEmail.Email,
				Phone:      dataAuthByEmail.Phone,
				AuthUserId: dataAuthByEmail.UserID,
				Message:    GenerateI18NMessage("FAILED_FOUND_DIFFERENT_USER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
			},
			{
				Email:      dataAuthByPhone.Email,
				Phone:      dataAuthByPhone.Phone,
				AuthUserId: dataAuthByPhone.UserID,
				Message:    GenerateI18NMessage("FAILED_FOUND_DIFFERENT_USER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
			},
		},
	}
}

func (input registrationNamedUserService) convertResponseSuccessNotFound() out.CheckEmailAndPhoneBeforeInsertResponse {
	return out.CheckEmailAndPhoneBeforeInsertResponse{
		Status: out.Status{
			Code:    strings.ToUpper(constanta.StatusMessage),
			Message: "SUCCESS_NOT_FOUND_DATA_MESSAGE",
		},
		Content: nil,
	}
}
