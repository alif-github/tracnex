package ReportService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"sort"
	"time"
)

func (input reportService) GetListReport(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct      in.GetListDataDTO
		searchByParam    []in.SearchByParam
		isMandatoryExist bool
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListReportValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	isMandatoryExist, _, err = input.validateAddParam(request, &inputStruct, &searchByParam)
	if err.Error != nil {
		return
	}

	if isMandatoryExist {
		output.Data.Content, _, _, err = input.doGetListReport(inputStruct, searchByParam, contextModel)
		if err.Error != nil {
			return
		}
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input reportService) doGetListReport(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, departmentID int64, actualDBResult []interface{}, err errorModel.ErrorModel) {
	var (
		db                  = serverconfig.ServerAttribute.DBConnection
		dbRedmine           = serverconfig.ServerAttribute.RedmineDBConnection
		scope               map[string]interface{}
		tempOutput          []out.ReportResponse
		resultDBInfraDevOps []repository.RedmineInfraModel
		resultDBDevTester   []repository.RedmineModel
	)

	//--- Get Data From Redmine
	departmentID, resultDBDevTester, resultDBInfraDevOps, err = input.doGetDataFromRedmine(inputStruct, &searchByParam)
	if err.Error != nil {
		return
	}

	//--- Validate Data Scope
	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	switch departmentID {
	case constanta.DeveloperDepartmentID, constanta.QAQCDepartmentID, constanta.UIUXDepartmentID:
		//--- Get Manhour
		var manhour []repository.HistoryTimeReportRedmineModel
		manhour, err = dao.RedmineDAO.GetManhourOnRedmine(dbRedmine, resultDBDevTester)
		if err.Error != nil {
			return
		}

		//--- Dev QA Get List Report
		actualDBResult, tempOutput, err = input.devQAGetListReport(inputStruct, searchByParam, contextModel, scope, manhour, departmentID)
		if err.Error != nil {
			return
		}
	case constanta.InfraDepartmentID, constanta.DevOpsDepartmentID:
		//--- Infra DevOps Get List
		actualDBResult, tempOutput, err = input.infraDevOpsGetListReport(inputStruct, searchByParam, contextModel, scope, resultDBInfraDevOps, departmentID)
		if err.Error != nil {
			return
		}
	default:
	}

	//-- Re-Arrange Report
	if tempOutput != nil {
		output, err = input.unionAllResultReport(db, tempOutput, inputStruct)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) devQAGetListReport(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel, scope map[string]interface{}, manhour []repository.HistoryTimeReportRedmineModel, departmentID int64) (actualDBResult []interface{}, tempOutput []out.ReportResponse, err errorModel.ErrorModel) {
	var (
		db                 = serverconfig.ServerAttribute.DBConnection
		backlogSearchBy    []in.SearchByParam
		actualSearchBy     []in.SearchByParam
		historySearchBy    []in.SearchByParam
		backlogDBResult    []interface{}
		historyDBResult    []interface{}
		isDataRedmineExist bool
	)

	for _, itemSearchBy := range searchByParam {
		//--- Print Values Redmine
		if itemSearchBy.SearchKey == "redmine" {
			isDataRedmineExist = true
			fmt.Println(fmt.Sprintf(`Ticket Number List -> %v`, itemSearchBy.SearchValue))
		}

		//--- Fill Search By On Each New Search By
		backlogSearchBy = append(backlogSearchBy, itemSearchBy)
		actualSearchBy = append(actualSearchBy, itemSearchBy)
		historySearchBy = append(historySearchBy, itemSearchBy)
	}

	//--- Get List Report Backlog ---> To Actual
	backlogDBResult, err = input.ReportDAO.GetListReport(db, inputStruct, backlogSearchBy, departmentID, 0, true, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	//--- Get List Report Actual ---> To Backlog
	if isDataRedmineExist {
		actualDBResult, historyDBResult, tempOutput, err = input.dataRedmineExistProcess(inputStruct, actualSearchBy, historySearchBy, contextModel, scope, manhour, departmentID)
		if err.Error != nil {
			return
		}
	}

	//--- Convert Report To Response
	tempOutput, err = input.convertModelToResponseGetListDevQA(backlogDBResult, actualDBResult, historyDBResult)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) infraDevOpsGetListReport(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, _ *applicationModel.ContextModel, scope map[string]interface{}, resultDBInfraDevOps []repository.RedmineInfraModel, departmentID int64) (actualDBResult []interface{}, tempOutput []out.ReportResponse, err errorModel.ErrorModel) {
	var (
		db           = serverconfig.ServerAttribute.DBConnection
		ticketNumber []int64
	)

	for _, itemResultDBInfraDevOps := range resultDBInfraDevOps {
		ticketNumber = append(ticketNumber, itemResultDBInfraDevOps.RedmineTicket.Int64)
	}

	//--- Reporting And Print Ticket Number
	fmt.Println(fmt.Sprintf(`Ticket Number List -> %v`, ticketNumber))
	actualDBResult, err = dao.ReportDAO.GetListReportForInfraDevOps(db, inputStruct, resultDBInfraDevOps, searchByParam, departmentID, 0, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	//--- Convert
	tempOutput, err = input.convertModelToResponseGetListInfraDevOps(actualDBResult)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) dataRedmineExistProcess(inputStruct in.GetListDataDTO, actualSearchBy, historySearchBy []in.SearchByParam, _ *applicationModel.ContextModel, scope map[string]interface{}, manhour []repository.HistoryTimeReportRedmineModel, departmentID int64) (actualDBResult, historyDBResult []interface{}, tempOutput []out.ReportResponse, err errorModel.ErrorModel) {
	var (
		fileName = "GetListReportService.go"
		funcName = "dataRedmineExistProcess"
		db       = serverconfig.ServerAttribute.DBConnection
	)

	getAll := in.GetListDataDTO{
		AbstractDTO: in.AbstractDTO{
			Page:    -99,
			Limit:   -99,
			OrderBy: inputStruct.OrderBy,
		},
	}

	//--- History
	historyDBResult, err = input.ReportDAO.GetListReportHistory(db, getAll, historySearchBy, manhour, departmentID, 0, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if historyDBResult != nil {
		for idx, itemHistoryDBResult := range historyDBResult {
			var (
				itemHistory = itemHistoryDBResult.(repository.ReportModel)
				dataHistory []repository.TicketRedmineModel
				sumTime     time.Duration
			)

			_ = json.Unmarshal([]byte(itemHistory.ActualHistoryManday.String), &dataHistory)
			for i := 0; i < len(dataHistory); i++ {
				for j := 0; j < len(dataHistory[i].History); j++ {
					var errorS error
					dataHistory[i].History[j].CreatedOn, errorS = time.Parse(constanta.DefaultDBTimeFormat, dataHistory[i].History[j].CreatedOnStr)
					if errorS != nil {
						fmt.Println(fmt.Sprintf(`[%s, %s] Convert Failed -> %s`, fileName, funcName, dataHistory[i].History[j].CreatedOnStr))
						err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
						return
					}
				}

				//--- Sort data by EndDate in ascending order
				sort.Slice(dataHistory[i].History, func(x, y int) bool {
					return dataHistory[i].History[x].CreatedOn.Before(dataHistory[i].History[y].CreatedOn)
				})
			}

			for i := 0; i < len(dataHistory); i++ {
				var (
					start, pause, end time.Time
					sumTimeTicket     time.Duration
				)

				for j := 0; j < len(dataHistory[i].History); j++ {
					if start.IsZero() && len(dataHistory[i].History)-(j+1) != 0 { //--- Add Start First
						if dataHistory[i].History[j].Subject != constanta.StartTicket {
							continue
						}
						start = dataHistory[i].History[j].CreatedOn
						continue
					}

					if len(dataHistory[i].History)-(j+1) == 0 && dataHistory[i].History[j].Subject != constanta.EndTicket { //--- Last Time
						sumTimeTicket = 0
						break
					}

					if dataHistory[i].History[j].Subject == constanta.PauseTicket { //--- Pause Time
						if pause.IsZero() {
							pause = dataHistory[i].History[j].CreatedOn
							sumTimeTicket += pause.Sub(start)
							start = time.Time{}
						}
						continue
					}

					if dataHistory[i].History[j].Subject == constanta.EndTicket { //--- End Time
						end = dataHistory[i].History[j].CreatedOn
						if (!pause.IsZero() && !start.IsZero()) || (!start.IsZero()) {
							sumTimeTicket += end.Sub(start)
							sumTime += sumTimeTicket
							break
						}

						if len(dataHistory[i].History)-(j+1) == 0 {
							sumTimeTicket = 0
							break
						}
					}
				}
			}

			//--- Hours
			itemHistory.ActualManday.Float64 = sumTime.Hours()
			historyDBResult[idx] = itemHistory
		}
	}

	//--- Get List Report Actual ---> To Backlog
	actualDBResult, err = input.ReportDAO.GetListReport(db, inputStruct, actualSearchBy, departmentID, 0, false, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) unionAllResultReport(db *sql.DB, tempOutput []out.ReportResponse, inputStruct in.GetListDataDTO) (output interface{}, err errorModel.ErrorModel) {
	var (
		fileName    = "GetListReportService.go"
		funcName    = "unionAllResultReport"
		resultDB    string
		resultUnion out.ResultsReportResponse
		errorS      error
	)

	if len(tempOutput) > 0 {
		resultDB, err = dao.ReportDAO.UnionAllResultReport(db, tempOutput, inputStruct)
		if err.Error != nil {
			return
		}

		errorS = json.Unmarshal([]byte(resultDB), &resultUnion)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
			return
		}
	}

	output = resultUnion
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) convertModelToResponseGetListDevQA(backlogDBResult, actualDBResult, historyDBResult []interface{}) (output []out.ReportResponse, err errorModel.ErrorModel) {
	var (
		fileName = "GetListReportService.go"
		funcName = "convertModelToResponseGetList"
		result   []out.ReportResponse
	)

	for _, backlogDBResultItem := range backlogDBResult {
		var itemBacklog = backlogDBResultItem.(repository.ReportModel)
		temp := out.ReportResponse{
			NIK:        itemBacklog.NIK.Int64,
			Name:       itemBacklog.Name.String,
			Department: itemBacklog.Department.String,
		}

		//--- Input Backlog Manday
		for _, historyDBResultItem := range historyDBResult {
			var itemHistory = historyDBResultItem.(repository.ReportModel)
			if itemBacklog.NIK.Int64 == itemHistory.NIK.Int64 && itemBacklog.Tracker.String == itemHistory.Tracker.String {
				temp.BacklogManday = itemHistory.ActualManday.Float64 / constanta.DaysDefaultMandays
				temp.BacklogTicket = itemHistory.RedmineNumber.String
				break
			}
		}

		switch itemBacklog.DepartmentID.Int64 {
		case constanta.DeveloperDepartmentID: //-- Developer
			var trackerDev out.TrackerDeveloper
			_ = json.Unmarshal([]byte(itemBacklog.RateStr.String), &trackerDev)
			temp.Tracker = constanta.TrackerTask
			temp.MandayRate = trackerDev.Task
		case constanta.QAQCDepartmentID: //-- Tester
			if util.IsStringEmpty(itemBacklog.Tracker.String) {
				var trackerTester out.TrackerQA
				_ = json.Unmarshal([]byte(itemBacklog.RateStr.String), &trackerTester)
				temp.Tracker = constanta.TrackerManual
				temp.MandayRate = trackerTester.Manual
				result = append(result, out.ReportResponse{
					NIK:        itemBacklog.NIK.Int64,
					Name:       itemBacklog.Name.String,
					Department: itemBacklog.Department.String,
					Tracker:    constanta.TrackerAuto,
					MandayRate: trackerTester.Automation,
				})
			} else {
				var (
					trackerTester    out.TrackerQA
					duplicateTracker string
					duplicateRate    float64
				)

				_ = json.Unmarshal([]byte(itemBacklog.RateStr.String), &trackerTester)
				temp.Tracker = itemBacklog.Tracker.String

				if temp.Tracker == constanta.TrackerAuto {
					temp.MandayRate = trackerTester.Automation

					//--- Manual
					duplicateTracker = constanta.TrackerManual
					duplicateRate = trackerTester.Manual
				} else if temp.Tracker == constanta.TrackerManual {
					temp.MandayRate = trackerTester.Manual

					//--- Automation
					duplicateTracker = constanta.TrackerAuto
					duplicateRate = trackerTester.Automation
				} else {
					err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, fmt.Sprintf(`wrong tracker [%s]`, temp.Name), constanta.Tracker, "")
					return
				}

				//--- Check
				var isExistDuplicateTracker bool
				for _, itemCheck := range backlogDBResult {
					var itemBacklogCheck = itemCheck.(repository.ReportModel)
					if itemBacklogCheck.Name.String == temp.Name && itemBacklogCheck.NIK.Int64 == temp.NIK && itemBacklogCheck.Tracker.String == duplicateTracker {
						isExistDuplicateTracker = true
					}
				}

				//--- Duplicate
				if !isExistDuplicateTracker {
					result = append(result, out.ReportResponse{
						NIK:        itemBacklog.NIK.Int64,
						Name:       itemBacklog.Name.String,
						Department: itemBacklog.Department.String,
						Tracker:    duplicateTracker,
						MandayRate: duplicateRate,
					})
				}
			}
		default:
		}

		for _, actualDBResultItem := range actualDBResult {
			itemActual := actualDBResultItem.(repository.ReportModel)
			if itemBacklog.NIK.Int64 == itemActual.NIK.Int64 && itemBacklog.Tracker.String == itemActual.Tracker.String {

				temp.ActualManday = itemActual.ActualManday.Float64 / constanta.DaysDefaultMandays
				temp.ActualTicket = itemActual.RedmineNumber.String
				temp.Manday = temp.ActualManday * temp.MandayRate

				break
			}
		}

		result = append(result, temp)
	}

	output = result
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) convertModelToResponseGetListInfraDevOps(actualDBResult []interface{}) (output []out.ReportResponse, err errorModel.ErrorModel) {
	var result []out.ReportResponse
	for _, actualDBResultItem := range actualDBResult {
		var (
			itemActual         = actualDBResultItem.(repository.ReportModel)
			trackerInfraDevOps out.TrackerInfraDevOps
		)

		temp := out.ReportResponse{
			NIK:        itemActual.NIK.Int64,
			Name:       itemActual.Name.String,
			Department: itemActual.Department.String,
			Tracker:    constanta.TrackerTask,
		}

		_ = json.Unmarshal([]byte(itemActual.RateStr.String), &trackerInfraDevOps)
		temp.MandayRate = trackerInfraDevOps.Task
		temp.ActualManday = itemActual.ActualManday.Float64 / constanta.DaysDefaultMandays

		if itemActual.RedmineNumber.String != "" && itemActual.RedmineNumber.String != "[null]" {
			temp.ActualTicket = itemActual.RedmineNumber.String
		}

		temp.Manday = temp.ActualManday * temp.MandayRate
		result = append(result, temp)
	}

	output = result
	err = errorModel.GenerateNonErrorModel()
	return
}
