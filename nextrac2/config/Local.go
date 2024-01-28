package config

import "time"

type LocalConfig struct {
	Configuration
	Server struct {
		Protocol      string `json:"protocol"`
		Ethernet      string `json:"ethernet"`
		AutoAddHost   bool   `json:"auto_add_host"`
		AutoAddClient bool   `json:"auto_add_client"`
		Host          string `json:"host"`
		Port          int    `json:"port"`
		Version       string `json:"version"`
		ResourceID    string `json:"resource_id"`
		PrefixPath    string `json:"prefix_path"`
		LogLevel      string `json:"log_level"`
	} `json:"server"`
	Postgresql struct {
		Address           string `json:"address"`
		DefaultSchema     string `json:"default_schema"`
		MaxOpenConnection int    `json:"max_open_connection"`
		MaxIdleConnection int    `json:"max_idle_connection"`
	} `json:"postgresql"`
	PostgresqlView struct {
		Address           string `json:"address"`
		DefaultSchema     string `json:"default_schema"`
		MaxOpenConnection int    `json:"max_open_connection"`
		MaxIdleConnection int    `json:"max_idle_connection"`
	} `json:"postgresql_view"`
	Redis struct {
		Host                   string `json:"host"`
		Port                   string `json:"port"`
		Db                     string `json:"db"`
		Password               string `json:"password"`
		Timeout                int    `json:"timeout"`
		RequestVolumeThreshold int    `json:"request_volume_threshold"`
		SleepWindow            int    `json:"sleep_window"`
		ErrorPercentThreshold  int    `json:"error_percent_threshold"`
		MaxConcurrentRequests  int    `json:"max_concurrent_requests"`
	} `json:"redis"`
	RedisSession struct {
		Db string `json:"db"`
	} `json:"redis_session"`
	ClientCredentials struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		SecretKey    string `json:"secret_key"`
		AuthUserID   int64  `json:"auth_user_id"`
	} `json:"client_credentials"`
	LogFile               []string             `json:"log_file"`
	JWTKey                JWTKey               `json:"jwt_key"`
	LanguageDirectoryPath string               `json:"language_directory_path"`
	NextracFrontend       NextracFrontend      `json:"nextrac_frontend"`
	AuthenticationServer  AuthenticationServer `json:"authentication_server"`
	NexcloudAPI           NexcloudAPI          `json:"nexcloud_api"`
	Nexdrive              Nexdrive             `json:"nexdrive"`
	Nexmile               Nexmile              `json:"nexmile"`
	Nextrade              Nextrade             `json:"nextrade"`
	Nexstar               Nexstar              `json:"nexstar"`
	Grochat               Grochat              `json:"grochat"`
	GrochatWS             GrochatWS            `json:"grochat_ws"`
	CommonPath            CommonPath           `json:"common_path"`
	Audit                 Audit                `json:"audit"`
	AlertServer           AlertServer          `json:"alert_server"`
	Azure                 struct {
		AccountName string `json:"account_name"`
		AccountKey  string `json:"account_key"`
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
	Email         Email         `json:"email"`
	Generator     Generator     `json:"generator"`
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

func (input LocalConfig) GetNexmile() Nexmile {
	return input.Nexmile
}

func (input LocalConfig) GetNextrade() Nextrade {
	return input.Nextrade
}

func (input LocalConfig) GetNexstar() Nexstar {
	return input.Nexstar
}

func (input LocalConfig) GetGenerator() Generator {
	return Generator{
		RootPath: input.Generator.RootPath,
		Path:     input.Generator.Path,
	}
}

func (input LocalConfig) GetServerProtocol() string {
	return input.Server.Protocol
}

func (input LocalConfig) GetServerEthernet() string {
	return input.Server.Ethernet
}

func (input LocalConfig) GetServerAutoAddHost() bool {
	return input.Server.AutoAddHost
}

func (input LocalConfig) GetServerHost() string {
	return input.Server.Host
}

func (input LocalConfig) GetServerPort() int {
	return input.Server.Port
}

func (input LocalConfig) GetServerVersion() string {
	return input.Server.Version
}

func (input LocalConfig) GetServerResourceID() string {
	return input.Server.ResourceID
}

func (input LocalConfig) GetPostgreSQLAddress() string {
	return input.Postgresql.Address
}

func (input LocalConfig) GetPostgreSQLDefaultSchema() string {
	return input.Postgresql.DefaultSchema
}

func (input LocalConfig) GetPostgreSQLMaxOpenConnection() int {
	return input.Postgresql.MaxOpenConnection
}

func (input LocalConfig) GetPostgreSQLMaxIdleConnection() int {
	return input.Postgresql.MaxIdleConnection
}

func (input LocalConfig) GetPostgreSQLAddressView() string {
	return input.PostgresqlView.Address
}

func (input LocalConfig) GetPostgreSQLDefaultSchemaView() string {
	return input.PostgresqlView.DefaultSchema
}

func (input LocalConfig) GetPostgreSQLMaxOpenConnectionView() int {
	return input.PostgresqlView.MaxOpenConnection
}

func (input LocalConfig) GetPostgreSQLMaxIdleConnectionView() int {
	return input.PostgresqlView.MaxIdleConnection
}

func (input LocalConfig) GetRedisHost() string {
	return input.Redis.Host
}

func (input LocalConfig) GetRedisPort() int {
	return convertStringParamToInt("Redis Port", input.Redis.Port)
}

func (input LocalConfig) GetRedisDB() int {
	return convertStringParamToInt("Redis DB", input.Redis.Db)
}

func (input LocalConfig) GetRedisSessionDB() int {
	return convertStringParamToInt("Redis DB", input.RedisSession.Db)
}

func (input LocalConfig) GetRedisPassword() string {
	return input.Redis.Password
}

func (input LocalConfig) GetRedisTimeout() int {
	return input.Redis.Timeout
}

func (input LocalConfig) GetRedisRequestVolumeThreshold() int {
	return input.Redis.RequestVolumeThreshold
}

func (input LocalConfig) GetRedisSleepWindow() int {
	return input.Redis.SleepWindow
}

func (input LocalConfig) GetRedisErrorPercentThreshold() int {
	return input.Redis.ErrorPercentThreshold
}

func (input LocalConfig) GetRedisMaxConcurrentRequests() int {
	return input.Redis.MaxConcurrentRequests
}

func (input LocalConfig) GetClientCredentialsClientID() string {
	return input.ClientCredentials.ClientID
}

func (input LocalConfig) GetClientCredentialsClientSecret() string {
	return input.ClientCredentials.ClientSecret
}

func (input LocalConfig) GetClientCredentialsSecretKey() string {
	return input.ClientCredentials.SecretKey
}

func (input LocalConfig) GetLogFile() []string {
	return input.LogFile
}

func (input LocalConfig) GetJWTToken() JWTKey {
	return input.JWTKey
}

func (input LocalConfig) GetLanguageDirectoryPath() string {
	return input.LanguageDirectoryPath
}

func (input LocalConfig) GetNextracFrontend() NextracFrontend {
	return input.NextracFrontend
}

func (input LocalConfig) GetAuthenticationServer() AuthenticationServer {
	return input.AuthenticationServer
}

func (input LocalConfig) GetNexcloudAPI() NexcloudAPI {
	return input.NexcloudAPI
}

func (input LocalConfig) GetNexdrive() Nexdrive {
	return input.Nexdrive
}

func (input LocalConfig) GetGrochat() Grochat {
	return input.Grochat
}

func (input LocalConfig) GetCommonPath() CommonPath {
	return input.CommonPath
}

func (input LocalConfig) GetAudit() Audit {
	return input.Audit
}

func (input LocalConfig) GetAlertServer() AlertServer {
	return input.AlertServer
}

func (input LocalConfig) GetServerPrefixPath() string {
	return input.Server.PrefixPath
}

func (input LocalConfig) GetAzure() Azure {
	return Azure{
		AccountName: input.Azure.AccountName,
		AccountKey:  input.Azure.AccountKey,
		Host:        input.Azure.Host,
		Suffix:      input.Azure.Suffix,
		IsActive:    input.Azure.IsActive,
	}
}

func (input LocalConfig) GetCDN() CDN {
	return input.Cdn
}

func (input LocalConfig) GetClientCredentialsAuthUserID() int64 {
	return input.ClientCredentials.AuthUserID
}

func (input LocalConfig) GetElasticSearchConnectionString() string {
	return input.ElasticSearch.ConnectionStr
}

func (input LocalConfig) GetServerAutoAddClient() bool {
	return input.Server.AutoAddClient
}

func (input LocalConfig) GetDataDirectory() DataDirectory {
	return input.DataDirectory
}

func (input LocalConfig) GetServerLogLevel() string {
	return input.Server.LogLevel
}

func (input LocalConfig) GetSchedulerStatus() Scheduler {
	return input.Scheduler
}

func (input LocalConfig) GetMasterData() MasterData {
	return input.MasterData
}

func (input LocalConfig) GetEmail() Email {
	return Email{
		Address:  input.Email.Address,
		Password: input.Email.Password,
		Host:     input.Email.Host,
		Port:     input.Email.Port,
	}
}

func (input LocalConfig) GetLocalTimezone() LocalTimeZone {
	zone, offset := time.Now().Zone()
	return LocalTimeZone{
		Zone:   zone,
		Offset: offset,
	}
}

func (input LocalConfig) GetEmailAddress() string {
	return input.Email.Address
}

func (input LocalConfig) GetEmailPassword() string {
	return input.Email.Password
}

func (input LocalConfig) GetEmailHost() string {
	return input.Email.Host
}

func (input LocalConfig) GetEmailPort() int {
	return input.Email.Port
}

func (input LocalConfig) GetRedmineDBAddress() string {
	return input.Redmine.Address
}

func (input LocalConfig) GetRedmineDBDefaultSchema() string {
	return input.Redmine.DefaultSchema
}

func (input LocalConfig) GetRedmineDBMaxOpenConnection() int {
	return input.Redmine.MaxOpenConnection
}

func (input LocalConfig) GetRedmineDBMaxIdleConnection() int {
	return input.Redmine.MaxIdleConnection
}

func (input LocalConfig) GetRedmineDBKeyAccess() string {
	return input.Redmine.AccessKey
}

func (input LocalConfig) GetRedmineAPI() ApiRedmine {
	return input.Redmine.Api
}

func (input LocalConfig) GetRedmineInfraDBHost() string {
	return input.RedmineInfra.Host
}

func (input LocalConfig) GetRedmineInfraDBPort() string {
	return input.RedmineInfra.Port
}

func (input LocalConfig) GetRedmineInfraDBUser() string {
	return input.RedmineInfra.User
}

func (input LocalConfig) GetRedmineInfraDBPass() string {
	return input.RedmineInfra.Pass
}

func (input LocalConfig) GetRedmineInfraDBDatabase() string {
	return input.RedmineInfra.Database
}

func (input LocalConfig) GetRedmineInfraDBKeyAccess() string {
	return input.RedmineInfra.AccessKey
}

func (input LocalConfig) GetDiscordLogToken() string {
	return input.DiscordLog.Token
}

func (input LocalConfig) GetDiscordLogChannelId() string {
	return input.DiscordLog.ChannelID
}

func (input LocalConfig) GetGroChat() Grochat {
	return input.Grochat
}

func (input LocalConfig) GetGroChatWS() GrochatWS {
	return input.GrochatWS
}
