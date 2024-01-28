package in

type PKCEClientMappingRequest struct {
	AbstractDTO
	ID             int64  `json:"id"`
	ParentClientID string `json:"parent_client_id"`
	ClientTypeID   int64  `json:"client_type_id"`
	CompanyID      string `json:"company_id"`
	BranchID       string `json:"branch_id"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ClientAlias    string `json:"client_alias"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
}
