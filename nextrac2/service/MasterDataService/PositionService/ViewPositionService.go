package PositionService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
)

func (input positionService) ViewPosition(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct master_data_request.PositionGetListRequest
	
	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateViewPosition)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewPosition(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input positionService) doViewPosition(inputStruct master_data_request.PositionGetListRequest, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	output, err = master_data_dao.GetViewPositionFromMasterData(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input positionService) ValidateViewPosition(inputStruct *master_data_request.PositionGetListRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}
