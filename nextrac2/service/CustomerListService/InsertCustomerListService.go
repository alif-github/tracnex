package CustomerListService

import (
	"database/sql"
	"errors"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input customerListService) InsertCustomer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "InsertCustomer"

	inputStruct, err := input.readBodyAndValidate(request, contextModel, input.validateInsertCustomer)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertCustomer, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code: 		util.GenerateConstantaI18n("OK", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: 	GenerateI18NMessage("SUCCESS_INSERT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListService) doInsertCustomer(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var customerModel []repository.CustomerListModel
	inputStruct := inputStructInterface.(in.CustomerListRequest)

	for _, branchDataElm := range inputStruct.BranchData {

		if branchDataElm.ImplementationAtStr != "" {
			branchDataElm.ImplementationAt, err = in.TimeStrToTime(branchDataElm.ImplementationAtStr, constanta.ImplementationAt)
			if err.Error != nil {
				return
			}
		}

		branchDataElm.ExpDateAt, err = in.TimeStrToTime(branchDataElm.ExpDateAtStr, constanta.ExpDate)
		if err.Error != nil {
			return
		}

		customerModel = append(customerModel, repository.CustomerListModel{
			CompanyID: 		sql.NullString{String: inputStruct.CompanyID},
			BranchID: 		sql.NullString{String: branchDataElm.BranchID},
			CompanyName: 	sql.NullString{String: branchDataElm.CompanyName},
			City: 			sql.NullString{String: branchDataElm.City},
			Implementer: 	sql.NullString{String: branchDataElm.Implementer},
			Implementation: sql.NullTime{Time: branchDataElm.ImplementationAt},
			Product: 		sql.NullString{String: branchDataElm.Product},
			Version: 		sql.NullString{String: branchDataElm.Version},
			LicenseType: 	sql.NullString{String: branchDataElm.LicenseType},
			UserAmount: 	sql.NullInt64{Int64: branchDataElm.UserOnLicense},
			ExpDate: 		sql.NullTime{Time: branchDataElm.ExpDateAt},
			CreatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			CreatedAt: 		sql.NullTime{Time: timeNow},
			CreatedClient: 	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			UpdatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedAt: 		sql.NullTime{Time: timeNow},
			UpdatedClient: 	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		})
	}

	id, err := dao.CustomerListDAO.InsertMultipleBranchCustomer(tx, customerModel)

	if err.Error != nil {
		return
	}

	for _, idCustomerElm := range id {
		dataAudit = append(dataAudit, repository.AuditSystemModel{
			TableName: 	sql.NullString{String: dao.CustomerListDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: idCustomerElm},
		})
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListService) DoInsertCustomerByImport(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	fileName := "InsertCustomerListService.go"
	funcName := "DoInsertCustomerByImport"
	inputStruct := inputStructInterface.(in.CustomerListImportRequest)

	//------- Repo model
	customerModel := repository.CustomerListModel{
		CompanyID: 		sql.NullString{String: inputStruct.CompanyID},
		BranchID: 		sql.NullString{String: inputStruct.BranchID},
		CompanyName: 	sql.NullString{String: inputStruct.CompanyName},
		City: 			sql.NullString{String: inputStruct.City},
		Implementer: 	sql.NullString{String: inputStruct.Implementer},
		Implementation: sql.NullTime{Time: inputStruct.Implementation},
		Product: 		sql.NullString{String: inputStruct.Product},
		Version: 		sql.NullString{String: inputStruct.Version},
		LicenseType: 	sql.NullString{String: inputStruct.LicenseType},
		UserAmount: 	sql.NullInt64{Int64: int64(inputStruct.UserAmount)},
		ExpDate: 		sql.NullTime{Time: inputStruct.ExpDate},
		CreatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedAt: 		sql.NullTime{Time: timeNow},
		CreatedClient: 	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedBy: 		sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt: 		sql.NullTime{Time: timeNow},
		UpdatedClient: 	sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	//------- Insert data customer list to DB
	var id int64
	id, err = dao.CustomerListDAO.InsertCustomerByImport(tx, customerModel)
	if err.Error != nil {
		return
	}

	if id < 1 {
		messageNewError := GenerateI18NMessage("FAILED_INSERT_CUSTOMER_MESSAGE", contextModel.AuthAccessTokenModel.Locale)
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errors.New(messageNewError))
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName: 	sql.NullString{String: dao.UserDAO.TableName},
		PrimaryKey:	sql.NullInt64{Int64: id},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListService) validateInsertCustomer(inputStruct *in.CustomerListRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertCustomer()
}
