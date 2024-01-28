package DistrictService

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type districtService struct {
	service.AbstractService
	service.GetListData
}

var DistrictService = districtService{}.New()

func (input districtService) New() (output districtService) {
	output.FileName = "DistrictService.go"
	output.ServiceName = "DISTRICT"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id", "name"}
	output.ValidSearchBy = []string{"province_id", "id", "code", "name"}
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.DistrictDataScope] = applicationModel.MappingScopeDB{
		View:  "d.id",
		Count: "d.id",
	}
	output.MappingScopeDB[constanta.ProvinceDataScope] = applicationModel.MappingScopeDB{
		View:  "p.id",
		Count: "p.id",
	}
	output.ListScope = input.SetListScope()
	return
}

func (input districtService) validateDataScope(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	output = service.ValidateScope(contextModel, []string{constanta.DistrictDataScope})
	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, "validateDataScope")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
