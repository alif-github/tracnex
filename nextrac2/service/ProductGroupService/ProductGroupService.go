package ProductGroupService

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

type productGroupService struct {
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var ProductGroupService = productGroupService{}.New()

func (input productGroupService) New() (output productGroupService) {
	output.FileName = "ProductGroupService.go"
	output.ServiceName = "PRODUCT_GROUP"
	output.ValidLimit = service.DefaultLimit
	output.ValidSearchBy = []string{"id","product_group_name"}
	output.ValidOrderBy = []string{
		"id",
		"product_group_name",
		"updated_at",
		"updated_by",
		"created_at",
		"updated_name",
	}

	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ProductGroupDataScope] = applicationModel.MappingScopeDB{
		View:  "pg.id",
		Count: "pg.id",
	}

	output.ListScope = input.SetListScope()

	return
}

func (input productGroupService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(productGroupRequest *in.ProductGroupRequest) errorModel.ErrorModel) (inputStruct in.ProductGroupRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		var errorS = json.Unmarshal([]byte(stringBody), &inputStruct)
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

func (input productGroupService) CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_product_group_product_group_name") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ProductGroup)
		}
	}

	return err
}

func (input productGroupService) validateDataScope(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	output = service.ValidateScope(contextModel, []string{constanta.ProductGroupDataScope})
	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, "validateDataScope")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}