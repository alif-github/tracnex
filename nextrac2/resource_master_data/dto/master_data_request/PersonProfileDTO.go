package master_data_request

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type PersonProfileGetListRequest struct {
	in.AbstractDTO
	ID       int64  `json:"id"`
	NIK      string `json:"nik"`
	NPWP     string `json:"npwp"`
	FistName string `json:"fist_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

type PersonProfileWriteRequest struct {
	ID                          int64     `json:"id"`
	PersonTitleID               int64     `json:"person_title_id"`
	Nik                         string    `json:"nik"`
	Npwp                        string    `json:"npwp"`
	FirstName                   string    `json:"first_name"`
	LastName                    string    `json:"last_name"`
	Sex                         string    `json:"sex"`
	Address1                    string    `json:"address_1"`
	Address2                    string    `json:"address_2"`
	Address3                    string    `json:"address_3"`
	Hamlet                      string    `json:"hamlet"`
	Neighbourhood               string    `json:"neighbourhood"`
	CountryID                   int64     `json:"country_id"`
	ProvinceID                  int64     `json:"province_id"`
	DistrictID                  int64     `json:"district_id"`
	SubDistrictID               int64     `json:"sub_district_id"`
	UrbanVillageID              int64     `json:"urban_village_id"`
	PostalCodeID                int64     `json:"postal_code_id"`
	IslandID                    int64     `json:"island_id"`
	PhoneCountryCode            string    `json:"phone_country_code"`
	Phone                       string    `json:"phone"`
	AlternativePhoneCountryCode string    `json:"alternative_phone_country_code"`
	AlternativePhone            string    `json:"alternative_phone"`
	Email                       string    `json:"email"`
	AlternativeEmail            string    `json:"alternative_email"`
	BirthPlace                  string    `json:"birth_place"`
	BirthDateStr                string    `json:"birth_date"`
	Occupation                  string    `json:"occupation"`
	SpouseFirstName             string    `json:"spouse_first_name"`
	SpouseLastName              string    `json:"spouse_last_name"`
	SpouseBirthDateStr          string    `json:"spouse_birth_date"`
	AnniversaryDateStr          string    `json:"anniversary_date"`
	EducationHistory            string    `json:"education_history"`
	JobExperienceHistory        string    `json:"job_experience_history"`
	Remark                      string    `json:"remark"`
	Status                      string    `json:"status"`
	UpdatedAt                   time.Time `json:"updated_at"`
	DeletedPhotoID              []int64   `json:"deleted_photo_id"`
	AuthUserID                  int64     `json:"auth_user_id"`
	ClientID                    string    `json:"client_id"`
	ListID                      []int64   `json:"list_id"`
	BirthDate                   time.Time
	SpouseBirthDate             time.Time
	AnniversaryDate             time.Time
	PersonProfileIDList         []int64 `json:"person_profile_id_list"`
	//AdditionalInfo              []model.AdditionalInformation `json:"additional_info"`
}

func (input *PersonProfileGetListRequest) ValidateView() (err errorModel.ErrorModel) {
	if input.ID < 1 {
		err = errorModel.GenerateEmptyFieldError("PersonProfileDTO.go", "ValidateView", constanta.ID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
