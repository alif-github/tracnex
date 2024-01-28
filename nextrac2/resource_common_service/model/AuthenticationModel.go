package model

import "nexsoft.co.id/nexcommon/util"

type AuthenticationModel struct {
	Oauth       OauthModel `json:"oauth"`
	Application []struct {
		Name    string `json:"name"`
		License string `json:"license"`
	} `json:"application"`
	Role AuthenticationRoleModel `json:"role"`
	Data AuthenticationDataModel `json:"data"`
}

type AuthenticationRoleModel struct {
	Role       string              `json:"role"`
	Permission map[string][]string `json:"permission"`
}

type AuthenticationDataModel struct {
	Group string              		`json:"group"`
	Scope map[string]interface{} 	`json:"scope"`
}

type OauthModel struct {
	ClientID   string                  `json:"client_id"`
	IsAdmin    bool                    `json:"is_admin"`
	UserID     string                  `json:"user_id"`
	ResourceID string                  `json:"resource_id"`
	Scope      string                  `json:"scope"`
	ClientInfo []AdditionalInformation `json:"client_info"`
	UserInfo   []AdditionalInformation `json:"user_info"`
}

type AdditionalInformation struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (input AuthenticationModel) String() string {
	return util.StructToJSON(input)
}
