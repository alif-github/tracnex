package TaskSchedulerService

import (
	"database/sql"
	"math"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
)

func (input taskSchedulerService) splitThread(totalData int64) (result int) {
	result = 5
	totalPage := int(math.Ceil(float64(totalData) / float64(constanta.TotalDataProductLicensePerChannel)))
	if totalPage < result {
		result = totalPage
	}
	return
}

func GenerateI18Message(messageID string, language string) (output string) {
	return util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.TaskSchedulerBundle, messageID, language, nil)
}

type responseDataForThread struct {
	Data interface{}
	Err  errorModel.ErrorModel
}

type requestDataForThread struct {
	Data interface{}
	Tx   *sql.Tx
}
