package PKCEClientMappingService

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

func (input pkceClientMappingService) ViewPKCEClientMapping(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.PKCERequest
	inputStruct, err = input.readPathParamViewPKCEClientMapping(request)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewPKCEClientMapping(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18Message("SUCCESS_VIEW_PKCE_CLIENT_MAPPING_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingService) doViewPKCEClientMapping(inputStruct in.PKCERequest, contextModel *applicationModel.ContextModel) (output out.PKCEClientMappingForDetail, err errorModel.ErrorModel) {
	var (
		fileName               = "ViewPKCEClientMappingService.go"
		funcName               = "doViewPKCEClientMapping"
		pkceClientMappingModel repository.PKCEClientMappingModel
		isOnlyHaveOwnAccess    bool
		pkceClientMappingOnDB  repository.ViewPKCEClientMappingModel
		scope                  map[string]interface{}
	)

	scope, err = input.validateDataScopePKCEClientMapping(contextModel)
	if err.Error != nil {
		return
	}

	pkceClientMappingModel = repository.PKCEClientMappingModel{ID: sql.NullInt64{Int64: inputStruct.ID}}
	_, isOnlyHaveOwnAccess = service.CheckIsOnlyHaveOwnPermission(*contextModel)
	if isOnlyHaveOwnAccess {
		pkceClientMappingModel.ClientID.String = contextModel.AuthAccessTokenModel.ClientID
	}

	pkceClientMappingOnDB, err = dao.PKCEClientMappingDAO.GetViewPKCEClientMapping(serverconfig.ServerAttribute.DBConnection, pkceClientMappingModel, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if pkceClientMappingOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, "pkce_client_mapping")
		return
	}

	output = input.convertPKCEModelToDTOOut(pkceClientMappingOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceClientMappingService) convertPKCEModelToDTOOut(pkceClientMappingModel repository.ViewPKCEClientMappingModel) out.PKCEClientMappingForDetail {
	return out.PKCEClientMappingForDetail{
		ID:          pkceClientMappingModel.ID.Int64,
		ClientID:    pkceClientMappingModel.ClientID.String,
		FirstName:   pkceClientMappingModel.FirstName.String,
		LastName:    pkceClientMappingModel.LastName.String,
		Username:    pkceClientMappingModel.Username.String,
		ClientType:  pkceClientMappingModel.ClientType.String,
		CompanyID:   pkceClientMappingModel.CompanyID.String,
		BranchID:    pkceClientMappingModel.BranchID.String,
		ClientAlias: pkceClientMappingModel.ClientAlias.String,
		UpdatedAt:   pkceClientMappingModel.UpdatedAt.Time,
		UpdatedBy:   pkceClientMappingModel.UpdatedBy.Int64,
		CreatedBy:   pkceClientMappingModel.CreatedBy.Int64,
		CreatedAt:   pkceClientMappingModel.CreatedAt.Time,
	}
}

func (input pkceClientMappingService) readPathParamViewPKCEClientMapping(request *http.Request) (output in.PKCERequest, err errorModel.ErrorModel) {
	var (
		fileName = "ViewPKCEClientMappingService.go"
		funcName = "readPathParamViewPKCEClientMapping"
		id       int64
	)

	id, err = service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	if id < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
		return
	}

	output.ID = id
	err = errorModel.GenerateNonErrorModel()
	return
}
