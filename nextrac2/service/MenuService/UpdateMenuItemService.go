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

func (input menuService) UpdateMenuItem(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "UpdateMenuItem"
	var inputStruct in.MenuRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateMenuItem)
	if err.Error != nil {return}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateMenuItem, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {return}

	output.Status = out.StatusResponse {
		Code: 		util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: 	GenerateI18NMessage("SUCCESS_UPDATE_ITEM_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input menuService) doUpdateMenuItem(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	fileName := "UpdateMenuItemService.go"
	funcName := "doUpdateMenuItem"

	inputStruct := inputStructInterface.(in.MenuRequest)
	var menuItemOnDB repository.MenuModel
	tableName := constanta.TableMenuItem

	menuItemOnDB, err  = dao.MenuDAO.GetMenuForUpdate(serverconfig.ServerAttribute.DBConnection, repository.MenuModel {
		ID: 		sql.NullInt64{Int64: inputStruct.ID},
		CreatedBy: 	sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	}, tableName)
	if err.Error != nil {return}

	if menuItemOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, "Menu Item")
		return
	}

	if menuItemOnDB.UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(fileName, funcName, "Menu Item")
		return
	}

	err = input.doCheckParentMenu(inputStruct.ServiceMenuID, constanta.TableMenuService)
	if err.Error != nil {return}

	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, tableName, menuItemOnDB.ID.Int64, contextModel.LimitedByCreatedBy)...)

	err = dao.MenuDAO.UpdateMenu(tx, repository.MenuModel {
		ID: 				sql.NullInt64{Int64: menuItemOnDB.ID.Int64},
		ServiceMenuID: 		sql.NullInt64{Int64: inputStruct.ServiceMenuID},
		Name: 				sql.NullString{String: inputStruct.Name},
		EnName: 			sql.NullString{String: inputStruct.EnName},
		Sequence: 			sql.NullInt64{Int64: inputStruct.Sequence},
		IconName: 			sql.NullString{String: inputStruct.IconName},
		Background: 		sql.NullString{String: inputStruct.Background},
		AvailableAction: 	sql.NullString{String: inputStruct.AvailableAction},
		MenuCode: 			sql.NullString{String: inputStruct.MenuCode},
		Status: 			sql.NullString{String: inputStruct.Status},
		Url: 				sql.NullString{String: inputStruct.Url},
		CreatedBy: 			sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
		UpdatedBy: 			sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: 		sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt: 			sql.NullTime{Time: timeNow},
	}, tableName)

	if err.Error != nil {return}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input menuService) validateUpdateMenuItem(inputStruct *in.MenuRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateMenuItem()
}