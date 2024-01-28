package out

import (
	"time"
)

type EmployeeRequestHistoryResponse struct {
	ID                 int64     `json:"id"`
	ReceiptNo		   string 	 `json:"receipt_no"`
	Description        string    `json:"description"`
	Date               []string  `json:"date"`
	RequestType        string    `json:"request_type"`
	Type               string    `json:"type"`
	TotalLeave         int64     `json:"total_leave"`
	Value              float64   `json:"value"`
	ApprovedValue      float64   `json:"approved_value"`
	Status             string    `json:"status"`
	VerifiedStatus     string    `json:"verified_status"`
	CancellationReason string    `json:"cancellation_reason"`
	Note               string    `json:"note"`
	Attachment         string    `json:"attachment"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}
