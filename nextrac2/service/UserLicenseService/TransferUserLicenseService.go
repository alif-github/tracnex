package UserLicenseService

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
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/LicenseConfigService"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"time"
)

func (input userLicenseService) TransferUserLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "TransferUserLicense"
	var inputStruct in.TransferUserLicenseRequest

	inputStruct, err = input.readBodyAndParamTransferKey(request, contextModel, input.validateTransferKey)
	if err.Error != nil {
		return
	}

	outputService, err := input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doTransferUserLicense, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	outputContent := make(map[string]interface{})
	outputContent["license_config_id"] = outputService

	output.Data.Content = outputContent

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_TRANSFER_LICENSE", contextModel.AuthAccessTokenModel.Locale),
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseService) doTransferUserLicense(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var userLicenseOnDB repository.UserLicenseModel
	var customerInstallationOnDB repository.CustomerInstallationModel
	var transferredLicenseConfig repository.LicenseConfigModel
	var inputLicenseConfig in.LicenseConfigRequest

	funcName := "doTransferUserLicense"
	inputStruct := inputStructInterface.(in.TransferUserLicenseRequest)

	if userLicenseOnDB, err = dao.UserLicenseDAO.GetActiveUserLicenseForUpdate(tx, repository.UserLicenseModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	}); err.Error != nil {
		return
	}

	if userLicenseOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.UserLicense)
		return
	}

	if userLicenseOnDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.UserLicense)
		return
	}
	userLicenseOnDB.QuotaLicense.Int64 = userLicenseOnDB.TotalLicense.Int64 - userLicenseOnDB.TotalActivated.Int64

	if inputStruct.NoOfUser > userLicenseOnDB.QuotaLicense.Int64 {
		err = errorModel.GenerateFieldFormatWithRuleError("AbstractDTO.go", "ValidateMinMaxString", "NEED_LESS_THAN", constanta.NumberOfUser, strconv.Itoa(int(userLicenseOnDB.QuotaLicense.Int64)))
		return
	}

	// Get license config
	if transferredLicenseConfig, err = dao.LicenseConfigDAO.GetLicenseConfigForTransferUserLicense(serverconfig.ServerAttribute.DBConnection, repository.LicenseConfigModel{
		ID: userLicenseOnDB.LicenseConfigId,
	}); err.Error != nil {
		return
	}

	if transferredLicenseConfig.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.LicenseConfig)
		return
	}

	// Validate installation id
	if customerInstallationOnDB, err = dao.CustomerInstallationDAO.GetCustomerInstallationByClientIDAndID(serverconfig.ServerAttribute.DBConnection, repository.CustomerInstallationForConfig{
		ID:               sql.NullInt64{Int64: inputStruct.InstallationID},
		ParentCustomerID: transferredLicenseConfig.ParentCustomerID,
		ClientTypeID:     userLicenseOnDB.ClientTypeId,
	}); err.Error != nil {
		return
	}

	if customerInstallationOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.CustomerInstallation)
		return
	}

	inputLicenseConfig = in.LicenseConfigRequest{
		InstallationID:   inputStruct.InstallationID,
		NoOfUser:         inputStruct.NoOfUser,
		IsUserConcurrent: transferredLicenseConfig.IsUserConcurrent.String,
		ProductValidFrom: transferredLicenseConfig.ProductValidFrom.Time,
		ProductValidThru: transferredLicenseConfig.ProductValidThru.Time,
	}

	if output, dataAudit, _, err = LicenseConfigService.LicenseConfigService.DoInsertLicenseConfig(tx, inputLicenseConfig, contextModel, timeNow); err.Error != nil {
		return
	}

	totalLicense := userLicenseOnDB.TotalLicense.Int64 - inputStruct.NoOfUser
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserLicenseDAO.TableName, userLicenseOnDB.ID.Int64, 0)...)
	if err = dao.UserLicenseDAO.UpdateTotalLicenseUserLicense(tx, repository.UserLicenseModel{
		ID:            userLicenseOnDB.ID,
		TotalLicense:  sql.NullInt64{Int64: totalLicense},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}); err.Error != nil {
		return
	}

	return
}

func (input userLicenseService) validateTransferKey(inputStruct *in.TransferUserLicenseRequest) errorModel.ErrorModel {
	return inputStruct.ValidateTransferredUser()
}
