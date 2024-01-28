package AuditMonitoringService

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
)

func (input auditMonitoringService) ViewAuditMonitoringData(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	inputStruct, err := input.readBodyAndValidate(request, contextModel, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewAuditMonitoringData(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input auditMonitoringService) doViewAuditMonitoringData(inputStruct in.AuditMonitoringRequest, contextModel *applicationModel.ContextModel) (result out.ViewAuditMonitoringResponse, err errorModel.ErrorModel) {
	funcName := "doViewAuditMonitoringData"
	auditModel := repository.AuditSystemModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	}

	auditModel.CreatedBy.Int64 = contextModel.LimitedByCreatedBy

	auditModel, err = dao.AuditSystemDAO.ViewAuditData(serverconfig.ServerAttribute.DBConnection, auditModel)
	if err.Error != nil {
		return
	}

	if auditModel.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	result = reformatRepositoryToDTOOut(auditModel)

	return
}

func reformatRepositoryToDTOOut(auditModel repository.AuditSystemModel) out.ViewAuditMonitoringResponse {
	var data map[string]interface{}
	_ = json.Unmarshal([]byte(auditModel.Data.String), &data)

	CensoringSecretData(auditModel.TableName.String, &data)

	temp := out.ViewAuditMonitoringResponse{
		ID:            auditModel.ID.Int64,
		TableName:     auditModel.TableName.String,
		UUIDKey:       auditModel.UUIDKey.String,
		Data:          data,
		PrimaryKey:    auditModel.PrimaryKey.Int64,
		Action:        auditModel.Action.Int32,
		CreatedBy:     auditModel.CreatedBy.Int64,
		CreatedClient: auditModel.CreatedClient.String,
		CreatedAt:     auditModel.CreatedAt.Time,
	}

	return temp
}

func (input auditMonitoringService) validateView(inputStruct *in.AuditMonitoringRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewAuditMonitoring()
}
