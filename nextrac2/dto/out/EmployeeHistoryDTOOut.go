package out

import "time"

type BundleEmployeeHistoryResponse struct {
	Locale      string                        `json:"locale"`
	ListHistory []EmployeeHistoryListResponse `json:"list_history"`
}

type EmployeeHistoryListResponse struct {
	ID                int64               `json:"id"`
	Editor            string              `json:"editor"`
	CreatedAt         time.Time           `json:"created_at"`
	DescriptionDetail []DescriptionDetail `json:"detail"`
}

type DescriptionDetail struct {
	KeyID  string `json:"key_id"`
	KeyEn  string `json:"key_en"`
	Before string `json:"before"`
	After  string `json:"after"`
}
