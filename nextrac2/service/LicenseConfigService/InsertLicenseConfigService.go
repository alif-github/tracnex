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
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input licenseConfigService) InsertLicenseConfig(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertLicenseConfig"
		inputStruct in.LicenseConfigRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAuditCustom(funcName, inputStruct, contextModel, input.DoInsertLicenseConfig, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- Additional function
	})

	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_LICENSE_CONFIG", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) DoInsertLicenseConfig(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, isServiceUpdate bool, err errorModel.ErrorModel) {
	var (
		inputStruct                 = inputStructInterface.(in.LicenseConfigRequest)
		licenseConfigModel          repository.LicenseConfigModel
		licenseConfigComponentModel []repository.LicenseConfigComponent
		idLicense                   int64
		scopeLimit                  map[string]interface{}
	)

	/*
		1. Create model from struct request
		2. Validate and get data customer installation by installation id ------------------> By : installation_id (parent_customer_id, customer_id, site_id, product_id, unique_id_1, unique_id_2)
		3. Validate and get data client mapping (client_id) by installation id -------------> By : installation_id (client_id)
		4. Validate and get data product by product id -------------------------------------> By : product_id (client_type_id, license_variant_id, license_type_id, deployment_method, max_offline_days, module_id_1 - module_id_10)
		5. Insert new license configuration with checksum ----------------------------------> Model license config model
		6. Get data component of product by product id -------------------------------------> By : product_id (product_id, component_id, component_value)
		7. Insert multiple license configuration product component with license config id --> Model license config component model
		8. Update customer installation in valid from and thru by installation id ----------> Model license config component model
	*/

	//--- Changing because target on product not client type directly
	input.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "pr.client_type_id",
		Count: "pr.client_type_id",
	}

	//--- Validate data scope
	scopeLimit, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	input.createModel(inputStruct, &licenseConfigModel, contextModel, timeNow)
	err = input.validateCustomerInstallationAndClientMapping(&licenseConfigModel, contextModel, scopeLimit)
	if err.Error != nil {
		return
	}

	idLicense, err = input.insertNewLicenseConfig(tx, &licenseConfigModel, &dataAudit)
	if err.Error != nil {
		return
	}

	output = idLicense
	err = input.getProductComponent(&licenseConfigComponentModel, licenseConfigModel, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	err = input.insertMultipleLicenceConfigProductComponent(tx, licenseConfigComponentModel, &dataAudit)
	if err.Error != nil {
		return
	}

	isServiceUpdate = true
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.CustomerInstallationDAO.TableName, licenseConfigModel.InstallationID.Int64, 0)...)
	err = dao.CustomerInstallationDAO.UpdateCustomerInstallationFromLicenseConfig(tx, licenseConfigModel)
	return
}

func (input licenseConfigService) insertMultipleLicenceConfigProductComponent(tx *sql.Tx, licenseConfigComponent []repository.LicenseConfigComponent, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {
	var idComponent []int64

	if len(licenseConfigComponent) != 0 {
		idComponent, err = dao.LicenseConfigProductComponentDAO.InsertLicenseConfigProductComponent(tx, licenseConfigComponent)
		if err.Error != nil {
			return
		}

		for _, itemIdComponent := range idComponent {
			*dataAudit = append(*dataAudit, repository.AuditSystemModel{
				TableName:  sql.NullString{String: dao.LicenseConfigProductComponentDAO.TableName},
				PrimaryKey: sql.NullInt64{Int64: itemIdComponent},
				Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
			})
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) insertNewLicenseConfig(tx *sql.Tx, licenseConfigModel *repository.LicenseConfigModel, dataAudit *[]repository.AuditSystemModel) (idResult int64, err errorModel.ErrorModel) {
	idResult, err = dao.LicenseConfigDAO.InsertLicenseConfig(tx, *licenseConfigModel)
	if err.Error != nil {
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.LicenseConfigDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idResult},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	licenseConfigModel.ID.Int64 = idResult
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) getProductComponent(licenseConfigComponentModel *[]repository.LicenseConfigComponent, licenseConfigModel repository.LicenseConfigModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	var productComponentModel []repository.ProductComponentModel

	productComponentModel, err = dao.ProductComponentDAO.GetProductComponentByIDProduct(serverconfig.ServerAttribute.DBConnection, licenseConfigModel.ProductID.Int64)
	if err.Error != nil {
		return
	}

	for _, itemProductComponentModel := range productComponentModel {
		*licenseConfigComponentModel = append(*licenseConfigComponentModel, repository.LicenseConfigComponent{
			LicenseConfigID: sql.NullInt64{Int64: licenseConfigModel.ID.Int64},
			ProductID:       sql.NullInt64{Int64: licenseConfigModel.ProductID.Int64},
			ComponentID:     sql.NullInt64{Int64: itemProductComponentModel.ComponentID.Int64},
			ComponentValue:  sql.NullString{String: itemProductComponentModel.ComponentValue.String},
			CreatedBy:       sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			CreatedClient:   sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			CreatedAt:       sql.NullTime{Time: timeNow},
			UpdatedBy:       sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedClient:   sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			UpdatedAt:       sql.NullTime{Time: timeNow},
		})
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) validateCustomerInstallationAndClientMapping(licenseConfigModel *repository.LicenseConfigModel, contextModel *applicationModel.ContextModel, scopeLimit map[string]interface{}) (err errorModel.ErrorModel) {
	var (
		fileName             = "InsertLicenseConfigService.go"
		funcName             = "validateCustomerInstallationAndClientMapping"
		resultInstallationDB repository.CustomerInstallationForConfig
		customerStr          = util2.GenerateConstantaI18n(constanta.Customer, contextModel.AuthAccessTokenModel.Locale, nil)
		installationIDStr    = util2.GenerateConstantaI18n(constanta.InstallationID, contextModel.AuthAccessTokenModel.Locale, nil)
	)

	if resultInstallationDB, err = dao.CustomerInstallationDAO.GetCustomerInstallationByIDJoinClientMappingAndProduct(serverconfig.ServerAttribute.DBConnection,
		repository.CustomerInstallationForConfig{ID: sql.NullInt64{Int64: licenseConfigModel.InstallationID.Int64}}, scopeLimit, input.MappingScopeDB); err.Error != nil {
		return
	}

	if resultInstallationDB.ID.Int64 < 1 {
		err = errorModel.GenerateClientIDNotFound(fileName, funcName, constanta.Installation, customerStr, installationIDStr)
		return
	}

	licenseConfigModel.ParentCustomerID.Int64 = resultInstallationDB.ParentCustomerID.Int64
	licenseConfigModel.CustomerID.Int64 = resultInstallationDB.CustomerID.Int64
	licenseConfigModel.SiteID.Int64 = resultInstallationDB.SiteID.Int64
	licenseConfigModel.ProductID.Int64 = resultInstallationDB.ProductID.Int64
	licenseConfigModel.UniqueID1.String = resultInstallationDB.UniqueID1.String
	licenseConfigModel.UniqueID2.String = resultInstallationDB.UniqueID2.String
	licenseConfigModel.ClientTypeID.Int64 = resultInstallationDB.ClientTypeID.Int64
	licenseConfigModel.LicenseVariantID.Int64 = resultInstallationDB.LicenseVariantID.Int64
	licenseConfigModel.LicenseTypeID.Int64 = resultInstallationDB.LicenseTypeID.Int64
	licenseConfigModel.DeploymentMethod.String = resultInstallationDB.DeploymentMethod.String
	licenseConfigModel.MaxOfflineDays.Int64 = resultInstallationDB.MaxOfflineDays.Int64
	licenseConfigModel.ModuleID1.Int64 = resultInstallationDB.ModuleID1.Int64
	licenseConfigModel.ModuleID2.Int64 = resultInstallationDB.ModuleID2.Int64
	licenseConfigModel.ModuleID3.Int64 = resultInstallationDB.ModuleID3.Int64
	licenseConfigModel.ModuleID4.Int64 = resultInstallationDB.ModuleID4.Int64
	licenseConfigModel.ModuleID5.Int64 = resultInstallationDB.ModuleID5.Int64
	licenseConfigModel.ModuleID6.Int64 = resultInstallationDB.ModuleID6.Int64
	licenseConfigModel.ModuleID7.Int64 = resultInstallationDB.ModuleID7.Int64
	licenseConfigModel.ModuleID8.Int64 = resultInstallationDB.ModuleID8.Int64
	licenseConfigModel.ModuleID9.Int64 = resultInstallationDB.ModuleID9.Int64
	licenseConfigModel.ModuleID10.Int64 = resultInstallationDB.ModuleID10.Int64
	licenseConfigModel.ClientID.String = resultInstallationDB.ClientID.String
	licenseConfigModel.IsUserConcurrent.String = "N"

	if resultInstallationDB.IsUserConcurrent.Bool {
		licenseConfigModel.IsUserConcurrent.String = "Y"
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) validateInsert(inputStruct *in.LicenseConfigRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertLicenseConfig()
}
