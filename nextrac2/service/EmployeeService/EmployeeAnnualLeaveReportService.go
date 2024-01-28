package EmployeeService

import (
	"database/sql"
	"fmt"
	"net/http"

	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
)

func (input employeeService) DownloadAnnualReport(request *http.Request, contextModel *applicationModel.ContextModel) (output []byte, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validSearchBy = []string{"first_name", "last_name", "el.type"}
		validOrderBy  = []string{"id", "created_at"}
		key = request.URL.Query().Get("key")
		keyword = request.URL.Query().Get("keyword")
		year = request.URL.Query().Get("year")
	)

	inputStruct, searchByParam, errModel = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, applicationModel.GetListEmployeeLeaveValidOperator, service.DefaultLimit)
	if errModel.Error != nil {
		return
	}

	results, errModel := dao.EmployeeLeaveDAO.GetEmployeeLeaveYearly(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, repository.EmployeeLeaveModel{
		SearchBy: sql.NullString{String:key},
		Keyword:  sql.NullString{String:keyword},
		IsYearly: sql.NullBool{Bool:true},
		Year:     sql.NullString{String:year},
	})

	if errModel.Error != nil {
		return
	}

	fmt.Println(results)

	return output, header, errorModel.ErrorModel{}
}