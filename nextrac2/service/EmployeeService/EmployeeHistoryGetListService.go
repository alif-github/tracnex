package EmployeeService

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"strconv"
	"strings"
	"time"
)

func (input employeeService) InitiateGetListEmployeeHistory(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		validSearchBy = []string{"id", "name", "primary_key"}
		validOrderBy  = []string{"created_at"}
		validLimit    = []int{3, 5}
		countData     interface{}
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, validSearchBy, applicationModel.GetListEmployeeHistoryValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateEmployeeHistory(searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  validOrderBy,
		ValidSearchBy: validOrderBy,
		ValidLimit:    validLimit,
		ValidOperator: applicationModel.GetListEmployeeHistoryValidOperator,
		CountData:     countData.(int),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) GetListEmployeeHistory(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
		validSearchBy = []string{"id", "name", "primary_key"}
		validOrderBy  = []string{"created_at"}
		validLimit    = []int{3, 5}
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, applicationModel.GetListEmployeeHistoryValidOperator, validLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListEmployeeHistory(inputStruct, searchByParam, *contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input employeeService) doInitiateEmployeeHistory(searchByParam []in.SearchByParam, _ *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var db = serverconfig.ServerAttribute.DBConnection
	output, err = dao.EmployeeDAO.GetCountEmployeeHistory(db, searchByParam)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) doGetListEmployeeHistory(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		db       = serverconfig.ServerAttribute.DBConnection
		dbResult []interface{}
	)

	dbResult, err = dao.EmployeeDAO.GetListEmployeeHistory(db, inputStruct, searchByParam)
	if err.Error != nil {
		return
	}

	output, err = input.convertModelToRspListHistory(dbResult, &contextModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) convertModelToRspListHistory(dbResult []interface{}, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		fileName      = "EmployeeHistoryGetListService.go"
		funcName      = "convertModelToRspListHistory"
		resultHistory []out.EmployeeHistoryListResponse
	)

	for _, dbResultItem := range dbResult {
		var (
			desc1             []out.DescriptionDetail
			desc2             []out.DescriptionDetail
			resultDescription []out.DescriptionDetail
			item              = dbResultItem.(repository.EmployeeHistoryModel)
		)

		_ = json.Unmarshal([]byte(item.Description1.String), &desc1)
		_ = json.Unmarshal([]byte(item.Description2.String), &desc2)

		//--- Get Main
		for i := 0; i < len(desc1); i++ {
			//--- Department
			if strings.ToLower(desc1[i].KeyID) == "departemen" {
				//--- Before Check
				if desc1[i].Before != "" {
					if err = input.getDepartmentName(&desc1[i].Before); err.Error != nil {
						return
					}
				}

				//--- After Check
				if desc1[i].After != "" {
					if err = input.getDepartmentName(&desc1[i].After); err.Error != nil {
						return
					}
				}
			}

			//--- Position
			if strings.ToLower(desc1[i].KeyID) == "posisi" {
				//--- Before Check
				if desc1[i].Before != "" {
					if err = input.getPositionName(&desc1[i].Before); err.Error != nil {
						return
					}
				}

				//--- After Check
				if desc1[i].After != "" {
					if err = input.getPositionName(&desc1[i].After); err.Error != nil {
						return
					}
				}
			}

			//--- Join Date
			if strings.ToLower(desc1[i].KeyID) == "tanggal bergabung" || strings.ToLower(desc1[i].KeyID) == "tanggal keluar" || strings.ToLower(desc1[i].KeyID) == "tanggal lahir" {
				var (
					timeBefore time.Time
					timeAfter  time.Time
					errors     error
				)

				//--- Before Check
				if desc1[i].Before != "" {
					timeBefore, errors = time.Parse(constanta.DefaultDBTimeFormat, desc1[i].Before)
					if errors != nil {
						err = errorModel.GenerateUnknownError(fileName, funcName, errors)
						return
					}
					desc1[i].Before = timeBefore.Format(constanta.DefaultInstallationTimeFormat)
				}

				//--- After Check
				if desc1[i].After != "" {
					timeAfter, errors = time.Parse(constanta.DefaultDBTimeFormat, desc1[i].After)
					if errors != nil {
						err = errorModel.GenerateUnknownError(fileName, funcName, errors)
						return
					}
					desc1[i].After = timeAfter.Format(constanta.DefaultInstallationTimeFormat)
				}
			}
		}

		//--- Get Level and Grade
		for j := 0; j < len(desc2); j++ {
			//--- Level
			if strings.ToLower(desc2[j].KeyID) == "level" {
				//--- Before Check
				if desc2[j].Before != "" {
					if err = input.getLevelName(&desc2[j].Before); err.Error != nil {
						return
					}
				}

				//--- After Check
				if desc2[j].After != "" {
					if err = input.getLevelName(&desc2[j].After); err.Error != nil {
						return
					}
				}
			}

			//--- Grade
			if strings.ToLower(desc2[j].KeyID) == "grade" {
				//--- Before Check
				if desc2[j].Before != "" {
					if err = input.getGradeName(&desc2[j].Before); err.Error != nil {
						return
					}
				}

				//--- After Check
				if desc2[j].After != "" {
					if err = input.getGradeName(&desc2[j].After); err.Error != nil {
						return
					}
				}
			}
		}

		resultDescription = append(resultDescription, desc1...)
		resultDescription = append(resultDescription, desc2...)

		resultHistory = append(resultHistory, out.EmployeeHistoryListResponse{
			ID:                item.ID.Int64,
			Editor:            item.Editor.String,
			CreatedAt:         item.CreatedAt.Time,
			DescriptionDetail: resultDescription,
		})
	}

	result = out.BundleEmployeeHistoryResponse{
		Locale:      contextModel.AuthAccessTokenModel.Locale,
		ListHistory: resultHistory,
	}

	return
}

func (input employeeService) getLevelName(level *string) (err errorModel.ErrorModel) {
	var (
		i  int
		r  repository.EmployeeLevelModel
		db = serverconfig.ServerAttribute.DBConnection
	)

	//--- Convert Number
	i, _ = strconv.Atoi(*level)

	if i > 0 {
		r, err = dao.EmployeeLevelGradeDAO.GetNameEmployeeLevelByIDWithoutDeleted(db, repository.EmployeeLevelModel{ID: sql.NullInt64{Int64: int64(i)}})
		if err.Error != nil {
			return
		}

		if r.Level.String != "" {
			*level = r.Level.String
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) getGradeName(grade *string) (err errorModel.ErrorModel) {
	var (
		i  int
		r  repository.EmployeeGradeModel
		db = serverconfig.ServerAttribute.DBConnection
	)

	//--- Convert Number
	i, _ = strconv.Atoi(*grade)

	if i > 0 {
		r, err = dao.EmployeeLevelGradeDAO.GetNameEmployeeGradeByIDWithoutDeleted(db, repository.EmployeeGradeModel{ID: sql.NullInt64{Int64: int64(i)}})
		if err.Error != nil {
			return
		}

		if r.Grade.String != "" {
			*grade = r.Grade.String
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) getDepartmentName(department *string) (err errorModel.ErrorModel) {
	var (
		i  int
		r  repository.DepartmentModel
		db = serverconfig.ServerAttribute.DBConnection
	)

	//--- Convert Number
	i, _ = strconv.Atoi(*department)

	if i > 0 {
		r, err = dao.DepartmentDAO.GetNameDepartmentByIDWithoutDeleted(db, repository.DepartmentModel{ID: sql.NullInt64{Int64: int64(i)}})
		if err.Error != nil {
			return
		}

		if r.Name.String != "" {
			*department = r.Name.String
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input employeeService) getPositionName(position *string) (err errorModel.ErrorModel) {
	var (
		i  int
		r  repository.EmployeePositionModel
		db = serverconfig.ServerAttribute.DBConnection
	)

	//--- Convert Number
	i, _ = strconv.Atoi(*position)

	if i > 0 {
		r, err = dao.EmployeePositionDAO.GetNamePositionByIDWithoutDeleted(db, repository.EmployeePositionModel{ID: sql.NullInt64{Int64: int64(i)}})
		if err.Error != nil {
			return
		}

		if r.Name.String != "" {
			*position = r.Name.String
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
