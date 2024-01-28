package EmployeeService

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

func TestEmployeeService_UpdateEmployee_ValidateEmployeeOnDB(t *testing.T) {
	// Open Mock
	ctrl := gomock.NewController(t)

	// Call Mock
	mocks := dao.NewMockEmployeeDAOInterface(ctrl)
	service := &employeeService{EmployeeDAO: mocks}
	timeNow := time.Now()

	cases := []struct {
		name          string
		outputMock1   repository.EmployeeModel
		outputMock2   errorModel.ErrorModel
		inputStruct   in.EmployeeRequest
		model         repository.EmployeeModel
		contextModel  applicationModel.ContextModel
		expectedError errorModel.ErrorModel
	}{
		{
			name: "Case 1 Positive",
			outputMock1: repository.EmployeeModel{
				ID:        sql.NullInt64{Int64: 1},
				UpdatedAt: sql.NullTime{Time: timeNow},
				CreatedBy: sql.NullInt64{Int64: 99},
			},
			outputMock2: errorModel.GenerateNonErrorModel(),
			inputStruct: in.EmployeeRequest{
				ID:           1,
				NIK:          99,
				RedmineId:    99,
				Name:         "NAME_1",
				DepartmentId: 1,
				UpdatedAt:    timeNow,
			},
			contextModel:  applicationModel.ContextModel{},
			expectedError: errorModel.GenerateNonErrorModel(),
		},
		{
			name:        "Case 2 Negative Not Found",
			outputMock2: errorModel.GenerateNonErrorModel(),
			inputStruct: in.EmployeeRequest{
				ID:           1,
				NIK:          99,
				RedmineId:    99,
				Name:         "NAME_1",
				DepartmentId: 1,
				UpdatedAt:    timeNow,
			},
			contextModel:  applicationModel.ContextModel{},
			expectedError: errorModel.GenerateUnknownDataError(service.FileName, "checkAndLockEmployeeOnDB", constanta.EmployeeConstanta),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mocks.EXPECT().GetEmployeeForUpdate(nil, repository.EmployeeModel{}).Return(tc.outputMock1, tc.outputMock2)
			err := service.checkAndLockEmployeeTimeSheetOnDB(tc.inputStruct, tc.model, &tc.contextModel)
			assert.NotNil(t, err, "Error Must Exist Either Positive")
			assert.Equal(t, tc.expectedError, err, "Error Must Equal")
		})
	}
}
