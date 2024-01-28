package unregister_user

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service/RegistrationNamedUserService"
	"nexsoft.co.id/nextrac2/test"
	"strconv"
	"testing"
)

type scenarioTableUnregisterUser struct {
	test.AbstractScenario
	contextModel applicationModel.ContextModel
}

func createScenarioUnregisterUser() []scenarioTableUnregisterUser {
	return []scenarioTableUnregisterUser{
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Skenario 1 Negative Case - PT. Makmur Sejahtera gagal melakukan unregister user nexmile",
				RequestBody: in.UnregisterNamedUserRequest{
					ID: 1,
				},
				Expected: errorModel.GenerateForbiddenClientCredentialAccess("UnregisterNamedUserService.go", "doUnregisterNamedUser"),
			},
			contextModel: applicationModel.ContextModel{
				AuthAccessTokenModel: model.AuthAccessTokenModel{
					RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
						ResourceUserID: 12,
					},
					ClientID: "TEST123",
					Locale:   constanta.IndonesianLanguage,
				},
			},
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Skenario 2 Positive Case - PT. Makmur Sejahtera sukses melakukan unregister user nexmile",
				RequestBody: in.UnregisterNamedUserRequest{
					ID: 1,
				},
				Expected: errorModel.GenerateNonErrorModel(),
			},
			contextModel: contextModelNexMile,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Skenario 3 Negative Case - PT. Maju Jaya gagal melakukan unregister user nechief mobile",
				RequestBody: in.UnregisterNamedUserRequest{
					ID: 2,
				},
				Expected: errorModel.GenerateForbiddenClientCredentialAccess("UnregisterNamedUserService.go", "doUnregisterNamedUser"),
			},
			contextModel: applicationModel.ContextModel{
				AuthAccessTokenModel: model.AuthAccessTokenModel{
					RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
						ResourceUserID: 12,
					},
					ClientID: "TEST123",
					Locale:   constanta.IndonesianLanguage,
				},
			},
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Skenario 4 Positive Case - PT. Maju Jaya sukses melakukan unregister user nexmile",
				RequestBody: in.UnregisterNamedUserRequest{
					ID: 2,
				},
				Expected: errorModel.GenerateNonErrorModel(),
			},
			contextModel: contextModelNexChiefMobile,
		},
	}
}

func TestUnregisterUser(t *testing.T) {
	usecases := createScenarioUnregisterUser()
	var errMessage string

	for _, usecase := range usecases {
		t.Run(usecase.Name, func(t *testing.T) {
			request, errs := http.NewRequest(http.MethodPost, "v1/nextrac/user-registration/unregister", nil)

			if errs != nil || request == nil {
				t.Log(errs)
				assert.FailNow(t, "Error build a request")
			}

			request = mux.SetURLVars(request, map[string]string{
				"id" : strconv.Itoa(int(usecase.RequestBody.(in.UnregisterNamedUserRequest).ID)),
			})

			reqBodyBtye, _ := json.Marshal(usecase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyBtye)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := RegistrationNamedUserService.RegistrationNamedUserService.UnregisterNamedUser(request, &usecase.contextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, usecase.Expected, err, errMessage)
		})
	}
}
