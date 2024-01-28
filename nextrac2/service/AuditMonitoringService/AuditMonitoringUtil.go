package AuditMonitoringService

import (
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
	"strings"
)

func GenerateI18NMessage(messageID string, language string) (output string) {
	return util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.AuditMonitoringBundle, messageID, language, nil)
}

func CensoringSecretData(tableName string, data *map[string]interface{}) {
	var tempData = *data
	for i := 0; i < len(config.ApplicationConfiguration.GetAudit().ListSecretData); i++ {
		splitCensoredData := strings.Split(strings.Trim(config.ApplicationConfiguration.GetAudit().ListSecretData[i], " "), ".")
		if len(splitCensoredData) > 1 {
			if tableName == splitCensoredData[0] {
				censoringData(&tempData, splitCensoredData[1])
			}
		} else {
			censoringData(&tempData, splitCensoredData[0])
		}
	}
	data = &tempData
}

func censoringData(data *map[string]interface{}, key string) {
	var tempData = *data

	if tempData[key] != nil {
		tempData[key] = "***********"
	}

	data = &tempData
}

func getMenuCodeConstanta(key string) string {
	listTableName := make(map[string]string)

	listTableName["nexsoft.profile-setting"] = "user"

	return listTableName[key]
}
