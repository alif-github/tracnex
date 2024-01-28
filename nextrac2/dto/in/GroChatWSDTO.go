package in

type LoginWS struct {
	Username string `json:"userName"`
	Password string `json:"password"`
}

type GroChatWSHandshake struct {
	Token    string `json:"token"`
	ClientId string `json:"client_id"`
	Sign     string `json:"sign"`
}
