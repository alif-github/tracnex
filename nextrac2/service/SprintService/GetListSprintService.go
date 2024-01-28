package SprintService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input sprintService) GetListSprint(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	inputStruct, _, err = input.ReadAndValidateGetListData(request, nil, input.ValidOrderBy, nil, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListSprint(inputStruct)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input sprintService) doGetListSprint(inputStruct in.GetListDataDTO) (output interface{}, err errorModel.ErrorModel) {
	var (
		sprint []string
		raw    []string
		db     = serverconfig.ServerAttribute.RedmineDBConnection
	)

	sprint, err = dao.SprintDAO.GetListSprintOnRedmine(db)
	if err.Error != nil {
		return
	}

	raw, err = dao.SprintDAO.ReArrangeDataSprint(db, inputStruct, sprint)
	if err.Error != nil {
		return
	}

	output = raw
	return
}

func (input sprintService) convertToDTOOut(dbResult []interface{}) (result []string) {
	for _, item := range dbResult {
		temp := item.(string)
		result = append(result, temp)
	}

	return result
}
