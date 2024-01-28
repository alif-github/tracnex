package EnumService

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

func (input enumService) ViewDetailEnum(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct in.EnumRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validationViewDetailEnum)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewDetailEnum(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input enumService) doViewDetailEnum(inputStruct in.EnumRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	var (
		listEnumOnDB  []interface{}
		searchByParam []in.SearchByParam
		userParam     in.GetListDataDTO
	)

	listEnumOnDB, err = dao.EnumDAO.GetListEnumLabel(serverconfig.ServerAttribute.DBConnection, inputStruct, userParam, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	result = input.convertModelToResponseGetListEnum(listEnumOnDB)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input enumService) validationViewDetailEnum(inputStruct *in.EnumRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}

func (input enumService) convertModelToResponseGetListEnum(dbResult []interface{}) (result []out.EnumLabelResponse) {
	for _, dbResultItem := range dbResult {
		item := dbResultItem.(repository.EnumModel)
		result = append(result, out.EnumLabelResponse{
			Value: item.EnumLabel.String,
		})
	}

	return
}
