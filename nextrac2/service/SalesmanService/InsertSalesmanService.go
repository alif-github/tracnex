package SalesmanService

import (
	"database/sql"
	"fmt"
	"net/http"
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
	"nexsoft.co.id/nextrac2/service/MasterDataService/DistrictService"
	"nexsoft.co.id/nextrac2/service/MasterDataService/ProvinceService"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"time"
)

func (input salesmanService) InsertSalesman(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertSalesman"
		inputStruct in.SalesmanRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertSalesman, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_SALESMAN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanService) doInsertSalesman(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName           = "InsertSalesmanService.go"
		funcName           = "doInsertSalesman"
		descScope          = "Scope for salesman on id "
		scopeName          = constanta.SalesmanDataScope
		inputStruct        = inputStructInterface.(in.SalesmanRequest)
		dataPersonTitleMDB out.PersonTitleResponse
		scopeLimit         map[string]interface{}
		salesmanID         int64
		dataScopeID        int64
		provinceOnDB       repository.ProvinceModel
		DistrictOnDB       repository.DistrictModel
		db                 = serverconfig.ServerAttribute.DBConnection
		internalToken 	   = resource_common_service.GenerateInternalToken(constanta.ResourceMasterData, 0, contextModel.AuthAccessTokenModel.ClientID, constanta.Issue, constanta.DefaultApplicationsLanguage)
	)

	scopeLimit, err = input.validateDataScopeSalesman(contextModel)
	if err.Error != nil {
		return
	}

	dataPersonTitleMDB, err = master_data_dao.ViewDetailPersonTitleFromMasterData(int(inputStruct.PersonTitleID), contextModel)
	if err.Error != nil {
		if err.Error.Error() == "E-4-MAD-SRV-004" {
			err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.PersonTitleID)
		}
		return
	}

	if dataPersonTitleMDB.ID < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.PersonTitleID)
		return
	}

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

	salesmanModel := repository.SalesmanModel{
		PersonTitleID: sql.NullInt64{Int64: inputStruct.PersonTitleID},
		PersonTitle:   sql.NullString{String: dataPersonTitleMDB.Title},
		Sex:           sql.NullString{String: inputStruct.Sex},
		Nik:           sql.NullString{String: inputStruct.Nik},
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
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	salesmanID, err = dao.SalesmanDAO.InsertSalesman(tx, salesmanModel)
	if err.Error != nil {
		err = checkDuplicateError(err)
		return
	}

	dataScopeModel := repository.DataScopeModel{
		Scope:         sql.NullString{String: fmt.Sprintf(`%s:%s`, scopeName, strconv.Itoa(int(salesmanID)))},
		Description:   sql.NullString{String: descScope + strconv.Itoa(int(salesmanID))},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	dataScopeID, err = dao.DataScopeDAO.InsertDataScope(tx, dataScopeModel)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.SalesmanDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: salesmanID},
	}, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.DataScopeDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: dataScopeID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanService) validateInsert(inputStruct *in.SalesmanRequest) errorModel.ErrorModel {
	return inputStruct.ValidationInsertSalesman()
}
