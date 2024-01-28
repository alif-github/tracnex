package ProductService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input productService) InitiateGetListProduct(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam    []in.SearchByParam
		count            int
		scope            map[string]interface{}
		deploymentMethod []in.DeploymentMethod
		newValueSearchBy []string
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListProductValidOperator)
	if err.Error != nil {
		return
	}

	deploymentMethod = append(deploymentMethod, in.DeploymentMethod{
		Code: "O",
		Name: "On Premise",
	}, in.DeploymentMethod{
		Code: "C",
		Name: "Cloud",
	}, in.DeploymentMethod{
		Code: "M",
		Name: "Mobile",
	})

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	count, err = dao.ProductDAO.GetCountProduct(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	for _, valueSearchBy := range input.ValidSearchBy {
		if valueSearchBy != "product_group_id" {
			newValueSearchBy = append(newValueSearchBy, valueSearchBy)
		}
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: newValueSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListProductValidOperator,
		EnumData:      deploymentMethod,
		CountData:     count,
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_GET_LIST_PRODUCT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) GetListProduct(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListProductValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListProduct(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_PRODUCT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) doGetListProduct(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
		scope    map[string]interface{}
		db       = serverconfig.ServerAttribute.DBConnection
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	dbResult, err = dao.ProductDAO.GetListProduct(db, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertToListDTOOut(dbResult)
	return
}

func (input productService) convertToListDTOOut(dbResult []interface{}) (result []out.ListProduct) {
	for _, dbResultItem := range dbResult {
		repo := dbResultItem.(repository.ProductModel)
		result = append(result, out.ListProduct{
			ID:                    repo.ID.Int64,
			ProductName:           repo.ProductName.String,
			ProductDescription:    repo.ProductDescription.String,
			ProductGroupName:      repo.ProductGroupName.String,
			ClientTypeName:        repo.ClientTypeName.String,
			LicenseVariantName:    repo.LicenseVariantName.String,
			LicenseTypeName:       repo.LicenseTypeName.String,
			ProductID:             repo.ProductID.String,
			ClientTypeID:          repo.ClientTypeID.Int64,
			ClientTypeDependantID: repo.ParentClientTypeID.Int64,
			UpdatedAt:             repo.UpdatedAt.Time,
		})
	}

	return result
}
