package out

import "time"

type APIProjectRedmineResponse struct {
	Projects   []Project `json:"projects"`
	TotalCount int64     `json:"total_count"`
	Offset     int64     `json:"offset"`
	limit      int64     `json:"limit"`
}

type Project struct {
	Id             int64     `json:"id"`
	Name           string    `json:"name"`
	Identifier     string    `json:"identifier"`
	Description    string    `json:"description"`
	Status         int64     `json:"status"`
	IsPublic       bool      `json:"is_public"`
	InheritMembers bool      `json:"inherit_members"`
	CreatedOn      time.Time `json:"created_on"`
	UpdatedOn      time.Time `json:"updated_on"`
}

type DropdownProject struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}
