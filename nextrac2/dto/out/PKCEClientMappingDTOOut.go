package out

import "time"

type PKCEClientMappingForList struct {
	ID             int64     `json:"id"`
	ParentClientID string    `json:"parent_client_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Username       string    `json:"username"`
	ClientType     string    `json:"client_type"`
	CompanyID      string    `json:"company_id"`
	BranchID       string    `json:"branch_id"`
	ClientAlias    string    `json:"client_alias"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      int64     `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      int64     `json:"updated_by"`
}

type PKCEClientMappingForDetail struct {
	ID          int64     `json:"id"`
	ClientID    string    `json:"client_id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Username    string    `json:"username"`
	ClientType  string    `json:"client_type"`
	CompanyID   string    `json:"company_id"`
	BranchID    string    `json:"branch_id"`
	ClientAlias string    `json:"client_alias"`
	CreatedAt   time.Time `json:"created_at"`
	CreatedBy   int64     `json:"created_by"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   int64     `json:"updated_by"`
}
