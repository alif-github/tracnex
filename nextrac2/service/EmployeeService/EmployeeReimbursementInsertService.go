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
	"strings"
	"time"
)

func (input employeeService) InsertEmployeeReimbursement(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	funcName := "InsertEmployeeReimbursement"

	body, errModel := input.readReimbursementFormData(request)
	if errModel.Error != nil {
		return
	}

	param, errModel := service.ReadParameterByPermissionAndName(constanta.NexTracParameterPermission, constanta.ParameterExpiredMedicalClaim, 0)
	if errModel.Error != nil {
		return
	}

	expiredMedicalClaim, _ := strconv.Atoi(param)

	if errModel = body.ValidateInsert(expiredMedicalClaim); errModel.Error != nil {
		return
	}

	_, errModel = input.InsertServiceWithAudit(funcName, body, contextModel, input.insertReimbursement, nil)
	if errModel.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input employeeService) insertReimbursement(tx *sql.Tx, body interface{}, contextModel *applicationModel.ContextModel, now time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	employeeReimbursementBody, _ := body.(*in.EmployeeReimbursementRequest)

	user, errModel := input.getUser(contextModel)
	if errModel.Error != nil {
		return
	}

	leaders, errModel := dao.EmployeeDAO.GetListByMemberId(serverconfig.ServerAttribute.DBConnection, user.EmployeeId.Int64)
	if errModel.Error != nil {
		return
	}

	employeeReimbursementBody.EmployeeId = user.EmployeeId.Int64

	errModel = input.validateReimbursementBody(employeeReimbursementBody, user)
	if errModel.Error != nil {
		return
	}

	employeeReimbursementBody.FileUploadId, errModel = input.uploadReimbursementAttachmentFile(tx, employeeReimbursementBody.Attachment, contextModel, now)
	if errModel.Error != nil {
		return
	}

	auditResult, errModel := input.doInsertReimbursement(tx, employeeReimbursementBody, contextModel, now)
	if errModel.Error != nil {
		return
	}

	/*
		Send Notification
	*/
	auditResult.Employee.Firstname = user.FirstName.String
	auditResult.Employee.Lastname = user.LastName.String
	date, _ := time.Parse(constanta.DefaultTimeFormat, now.Format(constanta.DefaultTimeFormat))

	input.addNotificationIntoAudit(&auditResult, &in.EmployeeNotification{
		IsMobileNotification:        true,
		IsRequestingForApproval:     true,
		IsRequestingForCancellation: false,
		IsCancellation:              false,
		EmployeeId:                  employeeReimbursementBody.EmployeeId,
		RequestType:                 constanta.ReimbursementType,
		Status:                      constanta.PendingRequestStatus,
		Date:                        []time.Time{date},
	})

	go input.sendReimbursementApprovalRequestNotifications(leaders, user, constanta.ReimbursementType, now)

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.FileUploadDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: employeeReimbursementBody.FileUploadId},
	})

	dataAudit = append(dataAudit, auditResult)
	return
}

func (input employeeService) sendReimbursementApprovalRequestNotifications(leaders []repository.EmployeeModel, user repository.UserModel, requestType string, now time.Time) {
	for _, leader := range leaders {
		/*
			Send Notification
		*/
		leaderName := fmt.Sprintf("%s %s", leader.FirstName.String, leader.LastName.String)
		employeeName := fmt.Sprintf("%s %s", user.FirstName.String, user.LastName.String)
		reqType := input.getEmailRequestType(requestType)
		createdAt := input.timeToString(now)

		message := fmt.Sprintf(constanta.ReimbursementApprovalRequestEmailBody, leaderName, employeeName, reqType, createdAt)
		input.sendNotificationToEmployee(leader.ClientId.String, message)

		/*
			Send Email
		*/
		go input.sendEmail(input.toMailAddress(leader.Email.String), "", constanta.HRISSubject, message)
	}
}

func (input employeeService) validateReimbursementBody(body *in.EmployeeReimbursementRequest, user repository.UserModel) errorModel.ErrorModel {
	if errModel := input.validateBenefitId(body.BenefitId, user); errModel.Error != nil {
		return errModel
	}

	if errModel := input.validateReimbursementAttachment(body.Attachment); errModel.Error != nil {
		return errModel
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) validateBenefitId(benefitId int64, user repository.UserModel) (errModel errorModel.ErrorModel) {
	funcName := "validateBenefitId"

	benefit, errModel := dao.BenefitDAO.GetMedicalBenefitByIdAndEmployeeLevelIdAndEmployeeGradeId(serverconfig.ServerAttribute.DBConnection, repository.Benefit{
		ID:              sql.NullInt64{Int64: benefitId},
		EmployeeLevelId: sql.NullInt64{Int64: user.EmployeeLevelId.Int64},
		EmployeeGradeId: sql.NullInt64{Int64: user.EmployeeGradeId.Int64},
	})
	if errModel.Error != nil {
		return
	}

	if !benefit.ID.Valid {
		errModel = errorModel.GenerateUnknownDataError(input.FileName, funcName, "benefit_id")
	}

	return
}

func (input employeeService) validateReimbursementAttachment(attachment *in.EmployeeReimbursementAttachment) errorModel.ErrorModel {
	funcName := "validateReimbursementAttachment"

	if attachment == nil {
		return errorModel.GenerateEmptyFieldError(input.FileName, funcName, "file")
	}

	var maxFileSize int64 = 5 * 1e6

	if attachment.FileHeader.Size > maxFileSize {
		return errorModel.GenerateFieldFormatWithRuleError(input.FileName, funcName, "REIMBURSEMENT_MAX_FILE_SIZE", "file", "")
	}

	if errModel := input.validateReimbursementFileExtension(attachment.FileHeader.Filename); errModel.Error != nil {
		return errModel
	}

	return errorModel.GenerateNonErrorModel()
}

func (input employeeService) validateReimbursementFileExtension(fileName string) errorModel.ErrorModel {
	funcName := "validateReimbursementFileExtension"

	allowedExtensions := []string{"pdf", "xls", "xlsx", "jpg", "jpeg", "png"}

	separatedFileParts := strings.Split(fileName, ".")
	fileExtension := separatedFileParts[len(separatedFileParts) - 1]

	for _, extension := range allowedExtensions {
		if fileExtension == extension {
			return errorModel.GenerateNonErrorModel()
		}
	}

	return errorModel.GenerateInvalidFileExtensionError(fileName, funcName, "file", allowedExtensions)
}

func (input employeeService) uploadReimbursementAttachmentFile(tx *sql.Tx, attachment *in.EmployeeReimbursementAttachment, contextModel *applicationModel.ContextModel, now time.Time) (fileUploadId int64, errModel errorModel.ErrorModel) {
	if attachment == nil {
		return 0, errorModel.GenerateNonErrorModel()
	}

	fileBytes, errModel := service.GetFileBytes(attachment.File)
	if errModel.Error != nil {
		return
	}

	var (
		container = constanta.ContainerEmployeeReimbursement + service.GetAzureDateContainer()
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
		Konektor:      sql.NullString{String: dao.EmployeeReimbursementDAO.TableName},
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

func (input employeeService) doInsertReimbursement(tx *sql.Tx, body *in.EmployeeReimbursementRequest, contextModel *applicationModel.ContextModel, now time.Time) (repository.AuditSystemModel, errorModel.ErrorModel) {
	authAccessToken := contextModel.AuthAccessTokenModel
	model := repository.EmployeeReimbursement{
		Name:           sql.NullString{String: body.Name},
		ReceiptNo:      sql.NullString{String: body.Receipt},
		BenefitId:      sql.NullInt64{Int64: body.BenefitId},
		Description:    sql.NullString{String: body.Description},
		Date:           sql.NullTime{Time: body.ReceiptDateTime},
		Value:          sql.NullFloat64{Float64: body.Value},
		Status:         sql.NullString{String: constanta.PendingReimbursementRequestType},
		VerifiedStatus: sql.NullString{String: constanta.PendingReimbursementVerification},
		FileUploadId:   sql.NullInt64{Int64: body.FileUploadId},
		EmployeeId:     sql.NullInt64{Int64: body.EmployeeId},
		CreatedBy:      sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedAt:      sql.NullTime{Time: now},
		CreatedClient:  sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:      sql.NullTime{Time: now},
		UpdatedBy:      sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedClient:  sql.NullString{String: authAccessToken.ClientID},
	}

	id, errModel := dao.EmployeeReimbursementDAO.InsertTx(tx, model)
	if errModel.Error != nil {
		return repository.AuditSystemModel{}, errModel
	}

	audit := repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.EmployeeReimbursementDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: id},
	}

	return audit, errModel
}

func (input employeeService) readReimbursementFormData(request *http.Request) (*in.EmployeeReimbursementRequest, errorModel.ErrorModel) {
	funcName := "readReimbursementFormData"

	err := request.ParseMultipartForm(constanta.SizeMaximumRequireFileImport)
	if err != nil {
		return nil, errorModel.GenerateInvalidRequestError(input.FileName, funcName, err)
	}

	body := &in.EmployeeReimbursementRequest{}

	content := request.FormValue("content")
	if err = json.Unmarshal([]byte(content), body); err != nil {
		return nil, errorModel.GenerateInvalidRequestError(input.FileName, funcName, err)
	}

	file, fileHeader, _ := request.FormFile("file")
	if file != nil {
		body.Attachment = &in.EmployeeReimbursementAttachment{
			File:       file,
			FileHeader: fileHeader,
		}
	}

	return body, errorModel.GenerateNonErrorModel()
}
