package ProvinceService

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type provinceService struct {
	service.AbstractService
	service.GetListData
}

var ProvinceService = provinceService{}.New()

func (input provinceService) New() (output provinceService) {
	output.FileName = "ProvinceService.go"
	output.ServiceName = "PROVINCE"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"id", "name"}
	output.ValidSearchBy = []string{"country_id", "id", "mdb_province_id", "code", "name"}
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ProvinceDataScope] = applicationModel.MappingScopeDB{
		View:  "province.id",
		Count: "province.id",
	}
	output.ListScope = input.SetListScope()
	return
}

func (input provinceService) validateDataScope(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	output = service.ValidateScope(contextModel, []string{constanta.ProvinceDataScope})
	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, "validateDataScope")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
