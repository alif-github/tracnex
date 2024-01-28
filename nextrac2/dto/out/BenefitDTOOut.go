package out

import "time"

type EmployeeReimbursementType struct {
	Id          int64     `json:"id"`
	BenefitName string    `json:"benefit_name"`
	BenefitType string    `json:"benefit_type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}