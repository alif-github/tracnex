package authentication_request

type ForgetPasswordDTOin struct {
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	PhoneMessage string `json:"phone_message"`
	EmailMessage string `json:"email_message"`
	EmailLink    string `json:"email_link"`
	ForgetCode   string `json:"forget_code"`
}

type ChangePasswordDTOin struct {
	UserID				int64	`json:"user_id"`
	OldPassword			string	`json:"old_password"`
	NewPassword			string	`json:"new_password"`
	VerifyNewPassword	string	`json:"verify_new_password"`
}