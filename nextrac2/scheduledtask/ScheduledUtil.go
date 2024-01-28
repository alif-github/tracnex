package scheduledtask

import (
"nexsoft.co.id/nextrac2/config"
"nexsoft.co.id/nextrac2/model/applicationModel"
"nexsoft.co.id/nextrac2/model/errorModel"
)

func SetLogger(fileName string, err errorModel.ErrorModel) (loggerModel applicationModel.LoggerModel) {
	loggerModel = applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())

	if err.Error != nil {
		loggerModel.Status = err.Code
		loggerModel.Code = err.Error.Error()
		loggerModel.Class = err.FileName
	} else {
		loggerModel.Status = err.Code
		loggerModel.Class = fileName
	}

	return
}
