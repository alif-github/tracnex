package license_config

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
	"nexsoft.co.id/nextrac2/service/LicenseConfigService"
	"nexsoft.co.id/nextrac2/test"
	"strconv"
	"testing"
)

type scenarioTableUpdateLicenseConfiguration struct {
	test.AbstractScenario
	contextModel applicationModel.ContextModel
}

func createScenarioTableUpdateLicenseConfig() []scenarioTableUpdateLicenseConfiguration {
	fileName := "UpdateLicenseConfigService.go"
	funcName := "doUpdateLicenseConfig"
	return []scenarioTableUpdateLicenseConfiguration{
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Failed Update License Configuration",
				RequestBody: in.LicenseConfigRequest{
					ID:           2,
					UpdatedAtStr: InitiateTestVar.timeNow.Format(constanta.DefaultTimeFormat),
				},
				Expected: errorModel.GenerateDataLockedError(fileName, funcName, constanta.LicenseConfig),
			},
			contextModel: *InitiateTestVar.ContextModel,
		},
	}
}

func TestUpdateLicenseConfiguration(t *testing.T) {
	useCaseTest := createScenarioTableUpdateLicenseConfig()
	var errMessage string

	for _, useCase := range useCaseTest {
		t.Run(useCase.Name, func(t *testing.T) {
			request, errS := http.NewRequest(http.MethodPut, "/v1/nextrac/licenseconfig/"+strconv.Itoa(int(useCase.RequestBody.(in.LicenseConfigRequest).ID)), nil)

			if errS != nil || request == nil {
				t.Log(errS)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyByte, _ := json.Marshal(useCase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyByte)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := LicenseConfigService.LicenseConfigService.UpdateLicenseConfig(request, &useCase.contextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, useCase.Expected, err, errMessage)
		})
	}
}
