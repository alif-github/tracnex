package out

import "time"

type StandarManhourResponse struct {
	ID         int64     `json:"id"`
	Case       string    `json:"case_name"`
	Department string    `json:"department"`
	Manhour    string    `json:"manhour"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type DetailStandarManhourResponse struct {
	ID           int64     `json:"id"`
	Case         string    `json:"case_name"`
	DepartmentID int64     `json:"department_id"`
	Department   string    `json:"department"`
	Manhour      string    `json:"manhour"`
	UpdatedAt    time.Time `json:"updated_at"`
}
