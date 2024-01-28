package authentication_request

type VerifyRequestDTO struct {
	Email		string `json:"email"`
	Username 	string `json:"username"`
	Password 	string `json:"password"`
}

