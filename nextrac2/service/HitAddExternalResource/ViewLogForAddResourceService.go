package HitAddExternalResource

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
)

func (input hitAddExternalResourceService) ViewLogForAddResource(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.AddResourceExternalRequest

	inputStruct, err = input.readBodyAndValidateForView(request, contextModel, input.validateViewLogForAddResource)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewLogForAddResource(inputStruct)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_LOG_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hitAddExternalResourceService) doViewLogForAddResource(inputStructInterface interface{}) (output interface{}, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.AddResourceExternalRequest)

	viewModel := repository.ViewClientRegistrationLogModel {
		ClientID: 	sql.NullString{String: inputStruct.ClientID},
	}

	var resultView repository.ViewClientRegistrationLogModel
	resultView, err = dao.PKCEClientMappingDAO.ViewDetailLogRegistrationJoinUser(serverconfig.ServerAttribute.DBConnection, viewModel)
	if err.Error != nil {
		return
	}

	output = out.ViewDetailLogForAddResourceDTOOut{
		ClientID: 		resultView.ClientID.String,
		AuthUserID: 	resultView.AuthUserID.Int64,
		ClientType: 	resultView.ClientType.String,
		Status: 		resultView.Status.String,
		FirstName: 		resultView.FirstName.String,
		LastName: 		resultView.LastName.String,
		Resource: 		resultView.Resource.String,
		UpdatedAt: 		resultView.UpdatedAt.String,
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input hitAddExternalResourceService) validateViewLogForAddResource(inputStruct *in.AddResourceExternalRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewLogForAddResource()
}