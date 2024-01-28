package in

type DistrictRequest struct {
	AbstractDTO
	ProvinceID int64  `json:"province_id"`
	ID         int64  `json:"id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
}
