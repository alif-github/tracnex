package EmployeeService

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/dto/out"
	"regexp"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type employeeReimbursementReportColumn struct {
	ColumnName string
	Width      float64
}

var defaultEmployeeReimbursementReportHeader = []employeeReimbursementReportColumn{
	{
		ColumnName: "No",
	},
	{
		ColumnName: "NIK",
		Width:      15,
	},
	{
		ColumnName: "Nama",
		Width:      35,
	},
	{
		ColumnName: "TglJoin",
		Width:      20,
	},
	{
		ColumnName: "Entitlements",
		Width:      20,
	},
	{
		ColumnName: "Total",
		Width:      20,
	},
	{
		ColumnName: "Sisa",
		Width:      20,
	},
}

var defaultEmployeeReimbursementReportDetailHeader = []employeeReimbursementReportColumn{
	{
		ColumnName: "No",
	},
	{
		ColumnName: "Nama",
		Width:      30,
	},
	{
		ColumnName: "Kwitansi",
		Width:      35,
	},
	{
		ColumnName: "Nominal",
		Width:      30,
	},
	{
		ColumnName: "Keterangan",
		Width:      30,
	},
}

func (input employeeService) DownloadReimbursementReport(request *http.Request, contextModel *applicationModel.ContextModel) (output []byte, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validSearchBy = []string{"first_name", "last_name", "id_card"}
		validOrderBy  = []string{"e.id", "e.created_at"}
	)

	inputStruct, searchByParam, errModel = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, applicationModel.GetListEmployeeReimbursementValidOperator, service.DefaultLimit)
	if errModel.Error != nil {
		return
	}

	model, errModel := input.getEmployeeReimbursementFilters(request)
	if errModel.Error != nil {
		return
	}

	model.IsFilter.Bool = false
	now := time.Now()
	year , _, _ := now.Date()

	if model.Year.String == "" && model.Month.String == ""{
		model.Year.String = strconv.Itoa(year)
	}

	if model.ReportType.String != "detail" && model.ReportType.String != "summary"{
		errModel = errorModel.GenerateFieldFormatWithRuleError(input.FileName, "DownloadReimbursementReport", "Report type hanya boleh diisi dengan summary atau detail", "report_type", "")
		return
	}

	results, errModel := dao.EmployeeReimbursementDAO.GetListEmployeeReimbursementReport(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, model)
	if errModel.Error != nil {
		return
	}

	return input.writeReimbursementIntoExcel(results, model)
}

func (input employeeService) writeReimbursementIntoExcel(listData []interface{}, model repository.EmployeeReimbursement) (output []byte, header map[string]string, errModel errorModel.ErrorModel) {
	var (
		fileName = "EmployeeReimbursementReportService.go"
		funcName = "writeAnnualLeaveIntoExcel"
		excel    = excelize.NewFile()
		sheet    = "SUMMARY " + model.Year.String
		colStart = byte('A')
		rowStart = 1
	)

	monthArr := []string{
		"JAN", "FEB", "MAR", "APRIL", "MAY", "JUNE",
		"JULY", "AUG", "SEP", "OKT", "NOV", "DES",
	}

	if model.ReportType.String == "detail"{
		for i:=0; i<len(monthArr);i++  {
			sheetName := monthArr[i]+" "+model.Year.String
			_, errs := excel.NewSheet(sheetName)
			if errs != nil {
				errModel = errorModel.GenerateUnknownError(fileName, funcName, errs)
				return
			}
			datas, _ := dao.EmployeeReimbursementDAO.GetDetailReportReimbursement(serverconfig.ServerAttribute.DBConnection, "2023", int64(i+1), model)

			//if len(datas) >= 1 {
				/*
					Write Header
				*/
				if errModel = input.writeReimbursementDetailHeader(excel, sheetName, colStart, strconv.Itoa(1)); errModel.Error != nil {
					return
				}

				/*
					Add header style
				*/
				if errModel = input.styleReimbursementDetailColumns(excel, sheetName, strconv.Itoa(1)); errModel.Error != nil {
					return
				}

				if errModel = input.writeReimbursementDataDetail(excel, sheetName, colStart, 2, datas); errModel.Error != nil {
					return
				}
			//}

		}
	}

	if model.ReportType.String == "summary"{
		_, errs := excel.NewSheet(sheet)
		if errs != nil {
			errModel = errorModel.GenerateUnknownError(fileName, funcName, errs)
			return
		}

		excel.SetActiveSheet(0)

		/*
			Write Header
		*/
		if errModel = input.writeReimbursementHeader(excel, sheet, colStart, strconv.Itoa(rowStart), model); errModel.Error != nil {
			return
		}

		/*
			Add header style
		*/

		addColumn := 0
		if model.Month.String != "" && model.Year.String != ""{
			addColumn = 1
		}

		if model.Month.String == "" && model.Year.String != ""{
			addColumn = 12
		}

		if errModel = input.styleReimbursementColumns(excel, sheet, strconv.Itoa(rowStart), addColumn); errModel.Error != nil {
			return
		}

		/*
			Write Data
		*/
		rowStart++
		if errModel = input.writeReimbursementData(excel, sheet, colStart, rowStart, listData, model.Month.String); errModel.Error != nil {
			return
		}
	}

	errs := excel.DeleteSheet("Sheet1")
	if errs != nil {
		errModel = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	buf, err := excel.WriteToBuffer()
	if err != nil {
		errModel = errorModel.GenerateUnknownError(input.FileName, funcName, err)
		return
	}

	output = buf.Bytes()

	title := ""
	if model.ReportType.String == "detail"{
		title = "DETAIL_MEDICAL_KARYAWAN_"+model.Year.String
	}else if model.ReportType.String == "summary"{
		title = "SUMMARY_MEDICAL_KARYAWAN_"+model.Year.String
	}

	header = make(map[string]string)
	header["Content-Disposition"] = fmt.Sprintf("attachment; filename=%s", title+".xlsx")
	header["Content-Type"] = "application/vnd.ms-excel"
	header["Content-Length"] = strconv.Itoa(len(output))

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) writeReimbursementHeader(excel *excelize.File, sheet string, colStart byte, rowStart string, model repository.EmployeeReimbursement) (errModel errorModel.ErrorModel) {
	funcName := "writeReimbursementHeader"

	monthArr := []string{
		"Jan", "Feb", "Mar", "April", "May", "June",
		"July", "Aug", "Sept", "Okt", "Nov", "Des",
	}

	i, _ := strconv.Atoi(model.Month.String)

	var defaultTemplateTemp []employeeReimbursementReportColumn
	defaultTemplateTemp = defaultEmployeeReimbursementReportHeader

	if model.Month.String != "" && model.Year.String != ""{
		defaultTemplateTemp = append(defaultTemplateTemp, employeeReimbursementReportColumn{
			ColumnName: monthArr[i-1],
            Width: 12,
		})
	}

	if model.Month.String == "" && model.Year.String != ""{
		for i:=0; i<len(monthArr);i++  {
			defaultTemplateTemp = append(defaultTemplateTemp, employeeReimbursementReportColumn{
				ColumnName: monthArr[i],
				Width: 12,
			})
		}
	}

	for _, col := range defaultTemplateTemp {
		cell := string(colStart) + rowStart
		yearInt, _ := strconv.Atoi(model.Year.String)
		lastYear := yearInt-1

		if col.ColumnName == "Sisa"{
			col.ColumnName = "Sisa "+ strconv.Itoa(lastYear)
		}

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

	defaultTemplateTemp = defaultEmployeeReimbursementReportHeader

	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) writeReimbursementData(excel *excelize.File, sheet string, colStart byte, rowStart int, listData []interface{}, month string) (errModel errorModel.ErrorModel) {

	var grandTotalArr[12]float64
	ids, _ := strconv.Atoi(month)
	for i, data := range listData {
		var report repository.EmployeeReimbursement
		report, _ = data.(repository.EmployeeReimbursement)

		for idx:=0; idx<12;idx++  {
			grandTotalArr[idx]+=report.MonthlyReportArr[idx]
		}

		row := []interface{}{}
		if month != ""{
			row = []interface{}{
				i + 1,
				report.IDCard.String,
				report.Firstname.String + " "+ report.Lastname.String,
				report.DateJoin.Time.Format("2006-02-01"),
				formatCurrency(report.CurrentMedicalValue.Float64),
				formatCurrency(report.MonthlyReport.Total.Float64),
				formatCurrency(report.LastMedicalValue.Float64),
				formatCurrency(report.MonthlyReportArr[ids-1]),
			}
		}else{
			row = []interface{}{
				i + 1,
				report.IDCard.String,
				report.Firstname.String + " "+ report.Lastname.String,
				report.DateJoin.Time.Format("2006-02-01"),
				formatCurrency(report.CurrentMedicalValue.Float64),
				formatCurrency(report.MonthlyReport.Total.Float64),
				formatCurrency(report.LastMedicalValue.Float64),
				formatCurrency(report.MonthlyReportArr[0]),
				formatCurrency(report.MonthlyReportArr[1]),
				formatCurrency(report.MonthlyReportArr[2]),
				formatCurrency(report.MonthlyReportArr[3]),
				formatCurrency(report.MonthlyReportArr[4]),
				formatCurrency(report.MonthlyReportArr[5]),
				formatCurrency(report.MonthlyReportArr[6]),
				formatCurrency(report.MonthlyReportArr[7]),
				formatCurrency(report.MonthlyReportArr[8]),
				formatCurrency(report.MonthlyReportArr[9]),
				formatCurrency(report.MonthlyReportArr[10]),
				formatCurrency(report.MonthlyReportArr[11]),
			}
		}


		if errModel = input.writeReimbursementRow(excel, sheet, colStart, rowStart, row, "header", ""); errModel.Error != nil {
			return
		}

		rowStart++
	}

	rowAdd := []interface{}{}
	if month != ""{
		rowAdd = []interface{}{
			"",
			"",
			"",
			"",
			"",
			"",
			"Total",
			formatCurrency(grandTotalArr[ids-1]),
		}
	}else{
		rowAdd = []interface{}{
			"",
			"",
			"",
			"",
			"",
			"",
			"Total",
			formatCurrency(grandTotalArr[0]),
			formatCurrency(grandTotalArr[1]),
			formatCurrency(grandTotalArr[2]),
			formatCurrency(grandTotalArr[3]),
			formatCurrency(grandTotalArr[4]),
			formatCurrency(grandTotalArr[5]),
			formatCurrency(grandTotalArr[6]),
			formatCurrency(grandTotalArr[7]),
			formatCurrency(grandTotalArr[8]),
			formatCurrency(grandTotalArr[9]),
			formatCurrency(grandTotalArr[10]),
			formatCurrency(grandTotalArr[11]),
		}
	}


	if errModel = input.writeReimbursementRow(excel, sheet, colStart, rowStart, rowAdd, "grandtotal", ""); errModel.Error != nil {
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) writeReimbursementRow(excel *excelize.File, sheet string, col byte, rowStart int, row []interface{}, typ string, tre string) (errModel errorModel.ErrorModel) {
	funcName := "writeReimbursementRow"

	styleAlign := "center"
	styleColor := "#FFFFFF00"
	isBold := false
	size := 11
	for i, item := range row {
		cell := fmt.Sprintf("%s%d", string(col), rowStart)

		if err := excel.SetCellValue(sheet, cell, item); err != nil {
			return
		}

		if i == 2 && typ == "header" {
			styleAlign = "left"
			isBold = true
			size = 11
		}else{
			styleAlign = "center"
			isBold = false
			size = 11
		}

		if i >= 7 && item != "-" && typ == "header" {
			styleColor = "#FF8C00"
			isBold = false
			size = 11
		}else if i == 6 && item == "-" && typ == "header"{
			styleColor = "#DC143C"
			isBold = false
			size = 11
		}else if item == "-" && typ == "header"{
			styleColor = "#FFFFFF00"
			isBold = false
			size = 11
		}

		if i >= 7 && typ == "grandtotal" {
			styleColor = "#FF8C00"
			isBold = true
			size = 13
		}else if i <= 6 && item != "-" && typ == "grandtotal"{
			styleColor = "#FFFFFF00"
			isBold = true
			size = 13
		}


		if (i == 2 || i ==3) && tre == "True" {
			styleColor = "#8FBC8F"
			isBold = true
			size = 12
		}
		//else if typ != "header"{
		//	styleColor = "#FFFFFF00"
		//	isBold = false
		//	size = 11
		//}

		if (i == 2 || i ==3) && tre == "Finish" {
			styleColor = "#DAA520"
			isBold = true
			size = 15
		}else if(i != 2 || i !=3) && tre == "Finish"{
			styleColor = "#FFFFFF00"
			isBold = false
			size = 11
		}

		if i == 1 && typ == "detail" {
			styleAlign = "left"
			isBold = true
		}

		/*
			Add Styles
		*/
		styleId, err := excel.NewStyle(&excelize.Style{
			Alignment: &excelize.Alignment{
				Horizontal: styleAlign,
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Pattern: 1,
				Color:   []string{styleColor},
			},
			Font: &excelize.Font{
				Bold: isBold,
				Size: float64(size),
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
		if typ == "grandtotal"{
			excel.MergeCell(sheet, "A"+strconv.Itoa(rowStart), "G"+strconv.Itoa(rowStart))
		}
		_ = excel.SetCellStyle(sheet, hCell, hCell, styleId)

		col++
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) styleReimbursementColumns(excel *excelize.File, sheet, rowStart string, addColumn int) (errModel errorModel.ErrorModel) {
	var (
		funcName = "styleEmployeeLeaveColumns"
		colStart = byte('A')
		colEnd   = colStart + byte(len(defaultEmployeeReimbursementReportHeader)+addColumn-1)
	)

	styleId, err := excel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#A9A9A9"},
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
			Size: 13,
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

func (input employeeService) writeReimbursementDetailHeader(excel *excelize.File, sheet string, colStart byte, rowStart string) (errModel errorModel.ErrorModel) {
	funcName := "writeReimbursementDetailHeader"

	for _, col := range defaultEmployeeReimbursementReportDetailHeader {
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

func (input employeeService) styleReimbursementDetailColumns(excel *excelize.File, sheet, rowStart string) (errModel errorModel.ErrorModel) {
	var (
		funcName = "styleReimbursementDetailColumns"
		colStart = byte('A')
		colEnd   = colStart + byte(len(defaultEmployeeReimbursementReportDetailHeader)-1)
	)

	styleId, err := excel.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Pattern: 1,
			Color:   []string{"#A9A9A9"},
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
			Size: 13,
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

func formatCurrency(amount float64)string {
	f := fmt.Sprintf("%f", amount)
	s := strings.Split(f, ".")

	sa := s[0]
	sb := s[1]

	sb = sb[0:2]

	r := regexp.MustCompile("(\\d+)(\\d{3})")
	for n := ""; n != sa; {
		n = sa
		sa = r.ReplaceAllString(sa, "$1.$2")
	}

	result := ""
	if amount <= 0 {
		result = "-"
	}else{
		result = fmt.Sprintf("%s", sa)
	}

	return result
}

func (input employeeService) writeReimbursementDataDetail(excel *excelize.File, sheet string, colStart byte, rowStart int, listData []repository.EmployeeReimbursement) (errModel errorModel.ErrorModel) {
	var total float64
	for i, data := range listData {
		var detail []out.EmployeeReimbursementForReport
		json.Unmarshal([]byte(data.Description.String), &detail)
		cDetail := len(detail)+1
		name := ""
		noStr := ""
		isTotal := ""
		var totalReimEmp float64

		for a:=0; a<cDetail;a++  {
			if a == 0{
				name = data.Firstname.String+" "+data.Lastname.String
				noStr = strconv.Itoa(int(i+1))
			}else{
              name = ""
              noStr = ""
			}

			var row  []interface{}

			if a == (cDetail-1){
				isTotal = "True"
				row = []interface{}{
					"", "",
					"Total",
					formatCurrency(totalReimEmp),
					"",
				}
			}else {
                totalReimEmp+=detail[a].ApprovedValue
                isTotal = ""
				row = []interface{}{
					noStr,
					name,
					detail[a].ReceiptNo,
					formatCurrency(detail[a].ApprovedValue),
					detail[a].Description,
				}
			}

			if errModel = input.writeReimbursementRow(excel, sheet, colStart, rowStart, row, "detail", isTotal); errModel.Error != nil {
				return
			}
			rowStart++
		}
        total+=totalReimEmp
	}

	rows := []interface{}{
		"",
		"",
		"GrandTotal",
		formatCurrency(total),
		"",
	}

    if errModel = input.writeReimbursementRow(excel, sheet, colStart, rowStart, rows, "detail", "Finish"); errModel.Error != nil {
           return
     }

	return errorModel.GenerateNonErrorModel()
}