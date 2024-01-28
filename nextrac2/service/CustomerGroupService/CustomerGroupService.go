package CustomerGroupService

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

type customerGroupService struct {
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var CustomerGroupService = customerGroupService{}.New()

func (input customerGroupService) New() (output customerGroupService) {
	output.FileName = "CustomerGroupService.go"
	output.ServiceName = "CUSTOMER_GROUP"
	output.ValidLimit = service.DefaultLimit
	output.ValidSearchBy = []string{"customer_group_name","customer_group_id"}
	output.ValidOrderBy = []string{
		"id",
		"customer_group_name",
		"customer_group_id",
		"updated_by",
		"created_at",
		"updated_at",
		"updated_name",
	}

	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.CustomerGroupDataScope] = applicationModel.MappingScopeDB{
		View:  "cg.id",
		Count: "cg.id",
	}

	output.ListScope = input.SetListScope()
	return
}

func (input customerGroupService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(customerGroupRequest *in.CustomerGroupRequest) errorModel.ErrorModel) (inputStruct in.CustomerGroupRequest, err errorModel.ErrorModel) {
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

func (input customerGroupService) CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_customergroup_customer_group_id") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.CustomerGroup)
		}
	}

	return err
}

func (input customerGroupService) validateDataScope(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	output = service.ValidateScope(contextModel, []string{constanta.CustomerGroupDataScope})
	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, "validateDataScope")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
