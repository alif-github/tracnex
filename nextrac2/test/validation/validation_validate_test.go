package validation

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
	"nexsoft.co.id/nextrac2/service/ValidationLicenseService"
	"nexsoft.co.id/nextrac2/test"
	"testing"
)

type scenarioTableValidation struct {
	test.AbstractScenario
	contextModel applicationModel.ContextModel
}

func createScenarioValidation() []scenarioTableValidation {
	return []scenarioTableValidation{
		// Skenario 1
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Negative Case - Tidak menemukan data Client Validation",
				RequestBody: in.ValidationLicenseRequest{
					AbstractDTO:  in.AbstractDTO{},
					ClientID:     "123",
					ClientTypeID: 1,
					SignatureKey: "280d9968c4154d698362087a91a80e1a",
					HwID:         "1",
					ProductDetail: []in.ProductEncryptDetail{
						{
							ProductKey:     "1",
							ProductEncrypt: "1",
						},
					},
				},
				Expected: errorModel.GenerateForbiddenAccessClientError("ValidationLicenseService.go", "doValidateLicense"),
			},
			contextModel: contextModelND,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Negative Case - Data Product License tidak ditemukan",
				RequestBody: in.ValidationLicenseRequest{
					AbstractDTO:  in.AbstractDTO{},
					ClientID:     "08181c991e6b409eb016cfaa365b439d",
					ClientTypeID: 1,
					SignatureKey: "280d9968c4154d698362087a91a80e1a",
					HwID:         "1",
					ProductDetail: []in.ProductEncryptDetail{
						{
							ProductKey:     "6",
							ProductEncrypt: "1",
						},
					},
				},
				Expected: errorModel.GenerateActivationLicenseError("ValidationLicenseService.go", "doValidateLicense", constanta.ActivationLicenseGetProductLicenseError, nil),
			},
			contextModel: contextModelND,
		},
		// Skenario 2
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Positive Case - Validasi License ND6 yang telah dibeli",
				RequestBody: in.ValidationLicenseRequest{
					AbstractDTO:  in.AbstractDTO{},
					ClientID:     "08181c991e6b409eb016cfaa365b439d",
					ClientTypeID: 1,
					SignatureKey: "280d9968c4154d698362087a91a80e1a",
					HwID:         "1",
					ProductDetail: []in.ProductEncryptDetail{
						{
							ProductKey:     "1",
							ProductEncrypt: "1",
						},
					},
				},
				Expected: errorModel.GenerateNonErrorModel(),
			},
			contextModel: contextModelND,
		},
		// Skenario 3
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Positive Case - Validasi License NexChief yang telah dibeli",
				RequestBody: in.ValidationLicenseRequest{
					AbstractDTO:  in.AbstractDTO{},
					ClientID:     "1a2b12faf6a345759ccffc500d609b52",
					ClientTypeID: 4,
					SignatureKey: "bb0734e85ba44b529611fd22668b6bad",
					HwID:         "4",
					ProductDetail: []in.ProductEncryptDetail{
						{
							ProductKey:     "2",
							ProductEncrypt: "2",
						},
					},
				},
				Expected: errorModel.GenerateNonErrorModel(),
			},
			contextModel: contextModelNexchief,
		},
	}
}

func TestValidation(t *testing.T) {
	usecases := createScenarioValidation()
	var errMessage string

	for _, usecase := range usecases {
		t.Run(usecase.Name, func(t *testing.T) {
			request, errs := http.NewRequest(http.MethodPut, "/v1/nextrac/validation-license", nil)
			if errs != nil || request == nil {
				t.Log(errs)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyBtye, _ := json.Marshal(usecase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyBtye)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := ValidationLicenseService.ValidationLicenseService.ValidateLicense(request, &usecase.contextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, usecase.Expected, err, errMessage)
		})
	}

}
