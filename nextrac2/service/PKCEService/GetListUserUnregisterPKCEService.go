package PKCEService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
)

func (input pkceService) GetListUserCustomForUnregister(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListUserPKCENexmileValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListPKCENexmile(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_CUSTOM_PKCE_UNREGIS_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) doGetListPKCENexmile(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	fileName := "GetListUserUnregisterPKCEService.go"
	funcName := "doGetListPKCENexmile"

	var resultGetList []interface{}

	resultGetList, err = dao.PKCEClientMappingDAO.GetListPKCEClientMappingByJoin(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	if len(resultGetList) < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.User)
		return
	}

	output = input.convertToDTO(resultGetList)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input pkceService) convertToDTO(data []interface{}) (output []out.ViewPKCEResponse) {
	for _, item := range data {
		pkceClientMappingModel := item.(repository.PKCEClientMappingModel)
		output = append(output, out.ViewPKCEResponse{
			UserID: 		pkceClientMappingModel.ID.Int64,
			ParentClientID: pkceClientMappingModel.ParentClientID.String,
			ClientID: 		pkceClientMappingModel.ClientID.String,
			Username: 		pkceClientMappingModel.Username.String,
			CreatedBy: 		pkceClientMappingModel.CreatedBy.Int64,
			UpdatedAt: 		pkceClientMappingModel.UpdatedAt.Time,
		})
	}

	return
}