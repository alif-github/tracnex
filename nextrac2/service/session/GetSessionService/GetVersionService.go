package GetSessionService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service/session"
	"nexsoft.co.id/nextrac2/util"
)

func (input getSessionService) GetSystemVersion(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	output.Data.Content = config.ApplicationConfiguration.GetServerVersion()
	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: session.GenerateLoginI18NMessage("GET_VERSION_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	return
}
