package BacklogService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"strconv"
	"time"
)

func (input backlogService) DeleteBacklogBySprint(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName      = "DeleteBacklogBySprint"
		validSearchBy = []string{"sprint"}
		searchByParam []in.SearchByParam
		inputStruct   in.BacklogRequest
	)

	_, searchByParam, err = input.ReadAndValidateDeleteListDataBacklog(request, validSearchBy, applicationModel.DeleteParentBacklogValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	input.setSprintFromSearchByParam(searchByParam, &inputStruct)

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteListBacklogBySprint, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input backlogService) setSprintFromSearchByParam(searchByParam []in.SearchByParam, inputStruct *in.BacklogRequest) {
	for _, searchBy := range searchByParam {
		inputStruct.Sprint = searchBy.SearchValue
	}
}

func (input backlogService) doDeleteListBacklogBySprint(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct                      = inputStructInterface.(in.BacklogRequest)
		listBacklogOnDBBySprintInterface []interface{}
		scope                            map[string]interface{}
		searchByParam                    []in.SearchByParam
		userParam                        = in.GetListDataDTO{
			AbstractDTO: in.AbstractDTO{
				Page:  -99,
				Limit: -99,
			},
		}
	)

	inputModel := repository.BacklogModel{
		Sprint:    sql.NullString{String: inputStruct.Sprint},
		CreatedBy: sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
	}

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	// get all backlog by sprint
	listBacklogOnDBBySprintInterface, err = dao.BacklogDAO.GetListDetailBacklogBySprint(serverconfig.ServerAttribute.DBConnection, userParam, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB, inputModel)
	if err.Error != nil {
		return
	}

	// delete all backlog
	for _, backlogOnDBInterface := range listBacklogOnDBBySprintInterface {
		var (
			item          repository.BacklogModel
			redmineNumber int64
		)

		item = backlogOnDBInterface.(repository.BacklogModel)
		redmineNumber, err = input.reformatRandomNumber2(item.RedmineNumber.Int64)
		if err.Error != nil {
			return
		}

		err = dao.BacklogDAO.DeleteBacklog(tx, repository.BacklogModel{
			ID:            sql.NullInt64{Int64: item.ID.Int64},
			RedmineNumber: sql.NullInt64{Int64: redmineNumber},
			UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedAt:     sql.NullTime{Time: timeNow},
			UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		})

		if err.Error != nil {
			return
		}
	}

	// Delete Backlog
	//dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.BacklogDAO.TableName, inputModel.ID.Int64, contextModel.LimitedByCreatedBy)...)

	//err = dao.BacklogDAO.DeleteBacklogBySprint(tx, inputModel, contextModel)
	//if err.Error != nil {
	//	return
	//}

	return
}

func (input backlogService) reformatRandomNumber2(inputNumber int64) (outputNumber int64, err errorModel.ErrorModel) {
	var (
		errs error
	)

	inputNumberString := strconv.FormatInt(inputNumber, 10)
	randomString2Digit := strconv.FormatInt(randomNumber2Digit(), 10)

	outputNumberNewString := inputNumberString + randomString2Digit
	outputNumber, errs = strconv.ParseInt(outputNumberNewString, 10, 64)
	if errs != nil {
		err = errorModel.GenerateSimpleErrorModel(400, "Error parsing uniq 2 digit")
		return
	}

	return
}
