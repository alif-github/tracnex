package config

import "time"

type DevelopmentConfig struct {
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
	Nextrade              Nextrade             `json:"nextrade"`
	Nexstar               Nexstar              `json:"nexstar"`
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

func (input DevelopmentConfig) GetServerProtocol() string {
	return input.Server.Protocol
}

func (input DevelopmentConfig) GetServerEthernet() string {
	return input.Server.Ethernet
}

func (input DevelopmentConfig) GetServerAutoAddHost() bool {
	return input.Server.AutoAddHost
}

func (input DevelopmentConfig) GetServerHost() string {
	return input.Server.Host
}

func (input DevelopmentConfig) GetServerPort() int {
	return convertStringParamToInt("Server Port", input.Server.Port)
}

func (input DevelopmentConfig) GetServerVersion() string {
	return input.Server.Version
}

func (input DevelopmentConfig) GetServerResourceID() string {
	return input.Server.ResourceID
}

func (input DevelopmentConfig) GetPostgreSQLAddress() string {
	return input.Postgresql.Address
}

func (input DevelopmentConfig) GetPostgreSQLDefaultSchema() string {
	return input.Postgresql.DefaultSchema
}

func (input DevelopmentConfig) GetPostgreSQLMaxOpenConnection() int {
	return input.Postgresql.MaxOpenConnection
}

func (input DevelopmentConfig) GetPostgreSQLMaxIdleConnection() int {
	return input.Postgresql.MaxIdleConnection
}

func (input DevelopmentConfig) GetPostgreSQLAddressView() string {
	return input.PostgresqlView.Address
}

func (input DevelopmentConfig) GetPostgreSQLDefaultSchemaView() string {
	return input.PostgresqlView.DefaultSchema
}

func (input DevelopmentConfig) GetPostgreSQLMaxOpenConnectionView() int {
	return input.PostgresqlView.MaxOpenConnection
}

func (input DevelopmentConfig) GetPostgreSQLMaxIdleConnectionView() int {
	return input.PostgresqlView.MaxIdleConnection
}

func (input DevelopmentConfig) GetRedisHost() string {
	return input.Redis.Host
}

func (input DevelopmentConfig) GetRedisPort() int {
	return convertStringParamToInt("Redis Port", input.Redis.Port)
}

func (input DevelopmentConfig) GetRedisDB() int {
	return convertStringParamToInt("Redis DB", input.Redis.Db)
}

func (input DevelopmentConfig) GetRedisSessionDB() int {
	return convertStringParamToInt("Redis DB", input.RedisSession.Db)
}

func (input DevelopmentConfig) GetRedisPassword() string {
	return input.Redis.Password
}

func (input DevelopmentConfig) GetRedisTimeout() int {
	return input.Redis.Timeout
}

func (input DevelopmentConfig) GetRedisRequestVolumeThreshold() int {
	return input.Redis.RequestVolumeThreshold
}

func (input DevelopmentConfig) GetRedisSleepWindow() int {
	return input.Redis.SleepWindow
}

func (input DevelopmentConfig) GetRedisErrorPercentThreshold() int {
	return input.Redis.ErrorPercentThreshold
}

func (input DevelopmentConfig) GetRedisMaxConcurrentRequests() int {
	return input.Redis.MaxConcurrentRequests
}

func (input DevelopmentConfig) GetClientCredentialsClientID() string {
	return input.ClientCredentials.ClientID
}

func (input DevelopmentConfig) GetClientCredentialsClientSecret() string {
	return input.ClientCredentials.ClientSecret
}

func (input DevelopmentConfig) GetClientCredentialsSecretKey() string {
	return input.ClientCredentials.SecretKey
}

func (input DevelopmentConfig) GetLogFile() []string {
	return input.LogFile
}

func (input DevelopmentConfig) GetJWTToken() JWTKey {
	return JWTKey{
		JWT:      input.JWTKey.JWT,
		Internal: input.JWTKey.Internal,
	}
}

func (input DevelopmentConfig) GetLanguageDirectoryPath() string {
	return input.LanguageDirectoryPath
}

func (input DevelopmentConfig) GetNextracFrontend() NextracFrontend {
	return input.NextracFrontend
}

func (input DevelopmentConfig) GetAuthenticationServer() AuthenticationServer {
	return input.AuthenticationServer
}

func (input DevelopmentConfig) GetNexcloudAPI() NexcloudAPI {
	return input.NexcloudAPI
}

func (input DevelopmentConfig) GetNexdrive() Nexdrive {
	return input.Nexdrive
}

func (input DevelopmentConfig) GetGrochat() Grochat {
	return input.Grochat
}

func (input DevelopmentConfig) GetNexmile() Nexmile {
	return input.Nexmile
}

func (input DevelopmentConfig) GetNextrade() Nextrade {
	return input.Nextrade
}

func (input DevelopmentConfig) GetNexstar() Nexstar {
	return input.Nexstar
}

func (input DevelopmentConfig) GetCommonPath() CommonPath {
	return input.CommonPath
}

func (input DevelopmentConfig) GetAudit() Audit {
	return input.Audit
}

func (input DevelopmentConfig) GetAlertServer() AlertServer {
	return input.AlertServer
}

func (input DevelopmentConfig) GetServerPrefixPath() string {
	return input.Server.PrefixPath
}

func (input DevelopmentConfig) GetAzure() Azure {
	return Azure{
		AccountName: input.Azure.AccountName,
		AccountKey:  input.Azure.AccountKey,
		Host:        input.Azure.Host,
		Suffix:      input.Azure.Suffix,
		IsActive:    input.Azure.IsActive,
	}
}

func (input DevelopmentConfig) GetCDN() CDN {
	return input.Cdn
}

func (input DevelopmentConfig) GetClientCredentialsAuthUserID() int64 {
	return input.ClientCredentials.AuthUserID
}

func (input DevelopmentConfig) GetElasticSearchConnectionString() string {
	return input.ElasticSearch.ConnectionStr
}

func (input DevelopmentConfig) GetServerAutoAddClient() bool {
	return input.Server.AutoAddClient
}

func (input DevelopmentConfig) GetDataDirectory() DataDirectory {
	return input.DataDirectory
}

func (input DevelopmentConfig) GetServerLogLevel() string {
	return input.Server.LogLevel
}

func (input DevelopmentConfig) GetSchedulerStatus() Scheduler {
	return input.Scheduler
}

func (input DevelopmentConfig) GetMasterData() MasterData {
	return input.MasterData
}

func (input DevelopmentConfig) GetGenerator() Generator {
	return Generator{
		RootPath: input.Generator.RootPath,
		Path:     input.Generator.Path,
	}
}

func (input DevelopmentConfig) GetEmail() Email {
	return Email{
		Address:  input.Email.Address,
		Password: input.Email.Password,
		Host:     input.Email.Host,
		Port:     input.Email.Port,
	}
}

func (input DevelopmentConfig) GetLocalTimezone() LocalTimeZone {
	zone, offset := time.Now().Zone()
	return LocalTimeZone{
		Zone:   zone,
		Offset: offset,
	}
}

func (input DevelopmentConfig) GetEmailAddress() string {
	return input.Email.Address
}

func (input DevelopmentConfig) GetEmailPassword() string {
	return input.Email.Password
}

func (input DevelopmentConfig) GetEmailHost() string {
	return input.Email.Host
}

func (input DevelopmentConfig) GetEmailPort() int {
	return input.Email.Port
}

func (input DevelopmentConfig) GetRedmineDBAddress() string {
	return input.Redmine.Address
}

func (input DevelopmentConfig) GetRedmineDBDefaultSchema() string {
	return input.Redmine.DefaultSchema
}

func (input DevelopmentConfig) GetRedmineDBMaxOpenConnection() int {
	return input.Redmine.MaxOpenConnection
}

func (input DevelopmentConfig) GetRedmineDBMaxIdleConnection() int {
	return input.Redmine.MaxIdleConnection
}

func (input DevelopmentConfig) GetRedmineDBKeyAccess() string {
	return input.Redmine.AccessKey
}

func (input DevelopmentConfig) GetRedmineAPI() ApiRedmine {
	return input.Redmine.Api
}

func (input DevelopmentConfig) GetRedmineInfraDBHost() string {
	return input.RedmineInfra.Host
}

func (input DevelopmentConfig) GetRedmineInfraDBPort() string {
	return input.RedmineInfra.Port
}

func (input DevelopmentConfig) GetRedmineInfraDBUser() string {
	return input.RedmineInfra.User
}

func (input DevelopmentConfig) GetRedmineInfraDBPass() string {
	return input.RedmineInfra.Pass
}

func (input DevelopmentConfig) GetRedmineInfraDBDatabase() string {
	return input.RedmineInfra.Database
}

func (input DevelopmentConfig) GetRedmineInfraDBKeyAccess() string {
	return input.RedmineInfra.AccessKey
}

func (input DevelopmentConfig) GetDiscordLogToken() string {
	return input.DiscordLog.Token
}

func (input DevelopmentConfig) GetDiscordLogChannelId() string {
	return input.DiscordLog.ChannelID
}

func (input DevelopmentConfig) GetGroChat() Grochat {
	return input.Grochat
}

func (input DevelopmentConfig) GetGroChatWS() GrochatWS {
	return input.GrochatWS
}
