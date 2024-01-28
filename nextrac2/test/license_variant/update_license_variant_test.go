package license_variant

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service/LicenseVariantService"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"testing"
)

type scenarioTable struct {
	name        string
	requestBody interface{}
	expected    errorModel.ErrorModel
}

func createRequestBodyForUpdate() []scenarioTable {
	fileName := "UpdateLicenseVariantService.go"
	funcName := "doUpdateLicenseVariant"
	return []scenarioTable{
		{
			name:        "SuccessUpdateDataLicenseVariant",
			requestBody: in.LicenseVariantRequest{
				ID: 90,
				LicenseVariantName: "Updated License",
				UpdatedAtStr: "2021-12-14T10:16:39.631007Z",
			},
			expected:    errorModel.GenerateNonErrorModel(),
		},
		{
			name:        "FailUpdateDataLicenseVariantLockData",
			requestBody: in.LicenseVariantRequest{
				ID: 92,
				LicenseVariantName: "Updated License",
				UpdatedAtStr: "2021-12-14T10:16:39.631007Z",
			},
			expected:    errorModel.GenerateDataLockedError(fileName, funcName, constanta.UpdatedAt),
		},
		{
			name:        "FailUpdateDataLicenseVariantDataNotFound",
			requestBody: in.LicenseVariantRequest{
				ID: 99,
				LicenseVariantName: "Updated License",
				UpdatedAtStr: "2021-12-14T10:16:39.631007Z",
			},
			expected:    errorModel.GenerateUnknownDataError(fileName, funcName, constanta.LicenseVariant),
		},
	}
}

func TestUpdateLicenseVariant(t *testing.T) {
	useCaseTest := createRequestBodyForUpdate()
	var errMessage string
	for _, useCase := range useCaseTest {
		t.Run(useCase.name, func(t *testing.T) {
			request, errs := http.NewRequest(http.MethodPut, "/v1/nextrac/license_variant", nil)
			if errs != nil || request == nil {
				t.Log(errs)
				assert.FailNow(t, "Error build a request")
			}
			request = mux.SetURLVars(request, map[string]string{
				"ID" : strconv.Itoa(int(useCase.requestBody.(in.LicenseVariantRequest).ID)),
			})
			reqBodyByte, _ := json.Marshal(useCase.requestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyByte)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := LicenseVariantService.LicenseVariantService.UpdateLicenseVariant(request, &contextModel)
			if err.Error != nil {
				errMessage = util2.GenerateI18NErrorMessage(err, contextModel.AuthAccessTokenModel.Locale)
			}else {
				errMessage = output.Status.Message
			}

			assert.Equal(t, useCase.expected, err, errMessage)
		})
	}
}

