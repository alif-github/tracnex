package grochat_request

type Invitation struct {
	Email         string `json:"email"`
	EmailMessage  string `json:"email_message"`
	URLInvitation string `json:"url_invitation"`
	ResourceId    string `json:"resource_id"`
	UserType      string `json:"userType"`
}
