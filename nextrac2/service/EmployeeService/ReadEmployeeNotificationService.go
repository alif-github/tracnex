package EmployeeService

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"time"
)

func (input employeeService) ReadEmployeeNotification(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	funcName := "GetListEmployeeNotification"

	output.Data.Content, errModel = input.ServiceWithDataAuditPreparedByService(funcName, nil, contextModel, input.readEmployeeNotification, func(i interface{}, model applicationModel.ContextModel) {})
	if errModel.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input employeeService) readEmployeeNotification(tx *sql.Tx, _ interface{}, contextModel *applicationModel.ContextModel, _ time.Time) (_ interface{}, auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	var (
		db          = serverconfig.ServerAttribute.DBConnection
		funcName    = "readEmployeeNotification"
		inputStruct = in.GetListDataDTO{
			AbstractDTO: in.AbstractDTO{
				Page:  -99,
				Limit: 0,
			},
		}
	)

	employee, errModel := dao.EmployeeDAO.GetByUserId(db, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	if !employee.ID.Valid {
		errModel = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employee")
		return
	}

	memberIdList := in.MemberList{}
	_ = json.Unmarshal([]byte(employee.Member.String), &memberIdList)

	auditList, errModel := dao.AuditSystemDAO.GetListEmployeeNotification(db, inputStruct, []in.SearchByParam{}, repository.EmployeeNotification{
		EmployeeId:     employee.ID.Int64,
		MemberIdList:   memberIdList.MemberID,
		FilterByIsRead: true,
		IsRead:         false,
	})
	if errModel.Error != nil {
		return
	}

	errModel = input.updateReadStatus(tx, auditList)
	return
}

func (input employeeService) updateReadStatus(tx *sql.Tx, auditList []interface{}) (errModel errorModel.ErrorModel) {
	for _, audit := range auditList {
		var notification in.EmployeeNotification

		data, _ := audit.(repository.AuditSystemModel)
		_ = json.Unmarshal([]byte(data.Description.String), &notification)

		notification.IsRead = true
		bytes, _ := json.Marshal(notification)
		description := string(bytes)

		errModel = dao.AuditSystemDAO.UpdateDescriptionByIdTx(tx, repository.AuditSystemModel{
			ID:            sql.NullInt64{Int64: data.ID.Int64},
			Description:   sql.NullString{String: description},
		})
		if errModel.Error != nil {
			return
		}
	}

	return errorModel.GenerateNonErrorModel()
}