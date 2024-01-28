package applicationModel

type DefaultOperator struct {
	DataType string   `json:"data_type"`
	Operator []string `json:"operator"`
}

type MultiDeleteDataDTOOut struct {
	IDSuccess []string            `json:"id_success"`
	IDFailed  []FailedMultiDelete `json:"id_failed"`
}

type FailedMultiDelete struct {
	ID       int64  `json:"id"`
	DataID   string `json:"data_id"`
	CausedBy string `json:"caused_by"`
}

var GetListCountryValidOperator map[string]DefaultOperator
var GetListProvinceValidOperator map[string]DefaultOperator
var GetListDistrictValidOperator map[string]DefaultOperator
var GetListUserValidOperator map[string]DefaultOperator
var GetListSubDistrictValidOperator map[string]DefaultOperator
var GetListUrbanVillageValidOperator map[string]DefaultOperator
var GetListIslandValidOperator map[string]DefaultOperator
var GetListPositionValidOperator map[string]DefaultOperator
var GetListPersonTitleValidOperator map[string]DefaultOperator
var GetListPostalCodeValidOperator map[string]DefaultOperator
var GetListCompanyTitleValidOperator map[string]DefaultOperator
var GetListCompanyProfileValidOperator map[string]DefaultOperator
var GetListBankValidOperator map[string]DefaultOperator
var GetListPersonProfileValidOperator map[string]DefaultOperator
var GetListVendorChannelValidOperator map[string]DefaultOperator
var GetListGroupOfDistributorValidOperator map[string]DefaultOperator
var GetListBrandValidOperator map[string]DefaultOperator
var GetListBrandOwnerValidOperator map[string]DefaultOperator
var GetListProductGroupHierarchyValidOperator map[string]DefaultOperator
var GetListPrincipalValidOperator map[string]DefaultOperator
var GetListRoleValidOperator map[string]DefaultOperator
var GetListDataGroupValidOperator map[string]DefaultOperator
var GetListProductCategoryValidOperator map[string]DefaultOperator
var GetListProductBrandValidOperator map[string]DefaultOperator
var GetListProductKeyAccountValidOperator map[string]DefaultOperator
var GetListDataScopeValidOperator map[string]DefaultOperator
var GetListProductValidOperator map[string]DefaultOperator
var GetListProductHistoryValidOperator map[string]DefaultOperator
var GetListContactPersonValidOperator map[string]DefaultOperator
var GetListJobProcessValidOperator map[string]DefaultOperator
var GetListAuditMonitoringValidOperator map[string]DefaultOperator
var GetListCronSchedulerValidOperator map[string]DefaultOperator
var GetListHostServerValidOperator map[string]DefaultOperator
var GetListRunningCornValidOperator map[string]DefaultOperator
var GetListParameterValidOperator map[string]DefaultOperator
var GetListDetailClientMappingValidOperator map[string]DefaultOperator
var GetListClientMappingValidOperator map[string]DefaultOperator
var GetListNexsoftRoleValidOperator map[string]DefaultOperator
var GetListPKCEClientMappingValidOperator map[string]DefaultOperator
var GetListUserPKCENexmileValidOperator map[string]DefaultOperator
var GetListCustomerValidOperator map[string]DefaultOperator
var GetListRegistrationLogValidOperator map[string]DefaultOperator
var GetListCustomerGroupValidOperator map[string]DefaultOperator
var GetListSalesmanValidOperator map[string]DefaultOperator
var GetListCustomerCategoryValidOperator map[string]DefaultOperator
var GetListProductGroupValidOperator map[string]DefaultOperator
var GetListLicenseVariantValidOperator map[string]DefaultOperator
var GetListLicenseTypeValidOperator map[string]DefaultOperator
var GetListModuleValidOperator map[string]DefaultOperator
var GetListComponentValidOperator map[string]DefaultOperator
var GetListCustomerTitleValidOperator map[string]DefaultOperator
var GetListClientTypeValidOperator map[string]DefaultOperator
var GetListMasterCustomerValidOperator map[string]DefaultOperator
var GetListMasterDistributorValidOperator map[string]DefaultOperator
var GetListProductLicenseValidOperator map[string]DefaultOperator
var GetListUserActiveValidOperator map[string]DefaultOperator
var GetListLicenseConfigValidOperator map[string]DefaultOperator
var GetListUserLicenseValidOperator map[string]DefaultOperator
var GetListUserRegistrationAdminValidOperator map[string]DefaultOperator
var GetListWhiteListDeviceValidOperator map[string]DefaultOperator
var GetListDepartmentValidOperator map[string]DefaultOperator
var GetListEmployeeValidOperator map[string]DefaultOperator
var GetListParentBacklogValidOperator map[string]DefaultOperator
var DeleteParentBacklogValidOperator map[string]DefaultOperator
var GetListDetailBacklogValidOperator map[string]DefaultOperator
var GetListReportValidOperator map[string]DefaultOperator
var GetListStandarManhourValidOperator map[string]DefaultOperator
var GetListStatusValidOperator map[string]DefaultOperator
var GetListProjectValidOperator map[string]DefaultOperator
var GetListReportHistoryValidOperator map[string]DefaultOperator
var GetListEmployeeContractValidOperator map[string]DefaultOperator
var GetListEmployeeRequestHistoryValidOperator map[string]DefaultOperator
var GetListEmployeeApprovalHistoryValidOperator map[string]DefaultOperator
var GetListEmployeeLeaveTypesValidOperator map[string]DefaultOperator
var GetListEmployeeReimbursementTypesValidOperator map[string]DefaultOperator
var GetListAbsentValidOperator map[string]DefaultOperator
var GetListEmployeeLeaveValidOperator map[string]DefaultOperator
var GetListEmployeeReimbursementValidOperator map[string]DefaultOperator
var GetListEmployeePositionValidOperator map[string]DefaultOperator
var GetListEmployeeMatrixValidOperator map[string]DefaultOperator
var GetListEmployeeHistoryValidOperator map[string]DefaultOperator
var GetListTodaysLeaveValidOperator map[string]DefaultOperator
var GetListEmployeeNotificationValidOperator map[string]DefaultOperator
var GetListJobReportAnnualLeaveValidOperator map[string]DefaultOperator
var GetListEmployeeLevelValidOperator map[string]DefaultOperator
var GetListEmployeeGradeValidOperator map[string]DefaultOperator
var GetListInternalCompanyValidOperator map[string]DefaultOperator

func InitiateDefaultOperator() {
	GetListCountryValidOperator = make(map[string]DefaultOperator)
	GetListCountryValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListProvinceValidOperator = make(map[string]DefaultOperator)
	GetListProvinceValidOperator["country_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListProvinceValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListProvinceValidOperator["mdb_province_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListProvinceValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListProvinceValidOperator["code"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListDistrictValidOperator = make(map[string]DefaultOperator)
	GetListDistrictValidOperator["province_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListDistrictValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListDistrictValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListDistrictValidOperator["code"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListSubDistrictValidOperator = make(map[string]DefaultOperator)
	GetListSubDistrictValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListSubDistrictValidOperator["code"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListSubDistrictValidOperator["district_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListUrbanVillageValidOperator = make(map[string]DefaultOperator)
	GetListUrbanVillageValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListUrbanVillageValidOperator["code"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListUrbanVillageValidOperator["sub_district_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListIslandValidOperator = make(map[string]DefaultOperator)
	GetListIslandValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListPositionValidOperator = make(map[string]DefaultOperator)
	GetListPositionValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListPersonTitleValidOperator = make(map[string]DefaultOperator)
	GetListPersonTitleValidOperator["title"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListPostalCodeValidOperator = make(map[string]DefaultOperator)
	GetListPostalCodeValidOperator["code"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListPostalCodeValidOperator["urban_village_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListCompanyTitleValidOperator = make(map[string]DefaultOperator)
	GetListCompanyTitleValidOperator["title"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListCompanyProfileValidOperator = make(map[string]DefaultOperator)
	GetListCompanyProfileValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListCompanyProfileValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListCompanyProfileValidOperator["npwp"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListBankValidOperator = make(map[string]DefaultOperator)
	GetListBankValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListVendorChannelValidOperator = make(map[string]DefaultOperator)
	GetListVendorChannelValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListPersonProfileValidOperator = make(map[string]DefaultOperator)
	GetListPersonProfileValidOperator["first_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListPersonProfileValidOperator["email"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListPersonProfileValidOperator["phone"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListPersonProfileValidOperator["nik"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListBrandValidOperator = make(map[string]DefaultOperator)
	GetListBrandValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListBrandOwnerValidOperator = make(map[string]DefaultOperator)
	GetListBrandOwnerValidOperator["brand_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListBrandOwnerValidOperator["company_profile_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListProductGroupHierarchyValidOperator = make(map[string]DefaultOperator)
	GetListProductGroupHierarchyValidOperator["principal_id"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListProductGroupHierarchyValidOperator["code"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListProductGroupHierarchyValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListProductGroupHierarchyValidOperator["level"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListProductGroupHierarchyValidOperator["parent_id"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListUserValidOperator = make(map[string]DefaultOperator)
	GetListUserValidOperator["nt_username"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListUserValidOperator["full_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListUserValidOperator["phone"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListUserValidOperator["email"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListPrincipalValidOperator = make(map[string]DefaultOperator)
	GetListPrincipalValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListRoleValidOperator = make(map[string]DefaultOperator)
	GetListRoleValidOperator["role_id"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListRoleValidOperator["description"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListProductCategoryValidOperator = make(map[string]DefaultOperator)
	GetListProductCategoryValidOperator["code"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListProductCategoryValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListProductBrandValidOperator = make(map[string]DefaultOperator)
	GetListProductBrandValidOperator["code"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListProductBrandValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListProductKeyAccountValidOperator = make(map[string]DefaultOperator)
	GetListProductKeyAccountValidOperator["code"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListProductKeyAccountValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListProductKeyAccountValidOperator["level"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListProductKeyAccountValidOperator["parent_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListDataScopeValidOperator = make(map[string]DefaultOperator)
	GetListDataScopeValidOperator["scope"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListProductValidOperator = make(map[string]DefaultOperator)
	GetListProductValidOperator["product_id"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListProductValidOperator["product_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListProductValidOperator["product_group_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListGroupOfDistributorValidOperator = make(map[string]DefaultOperator)
	GetListGroupOfDistributorValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListProductHistoryValidOperator = make(map[string]DefaultOperator)
	GetListProductHistoryValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListProductHistoryValidOperator["product_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListContactPersonValidOperator = make(map[string]DefaultOperator)
	GetListContactPersonValidOperator["email"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListContactPersonValidOperator["phone"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListContactPersonValidOperator["connector"] = DefaultOperator{DataType: "enum", Operator: []string{"eq"}}
	GetListContactPersonValidOperator["parent_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListContactPersonValidOperator["nik"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListJobProcessValidOperator = make(map[string]DefaultOperator)
	GetListJobProcessValidOperator["status"] = DefaultOperator{DataType: "enum", Operator: []string{"eq"}}
	GetListJobProcessValidOperator["job_id"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListAuditMonitoringValidOperator = make(map[string]DefaultOperator)
	GetListAuditMonitoringValidOperator["table_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListAuditMonitoringValidOperator["primary_key"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListAuditMonitoringValidOperator["created_by"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListAuditMonitoringValidOperator["created_client"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListAuditMonitoringValidOperator["menu_code"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListCronSchedulerValidOperator = make(map[string]DefaultOperator)
	GetListCronSchedulerValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListHostServerValidOperator = make(map[string]DefaultOperator)
	GetListHostServerValidOperator["host_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListRunningCornValidOperator = make(map[string]DefaultOperator)
	GetListRunningCornValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListParameterValidOperator = make(map[string]DefaultOperator)
	GetListParameterValidOperator["permission"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListParameterValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListDetailClientMappingValidOperator = make(map[string]DefaultOperator)

	GetListClientMappingValidOperator = make(map[string]DefaultOperator)
	GetListClientMappingValidOperator["success_status_nexcloud"] = DefaultOperator{DataType: "bool", Operator: []string{"eq"}}
	GetListClientMappingValidOperator["client_id"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListClientMappingValidOperator["client_alias"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListNexsoftRoleValidOperator = make(map[string]DefaultOperator)
	GetListNexsoftRoleValidOperator["role_id"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListNexsoftRoleValidOperator["status"] = DefaultOperator{DataType: "enum", Operator: []string{"eq"}}

	GetListPKCEClientMappingValidOperator = make(map[string]DefaultOperator)
	GetListPKCEClientMappingValidOperator["parent_client_id"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListPKCEClientMappingValidOperator["client_alias"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListUserPKCENexmileValidOperator = make(map[string]DefaultOperator)
	GetListUserPKCENexmileValidOperator["username"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListUserPKCENexmileValidOperator["parent_client_id"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListCustomerValidOperator = make(map[string]DefaultOperator)
	GetListCustomerValidOperator["company_id"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListCustomerValidOperator["branch_id"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListCustomerValidOperator["company_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListRegistrationLogValidOperator = make(map[string]DefaultOperator)
	GetListRegistrationLogValidOperator["success_status_auth"] = DefaultOperator{DataType: "bool", Operator: []string{"eq"}}
	GetListRegistrationLogValidOperator["success_status_nexcloud"] = DefaultOperator{DataType: "bool", Operator: []string{"eq"}}
	GetListRegistrationLogValidOperator["resource"] = DefaultOperator{DataType: "string", Operator: []string{"eq", "like"}}

	GetListCustomerGroupValidOperator = make(map[string]DefaultOperator)
	GetListCustomerGroupValidOperator["customer_group_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListCustomerGroupValidOperator["customer_group_id"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListDataGroupValidOperator = make(map[string]DefaultOperator)
	GetListDataGroupValidOperator["group_id"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListDataGroupValidOperator["description"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListDataGroupValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListCustomerCategoryValidOperator = make(map[string]DefaultOperator)
	GetListCustomerCategoryValidOperator["customer_category_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListCustomerCategoryValidOperator["customer_category_id"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListProductGroupValidOperator = make(map[string]DefaultOperator)
	GetListProductGroupValidOperator["product_group_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListProductGroupValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListSalesmanValidOperator = make(map[string]DefaultOperator)
	GetListSalesmanValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListSalesmanValidOperator["first_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListLicenseVariantValidOperator = make(map[string]DefaultOperator)
	GetListLicenseVariantValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListLicenseVariantValidOperator["license_variant_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListLicenseTypeValidOperator = make(map[string]DefaultOperator)
	GetListLicenseTypeValidOperator["license_type_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListLicenseTypeValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListModuleValidOperator = make(map[string]DefaultOperator)
	GetListModuleValidOperator["module_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListModuleValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListComponentValidOperator = make(map[string]DefaultOperator)
	GetListComponentValidOperator["component_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListComponentValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListCustomerTitleValidOperator = make(map[string]DefaultOperator)
	GetListCustomerTitleValidOperator["title"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListClientTypeValidOperator = make(map[string]DefaultOperator)
	GetListClientTypeValidOperator["client_type_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListClientTypeValidOperator["client_type"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}

	GetListMasterCustomerValidOperator = make(map[string]DefaultOperator)
	GetListMasterCustomerValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListMasterCustomerValidOperator["customer_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListMasterCustomerValidOperator["distributor_of"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListMasterCustomerValidOperator["province_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListMasterCustomerValidOperator["district_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListMasterCustomerValidOperator["customer_category_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListMasterCustomerValidOperator["customer_group_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListMasterCustomerValidOperator["salesman_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListMasterDistributorValidOperator = make(map[string]DefaultOperator)
	GetListMasterDistributorValidOperator["license_variant"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListMasterDistributorValidOperator["client_type"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListMasterDistributorValidOperator["updated_at"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListMasterDistributorValidOperator["updated_at_start"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListMasterDistributorValidOperator["updated_at_end"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListProductLicenseValidOperator = make(map[string]DefaultOperator)
	GetListProductLicenseValidOperator["customer_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListProductLicenseValidOperator["customer_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListLicenseConfigValidOperator = make(map[string]DefaultOperator)
	GetListLicenseConfigValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListLicenseConfigValidOperator["parent_customer_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListLicenseConfigValidOperator["customer_name"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListLicenseConfigValidOperator["distributor_of"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListLicenseConfigValidOperator["product_valid_from"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListLicenseConfigValidOperator["product_valid_thru"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListLicenseConfigValidOperator["product_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListLicenseConfigValidOperator["client_type_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListLicenseConfigValidOperator["province_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListLicenseConfigValidOperator["district_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListLicenseConfigValidOperator["unique_id_1"] = DefaultOperator{DataType: "char", Operator: []string{"like", "eq"}}
	GetListLicenseConfigValidOperator["license_status"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListLicenseConfigValidOperator["customer_group_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListLicenseConfigValidOperator["customer_category_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListLicenseConfigValidOperator["salesman_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListLicenseConfigValidOperator["allow_activation"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListUserLicenseValidOperator = make(map[string]DefaultOperator)
	GetListUserLicenseValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListUserLicenseValidOperator["customer_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListUserActiveValidOperator = make(map[string]DefaultOperator)
	GetListUserActiveValidOperator["username"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListUserRegistrationAdminValidOperator = make(map[string]DefaultOperator)
	GetListUserRegistrationAdminValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListUserRegistrationAdminValidOperator["company_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListWhiteListDeviceValidOperator = make(map[string]DefaultOperator)
	GetListWhiteListDeviceValidOperator["device"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListWhiteListDeviceValidOperator["description"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListDepartmentValidOperator = make(map[string]DefaultOperator)
	GetListDepartmentValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListDepartmentValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListEmployeeValidOperator = make(map[string]DefaultOperator)
	GetListEmployeeValidOperator["nik"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListEmployeeValidOperator["redmine_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListEmployeeValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListEmployeeValidOperator["department"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListEmployeeValidOperator["department_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListEmployeeValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListEmployeeValidOperator["id_card"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListEmployeeValidOperator["first_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListEmployeeValidOperator["is_timesheet"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListEmployeeValidOperator["is_redmine_check"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListParentBacklogValidOperator = make(map[string]DefaultOperator)
	GetListParentBacklogValidOperator["sprint"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListDetailBacklogValidOperator = make(map[string]DefaultOperator)
	GetListDetailBacklogValidOperator["sprint"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListDetailBacklogValidOperator["pic"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListDetailBacklogValidOperator["redmine_number"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListReportValidOperator = make(map[string]DefaultOperator)
	GetListReportValidOperator["sprint"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListReportValidOperator["department"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListReportValidOperator["id"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	DeleteParentBacklogValidOperator = make(map[string]DefaultOperator)
	DeleteParentBacklogValidOperator["sprint"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListStandarManhourValidOperator = make(map[string]DefaultOperator)
	GetListStandarManhourValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListStandarManhourValidOperator["case"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListStandarManhourValidOperator["department_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListProjectValidOperator = make(map[string]DefaultOperator)
	GetListProjectValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListStatusValidOperator = make(map[string]DefaultOperator)
	GetListStatusValidOperator["department_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListReportHistoryValidOperator = make(map[string]DefaultOperator)
	GetListReportHistoryValidOperator["department_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListReportHistoryValidOperator["created_at"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListEmployeeContractValidOperator = make(map[string]DefaultOperator)
	GetListEmployeeContractValidOperator["contract_no"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListEmployeeContractValidOperator["employee_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListEmployeeRequestHistoryValidOperator = make(map[string]DefaultOperator)

	GetListEmployeeApprovalHistoryValidOperator = make(map[string]DefaultOperator)
	GetListEmployeeApprovalHistoryValidOperator["status"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListEmployeeLeaveTypesValidOperator = make(map[string]DefaultOperator)
	GetListEmployeeReimbursementTypesValidOperator = make(map[string]DefaultOperator)

	GetListAbsentValidOperator = make(map[string]DefaultOperator)
	GetListAbsentValidOperator["id_card"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListAbsentValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListAbsentValidOperator["period"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListEmployeeLeaveValidOperator = make(map[string]DefaultOperator)
	GetListEmployeeLeaveValidOperator["first_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListEmployeeLeaveValidOperator["last_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListEmployeeLeaveValidOperator["el.type"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListEmployeeReimbursementValidOperator = make(map[string]DefaultOperator)
	GetListEmployeeReimbursementValidOperator["first_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListEmployeeReimbursementValidOperator["last_name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListEmployeeReimbursementValidOperator["id_card"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListEmployeePositionValidOperator = make(map[string]DefaultOperator)
	GetListEmployeePositionValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListEmployeePositionValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}

	GetListEmployeeHistoryValidOperator = make(map[string]DefaultOperator)
	GetListEmployeeHistoryValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListEmployeeHistoryValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListEmployeeHistoryValidOperator["primary_key"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListTodaysLeaveValidOperator = make(map[string]DefaultOperator)
	GetListTodaysLeaveValidOperator["id_card"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListTodaysLeaveValidOperator["name"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListTodaysLeaveValidOperator["department"] = DefaultOperator{DataType: "char", Operator: []string{"eq", "like"}}
	GetListTodaysLeaveValidOperator["type"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListJobReportAnnualLeaveValidOperator = make(map[string]DefaultOperator)
	GetListJobReportAnnualLeaveValidOperator["job_id"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListJobReportAnnualLeaveValidOperator["category"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}
	GetListJobReportAnnualLeaveValidOperator["created_at"] = DefaultOperator{DataType: "char", Operator: []string{"eq"}}

	GetListEmployeeLevelValidOperator = make(map[string]DefaultOperator)
	GetListEmployeeLevelValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListEmployeeGradeValidOperator = make(map[string]DefaultOperator)
	GetListEmployeeGradeValidOperator["id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}
	GetListEmployeeGradeValidOperator["level_id"] = DefaultOperator{DataType: "number", Operator: []string{"eq"}}

	GetListInternalCompanyValidOperator = make(map[string]DefaultOperator)
	GetListInternalCompanyValidOperator["company_name"] = DefaultOperator{DataType: "char", Operator: []string{"like"}}

	GetListEmployeeNotificationValidOperator = make(map[string]DefaultOperator)
}
