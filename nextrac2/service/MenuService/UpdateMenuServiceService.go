package MenuService

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
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input menuService) UpdateMenuService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "UpdateMenuService"
	var inputStruct in.MenuRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateMenuService)
	if err.Error != nil {return}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateMenuService, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {return}

	output.Status = out.StatusResponse {
		Code: 		util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: 	GenerateI18NMessage("SUCCESS_UPDATE_SERVICE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input menuService) doUpdateMenuService(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	fileName := "UpdateMenuServiceService.go"
	funcName := "doUpdateMenuService"

	inputStruct := inputStructInterface.(in.MenuRequest)
	var menuServiceOnDB repository.MenuModel
	tableName := constanta.TableMenuService

	menuServiceOnDB, err  = dao.MenuDAO.GetMenuForUpdate(serverconfig.ServerAttribute.DBConnection, repository.MenuModel{
		ID: 		sql.NullInt64{Int64: inputStruct.ID},
		CreatedBy: 	sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	}, tableName)
	if err.Error != nil {return}

	if menuServiceOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, "Menu Service")
		return
	}

	if menuServiceOnDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(fileName, funcName, "Menu Service")
		return
	}

	err = input.doCheckParentMenu(inputStruct.ParentMenuID, constanta.TableNameMenuParent)
	if err.Error != nil {return}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, tableName, menuServiceOnDB.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.MenuDAO.UpdateMenu(tx, repository.MenuModel{
		ID: 				sql.NullInt64{Int64: menuServiceOnDB.ID.Int64},
		ParentMenuID: 		sql.NullInt64{Int64: inputStruct.ParentMenuID},
		Name: 				sql.NullString{String: inputStruct.Name},
		EnName: 			sql.NullString{String: inputStruct.EnName},
		Sequence: 			sql.NullInt64{Int64: inputStruct.Sequence},
		IconName: 			sql.NullString{String: inputStruct.IconName},
		Background: 		sql.NullString{String: inputStruct.Background},
		AvailableAction: 	sql.NullString{String: inputStruct.AvailableAction},
		MenuCode: 			sql.NullString{String: inputStruct.MenuCode},
		Status: 			sql.NullString{String: inputStruct.Status},
		CreatedBy: 			sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
		UpdatedBy: 			sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: 		sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt: 			sql.NullTime{Time: timeNow},
	}, tableName)

	if err.Error != nil {return}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input menuService) validateUpdateMenuService(inputStruct *in.MenuRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateMenuService()
}