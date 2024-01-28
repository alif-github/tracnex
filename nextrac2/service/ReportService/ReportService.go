package ReportService

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"strings"
	"time"
)

type reportService struct {
	ReportDAO dao.ReportDAOInterface
	service.AbstractService
	service.GetListData
}

var ReportService = reportService{}.New()

func (input reportService) New() (output reportService) {
	output.FileName = "ReportService.go"
	output.ServiceName = constanta.ReportConstanta
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{"name", "nik"}
	output.ValidSearchBy = []string{
		"department",
		"sprint",
	}
	output.ReportDAO = dao.ReportDAO
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.EmployeeDataScope] = applicationModel.MappingScopeDB{
		View:  "e.id",
		Count: "e.id",
	}
	return
}

func (input reportService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ReportRequest) errorModel.ErrorModel) (inputStruct in.ReportRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidate"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	err = validation(&inputStruct)
	return
}

func (input reportService) doGetDataFromRedmine(inputStruct in.GetListDataDTO, searchByParam *[]in.SearchByParam) (departmentID int64, resultDBDevTester []repository.RedmineModel, resultDBInfraDevOps []repository.RedmineInfraModel, err errorModel.ErrorModel) {
	var (
		fileName             = "GetListReportService.go"
		funcName             = "doGetDataFromRedmine"
		db                   = serverconfig.ServerAttribute.DBConnection
		dbRedmine            = serverconfig.ServerAttribute.RedmineDBConnection
		dbRedmineInfraDevOps = serverconfig.ServerAttribute.RedmineInfraDBConnection
		isEmployeeExist      bool
		redmineID            []int64
	)

	for _, itemSearchByParam := range *searchByParam {
		//--- Check Department
		if itemSearchByParam.SearchKey == constanta.FilterDepartment {
			departmentIDTemp, errorS := strconv.Atoi(itemSearchByParam.SearchValue)
			if errorS != nil {
				err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
				return
			}
			departmentID = int64(departmentIDTemp)
		}

		//--- Check Is Employee Exist
		if itemSearchByParam.SearchKey == "id" {
			isEmployeeExist = true
		}
	}

	if !isEmployeeExist {
		redmineID, err = dao.EmployeeDAO.GetRedmineIDEmployeeByDepartmentID(db, departmentID)
		if err.Error != nil {
			return
		}
	}

	switch departmentID {
	case constanta.DeveloperDepartmentID, constanta.UIUXDepartmentID:
		resultDBDevTester, err = dao.RedmineDAO.GetTaskDeveloperOnRedmine(dbRedmine, inputStruct, searchByParam, redmineID)
		if err.Error != nil {
			return
		}
	case constanta.QAQCDepartmentID:
		resultDBDevTester, err = dao.RedmineDAO.GetTaskQAOnRedmine(dbRedmine, inputStruct, searchByParam, redmineID)
		if err.Error != nil {
			return
		}
	case constanta.InfraDepartmentID, constanta.DevOpsDepartmentID:
		resultDBInfraDevOps, err = dao.RedmineDAO.GetTaskInfraOnRedmine(dbRedmineInfraDevOps, inputStruct, departmentID, searchByParam, redmineID)
		if err.Error != nil {
			return
		}

		for i := 0; i < len(resultDBInfraDevOps); i++ {
			var (
				c             = resultDBInfraDevOps[i].Category.String
				cArray        = strings.Split(c, "|")
				resultManhour repository.StandarManhourModel
			)

			//--- Get Standar Manhour By Case
			if len(cArray) > 0 {
				resultManhour, err = dao.StandarManhourDAO.GetStandarManhourByCase(db, departmentID, cArray)
				if err.Error != nil {
					return
				}
			}

			//--- Input Standard Time
			resultDBInfraDevOps[i].Manhour.Float64 = resultManhour.Manhour.Float64
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) validateAddParam(request *http.Request, inputStruct *in.GetListDataDTO, searchBy *[]in.SearchByParam) (isMandatoryExist bool, departmentID int64, err errorModel.ErrorModel) {
	var (
		fileName          = "GetListReportService.go"
		funcName          = "validateAddParam"
		db                = serverconfig.ServerAttribute.DBConnection
		idxDeleted        = -1
		employeeIDBundle  []int64
		errorS            error
		isSprintExist     bool
		isDepartmentExist bool
		//sprintBundle     []string
	)

	employeeIDParam := service.GenerateQueryValue(request.URL.Query()[constanta.FilterEmployee])
	if !util.IsStringEmpty(employeeIDParam) {
		//-- Employee Unmarshal
		errorS = json.Unmarshal([]byte(employeeIDParam), &employeeIDBundle)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
			return
		}

		*searchBy = append(*searchBy, in.SearchByParam{
			DataType:       "char",
			SearchKey:      "id",
			SearchOperator: "eq",
			SearchValue:    employeeIDParam,
			SearchType:     constanta.Filter,
		})
	}

	for i := 0; i < len(*searchBy); i++ {
		if (*searchBy)[i].SearchKey == constanta.FilterSprint {
			var (
				timeStart time.Time
				timeEnd   time.Time
				sprint    []string
			)

			//-- Empty Then Error
			if util.IsStringEmpty((*searchBy)[i].SearchValue) {
				err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Sprint)
				return
			}

			//-- Split By Sign
			sprint = strings.Split((*searchBy)[i].SearchValue, "-")
			if len(sprint) > 2 {
				err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.SprintRule1, constanta.Sprint, "")
				return
			}

			//-- Test Time Start
			timeStart, errorS = time.Parse(constanta.DefaultTimeSprintFormat, sprint[0])
			if errorS != nil {
				err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.SprintRule2, constanta.SprintStart, "")
				return
			}

			//-- Test Time End
			timeEnd, errorS = time.Parse(constanta.DefaultTimeSprintFormat, sprint[1])
			if errorS != nil {
				err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.SprintRule2, constanta.SprintEnd, "")
				return
			}

			//-- Test Time Start Before Time End
			if !timeStart.Before(timeEnd) && !timeStart.Equal(timeEnd) {
				err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.SprintRule3, constanta.Sprint, "")
				return
			}

			inputStruct.UpdatedAtStart = timeStart
			inputStruct.UpdatedAtEnd = timeEnd.Add(24 * time.Hour)

			//-- Idx Deleted Check
			if idxDeleted > 0 {
				err = errorModel.GenerateUnknownError(fileName, funcName, errors.New("double key in filter [sprint]"))
				return
			}

			idxDeleted = i
			isSprintExist = true
		} else if (*searchBy)[i].SearchKey == constanta.FilterDepartment {
			var (
				departmentIDInt       int
				isDepartmentExistOnDB bool
			)

			//-- Department Can't Empty
			if util.IsStringEmpty((*searchBy)[i].SearchValue) {
				err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.DepartmentConstanta)
				return
			}

			//-- Should Be Converted To Int
			departmentIDInt, errorS = strconv.Atoi((*searchBy)[i].SearchValue)
			if errorS != nil {
				err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.DepartmentRule1, constanta.DepartmentConstanta, "")
				return
			}

			//-- Check ID Department
			isDepartmentExistOnDB, err = dao.DepartmentDAO.CheckIDDepartment(db, repository.DepartmentModel{ID: sql.NullInt64{Int64: int64(departmentIDInt)}})
			if err.Error != nil {
				return
			}

			//-- Err When Not Found In DB
			if !isDepartmentExistOnDB {
				err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.DepartmentConstanta)
				return
			}

			isDepartmentExist = true
			departmentID = int64(departmentIDInt)
		}

		//-- Check Mandatory Exist
		if len(*searchBy)-(i+1) == 0 && (isSprintExist || isDepartmentExist) {
			isMandatoryExist = true
		}
	}

	//-- Deleted Search By
	if idxDeleted >= 0 {
		*searchBy = append((*searchBy)[:idxDeleted], (*searchBy)[idxDeleted+1:]...)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) validateDataScope(contextModel *applicationModel.ContextModel) (scope map[string]interface{}, err errorModel.ErrorModel) {
	scope, err = input.ValidateMultipleDataScope(contextModel, []string{
		constanta.EmployeeDataScope,
	})

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) ServicePreparedDBCustomize(funcName string, db *sql.DB, inputStruct interface{}, contextModel *applicationModel.ContextModel, serve func(*sql.Tx, interface{}, *applicationModel.ContextModel, time.Time) (interface{}, errorModel.ErrorModel)) (output interface{}, err errorModel.ErrorModel) {
	var (
		errs    error
		tx      *sql.Tx
		timeNow = time.Now()
	)

	defer func() {
		if errs != nil || err.Error != nil {
			errs = tx.Rollback()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
			}
		} else {
			errs = tx.Commit()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
			}
		}
	}()

	tx, errs = db.Begin()
	if errs != nil {
		return
	}

	output, err = serve(tx, inputStruct, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
