package ClientService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
)

func (input clientService) InitiateRegistrationClient(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	result, errS := input.doInitiateRegistrationClient()
	if errS.Error != nil {
		return
	}

	output.Data.Content = out.InitiateClientTypeResponse {
		ClientTypeList: result,
	}
	output.Status = out.StatusResponse {
		Code: 		util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: 	GenerateI18NMessage("SUCCESS_INITIATE_LIST_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientService) doInitiateRegistrationClient() (result []out.ClientTypeResponse, err errorModel.ErrorModel) {
	result, err = dao.ClientTypeDAO.GetInitiateInsertClientService(serverconfig.ServerAttribute.DBConnection)

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}


