package TrackerService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input trackerService) GetListTracker(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	inputStruct, _, err = input.ReadAndValidateGetListData(request, nil, input.ValidOrderBy, nil, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListTracker(inputStruct)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input trackerService) doGetListTracker(inputStruct in.GetListDataDTO) (output interface{}, err errorModel.ErrorModel) {
	var dbResult []interface{}
	dbResult, err = dao.TrackerDAO.GetListTrackerOnRedmine(serverconfig.ServerAttribute.RedmineDBConnection, inputStruct)
	if err.Error != nil {
		return
	}

	output = input.convertToDTOOut(dbResult)
	return
}

func (input trackerService) convertToDTOOut(dbResult []interface{}) (result []string) {
	for _, item := range dbResult {
		temp := item.(string)
		result = append(result, temp)
	}

	return result
}
