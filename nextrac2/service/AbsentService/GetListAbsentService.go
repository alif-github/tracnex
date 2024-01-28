package AbsentService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input absentService) GetListAbsent(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListAbsentValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListAbsent(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input absentService) doGetListAbsent(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
		scope    map[string]interface{}
		db       = serverconfig.ServerAttribute.DBConnection
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	if err = input.PeriodCheck(&searchByParam); err.Error != nil {
		return
	}

	dbResult, err = dao.AbsentDAO.GetListAbsent(db, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertToListDTOOut(dbResult)
	return
}

func (input absentService) convertToListDTOOut(dbResult []interface{}) (result []out.ListAbsent) {
	for _, dbResultItem := range dbResult {
		repo := dbResultItem.(repository.AbsentModel)
		result = append(result, out.ListAbsent{
			IDCard:          repo.IDCard.String,
			Name:            repo.EmployeeName.String,
			NormalDays:      repo.NormalDays.Int64,
			ActualDays:      repo.ActualDays.Int64,
			Absent:          repo.Absent.Int64,
			Overdue:         repo.Overdue.Int64,
			LeaveEarly:      repo.LeaveEarly.Int64,
			Overtime:        repo.Overtime.Int64,
			NumbersOfLeave:  repo.NumberOfLeave.Int64,
			LeavingDuties:   repo.LeavingDuties.Int64,
			NumbersIn:       repo.NumbersIn.Int64,
			NumbersOut:      repo.NumbersOut.Int64,
			Scan:            repo.Scan.Int64,
			SickLeave:       repo.SickLeave.Int64,
			PaidLeave:       repo.PaidLeave.Int64,
			PermissionLeave: repo.PermissionLeave.Int64,
			WorkHours:       repo.WorkHours.Int64,
			PercentAbsent:   repo.PercentAbsent.Float64,
		})
	}

	return result
}
