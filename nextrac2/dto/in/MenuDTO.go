package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

const fileNameMenu = "MenuDTO.go"

type MenuRequest struct {
	AbstractDTO
	ID				int64		`json:"id"`
	ParentMenuID	int64		`json:"parent_menu_id"`
	ServiceMenuID	int64		`json:"service_menu_id"`
	Name			string		`json:"name"`
	EnName			string		`json:"en_name"`
	Sequence		int64		`json:"sequence"`
	IconName		string		`json:"icon_name"`
	Background		string		`json:"background"`
	AvailableAction	string		`json:"available_action"`
	MenuCode		string		`json:"menu_code"`
	Status			string		`json:"status"`
	Url				string		`json:"url"`
	CreatedBy		int64		`json:"created_by"`
	CreatedClient	string		`json:"created_client"`
	CreatedAtStr	string		`json:"created_at"`
	CreatedAt		time.Time
	UpdatedBy		int64		`json:"updated_by"`
	UpdatedClient	string		`json:"updated_client"`
	UpdatedAtStr	string		`json:"updated_at"`
	UpdatedAt		time.Time
}

func (input *MenuRequest) ValidateUpdateMenuItem() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdateMenuItem"

	err = input.doForUpdateService()
	if err.Error != nil {return}

	if input.ServiceMenuID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileNameMenu, funcName, "Service Menu ID")
	}

	if util.IsStringEmpty(input.Url) {
		return errorModel.GenerateEmptyFieldError(fileNameMenu, funcName, "URL")
	}

	err = input.validateMandatoryMenu()
	if err.Error != nil {return}

	return errorModel.GenerateNonErrorModel()
}

func (input *MenuRequest) ValidateUpdateMenuService() (err errorModel.ErrorModel) {
	funcName := "ValidateUpdateMenuService"

	err = input.doForUpdateService()
	if err.Error != nil {return}

	if input.ParentMenuID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileNameMenu, funcName, "Parent Menu ID")
	}

	err = input.validateMandatoryMenu()
	if err.Error != nil {return}

	return errorModel.GenerateNonErrorModel()
}

func (input *MenuRequest) ValidateUpdateMenuParent() (err errorModel.ErrorModel) {
	err = input.doForUpdateService()
	if err.Error != nil {return}

	err = input.validateMandatoryMenu()
	if err.Error != nil {return}

	return errorModel.GenerateNonErrorModel()
}

func (input MenuRequest) validateMandatoryMenu() (err errorModel.ErrorModel) {
	funcName := "validateMandatoryUpdateUser"
	var validationResult bool
	var errField string
	var additionalInfo string

	if util.IsStringEmpty(input.Name) {
		return errorModel.GenerateEmptyFieldError(fileNameMenu, funcName, "Name")
	}

	validationResult, errField, additionalInfo = util.IsNexsoftNameStandardValid(input.Name)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileNameMenu, funcName, errField, "Name", additionalInfo)
	}

	err = input.ValidateMinMaxString(input.Name, "Name", 1, 50)
	if err.Error != nil {
		return err
	}

	if util.IsStringEmpty(input.EnName) {
		return errorModel.GenerateEmptyFieldError(fileNameMenu, funcName, "En-Name")
	}

	validationResult, errField, additionalInfo = util.IsNexsoftNameStandardValid(input.EnName)
	if !validationResult {
		return errorModel.GenerateFieldFormatWithRuleError(fileNameMenu, funcName, errField, "En-Name", additionalInfo)
	}

	err = input.ValidateMinMaxString(input.Name, "En-Name", 1, 50)
	if err.Error != nil {
		return err
	}

	if input.Sequence < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileNameMenu, funcName, "sequence")
	}

	if !util.IsStringEmpty(input.IconName) {
		err = input.ValidateMinMaxString(input.IconName, "Icon Name", 1, 30)
		if err.Error != nil {
			return err
		}
	}

	if !util.IsStringEmpty(input.Background) {
		err = input.ValidateMinMaxString(input.Background, "Background", 1, 10)
		if err.Error != nil {
			return err
		}
	}

	if util.IsStringEmpty(input.AvailableAction) {
		return errorModel.GenerateEmptyFieldError(fileNameMenu, funcName, "Available Action")
	}

	err = input.ValidateMinMaxString(input.AvailableAction, "Available Action", 1, 256)
	if err.Error != nil {
		return err
	}

	if util.IsStringEmpty(input.MenuCode) {
		return errorModel.GenerateEmptyFieldError(fileNameMenu, funcName, "Menu Code")
	}

	err = input.ValidateMinMaxString(input.MenuCode, "Menu Code", 1, 256)
	if err.Error != nil {
		return err
	}

	if util.IsStringEmpty(input.Status) {
		return errorModel.GenerateEmptyFieldError(fileNameMenu, funcName, "Status")
	}

	err = input.ValidateMinMaxString(input.Status, "Status", 1, 256)
	if err.Error != nil {
		return err
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *MenuRequest) doForUpdateService() (err errorModel.ErrorModel) {
	funcName := "doForUpdateService"

	if input.ID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileNameMenu, funcName, constanta.ID)
	}

	if util.IsStringEmpty(input.UpdatedAtStr) {
		return errorModel.GenerateEmptyFieldError(fileNameMenu, funcName, constanta.UpdatedAt)
	}

	input.UpdatedAt, err = TimeStrToTime(input.UpdatedAtStr, constanta.UpdatedAt)
	if err.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}