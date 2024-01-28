package UserRegistrationService

import (
	"database/sql"
	"encoding/json"
	"net/http"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
)

type userRegistrationService struct {
	service.AbstractService
}

var UserRegistrationService = userRegistrationService{}.New()

func (input userRegistrationService) New() (output userRegistrationService) {
	output.FileName = "UserRegistrationService.go"
	return
}

func (input userRegistrationService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validate func(input *in.CheckLicenseNamedUserRequest) errorModel.ErrorModel) (inputStruct in.CheckLicenseNamedUserRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	if stringBody, err = input.ReadBody(request, contextModel); err.Error != nil {
		return
	}

	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	err = validate(&inputStruct)
	if err.Error != nil {
		return
	}

	return
}

func (input userRegistrationService) CheckLicenseNamedUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	fileName := "UserRegistrationService.go"
	funcName := "CheckLicenseNamedUser"
	var userLicenseOnDb, inputModel repository.UserLicenseModel
	var inputStruct in.CheckLicenseNamedUserRequest

	if inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateCheckLicense); err.Error != nil {
		return
	}

	inputModel = repository.UserLicenseModel{
		ClientID:     sql.NullString{String: inputStruct.ClientId, Valid: true},
		ClientTypeId: sql.NullInt64{Int64: inputStruct.ClientTypeID, Valid: true},
		UniqueId1:    sql.NullString{String: inputStruct.UniqueId1, Valid: true},
	}

	if !util2.IsStringEmpty(inputStruct.UniqueId2) {
		inputModel.UniqueId2.String = inputStruct.UniqueId2
		inputModel.UniqueId2.Valid = true
	} else {
		inputModel.UniqueId2.Valid = false
	}

	userLicenseOnDb, err = dao.UserLicenseDAO.CheckLicenseNamedUser(serverconfig.ServerAttribute.DBConnection, inputModel)
	if userLicenseOnDb.ID.Int64 < 1 {
		err = errorModel.GenerateUserLicenseNotFound(fileName, funcName)
		return
	}

	if userLicenseOnDb.TotalLicense.Int64 <= userLicenseOnDb.TotalActivated.Int64 {
		err = errorModel.GenerateUserLicenseFullFilled(fileName, funcName)
		return
	}

	output.Data.Content = input.convertLicenseNamedUserToDTOOut(userLicenseOnDb)
	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("OK", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_CHECK_LICENSE_NAMED_USER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userRegistrationService) convertLicenseNamedUserToDTOOut(userLicenseOnDB repository.UserLicenseModel) out.CheckUserLicenseNamedUserResponse {
	return out.CheckUserLicenseNamedUserResponse{
		ID:             userLicenseOnDB.ID.Int64,
		ProductKey:     userLicenseOnDB.ProductKey.String,
		TotalLicense:   userLicenseOnDB.TotalLicense.Int64,
		TotalActivated: userLicenseOnDB.TotalActivated.Int64,
		QuotaLicense:   userLicenseOnDB.QuotaLicense.Int64,
	}
}

func (input userRegistrationService) validateCheckLicense(inputStruct *in.CheckLicenseNamedUserRequest) (err errorModel.ErrorModel) {
	return inputStruct.ValidateCheckUserLicense()
}
