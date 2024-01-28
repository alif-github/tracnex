package repository

import "database/sql"

type EmployeeModel struct {
	ID                 sql.NullInt64
	NIK                sql.NullInt64
	IDCard             sql.NullString
	NPWP               sql.NullString
	FirstName          sql.NullString
	LastName           sql.NullString
	Email              sql.NullString
	Phone              sql.NullString
	Gender             sql.NullString
	PlaceOfBirth       sql.NullString
	DateOfBirth        sql.NullTime
	AddressResidence   sql.NullString
	AddressTax         sql.NullString
	DateJoin           sql.NullTime
	DateOut            sql.NullTime
	Religion           sql.NullString
	Type               sql.NullString
	Status             sql.NullString
	PositionID         sql.NullInt64
	Position           sql.NullString
	MaritalStatus      sql.NullString
	Education          sql.NullString
	MothersMaiden      sql.NullString
	NumberOfDependents sql.NullInt64
	Nationality        sql.NullString
	TaxMethod          sql.NullString
	ReasonResignation  sql.NullString
	Photo              sql.NullString
	IsLead             sql.NullBool
	Active             sql.NullBool
	EmployeeVariableID sql.NullInt64
	DepartmentId       sql.NullInt64
	FromDate           sql.NullTime
	ThruDate           sql.NullTime
	Information        sql.NullString
	ContractNo         sql.NullString
	Name               sql.NullString
	RedmineId          sql.NullInt64
	DepartmentName     sql.NullString
	CreatedAt          sql.NullTime
	CreatedBy          sql.NullInt64
	CreatedName        sql.NullString
	CreatedClient      sql.NullString
	UpdatedBy          sql.NullInt64
	UpdatedAt          sql.NullTime
	UpdatedName        sql.NullString
	UpdatedClient      sql.NullString
	Deleted            sql.NullBool
	IsUsed             sql.NullBool
	MandaysRate        sql.NullString
	BPJS               sql.NullString
	BPJSTk             sql.NullString
	LevelID            sql.NullInt64
	Level              sql.NullString
	GradeID            sql.NullInt64
	Grade              sql.NullString
	IsHaveMember       sql.NullBool
	IsHaveVariable     sql.NullBool
	Member             sql.NullString
	MemberList         []MemberList
	ClientId           sql.NullString
	FileUploadID       sql.NullInt64
}

type MemberList struct {
	ID        sql.NullInt64
	FirstName sql.NullString
	LastName  sql.NullString
}
