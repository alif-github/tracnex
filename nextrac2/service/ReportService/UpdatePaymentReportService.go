package ReportService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/redmine_request"
	"nexsoft.co.id/nextrac2/serverconfig"
	"sort"
	"strings"
	"time"
)

func (input reportService) UpdatePaymentReport(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName          = "UpdatePaymentReport"
		dbRedmineDev      = serverconfig.ServerAttribute.RedmineDBConnection
		inputStruct       in.GetListDataDTO
		searchByParam     []in.SearchByParam
		isMandatoryExist  bool
		actualResultDB    []interface{}
		outputGetList     interface{}
		successTicket     []int64
		errorTicket       []in.ErrorBundleReport
		departmentID      int64
		departmentIDParam int64
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListReportValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	isMandatoryExist, departmentIDParam, err = input.validateAddParam(request, &inputStruct, &searchByParam)
	if err.Error != nil {
		return
	}

	//--- If Not Developer Or QA-QC Returning
	if departmentIDParam != constanta.DeveloperDepartmentID && departmentIDParam != constanta.QAQCDepartmentID {
		output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
		err = errorModel.GenerateNonErrorModel()
		return
	}

	if isMandatoryExist {
		//--- Get Data (Main Process)
		outputGetList, departmentID, actualResultDB, err = input.doGetListReport(inputStruct, searchByParam, contextModel)
		if err.Error != nil {
			return
		}
	}

	//--- Update Sprint And Payment On Redmine
	_, err = input.ServicePreparedDBCustomize(funcName, dbRedmineDev, nil, contextModel, input.doUpdateSprintAndPaymentStatusOnRedmine)
	if err.Error != nil {
		return
	}

	if actualResultDB != nil {
		//--- Paid Ticket Redmine
		successTicket, errorTicket, err = input.paidTicketOnRedmine(actualResultDB, *contextModel)
		if err.Error != nil {
			return
		}
	}

	if len(successTicket) > 0 {
		err = input.updateStatusPaymentOnBacklog(outputGetList, departmentID, successTicket, contextModel)
		if err.Error != nil {
			return
		}
	}

	if len(errorTicket) > 0 {
		output.Other = errorTicket
	}

	content := make(map[string][]int64)
	content["success_ticket"] = successTicket

	output.Data.Content = content
	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) doUpdateStatusPaymentOnBacklog(tx *sql.Tx, dataInput interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		dataReportHistory = dataInput.(in.ReportHistory)
		resultID          int64
	)

	if err = dao.BacklogDAO.UpdatePaymentOnBacklog(tx, dataReportHistory.SuccessTicket, repository.BacklogModel{
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}); err.Error != nil {
		return
	}

	successTicketStrByte, _ := json.Marshal(dataReportHistory.SuccessTicket)
	successTicketStr := string(successTicketStrByte)
	if resultID, err = dao.ReportDAO.InsertToReportHistory(tx, repository.ReportHistoryModel{
		Data:          sql.NullString{String: dataReportHistory.DataActual},
		SuccessTicket: sql.NullString{String: successTicketStr},
		DepartmentID:  sql.NullInt64{Int64: dataReportHistory.DepartmentID},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
	}); err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: "report_history"},
		PrimaryKey: sql.NullInt64{Int64: resultID},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) doUpdateSprintAndPaymentStatusOnRedmine(tx *sql.Tx, _ interface{}, _ *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, err errorModel.ErrorModel) {
	err = input.doUpdatePaymentOnRedmine(tx, timeNow)
	if err.Error != nil {
		return
	}

	err = input.doUpdateSprintOnRedmine(tx, timeNow)
	if err.Error != nil {
		return
	}

	//--- Update For Testing
	//raw := "---\n- 20240106-20240120\n- 20231221-20240105\n- 20231206-20231220\n- 20231121-20231205\n"
	//err = dao.RedmineDAO.UpdateCustomFieldsOnRedmine(tx, 9, raw)
	//if err.Error != nil {
	//	return
	//}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) updateStatusPaymentOnBacklog(outputGetList interface{}, departmentID int64, successTicket []int64, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	var (
		funcName                = "updateStatusPaymentOnBacklog"
		tempOutputReportGetList = outputGetList.(out.ResultsReportResponse)
		dataGetList             string
		byteStrOutput           []byte
		reportHistory           in.ReportHistory
	)

	byteStrOutput, _ = json.Marshal(tempOutputReportGetList)
	dataGetList = string(byteStrOutput)
	reportHistory = in.ReportHistory{
		SuccessTicket: successTicket,
		DataActual:    dataGetList,
		DepartmentID:  departmentID,
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, reportHistory, contextModel, input.doUpdateStatusPaymentOnBacklog, func(i interface{}, model applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) doUpdatePaymentOnRedmine(tx *sql.Tx, timeNow time.Time) (err errorModel.ErrorModel) {
	var (
		fileName          = "UpdatePaymentReportService.go"
		funcName          = "doUpdatePaymentOnRedmine"
		dbRedmine         = serverconfig.ServerAttribute.RedmineDBConnection
		idCustomField     = constanta.IDPaymentOnRedmineCustomFields
		timeNowUpdated    = timeNow
		layoutPaid        = "2006-01"
		resultDB          string
		year, month, paid string
		isPaidUpdated     bool
		dataDate          []time.Time
	)

	//--- Update Preparation
	if timeNow.Month() == time.December {
		year = fmt.Sprintf(`%04d`, timeNowUpdated.Year()+1)
		month = fmt.Sprintf(`%02d`, time.January)
	} else {
		year = fmt.Sprintf(`%04d`, timeNowUpdated.Year())
		month = fmt.Sprintf(`%02d`, timeNowUpdated.Month()+1)
	}

	//--- Add This To List
	paid = fmt.Sprintf(`- PAID-%s-%s`, year, month)

	//--- Get Data From Paid
	resultDB, err = dao.RedmineDAO.GetCustomFieldsFromRedmine(dbRedmine, idCustomField)
	if err.Error != nil {
		return
	}

	//--- Check Data Updated
	splitData := strings.Split(resultDB, "\n")
	for _, itemSplitData := range splitData {
		//--- Continue
		if itemSplitData == "---" || itemSplitData == "" || itemSplitData == "- UNPAID" || itemSplitData == "- REJECTED" {
			continue
		}

		//--- Trim Paid
		date := strings.TrimLeft(itemSplitData, "- PAID-")
		timeParse, errorS := time.Parse(layoutPaid, date)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
			return
		}

		//--- Date Append
		dataDate = append(dataDate, timeParse)
	}

	//--- Sort data by Date in ascending order
	sort.Slice(dataDate, func(i, j int) bool {
		return dataDate[i].Before(dataDate[j])
	})

	//--- Latest Paid
	latestPaid := dataDate[len(dataDate)-1]

	//--- TimeNow Paid
	timeNowPaid := timeNow.Format(layoutPaid)
	timeNowPaidTime, errorS := time.Parse(layoutPaid, timeNowPaid)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	if latestPaid.Before(timeNowPaidTime) {
		for {
			//--- Add 1 Month Later
			latestPaid = latestPaid.AddDate(0, 1, 0)
			yearLatestPaid := fmt.Sprintf(`%04d`, latestPaid.Year())
			monthLatestPaid := fmt.Sprintf(`%02d`, latestPaid.Month())

			//--- Append To Status
			resultDB += fmt.Sprintf(`- PAID-%s-%s%s`, yearLatestPaid, monthLatestPaid, "\n")

			//--- If Equal Break
			if latestPaid.Equal(timeNowPaidTime) {
				break
			}
		}
	}

	//--- Check Data Is Updated Or Not
	splitData = strings.Split(resultDB, "\n")
	for _, itemSplitData := range splitData {
		if itemSplitData == paid {
			isPaidUpdated = true
		}
	}

	//--- Update To DB
	if !isPaidUpdated {
		resultDB += paid + "\n"
		fmt.Println("Starting Update To Redmine -> PAID PAYMENT")
		fmt.Println(fmt.Sprintf(`List PAID PAYMENT -> %s`, resultDB))
		return dao.RedmineDAO.UpdateCustomFieldsOnRedmine(tx, idCustomField, resultDB)
	}

	//--- Status Payment Has Updated
	fmt.Println("Redmine PAID PAYMENT Up To Date")
	fmt.Println(fmt.Sprintf(`List PAID PAYMENT Up To Date -> %s`, resultDB))
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) doUpdateSprintOnRedmine(tx *sql.Tx, timeNow time.Time) (err errorModel.ErrorModel) {
	var (
		fileName      = "UpdatePaymentReportService.go"
		funcName      = "doUpdateSprintOnRedmine"
		dbRedmine     = serverconfig.ServerAttribute.RedmineDBConnection
		idCustomField = constanta.IDSprintOnRedmineCustomFields
		resultDB      string
		dateRanges    []DateRange
		dateNew       []DateRange
		latestEndDate time.Time
	)

	//--- Get Data From Paid
	resultDB, err = dao.RedmineDAO.GetCustomFieldsFromRedmine(dbRedmine, idCustomField)
	if err.Error != nil {
		return
	}

	//--- Check Data Updated
	splitData := strings.Split(resultDB, "\n")
	for _, itemSplitData := range splitData {
		var (
			timeStart time.Time
			timeEnd   time.Time
			errorS    error
		)

		if itemSplitData == "---" || itemSplitData == "" {
			continue
		}

		d := strings.TrimLeft(itemSplitData, "- ")
		dateStr := strings.Split(d, "-")

		//--- Start Time
		timeStart, errorS = time.Parse(constanta.DefaultTimeSprintFormat, dateStr[0])
		if errorS != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
			return
		}

		//--- Start End
		timeEnd, errorS = time.Parse(constanta.DefaultTimeSprintFormat, dateStr[1])
		if errorS != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
			return
		}

		dateRanges = append(dateRanges, DateRange{
			StartDate: timeStart,
			EndDate:   timeEnd,
		})
	}

	//--- Sort data by EndDate in ascending order
	sort.Slice(dateRanges, func(i, j int) bool {
		return dateRanges[i].EndDate.Before(dateRanges[j].EndDate)
	})

	//--- Get the latest EndDate from existing data
	var latestStartDate time.Time
	if len(dateRanges) > 0 {
		latestEndDate = dateRanges[len(dateRanges)-1].EndDate
		latestStartDate = dateRanges[len(dateRanges)-1].StartDate
	}

	//--- Returning
	month2Later := timeNow.AddDate(0, 2, 0)
	fmt.Println("LastEndDateYear : ", latestEndDate.Year())
	fmt.Println("Month2LaterYear : ", month2Later.Year())
	fmt.Println("LastEndDateMonth : ", latestEndDate.Month())
	fmt.Println("Month2LaterMonth : ", month2Later.Month())
	if latestEndDate.Year() <= month2Later.Year() {
		if (latestEndDate.Year() == month2Later.Year()) && (latestEndDate.Month() >= month2Later.Month()) {
			fmt.Println(fmt.Sprintf(`Sprint Has Up To Date ---> %s-%s`, latestStartDate.Format(constanta.DefaultTimeSprintFormat), latestEndDate.Format(constanta.DefaultTimeSprintFormat)))
			return
		}
	} else {
		fmt.Println(fmt.Sprintf(`Sprint Has Up To Date ---> %s-%s`, latestStartDate.Format(constanta.DefaultTimeSprintFormat), latestEndDate.Format(constanta.DefaultTimeSprintFormat)))
		return
	}

	//--- TimeNow Paid
	//timeNowPaid := timeNow.Format(constanta.DefaultTimeSprintFormat)
	//timeNowPaidTime, errorS := time.Parse(constanta.DefaultTimeSprintFormat, timeNowPaid)
	//if errorS != nil {
	//	err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
	//	return
	//}

	for {
		var dateRange DateRange
		if (latestEndDate.Year() >= month2Later.Year() && latestEndDate.Month() >= month2Later.Month()) && latestEndDate.Day() == 20 {
			break
		}

		dateRange.StartDate = latestEndDate.AddDate(0, 0, 1)
		dateRange.EndDate = latestEndDate.AddDate(0, 0, 14)

		if fmt.Sprintf(`%02d`, latestEndDate.Day()) == "05" {
			dateRange.EndDate = time.Date(dateRange.EndDate.Year(), dateRange.EndDate.Month(), 20, 0, 0, 0, 0, time.UTC) //-- If 06 Then Create Day 20
		} else {
			dateRange.EndDate = time.Date(dateRange.EndDate.Year(), dateRange.EndDate.Month(), 05, 0, 0, 0, 0, time.UTC) //-- If 21 Then Create Day 05
		}

		latestEndDate = dateRange.EndDate
		dateNew = append(dateNew, dateRange)
	}

	//for i := 0; i < 4; i++ {
	//	var dateRange DateRange
	//	dateRange.StartDate = latestEndDate.AddDate(0, 0, 1)
	//	dateRange.EndDate = latestEndDate.AddDate(0, 0, 14)
	//
	//	if fmt.Sprintf(`%02d`, latestEndDate.Day()) == "05" {
	//		dateRange.EndDate = time.Date(dateRange.EndDate.Year(), dateRange.EndDate.Month(), 20, 0, 0, 0, 0, time.UTC) //-- If 06 Then Create Day 20
	//	} else {
	//		dateRange.EndDate = time.Date(dateRange.EndDate.Year(), dateRange.EndDate.Month(), 05, 0, 0, 0, 0, time.UTC) //-- If 21 Then Create Day 05
	//	}
	//
	//	latestEndDate = dateRange.EndDate
	//	dateNew = append(dateNew, dateRange)
	//}

	//--- Sort data last by EndDate in ascending order
	dateRanges = append(dateRanges, dateNew...)
	sort.Slice(dateRanges, func(i, j int) bool {
		return dateRanges[i].EndDate.Before(dateRanges[j].EndDate)
	})

	raw := "---\n"
	for i := len(dateRanges) - 1; i >= 0; i-- {
		raw += fmt.Sprintf(`- %s-%s`, dateRanges[i].StartDate.Format(constanta.DefaultTimeSprintFormat), dateRanges[i].EndDate.Format(constanta.DefaultTimeSprintFormat))
		raw += "\n"
	}

	//--- Update
	fmt.Println("Data Update Sprint On Redmine --->")
	fmt.Println(fmt.Sprintf(`Data Sprint -> %s`, raw))
	err = dao.RedmineDAO.UpdateCustomFieldsOnRedmine(tx, idCustomField, raw)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportService) paidTicketOnRedmine(actualResultDB []interface{}, contextModel applicationModel.ContextModel) (successTicket []int64, errorBundle []in.ErrorBundleReport, err errorModel.ErrorModel) {
	var (
		fileName = "UpdatePaymentReportService.go"
		funcName = "paidTicketOnRedmine"
		key      = config.ApplicationConfiguration.GetRedmineDBKeyAccess()
		api      = config.ApplicationConfiguration.GetRedmineAPI()
		timeNow  = time.Now()
		year     = fmt.Sprintf(`%04d`, timeNow.Year())
		month    = fmt.Sprintf(`%02d`, timeNow.Month())
		paid     = fmt.Sprintf(`PAID-%s-%s`, year, month)
	)

	fmt.Println("Starting To Paid On Redmine ---> ", paid)
	for _, itemActualResultDB := range actualResultDB {
		var (
			redmineNumber []int64
			data          repository.ReportModel
			errS          error
			bodyRequest   redmine_request.IssuePaidRedmineRequest
		)

		data = itemActualResultDB.(repository.ReportModel)
		errS = json.Unmarshal([]byte(data.RedmineNumber.String), &redmineNumber)
		if errS != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errS)
			return
		}

		bodyRequest = redmine_request.IssuePaidRedmineRequest{
			Issue: redmine_request.CustomFields{
				CustomField: []redmine_request.Fields{{ID: constanta.IDPaymentOnRedmineCustomFields, Value: paid}},
			},
		}

		for _, itemRedmineNumber := range redmineNumber {
			var (
				code   int
				domain string
				url    string
				errorS error
			)

			if itemRedmineNumber < 1 {
				continue
			}

			//-- Hit to redmine
			domain = api.Host + api.PathRedirect.UpdatePaid
			url = fmt.Sprintf(`%s/%d.json`, domain, itemRedmineNumber)
			code, _, errorS = common.HitPaidPaymentRedmineServer(key, url, bodyRequest, &contextModel)
			if errorS != nil {
				errorBundle = append(errorBundle, in.ErrorBundleReport{
					RedmineNumber: itemRedmineNumber,
					Message:       errorS.Error(),
				})
				continue
			}

			//--- If Failed Parsing To Error Category
			fmt.Println(fmt.Sprintf(`Update To Redmine Status ---> %d [#%d] %s`, code, itemRedmineNumber, url))
			fmt.Println(fmt.Sprintf(`Update To Redmine Status [Body] ---> %s`, util.StructToJSON(bodyRequest)))
			if code != http.StatusNoContent {
				if code == http.StatusNotFound {
					errorBundle = append(errorBundle, in.ErrorBundleReport{
						RedmineNumber: itemRedmineNumber,
						Message:       "Ticket Not Found",
					})
				} else {
					errorBundle = append(errorBundle, in.ErrorBundleReport{
						RedmineNumber: itemRedmineNumber,
						Message:       "Failed Unknown",
					})
				}
				continue
			}

			//-- Add Success Ticket
			successTicket = append(successTicket, itemRedmineNumber)
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

type DateRange struct {
	StartDate time.Time
	EndDate   time.Time
}
