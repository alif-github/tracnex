package test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
	"os"
	"strconv"
	"testing"
	"time"
)

func setServerConfiguration() {
	var err error
	config.GenerateConfiguration("testing")
	util.SetLogger(config.ApplicationConfiguration.GetLogFile())
	applicationModel.InitiateDefaultOperator()
	serverconfig.SetServerAttribute()

	//------------ error bundle
	//prefixPath := config.ApplicationConfiguration.GetServerPrefixPath()
	serverconfig.ServerAttribute.ErrorBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ErrorBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ErrorBundle.LoadMessageFile("../../i18n/common/error/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ErrorBundle.LoadMessageFile("../../i18n/common/error/id-ID.json")
	readError(err)
}

func dbMigrate() {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations_"+config.ApplicationConfiguration.GetPostgreSQLDefaultSchema(), "../../sql_migrations"),
	}
	if serverconfig.ServerAttribute.DBConnection != nil {
		n, err := migrate.Exec(serverconfig.ServerAttribute.DBConnection, "postgres", migrations, migrate.Up)
		if err != nil {
			logModel := applicationModel.GenerateLogModel("-", config.ApplicationConfiguration.GetServerResourceID())
			logModel.Message = err.Error()
			logModel.Status = 500
			util.LogError(logModel.ToLoggerObject())
			os.Exit(3)
		} else {
			logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
			logModel.Status = 200
			logModel.Message = "Applied " + strconv.Itoa(n) + " migrations!"
			util.LogInfo(logModel.ToLoggerObject())
		}
	} else {
		logModel := applicationModel.GenerateLogModel("-", config.ApplicationConfiguration.GetServerResourceID())
		logModel.Message = "null database"
		logModel.Status = 500
		util.LogError(logModel.ToLoggerObject())
		os.Exit(3)
	}
}

func InitAllConfiguration() (db *sql.DB) {
	// ---------------- Set Server Config
	setServerConfiguration()

	// ---------------- Migrate Auto
	dbMigrate()

	db = serverconfig.ServerAttribute.DBConnection

	return
}

func RollBackSchema(db *sql.DB) (err errorModel.ErrorModel) {
	err = util2.RollbackSchema(db)
	if err.Error != nil {
		return
	}

	return
}

func SetRequest(t *testing.T, db *sql.DB, method string, pathUrl string) (request *http.Request) {
	// ---------------- Set Request
	var errorS error
	request, errorS = http.NewRequest(method, pathUrl, nil)
	if errorS != nil || request == nil {
		RollBackSchema(db)
		t.Log(errorS)
		assert.FailNow(t, "Error build a request")
	}
	return
}


func setServerConfigurations() {
	var err error

	config.GenerateConfiguration("development")
	util.SetLogger(config.ApplicationConfiguration.GetLogFile())
	applicationModel.InitiateDefaultOperator()
	serverconfig.SetServerAttribute()

	//------------ error bundle
	//prefixPath := config.ApplicationConfiguration.GetLanguageDirectoryPath()
	serverconfig.ServerAttribute.ErrorBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ErrorBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ErrorBundle.LoadMessageFile("../../i18n/common/error/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ErrorBundle.LoadMessageFile("../../i18n/common/error/id-ID.json")
	readError(err)
}

func OpenTransactional(db *sql.DB) (tx *sql.Tx ,errS error) {
	if tx, errS = db.Begin(); errS != nil {
		fmt.Println(errS)
		return
	}
	return
}

func readError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}

func InitAllConfigurations() {
	// ---------------- Set Server Config
	setServerConfiguration()

	// ---------------- Migrate Auto
	dbMigrate()
	return
}

func SetClientType(dbc *sql.DB) error {
	stmt := `
	INSERT INTO nextrac2_testing.client_type (id, client_type, description)
	VALUES (4, 'nexChief', 'Client type NexChief'),
	    (5, 'nexMile', 'Client type NexMile test'),
	    (6, 'nexMile mobile', 'Client type NexMile mobile');
	`

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate test database tables")
	}

	return nil
}

func Truncate(dbc *sql.DB) error {
	stmt := `
		TRUNCATE TABLE 
		    client_mapping, client_credential, pkce_client_mapping,
		    customer_group, customer_category, salesman, customer,
		    product_group, license_type, license_variant, module, 
		    component, product, product_component, customer_site, 
		    customer_installation, license_configuration, license_configuration_productcomponent, 
		    product_license, user_license, user_registration_detail, user_registration,
		    client_mapping, customer_contact, nexmile_parameter ;
		
		ALTER SEQUENCE client_mapping_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE client_credential_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE pkce_client_mapping_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE customer_group_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE customer_category_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE salesman_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE customer_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE product_group_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE license_type_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE license_variant_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE module_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE component_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE product_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE product_component_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE customer_site_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE customer_installation_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE license_configuration_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE license_configuration_productcomponent_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE product_license_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE user_license_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE user_registration_detail_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE user_registration_pkey_seq RESTART WITH 1;
		ALTER SEQUENCE nexmile_parameter_pkey_seq RESTART WITH 1;`

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate test database tables")
	}

	return nil
}

func TruncateTx(dbc *sql.Tx) error {
	stmt := `
		TRUNCATE TABLE 
		    license_variant, license_type, customer_category, customer_group, salesman, customer, 
		    product, product_component, customer_installation, client_mapping, license_configuration, 
		    license_configuration_productcomponent, product_license, user_license, user_registration_detail, 
		    client_mapping, client_credential, customer_contact, customer_site, module, component;`

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate test database tables")
	}

	return nil
}

func StrToTime(timeStr string) (output time.Time) {
	output, _ = time.Parse(constanta.DefaultTimeFormat, timeStr)
	return
}

func SetDataWithTransactionalDB(contextModel applicationModel.ContextModel, serve func(*sql.Tx, applicationModel.ContextModel, time.Time) errorModel.ErrorModel) (err errorModel.ErrorModel) {
	var errs error
	var tx *sql.Tx

	timeNow := time.Now()
	funcName := "SetDataWithTransactionalDB"
	fileName := "testutil.go"

	defer func() {
		if errs != nil || err.Error != nil {
			errs = tx.Rollback()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
			}
		} else {
			errs = tx.Commit()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
			}
		}
	}()

	tx, errs = serverconfig.ServerAttribute.DBConnection.Begin()
	if errs != nil {
		return
	}

	return serve(tx, contextModel, timeNow)
}

func SetDataWithTransactionalDBWithOutput(contextModel applicationModel.ContextModel, serve func(*sql.Tx, applicationModel.ContextModel, time.Time) (interface{}, errorModel.ErrorModel)) (output interface{}, err errorModel.ErrorModel) {
	var errs error
	var tx *sql.Tx

	timeNow := time.Now()
	funcName := "SetDataWithTransactionalDBWithOutput"
	fileName := "testutil.go"

	defer func() {
		if errs != nil || err.Error != nil {
			errs = tx.Rollback()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
			}
		} else {
			errs = tx.Commit()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
			}
		}
	}()

	tx, errs = serverconfig.ServerAttribute.DBConnection.Begin()
	if errs != nil {
		return
	}

	return serve(tx, contextModel, timeNow)
}

func BeginDB(dbc *sql.DB) error {
	stmt := `BEGIN;`

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate test database tables")
	}

	return nil
}

func RollbackDB(dbc *sql.DB) error {
	stmt := `ROLLBACK;`

	if _, err := dbc.Exec(stmt); err != nil {
		return errors.Wrap(err, "truncate test database tables")
	}

	return nil
}