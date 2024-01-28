package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"strconv"
	"time"
)

type Absent struct {
	ID              int64   `json:"id"`
	EmployeeID      int64   `json:"employee_id"`
	IDCard          string  `json:"id_card"`
	AbsentID        int64   `json:"absent_id"`
	NormalDays      int64   `json:"normal_days"`
	ActualDays      int64   `json:"actual_days"`
	Absents         int64   `json:"absent"`
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
	PeriodStartStr  string  `json:"period_start"`
	PeriodEndStr    string  `json:"period_end"`
	CreatedBy       int64   `json:"created_by"`
	CreatedClient   string  `json:"created_client"`
	CreatedAtStr    string  `json:"created_at"`
	UpdatedBy       int64   `json:"updated_by"`
	UpdatedClient   string  `json:"updated_client"`
	UpdatedAtStr    string  `json:"updated_at"`
	PeriodStart     time.Time
	PeriodEnd       time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type AbsentMapping struct {
	IDCard          string `json:"id_card"`
	AbsentID        string `json:"absent_id"`
	NormalDays      string `json:"normal_days"`
	ActualDays      string `json:"actual_days"`
	Absents         string `json:"absent"`
	Overdue         string `json:"overdue"`
	LeaveEarly      string `json:"leave_early"`
	Overtime        string `json:"overtime"`
	NumbersOfLeave  string `json:"numbers_of_leave"`
	LeavingDuties   string `json:"leaving_duties"`
	NumbersIn       string `json:"numbers_in"`
	NumbersOut      string `json:"numbers_out"`
	Scan            string `json:"scan"`
	SickLeave       string `json:"sick_leave"`
	PaidLeave       string `json:"paid_leave"`
	PermissionLeave string `json:"permission_leave"`
	WorkHours       string `json:"work_hours"`
	PercentAbsent   string `json:"percent_absent"`
	PeriodStartStr  string `json:"period_start"`
	PeriodEndStr    string `json:"period_end"`
}

func (input Absent) ValidateInsertAbsent(mapping []AbsentMapping) (inputStruct []Absent, err errorModel.ErrorModel) {
	var (
		fileName = "AbsentDTO.go"
		funcName = "ValidateInsertAbsent"
	)

	for _, item := range mapping {
		var inputStructItem Absent

		//--- ID Card Check Empty
		if util.IsStringEmpty(item.IDCard) {
			continue
		}

		//--- ID Card [M]
		inputStructItem.IDCard = item.IDCard

		//--- Absent ID [M]
		inputStructItem.AbsentID, err = input.mandatoryStrInt(item.AbsentID, "Absent ID")
		if err.Error != nil {
			return
		}

		//--- Normal Days [M]
		inputStructItem.NormalDays, err = input.mandatoryStrInt(item.NormalDays, "Normal Days")
		if err.Error != nil {
			return
		}

		//--- Actual Days [O]
		inputStructItem.ActualDays, err = input.optionalStrInt(item.ActualDays)
		if err.Error != nil {
			return
		}

		//--- Absent Days [O]
		inputStructItem.Absents, err = input.optionalStrInt(item.Absents)
		if err.Error != nil {
			return
		}

		//--- Overdue Days [O]
		inputStructItem.Overdue, err = input.optionalStrInt(item.Overdue)
		if err.Error != nil {
			return
		}

		//--- Leave Early Days [O]
		inputStructItem.LeaveEarly, err = input.optionalStrInt(item.LeaveEarly)
		if err.Error != nil {
			return
		}

		//--- Overtime Days [O]
		inputStructItem.Overtime, err = input.optionalStrInt(item.Overtime)
		if err.Error != nil {
			return
		}

		//--- Numbers Of Leave Days [O]
		inputStructItem.NumbersOfLeave, err = input.optionalStrInt(item.NumbersOfLeave)
		if err.Error != nil {
			return
		}

		//--- Leaving Duties Days [O]
		inputStructItem.LeavingDuties, err = input.optionalStrInt(item.LeavingDuties)
		if err.Error != nil {
			return
		}

		//--- Numbers In Days [O]
		inputStructItem.NumbersIn, err = input.optionalStrInt(item.NumbersIn)
		if err.Error != nil {
			return
		}

		//--- Numbers Out Days [O]
		inputStructItem.NumbersOut, err = input.optionalStrInt(item.NumbersOut)
		if err.Error != nil {
			return
		}

		//--- Scan Days [O]
		inputStructItem.Scan, err = input.optionalStrInt(item.Scan)
		if err.Error != nil {
			return
		}

		//--- Sick Leave Days [O]
		inputStructItem.SickLeave, err = input.optionalStrInt(item.SickLeave)
		if err.Error != nil {
			return
		}

		//--- Paid Leave Days [O]
		inputStructItem.PaidLeave, err = input.optionalStrInt(item.PaidLeave)
		if err.Error != nil {
			return
		}

		//--- Permission Leave Days [O]
		inputStructItem.PermissionLeave, err = input.optionalStrInt(item.PermissionLeave)
		if err.Error != nil {
			return
		}

		//--- Work Hours Days [O]
		inputStructItem.WorkHours, err = input.optionalStrInt(item.WorkHours)
		if err.Error != nil {
			return
		}

		//--- Percent Absent ID [M]
		inputStructItem.PercentAbsent, err = input.mandatoryStrFloat(item.PercentAbsent, "Percent Absent")
		if err.Error != nil {
			return
		}

		//--- Period Start
		if util.IsStringEmpty(item.PeriodStartStr) {
			err = errorModel.GenerateEmptyFieldError(fileName, funcName, "Period Start")
			return
		}

		inputStructItem.PeriodStart, err = TimeStrToTimeWithTimeFormat(item.PeriodStartStr, "Period Start", "02-01-2006")
		if err.Error != nil {
			return
		}

		//--- Period End
		if util.IsStringEmpty(item.PeriodEndStr) {
			err = errorModel.GenerateEmptyFieldError(fileName, funcName, "Period End")
			return
		}

		inputStructItem.PeriodEnd, err = TimeStrToTimeWithTimeFormat(item.PeriodEndStr, "Period End", "02-01-2006")
		if err.Error != nil {
			return
		}

		//--- Append
		inputStruct = append(inputStruct, inputStructItem)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input Absent) mandatoryStrInt(value, fieldName string) (num int64, err errorModel.ErrorModel) {
	var (
		fileName = "AbsentDTO.go"
		funcName = "mandatoryStrInt"
	)

	//--- Mandatory Check Empty
	if util.IsStringEmpty(value) {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, fieldName)
		return
	}

	//--- Mandatory String Convert
	n, errs := strconv.Atoi(value)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	//--- Fill
	num = int64(n)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input Absent) mandatoryStrFloat(value, fieldName string) (num float64, err errorModel.ErrorModel) {
	var (
		fileName = "AbsentDTO.go"
		funcName = "mandatoryStrFloat"
		errs     error
	)

	//--- Mandatory Check Empty
	if util.IsStringEmpty(value) {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, fieldName)
		return
	}

	//--- Mandatory String Convert
	num, errs = strconv.ParseFloat(value, 64)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input Absent) optionalStrInt(value string) (num int64, err errorModel.ErrorModel) {
	var (
		fileName = "AbsentDTO.go"
		funcName = "optionalStrInt"
	)

	//--- Optional Check Empty
	if !util.IsStringEmpty(value) {

		n, errs := strconv.Atoi(value)
		if errs != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errs)
			return
		}

		num = int64(n)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
