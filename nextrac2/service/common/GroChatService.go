package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/grochat_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/grochat_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"strconv"
)

type groChatService struct {
	service.AbstractService
}

var GroChatService = groChatService{}.New()

func (input groChatService) New() (output groChatService) {
	output.FileName = "GroChatService.go"
	return
}

func (input groChatService) SendGroChatMessage(groChatRequest grochat_request.SendMessageGroChatRequest, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		fileName       = input.FileName
		funcName       = "SendGroChatMessage"
		groChat        = config.ApplicationConfiguration.GetGrochat()
		url            = groChat.Host + groChat.PathRedirect.SendMessage
		structResponse grochat_response.SendMessageGroChatResponse
	)

	internalToken := resource_common_service.GenerateInternalToken(constanta.GroChatResourceID, config.ApplicationConfiguration.GetClientCredentialsAuthUserID(), config.ApplicationConfiguration.GetClientCredentialsClientID(), config.ApplicationConfiguration.GetServerResourceID(), constanta.DefaultApplicationsLanguage)

	//--- For testing only, delete if unused
	//internalToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJsb2NhbGUiOiJpZC1JRCIsImNpZCI6IjUyZWI5ZjQyZmIxMDQwMjRiZDVkMDAxM2YwM2VjY2E0IiwicmVzb3VyY2UiOiJhdXRoIiwidmVyc2lvbiI6IjIuMC4wIiwidXNlcl9jbGllbnQiOiI1MmViOWY0MmZiMTA0MDI0YmQ1ZDAwMTNmMDNlY2NhNCIsImV4cCI6MTY3Nzk4NjE2OSwiaWF0IjoxNjc3NzI2OTY5LCJpc3MiOiJhdXRoIiwic3ViIjoiMzUxIn0.0B3fmptjjT4vxQf-mnL0IIEffxJlYfRTuXmqanxFZnOoEgHdnoGerFS_19ObhE5yEf681iR10d3LRTNGSsNcXw"

	statusCode, bodyResult, errs := common.HitSendMessageGroChatServer(internalToken, url, groChatRequest, contextModel)
	fmt.Println("Request Body API : ", groChatRequest)
	fmt.Println("Int Token : ", internalToken)
	if errs != nil {
		fmt.Println("Error Hit Send Message Grochat Server : ", err.Error.Error())
		err = errorModel.GenerateErrorModel(statusCode, errs.Error(), fileName, funcName, errs)
		return
	}
	fmt.Println("Result Hit HitSendMessageGroChatServer : ", bodyResult, "with status code")

	_ = json.Unmarshal([]byte(bodyResult), &structResponse)

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(structResponse.Message)
		err = errorModel.GenerateAuthenticationServerError(input.FileName, funcName, statusCode, strconv.Itoa(statusCode), causedBy)
		fmt.Println("Error Hit Send Message Grochat Server2 : ", err.Error.Error())
		return
	}

	return
}

func (input groChatService) GetGroChatMessage(inputStruct in.RegisterNamedUserRequest, OTP string, link string, purpose string) string {
	var (
		commonBundle = serverconfig.ServerAttribute.CommonServiceBundle
		param        = make(map[string]interface{})
		locale       = constanta.DefaultApplicationsLanguage
	)

	param[constanta.NameTableParam] = inputStruct.Firstname
	param[constanta.PurposeTableParam] = purpose

	switch inputStruct.ClientTypeID {
	case constanta.ResourceNexmileID:
		param[constanta.ClientTypeTableParam] = constanta.Nexmile
	case constanta.ResourceNexstarID:
		param[constanta.ClientTypeTableParam] = constanta.Nexstar
	default:
		param[constanta.ClientTypeTableParam] = " - "
	}

	param[constanta.CompanyIDTableParam] = inputStruct.UniqueID1
	param[constanta.CompanyNameTableParam] = inputStruct.CompanyName
	param[constanta.BranchIDTableParam] = inputStruct.UniqueID2
	param[constanta.BranchNameTableParam] = inputStruct.BranchName
	param[constanta.SalesmanIDTableParam] = inputStruct.SalesmanID
	param[constanta.UserTableParam] = inputStruct.UserID
	param[constanta.PasswordTableParam] = inputStruct.Password
	param[constanta.OTPTableParam] = OTP
	param[constanta.EmailTableParam] = inputStruct.Email
	param[constanta.ClientIDTableParam] = inputStruct.ClientID
	param[constanta.RegistrationIDTableParam] = inputStruct.AuthUserID
	param[constanta.LinkTableParam] = link

	return util.GenerateI18NServiceMessage(commonBundle, "OTP_MESSAGE_GROCHAT2", locale, param)
}
