package penambahan_user

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
	"nexsoft.co.id/nextrac2/service/ActivationUserNexmileService"
	"nexsoft.co.id/nextrac2/test"
	"testing"
)

type scenarioAktivasiUserNexmile struct {
	test.AbstractScenario
	contextModel applicationModel.ContextModel
}

func createScenarioAktivasiUserNexmile() (result []scenarioAktivasiUserNexmile) {
	validEmail := "ValidEmail@email.com"
	validPhone := "+62-897789666889"
	fileName := "ActivationUserNexmileService.go"
	result = append(result, scenarioAktivasiUserNexmile{
		AbstractScenario: test.AbstractScenario{
			Name:        "Forbidden Client Access",
			RequestBody: in.ActivationUserNexmileRequest{
				UserRegistrationDetailID: userRegistrationDetailIDToBeActive,
				ParentClientID:           "123",
				FirstName:                "Pertama",
				LastName:                 "Terakhir",
				AndroidID:                "ValidAndroID",
				ClientID:                 "TestClientIDs",
				UserID:                   "100",
				AuthUserID:               100,
				AliasName:                "ValidAlias",
				Email:                    validEmail,
				Phone:                    validPhone,
			},
			Expected:    errorModel.GenerateForbiddenAccessClientError(fileName, "activateUserNexmile"),
		},
		contextModel:     applicationModel.ContextModel{
			AuthAccessTokenModel: model.AuthAccessTokenModel{
				RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
					ResourceUserID: 12,
				},
				ClientID: "08181c991e6b409eb016cfaa365b439d",
				Locale:   constanta.IndonesianLanguage,
			},
		},
	})

	result = append(result, scenarioAktivasiUserNexmile{
		AbstractScenario: test.AbstractScenario{
			Name:        "Get User Registration Error",
			RequestBody: in.ActivationUserNexmileRequest{
				UserRegistrationDetailID: 3,
				ParentClientID:           "123",
				FirstName:                "Pertama",
				LastName:                 "Terakhir",
				AndroidID:                "ValidAndroID",
				ClientID:                 "TestClientID",
				UserID:                   "100",
				AuthUserID:               100,
				AliasName:                "ValidAlias",
				Email:                    validEmail,
				Phone:                    validPhone,
			},
			Expected:    errorModel.GenerateUnknownDataError(fileName, "activateUserNexmile", constanta.UserRegistrationDetailName),
		},
		contextModel:     applicationModel.ContextModel{
			AuthAccessTokenModel: model.AuthAccessTokenModel{
				RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
					ResourceUserID: 12,
				},
				ClientID: "123",
				Locale:   constanta.IndonesianLanguage,
			},
		},
	})

	result = append(result, scenarioAktivasiUserNexmile{
		AbstractScenario: test.AbstractScenario{
			Name:        "Success Activate User Nexmile",
			RequestBody: in.ActivationUserNexmileRequest{
				UserRegistrationDetailID: userRegistrationDetailIDToBeActive,
				ParentClientID:           "08181c991e6b409eb016cfaa365b439d",
				FirstName:                "Pertama",
				LastName:                 "Terakhir",
				AndroidID:                "ValidAndroID",
				ClientID:                 "TestClientID",
				UserID:                   "100",
				AuthUserID:               100,
				AliasName:                "ValidAlias",
				Email:                    validEmail,
				Phone:                    validPhone,
			},
			Expected:    errorModel.GenerateNonErrorModel(),
		},
		contextModel:     applicationModel.ContextModel{
			AuthAccessTokenModel: model.AuthAccessTokenModel{
				RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
					ResourceUserID: 12,
				},
				ClientID: "08181c991e6b409eb016cfaa365b439d",
				Locale:   constanta.IndonesianLanguage,
			},
		},
	})

	result = append(result, scenarioAktivasiUserNexmile{
		AbstractScenario: test.AbstractScenario{
			Name:        "Failed Get User License Data",
			RequestBody: in.ActivationUserNexmileRequest{
				UserRegistrationDetailID: userRegistrationDetailIDToBeError,
				ParentClientID:           "08181c991e6b409eb016cfaa365b439d",
				FirstName:                "Pertama1",
				LastName:                 "Terakhir1",
				AndroidID:                "ValidAndroID",
				ClientID:                 "TestClientID1",
				UserID:                   "101",
				AuthUserID:               101,
				AliasName:                "ValidAlias",
				Email:                    validEmail,
				Phone:                    validPhone,
			},
			Expected:    errorModel.GenerateActivationLicenseError(fileName, "getAvailableUserLicense", constanta.ActivationUserNexmileErrorUserLicense, nil),
		},
		contextModel:     applicationModel.ContextModel{
			AuthAccessTokenModel: model.AuthAccessTokenModel{
				RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
					ResourceUserID: 12,
				},
				ClientID: "08181c991e6b409eb016cfaa365b439d",
				Locale:   constanta.IndonesianLanguage,
			},
		},
	})
	
	return 
}

func TestActivateUserNexmile(t *testing.T) {
	usecaseTest := createScenarioAktivasiUserNexmile()
	var errMessage string

	for _, useCase := range usecaseTest {
		t.Run(useCase.Name, func(t *testing.T) {
			request, errs := http.NewRequest(http.MethodPut, "/v1/nextrac/activation-license", nil)
			if errs != nil || request == nil {
				t.Log(errs)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyBtye, _:= json.Marshal(useCase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyBtye)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := ActivationUserNexmileService.ActivationUserNexmileService.ActivateUserNexmile(request, &useCase.contextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, useCase.Expected, err, errMessage)
		})
	}
}
