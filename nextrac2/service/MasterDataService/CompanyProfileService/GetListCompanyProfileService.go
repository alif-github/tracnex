package CompanyProfileService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
	"strconv"
)

func (input companyProfileService) GetListCompanyProfile(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListCompanyProfileValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.DoGetListCompanyProfile(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input companyProfileService) DoGetListCompanyProfile(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var companyProfileStruct master_data_request.CompanyProfileGetListRequest

	companyProfileStruct.Page = inputStruct.Page
	companyProfileStruct.Limit = inputStruct.Limit
	companyProfileStruct.OrderBy = inputStruct.OrderBy

	for _, param := range searchByParam {
		switch param.SearchKey {
		case "id":
			idForSearch, _ := strconv.Atoi(param.SearchValue)
			companyProfileStruct.ID = int64(idForSearch)
			break
		case "name":
			companyProfileStruct.Name = param.SearchValue
			break
		case "npwp":
			companyProfileStruct.NPWP = param.SearchValue
			break
		}
	}

	output, err = master_data_dao.GetListCompanyProfileFromMasterData(companyProfileStruct, contextModel)
	return
}
