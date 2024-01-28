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
	"nexsoft.co.id/nextrac2/service"
	"strconv"
	"strings"
	"time"
)

type approvalRequestModel struct {
	Id                 int64
	EmployeeId         int64
	TotalLeave         int64
	MedicalValue       float64
	Status             string
	AllowanceId        int64
	AllowanceType      string
	RequestType        string
	CancellationReason string
	BenefitId          int64
	UpdatedAt          time.Time
	TableName          string
	DateList           []time.Time
	ClientID		   string
	EmployeeName	   string
	Email			   string
	CreatedAt		   time.Time
}

func (input employeeService) UpdateApprovalStatus(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	funcName := "UpdateApprovalStatus"

	body, errModel := input.readApprovalBody(request, contextModel)
	if errModel.Error != nil {
		return
	}

	if errModel = body.ValidateApprovalRequest(); errModel.Error != nil {
		return
	}

	output.Data.Content, errModel = input.ServiceWithDataAuditPreparedByService(funcName, body, contextModel, input.updateApprovalStatus, func(i interface{}, model applicationModel.ContextModel) {})
	if errModel.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input employeeService) updateApprovalStatus(tx *sql.Tx, data interface{}, contextModel *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	body, _ := data.(*in.EmployeeApprovalRequest)

	approvalRequest, errModel := input.getApprovalRequestByTypeAndId(tx, body.RequestType, body.Id)
	if errModel.Error != nil {
		return
	}

	if errModel = input.validateApprovalRequest(approvalRequest, body); errModel.Error != nil {
		return
	}

	auditData, errModel = input.updateRequest(tx, body, approvalRequest, contextModel, now)
	if errModel.Error != nil {
		return
	}

	return
}

func (input employeeService) updateRequest(tx *sql.Tx, body *in.EmployeeApprovalRequest, approvalRequest *approvalRequestModel, contextModel *applicationModel.ContextModel, now time.Time) (auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	var (
		benefitAudit []repository.AuditSystemModel
		benefitApprovalAudit []repository.AuditSystemModel
		annualLeaveApprovalAudit []repository.AuditSystemModel
		cancellationBenefitAudit []repository.AuditSystemModel
		cancellationAudit []repository.AuditSystemModel
	)

	benefit, errModel := dao.EmployeeBenefitsDAO.GetByEmployeeIdForUpdate(tx, approvalRequest.EmployeeId)
	if errModel.Error != nil {
		return
	}

	/*
		Reimbursement Approval
	*/
	if approvalRequest.TableName == dao.EmployeeReimbursementDAO.TableName {
		benefitApprovalAudit, errModel = input.updateReimbursementApprovalStatus(tx, body, approvalRequest, contextModel, now)
		if errModel.Error != nil {
			return
		}

		auditData = append(auditData, benefitApprovalAudit...)
		return
	}

	/*
		Leave Request Cancellation Approval
	*/
	if input.isCancellationRequest(approvalRequest) {
		cancellationBenefitAudit, errModel = input.restoreAnnualLeave(tx, body, approvalRequest, benefit, contextModel, now)
		if errModel.Error != nil {
			return
		}

		cancellationAudit, errModel = input.updateCancellationApprovalStatus(tx, body, approvalRequest, contextModel, now)
		if errModel.Error != nil {
			return
		}

		auditData = append(auditData, cancellationBenefitAudit...)
		auditData = append(auditData, cancellationAudit...)
		return
	}

	/*
		Leave Approval
	*/
	benefitAudit, errModel = input.updateAnnualLeaveBenefit(tx, body, approvalRequest, benefit, contextModel, now)
	if errModel.Error != nil {
		return
	}

	annualLeaveApprovalAudit, errModel = input.updateAnnualLeaveApprovalStatus(tx, body, approvalRequest, contextModel, now)
	if errModel.Error != nil {
		return
	}

	auditData = append(auditData, benefitAudit...)
	auditData = append(auditData, annualLeaveApprovalAudit...)
	return
}

func (input employeeService) isCancellationRequest(approvalRequest *approvalRequestModel) bool {
	return approvalRequest.Status == constanta.PendingCancellationRequestStatus
}

func (input employeeService) restoreAnnualLeave(tx *sql.Tx, body *in.EmployeeApprovalRequest, approvalRequest *approvalRequestModel, benefit repository.EmployeeBenefitsModel, contextModel *applicationModel.ContextModel, now time.Time) (auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	if body.Status != constanta.ApprovedRequestStatus {
		return nil, errorModel.GenerateNonErrorModel()
	}

	if !input.isAnnualLeave(approvalRequest.AllowanceType) {
		return nil, errorModel.GenerateNonErrorModel()
	}

	facilitiesActive, errModel := dao.EmployeeFacilitiesActiveDAO.GetByAllowanceIdAndEmployeeLevelIdAndEmployeeGradeIdTx(tx, repository.EmployeeFacilitiesActiveModel{
		AllowanceID: sql.NullInt64{Int64: approvalRequest.AllowanceId},
		LevelID:  	 sql.NullInt64{Int64: benefit.EmployeeLevelID.Int64},
		GradeID: 	 sql.NullInt64{Int64: benefit.EmployeeGradeID.Int64},
	})
	if errModel.Error != nil {
		return
	}

	value, _ := strconv.Atoi(facilitiesActive.Value.String)
	maxValue := int64(value)

	currentValue, lastValue := input.increaseCurrentAndLastAnnualLeave(benefit.CurrentAnnualLeave.Int64, benefit.LastAnnualLeave.Int64, approvalRequest.TotalLeave, maxValue)

	benefitModel := repository.EmployeeBenefitsModel{
		EmployeeID:         sql.NullInt64{Int64: approvalRequest.EmployeeId},
		CurrentAnnualLeave: sql.NullInt64{Int64: currentValue},
		LastAnnualLeave:    sql.NullInt64{Int64: lastValue},
		UpdatedAt:          sql.NullTime{Time: now},
		UpdatedBy:          sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:      sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, now, dao.EmployeeBenefitsDAO.TableName, benefitModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	errModel = dao.EmployeeBenefitsDAO.UpdateCurrentAnnualLeaveAndLastAnnualLeaveByEmployeeIdTx(tx, benefitModel)
	if errModel.Error != nil {
		return
	}

	errModel = dao.EmployeeHistoryLeaveDAO.UpdateCurrentAnnualLeaveAndLastAnnualLeaveByEmployeeIdTx(tx, benefitModel)
	return
}

func (input employeeService) updateCancellationApprovalStatus(tx *sql.Tx, body *in.EmployeeApprovalRequest, approvalRequest *approvalRequestModel, contextModel *applicationModel.ContextModel, now time.Time) (auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	status := constanta.CancelledRequestStatus

	if body.Status == constanta.RejectedRequestStatus {
		status = constanta.ApprovedRequestStatus
	}

	model := repository.EmployeeLeaveModel{
		ID: 		   sql.NullInt64{Int64: approvalRequest.Id},
		Status:        sql.NullString{String: status},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, now, dao.EmployeeLeaveDAO.TableName, approvalRequest.Id, contextModel.LimitedByCreatedBy)...)
	if errModel = dao.EmployeeLeaveDAO.UpdateStatusTx(tx, model); errModel.Error != nil {
		return
	}

	/*
		Send Notification
	*/
	input.addNotificationIntoAudit(&auditData[0], &in.EmployeeNotification{
		IsMobileNotification:        true,
		IsRequestingForApproval:     false,
		IsRequestingForCancellation: false,
		IsCancellation:              true,
		EmployeeId:                  approvalRequest.EmployeeId,
		RequestType:                 approvalRequest.RequestType,
		Status:                      body.Status,
		Date:                        approvalRequest.DateList,
	})

	employeeName := approvalRequest.EmployeeName
	requestType := input.getEmailRequestType(approvalRequest.RequestType)
	createdAt := input.timeToString(approvalRequest.CreatedAt)

	message := fmt.Sprintf(constanta.CancellationRequestApprovedEmailBody, employeeName, requestType, createdAt)
	if body.Status == constanta.RejectedRequestStatus {
		message = fmt.Sprintf(constanta.CancellationRequestRejectedEmailBody, employeeName, requestType, createdAt)
	}

	input.sendNotificationToEmployee(approvalRequest.ClientID, message)

	/*
		Send Email
	*/
	go input.sendEmail(input.toMailAddress(approvalRequest.Email), "", constanta.HRISSubject, message)
	return
}

func (input employeeService) isAnnualLeave(allowanceType string) bool {
	return strings.ToLower(allowanceType) == constanta.AnnualLeaveAllowanceTypeKeyword ||
			strings.ToLower(allowanceType) == constanta.CutiTahunanAllowanceTypeKeyword
}

func (input employeeService) updateReimbursementApprovalStatus(tx *sql.Tx, body *in.EmployeeApprovalRequest, approvalRequest *approvalRequestModel, contextModel *applicationModel.ContextModel, now time.Time) (auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	model := repository.EmployeeReimbursement{
		ID: 		   sql.NullInt64{Int64: approvalRequest.Id},
		Status:        sql.NullString{String: body.Status},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, now, dao.EmployeeReimbursementDAO.TableName, approvalRequest.Id, contextModel.LimitedByCreatedBy)...)
	if errModel = dao.EmployeeReimbursementDAO.UpdateStatusTx(tx, model); errModel.Error != nil {
		return
	}

	/*
		Send Notification
	*/
	input.addNotificationIntoAudit(&auditData[0], &in.EmployeeNotification{
		IsMobileNotification:        true,
		IsRequestingForApproval:     false,
		IsRequestingForCancellation: false,
		IsCancellation:              false,
		EmployeeId:                  approvalRequest.EmployeeId,
		RequestType:                 approvalRequest.RequestType,
		Status:                      body.Status,
		Date:                        []time.Time{approvalRequest.CreatedAt},
	})

	employeeName := approvalRequest.EmployeeName
	requestType := constanta.TypeMedicalEmail
	createdAt := input.timeToString(approvalRequest.CreatedAt)

	message := fmt.Sprintf(constanta.RequestApprovedEmailBody, employeeName, requestType, createdAt)

	if body.Status == constanta.RejectedRequestStatus {
		message = fmt.Sprintf(constanta.RequestRejectedEmailBody, employeeName, requestType, createdAt)
	}

	input.sendNotificationToEmployee(approvalRequest.ClientID, message)

	/*
		Send Email
	*/
	go input.sendEmail(input.toMailAddress(approvalRequest.Email), "", constanta.HRISSubject, message)

	return
}

func (input employeeService) updateReimbursementBenefit(tx *sql.Tx, approvalRequest *approvalRequestModel, benefit repository.EmployeeBenefitsModel, contextModel *applicationModel.ContextModel, now time.Time) (auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	lastValue := benefit.LastMedicalValue.Float64 - approvalRequest.MedicalValue
	currentValue := benefit.CurrentMedicalValue.Float64

	if lastValue < 0 {
		currentValue += lastValue
		lastValue = 0
	}

	benefitModel := repository.EmployeeBenefitsModel{
		EmployeeID:          sql.NullInt64{Int64: approvalRequest.EmployeeId},
		CurrentMedicalValue: sql.NullFloat64{Float64: currentValue},
		LastMedicalValue:    sql.NullFloat64{Float64: lastValue},
		UpdatedAt:           sql.NullTime{Time: now},
		UpdatedBy:           sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:       sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, now, dao.EmployeeBenefitsDAO.TableName, benefitModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	errModel = dao.EmployeeBenefitsDAO.UpdateCurrentMedicalValueAndLastMedicalValueByEmployeeIdTx(tx, benefitModel)
	return
}

func (input employeeService) updateAnnualLeaveApprovalStatus(tx *sql.Tx, body *in.EmployeeApprovalRequest, approvalRequest *approvalRequestModel, contextModel *applicationModel.ContextModel, now time.Time) (auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	model := repository.EmployeeLeaveModel{
		ID: 		   sql.NullInt64{Int64: approvalRequest.Id},
		Status:        sql.NullString{String: body.Status},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, now, dao.EmployeeLeaveDAO.TableName, approvalRequest.Id, contextModel.LimitedByCreatedBy)...)
	if errModel = dao.EmployeeLeaveDAO.UpdateStatusTx(tx, model); errModel.Error != nil {
		return
	}

	input.addNotificationIntoAudit(&auditData[0], &in.EmployeeNotification{
		IsMobileNotification:        true,
		IsRequestingForApproval:     false,
		IsRequestingForCancellation: false,
		IsCancellation:              false,
		EmployeeId:                  approvalRequest.EmployeeId,
		RequestType:                 approvalRequest.RequestType,
		Status:                      body.Status,
		Date:                        approvalRequest.DateList,
	})

	/*
		Send Notification
	*/
	employeeName := approvalRequest.EmployeeName
	requestType := input.getEmailRequestType(approvalRequest.RequestType)
	createdAt := input.timeToString(approvalRequest.CreatedAt)

	message := fmt.Sprintf(constanta.RequestApprovedEmailBody, employeeName, requestType, createdAt)

	if body.Status == constanta.RejectedRequestStatus {
		message = fmt.Sprintf(constanta.RequestRejectedEmailBody, employeeName, requestType, createdAt)
	}

	input.sendNotificationToEmployee(approvalRequest.ClientID, message)

	/*
		Send Email
	*/
	go input.sendEmail(input.toMailAddress(approvalRequest.Email), "", constanta.HRISSubject, message)
	return
}

func (input employeeService) updateAnnualLeaveBenefit(tx *sql.Tx, body *in.EmployeeApprovalRequest, approvalRequest *approvalRequestModel, benefit repository.EmployeeBenefitsModel, contextModel *applicationModel.ContextModel, now time.Time) (auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	if body.Status != constanta.ApprovedRequestStatus {
		return nil, errorModel.GenerateNonErrorModel()
	}

	if !input.isAnnualLeave(approvalRequest.AllowanceType) {
		return nil, errorModel.GenerateNonErrorModel()
	}

	currentValue, lastValue := input.decreaseCurrentAndLastAnnualLeave(benefit.CurrentAnnualLeave.Int64, benefit.LastAnnualLeave.Int64, approvalRequest.TotalLeave)

	benefitModel := repository.EmployeeBenefitsModel{
		EmployeeID:         sql.NullInt64{Int64: approvalRequest.EmployeeId},
		CurrentAnnualLeave: sql.NullInt64{Int64: currentValue},
		LastAnnualLeave:    sql.NullInt64{Int64: lastValue},
		UpdatedAt:          sql.NullTime{Time: now},
		UpdatedBy:          sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:      sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, now, dao.EmployeeBenefitsDAO.TableName, benefitModel.ID.Int64, contextModel.LimitedByCreatedBy)...)
	errModel = dao.EmployeeBenefitsDAO.UpdateCurrentAnnualLeaveAndLastAnnualLeaveByEmployeeIdTx(tx, benefitModel)
	if errModel.Error != nil {
		return
	}

	errModel = dao.EmployeeHistoryLeaveDAO.UpdateCurrentAnnualLeaveAndLastAnnualLeaveByEmployeeIdTx(tx, benefitModel)
	return
}

func (input employeeService) validateApprovalRequest(approvalRequest *approvalRequestModel, body *in.EmployeeApprovalRequest) errorModel.ErrorModel {
	funcName := "validateApprovalRequest"

	if body.UpdatedAt != approvalRequest.UpdatedAt {
		return errorModel.GenerateDataLockedError(input.FileName, funcName, "updated_at")
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) getApprovalRequestByTypeAndId(tx *sql.Tx, reqType string, id int64) (result *approvalRequestModel, errModel errorModel.ErrorModel) {
	funcName := "getApprovalRequestByTypeAndId"

	if input.isLeaveRequest(reqType) {
		var (
			leaveReq repository.EmployeeLeaveModel
			dateList []time.Time
		)

		leaveReq, errModel = dao.EmployeeLeaveDAO.GetByIdAndStatuses(tx, id, []string{constanta.PendingRequestStatus, constanta.PendingCancellationRequestStatus})
		if errModel.Error != nil {
			return
		}

		if !leaveReq.ID.Valid {
			return nil, errorModel.GenerateUnknownDataError(input.FileName, funcName, "approval_request")
		}

		_ = json.Unmarshal([]byte(leaveReq.Date.String), &dateList)

		return &approvalRequestModel{
			Id:                 leaveReq.ID.Int64,
			EmployeeId:         leaveReq.EmployeeId.Int64,
			TotalLeave:         leaveReq.Value.Int64,
			Status:             leaveReq.Status.String,
			AllowanceId:        leaveReq.AllowanceId.Int64,
			AllowanceType:      leaveReq.AllowanceType.String,
			RequestType:        leaveReq.Type.String,
			CancellationReason: leaveReq.CancellationReason.String,
			UpdatedAt:          leaveReq.UpdatedAt.Time,
			TableName:          dao.EmployeeLeaveDAO.TableName,
			DateList:           dateList,
			ClientID: 			leaveReq.ClientID.String,
			EmployeeName: 		fmt.Sprintf("%s %s", leaveReq.Firstname.String, leaveReq.Lastname.String),
			Email: 				leaveReq.Email.String,
			CreatedAt: 			leaveReq.CreatedAt.Time,
		}, errorModel.GenerateNonErrorModel()
	}

	reimbursementReq, errModel := dao.EmployeeReimbursementDAO.GetByIdForUpdate(tx, id, constanta.PendingRequestStatus)
	if errModel.Error != nil {
		return
	}

	if !reimbursementReq.ID.Valid {
		return nil, errorModel.GenerateUnknownDataError(input.FileName, funcName, "approval_request")
	}

	return &approvalRequestModel{
		Id:           reimbursementReq.ID.Int64,
		EmployeeId:   reimbursementReq.EmployeeId.Int64,
		MedicalValue: reimbursementReq.Value.Float64,
		Status:       reimbursementReq.Status.String,
		BenefitId:    reimbursementReq.BenefitId.Int64,
		UpdatedAt:    reimbursementReq.UpdatedAt.Time,
		TableName:    dao.EmployeeReimbursementDAO.TableName,
		EmployeeName: fmt.Sprintf("%s %s", reimbursementReq.Firstname.String, reimbursementReq.Lastname.String),
		Email: 		  reimbursementReq.Email.String,
		CreatedAt: 	  reimbursementReq.CreatedAt.Time,
		ClientID: 	  reimbursementReq.ClientId.String,
	}, errorModel.GenerateNonErrorModel()
}

func (input employeeService) decreaseCurrentAndLastAnnualLeave(currentAnnualLeave, lastAnnualLeave, totalLeave int64) (currentAnnualLeaveResult, lastAnnualLeaveResult int64) {
	lastAnnualLeaveResult = lastAnnualLeave - totalLeave
	currentAnnualLeaveResult = currentAnnualLeave

	if lastAnnualLeaveResult < 0 {
		currentAnnualLeaveResult += lastAnnualLeaveResult
		lastAnnualLeaveResult = 0
	}

	return
}

func (input employeeService) increaseCurrentAndLastAnnualLeave(currentAnnualLeave, lastAnnualLeave, totalLeave, max int64) (currentAnnualLeaveResult, lastAnnualLeaveResult int64) {
	currentAnnualLeaveResult = currentAnnualLeave + totalLeave
	lastAnnualLeaveResult = lastAnnualLeave

	if currentAnnualLeaveResult > max {
		lastAnnualLeaveResult += currentAnnualLeaveResult - max
		currentAnnualLeaveResult = max
	}

	return
}

func (input employeeService) isLeaveRequest(reqType string) bool {
	return reqType == constanta.LeaveType ||
		reqType == constanta.PermitType ||
		reqType == constanta.SickLeaveType
}

func (input employeeService) readApprovalBody(request *http.Request, contextModel *applicationModel.ContextModel) (body *in.EmployeeApprovalRequest, errModel errorModel.ErrorModel) {
	var (
		funcName   = "readApprovalBody"
		stringBody string
	)

	stringBody, errModel = input.ReadBody(request, contextModel)
	if errModel.Error != nil {
		return
	}

	body = &in.EmployeeApprovalRequest{}
	if err := json.Unmarshal([]byte(stringBody), body); err != nil {
		return nil, errorModel.GenerateInvalidRequestError(input.FileName, funcName, err)
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	body.Id = int64(id)
	return
}
