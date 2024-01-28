package out

import "time"

type ViewListDataGroupDTOOut struct {
	ID          int64     `json:"id"`
	GroupID     string    `json:"group_id"`
	Description string    `json:"description"`
	CreatedBy   int64     `json:"created_by"`
	CreatedName string    `json:"created_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   int64     `json:"updated_by"`
}

type ViewDetailDataGroupDTOOut struct {
	ID          int64     `json:"id"`
	GroupID     string    `json:"group_id"`
	Description string    `json:"description"`
	Scope       []Scopes  `json:"scope"`
	CreatedBy   int64     `json:"created_by"`
	UpdatedBy   int64     `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type InitiateInsertUpdateDataGroupDTOOut struct {
	Scopes []Scopes `json:"scopes"`
}

type Scopes struct {
	Key   string  `json:"key"`
	Menu  string  `json:"menu"`
	Scope []Scope `json:"scope"`
}

type Scope struct {
	Label     string `json:"label"`
	Value     string `json:"value"`
	IsChecked bool   `json:"is_checked"`
}

type ViewDetailDataGroupResponse struct {
	ID          int64          `json:"id"`
	GroupID     string         `json:"group_id"`
	Description string         `json:"description"`
	Scope       []DetailScopes `json:"scope"`
	CreatedBy   int64          `json:"created_by"`
	UpdatedBy   int64          `json:"updated_by"`
	UpdatedAt   time.Time      `json:"updated_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedName string         `json:"updated_name"`
}

type DetailScopes struct {
	Key   string        `json:"key"`
	Menu  string        `json:"menu"`
	Scope []DetailScope `json:"scope"`
}

type DetailScope struct {
	Label string     `json:"label"`
	Value ScopeValue `json:"value"`
}

type ScopeValue struct {
	ID       int64  `json:"id"`
	ParentID int64  `json:"parent_id"`
	Name     string `json:"name"`
}
