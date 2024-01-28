package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type PostalCodeRequest struct {
	AbstractDTO
	ID             int64  `json:"id"`
	UrbanVillageID int64  `json:"urban_village_id"`
	Code           string `json:"code"`
	Status         string `json:"status"`
	UpdatedAtStr   string `json:"updated_at"`
	UpdatedAt      time.Time
}

func (input PostalCodeRequest) ValidateView() (err errorModel.ErrorModel) {
	funcName := "ValidateView"
	if input.ID < 1 {
		err = errorModel.GenerateEmptyFieldError(PostalCodeDTOFileName, funcName, constanta.ID)
		return
	}
	return
}
