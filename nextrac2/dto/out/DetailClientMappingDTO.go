package out

type DetailClientMappingContent struct {
	ClientID        string                  `json:"client_id"`
	ClientTypeID    int64                   `json:"client_type_id"`
	AuthUserId      string                  `json:"auth_user_id"`
	UserName        string                  `json:"username"`
	CompanyId       string                  `json:"company_id"`
	BranchId        string                  `json:"branch_id"`
	Aliases         string                  `json:"aliases"`
}