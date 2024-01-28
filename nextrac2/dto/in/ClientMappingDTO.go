package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type ClientMappingRequest struct {
	AbstractDTO
	ClientTypeID int64               `json:"client_type_id"`
	ClientID     string              `json:"client_id"`
	CompanyData  []CompanyDataClient `json:"company_data"`
	ClientData   []ClientIDData      `json:"client_data"`
}

type CompanyDataClient struct {
	CompanyID  string             `json:"company_id"`
	BranchData []BranchDataClient `json:"branch_data"`
}

type BranchDataClient struct {
	ID             int64  `json:"id"`
	BranchName     string `json:"branch_name"`
	BranchID       string `json:"branch_id"`
	ClientAlias    string `json:"client_alias"`
	InstallationID int64  `json:"installation_id"`
	CustomerID     string `json:"customer_id"`
	SiteID         int64  `json:"site_id"`
	CreatedBy      int64  `json:"created_by"`
	CreatedClient  string `json:"created_client"`
	UpdatedBy      int64  `json:"updated_by"`
	UpdatedClient  string `json:"updated_client"`
	UpdatedAtStr   string `json:"updated_at_str"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type ClientIDData struct {
	ClientID string `json:"client_id"`
}

func (input *ClientMappingRequest) ValidateInsertClientMapping() errorModel.ErrorModel {
	return input.validateMandatoryClientMapping()
}

func (input *ClientMappingRequest) validateMandatoryClientMapping() errorModel.ErrorModel {
	var (
		fileName       = "ClientMappingDTO.go"
		funcName       = "validateMandatoryClientMapping"
		err            errorModel.ErrorModel
		arrayCompanyID []string
	)

	for _, companyDataElm := range input.CompanyData {
		if util.IsStringEmpty(companyDataElm.CompanyID) {
			return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CompanyID)
		}

		err = input.ValidateMinMaxString(companyDataElm.CompanyID, constanta.CompanyID, 1, 20)
		if err.Error != nil {
			return err
		}

		var arrayBranchID []string
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

			//--- Append to array
			arrayBranchID = append(arrayBranchID, branchDataElm.BranchID)
		}
		//------- Append to array
		arrayCompanyID = append(arrayCompanyID, companyDataElm.CompanyID)

		//------- Check unique request data
		err = input.UniqueStrArray(fileName, funcName, arrayBranchID)
		if err.Error != nil {
			return err
		}
	}

	//------- Check unique request data
	err = input.UniqueStrArray(fileName, funcName, arrayCompanyID)
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}

func (input ClientMappingRequest) ValidationForDetail() (err errorModel.ErrorModel) {
	fileName := "ClientMappingDTO.go"
	funcName := "validateMandatoryClientMapping"

	if input.ClientTypeID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, "client_type_id")
		return
	}

	if len(input.CompanyData) < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, "company_data")
		return
	}

	for _, companyData := range input.CompanyData {
		if util.IsStringEmpty(companyData.CompanyID) {
			err = errorModel.GenerateEmptyFieldError(fileName, funcName, "company_id")
			return
		}
		for _, branchData := range companyData.BranchData {
			if util.IsStringEmpty(branchData.BranchID) {
				err = errorModel.GenerateEmptyFieldError(fileName, funcName, "branch_id")
				return
			}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input *ClientMappingRequest) ValidationForGetClientMappingByClientID() (err errorModel.ErrorModel) {
	fileName := "ClientMappingDTO.go"
	funcName := "ValidationForGetClientMappingByClientID"

	if input.ClientTypeID < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, "client_type_id")
		return
	}

	if len(input.ClientData) < 1 {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, "client_data")
		return
	}

	for _, clientData := range input.ClientData {
		err = input.ValidateMinMaxString(clientData.ClientID, "constanta.ClientID", 1, 256)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
