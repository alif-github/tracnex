package service

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/mail"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/backgroundJobModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type AbstractService struct {
	ServiceName    string
	FileName       string
	Audit          bool
	MappingScopeDB map[string]applicationModel.MappingScopeDB
	ListScope      []string
}

func (input AbstractService) ReadBody(request *http.Request, contextModel *applicationModel.ContextModel) (string, errorModel.ErrorModel) {
	var stringBody string
	var errorS error

	funcName := "ReadBody"

	if request.Method != "GET" {
		stringBody, contextModel.LoggerModel.ByteIn, errorS = util.ReadBody(request)
		if errorS != nil {
			return "", errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		}

		if contextModel.IsSignatureCheck {
			digest := util.GenerateMessageDigest(stringBody)
			if !ValidateSignature(digest, contextModel.AuthAccessTokenModel.SignatureKey, request) {
				return "", errorModel.GenerateInvalidSignatureError(input.FileName, "ReadBody")
			}
		}
	}

	return stringBody, errorModel.GenerateNonErrorModel()
}

func (input AbstractService) InsertServiceWithAuditCustom(funcName string, inputStruct interface{}, contextModel *applicationModel.ContextModel, serve func(*sql.Tx, interface{}, *applicationModel.ContextModel, time.Time) (interface{}, []repository.AuditSystemModel, bool, errorModel.ErrorModel), additionalAfterCommit func(interface{}, applicationModel.ContextModel)) (output interface{}, err errorModel.ErrorModel) {
	return input.ServiceWithDataAuditGetByAuditServiceCustom(constanta.ActionAuditInsertConstanta, funcName, inputStruct, contextModel, serve, additionalAfterCommit)
}

func (input AbstractService) ServiceWithDataAuditGetByAuditServiceCustom(action int32, funcName string, inputStruct interface{}, contextModel *applicationModel.ContextModel, serve func(*sql.Tx, interface{}, *applicationModel.ContextModel, time.Time) (interface{}, []repository.AuditSystemModel, bool, errorModel.ErrorModel), additionalAfterCommit func(interface{}, applicationModel.ContextModel)) (output interface{}, err errorModel.ErrorModel) {
	var errs error
	var tx *sql.Tx
	var dataAudit []repository.AuditSystemModel
	var dataAuditSave []repository.AuditSystemModel
	var isServiceUpdate bool

	timeNow := time.Now()

	defer func() {
		if errs != nil || err.Error != nil {
			errs = tx.Rollback()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
			}
		} else {
			if isServiceUpdate {
				for i := 0; i < len(dataAudit); i++ {
					if dataAudit[i].Action.Int32 == constanta.ActionAuditInsertConstanta {
						dataAuditSave = append(dataAuditSave, GetAuditData(tx, action, *contextModel, timeNow, dataAudit[i].TableName.String, dataAudit[i].PrimaryKey.Int64, contextModel.LimitedByCreatedBy)...)
					} else {
						dataAuditSave = append(dataAuditSave, dataAudit[i])
					}
				}

				errs = tx.Commit()
				if errs != nil {
					err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
				} else {
					go additionalAfterCommit(output, *contextModel)
					input.ServiceWithAudit(dataAuditSave, *contextModel, err)
				}
			} else {
				for i := 0; i < len(dataAudit); i++ {
					dataAuditSave = append(dataAuditSave, GetAuditData(tx, action, *contextModel, timeNow, dataAudit[i].TableName.String, dataAudit[i].PrimaryKey.Int64, contextModel.LimitedByCreatedBy)...)
				}

				errs = tx.Commit()
				if errs != nil {
					err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
				} else {
					if output != nil {
						additionalAfterCommit(output, *contextModel)
					}
					input.ServiceWithAudit(dataAuditSave, *contextModel, err)
				}
			}
		}
	}()

	tx, errs = serverconfig.ServerAttribute.DBConnection.Begin()
	if errs != nil {
		return
	}

	output, dataAudit, isServiceUpdate, err = serve(tx, inputStruct, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input AbstractService) InsertServiceWithAudit(funcName string, inputStruct interface{}, contextModel *applicationModel.ContextModel, serve func(*sql.Tx, interface{}, *applicationModel.ContextModel, time.Time) (interface{}, []repository.AuditSystemModel, errorModel.ErrorModel), additionalAfterCommit func(interface{}, applicationModel.ContextModel)) (output interface{}, err errorModel.ErrorModel) {
	return input.ServiceWithDataAuditGetByAuditService(constanta.ActionAuditInsertConstanta, funcName, inputStruct, contextModel, serve, additionalAfterCommit)
}

func (input AbstractService) ServiceWithDataAuditGetByAuditService(action int32, funcName string, inputStruct interface{}, contextModel *applicationModel.ContextModel, serve func(*sql.Tx, interface{}, *applicationModel.ContextModel, time.Time) (interface{}, []repository.AuditSystemModel, errorModel.ErrorModel), additionalAfterCommit func(interface{}, applicationModel.ContextModel)) (output interface{}, err errorModel.ErrorModel) {
	var errs error
	var tx *sql.Tx
	var dataAudit []repository.AuditSystemModel
	var dataAuditSave []repository.AuditSystemModel

	timeNow := time.Now()

	defer func() {
		if errs != nil || err.Error != nil {
			errs = tx.Rollback()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
			}
		} else {
			for i := 0; i < len(dataAudit); i++ {
				auditList := GetAuditData(tx, action, *contextModel, timeNow, dataAudit[i].TableName.String, dataAudit[i].PrimaryKey.Int64, contextModel.LimitedByCreatedBy)
				if auditList != nil {
					auditList[0].Description.String = dataAudit[i].Description.String
				}

				dataAuditSave = append(dataAuditSave, auditList...)
			}

			errs = tx.Commit()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
			} else {
				if output != nil {
					additionalAfterCommit(output, *contextModel)
				}
				input.ServiceWithAudit(dataAuditSave, *contextModel, err)
			}
		}
	}()

	tx, errs = serverconfig.ServerAttribute.DBConnection.Begin()
	if errs != nil {
		return
	}

	output, dataAudit, err = serve(tx, inputStruct, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input AbstractService) ServiceWithDataAuditPreparedByService(funcName string, inputStruct interface{}, contextModel *applicationModel.ContextModel, serve func(*sql.Tx, interface{}, *applicationModel.ContextModel, time.Time) (interface{}, []repository.AuditSystemModel, errorModel.ErrorModel), additionalAfterCommit func(interface{}, applicationModel.ContextModel)) (output interface{}, err errorModel.ErrorModel) {
	var (
		errs          error
		tx            *sql.Tx
		dataAuditSave []repository.AuditSystemModel
		timeNow       = time.Now()
	)

	defer func() {
		if errs != nil || err.Error != nil {
			errs = tx.Rollback()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
			}
		} else {
			errs = tx.Commit()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
			} else {
				go additionalAfterCommit(output, *contextModel)
				input.ServiceWithAudit(dataAuditSave, *contextModel, err)
			}
		}
	}()

	tx, errs = serverconfig.ServerAttribute.DBConnection.Begin()
	if errs != nil {
		return
	}

	output, dataAuditSave, err = serve(tx, inputStruct, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input AbstractService) ServiceWithAudit(data []repository.AuditSystemModel, contextModel applicationModel.ContextModel, err errorModel.ErrorModel) {
	if err.Error != nil {
		return
	}
	if config.ApplicationConfiguration.GetAudit().IsActive {
		go input.saveAuditDataToDatabase(data, contextModel)
	} else {
		err = errorModel.GenerateInactiveAuditSystem(input.FileName, "ServiceWithAudit")
		input.LogError(err, contextModel)
	}
}

func (input AbstractService) saveAuditDataToDatabase(data []repository.AuditSystemModel, contextModel applicationModel.ContextModel) {
	var err errorModel.ErrorModel
	for i := 0; i < len(data); i++ {
		err = validateAuditSystemModel(data[i])
		if err.Error != nil {
			input.LogErrorAndAlert(err, contextModel)
			continue
		}
		err = dao.AuditSystemDAO.InsertAuditSystem(serverconfig.ServerAttribute.DBConnection, data[i])
		if err.Error != nil {
			input.LogError(err, contextModel)
		}
	}
}

func validateAuditSystemModel(data repository.AuditSystemModel) errorModel.ErrorModel {
	fileName := "AbstractService.go"
	funcName := "validateAuditSystemModel"
	if data.TableName.String == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.TableName)
	}
	if data.PrimaryKey.Int64 < 0 {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.PrimaryKey)
	}
	if data.Data.String == "" {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.Data)

	}
	if data.CreatedAt.Time.IsZero() {
		return errorModel.GenerateEmptyFieldError(fileName, funcName, constanta.CreatedAt)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input AbstractService) LogErrorAndAlert(err errorModel.ErrorModel, contextModel applicationModel.ContextModel) {
	go input.LogError(err, contextModel)
	go alertServer(err, contextModel)
}

func (input AbstractService) LogError(err errorModel.ErrorModel, contextModel applicationModel.ContextModel) {
	contextModel.LoggerModel.Status = err.Code
	if err.CausedBy != nil {
		err.Error = err.CausedBy
	}
	contextModel.LoggerModel.Message = util2.GenerateI18NErrorMessage(err, constanta.DefaultApplicationsLanguage)
	util.LogError(contextModel.LoggerModel.ToLoggerObject())
}

func alertServer(_ errorModel.ErrorModel, _ applicationModel.ContextModel) {
	//todo Alert Server
}

func GetAuditData(tx *sql.Tx, action int32, contextModel applicationModel.ContextModel, timeNow time.Time, tableName string, id int64, createdBy int64) []repository.AuditSystemModel {
	temp, err := dao.AuditSystemDAO.GetDataForAuditByIDTx(tx, action, contextModel, timeNow, tableName, id, createdBy)
	if err.Error != nil {
		loggerModel := contextModel.LoggerModel
		loggerModel.Status = 500
		loggerModel.Message = err.CausedBy.Error()
		util.LogError(loggerModel.ToLoggerObject())
	}

	return temp
}

func (input AbstractService) ServiceWithBackgroundProcess(db *sql.DB, isAlertWhenError bool, parentJob repository.JobProcessModel, childJob repository.JobProcessModel, task backgroundJobModel.ChildTask, contextModel applicationModel.ContextModel) {
	var (
		err, err2 errorModel.ErrorModel
		total     int
	)

	defer func() {
		timeNow := time.Now()
		if err.Error != nil {
			if isAlertWhenError {
				//todo save alert
				input.LogError(err, contextModel)
			}

			if parentJob.JobID.String != "" {
				parentJob.Status.String = constanta.JobProcessErrorStatus
				parentJob.UpdatedAt.Time = timeNow
				if parentJob.ContentDataOut.String == "" {
					parentJob.ContentDataOut.String = GetErrorMessage(err, contextModel)
				}

				err2 = dao.JobProcessDAO.UpdateErrorJobProcess(db, parentJob)
				if err2.Error != nil {
					input.LogError(err2, contextModel)
				}
			}

			childJob.Status.String = constanta.JobProcessErrorStatus
			childJob.UpdatedAt.Time = timeNow
			if childJob.ContentDataOut.String == "" {
				childJob.ContentDataOut.String = GetErrorMessage(err, contextModel)
			}

			err = dao.JobProcessDAO.UpdateErrorJobProcess(db, childJob)
			if err.Error != nil {
				input.LogError(err, contextModel)
			}
		} else {
			childJob.UpdatedAt.Time = timeNow
			err = dao.JobProcessDAO.UpdateJobProcessCounter(db, childJob)
			if err.Error != nil {
				input.LogError(err, contextModel)
			}

			if childJob.Total.Int32 == childJob.Counter.Int32 {
				successLog := fmt.Sprintf("Success do %s data %d from %d", childJob.Name.String, childJob.Counter.Int32, childJob.Total.Int32)
				LogMessage(successLog, 200)
			}
		}
	}()

	if parentJob.JobID.String != "" {
		childJob.ParentJobID = parentJob.JobID
		childJob.Level.Int32 = parentJob.Level.Int32 + 1
	} else {
		childJob.Level.Int32 = 0
	}

	childJob.Counter.Int32 = 0
	total, err = task.GetCountData(db, task.Data.SearchByParam, task.Data.IsCheckStatus, task.Data.CreatedBy)
	if err.Error != nil {
		return
	}

	childJob.Total.Int32 = int32(total)
	childJob.Parameter.String = util.StructToJSON(task.Data)

	//if childJob.Total.Int32 < 1 {
	//	LogMessage("No Data " + childJob.Name.String, 200)
	//	return
	//}

	err = dao.JobProcessDAO.InsertJobProcess(db, childJob)
	if err.Error != nil {
		return
	}

	//--- todo kill update job process every x minute if error found
	go input.DoUpdateJobEveryXMinute(constanta.UpdateLastUpdateTimeInMinute, childJob, contextModel)
	if task.DoJob != nil {
		err = task.DoJob(db, task.Data.Data, &childJob)
		if err.Error != nil {
			return
		}
	}

	if task.DoJobWithCtx != nil {
		err = task.DoJobWithCtx(db, task.Data.Data, &childJob, contextModel)
		if err.Error != nil {
			return
		}
	}

	if parentJob.JobID.String != "" {
		err = dao.JobProcessDAO.UpdateParentJobProcessCounter(db, childJob)
		if err.Error != nil {
			return
		}
	}
}

func (input AbstractService) ServiceWithChildBackgroundProcess(db *sql.DB, isAlertWhenError bool, listChildTask []backgroundJobModel.ChildTask, parentJob repository.JobProcessModel, contextModel applicationModel.ContextModel) {
	var (
		err     errorModel.ErrorModel
		timeNow = time.Now()
	)

	parentJob.Total.Int32 = int32(len(listChildTask))
	err = dao.JobProcessDAO.InsertJobProcess(db, parentJob)
	if err.Error != nil {
		return
	}

	go input.DoUpdateJobEveryXMinute(constanta.UpdateLastUpdateTimeInMinute, parentJob, contextModel)
	for i := 0; i < len(listChildTask); i++ {
		childJob := GetJobProcess(listChildTask[i], contextModel, timeNow)
		go input.ServiceWithBackgroundProcess(db, isAlertWhenError, parentJob, childJob, listChildTask[i], contextModel)
	}
}

func (input AbstractService) ServiceWithChildBackgroundProcessWithoutConcurrent(db *sql.DB, isAlertWhenError bool, listChildTask []backgroundJobModel.ChildTask, parentJob repository.JobProcessModel, contextModel applicationModel.ContextModel) {
	var err errorModel.ErrorModel

	parentJob.Total.Int32 = int32(len(listChildTask))

	err = dao.JobProcessDAO.InsertJobProcess(db, parentJob)
	if err.Error != nil {
		return
	}

	go input.DoUpdateJobEveryXMinute(constanta.UpdateLastUpdateTimeInMinute, parentJob, contextModel)

	timeNow := time.Now()
	for i := 0; i < len(listChildTask); i++ {
		childJob := GetJobProcess(listChildTask[i], contextModel, timeNow)
		input.ServiceWithBackgroundProcess(db, isAlertWhenError, parentJob, childJob, listChildTask[i], contextModel)
	}
}

func (input AbstractService) DoUpdateJobEveryXMinute(xMinute int, job repository.JobProcessModel, contextModel applicationModel.ContextModel) {
	var err errorModel.ErrorModel
	for true {
		time.Sleep(time.Duration(xMinute) * time.Minute)
		job, err = input.doUpdateJobUpdateAtOnDB(job, contextModel)
		if err.Error != nil {
			input.LogError(err, contextModel)
			return
		}

		if job.Status.String == constanta.JobProcessDoneStatus || job.Status.String == constanta.JobProcessErrorStatus {
			break
		}
	}
}

func (input AbstractService) doUpdateJobUpdateAtOnDB(job repository.JobProcessModel, contextModel applicationModel.ContextModel) (result repository.JobProcessModel, err errorModel.ErrorModel) {
	var (
		funcName = "doUpdateJobUpdateAtOnDB"
		tx       *sql.Tx
		errs     error
	)

	tx, errs = serverconfig.ServerAttribute.DBConnection.Begin()
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	defer func() {
		if errs != nil && err.Error != nil {
			_ = tx.Rollback()
			if errs != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
			}

			input.LogError(err, contextModel)
		} else {
			_ = tx.Commit()
		}
	}()

	result, err = dao.JobProcessDAO.GetJobProcessForUpdate(tx, job)
	if err.Error != nil {
		return
	}

	if result.Status.String == constanta.JobProcessDoneStatus || result.Status.String == constanta.JobProcessErrorStatus {
		return
	}

	job.UpdatedAt.Time = time.Now()
	err = dao.JobProcessDAO.UpdateJobProcessUpdateAt(tx, result)
	if err.Error != nil {
		return
	}

	return result, errorModel.GenerateNonErrorModel()
}

func (input AbstractService) DoAutoGenerateDataScope(tx *sql.Tx, fieldName string, id int64, contextModel *applicationModel.ContextModel, timeNow time.Time) (result int64, err errorModel.ErrorModel) {
	newScope := "nexsoft." + fieldName + ":" + strconv.Itoa(int(id))

	return dao.DataScopeDAO.InsertDataScope(tx, repository.DataScopeModel{
		Scope:         sql.NullString{String: newScope},
		Description:   sql.NullString{String: "Scope for " + fieldName + " on id " + strconv.Itoa(int(id))},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		Deleted:       sql.NullBool{},
	})
}
func (input AbstractService) SetListScope() (result []string) {
	for key := range input.MappingScopeDB {
		result = append(result, key)
	}
	return result
}

func (input AbstractService) GenerateDataScope(tx *sql.Tx, id int64, tableName string, scopeID string, userID int64, clientID string, timeNow time.Time) (dataAudit repository.AuditSystemModel, err errorModel.ErrorModel) {
	dataScope := repository.DataScopeModel{
		Scope:         sql.NullString{String: scopeID + ":" + strconv.Itoa(int(id))},
		Description:   sql.NullString{String: "Data Scope For Nextrac2 With ID " + strconv.Itoa(int(id)) + " for table : " + tableName},
		CreatedBy:     sql.NullInt64{Int64: userID},
		CreatedClient: sql.NullString{String: clientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: userID},
		UpdatedClient: sql.NullString{String: clientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	var dataScopeID int64
	dataScopeID, err = dao.DataScopeDAO.InsertDataScope(tx, dataScope)
	if err.Error != nil {
		if err.CausedBy != nil {
			if CheckDBError(err, "uq_scope") {
				err = errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.Scope)
				return
			} else if CheckDBError(err, "uq_datascope_scope") {
				err = errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.Scope)
				return
			}
		}
		return
	}

	dataAudit = repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.DataScopeDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: dataScopeID},
	}

	return
}

func (input AbstractService) TrimDTO(dto interface{}) {
	reflectType := reflect.TypeOf(dto).Elem()
	reflectValue := reflect.ValueOf(dto).Elem()
	for i := 0; i < reflectType.NumField(); i++ {
		currentField := reflectType.Field(i)
		currentValue := reflectValue.FieldByName(currentField.Name)
		if reflect.String.String() == currentField.Type.String() {
			currentValueString := currentValue.String()
			currentValue.SetString(strings.Trim(currentValueString, " "))
			if currentField.Tag.Get("auto_fix_name") == "true" {
				currentValue.SetString(strings.Title(currentValueString))
			}
		}
	}
}

func (input AbstractService) GetResponseMessage(messageID string, contextModel *applicationModel.ContextModel) (output out.StatusResponse) {
	paramOutput := make(map[string]interface{})
	paramOutput["SERVICE_NAME"] = util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.CommonServiceBundle, input.ServiceName, contextModel.AuthAccessTokenModel.Locale, nil)

	output.Code = util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil)
	output.Message = util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.CommonServiceBundle, messageID, contextModel.AuthAccessTokenModel.Locale, paramOutput)

	return output
}

func (input AbstractService) CheckUserLimitedByOwnAccess(contextModel *applicationModel.ContextModel, comparedID int64) (err errorModel.ErrorModel) {
	var (
		fileName = "AbstractService.go"
		funcName = "checkUserLimitedByOwnAccess"
	)

	//---------- Check Created By Limited
	if contextModel.LimitedByCreatedBy > 0 && (comparedID != contextModel.LimitedByCreatedBy) {
		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		return
	}

	return errorModel.GenerateNonErrorModel()
}

func (input AbstractService) ValidateDataScope(contextModel *applicationModel.ContextModel, keyScope string) (output map[string]interface{}, err errorModel.ErrorModel) {
	funcName := "ValidateDataScope"

	output = ValidateScope(contextModel, []string{keyScope})
	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input AbstractService) ValidateMultipleDataScope(contextModel *applicationModel.ContextModel, keyScope []string) (output map[string]interface{}, err errorModel.ErrorModel) {
	funcName := "ValidateMultipleDataScope"

	output = ValidateScope(contextModel, keyScope)
	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input AbstractService) SendMessageToEmail(validationResult repository.UserRegistrationDetailMapping, subject string, purpose string, userVerifyModel repository.UserVerificationModel,
	linkChannel string, linkQueryEmail func(string, repository.UserRegistrationDetailMapping, repository.UserVerificationModel) (string, errorModel.ErrorModel)) (err errorModel.ErrorModel) {

	var (
		funcName     = "SendMessageToEmail"
		templatePath = config.ApplicationConfiguration.GetDataDirectory().BaseDirectoryPath + config.ApplicationConfiguration.GetDataDirectory().Template
		template     = fmt.Sprintf(`%s/template-email-named_user.html`, templatePath)
		errorS       error
		clientType   string
	)

	var message = ""
	switch validationResult.PKCEClientMapping.ClientTypeID.Int64 {
	case constanta.ResourceNexmileID:
		clientType = constanta.Nexmile
	case constanta.ResourceNexstarID:
		clientType = constanta.Nexstar
	case constanta.ResourceNextradeID:
		clientType = constanta.Nextrade
	default:
		message = fmt.Sprintf(`%s %s %s`, constanta.Nexmile, constanta.Nexstar, constanta.Nextrade)
		err = errorModel.GenerateErrorEksternalClientTypeMustHave(input.FileName, funcName, message)
		return
	}

	subject += fmt.Sprintf(` %s`, clientType)
	reqEmail := util2.NewRequestMail(subject)
	reqEmail.To = []mail.Address{{Address: validationResult.UserRegistrationDetail.Email.String}}

	//--- Create Query Link
	linkChannel, err = linkQueryEmail(linkChannel, validationResult, userVerifyModel)
	if err.Error != nil {
		return
	}

	//--- Generate Email Template
	emailTemplateData := util2.TemplateDataActivationUserNexMile{
		Purpose:     purpose,
		Name:        validationResult.User.FirstName.String,
		ClientType:  clientType,
		UniqueID1:   validationResult.UserRegistrationDetail.UniqueID1.String,
		CompanyName: validationResult.PKCEClientMapping.CompanyName.String,
		UniqueID2:   validationResult.UserRegistrationDetail.UniqueID2.String,
		BranchName:  validationResult.PKCEClientMapping.BranchName.String,
		SalesmanID:  validationResult.UserRegistrationDetail.SalesmanID.String,
		UserID:      validationResult.UserRegistrationDetail.UserID.String,
		Password:    validationResult.UserRegistrationDetail.Password.String,
		OTP:         userVerifyModel.EmailCode.String,
		Email:       validationResult.UserRegistrationDetail.Email.String,
		ClientID:    validationResult.PKCEClientMapping.ClientID.String,
		AuthUserID:  validationResult.UserRegistrationDetail.AuthUserID.Int64,
		Link:        linkChannel,
	}

	//--- Get Template File
	template, errorS = filepath.Abs(template)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	//--- Generate Email Template
	errorS = reqEmail.GenerateEmailTemplate(template, emailTemplateData, constanta.IndonesianLanguage)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	//--- Send Email
	errorS = reqEmail.SendEmail()
	util2.SendMessageToEmail(validationResult.UserRegistrationDetail.Email.String, subject, message, applicationModel.LoggerModel{})
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	LogMessage(fmt.Sprintf(`Success Send Email Process`), 200)
	err = errorModel.GenerateNonErrorModel()
	return
}

func GetResponseMessages(messageID string, contextModel *applicationModel.ContextModel) (output out.StatusResponse) {
	paramOutput := make(map[string]interface{})
	paramOutput["SERVICE_NAME"] = util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.CommonServiceBundle, "GetResponseMessages", contextModel.AuthAccessTokenModel.Locale, nil)

	output.Code = util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil)
	output.Message = util2.GenerateI18NServiceMessage(serverconfig.ServerAttribute.CommonServiceBundle, messageID, contextModel.AuthAccessTokenModel.Locale, paramOutput)

	return output
}
