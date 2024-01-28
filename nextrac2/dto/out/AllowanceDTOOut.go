package out

import "time"

type EmployeeLeaveAllowance struct {
	Id            int64     `json:"id"`
	AllowanceName string    `json:"allowance_name"`
	AllowanceType string    `json:"allowance_type"`
	Value         string    `json:"value"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
