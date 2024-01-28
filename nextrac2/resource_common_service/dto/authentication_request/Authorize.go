package authentication_request

type AuthorizeRequestDTO struct {
	CodeChallenger string `json:"code_challenger"`
	ResponseType   string `json:"response_type"`
}
