package redmine_request

type IssuePaidRedmineRequest struct {
	Issue CustomFields `json:"issue"`
}

type CustomFields struct {
	CustomField []Fields `json:"custom_fields"`
}

type Fields struct {
	ID    int64  `json:"id"`
	Value string `json:"value"`
}
