package common

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/grochat_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/nexcloud_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/redmine_request"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	util2 "nexsoft.co.id/nextrac2/util"
	"os"
	"strings"
)

type MultiPartData struct {
	Data   string
	IsFile bool
}

func HitAPI(urlAddress string, header map[string][]string, body string, method string, contextModel applicationModel.ContextModel) (statusCode int, headerResult map[string][]string, bodyResult string, err error) {
	var (
		reqURL   *url.URL
		request  *http.Request
		response *http.Response
	)

	contextModel.LoggerModel.Class = "[HitAPI.go,HitAPI]"
	reqURL, err = url.Parse(urlAddress)
	if err != nil {
		contextModel.LoggerModel.Status = 400
		contextModel.LoggerModel.Message = "Redirect URL: " + urlAddress + " not Valid"
		util.LogError(contextModel.LoggerModel.ToLoggerObject())
		return
	}

	if header == nil {
		header = make(map[string][]string)
	}

	header[RequestIDConstanta] = []string{contextModel.LoggerModel.RequestID}
	header[SourceConstanta] = []string{contextModel.LoggerModel.Application}

	request = &http.Request{
		Method: strings.ToUpper(method),
		URL:    reqURL,
		Header: header,
		Body:   ioutil.NopCloser(strings.NewReader(body)),
	}

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		contextModel.LoggerModel.Status = 500
		contextModel.LoggerModel.Message = contextModel.LoggerModel.Message + " caused by : " + err.Error()
		util.LogError(contextModel.LoggerModel.ToLoggerObject())
		return
	}

	defer func() {
		err = response.Body.Close()
	}()

	bodyResultByte, _ := ioutil.ReadAll(response.Body)
	bodyResult = string(bodyResultByte)
	statusCode = response.StatusCode
	headerResult = response.Header

	contextModel.LoggerModel.Status = response.StatusCode
	contextModel.LoggerModel.Message = "Success Hit " + urlAddress + ", with Response " + bodyResult

	//util.LogInfo(contextModel.LoggerModel.ToLoggerObject())
	return
}

func HitAPINoHeader(urlAddress string, header map[string][]string, body string, method string, contextModel applicationModel.ContextModel) (statusCode int, headerResult map[string][]string, bodyResult string, err error) {
	var (
		reqURL   *url.URL
		request  *http.Request
		response *http.Response
	)

	contextModel.LoggerModel.Class = "[HitAPI.go,HitAPI]"
	reqURL, err = url.Parse(urlAddress)
	if err != nil {
		contextModel.LoggerModel.Status = 400
		contextModel.LoggerModel.Message = "Redirect URL: " + urlAddress + " not Valid"
		util.LogError(contextModel.LoggerModel.ToLoggerObject())
		return
	}

	request = &http.Request{
		Method: strings.ToUpper(method),
		URL:    reqURL,
		Header: header,
		Body:   ioutil.NopCloser(strings.NewReader(body)),
	}

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		contextModel.LoggerModel.Status = 500
		contextModel.LoggerModel.Message = contextModel.LoggerModel.Message + " caused by : " + err.Error()
		util.LogError(contextModel.LoggerModel.ToLoggerObject())
		return
	}

	defer func() {
		err = response.Body.Close()
	}()

	bodyResultByte, _ := ioutil.ReadAll(response.Body)
	bodyResult = string(bodyResultByte)
	statusCode = response.StatusCode
	headerResult = response.Header

	contextModel.LoggerModel.Status = response.StatusCode
	contextModel.LoggerModel.Message = "Success Hit " + urlAddress + ", with Response " + bodyResult

	return
}

func HitCheckTokenURL(checkTokenEndpoint string, token string, resourceID string, scope string, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{token}
	body := model.CheckURLTokenBody{
		ResourceID: resourceID,
		Scope:      scope,
	}

	statusCode, _, bodyResult, err = HitAPI(checkTokenEndpoint, header, util.StructToJSON(body), "POST", *contextModel)
	return
}

func HitAddClientResource(internalToken string, clientResourceEndpoint string, clientID string, resourceID string, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	body := model.AddClientResourceBody{
		ResourceID: resourceID,
		ClientID:   clientID,
	}

	statusCode, _, bodyResult, err = HitAPI(clientResourceEndpoint, header, util.StructToJSON(body), "POST", *contextModel)
	return
}

func HitRegisterUserAuthenticationServer(internalToken string, endpoint string, userDTO authentication_request.UserAuthenticationDTO, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	var (
		header   = make(map[string][]string)
		bodyByte []byte
		bodyReq  string
	)

	header[AuthorizationHeaderConstanta] = []string{internalToken}
	bodyByte, _ = util2.JSONMarshal(userDTO)
	bodyReq = string(bodyByte)
	statusCode, _, bodyResult, err = HitAPI(endpoint, header, bodyReq, "POST", *contextModel)
	return
}

func HitGetListUserAuthenticationServer(internalToken string, endpoint string, userDTO authentication_request.GetListUserDTO, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, util.StructToJSON(userDTO), "GET", *contextModel)
	return
}

func HitForgetPasswordAuthenticationServer(internalToken string, endpoint string, userDTO authentication_request.ForgetPasswordDTOin, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, util.StructToJSON(userDTO), "POST", *contextModel)
	return
}

func HitInitiateInsertUserAuthenticationServer(internalToken string, endpoint string, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, "", "GET", *contextModel)
	return
}

func HitCheckUserAuthenticationServer(internalToken string, endpoint string, userDTO authentication_request.UserAuthenticationDTO, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, util.StructToJSON(userDTO), "POST", *contextModel)
	return
}

func HitAddResourceClientAuthenticationServer(internalToken string, endpoint string, addResourceDTO authentication_request.AddResourceClient, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, util.StructToJSON(addResourceDTO), "POST", *contextModel)
	return
}

func HitRegisterClientAuthenticationServer(internalToken string, endpoint string, clientDTO authentication_request.ClientAuthentication, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, util.StructToJSON(clientDTO), "POST", *contextModel)
	return
}

func HitCheckClientUserAuthenticationServer(internalToken string, endpoint string, checkClientUserDTO authentication_request.CheckClientOrUser, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, util.StructToJSON(checkClientUserDTO), "POST", *contextModel)
	return
}

func HitUnregisterClientResourceNexcloud(internalToken string, endpoint string, unregisterClientDTO nexcloud_request.UnregisterClient, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, util.StructToJSON(unregisterClientDTO), "DELETE", *contextModel)
	return
}

func HitUpdateClientAuthenticationServer(internalToken string, endpoint string, clientDTO authentication_request.ClientUpdateRequest, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, util.StructToJSON(clientDTO), "PUT", *contextModel)
	return
}

func HitGetDetailClientAuthenticationServer(internalToken string, endpoint string, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, "", "GET", *contextModel)
	return
}

func HitChangePasswordAuthenticationServer(internalToken string, endpoint string, userDTO authentication_request.ChangePasswordDTOin, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, util.StructToJSON(userDTO), "POST", *contextModel)
	return
}

func HitPaidPaymentRedmineServer(accessKey string, endpoint string, issueDTO redmine_request.IssuePaidRedmineRequest, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[RedmineAPIKeyConstanta] = []string{accessKey}
	header[ContentTypeHeaderConstanta] = []string{constanta.ApplicationJSON}

	statusCode, _, bodyResult, err = HitAPINoHeader(endpoint, header, util.StructToJSON(issueDTO), "PUT", *contextModel)
	return
}

func HitGetClientAuthenticationServer(internalToken string, endpoint string, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	header := make(map[string][]string)
	header[AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, "", "GET", *contextModel)
	return
}

func HitMultipartFormDataRequest(urlAddress, method string, header map[string][]string, body map[string]MultiPartData, contextModel applicationModel.ContextModel) (statusCode int, headerResult map[string][]string, bodyResult string, err error) {
	var (
		bodyBuf         = &bytes.Buffer{}
		bodyWriter      = multipart.NewWriter(bodyBuf)
		arrString       []string
		part            io.Writer
		file            *os.File
		urlFile, reqURL *url.URL
		response        *http.Response
		request         *http.Request
	)

	contextModel.LoggerModel.Class = "[HitAPI.go,HitMultipartFormDataRequest]"

	reqURL, err = url.Parse(urlAddress)
	if err != nil {
		contextModel.LoggerModel.Status = 400
		contextModel.LoggerModel.Message = "Redirect URL: " + urlAddress + " not Valid"
		util.LogError(contextModel.LoggerModel.ToLoggerObject())
		return
	}

	if header == nil {
		header = make(map[string][]string)
	}

	header[RequestIDConstanta] = []string{contextModel.LoggerModel.RequestID}
	header[SourceConstanta] = []string{contextModel.LoggerModel.Application}

	defer func() {
		file.Close()
	}()

	for key, val := range body {
		if val.IsFile {
			urlFile, err = url.ParseRequestURI(val.Data)
			if err != nil {
				arrString = strings.Split(val.Data, "/")
				part, err = bodyWriter.CreateFormFile(key, arrString[len(arrString)-1])
				if err != nil {
					return
				}

				file, err = os.Open(val.Data)
				if err != nil {
					return
				}

				_, err = io.Copy(part, file)
				if err != nil {
					log.Println("err io copy")
					return
				}
			} else {
				err = bodyWriter.WriteField(key, urlFile.String())
				if err != nil {
					return
				}
			}
		} else {
			err = bodyWriter.WriteField(key, val.Data)
			if err != nil {
				return
			}
		}
	}
	err = bodyWriter.Close()
	if err != nil {
		return
	}

	request = &http.Request{
		Method: strings.ToUpper(method),
		URL:    reqURL,
		Header: header,
		Body:   ioutil.NopCloser(bodyBuf),
	}
	request.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		contextModel.LoggerModel.Status = 500
		contextModel.LoggerModel.Message = contextModel.LoggerModel.Message + " caused by : " + err.Error()
		util.LogError(contextModel.LoggerModel.ToLoggerObject())
		return
	}

	defer func() {
		err = response.Body.Close()
		return
	}()

	bodyResultByte, _ := ioutil.ReadAll(response.Body)
	bodyResult = string(bodyResultByte)
	statusCode = response.StatusCode
	headerResult = response.Header

	contextModel.LoggerModel.Status = response.StatusCode
	contextModel.LoggerModel.Message = "Success Hit " + urlAddress + ", with Response " + bodyResult

	util.LogInfo(contextModel.LoggerModel.ToLoggerObject())

	return
}

func HitSendMessageGroChatServer(internalToken string, endpoint string, groChatDTO grochat_request.SendMessageGroChatRequest, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	var (
		header   = make(map[string][]string)
		bodyByte []byte
		bodyReq  string
	)

	header[AuthorizationHeaderConstanta] = []string{internalToken}
	header[ContentTypeHeaderConstanta] = []string{"application/json"}
	bodyByte, _ = util2.JSONMarshal(groChatDTO)
	bodyReq = string(bodyByte)

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, bodyReq, "POST", *contextModel)
	return
}

func HitResendAuthenticationServer(internalToken string, endpoint string, authDTO authentication_request.ResendUserVerificationRequest, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	var (
		header   = make(map[string][]string)
		bodyByte []byte
		bodyReq  string
	)

	header[AuthorizationHeaderConstanta] = []string{internalToken}
	bodyByte, _ = util2.JSONMarshal(authDTO)
	bodyReq = string(bodyByte)

	statusCode, _, bodyResult, err = HitAPI(endpoint, header, bodyReq, "PUT", *contextModel)
	return
}
