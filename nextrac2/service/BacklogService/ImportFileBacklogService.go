package BacklogService

import (
	"database/sql"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"strconv"
	"strings"
)

func (input backlogService) UnmarshalFileBacklog(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		records     *excelize.File
		inputStruct in.ImportBacklogRequest
	)

	records, inputStruct, err = input.readBodyAndValidateFileCSVBacklog(request, inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doUnmarshalFileBacklog(records, inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input backlogService) doUnmarshalFileBacklog(records *excelize.File, inputStruct in.ImportBacklogRequest, contextModel *applicationModel.ContextModel) (listBacklog interface{}, err errorModel.ErrorModel) {
	var (
		fileName             = "ImportFileBacklogService.go"
		funcName             = "doUnmarshalFileBacklog"
		headerIndex          = 6
		skippedHeaderIndex   = 7
		listBacklogQAQC      []out.BacklogFromFileQaQcResponse
		listBacklogDeveloper []out.BacklogFromFileResponse
		scope                map[string]interface{}
	)

	mappingScopeDB := make(map[string]applicationModel.MappingScopeDB)
	mappingScopeDB[constanta.EmployeeDataScope] = applicationModel.MappingScopeDB{
		View:  "e.id",
		Count: "e.id",
	}

	createdBy := contextModel.LimitedByCreatedBy
	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	rows, errs := records.GetRows(records.GetSheetName(0))
	if errs != nil {
		if errs.Error() == "sheet backlog is not exist" {
			err = errorModel.GenerateSimpleErrorModel(400, "Tidak dapat menemukan sheet yang sesuai format")
		}

		return
	}

	// validasi form yang di input sesuai dengan department
	if inputStruct.DepartmentCode == constanta.DepartmentQAQC {
		if rows[6][0] != "Feature*" {
			err = errorModel.GenerateSimpleErrorModel(400, "File Import tidak sesuai dengen ketentuan QA/QC")
			return
		}
	}

	// validasi form yang di input sesuai dengan department
	if inputStruct.DepartmentCode == constanta.DepartmentDeveloper {
		if rows[6][0] != "Layer 1*" {
			err = errorModel.GenerateSimpleErrorModel(400, "File Import tidak sesuai dengen ketentuan Developer")
			return
		}
	}

	// VALIDATION DATA CAN BE PROCESS
	totalHeader := len(rows[headerIndex])
	rows = rows[skippedHeaderIndex:]
	for _, row := range rows {
		if len(row) != totalHeader {
			err = errorModel.GenerateSimpleErrorModel(400, "Data tidak dapat diproses, semua kolom belum terisi")
			return
		}
	}

	listPic := make(map[string]int64)

	// For File Dev Only
	if inputStruct.DepartmentCode == constanta.DepartmentDeveloper {
		for _, row := range rows {
			//redmineNumber, _ := strconv.ParseInt(row[5], 10, 64)
			statusCapitalize := strings.Title(row[11])
			mandays, _ := strconv.ParseFloat(row[12], 64)
			mandaysDone, _ := strconv.ParseFloat(row[12], 64)

			var (
				picId          int64
				resultEmployee repository.EmployeeModel
			)

			picName := strings.ToLower(row[10])
			if val, ok := listPic[picName]; ok {
				picId = val
			} else {
				resultEmployee, err = dao.EmployeeDAO.GetEmployeeIdByFullName(serverconfig.ServerAttribute.DBConnection, repository.EmployeeModel{
					Name: sql.NullString{String: picName},
				}, createdBy, scope, mappingScopeDB)

				if err.Error != nil {
					return
				}

				picId = resultEmployee.ID.Int64
				if picId < 1 {
					err = errorModel.GenerateUnknownDataError(fileName, funcName, fmt.Sprintf(`Name : %s`, picName))
					return
				}

				listPic[picName] = picId
			}

			if row != nil {
				backlog := out.BacklogFromFileResponse{
					DepartmentCode: constanta.DepartmentDeveloper,
					Layer1:         row[0],
					Layer2:         row[1],
					Layer3:         row[2],
					Layer4:         row[3],
					Layer5:         row[4],
					Redmine:        row[5],
					Sprint:         row[6],
					SprintName:     row[7],
					PicId:          picId,
					PicName:        resultEmployee.Name.String,
					Status:         statusCapitalize,
					Mandays:        mandays,
					MandaysDone:    mandaysDone,
					FlowChanged:    row[21],
					AdditionalData: row[22],
					Note:           row[23],
					Url:            row[24],
					Page:           row[25],
				}

				listBacklogDeveloper = append(listBacklogDeveloper, backlog)
			}
		}

		listBacklog = listBacklogDeveloper
	}

	if inputStruct.DepartmentCode == constanta.DepartmentQAQC {
		for _, row := range rows {
			feature, _ := strconv.ParseInt(row[0], 10, 64)
			//redmineNumber, _ := strconv.ParseInt(row[4], 10, 64)
			referenceTicket, _ := strconv.ParseInt(row[9], 10, 64)
			statusCapitalize := strings.Title(row[13])
			mandays, _ := strconv.ParseFloat(row[14], 64)
			mandaysDone, _ := strconv.ParseFloat(row[15], 64)

			var (
				picId          int64
				resultEmployee repository.EmployeeModel
			)

			picName := strings.ToLower(row[12])
			if val, ok := listPic[picName]; ok {
				picId = val
			} else {
				resultEmployee, err = dao.EmployeeDAO.GetEmployeeIdByFullName(serverconfig.ServerAttribute.DBConnection, repository.EmployeeModel{
					Name: sql.NullString{String: picName},
				}, createdBy, scope, mappingScopeDB)

				if err.Error != nil {
					return
				}

				picId = resultEmployee.ID.Int64
				if picId < 1 {
					err = errorModel.GenerateUnknownDataError(fileName, funcName, fmt.Sprintf(`Name : %s`, picName))
					return
				}

				listPic[picName] = picId
			}

			if row != nil {
				backlogQAQC := out.BacklogFromFileQaQcResponse{
					DepartmentCode:  constanta.DepartmentQAQC,
					Feature:         feature,
					Tracker:         row[1],
					Layer1:          row[2],
					Layer2:          row[3],
					Layer3:          row[4],
					Layer4:          row[5],
					Layer5:          row[6],
					Subject:         row[7],
					Redmine:         row[8],
					ReferenceTicket: referenceTicket,
					Sprint:          row[10],
					FormChanged:     "-",
					PicId:           picId,
					PicName:         resultEmployee.Name.String,
					Status:          statusCapitalize,
					Mandays:         mandays,
					MandaysDone:     mandaysDone,
					FlowChanged:     row[16],
					AdditionalData:  row[17],
					Note:            row[18],
					Url:             row[19],
					Page:            row[20],
				}

				listBacklogQAQC = append(listBacklogQAQC, backlogQAQC)
			}
		}

		listBacklog = listBacklogQAQC
	}

	return
}
