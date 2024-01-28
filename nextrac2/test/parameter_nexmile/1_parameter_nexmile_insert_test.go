package parameter_nexmile

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
	"nexsoft.co.id/nextrac2/service/NexmileParameterService"
	"nexsoft.co.id/nextrac2/test"
	"testing"
)

type scenarioTableInsertParameterNexmile struct {
	test.AbstractScenario
	contextModel     applicationModel.ContextModel
}

func createScenarioTableInsertParameterNexmile() []scenarioTableInsertParameterNexmile {
	fileName := "InsertNexmileParameterService.go"
	funcName := "doInsertNexmileParameter"
	return []scenarioTableInsertParameterNexmile{
		{
			AbstractScenario: test.AbstractScenario{
				Name:        "PKCE Client Mapping Not Found",
				RequestBody: in.NexmileParameterRequest{
					ClientID:      "12345678",
					ClientTypeID:  constanta.ResourceTestingNexmileID,
					UniqueID1:     "123",
					UniqueID2:     "123",
					ParameterData: []in.ParameterValue{
						{
							ParameterID: "1",
							Value:       "1",
						},
					},
				},
				Expected: errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientMappingClientID),
			},
			contextModel: contextModel,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name:        "Forbidden access Client",
				RequestBody: in.NexmileParameterRequest{
					ClientID:      "VALID1",
					ClientTypeID:  constanta.ResourceTestingNexmileID,
					UniqueID1:     "123",
					UniqueID2:     "123",
					ParameterData: []in.ParameterValue{
						{
							ParameterID: "1",
							Value:       "1",
						},
					},
				},
				Expected: errorModel.GenerateForbiddenAccessClientError(fileName, funcName),
			},
			contextModel: applicationModel.ContextModel{
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
			AbstractScenario: test.AbstractScenario{
				Name:        "Success Insert Nexmile Parameter",
				RequestBody: in.NexmileParameterRequest{
					ClientID:      "VALID1",
					ClientTypeID:  constanta.ResourceTestingNexmileID,
					UniqueID1:     "123",
					UniqueID2:     "123",
					ParameterData: []in.ParameterValue{
						{
							ParameterID: "1",
							Value:       "1",
						},
					},
				},
				Expected:    errorModel.GenerateNonErrorModel(),

			},
			contextModel: contextModel,
		},
	}
}

func TestInsertNexmileParameter(t *testing.T) {
	useCaseTest := createScenarioTableInsertParameterNexmile()
	var errMessage string

	for _, useCase := range useCaseTest {
		t.Run(useCase.Name, func(t *testing.T) {
			request, errS := http.NewRequest(http.MethodPost, "/v1/nextrac/nexmile-parameter/add", nil)

			if errS != nil || request == nil {
				t.Log(errS)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyByte, _ := json.Marshal(useCase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyByte)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := NexmileParameterService.NexmileParameterService.InsertNexmileParameter(request, &useCase.contextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, useCase.Expected, err, errMessage)
		})
	}
}
