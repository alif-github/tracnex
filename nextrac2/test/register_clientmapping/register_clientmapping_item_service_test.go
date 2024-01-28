package register_clientmapping

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service/ClientRegistrationNonOnPremiseService"
	"nexsoft.co.id/nextrac2/test"
	util2 "nexsoft.co.id/nextrac2/util"
	"testing"
)

func TestRegisterClientMappingClientCredentialValidation(t *testing.T) {
	scenario := []struct {
		name     string
		body     interface{}
		expected interface{}
	}{
		{
			name: "Test Negative Client Credential Failed ND6",
			body: in.ClientRegistrationNonOnPremiseRequest{
				ClientID:     "123",
				ClientTypeID: 1,
				DetailClient: []in.UniqueIDClient{
					{
						UniqueID1: "NS6024050001031",
						UniqueID2: "1468381449586",
					},
				},
			},
			expected: "E-4-TRAC-SRV-004",
		},
		{
			name: "Test Negative Customer Installation Not Found ND6",
			body: in.ClientRegistrationNonOnPremiseRequest{
				ClientID:     "08181c991e6b409eb016cfaa365b439d",
				ClientTypeID: 1,
				DetailClient: []in.UniqueIDClient{
					{
						UniqueID1: "NS30524050001031",
						UniqueID2: "1528381449586",
					},
				},
			},
			expected: "E-6-TRAC-SRV-011",
		},
		{
			name: "Test Positive ND6",
			body: in.ClientRegistrationNonOnPremiseRequest{
				ClientID:     "08181c991e6b409eb016cfaa365b439d",
				ClientTypeID: 1,
				DetailClient: []in.UniqueIDClient{
					{
						UniqueID1: "NS6024050001031",
						UniqueID2: "1468381449586",
					},
				},
			},
			expected: "",
		},
		{
			name: "Test Negative Client Credential Failed Nexchief",
			body: in.ClientRegistrationNonOnPremiseRequest{
				ClientID:     "123",
				ClientTypeID: 3,
				DetailClient: []in.UniqueIDClient{
					{
						UniqueID1: "NDI",
					},
				},
			},
			expected: "E-4-TRAC-SRV-004",
		},
		{
			name: "Test Negative Customer Installation Not Found Nexchief",
			body: in.ClientRegistrationNonOnPremiseRequest{
				ClientID:     "1a2b12faf6a345759ccffc500d609b52",
				ClientTypeID: 3,
				DetailClient: []in.UniqueIDClient{
					{
						UniqueID1: "NDI",
						UniqueID2: "ABC",
					},
				},
			},
			expected: "E-6-TRAC-SRV-011",
		},
		{
			name: "Test Positive Nexchief",
			body: in.ClientRegistrationNonOnPremiseRequest{
				ClientID:     "1a2b12faf6a345759ccffc500d609b52",
				ClientTypeID: 3,
				DetailClient: []in.UniqueIDClient{
					{
						UniqueID1: "NDI",
					},
				},
			},
			expected: "",
		},
	}

	for _, item := range scenario {
		t.Run(item.name, func(t *testing.T) {
			var errMessage string
			request := test.SetRequest(t, serverconfig.ServerAttribute.DBConnection, "POST", "/v1/nextrac/client-mapping")

			bodyByte, _ := json.Marshal(item.body.(in.ClientRegistrationNonOnPremiseRequest))
			bodyStr := string(bodyByte)

			request.Body = ioutil.NopCloser(bytes.NewBufferString(bodyStr))
			request.Header.Set("Content-Type", "application/json")

			_, _, errs := ClientRegistrationNonOnPremiseService.ClientRegistrationNonOnPremiseService.InsertClientRegistNonOnPremise(request, InitiateTestVar.ContextModel)
			if errs.Error != nil {
				errMessage = util2.GenerateI18NErrorMessage(errs, InitiateTestVar.ContextModel.AuthAccessTokenModel.Locale)
			}

			assert.Equal(t, item.expected.(string), errMessage)
		})
	}
}
