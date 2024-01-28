package EmployeeService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
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

func (input employeeService) InsertEmployeeLeave(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	funcName := "InsertEmployeeLeave"

	body, errModel := input.readEmployeeLeaveFormData(request)
	if errModel.Error != nil {
		return
	}

	if errModel = body.ValidateInsert(); errModel.Error != nil {
		return
	}

	_, errModel = input.InsertServiceWithAudit(funcName, body, contextModel, input.insertEmployeeLeave, nil)
	if errModel.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input employeeService) insertEmployeeLeave(tx *sql.Tx, body interface{}, contextModel *applicationModel.ContextModel, now time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	employeeLeaveBody, _ := body.(*in.EmployeeLeaveRequest)

	user, errModel := input.getUser(contextModel)
	if errModel.Error != nil {
		return
	}

	leaders, errModel := dao.EmployeeDAO.GetListByMemberId(serverconfig.ServerAttribute.DBConnection, user.EmployeeId.Int64)
	if errModel.Error != nil {
		return
	}

	employeeLeaveBody.EmployeeId = user.EmployeeId.Int64

	allowance, errModel := input.getAllowance(employeeLeaveBody.Type, employeeLeaveBody.AllowanceId, user)
	if errModel.Error != nil {
		return
	}

	/*
		Validate Max Leave
	*/
	if errModel = input.validateMaxLeave(employeeLeaveBody, allowance); errModel.Error != nil {
		return
	}

	employeeLeaveBody.AllowanceId = allowance.ID.Int64

	employeeLeaveBody.FileUploadId, errModel = input.uploadAttachmentFile(tx, employeeLeaveBody.Attachment, contextModel, now)
	if errModel.Error != nil {
		return
	}

	auditResult, errModel := input.insert(tx, employeeLeaveBody, contextModel, now)
	if errModel.Error != nil {
		return
	}

	/*
		Send Notification
	*/
	auditResult.Employee.Firstname = user.FirstName.String
	auditResult.Employee.Lastname = user.LastName.String

	input.addNotificationIntoAudit(&auditResult, &in.EmployeeNotification{
		IsMobileNotification:        true,
		IsRequestingForApproval:     true,
		IsRequestingForCancellation: false,
		IsCancellation:              false,
		EmployeeId:                  employeeLeaveBody.EmployeeId,
		RequestType:                 employeeLeaveBody.Type,
		Status:                      constanta.PendingRequestStatus,
		Date:                        employeeLeaveBody.DateList,
	})

	go input.sendLeaveApprovalRequestNotifications(leaders, user, employeeLeaveBody)

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName: sql.NullString{String: dao.FileUploadDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: employeeLeaveBody.FileUploadId},
	})

	dataAudit = append(dataAudit, auditResult)
	return
}

func (input employeeService) sendLeaveApprovalRequestNotifications(leaders []repository.EmployeeModel, user repository.UserModel, employeeLeaveBody *in.EmployeeLeaveRequest) {
	for _, leader := range leaders {
		/*
			Send Notification
		*/
		leaderName := fmt.Sprintf("%s %s", leader.FirstName.String, leader.LastName.String)
		employeeName := fmt.Sprintf("%s %s", user.FirstName.String, user.LastName.String)
		requestType := input.getEmailRequestType(employeeLeaveBody.Type)
		leaveDate := input.timeToString(employeeLeaveBody.DateList[0])

		if employeeLeaveBody.Type == constanta.LeaveType || employeeLeaveBody.Type == constanta.SickLeaveType {
			if len(employeeLeaveBody.DateList) > 1 {
				startDate := input.timeToString(employeeLeaveBody.DateList[0])
				endDate := input.timeToString(employeeLeaveBody.DateList[len(employeeLeaveBody.DateList) - 1])

				leaveDate = fmt.Sprintf("%s s/d %s", startDate, endDate)
			}
		}

		message := fmt.Sprintf(constanta.LeaveApprovalRequestEmailBody, leaderName, employeeName, requestType, leaveDate)
		input.sendNotificationToEmployee(leader.ClientId.String, message)

		/*
			Send Email
		*/
		go input.sendEmail(input.toMailAddress(leader.Email.String), "", constanta.RequestApprovalHRISSubject, message)
	}
}

func (input employeeService) validateMaxLeave(employeeLeaveBody *in.EmployeeLeaveRequest, allowance repository.Allowance) errorModel.ErrorModel {
	funcName := "validateMaxLeave"

	if employeeLeaveBody.Type != constanta.LeaveType || input.isAnnualLeave(allowance.AllowanceType.String) {
		return errorModel.GenerateNonErrorModel()
	}

	maxLeave, _ := strconv.Atoi(allowance.Value.String)
	totalLeave := len(employeeLeaveBody.DateList)

	if totalLeave > maxLeave {
		return errorModel.GenerateAmountOfLeaveExceedMaxLeaveError(input.FileName, funcName, maxLeave)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) getAllowance(leaveReqType string, allowanceId int64, user repository.UserModel) (allowance repository.Allowance, errModel errorModel.ErrorModel) {
	funcName := "getAllowance"

	if leaveReqType == constanta.LeaveType {
		allowance, errModel = dao.AllowanceDAO.GetAllowanceLeaveByIdAndEmployeeLevelIdAndEmployeeGradeId(serverconfig.ServerAttribute.DBConnection, repository.Allowance{
			ID:              sql.NullInt64{Int64: allowanceId},
			EmployeeLevelId: sql.NullInt64{Int64: user.EmployeeLevelId.Int64},
			EmployeeGradeId: sql.NullInt64{Int64: user.EmployeeGradeId.Int64},
		})
		if errModel.Error != nil {
			return
		}

		if !allowance.ID.Valid {
			errModel = errorModel.GenerateUnknownDataError(input.FileName, funcName, "allowance_id")
		}
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) validateAllowanceId(allowanceId int64, allowanceType string) (errModel errorModel.ErrorModel) {
	funcName := "validateAllowanceId"

	allowance, errModel := dao.AllowanceDAO.GetByIdAndAllowanceType(serverconfig.ServerAttribute.DBConnection, allowanceId, allowanceType)
	if errModel.Error != nil {
		return
	}

	if !allowance.ID.Valid {
		errModel = errorModel.GenerateUnknownDataError(input.FileName, funcName, "allowance_id")
	}
	return
}

func (input employeeService) getAllowanceType(leaveReqType string) string {
	var allowanceType string

	switch leaveReqType {
	case constanta.LeaveType:
		allowanceType = constanta.LeaveAllowanceType
		break
	case constanta.PermitType:
		allowanceType = constanta.PermitAllowanceType
		break
	case constanta.SickLeaveType:
		allowanceType = constanta.SickLeaveAllowanceType
		break
	}

	return allowanceType
}

func (input employeeService) getTotalLeave(dateList []time.Time, leaveReqType string) int {
	if leaveReqType == constanta.PermitType {
		return 1
	}

	return len(dateList)
}

func (input employeeService) getUser(contextModel *applicationModel.ContextModel) (user repository.UserModel, errModel errorModel.ErrorModel) {
	funcName := "getUser"

	user, errModel = dao.UserDAO.GetById(serverconfig.ServerAttribute.DBConnection, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	if !user.EmployeeId.Valid {
		errModel = errorModel.GenerateUnknownDataError(input.FileName, funcName, "employee")
	}

	return
}

func (input employeeService) uploadAttachmentFile(tx *sql.Tx, attachment *in.EmployeeLeaveAttachment, contextModel *applicationModel.ContextModel, now time.Time) (fileUploadId int64, errModel errorModel.ErrorModel) {
	if attachment == nil {
		return 0, errorModel.GenerateNonErrorModel()
	}

	fileBytes, errModel := service.GetFileBytes(attachment.File)
	if errModel.Error != nil {
		return
	}

	var (
		container = constanta.ContainerEmployeeLeave + service.GetAzureDateContainer()
	)

	files := []in.MultipartFileDTO{
		{
			FileContent: fileBytes,
			Filename:    attachment.FileHeader.Filename,
			Size:        attachment.FileHeader.Size,
			Host:        config.ApplicationConfiguration.GetCDN().Host,
			Path:        config.ApplicationConfiguration.GetCDN().RootPath,
			FileID:      contextModel.AuthAccessTokenModel.ResourceUserID,
		},
	}

	errModel = service.UploadFileToLocalCDN(container, &files, contextModel.AuthAccessTokenModel.ResourceUserID)
	if errModel.Error != nil {
		return
	}

	go service.UploadFileToAzure(&files)

	fileUploadModel := repository.FileUpload{
		FileName:      sql.NullString{String: files[0].Filename},
		FileSize:      sql.NullInt64{Int64: attachment.FileHeader.Size},
		Konektor:      sql.NullString{String: dao.EmployeeLeaveDAO.TableName},
		Host:          sql.NullString{String: files[0].Host},
		Path:          sql.NullString{String: files[0].Path},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
	}

	return dao.FileUploadDAO.InsertFileUploadInfoForBacklog(tx, fileUploadModel)
}

func (input employeeService) insert(tx *sql.Tx, body *in.EmployeeLeaveRequest, contextModel *applicationModel.ContextModel, now time.Time) (repository.AuditSystemModel, errorModel.ErrorModel) {
	totalLeave := int64(input.getTotalLeave(body.DateList, body.Type))
	dateList, _ := json.Marshal(body.DateList)

	authAccessToken := contextModel.AuthAccessTokenModel
	model := repository.EmployeeLeaveModel{
		Type: 		   sql.NullString{String: body.Type},
		Name:          sql.NullString{String: body.Name},
		AllowanceId:   sql.NullInt64{Int64: body.AllowanceId},
		Description:   sql.NullString{String: body.Description},
		Date:		   sql.NullString{String: string(dateList)},
		Value:         sql.NullInt64{Int64: totalLeave},
		Status:        sql.NullString{String: constanta.PendingRequestStatus},
		FileUploadId:  sql.NullInt64{Int64: body.FileUploadId},
		EmployeeId:    sql.NullInt64{Int64: body.EmployeeId},
		CreatedAt:     sql.NullTime{Time: now},
		CreatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}

	id, errModel := dao.EmployeeLeaveDAO.InsertTx(tx, model)
	if errModel.Error != nil {
		return repository.AuditSystemModel{}, errModel
	}

	audit := repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.EmployeeLeaveDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: id},
	}

	return audit, errModel
}

func (input employeeService) readEmployeeLeaveFormData(request *http.Request) (*in.EmployeeLeaveRequest, errorModel.ErrorModel) {
	funcName := "readFormData"

	err := request.ParseMultipartForm(constanta.SizeMaximumRequireFileImport)
	if err != nil {
		return nil, errorModel.GenerateInvalidRequestError(input.FileName, funcName, err)
	}

	body := &in.EmployeeLeaveRequest{}

	content := request.FormValue("content")
	if err = json.Unmarshal([]byte(content), body); err != nil {
		return nil, errorModel.GenerateInvalidRequestError(input.FileName, funcName, err)
	}

	file, fileHeader, _ := request.FormFile("file")
	if file != nil {
		body.Attachment = &in.EmployeeLeaveAttachment{
			File:       file,
			FileHeader: fileHeader,
		}
	}

	return body, errorModel.GenerateNonErrorModel()
}
