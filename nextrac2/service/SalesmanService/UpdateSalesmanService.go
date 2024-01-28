package SalesmanService

import (
	"database/sql"
	"net/http"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/MasterDataService/DistrictService"
	"nexsoft.co.id/nextrac2/service/MasterDataService/ProvinceService"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input salesmanService) UpdateSalesman(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateSalesman"
		inputStruct in.SalesmanRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateSalesman)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateSalesman, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("OK", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_UPDATE_SALESMAN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input salesmanService) doUpdateSalesman(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName        = "UpdateSalesmanService.go"
		funcName        = "doUpdateSalesman"
		inputStruct     = inputStructInterface.(in.SalesmanRequest)
		salesmanOnDB    repository.SalesmanModel
		dataResultTitle out.PersonTitleResponse
		scopeLimit      map[string]interface{}
		provinceOnDB    repository.ProvinceModel
		DistrictOnDB    repository.DistrictModel
		db              = serverconfig.ServerAttribute.DBConnection
		internalToken   = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, constanta.DefaultApplicationsLanguage)
	)

	scopeLimit, err = input.validateDataScopeSalesman(contextModel)
	if err.Error != nil {
		return
	}

	if util2.IsStringEmpty(inputStruct.Status) {
		inputStruct.Status = "A"
	}

	parameterSalesmanModel := repository.SalesmanModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		PersonTitleID: sql.NullInt64{Int64: inputStruct.PersonTitleID},
		Sex:           sql.NullString{String: inputStruct.Sex},
		FirstName:     sql.NullString{String: inputStruct.FirstName},
		LastName:      sql.NullString{String: inputStruct.LastName},
		Address:       sql.NullString{String: inputStruct.Address},
		Hamlet:        sql.NullString{String: inputStruct.Hamlet},
		Neighbourhood: sql.NullString{String: inputStruct.Neighbourhood},
		ProvinceID:    sql.NullInt64{Int64: inputStruct.ProvinceID},
		DistrictID:    sql.NullInt64{Int64: inputStruct.DistrictID},
		Phone:         sql.NullString{String: inputStruct.Phone},
		Email:         sql.NullString{String: inputStruct.Email},
		Status:        sql.NullString{String: inputStruct.Status},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	salesmanOnDB, err = dao.SalesmanDAO.GetSalesmanForUpdateDelete(serverconfig.ServerAttribute.DBConnection, parameterSalesmanModel, scopeLimit, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if salesmanOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.SalesmanID)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, salesmanOnDB.CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if salesmanOnDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(fileName, funcName, "Salesman")
		return
	}

	if inputStruct.UpdatedAt != salesmanOnDB.UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.UpdatedAt)
		return
	}

	dataResultTitle, err = master_data_dao.ViewDetailPersonTitleFromMasterData(int(parameterSalesmanModel.PersonTitleID.Int64), contextModel)
	if err.Error != nil {
		if err.Error.Error() == "E-4-MAD-SRV-004" {
			err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.PersonTitleID)
		}
		return
	}

	if dataResultTitle.ID < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.PersonTitleID)
		return
	}

	parameterSalesmanModel.PersonTitle.String = dataResultTitle.Title

	//--- Validate Province (Local Validation)
	provinceOnDB, err = dao.ProvinceDAO.GetProvinceForCustomer(db, repository.ProvinceModel{
		ID: sql.NullInt64{Int64: inputStruct.ProvinceID},
	}, scopeLimit, ProvinceService.ProvinceService.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if provinceOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Province)
		return
	}

	//--- MDB Validation
	_, err = master_data_dao.ViewDetailProvinceFromMasterData(int(provinceOnDB.MDBProvinceID.Int64), contextModel, internalToken)
	if err.Error != nil {
		return
	}

	//--- Validate District (Local Validation)
	DistrictOnDB, err = dao.DistrictDAO.GetDistrictWithProvinceID(db, repository.ListLocalDistrictModel{
		ID:         sql.NullInt64{Int64: inputStruct.DistrictID},
		ProvinceID: sql.NullInt64{Int64: inputStruct.ProvinceID},
	}, scopeLimit, DistrictService.DistrictService.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if DistrictOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.District)
		return
	}

	//--- MDB Validation
	_, err = master_data_dao.ViewDetailDistrictFromMasterData(int(DistrictOnDB.MdbDistrictID.Int64), contextModel, internalToken)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.SalesmanDAO.TableName, salesmanOnDB.ID.Int64, 0)...)
	err = dao.SalesmanDAO.UpdateSalesman(tx, parameterSalesmanModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanService) validateUpdateSalesman(inputStruct *in.SalesmanRequest) errorModel.ErrorModel {
	return inputStruct.ValidationUpdateSalesman()
}
