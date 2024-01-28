package out

import (
	"time"
)

type EmployeeApprovalHistory struct {
	ID                  int64     `json:"id"`
	Firstname           string    `json:"firstname"`
	Lastname            string    `json:"lastname"`
	IDCard              string    `json:"id_card"`
	Department          string    `json:"department"`
	RequestType         string    `json:"request_type"`
	Type                string    `json:"type"`
	ReceiptNo           string    `json:"receipt_no"`
	TotalLeave          int64     `json:"total_leave"`
	Value				float64   `json:"value"`
	Date                []string  `json:"date"`
	Status              string    `json:"status"`
	VerifiedStatus      string    `json:"verified_status"`
	ApprovedValue       float64   `json:"approved_value"`
	CancellationReason  string    `json:"cancellation_reason"`
	Note                string    `json:"note"`
	Attachment          string    `json:"attachment"`
	TotalRemainingLeave int64     `json:"total_remaining_leave"`
	Description         string    `json:"description"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
