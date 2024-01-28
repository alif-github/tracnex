package in

import "time"

type EmployeeNotification struct {
	IsMobileNotification        bool        `json:"is_mobile_notification"`
	IsRequestingForApproval     bool        `json:"is_requesting_for_approval"`
	IsRequestingForCancellation bool        `json:"is_requesting_for_cancellation"`
	IsCancellation              bool        `json:"is_cancellation"`
	IsVerified					bool 		`json:"is_verified"`
	EmployeeId                  int64       `json:"employee_id"`
	RequestType                 string      `json:"request_type"`
	Status                      string      `json:"status"`
	Date                        []time.Time `json:"date"`
	IsRead                      bool        `json:"is_read"`
}
