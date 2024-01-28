package DataGroupService

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
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input dataGroupService) DeleteDataGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteDataGroup"
		inputStruct in.DataGroupRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDelete)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteDataGroup, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- Function Additional
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupService) doDeleteDataGroup(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName    = "doInsertDataGroup"
		inputStruct = inputStructInterface.(in.DataGroupRequest)
		dataGroupDB repository.DataGroupModel
	)

	dataGroup := repository.DataGroupModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	dataGroup.CreatedBy.Int64 = contextModel.LimitedByCreatedBy
	dataGroupDB, err = dao.DataGroupDAO.GetDataGroupForDelete(tx, dataGroup)
	if err.Error != nil {
		return
	}

	if dataGroupDB.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ID)
		return
	}

	if dataGroupDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(input.FileName, funcName, constanta.DataGroup)
		return
	}

	if dataGroupDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(input.FileName, funcName, constanta.DataGroup)
		return
	}

	// ----------- Update for delete
	encodedStr, errorS := service.RandToken(10)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	dataGroup.GroupID.String = dataGroupDB.GroupID.String + encodedStr
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.DataGroupDAO.TableName, dataGroupDB.ID.Int64, 0)...)
	err = dao.DataGroupDAO.DeleteDataGroup(tx, dataGroup)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupService) validateDelete(inputStruct *in.DataGroupRequest) errorModel.ErrorModel {
	return inputStruct.ValidateDeleteDataGroup()
}
