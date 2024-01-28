package UserRegistrationAdminService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input userRegistrationAdminService) InitiateGetListUserRegistrationAdmin(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		countData     interface{}
		searchByParam []in.SearchByParam
		scope         map[string]interface{}
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListUserRegistrationAdminValidOperator)
	if err.Error != nil {
		return
	}

	scope, err = input.validateDataScopeUserRegistrationAdmin(contextModel)
	if err.Error != nil {
		return
	}

	countData, err = dao.UserRegistrationAdminDAO.GetCountUserRegistrationAdmin(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListUserRegistrationAdminValidOperator,
		CountData:     countData.(int),
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userRegistrationAdminService) GetListUserRegistrationAdmin(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListUserRegistrationAdminValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListUserRegistrationAdmin(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userRegistrationAdminService) doGetListUserRegistrationAdmin(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
		scope    map[string]interface{}
	)

	scope, err = input.validateDataScopeUserRegistrationAdmin(contextModel)
	if err.Error != nil {
		return
	}

	dbResult, err = dao.UserRegistrationAdminDAO.GetListUserRegistrationAdmin(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertToListDTOOut(dbResult)
	return
}

func (input userRegistrationAdminService) convertToListDTOOut(dbResult []interface{}) (result []out.ListUserRegistrationAdminResponse) {
	for _, dbResultItem := range dbResult {
		repo := dbResultItem.(repository.UserRegistrationAdminModel)
		result = append(result, out.ListUserRegistrationAdminResponse{
			ID:                 repo.ID.Int64,
			CustomerName:       repo.CustomerName.String,
			ParentCustomerName: repo.ParentCustomerName.String,
			CompanyID:          repo.UniqueID1.String,
			BranchID:           repo.UniqueID2.String,
			CompanyName:        repo.CompanyName.String,
			BranchName:         repo.BranchName.String,
			UserAdmin:          repo.UserAdmin.String,
			PasswordAdmin:      repo.PasswordAdmin.String,
		})
	}

	return result
}
