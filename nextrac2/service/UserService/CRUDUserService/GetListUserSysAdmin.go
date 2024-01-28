package CRUDUserService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input userService) GetListUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListUserValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListUser(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_LIST_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) InitiateUser(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		count         int
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListUserValidOperator)
	if err.Error != nil {
		return
	}

	count, err = dao.UserDAO.GetCountUser(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListUserValidOperator,
		EnumData:      nil,
		CountData:     count,
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_INITIATE_LIST_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) doGetListUser(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var resultGetList []interface{}

	resultGetList, err = dao.UserDAO.GetListUser(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	output = input.convertToDTO(resultGetList)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) convertToDTO(data []interface{}) (clients []out.GetListUserDTOOut) {
	for _, item := range data {
		client := item.(repository.ListUserModel)
		clients = append(clients, out.GetListUserDTOOut{
			ID:          client.ID.Int64,
			ClientID:    client.ClientID.String,
			AuthUserID:  client.AuthUserID.Int64,
			Firstname:   client.FirstName.String,
			Lastname:    client.LastName.String,
			Username:    client.Username.String,
			Email:       client.Email.String,
			Phone:       client.Phone.String,
			RoleID:      client.RoleID.String,
			GroupID:     client.GroupID.String,
			Status:      client.Status.String,
			Locale:      client.Locale.String,
			CreatedName: client.CreatedName.String,
			CreatedAt:   client.CreatedAt.Time,
			UpdatedAt:   client.UpdatedAt.Time,
			PlatformDevice:   client.PlatformDevice.String,
		})
	}

	return
}
