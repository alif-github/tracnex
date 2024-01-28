package dao

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type benefitDAO struct {
	AbstractDAO
}

var BenefitDAO = benefitDAO{}.New()

func (input benefitDAO) New() (output benefitDAO) {
	output.FileName = "BenefitDAO.go"
	output.TableName = "benefits"
	return
}

func (input benefitDAO) GetMedicalBenefitByIdAndEmployeeLevelIdAndEmployeeGradeId(db *sql.DB, model repository.Benefit) (result repository.Benefit, errModel errorModel.ErrorModel) {
	funcName := "GetMedicalBenefitByIdAndEmployeeLevelIdAndEmployeeGradeId"

	query := `SELECT 
				b.id, b.benefit_name, b.benefit_type  
			FROM ` + input.TableName + ` AS b
			LEFT JOIN employee_facilities_active AS efa
				ON b.id = efa.benefit_id
			WHERE 
				b.id = $1 AND
				b.active = TRUE AND
				efa.active = TRUE AND
				efa.employee_level_id = $2 AND 
				efa.employee_grade_id = $3 AND
				b.benefit_type ILIKE $4 AND 
				b.deleted = FALSE`

	params := []interface{}{
		model.ID.Int64,
		model.EmployeeLevelId.Int64,
		model.EmployeeGradeId.Int64,
		constanta.MedicalKeyword,
	}

	row := db.QueryRow(query, params...)
	err := row.Scan(
		&result.ID, &result.BenefitName, &result.BenefitType,
	)

	if err != nil && err != sql.ErrNoRows {
		errModel = errorModel.GenerateInternalDBServerError(input.FileName, funcName, err)
		return
	}

	errModel = errorModel.GenerateNonErrorModel()
	return
}

func (input benefitDAO) GetListByEmployeeLevelIdAndEmployeeGradeId(db *sql.DB, userParam in.GetListDataDTO, searchBy []in.SearchByParam, benefit repository.Benefit) (result []interface{}, errModel errorModel.ErrorModel) {
	query := `SELECT 
				b.id, b.benefit_name, b.benefit_type,
				b.description, b.created_at, b.updated_at 
			FROM ` + input.TableName + ` AS b 
			INNER JOIN employee_facilities_active AS efa 
				ON b.id = efa.benefit_id`

	addWhere := ` AND b.benefit_type ILIKE $1 
				  AND efa.employee_level_id = $2  
				  AND efa.employee_grade_id = $3 
				  AND efa.active = TRUE
				  AND efa.deleted = FALSE  
				  AND b.active = TRUE`

	params := []interface{}{
		constanta.MedicalKeyword,
		benefit.EmployeeLevelId.Int64,
		benefit.EmployeeGradeId.Int64,
	}

	return GetListDataDAO.GetListDataWithDefaultMustCheck(db, params, query, userParam, searchBy,
		func(rows *sql.Rows) (interface{}, error) {
			var model repository.Benefit

			err := rows.Scan(
				&model.ID, &model.BenefitName, &model.BenefitType,
				&model.Description, &model.CreatedAt, &model.UpdatedAt,
			)

			return model, err
		}, addWhere, DefaultFieldMustCheck{
			CreatedBy: FieldStatus{
				Value:     int64(0),
			},
			Deleted: FieldStatus{
				IsCheck: true,
				FieldName: "b.deleted",
			},
		})
}

func (input benefitDAO) InitiateListByBenefitType(db *sql.DB, searchBy []in.SearchByParam, benefit repository.Benefit) (result int, errModel errorModel.ErrorModel) {
	query := `SELECT 
				COUNT(*) 
			FROM ` + input.TableName + ` AS b 
			INNER JOIN employee_facilities_active AS efa 
				ON b.id = efa.benefit_id`

	addWhere := ` AND b.benefit_type ILIKE $1 
				  AND efa.employee_level_id = $2  
				  AND efa.employee_grade_id = $3 
				  AND efa.active = TRUE 
				  AND efa.deleted = FALSE 
				  AND b.active = TRUE`

	params := []interface{}{
		constanta.MedicalKeyword,
		benefit.EmployeeLevelId.Int64,
		benefit.EmployeeGradeId.Int64,
	}

	return GetListDataDAO.GetCountData(db, params, query, searchBy, addWhere, DefaultFieldMustCheck{
		CreatedBy: FieldStatus{
			Value:     int64(0),
		},
		Deleted: FieldStatus{
			IsCheck: true,
			FieldName: "b.deleted",
		},
	})
}
