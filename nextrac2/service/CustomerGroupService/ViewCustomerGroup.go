	package CustomerGroupService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input customerGroupService) ViewCustomerGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.CustomerGroupRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewCustomerGroup)
	if err.Error != nil {
		return
	}

	output.Data.Content, err =input.doViewCustomerGroup(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status =input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input customerGroupService) doViewCustomerGroup(inputStruct in.CustomerGroupRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	funcName := "doViewCustomerGroup"
	customerGroupModel := repository.CustomerGroupModel{
		ID:        sql.NullInt64{Int64: inputStruct.ID},
	}

	// Get scope
	scope, err := input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	customerGroupOnDB, err := dao.CustomerGroupDAO.ViewCustomerGroup(serverconfig.ServerAttribute.DBConnection, customerGroupModel, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if customerGroupOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, customerGroupOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.reformatModelToDTOView(customerGroupOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerGroupService) reformatModelToDTOView(customerGroupModel repository.CustomerGroupModel) out.CustomerGroupDetailResponse {
	return out.CustomerGroupDetailResponse{
		ID:                customerGroupModel.ID.Int64,
		CustomerGroupID:   customerGroupModel.CustomerGroupID.String,
		CustomerGroupName: customerGroupModel.CustomerGroupName.String,
		CreatedBy:         customerGroupModel.CreatedBy.Int64,
		CreatedAt:         customerGroupModel.CreatedAt.Time,
		UpdatedBy:         customerGroupModel.UpdatedBy.Int64,
		UpdatedAt:         customerGroupModel.UpdatedAt.Time,
		UpdatedName:       customerGroupModel.UpdatedName.String,
	}
}

func (input customerGroupService) validateViewCustomerGroup(inputStruct *in.CustomerGroupRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewCustomerGroup()
}