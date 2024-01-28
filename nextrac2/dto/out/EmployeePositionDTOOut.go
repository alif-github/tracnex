package out

import "time"

type ListEmployeePosition struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type DetailEmployeePosition struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CompanyName string    `json:"company_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedBy   int64     `json:"created_by"`
	CreatedName string    `json:"created_name"`
	UpdatedBy   int64     `json:"updated_by"`
	UpdatedName string    `json:"updated_name"`
}
