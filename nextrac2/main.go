package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	migrate "github.com/rubenv/sql-migrate"
	"golang.org/x/text/language"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/router"
	"nexsoft.co.id/nextrac2/scheduledtask/scheduledconfig"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"nexsoft.co.id/nextrac2/ws_handler/grochat_ws_handler"
	"os"
	"strconv"
	"strings"
)

func main() {
	var (
		arguments = "local"
		args      = os.Args
	)

	if len(args) > 1 {
		arguments = args[1]
	}

	config.GenerateConfiguration(arguments)

	util.SetLogger(config.ApplicationConfiguration.GetLogFile())
	applicationModel.InitiateDefaultOperator()

	serverconfig.SetServerAttribute()
	loadBundleI18N()

	dbMigrate()

	var (
		logModel   = applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
		versionNum int
		err        error
	)

	defer func() {
		if err != nil {
			logModel.Status = 500
			logModel.Message = err.Error()
			util.LogError(logModel.ToLoggerObject())
		}

		err = serverconfig.ServerAttribute.DBConnection.Close()
		if err != nil {
			logModel.Status = 500
			logModel.Message = "Failed to close NexTrac2 DB Connection " + err.Error()
			util.LogError(logModel.ToLoggerObject())
		}

		err = serverconfig.ServerAttribute.RedmineDBConnection.Close()
		if err != nil {
			logModel.Status = 500
			logModel.Message = "Failed to close Redmine DB Connection " + err.Error()
			util.LogError(logModel.ToLoggerObject())
		}

		err = serverconfig.ServerAttribute.RedmineInfraDBConnection.Close()
		if err != nil {
			logModel.Status = 500
			logModel.Message = "Failed to close Redmine Infra DB Connection " + err.Error()
			util.LogError(logModel.ToLoggerObject())
		}
	}()

	// Generate auto add host
	if config.ApplicationConfiguration.GetServerAutoAddHost() {
		service.GenerateIPAndHostNameService.GenerateIPAndServerID()
	}

	// Generate scheduler cron
	scheduledconfig.GenerateSchedulerCron(serverconfig.ServerAttribute.DBConnection)

	// Check server version for versioning permission
	if config.ApplicationConfiguration.GetServerVersion() == "" {
		err = errors.New(`version of app must be filled`)
		return
	} else {
		version := strings.Split(config.ApplicationConfiguration.GetServerVersion(), ".")
		if len(version) > 0 {
			versionNum, err = strconv.Atoi(version[0])
			if err != nil {
				return
			}
		}
	}

	// Initial role and menu
	errS := initialRolePermissionAndMenuActivated(serverconfig.ServerAttribute.DBConnection, versionNum)
	if errS.Error != nil {
		err = errors.New(errS.CausedBy.Error())
		return
	}

	logModel.Status = 200
	logModel.Message = "Server Start in port : " + strconv.Itoa(config.ApplicationConfiguration.GetServerPort())
	util.LogInfo(logModel.ToLoggerObject())
	util2.DiscordInfoServerRunning(arguments)

	/*
		GroChat WS
	*/
	startGroChatWS()
	//temporary disable because error "Failed to dial websocket : cannot upgrade connection"

	//Task.SchAddResourceNexcloud.StartMain()
	//fmt.Println(resource_common_service.GenerateInternalToken("auth", 0, "", "Testing auth", "id-ID"))
	router.APIController()
}

func dbMigrate() {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations", "./sql_migrations"),
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

func loadBundleI18N() {
	prefixPath := config.ApplicationConfiguration.GetLanguageDirectoryPath()

	//------------ error bundle
	serverconfig.ServerAttribute.ErrorBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ErrorBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err := serverconfig.ServerAttribute.ErrorBundle.LoadMessageFile(prefixPath + "/common/error/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.ErrorBundle.LoadMessageFile(prefixPath + "/common/error/id-ID.json")
	readError(err)

	//------------ common service
	serverconfig.ServerAttribute.CommonServiceBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.CommonServiceBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.CommonServiceBundle.LoadMessageFile(prefixPath + "/common_service/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.CommonServiceBundle.LoadMessageFile(prefixPath + "/common_service/id-ID.json")
	readError(err)

	//------------ user
	serverconfig.ServerAttribute.UserBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.UserBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.UserBundle.LoadMessageFile(prefixPath + "/user/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.UserBundle.LoadMessageFile(prefixPath + "/user/id-ID.json")
	readError(err)

	//------------ client
	serverconfig.ServerAttribute.ClientBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ClientBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ClientBundle.LoadMessageFile(prefixPath + "/client/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.ClientBundle.LoadMessageFile(prefixPath + "/client/id-ID.json")
	readError(err)

	//------------ customer list
	serverconfig.ServerAttribute.CustomerListBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.CustomerListBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.CustomerListBundle.LoadMessageFile(prefixPath + "/customer_list/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.CustomerListBundle.LoadMessageFile(prefixPath + "/customer_list/id-ID.json")
	readError(err)

	//------------ client mapping
	serverconfig.ServerAttribute.ClientMappingBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ClientMappingBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ClientMappingBundle.LoadMessageFile(prefixPath + "/client_mapping/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.ClientMappingBundle.LoadMessageFile(prefixPath + "/client_mapping/id-ID.json")
	readError(err)

	//------------ pkce user
	serverconfig.ServerAttribute.PKCEUserBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.PKCEUserBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.PKCEUserBundle.LoadMessageFile(prefixPath + "/pkce_user/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.PKCEUserBundle.LoadMessageFile(prefixPath + "/pkce_user/id-ID.json")
	readError(err)

	//------------ add resource external
	serverconfig.ServerAttribute.AddResourceExternalBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.AddResourceExternalBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.AddResourceExternalBundle.LoadMessageFile(prefixPath + "/add_resource_external/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.AddResourceExternalBundle.LoadMessageFile(prefixPath + "/add_resource_external/id-ID.json")
	readError(err)

	//------------ constanta
	serverconfig.ServerAttribute.ConstantaBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ConstantaBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ConstantaBundle.LoadMessageFile(prefixPath + "/common/constanta/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.ConstantaBundle.LoadMessageFile(prefixPath + "/common/constanta/id-ID.json")
	readError(err)

	//------------ activation user
	serverconfig.ServerAttribute.ActivationUserBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ActivationUserBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ActivationUserBundle.LoadMessageFile(prefixPath + "/activation_user/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.ActivationUserBundle.LoadMessageFile(prefixPath + "/activation_user/id-ID.json")
	readError(err)

	//------------ role
	serverconfig.ServerAttribute.RoleServiceBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.RoleServiceBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.RoleServiceBundle.LoadMessageFile(prefixPath + "/role/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.RoleServiceBundle.LoadMessageFile(prefixPath + "/role/id-ID.json")
	readError(err)

	//------------ session
	serverconfig.ServerAttribute.SessionBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.SessionBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.SessionBundle.LoadMessageFile(prefixPath + "/session/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.SessionBundle.LoadMessageFile(prefixPath + "/session/id-ID.json")
	readError(err)

	//------------ nexsoft role
	serverconfig.ServerAttribute.NexsoftRoleServiceBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.NexsoftRoleServiceBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.NexsoftRoleServiceBundle.LoadMessageFile(prefixPath + "/nexsoft_role/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.NexsoftRoleServiceBundle.LoadMessageFile(prefixPath + "/nexsoft_role/id-ID.json")
	readError(err)

	//------------ pkce client mapping
	serverconfig.ServerAttribute.PKCEClientMappingServiceBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.PKCEClientMappingServiceBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.PKCEClientMappingServiceBundle.LoadMessageFile(prefixPath + "/pkce_client_mapping/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.PKCEClientMappingServiceBundle.LoadMessageFile(prefixPath + "/pkce_client_mapping/id-ID.json")
	readError(err)

	//------------ import file customer
	serverconfig.ServerAttribute.ImportFileCustomerServiceBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ImportFileCustomerServiceBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ImportFileCustomerServiceBundle.LoadMessageFile(prefixPath + "/import_file/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ImportFileCustomerServiceBundle.LoadMessageFile(prefixPath + "/import_file/id-ID.json")
	readError(err)

	//------------ job process
	serverconfig.ServerAttribute.JobProcessBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.JobProcessBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.JobProcessBundle.LoadMessageFile(prefixPath + "/job_process/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.JobProcessBundle.LoadMessageFile(prefixPath + "/job_process/id-ID.json")
	readError(err)

	//------------ log
	serverconfig.ServerAttribute.ClientRegistrationLogServiceBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ClientRegistrationLogServiceBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ClientRegistrationLogServiceBundle.LoadMessageFile(prefixPath + "/client_registration_log/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ClientRegistrationLogServiceBundle.LoadMessageFile(prefixPath + "/client_registration_log/id-ID.json")
	readError(err)

	//------------ menu
	serverconfig.ServerAttribute.MenuBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.MenuBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.MenuBundle.LoadMessageFile(prefixPath + "/menu/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.MenuBundle.LoadMessageFile(prefixPath + "/menu/id-ID.json")
	readError(err)

	//------------ person title
	serverconfig.ServerAttribute.PersonTitleBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.PersonTitleBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.PersonTitleBundle.LoadMessageFile(prefixPath + "/person_title/en-US.json")
	readError(err)

	_, err = serverconfig.ServerAttribute.PersonTitleBundle.LoadMessageFile(prefixPath + "/person_title/id-ID.json")
	//------------ customer group
	serverconfig.ServerAttribute.CustomerGroupServiceBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.CustomerGroupServiceBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.CustomerGroupServiceBundle.LoadMessageFile(prefixPath + "/customer_group/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.CustomerGroupServiceBundle.LoadMessageFile(prefixPath + "/customer_group/id-ID.json")
	readError(err)

	//------------ data group
	serverconfig.ServerAttribute.DataGroupBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.DataGroupBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.DataGroupBundle.LoadMessageFile(prefixPath + "/data_group/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.DataGroupBundle.LoadMessageFile(prefixPath + "/data_group/id-ID.json")
	readError(err)

	//------------ salesman
	serverconfig.ServerAttribute.SalesmanServiceBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.SalesmanServiceBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.SalesmanServiceBundle.LoadMessageFile(prefixPath + "/salesman/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.SalesmanServiceBundle.LoadMessageFile(prefixPath + "/salesman/id-ID.json")
	readError(err)

	//------------ Province
	serverconfig.ServerAttribute.ProvinceBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ProvinceBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ProvinceBundle.LoadMessageFile(prefixPath + "/province/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ProvinceBundle.LoadMessageFile(prefixPath + "/province/id-ID.json")
	readError(err)

	//------------ District
	serverconfig.ServerAttribute.DistrictBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.DistrictBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.DistrictBundle.LoadMessageFile(prefixPath + "/district/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.DistrictBundle.LoadMessageFile(prefixPath + "/district/id-ID.json")
	readError(err)

	//------------ License Variant
	serverconfig.ServerAttribute.LicenseVariantServiceBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.LicenseVariantServiceBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.LicenseVariantServiceBundle.LoadMessageFile(prefixPath + "/license_variant/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.LicenseVariantServiceBundle.LoadMessageFile(prefixPath + "/license_variant/id-ID.json")
	readError(err)

	//------------ Product
	serverconfig.ServerAttribute.ProductBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ProductBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ProductBundle.LoadMessageFile(prefixPath + "/product/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ProductBundle.LoadMessageFile(prefixPath + "/product/id-ID.json")
	readError(err)

	//------------ Client Type
	serverconfig.ServerAttribute.ClientTypeServiceBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ClientTypeServiceBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ClientTypeServiceBundle.LoadMessageFile(prefixPath + "/client_type/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ClientTypeServiceBundle.LoadMessageFile(prefixPath + "/client_type/id-ID.json")
	readError(err)

	//------------ Customer Site
	serverconfig.ServerAttribute.CustomerSiteBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.CustomerSiteBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.CustomerSiteBundle.LoadMessageFile(prefixPath + "/customer_site/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.CustomerSiteBundle.LoadMessageFile(prefixPath + "/customer_site/id-ID.json")
	readError(err)

	//------------ Customer Installation
	serverconfig.ServerAttribute.CustomerInstallationBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.CustomerInstallationBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.CustomerInstallationBundle.LoadMessageFile(prefixPath + "/customer_installation/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.CustomerInstallationBundle.LoadMessageFile(prefixPath + "/customer_installation/id-ID.json")
	readError(err)

	//------------ License Config
	serverconfig.ServerAttribute.LicenseConfigBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.LicenseConfigBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.LicenseConfigBundle.LoadMessageFile(prefixPath + "/license_config/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.LicenseConfigBundle.LoadMessageFile(prefixPath + "/license_config/id-ID.json")
	readError(err)

	//------------ Product License
	serverconfig.ServerAttribute.ProductLicenseBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ProductLicenseBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ProductLicenseBundle.LoadMessageFile(prefixPath + "/product_license/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ProductLicenseBundle.LoadMessageFile(prefixPath + "/product_license/id-ID.json")
	readError(err)

	//------------ User License
	serverconfig.ServerAttribute.UserLicenseBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.UserLicenseBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.UserLicenseBundle.LoadMessageFile(prefixPath + "/user_license/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.UserLicenseBundle.LoadMessageFile(prefixPath + "/user_license/id-ID.json")
	readError(err)

	//------------ User Registration Detail
	serverconfig.ServerAttribute.UserRegistrationDetailBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.UserRegistrationDetailBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.UserRegistrationDetailBundle.LoadMessageFile(prefixPath + "/user_registration_detail/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.UserRegistrationDetailBundle.LoadMessageFile(prefixPath + "/user_registration_detail/id-ID.json")
	readError(err)

	//------------ Client Registration Non On Premise
	serverconfig.ServerAttribute.ClientRegistrationNonOnPremiseBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ClientRegistrationNonOnPremiseBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ClientRegistrationNonOnPremiseBundle.LoadMessageFile(prefixPath + "/client_registration_non_on_premise/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ClientRegistrationNonOnPremiseBundle.LoadMessageFile(prefixPath + "/client_registration_non_on_premise/id-ID.json")
	readError(err)

	//------------ Activation License
	serverconfig.ServerAttribute.ActivationLicenseBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ActivationLicenseBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ActivationLicenseBundle.LoadMessageFile(prefixPath + "/activation_license/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ActivationLicenseBundle.LoadMessageFile(prefixPath + "/activation_license/id-ID.json")
	readError(err)

	//------------ Validation License
	serverconfig.ServerAttribute.ValidationLicenseBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ValidationLicenseBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ValidationLicenseBundle.LoadMessageFile(prefixPath + "/validation_license/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ValidationLicenseBundle.LoadMessageFile(prefixPath + "/validation_license/id-ID.json")
	readError(err)

	//------------ Activation User Nexmile
	serverconfig.ServerAttribute.ActivationUserNexmileBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ActivationUserNexmileBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ActivationUserNexmileBundle.LoadMessageFile(prefixPath + "/activation_user_nexmile/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ActivationUserNexmileBundle.LoadMessageFile(prefixPath + "/activation_user_nexmile/id-ID.json")
	readError(err)

	//------------ Validation Named User Online
	serverconfig.ServerAttribute.ValidationNamedUserBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ValidationNamedUserBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ValidationNamedUserBundle.LoadMessageFile(prefixPath + "/validation_named_user/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ValidationNamedUserBundle.LoadMessageFile(prefixPath + "/validation_named_user/id-ID.json")
	readError(err)

	//------------ User Verification
	serverconfig.ServerAttribute.UserVerificationBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.UserVerificationBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.UserVerificationBundle.LoadMessageFile(prefixPath + "/user_verification/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.UserVerificationBundle.LoadMessageFile(prefixPath + "/user_verification/id-ID.json")
	readError(err)

	//------------ Task Scheduler
	serverconfig.ServerAttribute.TaskSchedulerBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.TaskSchedulerBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.TaskSchedulerBundle.LoadMessageFile(prefixPath + "/task_scheduler/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.TaskSchedulerBundle.LoadMessageFile(prefixPath + "/task_scheduler/id-ID.json")
	readError(err)

	//------------ Reset Password
	serverconfig.ServerAttribute.ForgetPasswordBundle = i18n.NewBundle(language.Indonesian)
	serverconfig.ServerAttribute.ForgetPasswordBundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	_, err = serverconfig.ServerAttribute.ForgetPasswordBundle.LoadMessageFile(prefixPath + "/forget_password/en-US.json")
	readError(err)
	_, err = serverconfig.ServerAttribute.ForgetPasswordBundle.LoadMessageFile(prefixPath + "/forget_password/id-ID.json")
	readError(err)
}

var fileNumber = 0

func readError(err error) {
	fileNumber++
	if err != nil {
		fmt.Println(err.Error() + " at file " + strconv.Itoa(fileNumber))
		os.Exit(3)
	}
}

func initialRolePermissionAndMenuActivated(db *sql.DB, version int) (err errorModel.ErrorModel) {
	var (
		role           = "super_user"
		permission     = fmt.Sprintf(`{"insert":["insert"],"view":["view"],"update":["update"],"delete":["delete"],"admin":["insert","update-own","view-own","changepassword-own"],"master":["insert","view","update","delete"],"nexsoft":["insert","update","delete","view"]}`)
		status         = constanta.StatusActive
		operatorMaster = "="
		operatorAdmin  = "="
		isExist        bool
		id             int64
	)

	if version == constanta.VersionRedesign {
		permission = fmt.Sprintf(`{"insert":["insert"],"view":["view"],"update":["update"],"delete":["delete"],"admin":["insert","update-own","view-own","changepassword-own"],"setup":["insert","view","update","delete"],"konsumen":["insert","view","update","delete"],"produk":["insert","view","update","delete"],"lisensi":["insert","view","update","delete"],"timesheet":["insert","view","update","delete"],"home":["view"]}`)
		operatorMaster = "<>"
	}

	// Get Role Permission
	id, isExist, err = dao.RoleDAO.GetRolePermission(db, role)
	if err.Error != nil {
		return
	}

	// Update Role Permission
	if isExist {
		err = dao.RoleDAO.UpdateRolePermission(db, id, permission)
		if err.Error != nil {
			return
		}
	}

	// Reset Menu
	err = dao.MenuDAO.ResetMenuRecursive(db)
	if err.Error != nil {
		return
	}

	// Update Menu Sync Version
	err = dao.MenuDAO.UpdateMenuRecursive(db, status, operatorMaster, operatorAdmin)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func startGroChatWS() {
	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
	groChatUser := config.ApplicationConfiguration.GetGroChatWS().User
	groChatWSHandler := grochat_ws_handler.NewGroChatWSHandler()
	//groChatWS := config.ApplicationConfiguration.GetGroChatWS()

	//conn, _ , err := websocket.DefaultDialer.Dial(groChatWS.Host + groChatWS.PathRedirect.WS, nil)
	//if err != nil {
		// handle error
		//logModel.Status = 500
		//logModel.Message = fmt.Sprintf("Failed to dial websocket : %s", err.Error())
		//util.LogError(logModel.ToLoggerObject())

		//os.Exit(1)
		//return
	//}
	//fmt.Println(connection.)
	//defer conn.Close()

	if err := groChatWSHandler.Dial(); err != nil {
		logModel.Status = 500
		logModel.Message = fmt.Sprintf("Failed to dial websocket : %s", err.Error())
		util.LogError(logModel.ToLoggerObject())

		os.Exit(1)
		return
	} else {
		fmt.Println("Connected to server websocket !")
	}

	if err := groChatWSHandler.Authenticate(groChatUser.Username, groChatUser.Password); err != nil {
		logModel.Status = 500
		logModel.Message = fmt.Sprintf("Authentication failed : %s", err.Error())
		util.LogError(logModel.ToLoggerObject())

		os.Exit(1)
		return
	}

	groChatWSHandler.Start()
	go groChatWSHandler.KeepAlive()

	/*
		Reconnection
	*/
	go func() {
		for {
			<-groChatWSHandler.ReconnectSignal()

			for {
				logModel.Status = 200
				logModel.Message = fmt.Sprintf("Reconnecting...")
				util.LogInfo(logModel.ToLoggerObject())

				if err := groChatWSHandler.Dial(); err != nil {
					continue
				}

				break
			}

			logModel.Status = 200
			logModel.Message = fmt.Sprintf("Connected")
			util.LogInfo(logModel.ToLoggerObject())

			groChatWSHandler.Start()
		}
	}()

	serverconfig.ServerAttribute.GroChatWSHandler = groChatWSHandler
}