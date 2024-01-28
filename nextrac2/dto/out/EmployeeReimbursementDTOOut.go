package out

import "time"

type EmployeeReimbursement struct {
	ID                  int64     `json:"id"`
	EmployeeId          int64     `json:"employee_id"`
	IDCard              string    `json:"id_card"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	FullName            string    `json:"full_name"`
	Department          string    `json:"department"`
	CurrentMedicalValue float64   `json:"current_medical_value"`
	ReceiptNo           string    `json:"receipt_no"`
	Value               float64   `json:"value"`
	Status              string    `json:"status"`
	VerifiedStatus      string    `json:"verified_status"`
	ApprovedValue       float64   `json:"approved_value"`
	Note                string    `json:"note"`
	Filename  			string 	  `json:"filename"`
	Attachment          string    `json:"attachment"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

type EmployeeReimbursementForReport struct {
	ReceiptNo           string    `json:"f1"`
	ApprovedValue       float64   `json:"f2"`
	Description         string    `json:"f3"`
}

type EmployeeReimbursementReportResponse struct {
	ID                  int64     `json:"id"`
	FirstName           string    `json:"first_name"`
	LastName            string    `json:"last_name"`
	FullName            string    `json:"full_name"`
	CurrentMedicalValue float64   `json:"current_medical_value"`
	TotalValue          float64   `json:"total_value"`
	LastMedicalValue    float64   `json:"last_medical_value"`
	January             float64   `json:"january"`
	February            float64   `json:"february"`
	March               float64   `json:"march"`
	April               float64   `json:"april"`
	May                 float64   `json:"may"`
	June                float64   `json:"june"`
	July                float64   `json:"july"`
	August              float64   `json:"august"`
	September           float64   `json:"september"`
	October             float64   `json:"october"`
	November            float64   `json:"november"`
	December            float64   `json:"december"`
}