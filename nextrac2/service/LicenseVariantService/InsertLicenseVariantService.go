package LicenseVariantService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input licenseVariantService) InsertLicenseVariant(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "InsertLicenseVariant"

	var inputStruct in.LicenseVariantRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertLicenseVariant, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_LICENSE_VARIANT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseVariantService) doInsertLicenseVariant(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.LicenseVariantRequest)

	licenseVariantModel := repository.LicenseVariantModel{
		LicenseVariantName: sql.NullString{String: inputStruct.LicenseVariantName},
		CreatedBy:          sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:      sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:          sql.NullTime{Time: timeNow},
		UpdatedBy:          sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:      sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:          sql.NullTime{Time: timeNow},
	}

	var licenseVariantID int64
	licenseVariantID, err = dao.LicenseVariantDAO.InsertLicenseVariant(tx, licenseVariantModel)
	if err.Error != nil {
		err = checkDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.LicenseVariantDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: licenseVariantID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseVariantService) validateInsert(inputStruct *in.LicenseVariantRequest) errorModel.ErrorModel {
	return inputStruct.ValidationInsertLicenseVariant()
}
