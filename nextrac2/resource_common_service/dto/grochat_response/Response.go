package grochat_response

type GroChatErrorResponse struct {
	Code        int    `json:"code"`
	Status      int    `json:"status"`
	Description string `json:"description"`
	Note		string `json:"note"`
}
