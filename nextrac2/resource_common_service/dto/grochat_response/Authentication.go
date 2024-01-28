package grochat_response

type GroChatData struct {
	NexSoft NexSoft `json:"nexsoft"`
}

type NexSoft struct {
	Payload Payload `json:"payload"`
}

type Payload struct {
	Data Data `json:"data"`
}

type Data struct {

}

type GroChatAuthenticationResponse struct {
	Data GroChatAuthenticationData `json:"data"`
}

type GroChatAuthenticationData struct {
	UserToken    string `json:"userToken"`
	RefreshToken string `json:"refreshToken"`
}

type GroChatAuthenticationNoteResponse struct {

}

type GroChatAuthenticationErrorResponse struct {
	Code        int    `json:"code"`
	Status      int    `json:"status"`
	Description string `json:"description"`
}