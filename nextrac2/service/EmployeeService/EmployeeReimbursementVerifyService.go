package EmployeeService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

type employeeReimbursementVerifyService struct {
	FileName string
	service.AbstractService
}

var EmployeeReimbursementVerifyService = employeeReimbursementVerifyService{FileName: "EmployeeReimbursementVerifyService.go"}

func (input employeeReimbursementVerifyService) VerifyReimbursement(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	gradeBody, err := input.readParamAndBody(request, context)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit("VerifyReimbursement", gradeBody, context, input.doUpdate, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", context)

	return
}

func (input employeeReimbursementVerifyService) doUpdate(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doUpdate"
	authAccessToken := context.AuthAccessTokenModel

	reimBody := body.(in.EmployeeReimbursementVerifyRequest)
	reimRepository := input.getReimRepository(reimBody, authAccessToken, now)
	if err.Error != nil {
		return
	}

	reimOnDB, err := dao.EmployeeReimbursementDAO.GetDetailReimbursementForVerify(tx, reimRepository.ID.Int64)
	if err.Error != nil {
		return
	}

	benefitOnDB, err := dao.EmployeeBenefitsDAO.GetDetailMedicalValueForVerify(serverconfig.ServerAttribute.DBConnection, reimOnDB.EmployeeId.Int64)
	if err.Error != nil {
		return
	}

	if reimOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "reimbursement")
		return
	}

	if benefitOnDB.ID.Int64 < 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employee")
		return
	}

	if reimRepository.ApprovedValue.Float64 > reimOnDB.Value.Float64 {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, "aproved_value melebihi value", "aproved_value", "")
	    return
	}

	lastValue := benefitOnDB.LastMedicalValue.Float64 - reimRepository.ApprovedValue.Float64
	currentValue := benefitOnDB.CurrentMedicalValue.Float64

	if lastValue < 0 {
		currentValue += lastValue
		lastValue = 0
	}

	if currentValue < 0 {
		err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, "sisa saldo kurang", "sisa saldo", "")
		return
	}

	err = input.validation(reimOnDB, reimBody)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *context, now, dao.EmployeeReimbursementDAO.TableName, reimOnDB.ID.Int64, context.LimitedByCreatedBy)...)

	err = dao.EmployeeReimbursementDAO.UpdateVerifiedStatus(tx, reimRepository)
	if err.Error != nil {
		return
	}

	err = dao.EmployeeBenefitsDAO.UpdateCurrentMedicalValueAndLastMedicalValueByEmployeeIdTx(tx, repository.EmployeeBenefitsModel{
		EmployeeID:          sql.NullInt64{Int64: reimOnDB.EmployeeId.Int64},
		CurrentMedicalValue: sql.NullFloat64{Float64: currentValue},
		LastMedicalValue:    sql.NullFloat64{Float64: lastValue},
		UpdatedAt:           sql.NullTime{Time: now},
		UpdatedBy:           sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedClient:       sql.NullString{String: authAccessToken.ClientID},
	})
	if err.Error != nil {
		return
	}

	/*
		Send Notification
	*/
	employeeName := fmt.Sprintf("%s %s", reimOnDB.Firstname.String, reimOnDB.Lastname.String)
	requestType := constanta.TypeMedicalEmail
	createdAt := input.timeToString(reimOnDB.CreatedAt.Time)

	message := fmt.Sprintf(constanta.RequestVerifiedEmailBody, employeeName, requestType, createdAt)

	if len(auditData) > 0 {
		input.addNotificationIntoAudit(&auditData[0], &in.EmployeeNotification{
			IsMobileNotification:        true,
			IsVerified: 				 true,
			EmployeeId:                  reimOnDB.EmployeeId.Int64,
			RequestType:                 constanta.ReimbursementType,
			Date:                        []time.Time{reimOnDB.CreatedAt.Time},
		})
	}

	input.sendNotificationToEmployee(reimOnDB.ClientId.String, message)

	/*
		Send Email
	*/
	go input.sendEmail(input.toMailAddress(reimOnDB.Email.String), "", constanta.HRISSubject, message)

	return
}

func (input employeeReimbursementVerifyService) readParamAndBody(request *http.Request, contextModel *applicationModel.ContextModel) (reimBody in.EmployeeReimbursementVerifyRequest, err errorModel.ErrorModel) {
	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	reimBody, bodySize, err := getReimbursementBody(request, input.FileName)
	reimBody.Id = id
	contextModel.LoggerModel.ByteIn = bodySize
	return
}

func (input employeeReimbursementVerifyService) getReimRepository(reim in.EmployeeReimbursementVerifyRequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.EmployeeReimbursement {
	status := ""
	if reim.ApprovedValue == 0{
		status = constanta.UnverifiedReimbursementVerification
	}else{
		status = constanta.VerifiedReimbursementVerification
	}
	return repository.EmployeeReimbursement{
		ID:            sql.NullInt64{Int64: reim.Id},
		ApprovedValue: sql.NullFloat64{Float64: reim.ApprovedValue},
		Note:          sql.NullString{String: reim.Note},
		EmployeeId:    sql.NullInt64{Int64:  reim.EmployeeId},
		VerifiedStatus: sql.NullString{String: status},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
	}
}

func (input employeeReimbursementVerifyService) validation(reimOnDB repository.EmployeeReimbursement, reimBody in.EmployeeReimbursementVerifyRequest) (err errorModel.ErrorModel) {
	err = reimBody.ValidateVerifyReimbursement(true)
	if err.Error != nil {
		return
	}
	err = service.OptimisticLock(reimOnDB.UpdatedAt.Time, reimBody.UpdatedAt, input.FileName, "reimbursement")
	return
}

func getReimbursementBody(request *http.Request, fileName string) (reimBody in.EmployeeReimbursementVerifyRequest, bodySize int, err errorModel.ErrorModel) {
	funcName := "getReimbursementBody"
	jsonString, bodySize, readError := util.ReadBody(request)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	readError = json.Unmarshal([]byte(jsonString), &reimBody)

	if readError != nil {
		err = errorModel.GenerateInvalidRequestError(fileName, funcName, readError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeReimbursementVerifyService) addNotificationIntoAudit(auditData *repository.AuditSystemModel, notification *in.EmployeeNotification) (title, body string) {
	if auditData == nil || notification == nil {
		return
	}

	bytes, _ := json.Marshal(notification)
	auditData.Description.String = string(bytes)

	return input.generateNotificationMessage(auditData.Employee.Firstname, auditData.Employee.Lastname, *notification)
}

func (input employeeReimbursementVerifyService) generateNotificationMessage(firstName, lastName string, notification in.EmployeeNotification) (title, body string) {
	requestType, date := input.getRequestTypeAndDate(notification)

	if notification.IsRequestingForApproval {
		title = constanta.EmployeeRequestPendingMessageTitle
		body = fmt.Sprintf(constanta.EmployeeRequestPendingMessageBody,
			fmt.Sprintf("%s %s", firstName, lastName),
			requestType,
			date)

		return
	}

	if notification.IsRequestingForCancellation {
		title = constanta.EmployeeCancelRequestPendingMessageTitle
		body = fmt.Sprintf(constanta.EmployeeCancelRequestPendingMessageBody,
			fmt.Sprintf("%s %s", firstName, lastName),
			requestType,
			date)

		return
	}

	if notification.IsCancellation {
		if notification.Status == constanta.ApprovedRequestStatus {
			title = fmt.Sprintf(constanta.EmployeeCancelRequestApprovedMessageTitle, requestType)
			body = fmt.Sprintf(constanta.EmployeeCancelRequestApprovedMessageBody,
				requestType,
				date)

			return
		}

		if notification.Status == constanta.RejectedRequestStatus {
			title = fmt.Sprintf(constanta.EmployeeCancelRequestRejectedMessageTitle, requestType)
			body = fmt.Sprintf(constanta.EmployeeCancelRequestRejectedMessageBody,
				requestType,
				date)

			return
		}
	}

	if notification.IsVerified {
		title = fmt.Sprintf(constanta.EmployeeRequestVerifiedMessageTitle, requestType)
		body = fmt.Sprintf(constanta.EmployeeRequestVerifiedMessageBody,
			requestType,
			date)

		return
	}

	if notification.Status == constanta.ApprovedRequestStatus {
		title = fmt.Sprintf(constanta.EmployeeRequestApprovedMessageTitle, requestType)
		body = fmt.Sprintf(constanta.EmployeeRequestApprovedMessageBody,
			requestType,
			date)

		return
	}

	if notification.Status == constanta.RejectedRequestStatus {
		title = fmt.Sprintf(constanta.EmployeeRequestRejectedMessageTitle, requestType)
		body = fmt.Sprintf(constanta.EmployeeRequestRejectedMessageBody,
			requestType,
			date)

		return
	}

	return "", ""
}

func (input employeeReimbursementVerifyService) getRequestTypeAndDate(notification in.EmployeeNotification) (requestType, date string) {
	startDate := notification.Date[0].Format("02 Jan 2006")
	endDate := ""

	if len(notification.Date) > 1 {
		endDate = fmt.Sprintf(" - %s", notification.Date[len(notification.Date) - 1].Format("02 Jan 2006"))
	}

	switch notification.RequestType {
	case constanta.LeaveType:
		requestType = constanta.LeaveTypeNotification
		date = startDate + endDate

		return
	case constanta.PermitType :
		requestType = constanta.PermitTypeNotification
		date = startDate

		return
	case constanta.SickLeaveType :
		requestType = constanta.SickLeaveTypeNotification
		date = startDate + endDate

		return
	case constanta.ReimbursementType :
		requestType = constanta.MedicalTypeNotification
		date = startDate

		return
	}

	return "", ""
}

func (input employeeReimbursementVerifyService) sendNotificationToEmployee(destinationId, message string) {
	groChatWSHandler := serverconfig.ServerAttribute.GroChatWSHandler

	_ = groChatWSHandler.SendNotification(destinationId, []string{message})
	return
}

func (input employeeReimbursementVerifyService) sendEmail(to []mail.Address, senderName string, subject string, body string) {
	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())

	reqEmail := util2.NewRequestMail(subject)
	reqEmail.SenderName = senderName
	reqEmail.To = to
	reqEmail.Body = body
	reqEmail.ContentType = "text/plain"

	if err := reqEmail.SendEmail(); err != nil {
		logModel.Status = 500
		logModel.Message = fmt.Sprintf("Failed to send email : %s", err.Error())

		util.LogError(logModel.ToLoggerObject())
		return
	}

	logModel.Status = 200
	logModel.Message = "Email sent"

	util.LogInfo(logModel.ToLoggerObject())
}

func (input employeeReimbursementVerifyService) toMailAddress(address string) []mail.Address {
	return []mail.Address{
		{
			Address: address,
		},
	}
}

func (input employeeReimbursementVerifyService) timeToString(t time.Time) string {
	months := []string{
		"Januari", "Februari", "Maret",
		"April",   "Mei", 	   "Juni",
		"Juli",    "Agustus",  "September",
		"Oktober", "November", "Desember",
	}

	month := months[t.Month() - 1]
	return fmt.Sprintf(t.Format("02 %s 2006"), month)
}