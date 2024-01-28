package ForgetPassword

import (
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
)

func GenerateI18NMessage(messageID string, language string) (output string) {
	return util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ForgetPasswordBundle, messageID, language, nil)
}