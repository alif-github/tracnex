package out

import "time"

type EmployeeNotification struct {
	ID                          int64       `json:"id"`
	IsRequestingForApproval     bool        `json:"is_requesting_for_approval"`
	IsRequestingForCancellation bool        `json:"is_requesting_for_cancellation"`
	IsCancellation              bool        `json:"is_cancellation"`
	IsVerified					bool 		`json:"is_verified"`
	EmployeeId                  int64       `json:"employee_id"`
	Name                        string      `json:"name"`
	RequestType                 string      `json:"request_type"`
	Status                      string      `json:"status"`
	Date                        []time.Time `json:"date"`
	MessageTitle                string      `json:"message_title"`
	MessageBody                 string      `json:"message_body"`
	IsRead						bool 		`json:"is_read"`
	CreatedAt                   time.Time   `json:"created_at"`
}
