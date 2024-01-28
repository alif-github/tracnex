package out

import "time"

type ViewParameterDTOOut struct {
	ID         int64     `json:"id"`
	Permission string    `json:"permission"`
	Name       string    `json:"name"`
	CreatedBy  int64     `json:"created_by"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ViewDetailParameterDTOOut struct {
	ID          int64     `json:"id"`
	Permission  string    `json:"permission"`
	Name        string    `json:"name"`
	Value       string    `json:"value"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ParameterForView struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Value       string    `json:"value"`
	UpdatedAt  time.Time `json:"updated_at"`
}
