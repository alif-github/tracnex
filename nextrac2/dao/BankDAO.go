package dao
//
//import (
//	"database/sql"
//	"nexsoft.co.id/nextrac2/dto/in"
//	"nexsoft.co.id/nextrac2/model/errorModel"
//	"nexsoft.co.id/nextrac2/repository"
//	"strconv"
//	"time"
//)
//
//type bankDAO struct {
//	AbstractDAO
//}
//
//var BankDAO = bankDAO{}.New()
//
//func (input bankDAO) New() (output bankDAO) {
//	output.FileName = "BankDAO.go"
//	output.TableName = "bank"
//	output.ElasticSearchIndex = "bank"
//	return
//}
//
//func (input bankDAO) InsertBank(tx *sql.Tx, userParam repository.BankModel, timeNow time.Time) (id int64, err errorModel.ErrorModel) {
//	funcName := "InsertBank"
//	query := "INSERT INTO bank(name, created_by, created_at, created_client, updated_by, updated_client, updated_at) " +
//		"VALUES ($1, $2, $3, $4, $5, $6, $7) returning id"
//
//	param := []interface{}{userParam.SocketID.String, userParam.CreatedBy.Int64, timeNow, userParam.CreatedClient.String, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, timeNow}
//
//	errorS := tx.QueryRow(query, param...).Scan(&id)
//	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
//		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
//		return
//	}
//
//	idStr := strconv.Itoa(int(id))
//	input.DoInsertAtElasticSearch(idStr, repository.BankElasticModel{
//		ID:        id,
//		SocketID:      userParam.SocketID.String,
//		Status:    "A",
//		CreatedBy: userParam.CreatedBy.Int64,
//	})
//	err = errorModel.GenerateNonErrorModel()
//	return
//}
//
//func (input bankDAO) GetBankForUpdate(db *sql.Tx, userParam repository.BankModel) (result repository.BankModel, err errorModel.ErrorModel) {
//	funcName := "GetBankForUpdate"
//	query := "SELECT id, updated_at, created_by FROM bank WHERE id = $1 AND deleted = FALSE "
//
//	param := []interface{}{userParam.ID.Int64}
//
//	if userParam.CreatedBy.Int64 > 0 {
//		query += " AND created_by = $2 "
//		param = append(param, userParam.CreatedBy.Int64)
//	}
//
//	query += " FOR UPDATE "
//
//	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.UpdatedAt, &result.CreatedBy)
//	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
//		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
//		return
//	}
//
//	err = errorModel.GenerateNonErrorModel()
//	return
//}
//
//func (input bankDAO) ViewBank(db *sql.DB, userParam repository.BankModel, isCheckStatus bool) (result repository.BankModel, err errorModel.ErrorModel) {
//	funcName := "ViewBank"
//	query := "SELECT id, name, status, created_by, updated_at FROM bank WHERE id = $1 AND deleted = FALSE "
//
//	param := []interface{}{userParam.ID.Int64}
//
//	if isCheckStatus {
//		query += "AND status = 'A' "
//	}
//
//	if userParam.CreatedBy.Int64 != 0 {
//		query += "AND created_by = $2 "
//		param = append(param, userParam.CreatedBy.Int64)
//	}
//
//	errorS := db.QueryRow(query, param...).Scan(&result.ID, &result.SocketID, &result.Status, &result.CreatedBy, &result.UpdatedAt)
//	if errorS != nil && errorS.Error() != "sql: no rows in result set" {
//		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
//		return
//	}
//	err = errorModel.GenerateNonErrorModel()
//	return
//}
//
//func (input bankDAO) UpdateBank(tx *sql.Tx, userParam repository.BankModel, timeNow time.Time) (err errorModel.ErrorModel) {
//	funcName := "UpdateBank"
//	query := "UPDATE bank set name = $1, status = $2, updated_by = $3, updated_client = $4, updated_at = $5 WHERE id = $6 "
//
//	param := []interface{}{userParam.SocketID.String, userParam.Status.String, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, timeNow, userParam.ID.Int64}
//
//	stmt, errorS := tx.Prepare(query)
//	if errorS != nil {
//		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
//		return
//	}
//
//	_, errorS = stmt.Exec(param...)
//
//	if errorS != nil {
//		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
//		return
//	}
//
//	id := strconv.Itoa(int(userParam.ID.Int64))
//	input.doUpdateAtElasticSearch(id, repository.BankElasticModel{
//		ID:        userParam.ID.Int64,
//		SocketID:      userParam.SocketID.String,
//		Status:    userParam.Status.String,
//		CreatedBy: userParam.CreatedBy.Int64,
//	})
//
//	err = errorModel.GenerateNonErrorModel()
//	return
//}
//
//func (input bankDAO) DeleteBank(tx *sql.Tx, userParam repository.BankModel, timeNow time.Time) (err errorModel.ErrorModel) {
//	funcName := "DeleteBank"
//	query := "UPDATE bank set deleted = TRUE, updated_by = $1, updated_client = $2, updated_at = $3 WHERE id = $4 "
//
//	param := []interface{}{userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, timeNow, userParam.ID.Int64}
//
//	stmt, errorS := tx.Prepare(query)
//	if errorS != nil {
//		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
//		return
//	}
//
//	_, errorS = stmt.Exec(param...)
//
//	if errorS != nil {
//		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
//		return
//	}
//
//	id := strconv.Itoa(int(userParam.ID.Int64))
//	input.doDeleteAtElasticSearch(id)
//
//	err = errorModel.GenerateNonErrorModel()
//	return
//}
//
//func (input bankDAO) GetListBank(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, isCheckStatus bool, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
//	if len(searchBy) > 0 {
//		result, err = input.listBankFromElastic(db, userParam, searchBy, isCheckStatus, createdBy)
//	} else {
//		result, err = input.getListBankFromDB(db, userParam, nil, isCheckStatus, createdBy)
//	}
//	return
//}
//
//func (input bankDAO) GetCountBank(db *sql.DB, searchBy []in.SearchByParam, isCheckStatus bool, createdBy int64) (result int, err errorModel.ErrorModel) {
//	if searchBy != nil {
//		return input.getCountBankFromElastic(searchBy, isCheckStatus, createdBy)
//	} else {
//		return input.getCountBankFromDB(db, nil, isCheckStatus, createdBy)
//	}
//}
//
////func (input bankDAO) listBankFromElastic(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, isCheckStatus bool, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
////	query := "SELECT " +
////		"	id, name, status, created_by, updated_at  " +
////		"FROM " +
////		"	bank "
////
////	result, err = GetListDataDAO.getlis(db, userParam, searchBy, query, input.ElasticSearchIndex, func(rows *sql.Rows) (interface{}, error) {
////		var temp repository.BankModel
////		errorS := rows.Scan(&temp.ID, &temp.SocketID, &temp.Status, &temp.CreatedBy, &temp.UpdatedAt)
////		return temp, errorS
////	}, DefaultFieldMustCheck{
////		ID:        FieldStatus{FieldName: "id"},
////		Deleted:   FieldStatus{FieldName: "deleted"},
////		Status:    FieldStatus{FieldName: "status", IsCheck: isCheckStatus},
////		CreatedBy: FieldStatus{FieldName: "created_by", Value: createdBy},
////	})
////
////	err = errorModel.GenerateNonErrorModel()
////	return
////}
//
//func (input bankDAO) getListBankFromDB(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, isCheckStatus bool, createdBy int64) (result []interface{}, err errorModel.ErrorModel) {
//	query :=
//		" SELECT " +
//			"	id, name, status, created_by, updated_at " +
//			"FROM " +
//			input.TableName + " "
//
//	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, query, userParam, searchByParam,
//		func(rows *sql.Rows) (interface{}, error) {
//			var temp repository.BankModel
//			errorS := rows.Scan(&temp.ID, &temp.SocketID, &temp.Status, &temp.CreatedBy, &temp.UpdatedAt)
//			return temp, errorS
//		}, "", DefaultFieldMustCheck{}.GetDefaultField(isCheckStatus, createdBy))
//}
//
//func (input bankDAO) getCountBankFromDB(db *sql.DB, searchByParam []in.SearchByParam, isCheckStatus bool, createdBy int64) (int, errorModel.ErrorModel) {
//	return GetListDataDAO.GetCountDataWithDefaultMustCheck(db, input.TableName, searchByParam, "", DefaultFieldMustCheck{}.GetDefaultField(isCheckStatus, createdBy))
//}
//
//func (input bankDAO) getCountBankFromElastic(searchBy []in.SearchByParam, isCheckStatus bool, createdBy int64) (result int, err errorModel.ErrorModel) {
//	funcName := "getCountBankFromElastic"
//
//	_, count, errs := input.doGetCountDataFromElastic(searchBy, isCheckStatus, createdBy)
//	if errs != nil {
//		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
//		return
//	}
//
//	result = count
//	err = errorModel.GenerateNonErrorModel()
//	return
//}
