package util

import (
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

//InputLog is user for usual logging that will save into main log file
//if there is an error it will save using logError method
func InputLog(err errorModel.ErrorModel, loggerModel applicationModel.LoggerModel) {
	if err.Error != nil {
		util.LogError(loggerModel.ToLoggerObject())
	} else {
		util.LogInfo(loggerModel.ToLoggerObject())
	}
}
