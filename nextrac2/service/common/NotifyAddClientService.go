package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
)

type notifyAddClientService struct {
	service.AbstractService
}

var NotifyAddClientService = notifyAddClientService{}.New()

func (input notifyAddClientService) New() (output notifyAddClientService) {
	output.FileName = "NotifyAddClientService.go"
	return
}

func (input notifyAddClientService) StartService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.NotifyAddClientDTOIn
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	err = input.hitAPIAuthenticationServerForAddResourceClient(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateCommonServiceBundleI18NMessage("SUCCESS_ADD_CLIENT_RESOURCE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input notifyAddClientService) hitAPIAuthenticationServerForAddResourceClient(inputStruct in.NotifyAddClientDTOIn, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	funcName := "hitAPIAuthenticationServerForAddResourceClient"
	internalToken := resource_common_service.GenerateInternalToken("auth", contextModel.AuthAccessTokenModel.AuthenticationServerUserID, contextModel.AuthAccessTokenModel.ClientID, config.ApplicationConfiguration.GetServerResourceID(), constanta.DefaultApplicationsLanguage)
	contextModel.LoggerModel.Class = "[" + input.FileName + "," + "hitAPIAuthenticationServerForAddResourceClient" + "]"
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	addResourceClientUrl := authenticationServer.Host + authenticationServer.PathRedirect.AddResourceClient

	statusCode, bodyResult, errorS := common.HitAddClientResource(internalToken, addResourceClientUrl, inputStruct.ClientID, config.ApplicationConfiguration.GetServerResourceID(), contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return err
	}

	if statusCode == 200 {
		var addClientResourceResult authentication_response.AddClientAuthenticationResponse
		_ = json.Unmarshal([]byte(bodyResult), &addClientResourceResult)
		contextModel.LoggerModel.Status = statusCode

		//if inputStruct.RoleName != "" {
		//	//todox check is role exist
		//	//todox if not exist, return error
		//} else {
		//	//todox role default
		//}

		//content := addClientResourceResult.Nexsoft.Payload.Data.Content
		//todox ask BA
		//_ = repository.UserModel{
		//	AuthUserID:                sql.NullInt64{Int64: content.UserID},
		//	BtlClientID:               sql.NullInt64{},
		//	TeamID:                    sql.NullInt64{},
		//	IsAdmin:                   sql.NullBool{},
		//	UserCode:                  sql.NullString{},
		//	ClientID:                  sql.NullString{String: content.ClientID},
		//	SignatureKey:              sql.NullString{String: content.SignatureKey},
		//	Locale:                    sql.NullString{String: content.Locale},
		//	Level:                     sql.NullInt32{},
		//	Status:                    sql.NullBool{},
		//	RegistrationDate:          sql.NullTime{Time: time.Now()},
		//	SaturdayStartLoginTime:    sql.NullString{},
		//	SaturdayEndLoginTime:      sql.NullString{},
		//	SundayStartLoginTime:      sql.NullString{},
		//	SundayEndLoginTime:        sql.NullString{},
		//	WorkingDaysStartLoginTime: sql.NullString{},
		//	WorkingDaysEndLoginTime:   sql.NullString{},
		//	AccountValidThru:          sql.NullTime{},
		//	LastToken:                 sql.NullTime{},
		//	CreatedBy:                 sql.NullInt64{},
		//	CreatedAt:                 sql.NullTime{},
		//	UpdatedBy:                 sql.NullInt64{},
		//	UpdatedAt:                 sql.NullTime{},
		//	Deleted:                   sql.NullBool{},
		//}

		fmt.Println(addClientResourceResult)
		util.LogInfo(contextModel.LoggerModel.ToLoggerObject())
		//todox save client to db with role
	} else {
		err = common.ReadAuthServerError("hitAPIAuthenticationServerForAddResourceClient", statusCode, bodyResult, contextModel)
		return
	}

	return errorModel.GenerateNonErrorModel()
}
