package out

type ListAbsent struct {
	IDCard          string  `json:"id_card"`
	Name            string  `json:"name"`
	NormalDays      int64   `json:"normal_days"`
	ActualDays      int64   `json:"actual_days"`
	Absent          int64   `json:"absent"`
	Overdue         int64   `json:"overdue"`
	LeaveEarly      int64   `json:"leave_early"`
	Overtime        int64   `json:"overtime"`
	NumbersOfLeave  int64   `json:"numbers_of_leave"`
	LeavingDuties   int64   `json:"leaving_duties"`
	NumbersIn       int64   `json:"numbers_in"`
	NumbersOut      int64   `json:"numbers_out"`
	Scan            int64   `json:"scan"`
	SickLeave       int64   `json:"sick_leave"`
	PaidLeave       int64   `json:"paid_leave"`
	PermissionLeave int64   `json:"permission_leave"`
	WorkHours       int64   `json:"work_hours"`
	PercentAbsent   float64 `json:"percent_absent"`
}
