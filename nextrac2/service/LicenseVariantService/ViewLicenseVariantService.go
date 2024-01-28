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
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input licenseVariantService) ViewLicenseVariant(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.LicenseVariantRequest

	inputStruct, err = input.readBodyAndValidateForView(request, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewLicenseVariant(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_LICENSE_VARIANT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseVariantService) doViewLicenseVariant(inputStruct in.LicenseVariantRequest, contextModel *applicationModel.ContextModel) (output out.LicenseVariantViewResponse, err errorModel.ErrorModel) {
	fileName := "ViewLicenseVariantService.go"
	funcName := "doViewLicenseVariant"
	var resultLicenseVariantOnDB repository.LicenseVariantModel

	licenseVariantModel := repository.LicenseVariantModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	}

	licenseVariantModel.CreatedBy.Int64 = 0

	resultLicenseVariantOnDB, err = dao.LicenseVariantDAO.ViewLicenseVariant(serverconfig.ServerAttribute.DBConnection, licenseVariantModel)
	if err.Error != nil {
		return
	}

	if resultLicenseVariantOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.LicenseVariant)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, resultLicenseVariantOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	output = input.convertLicenseVariantModelToDTOOut(resultLicenseVariantOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseVariantService) convertLicenseVariantModelToDTOOut(licenseVariantModel repository.LicenseVariantModel) out.LicenseVariantViewResponse {
	return out.LicenseVariantViewResponse{
		ID:                 licenseVariantModel.ID.Int64,
		LicenseVariantName: licenseVariantModel.LicenseVariantName.String,
		CreatedAt:          licenseVariantModel.CreatedAt.Time,
		UpdatedAt:          licenseVariantModel.UpdatedAt.Time,
		UpdatedName:        licenseVariantModel.UpdatedName.String,
	}
}

func (input licenseVariantService) validateView(inputStruct *in.LicenseVariantRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewLicenseVariant()
}
