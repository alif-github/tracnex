package authentication_response

import "nexsoft.co.id/nextrac2/resource_common_service/model"

type CheckClientOrUserResponse struct {
	Nexsoft		CheckClientOrUserBodyResponse	`json:"nexsoft"`
}

type CheckClientOrUserBodyResponse struct {
	Header 		model.HeaderResponse		`json:"header"`
	Payload		CheckClientOrUserPayload	`json:"payload"`
}

type CheckClientOrUserPayload struct {
	model.PayloadResponse
	Data	CheckClientOrUserData	`json:"data"`
}

type CheckClientOrUserData struct {
	Content		CheckClientOrUserContent	`json:"content"`
}

type CheckClientOrUserContent struct {
	IsExist					bool							`json:"is_exist"`
	AdditionalInformation	AdditionalInformationContent	`json:"additional_information"`
}

type AdditionalInformationContent struct {
	AliasName			string		`json:"alias_name"`
	ClientID			string		`json:"client_id"`
	ClientSecret		string		`json:"client_secret"`
	UserID				int64		`json:"user_id"`
	Username			string		`json:"username"`
	SignatureKey		string		`json:"signature_key"`
	GrantTypes			string		`json:"grant_types"`
	ResourceID			string		`json:"resource_id"`
	IPWhitelist			string		`json:"ip_whitelist"`
	UserStatus			int64		`json:"user_status"`
	Scope				string		`json:"scope"`
	Locale				string		`json:"locale"`
	ClientInformation	string		`json:"client_information"`
	UserInformation		string		`json:"user_information"`
}
