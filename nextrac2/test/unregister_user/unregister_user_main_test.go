package unregister_user

import (
	"database/sql"
	"fmt"
	"github.com/Azure/go-autorest/autorest/date"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/test"
	util2 "nexsoft.co.id/nextrac2/util"
	"os"
	"testing"
	"time"
)

func TestMain(main *testing.M) {
	os.Exit(testMain(main))
}

var masterDataID test.MasterDataTesting

var (
	contextModelNexMile, contextModelNexChiefMobile applicationModel.ContextModel
)

func testMain(main *testing.M) int {
	var err errorModel.ErrorModel
	var errS error
	fmt.Println("Start Testing Penambahan User Admin")

	// Set Configuration
	test.InitAllConfiguration()

	//Truncate function
	if errS = test.Truncate(serverconfig.ServerAttribute.DBConnection); errS != nil {
		fmt.Println("Gagal Truncate")
		return 1
	}

	//Set Database
	if err = test.SetDataWithTransactionalDB(applicationModel.ContextModel{}, setDatabase); err.Error != nil {
		fmt.Println(err)
		fmt.Println(util2.GenerateI18NErrorMessage(err, constanta.DefaultApplicationsLanguage))
		return 1
	}

	return main.Run()
}

func setDatabase(tx *sql.Tx, contextModel applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	var siteId int64

	if err = insertClientCredentials(timeNow, tx); err.Error != nil {
		return
	}

	if err = masterDataID.CreateInitiateData(tx, timeNow); err.Error != nil {
		return
	}

	// Insert Customer Site
	if siteId, err = dao.CustomerSiteDAO.InsertCustomerSite(tx, repository.CustomerInstallationModel{
		ParentCustomerID: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
		CustomerInstallationData: []repository.CustomerInstallationData{
			{
				CustomerID: sql.NullInt64{Int64: masterDataID.CustomerID},
			},
		},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
	}, 0); err.Error != nil {
		return
	}

	// NexMile
	if err = insertLicenseData(tx, repository.ClientMappingModel{
		ClientID:         sql.NullString{String: "98381c991e6b409eb016cfaa365k4cad"},
		SiteID:           sql.NullInt64{Int64: siteId},
		ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceND6ID},
		CustomerID:       sql.NullInt64{Int64: masterDataID.CustomerID},
		CompanyID:        sql.NullString{String: "123"},
		BranchID:         sql.NullString{String: "123"},
		ClientAlias:      sql.NullString{String: "PT. Makmur Sejahtera"},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
		ParentCustomerID: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
	}, &repository.PKCEClientMappingModel{
		ClientID:          sql.NullString{String: "VALID1"},
		ClientTypeID:      sql.NullInt64{Int64: constanta.ResourceTestingNexmileID},
		AuthUserID:        sql.NullInt64{Int64: 2},
		Username:          sql.NullString{String: "USERNAME"},
		ClientAlias:       sql.NullString{String: "ALIAS_VALID"},
		IsClientDependant: sql.NullString{String: "Y"},
		UpdatedBy:         sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:     sql.NullString{String: constanta.SystemClient},
		UpdatedAt:         sql.NullTime{Time: timeNow},
		CreatedBy:         sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:     sql.NullString{String: constanta.SystemClient},
		CreatedAt:         sql.NullTime{Time: timeNow},
	}, repository.CustomerInstallationModel{
		ParentCustomerID: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
		CustomerInstallationData: []repository.CustomerInstallationData{
			{
				CustomerID: sql.NullInt64{Int64: masterDataID.CustomerID},
				SiteID:     sql.NullInt64{Int64: siteId},
				Installation: []repository.CustomerInstallationDetail{
					{
						UniqueID1:          sql.NullString{String: "123"},
						ProductID:          sql.NullInt64{Int64: masterDataID.ProductNexmileID},
						UniqueID2:          sql.NullString{String: "123"},
						Remark:             sql.NullString{String: "VALID"},
						InstallationDate:   sql.NullTime{Time: timeNow},
						InstallationStatus: sql.NullString{String: "A"},
					},
				},
			},
		},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
	}, repository.LicenseConfigModel{
		NoOfUser:         sql.NullInt64{Int64: 20},
		ProductValidFrom: sql.NullTime{Time: date.Date{timeNow.Add(-time.Hour * 24)}.ToTime()},
		ProductValidThru: sql.NullTime{Time: date.Date{timeNow.Add((time.Hour * 24) * 100)}.ToTime()},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
		AllowActivation:  sql.NullString{String: "Y"},
		IsUserConcurrent: sql.NullString{String: "N"},
		ProductID:        sql.NullInt64{Int64: masterDataID.ProductNexmileID},
		ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceTestingNexmileID},
		LicenseVariantID: sql.NullInt64{Int64: masterDataID.LicenseVariantID},
		LicenseTypeID:    sql.NullInt64{Int64: masterDataID.LicenseTypeID},
		DeploymentMethod: sql.NullString{String: "O"},
		MaxOfflineDays:   sql.NullInt64{Int64: 1},
		ModuleID1:        sql.NullInt64{Int64: masterDataID.ModuleID},
		Component: []repository.ProductComponentModel{
			{
				ComponentID:    sql.NullInt64{Int64: masterDataID.ComponentID},
				ComponentValue: sql.NullString{String: "Valid"},
			},
		},
	}, []repository.LicenseConfigComponent{
		{
			LicenseConfigID: sql.NullInt64{},
			ProductID:       sql.NullInt64{Int64: masterDataID.ProductNexmileID},
			ComponentID:     sql.NullInt64{Int64: masterDataID.ComponentID},
			ComponentName:   sql.NullString{String: "Valid"},
			ComponentValue:  sql.NullString{String: "Valid"},
			CreatedBy:       sql.NullInt64{Int64: constanta.SystemID},
			CreatedAt:       sql.NullTime{Time: timeNow},
			CreatedClient:   sql.NullString{String: constanta.SystemClient},
			UpdatedBy:       sql.NullInt64{Int64: constanta.SystemID},
			UpdatedAt:       sql.NullTime{Time: timeNow},
			UpdatedClient:   sql.NullString{String: constanta.SystemClient},
		},
	}, repository.ProductLicenseModel{
		LicenseConfigId:        sql.NullInt64{},
		ProductKey:             sql.NullString{String: "3"},
		ProductEncrypt:         sql.NullString{String: "3"},
		ProductSignature:       sql.NullString{String: "3"},
		ClientId:               sql.NullString{String: "98381c991e6b409eb016cfaa365k4cad"},
		ClientSecret:           sql.NullString{String: "6bf54c4237964a3eb9637da1fb2c622a"},
		HWID:                   sql.NullString{String: "5"},
		ActivationDate:         sql.NullTime{Time: timeNow},
		LicenseStatus:          sql.NullInt32{Int32: 1},
		TerminationDescription: sql.NullString{},
		CreatedBy:              sql.NullInt64{Int64: constanta.SystemID},
		CreatedAt:              sql.NullTime{Time: timeNow},
		CreatedClient:          sql.NullString{String: constanta.SystemClient},
		UpdatedBy:              sql.NullInt64{Int64: constanta.SystemID},
		UpdatedAt:              sql.NullTime{Time: timeNow},
		UpdatedClient:          sql.NullString{String: constanta.SystemClient},
	}, repository.UserLicenseModel{
		ID:               sql.NullInt64{},
		ParentCustomerId: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
		CustomerId:       sql.NullInt64{Int64: masterDataID.CustomerID},
		TotalLicense:     sql.NullInt64{Int64: 20},
		TotalActivated:   sql.NullInt64{Int64: 1},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedAt:        sql.NullTime{Time: timeNow},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
	}, repository.UserRegistrationDetailModel{
		UserID:           sql.NullString{String: "1"},
		Password:         sql.NullString{String: "abc123"},
		SalesmanID:       sql.NullString{String: "1"},
		AndroidID:        sql.NullString{String: "5"},
		RegDate:          sql.NullTime{Time: timeNow},
		Status:           sql.NullString{String: "A"},
		Email:            sql.NullString{String: "123@gmail.com"},
		NoTelp:           sql.NullString{String: "12345"},
		SalesmanCategory: sql.NullString{String: "Salesman1"},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedAt:        sql.NullTime{Time: timeNow},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
	}); err.Error != nil {
		return
	}

	contextModelNexMile = applicationModel.ContextModel{
		AuthAccessTokenModel: model.AuthAccessTokenModel{
			RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
				ResourceUserID: 12,
			},
			ClientID: "98381c991e6b409eb016cfaa365k4cad",
			Locale:   constanta.IndonesianLanguage,
		},
	}

	// NexChief Mobile
	if err = insertLicenseData(tx, repository.ClientMappingModel{
		ClientID:         sql.NullString{String: "r3fb12faf6a348759ccffc500d609f31"},
		SiteID:           sql.NullInt64{Int64: siteId},
		ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceNexChiefID},
		CustomerID:       sql.NullInt64{Int64: masterDataID.CustomerID},
		CompanyID:        sql.NullString{String: "456"},
		BranchID:         sql.NullString{String: "456"},
		ClientAlias:      sql.NullString{String: "PT. Maju Jaya"},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
		ParentCustomerID: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
	}, &repository.PKCEClientMappingModel{
		ClientID:          sql.NullString{String: "VALID2"},
		ClientTypeID:      sql.NullInt64{Int64: constanta.ResourceTestingNexchiefMobileID},
		AuthUserID:        sql.NullInt64{Int64: 1},
		Username:          sql.NullString{String: "USERNAME"},
		ClientAlias:       sql.NullString{String: "ALIAS_VALID"},
		IsClientDependant: sql.NullString{String: "Y"},
		UpdatedBy:         sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:     sql.NullString{String: constanta.SystemClient},
		UpdatedAt:         sql.NullTime{Time: timeNow},
		CreatedBy:         sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:     sql.NullString{String: constanta.SystemClient},
		CreatedAt:         sql.NullTime{Time: timeNow},
	}, repository.CustomerInstallationModel{
		ParentCustomerID: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
		CustomerInstallationData: []repository.CustomerInstallationData{
			{
				CustomerID: sql.NullInt64{Int64: masterDataID.CustomerID},
				SiteID:     sql.NullInt64{Int64: siteId},
				Installation: []repository.CustomerInstallationDetail{
					{
						UniqueID1:          sql.NullString{String: "456"},
						ProductID:          sql.NullInt64{Int64: masterDataID.ProductNexchiefMobileID},
						UniqueID2:          sql.NullString{String: "456"},
						Remark:             sql.NullString{String: "VALID"},
						InstallationDate:   sql.NullTime{Time: timeNow},
						InstallationStatus: sql.NullString{String: "A"},
					},
				},
			},
		},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
	}, repository.LicenseConfigModel{
		NoOfUser:         sql.NullInt64{Int64: 25},
		ProductValidFrom: sql.NullTime{Time: date.Date{timeNow.Add(-time.Hour * 24)}.ToTime()},
		ProductValidThru: sql.NullTime{Time: date.Date{timeNow.Add((time.Hour * 24) * 100)}.ToTime()},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
		AllowActivation:  sql.NullString{String: "Y"},
		IsUserConcurrent: sql.NullString{String: "N"},
		ProductID:        sql.NullInt64{Int64: masterDataID.ProductNexchiefMobileID},
		ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceTestingNexchiefMobileID},
		LicenseVariantID: sql.NullInt64{Int64: masterDataID.LicenseVariantID},
		LicenseTypeID:    sql.NullInt64{Int64: masterDataID.LicenseTypeID},
		DeploymentMethod: sql.NullString{String: "O"},
		MaxOfflineDays:   sql.NullInt64{Int64: 1},
		ModuleID1:        sql.NullInt64{Int64: masterDataID.ModuleID},
		Component: []repository.ProductComponentModel{
			{
				ComponentID:    sql.NullInt64{Int64: masterDataID.ComponentID},
				ComponentValue: sql.NullString{String: "Valid"},
			},
		},
	}, []repository.LicenseConfigComponent{
		{
			LicenseConfigID: sql.NullInt64{},
			ProductID:       sql.NullInt64{Int64: masterDataID.ProductNexchiefMobileID},
			ComponentID:     sql.NullInt64{Int64: masterDataID.ComponentID},
			ComponentName:   sql.NullString{String: "Valid"},
			ComponentValue:  sql.NullString{String: "Valid"},
			CreatedBy:       sql.NullInt64{Int64: constanta.SystemID},
			CreatedAt:       sql.NullTime{Time: timeNow},
			CreatedClient:   sql.NullString{String: constanta.SystemClient},
			UpdatedBy:       sql.NullInt64{Int64: constanta.SystemID},
			UpdatedAt:       sql.NullTime{Time: timeNow},
			UpdatedClient:   sql.NullString{String: constanta.SystemClient},
		},
	}, repository.ProductLicenseModel{
		LicenseConfigId:        sql.NullInt64{},
		ProductKey:             sql.NullString{String: "4"},
		ProductEncrypt:         sql.NullString{String: "4"},
		ProductSignature:       sql.NullString{String: "4"},
		ClientId:               sql.NullString{String: "r3fb12faf6a348759ccffc500d609f31"},
		ClientSecret:           sql.NullString{String: "8kj40eb8063d4513beda8357948a7132"},
		HWID:                   sql.NullString{String: "6"},
		ActivationDate:         sql.NullTime{Time: timeNow},
		LicenseStatus:          sql.NullInt32{Int32: 1},
		TerminationDescription: sql.NullString{},
		CreatedBy:              sql.NullInt64{Int64: constanta.SystemID},
		CreatedAt:              sql.NullTime{Time: timeNow},
		CreatedClient:          sql.NullString{String: constanta.SystemClient},
		UpdatedBy:              sql.NullInt64{Int64: constanta.SystemID},
		UpdatedAt:              sql.NullTime{Time: timeNow},
		UpdatedClient:          sql.NullString{String: constanta.SystemClient},
	}, repository.UserLicenseModel{
		ID:               sql.NullInt64{},
		ParentCustomerId: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
		CustomerId:       sql.NullInt64{Int64: masterDataID.CustomerID},
		TotalLicense:     sql.NullInt64{Int64: 25},
		TotalActivated:   sql.NullInt64{Int64: 1},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedAt:        sql.NullTime{Time: timeNow},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
	}, repository.UserRegistrationDetailModel{
		UserID:           sql.NullString{String: "4"},
		Password:         sql.NullString{String: "acf314"},
		SalesmanID:       sql.NullString{String: "4"},
		AndroidID:        sql.NullString{String: "6"},
		RegDate:          sql.NullTime{Time: timeNow},
		Status:           sql.NullString{String: "A"},
		Email:            sql.NullString{String: "kb12@gmail.com"},
		NoTelp:           sql.NullString{String: "4321"},
		SalesmanCategory: sql.NullString{String: "Salesman4"},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedAt:        sql.NullTime{Time: timeNow},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
	}); err.Error != nil {
		return
	}

	contextModelNexChiefMobile = applicationModel.ContextModel{
		AuthAccessTokenModel: model.AuthAccessTokenModel{
			RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
				ResourceUserID: 12,
			},
			ClientID: "r3fb12faf6a348759ccffc500d609f31",
			Locale:   constanta.IndonesianLanguage,
		},
	}

	return
}

func insertClientCredentials(timeNow time.Time, tx *sql.Tx) (err errorModel.ErrorModel) {

	// NexMile
	if _, err = dao.ClientCredentialDAO.InsertClientCredential(tx, &repository.ClientCredentialModel{
		ClientID:      sql.NullString{String: "98381c991e6b409eb016cfaa365k4cad"},
		ClientSecret:  sql.NullString{String: "6bf54c4237964a3eb9637da1fb2c622a"},
		SignatureKey:  sql.NullString{String: "8k4d9968c4154d698362087a91a80l4a"},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// NexChief Mobile
	if _, err = dao.ClientCredentialDAO.InsertClientCredential(tx, &repository.ClientCredentialModel{
		ClientID:      sql.NullString{String: "r3fb12faf6a348759ccffc500d609f31"},
		ClientSecret:  sql.NullString{String: "8kj40eb8063d4513beda8357948a7132"},
		SignatureKey:  sql.NullString{String: "3100d9968c4154d698362087a91a80nkf"},
		UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient: sql.NullString{String: constanta.SystemClient},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient: sql.NullString{String: constanta.SystemClient},
		CreatedAt:     sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	return
}

func insertLicenseData(tx *sql.Tx, clientMappingData repository.ClientMappingModel, pkceClientMappingData *repository.PKCEClientMappingModel, customerInstallationData repository.CustomerInstallationModel, licenseConfigData repository.LicenseConfigModel, licenseConfigComponent []repository.LicenseConfigComponent, productLicenseData repository.ProductLicenseModel, userLicenseData repository.UserLicenseModel, userRegDetailData repository.UserRegistrationDetailModel) (err errorModel.ErrorModel) {

	// Insert Client Mapping
	clientMappingData.ID.Int64, err = dao.ClientMappingDAO.InsertClientMapping(tx, &clientMappingData)
	if err.Error != nil {
		return
	}

	// Insert Customer Installation
	if customerInstallationData.CustomerInstallationData[0].Installation[0].InstallationID.Int64, err = dao.CustomerInstallationDAO.InsertCustomerInstallationForTesting(tx, customerInstallationData, 0, 0, 0, clientMappingData.ID.Int64); err.Error != nil {
		return
	}

	// Insert PKCE Client Mapping
	pkceClientMappingData.ParentClientID.String = clientMappingData.ClientID.String
	pkceClientMappingData.CustomerID.Int64 = clientMappingData.CustomerID.Int64
	pkceClientMappingData.SiteID.Int64 = clientMappingData.SiteID.Int64
	pkceClientMappingData.CompanyID.String = clientMappingData.CompanyID.String
	pkceClientMappingData.BranchID.String = clientMappingData.BranchID.String

	_, err = dao.PKCEClientMappingDAO.InsertPKCEClientMapping(tx, pkceClientMappingData, true)
	if err.Error != nil {
		return
	}

	// Insert License Config Model
	licenseConfigData.InstallationID.Int64 = customerInstallationData.CustomerInstallationData[0].Installation[0].InstallationID.Int64
	licenseConfigData.ParentCustomerID = customerInstallationData.ParentCustomerID
	licenseConfigData.CustomerID = customerInstallationData.CustomerInstallationData[0].CustomerID
	licenseConfigData.UniqueID1 = customerInstallationData.CustomerInstallationData[0].Installation[0].UniqueID1
	licenseConfigData.UniqueID2 = customerInstallationData.CustomerInstallationData[0].Installation[0].UniqueID2

	licenseConfigData.SiteID = clientMappingData.SiteID
	licenseConfigData.ClientID = clientMappingData.ClientID

	productLicenseData.LicenseConfigId.Int64, err = dao.LicenseConfigDAO.InsertLicenseConfigForTesting(tx, licenseConfigData)
	if err.Error != nil {
		return
	}

	licenseConfigComponent[0].LicenseConfigID.Int64 = productLicenseData.LicenseConfigId.Int64
	_, err = dao.LicenseConfigProductComponentDAO.InsertLicenseConfigProductComponent(tx, licenseConfigComponent)
	if err.Error != nil {
		return
	}

	// Insert Product License
	userLicenseData.ProductLicenseID.Int64, err = dao.ProductLicenseDAO.InsertProductLicense(tx, productLicenseData)

	// Insert User License
	userLicenseData.SiteId = clientMappingData.SiteID
	userLicenseData.ClientID.String = clientMappingData.ClientID.String
	userLicenseData.UniqueId1.String = clientMappingData.CompanyID.String
	userLicenseData.UniqueId2.String = clientMappingData.BranchID.String

	userLicenseData.InstallationId.Int64 = licenseConfigData.InstallationID.Int64
	userLicenseData.ProductValidFrom.Time = licenseConfigData.ProductValidFrom.Time
	userLicenseData.ProductValidThru.Time = licenseConfigData.ProductValidThru.Time

	userRegDetailData.UserLicenseID.Int64, err = dao.UserLicenseDAO.InsertUserLicense(tx, userLicenseData)
	if err.Error != nil {
		return
	}

	// Insert ke user_reg_detail
	userRegDetailData.ParentCustomerID.Int64 = userLicenseData.ParentCustomerId.Int64
	userRegDetailData.CustomerID.Int64 = userLicenseData.CustomerId.Int64
	userRegDetailData.SiteID.Int64 = userLicenseData.SiteId.Int64
	userRegDetailData.InstallationID.Int64 = userLicenseData.InstallationId.Int64
	userRegDetailData.ProductValidThru.Time = userLicenseData.ProductValidThru.Time
	userRegDetailData.ProductValidFrom.Time = userLicenseData.ProductValidFrom.Time

	userRegDetailData.ClientID.String = pkceClientMappingData.ClientID.String
	userRegDetailData.UniqueID1.String = pkceClientMappingData.CompanyID.String
	userRegDetailData.UniqueID2.String = pkceClientMappingData.BranchID.String
	userRegDetailData.AuthUserID.Int64 = pkceClientMappingData.AuthUserID.Int64

	_, err = dao.UserRegistrationDetailDAO.InsertUserRegistrationDetail(tx, userRegDetailData, false)
	if err.Error != nil {
		return
	}

	return
}
