package NexmileParameterService

import (
	"database/sql"
	"fmt"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input nexmileParameterService) ViewNexmileParameter(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.NexmileParameterRequestForView

	fmt.Println("ViewNexmileParameter start")

	inputStruct, err = input.readBodyAndValidateForView(request, contextModel, input.validateNexmileParameterForView)
	if err.Error != nil {
		return
	}

	fmt.Println("ViewNexmileParameter read body")

	output.Data.Content, err = input.doViewNexmileParameter(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	fmt.Println("ViewNexmileParameter end")

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)

	return
}

func (input nexmileParameterService) validateNexmileParameterForView(inputStruct *in.NexmileParameterRequestForView) errorModel.ErrorModel {
	return inputStruct.ValidateViewNexmileParameter()
}

func (input nexmileParameterService) doViewNexmileParameter(inputStruct in.NexmileParameterRequestForView, contextModel *applicationModel.ContextModel) (output out.ViewNexmileParameterResponse, err errorModel.ErrorModel) {
	fileName := "ViewNexmileParameterService.go"
	funcName := "doViewNexmileParameter"

	var (
		nexmileParameterModel      repository.NexmileParameterModel
		parameterValueModelOnDB    []repository.ParameterValueModel
		userRegistrationDetailOnDB repository.UserRegistrationDetailModel
		userAuth                   authentication_response.UserAuthenticationResponse
		userRegistrationOnDB       repository.UserRegistrationAdminModel
		pkceClientMappingOnDB      repository.PKCEClientMappingModel
		clientMappingOnDB          repository.ClientMappingModel
		authUserID                 int64
	)

	// [Start] Step 1 - Get User Registration Detail
	userRegistrationDetailOnDB, err = dao.UserRegistrationDetailDAO.GetUserRegistrationDetailWithUserIDAndPassword(serverconfig.ServerAttribute.DBConnection, repository.UserRegistrationDetailModel{
		UserID:       sql.NullString{String: inputStruct.UserId},
		Password:     sql.NullString{String: inputStruct.Password},
		AndroidID:    sql.NullString{String: inputStruct.AndroidId},
		ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeId},
		AuthUserID:   sql.NullInt64{Int64: inputStruct.AuthUserID},
	})

	if err.Error != nil {
		return
	}

	if userRegistrationDetailOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.UserRegistrationDetailName)
		return
	}

	authUserID = userRegistrationDetailOnDB.AuthUserID.Int64
	// [END] Step 1 - Get User Registration Detail

	// Convert DTO to Model
	nexmileParameterModel = input.convertDTOToModel(inputStruct)

	// [Start] Step 2 - Parent Client Validation
	if util.IsStringEmpty(inputStruct.ClientID) {
		nexmileParameterModel.ClientID = userRegistrationDetailOnDB.ClientID
	}

	pkceClientMappingOnDB, clientMappingOnDB, err = dao.PKCEClientMappingDAO.GetPKCEClientWithClientIDAndUniqueID(serverconfig.ServerAttribute.DBConnection, repository.PKCEClientMappingModel{
		ClientID:     nexmileParameterModel.ClientID,
		CompanyID:    userRegistrationDetailOnDB.UniqueID1,
		BranchID:     userRegistrationDetailOnDB.UniqueID2,
		ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeId},
	})
	if err.Error != nil {
		return
	}

	fmt.Println(pkceClientMappingOnDB)
	fmt.Println(clientMappingOnDB)

	if clientMappingOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateClientValidationError(fileName, funcName)
		return
	}

	// ambil dari pkce_client_mapping -> parent_client_id = client_id yg login
	if pkceClientMappingOnDB.ClientID.String != contextModel.AuthAccessTokenModel.ClientID {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientMappingClientID)
		return
	}
	// [End] Step 2 - Parent Client Validation

	// [Start] Step 2.1 - User Validation (Auth)
	if inputStruct.AuthUserID > 0 {
		authUserID = inputStruct.AuthUserID
	}

	userAuth, err = resource_common_service.InternalGetUserByID(authUserID, contextModel)
	if err.Error != nil {
		return
	}
	if userAuth.Nexsoft.Payload.Data.Content.ClientID != nexmileParameterModel.ClientID.String {
		err = errorModel.GenerateUnknownAuthUserId(fileName, funcName)
		return
	}

	nexmileParameterModel.UniqueID1 = userRegistrationDetailOnDB.UniqueID1
	nexmileParameterModel.UniqueID2 = userRegistrationDetailOnDB.UniqueID2
	nexmileParameterModel.ClientID = clientMappingOnDB.ClientID
	nexmileParameterModel.ClientTypeID = clientMappingOnDB.ClientTypeID
	// [End] Step 2.1 - User Validation (Auth)

	// Step 3 - Pengambilan Data Nexmile Parameter
	parameterValueModelOnDB, err = dao.NexmileParameterDAO.GetFieldNexmileParameter(serverconfig.ServerAttribute.DBConnection, nexmileParameterModel)
	if err.Error != nil {
		return
	}

	// Step 4 - Pengambilan Data User Registrasi
	userRegistrationOnDB, err = dao.UserRegistrationAdminDAO.GetFieldForNexmileParameter(serverconfig.ServerAttribute.DBConnection, repository.UserRegistrationAdminModel{
		UniqueID1:       nexmileParameterModel.UniqueID1,
		UniqueID2:       nexmileParameterModel.UniqueID2,
		ClientTypeID:    pkceClientMappingOnDB.ClientTypeID,
		ClientMappingID: clientMappingOnDB.ID,
	})
	if err.Error != nil {
		return
	}

	fmt.Println(userRegistrationOnDB)

	if userRegistrationOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.UserRegistrationAdmin)
		return
	}

	// Convert To Response
	output = input.convertToResponseGetNexmileParameter(userRegistrationOnDB, userRegistrationDetailOnDB, parameterValueModelOnDB)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexmileParameterService) convertDTOToModel(inputStruct in.NexmileParameterRequestForView) (result repository.NexmileParameterModel) {
	return repository.NexmileParameterModel{
		ClientID:     sql.NullString{String: inputStruct.ClientID, Valid: !util.IsStringEmpty(inputStruct.ClientID)},
		AuthUserID:   sql.NullInt64{Int64: inputStruct.AuthUserID, Valid: !util2.IsFieldNumericEmpty(inputStruct.AuthUserID)},
		UserID:       sql.NullString{String: inputStruct.UserId, Valid: !util.IsStringEmpty(inputStruct.UserId)},
		AndroidID:    sql.NullString{String: inputStruct.AndroidId, Valid: !util.IsStringEmpty(inputStruct.AndroidId)},
		ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeId, Valid: !util2.IsFieldNumericEmpty(inputStruct.ClientTypeId)},
		Password:     sql.NullString{String: inputStruct.Password, Valid: !util.IsStringEmpty(inputStruct.Password)},
	}
}

func (input nexmileParameterService) convertToResponseGetNexmileParameter(userRegistration repository.UserRegistrationAdminModel, userRegistrationDetail repository.UserRegistrationDetailModel, parameterValueOnDB []repository.ParameterValueModel) out.ViewNexmileParameterResponse {
	var parameterValueOnDBResponse []out.ParameterValueResponse

	for _, itemParameterValue := range parameterValueOnDB {
		parameterValueOnDBResponse = append(parameterValueOnDBResponse, out.ParameterValueResponse{
			ParameterID:    itemParameterValue.ParameterID.String,
			ParameterValue: itemParameterValue.ParameterValue.String,
		})
	}

	return out.ViewNexmileParameterResponse{
		UniqueId1:         userRegistration.UniqueID1.String,
		CompanyName:       userRegistration.CompanyName.String,
		ProductValidThru:  userRegistrationDetail.ProductValidThru.Time,
		ProductValidaFrom: userRegistrationDetail.ProductValidFrom.Time,
		LicenseStatus:     userRegistrationDetail.LicenseStatus.Int64,
		PasswordAdmin:     userRegistration.PasswordAdmin.String,
		UniqueId2:         userRegistration.UniqueID2.String,
		UserAdmin:         userRegistration.UserAdmin.String,
		MaxOfflineDays:    userRegistrationDetail.MaxOfflineDays.Int64,
		Parameters:        parameterValueOnDBResponse,
	}
}
