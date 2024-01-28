package UserActivationService

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type userActivationService struct {
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var UserActivationService = userActivationService{}.New()

func (input userActivationService) New() (output userActivationService) {
	output.FileName = "UserActivationService.go"
	output.ValidLimit = []int{10, 20, 50, 100}
	output.ValidOrderBy = []string{"id"}
	return
}

func (input userActivationService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(userActivation *in.UserActivationRequest) errorModel.ErrorModel) (inputStruct in.UserActivationRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	code, _ := mux.Vars(request)["CODE"]
	email, _ := mux.Vars(request)["EMAIL"]
	username, _ := mux.Vars(request)["USERNAME"]

	if inputStruct.UserID == 0 {
		inputStruct.UserID = int64(id)
	}

	if inputStruct.EmailCode == "" {
		inputStruct.EmailCode = code
	}

	if inputStruct.Email == "" {
		inputStruct.Email = email
	}

	if inputStruct.Username == "" {
		inputStruct.Username = username
	}

	err = validation(&inputStruct)
	return
}

func (input userActivationService) HitActivationToAuth(inputStruct in.UserActivationRequest, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	funcName := "HitActivationEmailToAuth"
	var payloadMessage authentication_response.AuthenticationErrorResponse

	internalToken := resource_common_service.GenerateInternalToken("auth", 0, contextModel.AuthAccessTokenModel.ClientID, config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
	authConfig := config.ApplicationConfiguration.GetAuthenticationServer()
	authActivationUrl := authConfig.Host + authConfig.PathRedirect.InternalUser.Activation.Email

	header := make(map[string][]string)
	header[common.AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, errorS := common.HitAPI(authActivationUrl, header, util.StructToJSON(inputStruct), "POST", *contextModel)

	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &payloadMessage)

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(payloadMessage.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, payloadMessage.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}
