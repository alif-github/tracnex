package ReportAnnualLeaveService

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io/ioutil"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"os"
	"time"
)

func (input reportAnnualLeaveService) excelProduce(year1, year2 int, datas1, datas2 []repository.EmployeeLeaveReportModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (fileUploadID int64, err errorModel.ErrorModel) {
	var (
		fileName       = "ExcelReportAnnualLeaveService.go"
		funcName       = "excelProduce"
		file           *excelize.File
		headerMrgY     [][]interface{}
		headerMrgX     [][]interface{}
		header         [][]interface{}
		leaveYear1     = fmt.Sprintf(`Cuti %d`, year1)
		leaveYear2     = fmt.Sprintf(`Cuti %d`, year2)
		yearStr1       = fmt.Sprintf(`%d`, year1)
		yearStr2       = fmt.Sprintf(`%d`, year2)
		sheetName      = "Sheet1"
		db             = serverconfig.ServerAttribute.DBConnection
		fileInfo       os.FileInfo
		fileContent    []byte
		errs           error
		tx             *sql.Tx
		baseDirectory  = config.ApplicationConfiguration.GetDataDirectory().BaseDirectoryPath
		baseImportPath = config.ApplicationConfiguration.GetDataDirectory().ImportPath
		baseRootPath   = config.ApplicationConfiguration.GetDataDirectory().ReportLeavePath
		fileNamePath   = "report-leave.xlsx"
		directory      = baseDirectory + baseImportPath + baseRootPath + "/" + fileNamePath
	)

	tx, errs = db.Begin()
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	defer func() {
		if errs != nil || err.Error != nil {
			errs = tx.Rollback()
			if errs != nil {
				err = errorModel.GenerateUnknownError(fileName, funcName, errs)
			}
		} else {
			errs = tx.Commit()
			if errs != nil {
				err = errorModel.GenerateUnknownError(fileName, funcName, errs)
			}
		}
	}()

	headerMrgY = [][]interface{}{
		{"A1", "A2", "No"},
		{"B1", "B2", "NIK"},
		{"C1", "C2", "Nama"},
		{"D1", "D2", "Jabatan"},
		{"E1", "E2", "Divisi"},
		{"F1", "F2", "Mulai Kerja"},
		{"G1", "G2", "Selesai Probation"},
	}

	headerMrgX = [][]interface{}{
		{
			"I1", "J1", "K1",
			"L1", "M1", "N1",
			"O1", "P1", "Q1",
			"R1", "S1", "T1",
			leaveYear1,
		},
		{
			"AB1", "AC1", "AD1",
			"AE1", "AF1", "AG1",
			"AH1", "AI1", "AJ1",
			"AK1", "AL1", "AM1",
			leaveYear2,
		},
	}

	header = [][]interface{}{
		{"H1", "Jumlah"},
		{"H2", leaveYear1},
		{"I2", "Jan"},
		{"J2", "Feb"},
		{"K2", "Mar"},
		{"L2", "Apr"},
		{"M2", "Mei"},
		{"N2", "Juni"},
		{"O2", "Juli"},
		{"P2", "Agt"},
		{"Q2", "Sept"},
		{"R2", "Okt"},
		{"S2", "Nov"},
		{"T2", "Des"},
		{"U1", "Sisa"},
		{"U2", yearStr1},
		{"AA1", "Jumlah"},
		{"AA2", leaveYear2},
		{"AB2", "Jan"},
		{"AC2", "Feb"},
		{"AD2", "Mar"},
		{"AE2", "Apr"},
		{"AF2", "Mei"},
		{"AG2", "Juni"},
		{"AH2", "Juli"},
		{"AI2", "Agt"},
		{"AJ2", "Sept"},
		{"AK2", "Okt"},
		{"AL2", "Nov"},
		{"AM2", "Des"},
		{"AN1", "Sisa"},
		{"AN2", yearStr2},
	}

	//--- Open Excel Job
	file = excelize.NewFile()

	//--- Coordinate Y
	if err = input.excelHeader(file, sheetName, headerMrgY, true); err.Error != nil {
		return
	}

	//--- Coordinate X
	if err = input.excelHeader(file, sheetName, headerMrgX, true); err.Error != nil {
		return
	}

	//--- Non Coordinate
	if err = input.excelHeader(file, sheetName, header, false); err.Error != nil {
		return
	}

	//--- Content Year 1
	if err = input.excelContent(file, 'A', sheetName, datas1, false); err.Error != nil {
		return
	}

	//--- Content Year 2
	if err = input.excelContent(file, 'A', sheetName, datas2, true); err.Error != nil {
		return
	}

	//--- Create Directory
	dirPath := baseDirectory + baseImportPath + baseRootPath
	if errs = os.MkdirAll(dirPath, 0770); errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	//--- Save Excel
	errs = file.SaveAs(directory)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	//--- File Stat
	fileInfo, errs = os.Stat(directory)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	//--- File Content
	fileContent, errs = ioutil.ReadFile(directory)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	files := []in.MultipartFileDTO{
		{
			FileContent: fileContent,
			Filename:    fileInfo.Name(),
			Size:        fileInfo.Size(),
			Host:        config.ApplicationConfiguration.GetCDN().Host,
			Path:        config.ApplicationConfiguration.GetCDN().RootPath,
			FileID:      contextModel.AuthAccessTokenModel.ResourceUserID,
		},
	}

	//--- Upload File
	fileUploadID, err = input.uploadAttachmentFile(tx, files, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportAnnualLeaveService) excelHeader(file *excelize.File, sheetName string, header [][]interface{}, isMergedCell bool) (err errorModel.ErrorModel) {
	var (
		fileName = "ExcelReportAnnualLeaveService.go"
		funcName = "excelHeader"
	)

	for _, cell := range header {
		var (
			errs  error
			style int
		)

		if errs = file.SetCellValue(sheetName, cell[0].(string), cell[len(cell)-1]); errs != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errs)
			return
		}

		if isMergedCell {
			if errs = file.MergeCell(sheetName, cell[0].(string), cell[len(cell)-2].(string)); errs != nil {
				err = errorModel.GenerateUnknownError(fileName, funcName, errs)
				return
			}
		}

		style, errs = file.NewStyle(&excelize.Style{
			Border: []excelize.Border{
				{
					Type:  "top",
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
				{
					Type:  "bottom",
					Color: "#000000",
					Style: 1,
				},
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Pattern: 1,
				Color:   []string{"A9D08E"},
			},
			Font: &excelize.Font{
				Bold:   true,
				Family: "Batang",
				Size:   12,
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:   "center",
			},
		})

		if errs = file.SetCellStyle(sheetName, cell[0].(string), cell[len(cell)-2].(string), style); errs != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errs)
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportAnnualLeaveService) excelContent(file *excelize.File, startPoint int32, sheetName string, datas []repository.EmployeeLeaveReportModel, isDatas2 bool) (err errorModel.ErrorModel) {
	var (
		fileName  = "ExcelReportAnnualLeaveService.go"
		funcName  = "excelContent"
		rowNumber int64
		row       int
		font      = excelize.Font{Family: "Batang", Size: 12, Bold: false}
		fontName  = excelize.Font{Family: "Batang", Size: 12, Bold: true}
		align     = excelize.Alignment{Horizontal: "center", Vertical: "center"}
		alignName = excelize.Alignment{Horizontal: "left", Vertical: "center"}
	)

	for _, data := range datas {
		var (
			errs          error
			cell          string
			col           int32
			totOnMonth    int
			commentBundle []excelize.RichTextRun
		)

		if rowNumber != data.RowNumber.Int64 {
			if !isDatas2 {
				//--- Row Number
				if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.RowNumber.Int64, isDatas2, font, align); err.Error != nil {
					return
				}

				//--- NIK
				if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.IDCard.String, isDatas2, font, align); err.Error != nil {
					return
				}

				//--- Name
				if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.Name.String, isDatas2, fontName, alignName); err.Error != nil {
					return
				}

				//--- Position
				if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.Position.String, isDatas2, font, align); err.Error != nil {
					return
				}

				//--- Division
				if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.Department.String, isDatas2, font, align); err.Error != nil {
					return
				}

				//--- Join Date
				if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.DateJoin.String, isDatas2, font, align); err.Error != nil {
					return
				}

				//--- Probation Date
				if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.DateProbation.String, isDatas2, font, align); err.Error != nil {
					return
				}
			}

			//--- Leave Tot
			if data.TotalLeave.Int64 > 0 {
				if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.TotalLeave.Int64, isDatas2, font, align); err.Error != nil {
					return
				}
			} else {
				if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, nil, isDatas2, font, align); err.Error != nil {
					return
				}
			}

			//--- Leave Current
			col += 12
			if data.CurrentLeave.Int64 > 0 {
				if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.CurrentLeave.Int64, isDatas2, font, align); err.Error != nil {
					return
				}
			} else {
				if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, nil, isDatas2, font, align); err.Error != nil {
					return
				}
			}

			//--- State
			rowNumber = data.RowNumber.Int64
			row++
		}

		//--- Month Detail
		cell = fmt.Sprintf("%c%d", startPoint+int32(data.MonthLeave.Int64)+7, data.RowNumber.Int64+2)
		if isDatas2 {
			cell = fmt.Sprintf("%c%c%d", startPoint, startPoint+int32(data.MonthLeave.Int64), data.RowNumber.Int64+2)
		}

		for _, itemDate := range data.DetailLeaveList {
			var (
				comment string
				lg      = len(itemDate.Date)
			)

			if lg > 0 {
				if len(itemDate.Date) < 2 {
					comment = itemDate.Date[0].Time.Format(constanta.DefaultInstallationTimeFormat) + " : " + itemDate.Description.String + ".\n"
				} else {
					comment = itemDate.Date[0].Time.Format(constanta.DefaultInstallationTimeFormat) + " - " + itemDate.Date[len(itemDate.Date)-1].Time.Format(constanta.DefaultInstallationTimeFormat) + " : " + itemDate.Description.String + ".\n"
				}

				//--- Add Comment Bundle
				commentBundle = append(commentBundle, excelize.RichTextRun{Text: comment})
			}
			totOnMonth += lg
		}

		var style int
		if style, err = input.setCellStyleDefault(file, font, align); err.Error != nil {
			return
		}

		if totOnMonth > 0 {
			if errs = file.SetCellValue(sheetName, cell, totOnMonth); errs != nil {
				err = errorModel.GenerateUnknownError(fileName, funcName, errs)
				return
			}

			if errs = file.SetCellStyle(sheetName, cell, cell, style); errs != nil {
				err = errorModel.GenerateUnknownError(fileName, funcName, errs)
				return
			}

			//--- Add Comment
			_ = file.AddComment(sheetName, excelize.Comment{
				Cell:      cell,
				Paragraph: commentBundle,
			})
		} else {
			if errs = file.SetCellValue(sheetName, cell, nil); errs != nil {
				err = errorModel.GenerateUnknownError(fileName, funcName, errs)
				return
			}

			if errs = file.SetCellStyle(sheetName, cell, cell, style); errs != nil {
				err = errorModel.GenerateUnknownError(fileName, funcName, errs)
				return
			}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input reportAnnualLeaveService) setCellValueDefault(file *excelize.File, sheetName string, startPoint *int32, col *int32, row int, value interface{}, isDatas2 bool, font excelize.Font, align excelize.Alignment) (err errorModel.ErrorModel) {
	var (
		fileName = "ExcelReportAnnualLeaveService.go"
		funcName = "setCellValueDefault"
		cell     = fmt.Sprintf("%c%d", *startPoint+*col, row+3)
		errs     error
		style    int
	)

	if isDatas2 {
		cell = fmt.Sprintf("%c%c%d", *startPoint, *startPoint+*col, row+3)
	}

	if style, err = input.setCellStyleDefault(file, font, align); err.Error != nil {
		return
	}

	if errs = file.SetCellValue(sheetName, cell, value); errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	if errs = file.SetCellStyle(sheetName, cell, cell, style); errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	*col++
	return
}

func (input reportAnnualLeaveService) setCellStyleDefault(file *excelize.File, font excelize.Font, align excelize.Alignment) (style int, err errorModel.ErrorModel) {
	var (
		fileName = "ExcelReportAnnualLeaveService.go"
		funcName = "setCellStyleDefault"
		errs     error
	)

	style, errs = file.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{
				Type:  "top",
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
			{
				Type:  "bottom",
				Color: "#000000",
				Style: 1,
			},
		},
		Font:      &font,
		Alignment: &align,
	})
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
