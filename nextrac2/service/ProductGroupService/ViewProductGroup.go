package ProductGroupService

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

func (input productGroupService) ViewProductGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ProductGroupRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err =input.doViewProductGroup(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status =input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input productGroupService) doViewProductGroup(inputStruct in.ProductGroupRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	funcName := "doViewProductGroup"
	productGroupModel := repository.ProductGroupModel{
		ID:        sql.NullInt64{Int64: inputStruct.ID},
	}

	// Get scope
	scope, err := input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	productGroupOnDB, err := dao.ProductGroupDAO.ViewProductGroup(serverconfig.ServerAttribute.DBConnection, productGroupModel, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if productGroupOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, productGroupOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	result = input.reformatModelToDTO(productGroupOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productGroupService) reformatModelToDTO(modelOnDB repository.ProductGroupModel) out.ProductGroupDetailResponse {
	return out.ProductGroupDetailResponse{
		ID:               modelOnDB.ID.Int64,
		ProductGroupName: modelOnDB.ProductGroupName.String,
		CreatedBy:        modelOnDB.CreatedBy.Int64,
		CreatedAt:        modelOnDB.CreatedAt.Time,
		UpdatedBy:        modelOnDB.UpdatedBy.Int64,
		UpdatedAt:        modelOnDB.UpdatedAt.Time,
		UpdatedName:      modelOnDB.UpdatedName.String,
	}
}

func (input productGroupService) validateView(inputStruct *in.ProductGroupRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}
