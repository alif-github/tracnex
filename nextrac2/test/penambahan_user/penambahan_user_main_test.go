package penambahan_user

import (
	"database/sql"
	"fmt"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/jarcoal/httpmock"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
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
var dataMockUserNexmile []test.MockBodyAndResponseData
var (
	userRegistrationDetailIDToBeActive, userRegistrationDetailIDToBeError int64
)

//data repeat
var resourceIDData = "nc auth nexcare master"
var scopeData = "read write"


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
	httpmock.Activate()
	defer func() {
		httpmock.DeactivateAndReset()
	}()
	os.Exit(testMain(main))
}

func testMain(main *testing.M) int {
	var err errorModel.ErrorModel
	var errS error
	fmt.Println("Start Testing Penambahan User")

	defer func() {
		test.RollBackSchema(serverconfig.ServerAttribute.DBConnection)
	}()

	// Set Configuration
	test.InitAllConfiguration()

	//Truncate function
	if errS = test.SetClientType(serverconfig.ServerAttribute.DBConnection); errS != nil {
		fmt.Println("Gagal Buat Client Type")
		return 1
	}

	//Set Master Database
	if err = test.SetDataWithTransactionalDB(applicationModel.ContextModel{}, setMasterData); err.Error != nil {
		fmt.Println(err)
		fmt.Println(util2.GenerateI18NErrorMessage(err, constanta.DefaultApplicationsLanguage))
		return 1
	}

	// Set Data User Nexmile
	if err = setDataPenambahanUser(); err.Error != nil {
		return 1
	}

	// set mock user nexmile
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	checkClientUserUrl := authenticationServer.Host + authenticationServer.PathRedirect.InternalClient.CheckClientUser

	test.SetHttpMockResponseWithRequest(http.MethodPost, checkClientUserUrl, 200, dataMockUserNexmile, authentication_response.CheckClientOrUserResponse{
		Nexsoft: authentication_response.CheckClientOrUserBodyResponse{
			Payload: authentication_response.CheckClientOrUserPayload{
				Data: authentication_response.CheckClientOrUserData{
					Content: authentication_response.CheckClientOrUserContent{
						IsExist:               false,
						AdditionalInformation: authentication_response.AdditionalInformationContent{},
					},
				},
			},
		},
	})

	return main.Run()
}

func setMasterData(tx *sql.Tx, contextModel applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	if err = masterDataID.CreateInitiateData(tx, timeNow); err.Error != nil {
		return
	}
	return
}

func setDataPenambahanUser() (err errorModel.ErrorModel) {
	//Set Database Nexmile
	if err = test.SetDataWithTransactionalDB(applicationModel.ContextModel{}, setDataUserNexmile); err.Error != nil {
		fmt.Println(err)
		fmt.Println(util2.GenerateI18NErrorMessage(err, constanta.DefaultApplicationsLanguage))
		return
	}

	//Set Database Other Named User
	if err = test.SetDataWithTransactionalDB(applicationModel.ContextModel{}, setDataUserOtherNamedUser); err.Error != nil {
		fmt.Println(err)
		fmt.Println(util2.GenerateI18NErrorMessage(err, constanta.DefaultApplicationsLanguage))
		return
	}

	return
}

func setDataUserOtherNamedUser(tx *sql.Tx, contextModel applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	//var userLicenseID int64
	//var licenseConfigData repository.LicenseConfigModel

	// Create license with UserLicense
	if _, _, err = setLicenseData(tx, timeNow, modelSetLicenseData{
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
			ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceNexmileID},
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

	// insert mock regist user nexmile
	dataMockUserNexmile = append(dataMockUserNexmile, test.MockBodyAndResponseData{
		Body: authentication_request.CheckClientOrUser{ClientID: "59005771231e4caaa87c1636e3de186a"},
		Response: authentication_response.CheckClientOrUserResponse{
			Nexsoft: authentication_response.CheckClientOrUserBodyResponse{
				Payload: authentication_response.CheckClientOrUserPayload{
					Data: authentication_response.CheckClientOrUserData{
						Content: authentication_response.CheckClientOrUserContent{
							IsExist: true,
							AdditionalInformation: authentication_response.AdditionalInformationContent{
								AliasName:    "ValidAlias",
								ClientID:     "59005771231e4caaa87c1636e3de186a",
								ClientSecret: "7dg14c4237964a3eb9637da1fb2c897z",
								UserID:       103,
								Username:     "UserName5",
								SignatureKey: "ab3dd1354df9484aa38d81ffb5f37692",
								GrantTypes:   "code_pkce",
								ResourceID:   resourceIDData,
								UserStatus:   1,
								Scope:        scopeData,
								Locale:       constanta.DefaultApplicationsLanguage,
							},
						},
					},
				},
			},
		},
	})

	dataMockUserNexmile = append(dataMockUserNexmile, test.MockBodyAndResponseData{
		Body: authentication_request.CheckClientOrUser{ClientID: "62cda0a7242c497bbc502e4b33e87abc"},
		Response: authentication_response.CheckClientOrUserResponse{
			Nexsoft: authentication_response.CheckClientOrUserBodyResponse{
				Payload: authentication_response.CheckClientOrUserPayload{
					Data: authentication_response.CheckClientOrUserData{
						Content: authentication_response.CheckClientOrUserContent{
							IsExist: true,
							AdditionalInformation: authentication_response.AdditionalInformationContent{
								AliasName:    "ValidAlias",
								ClientID:     "62cda0a7242c497bbc502e4b33e87abc",
								ClientSecret: "c3ee8badff034b4aa9c90612441f2d72",
								UserID:       104,
								Username:     "UserName5",
								SignatureKey: "a2de8badff034b4aa9c90612441f2e71",
								GrantTypes:   "code_pkce",
								ResourceID:   resourceIDData,
								UserStatus:   1,
								Scope:        scopeData,
								Locale:       constanta.DefaultApplicationsLanguage,
							},
						},
					},
				},
			},
		},
	})

	return
}

func setDataUserNexmile(tx *sql.Tx, contextModel applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	var userLicenseID int64
	var licenseConfigData repository.LicenseConfigModel

	//Create license without userLicense
	if _, _, err = setLicenseData(tx, timeNow, modelSetLicenseData{
		clientCredentialModel: repository.ClientCredentialModel{
			ClientID:      sql.NullString{String: "1743b12b50074b0cb993d2a43badf36a"},
			ClientSecret:  sql.NullString{String: "3d065880853646849015b79c4f968f73"},
			SignatureKey:  sql.NullString{String: "83d8f10e24964b9da7b118673c8a847a"},
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
			CompanyID:        sql.NullString{String: "NS6024050001135"},
			BranchID:         sql.NullString{String: "1468381448675"},
			ClientAlias:      sql.NullString{String: "ND6 - PT. Eka Artha Buahaha 111"},
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
							Remark:             sql.NullString{String: "VALID NexDist"},
							InstallationDate:   sql.NullTime{Time: timeNow},
							InstallationStatus: sql.NullString{String: "A"},
							NoOfUser:           sql.NullInt64{Int64: 5},
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

	// Create license with UserLicense
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
			ClientTypeID:     sql.NullInt64{Int64: constanta.ResourceND6ID},
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

	// insert User Regis Detail
	if userRegistrationDetailIDToBeActive, err = insertUserRegistDetail(userLicenseID, tx, timeNow, modelInsertUserLicense{
		licenseConfigData:     licenseConfigData,
		inputStruct:           in.RegisterNamedUserRequest{
			ParentClientID:   licenseConfigData.ClientID.String,
			ClientID:         "TestClientID",
			ClientTypeID:     2,
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
		isRegisStatus:         true,
	}); err.Error != nil {
		return
	}

	// insert User Regis Detail
	if userRegistrationDetailIDToBeError, err = insertUserRegistDetail(userLicenseID, tx, timeNow, modelInsertUserLicense{
		licenseConfigData:     licenseConfigData,
		inputStruct:           in.RegisterNamedUserRequest{
			ParentClientID:   licenseConfigData.ClientID.String,
			ClientID:         "TestClientID1",
			ClientTypeID:     2,
			AuthUserID:       101,
			Firstname:        "Pertama1",
			Lastname:         "Terakhir1",
			Username:         "UsernameValid1",
			UserID:           "101",
			Password:         "PassValid1",
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
		isRegisStatus:         true,
	}); err.Error != nil {
		return
	}

	// insert mock regist user nexmile
	dataMockUserNexmile = append(dataMockUserNexmile, test.MockBodyAndResponseData{
		Body: authentication_request.CheckClientOrUser{ClientID: "08181c920e6b409eb016cfaa365b439d"},
		Response: authentication_response.CheckClientOrUserResponse{
			Nexsoft: authentication_response.CheckClientOrUserBodyResponse{
				Payload: authentication_response.CheckClientOrUserPayload{
					Data: authentication_response.CheckClientOrUserData{
						Content: authentication_response.CheckClientOrUserContent{
							IsExist: true,
							AdditionalInformation: authentication_response.AdditionalInformationContent{
								AliasName:    "ValidAlias",
								ClientID:     "08181c920e6b409eb016cfaa365b439d",
								ClientSecret: "6bf54c4237964a3eb9637da1fb2c622a",
								UserID:       102,
								Username:     "UsernameValid2",
								SignatureKey: "280d9968c4154d698362087a91a80e1a",
								GrantTypes:   "code_pkce",
								ResourceID:   resourceIDData,
								UserStatus:   1,
								Scope:        scopeData,
								Locale:       constanta.DefaultApplicationsLanguage,
							},
						},
					},
				},
			},
		},
	})

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

	dataMockUserNexmile = append(dataMockUserNexmile, test.MockBodyAndResponseData{
		Body: authentication_request.CheckClientOrUser{ClientID: userRegistrationModel.ClientID.String},
		Response: authentication_response.CheckClientOrUserResponse{
			Nexsoft: authentication_response.CheckClientOrUserBodyResponse{
				Payload: authentication_response.CheckClientOrUserPayload{
					Data: authentication_response.CheckClientOrUserData{
						Content: authentication_response.CheckClientOrUserContent{
							IsExist: true,
							AdditionalInformation: authentication_response.AdditionalInformationContent{
								AliasName:    userRegistrationModel.ClientAliases.String,
								ClientID:     userRegistrationModel.ClientID.String,
								ClientSecret: license.clientCredentialModel.ClientSecret.String,
								UserID:       userRegistrationModel.AuthUserID.Int64,
								Username:     userRegistrationModel.Username.String,
								SignatureKey: license.clientCredentialModel.SignatureKey.String,
								GrantTypes:   "code_pkce",
								ResourceID:   resourceIDData,
								UserStatus:   1,
								Scope:        scopeData,
								Locale:       constanta.DefaultApplicationsLanguage,
							},
						},
					},
				},
			},
		},
	})

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
