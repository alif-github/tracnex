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
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service/NexmileParameterService"
	"nexsoft.co.id/nextrac2/test"
	"testing"
)

type scenarioTableGetParameterNexmile struct {
	test.AbstractScenario
	contextModel applicationModel.ContextModel
	ExpPayload   interface{}
}

func createScenarioTableGetParameterNexmile() []scenarioTableGetParameterNexmile {
	fileName := "ViewNexmileParameterService.go"
	funcName := "doViewNexmileParameter"
	return []scenarioTableGetParameterNexmile{
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Client ID Invalid",
				RequestBody: in.NexmileParameterRequestForView{
					ClientID:     "12345678",
					AuthUserID:   123,
					UserId:       "1",
					AndroidId:    "1",
					ClientTypeId: 5,
				},
				Expected: errorModel.GenerateClientValidationError(fileName, funcName),
			},
			contextModel: contextModel,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Auth Server ID Not Found",
				RequestBody: in.NexmileParameterRequestForView{
					ClientID:     "VALID1",
					AuthUserID:   778,
					UserId:       "1",
					AndroidId:    "1",
					ClientTypeId: 5,
				},
				Expected: errorModel.GenerateUnknownAuthUserId(fileName, funcName),
			},
			contextModel: contextModel,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "User ID Different With User Registration Detail",
				RequestBody: in.NexmileParameterRequestForView{
					ClientID:     "VALID1",
					AuthUserID:   123,
					UserId:       "6",
					AndroidId:    "1",
					ClientTypeId: 5,
				},
				Expected: errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ParameterID),
			},
			contextModel: contextModel,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Success Get Nexmile Parameter, Error Status",
				RequestBody: in.NexmileParameterRequestForView{
					ClientID:     "VALID1",
					AuthUserID:   123,
					UserId:       "1",
					AndroidId:    "1",
					ClientTypeId: 5,
				},
				Expected: errorModel.GenerateNonErrorModel(),
			},
			contextModel: contextModel,
			ExpPayload: out.PayloadData{
				Content: out.ViewNexmileParameterResponse{
					UniqueId1:        "123",
					UniqueId2:        "123",
					CompanyName:      "PT. Makmur Sejahtera",
					PasswordAdmin:    "abc123",
					UserAdmin:        "user01",
					Parameters: []out.ParameterValueResponse{
						{
							ParameterID:    "1",
							ParameterValue: "1",
						},
					},
				}},
		},
	}
}

func TestGetNexmileParameter(t *testing.T) {
	useCaseTest := createScenarioTableGetParameterNexmile()
	var errMessage string

	for _, useCase := range useCaseTest {
		t.Run(useCase.Name, func(t *testing.T) {
			request, errS := http.NewRequest(http.MethodPost, "/v1/nextrac/nexmile-parameter", nil)

			if errS != nil || request == nil {
				t.Log(errS)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyByte, _ := json.Marshal(useCase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyByte)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := NexmileParameterService.NexmileParameterService.ViewNexmileParameter(request, &useCase.contextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, useCase.Expected, err, errMessage)
			if useCase.ExpPayload != nil {
				contentPayloadExpected := useCase.ExpPayload.(out.PayloadData).Content.(out.ViewNexmileParameterResponse)
				contentPayloadActual := output.Data.Content.(out.ViewNexmileParameterResponse)

				assert.Equal(t, contentPayloadExpected.UniqueId1, contentPayloadActual.UniqueId1, "Unique ID 1 Failed")
				assert.Equal(t, contentPayloadExpected.CompanyName, contentPayloadActual.CompanyName, "Company Name Failed")
				assert.Equal(t, contentPayloadExpected.PasswordAdmin, contentPayloadActual.PasswordAdmin, "Password Admin Failed")
				assert.Equal(t, contentPayloadExpected.UniqueId2, contentPayloadActual.UniqueId2, "Unique ID 2 Failed")
				assert.Equal(t, contentPayloadExpected.UserAdmin, contentPayloadActual.UserAdmin, "User Admin Failed")
				assert.Equal(t, len(contentPayloadExpected.Parameters), len(contentPayloadActual.Parameters), "Amount Parameter Different")
			}
		})
	}
}
