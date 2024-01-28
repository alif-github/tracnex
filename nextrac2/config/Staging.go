package config

import "time"

type StagingConfig struct {
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
		IsActive    bool   `json:"is_active"`
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
	DiscordLog struct {
		Token     string `json:"token"`
		ChannelID string `json:"channel_id"`
	} `json:"discord_log"`
}

func (input StagingConfig) GetServerProtocol() string {
	return input.Server.Protocol
}

func (input StagingConfig) GetServerEthernet() string {
	return input.Server.Ethernet
}

func (input StagingConfig) GetServerAutoAddHost() bool {
	return input.Server.AutoAddHost
}

func (input StagingConfig) GetServerHost() string {
	return input.Server.Host
}

func (input StagingConfig) GetServerPort() int {
	return convertStringParamToInt("Server Port", input.Server.Port)
}

func (input StagingConfig) GetServerVersion() string {
	return input.Server.Version
}

func (input StagingConfig) GetServerResourceID() string {
	return input.Server.ResourceID
}

func (input StagingConfig) GetPostgreSQLAddress() string {
	return input.Postgresql.Address
}

func (input StagingConfig) GetPostgreSQLDefaultSchema() string {
	return input.Postgresql.DefaultSchema
}

func (input StagingConfig) GetPostgreSQLMaxOpenConnection() int {
	return input.Postgresql.MaxOpenConnection
}

func (input StagingConfig) GetPostgreSQLMaxIdleConnection() int {
	return input.Postgresql.MaxIdleConnection
}

func (input StagingConfig) GetPostgreSQLAddressView() string {
	return input.PostgresqlView.Address
}

func (input StagingConfig) GetPostgreSQLDefaultSchemaView() string {
	return input.PostgresqlView.DefaultSchema
}

func (input StagingConfig) GetPostgreSQLMaxOpenConnectionView() int {
	return input.PostgresqlView.MaxOpenConnection
}

func (input StagingConfig) GetPostgreSQLMaxIdleConnectionView() int {
	return input.PostgresqlView.MaxIdleConnection
}

func (input StagingConfig) GetRedisHost() string {
	return input.Redis.Host
}

func (input StagingConfig) GetRedisPort() int {
	return convertStringParamToInt("Redis Port", input.Redis.Port)
}

func (input StagingConfig) GetRedisDB() int {
	return convertStringParamToInt("Redis DB", input.Redis.Db)
}

func (input StagingConfig) GetRedisSessionDB() int {
	return convertStringParamToInt("Redis DB", input.RedisSession.Db)
}

func (input StagingConfig) GetRedisPassword() string {
	return input.Redis.Password
}

func (input StagingConfig) GetRedisTimeout() int {
	return input.Redis.Timeout
}

func (input StagingConfig) GetRedisRequestVolumeThreshold() int {
	return input.Redis.RequestVolumeThreshold
}

func (input StagingConfig) GetRedisSleepWindow() int {
	return input.Redis.SleepWindow
}

func (input StagingConfig) GetRedisErrorPercentThreshold() int {
	return input.Redis.ErrorPercentThreshold
}

func (input StagingConfig) GetRedisMaxConcurrentRequests() int {
	return input.Redis.MaxConcurrentRequests
}

func (input StagingConfig) GetClientCredentialsClientID() string {
	return input.ClientCredentials.ClientID
}

func (input StagingConfig) GetClientCredentialsClientSecret() string {
	return input.ClientCredentials.ClientSecret
}

func (input StagingConfig) GetClientCredentialsSecretKey() string {
	return input.ClientCredentials.SecretKey
}

func (input StagingConfig) GetLogFile() []string {
	return input.LogFile
}

func (input StagingConfig) GetJWTToken() JWTKey {
	return JWTKey{
		JWT:      input.JWTKey.JWT,
		Internal: input.JWTKey.Internal,
	}
}

func (input StagingConfig) GetLanguageDirectoryPath() string {
	return input.LanguageDirectoryPath
}

func (input StagingConfig) GetNextracFrontend() NextracFrontend {
	return input.NextracFrontend
}

func (input StagingConfig) GetAuthenticationServer() AuthenticationServer {
	return input.AuthenticationServer
}

func (input StagingConfig) GetNexcloudAPI() NexcloudAPI {
	return input.NexcloudAPI
}

func (input StagingConfig) GetNexdrive() Nexdrive {
	return input.Nexdrive
}

func (input StagingConfig) GetGrochat() Grochat {
	return input.Grochat
}

func (input StagingConfig) GetNexmile() Nexmile {
	return input.Nexmile
}

func (input StagingConfig) GetNexstar() Nexstar {
	return input.Nexstar
}

func (input StagingConfig) GetNextrade() Nextrade {
	return input.Nextrade
}

func (input StagingConfig) GetCommonPath() CommonPath {
	return input.CommonPath
}

func (input StagingConfig) GetAudit() Audit {
	return input.Audit
}

func (input StagingConfig) GetAlertServer() AlertServer {
	return input.AlertServer
}

func (input StagingConfig) GetServerPrefixPath() string {
	return input.Server.PrefixPath
}

func (input StagingConfig) GetAzure() Azure {
	return Azure{
		AccountName: input.Azure.AccountName,
		AccountKey:  input.Azure.AccountKey,
		Host:        input.Azure.Host,
		Suffix:      input.Azure.Suffix,
	}
}

func (input StagingConfig) GetCDN() CDN {
	return input.Cdn
}

func (input StagingConfig) GetClientCredentialsAuthUserID() int64 {
	return input.ClientCredentials.AuthUserID
}

func (input StagingConfig) GetElasticSearchConnectionString() string {
	return input.ElasticSearch.ConnectionStr
}

func (input StagingConfig) GetServerAutoAddClient() bool {
	return input.Server.AutoAddClient
}

func (input StagingConfig) GetDataDirectory() DataDirectory {
	return input.DataDirectory
}

func (input StagingConfig) GetServerLogLevel() string {
	return input.Server.LogLevel
}

func (input StagingConfig) GetSchedulerStatus() Scheduler {
	return input.Scheduler
}

func (input StagingConfig) GetMasterData() MasterData {
	return input.MasterData
}

func (input StagingConfig) GetGenerator() Generator {
	return Generator{
		RootPath: input.Generator.RootPath,
		Path:     input.Generator.Path,
	}
}

func (input StagingConfig) GetEmail() Email {
	return Email{
		Address:  input.Email.Address,
		Password: input.Email.Password,
		Host:     input.Email.Host,
		Port:     input.Email.Port,
	}
}

func (input StagingConfig) GetLocalTimezone() LocalTimeZone {
	zone, offset := time.Now().Zone()
	return LocalTimeZone{
		Zone:   zone,
		Offset: offset,
	}
}

func (input StagingConfig) GetEmailAddress() string {
	return input.Email.Address
}

func (input StagingConfig) GetEmailPassword() string {
	return input.Email.Password
}

func (input StagingConfig) GetEmailHost() string {
	return input.Email.Host
}

func (input StagingConfig) GetEmailPort() int {
	return input.Email.Port
}

func (input StagingConfig) GetRedmineDBAddress() string {
	return input.Redmine.Address
}

func (input StagingConfig) GetRedmineDBDefaultSchema() string {
	return input.Redmine.DefaultSchema
}

func (input StagingConfig) GetRedmineDBMaxOpenConnection() int {
	return input.Redmine.MaxOpenConnection
}

func (input StagingConfig) GetRedmineDBMaxIdleConnection() int {
	return input.Redmine.MaxIdleConnection
}

func (input StagingConfig) GetRedmineDBKeyAccess() string {
	return input.Redmine.AccessKey
}

func (input StagingConfig) GetRedmineAPI() ApiRedmine {
	return input.Redmine.Api
}

func (input StagingConfig) GetRedmineInfraDBHost() string {
	return input.RedmineInfra.Host
}

func (input StagingConfig) GetRedmineInfraDBPort() string {
	return input.RedmineInfra.Port
}

func (input StagingConfig) GetRedmineInfraDBUser() string {
	return input.RedmineInfra.User
}

func (input StagingConfig) GetRedmineInfraDBPass() string {
	return input.RedmineInfra.Pass
}

func (input StagingConfig) GetRedmineInfraDBDatabase() string {
	return input.RedmineInfra.Database
}

func (input StagingConfig) GetRedmineInfraDBKeyAccess() string {
	return input.RedmineInfra.AccessKey
}

func (input StagingConfig) GetDiscordLogToken() string {
	return input.DiscordLog.Token
}

func (input StagingConfig) GetDiscordLogChannelId() string {
	return input.DiscordLog.ChannelID
}

func (input StagingConfig) GetGroChat() Grochat {
	return input.Grochat
}

func (input StagingConfig) GetGroChatWS() GrochatWS {
	return input.GrochatWS
}
