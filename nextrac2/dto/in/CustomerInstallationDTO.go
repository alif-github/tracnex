package in

import (
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type CustomerSiteInstallationRequest struct {
	AbstractDTO
	ID               int64                     `json:"id"`
	ParentCustomerID int64                     `json:"parent_customer_id"`
	CustomerSite     []CustomerSiteDataRequest `json:"customer_site"`
}

type CustomerSiteDataRequest struct {
	SiteID               int64                             `json:"site_id"`
	CustomerID           int64                             `json:"customer_id"`
	Action               int32                             `json:"action"`
	UpdatedAtStr         string                            `json:"updated_at"`
	CustomerInstallation []CustomerInstallationDataRequest `json:"customer_installation"`
	UpdatedAt            time.Time
}

type CustomerInstallationDataRequest struct {
	UniqueKey           int64  `json:"unique_key"`
	InstallationID      int64  `json:"installation_id"`
	ProductID           int64  `json:"product_id"`
	ParentClientTypeID  int64  `json:"parent_client_type_id"`
	Remark              string `json:"remark"`
	UniqueID1           string `json:"unique_id_1"`
	UniqueID2           string `json:"unique_id_2"`
	InstallationStatus  string `json:"installation_status"`
	InstallationDateStr string `json:"installation_date"`
	ProductValidFromStr string `json:"product_valid_from"`
	ProductValidThruStr string `json:"product_valid_thru"`
	UpdatedAtStr        string `json:"updated_at"`
	Action              int32  `json:"action"`
	InstallationDate    time.Time
	ProductValidFrom    time.Time
	ProductValidThru    time.Time
	UpdatedAt           time.Time
}

type CustomerInstallationDetailRequest struct {
	ParentCustomerID int64 `json:"parent_customer_id"`
	SiteID           int64 `json:"site_id"`
	ClientTypeID     int64 `json:"client_type_id"`
	IsLicense        bool  `json:"is_license"`
}

func (input CustomerInstallationDetailRequest) ValidateViewCustomerInstallation() (err errorModel.ErrorModel) {
	fileName := "CustomerInstallationDTO.go"
	funcName := "ValidateViewCustomerInstallation"

	if input.ParentCustomerID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ParentCustomerID)
	}

	if input.SiteID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.SiteID)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input CustomerSiteInstallationRequest) ValidateViewCustomerSite() (err errorModel.ErrorModel) {
	fileName := "CustomerInstallationDTO.go"
	funcName := "ValidateViewCustomerInstallation"

	if input.ParentCustomerID < 1 {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ParentCustomerID)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *CustomerSiteInstallationRequest) ValidateUpdateCustomerInstallation(contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		fileName = "CustomerInstallationDTO.go"
		funcName = "ValidateUpdateCustomerInstallation"
		number   = fmt.Sprintf(`no.`)
		msgStr   string
	)

	//--- Parent Customer ID
	if input.ParentCustomerID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ParentCustomerID)
	}

	//--- Customer Site
	for idx, itemCustomerSite := range input.CustomerSite {

		//--- Check Customer ID on Site
		if itemCustomerSite.CustomerID < 1 {
			custStr := util2.GenerateConstantaI18n(constanta.CustomerID, contextModel.AuthAccessTokenModel.Locale, nil)
			msgStr = fmt.Sprintf(`%s Site %s %d`, custStr, number, idx+1)
			return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, msgStr)
		}

		//--- Check Updated At on Site
		if !util.IsStringEmpty(itemCustomerSite.UpdatedAtStr) {
			uptStr := util2.GenerateConstantaI18n(constanta.UpdatedAt, contextModel.AuthAccessTokenModel.Locale, nil)
			msgStr = fmt.Sprintf(`%s Site %s %d`, uptStr, number, idx+1)
			input.CustomerSite[idx].UpdatedAt, err = TimeStrToTime(itemCustomerSite.UpdatedAtStr, msgStr)
			if err.Error != nil {
				return err
			}
		}

		//--- Check Action on Site
		if (itemCustomerSite.Action != int32(constanta.ActionInsertCode)) && (itemCustomerSite.Action != int32(constanta.ActionDeleteCode)) && (itemCustomerSite.Action != int32(constanta.ActionNoActionCode)) {
			actStr := util2.GenerateConstantaI18n(constanta.ActionCode, contextModel.AuthAccessTokenModel.Locale, nil)
			msgStr = fmt.Sprintf(`%s Site %s %d`, actStr, number, idx+1)
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.ActionCodeRegex1, msgStr, "")
		}

		//--- Customer Installation
		for idx2, itemCustomerInstallation := range itemCustomerSite.CustomerInstallation {
			msgSiteStr := fmt.Sprintf(`Site %s %d`, number, idx+1)

			if itemCustomerInstallation.ProductID < 1 {
				prdStr := util2.GenerateConstantaI18n(constanta.ProductID, contextModel.AuthAccessTokenModel.Locale, nil)
				msgStr = fmt.Sprintf(`%s Installation %s %d %s`, prdStr, number, idx2+1, msgSiteStr)
				return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, msgStr)
			}

			if !util.IsStringEmpty(itemCustomerInstallation.Remark) {
				rmkStr := util2.GenerateConstantaI18n(constanta.Remark, contextModel.AuthAccessTokenModel.Locale, nil)
				msgStr = fmt.Sprintf(`%s Installation %s %d %s`, rmkStr, number, idx2+1, msgSiteStr)
				err = input.ValidateMinMaxString(itemCustomerInstallation.Remark, msgStr, 1, 100)
				if err.Error != nil {
					return err
				}
			}

			unq1Str := util2.GenerateConstantaI18n(constanta.UniqueID1, contextModel.AuthAccessTokenModel.Locale, nil)
			msgStr = fmt.Sprintf(`%s Installation %s %d %s`, unq1Str, number, idx2+1, msgSiteStr)
			if util.IsStringEmpty(itemCustomerInstallation.UniqueID1) {
				err = errorModel.GenerateEmptyFieldError(fileName, funcName, msgStr)
				return
			}

			err = input.ValidateMinMaxString(itemCustomerInstallation.UniqueID1, msgStr, 3, 20)
			if err.Error != nil {
				return err
			}

			err = util2.ValidateSpecialCharacterTrimSpace(fileName, funcName, msgStr, itemCustomerInstallation.UniqueID1)
			if err.Error != nil {
				return err
			}

			if !util.IsStringEmpty(itemCustomerInstallation.UniqueID2) {
				unq2Str := util2.GenerateConstantaI18n(constanta.UniqueID2, contextModel.AuthAccessTokenModel.Locale, nil)
				msgStr = fmt.Sprintf(`%s Installation %s %d %s`, unq2Str, number, idx2+1, msgSiteStr)
				err = input.ValidateMinMaxString(itemCustomerInstallation.UniqueID2, msgStr, 1, 20)
				if err.Error != nil {
					return err
				}

				err = util2.ValidateSpecialCharacterTrimSpace(fileName, funcName, msgStr, itemCustomerInstallation.UniqueID2)
				if err.Error != nil {
					return err
				}
			}

			insDateStr := util2.GenerateConstantaI18n(constanta.InstallationDate, contextModel.AuthAccessTokenModel.Locale, nil)
			msgStr = fmt.Sprintf(`%s Installation %s %d %s`, insDateStr, number, idx2+1, msgSiteStr)
			if util.IsStringEmpty(itemCustomerInstallation.InstallationDateStr) {
				return errorModel.GenerateEmptyFieldError(fileName, funcName, msgStr)
			}

			input.CustomerSite[idx].CustomerInstallation[idx2].InstallationDate, err = TimeStrToTimeWithTimeFormat(itemCustomerInstallation.InstallationDateStr, msgStr, constanta.DefaultInstallationTimeFormat)
			if err.Error != nil {
				return err
			}

			if !util.IsStringEmpty(itemCustomerInstallation.UpdatedAtStr) {
				uptStr := util2.GenerateConstantaI18n(constanta.UpdatedAt, contextModel.AuthAccessTokenModel.Locale, nil)
				msgStr = fmt.Sprintf(`%s Installation %s %d %s`, uptStr, number, idx2+1, msgSiteStr)
				input.CustomerSite[idx].CustomerInstallation[idx2].UpdatedAt, err = TimeStrToTime(itemCustomerInstallation.UpdatedAtStr, msgStr)
				if err.Error != nil {
					return err
				}
			}

			if itemCustomerSite.Action == int32(constanta.ActionNoActionCode) {
				actStr := util2.GenerateConstantaI18n(constanta.ActionCode, contextModel.AuthAccessTokenModel.Locale, nil)
				msgStr = fmt.Sprintf(`%s Installation %s %d %s`, actStr, number, idx2+1, msgSiteStr)
				if (itemCustomerInstallation.Action != int32(constanta.ActionInsertCode)) && (itemCustomerInstallation.Action != int32(constanta.ActionDeleteCode)) && (itemCustomerInstallation.Action != int32(constanta.ActionUpdateCode)) {
					return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.ActionCodeRegex2, msgStr, "")
				}
			}
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input CustomerInstallationDataRequest) ValidateViewCustomerInstallationByIDInstallation() (err errorModel.ErrorModel) {
	fileName := "CustomerInstallationDTO.go"
	funcName := "ValidateViewCustomerInstallationByIDInstallation"

	if input.InstallationID < 1 {
		return errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.InstallationID)
	}

	return errorModel.GenerateNonErrorModel()
}
