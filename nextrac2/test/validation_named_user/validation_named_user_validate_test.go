package validation_named_user

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
	"nexsoft.co.id/nextrac2/service/ValidationNamedUserService"
	"nexsoft.co.id/nextrac2/test"
	"testing"
)

func createScenarioValidateUser() (result []test.AbstractScenarioWithContextModel) {
	result = append(result, test.AbstractScenarioWithContextModel{
		Name:         "Success Validate User",
		RequestBody:  in.ValidationNamedUserRequest{
			ClientId:    "62cda0a7242c497bbc502e4b33e87abc",
			AuthUserId:  100,
			UserId:      "100",
			UniqueId1:   "NS6024050001031",
			UniqueId2:   "1468381449586",
		},
		Expected:     errorModel.GenerateNonErrorModel(),
		ContextModel: applicationModel.ContextModel{
			AuthAccessTokenModel: model.AuthAccessTokenModel{
				RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
					ResourceUserID: 12,
				},
				ClientID: "62cda0a7242c497bbc502e4b33e87abc",
				Locale:   constanta.IndonesianLanguage,
			},
		},
	})
	
	return 
}

func TestValidateNamedUser(t *testing.T) {
	useCaseTest := createScenarioValidateUser()
	var errMessage string

	for _, useCase := range useCaseTest {
		t.Run(useCase.Name, func(t *testing.T) {
			request, errs := http.NewRequest(http.MethodPost, "/v1/nextrac/user-registration/validate", nil)
			if errs != nil || request == nil {
				t.Log(errs)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyBtye, _:= json.Marshal(useCase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyBtye)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := ValidationNamedUserService.ValidationNamedUserService.ValidateNamedUser(request, &useCase.ContextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, useCase.Expected, err, errMessage)
		})
	}
}
