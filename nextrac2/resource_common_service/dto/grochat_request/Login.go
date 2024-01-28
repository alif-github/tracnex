package grochat_request

type GroChatLoginRequest struct {
	RequestId         string `json:"requestID"`
	AuthorizationCode string `json:"authorization_code"`
	CodeVerifier      string `json:"code_verifier"`
	ResourceId        string `json:"resourceID"`
}