package in

import (
	"mime/multipart"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type EmployeeLeaveRequest struct {
	Name          string                   `json:"name"`
	Type          string                   `json:"type"`
	AllowanceId   int64                    `json:"allowance_id"`
	StrDateList   []string                 `json:"date"`
	StrLeaveDate  string                   `json:"leave_date"`
	StrLeaveTime  string                   `json:"leave_time"`
	StrReturnTime string                   `json:"return_time"`
	DateList      []time.Time              `json:"-"`
	Description   string                   `json:"description"`
	Value         int64                    `json:"value"`
	FileUploadId  int64                    `json:"-"`
	EmployeeId    int64                    `json:"-"`
	UpdatedAt     time.Time                `json:"updated_at"`
	Attachment    *EmployeeLeaveAttachment `json:"-"`
	AbstractDTO
}

type LeaveDetail struct {
	DateStr     []string `json:"date"`
	Type        string   `json:"type"`
	Description string   `json:"description"`
}

type EmployeeLeaveAttachment struct {
	File       multipart.File        `json:"-"`
	FileHeader *multipart.FileHeader `json:"-"`
}

type ReportAnnualLeave struct {
	YearStr  string `json:"year"`
	IsDetail bool   `json:"is_detail"`
	Year     []int64
}

func (input *ReportAnnualLeave) ValidateReportAnnualLeave() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeeLeaveDTO.go"
		funcName = "ValidateReportAnnualLeave"
	)

	if util.IsStringEmpty(input.YearStr) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, "Year")
		return
	}

	input.YearStr = strings.Trim(input.YearStr, " ")
	numYear, errs := strconv.Atoi(input.YearStr)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	input.Year = append(input.Year, int64(numYear-1), int64(numYear))
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *EmployeeLeaveRequest) ValidateInsert() errorModel.ErrorModel {
	fileName := "EmployeeLeaveDTO.go"
	funcName := "ValidateInsert"

	if input.Type == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "type")
	}

	if input.Type != constanta.LeaveType && input.Type != constanta.PermitType && input.Type != constanta.SickLeaveType {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "type")
	}

	if errModel := input.validateLeaveType(); errModel.Error != nil {
		return errModel
	}

	if errModel := input.validatePermitType(); errModel.Error != nil {
		return errModel
	}

	if errModel := input.validateSickLeaveType(); errModel.Error != nil {
		return errModel
	}

	if errModel := input.ValidateMinMaxString(input.Description, "description", 1, 256); errModel.Error != nil {
		return errModel
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *EmployeeLeaveRequest) validateLeaveType() errorModel.ErrorModel {
	fileName := "EmployeeLeaveDTO.go"
	funcName := "validateLeaveType"

	if input.Type != constanta.LeaveType {
		return errorModel.GenerateNonErrorModel()
	}

	if input.AllowanceId <= 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "allowance_id")
	}

	if len(input.StrDateList) == 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "date")
	}

	var tempTime *time.Time

	for _, date := range input.StrDateList {
		result, errModel := input.validateDate(date, constanta.LeaveType, "date")
		if errModel.Error != nil {
			return errModel
		}

		if tempTime != nil {
			if !result.After(*tempTime) {
				return errorModel.GenerateFormatFieldError(fileName, funcName, "date")
			}
		}

		tempTime = &result
		input.DateList = append(input.DateList, result)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *EmployeeLeaveRequest) validateSickLeaveType() errorModel.ErrorModel {
	fileName := "EmployeeLeaveDTO.go"
	funcName := "validateSickLeaveType"

	if input.Type != constanta.SickLeaveType {
		return errorModel.GenerateNonErrorModel()
	}

	var tempTime *time.Time

	for _, date := range input.StrDateList {
		result, errModel := input.validateDate(date, constanta.SickLeaveType, "date")
		if errModel.Error != nil {
			return errModel
		}

		if tempTime != nil {
			if !result.After(*tempTime) {
				return errorModel.GenerateFormatFieldError(fileName, funcName, "date")
			}
		}

		tempTime = &result
		input.DateList = append(input.DateList, result)
	}

	if input.Attachment == nil {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "file")
	}

	if input.Attachment.FileHeader.Size > 5*1e6 {
		return errorModel.GenerateFileSizeExceedsMaxLimitError(fileName, funcName, "file", 5, "MB")
	}

	if errModel := input.validateSickLeaveAttachmentExtension(); errModel.Error != nil {
		return errModel
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *EmployeeLeaveRequest) validateSickLeaveAttachmentExtension() errorModel.ErrorModel {
	fileName := "EmployeeLeaveDTO.go"
	funcName := "validateSickLeaveAttachmentExtension"

	allowedExtensions := []string{"pdf", "xls", "xlsx", "jpg", "jpeg", "png"}

	separatedFileParts := strings.Split(input.Attachment.FileHeader.Filename, ".")
	fileExtension := separatedFileParts[len(separatedFileParts)-1]

	for _, extension := range allowedExtensions {
		if fileExtension == extension {
			return errorModel.GenerateNonErrorModel()
		}
	}

	return errorModel.GenerateInvalidFileExtensionError(fileName, funcName, "file", allowedExtensions)
}

func (input *EmployeeLeaveRequest) validatePermitType() errorModel.ErrorModel {
	fileName := "EmployeeLeaveDTO.go"
	funcName := "validatePermitType"

	if input.Type != constanta.PermitType {
		return errorModel.GenerateNonErrorModel()
	}

	if input.StrLeaveDate == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "leave_date")
	}

	if input.StrLeaveTime == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "leave_time")
	}

	if input.StrReturnTime == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "return_time")
	}

	leaveDate, errModel := input.validateDate(input.StrLeaveDate, constanta.PermitType, "leave_date")
	if errModel.Error != nil {
		return errModel
	}

	reg := regexp.MustCompile("^([01]?[0-9]|2[0-3]):[0-5][0-9]$")
	if !reg.MatchString(input.StrLeaveTime) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "leave_time")
	}

	if !reg.MatchString(input.StrReturnTime) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "return_time")
	}

	separatedLeaveTime := strings.Split(input.StrLeaveTime, ":")
	leaveTimeHour, _ := strconv.Atoi(separatedLeaveTime[0])
	leaveTimeMinute, _ := strconv.Atoi(separatedLeaveTime[1])

	separatedReturnTime := strings.Split(input.StrReturnTime, ":")
	returnTimeHour, _ := strconv.Atoi(separatedReturnTime[0])
	returnTimeMinute, _ := strconv.Atoi(separatedReturnTime[1])

	leaveTime := time.Date(leaveDate.Year(), leaveDate.Month(), leaveDate.Day(), leaveTimeHour, leaveTimeMinute, 0, 0, leaveDate.Location())
	returnTime := time.Date(leaveDate.Year(), leaveDate.Month(), leaveDate.Day(), returnTimeHour, returnTimeMinute, 0, 0, leaveDate.Location())

	if leaveTime.After(returnTime) {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "leave_time")
	}

	input.DateList = append(input.DateList, leaveTime, returnTime)

	return errorModel.GenerateNonErrorModel()
}

func (input *EmployeeLeaveRequest) validateDate(strDate, leaveType, fieldName string) (result time.Time, errModel errorModel.ErrorModel) {
	fileName := "EmployeeLeaveDTO.go"
	funcName := "validateDate"

	date, err := time.Parse("2006-01-02", strDate)
	if err != nil {
		errModel = errorModel.GenerateFormatFieldError(fileName, funcName, fieldName)
		return
	}

	if input.isWeekend(date) {
		errModel = errorModel.GenerateDateContainsWeekendError(fileName, funcName, fieldName)
		return
	}

	now, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))

	if leaveType == constanta.LeaveType && !date.After(now) {
		errModel = errorModel.GenerateDateMustBeLaterThanCurrentDateError(fileName, funcName, fieldName)
		return
	}

	if leaveType == constanta.PermitType && date.Before(now) {
		errModel = errorModel.GenerateDateCannotBeLessThanCurrentDate(fileName, funcName, fieldName)
		return
	}

	return date, errorModel.GenerateNonErrorModel()
}

func (input *EmployeeLeaveRequest) isWeekend(date time.Time) bool {
	return date.Weekday() == time.Saturday || date.Weekday() == time.Sunday
}
