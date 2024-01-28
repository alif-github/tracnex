package license_variant

import (
	"bytes"
	"encoding/json"
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

func createRequestBodyForInsert() []scenarioTable {
	fileName1 := "AbstractDTO.go"
	fileName2 := "LicenseVariantDTO.go"
	funcName1 := "ValidateMinMaxString"
	funcName2 := "mandatoryValidation"

	return []scenarioTable{
		{
			name: "SuccessInsertDataLicenseVariant",
			requestBody: in.LicenseVariantRequest{
				LicenseVariantName: "Valid",
			},
			expected: errorModel.GenerateNonErrorModel(),
		},
		{
			name: "FailInsertDataLicenseVariantNameToLong",
			requestBody: in.LicenseVariantRequest{
				LicenseVariantName: "Failed Because Is Length Of Name License Variant To Long",
			},
			expected: errorModel.GenerateFieldFormatWithRuleError(fileName1, funcName1, "NEED_LESS_THAN", constanta.LicenseVariantName, strconv.Itoa(20)),
		},
		{
			name: "FailInsertDataLicenseVariantNameEmpty",
			expected: errorModel.GenerateEmptyFieldError(fileName1, funcName1, constanta.LicenseVariantName),
		},
		{
			name: "FailInsertDataLicenseVariantNameSpcCharacterExist",
			requestBody: in.LicenseVariantRequest{
				LicenseVariantName: "Valid%",
			},
			expected: errorModel.GenerateFieldFormatWithRuleError(fileName2, funcName2, constanta.ErrorSpecialCharacter, constanta.LicenseVariantName, ""),
		},
	}
}

func TestInsertLicenseVariant(t *testing.T) {
	useCaseTest := createRequestBodyForInsert()
	var errMessage string
	for _, useCase := range useCaseTest {
		t.Run(useCase.name, func(t *testing.T) {
			request, errs := http.NewRequest(http.MethodPost, "/v1/nextrac/license_variant", nil)
			if errs != nil || request == nil {
				t.Log(errs)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyByte, _ := json.Marshal(useCase.requestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyByte)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := LicenseVariantService.LicenseVariantService.InsertLicenseVariant(request, &contextModel)
			if err.Error != nil {
				errMessage = util2.GenerateI18NErrorMessage(err, contextModel.AuthAccessTokenModel.Locale)
			} else {
				errMessage = output.Status.Message
			}

			assert.Equal(t, useCase.expected, err, errMessage)
		})
	}
}
