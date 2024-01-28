package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"time"
)

type LicenseJournalRequest struct {
	AbstractDTO
	LastSyncStr string `json:"last_sync"`
	LastSync    time.Time
}

func (input *LicenseJournalRequest) ValidateInitiateLicenseJournal() errorModel.ErrorModel {
	var err errorModel.ErrorModel
	if !util.IsStringEmpty(input.LastSyncStr) {
		input.LastSync, err = TimeStrToTimeWithTimeFormat(input.LastSyncStr, "Last Sync", constanta.DefaultTimeFormat)
		if err.Error != nil {
			return err
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input *LicenseJournalRequest) ValidateGetListLicenseJournal() errorModel.ErrorModel {
	var (
		fileName = "LicenseJournalDTO.go"
		funcName = "ValidateGetListLicenseJournal"
	)

	//--- Validate Page
	if input.Page > 0 {
		if input.Limit < 1 {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_MORE_THAN", constanta.Limit, "0")
		}
	}

	//--- Validate Limit
	if input.Limit > 0 {
		if input.Page < 1 {
			return errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, "NEED_MORE_THAN", constanta.Page, "0")
		}
	}

	return input.ValidateInitiateLicenseJournal()
}
