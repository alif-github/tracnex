package PKCEClientMappingService

import (
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
)

func GenerateI18Message(messageID string, language string) (output string) {
	return util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.PKCEClientMappingServiceBundle, messageID, language, nil)
}