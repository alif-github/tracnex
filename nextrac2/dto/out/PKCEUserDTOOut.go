package out

import "time"

type PKCEResponse struct {
	UserID			int64			`json:"user_id"`
	ClientID		string			`json:"client_id"`
	Username		string			`json:"username"`
	ResourceList	[]ResourceList	`json:"resource_list"`
}

type ResourceList struct {
	ResourceID		string			`json:"resource_id"`
	Status			string			`json:"status"`
}

type ViewPKCEResponse struct {
	UserID			int64		`json:"user_id"`
	ParentClientID	string		`json:"parent_client_id"`
	ClientID		string		`json:"client_id"`
	Username		string		`json:"username"`
	CreatedBy		int64		`json:"created_by"`
	UpdatedAt		time.Time	`json:"updated_at"`
}

type ChangePasswordPKCEResponse struct {
	UserID			int64		`json:"user_id"`
	ClientID		string		`json:"client_id"`
	Username		string		`json:"username"`
}
