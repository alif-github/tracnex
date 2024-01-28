package out

type GroChatLoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}