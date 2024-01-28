package CustomerService

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/CustomerContactService"
	"time"
)

func (input customerService) DeleteCustomer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteCustomer"
		inputStruct in.CustomerRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDeleteCustomer)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteCustomer, func(interface{}, applicationModel.ContextModel) {
		//Additional Function
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input customerService) doDeleteCustomer(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName         = "DeleteCustomerService.go"
		funcName         = "doDeleteCustomer"
		inputStruct      = inputStructInterface.(in.CustomerRequest)
		scope            map[string]interface{}
		customerContacts []in.CustomerContactRequest
		detailResponses  out.CustomerErrorResponse
		customerOnDB     repository.CustomerModel
		tempAuditContact []repository.AuditSystemModel
		countChild       int
		dataChild        []out.DetailErrorCustomerResponse
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	inputModelCustomer := repository.CustomerModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	customerOnDB, err = dao.CustomerDAO.GetCustomerForDelete(tx, inputModelCustomer, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if customerOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Customer)
		return
	}

	if customerOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.Customer)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, customerOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if customerOnDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Customer)
		detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
		output = detailResponses
		return
	}

	//----------- Check Child Hierarchy
	countChild, dataChild, err = dao.CustomerDAO.GetCustomerHasChildHierarchy(serverconfig.ServerAttribute.DBConnection, repository.CustomerModel{ID: sql.NullInt64{Int64: customerOnDB.ID.Int64}})
	if err.Error != nil {
		return
	}

	if countChild > 0 {
		//----------- Child to detail
		var detail []string
		for _, dataChildItem := range dataChild {
			dataChildByte, errS := json.Marshal(dataChildItem)
			if errS == nil {
				detail = append(detail, string(dataChildByte))
			}
		}

		err = errorModel.GenerateErrorChildHierarchy(fileName, funcName, countChild, detail)
		return
	}

	if customerOnDB.IsParent.Bool {
		err = errorModel.GenerateDeleteParentCustomerError(input.FileName, funcName)
		return
	}

	inputModelCustomer.Npwp = customerOnDB.Npwp

	//----------- Update for delete
	encodedStr := service.RandTimeToken(constanta.RandTokenForDeleteLength, uuid.New().String())
	inputModelCustomer.Npwp.String += encodedStr

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.CustomerDAO.TableName, inputModelCustomer.ID.Int64, contextModel.LimitedByCreatedBy)...)
	err = dao.CustomerDAO.DeleteCustomer(tx, inputModelCustomer)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(customerOnDB.CustomerContactStr.String), &customerContacts)
	output, tempAuditContact, err = CustomerContactService.CustomerContacService.DoMultiDeleteCustomerContact(tx, customerContacts, contextModel, timeNow, false)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, tempAuditContact...)
	return
}

func (input customerService) validateDeleteCustomer(inputStruct *in.CustomerRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDelete()
}
