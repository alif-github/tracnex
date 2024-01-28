package parameter_nexmile

import (
	"database/sql"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/markbates/oncer"
	"nexsoft.co.id/nextrac2/config"
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
	httpmock.Activate()
	defer func() {
		httpmock.DeactivateAndReset()
	}()
	os.Exit(testMain(main))
}

var contextModel applicationModel.ContextModel
var masterDataID test.MasterDataTesting

func testMain(main *testing.M) int {
	var err errorModel.ErrorModel
	var errS error

	//----- Set Configuration
	test.InitAllConfigurations()
	fmt.Println(config.ApplicationConfiguration.GetPostgreSQLDefaultSchema())

	//----- Truncate function
	if errS = test.Truncate(serverconfig.ServerAttribute.DBConnection); errS != nil {
		fmt.Println("Gagal Truncate")
		return 1
	}

	test.SetMockAuthServerGetByAuthUserID(123, "VALID1")

	contextModel = applicationModel.ContextModel{
		AuthAccessTokenModel: model.AuthAccessTokenModel{
			RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
				ResourceUserID: 12,
			},
			ClientID: "98381c991e6b409eb016cfaa365k4cad",
			Locale:   constanta.IndonesianLanguage,
		},
	}

	//----- Set Database
	if err = test.SetDataWithTransactionalDB(contextModel, setDatabase); err.Error != nil {
		fmt.Println(err)
		fmt.Println(util2.GenerateI18NErrorMessage(err, constanta.DefaultApplicationsLanguage))
		return 1
	}

	return main.Run()
}

func setDatabase(tx *sql.Tx, contextModel applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	var siteID int64

	//----- Create Initiate Data
	if err = masterDataID.CreateInitiateData(tx, timeNow); err.Error != nil {
		return
	}

	//----- Create Model Credential
	dataClientCredential := []repository.ClientCredentialModel{
		{
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
		{
			ClientID:      sql.NullString{String: "1a2b12faf6a345759ccffc500d609b52"},
			ClientSecret:  sql.NullString{String: "47d40eb8063d4513beda8357948a1040"},
			SignatureKey:  sql.NullString{String: "bb0734e85ba44b529611fd22668b6bad"},
			UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient: sql.NullString{String: constanta.SystemClient},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient: sql.NullString{String: constanta.SystemClient},
			CreatedAt:     sql.NullTime{Time: timeNow},
		},
		{
			ClientID:      sql.NullString{String: "98381c991e6b409eb016cfaa365k4cad"},
			ClientSecret:  sql.NullString{String: "6bf54c4237964a3eb9637da1fb2c622a"},
			SignatureKey:  sql.NullString{String: "8k4d9968c4154d698362087a91a80l4a"},
			UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient: sql.NullString{String: constanta.SystemClient},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient: sql.NullString{String: constanta.SystemClient},
			CreatedAt:     sql.NullTime{Time: timeNow},
		},
		{
			ClientID:      sql.NullString{String: "r3fb12faf6a348759ccffc500d609f31"},
			ClientSecret:  sql.NullString{String: "8kj40eb8063d4513beda8357948a7132"},
			SignatureKey:  sql.NullString{String: "3100d9968c4154d698362087a91a80nkf"},
			UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient: sql.NullString{String: constanta.SystemClient},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient: sql.NullString{String: constanta.SystemClient},
			CreatedAt:     sql.NullTime{Time: timeNow},
		},
	}

	//----- Create Model Client Mapping
	dataClientMapping := []repository.ClientMappingModel{
		{
			ClientID:         sql.NullString{String: "98381c991e6b409eb016cfaa365k4cad"},
			ClientTypeID:     sql.NullInt64{Int64: 5},
			CustomerID:       sql.NullInt64{Int64: 3},
			CompanyID:        sql.NullString{String: "123"},
			BranchID:         sql.NullString{String: "123"},
			ClientAlias:      sql.NullString{String: "PT. Makmur Sejahtera"},
			UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:    sql.NullString{String: constanta.SystemClient},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:    sql.NullString{String: constanta.SystemClient},
			CreatedAt:        sql.NullTime{Time: timeNow},
			ParentCustomerID: sql.NullInt64{Int64: 1},
		},
		{
			ClientID:         sql.NullString{String: "r3fb12faf6a348759ccffc500d609f31"},
			ClientTypeID:     sql.NullInt64{Int64: 6},
			CustomerID:       sql.NullInt64{Int64: 3},
			CompanyID:        sql.NullString{String: "456"},
			BranchID:         sql.NullString{String: "456"},
			ClientAlias:      sql.NullString{String: "PT. Maju Jaya"},
			UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:    sql.NullString{String: constanta.SystemClient},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:    sql.NullString{String: constanta.SystemClient},
			CreatedAt:        sql.NullTime{Time: timeNow},
			ParentCustomerID: sql.NullInt64{Int64: 1},
		},
	}

	//----- Create Model PKCE Client Mapping
	dataPKCEClientMapping := []repository.PKCEClientMappingModel{
		repository.PKCEClientMappingModel{
			ParentClientID:    sql.NullString{String: "98381c991e6b409eb016cfaa365k4cad"},
			ClientID:          sql.NullString{String: "VALID1"},
			ClientTypeID:      sql.NullInt64{Int64: 5},
			AuthUserID:        sql.NullInt64{Int64: 2},
			Username:          sql.NullString{String: "USERNAME"},
			CustomerID:        sql.NullInt64{Int64: 3},
			SiteID:            sql.NullInt64{Int64: 1},
			CompanyID:         sql.NullString{String: "123"},
			BranchID:          sql.NullString{String: "123"},
			ClientAlias:       sql.NullString{String: "ALIAS_VALID"},
			IsClientDependant: sql.NullString{String: constanta.FlagStatusTrue},
			UpdatedBy:         sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:     sql.NullString{String: constanta.SystemClient},
			UpdatedAt:         sql.NullTime{Time: timeNow},
			CreatedBy:         sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:     sql.NullString{String: constanta.SystemClient},
			CreatedAt:         sql.NullTime{Time: timeNow},
		},
		{
			ParentClientID:    sql.NullString{String: "r3fb12faf6a348759ccffc500d609f31"},
			ClientID:          sql.NullString{String: "VALID2"},
			ClientTypeID:      sql.NullInt64{Int64: 6},
			AuthUserID:        sql.NullInt64{Int64: 1},
			Username:          sql.NullString{String: "USERNAME"},
			CustomerID:        sql.NullInt64{Int64: 3},
			SiteID:            sql.NullInt64{Int64: 1},
			CompanyID:         sql.NullString{String: "456"},
			BranchID:          sql.NullString{String: "456"},
			ClientAlias:       sql.NullString{String: "ALIAS_VALID"},
			IsClientDependant: sql.NullString{String: constanta.FlagStatusTrue},
			UpdatedBy:         sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:     sql.NullString{String: constanta.SystemClient},
			UpdatedAt:         sql.NullTime{Time: timeNow},
			CreatedBy:         sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:     sql.NullString{String: constanta.SystemClient},
			CreatedAt:         sql.NullTime{Time: timeNow},
		},
	}

	//----- Add Client Credential
	for _, credentialModel := range dataClientCredential {
		if _, err = dao.ClientCredentialDAO.InsertClientCredential(tx, &credentialModel); err.Error != nil {
			return
		}
	}

	//----- Add Customer Site
	if siteID, err = dao.CustomerSiteDAO.InsertCustomerSite(tx, repository.CustomerInstallationModel{
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

	for _, mappingModel := range dataClientMapping {
		var clientMappingIDTemp, productIDTemp,
		installationIDTemp, licenseConfigIDTemp,
		productLicenseIDTemp, userLicenseIDTemp int64

		clientMappingIDTemp, err = dao.ClientMappingDAO.InsertClientMapping(tx, &mappingModel)
		if err.Error != nil {
			return
		}

		if mappingModel.ClientTypeID.Int64 == 5 {
			productIDTemp = masterDataID.ProductNexmileID
		} else if mappingModel.ClientTypeID.Int64 == 6 {
			productIDTemp = masterDataID.ProductNexchiefMobileID
		}

		installationIDTemp, err = dao.CustomerInstallationDAO.InsertCustomerInstallationForTesting(tx, repository.CustomerInstallationModel{
			ParentCustomerID: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
			CustomerInstallationData: []repository.CustomerInstallationData{
				{
					SiteID:     sql.NullInt64{Int64: siteID},
					CustomerID: sql.NullInt64{Int64: masterDataID.CustomerID},
					Installation: []repository.CustomerInstallationDetail{
						{
							UniqueID1:          sql.NullString{String: mappingModel.CompanyID.String},
							UniqueID2:          sql.NullString{String: mappingModel.BranchID.String},
							ProductID:          sql.NullInt64{Int64: productIDTemp},
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
		}, 0, 0, 0, clientMappingIDTemp)
		if err.Error != nil {
			return
		}

		licenseConfigIDTemp, err = dao.LicenseConfigDAO.InsertLicenseConfigForTesting(tx, repository.LicenseConfigModel{
			InstallationID:   sql.NullInt64{Int64: installationIDTemp},
			ParentCustomerID: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
			CustomerID:       sql.NullInt64{Int64: masterDataID.CustomerID},
			SiteID:           sql.NullInt64{Int64: siteID},
			ClientID:         sql.NullString{String: mappingModel.ClientID.String},
			ProductID:        sql.NullInt64{Int64: productIDTemp},
			ClientTypeID:     sql.NullInt64{Int64: mappingModel.ClientTypeID.Int64},
			LicenseVariantID: sql.NullInt64{Int64: masterDataID.LicenseVariantID},
			LicenseTypeID:    sql.NullInt64{Int64: masterDataID.LicenseTypeID},
			DeploymentMethod: sql.NullString{String: "O"},
			AllowActivation:  sql.NullString{String: "Y"},
			NoOfUser:         sql.NullInt64{Int64: 1},
			IsUserConcurrent: sql.NullString{String: "N"},
			MaxOfflineDays:   sql.NullInt64{Int64: 1},
			UniqueID1:        sql.NullString{String: mappingModel.CompanyID.String},
			UniqueID2:        sql.NullString{String: mappingModel.BranchID.String},
			ProductValidFrom: sql.NullTime{Time: timeNow},
			ProductValidThru: sql.NullTime{Time: timeNow.AddDate(1, 0, 0)},
			ModuleID1:        sql.NullInt64{Int64: 1},
			UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:    sql.NullString{String: constanta.SystemClient},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:    sql.NullString{String: constanta.SystemClient},
			CreatedAt:        sql.NullTime{Time: timeNow},
		})
		if err.Error != nil {
			return
		}

		_, err = dao.LicenseConfigProductComponentDAO.InsertLicenseConfigProductComponent(tx, []repository.LicenseConfigComponent{
			{
				LicenseConfigID: sql.NullInt64{Int64: licenseConfigIDTemp},
				ProductID:       sql.NullInt64{Int64: productIDTemp},
				ComponentID:     sql.NullInt64{Int64: masterDataID.ComponentID},
				ComponentValue:  sql.NullString{String: "Valid"},
				UpdatedBy:       sql.NullInt64{Int64: constanta.SystemID},
				UpdatedClient:   sql.NullString{String: constanta.SystemClient},
				UpdatedAt:       sql.NullTime{Time: timeNow},
				CreatedBy:       sql.NullInt64{Int64: constanta.SystemID},
				CreatedClient:   sql.NullString{String: constanta.SystemClient},
				CreatedAt:       sql.NullTime{Time: timeNow},
			},
		})
		if err.Error != nil {
			return
		}

		modelProductLicense := repository.ProductLicenseModel{
			LicenseConfigId:  sql.NullInt64{Int64: licenseConfigIDTemp},
			ProductKey:       sql.NullString{String: "VALID"},
			ProductEncrypt:   sql.NullString{String: "VALID"},
			ProductSignature: sql.NullString{String: "VALID"},
			ClientId:         sql.NullString{String: mappingModel.ClientID.String},
			HWID:             sql.NullString{String: "VALID"},
			ActivationDate:   sql.NullTime{Time: timeNow},
			LicenseStatus:    sql.NullInt32{Int32: 1},
			UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:    sql.NullString{String: constanta.SystemClient},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:    sql.NullString{String: constanta.SystemClient},
			CreatedAt:        sql.NullTime{Time: timeNow},
		}

		for _, secretClient := range dataClientCredential {
			if secretClient.ClientID.String == mappingModel.ClientID.String {
				modelProductLicense.ClientSecret.String = secretClient.ClientSecret.String
			}
		}

		productLicenseIDTemp, err = dao.ProductLicenseDAO.InsertProductLicense(tx, modelProductLicense)
		if err.Error != nil {
			return
		}

		userLicenseIDTemp, err = dao.UserLicenseDAO.InsertUserLicense(tx, repository.UserLicenseModel{
			ProductLicenseID: sql.NullInt64{Int64: productLicenseIDTemp},
			ParentCustomerId: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
			CustomerId:       sql.NullInt64{Int64: masterDataID.CustomerID},
			SiteId:           sql.NullInt64{Int64: siteID},
			InstallationId:   sql.NullInt64{Int64: installationIDTemp},
			ClientID:         sql.NullString{String: mappingModel.ClientID.String},
			UniqueId1:        sql.NullString{String: mappingModel.CompanyID.String},
			UniqueId2:        sql.NullString{String: mappingModel.BranchID.String},
			ProductValidFrom: sql.NullTime{Time: timeNow},
			ProductValidThru: sql.NullTime{Time: timeNow.AddDate(1, 0, 0)},
			TotalLicense:     sql.NullInt64{Int64: 1},
			TotalActivated:   sql.NullInt64{Int64: 1},
			UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:    sql.NullString{String: constanta.SystemClient},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:    sql.NullString{String: constanta.SystemClient},
			CreatedAt:        sql.NullTime{Time: timeNow},
		})

		oncer.Do("run once pkce client mapping", func() {
			for _, mappingModelPKCEClientMapping := range dataPKCEClientMapping {
				_, err = dao.PKCEClientMappingDAO.InsertPKCEClientMapping(tx, &mappingModelPKCEClientMapping, true)
				if err.Error != nil {
					return
				}
			}
		})

		modelUserRegistrationDetail := repository.UserRegistrationDetailModel{
			UserLicenseID:    sql.NullInt64{Int64: userLicenseIDTemp},
			ParentCustomerID: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
			CustomerID:       sql.NullInt64{Int64: masterDataID.CustomerID},
			SiteID:           sql.NullInt64{Int64: siteID},
			InstallationID:   sql.NullInt64{Int64: installationIDTemp},
			UniqueID1:        sql.NullString{String: mappingModel.CompanyID.String},
			UniqueID2:        sql.NullString{String: mappingModel.BranchID.String},
			UserID:           sql.NullString{String: "1"},
			Password:         sql.NullString{String: "1"},
			SalesmanID:       sql.NullString{String: "VALID1"},
			AndroidID:        sql.NullString{String: "1"},
			RegDate:          sql.NullTime{Time: timeNow},
			Email:            sql.NullString{String: "VALID1"},
			NoTelp:           sql.NullString{String: "VALID1"},
			SalesmanCategory: sql.NullString{String: "VALID1"},
			ProductValidFrom: sql.NullTime{Time: timeNow},
			ProductValidThru: sql.NullTime{Time: timeNow.AddDate(1, 0, 0)},
			UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:    sql.NullString{String: constanta.SystemClient},
			UpdatedAt:        sql.NullTime{Time: timeNow},
			CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:    sql.NullString{String: constanta.SystemClient},
			CreatedAt:        sql.NullTime{Time: timeNow},
		}

		for _, mappingModelPKCEClientMapping2 := range dataPKCEClientMapping {
			if mappingModel.ClientID.String == mappingModelPKCEClientMapping2.ParentClientID.String {
				modelUserRegistrationDetail.ClientID.String = mappingModelPKCEClientMapping2.ClientID.String
				modelUserRegistrationDetail.AuthUserID.Int64 = mappingModelPKCEClientMapping2.AuthUserID.Int64
			}
		}

		_, err = dao.UserRegistrationAdminDAO.InsertUserRegistrationAdmin(tx, repository.UserRegistrationAdminModel{
			UniqueID1:       sql.NullString{String: mappingModel.CompanyID.String},
			UniqueID2:       sql.NullString{String: mappingModel.BranchID.String},
			CompanyName:     sql.NullString{String: mappingModel.ClientAlias.String},
			BranchName:      sql.NullString{String: "VALID"},
			UserAdmin:       sql.NullString{String: "user01"},
			PasswordAdmin:   sql.NullString{String: "abc123"},
			ClientTypeID:    sql.NullInt64{Int64: mappingModel.ClientTypeID.Int64},
			ClientMappingID: sql.NullInt64{Int64: clientMappingIDTemp},
			UpdatedBy:       sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient:   sql.NullString{String: constanta.SystemClient},
			UpdatedAt:       sql.NullTime{Time: timeNow},
			CreatedBy:       sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient:   sql.NullString{String: constanta.SystemClient},
			CreatedAt:       sql.NullTime{Time: timeNow},
		})
		if err.Error != nil {
			return
		}

		_, err = dao.UserRegistrationDetailDAO.InsertUserRegistrationDetail(tx, modelUserRegistrationDetail, false)
		if err.Error != nil {
			return
		}
	}

	return
}
