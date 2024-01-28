package product

import (
	"github.com/stretchr/testify/assert"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service/ProductService"
	"testing"
)

type testCase struct {
	name     string
	request  request
	expected errorModel.ErrorModel
}

type request struct {
	CollectionModule arrayCollModule
}

type arrayCollModule struct {
	arrayCollModuleItem []int64
}

func setRequestDuplicate() (arrayCollValueItem []arrayCollModule) {
	var case1, case2, case3, case4, case5 []int64

	case1 = []int64{1, 0, 0, 1, 2, 0, 0, 0, 0, 0}
	case2 = []int64{0, 0, 0, 0, 0, 0, 0, 2, 5, 2}
	case3 = []int64{1, 1, 1, 0, 0, 0, 3, 0, 1, 1}
	case4 = []int64{0, 0, 0, 1, 1, 1, 0, 0, 0, 0}
	case5 = []int64{0, 0, 0, 0, 0, 0, 1, 1, 0, 0}

	//----------------------- Appending
	arrayCollValueItem = append(arrayCollValueItem,
		arrayCollModule{arrayCollModuleItem: case1},
		arrayCollModule{arrayCollModuleItem: case2},
		arrayCollModule{arrayCollModuleItem: case3},
		arrayCollModule{arrayCollModuleItem: case4},
		arrayCollModule{arrayCollModuleItem: case5})

	return
}

func setRequestNotDuplicate() (arrayCollValueItem []arrayCollModule) {
	var case1, case2, case3, case4, case5 []int64

	case1 = []int64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	case2 = []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	case3 = []int64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	case4 = []int64{0, 0, 0, 1, 0, 0, 0, 0, 0, 0}
	case5 = []int64{0, 0, 0, 0, 0, 0, 0, 2, 0, 0}

	//----------------------- Appending
	arrayCollValueItem = append(arrayCollValueItem,
		arrayCollModule{arrayCollModuleItem: case1},
		arrayCollModule{arrayCollModuleItem: case2},
		arrayCollModule{arrayCollModuleItem: case3},
		arrayCollModule{arrayCollModuleItem: case4},
		arrayCollModule{arrayCollModuleItem: case5})

	return
}

func setOutputExpectedDuplicate() (err errorModel.ErrorModel) {
	err = errorModel.GenerateDataDuplicateInDTOError("ProductService.go", "TracAndCheckDuplicateModule")
	return
}

func setOutputExpectedNotDuplicate() (err errorModel.ErrorModel) {
	err = errorModel.GenerateNonErrorModel()
	return
}

func TestProductService_checkDuplicateSuccess(t *testing.T) {
	var dataInput []testCase

	arrayCase := setRequestDuplicate()
	 for _, valueArrayCase := range arrayCase {
		 dataInput = append(dataInput, testCase{
			 name:     "1. Check duplicate result, then error",
			 request:  request{CollectionModule: valueArrayCase},
			 expected: setOutputExpectedDuplicate(),
		 }, testCase{
			 name:     "2. Check duplicate result, then error",
			 request:  request{CollectionModule: valueArrayCase},
			 expected: setOutputExpectedDuplicate(),
		 }, testCase{
			 name:     "3. Check duplicate result, then error",
			 request:  request{CollectionModule: valueArrayCase},
			 expected: setOutputExpectedDuplicate(),
		 }, testCase{
			 name:     "4. Check duplicate result, then error",
			 request:  request{CollectionModule: valueArrayCase},
			 expected: setOutputExpectedDuplicate(),
		 }, testCase{
			 name:     "5. Check duplicate result, then error",
			 request:  request{CollectionModule: valueArrayCase},
			 expected: setOutputExpectedDuplicate(),
		 })
	 }

	for _, testCaseDataInput := range dataInput {
		t.Run(testCaseDataInput.name, func(t *testing.T) {
			result := ProductService.ProductService.TracAndCheckDuplicateModule(testCaseDataInput.request.CollectionModule.arrayCollModuleItem)
			assert.Equal(t, testCaseDataInput.expected, result, "Result Correct")
		})
	}
}

func TestProductService_checkNormalSuccess(t *testing.T) {
	var dataInput []testCase

	arrayCase := setRequestNotDuplicate()
	for _, valueArrayCase := range arrayCase {
		dataInput = append(dataInput, testCase{
			name:     "1. Check not duplicate result, then not error",
			request:  request{CollectionModule: valueArrayCase},
			expected: setOutputExpectedNotDuplicate(),
		}, testCase{
			name:     "2. Check not duplicate result, then not error",
			request:  request{CollectionModule: valueArrayCase},
			expected: setOutputExpectedNotDuplicate(),
		}, testCase{
			name:     "3. Check not duplicate result, then not error",
			request:  request{CollectionModule: valueArrayCase},
			expected: setOutputExpectedNotDuplicate(),
		}, testCase{
			name:     "4. Check not duplicate result, then not error",
			request:  request{CollectionModule: valueArrayCase},
			expected: setOutputExpectedNotDuplicate(),
		}, testCase{
			name:     "5. Check not duplicate result, then not error",
			request:  request{CollectionModule: valueArrayCase},
			expected: setOutputExpectedNotDuplicate(),
		})
	}

	for _, testCaseDataInput := range dataInput {
		t.Run(testCaseDataInput.name, func(t *testing.T) {
			result := ProductService.ProductService.TracAndCheckDuplicateModule(testCaseDataInput.request.CollectionModule.arrayCollModuleItem)
			assert.Equal(t, testCaseDataInput.expected, result, "Result Not Duplicate Correct")
		})
	}
}