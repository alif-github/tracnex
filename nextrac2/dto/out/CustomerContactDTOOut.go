package out

import "time"

type CustomerContactViewResponse struct {
	ID                 int64     `json:"id"`
	CustomerID         int64     `json:"customer_id"`
	MdbPersonProfileID int64     `json:"mdb_person_profile_id"`
	Nik                string    `json:"nik"`
	MdbPersonTitleID   int64     `json:"mdb_person_title_id"`
	PersonTitle        string    `json:"person_title"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Sex                string    `json:"sex"`
	Address            string    `json:"address"`
	Address2           string    `json:"address_2"`
	Address3           string    `json:"address_3"`
	Hamlet             string    `json:"hamlet"`
	Neighbourhood      string    `json:"neighbourhood"`
	ProvinceID         int64     `json:"province_id"`
	ProvinceName       string    `json:"province_name"`
	DistrictID         int64     `json:"district_id"`
	DistrictName       string    `json:"district_name"`
	Phone              string    `json:"phone"`
	Email              string    `json:"email"`
	MdbPositionID      int64     `json:"mdb_position_id"`
	PositionName       string    `json:"position_name"`
	Status             string    `json:"status"`
	CreatedBy          int64     `json:"created_by"`
	CreatedAt          time.Time `json:"created_at"`
	CreatedName        string    `json:"created_name"`
	UpdatedBy          int64     `json:"updated_by"`
	UpdatedAt          time.Time `json:"updated_at"`
	UpdatedName        string    `json:"updated_name"`
}
