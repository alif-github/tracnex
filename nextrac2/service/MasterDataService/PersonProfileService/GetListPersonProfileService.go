package PersonProfileService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
)

func (input personProfileService) GetListPersonProfile(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListPersonProfileValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.DoGetListPersonProfile(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input personProfileService) DoGetListPersonProfile(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var personProfileStruct master_data_request.PersonProfileGetListRequest

	personProfileStruct.Page = inputStruct.Page
	personProfileStruct.Limit = inputStruct.Limit
	personProfileStruct.Status = constanta.StatusActive

	for _, param := range searchByParam {
		switch param.SearchKey {
		case "first_name":
			personProfileStruct.FistName = param.SearchValue
			break
		case "email":
			personProfileStruct.Email = param.SearchValue
			break
		case "phone":
			personProfileStruct.Phone = param.SearchValue
			break
		case "nik":
			personProfileStruct.NIK = param.SearchValue
			break
		}
	}

	output, err = master_data_dao.GetListPersonProfileFromMasterData(personProfileStruct, contextModel)
	return
}
