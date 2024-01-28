package UserActivationService

import (
	"database/sql"
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
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input userActivationService) EmailActivation(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "EmailActivation"

	inputStruct, err := input.readBodyAndValidate(request, contextModel, input.validateEmailActivation)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doEmailActivation, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_ACTIVATION_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userActivationService) doEmailActivation(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName    = "EmailActivationService.go"
		funcName    = "doEmailActivation"
		inputStruct = inputStructInterface.(in.UserActivationRequest)
		resultID    repository.UserModel
	)

	err = input.HitActivationToAuth(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	resultID, err = dao.UserDAO.CheckIsAuthUserExistForUpdate(serverconfig.ServerAttribute.DBConnection, repository.UserModel{AuthUserID: sql.NullInt64{Int64: inputStruct.UserID}})
	if err.Error != nil {
		return
	}

	if resultID.ID.Int64 < 1 {
		err = errorModel.GenerateDataNotFound(funcName, fileName)
		return
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.UserDAO.TableName, resultID.ID.Int64, 0)...)
	err = input.doUpdateStatusUser(tx, inputStruct, timeNow)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userActivationService) doUpdateStatusUser(tx *sql.Tx, inputStruct in.UserActivationRequest, _ time.Time) (err errorModel.ErrorModel) {

	userTemp := repository.UserModel{
		Status:     sql.NullString{String: constanta.StatusActive},
		AuthUserID: sql.NullInt64{Int64: inputStruct.UserID},
	}

	err = dao.UserDAO.UpdateUserAfterActivation(tx, userTemp)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userActivationService) validateEmailActivation(inputStruct *in.UserActivationRequest) errorModel.ErrorModel {
	return inputStruct.ValidateActivationUser()
}
