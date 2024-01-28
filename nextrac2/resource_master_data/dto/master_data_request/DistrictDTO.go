package master_data_request

import (
	"nexsoft.co.id/nextrac2/dto/in"
	"time"
)

type DistrictRequest struct {
	in.AbstractDTO
	ID             int64     `json:"id"`
	CountryIDList  []int64   `json:"country_id_list"`
	UpdatedAtStart time.Time `json:"updated_at_start"`
}
