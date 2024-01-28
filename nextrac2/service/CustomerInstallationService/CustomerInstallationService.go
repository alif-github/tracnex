package CustomerInstallationService

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"time"
)

type customerInstallationService struct {
	service.AbstractService
	service.GetListData
}

var CustomerInstallationService = customerInstallationService{}.New()

func (input customerInstallationService) New() (output customerInstallationService) {
	var (
		parent = "_parent"
		child  = "_child"
	)

	output.FileName = "CustomerInstallationService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id"}
	output.ValidSearchBy = []string{"id"}

	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ProvinceDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cu.province_id",
		Count: "cu.province_id",
	}
	output.MappingScopeDB[constanta.DistrictDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cu.district_id",
		Count: "cu.district_id",
	}
	output.MappingScopeDB[constanta.SalesmanDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cu.salesman_id",
		Count: "cu.salesman_id",
	}
	output.MappingScopeDB[constanta.CustomerGroupDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cu.customer_group_id",
		Count: "cu.customer_group_id",
	}
	output.MappingScopeDB[constanta.CustomerCategoryDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cu.customer_category_id",
		Count: "cu.customer_category_id",
	}
	output.MappingScopeDB[constanta.ProvinceDataScope+child] = applicationModel.MappingScopeDB{
		View:  "cuc.province_id",
		Count: "cuc.province_id",
	}
	output.MappingScopeDB[constanta.DistrictDataScope+child] = applicationModel.MappingScopeDB{
		View:  "cuc.district_id",
		Count: "cuc.district_id",
	}
	output.MappingScopeDB[constanta.SalesmanDataScope+child] = applicationModel.MappingScopeDB{
		View:  "cuc.salesman_id",
		Count: "cuc.salesman_id",
	}
	output.MappingScopeDB[constanta.CustomerGroupDataScope+child] = applicationModel.MappingScopeDB{
		View:  "cuc.customer_group_id",
		Count: "cuc.customer_group_id",
	}
	output.MappingScopeDB[constanta.CustomerCategoryDataScope+child] = applicationModel.MappingScopeDB{
		View:  "cuc.customer_category_id",
		Count: "cuc.customer_category_id",
	}
	output.MappingScopeDB[constanta.ProductGroupDataScope] = applicationModel.MappingScopeDB{
		View:  "pd.product_group_id",
		Count: "pd.product_group_id",
	}
	output.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "pd.client_type_id",
		Count: "pd.client_type_id",
	}
	return
}

func (input customerInstallationService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.CustomerSiteInstallationRequest, contextModel *applicationModel.ContextModel) errorModel.ErrorModel) (inputStruct in.CustomerSiteInstallationRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidate"
		stringBody string
	)

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

	id, errorID := strconv.Atoi(mux.Vars(request)["ID"])
	if errorID != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorID)
		return
	}

	//--------------- Validate Parent Customer And Compare With Request Body
	if id < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(input.FileName, funcName, constanta.ParentCustomerID)
		return
	}

	if inputStruct.ParentCustomerID < 1 {
		err = errorModel.GenerateEmptyFieldOrZeroValueError(input.FileName, funcName, constanta.ParentCustomerID)
		return
	}

	if inputStruct.ParentCustomerID != int64(id) {
		err = errorModel.GenerateForbiddenAccessClientError(input.FileName, funcName)
		return
	}

	inputStruct.ParentCustomerID = int64(id)
	err = validation(&inputStruct, contextModel)
	return
}

func (input customerInstallationService) readBodyAndValidateForViewSite(request *http.Request, validation func(input *in.CustomerSiteInstallationRequest) errorModel.ErrorModel) (inputStruct in.CustomerSiteInstallationRequest, err errorModel.ErrorModel) {
	var (
		parentCustomerID int
		page             int
		limit            int
	)

	parentCustomerID, _ = strconv.Atoi(service.GenerateQueryValue(request.URL.Query()["parent_customer_id"]))
	page, _ = strconv.Atoi(service.GenerateQueryValue(request.URL.Query()["page"]))
	limit, _ = strconv.Atoi(service.GenerateQueryValue(request.URL.Query()["limit"]))

	inputStruct.ParentCustomerID = int64(parentCustomerID)
	inputStruct.Page = page
	inputStruct.Limit = limit

	err = inputStruct.ValidateInputPageLimitAndOrderBy(input.ValidLimit, input.ValidOrderBy)
	if err.Error != nil {
		return
	}

	err = validation(&inputStruct)
	return
}

func (input customerInstallationService) readBodyAndValidateForViewInstallation(request *http.Request, validation func(input *in.CustomerInstallationDetailRequest) errorModel.ErrorModel) (inputStruct in.CustomerInstallationDetailRequest, err errorModel.ErrorModel) {
	parentCustomerID, _ := strconv.Atoi(service.GenerateQueryValue(request.URL.Query()["parent_customer_id"]))
	siteID, _ := strconv.Atoi(service.GenerateQueryValue(request.URL.Query()["site_id"]))
	clientTypeID, _ := strconv.Atoi(service.GenerateQueryValue(request.URL.Query()["client_type_id"]))
	isLicense := service.GenerateQueryValue(request.URL.Query()["is_license"])

	inputStruct.ParentCustomerID = int64(parentCustomerID)
	inputStruct.SiteID = int64(siteID)
	inputStruct.ClientTypeID = int64(clientTypeID)
	inputStruct.IsLicense, _ = strconv.ParseBool(isLicense)

	err = validation(&inputStruct)
	return
}

func (input customerInstallationService) createModel(inputStruct in.CustomerSiteInstallationRequest, customerInstallationModel *repository.CustomerInstallationModel, contextModel *applicationModel.ContextModel, timeNow time.Time) {
	var customerSiteModel []repository.CustomerInstallationData

	for _, itemSite := range inputStruct.CustomerSite {
		var customerInstallationDetail []repository.CustomerInstallationDetail

		for _, itemInstallation := range itemSite.CustomerInstallation {
			customerInstallationDetail = append(customerInstallationDetail, repository.CustomerInstallationDetail{
				InstallationID:     sql.NullInt64{Int64: itemInstallation.InstallationID},
				ProductID:          sql.NullInt64{Int64: itemInstallation.ProductID},
				ParentClientTypeID: sql.NullInt64{Int64: itemInstallation.ParentClientTypeID},
				Remark:             sql.NullString{String: itemInstallation.Remark},
				UniqueID1:          sql.NullString{String: itemInstallation.UniqueID1},
				UniqueID2:          sql.NullString{String: itemInstallation.UniqueID2},
				InstallationStatus: sql.NullString{String: "A"},
				InstallationDate:   sql.NullTime{Time: itemInstallation.InstallationDate},
				ProductValidFrom:   sql.NullTime{Time: itemInstallation.ProductValidFrom},
				ProductValidThru:   sql.NullTime{Time: itemInstallation.ProductValidThru},
				Action:             sql.NullInt32{Int32: itemInstallation.Action},
				UpdatedAt:          sql.NullTime{Time: itemInstallation.UpdatedAt},
			})
		}

		customerSiteModel = append(customerSiteModel, repository.CustomerInstallationData{
			SiteID:       sql.NullInt64{Int64: itemSite.SiteID},
			CustomerID:   sql.NullInt64{Int64: itemSite.CustomerID},
			Installation: customerInstallationDetail,
			Action:       sql.NullInt32{Int32: itemSite.Action},
			UpdatedAt:    sql.NullTime{Time: itemSite.UpdatedAt},
		})
	}

	*customerInstallationModel = repository.CustomerInstallationModel{
		ParentCustomerID:         sql.NullInt64{Int64: inputStruct.ParentCustomerID},
		CustomerInstallationData: customerSiteModel,
		CreatedBy:                sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:            sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:                sql.NullTime{Time: timeNow},
		UpdatedBy:                sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:            sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:                sql.NullTime{Time: timeNow},
	}
}

func (input customerInstallationService) readBodyAndValidateForDetailInstallation(request *http.Request, validation func(input *in.CustomerInstallationDataRequest) errorModel.ErrorModel) (inputStruct in.CustomerInstallationDataRequest, err errorModel.ErrorModel) {

	installationID, _ := strconv.Atoi(mux.Vars(request)["INSTALLATIONID"])

	inputStruct.InstallationID = int64(installationID)

	err = validation(&inputStruct)
	return
}

func (input customerInstallationService) validateDataScope(contextModel *applicationModel.ContextModel) (newScope map[string]interface{}, err errorModel.ErrorModel) {
	var scope map[string]interface{}
	scope, err = input.validateMultipleDataScopeCustomerInstallation(contextModel)
	if err.Error != nil {
		return
	}

	newScope = input.restructureDataScopeCustomerInstallation(scope)
	return
}

func (input customerInstallationService) validateMultipleDataScopeCustomerInstallation(contextModel *applicationModel.ContextModel) (getScope map[string]interface{}, err errorModel.ErrorModel) {
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

func (input customerInstallationService) restructureDataScopeCustomerInstallation(scope map[string]interface{}) (newScope map[string]interface{}) {
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

func (input customerInstallationService) newMappingScopeDB() (mappingScopeDB map[string]applicationModel.MappingScopeDB) {
	parent := "_parent"
	mappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	mappingScopeDB[constanta.ProvinceDataScope+parent] = applicationModel.MappingScopeDB{View: "c.province_id"}
	mappingScopeDB[constanta.DistrictDataScope+parent] = applicationModel.MappingScopeDB{View: "c.district_id"}
	mappingScopeDB[constanta.SalesmanDataScope+parent] = applicationModel.MappingScopeDB{View: "c.salesman_id"}
	mappingScopeDB[constanta.CustomerGroupDataScope+parent] = applicationModel.MappingScopeDB{View: "c.customer_group_id"}
	mappingScopeDB[constanta.CustomerCategoryDataScope+parent] = applicationModel.MappingScopeDB{View: "c.customer_category_id"}
	return
}
