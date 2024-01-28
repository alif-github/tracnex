package master_data_request

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type PositionGetListRequest struct {
	in.AbstractDTO
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	UpdatedAtStr string `json:"updated_at"`
	UpdatedAt    time.Time
}

func (input *PositionGetListRequest) ValidateView() (err errorModel.ErrorModel) {
	if input.ID < 1 {
		err = errorModel.GenerateEmptyFieldError("PositionDTO.go", "ValidateView", constanta.ID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
