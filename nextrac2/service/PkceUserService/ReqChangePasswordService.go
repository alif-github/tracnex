package PkceUserService

import (
	"database/sql"
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
	"nexsoft.co.id/nextrac2/service/UserService/CRUDUserService"
	"nexsoft.co.id/nextrac2/util"
)

func (input reqChangePasswordService) RequestChangePassword(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ChangePassword

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateChangePassword)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doRequestChangePassword(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:        util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message:     GenerateI18NMessage("SUCCESS_CHANGE_PASSWORD_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reqChangePasswordService) doRequestChangePassword(inputStructInterface interface{}, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	fileName := "ReqChangePasswordService.go"
	funcName := "doRequestChangePassword"

	inputStruct := inputStructInterface.(in.ChangePassword)
	var tx *sql.Tx
	var errs error

	tx, errs = serverconfig.ServerAttribute.DBConnection.Begin()
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
		return
	}

	defer func() {
		_ = tx.Commit()
	}()

	//------ Check valid company and branch ID
	var clientMappingModel []repository.ClientMappingModel

	if inputStruct.CompanyID != "" && inputStruct.BranchID != "" {
		clientMappingModel = append(clientMappingModel, repository.ClientMappingModel {
			CompanyID: 	sql.NullString{String: inputStruct.CompanyID},
			BranchID: 	sql.NullString{String: inputStruct.BranchID},
			ClientID: 	sql.NullString{String: inputStruct.ParentClientID},
		})
		clientMappingModel, err = dao.ClientMappingDAO.CheckClientMapping(tx, clientMappingModel, true)
		if err.Error != nil {
			return
		}

		if len(clientMappingModel) < 1 {
			detail := util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.PKCEUserBundle, "DETAIL_ERROR_INVALID_ID_MESSAGE", contextModel.AuthAccessTokenModel.Locale, nil)
			err = errorModel.GenerateFailedChangePassword(fileName, funcName, []string{detail})
			return
		}
	}

	pkceClientMappingModel := repository.PKCEClientMappingModel {
		ClientTypeID: 	sql.NullInt64{Int64: inputStruct.ClientTypeID},
		Username: 		sql.NullString{String: inputStruct.Username},
		ParentClientID:	sql.NullString{String: inputStruct.ParentClientID},
		CompanyID: 		sql.NullString{String: inputStruct.CompanyID},
		BranchID: 		sql.NullString{String: inputStruct.BranchID},
	}

	//------ Change password with own permission
	pkceClientMappingModel.CreatedBy.Int64 = contextModel.LimitedByCreatedBy

	var resultPkceClientMappingModel repository.PKCEClientMappingModel
	resultPkceClientMappingModel, err = dao.PKCEClientMappingDAO.CheckPKCEClientMapping(tx, pkceClientMappingModel)
	if err.Error != nil {
		return
	}

	if resultPkceClientMappingModel.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Username + " : " + inputStruct.Username)
		return
	}

	userChangePasswordData := authentication_request.ChangePasswordDTOin{
		UserID: 			resultPkceClientMappingModel.AuthUserID.Int64,
		OldPassword: 		inputStruct.CurrentPassword,
		NewPassword: 		inputStruct.NewPassword,
		VerifyNewPassword: 	inputStruct.ConfirmationPassword,
	}

	//------ Change password to authentication server
	err = CRUDUserService.UserService.ChangePasswordToAuthenticationServer(userChangePasswordData, contextModel)
	if err.Error != nil {
		return
	}

	//------ Get list token by client ID
	var listToken []string
	listToken, err = dao.ClientTokenDAO.GetListTokenByClientID(tx, resultPkceClientMappingModel.ClientID.String)
	if err.Error != nil {
		return
	}

	//------ Delete token in redis
	go service.DeleteTokenFromRedis(listToken)

	//------ Delete token by client id
	err = dao.ClientTokenDAO.DeleteListTokenByClientID(tx, resultPkceClientMappingModel.ClientID.String)
	if err.Error != nil {
		return
	}

	outputTemp := out.ChangePasswordPKCEResponse{
		UserID: 	resultPkceClientMappingModel.AuthUserID.Int64,
		ClientID: 	resultPkceClientMappingModel.ClientID.String,
		Username: 	inputStruct.Username,
	}

	output = outputTemp
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reqChangePasswordService) validateChangePassword(inputStruct *in.ChangePassword) errorModel.ErrorModel {
	return inputStruct.ValidateReqChangePassword()
}
