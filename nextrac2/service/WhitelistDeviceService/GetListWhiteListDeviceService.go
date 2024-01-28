package WhitelistDeviceService

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

func (input whitelistDeviceService) InitiateWhiteListDevice(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		countData     interface{}
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListWhiteListDeviceValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateWhiteListDevice(searchByParam, contextModel)

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListWhiteListDeviceValidOperator,
		CountData:     countData.(int),
	}
	return
}

func (input whitelistDeviceService) doInitiateWhiteListDevice(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		createdBy int64 = contextModel.LimitedByCreatedBy
	)

	output, err = dao.WhiteListDevice.GetCountWhiteListDevice(serverconfig.ServerAttribute.DBConnection, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	return
}

func (input whitelistDeviceService) GetListWhiteListDevice(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListWhiteListDeviceValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListWhiteListDevice(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input whitelistDeviceService) doGetListWhiteListDevice(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
	)

	dbResult, err = dao.WhiteListDevice.GetListWhiteListDevice(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseGetList(dbResult)
	return
}

func (input whitelistDeviceService) convertModelToResponseGetList(dbResult []interface{}) (result []out.WhiteListDeviceResponse) {
	for _, dbResultItem := range dbResult {
		item := dbResultItem.(repository.WhiteListDeviceModel)
		result = append(result, out.WhiteListDeviceResponse{
			ID:          item.ID.Int64,
			Device:      item.Device.String,
			Description: item.Description.String,
			CreatedAt:   item.CreatedAt.Time,
			UpdatedAt:   item.UpdatedAt.Time,
			UpdatedName: item.UpdatedName.String,
		})
	}

	return result
}