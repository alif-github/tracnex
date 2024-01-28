package CustomerService

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
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/CustomerContactService"
	"time"
)

func (input customerService) UpdateCustomer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "UpdateCustomer"
	var (
		inputStruct in.CustomerRequest
		isUpdateAll = service.IsHaveAllPermission(contextModel)
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, func(input *in.CustomerRequest) errorModel.ErrorModel {
		return input.ValidateUpdate(isUpdateAll)
	})
	if err.Error != nil {
		return
	}

	additionalInfo, err := input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateCustomer, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		output.Other = additionalInfo
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input customerService) doUpdateCustomer(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName              = "doUpdateCustomer"
		scope                 map[string]interface{}
		detailResponses       out.CustomerErrorResponse
		companyProfileRequest master_data_request.CompanyProfileWriteRequest
		inputStruct           = inputStructInterface.(in.CustomerRequest)
		isUpdateAll 		  = service.IsHaveAllPermission(contextModel)
		internalToken 		  = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, constanta.DefaultApplicationsLanguage)
	)

	defer func() {
		// Get detail error
		if err.Error != nil {
			detailResponses.PreviousRequest = input.getPreviousPayload(inputStruct)
			output = detailResponses
		}
	}()

	// Get Scope
	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	customerOnDB, err := dao.CustomerDAO.GetCustomerForUpdate(tx, repository.CustomerModel{
		ID:        sql.NullInt64{Int64: inputStruct.ID},
		CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	}, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if customerOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Customer)
		detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
		return
	}

	err = input.validateBlockingData(inputStruct, customerOnDB, isUpdateAll)
	if err.Error != nil {
		detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
		return
	}

	//if customerOnDB.IsUsed.Bool {
	//	// validate MDB Company Profile ID if there is changed
	//	err = input.validateBlockingData(inputStruct, customerOnDB, isUpdateAll)
	//	if err.Error != nil {
	//		//err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.Customer)
	//		detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
	//		return
	//	}
	//
	//	// validate customer contact length
	//	//if len(inputStruct.CustomerContact) < 1 {
	//	//	err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.Customer)
	//	//	detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
	//	//	output = detailResponses
	//	//	return
	//	//}
	//}

	if customerOnDB.UpdatedAt.Time.Sub(inputStruct.UpdatedAt) != time.Duration(0) {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.Customer)
		detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
		return
	}

	err = input.validateDataForInsert(&inputStruct, contextModel, true, internalToken)
	if err.Error != nil {
		detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
		return
	}

	/*
		Set Company Profile Into MDB
	*/
	companyProfileRequest = input.convertToMDBCompanyProfileRequest(inputStruct)
	companyProfileRequest.ID = customerOnDB.MDBCompanyProfileID.Int64

	companyProfileId, err := input.setCompanyProfileIntoMDB(companyProfileRequest, contextModel)
	if err.Error != nil {
		detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
		return
	}

	inputStruct.MDBCompanyProfileID = companyProfileId

	/*
		Update Customer
	*/
	inputModel := input.convertDTOToModelUpdate(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	err = dao.CustomerDAO.UpdateCustomer(tx, inputModel)
	if err.Error != nil {
		err = input.checkDuplicateError(err)
		detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.CustomerDAO.TableName, customerOnDB.ID.Int64, contextModel.LimitedByCreatedBy)...)
	inputStruct.IsSuccess = true

	// Update Bulk Customer Contact
	tempOutput, auditCustContact, err := CustomerContactService.CustomerContacService.UpdateCustomerContactForCustomer(tx, inputStruct.CustomerContact, contextModel, timeNow, inputStruct.ID, inputStruct.MDBCompanyProfileID)
	if err.Error != nil {
		if tempOutput != nil {
			detailResponses = tempOutput.(out.CustomerErrorResponse)
		}
		return
	}

	dataAudit = append(dataAudit, auditCustContact...)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerService) validateBlockingData(inputStruct in.CustomerRequest, customerOnDB repository.CustomerModel, isUpdateAll bool) errorModel.ErrorModel {
	funcName := "validateBlockingData"

	isCustomerUsed := customerOnDB.IsUsed.Bool

	if isUpdateAll {
		return errorModel.GenerateNonErrorModel()
	}

	if inputStruct.Npwp != customerOnDB.Npwp.String {
		if !isCustomerUsed {
			return errorModel.GenerateCannotChangedError(input.FileName, funcName, constanta.NPWP)
		}

		return errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.Customer)
	}

	if inputStruct.ParentCustomerID != customerOnDB.ParentCustomerID.Int64 {
		if !isCustomerUsed {
			return errorModel.GenerateCannotChangedError(input.FileName, funcName, constanta.ParentCustomerID)
		}

		return errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.Customer)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input customerService) setCompanyProfileIntoMDB(companyProfileRequest master_data_request.CompanyProfileWriteRequest, contextModel *applicationModel.ContextModel) (companyProfileId int64, err errorModel.ErrorModel) {
	companyProfileId = companyProfileRequest.ID
	if companyProfileId == 0 {
		companyProfileId, err = master_data_dao.InsertCompanyProfile(companyProfileRequest, contextModel)
		return
	}
	err = master_data_dao.UpdateCompanyProfile(companyProfileRequest, contextModel)
	return
}

func (input customerService) convertDTOToModelUpdate(inputStruct in.CustomerRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.CustomerModel {
	return repository.CustomerModel{
		ID:                      sql.NullInt64{Int64: inputStruct.ID},
		IsPrincipal:             sql.NullBool{Bool: inputStruct.IsPrincipal},
		IsParent:                sql.NullBool{Bool: inputStruct.IsParent},
		ParentCustomerID:        sql.NullInt64{Int64: inputStruct.ParentCustomerID},
		MDBParentCustomerID:     sql.NullInt64{Int64: inputStruct.MDBParentCustomerID},
		MDBCompanyProfileID:     sql.NullInt64{Int64: inputStruct.MDBCompanyProfileID},
		Npwp:                    sql.NullString{String: inputStruct.Npwp},
		MDBCompanyTitleID:       sql.NullInt64{Int64: inputStruct.MDBCompanyTitleID},
		CompanyTitle:            sql.NullString{String: inputStruct.CompanyTitle},
		CustomerName:            sql.NullString{String: inputStruct.CustomerName},
		Address:                 sql.NullString{String: inputStruct.Address},
		Address2:                sql.NullString{String: inputStruct.Address2},
		Address3:                sql.NullString{String: inputStruct.Address3},
		Hamlet:                  sql.NullString{String: inputStruct.Hamlet},
		Neighbourhood:           sql.NullString{String: inputStruct.Neighbourhood},
		CountryID:               sql.NullInt64{Int64: inputStruct.CountryID},
		ProvinceID:              sql.NullInt64{Int64: inputStruct.ProvinceID},
		DistrictID:              sql.NullInt64{Int64: inputStruct.DistrictID},
		SubDistrictID:           sql.NullInt64{Int64: inputStruct.SubDistrictID},
		UrbanVillageID:          sql.NullInt64{Int64: inputStruct.UrbanVillageID},
		PostalCodeID:            sql.NullInt64{Int64: inputStruct.PostalCodeID},
		Longitude:               sql.NullFloat64{Float64: inputStruct.Longitude},
		Latitude:                sql.NullFloat64{Float64: inputStruct.Latitude},
		Phone:                   sql.NullString{String: inputStruct.Phone},
		AlternativePhone:        sql.NullString{String: inputStruct.AlternativePhone},
		Fax:                     sql.NullString{String: inputStruct.Fax},
		CompanyEmail:            sql.NullString{String: inputStruct.CompanyEmail},
		AlternativeCompanyEmail: sql.NullString{String: inputStruct.AlternativeCompanyEmail},
		CustomerSource:          sql.NullString{String: inputStruct.CustomerSource},
		TaxName:                 sql.NullString{String: inputStruct.TaxName},
		TaxAddress:              sql.NullString{String: inputStruct.TaxAddress},
		SalesmanID:              sql.NullInt64{Int64: inputStruct.SalesmanID},
		RefCustomerID:           sql.NullInt64{Int64: inputStruct.RefCustomerID},
		DistributorOF:           sql.NullString{String: inputStruct.DistributorOF},
		CustomerGroupID:         sql.NullInt64{Int64: inputStruct.CustomerGroupID},
		CustomerCategoryID:      sql.NullInt64{Int64: inputStruct.CustomerCategoryID},
		Status:                  sql.NullString{String: inputStruct.Status},
		UpdatedBy:               sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:               sql.NullTime{Time: timeNow},
		UpdatedClient:           sql.NullString{String: authAccessToken.ClientID},
	}
}

//func (input customerService) validateUpdateCustomer(inputStruct *in.CustomerRequest) errorModel.ErrorModel {
//	return inputStruct.ValidateUpdate()
//}
