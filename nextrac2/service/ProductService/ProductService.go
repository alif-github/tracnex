package ProductService

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"time"
)

type productService struct {
	service.AbstractService
	service.GetListData
}

var ProductService = productService{}.New()

func (input productService) New() (output productService) {
	output.FileName = "ProductService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"product_name",
		"product_id",
		"product_description",
		"product_group_name",
		"client_type",
		"license_variant_name",
		"license_type_name",
	}
	output.ValidSearchBy = []string{
		"product_id",
		"product_name",
		"product_group_id",
	}

	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "product.client_type_id",
		Count: "product.client_type_id",
	}
	output.MappingScopeDB[constanta.ProductGroupDataScope] = applicationModel.MappingScopeDB{
		View:  "product.product_group_id",
		Count: "product.product_group_id",
	}
	return
}

func (input productService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ProductRequest) errorModel.ErrorModel) (inputStruct in.ProductRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if !util.IsStringEmpty(stringBody) {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func (input productService) validateRelation(inputStruct in.ProductRequest, contextModel *applicationModel.ContextModel, scopeLimit map[string]interface{}) (err errorModel.ErrorModel) {
	var (
		fileName           = "InsertProductService.go"
		funcName           = "validateRelation"
		productGroupOnDB   repository.ProductGroupModel
		clientTypeOnDB     repository.ClientTypeModel
		licenseVariantOnDB repository.LicenseVariantModel
		licenseTypeOnDB    repository.LicenseTypeModel
		moduleOnDB         repository.ModuleModel
		componentOnDB      repository.ComponentModel
	)

	mappingScopeDB := make(map[string]applicationModel.MappingScopeDB)
	mappingScopeDB[constanta.ProductGroupDataScope] = applicationModel.MappingScopeDB{View: "pg.id", Count: "pg.id"}
	mappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{View: "ct.id", Count: "ct.id"}

	productGroupOnDB, err = dao.ProductGroupDAO.CheckIsProductGroupExist(serverconfig.ServerAttribute.DBConnection, repository.ProductGroupModel{ID: sql.NullInt64{Int64: inputStruct.ProductGroupID}}, scopeLimit, mappingScopeDB)
	if err.Error != nil {
		return
	}

	if productGroupOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ProductGroupID)
		return
	}

	clientTypeOnDB, err = dao.ClientTypeDAO.CheckClientTypeByIDWithScope(serverconfig.ServerAttribute.DBConnection, &repository.ClientTypeModel{ID: sql.NullInt64{Int64: inputStruct.ClientTypeID}}, scopeLimit, mappingScopeDB)
	if err.Error != nil {
		return
	}

	if clientTypeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.NewClientType)
		return
	}

	licenseVariantOnDB, err = dao.LicenseVariantDAO.CheckIsExistLicenseVariant(serverconfig.ServerAttribute.DBConnection, repository.LicenseVariantModel{ID: sql.NullInt64{Int64: inputStruct.LicenseVariantID}})
	if err.Error != nil {
		return
	}

	if licenseVariantOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.LicenseVariantID)
		return
	}

	licenseTypeOnDB, err = dao.LicenseTypeDAO.CheckLicenseTypeIsExist(serverconfig.ServerAttribute.DBConnection, repository.LicenseTypeModel{ID: sql.NullInt64{Int64: inputStruct.LicenseTypeID}})
	if err.Error != nil {
		return
	}

	if licenseTypeOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.LicenseTypeID)
		return
	}

	collectionModule := []int64{
		inputStruct.Module1, inputStruct.Module2, inputStruct.Module3,
		inputStruct.Module4, inputStruct.Module5, inputStruct.Module6,
		inputStruct.Module7, inputStruct.Module8, inputStruct.Module9,
		inputStruct.Module10}

	err = input.TracAndCheckDuplicateModule(collectionModule)
	if err.Error != nil {
		return
	}

	for index, valueModule := range collectionModule {
		if valueModule > 0 {
			moduleOnDB, err = dao.ModuleDAO.CheckModuleIsExist(serverconfig.ServerAttribute.DBConnection, repository.ModuleModel{ID: sql.NullInt64{Int64: valueModule}})
			if err.Error != nil {
				return
			}

			if moduleOnDB.ID.Int64 < 1 {
				err = errorModel.GenerateUnknownDataError(fileName, funcName, util2.GenerateConstantaI18n(constanta.Module, contextModel.AuthAccessTokenModel.Locale, nil)+" "+strconv.Itoa(index+1))
				return
			}
		} else {
			continue
		}
	}

	for indexComponent, valueComponent := range inputStruct.Component {
		componentOnDB, err = dao.ComponentDAO.CheckComponentIsExist(serverconfig.ServerAttribute.DBConnection, repository.ComponentModel{ID: sql.NullInt64{Int64: valueComponent.ComponentID}})
		if err.Error != nil {
			return
		}

		if componentOnDB.ID.Int64 < 1 {
			err = errorModel.GenerateUnknownDataError(fileName, funcName, util2.GenerateConstantaI18n(constanta.Component, contextModel.AuthAccessTokenModel.Locale, nil)+" "+strconv.Itoa(indexComponent+1))
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_product_product_name") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ProductName)
		} else if service.CheckDBError(err, "uq_product_productid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ProductID)
		}
	}
	return err
}

func (input productService) TracAndCheckDuplicateModule(collectionModule []int64) (err errorModel.ErrorModel) {
	funcName := "TracAndCheckDuplicateModule"

	var collIsolationModule []int64

	for _, valueArrayModule := range collectionModule {

		for _, valueCollIsolationModule := range collIsolationModule {
			if valueCollIsolationModule == valueArrayModule {
				err = errorModel.GenerateDuplicateErrorWithParam(input.FileName, funcName, constanta.Module)
				return
			}
		}

		if valueArrayModule > 0 {
			collIsolationModule = append(collIsolationModule, valueArrayModule)
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) createModelProduct(inputStruct in.ProductRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (parameterProductModel repository.ProductModel) {
	var productComponentModel []repository.ProductComponentModel
	for _, valueComponentModel := range inputStruct.Component {
		productComponentModel = append(productComponentModel, repository.ProductComponentModel{
			ID:             sql.NullInt64{Int64: valueComponentModel.ID},
			ComponentID:    sql.NullInt64{Int64: valueComponentModel.ComponentID},
			ComponentValue: sql.NullString{String: valueComponentModel.ComponentValue},
			UpdatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			UpdatedAt:      sql.NullTime{Time: timeNow},
			Deleted:        sql.NullBool{Bool: valueComponentModel.Deleted},
		})
	}

	parameterProductModel = repository.ProductModel{
		ID:                    sql.NullInt64{Int64: inputStruct.ID},
		ProductID:             sql.NullString{String: inputStruct.ProductID},
		ProductName:           sql.NullString{String: inputStruct.ProductName},
		ProductDescription:    sql.NullString{String: inputStruct.ProductDescription},
		ProductGroupID:        sql.NullInt64{Int64: inputStruct.ProductGroupID},
		ClientTypeID:          sql.NullInt64{Int64: inputStruct.ClientTypeID},
		IsLicense:             sql.NullBool{Bool: inputStruct.IsLicense},
		LicenseTypeID:         sql.NullInt64{Int64: inputStruct.LicenseTypeID},
		LicenseVariantID:      sql.NullInt64{Int64: inputStruct.LicenseVariantID},
		DeploymentMethod:      sql.NullString{String: inputStruct.DeploymentMethod},
		NoOfUser:              sql.NullInt64{Int64: inputStruct.NoOfUser},
		IsUserConcurrent:      sql.NullBool{Bool: inputStruct.IsUserConcurrent},
		MaxOfflineDays:        sql.NullInt64{Int64: inputStruct.MaxOfflineDays},
		Module1:               sql.NullInt64{Int64: inputStruct.Module1},
		Module2:               sql.NullInt64{Int64: inputStruct.Module2},
		Module3:               sql.NullInt64{Int64: inputStruct.Module3},
		Module4:               sql.NullInt64{Int64: inputStruct.Module4},
		Module5:               sql.NullInt64{Int64: inputStruct.Module5},
		Module6:               sql.NullInt64{Int64: inputStruct.Module6},
		Module7:               sql.NullInt64{Int64: inputStruct.Module7},
		Module8:               sql.NullInt64{Int64: inputStruct.Module8},
		Module9:               sql.NullInt64{Int64: inputStruct.Module9},
		Module10:              sql.NullInt64{Int64: inputStruct.Module10},
		CreatedBy:             sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:         sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:             sql.NullTime{Time: timeNow},
		UpdatedBy:             sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:         sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:             sql.NullTime{Time: timeNow},
		ProductComponentModel: productComponentModel,
	}

	return
}

func (input productService) validateDataScope(contextModel *applicationModel.ContextModel) (scope map[string]interface{}, err errorModel.ErrorModel) {
	scope, err = input.ValidateMultipleDataScope(contextModel, []string{
		constanta.ClientTypeDataScope,
		constanta.ProductGroupDataScope,
	})

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
