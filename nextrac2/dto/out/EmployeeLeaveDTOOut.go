package out

import "time"

type EmployeeLeave struct {
	ID            int64     `json:"id"`
	IDCard        string    `json:"id_card"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	FullName      string    `json:"full_name"`
	Department    string    `json:"department"`
	AllowanceName string    `json:"allowance_name"`
	Date          []string  `json:"date"`
	Type          string    `json:"type"`
	Value         int64     `json:"value"`
	Status		  string 	`json:"status"`
	LeaveTime     time.Time `json:"leave_time"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	CancellationReason string `json:"cancellation_reason"`
	Description   string    `json:"description"`
	CurrentAnnualLeave int64 `json:"current_annual_leave"`
	LastAnnualLeave    int64 `json:"last_annual_leave"`
	Filename  			string 	  `json:"filename"`
	Attachment          string    `json:"attachment"`
}

type EmployeeLeaveYearly struct {
	ID                     int64  `json:"id"`
	IDCard                 string `json:"id_card"`
	FirstName              string `json:"first_name"`
	LastName               string `json:"last_name"`
	FullName               string `json:"full_name"`
	Department             string `json:"department"`
	Level                  string `json:"level"`
	Grade                  string `json:"grade"`
	CurrentLeaveThisPeriod int64  `json:"current_leave_this_period"`
	LastLeaveBeforePeriod  int64  `json:"last_leave_before_period"`
	OwingLeave             int64  `json:"owing_leave"`
}

type ListTodaysLeave struct {
	IDCard     string   `json:"id_card"`
	Name       string   `json:"name"`
	Department string   `json:"department"`
	Date       []string `json:"date"`
	Type       string   `json:"type"`
}
