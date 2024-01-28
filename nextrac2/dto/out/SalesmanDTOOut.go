package out

import "time"

type ListSalesman struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	District  string    `json:"district"`
	Province  string    `json:"province"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ViewSalesman struct {
	ID            int64     `json:"id"`
	PersonTitleID int64     `json:"person_title_id"`
	PersonTitle   string    `json:"person_title"`
	Sex           string    `json:"sex"`
	Nik           string    `json:"nik"`
	Address       string    `json:"address"`
	Hamlet        string    `json:"hamlet"`
	Neighbourhood string    `json:"neighbourhood"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	ProvinceID    int64     `json:"province_id"`
	Province      string    `json:"province"`
	DistrictID    int64     `json:"district_id"`
	District      string    `json:"district"`
	Phone         string    `json:"phone"`
	Email         string    `json:"email"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UpdatedBy     int64     `json:"updated_by"`
	UpdatedName   string    `json:"updated_name"`
}

type GenderSalesman struct {
	Code       string `json:"code"`
	GenderName string `json:"gender_name"`
}
