package master_data_response

import (
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"time"
)

type CompanyProfileResponse struct {
	ID               int64       `json:"id"`
	CompanyTitleID   int64       `json:"company_title_id"`
	CompanyTitle     string      `json:"company_title"`
	Npwp             string      `json:"npwp"`
	Name             string      `json:"name"`
	Address1         string      `json:"address_1"`
	CountryID        int64       `json:"country_id"`
	CountryName      string      `json:"country_name"`
	DistrictID       int64       `json:"district_id"`
	DistrictName     string      `json:"district_name"`
	SubDistrictID    int64       `json:"sub_district_id"`
	SubDistrictName  string      `json:"sub_district_name"`
	UrbanVillageID   int64       `json:"urban_village_id"`
	UrbanVillageName string      `json:"urban_village_name"`
	PostalCodeID     int64       `json:"postal_code_id"`
	PostalCode       string      `json:"postal_code"`
	IslandID         int64       `json:"island_id"`
	IslandName       string      `json:"island_name"`
	Phone            string      `json:"phone"`
	Fax              string      `json:"fax"`
	Email            string      `json:"email"`
	AlternativeEmail string      `json:"alternative_email"`
	Logo             []PhotoList `json:"logo"`
	CompanyParent    int64       `json:"company_parent"`
	Role             string      `json:"role"`
	UpdatedAt        time.Time   `json:"updated_at"`
	CreatedBy        int64       `json:"created_by"`
}

type PhotoList struct {
	ID       int64  `json:"id"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	Host     string `json:"host"`
}

type ViewCompanyProfileResponse struct {
	ID                           int64                         `json:"id"`
	CompanyTitleID               int64                         `json:"company_title_id"`
	CompanyTitle                 string                        `json:"company_title"`
	NPWP                         string                        `json:"npwp"`
	Name                         string                        `json:"name"`
	AnonymizedName               string                        `json:"anonymized_name"`
	Address1                     string                        `json:"address_1"`
	Address2                     string                        `json:"address_2"`
	Address3                     string                        `json:"address_3"`
	Hamlet                       string                        `json:"hamlet"`
	Neighbourhood                string                        `json:"neighbourhood"`
	CountryID                    int64                         `json:"country_id"`
	CountryName                  string                        `json:"country_name"`
	ProvinceID                   int64                         `json:"province_id"`
	ProvinceName                 string                        `json:"province_name"`
	DistrictID                   int64                         `json:"district_id"`
	DistrictName                 string                        `json:"district_name"`
	SubDistrictID                int64                         `json:"sub_district_id"`
	SubDistrictName              string                        `json:"sub_district_name"`
	UrbanVillageID               int64                         `json:"urban_village_id"`
	UrbanVillageName             string                        `json:"urban_village_name"`
	PostalCodeID                 int64                         `json:"postal_code_id"`
	PostalCode                   string                        `json:"postal_code"`
	IslandID                     int64                         `json:"island_id"`
	IslandName                   string                        `json:"island_name"`
	Latitude                     float64                       `json:"latitude"`
	Longitude                    float64                       `json:"longitude"`
	Phone                        string                        `json:"phone"`
	Fax                          string                        `json:"fax"`
	Email                        string                        `json:"email"`
	AlternativeEmail             string                        `json:"alternative_email"`
	IsPKP                        string                        `json:"is_pkp"`
	IsBUMN                       string                        `json:"is_bumn"`
	IsPBF                        string                        `json:"is_pbf"`
	PBFLicenseNumber             string                        `json:"pbf_license_number"`
	PBFLicenseEndDate            time.Time                     `json:"pbf_license_end_date"`
	PBFTrusteeSipaNumber         string                        `json:"pbf_trustee_sipa_number"`
	PBFSipaLicenseEndDate        time.Time                     `json:"pbf_sipa_license_end_date"`
	AdditionalInfo               []model.AdditionalInformation `json:"additional_info"`
	CompanyParent                int64                         `json:"company_parent"`
	CompanyParentName            string                        `json:"company_parent_name"`
	Status                       string                        `json:"status"`
	UpdatedAt                    time.Time                     `json:"updated_at"`
	CreatedBy                    int64                         `json:"created_by"`
	CompanyUserDefinedCategory1  int64                         `json:"company_user_defined_category_1"`
	CompanyUserDefinedCategory2  int64                         `json:"company_user_defined_category_2"`
	CompanyUserDefinedCategory3  int64                         `json:"company_user_defined_category_3"`
	CompanyUserDefinedCategory4  int64                         `json:"company_user_defined_category_4"`
	CompanyUserDefinedCategory5  int64                         `json:"company_user_defined_category_5"`
	CompanyUserDefinedCategory6  int64                         `json:"company_user_defined_category_6"`
	CompanyUserDefinedCategory7  int64                         `json:"company_user_defined_category_7"`
	CompanyUserDefinedCategory8  int64                         `json:"company_user_defined_category_8"`
	CompanyUserDefinedCategory9  int64                         `json:"company_user_defined_category_9"`
	CompanyUserDefinedCategory10 int64                         `json:"company_user_defined_category_10"`
	CompanyUserDefinedCategory11 int64                         `json:"company_user_defined_category_11"`
	CompanyUserDefinedCategory12 int64                         `json:"company_user_defined_category_12"`
	CompanyUserDefinedCategory13 int64                         `json:"company_user_defined_category_13"`
	CompanyUserDefinedCategory14 int64                         `json:"company_user_defined_category_14"`
	CompanyUserDefinedCategory15 int64                         `json:"company_user_defined_category_15"`
	CompanyUserDefinedCategory16 int64                         `json:"company_user_defined_category_16"`
	CompanyUserDefinedCategory17 int64                         `json:"company_user_defined_category_17"`
	CompanyUserDefinedCategory18 int64                         `json:"company_user_defined_category_18"`
	CompanyUserDefinedCategory19 int64                         `json:"company_user_defined_category_19"`
	CompanyUserDefinedCategory20 int64                         `json:"company_user_defined_category_20"`
	Role                         string                        `json:"role"`
}

type MasterDataInsertedIDResponse struct {
	ID int64 `json:"id"`
}
