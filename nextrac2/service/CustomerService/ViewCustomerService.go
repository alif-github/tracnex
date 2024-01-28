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
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input customerService) ViewCustomer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.CustomerRequest
	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewCustomer)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewCustomer(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input customerService) doViewCustomer(inputStruct in.CustomerRequest, contextModel *applicationModel.ContextModel) (output out.CustomerViewResponse, err errorModel.ErrorModel) {
	var (
		funcName            = "doViewCustomer"
		db                  = serverconfig.ServerAttribute.DBConnection
		scope               map[string]interface{}
		customerOnDB        repository.CustomerModel
		customerContactOnDB []repository.CustomerContactModel
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	customerOnDB, err = dao.CustomerDAO.ViewCustomer(db, repository.CustomerModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	}, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if customerOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Customer)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, customerOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	customerContactOnDB, err = dao.CustomerContactDAO.GetCustomerContactByCustomerID(db, repository.CustomerContactModel{
		CustomerID: sql.NullInt64{Int64: inputStruct.ID},
	})
	if err.Error != nil {
		return
	}

	customerOnDB.CustomerContact = customerContactOnDB
	output = input.reformatModelToDTOView(customerOnDB)
	return
}

func (input customerService) reformatModelToDTOView(inputModel repository.CustomerModel) out.CustomerViewResponse {
	var customerContact []out.CustomerContactViewResponse

	for _, item := range inputModel.CustomerContact {
		customerContact = append(customerContact, out.CustomerContactViewResponse{
			ID:                 item.ID.Int64,
			CustomerID:         item.CustomerID.Int64,
			MdbPersonProfileID: item.MdbPersonProfileID.Int64,
			Nik:                item.Nik.String,
			MdbPersonTitleID:   item.MdbPersonTitle.Int64,
			PersonTitle:        item.PersonTitle.String,
			FirstName:          item.FirstName.String,
			LastName:           item.LastName.String,
			Sex:                item.Sex.String,
			Address:            item.Address.String,
			Address2:           item.Address2.String,
			Address3:           item.Address3.String,
			Hamlet:             item.Hamlet.String,
			Neighbourhood:      item.Neighbourhood.String,
			ProvinceID:         item.ProvinceID.Int64,
			ProvinceName:       item.ProvinceName.String,
			DistrictID:         item.DistrictID.Int64,
			DistrictName:       item.DistrictName.String,
			Phone:              item.Phone.String,
			Email:              item.Email.String,
			MdbPositionID:      item.MdbPositionID.Int64,
			PositionName:       item.PositionName.String,
			Status:             item.Status.String,
			CreatedBy:          item.CreatedBy.Int64,
			CreatedAt:          item.CreatedAt.Time,
			CreatedName:        item.CreatedName.String,
			UpdatedBy:          item.UpdatedBy.Int64,
			UpdatedAt:          item.UpdatedAt.Time,
			UpdatedName:        item.UpdatedName.String,
		})
	}

	return out.CustomerViewResponse{
		ID:                      inputModel.ID.Int64,
		IsPrincipal:             inputModel.IsPrincipal.Bool,
		IsParent:                inputModel.IsParent.Bool,
		ParentCustomerID:        inputModel.ParentCustomerID.Int64,
		ParentCustomerName:      inputModel.ParentCustomerName.String,
		MDBParentCustomerID:     inputModel.MDBParentCustomerID.Int64,
		MDBCompanyProfileID:     inputModel.MDBCompanyProfileID.Int64,
		Npwp:                    inputModel.Npwp.String,
		MDBCompanyTitleID:       inputModel.MDBCompanyTitleID.Int64,
		CompanyTitle:            inputModel.CompanyTitle.String,
		CustomerName:            inputModel.CustomerName.String,
		Address:                 inputModel.Address.String,
		Hamlet:                  inputModel.Hamlet.String,
		Neighbourhood:           inputModel.Neighbourhood.String,
		CountryID:               inputModel.CountryID.Int64,
		ProvinceID:              inputModel.ProvinceID.Int64,
		ProvinceName:            inputModel.ProvinceName.String,
		DistrictID:              inputModel.DistrictID.Int64,
		DistrictName:            inputModel.DistrictName.String,
		SubDistrictID:           inputModel.SubDistrictID.Int64,
		SubDistrictName:         inputModel.SubDistrictName.String,
		UrbanVillageID:          inputModel.UrbanVillageID.Int64,
		UrbanVillageName:        inputModel.UrbanVillageName.String,
		PostalCodeID:            inputModel.PostalCodeID.Int64,
		PostalCode:              inputModel.PostalCode.String,
		Longitude:               inputModel.Longitude.Float64,
		Latitude:                inputModel.Latitude.Float64,
		Phone:                   inputModel.Phone.String,
		AlternativePhone:        inputModel.AlternativePhone.String,
		Fax:                     inputModel.Fax.String,
		CompanyEmail:            inputModel.CompanyEmail.String,
		AlternativeCompanyEmail: inputModel.AlternativeCompanyEmail.String,
		CustomerSource:          inputModel.CustomerSource.String,
		TaxName:                 inputModel.TaxName.String,
		TaxAddress:              inputModel.TaxAddress.String,
		SalesmanID:              inputModel.SalesmanID.Int64,
		SalesmanName:            inputModel.SalesmanFirstName.String + " " + inputModel.SalesmanLastName.String,
		RefCustomerID:           inputModel.RefCustomerID.Int64,
		RefCustomerName:         inputModel.RefCustomerName.String,
		DistributorOF:           inputModel.DistributorOF.String,
		CustomerGroupID:         inputModel.CustomerGroupID.Int64,
		CustomerGroupName:       inputModel.CustomerGroupName.String,
		CustomerCategoryID:      inputModel.CustomerCategoryID.Int64,
		CustomerCategoryName:    inputModel.CustomerCategoryName.String,
		Status:                  inputModel.Status.String,
		CustomerContact:         customerContact,
		CreatedBy:               inputModel.CreatedBy.Int64,
		CreatedAt:               inputModel.CreatedAt.Time,
		CreatedName:             inputModel.CreatedName.String,
		UpdatedBy:               inputModel.UpdatedBy.Int64,
		UpdatedAt:               inputModel.UpdatedAt.Time,
		UpdatedName:             inputModel.UpdatedName.String,
		IsUsed:                  inputModel.IsUsed.Bool,
		Address2:                inputModel.Address2.String,
		Address3:                inputModel.Address3.String,
	}
}

func (input customerService) validateViewCustomer(inputStruct *in.CustomerRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}
