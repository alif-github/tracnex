package EmployeeService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
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
	"strconv"
	"time"
)

type cancellationRequest struct {
	Type        string
	EmployeeId  int64
	Firstname	string
	Lastname	string
	StrDateList string
	Date        time.Time
	Status      string
	UpdatedAt   time.Time
	Leaders		[]repository.EmployeeModel
}

func (input employeeService) CancelEmployeeRequest(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	funcName := "CancelEmployeeRequest"

	body, errModel := input.readCancellationBody(request, contextModel)
	if errModel.Error != nil {
		return
	}

	if errModel = body.ValidateCancellation(); errModel.Error != nil {
		return
	}

	output.Data.Content, errModel = input.ServiceWithDataAuditPreparedByService(funcName, body, contextModel, input.cancelEmployeeRequest, func(i interface{}, model applicationModel.ContextModel) {})
	if errModel.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input employeeService) cancelEmployeeRequest(tx *sql.Tx, data interface{}, contextModel *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	body, _ := data.(*in.EmployeeRequestHistory)

	cancelReq, errModel := input.getEmployeeRequest(tx, body)
	if errModel.Error != nil {
		return
	}

	cancelReq.Leaders, errModel = dao.EmployeeDAO.GetListByMemberId(serverconfig.ServerAttribute.DBConnection, cancelReq.EmployeeId)
	if errModel.Error != nil {
		return
	}

	if errModel = input.validateCancellationRequest(body.UpdatedAt, cancelReq.UpdatedAt); errModel.Error != nil {
		return
	}

	body.Status = constanta.PendingCancellationRequestStatus

	if cancelReq.Status == constanta.PendingRequestStatus {
		body.Status = constanta.CancelledRequestStatus
	}

	auditData, errModel = input.doCancelEmployeeRequest(tx, body, cancelReq, contextModel, now)
	return
}

func (input employeeService) doCancelEmployeeRequest(tx *sql.Tx, body *in.EmployeeRequestHistory, leaveReq cancellationRequest, contextModel *applicationModel.ContextModel, now time.Time) (auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	if leaveReq.Type == constanta.ReimbursementType {
		return input.cancelReimbursementRequest(tx, body, contextModel, now)
	}

	return input.cancelLeaveRequest(tx, body, leaveReq, contextModel, now)
}

func (input employeeService) cancelReimbursementRequest(tx *sql.Tx, body *in.EmployeeRequestHistory, contextModel *applicationModel.ContextModel, now time.Time) (auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	model := repository.EmployeeReimbursement{
		ID:                 sql.NullInt64{Int64: body.ID},
		Status: 		    sql.NullString{String: body.Status},
		CancellationReason: sql.NullString{String: body.CancellationReason},
		UpdatedAt:          sql.NullTime{Time: now},
		UpdatedBy:          sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:      sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, now, dao.EmployeeReimbursementDAO.TableName, body.ID, contextModel.LimitedByCreatedBy)...)

	errModel = dao.EmployeeReimbursementDAO.UpdateStatusAndCancellationReasonTx(tx, model)
	return
}

func (input employeeService) cancelLeaveRequest(tx *sql.Tx, body *in.EmployeeRequestHistory, leaveReq cancellationRequest, contextModel *applicationModel.ContextModel, now time.Time) (auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	model := repository.EmployeeLeaveModel{
		ID:                 sql.NullInt64{Int64: body.ID},
		Status: 		    sql.NullString{String: body.Status},
		CancellationReason: sql.NullString{String: body.CancellationReason},
		UpdatedAt:          sql.NullTime{Time: now},
		UpdatedBy:          sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:      sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, now, dao.EmployeeLeaveDAO.TableName, body.ID, contextModel.LimitedByCreatedBy)...)
	if errModel = dao.EmployeeLeaveDAO.UpdateStatusAndCancellationReasonTx(tx, model); errModel.Error != nil {
		return
	}

	if body.Status == constanta.PendingCancellationRequestStatus {
		var dateList []time.Time

		_ = json.Unmarshal([]byte(leaveReq.StrDateList), &dateList)

		auditData[0].Employee.Firstname = leaveReq.Firstname
		auditData[0].Employee.Lastname = leaveReq.Lastname

		input.addNotificationIntoAudit(&auditData[0], &in.EmployeeNotification{
			IsMobileNotification:        true,
			IsRequestingForApproval:     false,
			IsRequestingForCancellation: true,
			IsCancellation:              true,
			EmployeeId:                  leaveReq.EmployeeId,
			RequestType:                 leaveReq.Type,
			Status:                      body.Status,
			Date:                        dateList,
		})

		go input.sendCancellationRequestNotifications(leaveReq.Leaders, leaveReq, now)
	}
	return
}

func (input employeeService) sendCancellationRequestNotifications(leaders []repository.EmployeeModel, cancelReq cancellationRequest, now time.Time) {
	for _, leader := range leaders {
		/*
			Send Notification
		*/
		leaderName := fmt.Sprintf("%s %s", leader.FirstName.String, leader.LastName.String)
		employeeName := fmt.Sprintf("%s %s", cancelReq.Firstname, cancelReq.Lastname)
		requestType := input.getEmailRequestType(cancelReq.Type)
		createdAt := input.timeToString(now)

		message := fmt.Sprintf(constanta.CancellationRequestEmailBody, leaderName, employeeName, requestType, createdAt)
		input.sendNotificationToEmployee(leader.ClientId.String, message)

		/*
			Send Email
		*/
		go input.sendEmail(input.toMailAddress(leader.Email.String), "", constanta.CancelRequestApprovalHRISSubject, message)
	}
}


func (input employeeService) getEmployeeRequest(tx *sql.Tx, body *in.EmployeeRequestHistory) (result cancellationRequest, errModel errorModel.ErrorModel) {
	funcName := "getEmployeeRequest"

	/*
		Reimbursement
	*/
	if body.RequestType == constanta.ReimbursementType {
		var reimbursementReq repository.EmployeeReimbursement

		reimbursementReq, errModel = dao.EmployeeReimbursementDAO.GetByIdForUpdate(tx, body.ID, constanta.PendingRequestStatus)
		if errModel.Error != nil {
			return
		}

		if !reimbursementReq.ID.Valid {
			errModel = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employee_request")
			return
		}

		result = cancellationRequest{
			Type:       constanta.ReimbursementType,
			EmployeeId: reimbursementReq.EmployeeId.Int64,
			Firstname:  reimbursementReq.Firstname.String,
			Lastname:  	reimbursementReq.Lastname.String,
			Date:       reimbursementReq.Date.Time,
			Status:     reimbursementReq.Status.String,
			UpdatedAt:  reimbursementReq.UpdatedAt.Time,
		}
		return
	}

	/*
		Leave | Permit | Sick Leave
	*/
	leaveReq, errModel := dao.EmployeeLeaveDAO.GetByIdAndStatusesTx(tx, body.ID, []string{
		constanta.PendingRequestStatus,
		constanta.ApprovedRequestStatus,
	})
	if errModel.Error != nil {
		return
	}

	if !leaveReq.ID.Valid {
		errModel = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employee_request")
		return
	}

	if (leaveReq.Type.String == constanta.SickLeaveType && leaveReq.Status.String != constanta.PendingRequestStatus) || (leaveReq.Status.String == constanta.ApprovedRequestStatus && leaveReq.CancellationReason.String != "") {
		errModel = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employee_request")
		return
	}

	result = cancellationRequest{
		Type:        leaveReq.Type.String,
		EmployeeId:  leaveReq.EmployeeId.Int64,
		Firstname:   leaveReq.Firstname.String,
		Lastname:    leaveReq.Lastname.String,
		StrDateList: leaveReq.Date.String,
		Date:        time.Time{},
		Status:      leaveReq.Status.String,
		UpdatedAt:   leaveReq.UpdatedAt.Time,
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) validateCancellationRequest(updatedAtOnBody time.Time, updatedAtOnDB time.Time) errorModel.ErrorModel {
	funcName := "validateCancellationRequest"

	if updatedAtOnBody != updatedAtOnDB {
		return errorModel.GenerateDataLockedError(input.FileName, funcName, "updated_at")
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) readCancellationBody(request *http.Request, contextModel *applicationModel.ContextModel) (body *in.EmployeeRequestHistory, errModel errorModel.ErrorModel) {
	var (
		funcName   = "readCancellationBody"
		stringBody string
	)

	stringBody, errModel = input.ReadBody(request, contextModel)
	if errModel.Error != nil {
		return
	}

	body = &in.EmployeeRequestHistory{}
	if err := json.Unmarshal([]byte(stringBody), body); err != nil {
		return nil, errorModel.GenerateInvalidRequestError(input.FileName, funcName, err)
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	body.ID = int64(id)
	return
}
