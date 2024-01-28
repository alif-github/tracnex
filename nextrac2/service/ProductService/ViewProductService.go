package ProductService

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
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input productService) ViewProduct(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ProductRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.DoViewProduct(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_PRODUCT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) DoViewProduct(inputStruct in.ProductRequest, contextModel *applicationModel.ContextModel) (output out.ViewProduct, err errorModel.ErrorModel) {
	fileName := "ViewProductService"
	funcName := "DoViewProduct"
	var productOnDB repository.ProductModel
	var productComponentOnDB []repository.GetListProductComponent

	productModel := repository.ProductModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	}

	productModel.CreatedBy.Int64 = 0
	productOnDB, err = dao.ProductDAO.ViewProduct(serverconfig.ServerAttribute.DBConnection, productModel)
	if err.Error != nil {
		return
	}

	if productOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Product)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, productOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	productComponentOnDB, err = dao.ProductComponentDAO.GetListProductComponentByIDProduct(serverconfig.ServerAttribute.DBConnection, productModel)
	if err.Error != nil {
		return
	}

	output = input.convertProductModelToDTOOut(productOnDB, productComponentOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) convertProductModelToDTOOut(productModel repository.ProductModel, productComponentModel []repository.GetListProductComponent) out.ViewProduct {
	var component []out.ViewProductComponent

	for _, valueComponent := range productComponentModel {
		component = append(component, out.ViewProductComponent{
			ID:             valueComponent.ID.Int64,
			ComponentID:    valueComponent.ComponentID.Int64,
			ComponentName:  valueComponent.ComponentName.String,
			ComponentValue: valueComponent.ComponentValue.String,
		})
	}

	return out.ViewProduct{
		ID:                 productModel.ID.Int64,
		ProductID:          productModel.ProductID.String,
		ProductName:        productModel.ProductName.String,
		ProductDescription: productModel.ProductDescription.String,
		ProductGroupID:     productModel.ProductGroupID.Int64,
		ProductGroupName:   productModel.ProductGroupName.String,
		ClientTypeID:       productModel.ClientTypeID.Int64,
		ClientTypeName:     productModel.ClientTypeName.String,
		IsLicense:          productModel.IsLicense.Bool,
		LicenseVariantID:   productModel.LicenseVariantID.Int64,
		LicenseVariantName: productModel.LicenseVariantName.String,
		LicenseTypeID:      productModel.LicenseTypeID.Int64,
		LicenseTypeName:    productModel.LicenseTypeName.String,
		DeploymentMethod:   productModel.DeploymentMethod.String,
		NoOfUser:           productModel.NoOfUser.Int64,
		IsConcurrentUser:   productModel.IsUserConcurrent.Bool,
		MaxOfflineDays:     productModel.MaxOfflineDays.Int64,
		ModuleId1:          productModel.Module1.Int64,
		ModuleName1:        productModel.ModuleName1.String,
		ModuleId2:          productModel.Module2.Int64,
		ModuleName2:        productModel.ModuleName2.String,
		ModuleId3:          productModel.Module3.Int64,
		ModuleName3:        productModel.ModuleName3.String,
		ModuleId4:          productModel.Module4.Int64,
		ModuleName4:        productModel.ModuleName4.String,
		ModuleId5:          productModel.Module5.Int64,
		ModuleName5:        productModel.ModuleName5.String,
		ModuleId6:          productModel.Module6.Int64,
		ModuleName6:        productModel.ModuleName6.String,
		ModuleId7:          productModel.Module7.Int64,
		ModuleName7:        productModel.ModuleName7.String,
		ModuleId8:          productModel.Module8.Int64,
		ModuleName8:        productModel.ModuleName8.String,
		ModuleId9:          productModel.Module9.Int64,
		ModuleName9:        productModel.ModuleName9.String,
		ModuleId10:         productModel.Module10.Int64,
		ModuleName10:       productModel.ModuleName10.String,
		Component:          component,
		CreatedAt:          productModel.CreatedAt.Time,
		UpdatedAt:          productModel.UpdatedAt.Time,
		UpdatedName:        productModel.UpdatedName.String,
	}
}

func (input productService) validateView(inputStruct *in.ProductRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewProduct()
}
