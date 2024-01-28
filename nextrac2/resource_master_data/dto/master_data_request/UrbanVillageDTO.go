package master_data_request

import (
	"nexsoft.co.id/nextrac2/dto/in"
	"time"
)

type UrbanVillageRequest struct {
	in.AbstractDTO
	ID             int64     `json:"id"`
	UpdatedAtStart time.Time `json:"updated_at_start"`
}
