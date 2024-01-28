package PKCEClientMappingService

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

type updatePKCEClientNameService struct {
	service.AbstractService
}

var UpdatePKCEClientNameService = updatePKCEClientNameService{}.New()

func (input updatePKCEClientNameService) New() (output updatePKCEClientNameService) {
	output.FileName = "UpdatePKCEClientNameService.go"
	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "pkce_client_mapping.client_type_id",
		Count: "pkce_client_mapping.client_type_id",
	}
	return
}

func (input updatePKCEClientNameService) UpdateClientName(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		pkceClientMappingBody in.PKCERequest
		funcName              = "UpdateClientName"
	)

	pkceClientMappingBody, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateClientName)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, pkceClientMappingBody, contextModel, input.updateClientName, func(interface{}, applicationModel.ContextModel) {
		//--- func additional
	})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("OK", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18Message("SUCCESS_UPDATE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input updatePKCEClientNameService) updateClientName(tx *sql.Tx, body interface{}, contextModel *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		pkceClientBody       = body.(in.PKCERequest)
		pkceClientRepository = input.getPKCEClientRepository(pkceClientBody, contextModel.AuthAccessTokenModel, now)
		isOnlyHaveOwnAccess  bool
		pkceClientOnDB       repository.PKCEClientMappingModel
		funcName             = "updateClientName"
		scope                map[string]interface{}
	)

	scope, err = input.validateDataScopePKCEClientMapping(contextModel)
	if err.Error != nil {
		return
	}

	_, isOnlyHaveOwnAccess = service.CheckIsOnlyHaveOwnPermission(*contextModel)
	if isOnlyHaveOwnAccess {
		pkceClientRepository.ClientID.String = contextModel.AuthAccessTokenModel.ClientID
	}

	pkceClientOnDB, err = dao.PKCEClientMappingDAO.GetPKCEClientMappingForUpdateByType(tx, pkceClientRepository, constanta.Nexmile, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if pkceClientOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.PKCEClientMapping)
		return
	}

	if pkceClientOnDB.UpdatedAt.Time != pkceClientBody.UpdatedAt {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.PKCEClientMapping)
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, now, dao.PKCEClientMappingDAO.TableName, pkceClientOnDB.ID.Int64, 0)...)
	err = dao.PKCEClientMappingDAO.UpdatePKCECLientMapping(tx, pkceClientRepository)
	return
}

func (input updatePKCEClientNameService) getPKCEClientRepository(pkceClientBody in.PKCERequest, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.PKCEClientMappingModel {
	return repository.PKCEClientMappingModel{
		ID:            sql.NullInt64{Int64: pkceClientBody.ID},
		ClientAlias:   sql.NullString{String: pkceClientBody.ClientAlias},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input updatePKCEClientNameService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.PKCERequest) errorModel.ErrorModel) (inputStruct in.PKCERequest, err errorModel.ErrorModel) {
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

func (input updatePKCEClientNameService) validateUpdateClientName(inputStruct *in.PKCERequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateClientName()
}

func (input updatePKCEClientNameService) validateDataScopePKCEClientMapping(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	funcName := "validateDataScopePKCEClientMapping"

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
