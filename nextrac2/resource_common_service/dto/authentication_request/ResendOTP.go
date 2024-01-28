package authentication_request

type ResendUserVerificationRequest struct {
	UserID           int64  `json:"user_id"`
	Email            string `json:"email"`
	EmailLinkMessage string `json:"email_link_message"`
	EmailMessage     string `json:"email_message"`
}

type ResendUserVerificationMessageParam struct {
	Purpose      string `json:"purpose"`
	FirstName    string `json:"first_name"`
	ClientTypeID int64  `json:"client_type_id"`
	UniqueID1    string `json:"unique_id_1"`
	CompanyName  string `json:"company_name"`
	UniqueID2    string `json:"unique_id_2"`
	BranchName   string `json:"branch_name"`
	SalesmanID   string `json:"salesman_id"`
	UserID       string `json:"user_id"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	ClientID     string `json:"client_id"`
	AuthUserID   int64  `json:"auth_user_id"`
}
