package DataScopeService

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input dataScopeService) DoDeleteDataScope(tx *sql.Tx, inputModel repository.DataScopeModel, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var dataScopeOnDB repository.DataScopeModel
	dataScopeOnDB, err = dao.DataScopeDAO.GetDataScopeByScope(tx, inputModel)
	if err.Error != nil {
		return
	}

	dataScopeModel := repository.DataScopeModel{
		ID:            dataScopeOnDB.ID,
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
	}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.DataScopeDAO.TableName, dataScopeOnDB.ID.Int64, 0)...)
	err = dao.DataScopeDAO.DeleteDataScope(tx, dataScopeModel, timeNow)
	return
}
