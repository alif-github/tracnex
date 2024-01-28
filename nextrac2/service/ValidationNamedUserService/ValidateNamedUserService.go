package ValidationNamedUserService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service/UserService/CRUDUserService"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
)

func (input validationNamedUser) ValidateNamedUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ValidationNamedUserRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateRequestValidationNamedUser)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doValidateNamedUser(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VALIDATE_NAMED_USER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input validationNamedUser) validateRequestValidationNamedUser(inputStruct *in.ValidationNamedUserRequest) errorModel.ErrorModel {
	return inputStruct.ValidationDTONamedUser()
}

func (input validationNamedUser) doValidateNamedUser(inputStruct in.ValidationNamedUserRequest, contextModel *applicationModel.ContextModel) (output out.ValidationNamedUserResponse, err errorModel.ErrorModel) {
	var (
		fileName                    = "ValidateNamedUserService.go"
		funcName                    = "doValidateNamedUser"
		resourceTrac                = config.ApplicationConfiguration.GetServerResourceID()
		responseUserAuth            authentication_response.UserAuthenticationResponse
		userRegistrationDetailOnDb  repository.UserRegistrationDetailModel
		userRegistrationDetailModel repository.UserRegistrationDetailModel
	)

	// Convert into model for separate optional and mandatory field request
	userRegistrationDetailModel = input.convertInputToModel(inputStruct)

	// Get Data in User Registration Detail
	userRegistrationDetailOnDb, err = dao.UserRegistrationDetailDAO.GetDataForValidateNamedUser(serverconfig.ServerAttribute.DBConnection, userRegistrationDetailModel)
	if err.Error != nil {
		return
	}

	if userRegistrationDetailOnDb.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.UserRegistrationDetailID)
		return
	}

	// Cek resource nextrac
	responseUserAuth, err = CRUDUserService.HitAuthenticateServerForGetDetailUserAuth(in.UserRequest{ClientID: userRegistrationDetailOnDb.ClientID.String}, contextModel)
	if err.Error != nil {
		return
	}

	if !strings.Contains(responseUserAuth.Nexsoft.Payload.Data.Content.ResourceID, resourceTrac) {
		err = errorModel.GenerateUserAuthHasNoTracResource(fileName, funcName)
		return
	}

	output = input.convertToResponseValidateNamedUser(userRegistrationDetailOnDb)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input validationNamedUser) convertInputToModel(inputStruct in.ValidationNamedUserRequest) (result repository.UserRegistrationDetailModel) {
	return repository.UserRegistrationDetailModel{
		ClientID:   sql.NullString{String: inputStruct.ClientId, Valid: !util.IsStringEmpty(inputStruct.ClientId)},
		UniqueID1:  sql.NullString{String: inputStruct.UniqueId1, Valid: !util.IsStringEmpty(inputStruct.UniqueId1)},
		UniqueID2:  sql.NullString{String: inputStruct.UniqueId2, Valid: !util.IsStringEmpty(inputStruct.UniqueId2)},
		AuthUserID: sql.NullInt64{Int64: inputStruct.AuthUserId, Valid: !util2.IsFieldNumericEmpty(inputStruct.AuthUserId)},
		//UserID:     sql.NullString{String: inputStruct.UserId, Valid: !util.IsStringEmpty(inputStruct.UserId)},
	}
}

func (input validationNamedUser) convertToResponseValidateNamedUser(userRegistrationDetailOnDb repository.UserRegistrationDetailModel) out.ValidationNamedUserResponse {
	return out.ValidationNamedUserResponse{
		ProductValidFrom: userRegistrationDetailOnDb.ProductValidFrom.Time,
		ProductValidThru: userRegistrationDetailOnDb.ProductValidThru.Time,
		Status:           userRegistrationDetailOnDb.Status.String,
	}
}
