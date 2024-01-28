package master_data_request

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type CompanyProfileGetListRequest struct {
	in.AbstractDTO
	ID   int64  `json:"id"`
	NPWP string `json:"npwp"`
	Name string `json:"name"`
}

type CompanyProfileWriteRequest struct {
	ID                  int64     `json:"id"`
	CompanyTitleID      int64     `json:"company_title_id"`
	NPWP                string    `json:"npwp"`
	Name                string    `json:"name"`
	Address1            string    `json:"address_1"`
	Address2            string    `json:"address_2"`
	Address3            string    `json:"address_3"`
	Hamlet              string    `json:"hamlet"`
	Neighbourhood       string    `json:"neighbourhood"`
	ProvinceID          int64     `json:"province_id"`
	CountryID           int64     `json:"country_id"`
	DistrictID          int64     `json:"district_id"`
	SubDistrictID       int64     `json:"sub_district_id"`
	UrbanVillageID      int64     `json:"urban_village_id"`
	PostalCodeID        int64     `json:"postal_code_id"`
	Latitude            float64   `json:"latitude"`
	Longitude           float64   `json:"longitude"`
	PhoneCountryCode    string    `json:"phone_country_code"`
	Phone               string    `json:"phone"`
	PhoneFaxCountryCode string    `json:"phone_fax_country_code"`
	Fax                 string    `json:"fax"`
	Email               string    `json:"email"`
	AlternativeEmail    string    `json:"alternative_email"`
	CompanyParent       int64     `json:"company_parent"`
	UpdatedAtStr        string    `json:"-"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (input *CompanyProfileGetListRequest) ValidateView() (err errorModel.ErrorModel) {
	if input.ID < 1 {
		err = errorModel.GenerateEmptyFieldError("CompanyProfileDTO.go", "ValidateView", constanta.ID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
