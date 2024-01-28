package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type allowanceDAO struct {
	AbstractDAO
}

var AllowanceDAO = allowanceDAO{}.New()

func (input allowanceDAO) New() (output allowanceDAO) {
	output.FileName = "AllowanceDAO.go"
	output.TableName = "allowances"
	return
}

func (input allowanceDAO) GetByAllowanceType(db *sql.DB, allowanceType string) (result repository.Allowance, errModel errorModel.ErrorModel) {
	funcName := "GetByAllowanceType"

	query := `SELECT 
				id, allowance_name, allowance_type
			FROM ` + input.TableName + `
			WHERE 
				allowance_type = $1 AND 
				active = TRUE AND 
				deleted = FALSE`

	row := db.QueryRow(query, allowanceType)
	err := row.Scan(&result.ID, &result.AllowanceName, &result.AllowanceType)
	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input allowanceDAO) GetByIdAndAllowanceType(db *sql.DB, id int64, allowanceType string) (result repository.Allowance, errModel errorModel.ErrorModel) {
	funcName := "GetById"

	query := `SELECT 
				id, allowance_name, allowance_type
			FROM ` + input.TableName + `
			WHERE 
				id = $1 AND 
				allowance_type = $2 AND 
				active = TRUE AND 
				deleted = FALSE`

	row := db.QueryRow(query, id, allowanceType)
	err := row.Scan(&result.ID, &result.AllowanceName, &result.AllowanceType)
	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input allowanceDAO) GetAllowanceLeaveByIdAndEmployeeLevelIdAndEmployeeGradeId(db *sql.DB, model repository.Allowance) (result repository.Allowance, errModel errorModel.ErrorModel) {
	funcName := "GetAllowanceLeaveByIdAndEmployeeLevelIdAndEmployeeGradeId"

	query := `SELECT 
				al.id, al.allowance_name, al.allowance_type,
       			efa.value 
			FROM allowances AS al
			LEFT JOIN employee_facilities_active AS efa 
				ON al.id = efa.allowance_id
			WHERE 
				al.id = $1 AND 
				al.active = TRUE AND 
				efa.active = TRUE AND 
				efa.employee_level_id = $2 AND 
				efa.employee_grade_id = $3 AND 
				(
					al.allowance_type ILIKE $4 OR
					al.allowance_type ILIKE $5
				) AND 
				al.deleted = FALSE`

	params := []interface{}{
		model.ID.Int64,
		model.EmployeeLevelId.Int64,
		model.EmployeeGradeId.Int64,
		constanta.LeaveKeyword,
		constanta.CutiKeyword,
	}

	row := db.QueryRow(query, params...)
	err := row.Scan(
		&result.ID, &result.AllowanceName, &result.AllowanceType,
		&result.Value)
	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input allowanceDAO) InitiateListByEmployeeLevelIdAndEmployeeGradeId(db *sql.DB, searchBy []in.SearchByParam, allowance repository.Allowance) (result int, errModel errorModel.ErrorModel) {
	query := `SELECT 
				COUNT(*)
			FROM allowances AS al
			INNER JOIN employee_facilities_active AS efa
				ON al.id = efa.allowance_id`

	addWhere := ` AND (
					al.allowance_type ILIKE $1 OR 
					al.allowance_type ILIKE $2
				) AND 
				efa.employee_level_id = $3 AND 
				efa.employee_grade_id = $4 AND 
				efa.active = TRUE AND 
				efa.deleted = FALSE AND 
				al.active = TRUE`

	params := []interface{}{
		constanta.LeaveKeyword,
		constanta.CutiKeyword,
		allowance.EmployeeLevelId.Int64,
		allowance.EmployeeGradeId.Int64,
	}

	return GetListDataDAO.GetCountData(db, params, query, searchBy, addWhere, DefaultFieldMustCheck{
		CreatedBy: FieldStatus{
			Value:     int64(0),
		},
		Deleted: FieldStatus{
			IsCheck: true,
			FieldName: "al.deleted",
		},
	})
}

func (input allowanceDAO) GetListByEmployeeLevelIdAndEmployeeGradeId(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, allowance repository.Allowance) (result []interface{}, errModel errorModel.ErrorModel) {
	query := `SELECT 
				al.id, al.allowance_name, al.allowance_type, 
				efa.value, al.created_at, al.updated_at
			FROM allowances AS al
			INNER JOIN employee_facilities_active AS efa
				ON al.id = efa.allowance_id`

	addWhere := ` AND (
					al.allowance_type ILIKE $1 OR 
					al.allowance_type ILIKE $2
				) AND 
				efa.employee_level_id = $3 AND 
				efa.employee_grade_id = $4 AND 
				efa.active = TRUE AND 
				efa.deleted = FALSE AND 
				al.active = TRUE`

	params := []interface{}{
		constanta.LeaveKeyword,
		constanta.CutiKeyword,
		allowance.EmployeeLevelId.Int64,
		allowance.EmployeeGradeId.Int64,
	}

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var model repository.Allowance

			err := rows.Scan(
				&model.ID, &model.AllowanceName, &model.AllowanceType,
				&model.Value, &model.CreatedAt, &model.UpdatedAt,
			)

			return model, err
		}, addWhere, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{
				Value:     int64(0),
			},
			Deleted: FieldStatus{
				IsCheck: true,
				FieldName: "al.deleted",
			},
		})
}