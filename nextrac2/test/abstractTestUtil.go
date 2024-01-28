package test

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service/ProductService"
	"time"
)

type AbstractScenario struct {
	Name        string
	RequestBody interface{}
	Expected    interface{}
}

type AbstractScenarioWithContextModel struct {
	Name         string
	RequestBody  interface{}
	Expected     interface{}
	ContextModel applicationModel.ContextModel
}

type MasterDataTesting struct {
	CustomerID              int64
	ParentCustomerID        int64
	ProductND6ID            int64
	ProductNexmileID        int64
	ProductNexchiefID       int64
	ProductNexchiefMobileID int64
	LicenseVariantID        int64
	LicenseTypeID           int64
	ProductGroupID          int64
	ModuleID                int64
	ComponentID             int64
	CustomerCategoryID      int64
	CustomerGroupID         int64
	SalesmanID              int64
}

func (input *MasterDataTesting) CreateInitiateData(tx *sql.Tx, timeNow time.Time) (err errorModel.ErrorModel) {
	var productModel repository.ProductModel

	// ---------------- Insert License Variant
	if input.LicenseVariantID, err = dao.LicenseVariantDAO.InsertLicenseVariant(tx, repository.LicenseVariantModel{
		LicenseVariantName: sql.NullString{String: "variant test"},
		UpdatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:      sql.NullString{String: constanta.SystemClient},
		UpdatedAt:          sql.NullTime{Time: timeNow},
		CreatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:      sql.NullString{String: constanta.SystemClient},
		CreatedAt:          sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert License Type
	if input.LicenseTypeID, err = dao.LicenseTypeDAO.InsertLicenseType(tx, repository.LicenseTypeModel{
		LicenseTypeName: sql.NullString{String: "type test"},
		UpdatedBy:       sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:   sql.NullString{String: constanta.SystemClient},
		UpdatedAt:       sql.NullTime{Time: timeNow},
		CreatedBy:       sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:   sql.NullString{String: constanta.SystemClient},
		CreatedAt:       sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Customer Category
	if input.CustomerCategoryID, err = dao.CustomerCategoryDAO.InsertCustomerCategory(tx, repository.CustomerCategoryModel{
		CustomerCategoryID:   sql.NullString{String: "VALID"},
		CustomerCategoryName: sql.NullString{String: "VALID"},
		UpdatedBy:            sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:        sql.NullString{String: constanta.SystemClient},
		UpdatedAt:            sql.NullTime{Time: timeNow},
		CreatedBy:            sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:        sql.NullString{String: constanta.SystemClient},
		CreatedAt:            sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Customer Group
	if input.CustomerGroupID, err = dao.CustomerGroupDAO.InsertCustomerGroup(tx, repository.CustomerGroupModel{
		CustomerGroupID:   sql.NullString{String: "VALID"},
		CustomerGroupName: sql.NullString{String: "VALID"},
		UpdatedBy:         sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:     sql.NullString{String: constanta.SystemClient},
		UpdatedAt:         sql.NullTime{Time: timeNow},
		CreatedBy:         sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:     sql.NullString{String: constanta.SystemClient},
		CreatedAt:         sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Salesman
	if input.SalesmanID, err = dao.SalesmanDAO.InsertSalesman(tx, repository.SalesmanModel{
		PersonTitleID: sql.NullInt64{Int64: 1},
		PersonTitle:   sql.NullString{String: "VALID"},
		Nik:           sql.NullString{String: "VALID"},
		FirstName:     sql.NullString{String: "VALID"},
		LastName:      sql.NullString{String: "VALID"},
		Sex:           sql.NullString{String: "L"},
		Address:       sql.NullString{String: "VALID"},
		Phone:         sql.NullString{String: "VALID"},
		Email:         sql.NullString{String: "VALID"},
		Status:        sql.NullString{String: "A"},
		Hamlet:        sql.NullString{String: "VALID"},
		Neighbourhood: sql.NullString{String: "VALID"},
		ProvinceID:    sql.NullInt64{Int64: 1},
		DistrictID:    sql.NullInt64{Int64: 1},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Parent Customer
	if input.ParentCustomerID, err = dao.CustomerDAO.InsertCustomer(tx, repository.CustomerModel{
		IsPrincipal:             sql.NullBool{Bool: true},
		IsParent:                sql.NullBool{Bool: true},
		MDBCompanyProfileID:     sql.NullInt64{Int64: 1},
		MDBCompanyTitleID:       sql.NullInt64{Int64: 1},
		Npwp:                    sql.NullString{String: "VALID"},
		CompanyTitle:            sql.NullString{String: "VALID"},
		CustomerName:            sql.NullString{String: "Parent Customer"},
		Address:                 sql.NullString{String: "VALID"},
		Hamlet:                  sql.NullString{String: "VALID"},
		Neighbourhood:           sql.NullString{String: "VALID"},
		CountryID:               sql.NullInt64{Int64: 1},
		ProvinceID:              sql.NullInt64{Int64: 1},
		DistrictID:              sql.NullInt64{Int64: 1},
		SubDistrictID:           sql.NullInt64{Int64: 1},
		UrbanVillageID:          sql.NullInt64{Int64: 1},
		PostalCodeID:            sql.NullInt64{Int64: 1},
		Latitude:                sql.NullFloat64{Float64: 1},
		Longitude:               sql.NullFloat64{Float64: 1},
		Phone:                   sql.NullString{String: "VALID"},
		AlternativePhone:        sql.NullString{String: "VALID"},
		Fax:                     sql.NullString{String: "VALID"},
		CompanyEmail:            sql.NullString{String: "VALID"},
		AlternativeCompanyEmail: sql.NullString{String: "VALID"},
		CustomerSource:          sql.NullString{String: "VALID"},
		TaxName:                 sql.NullString{String: "VALID"},
		TaxAddress:              sql.NullString{String: "VALID"},
		SalesmanID:              sql.NullInt64{Int64: input.SalesmanID},
		DistributorOF:           sql.NullString{String: "VALID"},
		CustomerGroupID:         sql.NullInt64{Int64: input.CustomerGroupID},
		CustomerCategoryID:      sql.NullInt64{Int64: input.CustomerCategoryID},
		Status:                  sql.NullString{String: "A"},
		UpdatedBy:               sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:           sql.NullString{String: constanta.SystemClient},
		UpdatedAt:               sql.NullTime{Time: timeNow},
		CreatedBy:               sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:           sql.NullString{String: constanta.SystemClient},
		CreatedAt:               sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Customer
	if input.CustomerID, err = dao.CustomerDAO.InsertCustomer(tx, repository.CustomerModel{
		IsPrincipal:             sql.NullBool{Bool: true},
		IsParent:                sql.NullBool{Bool: false},
		ParentCustomerID:        sql.NullInt64{Int64: input.ParentCustomerID},
		MDBCompanyProfileID:     sql.NullInt64{Int64: 1},
		MDBCompanyTitleID:       sql.NullInt64{Int64: 1},
		Npwp:                    sql.NullString{String: "VALID"},
		CompanyTitle:            sql.NullString{String: "VALID"},
		CustomerName:            sql.NullString{String: "Customer"},
		Address:                 sql.NullString{String: "VALID"},
		Hamlet:                  sql.NullString{String: "VALID"},
		Neighbourhood:           sql.NullString{String: "VALID"},
		CountryID:               sql.NullInt64{Int64: 1},
		ProvinceID:              sql.NullInt64{Int64: 1},
		DistrictID:              sql.NullInt64{Int64: 1},
		SubDistrictID:           sql.NullInt64{Int64: 1},
		UrbanVillageID:          sql.NullInt64{Int64: 1},
		PostalCodeID:            sql.NullInt64{Int64: 1},
		Latitude:                sql.NullFloat64{Float64: 1},
		Longitude:               sql.NullFloat64{Float64: 1},
		Phone:                   sql.NullString{String: "VALID"},
		AlternativePhone:        sql.NullString{String: "VALID"},
		Fax:                     sql.NullString{String: "VALID"},
		CompanyEmail:            sql.NullString{String: "VALID"},
		AlternativeCompanyEmail: sql.NullString{String: "VALID"},
		CustomerSource:          sql.NullString{String: "VALID"},
		TaxName:                 sql.NullString{String: "VALID"},
		TaxAddress:              sql.NullString{String: "VALID"},
		SalesmanID:              sql.NullInt64{Int64: input.SalesmanID},
		DistributorOF:           sql.NullString{String: "VALID"},
		CustomerGroupID:         sql.NullInt64{Int64: input.CustomerGroupID},
		CustomerCategoryID:      sql.NullInt64{Int64: input.CustomerCategoryID},
		Status:                  sql.NullString{String: "A"},
		UpdatedBy:               sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:           sql.NullString{String: constanta.SystemClient},
		UpdatedAt:               sql.NullTime{Time: timeNow},
		CreatedBy:               sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:           sql.NullString{String: constanta.SystemClient},
		CreatedAt:               sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Product Group
	if input.ProductGroupID, err = dao.ProductGroupDAO.InsertProductGroup(tx, repository.ProductGroupModel{
		ProductGroupName: sql.NullString{String: "VALID"},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Module
	if input.ModuleID, err = dao.ModuleDAO.InsertModule(tx, repository.ModuleModel{
		ModuleName:    sql.NullString{String: "VALID"},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Component
	if input.ComponentID, err = dao.ComponentDAO.InsertComponent(tx, repository.ComponentModel{
		ComponentName: sql.NullString{String: "Valid"},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Product ND6
	productModel = repository.ProductModel{
		ProductID:          sql.NullString{String: "PROD_ND6"},
		ProductName:        sql.NullString{String: "PROD_ND6"},
		ProductDescription: sql.NullString{String: "VALID1"},
		ProductGroupID:     sql.NullInt64{Int64: input.ProductGroupID},
		ClientTypeID:       sql.NullInt64{Int64: constanta.ResourceND6ID},
		IsLicense:          sql.NullBool{Bool: true},
		LicenseVariantID:   sql.NullInt64{Int64: input.LicenseVariantID},
		LicenseTypeID:      sql.NullInt64{Int64: input.LicenseTypeID},
		DeploymentMethod:   sql.NullString{String: "O"},
		NoOfUser:           sql.NullInt64{Int64: 1},
		IsUserConcurrent:   sql.NullBool{Bool: true},
		MaxOfflineDays:     sql.NullInt64{Int64: 1},
		Module1:            sql.NullInt64{Int64: input.ModuleID},
		UpdatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:      sql.NullString{String: constanta.SystemClient},
		UpdatedAt:          sql.NullTime{Time: timeNow},
		CreatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:      sql.NullString{String: constanta.SystemClient},
		CreatedAt:          sql.NullTime{Time: timeNow},
		ProductComponentModel: []repository.ProductComponentModel{
			{
				ComponentID:    sql.NullInt64{Int64: input.ComponentID},
				ComponentValue: sql.NullString{String: "Valid"},
			},
		},
	}
	if err, input.ProductND6ID = input.insertProduct(tx, productModel); err.Error != nil {
		return
	}

	// ---------------- Insert Product Nexchief
	productModel = repository.ProductModel{
		ProductID:          sql.NullString{String: "PROD_NEXCHIEF"},
		ProductName:        sql.NullString{String: "PROD_NEXCHIEF"},
		ProductDescription: sql.NullString{String: "VALID1"},
		ProductGroupID:     sql.NullInt64{Int64: input.ProductGroupID},
		ClientTypeID:       sql.NullInt64{Int64: constanta.ResourceNexChiefID},
		IsLicense:          sql.NullBool{Bool: true},
		LicenseVariantID:   sql.NullInt64{Int64: input.LicenseVariantID},
		LicenseTypeID:      sql.NullInt64{Int64: input.LicenseTypeID},
		DeploymentMethod:   sql.NullString{String: "O"},
		NoOfUser:           sql.NullInt64{Int64: 1},
		IsUserConcurrent:   sql.NullBool{Bool: true},
		MaxOfflineDays:     sql.NullInt64{Int64: 1},
		Module1:            sql.NullInt64{Int64: input.ModuleID},
		UpdatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:      sql.NullString{String: constanta.SystemClient},
		UpdatedAt:          sql.NullTime{Time: timeNow},
		CreatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:      sql.NullString{String: constanta.SystemClient},
		CreatedAt:          sql.NullTime{Time: timeNow},
		ProductComponentModel: []repository.ProductComponentModel{
			{
				ComponentID:    sql.NullInt64{Int64: input.ComponentID},
				ComponentValue: sql.NullString{String: "Valid"},
			},
		},
	}
	if err, input.ProductNexchiefID = input.insertProduct(tx, productModel); err.Error != nil {
		return
	}

	// ---------------- Insert Product NexMile
	productModel = repository.ProductModel{
		ProductID:          sql.NullString{String: "PROD_NEXMILE"},
		ProductName:        sql.NullString{String: "PROD_NEXMILE"},
		ProductDescription: sql.NullString{String: "VALID1"},
		ProductGroupID:     sql.NullInt64{Int64: input.ProductGroupID},
		ClientTypeID:       sql.NullInt64{Int64: constanta.ResourceNexmileID},
		IsLicense:          sql.NullBool{Bool: true},
		LicenseVariantID:   sql.NullInt64{Int64: input.LicenseVariantID},
		LicenseTypeID:      sql.NullInt64{Int64: input.LicenseTypeID},
		DeploymentMethod:   sql.NullString{String: "O"},
		NoOfUser:           sql.NullInt64{Int64: 1},
		IsUserConcurrent:   sql.NullBool{Bool: false},
		MaxOfflineDays:     sql.NullInt64{Int64: 1},
		Module1:            sql.NullInt64{Int64: input.ModuleID},
		UpdatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:      sql.NullString{String: constanta.SystemClient},
		UpdatedAt:          sql.NullTime{Time: timeNow},
		CreatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:      sql.NullString{String: constanta.SystemClient},
		CreatedAt:          sql.NullTime{Time: timeNow},
		ProductComponentModel: []repository.ProductComponentModel{
			{
				ComponentID:    sql.NullInt64{Int64: input.ComponentID},
				ComponentValue: sql.NullString{String: "Valid"},
			},
		},
	}
	if err, input.ProductNexmileID = input.insertProduct(tx, productModel); err.Error != nil {
		return
	}

	// ---------------- Insert Product NexMileMobile
	productModel = repository.ProductModel{
		ProductID:          sql.NullString{String: "NEXMILE_MOBILE"},
		ProductName:        sql.NullString{String: "NEXMILE_MOBILE"},
		ProductDescription: sql.NullString{String: "VALID1"},
		ProductGroupID:     sql.NullInt64{Int64: input.ProductGroupID},
		ClientTypeID:       sql.NullInt64{Int64: constanta.ResourceTestingNexchiefMobileID},
		IsLicense:          sql.NullBool{Bool: true},
		LicenseVariantID:   sql.NullInt64{Int64: input.LicenseVariantID},
		LicenseTypeID:      sql.NullInt64{Int64: input.LicenseTypeID},
		DeploymentMethod:   sql.NullString{String: "O"},
		NoOfUser:           sql.NullInt64{Int64: 1},
		IsUserConcurrent:   sql.NullBool{Bool: false},
		MaxOfflineDays:     sql.NullInt64{Int64: 1},
		Module1:            sql.NullInt64{Int64: input.ModuleID},
		UpdatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:      sql.NullString{String: constanta.SystemClient},
		UpdatedAt:          sql.NullTime{Time: timeNow},
		CreatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:      sql.NullString{String: constanta.SystemClient},
		CreatedAt:          sql.NullTime{Time: timeNow},
		ProductComponentModel: []repository.ProductComponentModel{
			{
				ComponentID:    sql.NullInt64{Int64: input.ComponentID},
				ComponentValue: sql.NullString{String: "Valid"},
			},
		},
	}
	if err, input.ProductNexchiefMobileID = input.insertProduct(tx, productModel); err.Error != nil {
		return
	}

	return
}

func (input MasterDataTesting) insertProduct(tx *sql.Tx, productModel repository.ProductModel) (err errorModel.ErrorModel, insertedID int64) {
	productModel.ID.Int64, err = dao.ProductDAO.InsertProduct(tx, productModel)
	if err.Error != nil {
		err = ProductService.CheckDuplicateError(err)
		return
	}

	for i := 0; i < len(productModel.ProductComponentModel); i++ {
		productModel.ProductComponentModel[i].ID.Int64, err = dao.ProductComponentDAO.InsertSingleProductComponent(tx, productModel, productModel.ProductComponentModel[i])
		if err.Error != nil {
			return
		}
	}

	insertedID = productModel.ID.Int64
	return
}
