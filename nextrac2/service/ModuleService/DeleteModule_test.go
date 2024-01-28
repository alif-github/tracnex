package ModuleService

import (
	"database/sql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"testing"
	"time"
)

func TestModuleService_DeleteModule_ValidateModuleOnDB(t *testing.T) {
	// Open Mock
	ctrl := gomock.NewController(t)

	// Call Mock
	mocks := dao.NewMockModuleDAOInterface(ctrl)
	service := &moduleService{ModuleDAO: mocks}

	// Initial Set
	timeNow := time.Now()
	defaultFileName := service.FileName
	defaultFuncName := "validateModuleOnDB"
	defaultModuleModel := repository.ModuleModel{
		ID:         sql.NullInt64{Int64: 99},
		UpdatedAt:  sql.NullTime{Time: timeNow},
		CreatedBy:  sql.NullInt64{Int64: 1},
		IsUsed:     sql.NullBool{Bool: false},
		ModuleName: sql.NullString{String: "Module Mock 1"},
	}

	cases := []struct {
		name           string
		inputMock1     repository.ModuleModel
		outputMock1    repository.ModuleModel
		outputMock2    errorModel.ErrorModel
		inputService1  in.ModuleRequest
		inputService2  repository.ModuleModel
		inputService3  applicationModel.ContextModel
		expectedError  errorModel.ErrorModel
		expectedResult repository.ModuleModel
	}{
		{
			name:           "Case 1 Positive",
			inputMock1:     repository.ModuleModel{ID: sql.NullInt64{Int64: 1}},
			outputMock1:    defaultModuleModel,
			outputMock2:    errorModel.GenerateNonErrorModel(),
			inputService1:  in.ModuleRequest{UpdatedAt: timeNow},
			inputService2:  repository.ModuleModel{ID: sql.NullInt64{Int64: 1}},
			inputService3:  applicationModel.ContextModel{},
			expectedError:  errorModel.GenerateNonErrorModel(),
			expectedResult: defaultModuleModel,
		},
		{
			name:           "Case 2 Negative Not Found",
			inputMock1:     repository.ModuleModel{ID: sql.NullInt64{Int64: 1}},
			outputMock1:    repository.ModuleModel{},
			outputMock2:    errorModel.GenerateNonErrorModel(),
			inputService1:  in.ModuleRequest{UpdatedAt: timeNow},
			inputService2:  repository.ModuleModel{ID: sql.NullInt64{Int64: 1}},
			inputService3:  applicationModel.ContextModel{},
			expectedError:  errorModel.GenerateUnknownDataError(defaultFileName, defaultFuncName, constanta.Module),
			expectedResult: repository.ModuleModel{},
		},
		{
			name:       "Case 3 Negative Module Used",
			inputMock1: repository.ModuleModel{ID: sql.NullInt64{Int64: 1}},
			outputMock1: repository.ModuleModel{
				ID:         sql.NullInt64{Int64: 99},
				UpdatedAt:  sql.NullTime{Time: timeNow},
				CreatedBy:  sql.NullInt64{Int64: 1},
				IsUsed:     sql.NullBool{Bool: true},
				ModuleName: sql.NullString{String: "Module Mock 1"},
			},
			outputMock2:   errorModel.GenerateNonErrorModel(),
			inputService1: in.ModuleRequest{UpdatedAt: timeNow},
			inputService2: repository.ModuleModel{ID: sql.NullInt64{Int64: 1}},
			inputService3: applicationModel.ContextModel{},
			expectedError: errorModel.GenerateDataUsedError(defaultFileName, defaultFuncName, constanta.Module),
			expectedResult: repository.ModuleModel{
				ID:         sql.NullInt64{Int64: 99},
				UpdatedAt:  sql.NullTime{Time: timeNow},
				CreatedBy:  sql.NullInt64{Int64: 1},
				IsUsed:     sql.NullBool{Bool: true},
				ModuleName: sql.NullString{String: "Module Mock 1"},
			},
		},
		{
			name:       "Case 4 Negative Own Edit",
			inputMock1: repository.ModuleModel{ID: sql.NullInt64{Int64: 1}},
			outputMock1: repository.ModuleModel{
				ID:         sql.NullInt64{Int64: 99},
				UpdatedAt:  sql.NullTime{Time: timeNow},
				CreatedBy:  sql.NullInt64{Int64: 1},
				IsUsed:     sql.NullBool{Bool: true},
				ModuleName: sql.NullString{String: "Module Mock 1"},
			},
			outputMock2:   errorModel.GenerateNonErrorModel(),
			inputService1: in.ModuleRequest{UpdatedAt: timeNow},
			inputService2: repository.ModuleModel{ID: sql.NullInt64{Int64: 1}},
			inputService3: applicationModel.ContextModel{LimitedByCreatedBy: 99},
			expectedError: errorModel.GenerateForbiddenAccessClientError("AbstractService.go", "checkUserLimitedByOwnAccess"),
			expectedResult: repository.ModuleModel{
				ID:         sql.NullInt64{Int64: 99},
				UpdatedAt:  sql.NullTime{Time: timeNow},
				CreatedBy:  sql.NullInt64{Int64: 1},
				IsUsed:     sql.NullBool{Bool: true},
				ModuleName: sql.NullString{String: "Module Mock 1"},
			},
		},
		{
			name:       "Case 5 Lock Updated At",
			inputMock1: repository.ModuleModel{ID: sql.NullInt64{Int64: 1}},
			outputMock1: repository.ModuleModel{
				ID:         sql.NullInt64{Int64: 99},
				UpdatedAt:  sql.NullTime{Time: timeNow},
				CreatedBy:  sql.NullInt64{Int64: 1},
				IsUsed:     sql.NullBool{Bool: false},
				ModuleName: sql.NullString{String: "Module Mock 1"},
			},
			outputMock2:   errorModel.GenerateNonErrorModel(),
			inputService1: in.ModuleRequest{UpdatedAt: timeNow.Add(-1 * time.Hour)},
			inputService2: repository.ModuleModel{ID: sql.NullInt64{Int64: 1}},
			inputService3: applicationModel.ContextModel{},
			expectedError: errorModel.GenerateDataLockedError(defaultFileName, defaultFuncName, constanta.Module),
			expectedResult: repository.ModuleModel{
				ID:         sql.NullInt64{Int64: 99},
				UpdatedAt:  sql.NullTime{Time: timeNow},
				CreatedBy:  sql.NullInt64{Int64: 1},
				IsUsed:     sql.NullBool{Bool: false},
				ModuleName: sql.NullString{String: "Module Mock 1"},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mocks.EXPECT().GetModuleForUpdate(nil, tc.inputMock1).Return(tc.outputMock1, tc.outputMock2)
			result, err := service.validateModuleOnDB(nil, tc.inputService1, tc.inputService2, &tc.inputService3)
			assert.Equal(t, tc.expectedError, err, "Error Must Empty")
			assert.Equal(t, tc.expectedResult, result, "Result Must Equal")
		})
	}
}
