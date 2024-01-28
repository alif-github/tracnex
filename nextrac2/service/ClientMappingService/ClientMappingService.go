package ClientMappingService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
)

type clientMappingService struct {
	service.RegistrationPrepared
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var ClientMappingService = clientMappingService{}.New()

func (input clientMappingService) New() (output clientMappingService) {
	output.FileName = "ClientMappingService.go"
	output.ValidLimit = []int{10, 20, 50, 100}
	output.ValidOrderBy = []string{"id"}
	output.IdResourceAllowed = []int64{1}
	return
}

func (input clientMappingService) readBodyAndValidateInsertNewBranch(request *http.Request, contextModel *applicationModel.ContextModel, validation func(clientMappingRequest *in.ClientMappingRequest) errorModel.ErrorModel) (inputStruct in.ClientMappingRequest, err errorModel.ErrorModel) {
	var (
		funcName         = "readBodyAndValidateInsertNewBranch"
		stringBody       string
		isAllowed, isND6 bool
		preparedError    errorModel.ErrorModel
	)

	//---------- Read Body Request
	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	//---------- Unmarshal String Body to Main struct
	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	//---------- Init, Important Request Must Exist
	if inputStruct.ClientTypeID == 0 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(input.FileName, funcName, constanta.ClientTypeID)
		return
	}

	//---------- Check to DB, Client Type Exist On Table ?
	preparedError = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientTypeID)
	err = input.CheckIsClientTypeExist(inputStruct.ClientTypeID, preparedError)
	if err.Error != nil {
		return
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])

	for _, companyDataElm := range inputStruct.CompanyData {
		for _, branchDataElm := range companyDataElm.BranchData {
			if branchDataElm.ID == 0 {
				branchDataElm.ID = int64(id)
			}
		}
	}

	//---------- Check Client Type Allowing
	for _, idResourceItem := range input.IdResourceAllowed {
		if idResourceItem == inputStruct.ClientTypeID {
			isAllowed = true
			break
		}
	}

	//---------- Is Not Allowed, Then Forbidden to Access
	if !isAllowed {
		err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
		return
	}

	//---------- Login Must Be Same With Parent Client ID
	isND6 = inputStruct.ClientTypeID == constanta.ResourceND6ID
	if isND6 {
		if util.IsStringEmpty(inputStruct.ClientID) {
			err = errorModel.GenerateEmptyFieldError(input.FileName, funcName, constanta.ParentClientID)
			return
		}

		validationResult, errField, _ := util2.IsClientIDValid(inputStruct.ClientID)
		messageTemp := util2.GenerateConstantaI18n(errField, contextModel.AuthAccessTokenModel.Locale, nil)
		if !validationResult {
			err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, messageTemp, constanta.ClientID, "")
			return
		}

		if contextModel.AuthAccessTokenModel.ClientID != inputStruct.ClientID {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientID)
			return
		}
	}

	err = validation(&inputStruct)
	return
}
