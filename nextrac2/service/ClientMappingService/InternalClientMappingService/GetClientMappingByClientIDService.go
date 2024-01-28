package InternalClientMappingService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/ClientMappingService"
	"nexsoft.co.id/nextrac2/util"
)

type getClientMappingByClientIDService struct {
	service.AbstractService
	service.GetListData
}

var GetClientMappingByClientIDService = getClientMappingByClientIDService{}.New()

func (input getClientMappingByClientIDService) New() (output getClientMappingByClientIDService) {
	output.FileName = "GetClientMappingByClientIDService.go"
	output.ValidSearchBy = []string{}
	output.ValidOrderBy = []string{"client_mapping.created_at", "client_mapping.updated_at"}
	output.ValidLimit = service.DefaultLimit

	return
}

func (input getClientMappingByClientIDService) GetClientMappingsByClientID(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListDetailClientMappingValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	clientMappingBody,bodySize, err := ClientMappingService.GetClientMappingBodies(request, input.FileName)
	contextModel.LoggerModel.ByteIn = bodySize
	if err.Error != nil {
		return
	}

	err = clientMappingBody.ValidationForGetClientMappingByClientID()
	if err.Error != nil {
		return
	}

	err = input.validation(clientMappingBody)
	if err.Error != nil {
		return
	}

	clientMappingModel := input.convertToRepo(clientMappingBody)
	clients, err := dao.ClientMappingDAO.GetDetailClientMappingByClientID(serverconfig.ServerAttribute.DBConnection, clientMappingModel, inputStruct, searchByParam, "")
	if err.Error != nil {
		return
	}

	if clients == nil {
		err = errorModel.GenerateUnknownDataError(input.FileName, "GetClientMappingsByClientID", constanta.ClientMapping)
		return
	}

	output.Data.Content = input.convertToDTO(clients)
	output.Status =out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: ClientMappingService.GenerateI18Message("SUCCESS_VIEW_CLIENT_MAPPING_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input getClientMappingByClientIDService) convertToDTO(data []interface{}) (clients []out.ClientMappingResponse) {
	for _, item:= range data{
		client := item.(repository.CLientMappingDetailForViewModel)
		clients = append(clients, out.ClientMappingResponse{
			ID:           client.ID.Int64,
			ClientID:     client.ClientId.String,
			ClientTypeId: client.ClientTypeId.Int64,
			AuthUserId:   client.AuthUserId.Int64,
			Username:     client.Username.String,
			CompanyId:    client.CompanyId.String,
			BranchId:     client.BranchId.String,
			Aliases:      client.Aliases.String,
			SocketID:     client.SocketID.String,
			UpdatedAt:    client.UpdatedAt.Time,
			UpdatedBy:    client.UpdatedBy.Int64,
			CreatedAt:    client.CreatedAt.Time,
			CreatedBy:    client.CreatedBy.Int64,
		})
	}
	return
}

func (input getClientMappingByClientIDService) convertToRepo(clientMappingBody in.ClientMappingRequest) (clients repository.ClientMappingForDetailModel) {
	clients.ClientTypeID.Int64 = clientMappingBody.ClientTypeID

	for _, companyData := range clientMappingBody.ClientData {
		clients.ClientData = append(clients.ClientData, repository.ClientDataModel{ClientID: sql.NullString{String: companyData.ClientID}})
	}

	return
}

func (input getClientMappingByClientIDService) validation(inputStruct in.ClientMappingRequest) (err errorModel.ErrorModel) {
	funcName := "validation"

	result, err := dao.ClientTypeDAO.CheckClientTypeByID(serverconfig.ServerAttribute.DBConnection, &repository.ClientTypeModel{
		ID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
	})

	if err.Error != nil {
		return
	}

	if result.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ClientTypeID)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}