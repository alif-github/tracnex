package out

import "time"

type ClientMappingResponse struct {
	ID             int64     `json:"id"`
	ClientID       string    `json:"client_id"`
	ClientTypeId   int64     `json:"client_type_id"`
	AuthUserId     int64     `json:"auth_user_id"`
	Username       string    `json:"username"`
	CompanyId      string    `json:"company_id"`
	BranchId       string    `json:"branch_id"`
	SocketID       string    `json:"socket_id"`
	SocketPassword string    `json:"-"`
	Aliases        string    `json:"aliases"`
	CreatedAt      time.Time `json:"created_at"`
	CreatedBy      int64     `json:"created_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      int64     `json:"updated_by"`
}

type ClientMappingForView struct {
	ID                    int64     `json:"id"`
	ClientID              string    `json:"client_id"`
	SocketID              string    `json:"socket_id"`
	ClientType            string    `json:"client_type"`
	CompanyID             string    `json:"company_id"`
	BranchID              string    `json:"branch_id"`
	Aliases               string    `json:"aliases"`
	SuccessStatusAuth     bool      `json:"success_status_auth"`
	SuccessStatusNexcloud bool      `json:"success_status_nexcloud"`
	SuccessStatusNexdrive bool      `json:"success_status_nexdrive"`
	CreatedAt             time.Time `json:"created_at"`
	CreatedBy             int64     `json:"created_by"`
	UpdatedAt             time.Time `json:"updated_at"`
	UpdatedBy             int64     `json:"updated_by"`
}

type CompanyBranchErrorBulk struct {
	CompanyID    string `json:"company_id"`
	BranchID     string `json:"branch_id"`
	ErrorMessage string `json:"error_message"`
}
