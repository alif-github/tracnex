package CRUDUserService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	Login "nexsoft.co.id/nextrac2/service/session/Logout"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
)

func (input userService) KillSessionUserActive(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	keys := serverconfig.ServerAttribute.RedisClientSession.Keys(constanta.AllSessionUser).Val()

	for _, key := range keys {
		serverconfig.ServerAttribute.RedisClientSession.Del(key)

		token := strings.TrimPrefix(key, constanta.SessionUser)
		serverconfig.ServerAttribute.RedisClient.Del(token)

		err = Login.HitLogoutAuthenticationServer(token, contextModel)
		if err.Error != nil {
			return
		}
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_LIST_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
