package UserInvitationService

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
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
	"nexsoft.co.id/nextrac2/resource_common_service/dto/grochat_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/grochat_response"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/UserService"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input invitationService) SendInvitation(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, errModel errorModel.ErrorModel) {
	funcName := "SendInvitation"

	userInvitation, errModel := input.readBody(request, contextModel)
	if errModel.Error != nil {
		return
	}

	if errModel = userInvitation.ValidateInsert(); errModel.Error != nil {
		return
	}

	_, errModel = input.ServiceWithDataAuditPreparedByService(funcName, userInvitation, contextModel, input.sendInvitation, func(_ interface{}, _ applicationModel.ContextModel) {})
	if errModel.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: UserService.GenerateI18NMessage("SUCCESS_INSERT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	return
}

func (input invitationService) sendInvitation(tx *sql.Tx, data interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	var (
		body              in.UserInvitationRequest
		isInvitationExist bool
	)

	body, _ = data.(in.UserInvitationRequest)

	/*
		Validate Body
	*/
	if errModel = input.validateBody(tx, body); errModel.Error != nil {
		return
	}

	/*
		Request Internal Token
	*/
	internalToken := resource_common_service.GenerateInternalToken("chat", 0, "", config.ApplicationConfiguration.GetServerResourceID(), "id-ID")

	/*
		Send Invitation
	*/
	fe := config.ApplicationConfiguration.GetNextracFrontend()
	invitationReq := grochat_request.Invitation{
		Email:         body.Email,
		EmailMessage:  constanta.InvitationEmailBody,
		URLInvitation: fe.Host + fe.PathRedirect.Invitation,
		ResourceId:    config.ApplicationConfiguration.GetServerResourceID(),
		UserType:      constanta.GroChatUserTypeRegular,
	}

	invitationRes, errModel := input.requestSendInvitation(contextModel, internalToken, invitationReq)
	if errModel.Error != nil {
		return
	}

	/*
		Get existing user invitation
	*/
	userInvitation, errModel := dao.UserInvitationDAO.GetByEmailForUpdate(tx, body.Email)
	if errModel.Error != nil {
		return
	}

	isInvitationExist = userInvitation.Id.Valid
	body.Id = userInvitation.Id.Int64

	if isInvitationExist {
		/*
			Update User Invitation
		*/
		var auditUpdate []repository.AuditSystemModel

		auditUpdate, errModel = input.updateUserInvitation(tx, contextModel, timeNow, body, invitationRes)
		if errModel.Error != nil {
			return
		}

		dataAudit = append(dataAudit, auditUpdate...)
		return
	}

	/*
		Insert User Invitation
	*/
	auditInsert, errModel := input.insertUserInvitation(tx, contextModel, timeNow, body, invitationRes)
	if errModel.Error != nil {
		return
	}

	dataAudit = append(dataAudit, auditInsert)
	return
}

func (input invitationService) requestSendInvitation(contextModel *applicationModel.ContextModel, internalToken string, body grochat_request.Invitation) (*grochat_response.Invitation, errorModel.ErrorModel) {
	funcName := "requestLogin"

	groChatServer := config.ApplicationConfiguration.GetGrochat()

	path := fmt.Sprintf("%s%s", groChatServer.Host, groChatServer.PathRedirect.SendInvitation)

	headerRequest := make(map[string][]string)
	headerRequest[constanta.TokenHeaderNameConstanta] = []string{internalToken}
	headerRequest["Content-Type"] = []string{"application/json"}

	req := util.StructToJSON(body)

	/*
		Send Request
	*/
	statusCode, _, bodyResult, err := common.HitAPI(path, headerRequest, req, "POST", *contextModel)
	if err != nil {
		return nil, errorModel.GenerateUnknownError(input.FileName, funcName, err)
	}

	if statusCode != http.StatusOK {
		return nil, errorModel.GenerateUnknownError(input.FileName, funcName, errors.New(fmt.Sprintf("request failed : %d", statusCode)))
	}

	/*
		Generate Response
	*/
	response := &grochat_response.Invitation{}
	if err = json.Unmarshal([]byte(bodyResult), response); err != nil {
		return nil, errorModel.GenerateUnknownError(input.FileName, funcName, err)
	}

	if response.Note != "" {
		return nil, errorModel.GenerateAuthenticationServerError("GroChatServiceUtil.go", funcName, statusCode, "GROCHAT", errors.New(response.Note))
	}

	return response, errorModel.GenerateNonErrorModel()
}

func (input invitationService) insertUserInvitation(tx *sql.Tx, contextModel *applicationModel.ContextModel, now time.Time, body in.UserInvitationRequest, invitationRes *grochat_response.Invitation) (auditData repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	expiresOn := input.getExpiresOn(invitationRes.Data.ExpiredDate)

	id, errModel := dao.UserInvitationDAO.InsertTx(tx, repository.UserInvitation{
		InvitationCode: sql.NullString{String: invitationRes.Data.InvitationCode},
		Email:          sql.NullString{String: body.Email},
		ClientId:       sql.NullString{String: invitationRes.Data.ClientId},
		RoleId:         sql.NullInt64{Int64: body.RoleId},
		DataGroupId:    sql.NullInt64{Int64: body.DataGroupId},
		ExpiresOn:      sql.NullTime{Time: expiresOn},
		CreatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:      sql.NullTime{Time: now},
		UpdatedAt:      sql.NullTime{Time: now},
		CreatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
	})
	if errModel.Error != nil {
		return
	}

	auditData = repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.UserInvitationDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: id},
	}

	return auditData, errorModel.GenerateNonErrorModel()
}

func (input invitationService) updateUserInvitation(tx *sql.Tx, contextModel *applicationModel.ContextModel, now time.Time, body in.UserInvitationRequest, invitationRes *grochat_response.Invitation) (auditData []repository.AuditSystemModel, errModel errorModel.ErrorModel) {
	expiresOn := input.getExpiresOn(invitationRes.Data.ExpiredDate)

	model := repository.UserInvitation{
		Id: 			sql.NullInt64{Int64: body.Id},
		InvitationCode: sql.NullString{String: invitationRes.Data.InvitationCode},
		Email:          sql.NullString{String: body.Email},
		ClientId:       sql.NullString{String: invitationRes.Data.ClientId},
		RoleId:         sql.NullInt64{Int64: body.RoleId},
		DataGroupId:    sql.NullInt64{Int64: body.DataGroupId},
		ExpiresOn:      sql.NullTime{Time: expiresOn},
		CreatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedClient:  sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:      sql.NullTime{Time: now},
		UpdatedAt:      sql.NullTime{Time: now},
		CreatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedBy:      sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, now, dao.UserDAO.TableName, body.Id, 0)...)

	errModel = dao.UserInvitationDAO.UpdateByIdTx(tx, model)
	return
}

func (input invitationService) validateBody(tx *sql.Tx, body in.UserInvitationRequest) (errModel errorModel.ErrorModel) {
	funcName := "validateBody"

	role, errModel := dao.RoleDAO.GetByIdTx(tx, body.RoleId)
	if errModel.Error != nil {
		return
	}

	if !role.ID.Valid {
		return errorModel.GenerateUnknownDataError(input.FileName, funcName, "role_id")
	}

	dataGroup, errModel := dao.DataGroupDAO.GetByIdTx(tx, body.DataGroupId)
	if errModel.Error != nil {
		return
	}

	if !dataGroup.ID.Valid {
		return errorModel.GenerateUnknownDataError(input.FileName, funcName, "data_group_id")
	}

	return errorModel.GenerateNonErrorModel()
}

func (input invitationService) getExpiresOn(t string) time.Time {
	result, _ := time.Parse("02-01-2006 15:04:05", t)
	return result
}
