package dao

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type fileUploadDAO struct {
	AbstractDAO
}

var FileUploadDAO = fileUploadDAO{}.New()

func (input fileUploadDAO) New() (output fileUploadDAO) {
	output.FileName = "FileUpload.go"
	output.TableName = "file_upload"
	return
}

func (input fileUploadDAO) InsertFileUploadInfoForBacklog(tx *sql.Tx, userParam repository.FileUpload) (id int64, err errorModel.ErrorModel) {
	funcName := "InsertFileUploadInfo"
	var param []interface{}

	query := fmt.Sprintf(`INSERT INTO %s(
		file_name, category, connector, 
		parent_id, host, path, 
		created_by, created_client, created_at, 
		updated_by, updated_client, updated_at) 
		VALUES ($1, $2, $3, 
		$4, $5, $6, 
		$7, $8, $9, 
		$10, $11, $12) RETURNING id`, input.TableName)

	param = append(param,
		userParam.FileName.String, userParam.Category.String, userParam.Konektor.String,
		userParam.ParentID.Int64, userParam.Host.String, userParam.Path.String,
		userParam.CreatedBy.Int64, userParam.CreatedClient.String, userParam.CreatedAt.Time,
		userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time)

	results := tx.QueryRow(query, param...)

	errs := results.Scan(&id)
	if errs != nil && errs.Error() != sql.ErrNoRows.Error() {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) UpdateFileUpload(db *sql.Tx, userParam repository.FileUpload) (err errorModel.ErrorModel) {
	var (
		funcName = "UpdateFileUpload"
		query    string
		params   []interface{}
	)

	query = fmt.Sprintf(`
		UPDATE %s SET 
		updated_by = $1, updated_client = $2, updated_at = $3, 
		category = $4, connector = $5, parent_id = $6
		WHERE id = $7 `,
		input.TableName)

	params = append(params, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time)

	if userParam.Category.String != "" {
		params = append(params, userParam.Category.String)
	} else {
		params = append(params, nil)
	}

	if userParam.Konektor.String != "" {
		params = append(params, userParam.Konektor.String)
	} else {
		params = append(params, nil)
	}

	if userParam.ParentID.Int64 > 0 {
		params = append(params, userParam.ParentID.Int64)
	} else {
		params = append(params, nil)
	}

	params = append(params, userParam.ID.Int64)

	stmt, errs := db.Prepare(query)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	_, errs = stmt.Exec(params...)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) GetCountFileUpload(db *sql.DB, userParam repository.FileUpload) (count int64, err errorModel.ErrorModel) {
	var (
		funcName = "GetCountFileUpload"
		params   []interface{}
		results  *sql.Row
		dbError  error
		query    string
	)

	query = fmt.Sprintf(`SELECT count(id) FROM %s WHERE parent_id = $1 AND deleted = FALSE`, input.TableName)

	params = []interface{}{userParam.ParentID.Int64}
	results = db.QueryRow(query, params...)
	dbError = results.Scan(&count)

	if dbError != nil && dbError != sql.ErrNoRows {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) GetFileUpload(db *sql.DB, userParam repository.FileUpload) (result []repository.FileUpload, err errorModel.ErrorModel) {
	var (
		funcName = "GetFileUpload"
		params   []interface{}
		rows     *sql.Rows
		dbError  error
		query    string
	)

	query = fmt.Sprintf(`SELECT id, file_name, category, host, path, updated_at 
		FROM %s WHERE parent_id = $1 AND connector = $2 AND deleted = FALSE`, input.TableName)

	params = []interface{}{userParam.ParentID.Int64, userParam.Konektor.String}
	rows, dbError = db.Query(query, params...)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	var resultTemp []interface{}
	resultTemp, err = RowsCatchResult(rows, input.resultRowsInput)
	if err.Error != nil {
		return
	}

	for _, itemResultTemp := range resultTemp {
		result = append(result, itemResultTemp.(repository.FileUpload))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) resultRowsInput(rows *sql.Rows) (resultTemp interface{}, err errorModel.ErrorModel) {
	var (
		funcName = "resultRowsInput"
		errorS   error
		temp     repository.FileUpload
	)

	errorS = rows.Scan(&temp.ID, &temp.FileName, &temp.Category, &temp.Host, &temp.Path, &temp.UpdatedAt)
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	resultTemp = temp
	return
}

func (input fileUploadDAO) DeleteFileUpload(db *sql.Tx, userParam repository.FileUpload) (err errorModel.ErrorModel) {
	funcName := "DeleteFileUpload"

	query := fmt.Sprintf(`UPDATE %s SET deleted = TRUE, updated_at = $1 
		WHERE `, input.TableName)

	params := []interface{}{userParam.UpdatedAt.Time}

	if userParam.ParentID.Int64 > 0 {
		query += " parent_id = $2 "
		params = append(params, userParam.ParentID.Int64)
	} else {
		query += " id = $2 "
		params = append(params, userParam.ID.Int64)
	}

	if userParam.Category.String != "" {
		query += " AND category = $3 "
		params = append(params, userParam.Category.String)
	}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	_, errs = stmt.Exec(params...)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) DeleteFileUploadByID(db *sql.Tx, userParam repository.FileUpload) (err errorModel.ErrorModel) {
	funcName := "DeleteFileUpload"

	query := fmt.Sprintf(`UPDATE %s 
		SET 
		deleted = TRUE, updated_at = $1, updated_by = $2, 
		updated_client = $3 
		WHERE id = $4 `, input.TableName)

	params := []interface{}{userParam.UpdatedAt.Time, userParam.UpdatedBy.Int64, userParam.UpdatedClient.String,
		userParam.ID.Int64}

	if userParam.Category.String != "" {
		query += " AND category = $5 "
		params = append(params, userParam.Category.String)
	}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	_, errs = stmt.Exec(params...)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) DeleteFileUploadForList(db *sql.Tx, userParam repository.FileUpload) (err errorModel.ErrorModel) {
	funcName := "DeleteFileUploadForList"

	query := fmt.Sprintf(`UPDATE %s SET deleted = TRUE, updated_at = $1 
		WHERE parent_id = $2 AND connector = $3 `, input.TableName)

	params := []interface{}{
		userParam.UpdatedAt.Time,
		userParam.ParentID.Int64,
		userParam.Konektor.String,
	}

	if userParam.Category.String != "" {
		query += " AND category = $4 "
		params = append(params, userParam.Category.String)
	}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	_, errs = stmt.Exec(params...)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) UpdateTimeUpdatedAtFileUpload(db *sql.Tx, userParam repository.FileUpload) (err errorModel.ErrorModel) {
	funcName := "UpdateTimeUpdatedAtFileUpload"

	query := fmt.Sprintf(`UPDATE %s SET 
		updated_by = $1, updated_client = $2, updated_at = $3 
		WHERE parent_id = $4 AND deleted = FALSE `, input.TableName)

	params := []interface{}{userParam.UpdatedBy.Int64, userParam.UpdatedClient.String, userParam.UpdatedAt.Time,
		userParam.ParentID.Int64}
	stmt, errs := db.Prepare(query)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	_, errs = stmt.Exec(params...)
	if errs != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) GetFileUploadForReminder(db *sql.Tx, userParam repository.FileUpload) (result []repository.FileUpload, err errorModel.ErrorModel) {
	var (
		funcName = "GetFileUploadForReminder"
		params   []interface{}
		rows     *sql.Rows
		dbError  error
		query    string
	)

	query = fmt.Sprintf(`SELECT id, file_name, category, host, path, updated_at 
		FROM %s WHERE parent_id = $1 AND connector = $2 AND deleted = FALSE`, input.TableName)

	params = []interface{}{userParam.ParentID.Int64, userParam.Konektor.String}
	rows, dbError = db.Query(query, params...)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	var resultTemp []interface{}
	resultTemp, err = RowsCatchResult(rows, input.resultRowsInput)
	if err.Error != nil {
		return
	}

	for _, itemResultTemp := range resultTemp {
		result = append(result, itemResultTemp.(repository.FileUpload))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) GetFileUploadForDelete(db *sql.DB, userParam repository.FileUpload) (results []repository.FileUpload, err errorModel.ErrorModel) {
	var (
		funcName = "GetFileUploadForDelete"
		index    = 2
	)

	query := fmt.Sprintf(`SELECT id, parent_id, updated_at, updated_by
		FROM %s 
		WHERE deleted = FALSE AND connector = $1`,
		input.TableName)

	params := []interface{}{
		userParam.Konektor.String,
	}

	if userParam.ParentID.Int64 > 0 {
		query += fmt.Sprintf(" AND parent_id = $%d ", index)
		params = append(params, userParam.ParentID.Int64)
		index++
	}
	if userParam.ID.Int64 > 0 {
		query += fmt.Sprintf(" AND id = $%d ", index)
		params = append(params, userParam.ID.Int64)
		index++
	}

	rows, errorDB := db.Query(query, params...)
	if errorDB != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
		return
	}

	if rows != nil {
		defer func() {
			errorDB = rows.Close()
			if errorDB != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
				return
			}
		}()
		for rows.Next() {
			var result repository.FileUpload
			errorDB = rows.Scan(
				&result.ID,
				&result.ParentID,
				&result.UpdatedAt,
				&result.UpdatedBy)
			if errorDB != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
				return
			}
			results = append(results, result)
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorDB)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) GetFileUploadWithCategory(db *sql.DB, userParam repository.FileUpload) (result []repository.FileUpload, err errorModel.ErrorModel) {
	var (
		funcName = "GetFileUploadWithCategory"
		params   []interface{}
		rows     *sql.Rows
		dbError  error
		query    string
	)

	query = fmt.Sprintf(`SELECT id, file_name, category, host, path, 
		updated_at 
		FROM %s 
		WHERE 
		parent_id = $1 AND connector = $2 AND category = $3 AND 
		deleted = FALSE `, input.TableName)

	params = []interface{}{userParam.ParentID.Int64, userParam.Konektor.String, userParam.Category.String}

	if userParam.ID.Int64 > 0 {
		query += fmt.Sprintf(` AND id = $4 `)
		params = append(params, userParam.ID.Int64)
	}

	rows, dbError = db.Query(query, params...)
	if dbError != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}

	var resultTemp []interface{}
	resultTemp, err = RowsCatchResult(rows, input.resultRowsInput)
	if err.Error != nil {
		return
	}

	for _, itemResultTemp := range resultTemp {
		result = append(result, itemResultTemp.(repository.FileUpload))
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) GetFileByParentID(db *sql.DB, parentId int64) (file repository.FileUpload, err errorModel.ErrorModel) {
	funcName := "GetPhotoByParentID"
	query := "SELECT " +
		"	id, host, path FROM file_upload " +
		" WHERE parent_id = $1 AND deleted=FALSE LIMIT 1 "

	param := []interface{}{parentId}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&file.ID, &file.Host, &file.Path)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) GetListFileUploadAndJobProcess(db *sql.DB, userParam in.GetListDataDTO, searchByParam []in.SearchByParam, createdBy int64, connector string) (result []interface{}, err errorModel.ErrorModel) {
	var query string
	query = fmt.Sprintf(`
		SELECT 
		    j.job_id, 
		    j.status, 
		    ((j.counter/j.total)*100) as process_persent, 
		    CONCAT(f.host, f.path) as file_url, 
		    j.created_at, 
		    u.nt_username 
		FROM %s j
		LEFT JOIN %s f ON j.id = f.parent_id AND f.connector = '%s'
		LEFT JOIN "%s" u ON j.created_by = u.id `,
		JobProcessDAO.TableName, FileUploadDAO.TableName, connector,
		UserDAO.TableName)

	//--- Search By Param
	for i := 0; i < len(searchByParam); i++ {
		if searchByParam[i].SearchKey == "category" {
			searchByParam[i].SearchKey = "f." + searchByParam[i].SearchKey
			continue
		}
		searchByParam[i].SearchKey = "j." + searchByParam[i].SearchKey
	}

	//--- Order By
	userParam.OrderBy = "j." + userParam.OrderBy
	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, []interface{}{}, query, userParam, searchByParam,
		func(rows *sql.Rows) (interface{}, error) {
			var temp repository.GetListFileUploadJobProcess
			dbError := rows.Scan(
				&temp.JobID, &temp.Status, &temp.Progress,
				&temp.FileUrl, &temp.CreatedAt, &temp.CreatedName,
			)
			return temp, dbError
		}, "", DefaultFieldMustCheck{
			ID:        FieldStatus{FieldName: "j.id"},
			Deleted:   FieldStatus{FieldName: "j.deleted"},
			CreatedBy: FieldStatus{FieldName: "j.created_by", Value: createdBy},
		})
}

func (input fileUploadDAO) UpdateFileUploads(db *sql.Tx, userParam in.MultipartFileDTO) (err errorModel.ErrorModel) {
	funcName := "UpdateFileUploads"
	query := "UPDATE file_upload SET host = $1 WHERE id = $2"
	param := []interface{}{userParam.Host, userParam.FileID}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		fmt.Println("Error Query Azure Prepare : ", errs.Error())
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		fmt.Println("Error Query Azure Exec : ", errs.Error())
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) GetFileByParentIDAndCategory(db *sql.DB, parent_id int64, category string) (file repository.FileUpload, err errorModel.ErrorModel) {
	funcName := "GetFileByParentIDAndCategory"
	query := "SELECT " +
		"	id, host, path FROM file_upload " +
		" WHERE parent_id = $1 AND connector = $2 AND deleted=FALSE LIMIT 1 "

	param := []interface{}{parent_id, category}

	results := db.QueryRow(query, param...)
	dbError := results.Scan(&file.ID, &file.Host, &file.Path)

	if dbError != nil && dbError.Error() != "sql: no rows in result set" {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, dbError)
		return
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input fileUploadDAO) DeleteFileUploadByParentID(db *sql.Tx, parentId int64, category string, connector string) (err errorModel.ErrorModel) {
	funcName := "DeleteFileUpload"
	query := "DELETE FROM file_upload WHERE parent_id = $1 AND category = $2 AND connector = $3 "
	param := []interface{}{parentId, category, connector}

	stmt, errs := db.Prepare(query)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	_, errs = stmt.Exec(param...)
	if errs != nil {
		return errorModel.GenerateInternalDBServerError(input.FileName, funcName, errs)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
