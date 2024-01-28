package out

type LoginResponse struct {
	Code     int64  `json:"code"`
	Status   int64  `json:"status"`
	Data     DataLoginResponse `json:"data"`
}

type DataLoginResponse struct {
	NexSoft      APIResponse       `json:"nexsoft"`
	UserModel    UserModelResponse `json:"userModel"`
	UserToken    string            `json:"user_token"`
	RefreshToken string            `json:"refresh_token"`
}

type UserModelResponse struct {
	Resend interface{}  `json:"resend"`
	Auth  AuthResponse  `json:"auth"`
}

type AuthResponse struct {
	SignatureKey string  `json:"signature_key"`
	ClientId     string  `json:"client_id"`
	UserId       int64  `json:"user_id"`
}

type GroWSResponseSuccess struct {
	Code    int64   `json:"code"`
	Status  int64   `json:"status"`
	Data    GroWSResponseSuccessData `json:"data"`
}

type GroWSResponseSuccessData struct {
	Success        bool                 `json:"success"`
	HeaderGro      HeaderGroResponse    `json:"header"`
	PayloadGro     PayloadGroResponse `json:"payload"`
}

type HeaderGroResponse struct {
	RequestId  string   `json:"request_id"`
	Version    string   `json:"version"`
	Timestamp  string   `json:"timestamp"`
}

type PayloadGroResponse struct {
	UserId      int64     `json:"user_id"`
	ClientId    string    `json:"client_id"`
	TypeDeviceId int64    `json:"type_device_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SignUuid    string    `json:"sign_uuid"`
}