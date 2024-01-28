package LicenseConfigService

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
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input licenseConfigService) UpdateLicenseConfig(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateLicenseConfig"
		inputStruct in.LicenseConfigRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdate)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateLicenseConfig, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_UPDATE_LICENSE_CONFIG_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) doUpdateLicenseConfig(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName          = "UpdateLicenseConfigService.go"
		funcName          = "doUpdateLicenseConfig"
		inputStruct       = inputStructInterface.(in.LicenseConfigRequest)
		licenseConfigOnDB repository.LicenseConfigModel
		scopeLimit        map[string]interface{}
	)

	licenseConfigModel := repository.LicenseConfigModel{
		ID:              sql.NullInt64{Int64: inputStruct.ID},
		AllowActivation: sql.NullString{String: "Y"},
		UpdatedBy:       sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:   sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:       sql.NullTime{Time: timeNow},
	}

	licenseConfigModel.CreatedBy.Int64 = 0

	//--- Changing because target on product not client type directly
	input.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "pr.client_type_id",
		Count: "pr.client_type_id",
	}

	//--- Validate data scope
	scopeLimit, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	licenseConfigOnDB, err = dao.LicenseConfigDAO.GetLicenseConfigForUpdate(serverconfig.ServerAttribute.DBConnection, licenseConfigModel, scopeLimit, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if licenseConfigOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
		return
	}

	if licenseConfigOnDB.AllowActivation.String == "Y" {
		err = errorModel.GenerateLicenseHasBeenActivated(fileName, funcName)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, licenseConfigOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if licenseConfigOnDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.LicenseConfig)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.LicenseConfigDAO.TableName, licenseConfigOnDB.ID.Int64, 0)...)
	err = dao.LicenseConfigDAO.ActivatingLicenseConfig(tx, licenseConfigModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) validateUpdate(inputStruct *in.LicenseConfigRequest) errorModel.ErrorModel {
	return inputStruct.ValidationUpdateLicenseConfig()
}
