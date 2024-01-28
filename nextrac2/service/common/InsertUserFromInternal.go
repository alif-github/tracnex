package common

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"strings"
	"time"
)

type insertUserFromInternalService struct {
	service.AbstractService
}

var InsertUserFromInternalService = insertUserFromInternalService{}.New()

func (input insertUserFromInternalService) New() (output insertUserFromInternalService) {
	output.FileName = "InsertUserFromInternal.go"
	return
}

func checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_user_clientid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ClientID)
		} else if service.CheckDBError(err, "uq_user_authuserid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.AuthUserID)
		}
	}
	return err
}

func (input insertUserFromInternalService) StartService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.UserRequest

	inputStruct, err = input.readBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService("InsertInternalUser", inputStruct, contextModel, DoInsertUser, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func DoInsertUser(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "doInsertUser"

	inputStruct := inputStructInterface.(in.UserRequest)
	var roleModel repository.RoleModel
	var groupModel repository.DataGroupModel

	roleModel, err = dao.RoleDAO.GetRoleByName(serverconfig.ServerAttribute.DBConnection, repository.RoleModel{RoleID: sql.NullString{String: inputStruct.Role}})
	if err.Error != nil {
		return
	}

	if roleModel.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError("InsertUserFromInternal.go", funcName, constanta.Role)
		return
	}

	groupModel, err = dao.DataGroupDAO.GetRoleByName(tx, repository.DataGroupModel{GroupID: sql.NullString{String: inputStruct.Group}})
	if err.Error != nil {
		return
	}

	if groupModel.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError("InsertUserFromInternal.go", funcName, constanta.DataGroup)
		return
	}

	clientRoleScope := repository.ClientRoleScopeModel{
		RoleID:  roleModel.ID,
		GroupID: groupModel.ID,
	}

	statusCode, bodyResult, errs := hitAuthenticationServer(inputStruct, contextModel)
	if errs != nil {
		err = errorModel.GenerateUnknownError("InsertUserFromInternal.go", funcName, errs)
		return
	}

	if statusCode == 200 {
		var registerUserResponse authentication_response.RegisterUserAuthenticationResponse
		_ = json.Unmarshal([]byte(bodyResult), &registerUserResponse)

		dataAudit, err = saveUserToDB(tx, inputStruct, registerUserResponse, clientRoleScope, contextModel, timeNow)
		if err.Error != nil {
			return
		}

		contextModel.LoggerModel.Status = statusCode
	} else {
		err = common.ReadAuthServerError("DoInsertUser", statusCode, bodyResult, contextModel)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func hitAuthenticationServer(inputStruct in.UserRequest, contextModel *applicationModel.ContextModel) (statusCode int, bodyResult string, err error) {
	registerAuthentication := authentication_request.UserAuthenticationDTO{
		Username:              inputStruct.Username,
		Password:              inputStruct.Password,
		FirstName:             inputStruct.FirstName,
		LastName:              inputStruct.LastName,
		Email:                 inputStruct.Email,
		CountryCode:           inputStruct.CountryCode,
		Phone:                 inputStruct.Phone,
		Device:                inputStruct.Device,
		Locale:                inputStruct.Locale,
		EmailMessage:          getEmailMessage(inputStruct),
		IPWhitelist:           inputStruct.IPWhitelist,
		EmailLinkMessage:      config.ApplicationConfiguration.GetNextracFrontend().Host + config.ApplicationConfiguration.GetNextracFrontend().PathRedirect.VerifyUserPath,
		PhoneMessage:          "test",
		ResourceID:            config.ApplicationConfiguration.GetServerResourceID(),
		AdditionalInformation: inputStruct.AdditionalInformation,
	}

	if !inputStruct.IsAdmin {
		hostSysUser := config.ApplicationConfiguration.GetNextracFrontend().Host
		pathSysUser := strings.Replace(config.ApplicationConfiguration.GetNextracFrontend().PathRedirect.VerifyUserPath, "/nexsoft-admin", "", 1)
		registerAuthentication.EmailLinkMessage = hostSysUser + pathSysUser
	}

	internalToken := resource_common_service.GenerateInternalToken("auth", contextModel.AuthAccessTokenModel.AuthenticationServerUserID, contextModel.AuthAccessTokenModel.ClientID, config.ApplicationConfiguration.GetServerResourceID(), constanta.DefaultApplicationsLanguage)
	contextModel.LoggerModel.Class = "[" + "InsertUserFromInternal.go" + "," + "doInsertUser" + "]"
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	registerUserUrl := authenticationServer.Host + authenticationServer.PathRedirect.InternalUser.CrudUser

	statusCode, bodyResult, err = common.HitRegisterUserAuthenticationServer(internalToken, registerUserUrl, registerAuthentication, contextModel)
	return
}

func saveUserToDB(tx *sql.Tx, inputStruct in.UserRequest, registerUserResponse authentication_response.RegisterUserAuthenticationResponse, clientRoleScope repository.ClientRoleScopeModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	content := registerUserResponse.Nexsoft.Payload.Data.Content
	userModel := repository.UserModel{
		ClientID:       sql.NullString{String: content.ClientID},
		AuthUserID:     sql.NullInt64{Int64: content.UserID},
		Locale:         sql.NullString{String: inputStruct.Locale},
		SignatureKey:   sql.NullString{String: content.SignatureKey},
		AdditionalInfo: sql.NullString{String: inputStruct.AdditionalInformationString()},
		CreatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedAt:      sql.NullTime{Time: timeNow},
		CreatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedAt:      sql.NullTime{Time: timeNow},
		UpdatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
	}

	var id int64
	id, err = dao.UserDAO.InsertUser(tx, userModel)
	if err.Error != nil {
		err = checkDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.UserDAO.TableName, id, contextModel.LimitedByCreatedBy)...)

	clientRoleScope.ClientID.String = content.ClientID
	clientRoleScope.CreatedBy = userModel.CreatedBy
	clientRoleScope.CreatedAt = userModel.CreatedAt
	clientRoleScope.CreatedClient = userModel.CreatedClient
	clientRoleScope.UpdatedBy = userModel.UpdatedBy
	clientRoleScope.UpdatedAt = userModel.UpdatedAt
	clientRoleScope.UpdatedClient = userModel.UpdatedClient

	id, err = dao.ClientRoleScopeDAO.InsertClientRoleScope(tx, clientRoleScope)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.ClientRoleScopeDAO.TableName, id, contextModel.LimitedByCreatedBy)...)

	//personProfileModel := repository.PersonProfileModel{
	//	AuthUserID:     sql.NullInt64{Int64: content.UserID},
	//	FirstName:      sql.NullString{String: inputStruct.FirstName},
	//	LastName:       sql.NullString{String: inputStruct.LastName},
	//	Phone:          sql.NullString{String: inputStruct.CountryCode + "-" + inputStruct.Phone},
	//	Email:          sql.NullString{String: inputStruct.Email},
	//	AdditionalInfo: userModel.AdditionalInfo,
	//	CreatedBy:      userModel.CreatedBy,
	//	CreatedAt:      userModel.CreatedAt,
	//	CreatedClient:  userModel.CreatedClient,
	//	UpdatedBy:      userModel.UpdatedBy,
	//	UpdatedAt:      userModel.UpdatedAt,
	//	UpdatedClient:  userModel.UpdatedClient,
	//}
	//
	//var personProfileDB repository.PersonProfileModel
	//
	//personProfileDB, err = dao.PersonProfileDAO.GetPersonProfileByPhoneAndEmail(tx, personProfileModel)
	//if err.Error != nil {
	//	return
	//}
	//
	//if personProfileDB.ID.Int64 == 0 {
	//	id, err = dao.PersonProfileDAO.InsertUserPersonProfile(tx, personProfileModel)
	//	if err.Error != nil {
	//		return
	//	}
	//	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditInsertConstanta, *contextModel, timeNow, dao.PersonProfileDAO.TableName, id, 0)...)
	//
	//} else {
	//	if personProfileDB.AuthUserID.Int64 == 0 {
	//		personProfileDB.UpdatedAt = personProfileModel.UpdatedAt
	//		personProfileDB.UpdatedClient = personProfileModel.UpdatedClient
	//		personProfileDB.UpdatedBy = personProfileModel.UpdatedBy
	//		personProfileDB.AuthUserID = personProfileModel.AuthUserID
	//		err = dao.PersonProfileDAO.UpdatePersonProfileAuthUserID(tx, personProfileDB)
	//		if err.Error != nil {
	//			return
	//		}
	//		dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.PersonProfileDAO.TableName, personProfileDB.ID.Int64, contextModel.LimitedByCreatedBy)...)
	//
	//	}
	//}

	err = errorModel.GenerateNonErrorModel()
	return
}

func getEmailMessage(inputStruct in.UserRequest) string {
	param := make(map[string]interface{})
	param["USERNAME"] = inputStruct.Username
	param["EMAIL"] = inputStruct.Email
	param["PHONE"] = inputStruct.Phone
	param["ACTIVATION_CODE"] = "{{.ACTIVATION_CODE}}"
	param["ACTIVATION_LINK"] = "{{.ACTIVATION_LINK}}"
	param["USER_ID"] = "{{.USER_ID}}"

	return util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.CommonServiceBundle, "VERIFY_EMAIL_MESSAGE", inputStruct.Locale, param)
}

func (input insertUserFromInternalService) readBody(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.UserRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	err = inputStruct.ValidateInternalInsertUser()

	return
}
