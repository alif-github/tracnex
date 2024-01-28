package out

import "time"

type ViewDetailRoleDTOOut struct {
	ID          int64         `json:"id"`
	RoleID      string        `json:"role_id"`
	Description string        `json:"description"`
	Permissions []Permissions `json:"permission"`
	CreatedBy   int64         `json:"created_by"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedBy   int64         `json:"updated_by"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type ViewListRoleDTOOut struct {
	ID          int64     `json:"id"`
	RoleID      string    `json:"role_id"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedName string    `json:"created_name"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type InitiateInsertUpdateRoleDTOOut struct {
	Permissions []Permissions `json:"permissions"`
}

type Permissions struct {
	Key        string       `json:"key"`
	Menu       string       `json:"menu"`
	Permission []Permission `json:"permission"`
}

type Permission struct {
	Label     string `json:"label"`
	Value     string `json:"value"`
	IsChecked bool   `json:"is_checked"`
}
