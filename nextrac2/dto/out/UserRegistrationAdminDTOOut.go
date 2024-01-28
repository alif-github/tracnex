package out

type ListUserRegistrationAdminResponse struct {
	ID                 int64  `json:"id"`
	CustomerName       string `json:"customer_name"`
	ParentCustomerName string `json:"parent_customer_name"`
	CompanyID          string `json:"company_id"`
	BranchID           string `json:"branch_id"`
	CompanyName        string `json:"company_name"`
	BranchName         string `json:"branch_name"`
	UserAdmin          string `json:"user_admin"`
	PasswordAdmin      string `json:"password_admin"`
}

type DetailUserRegistrationAdminResponse struct {
	ID                 int64  `json:"id"`
	ParentCustomerId   int64  `json:"parent_customer_id"`
	ParentCustomerName string `json:"parent_customer_name"`
	CustomerId         int64  `json:"customer_id"`
	SiteId             int64  `json:"site_id"`
	CustomerName       string `json:"customer_name"`
	CompanyId          string `json:"company_id"`
	BranchId           string `json:"branch_id"`
	CompanyName        string `json:"company_name"`
	BranchName         string `json:"branch_name"`
	UserAdmin          string `json:"user_admin"`
	PasswordAdmin      string `json:"password_admin"`
}
