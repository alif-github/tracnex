package EmployeeFacilitiesActiveService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"time"
)

type employeeMatrixDeleteService struct {
	FileName string
	service.AbstractService
}

var EmployeeMatrixDeleteService = employeeMatrixDeleteService{FileName: "EmployeeMatrixDeleteService.go"}

func (input employeeMatrixDeleteService) DeleteEmployeeMatrix(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "DeleteEmployeeMatrix"

	params := request.URL.Query()
	levelId := params.Get("levelID")
	gradeId := params.Get("gradeID")

	if levelId == ""{
		err = errorModel.GenerateEmptyFieldError(input.FileName, funcName, "levelID")
		return
	}

	if gradeId == ""{
		err = errorModel.GenerateEmptyFieldError(input.FileName, funcName, "gradeID")
		return
	}

	levelID, errs := strconv.ParseInt(levelId, 10, 64)
	if errs!= nil {
		err = errorModel.GenerateFormatFieldError(input.FileName, funcName, "levekID")
		return
	}
	gradeID, errs := strconv.ParseInt(gradeId, 10, 64)
	if errs!= nil {
		err = errorModel.GenerateFormatFieldError(input.FileName, funcName, "gradeID")
		return
	}

	count, err := dao.EmployeeFacilitiesActiveDAO.GetDetailEmployeeMatrixForUpdate(serverconfig.ServerAttribute.DBConnection, levelID, gradeID)
	if err.Error != nil {
		return
	}

	if count == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "matrix")
		return
	}

	matrixRepo := input.getEmpMatrixRepository(levelID, gradeID, context.AuthAccessTokenModel)

	err = dao.EmployeeFacilitiesActiveDAO.DeleteEmployeeMatrix(serverconfig.ServerAttribute.DBConnection, matrixRepo)
	if err.Error != nil {
		return
	}

	output.Status = service.GetResponseMessages("SUCCESS_DELETE_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeMatrixDeleteService) getEmpMatrixRepository(level int64, grade int64, authAccessToken model2.AuthAccessTokenModel) repository.EmployeeFacilitiesActiveModel {
	tm := time.Now()
	return repository.EmployeeFacilitiesActiveModel{
		LevelID:       sql.NullInt64{Int64: level},
		GradeID:       sql.NullInt64{Int64: grade},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: tm},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		Deleted:       sql.NullBool{Bool: true},
	}
}
