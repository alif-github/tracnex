package license_variant

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/test"
	"os"
	"testing"
	"time"
)

var tx *sql.Tx
var contextModel applicationModel.ContextModel

func TestMain(m *testing.M) {
	os.Exit(testMain(m))
}

func testMain(m *testing.M) int {
	//BEFORE
	var err, errS error
	fmt.Println("Start Testing License Variant")

	// Set ContextModel
	contextModel = applicationModel.ContextModel{
		AuthAccessTokenModel: model.AuthAccessTokenModel{
			RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
				ResourceUserID: 12,
			},
			ClientID: "3e3cb40e14d645eb8783f53a30c822d4",
			Locale:   constanta.IndonesianLanguage,
		},
	}

	// Set Configuration
	fmt.Println("Init config")
	test.InitAllConfiguration()
	if tx, errS = serverconfig.ServerAttribute.DBConnection.Begin(); errS != nil {
		fmt.Println(errS)
		return 1
	}

	test.Truncate(serverconfig.ServerAttribute.DBConnection)

	//Set Database
	fmt.Println("Open Connection DB Connection")
	if err = setDatabase(serverconfig.ServerAttribute.DBConnection); err != nil {
		fmt.Println(err.Error())
		return 1
	}


	//Truncate function
	fmt.Println("truncate DB")
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	return m.Run()
}

func setDatabase(db *sql.DB) (err error) {
	timeNow := time.Now()
	dataLicenseVariant := []repository.LicenseVariantModel{
		{
			ID:                 sql.NullInt64{Int64: 90},
			LicenseVariantName: sql.NullString{String: "License Test"},
			CreatedBy:          sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			CreatedClient:      sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			CreatedAt:          sql.NullTime{Time: test.StrToTime("2021-12-14T10:16:39.631007Z")},
			UpdatedBy:          sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedClient:      sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			UpdatedAt:          sql.NullTime{Time: test.StrToTime("2021-12-14T10:16:39.631007Z")},
		},
		{
			ID:                 sql.NullInt64{Int64: 92},
			LicenseVariantName: sql.NullString{String: "License Test 2"},
			CreatedBy:          sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			CreatedClient:      sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			CreatedAt:          sql.NullTime{Time: timeNow},
			UpdatedBy:          sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedClient:      sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			UpdatedAt:          sql.NullTime{Time: timeNow},
		},
	}

	if err = dao.LicenseVariantDAO.InsertDataForTesting(db, dataLicenseVariant); err != nil {
		return
	}

	return nil
}
