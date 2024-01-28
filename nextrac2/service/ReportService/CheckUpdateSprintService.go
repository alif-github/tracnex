package ReportService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input reportService) CheckUpdateSprint(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName     = "CheckUpdateSprint"
		dbRedmineDev = serverconfig.ServerAttribute.RedmineDBConnection
	)

	_, err = input.ServicePreparedDBCustomize(funcName, dbRedmineDev, nil, contextModel, input.doUpdateSprintAndPaymentStatusOnRedmine)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
