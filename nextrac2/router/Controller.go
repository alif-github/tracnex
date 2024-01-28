package router

import (
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/endpoint/BacklogEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/CompanyEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/DashboardEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/DepartmentEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/EmployeeAllowanceEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/EmployeeContractEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/EmployeeEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/EmployeeFacilitiesActiveEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/EmployeeGradeEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/EmployeeHistoryEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/EmployeeLevelEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/EmployeeMasterBenefitEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/EmployeePositionEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/EnumEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/MigrationEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ParameterEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ReportAnnualLeaveEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ReportEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/StandarManhourEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/WhiteListDeviceEndpoint"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/endpoint"
	"nexsoft.co.id/nextrac2/endpoint/ActivationLicenseEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ActivationUserNexmileEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/AddResourceExternalEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/AuditMonitoringEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ClientCredentialEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ClientMappingEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ClientRegistrationLogEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ClientTypeEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/CommonEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/CompanyProfileEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/CompanyTitleEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ComponentEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/CustomerCategoryEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/CustomerEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/CustomerGroupEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/CustomerInstallationEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/CustomerSIteEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/DataGroupEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/DistrictEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/DownloadLogEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/GenerateInternalTokenEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ImportFileCustomerListEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/JobProccessEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/LicenseConfigEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/LicenseJournalEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/LicenseTypeEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/LicenseVariantEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/MasterCustomerEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/MenuEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ModuleEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/NexmileParameter"
	"nexsoft.co.id/nextrac2/endpoint/NexsoftRoleEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/PKCEClientMappingEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/PersonProfileEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/PersonTitleEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/PkceUserEndPoint"
	"nexsoft.co.id/nextrac2/endpoint/PositionEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/PostalCodeEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ProductEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ProductGroupEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ProductLicenseEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ProvinceEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/RegisterClientEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/RegisterUserEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/RoleEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/SalesmanEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/SessionEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/SubDistrictEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/UrbanVillageEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/UserActivationEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/UserEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/UserLicenseEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/UserRegistrationAdminEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/UserRegistrationDetailEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/UserVerificationEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ValidationLicenseEndpoint"
	"nexsoft.co.id/nextrac2/endpoint/ValidationNamedUserEndpoint"
)

func APIController() {
	handler := mux.NewRouter()
	prefixPath := config.ApplicationConfiguration.GetServerPrefixPath()
	if prefixPath != "" {
		prefixPath = "/" + prefixPath
	}

	//--------- Internal Common
	handler.HandleFunc(config.ApplicationConfiguration.GetCommonPath().ResourceClients, CommonEndpoint.InternalCommonEndpoint.InternalRegisterUserEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc(config.ApplicationConfiguration.GetCommonPath().NotifyDeletedToken, CommonEndpoint.InternalCommonEndpoint.DeleteTokenEndpoint).Methods("DELETE", "OPTIONS")

	//--------- Session User
	handler.HandleFunc("/v1"+prefixPath+"/session/login/authorize", SessionEndpoint.LoginEndpoint.AuthorizeEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/session/login/verify", SessionEndpoint.LoginEndpoint.VerifyEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/session/login/token", SessionEndpoint.LoginEndpoint.TokenEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/session/login/nexmile", SessionEndpoint.LoginEndpoint.LoginNexmileEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/session/login/grochat", SessionEndpoint.LoginEndpoint.LoginGroChatEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/session/refresh/token", SessionEndpoint.RefreshTokenEndpoint.RefreshTokenEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/users/verify", SessionEndpoint.GetSessionEndpoint.GetSessionEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/users/verify/datetime", SessionEndpoint.GetSessionEndpoint.GetCurrentDatetimeEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/users/verify/dashboard", SessionEndpoint.GetSessionEndpoint.GetDashboardViewEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/session/logout", SessionEndpoint.LogoutEndpoint.LogoutEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/application-version", SessionEndpoint.GetSessionEndpoint.GetApplicationVersionEndpoint).Methods("GET", "OPTIONS")

	//--------- Session Admin
	handler.HandleFunc("/v1"+prefixPath+"/admin/session/login/token", SessionEndpoint.LoginEndpoint.TokenAdminEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/users/verify", SessionEndpoint.GetSessionEndpoint.GetAdminSessionEndpoint).Methods("GET", "OPTIONS")

	//--------- User Sys User
	handler.HandleFunc("/v1"+prefixPath+"/user/{ID}", UserEndpoint.UserSysUserEndpoint.UserEndpointWithParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user/view_profile", UserEndpoint.UserSysUserEndpoint.ProfileSettingSysUser).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user/update_profile", UserEndpoint.UserSysUserEndpoint.ProfileSettingSysUser).Methods("PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user/change_password/{ID}", UserEndpoint.UserSysUserEndpoint.ChangePasswordSysUser).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user/otp/resend/{ID}", UserEndpoint.UserSysUserEndpoint.ResendVerificationCodeSysUser).Methods(http.MethodPost, http.MethodOptions)

	//--------- User Sys Admin
	handler.HandleFunc("/v1"+prefixPath+"/admin/user", UserEndpoint.UserSysAdminEndpoint.UserEndpointWithoutParam).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/user/initiate", UserEndpoint.UserSysAdminEndpoint.InitiateUserParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/user/check", UserEndpoint.UserSysAdminEndpoint.CheckUserAuthEndpointWithoutParam).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/user/login", UserEndpoint.UserSysAdminEndpoint.SysUserLoginEndpointWithoutParam).Methods(http.MethodGet, http.MethodOptions)
	handler.HandleFunc("/v1"+prefixPath+"/admin/user/redis/logout", UserEndpoint.UserSysAdminEndpoint.SysUserLoginEndpointWithoutParam).Methods(http.MethodPost, http.MethodOptions)
	handler.HandleFunc("/v1"+prefixPath+"/admin/user/check/username", UserEndpoint.UserSysAdminEndpoint.CheckUsernameAuthEndpointWithoutParam).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/user/view_profile", UserEndpoint.UserSysAdminEndpoint.ProfileAdminUser).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/user/update_profile", UserEndpoint.UserSysAdminEndpoint.ProfileAdminUser).Methods("PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/user/invitation", UserEndpoint.UserSysAdminEndpoint.Invitation).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/user/change_password/{ID}", UserEndpoint.UserSysAdminEndpoint.ChangePasswordAdminUser).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/user/otp/resend/{ID}", UserEndpoint.UserSysAdminEndpoint.UserVerificationEndpointWithParam).Methods(http.MethodPost, http.MethodOptions)
	handler.HandleFunc("/v1"+prefixPath+"/admin/user/{ID}", UserEndpoint.UserSysAdminEndpoint.UserEndpointWithParam).Methods("DELETE", "PUT", "GET", "OPTIONS")

	//--------- Customer Sys Admin
	handler.HandleFunc("/v1"+prefixPath+"/admin/customer", CustomerEndpoint.CustomerEndpoint.CustomerEndpointWithoutParam).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/customer/initiate", CustomerEndpoint.CustomerEndpoint.InitiateCustomer).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/customer/{ID}", CustomerEndpoint.CustomerEndpoint.CustomerEndpointWithParam).Methods("GET", "OPTIONS")

	//--------- Registration Client Credential
	handler.HandleFunc("/v1"+prefixPath+"/auth/registration/client", RegisterClientEndpoint.RegisterClientEndpoint.RegisterClientEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/auth/initiate/client", RegisterClientEndpoint.RegisterClientEndpoint.InitiateClientEndpoint).Methods("GET", "OPTIONS")

	//--------- Client Additional
	handler.HandleFunc("/v1"+prefixPath+"/auth/registration/branch", ClientMappingEndpoint.ClientMappingEndpoint.InsertNewBranchToClientMappingEndPoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/auth/registration/nexcloud", AddResourceExternalEndpoint.AddResourceExternalEndpoint.AddResourceNexcloudEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/auth/registration/nexdrive", AddResourceExternalEndpoint.AddResourceExternalEndpoint.AddResourceNexdriveEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/auth/registration/nexcloud/{CLIENT}", AddResourceExternalEndpoint.AddResourceExternalEndpoint.ViewLogForAddResourceExternalEndpoint).Methods("GET", "OPTIONS")

	//--------- Registration User PKCE
	handler.HandleFunc("/v1"+prefixPath+"/auth/registration/pkce", RegisterUserEndpoint.RegisterUserPKCEEndpoint.RegistrationUserPKCEEndpoint).Methods("POST", "OPTIONS")
	handler.Path("/v1/user/verify").Queries("activation_code", "{CODE}", "user_id", "{ID}", "email", "{EMAIL}", "username", "{USERNAME}").HandlerFunc(UserActivationEndpoint.UserActivationEndpoint.UserActivationEndpoint).Methods("POST", "OPTIONS")

	//--------- Un-registration User PKCE
	handler.HandleFunc("/v1"+prefixPath+"/auth/registration/pkce/unregister", RegisterUserEndpoint.RegisterUserPKCEEndpoint.UnregisterUserPKCEEndpoint).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/auth/registration/pkce/unregister/{USERNAME}", RegisterUserEndpoint.RegisterUserPKCEEndpoint.ViewForUnregisterPKCEEndpoint).Methods("GET", "OPTIONS")

	//--------- Nexmile req change password
	handler.HandleFunc("/v1"+prefixPath+"/auth/registration/pkce/passwords", PkceUserEndPoint.ChangePasswordEndPoint.ChangePasswordEndPoint).Methods("POST", "OPTIONS")

	//--------- Detail client mapping
	handler.HandleFunc("/v1/internal/users/client_mapping", ClientMappingEndpoint.ClientMappingEndpoint.DetailMappingEndPoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1/internal/users/client_mapping_by_clientid", ClientMappingEndpoint.ClientMappingEndpoint.DetailMappingEndPointByClientID).Methods("GET", "OPTIONS")

	//---------- Health
	handler.HandleFunc("/v1"+prefixPath+"/health", endpoint.HealthEndpoint.GetHealthStatus).Methods("GET", "OPTIONS")
	handler.Handle("/metrics", promhttp.Handler()).Methods("GET", "OPTIONS")
	handler.Handle("/health/prometheus", promhttp.Handler()).Methods("GET", "OPTIONS")

	//---------- Client Mapping
	handler.HandleFunc("/v1"+prefixPath+"/auth/registration/usersocketid", ClientMappingEndpoint.ClientMappingEndpoint.ClientMappingChangeSocketIDEndpoint).Methods("PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/client_mapping", ClientMappingEndpoint.ClientMappingEndpoint.UIClientMappingWithoutPathParamEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/client_mapping/initiate", ClientMappingEndpoint.ClientMappingEndpoint.UIGetListCLientMappingInitiateEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/client_mapping/change_name/{ID}", ClientMappingEndpoint.ClientMappingEndpoint.UIClientMappingChangeNameND6Endpoint).Methods("PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/client_mapping/{ID}", ClientMappingEndpoint.ClientMappingEndpoint.UIClientMappingWithPathParamEndpoint).Methods("GET", "OPTIONS")

	//---------- PKCE Client Mapping
	handler.HandleFunc("/v1"+prefixPath+"/pkce_client_mapping", PKCEClientMappingEndpoint.PKCEClientMappingEndpoints.UIPKCEClientMappingWithoutPathParamEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/pkce_client_mapping/initiate", PKCEClientMappingEndpoint.PKCEClientMappingEndpoints.UIGetListPKCEClientMappingInitiateEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/pkce_client_mapping/change_name/{ID}", PKCEClientMappingEndpoint.PKCEClientMappingEndpoints.ChangeNamePCKEClientMappingEndpoint).Methods("PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/pkce_client_mapping/{ID}", PKCEClientMappingEndpoint.PKCEClientMappingEndpoints.UIPKCEClientMappingWithPathParamEndpoint).Methods("GET", "OPTIONS")

	//---------- Role
	handler.HandleFunc("/v1"+prefixPath+"/admin/role", RoleEndpoint.RoleEndpoint.RoleWithoutParam).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/role/initiate", RoleEndpoint.RoleEndpoint.InitiateGetListRole).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/role/{ID}", RoleEndpoint.RoleEndpoint.RoleWithParam).Methods("DELETE", "GET", "PUT", "OPTIONS")

	//---------- Nexsoft Role
	handler.HandleFunc("/v1"+prefixPath+"/admin/nexsoft_role", NexsoftRoleEndpoint.NexsoftRoleEndpoint.NexsoftRoleWithoutPathParamEndpoint).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/nexsoft_role/initiate", NexsoftRoleEndpoint.NexsoftRoleEndpoint.InitiateNexsoftRoleEnpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/nexsoft_role/{ID}", NexsoftRoleEndpoint.NexsoftRoleEndpoint.NexsoftRoleWithPathParamEndpoint).Methods("PUT", "GET", "DELETE", "OPTIONS")

	//---------- Help
	handler.HandleFunc("/v1"+prefixPath+"/help/internal_token", GenerateInternalTokenEndpoint.GenerateInTokenEndpoint.GenerateInternalToken).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/help/client_credential", ClientCredentialEndpoint.ClientCredentialEndpoint.ClientCredentialHelping).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/help/user/{ID}", UserEndpoint.UserSysAdminEndpoint.UserEndpointHelpingDeleted).Methods("PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/help/client/{ID}", UserEndpoint.UserSysAdminEndpoint.UserEndpointHelpingFirstName).Methods("PUT", "OPTIONS")

	//---------- Import File
	handler.HandleFunc("/v1"+prefixPath+"/import", ImportFileCustomerListEndpoint.ImportFileCustomerListEndpoint.ImportFileCustomerListEndpointWithoutParam).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/import/validate", ImportFileCustomerListEndpoint.ImportFileCustomerListEndpoint.ImportAndValidateDataEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/import/initiate", ImportFileCustomerListEndpoint.ImportFileCustomerListEndpoint.InitiateImportEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/import/confirm", ImportFileCustomerListEndpoint.ImportFileCustomerListEndpoint.ConfirmImportEndpoint).Methods("POST", "OPTIONS")

	//---------- Job Process
	handler.HandleFunc("/v1"+prefixPath+"/job_process", JobProccessEndpoint.JobProcessEndpoint.JobProcessEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/job_process/initiate", JobProccessEndpoint.JobProcessEndpoint.InitiateGetListJobProcess).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/job_process/{ID}", JobProccessEndpoint.JobProcessEndpoint.JobProcessEndpointWithParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/synchronize/regional", JobProccessEndpoint.JobProcessEndpoint.SynchronizeRegionalDataEndpoint).Methods("POST", "OPTIONS")

	//---------- Log Registration
	handler.HandleFunc("/v1"+prefixPath+"/registration_log", ClientRegistrationLogEndpoint.ClientRegistrationLogEndpoint.ClientRegistrationLogWithoutParam).Methods("GET", "OPTIONS")

	//---------- Menu
	handler.HandleFunc("/v1"+prefixPath+"/menu/parent", MenuEndpoint.MenuEndpoint.ServiceMenuParentSysUserWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/menu/parent_sysadmin", MenuEndpoint.MenuEndpoint.ServiceMenuParentSysAdminWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/menu/parent/{ID}", MenuEndpoint.MenuEndpoint.ServiceMenuParentSysUserWithParam).Methods("PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/menu/parent_sysadmin/{ID}", MenuEndpoint.MenuEndpoint.ServiceMenuParentSysAdminWithParam).Methods("PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/menu/service/{ID}", MenuEndpoint.MenuEndpoint.ServiceMenuServiceWithParam).Methods("GET", "PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/menu/item/{ID}", MenuEndpoint.MenuEndpoint.ServiceMenuItemWithParam).Methods("GET", "PUT", "OPTIONS")

	//---------- Person Title
	handler.HandleFunc("/v1"+prefixPath+"/persontitle", PersonTitleEndpoint.PersonTitleEndpoint.PersonTitleEndpointWithoutParam).Methods("GET", "OPTIONS")

	//---------- Data Group
	handler.HandleFunc("/v1"+prefixPath+"/admin/data_group", DataGroupEndpoint.DataGroupEndpoint.DataGroupEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/data_group/initiate", DataGroupEndpoint.DataGroupEndpoint.InitiateDataGroupEndpoint).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/data_group/{ID}", DataGroupEndpoint.DataGroupEndpoint.DataGroupEndpointWithPathParam).Methods("GET", "PUT", "DELETE", "OPTIONS")

	//---------- Data Scope
	handler.HandleFunc("/v1"+prefixPath+"/admin/data_scope", DataGroupEndpoint.DataGroupEndpoint.InsertHelperDataScopeEndpoint).Methods("POST", "OPTIONS")

	//---------- Customer Group
	handler.HandleFunc("/v1"+prefixPath+"/admin/customer_group", CustomerGroupEndpoint.CustomerGroupEndpoint.CustomerGroupAdminEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customer_group", CustomerGroupEndpoint.CustomerGroupEndpoint.CustomerGroupEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customer_group/initiate", CustomerGroupEndpoint.CustomerGroupEndpoint.InitiateCustomerGroupEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customer_group/{ID}", CustomerGroupEndpoint.CustomerGroupEndpoint.CustomerGroupEndpointWithPathParam).Methods("GET", "PUT", "DELETE", "OPTIONS")

	//---------- Salesman
	handler.HandleFunc("/v1"+prefixPath+"/admin/salesman", SalesmanEndpoint.SalesmanEndpoint.SalesmanWithoutParamAdmin).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/salesman", SalesmanEndpoint.SalesmanEndpoint.SalesmanWithoutParam).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/salesman/initiate", SalesmanEndpoint.SalesmanEndpoint.InitiateSalesman).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/salesman/{ID}", SalesmanEndpoint.SalesmanEndpoint.SalesmanWithParam).Methods("PUT", "GET", "DELETE", "OPTIONS")

	//---------- Customer Category
	handler.HandleFunc("/v1"+prefixPath+"/admin/customer_category", CustomerCategoryEndpoint.CustomerCategoryEndpoint.CustomerCategoryAdminEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customer_category", CustomerCategoryEndpoint.CustomerCategoryEndpoint.CustomerCategoryEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customer_category/initiate", CustomerCategoryEndpoint.CustomerCategoryEndpoint.InitiateCustomerCategoryEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customer_category/{ID}", CustomerCategoryEndpoint.CustomerCategoryEndpoint.CustomerCategoryEndpointWithPathParam).Methods("GET", "PUT", "DELETE", "OPTIONS")

	//---------- Province
	handler.HandleFunc("/v1"+prefixPath+"/admin/province", ProvinceEndpoint.ProvinceEndpoint.ProvinceAdminEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/province/last-sync", ProvinceEndpoint.ProvinceEndpoint.ResetLastSyncProvinceEndpoint).Methods("PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/province", ProvinceEndpoint.ProvinceEndpoint.ProvinceEndpointWithoutParam).Methods("GET", "OPTIONS")

	//---------- District
	handler.HandleFunc("/v1"+prefixPath+"/admin/district", DistrictEndpoint.DistrictEndpoint.DistrictAdminEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/district", DistrictEndpoint.DistrictEndpoint.DistrictEndpointWithoutParam).Methods("GET", "OPTIONS")

	//---------- Product Group
	handler.HandleFunc("/v1"+prefixPath+"/admin/product_group", ProductGroupEndpoint.ProductGroupEndpoint.ProductGroupAdminEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/product_group", ProductGroupEndpoint.ProductGroupEndpoint.ProductGroupEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/product_group/initiate", ProductGroupEndpoint.ProductGroupEndpoint.InitiateProductGroupEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/product_group/{ID}", ProductGroupEndpoint.ProductGroupEndpoint.ProductGroupEndpointWithPathParam).Methods("GET", "PUT", "DELETE", "OPTIONS")

	//---------- License Variant
	handler.HandleFunc("/v1"+prefixPath+"/license_variant", LicenseVariantEndpoint.LicenseVariantEndpoint.LicenseVariantWithoutParam).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/license_variant/initiate", LicenseVariantEndpoint.LicenseVariantEndpoint.InitiateLicenseVariant).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/license_variant/{ID}", LicenseVariantEndpoint.LicenseVariantEndpoint.LicenseVariantWithParam).Methods("PUT", "GET", "DELETE", "OPTIONS")

	//--------- Session Client
	handler.HandleFunc("/v1"+prefixPath+"/token", SessionEndpoint.LoginEndpoint.TokenClientEndpoint).Methods("POST", "OPTIONS")

	//---------- License Type
	handler.HandleFunc("/v1"+prefixPath+"/license_type", LicenseTypeEndpoint.LicenseTypeEndpoint.LicenseTypeEndpointWithoutParam).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/license_type/initiate", LicenseTypeEndpoint.LicenseTypeEndpoint.InitiateLicenseTypeEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/license_type/{ID}", LicenseTypeEndpoint.LicenseTypeEndpoint.LicenseTypeEndpointWithPathParam).Methods("GET", "PUT", "DELETE", "OPTIONS")

	//---------- Module
	handler.HandleFunc("/v1"+prefixPath+"/module", ModuleEndpoint.ModuleEndpoint.ModuleEndpointWithoutParam).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/module/initiate", ModuleEndpoint.ModuleEndpoint.InitiateModuleEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/module/{ID}", ModuleEndpoint.ModuleEndpoint.ModuleEndpointWithPathParam).Methods("GET", "PUT", "DELETE", "OPTIONS")

	//---------- Component
	handler.HandleFunc("/v1"+prefixPath+"/component", ComponentEndpoint.ComponentEndpoint.ComponentEndpointWithoutParam).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/component/initiate", ComponentEndpoint.ComponentEndpoint.InitiateComponentEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/component/{ID}", ComponentEndpoint.ComponentEndpoint.ComponentEndpointWithPathParam).Methods("GET", "PUT", "DELETE", "OPTIONS")

	//---------- Product
	handler.HandleFunc("/v1"+prefixPath+"/product", ProductEndpoint.ProductEndpoint.ProductWithoutParam).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/product/initiate", ProductEndpoint.ProductEndpoint.InitiateProduct).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/product/{ID}", ProductEndpoint.ProductEndpoint.ProductWithParam).Methods("GET", "DELETE", "PUT", "OPTIONS")

	//----------------------------------------------------------------- Master-data section ----------------------------------------------------------------------------------
	//---------- Position
	handler.HandleFunc("/v1"+prefixPath+"/position", PositionEndpoint.PositionEndpoint.PositionEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/position/{ID}", PositionEndpoint.PositionEndpoint.PositionEndpointWithPathParam).Methods("GET", "OPTIONS")

	//---------- Company Title
	handler.HandleFunc("/v1"+prefixPath+"/companytitle", CompanyTitleEndpoint.CompanyTitleEndpoint.CompanyTitleEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/companytitle/{ID}", CompanyTitleEndpoint.CompanyTitleEndpoint.CompanyTitleEndpointWithPathParam).Methods("GET", "OPTIONS")

	//---------- Company Profile
	handler.HandleFunc("/v1"+prefixPath+"/companyprofile", CompanyProfileEndpoint.CompanyProfileEndpoint.CompanyProfileEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/companyprofile/{ID}", CompanyProfileEndpoint.CompanyProfileEndpoint.CompanyProfileEndpointWithPathParam).Methods("GET", "OPTIONS")

	//---------- Person Profile
	handler.HandleFunc("/v1"+prefixPath+"/personprofile", PersonProfileEndpoint.PersonProfileEndpoint.PersonProfileEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/personprofile/{ID}", PersonProfileEndpoint.PersonProfileEndpoint.PersonProfileEndpointWithPathParam).Methods("GET", "OPTIONS")

	//---------- Client Type
	handler.HandleFunc("/v1"+prefixPath+"/admin/clienttype", ClientTypeEndpoint.ClientTypeEndpoint.ClientTypeAdminEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/clienttype", ClientTypeEndpoint.ClientTypeEndpoint.ClientTypeEndpointWithoutParam).Methods("GET", "POST", "PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/clienttype/initiate", ClientTypeEndpoint.ClientTypeEndpoint.InitiateClientTypeEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/clienttype/{ID}", ClientTypeEndpoint.ClientTypeEndpoint.ClientTypeEndpointWithParam).Methods("PUT", "GET", "DELETE", "OPTIONS")

	//---------- Urban Village
	handler.HandleFunc("/v1"+prefixPath+"/urbanvillage", UrbanVillageEndpoint.UrbanVillageEndpoint.UrbanVillageEndpointWhitoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/urbanvillage/initiate", UrbanVillageEndpoint.UrbanVillageEndpoint.InitiateUrbanVillageEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/urbanvillage/{ID}", UrbanVillageEndpoint.UrbanVillageEndpoint.UrbanVillageEndpointWhitPathParam).Methods("GET", "OPTIONS")

	//---------- Sub District
	handler.HandleFunc("/v1"+prefixPath+"/subdistrict", SubDistrictEndpoint.SubDistrictEndpoint.SubDistrictEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/subdistrict/initiate", SubDistrictEndpoint.SubDistrictEndpoint.InitiateSubDistrictEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/subdistrict/{ID}", SubDistrictEndpoint.SubDistrictEndpoint.SubDistrictEndpointWhithPathParam).Methods("GET", "OPTIONS")

	//---------- Customer Site
	handler.HandleFunc("/v1"+prefixPath+"/customersite", CustomerSIteEndpoint.CustomerSiteEndpoint.CustomerSiteEndpointWithoutParam).Methods("POST", "OPTIONS")

	//---------- Customer Installation
	handler.HandleFunc("/v1"+prefixPath+"/customerinstallation/{ID}", CustomerInstallationEndpoint.CustomerInstallationEndpoint.CustomerInstallationEndpointWithParam).Methods("PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customerinstallation/site", CustomerInstallationEndpoint.CustomerInstallationEndpoint.CustomerSiteInSiteInstallationEndpointWithParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customerinstallation/installation", CustomerInstallationEndpoint.CustomerInstallationEndpoint.CustomerInstallationInSiteInstallationEndpointWithParam).Methods("GET", "OPTIONS")

	//---------- Postal Code
	handler.HandleFunc("/v1"+prefixPath+"/postalcode", PostalCodeEndpoint.PostalCodeEndpoint.PostalCodeEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/postalcode/initiate", PostalCodeEndpoint.PostalCodeEndpoint.InitiatePostalCodeEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/postalcode/{ID}", PostalCodeEndpoint.PostalCodeEndpoint.PostalCodeEndpointWithPathParam).Methods("GET", "OPTIONS")

	//---------- Customer
	handler.HandleFunc("/v1"+prefixPath+"/customer", MasterCustomerEndpoint.MasterCustomerEndpoint.CustomerEndpointWithoutPathParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customer/initiate", MasterCustomerEndpoint.MasterCustomerEndpoint.InitiateCustomerEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customer/{ID}", MasterCustomerEndpoint.MasterCustomerEndpoint.CustomerEndpointWithPathParam).Methods("GET", "PUT", "DELETE", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customer-nonparent", MasterCustomerEndpoint.MasterCustomerEndpoint.GetListCustomerNonParentEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customer-nonparent/initiate", MasterCustomerEndpoint.MasterCustomerEndpoint.InitiateCustomerNonParentEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customer-parent", MasterCustomerEndpoint.MasterCustomerEndpoint.GetListCustomerParentEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/customer-parent/initiate", MasterCustomerEndpoint.MasterCustomerEndpoint.InitiateCustomerParentEndpoint).Methods("GET", "OPTIONS")

	//---------- License Config
	handler.HandleFunc("/v1"+prefixPath+"/licenseconfig", LicenseConfigEndpoint.LicenseConfigEndpoint.LicenseConfigWithoutParam).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/licenseconfig/extend", LicenseConfigEndpoint.LicenseConfigEndpoint.InsertMultipleLicenseConfig).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/licenseconfig/initiate", LicenseConfigEndpoint.LicenseConfigEndpoint.InitiateLicenseConfig).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/licenseconfig/all", LicenseConfigEndpoint.LicenseConfigEndpoint.SelectAllLicenseConfig).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/licenseconfig/{ID}", LicenseConfigEndpoint.LicenseConfigEndpoint.LicenseConfigWithParam).Methods("DELETE", "PUT", "GET", "OPTIONS")

	//---------- Installation Detail
	handler.HandleFunc("/v1"+prefixPath+"/installation-detail/{INSTALLATIONID}", CustomerInstallationEndpoint.CustomerInstallationEndpoint.DetailInstallationEndpointWithParam).Methods("GET", "OPTIONS")

	//---------- Product License
	handler.HandleFunc("/v1"+prefixPath+"/product-license", ProductLicenseEndpoint.ProductLicenseEndpoint.ProductLicenseWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/product-license/initiate", ProductLicenseEndpoint.ProductLicenseEndpoint.InitiateProductLicenseEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/product-license/decrypt/{id}", ProductLicenseEndpoint.ProductLicenseEndpoint.DecryptProductLicenseEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/product-license/{id}", ProductLicenseEndpoint.ProductLicenseEndpoint.DetailProductLicenseWithParam).Methods("GET", "PUT", "OPTIONS")

	//---------- User License
	handler.HandleFunc("/v1"+prefixPath+"/user-license", UserLicenseEndpoint.UserLicenseEndpoint.UserLicenseWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user-license/initiate", UserLicenseEndpoint.UserLicenseEndpoint.InitiateUserLicenseEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user-license/view", UserLicenseEndpoint.UserLicenseEndpoint.ViewDetailUserLicenseEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user-license/view/initiate", UserLicenseEndpoint.UserLicenseEndpoint.InitiateViewDetailUserLicenseEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user-license/transfer-key/{ID}", UserLicenseEndpoint.UserLicenseEndpoint.TransferUserLicenseEndpointWithPathParam).Methods("PUT", "GET", "OPTIONS")

	//---------- New Registration User Nexmile and Nexstar
	//handler.HandleFunc("/v1"+prefixPath+"/user-registration/register", UserRegistrationDetailEndpoint.UserRegistrationDetailEndpoint.RegisterAndActivateUserNexmileNexstar).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user-registration/resend", UserRegistrationDetailEndpoint.UserRegistrationDetailEndpoint.ResendActivationNexmileNexstar).Methods("PUT", "OPTIONS")

	//---------- User Registration Detail (License Named User)
	handler.HandleFunc("/v1"+prefixPath+"/user-registration/check", UserRegistrationDetailEndpoint.UserRegistrationDetailEndpoint.CheckLicenseNamedUser).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user-registration/add", UserRegistrationDetailEndpoint.UserRegistrationDetailEndpoint.UserRegistrationDetailNamedUserWithoutParam).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user-registration/register_user", UserRegistrationDetailEndpoint.UserRegistrationDetailEndpoint.UserRegistrationDetailNamedUserClientMappingWithoutParam).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user-registration/register-user/check", UserRegistrationDetailEndpoint.UserRegistrationDetailEndpoint.UserRegistrationDetailNamedUserCheckWithoutParam).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user-registration/unregister/{id}", UserRegistrationDetailEndpoint.UserRegistrationDetailEndpoint.UserRegistrationDetailNamedUserWithParam).Methods("PUT", "OPTIONS")

	//---------- Register Client Non On Premise
	handler.HandleFunc("/v1"+prefixPath+"/client-mapping", RegisterClientEndpoint.RegisterClientEndpoint.RegisterClientNonOnPremiseEndpoint).Methods("POST", "OPTIONS")

	//---------- Activation License
	handler.HandleFunc("/v1"+prefixPath+"/activation-license", ActivationLicenseEndpoint.ActivationLicenseEndpoint.ActivationLicenseEndpointWithoutPathParam).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/activation-license-nexmile", ActivationUserNexmileEndpoint.ActivationUserNexmileEndpoint.ActivateUserNexmileEndpoint).Methods("PUT", "OPTIONS")

	//---------- Validation License
	handler.HandleFunc("/v1"+prefixPath+"/validation-license", ValidationLicenseEndpoint.ValidationLicenseEndpoint.ValidateLicenseEndpoint).Methods("PUT", "OPTIONS")

	//---------- Update License HWID
	handler.HandleFunc("/v1"+prefixPath+"/update-license-hwid", ProductLicenseEndpoint.ProductLicenseEndpoint.UpdateProductLicenseHWIDEndpoint).Methods("PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/update-license-hwid/internal", ProductLicenseEndpoint.ProductLicenseEndpoint.UpdateProductLicenseHWIDByPassEndpoint).Methods("PUT", "OPTIONS")

	//---------- User Registration Admin
	handler.HandleFunc("/v1"+prefixPath+"/user-registration-admin", UserRegistrationAdminEndpoint.UserRegistrationAdminEndpoint.UserRegistrationAdminWithoutParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user-registration-admin/initiate", UserRegistrationAdminEndpoint.UserRegistrationAdminEndpoint.UserRegistrationAdminInitiate).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user-registration-admin/{id}", UserRegistrationAdminEndpoint.UserRegistrationAdminEndpoint.UserRegistrationAdminWithParam).Methods("GET", "OPTIONS")

	//---------- Validation Named User Online
	handler.HandleFunc("/v1"+prefixPath+"/user-registration/validate", ValidationNamedUserEndpoint.ValidationNamedUserEndpoint.ValidateNamedUserEndpoint).Methods("POST", "OPTIONS")

	//---------- Nexmile Parameter
	handler.HandleFunc("/v1"+prefixPath+"/nexmile-parameter", NexmileParameter.NexmileParametersEndpoint.GetNexmileParameterWithoutParam).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/nexmile-parameter/add", NexmileParameter.NexmileParametersEndpoint.AddNexmileParameterWithoutParam).Methods("POST", "OPTIONS")

	//---------- User Nexmile Verification
	handler.HandleFunc("/v1"+prefixPath+"/user-verification", UserVerificationEndpoint.UserVerificationEndpoint.UserVerificationWithoutParam).Methods("PUT", "OPTIONS")

	//--------- Audit Monitoring
	handler.HandleFunc("/v1"+prefixPath+"/audit-monitoring", AuditMonitoringEndpoint.AuditMonitoringEndpoint.AuditMonitoringEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/admin/audit-monitoring", AuditMonitoringEndpoint.AuditMonitoringEndpoint.AuditMonitoringEndpointWithoutParam).Methods("GET", "OPTIONS")

	//--------- Download File Log Sysadmin
	handler.HandleFunc("/v1"+prefixPath+"/admin/logfile/download", DownloadLogEndpoint.DownloadLogEndpoint.DownloadLogEndpoint).Methods("GET", "OPTIONS")

	//--------- Session Client
	handler.HandleFunc("/v1"+prefixPath+"/token", SessionEndpoint.LoginEndpoint.TokenClientEndpoint).Methods("POST", "OPTIONS")

	//--------- API Whitelist Forget User
	handler.HandleFunc("/v1"+prefixPath+"/user/forget", UserEndpoint.UserEndpoint.UserEndpointWithoutParamResetPassword).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/user/forget/change-password", UserEndpoint.UserEndpoint.UserEndpointWithoutParamChangePassword).Methods("POST", "OPTIONS")

	//-**********************************************************************************************************************************************************************************************************************************
	//--------- [Internal API]
	//-**********************************************************************************************************************************************************************************************************************************
	internal := handler.PathPrefix("/v1" + prefixPath + "/internal").Subrouter()

	//--------- [Internal] License Journal
	internal.HandleFunc("/license-journal", LicenseJournalEndpoint.LicenseJournalEndpoint.ListLicenseJournalEksternalEndpoint).Methods("POST", "OPTIONS")
	internal.HandleFunc("/license-journal/count", LicenseJournalEndpoint.LicenseJournalEndpoint.InitiateLicenseJournalEksternalEndpoint).Methods("POST", "OPTIONS")

	//--------- [Internal] Customer
	internal.HandleFunc("/customer", MasterCustomerEndpoint.MasterCustomerEndpoint.InternalGetListCustomerEndpoint).Methods("POST", "OPTIONS")
	internal.HandleFunc("/distributor", MasterCustomerEndpoint.MasterCustomerEndpoint.InternalGetListDistributorEndpoint).Methods("POST", "OPTIONS")
	internal.HandleFunc("/customer/count", MasterCustomerEndpoint.MasterCustomerEndpoint.InternalCountCustomerEndpoint).Methods("POST", "OPTIONS")
	//-**********************************************************************************************************************************************************************************************************************************

	//---------- White List CRUD
	handler.HandleFunc("/v1"+prefixPath+"/whitelist-device", WhiteListDeviceEndpoint.WhiteListDeviceEndpoint.WhiteListDeviceEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/whitelist-device/initiate", WhiteListDeviceEndpoint.WhiteListDeviceEndpoint.InitiateGetListWhiteListDeviceEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/whitelist-device/{ID}", WhiteListDeviceEndpoint.WhiteListDeviceEndpoint.WhiteListDeviceEndpointWithParam).Methods("PUT", "GET", "DELETE", "OPTIONS")

	//--------- CS Report Bug
	handler.HandleFunc("/v1"+prefixPath+"/internal-update-hwid", ProductLicenseEndpoint.ProductLicenseEndpoint.UpdateProductLicenseHWIDByPassEndpoint).Methods("PUT", "OPTIONS")

	//--------- Department
	department := handler.PathPrefix("/v1" + prefixPath + "/department").Subrouter()
	department.HandleFunc("", DepartmentEndpoint.DepartmentEndpoint.DepartmentEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	department.HandleFunc("/initiate", DepartmentEndpoint.DepartmentEndpoint.InitiateDepartmentEndpointWithoutParam).Methods("GET", "OPTIONS")
	department.HandleFunc("/{ID}", DepartmentEndpoint.DepartmentEndpoint.DepartmentEndpointWithParam).Methods("PUT", "GET", "OPTIONS")

	//--------- Employee & Employee Timesheet
	e := handler.PathPrefix("/v1" + prefixPath + "/employee").Subrouter()
	et := e.PathPrefix("/timesheet").Subrouter()

	//---------- Employee Admin
	handler.HandleFunc("/v1"+prefixPath+"/admin/employee", EmployeeEndpoint.EmployeeEndpoint.EmployeeAdminEndpointWithoutParam).Methods("GET", "OPTIONS")

	//--------- Employee Leave
	employeeLeave := e.PathPrefix("/leave").Subrouter()
	employeeLeave.HandleFunc("", EmployeeEndpoint.EmployeeEndpoint.EmployeeLeaveEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	employeeLeave.HandleFunc("/initiate", EmployeeEndpoint.EmployeeEndpoint.InitiateGetListEmployeeLeave).Methods("GET", "POST", "OPTIONS")
	employeeLeave.HandleFunc("/remaining-leave", EmployeeEndpoint.EmployeeEndpoint.GetRemainingLeave).Methods("GET", "OPTIONS")
	employeeLeave.HandleFunc("/download", EmployeeEndpoint.EmployeeEndpoint.DownloadEmployeeLeaveReport).Methods("GET", "OPTIONS")
	employeeLeave.HandleFunc("/types", EmployeeEndpoint.EmployeeEndpoint.GetEmployeeLeaveTypes).Methods("GET", "OPTIONS")
	employeeLeave.HandleFunc("/types/initiate", EmployeeEndpoint.EmployeeEndpoint.InitiateGetEmployeeLeaveTypes).Methods("GET", "OPTIONS")
	employeeLeave.HandleFunc("/annual", EmployeeEndpoint.EmployeeEndpoint.GetListEmployeeLeaveYearlyEndpoint).Methods("GET", "OPTIONS")
	employeeLeave.HandleFunc("/annual/download", EmployeeEndpoint.EmployeeEndpoint.DownloadAnnualLeaveReport).Methods("GET", "OPTIONS")
	employeeLeave.HandleFunc("/annual/initiate", EmployeeEndpoint.EmployeeEndpoint.IntiateEmployeeLeaveYearlyEndpoint).Methods("GET", "OPTIONS")
	employeeLeave.HandleFunc("/{ID}", EmployeeEndpoint.EmployeeEndpoint.EmployeeLeaveEndpointWithParam).Methods("GET", "OPTIONS")

	//--------- Employee Reimbursement
	employeeReimbursement := e.PathPrefix("/reimbursement").Subrouter()
	employeeReimbursement.HandleFunc("", EmployeeEndpoint.EmployeeEndpoint.EmployeeReimbursementEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	employeeReimbursement.HandleFunc("/initiate", EmployeeEndpoint.EmployeeEndpoint.InitiateGetListEmployeeReimbursement).Methods("GET", "OPTIONS")
	employeeReimbursement.HandleFunc("/remaining-balance", EmployeeEndpoint.EmployeeEndpoint.GetMedicalRemainingBalance).Methods("GET", "OPTIONS")
	employeeReimbursement.HandleFunc("/types", EmployeeEndpoint.EmployeeEndpoint.GetEmployeeReimbursementTypes).Methods("GET", "OPTIONS")
	employeeReimbursement.HandleFunc("/types/initiate", EmployeeEndpoint.EmployeeEndpoint.InitiateGetEmployeeReimbursementTypes).Methods("GET", "OPTIONS")
	employeeReimbursement.HandleFunc("/verify/{ID}", EmployeeEndpoint.EmployeeEndpoint.VerifyReimbursementEnpoint).Methods("PUT", "OPTIONS")
	employeeReimbursement.HandleFunc("/download", EmployeeEndpoint.EmployeeEndpoint.DownloadReimbursementReport).Methods("GET", "OPTIONS")
	employeeReimbursement.HandleFunc("/report", EmployeeEndpoint.EmployeeEndpoint.GetListEmployeeReimbursementReport).Methods("GET", "OPTIONS")
	employeeReimbursement.HandleFunc("/report/initiate", EmployeeEndpoint.EmployeeEndpoint.InitiateGetListEmployeeReimbursementReport).Methods("GET", "OPTIONS")

	//--------- Employee Notification
	employeeNotification := e.PathPrefix("/notification").Subrouter()
	employeeNotification.HandleFunc("", EmployeeEndpoint.EmployeeEndpoint.EmployeeNotificationEndpointWithoutParam).Methods("GET", "OPTIONS")
	employeeNotification.HandleFunc("/read", EmployeeEndpoint.EmployeeEndpoint.ReadNotificationEndpoint).Methods("PUT", "OPTIONS")

	//--------- Employee Profile
	e.HandleFunc("", EmployeeEndpoint.EmployeeEndpoint.EmployeeEndpointWithoutParam).Methods("POST", "GET", "OPTIONS")
	e.HandleFunc("/initiate", EmployeeEndpoint.EmployeeEndpoint.InitiateEmployeeEndpointWithoutParam).Methods("GET", "OPTIONS")
	et.HandleFunc("", EmployeeEndpoint.EmployeeEndpoint.EmployeeTimeSheetEndpointWithoutParam).Methods("GET", "OPTIONS")
	e.HandleFunc("/{ID}", EmployeeEndpoint.EmployeeEndpoint.EmployeeEndpointWithParam).Methods("DELETE", "PUT", "GET", "OPTIONS")

	//--------- Employee Timesheet
	et.HandleFunc("/initiate", EmployeeEndpoint.EmployeeEndpoint.InitiateEmployeeTimeSheetEndpointWithoutParam).Methods("GET", "OPTIONS")
	et.HandleFunc("/check", EmployeeEndpoint.EmployeeEndpoint.EmployeeTimeSheetCheckEndpointWithoutParam).Methods("POST", "OPTIONS")
	et.HandleFunc("/{ID}", EmployeeEndpoint.EmployeeEndpoint.EmployeeTimeSheetEndpointWithParam).Methods("PUT", "GET", "OPTIONS")

	//--------- Employee History
	employeeHistory := e.PathPrefix("/history").Subrouter()
	employeeHistory.HandleFunc("/request", EmployeeEndpoint.EmployeeEndpoint.EmployeeRequestHistoryEndpointWithoutParam).Methods("GET", "OPTIONS")
	employeeHistory.HandleFunc("/request/initiate", EmployeeEndpoint.EmployeeEndpoint.InitiateEmployeeRequestHistoryEndpoint).Methods("GET", "OPTIONS")
	employeeHistory.HandleFunc("/request/cancellation/{ID}", EmployeeEndpoint.EmployeeEndpoint.CancelEmployeeRequestEndpoint).Methods("PUT", "OPTIONS")
	employeeHistory.HandleFunc("/approval", EmployeeEndpoint.EmployeeEndpoint.EmployeeApprovalHistoryEndpointWithoutParam).Methods("GET", "OPTIONS")
	employeeHistory.HandleFunc("/approval/initiate", EmployeeEndpoint.EmployeeEndpoint.InitiateEmployeeApprovalHistoryEndpoint).Methods("GET", "OPTIONS")
	employeeHistory.HandleFunc("/approval/status/{ID}", EmployeeEndpoint.EmployeeEndpoint.UpdateStatusEndpointWithoutParam).Methods("PUT", "OPTIONS")

	//--------- EmployeeGrade
	handler.HandleFunc("/v1"+prefixPath+"/employee-grade", EmployeeGradeEndpoint.EmployeeGradeEndpoint.EmployeeGradeEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-grade/initiate", EmployeeGradeEndpoint.EmployeeGradeEndpoint.InitiateEmployeeGradeEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-grade/{ID}", EmployeeGradeEndpoint.EmployeeGradeEndpoint.EmployeeGradeEndpointWithParam).Methods("DELETE", "PUT", "OPTIONS")

	//--------- EmployeeLevel
	handler.HandleFunc("/v1"+prefixPath+"/employee-level", EmployeeLevelEndpoint.EmployeeLevelEndpoint.EmployeeLevelEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-level/initiate", EmployeeLevelEndpoint.EmployeeLevelEndpoint.InitiateEmployeeLevelEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-level/{ID}", EmployeeLevelEndpoint.EmployeeLevelEndpoint.EmployeeLevelEndpointWithParam).Methods("DELETE", "PUT", "OPTIONS")

	//--------- EmployeeGradeMatrix
	handler.HandleFunc("/v1"+prefixPath+"/employee-grade-matrix", EmployeeGradeEndpoint.EmployeeGradeEndpoint.EmployeeGradeMatrixEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-grade-matrix/initiate", EmployeeGradeEndpoint.EmployeeGradeEndpoint.InitiateEmployeeGradeMatrixEndpoint).Methods("GET", "OPTIONS")

	//--------- EmployeeLevelMatrix
	handler.HandleFunc("/v1"+prefixPath+"/employee-level-matrix", EmployeeLevelEndpoint.EmployeeLevelEndpoint.EmployeeLevelMatrixEndpointWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-level-matrix/initiate", EmployeeLevelEndpoint.EmployeeLevelEndpoint.InitiateEmployeeLevelMatrixEndpoint).Methods("GET", "OPTIONS")

	//--------- EmployeeAllowance
	handler.HandleFunc("/v1"+prefixPath+"/employee-allowance", EmployeeAllowanceEndpoint.EmployeeAllowanceEndpoint.EmployeeAllowanceEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-allowance/initiate", EmployeeAllowanceEndpoint.EmployeeAllowanceEndpoint.InitiateEmployeeAllowanceEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-allowance/{ID}", EmployeeAllowanceEndpoint.EmployeeAllowanceEndpoint.EmployeeAllowanceEndpointWithParam).Methods("DELETE", "PUT", "OPTIONS")

	//--------- EmployeeBenefit
	handler.HandleFunc("/v1"+prefixPath+"/employee-benefit", EmployeeMasterBenefitEndpoint.EmployeeBenefitEndpoint.EmployeeBenefitEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-benefit/initiate", EmployeeMasterBenefitEndpoint.EmployeeBenefitEndpoint.InitiateEmployeeBenefitEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-benefit/{ID}", EmployeeMasterBenefitEndpoint.EmployeeBenefitEndpoint.EmployeeBenefitEndpointWithParam).Methods("DELETE", "PUT", "OPTIONS")

	//--------- EmployeeMatrix
	handler.HandleFunc("/v1"+prefixPath+"/employee-matrix", EmployeeFacilitiesActiveEndpoint.EmployeeMatrixEndpoint.EmployeeMatrixEndpointWithoutParam).Methods("GET", "POST", "DELETE", "PUT", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-matrix/detail", EmployeeFacilitiesActiveEndpoint.EmployeeMatrixEndpoint.EmployeeMatrixEndpointWithParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-matrix/initiate", EmployeeFacilitiesActiveEndpoint.EmployeeMatrixEndpoint.InitiateEmployeeMatrixEndpoint).Methods("GET", "OPTIONS")

	//--------- EmployeeHistory
	handler.HandleFunc("/v1"+prefixPath+"/employee-history", EmployeeHistoryEndpoint.EmployeeHistoryEndpoint.EmployeeHistoryWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-history/initiate", EmployeeHistoryEndpoint.EmployeeHistoryEndpoint.InitiateEmployeeHistoryEndpoint).Methods("GET", "OPTIONS")

	//--------- Parameter
	handler.HandleFunc("/v1"+prefixPath+"/parameter", ParameterEndpoint.ParameterEndpoint.ParameterEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	//handler.HandleFunc("/v1"+prefixPath+"/parameter/{ID}", ParameterEndpoint.ParameterEndpoint.ParameterEndpointWithoutParam).Methods("GET", "OPTIONS")

	//--------- Company
	handler.HandleFunc("/v1"+prefixPath+"/company", CompanyEndpoint.CompanyEndpoint.CompanyEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/company/initiate", CompanyEndpoint.CompanyEndpoint.InitiateCompanyEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/company/{ID}", CompanyEndpoint.CompanyEndpoint.CompanyEndpointWithParam).Methods("GET", "PUT", "DELETE", "OPTIONS")

	//--------- Backlog
	backlog := handler.PathPrefix("/v1" + prefixPath + "/backlog").Subrouter()
	backlog.HandleFunc("", BacklogEndpoint.BacklogEndpoint.BacklogEndpointWithoutParam).Methods("DELETE", "POST", "GET", "OPTIONS")
	backlog.HandleFunc("/initiate", BacklogEndpoint.BacklogEndpoint.InitiateBacklogEndpointWithoutParam).Methods("GET", "OPTIONS")
	backlog.HandleFunc("/import", BacklogEndpoint.BacklogEndpoint.ImportFileBacklogEndpointWithoutParam).Methods("POST", "OPTIONS")
	backlog.HandleFunc("/detail", BacklogEndpoint.BacklogEndpoint.BacklogDetailEndpointWithoutParam).Methods("POST", "GET", "OPTIONS")
	backlog.HandleFunc("/detail/update-status", BacklogEndpoint.BacklogEndpoint.UpdateStatusBacklogEndpointWithoutParam).Methods("POST", "GET", "OPTIONS")
	backlog.HandleFunc("/detail/status", BacklogEndpoint.BacklogEndpoint.StatusBacklogDetailEndpointWithoutParam).Methods("GET", "OPTIONS")
	backlog.HandleFunc("/detail/initiate", BacklogEndpoint.BacklogEndpoint.InitiateBacklogDetailEndpointWithoutParam).Methods("POST", "GET", "OPTIONS")
	backlog.HandleFunc("/detail/{ID}", BacklogEndpoint.BacklogEndpoint.BacklogDetailEndpointWithParam).Methods("POST", "GET", "PUT", "DELETE", "OPTIONS")
	backlog.HandleFunc("/ddl/project", BacklogEndpoint.BacklogEndpoint.DropDownSearchProjectEndpoint).Methods("GET", "OPTIONS")
	backlog.HandleFunc("/ddl/sprint", BacklogEndpoint.BacklogEndpoint.DropDownListSprintEndpoint).Methods("GET", "OPTIONS")
	backlog.HandleFunc("/ddl/tracker", BacklogEndpoint.BacklogEndpoint.DropDownListTrackerEndpoint).Methods("GET", "OPTIONS")
	backlog.HandleFunc("/ddl/employee", BacklogEndpoint.BacklogEndpoint.DropDownListEmployeeEndpoint).Methods("GET", "OPTIONS")
	backlog.HandleFunc("/ddl/department", BacklogEndpoint.BacklogEndpoint.DropDownListDepartmentEndpoint).Methods("GET", "OPTIONS")

	//--------- Report
	report := handler.PathPrefix("/v1" + prefixPath + "/report").Subrouter()
	report.HandleFunc("", ReportEndpoint.ReportEndpoint.ReportEndpointWithoutParam).Methods("GET", "OPTIONS")
	report.HandleFunc("/initiate", ReportEndpoint.ReportEndpoint.InitiateReportEndpoint).Methods("GET", "OPTIONS")
	report.HandleFunc("/download", ReportEndpoint.ReportEndpoint.DownloadReportEndpointWithoutParam).Methods("GET", "OPTIONS")
	report.HandleFunc("/paid", ReportEndpoint.ReportEndpoint.PaidPaymentReportEndpoint).Methods("GET", "OPTIONS")
	report.HandleFunc("/reset", ReportEndpoint.ReportEndpoint.HelperSetDefaultRedmineSprintPaid).Methods("POST", "OPTIONS")
	report.HandleFunc("/history", ReportEndpoint.ReportEndpoint.HelperGetListPaymentHistory).Methods("GET", "OPTIONS")
	report.HandleFunc("/history/{ID}", ReportEndpoint.ReportEndpoint.HelperViewDetailPaymentHistory).Methods("GET", "OPTIONS")
	report.HandleFunc("/ddl/employee", ReportEndpoint.ReportEndpoint.DropDownListEmployeeEndpoint).Methods("GET", "OPTIONS")
	report.HandleFunc("/ddl/department", ReportEndpoint.ReportEndpoint.DropDownListDepartmentEndpoint).Methods("GET", "OPTIONS")

	//---------- Update Sprint
	//handler.HandleFunc("/v1"+prefixPath+"/report-try", ReportEndpoint.ReportEndpoint.PaidPaymentReportEndpointTry).Methods("GET", "OPTIONS")

	//--------- Standar Manhour
	manhour := handler.PathPrefix("/v1" + prefixPath + "/manhour").Subrouter()
	manhour.HandleFunc("", StandarManhourEndpoint.StandarManhourEndpoint.StandarManhourWithoutParam).Methods("POST", "GET", "OPTIONS")
	manhour.HandleFunc("/initiate", StandarManhourEndpoint.StandarManhourEndpoint.InitiateStandarManhour).Methods("GET", "OPTIONS")
	manhour.HandleFunc("/{ID}", StandarManhourEndpoint.StandarManhourEndpoint.StandarManhourWithParam).Methods("PUT", "DELETE", "GET", "OPTIONS")

	//--------- Reset Gorp Migrations
	migration := handler.PathPrefix("/v1" + prefixPath + "/admin/migration").Subrouter()
	migration.HandleFunc("/reset", MigrationEndpoint.MigrationEndpoint.MigrationWithoutParam).Methods("POST", "OPTIONS")

	//--------- DDL Enum
	ddlCustom := handler.PathPrefix("/v1" + prefixPath + "/ddl/custom").Subrouter()
	ddlCustom.HandleFunc("", EnumEndpoint.EnumEndpoint.EnumEndpointWithoutParam).Methods("POST", "OPTIONS")

	//--------- Employee Contract
	contract := handler.PathPrefix("/v1" + prefixPath + "/contract").Subrouter()
	contract.HandleFunc("", EmployeeContractEndpoint.EmployeeContractEndpoint.EmployeeContractEndpointWithoutParam).Methods("POST", "GET", "OPTIONS")
	contract.HandleFunc("/initiate", EmployeeContractEndpoint.EmployeeContractEndpoint.InitiateEmployeeContractEndpoint).Methods("GET", "OPTIONS")
	contract.HandleFunc("/{ID}", EmployeeContractEndpoint.EmployeeContractEndpoint.EmployeeContractEndpointWithParam).Methods("GET", "PUT", "DELETE", "OPTIONS")

	//--------- Absent
	handler.HandleFunc("/v1"+prefixPath+"/absent", DashboardEndpoint.DashboardEndpoint.AbsentDashboardWithoutParam).Methods("POST", "GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/absent/initiate", DashboardEndpoint.DashboardEndpoint.InitiateAbsentDashboardEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/absent/period", DashboardEndpoint.DashboardEndpoint.GetListAbsentPeriodDashboardEndpoint).Methods("GET", "OPTIONS")

	//--------- Today Leave
	handler.HandleFunc("/v1"+prefixPath+"/today-leave", DashboardEndpoint.DashboardEndpoint.LeaveDashboardWithoutParam).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/today-leave/initiate", DashboardEndpoint.DashboardEndpoint.InitiateTodayLeaveEndpoint).Methods("GET", "OPTIONS")

	//--------- Dashboard
	dashboard := handler.PathPrefix("/v1" + prefixPath + "/dashboard").Subrouter()
	dashboard.HandleFunc("/reimbursement", DashboardEndpoint.DashboardEndpoint.ReimbursementDashboardPanelWithoutParam).Methods("GET", "OPTIONS")
	dashboard.HandleFunc("/leave", DashboardEndpoint.DashboardEndpoint.LeaveDashboardPanelWithoutParam).Methods("GET", "OPTIONS")
	dashboard.HandleFunc("/absent", DashboardEndpoint.DashboardEndpoint.AbsentDashboardPanelWithoutParam).Methods("GET", "OPTIONS")

	//--------- Employee Position
	handler.HandleFunc("/v1"+prefixPath+"/employee-position", EmployeePositionEndpoint.EmployeePositionEndpoint.EmployeePositionEndpointWithoutParam).Methods("GET", "POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-position/initiate", EmployeePositionEndpoint.EmployeePositionEndpoint.InitiateEmployeePositionEndpoint).Methods("GET", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/employee-position/{ID}", EmployeePositionEndpoint.EmployeePositionEndpoint.EmployeePositionEndpointWithParam).Methods("GET", "DELETE", "PUT", "OPTIONS")

	//--------- Download Report Annual Leave
	handler.HandleFunc("/v1"+prefixPath+"/report-leave/download", ReportAnnualLeaveEndpoint.ReportAnnualLeaveEndpoint.ReportAnnualLeaveEndpoint).Methods("POST", "OPTIONS")
	handler.HandleFunc("/v1"+prefixPath+"/report-leave/job", ReportAnnualLeaveEndpoint.ReportAnnualLeaveEndpoint.ReportAnnualLeaveJobListEndpoint).Methods("GET", "OPTIONS")

	handler.Use(Middleware)
	fmt.Print(http.ListenAndServe(config.ApplicationConfiguration.GetServerHost()+":"+strconv.Itoa(config.ApplicationConfiguration.GetServerPort()), handler))
}
