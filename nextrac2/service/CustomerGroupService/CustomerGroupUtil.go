package CustomerGroupService

import (
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
)

func GenerateI18NMessage(messageID string, language string) (output string) {
	return util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.CustomerGroupServiceBundle, messageID, language, nil)
}