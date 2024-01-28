package ClientCredentialService

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
	"time"
)

func (input clientCredentialService) InsertClientCredential(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "InsertClientCredential"
	var inputStruct in.ClientCredential

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateInsertClientCredential)
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

func (input clientCredentialService) doInsertLicenseType(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.ClientCredential)
	inputModel := input.convertDTOToModel(inputStruct, timeNow)

	insertedID, err := dao.ClientCredentialDAO.InsertClientCredential(tx, &inputModel)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.ClientCredentialDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: insertedID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientCredentialService) convertDTOToModel(inputStruct in.ClientCredential, timeNow time.Time) repository.ClientCredentialModel {
	return repository.ClientCredentialModel{
		ClientID:      sql.NullString{String: inputStruct.ClientID},
		ClientSecret:  sql.NullString{String: inputStruct.ClientSecret},
		SignatureKey:  sql.NullString{String: inputStruct.SignatureKey},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
	}
}

func (input clientCredentialService) ValidateInsertClientCredential(inputStruct in.ClientCredential) errorModel.ErrorModel {
	return inputStruct.MandatoryValidationClientCredential()
}
