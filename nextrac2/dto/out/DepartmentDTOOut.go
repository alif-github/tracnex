package out

import "time"

type ListDepartmentDTOOut struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ViewDepartmentResponse struct {
	ID             int64     `json:"id"`
	DepartmentName string    `json:"department_name"`
	Description    string    `json:"description"`
	UpdatedAt      time.Time `json:"updated_at"`
}
