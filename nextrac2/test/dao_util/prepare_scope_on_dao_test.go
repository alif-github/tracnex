package dao_util

import (
	"github.com/stretchr/testify/assert"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"testing"
)

type testCase struct {
	name     string
	request  request
	expected string
}

type request struct {
	scope           preparationInitial
	additionalWhere string
	idxStart        int
	keyDataScope    string
	isView          bool
}

type preparationInitial struct {
	scopeLimit map[string]interface{}
	scopeDB    map[string]applicationModel.MappingScopeDB
}

func casePreparation(inputInterface []interface{}) (requestTemp request) {
	var preparationInit preparationInitial
	var additionalWhereTemp string

	scopeDBTemp := make(map[string]applicationModel.MappingScopeDB)
	scopeDBTemp["nexsoft.valid"] = applicationModel.MappingScopeDB{
		View:  "id",
		Count: "id",
	}

	//------------- Case
	scopeLimitTemp := make(map[string]interface{})
	scopeLimitTemp["nexsoft.valid"] = inputInterface

	preparationInit.scopeDB = scopeDBTemp
	preparationInit.scopeLimit = scopeLimitTemp

	requestTemp = request{
		scope:           preparationInit,
		additionalWhere: additionalWhereTemp,
		idxStart:        1,
		keyDataScope:    "nexsoft.valid",
		isView:          true,
	}

	return
}

func TestDAOUtil_checkResultSuccess(t *testing.T) {
	var dataInput []testCase

	dataInput = append(dataInput, testCase{
		name:     "Test 1 : Data Result Sukses",
		request:  casePreparation([]interface{}{"1", "2", "3"}),
		expected: " id IN (1,2,3)",
	}, testCase{
		name:     "Test 2 : Data Result Sukses",
		request:  casePreparation([]interface{}{"all"}),
		expected: "",
	}, testCase{
		name:     "Test 3 : Data Result Sukses",
		request:  casePreparation([]interface{}{"2"}),
		expected: " id IN (2)",
	})

	for _, testCaseDataInput := range dataInput {
		t.Run(testCaseDataInput.name, func(t *testing.T) {
			additionalWhere := testCaseDataInput.request.additionalWhere
			dao.PrepareScopeOnDAO(testCaseDataInput.request.scope.scopeLimit,
				testCaseDataInput.request.scope.scopeDB, &additionalWhere,
				testCaseDataInput.request.idxStart, testCaseDataInput.request.keyDataScope, testCaseDataInput.request.isView)

			assert.Equal(t, testCaseDataInput.expected, additionalWhere, "Result Correct")
		})
	}
}
