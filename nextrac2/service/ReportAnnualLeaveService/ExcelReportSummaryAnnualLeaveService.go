package ReportAnnualLeaveService

import (
	"database/sql"
	"github.com/360EntSecGroup-Skylar/excelize"
	"io/ioutil"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"os"
	"time"
)

func (input reportAnnualLeaveService) excelSummaryProduce(data []out.EmployeeLeaveYearly, contextModel *applicationModel.ContextModel, timeNow time.Time) (fileUploadID int64, err errorModel.ErrorModel) {
	var (
		fileName       = "ExcelReportSummaryAnnualLeaveService.go"
		funcName       = "excelSummaryProduce"
		file           *excelize.File
		header         [][]interface{}
		sheetName      = "Sheet1"
		errs           error
		tx             *sql.Tx
		fileInfo       os.FileInfo
		fileContent    []byte
		db             = serverconfig.ServerAttribute.DBConnection
		baseDirectory  = config.ApplicationConfiguration.GetDataDirectory().BaseDirectoryPath
		baseImportPath = config.ApplicationConfiguration.GetDataDirectory().ImportPath
		baseRootPath   = config.ApplicationConfiguration.GetDataDirectory().ReportLeavePath
		fileNamePath   = "report-summary-leave.xlsx"
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

	header = [][]interface{}{
		{"A1", "A2", "NIK"},
		{"B1", "B2", "Nama"},
		{"C1", "C2", "Department"},
		{"D1", "D2", "Level/Grade"},
		{"E1", "E2", "Cuti Saat Ini"},
		{"F1", "F2", "Cuti Sebelum"},
		{"G1", "G2", "Cuti Terhutang"},
	}

	//--- Open Excel Job
	file = excelize.NewFile()

	//--- Header
	if err = input.excelHeader(file, sheetName, header, true); err.Error != nil {
		return
	}

	//--- Content Year
	if err = input.excelSummaryContent(file, 'A', sheetName, data); err.Error != nil {
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

func (input reportAnnualLeaveService) excelSummaryContent(file *excelize.File, startPoint int32, sheetName string, datas []out.EmployeeLeaveYearly) (err errorModel.ErrorModel) {
	var (
		row       int
		isDatas2  bool
		font      = excelize.Font{Family: "Batang", Size: 12, Bold: false}
		fontName  = excelize.Font{Family: "Batang", Size: 12, Bold: true}
		align     = excelize.Alignment{Horizontal: "center", Vertical: "center"}
		alignName = excelize.Alignment{Horizontal: "left", Vertical: "center"}
	)

	for _, data := range datas {
		var col int32

		//--- NIK
		if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.IDCard, isDatas2, font, align); err.Error != nil {
			return
		}

		//--- Name
		if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.FullName, isDatas2, fontName, alignName); err.Error != nil {
			return
		}

		//--- Department
		if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.Department, isDatas2, font, align); err.Error != nil {
			return
		}

		//--- Level And Grade
		levelGrade := data.Level + "/" + data.Grade
		if levelGrade == "/" {
			levelGrade = ""
		}

		if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, levelGrade, isDatas2, font, align); err.Error != nil {
			return
		}

		//--- Current Leave
		if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.CurrentLeaveThisPeriod, isDatas2, font, align); err.Error != nil {
			return
		}

		//--- Leave Before
		if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.LastLeaveBeforePeriod, isDatas2, font, align); err.Error != nil {
			return
		}

		//--- Owing Leave
		if err = input.setCellValueDefault(file, sheetName, &startPoint, &col, row, data.OwingLeave, isDatas2, font, align); err.Error != nil {
			return
		}

		row++
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
