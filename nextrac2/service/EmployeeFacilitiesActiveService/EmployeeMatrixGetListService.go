package EmployeeFacilitiesActiveService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
)

type employeeMatrixGetListService struct {
	service.GetListData
	FileName string
}

var EmployeeMatrixGetListService = employeeMatrixGetListService{}.New()

func (input employeeMatrixGetListService) New() (output employeeMatrixGetListService) {
	output.FileName = "EmployeeMatrixGetListService.go"
	output.ValidSearchBy = []string{""}
	output.ValidOrderBy = []string{"employee_facilities_active.employee_level_id", "employee_facilities_active.employee_grade_id", "lvl.level", "grade.grade"}
	output.ValidLimit = service.DefaultLimit
	return
}

func (input employeeMatrixGetListService) GetEmployeeMatrix(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam
	var createdBy int64 = 0

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListEmployeeMatrixValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	createdBy, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*context)
	if !isOnlyHaveOwnAccess {
		createdBy = 0
	}

	params := request.URL.Query()
	key := params.Get("key")
	keyword := params.Get("keyword")

	matrixs, err := dao.EmployeeFacilitiesActiveDAO.GetListMatrixs(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, createdBy, key, keyword)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.convertRepoToDTO(matrixs)
	output.Status = service.GetResponseMessages("SUCCESS_GET_LIST_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeMatrixGetListService) InitiateEmployeeMatrix(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	_, _, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListEmployeeValidOperator)
	if err.Error != nil {
		return
	}

	params := request.URL.Query()
	key := params.Get("key")
	keyword := params.Get("keyword")

	countData, err := dao.EmployeeFacilitiesActiveDAO.GetCountEmployeeMatrix(serverconfig.ServerAttribute.DBConnection, key, keyword)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		CountData:     int(countData),
		ValidOperator: applicationModel.GetListEmployeeMatrixValidOperator,
		ValidSearchParam: []out.SearchByParam{
			{
				Key:   "lvl.level",
				Value: "Level",
			},
			{
				Key:   "grade.grade",
				Value: "Grade",
			},
		},
	}
	output.Status = service.GetResponseMessages("SUCCESS_INITIATE_MESSAGE", context)

	return
}

func (input employeeMatrixGetListService) convertRepoToDTO(data []interface{}) (matrixs []out.EmployeeMatrixForView) {
	for _, item := range data {
		matrix := item.(repository.EmployeeFacilitiesActiveModel)
		matrixs = append(matrixs, out.EmployeeMatrixForView{
			Grade:   matrix.Grade.String,
			GradeID: matrix.GradeID.Int64,
			Level:   matrix.Level.String,
			LevelID: matrix.LevelID.Int64,
		})
	}
	return
}
