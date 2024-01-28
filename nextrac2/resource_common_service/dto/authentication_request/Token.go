package authentication_request

type TokenRequestDTO struct {
	CodeVerifier      string `json:"code_verifier"`
	AuthorizationCode string `json:"authorization_code"`
}
