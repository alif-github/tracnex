package in

type GenerateInternalTokenRequestDTO struct {
	AuthUserID	int64		`json:"auth_user_id"`
	ClientID	string		`json:"client_id"`
	Issuer		string		`json:"issuer"`
	Locale		string		`json:"locale"`
	Destination	string		`json:"destination"`
}