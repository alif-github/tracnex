package RegistrationNamedUserService

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
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input registrationNamedUserService) UnregisterNamedUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "UnregisterNamedUser"
	var userParam in.UnregisterNamedUserRequest

	if userParam, err = input.readBodyAndValidateForUnregisterRequest(request, input.validateUnregisterNamedUser); err.Error != nil {
		return
	}

	if output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, userParam, contextModel, input.doUnregisterNamedUser, func(interface{}, applicationModel.ContextModel) {}); err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_UNREGISTER_LICENSE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input registrationNamedUserService) validateUnregisterNamedUser(inputStruct *in.UnregisterNamedUserRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUnregisterNamedUser()
}

func (input registrationNamedUserService) doUnregisterNamedUser(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doUnregisterNamedUser"
	fileName := "UnregisterNamedUserService.go"

	inputStruct := inputStructInterface.(in.UnregisterNamedUserRequest)
	var userRegDetailOnDB repository.UserRegistrationDetailModel

	userRegDetailModel := repository.UserRegistrationDetailModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		Status:        sql.NullString{String: constanta.NonactiveUser},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	// Get ParentClientId For Client Validation
	var userLicenseOnDB repository.UserLicenseModel
	if userLicenseOnDB, err = dao.UserLicenseDAO.GetFieldForValidationUnregister(serverconfig.ServerAttribute.DBConnection, userRegDetailModel); err.Error != nil {
		return
	}

	if userLicenseOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.UserRegistrationDetailID)
		return
	}

	// Step 1 - Client Validation
	if userLicenseOnDB.ClientID.String != contextModel.AuthAccessTokenModel.ClientID {
		err = errorModel.GenerateForbiddenClientCredentialAccess(fileName, funcName)
		return
	}

	if userRegDetailOnDB, err = dao.UserRegistrationDetailDAO.GetUserRegistrationDetailForUnregister(serverconfig.ServerAttribute.DBConnection, userRegDetailModel); err.Error != nil {
		return
	}

	// Step 2 - Validation Total License Named User
	if userLicenseOnDB.TotalActivated.Int64 < 1 {
		err = errorModel.GenerateTotalActivatedZeroValue(fileName, funcName)
		return
	}

	// Step 3 - Check Status
	if userRegDetailOnDB.Status.String == "R" {
		err = errorModel.GenerateLicenseHasNotBeenActivated(fileName, funcName)
		return
	} else if userRegDetailOnDB.Status.String == "N" {
		err = errorModel.GenerateLicenseHasBeenDeactivated(fileName, funcName)
		return
	}

	// Step 4 - Deactivated License Named User
	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserRegistrationDetailDAO.TableName, userRegDetailOnDB.ID.Int64, contextModel.LimitedByCreatedBy)...)

	if err = dao.UserRegistrationDetailDAO.UnregisterNamedUser(tx, userRegDetailModel); err.Error != nil {
		return
	}

	// Step 5 - Reduce Total License
	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserLicenseDAO.TableName, userRegDetailOnDB.UserLicenseID.Int64, contextModel.LimitedByCreatedBy)...)

	userRegDetailModel.UserLicenseID = userRegDetailOnDB.UserLicenseID

	if err = dao.UserLicenseDAO.ReduceTotalLicense(tx, userRegDetailModel); err.Error != nil {
		return
	}

	// Step 6 - Deactivated User
	userOnDB, err := dao.UserDAO.GetUserByClientID(serverconfig.ServerAttribute.DBConnection, repository.UserModel{
		ClientID: userRegDetailOnDB.ClientID,
	})
	if err.Error != nil {
		return
	}

	if userOnDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.User)
		return
	}

	err = resource_common_service.InternalDeleteClientByClientID(userRegDetailOnDB.ClientID.String, contextModel)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserDAO.TableName, userOnDB.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.UserDAO.UpdateUserStatus(tx, repository.UserModel{
		AuthUserID:    userRegDetailOnDB.AuthUserID,
		UpdatedClient: userRegDetailModel.UpdatedClient,
		UpdatedAt:     userRegDetailModel.UpdatedAt,
		UpdatedBy:     userRegDetailModel.UpdatedBy,
		Status:        sql.NullString{String: constanta.StatusNonActive},
	})
	if err.Error != nil {
		return
	}

	// Step 7 - Delete data user verification
	userVerifOnDB, err := dao.UserVerificationDAO.GetUserVerificationForUnregister(serverconfig.ServerAttribute.DBConnection, repository.UserVerificationModel{
		UserRegistrationDetailID: sql.NullInt64{Int64: inputStruct.ID},
	})
	if err.Error != nil {
		return
	}

	err = dao.UserVerificationDAO.HardDeleteUserVerification(tx, userVerifOnDB)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
