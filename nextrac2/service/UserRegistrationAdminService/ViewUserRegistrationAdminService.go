package UserRegistrationAdminService

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"strconv"
)

func (input userRegistrationAdminService) ViewUserRegistrationAdmin(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.UserRegistrationAdminRequest

	if inputStruct, err = input.readBodyAndValidateForViewUserRegistrationAdmin(request, input.validateStructForViewUserRegistrationAdmin); err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewUserRegistrationAdmin(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input userRegistrationAdminService) readBodyAndValidateForViewUserRegistrationAdmin(request *http.Request, validation func(input *in.UserRegistrationAdminRequest) errorModel.ErrorModel) (inputStruct in.UserRegistrationAdminRequest, err errorModel.ErrorModel) {
	id, _ := strconv.Atoi(mux.Vars(request)["id"])

	inputStruct.ID = int64(id)

	err = validation(&inputStruct)

	return
}

func (input userRegistrationAdminService) validateStructForViewUserRegistrationAdmin(inputStruct *in.UserRegistrationAdminRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewUserRegistrationAdmin()
}

func (input userRegistrationAdminService) doViewUserRegistrationAdmin(inputStruct in.UserRegistrationAdminRequest, contextModel *applicationModel.ContextModel) (output out.DetailUserRegistrationAdminResponse, err errorModel.ErrorModel) {
	var (
		fileName                    = "ViewUserRegistrationAdminService.go"
		funcName                    = "doViewUserRegistrationAdmin"
		detailUserRegistrationAdmin repository.UserRegistrationAdminModel
		scope                       map[string]interface{}
		mappingDB                   = make(map[string]applicationModel.MappingScopeDB)
	)

	mappingDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "ur.client_type_id",
		Count: "ur.client_type_id",
	}

	scope, err = input.validateDataScopeUserRegistrationAdmin(contextModel)
	if err.Error != nil {
		return
	}

	if detailUserRegistrationAdmin, err = dao.UserRegistrationAdminDAO.ViewDetailUserRegistrationAdmin(serverconfig.ServerAttribute.DBConnection, repository.UserRegistrationAdminModel{ID: sql.NullInt64{Int64: inputStruct.ID}}, contextModel.LimitedByCreatedBy, scope, mappingDB); err.Error != nil {
		return
	}

	if detailUserRegistrationAdmin.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.UserRegistrationID)
		return
	}

	output = out.DetailUserRegistrationAdminResponse{
		ID:                 detailUserRegistrationAdmin.ID.Int64,
		ParentCustomerId:   detailUserRegistrationAdmin.ParentCustomerId.Int64,
		ParentCustomerName: detailUserRegistrationAdmin.ParentCustomerName.String,
		CustomerId:         detailUserRegistrationAdmin.CustomerId.Int64,
		SiteId:             detailUserRegistrationAdmin.SiteId.Int64,
		CustomerName:       detailUserRegistrationAdmin.CustomerName.String,
		CompanyId:          detailUserRegistrationAdmin.UniqueID1.String,
		BranchId:           detailUserRegistrationAdmin.UniqueID2.String,
		CompanyName:        detailUserRegistrationAdmin.CompanyName.String,
		BranchName:         detailUserRegistrationAdmin.BranchName.String,
		UserAdmin:          detailUserRegistrationAdmin.UserAdmin.String,
		PasswordAdmin:      detailUserRegistrationAdmin.PasswordAdmin.String,
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
