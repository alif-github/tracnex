package CustomerContactService

import (
	"database/sql"
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

func (input customerContactService) DoMultiDeleteCustomerContact(tx *sql.Tx, inputStruct []in.CustomerContactRequest, contextModel *applicationModel.ContextModel, timeNow time.Time, isUpdate bool) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var tempAudit []repository.AuditSystemModel
	var detailResponses out.CustomerErrorResponse

	for i := 0; i < len(inputStruct); i++ {

		output, tempAudit, err = input.doDeleteCustomerContact(tx, inputStruct[i], contextModel, timeNow, isUpdate)
		if err.Error != nil {
			detailResponses.InsertCustomerResponse.InsertCustomerContactDetail = getErrorMessage(err, contextModel, inputStruct[i].Nik)
			output = detailResponses
			return
		}

		dataAudit = append(dataAudit, tempAudit...)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerContactService) doDeleteCustomerContact(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time, isUpdate bool) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doDeleteCustomerContact"
	inputStruct := inputStructInterface.(in.CustomerContactRequest)

	if inputStruct.ID < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.CustomerContact)
		return
	}

	inputModel := repository.CustomerContactModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	custContactOnDB, err := dao.CustomerContactDAO.GetCustomerContactForUpdate(tx, repository.CustomerContactModel{
		ID:         sql.NullInt64{Int64: inputStruct.ID},
		CustomerID: sql.NullInt64{Int64: inputStruct.CustomerID},
	})
	if err.Error != nil {
		return
	}

	if custContactOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.CustomerContact)
		return
	}

	inputModel.Nik = custContactOnDB.Nik

	// ----------- Update for delete
	encodedStr, errorS := service.RandToken(constanta.RandTokenForDeleteLength)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	if isUpdate {
		inputModel.Nik.String += encodedStr
	}

	err = dao.CustomerContactDAO.DeleteCustomerContact(tx, inputModel)
	if err.Error != nil {
		err.AdditionalInformation = append(err.AdditionalInformation, " DELETE NIK : "+inputStruct.Nik)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.CustomerContactDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerContactService) validateDeleteCustomerContact(inputStruct *in.CustomerContactRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}
