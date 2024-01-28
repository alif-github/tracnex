package ProductLicenseService

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input productLicenseService) DoInsertProductLicenseBulk(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.([]in.ProductLicenseRequest)
	inputModel := input.convertDTOToModelInsertBulk(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
	var insertedID []int64

	// todo validate product license item
	for _, item := range inputStruct {
		err = input.validationDataInsertBulkProductLicense(item)
		if err.Error != nil {
			return
		}
	}

	insertedID, err = dao.ProductLicenseDAO.InsertBulkProductLicense(tx, inputModel)
	if err.Error != nil {
		return
	}

	for _, id := range insertedID {
		dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.ProductLicenseDAO.TableName, id, contextModel.LimitedByCreatedBy)...)
	}

	return
}

func (input productLicenseService) validationDataInsertBulkProductLicense(inputStruct in.ProductLicenseRequest) (err errorModel.ErrorModel) {
	funcName := "validationDataInsertBulkProductLicense"
	var isDataExist bool

	isDataExist, err = dao.LicenseConfigDAO.IsLicenseConfigExist(serverconfig.ServerAttribute.DBConnection, repository.LicenseConfigModel{
		ID: sql.NullInt64{Int64: inputStruct.LicenseConfigID},
	})
	if err.Error != nil {
		return
	}
	if !isDataExist {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.LicenseConfig)
		return
	}

	isDataExist, err = dao.ClientCredentialDAO.IsClientCredentialExist(serverconfig.ServerAttribute.DBConnection, repository.ClientCredentialModel{
		ClientID:     sql.NullString{String: inputStruct.ClientID},
		ClientSecret: sql.NullString{String: inputStruct.ClientSecret},
	})
	if err.Error != nil {
		return
	}
	if !isDataExist {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) convertDTOToModelInsertBulk(inputStruct []in.ProductLicenseRequest, authAccessModel model.AuthAccessTokenModel, timeNow time.Time) (result []repository.ProductLicenseModel) {
	for _, item := range inputStruct {
		result = append(result, repository.ProductLicenseModel{
			LicenseConfigId: sql.NullInt64{Int64: item.LicenseConfigID},
			ProductKey:      sql.NullString{String: item.ProductKey},
			ProductEncrypt:  sql.NullString{String: item.ProductEncrypt},
			ClientId:        sql.NullString{String: item.ClientID},
			ClientSecret:    sql.NullString{String: item.ClientSecret},
			HWID:            sql.NullString{String: item.Hwid},
			ActivationDate:  sql.NullTime{Time: item.ActivationDate},
			LicenseStatus:   sql.NullInt32{Int32: int32(item.LicenseStatus)},
			CreatedBy:       sql.NullInt64{Int64: authAccessModel.ResourceUserID},
			CreatedAt:       sql.NullTime{Time: timeNow},
			CreatedClient:   sql.NullString{String: authAccessModel.ClientID},
			UpdatedBy:       sql.NullInt64{Int64: authAccessModel.ResourceUserID},
			UpdatedAt:       sql.NullTime{Time: timeNow},
			UpdatedClient:   sql.NullString{String: authAccessModel.ClientID},
		})
	}
	return result
}
