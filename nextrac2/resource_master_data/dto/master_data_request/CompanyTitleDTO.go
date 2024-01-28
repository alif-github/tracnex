package master_data_request

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type CompanyTitleRequest struct {
	in.AbstractDTO
	ID           int64  `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	UpdatedAtStr string `json:"updated_at"`
	UpdatedAt    time.Time
}

func (input *CompanyTitleRequest) ValidateView() (err errorModel.ErrorModel) {
	if input.ID < 1 {
		err = errorModel.GenerateEmptyFieldError("CompanyTitleDTO.go", "ValidateView", constanta.ID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
