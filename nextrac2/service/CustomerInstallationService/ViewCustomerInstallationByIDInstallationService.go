package CustomerInstallationService

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

func (input customerInstallationService) ViewCustomerSiteInstallationByInstallationIDService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.CustomerInstallationDataRequest

	inputStruct, err = input.readBodyAndValidateForDetailInstallation(request, input.validateViewCustomerSiteInstallation)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.DoViewDetailByInstallationID(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_CUSTOMER_INSTALLATION_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) DoViewDetailByInstallationID(inputStruct in.CustomerInstallationDataRequest, contextModel *applicationModel.ContextModel) (output out.CustomerInstallationDetailListForConfig, err errorModel.ErrorModel) {
	var (
		fileName                 = "ViewCustomerInstallationByIDInstallationService.go"
		funcName                 = "DoViewDetailByInstallationID"
		componentList            []repository.ProductComponentModel
		resultDetailInstallation repository.CustomerInstallationDetailConfig
		listProductComponent     []out.ListProductComponent
		resultClientMapping      repository.ClientMappingModel
		customerStr              = util2.GenerateConstantaI18n(constanta.Customer, contextModel.AuthAccessTokenModel.Locale, nil)
		installationIDStr        = util2.GenerateConstantaI18n(constanta.InstallationID, contextModel.AuthAccessTokenModel.Locale, nil)
		db                       = serverconfig.ServerAttribute.DBConnection
		scope                    map[string]interface{}
		mappingScopeDB           map[string]applicationModel.MappingScopeDB
	)

	mappingScopeDB = input.newMappingScopeDBForViewInstallation()

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	resultDetailInstallation, err = dao.CustomerInstallationDAO.ViewDetailInstallationByInstallationID(db, repository.CustomerInstallationDetailConfig{InstallationID: sql.NullInt64{Int64: inputStruct.InstallationID}}, scope, mappingScopeDB)
	if err.Error != nil {
		return
	}

	if resultDetailInstallation.InstallationID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Installation)
		return
	}

	if resultDetailInstallation.ClientMappingID.Int64 < 1 {
		err = errorModel.GenerateClientIDNotFound(fileName, funcName, constanta.Installation, customerStr, installationIDStr)
		return
	}

	resultClientMapping, err = dao.ClientMappingDAO.GetClientIDWithID(db, repository.ClientMappingModel{ID: sql.NullInt64{Int64: resultDetailInstallation.ClientMappingID.Int64}})
	if err.Error != nil {
		return
	}

	if resultClientMapping.ID.Int64 < 1 {
		err = errorModel.GenerateClientIDNotFound(fileName, funcName, constanta.Installation, customerStr, installationIDStr)
		return
	}

	resultDetailInstallation.ClientID.String = resultClientMapping.ClientID.String
	componentList, err = dao.ProductComponentDAO.GetProductComponentByIDProduct(db, resultDetailInstallation.ProductID.Int64)
	if err.Error != nil {
		return
	}

	resultDetailInstallation.Component = componentList
	for _, itemComponent := range resultDetailInstallation.Component {
		listProductComponent = append(listProductComponent, out.ListProductComponent{
			ComponentName:  itemComponent.ComponentName.String,
			ComponentValue: itemComponent.ComponentValue.String,
		})
	}

	result := out.CustomerInstallationDetailListForConfig{
		InstallationID:     resultDetailInstallation.InstallationID.Int64,
		ParentCustomerID:   resultDetailInstallation.ParentCustomerID.Int64,
		ParentCustomer:     resultDetailInstallation.ParentCustomer.String,
		CustomerID:         resultDetailInstallation.CustomerID.Int64,
		Customer:           resultDetailInstallation.Customer.String,
		SiteID:             resultDetailInstallation.SiteID.Int64,
		ProductID:          resultDetailInstallation.ProductID.Int64,
		ProductName:        resultDetailInstallation.ProductName.String,
		ClientID:           resultDetailInstallation.ClientID.String,
		LicenseVariantName: resultDetailInstallation.LicenseVariantName.String,
		LicenseTypeName:    resultDetailInstallation.LicenseTypeName.String,
		DeploymentMethod:   resultDetailInstallation.DeploymentMethod.String,
		NoOfUser:           resultDetailInstallation.NoOfUser.Int64,
		IsUserConcurrent:   resultDetailInstallation.IsUserConcurrent.Bool,
		UniqueID1:          resultDetailInstallation.UniqueID1.String,
		UniqueID2:          resultDetailInstallation.UniqueID2.String,
		ProductValidFrom:   resultDetailInstallation.ProductValidFrom.Time,
		ProductValidThru:   resultDetailInstallation.ProductValidThru.Time,
		MaxOfflineDays:     resultDetailInstallation.MaxOfflineDays.Int64,
		ClientTypeName:     resultDetailInstallation.ClientType.String,
		ModuleID1:          resultDetailInstallation.ModuleIDName1.String,
		ModuleID2:          resultDetailInstallation.ModuleIDName2.String,
		ModuleID3:          resultDetailInstallation.ModuleIDName3.String,
		ModuleID4:          resultDetailInstallation.ModuleIDName4.String,
		ModuleID5:          resultDetailInstallation.ModuleIDName5.String,
		ModuleID6:          resultDetailInstallation.ModuleIDName6.String,
		ModuleID7:          resultDetailInstallation.ModuleIDName7.String,
		ModuleID8:          resultDetailInstallation.ModuleIDName8.String,
		ModuleID9:          resultDetailInstallation.ModuleIDName9.String,
		ModuleID10:         resultDetailInstallation.ModuleIDName10.String,
		ProductComponent:   listProductComponent,
	}

	output = result
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) validateViewCustomerSiteInstallation(inputStruct *in.CustomerInstallationDataRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewCustomerInstallationByIDInstallation()
}

func (input customerInstallationService) newMappingScopeDBForViewInstallation() (mappingScopeDB map[string]applicationModel.MappingScopeDB) {
	var (
		parent = "_parent"
		child  = "_child"
	)

	mappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	mappingScopeDB[constanta.ProvinceDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cup.province_id",
		Count: "cup.province_id",
	}
	mappingScopeDB[constanta.DistrictDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cup.district_id",
		Count: "cup.district_id",
	}
	mappingScopeDB[constanta.SalesmanDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cup.salesman_id",
		Count: "cup.salesman_id",
	}
	mappingScopeDB[constanta.CustomerGroupDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cup.customer_group_id",
		Count: "cup.customer_group_id",
	}
	mappingScopeDB[constanta.CustomerCategoryDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cup.customer_category_id",
		Count: "cup.customer_category_id",
	}
	mappingScopeDB[constanta.ProvinceDataScope+child] = applicationModel.MappingScopeDB{
		View:  "cuc.province_id",
		Count: "cuc.province_id",
	}
	mappingScopeDB[constanta.DistrictDataScope+child] = applicationModel.MappingScopeDB{
		View:  "cuc.district_id",
		Count: "cuc.district_id",
	}
	mappingScopeDB[constanta.SalesmanDataScope+child] = applicationModel.MappingScopeDB{
		View:  "cuc.salesman_id",
		Count: "cuc.salesman_id",
	}
	mappingScopeDB[constanta.CustomerGroupDataScope+child] = applicationModel.MappingScopeDB{
		View:  "cuc.customer_group_id",
		Count: "cuc.customer_group_id",
	}
	mappingScopeDB[constanta.CustomerCategoryDataScope+child] = applicationModel.MappingScopeDB{
		View:  "cuc.customer_category_id",
		Count: "cuc.customer_category_id",
	}
	mappingScopeDB[constanta.ProductGroupDataScope] = applicationModel.MappingScopeDB{
		View:  "pr.product_group_id",
		Count: "pr.product_group_id",
	}
	mappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "ct.id",
		Count: "ct.id",
	}
	return
}
