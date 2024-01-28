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
	util2 "nexsoft.co.id/nextrac2/util"
	"testing"
)

type scenarioTableInsertLicenseConfiguration struct {
	test.AbstractScenario
	contextModel applicationModel.ContextModel
}

func createScenarioTableInsertLicenseConfig() []scenarioTableInsertLicenseConfiguration {
	fileName := "InsertLicenseConfigService.go"
	funcName := "validateCustomerInstallationAndClientMapping"
	fileNameDTO := "LicenseConfigDTO.go"
	funcNameDTO := "ValidateInsertLicenseConfig"
	customerStr := util2.GenerateConstantaI18n(constanta.Customer, InitiateTestVar.ContextModel.AuthAccessTokenModel.Locale, nil)
	installationIDStr := util2.GenerateConstantaI18n(constanta.InstallationID, InitiateTestVar.ContextModel.AuthAccessTokenModel.Locale, nil)
	return []scenarioTableInsertLicenseConfiguration{
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Negative Case Installation ID Not Found",
				RequestBody: in.LicenseConfigRequest{
					InstallationID:      99,
					NoOfUser:            1,
					ProductValidFromStr: "2022-01-20",
					ProductValidThruStr: "2023-01-20",
				},
				Expected: errorModel.GenerateClientIDNotFound(fileName, funcName, constanta.Installation, customerStr, installationIDStr),
			},
			contextModel: *InitiateTestVar.ContextModel,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Negative Case Installation ID 0",
				RequestBody: in.LicenseConfigRequest{
					InstallationID:      0,
					NoOfUser:            1,
					ProductValidFromStr: "2022-01-20",
					ProductValidThruStr: "2023-01-20",
				},
				Expected: errorModel.GenerateEmptyFieldOrZeroValueError(fileNameDTO, funcNameDTO, constanta.InstallationID),
			},
			contextModel: *InitiateTestVar.ContextModel,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Negative Case Product Valid Thru Before Valid From",
				RequestBody: in.LicenseConfigRequest{
					InstallationID:      InitiateTestVar.installationID,
					NoOfUser:            1,
					ProductValidFromStr: "2023-01-20",
					ProductValidThruStr: "2022-01-20",
				},
				Expected: errorModel.GenerateDateValidateFromThru(fileNameDTO, funcNameDTO, "E-6-TRAC-SRV-016"),
			},
			contextModel: *InitiateTestVar.ContextModel,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Negative Case No Of User Less Than 1",
				RequestBody: in.LicenseConfigRequest{
					InstallationID:      InitiateTestVar.installationID,
					NoOfUser:            0,
					ProductValidFromStr: "2022-01-20",
					ProductValidThruStr: "2023-01-20",
				},
				Expected: errorModel.GenerateEmptyFieldOrZeroValueError(fileNameDTO, funcNameDTO, constanta.NumberOfUser),
			},
			contextModel: *InitiateTestVar.ContextModel,
		},
		{
			AbstractScenario: test.AbstractScenario{
				Name: "Success Insert License Configuration",
				RequestBody: in.LicenseConfigRequest{
					InstallationID:      InitiateTestVar.installationID,
					NoOfUser:            1,
					ProductValidFromStr: "2022-01-20",
					ProductValidThruStr: "2023-01-20",
				},
				Expected: errorModel.GenerateNonErrorModel(),
			},
			contextModel: *InitiateTestVar.ContextModel,
		},
	}
}

func TestAddNewLicenseConfiguration(t *testing.T) {
	useCaseTest := createScenarioTableInsertLicenseConfig()
	var errMessage string

	for _, useCase := range useCaseTest {
		t.Run(useCase.Name, func(t *testing.T) {
			request, errS := http.NewRequest(http.MethodPost, "/v1/nextrac/licenseconfig", nil)

			if errS != nil || request == nil {
				t.Log(errS)
				assert.FailNow(t, "Error build a request")
			}

			reqBodyByte, _ := json.Marshal(useCase.RequestBody)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(string(reqBodyByte)))
			request.Header.Set("Content-Type", "application/json")

			output, _, err := LicenseConfigService.LicenseConfigService.InsertLicenseConfig(request, &useCase.contextModel)
			errMessage = util.StructToJSON(output)

			assert.Equal(t, useCase.Expected, err, errMessage)
		})
	}
}
