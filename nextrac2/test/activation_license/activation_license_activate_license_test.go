package activation_license

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service/ActivationLicenseService"
	"nexsoft.co.id/nextrac2/test"
	"testing"
)

type scenarioTableActivateLicense struct {
	test.AbstractScenario
	contextModel applicationModel.ContextModel
}

func createScenarioActivateTest() []scenarioTableActivateLicense {
	fileName := "ActivationLicenseService.go"
	return []scenarioTableActivateLicense{
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Forbidden access context Model",
				RequestBody: in.ActivationLicenseRequest{
					AbstractDTO:  in.AbstractDTO{},
					ClientID:     "123",
					ClientTypeID: 1,
					SignatureKey: "280d9968c4154d698362087a91a80e1a",
					Hwid:         "1",
					DetailClient: []in.UniqueIDClient{
						{
							UniqueID1: "NS6024050001031",
							UniqueID2: "1468381449586",
						},
					},
				},
				Expected:     errorModel.GenerateForbiddenAccessClientError(fileName, "doActivateLicense"),
			},
			contextModel: contextModelND,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Client credential tidak valid",
				RequestBody: in.ActivationLicenseRequest{
					AbstractDTO:  in.AbstractDTO{},
					ClientID:     "123",
					ClientTypeID: 1,
					SignatureKey: "280d9968c4154d698362087a91a80e1a",
					Hwid:         "1",
					DetailClient: []in.UniqueIDClient{
						{
							UniqueID1: "NS6024050001031",
							UniqueID2: "1468381449586",
						},
					},
				},
				Expected:     errorModel.GenerateUnknownDataError(fileName, "doActivateLicense", constanta.ClientID),
			},contextModel: applicationModel.ContextModel{
				AuthAccessTokenModel: model.AuthAccessTokenModel{
					RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
						ResourceUserID: 12,
					},
					ClientID: "123",
					Locale:   constanta.IndonesianLanguage,
				},
			},
		},
		{
			AbstractScenario : test.AbstractScenario{
				Name: "Data License Config Tidak Ditemukan",
				RequestBody: in.ActivationLicenseRequest{
					AbstractDTO:  in.AbstractDTO{},
					ClientID:     "08181c991e6b409eb016cfaa365b439d",
					ClientTypeID: 1,
					SignatureKey: "280d9968c4154d698362087a91a80e1a",
					Hwid:         "1",
					DetailClient: []in.UniqueIDClient{
						{
							UniqueID1: "NSS6024050001031",
							UniqueID2: "1468381449586",
						},
					},
				},
				Expected:     errorModel.GenerateActivationLicenseError(fileName, "doActivateLicense", constanta.ActivationLicenseGetLicenseConfigError, nil),
			},
			contextModel: contextModelND,
		},
		{
			AbstractScenario : test.AbstractScenario{
				Name: "Aktivasi License ND6 Berhasil",
				RequestBody: in.ActivationLicenseRequest{
					AbstractDTO:  in.AbstractDTO{},
					ClientID:     "08181c991e6b409eb016cfaa365b439d",
					ClientTypeID: 1,
					SignatureKey: "280d9968c4154d698362087a91a80e1a",
					Hwid:         "1",
					DetailClient: []in.UniqueIDClient{
						{
							UniqueID1: "NS6024050001031",
							UniqueID2: "1468381449586",
						},
					},
				},
				Expected:     errorModel.GenerateNonErrorModel(),
			},
			contextModel: contextModelND,
		},
		{
			AbstractScenario : test.AbstractScenario{
				Name: "Aktivasi License NexChief Berhasil",
				RequestBody: in.ActivationLicenseRequest{
					AbstractDTO:  in.AbstractDTO{},
					ClientID:     "1a2b12faf6a345759ccffc500d609b52",
					ClientTypeID: 4,
					SignatureKey: "3100d9968c4154d698362087a91a80e1a",
					Hwid:         "4",
					DetailClient: []in.UniqueIDClient{
						{
							UniqueID1: "NDI",
							UniqueID2: "",
						},
					},
				},
				Expected:     errorModel.GenerateNonErrorModel(),
			},
			contextModel: contextModelNexchief,
		},
	}
}

func TestActivateLicense(t *testing.T) {
	usecaseTest := createScenarioActivateTest()
	var errMessage string

	for _, usecase := range usecaseTest {
		t.Run(usecase.Name, func(t *testing.T) {
			request, errs := http.NewRequest(http.MethodPut, "/v1/nextrac/activation-license", nil)
			if errs != nil || request == nil {
				t.Log(errs)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyBtye, _:= json.Marshal(usecase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyBtye)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := ActivationLicenseService.ActivationLicenseService.ActivateLicense(request, &usecase.contextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, usecase.Expected, err, errMessage)
		})
	}
}