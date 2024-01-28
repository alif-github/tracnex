package master_data_response

import (
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"time"
)

type GetListPersonProfileResponse struct {
	ID              int64                         `json:"id"`
	AuthUserID      int64                         `json:"auth_user_id"`
	ClientID        string                        `json:"client_id"`
	Username        string                        `json:"username"`
	AuthEmail       string                        `json:"auth_email"`
	AuthPhone       string                        `json:"auth_phone"`
	PersonTitleID   int64                         `json:"person_title_id"`
	PersonTitleName string                        `json:"person_title_name"`
	Nik             string                        `json:"nik"`
	Npwp            string                        `json:"npwp"`
	FirstName       string                        `json:"first_name"`
	LastName        string                        `json:"last_name"`
	Sex             string                        `json:"sex"`
	Email           string                        `json:"email"`
	Phone           string                        `json:"phone"`
	AdditionalInfo  []model.AdditionalInformation `json:"additional_info"`
	Photo           []PhotoList                   `json:"photo"`
	Status          string                        `json:"status"`
	CreatedBy       int64                         `json:"created_by"`
	UpdatedAt       time.Time                     `json:"updated_at"`
	AuthUpdatedAt   time.Time                     `json:"auth_updated_at"`
}

type ViewPersonProfileResponse struct {
	ID                   int64                         `json:"id"`
	AuthUserID           int64                         `json:"auth_user_id"`
	ClientID             string                        `json:"client_id"`
	PersonTitleID        int64                         `json:"person_title_id"`
	PersonTitleName      string                        `json:"person_title_name"`
	Nik                  string                        `json:"nik"`
	Npwp                 string                        `json:"npwp"`
	FirstName            string                        `json:"first_name"`
	LastName             string                        `json:"last_name"`
	Sex                  string                        `json:"sex"`
	Address1             string                        `json:"address_1"`
	Address2             string                        `json:"address_2"`
	Address3             string                        `json:"address_3"`
	Hamlet               string                        `json:"hamlet"`
	Neighbourhood        string                        `json:"neighbourhood"`
	CountryID            int64                         `json:"country_id"`
	CountryName          string                        `json:"country_name"`
	ProvinceID           int64                         `json:"province_id"`
	ProvinceName         string                        `json:"province_name"`
	DistrictID           int64                         `json:"district_id"`
	DistrictName         string                        `json:"district_name"`
	SubDistrictID        int64                         `json:"sub_district_id"`
	SubDistrictName      string                        `json:"sub_district_name"`
	UrbanVillageID       int64                         `json:"urban_village_id"`
	UrbanVillageName     string                        `json:"urban_village_name"`
	PostalCodeID         int64                         `json:"postal_code_id"`
	PostalCode           string                        `json:"postal_code"`
	IslandID             int64                         `json:"island_id"`
	IslandName           string                        `json:"island_name"`
	Phone                string                        `json:"phone"`
	AlternativePhone     string                        `json:"alternative_phone"`
	Email                string                        `json:"email"`
	Photo                []PhotoList                   `json:"photo"`
	AlternativeEmail     string                        `json:"alternative_email"`
	BirthPlace           string                        `json:"birth_place"`
	BirthDate            time.Time                     `json:"birth_date"`
	Occupation           string                        `json:"occupation"`
	SpouseFirstName      string                        `json:"spouse_first_name"`
	SpouseLastName       string                        `json:"spouse_last_name"`
	SpouseBirthDate      time.Time                     `json:"spouse_birth_date"`
	AnniversaryDate      time.Time                     `json:"anniversary_date"`
	EducationHistory     string                        `json:"education_history"`
	JobExperienceHistory string                        `json:"job_experience_history"`
	Remark               string                        `json:"remark"`
	Status               string                        `json:"status"`
	AdditionalInfo       []model.AdditionalInformation `json:"additional_info"`
	CreatedBy            int64                         `json:"created_by"`
	UpdatedAt            time.Time                     `json:"updated_at"`
}