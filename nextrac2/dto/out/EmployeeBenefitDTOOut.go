package out

type EmployeeAnnualLeaveBenefit struct {
	TotalLeave         int64 `json:"total_leave"`
	CurrentAnnualLeave int64 `json:"current_annual_leave"`
	LastAnnualLeave    int64 `json:"last_annual_leave"`
	NegativeLeave	   int64 `json:"negative_leave"`
}

type EmployeeReimbursementBenefit struct {
	CurrentMedicalValue float64 `json:"current_medical_value"`
}
