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
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	"nexsoft.co.id/nextrac2/util"
)

func (input userService) ChangePasswordUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ChangePasswordUserRequestDTO

	inputStruct, err = input.readBodyAndValidateForChangePasswordUser(request, contextModel, input.validateChangePasswordUser)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doChangePasswordUser(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_CHANGE_PASSWORD_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) doChangePasswordUser(inputStructInterface interface{}, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		fileName       = "ChangePassword.go"
		funcName       = "doChangePasswordUser"
		inputStruct    = inputStructInterface.(in.ChangePasswordUserRequestDTO)
		resultDataUser repository.UserModel
		modelDataUser  repository.UserModel
		tx             *sql.Tx
		errs           error
		listToken      []string
	)

	tx, errs = serverconfig.ServerAttribute.DBConnection.Begin()
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
		return
	}

	defer func() {
		_ = tx.Commit()
	}()

	modelDataUser = repository.UserModel{ID: sql.NullInt64{Int64: inputStruct.ID}}
	if contextModel.LimitedByCreatedBy > 0 {
		if contextModel.AuthAccessTokenModel.ResourceUserID != inputStruct.ID {
			err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
			return
		}
	}

	resultDataUser, err = dao.UserDAO.GetUserForUpdate(tx, modelDataUser)
	if err.Error != nil {
		return
	}

	if resultDataUser.Username.String != inputStruct.Username {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, fmt.Sprintf(`%s: %s`, constanta.Username, inputStruct.Username))
		return
	}

	userChangePasswordData := authentication_request.ChangePasswordDTOin{
		UserID:            resultDataUser.AuthUserID.Int64,
		OldPassword:       inputStruct.OldPassword,
		NewPassword:       inputStruct.NewPassword,
		VerifyNewPassword: inputStruct.VerifyNewPassword,
	}

	err = input.ChangePasswordToAuthenticationServer(userChangePasswordData, contextModel)
	if err.Error != nil {
		return
	}

	//--- Get list token by client ID
	listToken, err = dao.ClientTokenDAO.GetListTokenByClientID(tx, resultDataUser.ClientID.String)
	if err.Error != nil {
		return
	}

	//--- Delete token in redis
	go service.DeleteTokenFromRedis(listToken)

	//--- Delete token by client id
	err = dao.ClientTokenDAO.DeleteListTokenByClientID(tx, resultDataUser.ClientID.String)
	if err.Error != nil {
		return
	}

	outputTemp := out.ChangePasswordPKCEResponse{
		UserID:   resultDataUser.ID.Int64,
		ClientID: resultDataUser.ClientID.String,
		Username: resultDataUser.Username.String,
	}

	output = outputTemp
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) validateChangePasswordUser(inputStruct *in.ChangePasswordUserRequestDTO) errorModel.ErrorModel {
	return inputStruct.ValidateChangePassword()
}
