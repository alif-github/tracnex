package grochat_ws_handler

import "encoding/json"

type LoginWS struct {
	Username string `json:"userName"`
	Password string `json:"password"`
}

type GroChatWSHandshake struct {
	Token    string `json:"token"`
	ClientId string `json:"client_id"`
	Sign     string `json:"sign"`
}

type GroChatWSKeepAlive struct {
	Type string
}

type LoginResponse struct {
	Code     int64  `json:"code"`
	Status   int64  `json:"status"`
	Data     DataLoginResponse `json:"data"`
}

type DataLoginResponse struct {
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

/*
	Message
*/
type MessageType struct {
	Type string `json:"type"`
}

type ChatMessage struct {
	Type    string  `json:"type"`
	Message Message `json:"message"`
}

func (c *ChatMessage) Bytes() []byte {
	bytes, _ := json.Marshal(c)
	return bytes
}

type Message struct {
	SourceId      string        `json:"source_id"`
	DestinationId string        `json:"destination_id"`
	MessageDetail MessageDetail `json:"message_detail"`
	Identity      Identity      `json:"identity"`
}

type Identity struct {
	ClientId string `json:"client_id"`
	Sign     string `json:"sign"`
}

type MessageDetail struct {
	Guarantee    bool         `json:"guarantee"`
	Type         string       `json:"type"`
	MessageId    string       `json:"message_id"`
	MessageModel MessageModel `json:"message_model"`
}

type MessageModel struct {
	AdditionalKey string      `json:"additional_key"`
	Content       ContentText `json:"content"`
	CreatedAt     int64       `json:"created_at"`
	DeletedType   int         `json:"deleted_type"`
	IsBroadcast   int         `json:"is_broadcast"`
	IsEncrypted   string      `json:"is_encrypted"`
	IsForward     int         `json:"is_forward"`
	LocalId       string      `json:"local_id"`
	RoomId        string      `json:"room_id"`
	TypeMessageId int         `json:"type_message_id"`
	VectorX       int         `json:"vector_x"`
	VectorY       int         `json:"vector_y"`
}

type ContentText struct {
	Text string `json:"text"`
}
