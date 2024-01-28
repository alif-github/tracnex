package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"strings"
	"time"
)

type EmployeeContractRequest struct {
	AbstractDTO
	ID           int64  `json:"id"`
	Contract     string `json:"contract"`
	Information  string `json:"information"`
	EmployeeID   int64  `json:"employee_id"`
	FromDateStr  string `json:"from_date"`
	ThruDateStr  string `json:"thru_date"`
	UpdatedAtStr string `json:"updated_at"`
	FromDate     time.Time
	ThruDate     time.Time
	UpdatedAt    time.Time
}

func (input *EmployeeContractRequest) ValidateInsert() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeeContractDTO.go"
		funcName = "ValidateInsert"
	)

	//--- Employee ID Check Empty
	if input.EmployeeID < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.EmployeeID)
		return
	}

	return input.mandatoryFieldValidation(fileName, funcName)
}

func (input *EmployeeContractRequest) ValidateView() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeeContractDTO.go"
		funcName = "ValidateView"
	)

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *EmployeeContractRequest) ValidateDelete() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeeContractDTO.go"
		funcName = "ValidateDelete"
	)

	return input.validationForUpdateAndDelete(fileName, funcName)
}

func (input *EmployeeContractRequest) ValidateUpdate() (err errorModel.ErrorModel) {
	var (
		fileName = "EmployeeContractDTO.go"
		funcName = "ValidateUpdate"
	)

	if err = input.validationForUpdateAndDelete(fileName, funcName); err.Error != nil {
		return
	}

	return input.mandatoryFieldValidation(fileName, funcName)
}

func (input *EmployeeContractRequest) validationForUpdateAndDelete(fileName, funcName string) (err errorModel.ErrorModel) {
	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	if util.IsStringEmpty(input.UpdatedAtStr) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *EmployeeContractRequest) mandatoryFieldValidation(fileName, funcName string) (err errorModel.ErrorModel) {
	//--- Contract Check Empty
	if util.IsStringEmpty(input.Contract) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ContractNo)
		return
	}

	//--- Trim Contract Check Empty
	sc := strings.Trim(input.Contract, " ")
	if util.IsStringEmpty(sc) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ContractNo)
		return
	}

	//--- Contract Check Min Max
	err = input.ValidateMinMaxString(input.Contract, constanta.ContractNo, 1, 50)
	if err.Error != nil {
		return
	}

	//--- Information Check Empty
	if !util.IsStringEmpty(input.Information) {
		//--- Trim Information Check Empty
		si := strings.Trim(input.Information, " ")
		if util.IsStringEmpty(si) {
			err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Information)
			return
		}

		//--- Information Check Min Max
		err = input.ValidateMinMaxString(input.Information, constanta.Information, 1, 256)
		if err.Error != nil {
			return
		}
	}

	//--- From Date Check Empty
	if util.IsStringEmpty(input.FromDateStr) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.FromDate)
		return
	}

	//--- Date Join Check Format Time
	input.FromDate, err = TimeStrToTimeWithTimeFormat(input.FromDateStr, constanta.FromDate, constanta.DefaultInstallationTimeFormat)
	if err.Error != nil {
		return
	}

	//--- Thru Date Check Empty
	if util.IsStringEmpty(input.ThruDateStr) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ThruDate)
		return
	}

	//--- Date Out Check Format Time
	input.ThruDate, err = TimeStrToTimeWithTimeFormat(input.ThruDateStr, constanta.ThruDate, constanta.DefaultInstallationTimeFormat)
	if err.Error != nil {
		return
	}

	//--- Thru Date After From Date Check
	if input.ThruDate.Before(input.FromDate) || input.ThruDate.Equal(input.FromDate) {
		err = errorModel.GenerateDateValidateFromThru(fileName, funcName, "E-6-TRAC-SRV-016")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
