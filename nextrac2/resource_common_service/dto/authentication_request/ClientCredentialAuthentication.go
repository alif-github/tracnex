package authentication_request

type ClientCredentialAuthentication struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	WantedScope  string `json:"wanted_scope"`
}
