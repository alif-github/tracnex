package CustomerCategoryService

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

type customerCategoryService struct {
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var CustomerCategoryService = customerCategoryService{}.New()

func (input customerCategoryService) New() (output customerCategoryService) {
	output.FileName = "CustomerCategoryService.go"
	output.ServiceName = "CUSTOMER_CATEGORY"
	output.ValidLimit = service.DefaultLimit
	output.ValidSearchBy = []string{"customer_category_name","customer_category_id"}
	output.ValidOrderBy = []string{
		"id",
		"customer_category_name",
		"customer_category_id",
		"updated_at",
		"created_at",
		"updated_by",
		"updated_name",
	}

	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.CustomerCategoryDataScope] = applicationModel.MappingScopeDB{
		View:  "cc.id",
		Count: "cc.id",
	}

	output.ListScope = input.SetListScope()
	return
}

func (input customerCategoryService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(customerCategoryRequest *in.CustomerCategoryRequest) errorModel.ErrorModel) (inputStruct in.CustomerCategoryRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateErrorFormatJSON(input.FileName, "readBodyAndValidate", errorS)
			return
		}
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	inputStruct.ID = int64(id)

	err = validation(&inputStruct)
	return
}

func (input customerCategoryService) CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_customercategory_customer_category_id") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.CustomerCategory)
		}
	}

	return err
}

func (input customerCategoryService) validateDataScope(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	output = service.ValidateScope(contextModel, []string{constanta.CustomerCategoryDataScope})
	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, "validateDataScope")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}


