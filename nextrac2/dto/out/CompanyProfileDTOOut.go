package out

import "time"

type ViewCompanyProfileResponse struct {
	ID                int64     `json:"id"`
	Npwp              string    `json:"npwp"`
	Status            string    `json:"status"`
	CompanyTitleID    int64     `json:"company_title_id"`
	CompanyTitle      string    `json:"company_title"`
	Name              string    `json:"name"`
	CompanyParent     int64     `json:"company_parent"`
	CompanyParentName string    `json:"company_parent_name"`
	CustomerParentID  int64     `json:"customer_parent_id"`
	IsParentCompany   bool      `json:"is_parent_company"`
	Address           string    `json:"address"`
	Address2          string    `json:"address_2"`
	Address3          string    `json:"address_3"`
	Hamlet            string    `json:"hamlet"`
	Neighbourhood     string    `json:"neighbourhood"`
	ProvinceID        int64     `json:"province_id"`
	ProvinceName      string    `json:"province_name"`
	DistrictID        int64     `json:"district_id"`
	DistrictName      string    `json:"district_name"`
	SubDistrictID     int64     `json:"sub_district_id"`
	SubDistrictName   string    `json:"sub_district_name"`
	UrbanVillageID    int64     `json:"urban_village_id"`
	UrbanVillageName  string    `json:"urban_village_name"`
	PostalCodeID      int64     `json:"postal_code_id"`
	PostalCode        string    `json:"postal_code"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	Phone             string    `json:"phone"`
	Fax               string    `json:"fax"`
	Email             string    `json:"email"`
	AlternativeEmail  string    `json:"alternative_email"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedBy         int64     `json:"created_by"`
}
