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
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input licenseTypeService) DeleteLicenseType(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "DeleteLicenseType"
	var inputStruct in.LicenseTypeRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDeleteLicenseType)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteLicenseType, func(interface{}, applicationModel.ContextModel) {
		// additional function
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input licenseTypeService) doDeleteLicenseType(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doDeleteLicenseType"
	inputStruct := inputStructInterface.(in.LicenseTypeRequest)

	inputModel := repository.LicenseTypeModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

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

	// ----------- Update for delete
	encodedStr, errorS := service.RandToken(8)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	inputModel.LicenseTypeName.String = licenseTypeOnDB.LicenseTypeName.String + encodedStr

	err = dao.LicenseTypeDAO.DeleteLicenseType(tx, inputModel)
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.LicenseTypeDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	return
}

func (input licenseTypeService) validateDeleteLicenseType(inputStruct *in.LicenseTypeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}
