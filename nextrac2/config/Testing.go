package config

type TestingConfig struct {
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
	Grochat               Grochat              `json:"grochat"`
	GrochatWS             GrochatWS            `json:"grochat_ws"`
	Nexmile               Nexmile              `json:"nexmile"`
	Nextrade              Nextrade             `json:"nextrade"`
	Nexstar               Nexstar              `json:"nexstar"`
	CommonPath            CommonPath           `json:"common_path"`
	Audit                 Audit                `json:"audit"`
	AlertServer           AlertServer          `json:"alert_server"`
	Azure                 struct {
		AccountName string `json:"account_name"`
		AccountKey  string `json:"account_key"`
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

func (input TestingConfig) GetServerProtocol() string {
	return input.Server.Protocol
}

func (input TestingConfig) GetServerEthernet() string {
	return input.Server.Ethernet
}

func (input TestingConfig) GetServerAutoAddHost() bool {
	return input.Server.AutoAddHost
}

func (input TestingConfig) GetServerHost() string {
	return input.Server.Host
}

func (input TestingConfig) GetServerPort() int {
	return input.Server.Port
}

func (input TestingConfig) GetServerVersion() string {
	return input.Server.Version
}

func (input TestingConfig) GetServerResourceID() string {
	return input.Server.ResourceID
}

func (input TestingConfig) GetPostgreSQLAddress() string {
	return input.Postgresql.Address
}

func (input TestingConfig) GetPostgreSQLDefaultSchema() string {
	return input.Postgresql.DefaultSchema
}

func (input TestingConfig) GetPostgreSQLMaxOpenConnection() int {
	return input.Postgresql.MaxOpenConnection
}

func (input TestingConfig) GetPostgreSQLMaxIdleConnection() int {
	return input.Postgresql.MaxIdleConnection
}

func (input TestingConfig) GetPostgreSQLAddressView() string {
	return input.PostgresqlView.Address
}

func (input TestingConfig) GetPostgreSQLDefaultSchemaView() string {
	return input.PostgresqlView.DefaultSchema
}

func (input TestingConfig) GetPostgreSQLMaxOpenConnectionView() int {
	return input.PostgresqlView.MaxOpenConnection
}

func (input TestingConfig) GetPostgreSQLMaxIdleConnectionView() int {
	return input.PostgresqlView.MaxIdleConnection
}

func (input TestingConfig) GetRedisHost() string {
	return input.Redis.Host
}

func (input TestingConfig) GetRedisPort() int {
	return convertStringParamToInt("Redis Port", input.Redis.Port)
}

func (input TestingConfig) GetRedisDB() int {
	return convertStringParamToInt("Redis DB", input.Redis.Db)
}

func (input TestingConfig) GetRedisSessionDB() int {
	return convertStringParamToInt("Redis DB", input.RedisSession.Db)
}

func (input TestingConfig) GetRedisPassword() string {
	return input.Redis.Password
}

func (input TestingConfig) GetRedisTimeout() int {
	return input.Redis.Timeout
}

func (input TestingConfig) GetRedisRequestVolumeThreshold() int {
	return input.Redis.RequestVolumeThreshold
}

func (input TestingConfig) GetRedisSleepWindow() int {
	return input.Redis.SleepWindow
}

func (input TestingConfig) GetRedisErrorPercentThreshold() int {
	return input.Redis.ErrorPercentThreshold
}

func (input TestingConfig) GetRedisMaxConcurrentRequests() int {
	return input.Redis.MaxConcurrentRequests
}

func (input TestingConfig) GetClientCredentialsClientID() string {
	return input.ClientCredentials.ClientID
}

func (input TestingConfig) GetClientCredentialsClientSecret() string {
	return input.ClientCredentials.ClientSecret
}

func (input TestingConfig) GetClientCredentialsSecretKey() string {
	return input.ClientCredentials.SecretKey
}

func (input TestingConfig) GetLogFile() []string {
	return input.LogFile
}

func (input TestingConfig) GetJWTToken() JWTKey {
	return input.JWTKey
}

func (input TestingConfig) GetLanguageDirectoryPath() string {
	return input.LanguageDirectoryPath
}

func (input TestingConfig) GetNextracFrontend() NextracFrontend {
	return input.NextracFrontend
}

func (input TestingConfig) GetAuthenticationServer() AuthenticationServer {
	return input.AuthenticationServer
}

func (input TestingConfig) GetNexcloudAPI() NexcloudAPI {
	return input.NexcloudAPI
}

func (input TestingConfig) GetNexdrive() Nexdrive {
	return input.Nexdrive
}

func (input TestingConfig) GetGrochat() Grochat {
	return input.Grochat
}

func (input TestingConfig) GetNexmile() Nexmile {
	return input.Nexmile
}

func (input TestingConfig) GetNextrade() Nextrade {
	return input.Nextrade
}

func (input TestingConfig) GetNexstar() Nexstar {
	return input.Nexstar
}

func (input TestingConfig) GetCommonPath() CommonPath {
	return input.CommonPath
}

func (input TestingConfig) GetAudit() Audit {
	return input.Audit
}

func (input TestingConfig) GetAlertServer() AlertServer {
	return input.AlertServer
}

func (input TestingConfig) GetServerPrefixPath() string {
	return input.Server.PrefixPath
}

func (input TestingConfig) GetAzure() Azure {
	return Azure{
		AccountName: input.Azure.AccountName,
		AccountKey:  input.Azure.AccountKey,
		Host:        input.Azure.Host,
		Suffix:      input.Azure.Suffix,
	}
}

func (input TestingConfig) GetCDN() CDN {
	return input.Cdn
}

func (input TestingConfig) GetClientCredentialsAuthUserID() int64 {
	return input.ClientCredentials.AuthUserID
}

func (input TestingConfig) GetElasticSearchConnectionString() string {
	return input.ElasticSearch.ConnectionStr
}

func (input TestingConfig) GetServerAutoAddClient() bool {
	return input.Server.AutoAddClient
}

func (input TestingConfig) GetDataDirectory() DataDirectory {
	return input.DataDirectory
}

func (input TestingConfig) GetServerLogLevel() string {
	return input.Server.LogLevel
}

func (input TestingConfig) GetSchedulerStatus() Scheduler {
	return input.Scheduler
}

func (input TestingConfig) GetMasterData() MasterData {
	return input.MasterData
}

func (input TestingConfig) GetGenerator() Generator {
	return Generator{
		RootPath: input.Generator.RootPath,
		Path:     input.Generator.Path,
	}
}

func (input TestingConfig) GetRedmineDBAddress() string {
	return input.Redmine.Address
}

func (input TestingConfig) GetRedmineDBDefaultSchema() string {
	return input.Redmine.DefaultSchema
}

func (input TestingConfig) GetRedmineDBMaxOpenConnection() int {
	return input.Redmine.MaxOpenConnection
}

func (input TestingConfig) GetRedmineDBMaxIdleConnection() int {
	return input.Redmine.MaxIdleConnection
}

func (input TestingConfig) GetRedmineDBKeyAccess() string {
	return input.Redmine.AccessKey
}

func (input TestingConfig) GetRedmineAPI() ApiRedmine {
	return input.Redmine.Api
}

func (input TestingConfig) GetRedmineInfraDBHost() string {
	return input.RedmineInfra.Host
}

func (input TestingConfig) GetRedmineInfraDBPort() string {
	return input.RedmineInfra.Port
}

func (input TestingConfig) GetRedmineInfraDBUser() string {
	return input.RedmineInfra.User
}

func (input TestingConfig) GetRedmineInfraDBPass() string {
	return input.RedmineInfra.Pass
}

func (input TestingConfig) GetRedmineInfraDBDatabase() string {
	return input.RedmineInfra.Database
}

func (input TestingConfig) GetRedmineInfraDBKeyAccess() string {
	return input.RedmineInfra.AccessKey
}

func (input TestingConfig) GetGroChat() Grochat {
	return input.Grochat
}

func (input TestingConfig) GetGroChatWS() GrochatWS {
	return input.GrochatWS
}
