package CustomerCategoryService

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

func (input customerCategoryService) ViewCustomerCategory(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.CustomerCategoryRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err =input.doViewCustomerCategory(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status =input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input customerCategoryService) doViewCustomerCategory(inputStruct in.CustomerCategoryRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	funcName := "doViewCustomerCategory"
	customerCatgModel := repository.CustomerCategoryModel{
		ID:        sql.NullInt64{Int64: inputStruct.ID},
	}

	// Get scope
	scope, err := input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	customerCatgOnDB, err := dao.CustomerCategoryDAO.ViewCustomerCategory(serverconfig.ServerAttribute.DBConnection, customerCatgModel, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if customerCatgOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, customerCatgOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.reformatModelToDTOView(customerCatgOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerCategoryService) reformatModelToDTOView(customerCatgModel repository.CustomerCategoryModel) out.CustomerCategoryDetailResponse {
	return out.CustomerCategoryDetailResponse{
		ID:                   customerCatgModel.ID.Int64,
		CustomerCategoryID:   customerCatgModel.CustomerCategoryID.String,
		CustomerCategoryName: customerCatgModel.CustomerCategoryName.String,
		CreatedBy:            customerCatgModel.CreatedBy.Int64,
		CreatedAt:            customerCatgModel.CreatedAt.Time,
		UpdatedBy:            customerCatgModel.UpdatedBy.Int64,
		UpdatedAt:            customerCatgModel.UpdatedAt.Time,
		UpdatedName:          customerCatgModel.UpdatedName.String,
	}
}

func (input customerCategoryService) validateView(inputStruct *in.CustomerCategoryRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}