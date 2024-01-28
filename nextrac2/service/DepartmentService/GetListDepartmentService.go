package DepartmentService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"strings"
)

func (input departmentService) GetListDepartment(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	//--- Reset When DDL
	path := request.URL.Path
	isDDL := strings.Contains(path, "ddl")
	if isDDL {
		contextModel.LimitedByCreatedBy = 0
	}

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListDepartmentValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListDepartment(inputStruct, searchByParam)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input departmentService) doGetListDepartment(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
	)

	dbResult, err = dao.DepartmentDAO.GetListDepartment(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, 0)
	if err.Error != nil {
		return
	}

	output = input.convertToDTOOut(dbResult)
	return
}

func (input departmentService) convertToDTOOut(dbResult []interface{}) (result []out.ListDepartmentDTOOut) {
	for _, item := range dbResult {
		repo := item.(repository.DepartmentModel)
		result = append(result, out.ListDepartmentDTOOut{
			ID:          repo.ID.Int64,
			Name:        repo.Name.String,
			Description: repo.Description.String,
			UpdatedAt:   repo.UpdatedAt.Time,
		})
	}
	return result
}
