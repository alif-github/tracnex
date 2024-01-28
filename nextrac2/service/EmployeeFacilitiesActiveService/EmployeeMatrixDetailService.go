package EmployeeFacilitiesActiveService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type employeeMatrixDetailService struct {
	FileName string
	service.AbstractService
}

var EmployeeMatrixDetailService = employeeMatrixDetailService{FileName: "EmployeeMatrixDetailService.go"}

func (input employeeMatrixDetailService) DetailEmployeeMatrix(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "DetailEmployeeMatrix"

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
		err = errorModel.GenerateFormatFieldError(input.FileName, funcName, "levelID")
		return
	}
	gradeID, errs := strconv.ParseInt(gradeId, 10, 64)
	if errs!= nil {
		err = errorModel.GenerateFormatFieldError(input.FileName, funcName, "gradeID")
		return
	}

	matrixs, err := dao.EmployeeFacilitiesActiveDAO.GetDetailEmployeeMatrix(serverconfig.ServerAttribute.DBConnection, levelID, gradeID)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.getEmpMatrixRepository(matrixs, levelID, gradeID)

	output.Status = service.GetResponseMessages("SUCCESS_GET_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeMatrixDetailService) getEmpMatrixRepository(matrix repository.EmployeeFacilitiesActiveModel, levelID int64, gradeID int64) out.EmployeeMatrixResponse {
	var resultAllowance []out.EmployeeAllowanceResponse
	var resultBenefit   []out.EmployeeBenefitMasterResponse
	allowances, _:= dao.EmployeeAllowanceDAO.GetAllowanceForDetail(serverconfig.ServerAttribute.DBConnection)
	benefits, _:= dao.EmployeeMasterBenefitDAO.GetBenefitForDetail(serverconfig.ServerAttribute.DBConnection)
	matrixs, _:= dao.EmployeeFacilitiesActiveDAO.GetMatrixForDetail(serverconfig.ServerAttribute.DBConnection, gradeID, levelID)

	for _, item := range allowances {
		resultAllowance = append(resultAllowance, out.EmployeeAllowanceResponse{
			ID:                item.ID.Int64,
			AllowanceName:     item.AllowanceName.String,
			AllowanceType:     item.Type.String,
			UpdatedAt:         item.UpdatedAt.Time,
		})
	}

	for _, benefit := range benefits {
		resultBenefit = append(resultBenefit, out.EmployeeBenefitMasterResponse{
			ID:                benefit.ID.Int64,
			BenefitName:       benefit.BenefitName.String,
			BenefitType:       benefit.BenefitType.String,
			UpdatedAt:         benefit.UpdatedAt.Time,
		})
	}

	if gradeID != 0 && levelID != 0{
		for i:=0; i<len(matrixs); i++  {
			for a:=0; a<len(resultAllowance);a++  {
				if matrixs[i].AllowanceID.Int64 == resultAllowance[a].ID{
					resultAllowance[a].Active = matrixs[i].Active.Bool
					resultAllowance[a].Value = matrixs[i].Value.String
					break
				}
			}

			for a:=0; a<len(resultBenefit);a++  {
				if matrixs[i].BenefitID.Int64 == resultBenefit[a].ID{
					resultBenefit[a].Active = matrixs[i].Active.Bool
					resultBenefit[a].Value = matrixs[i].Value.String
					break
				}
			}
		}
	}

	return out.EmployeeMatrixResponse{
		LevelID:       matrix.LevelID.Int64,
		Level:         matrix.Level.String,
		GradeID:       matrix.GradeID.Int64,
		Grade:         matrix.Grade.String,
		AllowanceList: resultAllowance,
		BenefitList:   resultBenefit,
	}
}
