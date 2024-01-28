package ResendOTPService

import (
	"net/url"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"strconv"
)

func (input resendOTPService) linkQueryAuth(linkStr string, dataOnDB repository.UserRegistrationDetailMapping) (linkQuery string, err errorModel.ErrorModel) {
	var (
		queryUrlEmailLink url.Values
		urlEmailLink      *url.URL
	)

	queryUrlEmailLink, urlEmailLink, err = input.linkQuery(linkStr, dataOnDB)
	if err.Error != nil {
		return
	}

	urlEmailLink.RawQuery = queryUrlEmailLink.Encode()
	linkQuery = urlEmailLink.String()

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input resendOTPService) linkQueryEmail(linkStr string, dataOnDB repository.UserRegistrationDetailMapping, userVerifyModel repository.UserVerificationModel) (linkQuery string, err errorModel.ErrorModel) {
	var (
		queryUrlEmailLink url.Values
		urlEmailLink      *url.URL
	)

	queryUrlEmailLink, urlEmailLink, err = input.linkQuery(linkStr, dataOnDB)
	if err.Error != nil {
		return
	}

	queryUrlEmailLink.Set(constanta.ActivationCodeQueryParam, userVerifyModel.EmailCode.String)
	queryUrlEmailLink.Set(constanta.EmailQueryParam, dataOnDB.UserRegistrationDetail.Email.String)
	queryUrlEmailLink.Set(constanta.UserIDQueryParam, strconv.Itoa(int(dataOnDB.UserRegistrationDetail.AuthUserID.Int64)))
	urlEmailLink.RawQuery = queryUrlEmailLink.Encode()
	linkQuery = urlEmailLink.String()

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input resendOTPService) linkQueryGroChat(linkStr string, dataOnDB repository.UserRegistrationDetailMapping, userVerifyModel repository.UserVerificationModel) (linkQuery string, err errorModel.ErrorModel) {
	var (
		queryUrlEmailLink url.Values
		urlEmailLink      *url.URL
	)

	queryUrlEmailLink, urlEmailLink, err = input.linkQuery(linkStr, dataOnDB)
	if err.Error != nil {
		return
	}

	queryUrlEmailLink.Set(constanta.OTPQueryParam, userVerifyModel.PhoneCode.String)
	queryUrlEmailLink.Set(constanta.PhoneQueryParam, dataOnDB.UserRegistrationDetail.NoTelp.String)
	queryUrlEmailLink.Set(constanta.UserIDQueryParam, strconv.Itoa(int(dataOnDB.UserRegistrationDetail.AuthUserID.Int64)))
	urlEmailLink.RawQuery = queryUrlEmailLink.Encode()
	linkQuery = urlEmailLink.String()

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input resendOTPService) linkQuery(linkStr string, dataOnDB repository.UserRegistrationDetailMapping) (queryUrlEmailLink url.Values, urlEmailLink *url.URL, err errorModel.ErrorModel) {
	var (
		fileName = "ResendOTPServiceUtil.go"
		funcName = "linkQuery"
		errS     error
	)

	urlEmailLink, errS = url.Parse(linkStr)
	if errS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errS)
		return
	}

	queryUrlEmailLink = urlEmailLink.Query()

	//--- Client ID Hold
	//queryUrlEmailLink.Set(constanta.ClientIDQueryParam, dataOnDB.PKCEClientMapping.ClientID.String)
	queryUrlEmailLink.Set(constanta.UniqueID1QueryParam, dataOnDB.UserRegistrationDetail.UniqueID1.String)
	queryUrlEmailLink.Set(constanta.UniqueID2QueryParam, dataOnDB.UserRegistrationDetail.UniqueID2.String)
	queryUrlEmailLink.Set(constanta.SalesmanIDQueryParam, dataOnDB.UserRegistrationDetail.SalesmanID.String)
	queryUrlEmailLink.Set(constanta.UserQueryParam, dataOnDB.UserRegistrationDetail.UserID.String)

	err = errorModel.GenerateNonErrorModel()
	return
}
