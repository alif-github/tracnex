package config

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"os"
	"strconv"
)

type ProductionConfig struct {
	Configuration
	Server struct {
		Protocol      string `json:"protocol"`
		Ethernet      string `json:"ethernet"`
		AutoAddHost   bool   `json:"auto_add_host"`
		AutoAddClient bool   `json:"auto_add_client"`
		Host          string `envconfig:"NEXTRAC2_HOST"`
		Port          string `envconfig:"NEXTRAC2_PORT"`
		Version       string `json:"version"`
		ResourceID    string `envconfig:"NEXTRAC2_RESOURCE_ID"`
		PrefixPath    string `json:"prefix_path"`
		LogLevel      string `json:"log_level"`
	} `json:"server"`
	Postgresql struct {
		Address           string `envconfig:"NEXTRAC2_DB_CONNECTION"`
		DefaultSchema     string `envconfig:"NEXTRAC2_DB_PARAM"`
		MaxOpenConnection int    `json:"max_open_connection"`
		MaxIdleConnection int    `json:"max_idle_connection"`
	} `json:"postgresql"`
	PostgresqlView struct {
		Address           string `envconfig:"NEXTRAC2_DB_CONNECTION_VIEW"`
		DefaultSchema     string `envconfig:"NEXTRAC2_DB_PARAM_VIEW"`
		MaxOpenConnection int    `json:"max_open_connection"`
		MaxIdleConnection int    `json:"max_idle_connection"`
	} `json:"postgresql_view"`
	Redis struct {
		Host                   string `envconfig:"NEXTRAC2_REDIS_HOST"`
		Port                   string `envconfig:"NEXTRAC2_REDIS_PORT"`
		Db                     string `envconfig:"NEXTRAC2_REDIS_DB"`
		Password               string `envconfig:"NEXTRAC2_REDIS_PASSWORD"`
		Timeout                int    `json:"timeout"`
		RequestVolumeThreshold int    `json:"request_volume_threshold"`
		SleepWindow            int    `json:"sleep_window"`
		ErrorPercentThreshold  int    `json:"error_percent_threshold"`
		MaxConcurrentRequests  int    `json:"max_concurrent_requests"`
	} `json:"redis"`
	RedisSession struct {
		Db string `envconfig:"NEXTRAC2_REDIS_SESSION_DB"`
	} `json:"redis_session"`
	ClientCredentials struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `envconfig:"NEXTRAC2_CLIENT_SECRET"`
		SecretKey    string `envconfig:"NEXTRAC2_SIGNATURE_KEY"`
		AuthUserID   int64  `json:"auth_user_id"`
	} `json:"client_credentials"`
	LogFile []string `json:"log_file"`
	JWTKey  struct {
		JWT      string `envconfig:"NEXTRAC2_JWT_KEY"`
		Internal string `envconfig:"NEXTRAC2_INTERNAL_KEY"`
	} `json:"jwt_key"`
	LanguageDirectoryPath string               `json:"language_directory_path"`
	NextracFrontend       NextracFrontend      `json:"nextrac_frontend"`
	AuthenticationServer  AuthenticationServer `json:"authentication_server"`
	NexcloudAPI           NexcloudAPI          `json:"nexcloud_api"`
	Nexdrive              Nexdrive             `json:"nexdrive"`
	Grochat               Grochat              `json:"grochat"`
	GrochatWS             GrochatWS            `json:"grochat_ws"`
	Nexmile               Nexmile              `json:"nexmile"`
	Nexstar               Nexstar              `json:"nexstar"`
	Nextrade              Nextrade             `json:"nextrade"`
	CommonPath            CommonPath           `json:"common_path"`
	Audit                 Audit                `json:"audit"`
	AlertServer           AlertServer          `json:"alert_server"`
	Azure                 struct {
		AccountName string `envconfig:"AZURE_ACCOUNT_NAME"`
		AccountKey  string `envconfig:"AZURE_ACCOUNT_KEY"`
		Host        string `json:"host"`
		Suffix      string `json:"suffix"`
	} `json:"azure"`
	Cdn           CDN `json:"cdn"`
	ElasticSearch struct {
		ConnectionStr string `json:"connection_str"`
	} `json:"elastic_search"`
	DataDirectory DataDirectory `json:"data_directory"`
	Scheduler     Scheduler     `json:"scheduler"`
	MasterData    MasterData    `json:"master_data"`
	Generator     Generator     `json:"generator"`
	Email         Email         `json:"email"`
	Redmine       struct {
		Address           string     `json:"address"`
		DefaultSchema     string     `json:"default_schema"`
		MaxOpenConnection int        `json:"max_open_connection"`
		MaxIdleConnection int        `json:"max_idle_connection"`
		AccessKey         string     `json:"api_access_key"`
		Api               ApiRedmine `json:"api"`
	} `json:"redmine"`
	RedmineInfra struct {
		Host      string `json:"host"`
		Port      string `json:"port"`
		User      string `json:"username"`
		Pass      string `json:"password"`
		Database  string `json:"database"`
		AccessKey string `json:"api_access_key"`
	} `json:"redmine_infra"`
}

func (input ProductionConfig) GetServerProtocol() string {
	return input.Server.Protocol
}

func (input ProductionConfig) GetServerEthernet() string {
	return input.Server.Ethernet
}

func (input ProductionConfig) GetServerAutoAddHost() bool {
	return input.Server.AutoAddHost
}

func (input ProductionConfig) GetServerHost() string {
	return input.Server.Host
}

func (input ProductionConfig) GetServerPort() int {
	return convertStringParamToInt("Server Port", input.Server.Port)
}

func (input ProductionConfig) GetServerVersion() string {
	return input.Server.Version
}

func (input ProductionConfig) GetServerResourceID() string {
	return input.Server.ResourceID
}

func (input ProductionConfig) GetPostgreSQLAddress() string {
	return input.Postgresql.Address
}

func (input ProductionConfig) GetPostgreSQLDefaultSchema() string {
	return input.Postgresql.DefaultSchema
}

func (input ProductionConfig) GetPostgreSQLMaxOpenConnection() int {
	return input.Postgresql.MaxOpenConnection
}

func (input ProductionConfig) GetPostgreSQLMaxIdleConnection() int {
	return input.Postgresql.MaxIdleConnection
}

func (input ProductionConfig) GetPostgreSQLAddressView() string {
	return input.PostgresqlView.Address
}

func (input ProductionConfig) GetPostgreSQLDefaultSchemaView() string {
	return input.PostgresqlView.DefaultSchema
}

func (input ProductionConfig) GetPostgreSQLMaxOpenConnectionView() int {
	return input.PostgresqlView.MaxOpenConnection
}

func (input ProductionConfig) GetPostgreSQLMaxIdleConnectionView() int {
	return input.PostgresqlView.MaxIdleConnection
}

func (input ProductionConfig) GetRedisHost() string {
	return input.Redis.Host
}

func (input ProductionConfig) GetRedisPort() int {
	return convertStringParamToInt("Redis Port", input.Redis.Port)
}

func (input ProductionConfig) GetRedisDB() int {
	return convertStringParamToInt("Redis DB", input.Redis.Db)
}

func (input ProductionConfig) GetRedisSessionDB() int {
	return convertStringParamToInt("Redis DB", input.RedisSession.Db)
}

func (input ProductionConfig) GetRedisPassword() string {
	return input.Redis.Password
}

func (input ProductionConfig) GetRedisTimeout() int {
	return input.Redis.Timeout
}

func (input ProductionConfig) GetRedisRequestVolumeThreshold() int {
	return input.Redis.RequestVolumeThreshold
}

func (input ProductionConfig) GetRedisSleepWindow() int {
	return input.Redis.SleepWindow
}

func (input ProductionConfig) GetRedisErrorPercentThreshold() int {
	return input.Redis.ErrorPercentThreshold
}

func (input ProductionConfig) GetRedisMaxConcurrentRequests() int {
	return input.Redis.MaxConcurrentRequests
}

func (input ProductionConfig) GetClientCredentialsClientID() string {
	return input.ClientCredentials.ClientID
}

func (input ProductionConfig) GetClientCredentialsClientSecret() string {
	return input.ClientCredentials.ClientSecret
}

func (input ProductionConfig) GetClientCredentialsSecretKey() string {
	return input.ClientCredentials.SecretKey
}

func (input ProductionConfig) GetLogFile() []string {
	return input.LogFile
}

func (input ProductionConfig) GetJWTToken() JWTKey {
	return JWTKey{
		JWT:      input.JWTKey.JWT,
		Internal: input.JWTKey.Internal,
	}
}

func (input ProductionConfig) GetLanguageDirectoryPath() string {
	return input.LanguageDirectoryPath
}

func (input ProductionConfig) GetNextracFrontend() NextracFrontend {
	return input.NextracFrontend
}

func (input ProductionConfig) GetAuthenticationServer() AuthenticationServer {
	return input.AuthenticationServer
}

func (input ProductionConfig) GetNexcloudAPI() NexcloudAPI {
	return input.NexcloudAPI
}

func (input ProductionConfig) GetNexdrive() Nexdrive {
	return input.Nexdrive
}
func (input ProductionConfig) GetGrochat() Grochat {
	return input.Grochat
}

func (input ProductionConfig) GetNexmile() Nexmile {
	return input.Nexmile
}

func (input ProductionConfig) GetNexstar() Nexstar {
	return input.Nexstar
}

func (input ProductionConfig) GetNextrade() Nextrade {
	return input.Nextrade
}

func (input ProductionConfig) GetCommonPath() CommonPath {
	return input.CommonPath
}

func (input ProductionConfig) GetAudit() Audit {
	return input.Audit
}

func (input ProductionConfig) GetAlertServer() AlertServer {
	return input.AlertServer
}

func (input ProductionConfig) GetServerPrefixPath() string {
	return input.Server.PrefixPath
}

func (input ProductionConfig) GetAzure() Azure {
	return Azure{
		AccountName: input.Azure.AccountName,
		AccountKey:  input.Azure.AccountKey,
		Host:        input.Azure.Host,
		Suffix:      input.Azure.Suffix,
	}
}

func (input ProductionConfig) GetCDN() CDN {
	return input.Cdn
}

func (input ProductionConfig) GetClientCredentialsAuthUserID() int64 {
	return input.ClientCredentials.AuthUserID
}

func convertStringParamToInt(key string, value string) int {
	intPort, err := strconv.Atoi(value)
	if err != nil {
		logModel := applicationModel.GenerateLogModel("-", "-")
		logModel.Message = "Invalid " + key + " : " + err.Error()
		logModel.Status = 500
		util.LogError(logModel.ToLoggerObject())
		os.Exit(3)
	}
	return intPort
}

func (input ProductionConfig) GetElasticSearchConnectionString() string {
	return input.ElasticSearch.ConnectionStr
}

func (input ProductionConfig) GetServerAutoAddClient() bool {
	return input.Server.AutoAddClient
}

func (input ProductionConfig) GetDataDirectory() DataDirectory {
	return input.DataDirectory
}

func (input ProductionConfig) GetServerLogLevel() string {
	return input.Server.LogLevel
}

func (input ProductionConfig) GetSchedulerStatus() Scheduler {
	return input.Scheduler
}

func (input ProductionConfig) GetMasterData() MasterData {
	return input.MasterData
}

func (input ProductionConfig) GetGenerator() Generator {
	return Generator{
		RootPath: input.Generator.RootPath,
		Path:     input.Generator.Path,
	}
}

func (input ProductionConfig) GetEmailAddress() string {
	return input.Email.Address
}

func (input ProductionConfig) GetEmailPassword() string {
	return input.Email.Password
}

func (input ProductionConfig) GetEmailHost() string {
	return input.Email.Host
}

func (input ProductionConfig) GetEmailPort() int {
	return input.Email.Port
}

func (input ProductionConfig) GetEmail() Email {
	return input.Email
}

func (input ProductionConfig) GetRedmineDBAddress() string {
	return input.Redmine.Address
}

func (input ProductionConfig) GetRedmineDBDefaultSchema() string {
	return input.Redmine.DefaultSchema
}

func (input ProductionConfig) GetRedmineDBMaxOpenConnection() int {
	return input.Redmine.MaxOpenConnection
}

func (input ProductionConfig) GetRedmineDBMaxIdleConnection() int {
	return input.Redmine.MaxIdleConnection
}

func (input ProductionConfig) GetRedmineDBKeyAccess() string {
	return input.Redmine.AccessKey
}

func (input ProductionConfig) GetRedmineAPI() ApiRedmine {
	return input.Redmine.Api
}

func (input ProductionConfig) GetRedmineInfraDBHost() string {
	return input.RedmineInfra.Host
}

func (input ProductionConfig) GetRedmineInfraDBPort() string {
	return input.RedmineInfra.Port
}

func (input ProductionConfig) GetRedmineInfraDBUser() string {
	return input.RedmineInfra.User
}

func (input ProductionConfig) GetRedmineInfraDBPass() string {
	return input.RedmineInfra.Pass
}

func (input ProductionConfig) GetRedmineInfraDBDatabase() string {
	return input.RedmineInfra.Database
}

func (input ProductionConfig) GetRedmineInfraDBKeyAccess() string {
	return input.RedmineInfra.AccessKey
}

func (input ProductionConfig) GetGroChat() Grochat {
	return input.Grochat
}

func (input ProductionConfig) GetGroChatWS() GrochatWS {
	return input.GrochatWS
}