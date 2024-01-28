package out

type ViewDetailPersonProfileResponse struct {
	ID              int64  `json:"id"`
	Nik             string `json:"nik"`
	PersonTitleID   int64  `json:"person_title_id"`
	PersonTitleName string `json:"person_title_name"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Sex             string `json:"sex"`
	Address         string `json:"address"`
	Address2        string `json:"address_2"`
	Address3        string `json:"address_3"`
	Hamlet          string `json:"hamlet"`
	Neighbourhood   string `json:"neighbourhood"`
	ProvinceID      int64  `json:"province_id"`
	ProvinceName    string `json:"province_name"`
	DistrictID      int64  `json:"district_id"`
	DistrictName    string `json:"district_name"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
}
