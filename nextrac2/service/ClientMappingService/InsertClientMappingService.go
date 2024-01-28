package ClientMappingService

import (
	"database/sql"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service/ClientService"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input clientMappingService) InsertNewBranchToClientMapping(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertNewBranchToClientMapping"
		inputStruct in.ClientMappingRequest
	)

	inputStruct, err = input.readBodyAndValidateInsertNewBranch(request, contextModel, input.validateInsertClientMapping)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertNewBranchToClientMapping, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18Message("SUCCESS_INSERT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingService) doInsertNewBranchToClientMapping(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct          = inputStructInterface.(in.ClientMappingRequest)
		modelClientMapping   []repository.ClientMappingModel
		custInstallationOnDB []repository.CustomerInstallationForConfig
		result               interface{}
		finalStruct          interface{}
		dataHeader           repository.ClientMappingModel
		idClientMapping      []int64
	)

	//---------- Check data customer installation
	custInstallationOnDB, err = input.checkDataInCustomerInstallation(inputStruct)
	if err.Error != nil {
		return
	}

	//---------- Check feedback data is empty?
	//if len(inputStruct.CompanyData) == 0 {
	//	detail := util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ClientMappingBundle, "DETAIL_ERROR_DATA_HAS_BEEN_REGISTERED", contextModel.AuthAccessTokenModel.Locale, nil)
	//	err = errorModel.GenerateInvalidRegistrationNewBranch(fileName, funcName, []string{detail})
	//	return
	//}

	//---------- Check and get data to customer
	finalStruct, err = input.checkDataInCustomer(tx, inputStruct, contextModel, timeNow)
	if err.Error != nil {
		fmt.Println("Error checkDataInCustomer : ", finalStruct)
		return
	}

	inputStruct = finalStruct.(in.ClientMappingRequest)
	//---------- Check data client mapping
	result, err = input.checkDataInClientMapping(tx, inputStruct, contextModel)
	if err.Error != nil {
		fmt.Println("Error dataHeader : ", dataHeader)
		return
	}

	dataStructForInsert := result.(in.ClientMappingRequest)
	//if len(inputStruct.CompanyData) == 0 {
	//	err = errorModel.GenerateNonErrorModel()
	//	return
	//}

	if len(dataStructForInsert.CompanyData) != 0 {
		//---------- Get data header client mapping
		dataHeader, err = input.GetDataHeaderInClientMapping(tx, inputStruct)
		if err.Error != nil {
			return
		}

		for _, companyDataElm := range dataStructForInsert.CompanyData {
			for _, branchDataElm := range companyDataElm.BranchData {
				modelClientMapping = append(modelClientMapping, repository.ClientMappingModel{
					ClientID:      sql.NullString{String: dataStructForInsert.ClientID},
					ClientTypeID:  sql.NullInt64{Int64: dataStructForInsert.ClientTypeID},
					CompanyID:     sql.NullString{String: companyDataElm.CompanyID},
					BranchID:      sql.NullString{String: branchDataElm.BranchID},
					ClientAlias:   sql.NullString{String: branchDataElm.ClientAlias},
					SocketID:      sql.NullString{String: dataHeader.SocketID.String},
					CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
					CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
					CreatedAt:     sql.NullTime{Time: timeNow},
					UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
					UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
					UpdatedAt:     sql.NullTime{Time: timeNow},
				})
			}
		}

		idClientMapping, err = dao.ClientMappingDAO.InsertMultipleClientMapping(tx, modelClientMapping)
		if err.Error != nil {
			fmt.Println("Error Insert Client Mapping : ", err.Error.Error())
			return
		}

		for _, idClientMappingElm := range idClientMapping {
			dataAudit = append(dataAudit, repository.AuditSystemModel{
				TableName:  sql.NullString{String: dao.ClientMappingDAO.TableName},
				PrimaryKey: sql.NullInt64{Int64: idClientMappingElm},
			})
		}
	}

	//--- Check must update customer installation
	tempDataAudit, err := ClientService.ClientService.CheckDataMustUpdateInCustomerInstallation(tx, custInstallationOnDB, inputStruct.ClientID, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, tempDataAudit...)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingService) checkDataInCustomerInstallation(inputStruct in.ClientMappingRequest) (output []repository.CustomerInstallationForConfig, err errorModel.ErrorModel) {
	var (
		model []repository.CustomerInstallationDetail
		db    = serverconfig.ServerAttribute.DBConnection
	)

	for _, companyDataElm := range inputStruct.CompanyData {
		for _, branchDataElm := range companyDataElm.BranchData {
			model = append(model, repository.CustomerInstallationDetail{
				UniqueID1: sql.NullString{String: companyDataElm.CompanyID},
				UniqueID2: sql.NullString{String: branchDataElm.BranchID},
			})
		}
	}

	output, err = dao.CustomerInstallationDAO.GetCustomerInstallationByUniqueID(db, model, true)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingService) checkDataInClientMapping(tx *sql.Tx, inputStruct in.ClientMappingRequest, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		companyData              []in.CompanyData
		resultCheckClientMapping interface{}
	)

	for _, companyDataElm := range inputStruct.CompanyData {
		for _, branchDataElm := range companyDataElm.BranchData {
			var branchData []in.BranchData
			branchData = append(branchData, in.BranchData{
				BranchID:    branchDataElm.BranchID,
				ClientAlias: branchDataElm.ClientAlias,
			})
			companyData = append(companyData, in.CompanyData{
				CompanyID:  companyDataElm.CompanyID,
				BranchData: branchData,
			})
		}
	}

	clientRequestStruct := in.ClientRequest{
		ClientTypeID: inputStruct.ClientTypeID,
		CompanyData:  companyData,
	}

	//--- Memisahkan data yang telah ada di client mapping
	resultCheckClientMapping, err = ClientService.ClientService.DoCheckClientMappingSpecialInsertNewBranch(tx, clientRequestStruct, false, contextModel)
	if err.Error != nil {
		return
	}

	if resultCheckClientMapping == nil {
		output = in.ClientMappingRequest{}
		err = errorModel.GenerateNonErrorModel()
		return
	}

	//--- Casting to client request
	var (
		resultStructFromClientMapping = resultCheckClientMapping.(in.ClientRequest)
		companyDataClient             []in.CompanyDataClient
	)

	for _, newCompanyDataElm := range resultStructFromClientMapping.CompanyData {
		for _, newBranchDataElm := range newCompanyDataElm.BranchData {
			var branchDataClient []in.BranchDataClient
			branchDataClient = append(branchDataClient, in.BranchDataClient{
				BranchID:    newBranchDataElm.BranchID,
				ClientAlias: newBranchDataElm.ClientAlias,
			})

			companyDataClient = append(companyDataClient, in.CompanyDataClient{
				CompanyID:  newCompanyDataElm.CompanyID,
				BranchData: branchDataClient,
			})
		}
	}

	clientMappingRequest := in.ClientMappingRequest{
		ClientTypeID: inputStruct.ClientTypeID,
		ClientID:     inputStruct.ClientID,
		CompanyData:  companyDataClient,
	}

	output = clientMappingRequest
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingService) checkDataInCustomer(tx *sql.Tx, inputStruct in.ClientMappingRequest,
	contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, err errorModel.ErrorModel) {

	var (
		fileName            = "InsertClientMappingService.go"
		funcName            = "checkDataInCustomer"
		companyData         []in.CompanyData
		resultCheckCustomer interface{}
	)

	for _, companyDataElm := range inputStruct.CompanyData {
		for _, branchDataElm := range companyDataElm.BranchData {

			var branchData []in.BranchData
			branchData = append(branchData, in.BranchData{
				BranchName: branchDataElm.BranchName,
				BranchID:   branchDataElm.BranchID,
			})

			companyData = append(companyData, in.CompanyData{
				CompanyID:  companyDataElm.CompanyID,
				BranchData: branchData,
			})
		}
	}

	clientRequestStruct := in.ClientRequest{
		ClientTypeID: inputStruct.ClientTypeID,
		CompanyData:  companyData,
	}

	//------ Prepared error when result == nil (not listed in table customer)
	detail := util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ClientMappingBundle, "DETAIL_INVALID_BRANCH_COMPANY_MESSAGE", contextModel.AuthAccessTokenModel.Locale, nil)
	errors := errorModel.GenerateInvalidRegistrationNewBranch(fileName, funcName, []string{detail})

	resultCheckCustomer, err = ClientService.ClientService.DoCheckCustomerExist(tx, clientRequestStruct, errors, contextModel, timeNow, false)
	if err.Error != nil {
		return
	}

	//------ Casting to client request
	var (
		resultStructFromCustomer = resultCheckCustomer.(in.ClientRequest)
		companyDataClient        []in.CompanyDataClient
	)

	for _, newCompanyDataElm := range resultStructFromCustomer.CompanyData {
		for _, newBranchDataElm := range newCompanyDataElm.BranchData {
			var branchDataClient []in.BranchDataClient
			branchDataClient = append(branchDataClient, in.BranchDataClient{
				BranchID:    newBranchDataElm.BranchID,
				ClientAlias: newBranchDataElm.ClientAlias,
			})

			companyDataClient = append(companyDataClient, in.CompanyDataClient{
				CompanyID:  newCompanyDataElm.CompanyID,
				BranchData: branchDataClient,
			})
		}
	}

	clientMappingRequest := in.ClientMappingRequest{
		ClientTypeID: inputStruct.ClientTypeID,
		ClientID:     inputStruct.ClientID,
		CompanyData:  companyDataClient,
	}

	output = clientMappingRequest
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingService) GetDataHeaderInClientMapping(tx *sql.Tx, inputStruct in.ClientMappingRequest) (result repository.ClientMappingModel, err errorModel.ErrorModel) {
	fileName := "InsertClientMappingService.go"
	funcName := "GetDataHeaderInClientMapping"

	result, err = dao.ClientMappingDAO.GetDataHeaderClientMapping(tx, repository.ClientMappingModel{ClientID: sql.NullString{String: inputStruct.ClientID}})
	if err.Error != nil {
		return
	}

	if result.ClientID.String == "" {
		err = errorModel.GenerateDataNotFound(fileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientMappingService) validateInsertClientMapping(inputStruct *in.ClientMappingRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertClientMapping()
}
