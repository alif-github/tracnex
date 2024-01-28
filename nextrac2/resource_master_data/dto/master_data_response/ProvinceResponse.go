package master_data_response

import (
	"nexsoft.co.id/nextrac2/dto/out"
	"time"
)

type GetListProvince struct {
	Nexsoft NexsoftListProvinceResponse `json:"nexsoft"`
}

type NexsoftListProvinceResponse struct {
	Header  out.Header          `json:"header"`
	Payload PayloadListProvince `json:"payload"`
}

type PayloadListProvince struct {
	Status out.StatusResponse `json:"status"`
	Data   ListProvince       `json:"data"`
	Other  interface{}        `json:"other"`
}

type ListProvince struct {
	Meta    interface{}        `json:"meta"`
	Content []ProvinceResponse `json:"content"`
}

//--------------View response

type ViewProvince struct {
	Nexsoft NexsoftViewProvinceResponse `json:"nexsoft"`
}

type NexsoftViewProvinceResponse struct {
	Header  out.Header          `json:"header"`
	Payload PayloadViewProvince `json:"payload"`
}

type PayloadViewProvince struct {
	Status out.StatusResponse `json:"status"`
	Data   DataViewProvince   `json:"data"`
	Other  interface{}        `json:"other"`
}

type DataViewProvince struct {
	Meta    interface{}      `json:"meta"`
	Content ProvinceResponse `json:"content"`
}

//-------------- Content

type ProvinceResponse struct {
	ID          int64     `json:"id"`
	CountryID   int64     `json:"country_id"`
	CountryName string    `json:"country_name"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	CreatedBy   int64     `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}
