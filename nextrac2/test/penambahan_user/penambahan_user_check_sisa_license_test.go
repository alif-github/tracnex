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
	"nexsoft.co.id/nextrac2/service/UserRegistrationService"
	"nexsoft.co.id/nextrac2/test"
	"testing"
)

func createScenarionCekSisaLicense() (result []test.AbstractScenarioWithContextModel) {
	result = append(result, test.AbstractScenarioWithContextModel{
		Name:         "Success Check Sisa License",
		RequestBody:  in.CheckLicenseNamedUserRequest{
			ClientId:     "c833660f2d254027b527da2320d39f14",
			ClientTypeID: constanta.ResourceNexmileID,
			UniqueId1:    "NS6084010002596",
			UniqueId2:    "1596128276342",
		},
		Expected:     errorModel.GenerateNonErrorModel(),
		ContextModel: applicationModel.ContextModel{
			AuthAccessTokenModel: model.AuthAccessTokenModel{
				RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
					ResourceUserID: 12,
				},
				ClientID: "08181c991e6b409eb016cfaa365b439d",
				Locale:   constanta.IndonesianLanguage,
			},
		},
	})

	result = append(result, test.AbstractScenarioWithContextModel{
		Name:         "Fail Check Sisa License",
		RequestBody:  in.CheckLicenseNamedUserRequest{
			ClientId:     "1743b12b50074b0cb993d2a43badf36a",
			ClientTypeID: constanta.ResourceND6ID,
			UniqueId1:    "NS6024050001135",
			UniqueId2:    "1468381448675",
		},
		Expected:     errorModel.GenerateUserLicenseNotFound("UserRegistrationService.go", "CheckLicenseNamedUser"),
		ContextModel: applicationModel.ContextModel{
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

func TestCheckRemainingLicense(t *testing.T) {
	useCaseTest := createScenarionCekSisaLicense()
	var errMessage string

	for _, useCase := range useCaseTest {
		t.Run(useCase.Name, func(t *testing.T) {
			request, errs := http.NewRequest(http.MethodPost, "/v1/nextrac/user-registration/check", nil)
			if errs != nil || request == nil {
				t.Log(errs)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyBtye, _:= json.Marshal(useCase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyBtye)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := UserRegistrationService.UserRegistrationService.CheckLicenseNamedUser(request, &useCase.ContextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, useCase.Expected, err, errMessage)
		})
	}
}