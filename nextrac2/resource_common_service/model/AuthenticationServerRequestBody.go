package model

type CheckURLTokenBody struct {
	ResourceID string `json:"resource_id"`
	Scope      string `json:"scope"`
}

type AddClientResourceBody struct {
	ResourceID string `json:"resource_id"`
	ClientID   string `json:"client_id"`
}

type RefreshTokenBody struct {
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}
