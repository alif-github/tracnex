package PersonTitleService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
	"nexsoft.co.id/nextrac2/util"
)

func (input personTitleService) GetListPersonTitle(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListPersonTitleValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListPersonTitle(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_PERSON_TITLE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input personTitleService) doGetListPersonTitle(inputStruct in.GetListDataDTO, _ []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var personTitleStruct in.PersonTitleRequest
	var pageLimit in.AbstractDTO

	pageLimit = in.AbstractDTO{
		Page:    inputStruct.Page,
		Limit:   inputStruct.Limit,
		OrderBy: inputStruct.OrderBy,
	}

	personTitleStruct = in.PersonTitleRequest{AbstractDTO: pageLimit}

	output, err = master_data_dao.GetListPersonTitleFromMasterData(personTitleStruct, contextModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

//func (input personTitleService) validateListPersonTitle(inputStruct *in.PersonTitleRequest) errorModel.ErrorModel {
//	return inputStruct.ValidateInputPageLimitAndOrderBy(input.ValidLimit, input.ValidOrderBy)
//}