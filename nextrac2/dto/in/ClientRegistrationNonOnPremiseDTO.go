package in

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
)

type ClientRegistrationNonOnPremiseRequest struct {
	ClientID     string           `json:"client_id"`
	ClientTypeID int64            `json:"client_type_id"`
	DetailClient []UniqueIDClient `json:"detail_client"`
}

type UniqueIDClient struct {
	UniqueID1 string `json:"unique_id_1"`
	UniqueID2 string `json:"unique_id_2"`
}

func (input ClientRegistrationNonOnPremiseRequest) ValidateInsertClientRegistNonOnPremise() (err errorModel.ErrorModel) {
	var (
		fileName = "ClientRegistrationNonOnPremiseDTO.go"
		funcName = "ValidateInsertClientRegistrationNonOnPremise"
	)

	if util.IsStringEmpty(input.ClientID) {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.ClientID)
		return
	}

	err = util2.ValidateMinMaxString(input.ClientID, constanta.ClientID, 1, 256)
	if err.Error != nil {
		return
	}

	if input.ClientTypeID < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.ClientTypeID)
		return
	}

	if len(input.DetailClient) < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(fileName, funcName, constanta.DetailClient)
		return
	}

	for _, itemUnique := range input.DetailClient {

		//--- Unique ID 1 is Mandatory
		if util.IsStringEmpty(itemUnique.UniqueID1) {
			err = errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.UniqueID1)
			return
		}

		err = util2.ValidateMinMaxString(itemUnique.UniqueID1, constanta.UniqueID1, 1, 20)
		if err.Error != nil {
			return
		}

		//--- Unique ID 2 is Optional
		if !util.IsStringEmpty(itemUnique.UniqueID2) {
			err = util2.ValidateMinMaxString(itemUnique.UniqueID2, constanta.UniqueID2, 1, 20)
			if err.Error != nil {
				return
			}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
