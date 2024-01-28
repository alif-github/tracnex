package LicenseVariantService

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
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input licenseVariantService) UpdateLicenseVariant(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "UpdateLicenseVariant"
	var inputStruct in.LicenseVariantRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateLicenseVariant)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdatedLicenseVariant, func(interface{}, applicationModel.ContextModel) {
		// additional function
	})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("OK", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_UPDATE_LICENSE_VARIANT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	return
}

func (input licenseVariantService) doUpdatedLicenseVariant(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	fileName := "UpdateLicenseVariantService.go"
	funcName := "doUpdateLicenseVariant"

	inputStruct := inputStructInterface.(in.LicenseVariantRequest)
	var licenseVariantOnDB repository.LicenseVariantModel

	licenseVariantModel := repository.LicenseVariantModel{
		ID:                 sql.NullInt64{Int64: inputStruct.ID},
		LicenseVariantName: sql.NullString{String: inputStruct.LicenseVariantName},
		UpdatedBy:          sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:      sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:          sql.NullTime{Time: timeNow},
	}

	licenseVariantModel.CreatedBy.Int64 = 0

	licenseVariantOnDB, err = dao.LicenseVariantDAO.GetLicenseVariantForUpdate(serverconfig.ServerAttribute.DBConnection, licenseVariantModel)
	if err.Error != nil {
		return
	}

	if licenseVariantOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.LicenseVariant)
		return
	}

	if contextModel.LimitedByCreatedBy > 0 && contextModel.LimitedByCreatedBy != licenseVariantOnDB.CreatedBy.Int64 {
		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		return
	}

	if licenseVariantOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(fileName, funcName, constanta.LicenseVariant)
		return
	}

	licenseVariantModel.CreatedBy.Int64 = contextModel.LimitedByCreatedBy

	if inputStruct.UpdatedAt != licenseVariantOnDB.UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.UpdatedAt)
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.LicenseVariantDAO.TableName, licenseVariantOnDB.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.LicenseVariantDAO.UpdateLicenseVariant(tx, licenseVariantModel)
	if err.Error != nil {
		err = checkDuplicateError(err)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseVariantService) validateUpdateLicenseVariant(inputStruct *in.LicenseVariantRequest) errorModel.ErrorModel {
	return inputStruct.ValidationUpdateLicenseVariant()
}
