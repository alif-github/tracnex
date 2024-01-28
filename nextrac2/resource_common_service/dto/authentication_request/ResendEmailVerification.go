package authentication_request

type ResendEmailVerificationRequest struct {
	UserId int64  `json:"user_id"`
	Email  string `json:"email"`
}
