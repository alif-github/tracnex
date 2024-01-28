package penambahan_user_admin

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
	"nexsoft.co.id/nextrac2/service/UserRegistrationAdminService"
	"nexsoft.co.id/nextrac2/test"
	"testing"
)

type scenarioTablePenambahanUserAdmin struct {
	test.AbstractScenario
	contextModel applicationModel.ContextModel
}

func createScenarioPenambahanUserAdmin() []scenarioTablePenambahanUserAdmin {
	return []scenarioTablePenambahanUserAdmin{
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Skenario 1 Negative Case - PT. Makmur Sejahtera (Nexmile) melakukan penambahan user admin dengan client id yang dikirimkan berbeda dengan client id yang login di resource",
				RequestBody: in.UserRegistrationAdminRequest{
					UniqueID1:     "123",
					UniqueID2:     "123",
					UserAdmin:     "user01",
					PasswordAdmin: "abc123",
					CompanyName:   "PT. Makmur Sejahtera",
					BranchName:    "ABC",
					ClientID:      "891",
					ClientTypeID:  5,
				},
				Expected: errorModel.GenerateForbiddenAccessClientError("UserRegistrationAdminService.go", "readBodyAndValidate"),
			},
			contextModel: contextModelNexMile,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Skenario 1 Negative Case - PT. Makmur Sejahtera (Nexmile) gagal melakukan penambahan user admin karena client id yang dicari tidak ditemukan pada tabel client mapping",
				RequestBody: in.UserRegistrationAdminRequest{
					UniqueID1:     "123",
					UniqueID2:     "123",
					UserAdmin:     "user01",
					PasswordAdmin: "abc123",
					CompanyName:   "PT. Makmur Sejahtera",
					BranchName:    "ABC",
					ClientID:      "891",
					ClientTypeID:  5,
				},
				Expected:  errorModel.GenerateUnknownDataError("InsertUserRegistrationAdminService.go", "doInsertUserRegistrationAdmin", constanta.ClientMappingClientID),
			},
			contextModel: applicationModel.ContextModel{
				AuthAccessTokenModel: model.AuthAccessTokenModel{
					RedisAuthAccessTokenModel: model.RedisAuthAccessTokenModel{
						ResourceUserID: 12,
					},
					ClientID: "891",
					Locale:   constanta.IndonesianLanguage,
				},
			},
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Skenario 2 Positive Case - PT. Makmur Sejahtera (Nexmile) melakukan penambahan user admin",
				RequestBody: in.UserRegistrationAdminRequest{
					UniqueID1:     "123",
					UniqueID2:     "123",
					UserAdmin:     "user01",
					PasswordAdmin: "abc123",
					CompanyName:   "PT. Makmur Sejahtera",
					BranchName:    "ABC",
					ClientID:      "98381c991e6b409eb016cfaa365k4cad",
					ClientTypeID:  5,
				},
				Expected: errorModel.GenerateNonErrorModel(),
			},
			contextModel: contextModelNexMile,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Skenario 3 Positive Case - PT. Maju Jaya (NexChief Mobile) melakukan penambahan user admin",
				RequestBody: in.UserRegistrationAdminRequest{
					UniqueID1:     "456",
					UniqueID2:     "456",
					UserAdmin:     "user02",
					PasswordAdmin: "abc456",
					CompanyName:   "PT. Maju Jaya",
					BranchName:    "DEF",
					ClientID:      "r3fb12faf6a348759ccffc500d609f31",
					ClientTypeID:  6,
				},
				Expected: errorModel.GenerateNonErrorModel(),
			},
			contextModel: contextModelNexChiefMobile,
		},
	}
}

func TestPenambahanUserAdmin(t *testing.T) {
	usecases := createScenarioPenambahanUserAdmin()
	var errMessage string

	for _, usecase := range usecases {
		t.Run(usecase.Name, func(t *testing.T) {
			request, errs := http.NewRequest(http.MethodPost, "v1/nextrac/user-registration-admin", nil)
			if errs != nil || request == nil {
				t.Log(errs)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyBtye, _ := json.Marshal(usecase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyBtye)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := UserRegistrationAdminService.UserRegistrationAdminService.InsertUserRegistrationAdmin(request, &usecase.contextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, usecase.Expected, err, errMessage)
		})
	}
}
