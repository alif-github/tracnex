package EmployeeService

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/mail"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"strconv"
	"time"
)

type employeeService struct {
	EmployeeDAO dao.EmployeeDAOInterface
	service.AbstractService
	service.GetListData
}

var EmployeeService = employeeService{}.New()

func (input employeeService) New() (output employeeService) {
	output.FileName = "EmployeeService.go"
	output.ServiceName = constanta.EmployeeConstanta
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"nik",
		"redmine_id",
		"name",
		"department",
		"updated_at",
	}
	output.ValidSearchBy = []string{
		"nik",
		"redmine_id",
		"name",
		"department",
		"department_id",
		"is_timesheet",
		"is_redmine_check",
	}
	output.EmployeeDAO = dao.EmployeeDAO
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.EmployeeDataScope] = applicationModel.MappingScopeDB{
		View:  "e.id",
		Count: "e.id",
	}
	return
}

func (input employeeService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.EmployeeRequest) errorModel.ErrorModel) (inputStruct in.EmployeeRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if stringBody != "" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func (input employeeService) readBodyWithFileAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.EmployeeRequest) errorModel.ErrorModel) (inputStruct in.EmployeeRequest, files []in.MultipartFileDTO, err errorModel.ErrorModel) {
	var (
		fileName  = input.FileName
		funcName  = "readBodyWithFileAndValidate"
		totalSize int64
		errs      error
		content   string
	)

	errs = request.ParseMultipartForm(32 << 20)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	content = request.FormValue("content")
	byteContent := []byte(content)
	if errs = json.Unmarshal(byteContent, &inputStruct); errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	files, totalSize, err = service.ReadFileWithMultipart(request, 1, input.imageValidation)
	if err.Error != nil {
		return
	}

	contextModel.LoggerModel.ByteIn = int(totalSize) + len(byteContent)
	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func (input employeeService) imageValidation(dto in.MultipartFileDTO) errorModel.ErrorModel {
	var (
		funcName   = "imageValidation"
		errs       error
		additional string
	)

	errs, additional = util2.IsFileImage(dto.FileContent, constanta.CompanyProfileMaximumPhotoSize)
	if errs != nil {
		return errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, errs.Error(), "Photo", additional)
	}
	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) validateDataScope(contextModel *applicationModel.ContextModel) (scope map[string]interface{}, err errorModel.ErrorModel) {
	scope, err = input.ValidateMultipleDataScope(contextModel, []string{
		constanta.EmployeeDataScope,
	})

	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) addNotificationIntoAudit(auditData *repository.AuditSystemModel, notification *in.EmployeeNotification) (title, body string) {
	if auditData == nil || notification == nil {
		return
	}

	bytes, _ := json.Marshal(notification)
	auditData.Description.String = string(bytes)

	return input.generateNotificationMessage(auditData.Employee.Firstname, auditData.Employee.Lastname, *notification)
}

func (input employeeService) generateNotificationMessage(firstName, lastName string, notification in.EmployeeNotification) (title, body string) {
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

func (input employeeService) getRequestTypeAndDate(notification in.EmployeeNotification) (requestType, date string) {
	startDate := notification.Date[0].Format("02 Jan 2006")
	endDate := ""

	if len(notification.Date) > 1 {
		endDate = fmt.Sprintf(" - %s", notification.Date[len(notification.Date)-1].Format("02 Jan 2006"))
	}

	switch notification.RequestType {
	case constanta.LeaveType:
		requestType = constanta.LeaveTypeNotification
		date = startDate + endDate

		return
	case constanta.PermitType:
		requestType = constanta.PermitTypeNotification
		date = startDate

		return
	case constanta.SickLeaveType:
		requestType = constanta.SickLeaveTypeNotification
		date = startDate + endDate

		return
	case constanta.ReimbursementType:
		requestType = constanta.MedicalTypeNotification
		date = startDate

		return
	}

	return "", ""
}

func (input employeeService) sendNotificationToLeaders(leaders []repository.EmployeeModel, messageTitle, messageBody string) errorModel.ErrorModel {
	groChatWSHandler := serverconfig.ServerAttribute.GroChatWSHandler

	for _, leader := range leaders {
		if leader.ClientId.String == "" {
			continue
		}

		message := fmt.Sprintf("%s\n\n%s", messageTitle, messageBody)
		_ = groChatWSHandler.SendNotification(leader.ClientId.String, []string{message})
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) sendNotificationToEmployee(destinationId, message string) {
	groChatWSHandler := serverconfig.ServerAttribute.GroChatWSHandler

	_ = groChatWSHandler.SendNotification(destinationId, []string{message})
	return
}

func (input employeeService) toMailAddress(address string) []mail.Address {
	return []mail.Address{
		{
			Address: address,
		},
	}
}

func (input employeeService) sendEmail(to []mail.Address, senderName string, subject string, body string) {
	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())

	reqEmail := util.NewRequestMail(subject)
	reqEmail.SenderName = senderName
	reqEmail.To = to
	reqEmail.Body = body
	reqEmail.ContentType = "text/plain"

	if err := reqEmail.SendEmail(); err != nil {
		logModel.Status = 500
		logModel.Message = fmt.Sprintf("Failed to send email : %s", err.Error())

		util2.LogError(logModel.ToLoggerObject())
		return
	}

	logModel.Status = 200
	logModel.Message = "Email sent"

	util2.LogInfo(logModel.ToLoggerObject())
}

func (input employeeService) getEmailRequestType(requestType string) string {
	switch requestType {
	case constanta.LeaveType:
		return constanta.TypeLeaveEmail
	case constanta.PermitType:
		return constanta.TypePermitEmail
	case constanta.SickLeaveType:
		return constanta.TypeSickLeaveEmail
	case constanta.ReimbursementType:
		return constanta.TypeMedicalEmail
	}

	return ""
}

func (input employeeService) timeToString(t time.Time) string {
	months := []string{
		"Januari", "Februari", "Maret",
		"April", "Mei", "Juni",
		"Juli", "Agustus", "September",
		"Oktober", "November", "Desember",
	}

	month := months[t.Month()-1]
	return fmt.Sprintf(t.Format("02 %s 2006"), month)
}
