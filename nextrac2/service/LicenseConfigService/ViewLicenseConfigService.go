package LicenseConfigService

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

func (input licenseConfigService) ViewDetailLicenseConfigService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.LicenseConfigRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewLicenseConfig)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewLicenseConfig(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_LICENSE_CONFIG_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) doViewLicenseConfig(inputStruct in.LicenseConfigRequest, contextModel *applicationModel.ContextModel) (output out.ViewDetailLicenseConfig, err errorModel.ErrorModel) {
	var (
		fileName                = "ViewLicenseConfigService.go"
		funcName                = "doViewLicenseConfig"
		resultOnDBLicenseConfig repository.LicenseConfigModel
		componentList           []repository.ProductComponentModel
		scope                   map[string]interface{}
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	resultOnDBLicenseConfig, err = dao.LicenseConfigDAO.ViewDetailLicenseConfig(serverconfig.ServerAttribute.DBConnection, repository.LicenseConfigModel{ID: sql.NullInt64{Int64: inputStruct.ID}, CreatedBy: sql.NullInt64{Int64: 0}}, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if resultOnDBLicenseConfig.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, resultOnDBLicenseConfig.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	componentList, err = dao.ProductComponentDAO.GetProductComponentByIDProduct(serverconfig.ServerAttribute.DBConnection, resultOnDBLicenseConfig.ProductID.Int64)
	if err.Error != nil {
		return
	}

	output = input.reformatToDTO(resultOnDBLicenseConfig, componentList)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) validateViewLicenseConfig(inputStruct *in.LicenseConfigRequest) errorModel.ErrorModel {
	return inputStruct.ValidationViewLicenseConfig()
}

func (input licenseConfigService) reformatToDTO(resultOnDBLicenseConfig repository.LicenseConfigModel, componentList []repository.ProductComponentModel) (result out.ViewDetailLicenseConfig) {
	var listProductComponent []out.ListProductComponent

	resultOnDBLicenseConfig.Component = componentList
	for _, itemComponent := range resultOnDBLicenseConfig.Component {
		listProductComponent = append(listProductComponent, out.ListProductComponent{
			ComponentName:  itemComponent.ComponentName.String,
			ComponentValue: itemComponent.ComponentValue.String,
		})
	}

	result.ID = resultOnDBLicenseConfig.ID.Int64
	result.InstallationID = resultOnDBLicenseConfig.InstallationID.Int64
	result.ParentCustomerID = resultOnDBLicenseConfig.ParentCustomerID.Int64
	result.ParentCustomer = resultOnDBLicenseConfig.ParentCustomer.String
	result.CustomerID = resultOnDBLicenseConfig.CustomerID.Int64
	result.Customer = resultOnDBLicenseConfig.Customer.String
	result.SiteID = resultOnDBLicenseConfig.SiteID.Int64
	result.ProductID = resultOnDBLicenseConfig.ProductID.Int64
	result.ProductName = resultOnDBLicenseConfig.ProductName.String
	result.ClientID = resultOnDBLicenseConfig.ClientID.String
	result.LicenseVariantName = resultOnDBLicenseConfig.LicenseVariantName.String
	result.LicenseTypeName = resultOnDBLicenseConfig.LicenseTypeName.String
	result.DeploymentMethod = resultOnDBLicenseConfig.DeploymentMethod.String
	result.NoOfUser = resultOnDBLicenseConfig.NoOfUser.Int64
	result.IsUserConcurrent = resultOnDBLicenseConfig.IsUserConcurrent.String
	result.UniqueID1 = resultOnDBLicenseConfig.UniqueID1.String
	result.UniqueID2 = resultOnDBLicenseConfig.UniqueID2.String
	result.ProductValidFrom = resultOnDBLicenseConfig.ProductValidFrom.Time.Format(constanta.DefaultInstallationTimeFormat)
	result.ProductValidThru = resultOnDBLicenseConfig.ProductValidThru.Time.Format(constanta.DefaultInstallationTimeFormat)
	result.MaxOfflineDays = resultOnDBLicenseConfig.MaxOfflineDays.Int64
	result.ClientTypeName = resultOnDBLicenseConfig.ClientType.String
	result.AllowActivation = resultOnDBLicenseConfig.AllowActivation.String
	result.ModuleID1 = resultOnDBLicenseConfig.ModuleIDName1.String
	result.ModuleID2 = resultOnDBLicenseConfig.ModuleIDName2.String
	result.ModuleID3 = resultOnDBLicenseConfig.ModuleIDName3.String
	result.ModuleID4 = resultOnDBLicenseConfig.ModuleIDName4.String
	result.ModuleID5 = resultOnDBLicenseConfig.ModuleIDName5.String
	result.ModuleID6 = resultOnDBLicenseConfig.ModuleIDName6.String
	result.ModuleID7 = resultOnDBLicenseConfig.ModuleIDName7.String
	result.ModuleID8 = resultOnDBLicenseConfig.ModuleIDName8.String
	result.ModuleID9 = resultOnDBLicenseConfig.ModuleIDName9.String
	result.ModuleID10 = resultOnDBLicenseConfig.ModuleIDName10.String
	result.CreatedAt = resultOnDBLicenseConfig.CreatedAt.Time
	result.UpdatedAt = resultOnDBLicenseConfig.UpdatedAt.Time
	result.UpdatedName = resultOnDBLicenseConfig.UpdatedName.String
	result.ProductComponent = listProductComponent

	return
}

func (input licenseConfigService) validateMultipleDataScopeLicenseConfig(contextModel *applicationModel.ContextModel) (getScope map[string]interface{}, err errorModel.ErrorModel) {
	getScope, err = input.ValidateMultipleDataScope(contextModel, []string{
		constanta.ProvinceDataScope,
		constanta.DistrictDataScope,
		constanta.SalesmanDataScope,
		constanta.CustomerGroupDataScope,
		constanta.CustomerCategoryDataScope,
		constanta.ProductGroupDataScope,
		constanta.ClientTypeDataScope,
	})

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) restructureDataScopeLicenseConfig(scope map[string]interface{}) (newScope map[string]interface{}) {
	var (
		parent = "_parent"
		child  = "_child"
	)

	if len(scope) > 0 {
		newScope = make(map[string]interface{})
		newScope[constanta.ProvinceDataScope+parent] = scope[constanta.ProvinceDataScope]
		newScope[constanta.DistrictDataScope+parent] = scope[constanta.DistrictDataScope]
		newScope[constanta.SalesmanDataScope+parent] = scope[constanta.SalesmanDataScope]
		newScope[constanta.CustomerGroupDataScope+parent] = scope[constanta.CustomerGroupDataScope]
		newScope[constanta.CustomerCategoryDataScope+parent] = scope[constanta.CustomerCategoryDataScope]
		newScope[constanta.ProvinceDataScope+child] = scope[constanta.ProvinceDataScope]
		newScope[constanta.DistrictDataScope+child] = scope[constanta.DistrictDataScope]
		newScope[constanta.SalesmanDataScope+child] = scope[constanta.SalesmanDataScope]
		newScope[constanta.CustomerGroupDataScope+child] = scope[constanta.CustomerGroupDataScope]
		newScope[constanta.CustomerCategoryDataScope+child] = scope[constanta.CustomerCategoryDataScope]
		newScope[constanta.ProductGroupDataScope] = scope[constanta.ProductGroupDataScope]
		newScope[constanta.ClientTypeDataScope] = scope[constanta.ClientTypeDataScope]
	}

	return
}
