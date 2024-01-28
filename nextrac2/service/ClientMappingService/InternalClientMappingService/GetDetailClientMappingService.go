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

type getDetailClientMappingService struct {
	service.AbstractService
	service.GetListData
}

var GetDetailClientMappingService = getDetailClientMappingService{}.New()

func (input getDetailClientMappingService) New() (output getDetailClientMappingService) {
	output.FileName = "GetDetailClientMappingService.go"
	output.ValidSearchBy = []string{}
	output.ValidOrderBy = []string{"client_mapping.created_at", "client_mapping.updated_at"}
	output.ValidLimit = service.DefaultLimit
	return
}

func (input getDetailClientMappingService) GetDetailClientMappings(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var clientId string
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

	err = clientMappingBody.ValidationForDetail()
	if err.Error != nil {
		return
	}

	err = input.validation(clientMappingBody)
	if err.Error != nil {
		return
	}

	//_, isOnlyHaveOwnPermission := service.CheckIsOnlyHaveOwnPermission(*contextModel)
	//if isOnlyHaveOwnPermission {
	//	clientId = contextModel.AuthAccessTokenModel.ClientID
	//}

	clientMappingRepository := input.convertToRepo(clientMappingBody)
	clients, err :=dao.ClientMappingDAO.GetDetailClientMapping(serverconfig.ServerAttribute.DBConnection, clientMappingRepository, inputStruct, searchByParam, clientId)
	if err.Error != nil {
		return
	}

	if clients == nil {
		err = errorModel.GenerateUnknownDataError(input.FileName, "GetDetailClientMappings", constanta.ClientMapping)
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

func (input getDetailClientMappingService) convertToDTO(data []interface{}) (clients []out.ClientMappingResponse) {
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

func (input getDetailClientMappingService) convertToRepo(clientMappingBody in.ClientMappingRequest) (clients repository.ClientMappingForDetailModel) {
	clients.ClientTypeID.Int64 = clientMappingBody.ClientTypeID

	var companies []repository.CompanyDataModel
	for _, companyData := range clientMappingBody.CompanyData {
		var branches []repository.BranchDataModel
		for _, branchData := range companyData.BranchData {
			branches = append(branches, repository.BranchDataModel{BranchID: sql.NullString{String: branchData.BranchID}})
		}
		companies = append(companies, repository.CompanyDataModel{
			CompanyID:  sql.NullString{String: companyData.CompanyID},
			BranchData: branches,
		})
	}
	clients.CompanyData = companies
	return
}

func (input getDetailClientMappingService) validation(inputStruct in.ClientMappingRequest) (err errorModel.ErrorModel) {
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


