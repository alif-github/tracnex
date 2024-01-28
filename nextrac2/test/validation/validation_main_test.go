package validation

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
	"nexsoft.co.id/nextrac2/service/ProductService"
	"nexsoft.co.id/nextrac2/test"
	util2 "nexsoft.co.id/nextrac2/util"
	"os"
	"testing"
	"time"
)

func TestMain(main *testing.M) {
	os.Exit(testMain(main))
}

var tx *sql.Tx
var (
	contextModelND, contextModelNexchief, contextModelNexMile, contextModelNexChiefMobile applicationModel.ContextModel
)

func testMain(main *testing.M) int {
	var err errorModel.ErrorModel
	var errS error

	fmt.Println("Start Testing Activation")

	// Set Configuration
	test.InitAllConfiguration()

	errS = test.Truncate(serverconfig.ServerAttribute.DBConnection)
	if errS != nil {
		return 1
	}

	// Set Database
	if err = test.SetDataWithTransactionalDB(applicationModel.ContextModel{}, setDatabase); err.Error != nil {
		fmt.Println(err)
		fmt.Println(util2.GenerateI18NErrorMessage(err, constanta.DefaultApplicationsLanguage))
		return 1
	}

	// Open Transactional
	if tx, errS = test.OpenTransactional(serverconfig.ServerAttribute.DBConnection); errS != nil {
		return 1
	}

	defer func() {
		if err.Error != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	return main.Run()
}

func setDatabase(tx *sql.Tx, contextModel applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {

	var (
		licenseVariantID, licenseTypeID, productGroupID, parentCustomerID, customerID, moduleID, componentID int64
	)

	if err = insertClientCredentials(timeNow, tx); err.Error != nil {
		return
	}

	if err,
		licenseVariantID,
		licenseTypeID,
		productGroupID,
		parentCustomerID,
		customerID,
		moduleID,
		componentID = insertMasterData(timeNow, tx); err.Error != nil {
		return
	}

	// Insert License ND6 (ID Client Type 1)
	if err = insertLicenseData(timeNow, tx, repository.ProductModel{
		ProductID:          sql.NullString{String: "VALID1"},
		ProductName:        sql.NullString{String: "VALID1"},
		ProductDescription: sql.NullString{String: "VALID1"},
		ProductGroupID:     sql.NullInt64{Int64: productGroupID},
		ClientTypeID:       sql.NullInt64{Int64: constanta.ResourceND6ID},
		IsLicense:          sql.NullBool{Bool: false},
		LicenseVariantID:   sql.NullInt64{Int64: licenseVariantID},
		LicenseTypeID:      sql.NullInt64{Int64: licenseTypeID},
		DeploymentMethod:   sql.NullString{String: "O"},
		NoOfUser:           sql.NullInt64{Int64: 1},
		IsUserConcurrent:   sql.NullBool{Bool: false},
		MaxOfflineDays:     sql.NullInt64{Int64: 1},
		Module1:            sql.NullInt64{Int64: moduleID},
		UpdatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:      sql.NullString{String: constanta.SystemClient},
		UpdatedAt:          sql.NullTime{Time: timeNow},
		CreatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:      sql.NullString{String: constanta.SystemClient},
		CreatedAt:          sql.NullTime{Time: timeNow},
		ProductComponentModel: []repository.ProductComponentModel{
			repository.ProductComponentModel{
				ComponentID:    sql.NullInt64{Int64: componentID},
				ComponentValue: sql.NullString{String: "VALID"},
			},
		},
	}, repository.ClientMappingModel{
		ClientID:         sql.NullString{String: "08181c991e6b409eb016cfaa365b439d"},
		ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceND6ID},
		CustomerID:       sql.NullInt64{Int64: customerID},
		CompanyID:        sql.NullString{String: "NS6024050001031"},
		BranchID:         sql.NullString{String: "1468381449586"},
		ClientAlias:      sql.NullString{String: "ND6 - PT. Eka Artha Buanas 111"},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
		ParentCustomerID: sql.NullInt64{Int64: parentCustomerID},
	}, repository.CustomerInstallationModel{
		ParentCustomerID: sql.NullInt64{Int64: parentCustomerID},
		CustomerInstallationData: []repository.CustomerInstallationData{
			{
				CustomerID: sql.NullInt64{Int64: customerID},
				Installation: []repository.CustomerInstallationDetail{
					{
						UniqueID1:          sql.NullString{String: "NS6024050001031"},
						UniqueID2:          sql.NullString{String: "1468381449586"},
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
	}, repository.ProductLicenseModel{
		LicenseConfigId:        sql.NullInt64{},
		ProductKey:             sql.NullString{String: "1"},
		ProductEncrypt:         sql.NullString{String: "1"},
		ProductSignature:       sql.NullString{String: "1"},
		ClientId:               sql.NullString{String: "08181c991e6b409eb016cfaa365b439d"},
		ClientSecret:           sql.NullString{String: "6bf54c4237964a3eb9637da1fb2c622a"},
		HWID:                   sql.NullString{String: "1"},
		ActivationDate:         sql.NullTime{Time: timeNow},
		LicenseStatus:          sql.NullInt32{Int32: 1},
		TerminationDescription: sql.NullString{},
		CreatedBy:              sql.NullInt64{Int64: constanta.SystemID},
		CreatedAt:              sql.NullTime{Time: timeNow},
		CreatedClient:          sql.NullString{String: constanta.SystemClient},
		UpdatedBy:              sql.NullInt64{Int64: constanta.SystemID},
		UpdatedAt:              sql.NullTime{Time: timeNow},
		UpdatedClient:          sql.NullString{String: constanta.SystemClient},
	}); err.Error != nil {
		return
	}

	contextModelND = applicationModel.ContextModel{
		AuthAccessTokenModel: model.AuthAccessTokenModel{
			RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
				ResourceUserID: 12,
			},
			ClientID: "08181c991e6b409eb016cfaa365b439d",
			Locale:   constanta.IndonesianLanguage,
		},
	}

	// Insert License NexChief (4)
	if err = insertLicenseData(timeNow, tx, repository.ProductModel{
		ProductID:          sql.NullString{String: "VALID2"},
		ProductName:        sql.NullString{String: "VALID2"},
		ProductDescription: sql.NullString{String: "VALID2"},
		ProductGroupID:     sql.NullInt64{Int64: productGroupID},
		ClientTypeID:       sql.NullInt64{Int64: constanta.ResourceNexChiefID},
		IsLicense:          sql.NullBool{Bool: true},
		LicenseVariantID:   sql.NullInt64{Int64: licenseVariantID},
		LicenseTypeID:      sql.NullInt64{Int64: licenseTypeID},
		DeploymentMethod:   sql.NullString{String: "O"},
		NoOfUser:           sql.NullInt64{Int64: 1},
		IsUserConcurrent:   sql.NullBool{Bool: true},
		MaxOfflineDays:     sql.NullInt64{Int64: 1},
		Module1:            sql.NullInt64{Int64: moduleID},
		UpdatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:      sql.NullString{String: constanta.SystemClient},
		UpdatedAt:          sql.NullTime{Time: timeNow},
		CreatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:      sql.NullString{String: constanta.SystemClient},
		CreatedAt:          sql.NullTime{Time: timeNow},
		ProductComponentModel: []repository.ProductComponentModel{
			repository.ProductComponentModel{
				ComponentID:    sql.NullInt64{Int64: componentID},
				ComponentValue: sql.NullString{String: "Valid"},
			},
		},
	}, repository.ClientMappingModel{
		ClientID:         sql.NullString{String: "1a2b12faf6a345759ccffc500d609b52"},
		ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceNexChiefID},
		CustomerID:       sql.NullInt64{Int64: customerID},
		CompanyID:        sql.NullString{String: "NDI"},
		ClientAlias:      sql.NullString{String: "PT. Manohara Asri"},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
		ParentCustomerID: sql.NullInt64{Int64: parentCustomerID},
	}, repository.CustomerInstallationModel{
		ParentCustomerID: sql.NullInt64{Int64: parentCustomerID},
		CustomerInstallationData: []repository.CustomerInstallationData{
			{
				CustomerID: sql.NullInt64{Int64: customerID},
				Installation: []repository.CustomerInstallationDetail{
					{
						UniqueID1:          sql.NullString{String: "NDI"},
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
	}, repository.ProductLicenseModel{
		LicenseConfigId:        sql.NullInt64{},
		ProductKey:             sql.NullString{String: "2"},
		ProductEncrypt:         sql.NullString{String: "2"},
		ProductSignature:       sql.NullString{String: "2"},
		ClientId:               sql.NullString{String: "1a2b12faf6a345759ccffc500d609b52"},
		ClientSecret:           sql.NullString{String: "47d40eb8063d4513beda8357948a1040"},
		HWID:                   sql.NullString{String: "4"},
		ActivationDate:         sql.NullTime{Time: timeNow},
		LicenseStatus:          sql.NullInt32{Int32: 1},
		TerminationDescription: sql.NullString{},
		CreatedBy:              sql.NullInt64{Int64: constanta.SystemID},
		CreatedAt:              sql.NullTime{Time: timeNow},
		CreatedClient:          sql.NullString{String: constanta.SystemClient},
		UpdatedBy:              sql.NullInt64{Int64: constanta.SystemID},
		UpdatedAt:              sql.NullTime{Time: timeNow},
		UpdatedClient:          sql.NullString{String: constanta.SystemClient},
	}); err.Error != nil {
		return
	}

	contextModelNexchief = applicationModel.ContextModel{
		AuthAccessTokenModel: model.AuthAccessTokenModel{
			RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
				ResourceUserID: 12,
			},
			ClientID: "1a2b12faf6a345759ccffc500d609b52",
			Locale:   constanta.IndonesianLanguage,
		},
	}

	// Insert License NexMile (5)
	if err = insertLicenseData(timeNow, tx, repository.ProductModel{
		ProductID:          sql.NullString{String: "VALID3"},
		ProductName:        sql.NullString{String: "VALID3"},
		ProductDescription: sql.NullString{String: "VALID3"},
		ProductGroupID:     sql.NullInt64{Int64: productGroupID},
		ClientTypeID:       sql.NullInt64{Int64: constanta.ResourceTestingNexmileID},
		IsLicense:          sql.NullBool{Bool: true},
		LicenseVariantID:   sql.NullInt64{Int64: licenseVariantID},
		LicenseTypeID:      sql.NullInt64{Int64: licenseTypeID},
		DeploymentMethod:   sql.NullString{String: "O"},
		NoOfUser:           sql.NullInt64{Int64: 1},
		IsUserConcurrent:   sql.NullBool{Bool: true},
		MaxOfflineDays:     sql.NullInt64{Int64: 1},
		Module1:            sql.NullInt64{Int64: moduleID},
		UpdatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:      sql.NullString{String: constanta.SystemClient},
		UpdatedAt:          sql.NullTime{Time: timeNow},
		CreatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:      sql.NullString{String: constanta.SystemClient},
		CreatedAt:          sql.NullTime{Time: timeNow},
		ProductComponentModel: []repository.ProductComponentModel{
			repository.ProductComponentModel{
				ComponentID:    sql.NullInt64{Int64: componentID},
				ComponentValue: sql.NullString{String: "Valid"},
			},
		},
	}, repository.ClientMappingModel{
		ClientID:         sql.NullString{String: "98381c991e6b409eb016cfaa365k4cad"},
		ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceTestingNexmileID},
		CustomerID:       sql.NullInt64{Int64: customerID},
		CompanyID:        sql.NullString{String: "123"},
		ClientAlias:      sql.NullString{String: "PT. Makmur Sejahtera"},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
		ParentCustomerID: sql.NullInt64{Int64: parentCustomerID},
	}, repository.CustomerInstallationModel{
		ParentCustomerID: sql.NullInt64{Int64: parentCustomerID},
		CustomerInstallationData: []repository.CustomerInstallationData{
			{
				CustomerID: sql.NullInt64{Int64: customerID},
				Installation: []repository.CustomerInstallationDetail{
					{
						UniqueID1:          sql.NullString{String: "123"},
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

	// Insert License NexChief Mobile (6)
	if err = insertLicenseData(timeNow, tx, repository.ProductModel{
		ProductID:          sql.NullString{String: "VALID4"},
		ProductName:        sql.NullString{String: "VALID4"},
		ProductDescription: sql.NullString{String: "VALID4"},
		ProductGroupID:     sql.NullInt64{Int64: productGroupID},
		ClientTypeID:       sql.NullInt64{Int64: constanta.ResourceTestingNexchiefMobileID},
		IsLicense:          sql.NullBool{Bool: true},
		LicenseVariantID:   sql.NullInt64{Int64: licenseVariantID},
		LicenseTypeID:      sql.NullInt64{Int64: licenseTypeID},
		DeploymentMethod:   sql.NullString{String: "O"},
		NoOfUser:           sql.NullInt64{Int64: 1},
		IsUserConcurrent:   sql.NullBool{Bool: true},
		MaxOfflineDays:     sql.NullInt64{Int64: 1},
		Module1:            sql.NullInt64{Int64: moduleID},
		UpdatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:      sql.NullString{String: constanta.SystemClient},
		UpdatedAt:          sql.NullTime{Time: timeNow},
		CreatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:      sql.NullString{String: constanta.SystemClient},
		CreatedAt:          sql.NullTime{Time: timeNow},
		ProductComponentModel: []repository.ProductComponentModel{
			repository.ProductComponentModel{
				ComponentID:    sql.NullInt64{Int64: componentID},
				ComponentValue: sql.NullString{String: "Valid"},
			},
		},
	}, repository.ClientMappingModel{
		ClientID:         sql.NullString{String: "r3fb12faf6a348759ccffc500d609f31"},
		ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceTestingNexchiefMobileID},
		CustomerID:       sql.NullInt64{Int64: customerID},
		CompanyID:        sql.NullString{String: "456"},
		ClientAlias:      sql.NullString{String: "PT. Maju Jaya"},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
		ParentCustomerID: sql.NullInt64{Int64: parentCustomerID},
	}, repository.CustomerInstallationModel{
		ParentCustomerID: sql.NullInt64{Int64: parentCustomerID},
		CustomerInstallationData: []repository.CustomerInstallationData{
			{
				CustomerID: sql.NullInt64{Int64: customerID},
				Installation: []repository.CustomerInstallationDetail{
					{
						UniqueID1:          sql.NullString{String: "456"},
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

	err = errorModel.GenerateNonErrorModel()
	return
}

func insertClientCredentials(timeNow time.Time, tx *sql.Tx) (err errorModel.ErrorModel) {

	if _, err = dao.ClientCredentialDAO.InsertClientCredential(tx, &repository.ClientCredentialModel{
		ClientID:      sql.NullString{String: "08181c991e6b409eb016cfaa365b439d"},
		ClientSecret:  sql.NullString{String: "6bf54c4237964a3eb9637da1fb2c622a"},
		SignatureKey:  sql.NullString{String: "280d9968c4154d698362087a91a80e1a"},
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
		ClientID:      sql.NullString{String: "1a2b12faf6a345759ccffc500d609b52"},
		ClientSecret:  sql.NullString{String: "47d40eb8063d4513beda8357948a1040"},
		SignatureKey:  sql.NullString{String: "bb0734e85ba44b529611fd22668b6bad"},
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

func insertMasterData(timeNow time.Time, tx *sql.Tx) (err errorModel.ErrorModel, licenseVariantID, licenseTypeID, productGroupID, parentCustomerID, customerID, moduleID, componentID int64) {

	var (
		customerCategoryID, customerGroupID, salesmanID int64
	)

	// Insert Data Master License Variant
	if licenseVariantID, err = dao.LicenseVariantDAO.InsertLicenseVariant(tx, repository.LicenseVariantModel{
		LicenseVariantName: sql.NullString{String: "VALID"},
		UpdatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:      sql.NullString{String: constanta.SystemClient},
		UpdatedAt:          sql.NullTime{Time: timeNow},
		CreatedBy:          sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:      sql.NullString{String: constanta.SystemClient},
		CreatedAt:          sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// Insert Data Master License Type
	if licenseTypeID, err = dao.LicenseTypeDAO.InsertLicenseType(tx, repository.LicenseTypeModel{
		LicenseTypeName: sql.NullString{String: "VALID"},
		UpdatedBy:       sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:   sql.NullString{String: constanta.SystemClient},
		UpdatedAt:       sql.NullTime{Time: timeNow},
		CreatedBy:       sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:   sql.NullString{String: constanta.SystemClient},
		CreatedAt:       sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	// Insert Data Master Customer Category
	if customerCategoryID, err = dao.CustomerCategoryDAO.InsertCustomerCategory(tx, repository.CustomerCategoryModel{
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

	// Insert Data Master Customer Group
	if customerGroupID, err = dao.CustomerGroupDAO.InsertCustomerGroup(tx, repository.CustomerGroupModel{
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

	// Insert Data Master Salesman
	if salesmanID, err = dao.SalesmanDAO.InsertSalesman(tx, repository.SalesmanModel{
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

	// Insert Data Master Parent Customer
	if parentCustomerID, err = dao.CustomerDAO.InsertCustomer(tx, repository.CustomerModel{
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
		SalesmanID:              sql.NullInt64{Int64: salesmanID},
		DistributorOF:           sql.NullString{String: "VALID"},
		CustomerGroupID:         sql.NullInt64{Int64: customerGroupID},
		CustomerCategoryID:      sql.NullInt64{Int64: customerCategoryID},
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

	// Insert Data Master Customer
	if customerID, err = dao.CustomerDAO.InsertCustomer(tx, repository.CustomerModel{
		IsPrincipal:             sql.NullBool{Bool: true},
		IsParent:                sql.NullBool{Bool: false},
		ParentCustomerID:        sql.NullInt64{Int64: parentCustomerID},
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
		SalesmanID:              sql.NullInt64{Int64: salesmanID},
		DistributorOF:           sql.NullString{String: "VALID"},
		CustomerGroupID:         sql.NullInt64{Int64: customerGroupID},
		CustomerCategoryID:      sql.NullInt64{Int64: customerCategoryID},
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

	// Insert Data Master Product Group
	if productGroupID, err = dao.ProductGroupDAO.InsertProductGroup(tx, repository.ProductGroupModel{
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

	// Insert Data Master Module
	if moduleID, err = dao.ModuleDAO.InsertModule(tx, repository.ModuleModel{
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

	// Insert Data Master Component
	if componentID, err = dao.ComponentDAO.InsertComponent(tx, repository.ComponentModel{
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

	return
}

func insertLicenseData(timeNow time.Time, tx *sql.Tx, productModel repository.ProductModel, clientMappingData repository.ClientMappingModel, customerInstallationData repository.CustomerInstallationModel, productLicenseData repository.ProductLicenseModel) (err errorModel.ErrorModel) {

	// Insert Product
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
	customerInstallationData.CustomerInstallationData[0].Installation[0].ProductID = productModel.ID

	// ---------------- Insert Customer Site
	if customerInstallationData.CustomerInstallationData[0].SiteID.Int64, err = dao.CustomerSiteDAO.InsertCustomerSite(tx, customerInstallationData, 0); err.Error != nil {
		return
	}
	clientMappingData.SiteID = customerInstallationData.CustomerInstallationData[0].SiteID

	// ---------------- Insert Client Mapping
	clientMappingData.ID.Int64, err = dao.ClientMappingDAO.InsertClientMapping(tx, &clientMappingData)
	if err.Error != nil {
		return
	}

	// ---------------- Insert Customer Installation
	if customerInstallationData.CustomerInstallationData[0].Installation[0].InstallationID.Int64, err = dao.CustomerInstallationDAO.InsertCustomerInstallationForTesting(tx, customerInstallationData, 0, 0, 0, clientMappingData.ID.Int64); err.Error != nil {
		return
	}

	licenseConfigData := repository.LicenseConfigModel{
		InstallationID:   customerInstallationData.CustomerInstallationData[0].Installation[0].InstallationID,
		NoOfUser:         sql.NullInt64{Int64: 25},
		ProductValidFrom: sql.NullTime{Time: date.Date{timeNow.Add(-time.Hour * 24)}.ToTime()},
		ProductValidThru: sql.NullTime{Time: date.Date{timeNow.Add((time.Hour * 24) * 100)}.ToTime()},
		UpdatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		UpdatedClient:    sql.NullString{String: constanta.SystemClient},
		UpdatedAt:        sql.NullTime{Time: timeNow},
		CreatedBy:        sql.NullInt64{Int64: constanta.SystemID},
		CreatedClient:    sql.NullString{String: constanta.SystemClient},
		CreatedAt:        sql.NullTime{Time: timeNow},
		ParentCustomerID: customerInstallationData.ParentCustomerID,
		CustomerID:       customerInstallationData.CustomerInstallationData[0].CustomerID,
		SiteID:           customerInstallationData.CustomerInstallationData[0].SiteID,
		ClientID:         clientMappingData.ClientID,
		ProductID:        productModel.ID,
		ClientTypeID:     productModel.ClientTypeID,
		LicenseVariantID: productModel.LicenseVariantID,
		LicenseTypeID:    productModel.LicenseTypeID,
		DeploymentMethod: productModel.DeploymentMethod,
		MaxOfflineDays:   productModel.MaxOfflineDays,
		UniqueID1:        customerInstallationData.CustomerInstallationData[0].Installation[0].UniqueID1,
		UniqueID2:        customerInstallationData.CustomerInstallationData[0].Installation[0].UniqueID2,
		ModuleID1:        productModel.Module1,
		AllowActivation:  sql.NullString{String: "Y"},
		Component:        productModel.ProductComponentModel,
	}

	if productModel.IsUserConcurrent.Bool {
		licenseConfigData.IsUserConcurrent.String = "Y"
	} else {
		licenseConfigData.IsUserConcurrent.String = "N"
	}

	productLicenseData.LicenseConfigId.Int64, err = dao.LicenseConfigDAO.InsertLicenseConfigForTesting(tx, licenseConfigData)
	if err.Error != nil {
		return
	}

	// Product License
	_, err = dao.ProductLicenseDAO.InsertProductLicense(tx, productLicenseData)

	return
}
