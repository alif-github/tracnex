package ClientTypeService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
)

func (input clientTypeService) InitiateClientType(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData interface{}

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListClientTypeValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateClientType(searchByParam, contextModel)

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_GET_LIST_CLIENT_TYPE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListClientTypeValidOperator,
		CountData:     countData.(int),
	}
	return
}

func (input clientTypeService) doInitiateClientType(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		createdBy int64
	)

	output = 0
	createdBy = contextModel.LimitedByCreatedBy

	for i := 0; i < len(searchByParam); i++ {
		if searchByParam[i].SearchKey == "client_type_id" {
			searchByParam[i].SearchKey = "id"
		}
	}

	output, err = dao.ClientTypeDAO.GetCountClientType(serverconfig.ServerAttribute.DBConnection, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	return
}
