package authentication_request

type CheckUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	ClientID string `json:"client_id"`
}
