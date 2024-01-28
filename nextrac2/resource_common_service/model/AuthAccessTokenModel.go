package model

type AuthAccessTokenModel struct {
	RedisAuthAccessTokenModel
	IsAdmin                    bool   `json:"is_admin"`
	ClientID                   string `json:"cid"`
	AuthenticationServerUserID int64  `json:"aid"`
	Locale                     string `json:"lang"`
	Scope                      string `json:"scp"`
}

func (input AuthAccessTokenModel) ConvertToRedisModel() RedisAuthAccessTokenModel {
	return RedisAuthAccessTokenModel{
		ResourceUserID: input.ResourceUserID,
		Authentication: input.Authentication,
		IPWhiteList:    input.IPWhiteList,
		SignatureKey:   input.SignatureKey,
		Locale:         input.Locale,
	}
}

type RedisAuthAccessTokenModel struct {
	ResourceUserID int64  `json:"rid"`
	Authentication string `json:"auth"`
	IPWhiteList    string `json:"ipl"`
	SignatureKey   string `json:"sign"`
	Locale         string `json:"locale"`
}
