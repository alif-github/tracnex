package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type MultiDeleteRequest struct {
	AbstractDTO
	DeletedID []DeletedID `json:"deleted_id"`
}

type DeletedID struct {
	ID           int64  `json:"id"`
	UpdatedAtStr string `json:"updated_at"`
	UpdatedAt    time.Time
}

func (input *MultiDeleteRequest) ValidateMultiDelete(index int) (err errorModel.ErrorModel) {
	fileName := "MultiDeleteDTO.go"
	funcName := "ValidateMultiDelete"

	if input.DeletedID[index].ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	input.DeletedID[index].UpdatedAt, err = TimeStrToTime(input.DeletedID[index].UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}
