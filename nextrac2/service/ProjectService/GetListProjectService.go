package ProjectService

import (
	"encoding/json"
	"net/http"
	"net/url"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"strconv"
)

func (input projectService) GetListProjectByRedmineAPI(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListBankValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListProjectByRedmineAPI(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input projectService) doGetListProjectByRedmineAPI(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		result []string
	)

	responses, err := HitAPIProjectRedmine(inputStruct, searchByParam, contextModel)

	for _, response := range responses.Projects {
		result = append(result, response.Name)
	}

	output = result
	return
}

func HitAPIProjectRedmine(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (structResponse out.APIProjectRedmineResponse, err errorModel.ErrorModel) {
	var (
		fileName = "GetListProjectService.go"
		funcName = "HitAPIProjectRedmine"
		token    = "1e745f89c9b66e9a6418460cefa5a253c1c3bc3d"
		baseUrl  = "http://redmine.nexcloud.id/projects.json"
	)

	header := make(map[string][]string)
	header["X-Redmine-API-Key"] = []string{token}

	parameters := url.Values{}
	for _, search := range searchByParam {
		switch search.SearchKey {
		case "name":
			parameters.Add(search.SearchKey, "~"+search.SearchValue)
		}
	}

	if inputStruct.Limit > 0 {
		parameters.Add("limit", strconv.Itoa(inputStruct.Limit))
	}

	if inputStruct.Page > 0 {
		parameters.Add("page", strconv.Itoa(inputStruct.Page))
	}

	u, errs := url.Parse(baseUrl)
	if errs != nil {
		err = errorModel.GenerateSimpleErrorModel(400, "Error parsing URL")
		return
	}

	u.RawQuery = parameters.Encode()

	statusCode, _, bodyResult, errs := common.HitAPI(u.String(), header, "", "GET", *contextModel)
	if errs != nil {
		err = errorModel.GenerateErrorModel(statusCode, errs.Error(), fileName, funcName, errs)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &structResponse)

	return
}
