package PKCEClientMappingService

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type pkceClientMappingService struct {
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var PKCEClientMappingService = pkceClientMappingService{}.New()

func (input pkceClientMappingService) New() (output pkceClientMappingService) {
	output.FileName = "PKCEClientMappingService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidSearchBy = []string{"parent_client_id", "client_alias"}
	output.ValidOrderBy = []string{
		"parent_client_id",
		"nt_username",
		"client_type",
		"company_id",
		"branch_id",
		"client_alias",
	}
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "pkce_client_mapping.client_type_id",
		Count: "pkce_client_mapping.client_type_id",
	}
	return
}

func (input pkceClientMappingService) validateDataScopePKCEClientMapping(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	funcName := "validateDataScopePKCEClientMapping"

	output = service.ValidateScope(contextModel, []string{
		constanta.ClientTypeDataScope,
	})

	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
