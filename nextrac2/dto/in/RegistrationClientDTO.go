package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type ClientRequest struct {
	AbstractDTO
	URLRequest
	BodyRequest
	ClientTypeID int64         `json:"client_type_id"`
	ClientName   string        `json:"client_name"`
	SocketID     string        `json:"socket_id"`
	CompanyData  []CompanyData `json:"company_data"`
}

type CompanyData struct {
	ID          int64        `json:"id"`
	CompanyID   string       `json:"company_id"`
	CompanyName string       `json:"company_name"`
	BranchData  []BranchData `json:"branch_data"`
}

type BranchData struct {
	ClientID     string `json:"client_id"`
	ClientAlias  string `json:"client_alias"`
	BranchName   string `json:"branch_name"`
	BranchID     string `json:"branch_id"`
	UpdatedAtStr string `json:"updated_at"`
	UpdatedAt    time.Time
}

func (input *ClientRequest) ValidateRegisClient() errorModel.ErrorModel {
	err := input.validateMandatoryClient()
	if err.Error != nil {
		return err
	}

	if input.SocketID != "" {
		err = input.ValidateMinMaxString(input.SocketID, constanta.SocketUser, 1, 50)
		if err.Error != nil {
			return err
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *ClientRequest) validateMandatoryClient() errorModel.ErrorModel {
	fileName := "RegistrationClientDTOIn.go"
	funcName := "validateMandatoryClient"
	var validationResult bool
	var err errorModel.ErrorModel
	var errField string
	var additionalInfo string
	var arrayCompanyAndBranchID []string

	if util.IsStringEmpty(input.ClientName) {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientName)
	}

	validationResult, errField, additionalInfo = util.IsNexsoftNameStandardValid(input.ClientName)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, errField, constanta.ClientName, additionalInfo)
	}

	err = input.ValidateMinMaxString(input.ClientName, constanta.ClientName, 1, 50)
	if err.Error != nil {
		return err
	}

	if len(input.CompanyData) < 1 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CompanyID)
	}

	for _, companyDataElm := range input.CompanyData {
		if util.IsStringEmpty(companyDataElm.CompanyID) {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CompanyID)
		}

		err = input.ValidateMinMaxString(companyDataElm.CompanyID, constanta.CompanyID, 1, 20)
		if err.Error != nil {
			return err
		}

		if len(companyDataElm.BranchData) < 1 {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.BranchID)
		}

		for _, branchDataElm := range companyDataElm.BranchData {
			//--- Branch Name
			//if util.IsStringEmpty(branchDataElm.BranchName) {
			//	return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.BranchName)
			//}
			//
			//err = input.ValidateMinMaxString(branchDataElm.BranchName, constanta.BranchName, 1, 50)
			//if err.Error != nil {
			//	return err
			//}

			//--- Branch Name
			if !util.IsStringEmpty(branchDataElm.BranchName) {
				err = input.ValidateMinMaxString(branchDataElm.BranchName, constanta.BranchName, 1, 50)
				if err.Error != nil {
					return err
				}
			}

			//--- Branch ID
			if util.IsStringEmpty(branchDataElm.BranchID) {
				return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.BranchID)
			}

			err = input.ValidateMinMaxString(branchDataElm.BranchID, constanta.BranchID, 1, 20)
			if err.Error != nil {
				return err
			}

			//------ Append company id + branch id
			arrayCompanyAndBranchID = append(arrayCompanyAndBranchID, companyDataElm.CompanyID+branchDataElm.BranchID)
		}
	}

	//------ Check unique request data
	err = input.UniqueStrArray(fileName, funcName, arrayCompanyAndBranchID)
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}
