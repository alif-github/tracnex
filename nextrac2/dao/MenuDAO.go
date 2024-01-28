package dao

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type menuDAO struct {
	AbstractDAO
}

var MenuDAO = menuDAO{}.New()
var sysUserTableName = "menu_item"
var sysAdminTableName = "menu_sysadmin"

func (input menuDAO) New() (output menuDAO) {
	output.FileName = "MenuDAO.go"
	return
}

func (input menuDAO) ViewMenuListAdmin(db *sql.DB) (result []out.ParentMenu, err errorModel.ErrorModel) {
	funcName := "ViewMenuList"
	var jsonResult sql.NullString

	query :=
		"SELECT " +
			"	json_agg(a) as menu_parent " +
			"FROM " +
			"	(SELECT " +
			"		parent.id, parent.name, parent.en_name, " +
			"		parent.sequence, parent.icon_name, parent.background, " +
			"		parent.available_action, parent.menu_code " +
			"	FROM " +
			"		menu_sysadmin parent " +
			"	WHERE " +
			"		parent.status = 'A' AND " +
			"		parent.deleted = FALSE ORDER BY sequence)a"

	results := db.QueryRow(query)
	errS := results.Scan(&jsonResult)

	if errS != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errS)
	}

	_ = json.Unmarshal([]byte(jsonResult.String), &result)

	return
}

func (input menuDAO) ViewMenuList(db *sql.DB) (result []out.ParentMenu, err errorModel.ErrorModel) {
	funcName := "ViewMenuList"
	var jsonResult sql.NullString

	query :=
		"SELECT " +
			"	json_agg(a) as menu_parent " +
			"FROM " +
			"	(SELECT " +
			"		parent.id, parent.name, parent.en_name, " +
			"		parent.sequence, parent.icon_name, parent.background, " +
			"		parent.available_action, parent.menu_code, " +
			"		(SELECT " +
			"			json_agg(a) as menu_service " +
			"		FROM " +
			"			(SELECT " +
			"				service.id, service.parent_menu_id, service.name, " +
			"				service.en_name, service.sequence, service.icon_name, " +
			"				service.menu_code, service.background, service.available_action, " +
			"				(SELECT " +
			"					json_agg(a) as menu_item  " +
			"				FROM " +
			"					(SELECT " +
			"						id, service_menu_id, name, " +
			"						en_name, sequence, icon_name, " +
			"						background, menu_code, url, " +
			"						available_action " +
			"					FROM " +
			"						menu_item " +
			"					WHERE " +
			"						service_menu_id = service.id AND " +
			"						status = 'A' AND " +
			"						deleted = FALSE ORDER BY sequence)a) " +
			"			FROM " +
			"				service_menu service " +
			"			WHERE " +
			"				parent_menu_id = parent.id AND " +
			"				status = 'A' AND " +
			"				deleted = FALSE ORDER BY sequence)a) " +
			"	FROM " +
			"		parent_menu parent " +
			"	WHERE " +
			"		parent.status = 'A' AND " +
			"		parent.deleted = FALSE ORDER BY sequence)a"

	results := db.QueryRow(query)
	errS := results.Scan(&jsonResult)

	if errS != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errS)
	}

	_ = json.Unmarshal([]byte(jsonResult.String), &result)

	return
}

func (input menuDAO) ViewParentMenuList(db *sql.DB, isSysadmin bool) (result []out.MenuList, err errorModel.ErrorModel) {
	funcName := "ViewParentMenuList"
	var jsonResult sql.NullString
	var tableName string

	if isSysadmin {
		tableName = "menu_sysadmin"
	} else {
		tableName = "parent_menu"
	}

	query :=
		"	SELECT " +
			"	json_agg(parent_menu) " +
			"FROM " +
			"	(SELECT " +
			"		id, name, en_name, " +
			"		menu_code, available_action " +
			"	FROM " +
			"		" + tableName + "" +
			"	WHERE " +
			"		status = 'A' AND deleted = FALSE" +
			"		ORDER BY sequence) " +
			"	AS " +
			"		parent_menu"

	results := db.QueryRow(query)
	errS := results.Scan(&jsonResult)

	if errS != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errS)
	}

	errS = json.Unmarshal([]byte(jsonResult.String), &result)
	if errS != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errS)
	}
	return
}

func (input menuDAO) ViewServiceMenuList(db *sql.DB, parentMenuID int) (result []out.MenuList, err errorModel.ErrorModel) {
	funcName := "ViewServiceMenuList"
	var jsonResult sql.NullString
	query :=
		"	SELECT " +
			"	json_agg(service_menu) " +
			"FROM " +
			"	(SELECT " +
			"		id, name, en_name, menu_code, available_action " +
			"	FROM " +
			"		service_menu" +
			"	WHERE " +
			"		parent_menu_id = $1 AND" +
			"		status = 'A' AND deleted = FALSE" +
			"		ORDER BY sequence) " +
			"	AS " +
			"		service_menu"

	results := db.QueryRow(query, parentMenuID)
	errS := results.Scan(&jsonResult)

	if errS != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errS)
	}

	_ = json.Unmarshal([]byte(jsonResult.String), &result)

	return
}

func (input menuDAO) ViewMenuItemList(db *sql.DB, serviceMenuID int) (result []out.NewMenuItemList, err errorModel.ErrorModel) {
	funcName := "ViewMenuItemList"
	var jsonResult sql.NullString
	query :=
		"	SELECT " +
			"	json_agg(menu_item) " +
			"FROM " +
			"	(SELECT " +
			"		id, name, en_name, menu_code, available_action, updated_at " +
			"	FROM " +
			"		menu_item " +
			"	WHERE " +
			"		service_menu_id = $1 AND " +
			"		status = 'A' AND deleted = FALSE" +
			"		ORDER BY sequence) " +
			"	AS " +
			"		menu_item "

	results := db.QueryRow(query, serviceMenuID)
	errS := results.Scan(&jsonResult)

	if errS != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errS)
	}

	_ = json.Unmarshal([]byte(jsonResult.String), &result)

	return
}

func (input menuDAO) GetMenuForUpdate(db *sql.DB, userParam repository.MenuModel, tableNameInput string) (result repository.MenuModel, err errorModel.ErrorModel) {
	fileName := "MenuDAO.go"
	funcName := "GetMenuSysAdmin"
	var tableName string

	switch tableNameInput {
	case constanta.TableNameMenuSysAdmin:
		tableName = constanta.TableNameMenuSysAdmin
	case constanta.TableNameMenuParent:
		tableName = constanta.TableNameMenuParent
	case constanta.TableMenuService:
		tableName = constanta.TableMenuService
	case constanta.TableMenuItem:
		tableName = constanta.TableMenuItem
	default:
		tableName = constanta.TableNameMenuSysAdmin
	}

	query := " SELECT " +
		"	id, updated_at " +
		"	FROM " + tableName + " " +
		"	WHERE id = $1 AND deleted = FALSE "

	var params []interface{}
	params = append(params, userParam.ID.Int64)

	if userParam.CreatedBy.Int64 > 0 {
		query += " AND created_by = $2 "
		params = append(params, userParam.CreatedBy.Int64)
	}

	query += " FOR UPDATE"

	results := db.QueryRow(query, params...)
	errS := results.Scan(&result.ID, &result.UpdatedAt)

	if errS != nil && errS.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input menuDAO) UpdateMenu(db *sql.Tx, userParam repository.MenuModel, tableNameInput string) (err errorModel.ErrorModel) {
	funcName := "UpdateMenu"
	var tableName string

	switch tableNameInput {
	case constanta.TableNameMenuSysAdmin:
		tableName = constanta.TableNameMenuSysAdmin
	case constanta.TableNameMenuParent:
		tableName = constanta.TableNameMenuParent
	case constanta.TableMenuService:
		tableName = constanta.TableMenuService
	case constanta.TableMenuItem:
		tableName = constanta.TableMenuItem
	default:
		tableName = constanta.TableNameMenuSysAdmin
	}

	query := "UPDATE " + tableName + " " +
		"SET " +
		"name = $1, en_name = $2, sequence = $3, " +
		"available_action = $4, menu_code = $5, status = $6, " +
		"updated_by = $7, updated_client = $8, updated_at = $9, " +
		"icon_name = $10, background = $11"

	params := []interface{}{
		userParam.Name.String, userParam.EnName.String, userParam.Sequence.Int64,
		userParam.AvailableAction.String, userParam.MenuCode.String, userParam.Status.String,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
	}

	if userParam.IconName.String != "" {
		params = append(params, userParam.IconName.String)
	} else {
		params = append(params, nil)
	}

	if userParam.Background.String != "" {
		params = append(params, userParam.Background.String)
	} else {
		params = append(params, nil)
	}

	if tableName == constanta.TableMenuService {
		query += ", parent_menu_id = $12 WHERE " +
			"id = $13 "
		params = append(params, userParam.ParentMenuID.Int64)
	} else if tableName == constanta.TableMenuItem {
		query += ", service_menu_id = $12, url = $13 WHERE " +
			"id = $14 "
		params = append(params, userParam.ServiceMenuID.Int64, userParam.Url.String)
	} else {
		query += " WHERE " +
			"id = $12 "
	}

	params = append(params, userParam.ID.Int64)

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(params...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	return
}

func (input menuDAO) ViewMenuListForPermission(db *sql.DB) (result []out.ParentMenuList, err errorModel.ErrorModel) {
	funcName := "ViewMenuListForPermission"
	var jsonResult sql.NullString

	query :=
		"SELECT " +
			"	json_agg(a) as menu_parent " +
			"FROM " +
			"	(SELECT " +
			"		parent.id, parent.name, parent.en_name, " +
			"		parent.available_action, parent.menu_code, " +
			"		parent.icon_name, " +
			"		(SELECT " +
			"			json_agg(a) as items " +
			"		FROM " +
			"			(SELECT " +
			"				service.id, service.name, service.en_name, " +
			"				service.menu_code, service.available_action, " +
			"				service.icon_name, " +
			"				(SELECT " +
			"					json_agg(a) as items " +
			"				FROM " +
			"					(SELECT " +
			"						id, name, en_name, " +
			"						sequence, menu_code, available_action, " +
			"						icon_name " +
			"					FROM " +
			"						menu_item " +
			"					WHERE " +
			"						service_menu_id = service.id AND " +
			"						status = 'A' AND " +
			"						deleted = FALSE ORDER BY sequence)a) " +
			"			FROM " +
			"				service_menu service " +
			"			WHERE " +
			"				parent_menu_id = parent.id AND " +
			"				status = 'A' AND " +
			"				deleted = FALSE ORDER BY sequence)a) " +
			"	FROM " +
			"		parent_menu parent " +
			"	WHERE " +
			"		parent.status = 'A' AND " +
			"		parent.deleted = FALSE ORDER BY sequence)a"

	results := db.QueryRow(query)
	errS := results.Scan(&jsonResult)

	if errS != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errS)
	}

	_ = json.Unmarshal([]byte(jsonResult.String), &result)

	return
}

func (input menuDAO) ViewMenuListAdminForPermission(db *sql.DB) (result []out.ParentMenuList, err errorModel.ErrorModel) {
	var (
		funcName   = "ViewMenuListAdminForPermission"
		jsonResult sql.NullString
		query      string
	)

	query = fmt.Sprintf(`SELECT 
			json_agg(a) as menu_parent 
			FROM 
				(SELECT 
					parent.id, parent.name, parent.en_name, 
					parent.available_action, parent.menu_code, parent.icon_name 
				FROM 
					menu_sysadmin parent 
				WHERE 
					parent.status = 'A' AND 
					parent.deleted = FALSE ORDER BY sequence)a`)

	results := db.QueryRow(query)
	errS := results.Scan(&jsonResult)

	if errS != nil {
		return result, errorModel.GenerateInternalDBServerError(input.FileName, funcName, errS)
	}

	_ = json.Unmarshal([]byte(jsonResult.String), &result)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input menuDAO) GetTableNameWithMenuCode(db *sql.DB, userParam repository.MenuModel, isAdmin bool) (result string, err errorModel.ErrorModel) {
	funcName := "GetTableNameWithMenuCode"
	var tableName string

	if isAdmin {
		tableName = sysAdminTableName
	} else {
		tableName = sysUserTableName
	}

	query := fmt.Sprintf(
		`SELECT 
		table_name
	FROM %s
	WHERE 
		menu_code = $1 AND status = 'A' AND deleted = FALSE `, tableName)

	params := []interface{}{userParam.MenuCode.String}

	dbError := db.QueryRow(query, params...).Scan(&result)

	if dbError != nil && dbError.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	return
}

func (input menuDAO) ResetMenuRecursive(db *sql.DB) (err errorModel.ErrorModel) {
	var (
		fileName = input.FileName
		funcName = "ResetMenuRecursive"
	)

	query := fmt.Sprintf(`
		with parent_menu_update as (
			update parent_menu set status = 'N' 
			returning parent_menu.id 
		), service_menu_update as (
			update service_menu set status = 'N' 
			where parent_menu_id in (select id from parent_menu_update) 
			returning service_menu.id
		)
		update menu_item set status = 'N' 
		where service_menu_id in (select id from service_menu_update);
		`)

	stmt, errs := db.Prepare(query)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
		return
	}

	_, errs = stmt.Exec()
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input menuDAO) UpdateMenuRecursive(db *sql.DB, status, operatorMaster, operatorAdmin string) (err errorModel.ErrorModel) {
	var (
		fileName   = input.FileName
		funcName   = "UpdateMenuRecursive"
		codeMaster = "master"
		codeAdmin  = "admin"
	)

	query := fmt.Sprintf(`
		with parent_menu_update as (
			update parent_menu set status = '%s' 
			where menu_code %s '%s' or menu_code %s '%s' 
			returning parent_menu.id 
		), service_menu_update as (
			update service_menu set status = '%s' 
			where parent_menu_id in (select id from parent_menu_update) 
			returning service_menu.id
		)
		update menu_item set status = '%s' 
		where service_menu_id in (select id from service_menu_update);
		`,
		status, operatorAdmin, codeAdmin,
		operatorMaster, codeMaster, status, status)

	stmt, errs := db.Prepare(query)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
		return
	}

	_, errs = stmt.Exec()
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
