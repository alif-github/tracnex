package UserRegistrationAdminService

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"time"
)

type userRegistrationAdminService struct {
	service.AbstractService
	service.GetListData
}

var UserRegistrationAdminService = userRegistrationAdminService{}.New()

func (input userRegistrationAdminService) New() (output userRegistrationAdminService) {
	output.FileName = "UserRegistrationAdminService.go"
	output.ServiceName = "USER_REGISTRATION_ADMIN"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"customer_name",
		"parent_customer_name",
		"company_id",
		"branch_id",
		"company_name",
		"branch_name",
		"user_admin",
		"password_admin",
	}
	output.ValidSearchBy = []string{"company_name", "id"}
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "user_registration.client_type_id",
		Count: "user_registration.client_type_id",
	}
	return
}

func (input userRegistrationAdminService) readBodyAndValidateForView(request *http.Request, validation func(input *in.LicenseVariantRequest) errorModel.ErrorModel) (inputStruct in.LicenseVariantRequest, err errorModel.ErrorModel) {

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func (input userRegistrationAdminService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.UserRegistrationAdminRequest) errorModel.ErrorModel) (inputStruct in.UserRegistrationAdminRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidate"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	err = validation(&inputStruct)
	if err.Error != nil {
		return
	}

	if contextModel.AuthAccessTokenModel.ClientID != inputStruct.ClientID {
		err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
		return
	}

	return
}

func (input userRegistrationAdminService) convertDTOToModel(inputStruct in.UserRegistrationAdminRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) repository.UserRegistrationAdminModel {
	return repository.UserRegistrationAdminModel{
		UniqueID1:     sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:     sql.NullString{String: inputStruct.UniqueID2},
		UserAdmin:     sql.NullString{String: inputStruct.UserAdmin},
		PasswordAdmin: sql.NullString{String: inputStruct.PasswordAdmin},
		CompanyName:   sql.NullString{String: inputStruct.CompanyName},
		BranchName:    sql.NullString{String: inputStruct.BranchName},
		ClientID:      sql.NullString{String: inputStruct.ClientID},
		ClientTypeID:  sql.NullInt64{Int64: inputStruct.ClientTypeID},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}
}

func (input userRegistrationAdminService) validateDataScopeUserRegistrationAdmin(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	funcName := "validateDataScopeUserRegistrationAdmin"

	output = service.ValidateScope(contextModel, []string{
		constanta.ClientTypeDataScope,
	})

	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
