package license_config

import (
	"database/sql"
	"errors"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/markbates/oncer"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service/ProductService"
	"nexsoft.co.id/nextrac2/test"
	"os"
	"testing"
	"time"
)

var InitiateTestVar initiateTest

type initiateTest struct {
	Tx              *sql.Tx
	ContextModel    *applicationModel.ContextModel
	masterDataID    test.MasterDataTesting
	installationID  int64
	licenseConfigID int64
	timeNow         time.Time
}

type modelSetLicenseData struct {
	clientCredentialModel    repository.ClientCredentialModel
	clientMappingData        repository.ClientMappingModel
	customerInstallationData repository.CustomerInstallationModel
}

func TestMain(m *testing.M) {
	var err error

	defer func() {
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

	//------------------ Set The Database
	errs := setDataOnDatabaseLocal(*InitiateTestVar.ContextModel)
	if errs.Error != nil {
		return errors.New("failure in database")
	}

	return nil
}

func setDataOnDatabaseLocal(contextModel applicationModel.ContextModel) (err errorModel.ErrorModel) {
	test.SetDataWithTransactionalDB(contextModel, doSetDataInitiateOnDatabaseLocal)
	return test.SetDataWithTransactionalDB(contextModel, doSetDataOnDatabaseLocal)
}

func doSetDataInitiateOnDatabaseLocal(tx *sql.Tx, _ applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	clientTypeModelInit := []repository.ClientTypeModel{
		{
			ID:            sql.NullInt64{Int64: 4},
			ClientType:    sql.NullString{String: "nexChief"},
			Description:   sql.NullString{String: "Client Type Nexchief"},
			UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient: sql.NullString{String: constanta.SystemClient},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient: sql.NullString{String: constanta.SystemClient},
			CreatedAt:     sql.NullTime{Time: timeNow},
		},
		{
			ID:            sql.NullInt64{Int64: 5},
			ClientType:    sql.NullString{String: "nexMile"},
			Description:   sql.NullString{String: "Client Type Nexmile"},
			UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient: sql.NullString{String: constanta.SystemClient},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient: sql.NullString{String: constanta.SystemClient},
			CreatedAt:     sql.NullTime{Time: timeNow},
		},
		{
			ID:            sql.NullInt64{Int64: 6},
			ClientType:    sql.NullString{String: "nexChief mobile"},
			Description:   sql.NullString{String: "Client Type Nexchief Mobile"},
			UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient: sql.NullString{String: constanta.SystemClient},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient: sql.NullString{String: constanta.SystemClient},
			CreatedAt:     sql.NullTime{Time: timeNow},
		},
	}

	for _, itemClientType := range clientTypeModelInit {
		_, err = dao.ClientTypeDAO.InsertClientType(tx, itemClientType)
		if err.Error != nil {
			return
		}
	}

	if err = InitiateTestVar.masterDataID.CreateInitiateData(tx, timeNow); err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func doSetDataOnDatabaseLocal(tx *sql.Tx, _ applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	masterDataID := InitiateTestVar.masterDataID

	_, _, err = setLicenseData(tx, timeNow, modelSetLicenseData{
		clientCredentialModel: repository.ClientCredentialModel{
			ClientID:      sql.NullString{String: "c833660f2d254027b527da2320d39f14"},
			ClientSecret:  sql.NullString{String: "5663f56c6d1f4a91a5845c14adf6692e"},
			SignatureKey:  sql.NullString{String: "f76629da229e40cba31872cbaee2259e"},
			UpdatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			UpdatedClient: sql.NullString{String: constanta.SystemClient},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			CreatedBy:     sql.NullInt64{Int64: constanta.SystemID},
			CreatedClient: sql.NullString{String: constanta.SystemClient},
			CreatedAt:     sql.NullTime{Time: timeNow},
		},
		clientMappingData: repository.ClientMappingModel{
			ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceND6ID},
			CustomerID:       sql.NullInt64{Int64: masterDataID.CustomerID},
			CompanyID:        sql.NullString{String: "NS6084010002596"},
			BranchID:         sql.NullString{String: "1596128276342"},
			ClientAlias:      sql.NullString{String: "ND6 - PT. Eka Artha Buanas 151"},
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
							ProductID:          sql.NullInt64{Int64: masterDataID.ProductND6ID},
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
	})
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func setLicenseData(tx *sql.Tx, timeNow time.Time, data modelSetLicenseData) (userLicenseID int64, licenseConfigData repository.LicenseConfigModel, err errorModel.ErrorModel) {
	var productData out.ViewProduct

	// ---------------- Insert Client Credential
	if _, err = dao.ClientCredentialDAO.InsertClientCredential(tx, &data.clientCredentialModel); err.Error != nil {
		return
	}

	// ---------------- Insert Customer Site
	if data.clientMappingData.SiteID.Int64, err = dao.CustomerSiteDAO.InsertCustomerSite(tx, data.customerInstallationData, 0); err.Error != nil {
		return
	}

	// ---------------- Insert Client Mapping
	data.clientMappingData.ClientID = data.clientCredentialModel.ClientID
	data.clientMappingData.ID.Int64, err = dao.ClientMappingDAO.InsertClientMapping(tx, &data.clientMappingData)
	if err.Error != nil {
		return
	}

	// ---------------- Insert Customer Installation
	var queue int64
	data.customerInstallationData.CustomerInstallationData[0].SiteID.Int64 = data.clientMappingData.SiteID.Int64

	for idx := range data.customerInstallationData.CustomerInstallationData[0].Installation {
		// ---------------- Get Product Data
		productData, err = ProductService.ProductService.DoViewProduct(in.ProductRequest{
			ID: data.customerInstallationData.CustomerInstallationData[0].Installation[idx].ProductID.Int64,
		}, &applicationModel.ContextModel{})

		// ---------------- Insert Customer Installation
		data.customerInstallationData.CustomerInstallationData[0].Installation[idx].UniqueID1 = data.clientMappingData.CompanyID
		data.customerInstallationData.CustomerInstallationData[0].Installation[idx].UniqueID2 = data.clientMappingData.BranchID

		for i := 0; i < 2; i++ {
			if data.customerInstallationData.CustomerInstallationData[0].Installation[idx].InstallationID.Int64, err = dao.CustomerInstallationDAO.InsertCustomerInstallationForTesting(tx, data.customerInstallationData, 0, idx, queue, data.clientMappingData.ID.Int64); err.Error != nil {
				return
			}

			oncer.Do("First Init For First Data", func() {
				InitiateTestVar.installationID = data.customerInstallationData.CustomerInstallationData[0].Installation[idx].InstallationID.Int64
			})

			if i == 1 {
				if userLicenseID, licenseConfigData, err = insertLicense(tx, timeNow, data, productData, idx); err.Error != nil {
					return
				}
				InitiateTestVar.timeNow = timeNow
			}
			queue++
		}
	}

	return
}

func insertLicense(tx *sql.Tx, timeNow time.Time, data modelSetLicenseData, productData out.ViewProduct, idx int) (userLicenseID int64, licenseConfigData repository.LicenseConfigModel, err errorModel.ErrorModel) {

	// ------- Insert License Config
	licenseConfigData = createLicenseConfigModel(timeNow, data, productData, idx)
	InitiateTestVar.licenseConfigID, err = dao.LicenseConfigDAO.InsertLicenseConfigForTesting(tx, licenseConfigData)
	if err.Error != nil {
		return
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
