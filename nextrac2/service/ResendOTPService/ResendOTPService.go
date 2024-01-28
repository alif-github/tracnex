package ResendOTPService

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/grochat_request"
	"nexsoft.co.id/nextrac2/service"
	common2 "nexsoft.co.id/nextrac2/service/common"
	"strings"
)

type resendOTPService struct {
	service.AbstractService
	service.GetListData
}

var ResendOTPService = resendOTPService{}.New()

func (input resendOTPService) New() (output resendOTPService) {
	output.FileName = "ResendOTPService.go"
	output.ServiceName = constanta.ResendOTPService
	return
}

func (input resendOTPService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.UserVerificationRequest) errorModel.ErrorModel) (inputStruct in.UserVerificationRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidate"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	err = validation(&inputStruct)
	return
}

func (input resendOTPService) sendMessageToGrochat(dataOnDBUserRegDetail repository.UserRegistrationDetailMapping, userVerifyModel repository.UserVerificationModel, contextModel *applicationModel.ContextModel, linkChannelGroChat string) (err errorModel.ErrorModel) {
	var (
		otp             = userVerifyModel.PhoneCode.String
		groMessageParam in.RegisterNamedUserRequest
		groChatRequest  grochat_request.SendMessageGroChatRequest
	)

	//--- For testing only, delete if unused
	//dataOnDBUserRegDetail.UserRegistrationDetail.NoTelp.String = "+62-88888888882"

	dataOnDBUserRegDetail.UserRegistrationDetail.NoTelp.String = strings.ReplaceAll(dataOnDBUserRegDetail.UserRegistrationDetail.NoTelp.String, "-", "")
	groMessageParam = in.RegisterNamedUserRequest{
		ClientID:     dataOnDBUserRegDetail.PKCEClientMapping.ClientID.String,
		Firstname:    dataOnDBUserRegDetail.User.FirstName.String,
		UniqueID1:    dataOnDBUserRegDetail.UserRegistrationDetail.UniqueID1.String,
		CompanyName:  dataOnDBUserRegDetail.PKCEClientMapping.CompanyName.String,
		UniqueID2:    dataOnDBUserRegDetail.UserRegistrationDetail.UniqueID2.String,
		BranchName:   dataOnDBUserRegDetail.PKCEClientMapping.BranchName.String,
		SalesmanID:   dataOnDBUserRegDetail.UserRegistrationDetail.SalesmanID.String,
		UserID:       dataOnDBUserRegDetail.UserRegistrationDetail.UserID.String,
		Password:     dataOnDBUserRegDetail.UserRegistrationDetail.Password.String,
		ClientTypeID: dataOnDBUserRegDetail.PKCEClientMapping.ClientTypeID.Int64,
		NoTelp:       dataOnDBUserRegDetail.UserRegistrationDetail.NoTelp.String,
		Email:        dataOnDBUserRegDetail.UserRegistrationDetail.Email.String,
		AuthUserID:   dataOnDBUserRegDetail.UserRegistrationDetail.AuthUserID.Int64,
	}

	linkChannelGroChat, err = input.linkQueryGroChat(linkChannelGroChat, dataOnDBUserRegDetail, userVerifyModel)
	if err.Error != nil {
		return
	}

	groChatRequest = grochat_request.SendMessageGroChatRequest{
		Data: grochat_request.DetailData{
			PhoneNumber: groMessageParam.NoTelp,
			MessageContent: grochat_request.MessageContent{
				Message: common2.GroChatService.GetGroChatMessage(groMessageParam, otp, linkChannelGroChat, fmt.Sprintf(` (%s) `, constanta.Resend)),
			},
		},
	}

	//--- Get Default
	groChatRequest.GetDefault(&groChatRequest)
	return common2.GroChatService.SendGroChatMessage(groChatRequest, contextModel)
}
