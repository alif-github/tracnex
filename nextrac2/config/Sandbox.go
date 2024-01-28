package config

import "time"

type SandboxConfig struct {
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
	DiscordLog struct {
		Token     string `json:"token"`
		ChannelID string `json:"channel_id"`
	} `json:"discord_log"`
}

func (input SandboxConfig) GetServerProtocol() string {
	return input.Server.Protocol
}

func (input SandboxConfig) GetServerEthernet() string {
	return input.Server.Ethernet
}

func (input SandboxConfig) GetServerAutoAddHost() bool {
	return input.Server.AutoAddHost
}

func (input SandboxConfig) GetServerHost() string {
	return input.Server.Host
}

func (input SandboxConfig) GetServerPort() int {
	return convertStringParamToInt("Server Port", input.Server.Port)
}

func (input SandboxConfig) GetServerVersion() string {
	return input.Server.Version
}

func (input SandboxConfig) GetServerResourceID() string {
	return input.Server.ResourceID
}

func (input SandboxConfig) GetPostgreSQLAddress() string {
	return input.Postgresql.Address
}

func (input SandboxConfig) GetPostgreSQLDefaultSchema() string {
	return input.Postgresql.DefaultSchema
}

func (input SandboxConfig) GetPostgreSQLMaxOpenConnection() int {
	return input.Postgresql.MaxOpenConnection
}

func (input SandboxConfig) GetPostgreSQLMaxIdleConnection() int {
	return input.Postgresql.MaxIdleConnection
}

func (input SandboxConfig) GetPostgreSQLAddressView() string {
	return input.PostgresqlView.Address
}

func (input SandboxConfig) GetPostgreSQLDefaultSchemaView() string {
	return input.PostgresqlView.DefaultSchema
}

func (input SandboxConfig) GetPostgreSQLMaxOpenConnectionView() int {
	return input.PostgresqlView.MaxOpenConnection
}

func (input SandboxConfig) GetPostgreSQLMaxIdleConnectionView() int {
	return input.PostgresqlView.MaxIdleConnection
}

func (input SandboxConfig) GetRedisHost() string {
	return input.Redis.Host
}

func (input SandboxConfig) GetRedisPort() int {
	return convertStringParamToInt("Redis Port", input.Redis.Port)
}

func (input SandboxConfig) GetRedisDB() int {
	return convertStringParamToInt("Redis DB", input.Redis.Db)
}

func (input SandboxConfig) GetRedisSessionDB() int {
	return convertStringParamToInt("Redis DB", input.RedisSession.Db)
}

func (input SandboxConfig) GetRedisPassword() string {
	return input.Redis.Password
}

func (input SandboxConfig) GetRedisTimeout() int {
	return input.Redis.Timeout
}

func (input SandboxConfig) GetRedisRequestVolumeThreshold() int {
	return input.Redis.RequestVolumeThreshold
}

func (input SandboxConfig) GetRedisSleepWindow() int {
	return input.Redis.SleepWindow
}

func (input SandboxConfig) GetRedisErrorPercentThreshold() int {
	return input.Redis.ErrorPercentThreshold
}

func (input SandboxConfig) GetRedisMaxConcurrentRequests() int {
	return input.Redis.MaxConcurrentRequests
}

func (input SandboxConfig) GetClientCredentialsClientID() string {
	return input.ClientCredentials.ClientID
}

func (input SandboxConfig) GetClientCredentialsClientSecret() string {
	return input.ClientCredentials.ClientSecret
}

func (input SandboxConfig) GetClientCredentialsSecretKey() string {
	return input.ClientCredentials.SecretKey
}

func (input SandboxConfig) GetLogFile() []string {
	return input.LogFile
}

func (input SandboxConfig) GetJWTToken() JWTKey {
	return JWTKey{
		JWT:      input.JWTKey.JWT,
		Internal: input.JWTKey.Internal,
	}
}

func (input SandboxConfig) GetLanguageDirectoryPath() string {
	return input.LanguageDirectoryPath
}

func (input SandboxConfig) GetNextracFrontend() NextracFrontend {
	return input.NextracFrontend
}

func (input SandboxConfig) GetAuthenticationServer() AuthenticationServer {
	return input.AuthenticationServer
}

func (input SandboxConfig) GetNexcloudAPI() NexcloudAPI {
	return input.NexcloudAPI
}

func (input SandboxConfig) GetNexdrive() Nexdrive {
	return input.Nexdrive
}

func (input SandboxConfig) GetGrochat() Grochat {
	return input.Grochat
}

func (input SandboxConfig) GetNexmile() Nexmile {
	return input.Nexmile
}

func (input SandboxConfig) GetNexstar() Nexstar {
	return input.Nexstar
}

func (input SandboxConfig) GetNextrade() Nextrade {
	return input.Nextrade
}

func (input SandboxConfig) GetCommonPath() CommonPath {
	return input.CommonPath
}

func (input SandboxConfig) GetAudit() Audit {
	return input.Audit
}

func (input SandboxConfig) GetAlertServer() AlertServer {
	return input.AlertServer
}

func (input SandboxConfig) GetServerPrefixPath() string {
	return input.Server.PrefixPath
}

func (input SandboxConfig) GetAzure() Azure {
	return Azure{
		AccountName: input.Azure.AccountName,
		AccountKey:  input.Azure.AccountKey,
		Host:        input.Azure.Host,
		Suffix:      input.Azure.Suffix,
	}
}

func (input SandboxConfig) GetCDN() CDN {
	return input.Cdn
}

func (input SandboxConfig) GetClientCredentialsAuthUserID() int64 {
	return input.ClientCredentials.AuthUserID
}

func (input SandboxConfig) GetElasticSearchConnectionString() string {
	return input.ElasticSearch.ConnectionStr
}

func (input SandboxConfig) GetServerAutoAddClient() bool {
	return input.Server.AutoAddClient
}

func (input SandboxConfig) GetDataDirectory() DataDirectory {
	return input.DataDirectory
}

func (input SandboxConfig) GetServerLogLevel() string {
	return input.Server.LogLevel
}

func (input SandboxConfig) GetSchedulerStatus() Scheduler {
	return input.Scheduler
}

func (input SandboxConfig) GetMasterData() MasterData {
	return input.MasterData
}

func (input SandboxConfig) GetGenerator() Generator {
	return Generator{
		RootPath: input.Generator.RootPath,
		Path:     input.Generator.Path,
	}
}

func (input SandboxConfig) GetEmailAddress() string {
	return input.Email.Address
}

func (input SandboxConfig) GetEmailPassword() string {
	return input.Email.Password
}

func (input SandboxConfig) GetEmailHost() string {
	return input.Email.Host
}

func (input SandboxConfig) GetEmailPort() int {
	return input.Email.Port
}

func (input SandboxConfig) GetEmail() Email {
	return Email{
		Address:  input.Email.Address,
		Password: input.Email.Password,
		Host:     input.Email.Host,
		Port:     input.Email.Port,
	}
}

func (input SandboxConfig) GetLocalTimezone() LocalTimeZone {
	zone, offset := time.Now().Zone()
	return LocalTimeZone{
		Zone:   zone,
		Offset: offset,
	}
}

func (input SandboxConfig) GetRedmineDBAddress() string {
	return input.Redmine.Address
}

func (input SandboxConfig) GetRedmineDBDefaultSchema() string {
	return input.Redmine.DefaultSchema
}

func (input SandboxConfig) GetRedmineDBMaxOpenConnection() int {
	return input.Redmine.MaxOpenConnection
}

func (input SandboxConfig) GetRedmineDBMaxIdleConnection() int {
	return input.Redmine.MaxIdleConnection
}

func (input SandboxConfig) GetRedmineDBKeyAccess() string {
	return input.Redmine.AccessKey
}

func (input SandboxConfig) GetRedmineAPI() ApiRedmine {
	return input.Redmine.Api
}

func (input SandboxConfig) GetRedmineInfraDBHost() string {
	return input.RedmineInfra.Host
}

func (input SandboxConfig) GetRedmineInfraDBPort() string {
	return input.RedmineInfra.Port
}

func (input SandboxConfig) GetRedmineInfraDBUser() string {
	return input.RedmineInfra.User
}

func (input SandboxConfig) GetRedmineInfraDBPass() string {
	return input.RedmineInfra.Pass
}

func (input SandboxConfig) GetRedmineInfraDBDatabase() string {
	return input.RedmineInfra.Database
}

func (input SandboxConfig) GetRedmineInfraDBKeyAccess() string {
	return input.RedmineInfra.AccessKey
}

func (input SandboxConfig) GetDiscordLogToken() string {
	return input.DiscordLog.Token
}

func (input SandboxConfig) GetDiscordLogChannelId() string {
	return input.DiscordLog.ChannelID
}

func (input SandboxConfig) GetGroChat() Grochat {
	return input.Grochat
}

func (input SandboxConfig) GetGroChatWS() GrochatWS {
	return input.GrochatWS
}
