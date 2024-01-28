package UserRegistrationAdminService

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
	"time"
)

func (input userRegistrationAdminService) InsertUserRegistrationAdmin(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "InsertUserRegistrationAdmin"
	var inputStruct in.UserRegistrationAdminRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateInsertUserRegistrationAdmin)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertUserRegistrationAdmin, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input userRegistrationAdminService) doInsertUserRegistrationAdmin(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	fileName := "InsertUserRegistrationAdminService.go"
	funcName := "doInsertUserRegistrationAdmin"
	var insertedID int64
	var clientMappingData repository.ClientMappingModel
	inputStruct := inputStructInterface.(in.UserRegistrationAdminRequest)
	inputModel := input.convertDTOToModel(inputStruct, contextModel, timeNow)

	clientMappingData, err = dao.ClientMappingDAO.CheckClientMappingWithClientID(serverconfig.ServerAttribute.DBConnection, repository.ClientMappingModel{
		ClientID:  sql.NullString{String: inputModel.ClientID.String},
		CompanyID: sql.NullString{String: inputModel.UniqueID1.String},
		BranchID:  sql.NullString{String: inputModel.UniqueID2.String},
	})
	if err.Error != nil {
		return
	}

	if clientMappingData.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientMappingClientID)
		return
	}

	if clientMappingData.ClientTypeID.Int64 != inputModel.ClientTypeID.Int64 {
		err = errorModel.GenerateDifferentRequestAndDBResult(fileName, funcName, "Request", "Database")
		return
	}

	inputModel.CustomerId.Int64 = clientMappingData.CustomerID.Int64
	inputModel.ParentCustomerId.Int64 = clientMappingData.ParentCustomerID.Int64
	inputModel.SiteId.Int64 = clientMappingData.SiteID.Int64
	inputModel.ClientMappingID.Int64 = clientMappingData.ID.Int64

	insertedID, err = dao.UserRegistrationAdminDAO.InsertUserRegistrationAdmin(tx, inputModel)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.UserRegistrationAdminDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: insertedID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userRegistrationAdminService) ValidateInsertUserRegistrationAdmin(inputStruct *in.UserRegistrationAdminRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertUserRegistrationAdmin()
}
