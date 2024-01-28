package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type CustomerListRequest struct {
	AbstractDTO
	ID				int64					`json:"id"`
	CompanyID		string					`json:"company_id"`
	BranchData		[]CustomerBranchData	`json:"branch_data"`
}

type CustomerBranchData struct {
	ID					int64		`json:"id"`
	BranchID			string		`json:"branch_id"`
	CompanyName			string		`json:"company_name"`
	City				string		`json:"city"`
	Implementer			string		`json:"implementer"`
	ImplementationAtStr	string		`json:"implementation_at"`
	Product				string		`json:"product"`
	Version				string		`json:"version"`
	LicenseType			string		`json:"license_type"`
	UserOnLicense		int64		`json:"user_on_license"`
	ExpDateAtStr		string		`json:"exp_date_at"`
	ImplementationAt	time.Time
	ExpDateAt			time.Time
}

type CustomerListImportRequest struct {
	CompanyID		string		`json:"company_id"`
	BranchID		string		`json:"branch_id"`
	CompanyName		string		`json:"company_name"`
	City			string		`json:"city"`
	Implementer		string		`json:"implementer"`
	Implementation	time.Time	`json:"implementation"`
	Product			string		`json:"product"`
	Version			string		`json:"version"`
	LicenseType		string		`json:"license_type"`
	UserAmount		int			`json:"user_amount"`
	ExpDate			time.Time	`json:"exp_date"`
}

func (input *CustomerListRequest) ValidateInsertCustomer() errorModel.ErrorModel {
	fileName := "CustomerListDTO.go"
	funcName := "ValidateInsertCustomer"

	var err errorModel.ErrorModel
	var validationResult bool

	err = input.validateMandatoryCustomer()
	if err.Error != nil {
		return err
	}

	for _, branchDataElm := range input.BranchData {

		validationResult = util.IsStringEmpty(branchDataElm.CompanyName)
		if validationResult {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CompanyName)
		}

		err = input.ValidateMinMaxString(branchDataElm.CompanyName, constanta.CompanyName, 1, 50)
		if err.Error != nil {
			return err
		}

		if branchDataElm.City != "" {
			err = input.ValidateMinMaxString(branchDataElm.City, constanta.City, 1, 50)
			if err.Error != nil {
				return err
			}
		}

		if branchDataElm.Implementer != "" {
			err = input.ValidateMinMaxString(branchDataElm.Implementer, constanta.Implementer, 1, 50)
			if err.Error != nil {
				return err
			}
		}

		if branchDataElm.ImplementationAtStr != "" {
			branchDataElm.ImplementationAt, err = TimeStrToTime(branchDataElm.ImplementationAtStr, constanta.ImplementationAt)
			if err.Error != nil {
				return err
			}
		}

		validationResult = util.IsStringEmpty(branchDataElm.Product)
		if validationResult {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Product)
		}

		err = input.ValidateMinMaxString(branchDataElm.Product, constanta.Product, 1, 50)
		if err.Error != nil {
			return err
		}

		if branchDataElm.Version != "" {
			err = input.ValidateMinMaxString(branchDataElm.Version, constanta.Version, 1, 50)
			if err.Error != nil {
				return err
			}
		}

		if branchDataElm.LicenseType != "" {
			err = input.ValidateMinMaxString(branchDataElm.LicenseType, constanta.LicenseType, 1, 50)
			if err.Error != nil {
				return err
			}
		}

		validationResult = util.IsStringEmpty(branchDataElm.ExpDateAtStr)
		if validationResult {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ExpDate)
		}

		branchDataElm.ExpDateAt, err = TimeStrToTime(branchDataElm.ExpDateAtStr, constanta.ExpDate)
		if err.Error != nil {
			return err
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *CustomerListRequest) validateMandatoryCustomer() errorModel.ErrorModel {
	fileName := "CustomerListDTO.go"
	funcName := "validateMandatoryCustomer"
	var validationResult bool
	var errField string
	var additionalInfo string
	var err errorModel.ErrorModel

	validationResult = util.IsStringEmpty(input.CompanyID)
	if validationResult {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CompanyID)
	}

	validationResult, errField, additionalInfo = IsOnlyWordCharacterValid(input.CompanyID)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, "company_id", additionalInfo)
	}

	err = input.ValidateMinMaxString(input.CompanyID, constanta.CompanyID, 1, 20)
	if err.Error != nil {
		return err
	}

	for _, branchDataElm := range input.BranchData {

		validationResult = util.IsStringEmpty(branchDataElm.BranchID)
		if validationResult {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.BranchID)
		}

		err = input.ValidateMinMaxString(branchDataElm.BranchID, constanta.BranchID, 1, 20)
		if err.Error != nil {
			return err
		}

		validationResult, errField, additionalInfo = IsOnlyDigitValid(branchDataElm.BranchID)
		if !validationResult {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, "branch_id", additionalInfo)
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *CustomerListRequest) ValidateViewCustomer() errorModel.ErrorModel {
	fileName := "CustomerListDTO.go"
	funcName := "ValidateViewCustomer"

	if input.ID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
	}

	return errorModel.GenerateNonErrorModel()
}