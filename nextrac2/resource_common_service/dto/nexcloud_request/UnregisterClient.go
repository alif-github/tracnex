package nexcloud_request

type UnregisterClient struct {
	ClientID	string	`json:"client_id"`
	UpdatedAt	string	`json:"updated_at"`
}