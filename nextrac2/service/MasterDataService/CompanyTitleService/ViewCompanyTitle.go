package CompanyTitleService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
)

func (input companyTitleService) ViewCompanyTitle(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct master_data_request.CompanyTitleRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateViewCompanyTitle)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.DoViewCompanyTitle(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input companyTitleService) DoViewCompanyTitle(inputStruct master_data_request.CompanyTitleRequest, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	output, err = master_data_dao.GetViewCompanyTitleFromMasterData(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input companyTitleService) ValidateViewCompanyTitle(inputStruct *master_data_request.CompanyTitleRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}