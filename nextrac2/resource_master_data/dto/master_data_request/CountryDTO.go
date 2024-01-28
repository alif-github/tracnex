package master_data_request

import (
	"nexsoft.co.id/nextrac2/dto/in"
	"time"
)

type CountryRequest struct {
	in.AbstractDTO
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	UpdatedAtStart time.Time `json:"updated_at_start"`
}
