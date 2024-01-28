package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"github.com/tkanos/gonfig"
	"os"
)

var ApplicationConfiguration Configuration

type Configuration interface {
	GetServerHost() string
	GetServerProtocol() string
	GetServerEthernet() string
	GetServerAutoAddHost() bool
	GetServerPort() int
	GetServerVersion() string
	GetServerResourceID() string
	GetServerPrefixPath() string
	GetPostgreSQLAddress() string
	GetPostgreSQLDefaultSchema() string
	GetPostgreSQLMaxOpenConnection() int
	GetPostgreSQLMaxIdleConnection() int
	GetPostgreSQLAddressView() string
	GetPostgreSQLDefaultSchemaView() string
	GetPostgreSQLMaxOpenConnectionView() int
	GetPostgreSQLMaxIdleConnectionView() int
	GetRedisHost() string
	GetRedisPort() int
	GetRedisDB() int
	GetRedisSessionDB() int
	GetRedisPassword() string
	GetRedisTimeout() int
	GetRedisRequestVolumeThreshold() int
	GetRedisSleepWindow() int
	GetRedisErrorPercentThreshold() int
	GetRedisMaxConcurrentRequests() int
	GetClientCredentialsClientID() string
	GetClientCredentialsClientSecret() string
	GetClientCredentialsSecretKey() string
	GetClientCredentialsAuthUserID() int64
	GetLogFile() []string
	GetJWTToken() JWTKey
	GetLanguageDirectoryPath() string
	GetNextracFrontend() NextracFrontend
	GetAuthenticationServer() AuthenticationServer
	GetNexcloudAPI() NexcloudAPI
	GetNexdrive() Nexdrive
	GetGrochat() Grochat
	GetNexmile() Nexmile
	GetNextrade() Nextrade
	GetNexstar() Nexstar
	GetCommonPath() CommonPath
	GetAudit() Audit
	GetAlertServer() AlertServer
	GetAzure() Azure
	GetCDN() CDN
	GetElasticSearchConnectionString() string
	GetServerAutoAddClient() bool
	GetDataDirectory() DataDirectory
	GetServerLogLevel() string
	GetSchedulerStatus() Scheduler
	GetMasterData() MasterData
	GetEmailAddress() string
	GetEmailPassword() string
	GetEmailHost() string
	GetEmailPort() int
	GetEmail() Email
	GetGenerator() Generator
	GetLocalTimezone() LocalTimeZone
	GetRedmineDBAddress() string
	GetRedmineDBDefaultSchema() string
	GetRedmineDBMaxOpenConnection() int
	GetRedmineDBMaxIdleConnection() int
	GetRedmineDBKeyAccess() string
	GetRedmineAPI() ApiRedmine
	GetRedmineInfraDBHost() string
	GetRedmineInfraDBPort() string
	GetRedmineInfraDBUser() string
	GetRedmineInfraDBPass() string
	GetRedmineInfraDBDatabase() string
	GetRedmineInfraDBKeyAccess() string
	GetDiscordLogToken() string
	GetDiscordLogChannelId() string
	GetGroChat() Grochat
	GetGroChatWS() GrochatWS
}

type AuthenticationServer struct {
	Host         string                     `json:"host"`
	PathRedirect AuthenticationPathRedirect `json:"path_redirect"`
}

type NexcloudAPI struct {
	Host         string                  `json:"host"`
	PathRedirect NexcloudAPIPathRedirect `json:"path_redirect"`
}

type Nexdrive struct {
	Host         string                  `json:"host"`
	PathRedirect NexdriveAPIPathRedirect `json:"path_redirect"`
}

type Grochat struct {
	Host         string                 `json:"host"`
	PathRedirect GrochatAPIPathRedirect `json:"path_redirect"`
}

type GrochatWS struct {
	Host         string                `json:"host"`
	PathRedirect GrochatWSPathRedirect `json:"path_redirect"`
	User         struct {
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"user"`
}

type Nexmile struct {
	Host         string                 `json:"host"`
	PathRedirect NexmileAPIPathRedirect `json:"path_redirect"`
}

type Nextrade struct {
	Host         string                  `json:"host"`
	PathRedirect NextradeAPIPathRedirect `json:"path_redirect"`
}

type Nexstar struct {
	Host         string                 `json:"host"`
	PathRedirect NexstarAPIPathRedirect `json:"path_redirect"`
}

type AuthenticationPathRedirect struct {
	CheckToken        string `json:"check_token"`
	AddResourceClient string `json:"add_resource_client"`
	CheckUser         string `json:"check_user"`
	GetUser           string `json:"get_user"`
	InternalUser      struct {
		CrudUser       string `json:"crud_user"`
		CheckUser      string `json:"check_user"`
		Initiate       string `json:"initiate"`
		ChangePassword string `json:"change_password"`
		Forget         struct {
			Phone          string `json:"phone"`
			Email          string `json:"email"`
			ChangePassword struct {
				Email string `json:"email"`
			} `json:"change_password"`
		} `json:"forget"`
		Activation struct {
			Phone string `json:"phone"`
			Email string `json:"email"`
		} `json:"activation"`
		ResendActivation struct {
			Phone string `json:"phone"`
			Email string `json:"email"`
		} `json:"resend_activation"`
	} `json:"internal_user"`
	InternalClient struct {
		CrudClient       string `json:"crud_client"`
		AddResourceAdmin string `json:"add_resource_admin"`
		CheckClientUser  string `json:"check_client_user"`
	} `json:"internal_client"`
	Authorize string `json:"authorize"`
	Verify    string `json:"verify"`
	Token     string `json:"token"`
	Logout    string `json:"logout"`
}

type NexcloudAPIPathRedirect struct {
	AddResourceClient string `json:"add_resource_client"`
	CrudClient        string `json:"crud_client"`
}

type NexdriveAPIPathRedirect struct {
	AddResourceClient string `json:"add_resource_client"`
}

type GrochatAPIPathRedirect struct {
	SendMessage    string `json:"send_message"`
	Authentication string `json:"authentication"`
	UserDetail     string `json:"user_detail"`
	Login          string `json:"login"`
	SignId         string `json:"sign_id"`
	SendInvitation string `json:"send_invitation"`
	PWAInvitation  string `json:"pwa_invitation"`
}

type GrochatWSPathRedirect struct {
	WS string `json:"ws"`
}

type NexmileAPIPathRedirect struct {
	ActivationUser string `json:"activation_user"`
}

type NextradeAPIPathRedirect struct {
	ActivationUser string `json:"activation_user"`
}

type NexstarAPIPathRedirect struct {
	ActivationUser string `json:"activation_user"`
}

type NextracFrontend struct {
	Host         string              `json:"host"`
	PathRedirect NextracPathRedirect `json:"path_redirect"`
}

type NextracPathRedirect struct {
	ResetPasswordPath   string `json:"reset_password_path"`
	VerifyUserPath      string `json:"verify_user_path"`
	TodoList            string `json:"todo_list"`
	AccountRegistration string `json:"account_registration"`
	Invitation			string `json:"invitation"`
	ProfileCompletion	string `json:"profile_completion"`
}

type CommonPath struct {
	ResourceClients    string `json:"resource_clients"`
	NotifyDeletedToken string `json:"notify_deleted_token"`
}

type Audit struct {
	IsActive       bool     `json:"is_active"`
	ListSecretData []string `json:"list_secret_data"`
}

type JWTKey struct {
	JWT      string `json:"jwt"`
	Internal string `json:"internal"`
}

type AlertServer struct {
	Host         string              `json:"host"`
	PathRedirect NextracPathRedirect `json:"path_redirect"`
}

type AlertServerPathRedirect struct {
	Alert string `json:"alert"`
}

type Azure struct {
	AccountName string
	AccountKey  string
	Host        string
	Suffix      string
	IsActive    bool
}

type CDN struct {
	RootPath string `json:"root_path"`
	Host     string `json:"host"`
	Suffix   string `json:"suffix"`
}

type Scheduler struct {
	IsActive  bool `json:"is_active"`
	IsTesting bool `json:"is_testing"`
}

type MasterData struct {
	Host         string                    `json:"host"`
	PathRedirect MasterDataPathRedirectApi `json:"path_redirect"`
}

type MasterDataPathRedirectApi struct {
	PersonTitle    PersonTitleMasterData    `json:"person_title"`
	Province       ProvinceMasterData       `json:"province"`
	District       DistrictMasterData       `json:"district"`
	Position       PositionMasterData       `json:"position"`
	CompanyTitle   CompanyTitleMasterData   `json:"company_title"`
	CompanyProfile CompanyProfileMasterData `json:"company_profile"`
	PersonProfile  PersonProfileMasterData  `json:"person_profile"`
	ContactPerson  ContactPersonMasterData  `json:"contact_person"`
	SubDistrict    SubDistrictMasterData    `json:"sub_district"`
	UrbanVillage   UrbanVillageMasterData   `json:"urban_village"`
	PostalCode     PostalCodeMasterData     `json:"postal_code"`
	Country        CountryMasterData        `json:"country"`
}

type PersonTitleMasterData struct {
	GetList string `json:"get_list"`
	View    string `json:"view"`
}

type ProvinceMasterData struct {
	GetList string `json:"get_list"`
	View    string `json:"view"`
}

type DistrictMasterData struct {
	GetList string `json:"get_list"`
	View    string `json:"view"`
}

type PositionMasterData struct {
	GetList string `json:"get_list"`
	View    string `json:"view"`
}

type CompanyTitleMasterData struct {
	GetList string `json:"get_list"`
	View    string `json:"view"`
}

type CompanyProfileMasterData struct {
	GetList  string `json:"get_list"`
	View     string `json:"view"`
	Validate string `json:"validate"`
}

type PersonProfileMasterData struct {
	GetList  string `json:"get_list"`
	View     string `json:"view"`
	Validate string `json:"validate"`
}

type ContactPersonMasterData struct {
	GetList  string `json:"get_list"`
	BaseUrl  string `json:"base_url"`
	Validate string `json:"validate"`
}

type SubDistrictMasterData struct {
	GetList string `json:"get_list"`
	Count   string `json:"count"`
	View    string `json:"view"`
}

type UrbanVillageMasterData struct {
	GetList string `json:"get_list"`
	Count   string `json:"count"`
	View    string `json:"view"`
}

type PostalCodeMasterData struct {
	GetList string `json:"get_list"`
	Count   string `json:"count"`
	View    string `json:"view"`
}

type CountryMasterData struct {
	GetList string `json:"get_list"`
	Count   string `json:"count"`
}

func GenerateConfiguration(arguments string) {
	var (
		err      error
		enviName = os.Getenv("nextrac2CoreConfig")
		strError = "Error get config 2 -> "
	)

	fmt.Println(fmt.Sprintf(`Argument [%s]`, arguments))
	if arguments == "sandbox" {
		var (
			temp     = SandboxConfig{}
			filename = "config_sandbox.json"
		)

		err = gonfig.GetConf(enviName+"/"+filename, &temp)
		if err != nil {
			fmt.Print("Error get config sandbox -> ", err)
			os.Exit(2)
		}

		err = envconfig.Process(enviName+"/"+filename, &temp)
		if err != nil {
			fmt.Print(strError, err)
			os.Exit(2)
		}

		ApplicationConfiguration = &temp
	} else if arguments == "staging" {
		var (
			temp     = StagingConfig{}
			filename = "config_staging.json"
		)

		err = gonfig.GetConf(enviName+"/"+filename, &temp)
		if err != nil {
			fmt.Print("Error get config staging -> ", err)
			os.Exit(2)
		}

		err = envconfig.Process(enviName+"/"+filename, &temp)
		if err != nil {
			fmt.Print(strError, err)
			os.Exit(2)
		}

		ApplicationConfiguration = &temp
	} else if arguments == "production" {
		var (
			temp     = ProductionConfig{}
			filename = "config_production.json"
		)

		err = gonfig.GetConf(enviName+"/"+filename, &temp)
		if err != nil {
			fmt.Print("Error get config production -> ", err)
			os.Exit(2)
		}

		err = envconfig.Process(enviName+"/"+filename, &temp)
		if err != nil {
			fmt.Print(strError, err)
			os.Exit(2)
		}

		ApplicationConfiguration = &temp
	} else if arguments == "testing" {
		var (
			temp     = TestingConfig{}
			filename = "config_testing.json"
		)

		err = gonfig.GetConf(enviName+"/"+filename, &temp)
		if err != nil {
			fmt.Print("Error get config testing -> ", err)
			os.Exit(2)
		}

		err = envconfig.Process(enviName+"/"+filename, &temp)
		if err != nil {
			fmt.Print(strError, err)
			os.Exit(2)
		}

		ApplicationConfiguration = &temp
	} else if arguments == "development" {
		var (
			temp     = DevelopmentConfig{}
			filename = "config_development.json"
		)

		err = gonfig.GetConf(enviName+"/"+filename, &temp)
		if err != nil {
			fmt.Print("Error get config development -> ", err)
			os.Exit(2)
		}

		err = envconfig.Process(enviName+"/"+filename, &temp)
		if err != nil {
			fmt.Print(strError, err)
			os.Exit(2)
		}

		ApplicationConfiguration = &temp
	} else {
		var temp = LocalConfig{}
		err = gonfig.GetConf(enviName+"/config_local.json", &temp)
		if err != nil {
			fmt.Print("Error get config local -> ", err)
			os.Exit(2)
		}

		ApplicationConfiguration = &temp
	}

	if err != nil {
		fmt.Print("Error get config -> ", err)
		os.Exit(2)
	}
}

type DataDirectory struct {
	BaseDirectoryPath string `json:"base_directory_path"`
	ImportPath        string `json:"import_path"`
	CustomerPath      string `json:"customer_path"`
	DonePath          string `json:"done_path"`
	ProcessPath       string `json:"process_path"`
	FailedPath        string `json:"failed_path"`
	InboundPath       string `json:"inbound_path"`
	KeyFile           string `json:"key_file"`
	KeyContent        string `json:"key_content"`
	Template          string `json:"template"`
	Backlog           string `json:"backlog"`
	ReportLeavePath   string `json:"report_leave_path"`
}

type Generator struct {
	RootPath string `json:"root_path"`
	Path     string `json:"path"`
}

type Email struct {
	Address  string `json:"address"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
}

type LocalTimeZone struct {
	Zone   string
	Offset int
}

type ApiRedmine struct {
	Host         string `json:"host"`
	PathRedirect struct {
		UpdatePaid string `json:"update_paid"`
	} `json:"path_redirect"`
}
