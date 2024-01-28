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
	"nexsoft.co.id/nextrac2/service/RegistrationNamedUserService"
	"nexsoft.co.id/nextrac2/test"
	"testing"
)

type scenarioRegistrasiUserNexmile struct {
	test.AbstractScenario
	contextModel applicationModel.ContextModel
}

func createScenarioRegistrasiUserNexmile() (result []scenarioRegistrasiUserNexmile) {
	fileName := "RegistrationNamedUserService.go"
	validEmail := "ValidEmail@email.com"
	validPhone := "+62-897789666889"

	result = append(result, scenarioRegistrasiUserNexmile{
		AbstractScenario: test.AbstractScenario{
			Name:        "Forbidden Access Client ID ",
			RequestBody: in.RegisterNamedUserRequest{
				ParentClientID:   "123",
				ClientID:         "08181c920e6b409eb016cfaa365b439d",
				ClientTypeID:     2,
				AuthUserID:       102,
				Firstname:        "Pertama2",
				Lastname:         "Terakhir2",
				Username:         "UsernameValid2",
				UserID:           "102",
				Password:         "PassValid2",
				ClientAliases:    "ValidAlias",
				SalesmanID:       "ValidSalesID",
				AndroidID:        "ValidAndroID",
				Email:            validEmail,
				UniqueID1:        "NS6024050001031",
				UniqueID2:        "1468381449586",
				NoTelp:           validPhone,
				SalesmanCategory: "ValidCtg",
			},
			Expected:    errorModel.GenerateForbiddenClientCredentialAccess(fileName, "readBodyAndValidate"),
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

	result = append(result, scenarioRegistrasiUserNexmile{
		AbstractScenario: test.AbstractScenario{
			Name:        "Client Not Valid on Auth",
			RequestBody: in.RegisterNamedUserRequest{
				ParentClientID:   "08181c991e6b409eb016cfaa365b438d",
				ClientID:         "12381c991e6b409eb016cfaa365b438d",
				ClientTypeID:     2,
				AuthUserID:       102,
				Firstname:        "Pertama2",
				Lastname:         "Terakhir2",
				Username:         "UsernameValid2",
				UserID:           "102",
				Password:         "PassValid2",
				ClientAliases:    "ValidAlias",
				SalesmanID:       "ValidSalesID",
				AndroidID:        "ValidAndroID",
				Email:            validEmail,
				UniqueID1:        "NS6024050001031",
				UniqueID2:        "1468381449586",
				NoTelp:           validPhone,
				SalesmanCategory: "ValidCtg",
			},
			Expected:    errorModel.GenerateUnknownDataError("InsertNamedUserService.go", "checkDataToAuthServer", constanta.ClientMappingClientID),
		},
		contextModel:     applicationModel.ContextModel{
			AuthAccessTokenModel: model.AuthAccessTokenModel{
				RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
					ResourceUserID: 12,
				},
				ClientID: "08181c991e6b409eb016cfaa365b438d",
				Locale:   constanta.IndonesianLanguage,
			},
		},
	})

	result = append(result, scenarioRegistrasiUserNexmile{
		AbstractScenario: test.AbstractScenario{
			Name:        "Cannot find user license",
			RequestBody: in.RegisterNamedUserRequest{
				ParentClientID:   "08181c991e6b409eb016cfaa365b439d",
				ClientID:         "08181c920e6b409eb016cfaa365b439d",
				ClientTypeID:     2,
				AuthUserID:       102,
				Firstname:        "Pertama2",
				Lastname:         "Terakhir2",
				Username:         "UsernameValid2",
				UserID:           "102",
				Password:         "PassValid2",
				ClientAliases:    "ValidAlias",
				SalesmanID:       "ValidSalesID",
				AndroidID:        "ValidAndroID",
				Email:            validEmail,
				UniqueID1:        "NS6024050001032",
				UniqueID2:        "1468381449582",
				NoTelp:           validPhone,
				SalesmanCategory: "ValidCtg",
			},
			Expected:    errorModel.GenerateNotFoundActiveLicense("InsertNamedUserClientMappingService.go", "doInsertNamedUserClientMapping"),
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

	result = append(result, scenarioRegistrasiUserNexmile{
		AbstractScenario: test.AbstractScenario{
			Name:        "Success Regist User Nexmile",
			RequestBody: in.RegisterNamedUserRequest{
				ParentClientID:   "08181c991e6b409eb016cfaa365b439d",
				ClientID:         "08181c920e6b409eb016cfaa365b439d",
				ClientTypeID:     2,
				AuthUserID:       102,
				Firstname:        "Pertama2",
				Lastname:         "Terakhir2",
				Username:         "UsernameValid2",
				UserID:           "102",
				Password:         "PassValid2",
				ClientAliases:    "ValidAlias",
				SalesmanID:       "ValidSalesID",
				AndroidID:        "ValidAndroID",
				Email:            validEmail,
				UniqueID1:        "NS6024050001031",
				UniqueID2:        "1468381449586",
				NoTelp:           validPhone,
				SalesmanCategory: "ValidCtg",
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

	return
}

func TestRegisterUserNexmile(t *testing.T) {
	useCaseTest := createScenarioRegistrasiUserNexmile()
	var errMessage string

	for _, useCase := range useCaseTest {
		t.Run(useCase.Name, func(t *testing.T) {
			request, errs := http.NewRequest(http.MethodPut, "/v1/nextrac/user-registration/add-nexmile", nil)
			if errs != nil || request == nil {
				t.Log(errs)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyBtye, _:= json.Marshal(useCase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyBtye)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := RegistrationNamedUserService.RegistrationNamedUserService.InsertNamedUserClientMapping(request, &useCase.contextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, useCase.Expected, err, errMessage)
		})
	}
}