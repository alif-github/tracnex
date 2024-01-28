package ClientRegistrationLogService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input clientRegistrationLogService) GetListLog(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListRegistrationLogValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListLog(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code: 		util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: 	GenerateI18NMessage("SUCCESS_GET_LIST_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientRegistrationLogService) doGetListLog(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var resultGetList []interface{}

	resultGetList, err = dao.ClientRegistrationLogDAO.GetListLog(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	output = input.convertToDTO(resultGetList)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientRegistrationLogService) convertToDTO(data []interface{}) (logs []out.GetListLogDTOOut) {
	for _, item := range data {
		log := item.(repository.ListClientRegistrationLogModel)
		logs = append(logs, out.GetListLogDTOOut {
			ID: 					log.ID.Int64,
			ClientID: 				log.ClientID.String,
			ClientTypeID: 			log.ClientTypeID.Int64,
			SuccessStatusAuth: 		log.SuccessStatusAuth.Bool,
			SuccessStatusNexcloud: 	log.SuccessStatusNexcloud.Bool,
			Resource: 				log.Resource.String,
		})
	}

	return
}