package out

import "time"

type ParentBacklogResponse struct {
	Sprint       string  `json:"sprint"`
	TotalMandays float64 `json:"total_mandays"`
}

type BacklogFromFileResponse struct {
	DepartmentCode string  `json:"department_code"`
	Layer1         string  `json:"layer_1"`
	Layer2         string  `json:"layer_2"`
	Layer3         string  `json:"layer_3"`
	Layer4         string  `json:"layer_4"`
	Layer5         string  `json:"layer_5"`
	Redmine        string  `json:"redmine_number"`
	Sprint         string  `json:"sprint"`
	SprintName     string  `json:"sprint_name"`
	PicId          int64   `json:"pic_id"`
	PicName        string  `json:"pic_name"`
	Status         string  `json:"status"`
	Mandays        float64 `json:"mandays"`
	MandaysDone    float64 `json:"mandays_done"`
	FlowChanged    string  `json:"flow_changed"`
	AdditionalData string  `json:"additional_data"`
	Note           string  `json:"note"`
	Url            string  `json:"url"`
	Page           string  `json:"page"`
}

type BacklogFromFileQaQcResponse struct {
	DepartmentCode  string  `json:"department_code"`
	Feature         int64   `json:"feature"`
	Layer1          string  `json:"layer_1"`
	Layer2          string  `json:"layer_2"`
	Layer3          string  `json:"layer_3"`
	Layer4          string  `json:"layer_4"`
	Layer5          string  `json:"layer_5"`
	Subject         string  `json:"subject"`
	Redmine         string  `json:"redmine_number"`
	ReferenceTicket int64   `json:"reference_ticket"`
	Sprint          string  `json:"sprint"`
	PicId           int64   `json:"pic_id"`
	PicName         string  `json:"pic_name"`
	Status          string  `json:"status"`
	Mandays         float64 `json:"mandays"`
	MandaysDone     float64 `json:"mandays_done"`
	Tracker         string  `json:"tracker"`
	FormChanged     string  `json:"form_changed"`
	FlowChanged     string  `json:"flow_changed"`
	AdditionalData  string  `json:"additional_data"`
	Note            string  `json:"note"`
	Url             string  `json:"url"`
	Page            string  `json:"page"`
}

type DetailBacklogResponse struct {
	ID             int64     `json:"id"`
	Layer1         string    `json:"layer_1"`
	Layer2         string    `json:"layer_2"`
	Layer3         string    `json:"layer_3"`
	Redmine        int64     `json:"redmine"`
	Sprint         string    `json:"sprint"`
	Pic            string    `json:"pic"`
	Status         string    `json:"status"`
	Mandays        float64   `json:"mandays"`
	DepartmentCode string    `json:"department_code"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type ViewDetailBacklogResponse struct {
	ID              int64     `json:"id"`
	Layer1          string    `json:"layer_1"`
	Layer2          string    `json:"layer_2"`
	Layer3          string    `json:"layer_3"`
	Layer4          string    `json:"layer_4"`
	Layer5          string    `json:"layer_5"`
	Subject         string    `json:"subject"`
	Feature         int64     `json:"feature"`
	Tracker         string    `json:"tracker"`
	ReferenceTicket int64     `json:"reference_ticket"`
	Description     string    `json:"description"`
	RedmineNumber   string    `json:"redmine_number"`
	Sprint          string    `json:"sprint"`
	SprintName      string    `json:"sprint_name"`
	PicId           int64     `json:"pic_id"`
	Pic             string    `json:"pic"`
	Status          string    `json:"status"`
	Mandays         float64   `json:"mandays"`
	MandaysDone     float64   `json:"mandays_done"`
	FlowChanged     string    `json:"flow_changed"`
	AdditionalData  string    `json:"additional_data"`
	DepartmentId    int64     `json:"department_id"`
	DepartmentName  string    `json:"department_name"`
	Note            string    `json:"note"`
	Url             string    `json:"url"`
	UrlFile         string    `json:"url_file"`
	Page            string    `json:"page"`
	UpdatedAt       time.Time `json:"updated_at"`
	CreatedName     string    `json:"created_name"`
	UpdatedName     string    `json:"updated_name"`
	CreatedBy       int64     `json:"created_by"`
	CreatedAt       time.Time `json:"created_at"`
}
