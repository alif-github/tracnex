package ProcessFileListCustomerService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/util"
	"strings"
)

type initiateImport struct {
	Message			map[string]string	`json:"message"`
	ListDataType	[]string			`json:"list_data_type"`
}

func (input importService) InitiateImportData(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	output.Data.Content = input.getInitiateImportMessage()
	output.Status = out.StatusResponse {
		Code: 		util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message:	GenerateI18NMessage("SUCCESS_INITIATE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input importService) getInitiateImportMessage() interface{} {
	message := make(map[string]string)
	var listDataType []string
	for key := range input.AvailableTypeImport {
		message[key] = 	messageEscape.Replace(input.AvailableTypeImport[key].Message)
		listDataType = 	append(listDataType, key)
	}

	return initiateImport {
		Message: 		message,
		ListDataType: 	listDataType,
	}
}

var messageEscape = strings.NewReplacer (
	"\t", "",
	"\n", "<br/>",
)