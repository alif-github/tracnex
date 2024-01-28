package out

type GetListLogDTOOut struct {
	ID                    int64  `json:"id"`
	ClientID              string `json:"client_id"`
	ClientTypeID          int64  `json:"client_type_id"`
	SuccessStatusAuth     bool   `json:"success_status_auth"`
	SuccessStatusNexcloud bool   `json:"success_status_nexcloud"`
	Resource              string `json:"resource"`
}

type DetailErrorRegistrationClientID struct {
	CompanyId    string `json:"company_id"`
	BranchID     string `json:"branch_id"`
	ErrorMessage string `json:"error_message"`
}
