package ClientMappingService

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"strconv"
	"time"
)

type uiUpdateClientNameService struct {
	service.AbstractService
}

var UIUpdateClientNameService = uiUpdateClientNameService{}.New()

func (input uiUpdateClientNameService) New() (output uiUpdateClientNameService) {
	output.FileName = "UIUpdateClientNameService.go"
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "client_mapping.client_type_id",
		Count: "client_mapping.client_type_id",
	}

	return
}

func (input uiUpdateClientNameService) UpdateClientName(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		clientMappingBody in.ClientMappingForUIRequest
		funcName          = "UpdateClientName"
	)

	clientMappingBody, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateClientName)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, clientMappingBody, contextModel, input.updateClientName, func(interface{}, applicationModel.ContextModel) {
		//--- func additional
	})

	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("OK", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18Message("SUCCESS_UPDATE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input uiUpdateClientNameService) updateClientName(tx *sql.Tx, body interface{}, contextModel *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		clientMappingBody       = body.(in.ClientMappingForUIRequest)
		clientMappingRepository = input.getClientmappingRepository(clientMappingBody, contextModel.AuthAccessTokenModel, now)
		isOnlyHaveOwnAccess     bool
		clientMappingOnDB       repository.ClientMappingModel
		funcName                = "updateClientName"
		scope                   map[string]interface{}
	)

	scope, err = input.validateDataScopeClientMapping(contextModel)
	if err.Error != nil {
		return
	}

	_, isOnlyHaveOwnAccess = service.CheckIsOnlyHaveOwnPermission(*contextModel)
	if isOnlyHaveOwnAccess {
		clientMappingRepository.ClientID.String = contextModel.AuthAccessTokenModel.ClientID
	}

	clientMappingOnDB, err = dao.ClientMappingDAO.GetClientMappingForUpdateByType(tx, clientMappingRepository, constanta.ND6, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if clientMappingOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientMapping)
		return
	}

	if clientMappingOnDB.UpdatedAt.Time != clientMappingBody.UpdatedAt {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.ClientMapping)
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, now, dao.ClientMappingDAO.TableName, clientMappingOnDB.ID.Int64, 0)...)
	err = dao.ClientMappingDAO.UpdateClientName(tx, clientMappingRepository)
	return
}

func (input uiUpdateClientNameService) getClientmappingRepository(clientMappingBody in.ClientMappingForUIRequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.ClientMappingModel {
	return repository.ClientMappingModel{
		ID:            sql.NullInt64{Int64: clientMappingBody.ID},
		ClientAlias:   sql.NullString{String: clientMappingBody.ClientAlias},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input uiUpdateClientNameService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ClientMappingForUIRequest) errorModel.ErrorModel) (inputStruct in.ClientMappingForUIRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidate"
		stringBody string
	)

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

func (input uiUpdateClientNameService) validateUpdateClientName(inputStruct *in.ClientMappingForUIRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateClientName()
}

func (input uiUpdateClientNameService) validateDataScopeClientMapping(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	funcName := "validateDataScopeClientMapping"

	output = service.ValidateScope(contextModel, []string{
		constanta.ClientTypeDataScope,
	})

	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
