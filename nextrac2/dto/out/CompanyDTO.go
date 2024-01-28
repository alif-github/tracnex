package out

import "time"

type CompanyForViewResponse struct {
	ID                   int64     `json:"id"`
	CompanyTitle         string    `json:"company_title"`
	CompanyName          string    `json:"company_name"`
	PhotoIcon            string    `json:"photo_icon"`
	Address              string    `json:"address"`
	Telephone            string    `json:"telephone"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            int64     `json:"updated_by"`
}

type CompanyDetailResponse struct {
	ID                   int64     `json:"id"`
	CompanyTitle         string    `json:"company_title"`
	CompanyName          string    `json:"company_name"`
	PhotoIcon            string    `json:"photo_icon"`
	Address              string    `json:"address"`
	Address2             string    `json:"address2"`
	Neighbourhood         string    `json:"neighbourhood"`
	Hamlet               string    `json:"hamlet"`
	ProvinceId           int64     `json:"province_id"`
	ProvinceName         string    `json:"province_name"`
	DistrictId           int64     `json:"district_id"`
	DistrictName         string    `json:"district_name"`
	SubDistrictId        int64     `json:"sub_district_id"`
	SubDistrictName      string    `json:"sub_district_name"`
	VillageId            int64     `json:"urban_village_id"`
	Village              string    `json:"urban_village_name"`
	PostalCodeId         int64     `json:"postal_code_id"`
	PostalCode           string    `json:"postal_code_name"`
	Longitude            string    `json:"longitude"`
	Latitude             string    `json:"latitude"`
	Telephone            string    `json:"telephone"`
	TelephoneAlternate   string    `json:"alternate_telephone"`
	Fax                   string   `json:"fax"`
	Email                string     `json:"email"`
	AlternateEmail       string     `json:"alternate_email"`
	Npwp                 string     `json:"npwp"`
	TaxName              string      `json:"tax_name"`
	TaxAddress           string     `json:"tax_address"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	UpdatedBy            int64     `json:"updated_by"`
}
