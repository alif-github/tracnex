package in

type EmployeeMatrixRequest struct {
	AbstractDTO
	ID                    int64     `json:"id"`
	LevelID               int64     `json:"level_id"`
	GradeID               int64     `json:"grade_id"`
	AllowanceList         []EmployeeAllowanceRequest `json:"allowance_list"`
	BenefitList           []EmpBenefitRequest `json:"benefit_list"`
}
