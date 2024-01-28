package validation_named_user

import (
	"database/sql"
	"fmt"
	"github.com/Azure/go-autorest/autorest/date"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service/ProductService"
	"nexsoft.co.id/nextrac2/test"
	util2 "nexsoft.co.id/nextrac2/util"
	"os"
	"strconv"
	"testing"
	"time"
)

var masterDataID test.MasterDataTesting

type modelSetLicenseData struct {
	clientCredentialModel    repository.ClientCredentialModel
	clientMappingData        repository.ClientMappingModel
	customerInstallationData repository.CustomerInstallationModel
}

type modelInsertUserLicense struct {
	licenseConfigData     repository.LicenseConfigModel
	inputStruct           in.RegisterNamedUserRequest
	clientCredentialModel repository.ClientCredentialModel
	isRegisStatus         bool
}

func TestMain(main *testing.M) {
	os.Exit(testMain(main))
}

func testMain(main *testing.M) int {
	var err errorModel.ErrorModel
	var errS error
	fmt.Println("Start Testing Validation User")

	// Set Configuration
	test.InitAllConfiguration()

	//Truncate function
	if errS = test.Truncate(serverconfig.ServerAttribute.DBConnection); errS != nil {
		fmt.Println("Gagal Buat Client Type")
		return 1
	}

	//Set Master Database
	if err = test.SetDataWithTransactionalDB(applicationModel.ContextModel{}, setMasterData); err.Error != nil {
		fmt.Println(err)
		fmt.Println(util2.GenerateI18NErrorMessage(err, constanta.DefaultApplicationsLanguage))
		return 1
	}

	//Set Data Validation User
	if err = test.SetDataWithTransactionalDB(applicationModel.ContextModel{}, setDataValidationUser); err.Error != nil {
		fmt.Println(err)
		fmt.Println(util2.GenerateI18NErrorMessage(err, constanta.DefaultApplicationsLanguage))
		return 1
	}

	return main.Run()
}

func setMasterData(tx *sql.Tx, contextModel applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	if err = masterDataID.CreateInitiateData(tx, timeNow); err.Error != nil {
		return
	}
	return
}

func setDataValidationUser(tx *sql.Tx, contextModel applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	var userLicenseID int64
	var licenseConfigData repository.LicenseConfigModel

	// Insert license Data
	if userLicenseID, licenseConfigData, err = setLicenseData(tx, timeNow, modelSetLicenseData{
		clientCredentialModel: repository.ClientCredentialModel{
			ClientID:      sql.NullString{String: "08181c991e6b409eb016cfaa365b439d"},
			ClientSecret:  sql.NullString{String: "6bf54c4237964a3eb9637da1fb2c622a"},
			SignatureKey:  sql.NullString{String: "280d9968c4154d698362087a91a80e1a"},
			UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient: sql.NullString{String: constanta.SystemClient},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient: sql.NullString{String: constanta.SystemClient},
			CreatedAt:     sql.NullTime{Time: timeNow},
		},
		clientMappingData: repository.ClientMappingModel{
			ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceNexmileID},
			CustomerID:       sql.NullInt64{Int64: masterDataID.CustomerID},
			CompanyID:        sql.NullString{String: "NS6024050001031"},
			BranchID:         sql.NullString{String: "1468381449586"},
			ClientAlias:      sql.NullString{String: "ND6 - PT. Eka Artha Buanas 111"},
			UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:    sql.NullString{String: constanta.SystemClient},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:    sql.NullString{String: constanta.SystemClient},
			CreatedAt:        sql.NullTime{Time: timeNow},
			ParentCustomerID: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
		},
		customerInstallationData: repository.CustomerInstallationModel{
			ParentCustomerID: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
			CustomerInstallationData: []repository.CustomerInstallationData{
				{
					CustomerID: sql.NullInt64{Int64: masterDataID.CustomerID},
					Installation: []repository.CustomerInstallationDetail{
						{
							ProductID:          sql.NullInt64{Int64: masterDataID.ProductNexmileID},
							Remark:             sql.NullString{String: "VALID"},
							InstallationDate:   sql.NullTime{Time: timeNow},
							InstallationStatus: sql.NullString{String: "A"},
							NoOfUser:           sql.NullInt64{Int64: 1},
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
		},
	}); err.Error != nil {
		return
	}

	// Insert User Registration Detail
	if _, err = insertUserRegistDetail(userLicenseID, tx, timeNow, modelInsertUserLicense{
		licenseConfigData: licenseConfigData,
		inputStruct: in.RegisterNamedUserRequest{
			ParentClientID:   licenseConfigData.ClientID.String,
			ClientID:         "62cda0a7242c497bbc502e4b33e87abc",
			ClientTypeID:     constanta.ResourceNexmileID,
			AuthUserID:       100,
			Firstname:        "Pertama",
			Lastname:         "Terakhir",
			Username:         "UsernameValid",
			UserID:           "100",
			Password:         "PassValid",
			ClientAliases:    "ValidAlias",
			SalesmanID:       "ValidSalesID",
			AndroidID:        "ValidAndroID",
			Email:            "ValidEmail@email.com",
			UniqueID1:        licenseConfigData.UniqueID1.String,
			UniqueID2:        licenseConfigData.UniqueID2.String,
			NoTelp:           "+62-897789666889",
			SalesmanCategory: "ValidCtg",
		},
		clientCredentialModel: repository.ClientCredentialModel{
			ClientID:     sql.NullString{String: "08181c991e6b409eb016cfaa365b439d"},
			ClientSecret: sql.NullString{String: "6bf54c4237964a3eb9637da1fb2c622a"},
			SignatureKey: sql.NullString{String: "280d9968c4154d698362087a91a80e1a"},
		},
		isRegisStatus: false,
	}); err.Error != nil {
		return
	}

	return
}

func setLicenseData(tx *sql.Tx, timeNow time.Time, data modelSetLicenseData) (userLicenseID int64, licenseConfigData repository.LicenseConfigModel, err errorModel.ErrorModel) {
	var productData out.ViewProduct
	var siteID int64

	// Insert Client Credential
	if _, err = dao.ClientCredentialDAO.InsertClientCredential(tx, &data.clientCredentialModel); err.Error != nil {
		return
	}

	// ---------------- Insert Customer Site
	if siteID, err = dao.CustomerSiteDAO.InsertCustomerSite(tx, data.customerInstallationData, 0); err.Error != nil {
		return
	}
	data.clientMappingData.SiteID.Int64 = siteID

	// ---------------- Insert Client Mapping
	data.clientMappingData.ClientID = data.clientCredentialModel.ClientID
	data.clientMappingData.ID.Int64, err = dao.ClientMappingDAO.InsertClientMapping(tx, &data.clientMappingData)
	if err.Error != nil {
		return
	}

	// ---------------- Insert Customer Installation
	data.customerInstallationData.CustomerInstallationData[0].SiteID.Int64 = siteID
	var queue int64
	for idx := range data.customerInstallationData.CustomerInstallationData[0].Installation {
		// get product
		productData, err = ProductService.ProductService.DoViewProduct(in.ProductRequest{
			ID: data.customerInstallationData.CustomerInstallationData[0].Installation[idx].ProductID.Int64,
		}, &applicationModel.ContextModel{})

		// insert customer Installation
		data.customerInstallationData.CustomerInstallationData[0].Installation[idx].UniqueID1 = data.clientMappingData.CompanyID
		data.customerInstallationData.CustomerInstallationData[0].Installation[idx].UniqueID2 = data.clientMappingData.BranchID
		if data.customerInstallationData.CustomerInstallationData[0].Installation[idx].InstallationID.Int64, err = dao.CustomerInstallationDAO.InsertCustomerInstallationForTesting(tx, data.customerInstallationData, 0, idx, queue, data.clientMappingData.ID.Int64); err.Error != nil {
			return
		}

		if userLicenseID, licenseConfigData, err = insertLicense(tx, timeNow, data, productData, idx); err.Error != nil {
			return
		}
		queue++
	}

	return
}

func insertLicense(tx *sql.Tx, timeNow time.Time, data modelSetLicenseData, productData out.ViewProduct, idx int) (userLicenseID int64, licenseConfigData repository.LicenseConfigModel, err errorModel.ErrorModel) {
	var (
		licenseConfigID,
		productLicenseID int64
	)

	// insert license config
	licenseConfigData = createLicenseConfigModel(timeNow, data, productData, idx)
	licenseConfigID, err = dao.LicenseConfigDAO.InsertLicenseConfigForTesting(tx, licenseConfigData)
	if err.Error != nil {
		return
	}

	//insert product license
	if productLicenseID, err = dao.ProductLicenseDAO.InsertProductLicense(tx, repository.ProductLicenseModel{
		LicenseConfigId:  sql.NullInt64{Int64: licenseConfigID},
		ProductKey:       sql.NullString{String: data.clientCredentialModel.ClientID.String + strconv.Itoa(int(licenseConfigID)) + strconv.Itoa(idx)},
		ProductEncrypt:   sql.NullString{String: data.clientCredentialModel.ClientSecret.String + strconv.Itoa(int(licenseConfigID)) + strconv.Itoa(idx)},
		ProductSignature: sql.NullString{String: data.clientCredentialModel.SignatureKey.String + strconv.Itoa(int(licenseConfigID)) + strconv.Itoa(idx)},
		ClientId:         data.clientCredentialModel.ClientID,
		ClientSecret:     data.clientCredentialModel.ClientSecret,
		HWID:             sql.NullString{String: "HWID" + strconv.Itoa(int(licenseConfigID))},
		ActivationDate:   sql.NullTime{Time: timeNow},
		LicenseStatus:    sql.NullInt32{Int32: constanta.ProductLicenseStatusActive},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	if !productData.IsConcurrentUser {
		if userLicenseID, err = dao.UserLicenseDAO.InsertUserLicense(tx, repository.UserLicenseModel{
			ProductLicenseID: sql.NullInt64{Int64: productLicenseID},
			ParentCustomerId: licenseConfigData.ParentCustomerID,
			CustomerId:       licenseConfigData.CustomerID,
			SiteId:           licenseConfigData.SiteID,
			InstallationId:   licenseConfigData.InstallationID,
			ClientID:         licenseConfigData.ClientID,
			UniqueId1:        licenseConfigData.UniqueID1,
			UniqueId2:        licenseConfigData.UniqueID2,
			ProductValidFrom: licenseConfigData.ProductValidFrom,
			ProductValidThru: licenseConfigData.ProductValidThru,
			TotalLicense:     licenseConfigData.NoOfUser,
			UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:    sql.NullString{String: constanta.SystemClient},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:    sql.NullString{String: constanta.SystemClient},
			CreatedAt:        sql.NullTime{Time: timeNow},
		}); err.Error != nil {
			return
		}
	}
	return
}

func createLicenseConfigModel(timeNow time.Time, data modelSetLicenseData, productData out.ViewProduct, idxSite int) (licenseConfigData repository.LicenseConfigModel) {
	licenseConfigData = repository.LicenseConfigModel{
		InstallationID:   data.customerInstallationData.CustomerInstallationData[0].Installation[idxSite].InstallationID,
		NoOfUser:         data.customerInstallationData.CustomerInstallationData[0].Installation[idxSite].NoOfUser,
		ProductValidFrom: sql.NullTime{Time: date.Date{timeNow.Add(-time.Hour * 24)}.ToTime()},
		ProductValidThru: sql.NullTime{Time: date.Date{timeNow.Add((time.Hour * 24) * 100)}.ToTime()},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
		ParentCustomerID: data.customerInstallationData.ParentCustomerID,
		CustomerID:       data.customerInstallationData.CustomerInstallationData[0].CustomerID,
		SiteID:           data.customerInstallationData.CustomerInstallationData[0].SiteID,
		ClientID:         data.clientMappingData.ClientID,
		ProductID:        sql.NullInt64{Int64: productData.ID},
		ClientTypeID:     sql.NullInt64{Int64: productData.ClientTypeID},
		LicenseVariantID: sql.NullInt64{Int64: productData.LicenseVariantID},
		LicenseTypeID:    sql.NullInt64{Int64: productData.LicenseTypeID},
		DeploymentMethod: sql.NullString{String: productData.DeploymentMethod},
		MaxOfflineDays:   sql.NullInt64{Int64: productData.MaxOfflineDays},
		UniqueID1:        data.customerInstallationData.CustomerInstallationData[0].Installation[idxSite].UniqueID1,
		UniqueID2:        data.customerInstallationData.CustomerInstallationData[0].Installation[idxSite].UniqueID2,
		ModuleID1:        sql.NullInt64{Int64: productData.ModuleId1},
		AllowActivation:  sql.NullString{String: "Y"},
		Component: []repository.ProductComponentModel{
			{
				ID:             sql.NullInt64{Int64: productData.Component[0].ID},
				ComponentID:    sql.NullInt64{Int64: productData.Component[0].ComponentID},
				ComponentValue: sql.NullString{String: productData.Component[0].ComponentValue},
			},
		},
	}

	if productData.IsConcurrentUser {
		licenseConfigData.IsUserConcurrent.String = "Y"
	} else {
		licenseConfigData.IsUserConcurrent.String = "N"
	}

	return
}

func insertUserRegistDetail(userLicenseID int64, tx *sql.Tx, timeNow time.Time, license modelInsertUserLicense) (userRegistrationID int64, err errorModel.ErrorModel) {
	userRegistrationModel := repository.UserRegistrationDetailModel{
		UserLicenseID:    sql.NullInt64{Int64: userLicenseID},
		ParentCustomerID: license.licenseConfigData.ParentCustomerID,
		CustomerID:       license.licenseConfigData.CustomerID,
		SiteID:           license.licenseConfigData.SiteID,
		InstallationID:   license.licenseConfigData.InstallationID,
		ParentClientID:   license.licenseConfigData.ClientID,
		ClientID:         sql.NullString{String: license.inputStruct.ClientID},
		UniqueID1:        license.licenseConfigData.UniqueID1,
		UniqueID2:        license.licenseConfigData.UniqueID2,
		AuthUserID:       sql.NullInt64{Int64: license.inputStruct.AuthUserID},
		Firstname:        sql.NullString{String: license.inputStruct.Firstname},
		Lastname:         sql.NullString{String: license.inputStruct.Lastname},
		Username:         sql.NullString{String: license.inputStruct.Username},
		UserID:           sql.NullString{String: license.inputStruct.UserID},
		Password:         sql.NullString{String: license.inputStruct.Password},
		ClientAliases:    sql.NullString{String: license.inputStruct.ClientAliases},
		SalesmanID:       sql.NullString{String: license.inputStruct.SalesmanID},
		AndroidID:        sql.NullString{String: license.inputStruct.AndroidID},
		RegDate:          sql.NullTime{Time: timeNow},
		Email:            sql.NullString{String: license.inputStruct.Email},
		NoTelp:           sql.NullString{String: license.inputStruct.NoTelp},
		SalesmanCategory: sql.NullString{String: license.inputStruct.SalesmanCategory},
		ProductValidFrom: license.licenseConfigData.ProductValidFrom,
		ProductValidThru: license.licenseConfigData.ProductValidThru,
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
	}

	// insert User Regis Detail
	if userRegistrationID, err = dao.UserRegistrationDetailDAO.InsertUserRegistrationDetail(tx, userRegistrationModel, license.isRegisStatus); err.Error != nil {
		return
	}

	if !license.isRegisStatus {
		_, err = dao.PKCEClientMappingDAO.InsertPKCEClientMapping(tx, &repository.PKCEClientMappingModel{
			ParentClientID: license.licenseConfigData.ClientID,
			ClientID:       userRegistrationModel.ClientID,
			ClientTypeID:   license.licenseConfigData.ClientTypeID,
			AuthUserID:     userRegistrationModel.AuthUserID,
			Username:       userRegistrationModel.Username,
			InstallationID: userRegistrationModel.InstallationID,
			CustomerID:     userRegistrationModel.CustomerID,
			SiteID:         userRegistrationModel.SiteID,
			CompanyID:      userRegistrationModel.UniqueID1,
			BranchID:       userRegistrationModel.UniqueID2,
			ClientAlias:    userRegistrationModel.ClientAliases,
			UpdatedBy:      sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:  sql.NullString{String: constanta.SystemClient},
			UpdatedAt:      sql.NullTime{Time: timeNow},
			CreatedBy:      sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:  sql.NullString{String: constanta.SystemClient},
			CreatedAt:      sql.NullTime{Time: timeNow},
		}, true)
	}

	return
}
