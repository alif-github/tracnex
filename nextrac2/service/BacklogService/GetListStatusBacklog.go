package BacklogService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
)

func (input backlogService) GetListStatusBacklog(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct         in.GetListDataDTO
		searchByParam       []in.SearchByParam
		statusValidOrderBy  = []string{"department_id"}
		statusValidSearchBy = []string{"department_id"}
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, statusValidSearchBy, statusValidOrderBy, applicationModel.GetListStatusValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListStatusBacklog(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)

	return
}

func (input backlogService) doGetListStatusBacklog(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		statusDev = constanta.StatusAllowedDeveloper
		statusQA  = constanta.StatusAllowedQA
		statusAll []string
	)

	if searchByParam != nil {
		for _, filter := range searchByParam {
			if filter.SearchValue == "1" {
				output = statusDev
				return
			} else if filter.SearchValue == "2" {
				output = statusQA
				return
			} else {
				return
			}
		}
	}

	// Combine All Status
	uniqueMap := make(map[string]string)
	for _, dev := range statusDev {
		uniqueMap[dev] = dev
	}

	for _, qa := range statusQA {
		uniqueMap[qa] = qa
	}

	// Mengonversi map kembali ke array
	for num := range uniqueMap {
		statusAll = append(statusAll, num)
	}

	output = statusAll
	return
}
