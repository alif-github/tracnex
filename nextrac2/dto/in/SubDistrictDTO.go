package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type SubDistrictRequest struct {
	AbstractDTO
	ID           int64  `json:"id"`
	DistrictID   int64  `json:"district_id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	UpdatedAtStr string `json:"updated_at"`
	UpdatedAt    time.Time
}

func (input SubDistrictRequest) ValidateView() (err errorModel.ErrorModel) {
	funcName := "ValidateView"
	if input.ID < 1 {
		err = errorModel.GenerateEmptyFieldError(SubDistrictDTOFileName, funcName, constanta.ID)
		return
	}
	return
}
