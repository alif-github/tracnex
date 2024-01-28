package grochat_response

type GroChatUserDetailResponse struct {
	Data *GroChatUserDetailData `json:"data"`
}

type GroChatUserDetailData struct {
	GroChat    GroChatUserDetail `json:"grochatData"`
	AuthServer AuthUserDetail    `json:"authServer"`
}

type GroChatUserDetail struct {
	Auth GroChatUserDetailAuth `json:"auth"`
}

type GroChatUserDetailAuth struct {
	UserId   int64  `json:"user_id"`
	ClientId string `json:"client_id"`
}

type AuthUserDetail struct {
	NexSoft AuthUserDetailNexSoft `json:"nexsoft"`
}

type AuthUserDetailNexSoft struct {
	Payload AuthUserDetailPayload `json:"payload"`
}

type AuthUserDetailPayload struct {
	Data AuthUserDetailData `json:"data"`
}

type AuthUserDetailData struct {
	Content AuthUserDetailContent `json:"content"`
}

type AuthUserDetailContent struct {
	AliasName string `json:"alias_name"`
	UserId    int64  `json:"user_id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	ClientId  string `json:"client_id"`
	Locale    string `json:"locale"`
}