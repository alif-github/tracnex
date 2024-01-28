package UserInvitationService

import (
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
)

type invitationService struct {
	service.AbstractService
	FileName string
}

var InvitationService = invitationService{}.New()

func (input invitationService) New() (output invitationService) {
	output.FileName = "UserInvitationService.go"
	return
}

func (input invitationService) readBody(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.UserInvitationRequest, errModel errorModel.ErrorModel) {
	stringBody, errModel := input.ReadBody(request, contextModel)
	if errModel.Error != nil {
		return
	}

	if err := json.Unmarshal([]byte(stringBody), &inputStruct); err != nil {
		errModel = errorModel.GenerateInvalidRequestError(input.FileName, "readBody", err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}
