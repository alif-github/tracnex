package register_clientmapping

import (
	"database/sql"
	"errors"
	"github.com/jarcoal/httpmock"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/test"
	"os"
	"testing"
	"time"
)

var InitiateTestVar initiateTest

type initiateTest struct {
	Tx           *sql.Tx
	ContextModel *applicationModel.ContextModel
}

func TestMain(m *testing.M) {
	var err error

	httpmock.Activate()
	defer func() {
		httpmock.DeactivateAndReset()
		test.RollBackSchema(serverconfig.ServerAttribute.DBConnection)
	}()

	err = doTestMain()
	if err != nil {
		os.Exit(1)
	}

	m.Run()
}

func doTestMain() (err error) {
	var tx *sql.Tx

	//------------------ Open All Configuration
	db := test.InitAllConfiguration()

	//------------------ Setting Mocking Auth Server
	modelAuth := []struct {
		clientID     string
		clientSecret string
		signatureKey string
		aliasName    string
	}{
		{
			clientID:     "08181c991e6b409eb016cfaa365b439d",
			clientSecret: "6bf54c4237964a3eb9637da1fb2c622a",
			signatureKey: "d17c040f6287425fafbe7a15315fb31a",
			aliasName:    "Pt Eka Artha",
		},
		{
			clientID:     "1a2b12faf6a345759ccffc500d609b52",
			clientSecret: "47d40eb8063d4513beda8357948a1040",
			signatureKey: "3100d9968c4154d698362087a91a80e1a",
			aliasName:    "Pt Manohara Asri",
		},
	}

	//------------------ Setting Mocking Auth Server
	for idx := range modelAuth {
		test.SetMockAuthServerGetCredential(modelAuth[idx].clientID, modelAuth[idx].clientSecret, modelAuth[idx].signatureKey, modelAuth[idx].aliasName)
	}

	//------------------ Open Tx Database
	tx, err = db.Begin()
	if err != nil {
		return err
	}

	//------------------ Set ContextModel
	contextModel := &applicationModel.ContextModel{
		AuthAccessTokenModel: model.AuthAccessTokenModel{
			RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
				ResourceUserID: 12,
			},
			ClientID: "3e3cb40e14d645eb8783f53a30c822d4",
			Locale:   constanta.IndonesianLanguage,
		},
	}

	//------------------ Set To General Model
	InitiateTestVar = initiateTest{Tx: tx, ContextModel: contextModel}

	//------------------ Set TimeNow
	timeNow := time.Now()

	//------------------ Set The Database
	errs := setDataOnDatabaseLocal(timeNow)
	if errs.Error != nil {
		return errors.New("failure in database")
	}

	return nil
}

func setDataOnDatabaseLocal(timeNow time.Time) (err errorModel.ErrorModel) {
	dataIDDatabase := struct {
		licenseVariantID, licenseTypeID, customerCategoryID,
		customerGroupID, salesmanID, parentCustomerID,
		customerID, productGroupID, moduleID,
		productIDND6, productIDNc, siteID, ncClientTypeID int64
	}{}

	// ---------------- Insert Client Credential
	if _, err = dao.ClientCredentialDAO.InsertClientCredential(InitiateTestVar.Tx, &repository.ClientCredentialModel{
		ClientID:      sql.NullString{String: "08181c991e6b409eb016cfaa365b439d"},
		ClientSecret:  sql.NullString{String: "6bf54c4237964a3eb9637da1fb2c622a"},
		SignatureKey:  sql.NullString{String: "d17c040f6287425fafbe7a15315fb31a"},
		CreatedBy:     sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Client Credential
	if _, err = dao.ClientCredentialDAO.InsertClientCredential(InitiateTestVar.Tx, &repository.ClientCredentialModel{
		ClientID:      sql.NullString{String: "1a2b12faf6a345759ccffc500d609b52"},
		ClientSecret:  sql.NullString{String: "47d40eb8063d4513beda8357948a1040"},
		SignatureKey:  sql.NullString{String: "3100d9968c4154d698362087a91a80e1a"},
		CreatedBy:     sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Client Type (Nexchief)
	if dataIDDatabase.ncClientTypeID, err = dao.ClientTypeDAO.InsertClientType(InitiateTestVar.Tx, repository.ClientTypeModel{
		ClientType:    sql.NullString{String: "Nexchief"},
		Description:   sql.NullString{String: "Client type Nexchief"},
		CreatedBy:     sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert License Variant
	if dataIDDatabase.licenseVariantID, err = dao.LicenseVariantDAO.InsertLicenseVariant(InitiateTestVar.Tx, repository.LicenseVariantModel{
		LicenseVariantName: sql.NullString{String: "VALID"},
		CreatedBy:          sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:      sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:          sql.NullTime{Time: timeNow},
		UpdatedBy:          sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:      sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:          sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert License Type
	if dataIDDatabase.licenseTypeID, err = dao.LicenseTypeDAO.InsertLicenseType(InitiateTestVar.Tx, repository.LicenseTypeModel{
		LicenseTypeName: sql.NullString{String: "VALID"},
		CreatedBy:       sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:   sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:       sql.NullTime{Time: timeNow},
		UpdatedBy:       sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:   sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:       sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Customer Category
	if dataIDDatabase.customerCategoryID, err = dao.CustomerCategoryDAO.InsertCustomerCategory(InitiateTestVar.Tx, repository.CustomerCategoryModel{
		CustomerCategoryID:   sql.NullString{String: "VALID"},
		CustomerCategoryName: sql.NullString{String: "VALID"},
		CreatedBy:            sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:        sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:            sql.NullTime{Time: timeNow},
		UpdatedBy:            sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:        sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:            sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Customer Group
	if dataIDDatabase.customerGroupID, err = dao.CustomerGroupDAO.InsertCustomerGroup(InitiateTestVar.Tx, repository.CustomerGroupModel{
		CustomerGroupID:   sql.NullString{String: "VALID"},
		CustomerGroupName: sql.NullString{String: "VALID"},
		CreatedBy:         sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:     sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:         sql.NullTime{Time: timeNow},
		UpdatedBy:         sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:     sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:         sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Salesman
	if dataIDDatabase.salesmanID, err = dao.SalesmanDAO.InsertSalesman(InitiateTestVar.Tx, repository.SalesmanModel{
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
		CreatedBy:     sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Parent Customer
	if dataIDDatabase.parentCustomerID, err = dao.CustomerDAO.InsertCustomer(InitiateTestVar.Tx, repository.CustomerModel{
		IsPrincipal:             sql.NullBool{Bool: true},
		IsParent:                sql.NullBool{Bool: true},
		MDBCompanyProfileID:     sql.NullInt64{Int64: 1},
		MDBCompanyTitleID:       sql.NullInt64{Int64: 1},
		Npwp:                    sql.NullString{String: "VALID"},
		CompanyTitle:            sql.NullString{String: "VALID"},
		CustomerName:            sql.NullString{String: "VALID"},
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
		SalesmanID:              sql.NullInt64{Int64: dataIDDatabase.salesmanID},
		DistributorOF:           sql.NullString{String: "VALID"},
		CustomerGroupID:         sql.NullInt64{Int64: dataIDDatabase.customerGroupID},
		CustomerCategoryID:      sql.NullInt64{Int64: dataIDDatabase.customerCategoryID},
		Status:                  sql.NullString{String: "A"},
		CreatedBy:               sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:           sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:               sql.NullTime{Time: timeNow},
		UpdatedBy:               sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:           sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:               sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Customer
	if dataIDDatabase.customerID, err = dao.CustomerDAO.InsertCustomer(InitiateTestVar.Tx, repository.CustomerModel{
		IsPrincipal:             sql.NullBool{Bool: true},
		IsParent:                sql.NullBool{Bool: false},
		ParentCustomerID:        sql.NullInt64{Int64: dataIDDatabase.parentCustomerID},
		MDBCompanyProfileID:     sql.NullInt64{Int64: 1},
		MDBCompanyTitleID:       sql.NullInt64{Int64: 1},
		Npwp:                    sql.NullString{String: "VALID"},
		CompanyTitle:            sql.NullString{String: "VALID"},
		CustomerName:            sql.NullString{String: "VALID"},
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
		SalesmanID:              sql.NullInt64{Int64: dataIDDatabase.salesmanID},
		DistributorOF:           sql.NullString{String: "VALID"},
		CustomerGroupID:         sql.NullInt64{Int64: dataIDDatabase.customerGroupID},
		CustomerCategoryID:      sql.NullInt64{Int64: dataIDDatabase.customerCategoryID},
		Status:                  sql.NullString{String: "A"},
		CreatedBy:               sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:           sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:               sql.NullTime{Time: timeNow},
		UpdatedBy:               sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:           sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:               sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Product Group
	if dataIDDatabase.productGroupID, err = dao.ProductGroupDAO.InsertProductGroup(InitiateTestVar.Tx, repository.ProductGroupModel{
		ProductGroupName: sql.NullString{String: "VALID"},
		CreatedBy:        sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:    sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:        sql.NullTime{Time: timeNow},
		UpdatedBy:        sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:    sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Module
	if dataIDDatabase.moduleID, err = dao.ModuleDAO.InsertModule(InitiateTestVar.Tx, repository.ModuleModel{
		ModuleName:    sql.NullString{String: "VALID"},
		CreatedBy:     sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Product
	if dataIDDatabase.productIDND6, err = dao.ProductDAO.InsertProduct(InitiateTestVar.Tx, repository.ProductModel{
		ProductID:          sql.NullString{String: "VALID1"},
		ProductName:        sql.NullString{String: "VALID1"},
		ProductDescription: sql.NullString{String: "VALID1"},
		ProductGroupID:     sql.NullInt64{Int64: dataIDDatabase.productGroupID},
		ClientTypeID:       sql.NullInt64{Int64: 1},
		IsLicense:          sql.NullBool{Bool: true},
		LicenseVariantID:   sql.NullInt64{Int64: dataIDDatabase.licenseVariantID},
		LicenseTypeID:      sql.NullInt64{Int64: dataIDDatabase.licenseTypeID},
		DeploymentMethod:   sql.NullString{String: "O"},
		NoOfUser:           sql.NullInt64{Int64: 1},
		IsUserConcurrent:   sql.NullBool{Bool: true},
		MaxOfflineDays:     sql.NullInt64{Int64: 1},
		Module1:            sql.NullInt64{Int64: dataIDDatabase.moduleID},
		CreatedBy:          sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:      sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:          sql.NullTime{Time: timeNow},
		UpdatedBy:          sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:      sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:          sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// ---------------- Insert Product
	if dataIDDatabase.productIDNc, err = dao.ProductDAO.InsertProduct(InitiateTestVar.Tx, repository.ProductModel{
		ProductID:          sql.NullString{String: "VALID2"},
		ProductName:        sql.NullString{String: "VALID2"},
		ProductDescription: sql.NullString{String: "VALID2"},
		ProductGroupID:     sql.NullInt64{Int64: dataIDDatabase.productGroupID},
		ClientTypeID:       sql.NullInt64{Int64: dataIDDatabase.ncClientTypeID},
		IsLicense:          sql.NullBool{Bool: true},
		LicenseVariantID:   sql.NullInt64{Int64: dataIDDatabase.licenseVariantID},
		LicenseTypeID:      sql.NullInt64{Int64: dataIDDatabase.licenseTypeID},
		DeploymentMethod:   sql.NullString{String: "O"},
		NoOfUser:           sql.NullInt64{Int64: 1},
		IsUserConcurrent:   sql.NullBool{Bool: true},
		MaxOfflineDays:     sql.NullInt64{Int64: 1},
		Module1:            sql.NullInt64{Int64: dataIDDatabase.moduleID},
		CreatedBy:          sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:      sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:          sql.NullTime{Time: timeNow},
		UpdatedBy:          sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:      sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:          sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	customerInstallationModel := repository.CustomerInstallationModel{
		ParentCustomerID: sql.NullInt64{Int64: dataIDDatabase.parentCustomerID},
		CustomerInstallationData: []repository.CustomerInstallationData{
			{
				CustomerID: sql.NullInt64{Int64: dataIDDatabase.customerID},
				Installation: []repository.CustomerInstallationDetail{
					{
						ProductID:          sql.NullInt64{Int64: dataIDDatabase.productIDND6},
						UniqueID1:          sql.NullString{String: "NS6024050001031"},
						UniqueID2:          sql.NullString{String: "1468381449586"},
						Remark:             sql.NullString{String: "VALID"},
						InstallationDate:   sql.NullTime{Time: timeNow},
						InstallationStatus: sql.NullString{String: "A"},
					},
					{
						ProductID:          sql.NullInt64{Int64: dataIDDatabase.productIDNc},
						UniqueID1:          sql.NullString{String: "NDI"},
						Remark:             sql.NullString{String: "VALID"},
						InstallationDate:   sql.NullTime{Time: timeNow},
						InstallationStatus: sql.NullString{String: "A"},
					},
				},
			},
		},
		CreatedBy:     sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: InitiateTestVar.ContextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: InitiateTestVar.ContextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	// ---------------- Insert Customer Site
	if dataIDDatabase.siteID, err = dao.CustomerSiteDAO.InsertCustomerSite(InitiateTestVar.Tx, customerInstallationModel, 0); err.Error != nil {
		return
	}

	customerInstallationModel.CustomerInstallationData[0].SiteID.Int64 = dataIDDatabase.siteID

	// ---------------- Insert Customer Installation
	var queue int64
	for idx := range customerInstallationModel.CustomerInstallationData[0].Installation {
		if _, err = dao.CustomerInstallationDAO.InsertCustomerInstallation(InitiateTestVar.Tx, customerInstallationModel, 0, idx, queue); err.Error != nil {
			return
		}
		queue++
	}

	defer func() {
		if err.Error != nil {
			errorS := InitiateTestVar.Tx.Rollback()
			if errorS != nil {
				err = errorModel.GenerateUnknownError("Testing", "Testing", errorS)
				return
			}
			return
		} else {
			errorS := InitiateTestVar.Tx.Commit()
			if errorS != nil {
				err = errorModel.GenerateUnknownError("Testing", "Testing", errorS)
				return
			}
		}
	}()

	return
}
