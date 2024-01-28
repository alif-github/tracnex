package ClientMappingService

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
)

type uiViewClientMappingService struct {
	FileName string
}

var UIViewClientMappingService = uiViewClientMappingService{FileName: "UIViewClientMappingService.go"}

func (input uiViewClientMappingService) ViewClientMapping(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var clientMappingBody in.ClientMappingForUIRequest
	clientMappingBody, err = input.readPathParamAndValidate(request)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewClientMapping(clientMappingBody, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18Message("SUCCESS_VIEW_CLIENT_MAPPING_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	return
}

func (input uiViewClientMappingService) doViewClientMapping(clientMappingBody in.ClientMappingForUIRequest, contextModel *applicationModel.ContextModel) (result out.ClientMappingForView, err errorModel.ErrorModel) {
	var (
		funcName          = "doViewClientMapping"
		clientMappingOnDB repository.ClientMappingForViewModel
		db                = serverconfig.ServerAttribute.DBConnection
		scope             map[string]interface{}
	)

	scope, err = input.validateDataScopeClientMapping(contextModel)
	if err.Error != nil {
		return
	}

	mappingScopeDB := make(map[string]applicationModel.MappingScopeDB)
	mappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View:  "client_mapping.client_type_id",
		Count: "client_mapping.client_type_id",
	}

	clientMappingModel := repository.ClientMappingModel{
		ID:        sql.NullInt64{Int64: clientMappingBody.ID},
		CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	}

	clientMappingOnDB, err = dao.ClientMappingDAO.ViewClientMapping(db, clientMappingModel, scope, mappingScopeDB)
	if err.Error != nil {
		return
	}

	if clientMappingOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientMapping)
		return
	}

	result = input.convertDAOToDTO(clientMappingOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input uiViewClientMappingService) convertDAOToDTO(clientMappingModel repository.ClientMappingForViewModel) (output out.ClientMappingForView) {
	return out.ClientMappingForView{
		ID:                    clientMappingModel.ID.Int64,
		ClientID:              clientMappingModel.ClientID.String,
		SocketID:              clientMappingModel.SocketID.String,
		ClientType:            clientMappingModel.ClientType.String,
		CompanyID:             clientMappingModel.CompanyID.String,
		BranchID:              clientMappingModel.BranchID.String,
		Aliases:               clientMappingModel.Aliases.String,
		UpdatedAt:             clientMappingModel.UpdatedAt.Time,
		UpdatedBy:             clientMappingModel.UpdatedBy.Int64,
		CreatedBy:             clientMappingModel.CreatedBy.Int64,
		CreatedAt:             clientMappingModel.CreatedAt.Time,
		SuccessStatusAuth:     clientMappingModel.SuccessStatusAuth.Bool,
		SuccessStatusNexcloud: clientMappingModel.SuccessStatusNexcloud.Bool,
		SuccessStatusNexdrive: clientMappingModel.SuccessStatusNexdrive.Bool,
	}
}

func (input uiViewClientMappingService) readPathParamAndValidate(request *http.Request) (clientMappingBody in.ClientMappingForUIRequest, err errorModel.ErrorModel) {
	var id int64
	id, err = readPathParam(request)
	if err.Error != nil {
		return
	}

	clientMappingBody.ID = id
	err = clientMappingBody.ValidateViewCLientMapping()
	return
}

func (input uiViewClientMappingService) validateDataScopeClientMapping(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
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
