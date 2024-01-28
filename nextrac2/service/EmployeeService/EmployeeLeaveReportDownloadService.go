package EmployeeService

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/360EntSecGroup-Skylar/excelize"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
)

type employeeLeaveColumn struct {
	ColumnName string
	Width      float64
}

var defaultEmployeeLeaveCSVHeader = []employeeLeaveColumn{
	{
		ColumnName: "No",
	},
	{
		ColumnName: "NIK",
		Width:      15,
	},
	{
		ColumnName: "Nama",
		Width:      25,
	},
	{
		ColumnName: "Departemen",
		Width:      25,
	},
	{
		ColumnName: "Jenis Pengajuan",
	},
	{
		ColumnName: "Tanggal Diajukan",
		Width:      20,
	},
	{
		ColumnName: "Tanggal Absensi",
		Width:      30,
	},
	{
		ColumnName: "Total Absensi",
		Width:      20,
	},
}

func (input employeeService) DownloadEmployeeLeaveReport(request *http.Request, contextModel *applicationModel.ContextModel) (output []byte, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validSearchBy = []string{"first_name", "last_name", "el.type"}
		validOrderBy  = []string{"id", "created_at"}
	)

	inputStruct, searchByParam, errModel = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, applicationModel.GetListEmployeeLeaveValidOperator, service.DefaultLimit)
	if errModel.Error != nil {
		return
	}

	employee, errModel := dao.EmployeeDAO.GetByUserId(serverconfig.ServerAttribute.DBConnection, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	employeeLeaveFilter, errModel := input.getEmployeeLeaveFilter(request, employee)
	if errModel.Error != nil {
		return
	}

	results, errModel := dao.EmployeeLeaveDAO.GetListEmployeeLeave(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, employeeLeaveFilter)
	if errModel.Error != nil {
		return
	}

	return input.writeEmployeeLeaveIntoExcel(results)
}

func (input employeeService) writeEmployeeLeaveIntoExcel(employeeLeaveList []interface{}) (output []byte, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		fileName = "EmployeeLeaveReportDownloadService.go"
		funcName = "writeEmployeeLeaveIntoExcel"
		excel    = excelize.NewFile()
		sheet    = constanta.EmployeeLeaveReportSheetName
		colStart = byte('A')
		rowStart = 1
	)

	sheetId, errs := excel.NewSheet(sheet)
	if errs != nil {
		errModel = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	errs = excel.DeleteSheet("Sheet1")
	if errs != nil {
		errModel = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	excel.SetActiveSheet(sheetId)

	/*
		Write Header
	*/
	if errModel = input.writeEmployeeLeaveHeader(excel, sheet, colStart, strconv.Itoa(rowStart)); errModel.Error != nil {
		return
	}

	/*
		Add header style
	*/
	if errModel = input.styleEmployeeLeaveColumns(excel, sheet, strconv.Itoa(rowStart)); errModel.Error != nil {
		return
	}

	/*
		Write Data
	*/
	rowStart++
	if errModel = input.writeEmployeeLeaveData(excel, sheet, colStart, rowStart, employeeLeaveList); errModel.Error != nil {
		return
	}

	buf, err := excel.WriteToBuffer()
	if err != nil {
		errModel = errorModel.GenerateUnknownError(input.FileName, funcName, err)
		return
	}

	output = buf.Bytes()

	header = make(map[string]string)
	header["Content-Disposition"] = fmt.Sprintf("attachment; filename=%s", constanta.EmployeeLeaveReportFileName)
	header["Content-Type"] = "application/vnd.ms-excel"
	header["Content-Length"] = strconv.Itoa(len(output))

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) writeEmployeeLeaveHeader(excel *excelize.File, sheet string, colStart byte, rowStart string) (errModel errorModel.ErrorModel) {
	funcName := "writeEmployeeLeaveHeader"

	for _, col := range defaultEmployeeLeaveCSVHeader {
		cell := string(colStart) + rowStart

		if err := excel.SetCellValue(sheet, cell, col.ColumnName); err != nil {
			errModel = errorModel.GenerateUnknownError(input.FileName, funcName, err)
			return
		}

		/*
			Set Col Width
		*/
		input.setColumnWidth(excel, sheet, string(colStart), col.ColumnName, col.Width)
		colStart++
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) writeEmployeeLeaveData(excel *excelize.File, sheet string, colStart byte, rowStart int, employeeLeaveList []interface{}) (errModel errorModel.ErrorModel) {
	for i, data := range employeeLeaveList {
		var (
			dateList      []string
			employeeLeave repository.EmployeeLeaveModel
		)

		employeeLeave, _ = data.(repository.EmployeeLeaveModel)

		_ = json.Unmarshal([]byte(employeeLeave.StrDateList.String), &dateList)

		absenceDate, absenceTotal := input.getAbsence(employeeLeave.Type.String, dateList)

		row := []interface{}{
			i + 1,
			employeeLeave.IDCard.String,
			fmt.Sprintf("%s %s", employeeLeave.Firstname.String, employeeLeave.Lastname.String),
			employeeLeave.Department.String,
			input.getLeaveAlias(employeeLeave.Type.String),
			employeeLeave.CreatedAt.Time.Format("02 Jan 2006"),
			absenceDate,
			absenceTotal,
		}

		if errModel = input.writeEmployeeLeaveRow(excel, sheet, colStart, rowStart, row); errModel.Error != nil {
			return
		}

		rowStart++
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) writeEmployeeLeaveRow(excel *excelize.File, sheet string, col byte, rowStart int, row []interface{}) (errModel errorModel.ErrorModel) {
	funcName := "writeEmployeeLeaveRow"

	for _, item := range row {
		cell := fmt.Sprintf("%s%d", string(col), rowStart)

		if err := excel.SetCellValue(sheet, cell, item); err != nil {
			return
		}

		/*
			Add Styles
		*/
		styleId, err := excel.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: "center",
			},
			Border: []excelize.Border{
				{
					Type:  "top",
					Color: "#000000",
					Style: 1,
				},
				{
					Type:  "bottom",
					Color: "#000000",
					Style: 1,
				},
				{
					Type:  "left",
					Color: "#000000",
					Style: 1,
				},
				{
					Type:  "right",
					Color: "#000000",
					Style: 1,
				},
			},
		})
		if err != nil {
			errModel = errorModel.GenerateUnknownError(input.FileName, funcName, err)
			return
		}

		hCell := fmt.Sprintf("%s%d", string(col), rowStart)
		_ = excel.SetCellStyle(sheet, hCell, hCell, styleId)

		col++
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) getLeaveAlias(leaveType string) string {
	switch leaveType {
	case constanta.LeaveType:
		return constanta.LeaveTypeAlias
	case constanta.PermitType:
		return constanta.PermitTypeAlias
	case constanta.SickLeaveType:
		return constanta.SickLeaveTypeAlias
	}

	return ""
}

func (input employeeService) getAbsence(leaveType string, dateList []string) (date string, absenceTotal string) {
	/*
		Permit
	*/
	if leaveType == constanta.PermitType {
		startDate, _ := time.Parse(constanta.DefaultTimeFormat, dateList[0])
		endDate, _ := time.Parse(constanta.DefaultTimeFormat, dateList[1])

		diff := endDate.Sub(startDate)

		date = startDate.Format("02 Jan 2006")

		absenceTotal = strconv.FormatFloat(math.Round(diff.Hours() * 10000)/10000, 'f', -1, 64)
		absenceTotal = fmt.Sprintf("%s Jam", absenceTotal)
		return
	}

	/*
		Leave Or Sick Leave
	*/
	strStartDate := ""
	strEndDate := ""

	startDate, _ := time.Parse(constanta.DefaultTimeFormat, dateList[0])
	strStartDate = startDate.Format("02 Jan 2006")

	if len(dateList) > 1 {
		endDate, _ := time.Parse(constanta.DefaultTimeFormat, dateList[len(dateList)-1])
		strEndDate = fmt.Sprintf(" - %s", endDate.Format("02 Jan 2006"))
	}

	return strStartDate + strEndDate, fmt.Sprintf("%d Hari", len(dateList))
}

func (input employeeService) setColumnWidth(excel *excelize.File, sheet, colStart, value string, width float64) {
	result := width

	if width == 0 {
		result = float64(utf8.RuneCountInString(value)) + 2
	}

	result += 0.78

	_ = excel.SetColWidth(sheet, colStart, colStart, result)
}

func (input employeeService) styleEmployeeLeaveColumns(excel *excelize.File, sheet, rowStart string) (errModel errorModel.ErrorModel) {
	var (
		funcName = "styleEmployeeLeaveColumns"
		colStart = byte('A')
		colEnd   = colStart + byte(len(defaultEmployeeLeaveCSVHeader)-1)
	)

	styleId, err := excel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#a0f78b"},
		},
		Border: []excelize.Border{
			{
				Type:  "top",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "left",
				Color: "#000000",
				Style: 1,
			},
			{
				Type:  "right",
				Color: "#000000",
				Style: 1,
			},
		},
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		errModel = errorModel.GenerateUnknownError(input.FileName, funcName, err)
		return
	}

	hCell := fmt.Sprintf("%s%s", string(colStart), rowStart)
	vCell := fmt.Sprintf("%s%s", string(colEnd), rowStart)

	_ = excel.SetCellStyle(sheet, hCell, vCell, styleId)
	return
}

func (input employeeService) generateEmployeeLeaveCSVResult(data []interface{}) (output [][]string) {
	//output = append(output, defaultEmployeeLeaveCSVHeader)
	//
	//for _, item := range data {
	//	employeeLeave, _ := item.(repository.EmployeeLeaveModel)
	//}

	return
}
