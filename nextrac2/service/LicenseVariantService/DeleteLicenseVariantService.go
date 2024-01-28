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
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input licenseVariantService) DeleteLicenseVariant(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "DeleteLicenseVariant"
	var inputStruct in.LicenseVariantRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDelete)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteLicenseVariant, func(_ interface{}, _ applicationModel.ContextModel) {
		// additional function
	})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_DELETE_LICENSE_VARIANT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseVariantService) doDeleteLicenseVariant(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	fileName := "DeleteLicenseVariantService.go"
	funcName := "doDeleteLicenseVariant"

	inputStruct := inputStructInterface.(in.LicenseVariantRequest)
	var licenseVariantOnDB repository.LicenseVariantModel

	licenseVariantModel := repository.LicenseVariantModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
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

	if licenseVariantOnDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.LicenseVariant)
		return
	}

	// ----------- Update for delete
	encodedStr, errorS := service.RandToken(8)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	licenseVariantModel.LicenseVariantName.String = licenseVariantOnDB.LicenseVariantName.String + encodedStr

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.LicenseVariantDAO.TableName, licenseVariantOnDB.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.LicenseVariantDAO.DeleteLicenseVariant(tx, licenseVariantModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseVariantService) validateDelete(inputStruct *in.LicenseVariantRequest) errorModel.ErrorModel {
	return inputStruct.ValidationDeleteLicenseVariant()
}
