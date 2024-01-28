package out

import "time"

//--------------Get List Response

type GetListDistrict struct {
	Nexsoft NexsoftListDistrictResponse `json:"nexsoft"`
}

type NexsoftListDistrictResponse struct {
	Header  Header              `json:"header"`
	Payload PayloadListDistrict `json:"payload"`
}

type PayloadListDistrict struct {
	Status StatusResponse `json:"status"`
	Data   ListDistrict   `json:"data"`
	Other  interface{}    `json:"other"`
}

type ListDistrict struct {
	Meta    interface{}        `json:"meta"`
	Content []DistrictResponse `json:"content"`
}

//--------------View response

type ViewDistrict struct {
	Nexsoft NexsoftViewDistrictResponse `json:"nexsoft"`
}

type NexsoftViewDistrictResponse struct {
	Header  Header              `json:"header"`
	Payload PayloadViewDistrict `json:"payload"`
}

type PayloadViewDistrict struct {
	Status StatusResponse   `json:"status"`
	Data   DataViewDistrict `json:"data"`
	Other  interface{}      `json:"other"`
}

type DataViewDistrict struct {
	Meta    interface{}      `json:"meta"`
	Content DistrictResponse `json:"content"`
}

//-------------- Content

type DistrictResponse struct {
	ID           int64     `json:"id"`
	ProvinceID   int64     `json:"province_id"`
	ProvinceName string    `json:"province_name"`
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	CreatedBy    string    `json:"created_by"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type DistrictLocalResponse struct {
	ID         int64  `json:"id"`
	ProvinceID int64  `json:"province_id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
}
