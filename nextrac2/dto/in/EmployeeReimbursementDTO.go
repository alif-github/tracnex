package in

import (
	"fmt"
	"mime/multipart"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"strings"
	"time"
)

type EmployeeReimbursementRequest struct {
	Name            string                           `json:"name"`
	BenefitId       int64                            `json:"benefit_id"`
	Receipt         string                           `json:"receipt"`
	Value           float64                          `json:"value"`
	Description     string                           `json:"description"`
	ReceiptDate     string                           `json:"receipt_date"`
	ReceiptDateTime time.Time                        `json:"-"`
	EmployeeId      int64                            `json:"-"`
	FileUploadId    int64                            `json:"-"`
	Attachment      *EmployeeReimbursementAttachment `json:"-"`
	AbstractDTO
}

type EmployeeReimbursementVerifyRequest struct {
	Id                  int64      `json:"-"`
	IdCard              string     `json:"id_card"`
	EmployeeId          int64      `json:"employee_id"`
	CurrentMedicalValue float64    `json:"current_medical_value"`
	ApprovedValue       float64    `json:"approved_value"`
	Note                string     `json:"note"`
	UpdatedAtStr        string     `json:"updated_at"`
	UpdatedAt           time.Time
	AbstractDTO
}

type EmployeeReimbursementAttachment struct {
	File       multipart.File        `json:"-"`
	FileHeader *multipart.FileHeader `json:"-"`
}

func (input *EmployeeReimbursementRequest) ValidateInsert(expiredMedicalClaimParameter int) errorModel.ErrorModel {
	fileName := "EmployeeReimbursementDTO.go"
	funcName := "ValidateInsert"

	input.Receipt = strings.Trim(input.Receipt, " ")
	input.Description = strings.Trim(input.Description, " ")

	/*
		Benefit ID
	*/
	if input.BenefitId <= 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "benefit_id")
	}

	/*
		Receipt No
	*/
	if errModel := input.ValidateMinMaxString(input.Receipt, "receipt", 1, 256); errModel.Error != nil {
		return errModel
	}

	/*
		Description
	*/
	if errModel := input.ValidateMinMaxString(input.Description, "description", 0, 256); errModel.Error != nil {
		return errModel
	}

	/*
		Value
	*/
	if input.Value <= 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "value")
	}

	if input.Value >= 10000000000 {
		return errorModel.GenerateValueMustLessThanError(fileName, funcName, "value", "Rp. 10.000.000.000")
	}

	if errModel := input.validateDecimalPlaces(input.Value, 2, "value"); errModel.Error != nil {
		return errModel
	}

	/*
		Receipt Date
	*/
	if input.ReceiptDate == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, "receipt_date")
	}

	receiptDateTime, err := time.Parse("2006-01-02", input.ReceiptDate)
	if err != nil {
		return errorModel.GenerateFormatFieldError(fileName, funcName, "receipt_date")
	}

	now, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	diff := int(now.Sub(receiptDateTime).Hours() / 24)

	if diff > expiredMedicalClaimParameter {
		return errorModel.GenerateDataExpiredError(fileName, funcName, "receipt_date")
	}

	input.ReceiptDateTime = receiptDateTime

	return errorModel.GenerateNonErrorModel()
}

func (input *EmployeeReimbursementRequest) validateDecimalPlaces(value float64, maxDecimalPlaces int, fieldName string) errorModel.ErrorModel {
	fileName := "EmployeeReimbursementDTO.go"
	funcName := "validateDecimalPlaces"

	strValue := fmt.Sprintf("%.3f", value)
	separatedValues := strings.Split(strValue, ".")

	if len(separatedValues) < 2 {
		return errorModel.GenerateNonErrorModel()
	}

	decimals := strings.TrimRight(separatedValues[1], "0")

	if len(decimals) > maxDecimalPlaces {
		return errorModel.GenerateMaxDecimalPlacesError(fileName, funcName, fieldName, maxDecimalPlaces)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *EmployeeReimbursementVerifyRequest) ValidateVerifyReimbursement(isUpdate bool) (err errorModel.ErrorModel) {
	funcName := "ValidateVerifyReimbursement"
	fileName := "EmployeeReimbursementDTO.go"

    if input.Note == ""{
    	return errorModel.GenerateEmptyFieldError(fileName, funcName, "note")
	}

	err = input.ValidateMinMaxString(input.Note, "note", 1, 256)
	if err.Error != nil {
		return
	}

	if isUpdate {
		if util.IsStringEmpty(input.UpdatedAtStr) {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
		}

		input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
		if err.Error != nil {
			return
		}
	}

	return errorModel.GenerateNonErrorModel()
}