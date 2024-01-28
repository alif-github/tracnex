package AbsentService

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
	"nexsoft.co.id/nextrac2/service"
	"strings"
	"time"
)

func (input absentService) UploadFileAbsent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		fileName     = "UploadFileAbsentService.go"
		funcName     = "UploadFileAbsent"
		records      *excelize.File
		inputStruct  []in.Absent
		inputMapping []in.AbsentMapping
	)

	//--- Read and Handle Upload Excel
	records, err = input.readAndHandleUploadExcel(request, "file-absent")
	if err.Error != nil {
		return
	}

	//--- Get Rows
	rows, errors := records.GetRows("Table 1")
	if errors != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errors)
		return
	}

	rows[0][0] = strings.Trim(rows[0][0], "Dari ")
	s := strings.Split(rows[0][0], " s/d ")
	periodStart := s[0]
	periodEnd := s[1]

	for i, row := range rows {
		if i < 2 {
			continue //--- Header Continue
		}

		inputMapping = append(inputMapping, in.AbsentMapping{
			AbsentID:        row[2],
			IDCard:          row[3],
			NormalDays:      row[4],
			ActualDays:      row[5],
			Absents:         row[6],
			Overdue:         row[7],
			LeaveEarly:      row[8],
			Overtime:        row[9],
			NumbersOfLeave:  row[10],
			LeavingDuties:   row[11],
			NumbersIn:       row[12],
			NumbersOut:      row[13],
			Scan:            row[14],
			SickLeave:       row[15],
			PaidLeave:       row[16],
			PermissionLeave: row[17],
			WorkHours:       row[18],
			PercentAbsent:   row[19],
			PeriodStartStr:  periodStart,
			PeriodEndStr:    periodEnd,
		})
	}

	//--- Validate Insert
	inputStruct, err = in.Absent{}.ValidateInsertAbsent(inputMapping)
	if err.Error != nil {
		return
	}

	//--- Audit System
	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doInsertToDBAbsent, func(i interface{}, model applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input absentService) doInsertToDBAbsent(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct []in.Absent
		modelAdd    []repository.AbsentModel
		modelEdit   []repository.AbsentModel
		resultID    []int64
	)

	if inputStructInterface != nil {
		inputStruct = inputStructInterface.([]in.Absent)
	}

	//--- Create Model And Check NIK
	modelAdd, modelEdit, err = input.createModel(inputStruct, *contextModel, timeNow)
	if err.Error != nil {
		return
	}

	if len(modelEdit) > 0 {
		//--- Console Log
		fmt.Println(fmt.Sprintf(`Edit Absent -> %d`, len(modelEdit)))

		//--- Audit Edit
		for _, iEdit := range modelEdit {
			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.AbsentDAO.TableName, iEdit.ID.Int64, 0)...)
		}

		//--- Update Absents
		err = dao.AbsentDAO.UpdateMultipleAbsent(tx, modelEdit)
		if err.Error != nil {
			return
		}
	}

	if len(modelAdd) > 0 {
		//--- Console Log
		fmt.Println(fmt.Sprintf(`Add Absent -> %d`, len(modelAdd)))

		//--- Insert Absents
		resultID, err = dao.AbsentDAO.InsertMultipleAbsent(tx, modelAdd)
		if err.Error != nil {
			return
		}

		//--- Data Audit
		for _, i := range resultID {
			dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.AbsentDAO.TableName, i, 0)...)
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input absentService) createModel(item []in.Absent, cm applicationModel.ContextModel, timeNow time.Time) (modelAdd, modelEdit []repository.AbsentModel, err errorModel.ErrorModel) {
	var (
		fileName = "UploadFileAbsentService.go"
		funcName = "createdModel"
		db       = serverconfig.ServerAttribute.DBConnection
	)

	for _, i := range item {
		var (
			employeeOnDB repository.EmployeeModel
			absentOnDB   repository.AbsentModel
			isExist      bool
		)

		//--- ID Card Check
		if i.IDCard == "" {
			continue //--- ID Card Empty, Continue
		}

		//--- Check Employee By NIK
		isExist, employeeOnDB, err = dao.EmployeeDAO.CheckEmployeeByNIKParamOnly(db, i.IDCard)
		if err.Error != nil {
			return
		}

		//--- If Not Exist Then Error
		if !isExist {
			err = errorModel.GenerateUnknownDataError(fileName, funcName, fmt.Sprintf(`NIK - %s`, i.IDCard))
			return
		}

		//--- Check Action Data
		absentOnDB, err = dao.AbsentDAO.GetAbsentID(db, repository.AbsentModel{
			AbsentID:    sql.NullInt64{Int64: i.AbsentID},
			PeriodStart: sql.NullTime{Time: i.PeriodStart},
			PeriodEnd:   sql.NullTime{Time: i.PeriodEnd},
		})
		if err.Error != nil {
			return
		}

		if absentOnDB.ID.Int64 > 0 {
			//--- Create Model Edit
			modelEdit = append(modelEdit, repository.AbsentModel{
				ID:              sql.NullInt64{Int64: absentOnDB.ID.Int64},
				EmployeeID:      sql.NullInt64{Int64: employeeOnDB.ID.Int64},
				AbsentID:        sql.NullInt64{Int64: i.AbsentID},
				NormalDays:      sql.NullInt64{Int64: i.NormalDays},
				ActualDays:      sql.NullInt64{Int64: i.ActualDays},
				Absent:          sql.NullInt64{Int64: i.Absents},
				Overdue:         sql.NullInt64{Int64: i.Overdue},
				LeaveEarly:      sql.NullInt64{Int64: i.LeaveEarly},
				Overtime:        sql.NullInt64{Int64: i.Overtime},
				NumberOfLeave:   sql.NullInt64{Int64: i.NumbersOfLeave},
				LeavingDuties:   sql.NullInt64{Int64: i.LeavingDuties},
				NumbersIn:       sql.NullInt64{Int64: i.NumbersIn},
				NumbersOut:      sql.NullInt64{Int64: i.NumbersOut},
				Scan:            sql.NullInt64{Int64: i.Scan},
				SickLeave:       sql.NullInt64{Int64: i.SickLeave},
				PaidLeave:       sql.NullInt64{Int64: i.PaidLeave},
				PermissionLeave: sql.NullInt64{Int64: i.PermissionLeave},
				WorkHours:       sql.NullInt64{Int64: i.WorkHours},
				PercentAbsent:   sql.NullFloat64{Float64: i.PercentAbsent},
				PeriodStart:     sql.NullTime{Time: i.PeriodStart},
				PeriodEnd:       sql.NullTime{Time: i.PeriodEnd},
				UpdatedBy:       sql.NullInt64{Int64: cm.AuthAccessTokenModel.ResourceUserID},
				UpdatedClient:   sql.NullString{String: cm.AuthAccessTokenModel.ClientID},
				UpdatedAt:       sql.NullTime{Time: timeNow},
			})

			//--- Continue When Finished
			continue
		}

		//--- Create Model Add
		modelAdd = append(modelAdd, repository.AbsentModel{
			EmployeeID:      sql.NullInt64{Int64: employeeOnDB.ID.Int64},
			AbsentID:        sql.NullInt64{Int64: i.AbsentID},
			NormalDays:      sql.NullInt64{Int64: i.NormalDays},
			ActualDays:      sql.NullInt64{Int64: i.ActualDays},
			Absent:          sql.NullInt64{Int64: i.Absents},
			Overdue:         sql.NullInt64{Int64: i.Overdue},
			LeaveEarly:      sql.NullInt64{Int64: i.LeaveEarly},
			Overtime:        sql.NullInt64{Int64: i.Overtime},
			NumberOfLeave:   sql.NullInt64{Int64: i.NumbersOfLeave},
			LeavingDuties:   sql.NullInt64{Int64: i.LeavingDuties},
			NumbersIn:       sql.NullInt64{Int64: i.NumbersIn},
			NumbersOut:      sql.NullInt64{Int64: i.NumbersOut},
			Scan:            sql.NullInt64{Int64: i.Scan},
			SickLeave:       sql.NullInt64{Int64: i.SickLeave},
			PaidLeave:       sql.NullInt64{Int64: i.PaidLeave},
			PermissionLeave: sql.NullInt64{Int64: i.PermissionLeave},
			WorkHours:       sql.NullInt64{Int64: i.WorkHours},
			PercentAbsent:   sql.NullFloat64{Float64: i.PercentAbsent},
			PeriodStart:     sql.NullTime{Time: i.PeriodStart},
			PeriodEnd:       sql.NullTime{Time: i.PeriodEnd},
			CreatedBy:       sql.NullInt64{Int64: cm.AuthAccessTokenModel.ResourceUserID},
			CreatedClient:   sql.NullString{String: cm.AuthAccessTokenModel.ClientID},
			CreatedAt:       sql.NullTime{Time: timeNow},
			UpdatedBy:       sql.NullInt64{Int64: cm.AuthAccessTokenModel.ResourceUserID},
			UpdatedClient:   sql.NullString{String: cm.AuthAccessTokenModel.ClientID},
			UpdatedAt:       sql.NullTime{Time: timeNow},
		})
	}
	return
}
