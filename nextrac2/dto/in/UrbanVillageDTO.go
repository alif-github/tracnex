package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type UrbanVillageRequest struct {
	AbstractDTO
	ID            int64  `json:"id"`
	SubDistrictID int64  `json:"sub_district_id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Status        string `json:"status"`
	UpdatedAtStr  string `json:"updated_at"`
	UpdatedAt     time.Time
}

func (input UrbanVillageRequest) ValidateView() (err errorModel.ErrorModel) {
	funcName := "ValidateView"
	if input.ID < 1 {
		err = errorModel.GenerateEmptyFieldError(UrbanVillageDTOFileName, funcName, constanta.ID)
		return
	}
	return
}

func (input UrbanVillageRequest) ValidateGetList() (err errorModel.ErrorModel) {
	funcName := "ValidateGetList"

	if input.SubDistrictID < 1 {
		err = errorModel.GenerateEmptyFieldError(UrbanVillageDTOFileName, funcName, constanta.SubDistrictID)
		return
	}

	return
}