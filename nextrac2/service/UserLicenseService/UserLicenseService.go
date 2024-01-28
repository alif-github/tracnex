package UserLicenseService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type userLicenseService struct {
	service.AbstractService
	service.GetListData
}

var UserLicenseService = userLicenseService{}.New()

func (input userLicenseService) New() (output userLicenseService) {
	output.FileName = "UserLicenseService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"license_config_id",
		"customer_name",
		"unique_id_1",
		"unique_id_2",
		"installation_id",
		"total_license",
		"total_active",
	}
	output.ValidSearchBy = []string{
		"customer_name",
		"id",
	}

	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ProvinceDataScope] = applicationModel.MappingScopeDB{
		View:  "cu.province_id",
		Count: "cu.province_id",
	}

	output.MappingScopeDB[constanta.DistrictDataScope] = applicationModel.MappingScopeDB{
		View:  "cu.district_id",
		Count: "cu.district_id",
	}

	output.MappingScopeDB[constanta.CustomerGroupDataScope] = applicationModel.MappingScopeDB{
		View:  "cu.customer_group_id",
		Count: "cu.customer_group_id",
	}

	output.MappingScopeDB[constanta.CustomerCategoryDataScope] = applicationModel.MappingScopeDB{
		View:  "cu.customer_category_id",
		Count: "cu.customer_category_id",
	}

	output.MappingScopeDB[constanta.SalesmanDataScope] = applicationModel.MappingScopeDB{
		View:  "cu.salesman_id",
		Count: "cu.salesman_id",
	}

	output.MappingScopeDB[constanta.ProductGroupDataScope] = applicationModel.MappingScopeDB{
		View:  "pg.id",
		Count: "pg.id",
	}

	output.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "ct.id",
		Count: "ct.id",
	}

	return
}

func (input userLicenseService) readBodyCustomUserLicenseView(request *http.Request) (inputStruct in.GetListDataDTO, userLicenseViewStruct in.ViewUserLicenseRequest) {
	var userLicenseId int

	inputStruct.Page, _ = strconv.Atoi(service.GenerateQueryValue(request.URL.Query()["page"]))
	inputStruct.Limit, _ = strconv.Atoi(service.GenerateQueryValue(request.URL.Query()["limit"]))
	inputStruct.Filter = service.GenerateQueryValue(request.URL.Query()["filter"])
	inputStruct.Search = service.GenerateQueryValue(request.URL.Query()["search"])
	inputStruct.OrderBy = service.GenerateQueryValue(request.URL.Query()["order"])
	userLicenseId, _ = strconv.Atoi(service.GenerateQueryValue(request.URL.Query()["id"]))
	userLicenseViewStruct.UserLicenseId = int64(userLicenseId)

	return
}

func (input userLicenseService) readCountDataUserRegistrationDetail(request *http.Request, validSearchBy []string, validOperator map[string]applicationModel.DefaultOperator) (inputStruct in.GetListDataDTO, userLicenseViewStruct in.ViewUserLicenseRequest, searchByParam []in.SearchByParam, err errorModel.ErrorModel) {
	inputStruct, userLicenseViewStruct = input.readBodyCustomUserLicenseView(request)
	searchByParam, err = inputStruct.ValidateGetCountData(validSearchBy, validOperator)
	return
}

func (input userLicenseService) readGetListUserRegistrationDetail(request *http.Request, validSearchBy []string, validOrderBy []string, validOperator map[string]applicationModel.DefaultOperator, validLimit []int) (inputStruct in.GetListDataDTO, userLicenseViewStruct in.ViewUserLicenseRequest, searchByParam []in.SearchByParam, err errorModel.ErrorModel) {
	inputStruct, userLicenseViewStruct = input.readBodyCustomUserLicenseView(request)
	searchByParam, err = inputStruct.ValidateGetListData(validSearchBy, validOrderBy, validOperator, validLimit)
	return
}

func (input userLicenseService) readBodyAndParamTransferKey(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.TransferUserLicenseRequest) errorModel.ErrorModel) (inputStruct in.TransferUserLicenseRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndParamTransferKey"

	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}