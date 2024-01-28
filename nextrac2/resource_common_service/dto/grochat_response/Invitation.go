package grochat_response

type Invitation struct {
	Code        int            `json:"code"`
	Description string         `json:"description"`
	Data        InvitationData `json:"data"`
	Note        string         `json:"note"`
}

type InvitationData struct {
	InvitationCode string `json:"invitation_code"`
	ExpiredDate    string `json:"expired_date"`
	ResourceId     string `json:"resource_id"`
	ClientId       string `json:"client_id"`
}
