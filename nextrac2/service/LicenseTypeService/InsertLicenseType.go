package LicenseTypeService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"time"
)

func (input licenseTypeService) InsertLicenseType(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "InsertLicenseType"
	var inputStruct in.LicenseTypeRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateInsertLicenseType)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertLicenseType, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input licenseTypeService) doInsertLicenseType(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.LicenseTypeRequest)
	inputModel := input.convertDTOToModel(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	insertedID, err := dao.LicenseTypeDAO.InsertLicenseType(tx, inputModel)
	if err.Error != nil {
		err = input.checkDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName: sql.NullString{String: dao.LicenseTypeDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: insertedID},
	})

	err = errorModel.GenerateNonErrorModel()
	return 
}

func (input licenseTypeService) convertDTOToModel(inputStruct in.LicenseTypeRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.LicenseTypeModel {
	return repository.LicenseTypeModel{
		LicenseTypeName: sql.NullString{String: inputStruct.LicenseTypeName},
		LicenseTypeDesc: sql.NullString{String: inputStruct.LicenseTypeDesc},
		CreatedBy:       sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedAt:       sql.NullTime{Time: timeNow},
		CreatedClient:   sql.NullString{String: authAccessToken.ClientID},
		UpdatedBy:       sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:       sql.NullTime{Time: timeNow},
		UpdatedClient:   sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input licenseTypeService) ValidateInsertLicenseType(inputStruct *in.LicenseTypeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsert()
}