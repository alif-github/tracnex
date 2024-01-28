package CustomerService

import (
	"database/sql"
	"fmt"
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
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service/CustomerCategoryService"
	"nexsoft.co.id/nextrac2/service/CustomerContactService"
	"nexsoft.co.id/nextrac2/service/CustomerGroupService"
	"nexsoft.co.id/nextrac2/service/MasterDataService/CompanyProfileService"
	"nexsoft.co.id/nextrac2/service/SalesmanService"
	"nexsoft.co.id/nextrac2/util"
	"strings"
	"time"
)

func (input customerService) InsertCustomer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName        = "InsertCustomer"
		inputStruct     in.CustomerRequest
		detailResponses out.CustomerErrorResponse
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsertCustomer)
	if err.Error != nil {
		detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
		detailResponses.PreviousRequest = input.getPreviousPayload(inputStruct)
		output.Other = detailResponses
		return
	}

	for i := 0; i < len(inputStruct.CustomerContact); i++ {
		err, detailResponses = CustomerContactService.CustomerContacService.ValidateInsertBulkCustomerContact(&inputStruct.CustomerContact[i], contextModel)
		if err.Error != nil {
			output.Other = detailResponses
			return
		}
	}

	additionalInfo, err := input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertCustomer, nil)
	if err.Error != nil {
		if additionalInfo != nil {
			output.Other = additionalInfo
		}
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input customerService) doInsertCustomer(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (additionalInfo interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		customerContactStruct []in.CustomerContactRequest
		tempAuditData         []repository.AuditSystemModel
		detailResponses       out.CustomerErrorResponse
		companyProfileReq     master_data_request.CompanyProfileWriteRequest
		mdbInsertedID         int64
		isNewCompanyProfile   bool
		inputStruct           = inputStructInterface.(in.CustomerRequest)
		internalToken 		  = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, constanta.DefaultApplicationsLanguage)
	)

	defer func() {
		if err.Error != nil {
			detailResponses.PreviousRequest = input.getPreviousPayload(inputStruct)
			additionalInfo = detailResponses
		}
	}()

	err = input.validateDataForInsert(&inputStruct, contextModel, false, internalToken)
	if err.Error != nil {
		detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
		return
	}

	//--- Company Profile section
	companyProfileReq = input.convertToMDBCompanyProfileRequest(inputStruct)
	if inputStruct.MDBCompanyProfileID == 0 {
		isNewCompanyProfile = true
		mdbInsertedID, err = master_data_dao.InsertCompanyProfile(companyProfileReq, contextModel)
		if err.Error != nil {
			detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
			return
		}

		inputStruct.MDBCompanyProfileID = mdbInsertedID
	} else {
		isNewCompanyProfile = false
		mdbInsertedID = inputStruct.MDBCompanyProfileID
		err = master_data_dao.UpdateCompanyProfile(companyProfileReq, contextModel)
		if err.Error != nil {
			detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
			return
		}
	}

	// Insert Customer
	inputModel := input.convertDTOToModelInsert(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	insertedCustomerID, err := dao.CustomerDAO.InsertCustomer(tx, inputModel)
	if err.Error != nil {
		err = input.checkDuplicateError(err)
		detailResponses.InsertCustomerResponse.InsertCustomerDetail = getErrorMessage(err, contextModel, inputStruct.Npwp)
		return
	}

	// Tambah flag untuk customer
	inputStruct.IsSuccess = true

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.CustomerDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: insertedCustomerID},
	})

	// Insert bulk Customer Contact
	if len(inputStruct.CustomerContact) > 0 {
		for i := 0; i < len(inputStruct.CustomerContact); i++ {
			inputStruct.CustomerContact[i].CustomerID = insertedCustomerID
			inputStruct.CustomerContact[i].MdbCompanyProfileID = mdbInsertedID
			customerContactStruct = append(customerContactStruct, inputStruct.CustomerContact[i])
		}

		detailResponses, tempAuditData, err = CustomerContactService.CustomerContacService.InsertBulkCustomerContactFromCustomer(tx, inputStruct.CustomerContact, contextModel, timeNow, false, isNewCompanyProfile)
		dataAudit = append(dataAudit, tempAuditData...)
		if err.Error != nil {
			return
		}

	}
	return
}

func (input customerService) convertDTOToModelInsert(inputStruct in.CustomerRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.CustomerModel {
	return repository.CustomerModel{
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
		CreatedBy:               sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedAt:               sql.NullTime{Time: timeNow},
		CreatedClient:           sql.NullString{String: authAccessToken.ClientID},
		UpdatedBy:               sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:               sql.NullTime{Time: timeNow},
		UpdatedClient:           sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input customerService) validateDataForInsert(inputStruct *in.CustomerRequest, contextModel *applicationModel.ContextModel, isForUpdate bool, internalToken string) (err errorModel.ErrorModel) {
	var (
		funcName     = "validateDataForInsert"
		customerOnDB repository.CustomerModel
	)

	//--- Get data scope
	scope, err := input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	//--- Validate MDB Data
	err = input.validateMDBForInsert(inputStruct, contextModel, isForUpdate, scope, internalToken)
	if err.Error != nil {
		return
	}

	//--- Validate Customer ID
	if inputStruct.ParentCustomerID > 0 {
		customerOnDB, err = dao.CustomerDAO.GetCustomerParentForValidate(serverconfig.ServerAttribute.DBConnection, repository.CustomerModel{
			ID:        sql.NullInt64{Int64: inputStruct.ParentCustomerID},
			CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
		}, scope, input.MappingScopeDB)
		if err.Error != nil {
			return
		}

		if customerOnDB.ID.Int64 < 1 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ParentCustomerID)
			return
		}

		inputStruct.MDBParentCustomerID = customerOnDB.MDBCompanyProfileID.Int64
	}

	//--- Validate RefCustomerID
	if !util.IsDataEmpty(inputStruct.RefCustomerID) {
		err = input.checkOptionalData(repository.CustomerModel{
			ID:        sql.NullInt64{Int64: inputStruct.RefCustomerID},
			CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
		}, scope)
		if err.Error != nil {
			return
		}
	}

	//--- Validate CustomerGroupID
	if !util.IsDataEmpty(inputStruct.CustomerGroupID) {
		if err = input.checkOptionalData(repository.CustomerGroupModel{
			ID:        sql.NullInt64{Int64: inputStruct.CustomerGroupID},
			CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
		}, scope); err.Error != nil {
			return
		}
	}

	//--- Validate CustomerCategoryID
	if !util.IsDataEmpty(inputStruct.CustomerCategoryID) {
		if err = input.checkOptionalData(repository.CustomerCategoryModel{
			ID:        sql.NullInt64{Int64: inputStruct.CustomerCategoryID},
			CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
		}, scope); err.Error != nil {
			return
		}
	}

	//--- Validate Salesman ID
	if !util.IsDataEmpty(inputStruct.SalesmanID) {
		if err = input.checkOptionalData(repository.SalesmanModel{
			ID:        sql.NullInt64{Int64: inputStruct.SalesmanID},
			CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
		}, scope); err.Error != nil {
			return
		}
	}

	return
}

func (input customerService) checkOptionalData(inputStruct interface{}, scope map[string]interface{}) (err errorModel.ErrorModel) {
	db := serverconfig.ServerAttribute.DBConnection
	funcName := "checkOptionalData"
	var isDataExist bool
	var fieldName string
	switch inputStruct.(type) {
	case repository.CustomerModel:
		isDataExist, err = dao.CustomerDAO.IsExistCustomerForInsert(db, inputStruct.(repository.CustomerModel))
		fieldName = constanta.RefCustomerID
		break
	case repository.CustomerGroupModel:
		isDataExist, err = dao.CustomerGroupDAO.IsExistCustomerGroupForInsert(db, inputStruct.(repository.CustomerGroupModel), scope, CustomerGroupService.CustomerGroupService.MappingScopeDB)
		fieldName = constanta.CustomerGroup
		break
	case repository.CustomerCategoryModel:
		isDataExist, err = dao.CustomerCategoryDAO.IsExistCustomerCategory(db, inputStruct.(repository.CustomerCategoryModel), scope, CustomerCategoryService.CustomerCategoryService.MappingScopeDB)
		fieldName = constanta.CustomerCategory
		break
	case repository.SalesmanModel:
		isDataExist, err = dao.SalesmanDAO.GetSalesmanForInsert(db, inputStruct.(repository.SalesmanModel), scope, SalesmanService.SalesmanService.MappingScopeDB)
		fieldName = constanta.SalesmanID
		break
	}
	if err.Error != nil {
		return
	}
	if !isDataExist {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, fieldName)
		return
	}

	return
}

func (input customerService) validateMDBForInsert(inputStruct *in.CustomerRequest, contextModel *applicationModel.ContextModel, isForUpdate bool, scope map[string]interface{}, internalToken string) (err errorModel.ErrorModel) {
	var (
		funcName               = "validateMDBForInsert"
		dataMDB                interface{}
		companyProfileResponse out.ViewCompanyProfileResponse
	)

	//--- Validate company profile on MDB
	if !isForUpdate {

		if inputStruct.MDBCompanyProfileID > 0 {
			//--- Company Profile
			var reqMasterData = master_data_request.CompanyProfileGetListRequest{
				ID: inputStruct.MDBCompanyProfileID,
			}

			//--- View Company Profile
			if dataMDB, err = CompanyProfileService.CompanyProfileService.DoViewCompanyProfile(
				reqMasterData, contextModel); err.Error != nil && err.Error.Error() == constanta.ErrorMDBDataNotFound {
				err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.CompanyProfile)
				return
			}

			companyProfileResponse, _ = dataMDB.(out.ViewCompanyProfileResponse)
			if companyProfileResponse.Npwp != inputStruct.Npwp {
				err = errorModel.GenerateDifferentRequestAndDBResult(input.FileName, funcName, constanta.NPWP, constanta.MDBNPWP)
				return
			}

			//--- Validate company title and regional on MDB
			err = input.validateMDBComponent(inputStruct, scope, contextModel, internalToken)
			if err.Error != nil {
				return
			}

			err = errorModel.GenerateNonErrorModel()
			return
		}

		if inputStruct.IsParent {
			var (
				dataListCompanyProfile []master_data_response.CompanyProfileResponse
				companyProfileList     master_data_request.CompanyProfileGetListRequest
			)

			companyProfileList.Page = 1
			companyProfileList.Limit = 100
			companyProfileList.OrderBy = constanta.ID
			companyProfileList.NPWP = inputStruct.Npwp

			if dataListCompanyProfile, err = master_data_dao.GetListCompanyProfileFromMasterData(companyProfileList, contextModel); err.Error != nil {
				err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.CompanyProfile)
				return
			}

			fmt.Println(fmt.Sprintf("[%s] -> Data Get List Company Profile -> %+v", funcName, dataListCompanyProfile))
			if len(dataListCompanyProfile) > 0 {
				for _, itemCompanyProfile := range dataListCompanyProfile {
					//--- View Company Profile
					if itemCompanyProfile.CompanyParent < 1 {
						fmt.Println(fmt.Sprintf("[%s] -> Company Parent : %d", funcName, itemCompanyProfile.CompanyParent))
						dataMDB, err = CompanyProfileService.CompanyProfileService.DoViewCompanyProfile(master_data_request.CompanyProfileGetListRequest{ID: itemCompanyProfile.ID}, contextModel)
						if err.Error != nil && err.Error.Error() == constanta.ErrorMDBDataNotFound {
							err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.CompanyProfile)
							return
						}

						companyProfileResponse, _ = dataMDB.(out.ViewCompanyProfileResponse)
						inputStruct.MDBCompanyProfileID = companyProfileResponse.ID
						return
					}
				}
			}
		}
	}

	fmt.Println(fmt.Sprintf("[%s] -> Sucess Validate company profile on MDB", funcName))

	//--- Validate company title and regional on MDB
	err = input.validateMDBComponent(inputStruct, scope, contextModel, internalToken)
	if err.Error != nil {
		return
	}

	fmt.Println(fmt.Sprintf("[%s] -> End", funcName))
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerService) validateInsertCustomer(inputStruct *in.CustomerRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsert()
}

func (input customerService) convertToMDBCompanyProfileRequest(inputStruct in.CustomerRequest) (output master_data_request.CompanyProfileWriteRequest) {
	var (
		arrPhone []string
	)

	output = master_data_request.CompanyProfileWriteRequest{
		ID:               inputStruct.MDBCompanyProfileID,
		CompanyTitleID:   inputStruct.MDBCompanyTitleID,
		NPWP:             inputStruct.Npwp,
		Name:             inputStruct.CustomerName,
		Address1:         inputStruct.Address,
		Address2:         inputStruct.Address2,
		Address3:         inputStruct.Address3,
		Hamlet:           inputStruct.Hamlet,
		Neighbourhood:    inputStruct.Neighbourhood,
		ProvinceID:       inputStruct.MDBProvinceID,
		CountryID:        inputStruct.CountryID,
		DistrictID:       inputStruct.MDBDistrictID,
		SubDistrictID:    inputStruct.MDBSubDistrictID,
		UrbanVillageID:   inputStruct.MDBUrbanVillageID,
		PostalCodeID:     inputStruct.MDBPostalCodeID,
		Latitude:         inputStruct.Latitude,
		Longitude:        inputStruct.Longitude,
		Email:            inputStruct.CompanyEmail,
		AlternativeEmail: inputStruct.AlternativeCompanyEmail,
		CompanyParent:    inputStruct.MDBParentCustomerID,
	}

	// Split phone
	arrPhone = strings.Split(inputStruct.Phone, "-")
	output.PhoneCountryCode = arrPhone[0]
	output.Phone = arrPhone[1]

	// Split Fax
	arrPhone = strings.Split(inputStruct.Fax, "-")
	if len(arrPhone) > 1 {
		output.PhoneFaxCountryCode = arrPhone[0]
		output.Fax = arrPhone[1]
	}

	return
}
