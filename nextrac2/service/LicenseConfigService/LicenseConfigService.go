package LicenseConfigService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service"
)

type licenseConfigService struct {
	service.AbstractService
	service.GetListData
}

var LicenseConfigService = licenseConfigService{}.New()

func (input licenseConfigService) New() (output licenseConfigService) {
	var (
		parent = "_parent"
		child  = "_child"
	)

	output.FileName = "LicenseConfigService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"product_valid_thru",
		"id",
		"customer_name",
		"product_name",
		"license_variant_name",
		"license_type_name",
		"unique_id_1",
		"unique_id_2",
		"installation_id",
		"product_valid_from",
		"allow_activation",
	}
	output.ValidSearchBy = []string{
		"id",
		"parent_customer_id",
		"customer_name",
		"distributor_of",
		"product_valid_from",
		"product_valid_thru",
		"product_id",
		"client_type_id",
		"province_id",
		"district_id",
		"unique_id_1",
		"license_status",
		"customer_group_id",
		"customer_category_id",
		"salesman_id",
		"allow_activation",
	}

	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ProvinceDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cup.province_id",
		Count: "cup.province_id",
	}
	output.MappingScopeDB[constanta.DistrictDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cup.district_id",
		Count: "cup.district_id",
	}
	output.MappingScopeDB[constanta.SalesmanDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cup.salesman_id",
		Count: "cup.salesman_id",
	}
	output.MappingScopeDB[constanta.CustomerGroupDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cup.customer_group_id",
		Count: "cup.customer_group_id",
	}
	output.MappingScopeDB[constanta.CustomerCategoryDataScope+parent] = applicationModel.MappingScopeDB{
		View:  "cup.customer_category_id",
		Count: "cup.customer_category_id",
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
		View:  "pr.product_group_id",
		Count: "pr.product_group_id",
	}
	output.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "lc.client_type_id",
		Count: "lc.client_type_id",
	}

	return
}

func (input licenseConfigService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.LicenseConfigRequest) errorModel.ErrorModel) (inputStruct in.LicenseConfigRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if request.Method != "GET" {
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

func (input licenseConfigService) createModel(inputStruct in.LicenseConfigRequest, licenseConfigModel *repository.LicenseConfigModel, contextModel *applicationModel.ContextModel, timeNow time.Time) {
	*licenseConfigModel = repository.LicenseConfigModel{
		InstallationID:   sql.NullInt64{Int64: inputStruct.InstallationID},
		NoOfUser:         sql.NullInt64{Int64: inputStruct.NoOfUser},
		ProductValidFrom: sql.NullTime{Time: inputStruct.ProductValidFrom},
		ProductValidThru: sql.NullTime{Time: inputStruct.ProductValidThru},
		CreatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:        sql.NullTime{Time: timeNow},
		UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
	}
}

func (input licenseConfigService) createModelForInsertMultiple(inputStruct in.LicenseConfigMultipleRequest, licenseConfigModel *repository.LicenseConfigModel, mapLicenseConfigID map[int64]bool, contextModel *applicationModel.ContextModel, timeNow time.Time) (err errorModel.ErrorModel) {
	var (
		funcName             = "createModelForInsertMultiple"
		licenseConfigIDModel []repository.LicenseConfigIDsModel
	)

	for _, item := range inputStruct.LicenseConfigID {
		if item < 1 {
			err = errorModel.GenerateEmptyFieldOrZeroValueError(input.FileName, funcName, constanta.ID)
			return
		}

		_, ok := mapLicenseConfigID[item]
		if ok {
			err = errorModel.GenerateDuplicateErrorWithParam(input.FileName, funcName, fmt.Sprintf(`ID: %d`, item))
			return
		}

		licenseConfigIDModel = append(licenseConfigIDModel, repository.LicenseConfigIDsModel{ID: sql.NullInt64{Int64: item}})
		mapLicenseConfigID[item] = true
	}

	*licenseConfigModel = repository.LicenseConfigModel{
		ProductValidThru: sql.NullTime{Time: inputStruct.ProductValidThru},
		LicenseConfigIDs: licenseConfigIDModel,
		CreatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:        sql.NullTime{Time: timeNow},
		UpdatedBy:        sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:    sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:        sql.NullTime{Time: timeNow},
	}

	return
}

func (input licenseConfigService) validateDataScope(contextModel *applicationModel.ContextModel) (newScope map[string]interface{}, err errorModel.ErrorModel) {
	var scope map[string]interface{}

	scope, err = input.validateMultipleDataScopeLicenseConfig(contextModel)
	if err.Error != nil {
		return
	}

	newScope = input.restructureDataScopeLicenseConfig(scope)
	return
}
