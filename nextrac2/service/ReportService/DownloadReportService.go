package ReportService

import (
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"strconv"
	"strings"
	"time"
)

func (input reportService) DownloadListReport(request *http.Request, contextModel *applicationModel.ContextModel) (output [][]string, header map[string]string, err errorModel.ErrorModel) {
	var (
		fileName         = "DownloadReportService.go"
		funcName         = "DownloadListReport"
		inputStruct      in.GetListDataDTO
		searchByParam    []in.SearchByParam
		data             interface{}
		departmentID     int64
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
		data, departmentID, _, err = input.doGetListReport(inputStruct, searchByParam, contextModel)
		if err.Error != nil {
			return
		}
	}

	if data != nil {
		var (
			timeNow     = time.Now()
			timeNowStr  = timeNow.Format(fmt.Sprintf(`200601021504`))
			title       = "report-"
			//currency    = "Rp."
			dataOnModel = data.(out.ResultsReportResponse).Results
			headerCSV   []string
		)

		//-- Prepare Default Header
		headerDefault := []string{
			"NIK",
			"Nama",
			"Department",
			"Backlog Mandays",
			"Actual Mandays",
			"Manday Rate",
			"Mandays",
		}

		//-- Set header
		header = make(map[string]string)

		//-- Set title and attachment
		switch departmentID {
		case constanta.DeveloperDepartmentID:
			title += fmt.Sprintf(`developer-%s.csv`, timeNowStr)
			headerCSV = headerDefault
		case constanta.QAQCDepartmentID:
			title += fmt.Sprintf(`tester-%s.csv`, timeNowStr)
			headerCSV = []string{
				"NIK",
				"Nama",
				"Department",
				"Backlog Mandays Automation",
				"Actual Mandays Automation",
				"Manday Rate Automation",
				"Mandays Automation",
				"Backlog Mandays Manual",
				"Actual Mandays Manual",
				"Manday Rate Manual",
				"Mandays Manual",
				"Total Mandays",
			}
		case constanta.InfraDepartmentID:
			title += fmt.Sprintf(`infra-%s.csv`, timeNowStr)
			headerCSV = []string{
				"NIK",
				"Nama",
				"Department",
				"Backlog Mandays",
				"Manday Rate",
				"Mandays",
			}
		case constanta.DevOpsDepartmentID:
			title += fmt.Sprintf(`devops-%s.csv`, timeNowStr)
			headerCSV = []string{
				"NIK",
				"Nama",
				"Department",
				"Backlog Mandays",
				"Manday Rate",
				"Mandays",
			}
		case constanta.UIUXDepartmentID:
			title += fmt.Sprintf(`uiux-%s.csv`, timeNowStr)
			headerCSV = headerDefault
		default:
		}

		header["Content-Disposition"] = fmt.Sprintf(`attachment;filename=%s`, title)
		output = append(output, headerCSV)
		for _, itemDataOnModel := range dataOnModel {
			var row []string
			row = append(row,
				strconv.Itoa(int(itemDataOnModel.NIK)),
				itemDataOnModel.Name,
				itemDataOnModel.Department,
			)

			switch departmentID {
			case constanta.DeveloperDepartmentID, constanta.UIUXDepartmentID:
				row = append(row,

					//-- Task
					strings.ReplaceAll(fmt.Sprintf(`%.4f`, itemDataOnModel.Detail[0].BacklogManday), ".", ","),
					strings.ReplaceAll(fmt.Sprintf(`%.4f`, itemDataOnModel.Detail[0].ActualManday), ".", ","),
					strings.ReplaceAll(fmt.Sprintf(`%.2f`, itemDataOnModel.Detail[0].MandayRate), ".", ","),
					strings.ReplaceAll(fmt.Sprintf(`%.2f`, itemDataOnModel.Detail[0].Manday), ".", ","),
				)
			case constanta.InfraDepartmentID, constanta.DevOpsDepartmentID:
				row = append(row,

					//-- Task
					strings.ReplaceAll(fmt.Sprintf(`%.4f`, itemDataOnModel.Detail[0].BacklogManday), ".", ","),
					strings.ReplaceAll(fmt.Sprintf(`%.2f`, itemDataOnModel.Detail[0].MandayRate), ".", ","),
					strings.ReplaceAll(fmt.Sprintf(`%.2f`, itemDataOnModel.Detail[0].Manday), ".", ","),
				)
			case constanta.QAQCDepartmentID:
				var automationIndex, manualIndex int
				for idx, itemDataOnModelDetail := range itemDataOnModel.Detail {

					//-- Automation
					if itemDataOnModelDetail.Tracker == constanta.TrackerAuto {
						automationIndex = idx
					}

					//-- Manual
					if itemDataOnModelDetail.Tracker == constanta.TrackerManual {
						manualIndex = idx
					}
				}

				row = append(row,

					//-- Automation
					strings.ReplaceAll(fmt.Sprintf(`%.4f`, itemDataOnModel.Detail[automationIndex].BacklogManday), ".", ","),
					strings.ReplaceAll(fmt.Sprintf(`%.4f`, itemDataOnModel.Detail[automationIndex].ActualManday), ".", ","),
					strings.ReplaceAll(fmt.Sprintf(`%.2f`, itemDataOnModel.Detail[automationIndex].MandayRate), ".", ","),
					strings.ReplaceAll(fmt.Sprintf(`%.2f`, itemDataOnModel.Detail[automationIndex].Manday), ".", ","),

					//-- Manual
					strings.ReplaceAll(fmt.Sprintf(`%.4f`, itemDataOnModel.Detail[manualIndex].BacklogManday), ".", ","),
					strings.ReplaceAll(fmt.Sprintf(`%.4f`, itemDataOnModel.Detail[manualIndex].ActualManday), ".", ","),
					strings.ReplaceAll(fmt.Sprintf(`%.2f`, itemDataOnModel.Detail[manualIndex].MandayRate), ".", ","),
					strings.ReplaceAll(fmt.Sprintf(`%.2f`, itemDataOnModel.Detail[manualIndex].Manday), ".", ","),

					//-- Total
					strings.ReplaceAll(fmt.Sprintf(`%.2f`, itemDataOnModel.TotalManday), ".", ","),
				)
			default:
			}
			output = append(output, row)
		}

		return
	}

	err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ReportConstanta)
	return
}
