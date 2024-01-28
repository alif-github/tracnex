package in

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type CustomerSiteRequest struct {
	AbstractDTO
	ID               int64  `json:"id"`
	ParentCustomerID int64  `json:"parent_customer_id"`
	CustomerID       int64  `json:"customer_id"`
	UpdatedAtStr     string `json:"updated_at"`
	UpdatedAt        time.Time
}

func (input *CustomerSiteRequest) ValidateInsertCustomerSite() errorModel.ErrorModel {
	return input.validateMandatory()
}

func (input *CustomerSiteRequest) validateMandatory() errorModel.ErrorModel {
	fileName := "CustomerSiteDTO.go"
	funcName := "validateMandatory"

	if input.ParentCustomerID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ParentCustomerID)
	}

	if input.CustomerID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.CustomerID)
	}

	return errorModel.GenerateNonErrorModel()
}