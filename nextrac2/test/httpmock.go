package test

import (
	"github.com/jarcoal/httpmock"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"strconv"
)

type MockBodyAndResponseData struct {
	Body     interface{}
	Response interface{}
}

func SetMockAuthServerGetCredential(clientID string, clientSecret string, signKey string, alias string) {
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	httpmock.RegisterResponder("GET", authenticationServer.Host+authenticationServer.PathRedirect.InternalClient.CrudClient+"/"+clientID,
		func(_ *http.Request) (*http.Response, error) {
			resp, errMock := httpmock.NewJsonResponse(200, map[string]interface{}{
				"nexsoft": map[string]interface{}{
					"payload": map[string]interface{}{
						"data": map[string]interface{}{
							"content": map[string]interface{}{
								"client_id":     clientID,
								"client_secret": clientSecret,
								"signature_key": signKey,
								"alias_name":    alias,
							},
						},
					},
				},
			})
			if errMock != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})
}

func SetHttpMockResponseWithRequest(method, pathUrl string, codeStatus int, mockData []MockBodyAndResponseData, responseError interface{})  {
	httpmock.RegisterResponder(method, pathUrl,
		func(req *http.Request) (resp *http.Response, errMock error) {
			var stringBody, strExpectedBody string
			var errorS error

			stringBody, _, errorS = util.ReadBody(req)
			if errorS != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}

			for _, data := range mockData {
				strExpectedBody = util.StructToJSON(data.Body)
				if strExpectedBody == stringBody {
					resp, errMock = httpmock.NewJsonResponse(codeStatus, data.Response)
					if errMock != nil {
						return httpmock.NewStringResponse(500, ""), nil
					}
					return resp, nil
				}
			}

			resp, errMock = httpmock.NewJsonResponse(codeStatus, responseError)
			if errMock != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})
}

func SetHttpMockResponse(method, pathUrl string, codeStatus int, bodyResponse interface{})  {
	httpmock.RegisterResponder(method, pathUrl,
		func(_ *http.Request) (*http.Response, error) {
			resp, errMock := httpmock.NewJsonResponse(codeStatus, bodyResponse)
			if errMock != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})
}

func SetHttpMockWithErrorResponse(method, pathUrl string, successCode, errorCode int, successResponse, errorResponse interface{})  {
	httpmock.RegisterResponder(method, pathUrl,
		func(_ *http.Request) (*http.Response, error) {
			resp, errMock := httpmock.NewJsonResponse(successCode, successResponse)
			if errMock != nil {
				return httpmock.NewStringResponse(errorCode, util.StructToJSON(errorResponse)), nil
			}
			return resp, nil
		})
}

func SetMockAuthServerGetByAuthUserID(authUserID int, clientID string) {
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	httpmock.RegisterResponder("GET", authenticationServer.Host + authenticationServer.PathRedirect.InternalUser.CrudUser+"/"+strconv.Itoa(authUserID),
		func(_ *http.Request) (*http.Response, error) {
			resp, errMock := httpmock.NewJsonResponse(200, map[string]interface{}{
				"nexsoft": map[string]interface{}{
					"payload": map[string]interface{}{
						"data": map[string]interface{}{
							"content": map[string]interface{}{
								"client_id":     clientID,
							},
						},
					},
				},
			})
			if errMock != nil {
				return httpmock.NewStringResponse(500, ""), nil
			}
			return resp, nil
		})
}