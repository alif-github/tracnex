package client_service

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service/ClientService"
	"testing"
)

type testCase struct {
	name     string
	request  request
	expected in.ClientRequest
}

type request struct {
	InputStruct in.ClientRequest
	ModelStruct []repository.ClientMappingModel
}

func setRequest() (requestSet request) {
	requestSet.InputStruct = setInputScenario()
	requestSet.ModelStruct = setComparationScenario()
	return
}

func setOutputExpected() in.ClientRequest {
	var userBranchDataRequest []in.BranchData
	var userCompanyDataRequst []in.CompanyData

	userBranchDataRequest = append(userBranchDataRequest,
		in.BranchData{
			BranchID: "VALID02",
		})

	userCompanyDataRequst = append(userCompanyDataRequst,
		in.CompanyData{
			CompanyID:  "VALID01",
			BranchData: userBranchDataRequest,
		})

	userClientRequest := in.ClientRequest{
		ClientTypeID: 1,
		ClientName:   "VALID04",
		SocketID:     "VALID05",
		CompanyData:  userCompanyDataRequst,
	}

	return userClientRequest
}

func setComparationScenario() (comparationModel []repository.ClientMappingModel) {
	comparationModel = append(comparationModel, repository.ClientMappingModel{
		CompanyID: sql.NullString{String: "VALID01"},
		BranchID:  sql.NullString{String: "VALID03"},
	})

	return
}

func setInputScenario() in.ClientRequest {
	var userBranchDataRequest []in.BranchData
	var userCompanyDataRequst []in.CompanyData

	userBranchDataRequest = append(userBranchDataRequest,
		in.BranchData{
			BranchID: "VALID03",
		}, in.BranchData{
			BranchID: "VALID02",
		})

	userCompanyDataRequst = append(userCompanyDataRequst,
		in.CompanyData{
			CompanyID:  "VALID01",
			BranchData: userBranchDataRequest,
		})

	userClientRequest := in.ClientRequest{
		ClientTypeID: 1,
		ClientName:   "VALID04",
		SocketID:     "VALID05",
		CompanyData:  userCompanyDataRequst,
	}

	return userClientRequest
}

func TestClientService_removeDataSuccess(t *testing.T) {
	var dataInputSubTask []testCase

	dataInputSubTask = append(dataInputSubTask,
		testCase{
			name:     "Check nullable result",
			request:  setRequest(),
			expected: setOutputExpected(),
		}, testCase{
			name:     "Comparation actual dan expected",
			request:  setRequest(),
			expected: setOutputExpected(),
		})

	for index, testCaseDataInput := range dataInputSubTask {
		if index == 0 {
			t.Run(testCaseDataInput.name, func(t *testing.T) {
				result := ClientService.ClientService.RemoveDataRegistered(testCaseDataInput.request.InputStruct, testCaseDataInput.request.ModelStruct)
				assert.NotNil(t, result.(in.ClientRequest), "Result remove process is nil")
			})
		} else if index == 1 {
			t.Run(testCaseDataInput.name, func(t *testing.T) {
				result := ClientService.ClientService.RemoveDataRegistered(testCaseDataInput.request.InputStruct, testCaseDataInput.request.ModelStruct)
				assert.Equal(t, testCaseDataInput.expected, result, "Result is not same")
			})
		}
	}
}
