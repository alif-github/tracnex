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

func createScenarioRegistOtherNamedUser() (result []test.AbstractScenarioWithContextModel) {
	email := "email5@email.com"
	phoneData := "+62-87998767787"
	result = append(result, test.AbstractScenarioWithContextModel{
		Name:         "Failed Check Auth Server",
		RequestBody:  in.RegisterNamedUserRequest{
			ParentClientID:   "c833660f2d254027b527da2320d39f14",
			ClientID:         "59245771231e4caaa87c1636e3de186a",
			ClientTypeID:     constanta.ResourceNexmileID,
			AuthUserID:       103,
			Firstname:        "Pertama5",
			Lastname:         "Akhir5",
			Username:         "UserName5",
			UserID:           "103",
			Password:         "Password5",
			ClientAliases:    "ClientAlias5",
			SalesmanID:       "SalesID5",
			AndroidID:        "AndroidID",
			Email:            email,
			UniqueID1:        "NS6084010002596",
			UniqueID2:        "1596128276342",
			NoTelp:           phoneData,
			SalesmanCategory: "SalesCTG",
		},
		Expected:     errorModel.GenerateUnknownDataError("InsertNamedUserService.go", "checkDataToAuthServer", constanta.ClientMappingClientID),
		ContextModel: applicationModel.ContextModel{
			AuthAccessTokenModel: model.AuthAccessTokenModel{
				RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
					ResourceUserID: 12,
				},
				ClientID: "c833660f2d254027b527da2320d39f14",
				Locale:   constanta.IndonesianLanguage,
			},
		},
	})

	result = append(result, test.AbstractScenarioWithContextModel{
		Name:         "Success Regist Other Named User",
		RequestBody:  in.RegisterNamedUserRequest{
			ParentClientID:   "c833660f2d254027b527da2320d39f14",
			ClientID:         "59005771231e4caaa87c1636e3de186a",
			ClientTypeID:     constanta.ResourceNexmileID,
			AuthUserID:       103,
			Firstname:        "Pertama5",
			Lastname:         "Akhir5",
			Username:         "UserName5",
			UserID:           "103",
			Password:         "Password5",
			ClientAliases:    "ClientAlias5",
			SalesmanID:       "SalesID5",
			AndroidID:        "AndroidID",
			Email:            email,
			UniqueID1:        "NS6084010002596",
			UniqueID2:        "1596128276342",
			NoTelp:           phoneData,
			SalesmanCategory: "SalesCTG",
		},
		Expected:     errorModel.GenerateNonErrorModel(),
		ContextModel: applicationModel.ContextModel{
			AuthAccessTokenModel: model.AuthAccessTokenModel{
				RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
					ResourceUserID: 12,
				},
				ClientID: "c833660f2d254027b527da2320d39f14",
				Locale:   constanta.IndonesianLanguage,
			},
		},
	})

	result = append(result, test.AbstractScenarioWithContextModel{
		Name:         "Failed check license quota",
		RequestBody:  in.RegisterNamedUserRequest{
			ParentClientID:   "c833660f2d254027b527da2320d39f14",
			ClientID:         "62cda0a7242c497bbc502e4b33e87abc",
			ClientTypeID:     constanta.ResourceNexmileID,
			AuthUserID:       104,
			Firstname:        "Pertama5",
			Lastname:         "Akhir5",
			Username:         "UserName5",
			UserID:           "104",
			Password:         "Password5",
			ClientAliases:    "ClientAlias5",
			SalesmanID:       "SalesID5",
			AndroidID:        "AndroidID",
			Email:            email,
			UniqueID1:        "NS6084010002596",
			UniqueID2:        "1596128276342",
			NoTelp:           phoneData,
			SalesmanCategory: "SalesCTG",
		},
		Expected:     errorModel.GenerateUserLicenseNotFound("InsertNamedUserService.go", "checkLicenseQuota"),
		ContextModel: applicationModel.ContextModel{
			AuthAccessTokenModel: model.AuthAccessTokenModel{
				RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
					ResourceUserID: 12,
				},
				ClientID: "c833660f2d254027b527da2320d39f14",
				Locale:   constanta.IndonesianLanguage,
			},
		},
	})
	return
}

func TestRegistOtherNamedUser(t *testing.T) {
	useCaseTest := createScenarioRegistOtherNamedUser()
	var errMessage string

	for _, useCase := range useCaseTest {
		t.Run(useCase.Name, func(t *testing.T) {
			request, errs := http.NewRequest(http.MethodPost, "/v1/nextrac/user-registration/add", nil)
			if errs != nil || request == nil {
				t.Log(errs)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyBtye, _:= json.Marshal(useCase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyBtye)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := RegistrationNamedUserService.RegistrationNamedUserService.InsertNamedUser(request, &useCase.ContextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, useCase.Expected, err, errMessage)
		})
	}
}
