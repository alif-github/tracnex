package CustomerSiteService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input customerSiteService) InsertCustomerSite(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "InsertCustomerSite"

	var inputStruct in.CustomerSiteRequest
	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertCustomerSite, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code: util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_CUSTOMER_SITE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerSiteService) doInsertCustomerSite(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.CustomerSiteRequest)
	var customerSiteModel repository.CustomerSiteModel
	var customerID []int64

	customerSiteModel = repository.CustomerSiteModel{
		CreatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:        sql.NullTime{Time: timeNow},
		UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
	}

	customerID = append(customerID, inputStruct.ParentCustomerID, inputStruct.CustomerID)
	err = input.validateRelation(customerID, &customerSiteModel)
	if err.Error != nil {
		return
	}

	var customerSiteID int64
	//customerSiteID, err = dao.CustomerSiteDAO.InsertCustomerSite(tx, customerSiteModel)
	//if err.Error != nil {
	//	return
	//}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.CustomerSiteDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: customerSiteID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerSiteService) validateInsert(inputStruct *in.CustomerSiteRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertCustomerSite()
}
