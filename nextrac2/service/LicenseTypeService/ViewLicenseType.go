package LicenseTypeService

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
)

func (input licenseTypeService) ViewLicenseType(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.LicenseTypeRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewLicenseType)
	if err.Error != nil {
		return
	}

	output.Data.Content, err =input.doViewLicenseType(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status =input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input licenseTypeService) doViewLicenseType(inputStruct in.LicenseTypeRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	funcName := "doViewLicenseType"

	licenseTypeOnDB, err := dao.LicenseTypeDAO.ViewLicenseType(serverconfig.ServerAttribute.DBConnection, repository.LicenseTypeModel{
		ID:        sql.NullInt64{Int64: inputStruct.ID},
	})

	if licenseTypeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.LicenseType)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, licenseTypeOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.convertModelToResponseDetail(licenseTypeOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseTypeService) convertModelToResponseDetail(inputModel repository.LicenseTypeModel) out.LicenseTypeResponse {
	return out.LicenseTypeResponse{
		ID:              inputModel.ID.Int64,
		LicenseTypeName: inputModel.LicenseTypeName.String,
		LicenseTypeDesc: inputModel.LicenseTypeDesc.String,
		CreatedAt:       inputModel.CreatedAt.Time,
		UpdatedAt:       inputModel.UpdatedAt.Time,
		UpdatedBy:       inputModel.UpdatedBy.Int64,
		UpdatedName:     inputModel.UpdatedName.String,
	}
}

func (input licenseTypeService) validateViewLicenseType(inputStruct *in.LicenseTypeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}