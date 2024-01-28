package serverconfig

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dgrr/fastws"
	"nexsoft.co.id/nextrac2/ws_handler/grochat_ws_handler"
	"os"
	"strconv"

	"github.com/Azure/azure-pipeline-go/pipeline"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/bukalapak/go-redis"
	"github.com/go-co-op/gocron"
	_ "github.com/go-sql-driver/mysql"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/olivere/elastic.v7"
	"nexsoft.co.id/nexcommon/db/dbconfig"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/model/applicationModel"
)

var ServerAttribute serverAttribute

type serverAttribute struct {
	Version                              string
	DBConnection                         *sql.DB
	DBConnectionView                     *sql.DB
	RedmineDBConnection                  *sql.DB
	RedmineInfraDBConnection             *sql.DB
	GroChatWSHandler					 *grochat_ws_handler.GroChatWSHandler
	GroChatWSConn                  		 *fastws.Conn
	GroChatWSReconnectSignal	   		 chan bool
	AzurePipeline                        pipeline.Pipeline
	RedisClient                          *redis.Client
	RedisClientSession                   *redis.Client
	ElasticClient                        *elastic.Client
	ErrorBundle                          *i18n.Bundle
	NotifyAddClientBundle                *i18n.Bundle
	AuditMonitoringBundle                *i18n.Bundle
	SessionBundle                        *i18n.Bundle
	ConstantaBundle                      *i18n.Bundle
	CommonServiceBundle                  *i18n.Bundle
	RoleServiceBundle                    *i18n.Bundle
	CountryBundle                        *i18n.Bundle
	ProvinceBundle                       *i18n.Bundle
	DistrictBundle                       *i18n.Bundle
	SubDistrictBundle                    *i18n.Bundle
	UrbanVillageBundle                   *i18n.Bundle
	UserBundle                           *i18n.Bundle
	PostalCodeBundle                     *i18n.Bundle
	IslandBundle                         *i18n.Bundle
	JobProcessBundle                     *i18n.Bundle
	PersonTitleBundle                    *i18n.Bundle
	CompanyTitleBundle                   *i18n.Bundle
	PositionBundle                       *i18n.Bundle
	BankBundle                           *i18n.Bundle
	PersonProfileBundle                  *i18n.Bundle
	CompanyProfileBundle                 *i18n.Bundle
	GroupOfDistributorBundle             *i18n.Bundle
	PrincipalBundle                      *i18n.Bundle
	VendorChannelBundle                  *i18n.Bundle
	ContactPersonBundle                  *i18n.Bundle
	BrandBundle                          *i18n.Bundle
	BrandOwnerBundle                     *i18n.Bundle
	ProductGroupHierarchyBundle          *i18n.Bundle
	ProductCategoryBundle                *i18n.Bundle
	ProductBrandBundle                   *i18n.Bundle
	ProductKeyAccountBundle              *i18n.Bundle
	ProductBundle                        *i18n.Bundle
	ProductHistoryBundle                 *i18n.Bundle
	DataScopeBundle                      *i18n.Bundle
	DataGroupBundle                      *i18n.Bundle
	CronSchedulerBundle                  *i18n.Bundle
	MenuBundle                           *i18n.Bundle
	ClientBundle                         *i18n.Bundle
	CustomerListBundle                   *i18n.Bundle
	ClientMappingBundle                  *i18n.Bundle
	PKCEUserBundle                       *i18n.Bundle
	AddResourceExternalBundle            *i18n.Bundle
	ActivationUserBundle                 *i18n.Bundle
	NexsoftRoleServiceBundle             *i18n.Bundle
	PKCEClientMappingServiceBundle       *i18n.Bundle
	ImportFileCustomerServiceBundle      *i18n.Bundle
	ClientRegistrationLogServiceBundle   *i18n.Bundle
	SalesmanServiceBundle                *i18n.Bundle
	CustomerGroupServiceBundle           *i18n.Bundle
	LicenseVariantServiceBundle          *i18n.Bundle
	ClientTypeServiceBundle              *i18n.Bundle
	CustomerSiteBundle                   *i18n.Bundle
	CustomerInstallationBundle           *i18n.Bundle
	LicenseConfigBundle                  *i18n.Bundle
	ProductLicenseBundle                 *i18n.Bundle
	UserLicenseBundle                    *i18n.Bundle
	ClientRegistrationNonOnPremiseBundle *i18n.Bundle
	ActivationLicenseBundle              *i18n.Bundle
	UserRegistrationDetailBundle         *i18n.Bundle
	UserRegistrationBundle               *i18n.Bundle
	ValidationLicenseBundle              *i18n.Bundle
	ActivationUserNexmileBundle          *i18n.Bundle
	ValidationNamedUserBundle            *i18n.Bundle
	UserVerificationBundle               *i18n.Bundle
	TodolistBundle                       *i18n.Bundle
	AccountRegistrationBundle            *i18n.Bundle
	RemarkBundle                         *i18n.Bundle
	ForgetPasswordBundle                 *i18n.Bundle
	CronScheduler                        *gocron.Scheduler
	TaskSchedulerBundle                  *i18n.Bundle
	DiscordConn                          *discordgo.Session
}

func SetServerAttribute() {
	dbParam := config.ApplicationConfiguration.GetPostgreSQLDefaultSchema()
	dbConnection := config.ApplicationConfiguration.GetPostgreSQLAddress()
	dbMaxOpenConnection := config.ApplicationConfiguration.GetPostgreSQLMaxOpenConnection()
	dbMaxIdleConnection := config.ApplicationConfiguration.GetPostgreSQLMaxIdleConnection()
	ServerAttribute.DBConnection = dbconfig.GetDbConnection(dbParam, dbConnection, dbMaxOpenConnection, dbMaxIdleConnection)

	dbParamView := config.ApplicationConfiguration.GetPostgreSQLDefaultSchema()
	dbConnectionView := config.ApplicationConfiguration.GetPostgreSQLAddressView()
	dbMaxOpenConnectionView := config.ApplicationConfiguration.GetPostgreSQLMaxOpenConnectionView()
	dbMaxIdleConnectionView := config.ApplicationConfiguration.GetPostgreSQLMaxIdleConnectionView()
	ServerAttribute.DBConnectionView = dbconfig.GetDbConnection(dbParamView, dbConnectionView, dbMaxOpenConnectionView, dbMaxIdleConnectionView)

	//--- Redmine Open Conn
	dbRedmineConnectionView := config.ApplicationConfiguration.GetRedmineDBAddress()
	dbRedmineMaxOpenConnectionView := config.ApplicationConfiguration.GetRedmineDBMaxOpenConnection()
	dbRedmineMaxIdleConnectionView := config.ApplicationConfiguration.GetRedmineDBMaxIdleConnection()
	instance, errOpenRedmine := sql.Open("pgx", dbRedmineConnectionView)
	if errOpenRedmine != nil {
		fmt.Println("error open connect to DB Redmine ", errOpenRedmine)
		fmt.Println(fmt.Sprintf(`connect failed to DB Redmine %v`, dbRedmineConnectionView))
		instance = nil
		os.Exit(1)
		return
	}

	instance.SetMaxOpenConns(dbRedmineMaxOpenConnectionView)
	instance.SetMaxIdleConns(dbRedmineMaxIdleConnectionView)

	fmt.Println("Connected to the database Redmine !")
	ServerAttribute.RedmineDBConnection = instance
	//--- End Of Open Conn

	//--- Redmine Infra Open Conn
	dbRedmineInfraConnection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.ApplicationConfiguration.GetRedmineInfraDBUser(),
		config.ApplicationConfiguration.GetRedmineInfraDBPass(), config.ApplicationConfiguration.GetRedmineInfraDBHost(),
		config.ApplicationConfiguration.GetRedmineInfraDBPort(), config.ApplicationConfiguration.GetRedmineInfraDBDatabase())
	dbConn, errOpenInfraRedmine := sql.Open("mysql", dbRedmineInfraConnection)
	if errOpenInfraRedmine != nil {
		fmt.Println("Error open connect to DB Infra Redmine : ", errOpenInfraRedmine)
		fmt.Println(fmt.Sprintf(`Open connect failed to DB %v`, dbRedmineInfraConnection))
		dbConn = nil
		os.Exit(1)
		return
	}

	fmt.Println("Connected to the database Redmine Infra !")
	ServerAttribute.RedmineInfraDBConnection = dbConn
	//--- End Of Open Conn

	redisHost := config.ApplicationConfiguration.GetRedisHost()
	redisDB := config.ApplicationConfiguration.GetRedisDB()
	redisSessionDB := config.ApplicationConfiguration.GetRedisSessionDB()
	redisPassword := config.ApplicationConfiguration.GetRedisPassword()
	redisPort := config.ApplicationConfiguration.GetRedisPort()
	redisTimeout := config.ApplicationConfiguration.GetRedisTimeout()
	redisVolumeThreshold := config.ApplicationConfiguration.GetRedisRequestVolumeThreshold()
	redisSleepWindow := config.ApplicationConfiguration.GetRedisSleepWindow()
	redisErrorPercentThreshold := config.ApplicationConfiguration.GetRedisErrorPercentThreshold()
	redisMaxConcurrentRequest := config.ApplicationConfiguration.GetRedisMaxConcurrentRequests()

	optCB := &hystrix.CommandConfig{
		Timeout:                redisTimeout,
		RequestVolumeThreshold: redisVolumeThreshold,
		SleepWindow:            redisSleepWindow,
		ErrorPercentThreshold:  redisErrorPercentThreshold,
		MaxConcurrentRequests:  redisMaxConcurrentRequest,
	}

	ServerAttribute.RedisClient = getRedisClient(redisHost, redisPort, redisDB, redisPassword, optCB)
	ServerAttribute.RedisClientSession = getRedisClient(redisHost, redisPort, redisSessionDB, redisPassword, optCB)
	ServerAttribute.Version = config.ApplicationConfiguration.GetServerVersion()

	//-- Initiate Cron Job
	//ServerAttribute.CronScheduler, _ = gocron.NewScheduler(time.FixedZone(config.ApplicationConfiguration.GetLocalTimezone().Zone, config.ApplicationConfiguration.GetLocalTimezone().Offset))
	//if err != nil {
	//	println("Error Create Cron Scheduler : ", err.Error())
	//	return
	//}

	//ServerAttribute.CronScheduler.StartAsync()

	fmt.Println(fmt.Sprintf(`Azure path -> %s%s`, config.ApplicationConfiguration.GetAzure().Host, config.ApplicationConfiguration.GetAzure().Suffix))
	credential, err := azblob.NewSharedKeyCredential(config.ApplicationConfiguration.GetAzure().AccountName, config.ApplicationConfiguration.GetAzure().AccountKey)
	if err != nil {
		//-- Print Log
		fmt.Println(fmt.Sprintf(`Error Credential Azure -> %s | Azure ACC NAME -> %s | Azure ACC KEY -> %s`, err.Error(), config.ApplicationConfiguration.GetAzure().AccountName, config.ApplicationConfiguration.GetAzure().AccountKey))
		logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
		logModel.Message = "Invalid credentials with error: " + err.Error()
		logModel.Status = 500
		util.LogError(logModel.ToLoggerObject())
		os.Exit(1)
		return
	}

	ServerAttribute.AzurePipeline = azblob.NewPipeline(credential, azblob.PipelineOptions{})
	ServerAttribute.ElasticClient, err = elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))
	if err != nil {
		logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
		logModel.Message = "Invalid Connect to Elastic Search with message : " + err.Error()
		logModel.Status = 500
		util.LogError(logModel.ToLoggerObject())
		os.Exit(1)
		return
	}

	var (
		errDiscord error
		token      = config.ApplicationConfiguration.GetDiscordLogToken()
	)

	fmt.Println("Token : ", token)
	ServerAttribute.DiscordConn, errDiscord = discordgo.New(token)
	if errDiscord != nil {
		fmt.Println("Discord Connecting Failed !")
		os.Exit(1)
		return
	}

	//errDiscord = ServerAttribute.DiscordConn.Open()
	//if errDiscord != nil {
	//	fmt.Println("Discord Open Failed !", errDiscord)
	//	os.Exit(1)
	//	return
	//}
	//ServerAttribute.GroChatWSConn = dialGroChatWS()
	//ServerAttribute.GroChatWSReconnectSignal = make(chan bool)
}

//func dialGroChatWS() *fastws.Conn {
//	conn, err := fastws.Dial(config.ApplicationConfiguration.GetGroChatWS().Host)
//	if err != nil {
//		logModel := model.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion())
//		logModel.Status = 500
//		logModel.Message = fmt.Sprintf("WebSocket dial error : %s", err.Error())
//
//		util.LogError(logModel.ToLoggerObject())
//		os.Exit(1)
//	}
//
//	conn.WriteTimeout = 0
//	conn.ReadTimeout = 0
//
//	return conn
//}

func getRedisClient(host string, port int, db int, password string, optCB *hystrix.CommandConfig) *redis.Client {
	redisAddress := host + ":" + strconv.Itoa(port)
	opts := &redis.Options{
		CircuitBreaker: optCB,
		Addr:           redisAddress,
		Password:       password,
		DB:             db,
	}

	return redis.NewClient(opts)
}
