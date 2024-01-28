package AbsentService

import (
	"database/sql"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io"
	"mime/multipart"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strings"
	"time"
)

type absentService struct {
	service.AbstractService
	service.GetListData
}

var AbsentService = absentService{}.New()

func (input absentService) New() (output absentService) {
	output.ServiceName = "Absent"
	output.FileName = "AbsentService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidSearchBy = []string{
		"id_card",
		"name",
		"period",
	}
	output.ValidOrderBy = []string{
		"id_card",
		"name",
	}
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.EmployeeDataScope] = applicationModel.MappingScopeDB{
		View:  "a.employee_id",
		Count: "a.employee_id",
	}
	return
}

func (input absentService) readAndHandleUploadExcel(request *http.Request, key string) (xlsx *excelize.File, err errorModel.ErrorModel) {
	var (
		funcName = "readAndHandleUploadExcel"
		errors   error
		file     multipart.File
		excel    []byte
	)

	//--- Get File Excel From Form Data
	errors = request.ParseMultipartForm(10 << 20)
	if errors != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errors)
		return
	}

	file, _, errors = request.FormFile(key)
	if errors != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errors)
		return
	}

	defer func() {
		_ = file.Close()
	}()

	//--- Read Excel
	excel, errors = io.ReadAll(file)
	if errors != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errors)
		return
	}

	//--- Excel Open File
	xlsx, errors = excelize.OpenReader(strings.NewReader(string(excel)))
	if errors != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errors)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input absentService) validateDataScope(contextModel *applicationModel.ContextModel) (scope map[string]interface{}, err errorModel.ErrorModel) {
	scope, err = input.ValidateMultipleDataScope(contextModel, []string{
		constanta.EmployeeDataScope,
	})

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input absentService) PeriodCheck(searchByParam *[]in.SearchByParam) (err errorModel.ErrorModel) {
	var (
		funcName      = "PeriodCheck"
		db            = serverconfig.ServerAttribute.DBConnection
		isPeriodExist bool
		periodLast    repository.AbsentModel
	)

	//--- Get Last Period Absent
	periodLast, err = dao.AbsentDAO.GetLastPeriodAbsent(db)
	if err.Error != nil {
		return
	}

	if len(*searchByParam) < 1 {
		//--- Add New Period
		if !periodLast.PeriodStart.Time.IsZero() && !periodLast.PeriodEnd.Time.IsZero() {
			input.addedNewPeriodSearchByParam(periodLast, searchByParam)
		}
	} else {
		//--- Search By Param Period
		for idx, item := range *searchByParam {
			if item.SearchKey == "period" {
				var (
					startTime time.Time
					endTime   time.Time
					errors    error
				)

				isPeriodExist = true
				p := strings.Split(item.SearchValue, "-")
				if len(p) != 2 {
					err = errorModel.GenerateFormatFieldError(input.FileName, funcName, "Period")
					return
				}

				//--- Start Time
				startTime, errors = time.Parse(constanta.DefaultTimeSprintFormat, p[0])
				if errors != nil {
					err = errorModel.GenerateUnknownError(input.FileName, funcName, errors)
					return
				}

				//--- End Time
				endTime, errors = time.Parse(constanta.DefaultTimeSprintFormat, p[1])
				if errors != nil {
					err = errorModel.GenerateUnknownError(input.FileName, funcName, errors)
					return
				}

				//--- Delete Period
				*searchByParam = append((*searchByParam)[:idx], (*searchByParam)[idx+1:]...)
				input.addedNewPeriodSearchByParam(repository.AbsentModel{
					PeriodStart: sql.NullTime{Time: startTime},
					PeriodEnd:   sql.NullTime{Time: endTime},
				}, searchByParam)
				break
			}
		}

		if !isPeriodExist {
			//--- Add New Period
			if !periodLast.PeriodStart.Time.IsZero() && !periodLast.PeriodEnd.Time.IsZero() {
				input.addedNewPeriodSearchByParam(periodLast, searchByParam)
			}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
