package DepartmentService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input departmentService) ViewDepartment(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.DepartmentRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewDepartment)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewDepartment(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input departmentService) validateViewDepartment(inputStruct *in.DepartmentRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}

func (input departmentService) doViewDepartment(inputStruct in.DepartmentRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		funcName = "doViewDepartment"
	)

	departmentOnDB, err := dao.DepartmentDAO.ViewDepartment(serverconfig.ServerAttribute.DBConnection, repository.DepartmentModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	})
	if err.Error != nil {
		return
	}

	if departmentOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	result = input.convertModelToResponseDetail(departmentOnDB)

	return
}

func (input departmentService) convertModelToResponseDetail(inputModel repository.DepartmentModel) out.ViewDepartmentResponse {
	return out.ViewDepartmentResponse{
		ID:             inputModel.ID.Int64,
		DepartmentName: inputModel.Name.String,
		Description:    inputModel.Description.String,
		UpdatedAt:      inputModel.UpdatedAt.Time,
	}
}
