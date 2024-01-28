package penambahan_user_admin

import (
	"database/sql"
	"fmt"
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
	contextModelNexMile, contextModelNexChiefMobile  applicationModel.ContextModel
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

	if err = insertClientCredentials(timeNow, tx); err.Error != nil {
		return
	}

	if err = masterDataID.CreateInitiateData(tx, timeNow); err.Error != nil {
		return
	}

	// insert Data Pendukung User Admin Nexmile
	if err = insertUserAdminData(tx, masterDataID, repository.ClientMappingModel{
		ClientID:         sql.NullString{String: "98381c991e6b409eb016cfaa365k4cad"},
		ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceTestingNexmileID},
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
	}, repository.CustomerInstallationModel{
		ParentCustomerID: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
		CustomerInstallationData: []repository.CustomerInstallationData{
			{
				CustomerID: sql.NullInt64{Int64: masterDataID.CustomerID},
				Installation: []repository.CustomerInstallationDetail{
					{
						UniqueID1:          sql.NullString{String: "123"},
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

	// insert Data Pendukung User Admin NexChief Mobile (6)
	if err = insertUserAdminData(tx, masterDataID, repository.ClientMappingModel{
		ClientID:         sql.NullString{String: "r3fb12faf6a348759ccffc500d609f31"},
		ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceTestingNexchiefMobileID},
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
	}, repository.CustomerInstallationModel{
		ParentCustomerID: sql.NullInt64{Int64: masterDataID.ParentCustomerID},
		CustomerInstallationData: []repository.CustomerInstallationData{
			{
				CustomerID: sql.NullInt64{Int64: masterDataID.CustomerID},
				Installation: []repository.CustomerInstallationDetail{
					{
						UniqueID1:          sql.NullString{String: "456"},
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

func insertUserAdminData( tx *sql.Tx, masterDataID test.MasterDataTesting, clientMappingData repository.ClientMappingModel, customerInstallationData repository.CustomerInstallationModel) (err errorModel.ErrorModel) {


	customerInstallationData.CustomerInstallationData[0].Installation[0].ProductID.Int64 = masterDataID.ProductNexmileID

	// Insert Customer Site
	if customerInstallationData.CustomerInstallationData[0].SiteID.Int64, err = dao.CustomerSiteDAO.InsertCustomerSite(tx, customerInstallationData, 0); err.Error != nil {
		return
	}

	clientMappingData.SiteID = customerInstallationData.CustomerInstallationData[0].SiteID

	// Insert Client Mapping
	clientMappingData.ID.Int64, err = dao.ClientMappingDAO.InsertClientMapping(tx, &clientMappingData)
	if err.Error != nil {
		return
	}

	// ---------------- Insert Customer Installation
	if customerInstallationData.CustomerInstallationData[0].Installation[0].InstallationID.Int64, err = dao.CustomerInstallationDAO.InsertCustomerInstallationForTesting(tx, customerInstallationData, 0, 0, 0, clientMappingData.ID.Int64); err.Error != nil {
		return
	}

	return
}
