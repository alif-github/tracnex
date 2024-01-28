package out

import "time"

type GetListEmployeeTimeSheetResponse struct {
	ID                    int64     `json:"id"`
	Nik                   int64     `json:"nik"`
	RedmineId             int64     `json:"redmine_id"`
	Name                  string    `json:"name"`
	DepartmentName        string    `json:"department_name"`
	MandaysRate           float64   `json:"mandays_rate"`
	MandaysRateAutomation float64   `json:"mandays_rate_automation"`
	MandaysRateManual     float64   `json:"mandays_rate_manual"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	UpdatedBy             int64     `json:"updated_by"`
	UpdatedName           string    `json:"updated_name"`
}

type GetListEmployeeResponse struct {
	ID           int64     `json:"id"`
	NIK          string    `json:"nik"`
	Name         string    `json:"name"`
	DepartmentID int64     `json:"department_id"`
	Department   string    `json:"department"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	UpdatedBy    int64     `json:"updated_by"`
	UpdatedName  string    `json:"updated_name"`
}

type GetListEmployeeForDDLResponse struct {
	ID             int64     `json:"id"`
	Nik            string    `json:"nik"`
	RedmineId      int64     `json:"redmine_id"`
	Name           string    `json:"name"`
	DepartmentName string    `json:"department_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	UpdatedBy      int64     `json:"updated_by"`
	UpdatedName    string    `json:"updated_name"`
}

type ViewEmployeeTimeSheetResponse struct {
	ID                    int64     `json:"id"`
	IDCard                string    `json:"id_card"`
	RedmineId             int64     `json:"redmine_id"`
	Name                  string    `json:"name"`
	DepartmentID          int64     `json:"department_id"`
	DepartmentName        string    `json:"department_name"`
	MandaysRate           float64   `json:"mandays_rate"`
	MandaysRateAutomation float64   `json:"mandays_rate_automation"`
	MandaysRateManual     float64   `json:"mandays_rate_manual"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	CreatedName           string    `json:"created_name"`
	UpdatedName           string    `json:"updated_name"`
}

type ViewEmployeeResponse struct {
	ID                 int64        `json:"id"`
	IDCard             string       `json:"id_card"`
	NPWP               string       `json:"npwp"`
	FirstName          string       `json:"first_name"`
	LastName           string       `json:"last_name"`
	Email              string       `json:"email"`
	Phone              string       `json:"phone"`
	Gender             string       `json:"gender"`
	PlaceOfBirth       string       `json:"place_of_birth"`
	DateOfBirth        time.Time    `json:"date_of_birth"`
	AddressResidence   string       `json:"address_residence"`
	AddressTax         string       `json:"address_tax"`
	DateJoin           time.Time    `json:"date_join"`
	DateOut            time.Time    `json:"date_out"`
	Religion           string       `json:"religion"`
	Type               string       `json:"type"`
	Status             string       `json:"status"`
	PositionID         int64        `json:"position_id"`
	Position           string       `json:"position"`
	MaritalStatus      string       `json:"marital_status"`
	Education          string       `json:"education"`
	MothersMaiden      string       `json:"mothers_maiden"`
	NumberOfDependents int64        `json:"number_of_dependents"`
	Nationality        string       `json:"nationality"`
	TaxMethod          string       `json:"tax_method"`
	ReasonResignation  string       `json:"reason_resignation"`
	Photo              string       `json:"photo"`
	DepartmentID       int64        `json:"department_id"`
	DepartmentName     string       `json:"department_name"`
	CreatedAt          time.Time    `json:"created_at"`
	UpdatedAt          time.Time    `json:"updated_at"`
	CreatedName        string       `json:"created_name"`
	UpdatedName        string       `json:"updated_name"`
	BPJSNo             string       `json:"bpjs_no"`
	BPJSTkNo           string       `json:"bpjs_tk_no"`
	IsHaveMember       bool         `json:"is_have_member"`
	MemberList         []MemberList `json:"member"`
	LevelID            int64        `json:"level_id"`
	Level              string       `json:"level"`
	GradeID            int64        `json:"grade_id"`
	Grade              string       `json:"grade"`
	Active             bool         `json:"active"`
}

type MemberList struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ViewEmployeeByNIKResponse struct {
	IDCard    string `json:"id_card"`
	RedmineId int64  `json:"redmine_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type EmployeeLevelResponse struct {
	ID          int64     `json:"id"`
	Level       string    `json:"level"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EmployeeGradeResponse struct {
	ID          int64     `json:"id"`
	Grade       string    `json:"grade"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EmployeeAllowanceResponse struct {
	ID            int64     `json:"id"`
	AllowanceName string    `json:"allowance_name"`
	AllowanceType string    `json:"allowance_type"`
	Value         string    `json:"value"`
	Active        bool      `json:"active"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type EmployeeBenefitMasterResponse struct {
	ID          int64     `json:"id"`
	BenefitName string    `json:"benefit_name"`
	BenefitType string    `json:"benefit_type"`
	Value       string    `json:"value"`
	Active      bool      `json:"active"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EmployeeMatrixResponse struct {
	LevelID       int64                           `json:"level_id"`
	Level         string                          `json:"level"`
	GradeID       int64                           `json:"grade_id"`
	Grade         string                          `json:"grade"`
	AllowanceList []EmployeeAllowanceResponse     `json:"allowance_list"`
	BenefitList   []EmployeeBenefitMasterResponse `json:"benefit_list"`
}

type EmployeeMatrixForView struct {
    Id      int64  `json:"id"`
	LevelID int64  `json:"level_id"`
	Level   string `json:"level"`
	GradeID int64  `json:"grade_id"`
	Grade   string `json:"grade"`
}
