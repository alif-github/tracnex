package out

import (
	"nexsoft.co.id/nextrac2/repository"
	"time"
)

type ListJobProcessResponse struct {
	Level     int       `json:"level"`
	JobID     string    `json:"job_id"`
	Group     string    `json:"group"`
	Type      string    `json:"type"`
	Name      string    `json:"name"`
	Counter   int       `json:"counter"`
	Total     int       `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Duration  float64   `json:"duration"`
}

type ViewJobProcessResponse struct {
	Level           int                        `json:"level"`
	JobID           string                     `json:"job_id"`
	Group           string                     `json:"group"`
	Type            string                     `json:"type"`
	Name            string                     `json:"name"`
	UrlIn           string                     `json:"url_in"`
	FileNameIn      string                     `json:"file_name_in"`
	ContentDataOut  string                     `json:"content_data_out"`
	Counter         int                        `json:"counter"`
	Total           int                        `json:"total"`
	Status          string                     `json:"status"`
	CreatedAt       time.Time                  `json:"created_at"`
	UpdatedAt       time.Time                  `json:"updated_at"`
	Duration        float64                    `json:"duration"`
	Percentage      float64                    `json:"percentage"`
	ChildJobProcess repository.ChildJobProcess `json:"child_job_process"`
}

type ViewConfirmImportJobProcessResponse struct {
	JobID string `json:"job_id"`
}

type GetListFileUploadJobProcess struct {
	JobID       string    `json:"job_id"`
	Status      string    `json:"status"`
	Progress    float64   `json:"progress"`
	FileUrl     string    `json:"file_url"`
	CreatedName string    `json:"created_name"`
	CreatedAt   time.Time `json:"created_at"`
}
