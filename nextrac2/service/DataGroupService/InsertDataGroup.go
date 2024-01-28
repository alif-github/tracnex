package DataGroupService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"time"
)

func (input dataGroupService) InsertDataGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertDataGroup"
		inputStruct in.DataGroupRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertDataGroup, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- Function Additional
	})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_DATA_GROUP_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupService) doInsertDataGroup(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName    = "doInsertDataGroup"
		inputStruct = inputStructInterface.(in.DataGroupRequest)
		dataGroupID int64
		countScope  int
		scope       map[string][]string
	)

	countScope, err = dao.DataScopeDAO.CheckIsScopeValid(serverconfig.ServerAttribute.DBConnection, inputStruct.Scope)
	if countScope != len(inputStruct.Scope) {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.DataScope)
		return
	}

	scope = service.GenerateHashMapPermissionAndDataScope(inputStruct.Scope, true, true)
	err = input.validateScopeForInsert(scope)
	if err.Error != nil {
		return
	}

	if len(scope[constanta.DistrictDataScope]) > 0 && !common.ValidateStringContainInStringArray(scope[constanta.DistrictDataScope], "all") {
		if len(scope[constanta.ProvinceDataScope]) > 0 && !common.ValidateStringContainInStringArray(scope[constanta.ProvinceDataScope], "all") {
			countScope, err = dao.DistrictDAO.GetCountProvinceOnDistrict(serverconfig.ServerAttribute.DBConnection, convertArrStringToArrInt(scope[constanta.DistrictDataScope]))
			if countScope != len(scope[constanta.DistrictDataScope]) {
				err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ProvinceID)
				return
			}
		}
	}

	dataGroup := repository.DataGroupModel{
		GroupID:       sql.NullString{String: inputStruct.GroupID},
		Description:   sql.NullString{String: inputStruct.Description},
		Scope:         sql.NullString{String: util.StructToJSON(scope)},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	dataGroupID, err = dao.DataGroupDAO.InsertDataGroup(tx, dataGroup)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.DataGroupDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: dataGroupID},
	})

	output = dataGroupID
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupService) validateScopeForInsert(scope map[string][]string) (err errorModel.ErrorModel) {
	for key, value := range scope {
		if common.ValidateStringContainInStringArray(scope[key], "all") && len(value) > 1 {
			err = errorModel.GenerateFieldFormatWithRuleError(input.FileName, "validateScopeForInsert", "SCOPE_INSERT_VALIDATION", constanta.Scope, "")
			return
		}
	}

	return
}

func (input dataGroupService) validateScopeData(scope map[string][]string) (err errorModel.ErrorModel) {

	if len(scope[constanta.ProvinceDataScope]) > 0 && !common.ValidateStringContainInStringArray(scope[constanta.ProvinceDataScope], "all") {
		scope[constanta.DistrictDataScope] = []string{"all"}
	} else if common.ValidateStringContainInStringArray(scope[constanta.ProvinceDataScope], "all") {
		scope[constanta.DistrictDataScope] = []string{"all"}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupService) validateInsert(inputStruct *in.DataGroupRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertDataGroup()
}

func convertArrStringToArrInt(arrString []string) []int {
	var arrInt []int
	for _, value := range arrString {
		id, _ := strconv.Atoi(value)
		arrInt = append(arrInt, id)
	}

	return arrInt
}
