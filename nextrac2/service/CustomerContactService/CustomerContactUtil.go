package CustomerContactService

import (
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	util2 "nexsoft.co.id/nextrac2/util"
)

func getErrorMessage(err errorModel.ErrorModel, contextModel *applicationModel.ContextModel, identityCode string) (output out.CustomerContactErrorStatus) {
	errCode := err.Error.Error()
	errMessage := util2.GenerateI18NErrorMessage(err, contextModel.AuthAccessTokenModel.Locale)
	if errMessage == errCode {
		if err.CausedBy != nil {
			errMessage = err.CausedBy.Error()
		}
	}
	return out.CustomerContactErrorStatus{
		Nik:     identityCode,
		Status:  constanta.JobProcessErrorStatus,
		Message: errMessage,
	}
}
