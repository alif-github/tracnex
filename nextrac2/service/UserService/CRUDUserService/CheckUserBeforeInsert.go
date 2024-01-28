package CRUDUserService

import (
	"database/sql"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input userService) CheckUserAuth(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.UserRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateCheckUserAuth)
	if err.Error != nil {
		return
	}

	output, err = input.doCheckUserAuth(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) validateCheckUserAuth(inputStruct *in.UserRequest) errorModel.ErrorModel {
	return inputStruct.ValidationCheckUserAuth()
}

func (input userService) doCheckUserAuth(inputStructInterface interface{}, contextModel *applicationModel.ContextModel) (output out.Payload, err errorModel.ErrorModel) {
	var (
		fileName              = "CheckUserBeforeInsert.go"
		funcName              = "doCheckUserAuth"
		checkUserRequest      = input.setRequestForCheckUserAuth(inputStructInterface)
		checkEmailUserRequest = input.setRequestForCheckUserAuthByEmail(inputStructInterface)
		checkPhoneUserRequest = input.setRequestForCheckUserAuthByPhone(inputStructInterface)
		authUser              authentication_response.UserContent
		userTrac              repository.UserModel
	)

	userAuthByEmailAndPhone, err := HitAuthenticateServerForGetDetailUserAuth(checkUserRequest, contextModel)
	if err.Error != nil && err.CausedBy.Error() != constanta.AuthenticationDataNotFound {
		return
	}

	if userAuthByEmailAndPhone.Nexsoft.Payload.Data.Content.UserID < 1 {
		userAuthByEmail, errs := HitAuthenticateServerForGetDetailUserAuth(checkEmailUserRequest, contextModel)
		dataUserAuthByEmail := userAuthByEmail.Nexsoft.Payload.Data.Content
		if errs.Error != nil && errs.CausedBy.Error() != constanta.AuthenticationDataNotFound {
			err = errs
			return
		}

		userAuthByPhone, errs := HitAuthenticateServerForGetDetailUserAuth(checkPhoneUserRequest, contextModel)
		dataUserAuthByPhone := userAuthByPhone.Nexsoft.Payload.Data.Content
		if errs.Error != nil && errs.CausedBy.Error() != constanta.AuthenticationDataNotFound {
			err = errs
			return
		}

		if dataUserAuthByEmail.UserID > 0 && dataUserAuthByPhone.UserID > 0 {
			err = errorModel.GenerateBothEmailAndPhoneAlreadyRegisteredAuth(input.FileName, funcName)
			return
		}

		if dataUserAuthByEmail.UserID > 0 {
			userAuthByEmailAndPhone = userAuthByEmail
		} else if dataUserAuthByPhone.UserID > 0 {
			userAuthByEmailAndPhone = userAuthByPhone
		}
	}

	authUser = userAuthByEmailAndPhone.Nexsoft.Payload.Data.Content

	if authUser.UserID > 0 {
		userTrac, err = dao.UserDAO.CheckIsUserExists(serverconfig.ServerAttribute.DBConnection, repository.UserModel{
			AuthUserID: sql.NullInt64{Int64: authUser.UserID},
		})
		if err.Error != nil {
			return
		}

		if userTrac.ID.Int64 > 0 {
			err = errorModel.GenerateEmailAndPhoneAlreadyRegisteredNextrac(fileName, funcName)
			return
		}

		output.Data.Content = userAuthByEmailAndPhone.Nexsoft.Payload.Data.Content
		output.Status = out.StatusResponse{
			Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
			Message: userService2.GenerateI18NMessage("SUCCESS_USER_ALREADY_REGISTERED_AUTH", contextModel.AuthAccessTokenModel.Locale),
		}

		err = errorModel.GenerateNonErrorModel()
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_USER_NOT_FOUND_AUTH", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) setRequestForCheckUserAuth(inputStructInterface interface{}) (authRequest in.UserRequest) {
	inputStruct := inputStructInterface.(in.UserRequest)
	return in.UserRequest{
		Email: inputStruct.Email,
		Phone: inputStruct.Phone,
	}
}

func (input userService) setRequestForCheckSignatureUserAuth(inputStructInterface interface{}) (authRequest in.UserRequest) {
	var (
		inputStruct = inputStructInterface.(in.UserRequest)
		phoneNumber = fmt.Sprintf(`%s-%s`, inputStruct.CountryCode, inputStruct.Phone)
	)

	return in.UserRequest{
		Email: inputStruct.Email,
		Phone: phoneNumber,
	}
}

func (input userService) setRequestForCheckUserAuthByEmail(inputStructInterface interface{}) (authRequest in.UserRequest) {
	inputStruct := inputStructInterface.(in.UserRequest)
	return in.UserRequest{
		Email: inputStruct.Email,
	}
}

func (input userService) setRequestForCheckUserAuthByPhone(inputStructInterface interface{}) (authRequest in.UserRequest) {
	inputStruct := inputStructInterface.(in.UserRequest)
	return in.UserRequest{
		Phone: inputStruct.Phone,
	}
}
