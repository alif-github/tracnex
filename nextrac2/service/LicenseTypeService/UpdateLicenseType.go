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
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input licenseTypeService) UpdateLicenseType(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel)  {
	funcName := "UpdateLicenseType"
	var inputStruct in.LicenseTypeRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateLiceseType)
	if err.Error != nil {
		return
	}

	output.Data.Content, err =input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateLicenseType, func( interface{}, applicationModel.ContextModel) {
		// additional function
	})
	if err.Error != nil {
		return
	}

	output.Status =input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)

	return
}

func (input licenseTypeService) doUpdateLicenseType(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doUpdateLicenseType"
	inputStruct := inputStructInterface.(in.LicenseTypeRequest)
	inputModel := input.convertDTOToModelUpdate(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	licenseTypeOnDB, err := dao.LicenseTypeDAO.GetLicenseTypeForUpdate(tx, repository.LicenseTypeModel{
		ID: inputModel.ID,
	})

	if err.Error != nil {
		return
	}
	if licenseTypeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.LicenseType)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, licenseTypeOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if licenseTypeOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.LicenseType)
		return
	}

	if licenseTypeOnDB.UpdatedAt.Time.Unix() != inputStruct.UpdatedAt.Unix() {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.LicenseType)
		return
	}

	err = dao.LicenseTypeDAO.UpdateLicenseType(tx, inputModel)
	if err.Error != nil {
		err = input.checkDuplicateError(err)
		return
	}
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.LicenseTypeDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	return
}

func (input licenseTypeService) convertDTOToModelUpdate(inputStruct in.LicenseTypeRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.LicenseTypeModel {
	return repository.LicenseTypeModel{
		ID:              sql.NullInt64{Int64: inputStruct.ID},
		LicenseTypeName: sql.NullString{String: inputStruct.LicenseTypeName},
		LicenseTypeDesc: sql.NullString{String: inputStruct.LicenseTypeDesc},
		UpdatedBy:       sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:       sql.NullTime{Time: timeNow},
		UpdatedClient:   sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input licenseTypeService) validateUpdateLiceseType(inputStruct *in.LicenseTypeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdate()
}
