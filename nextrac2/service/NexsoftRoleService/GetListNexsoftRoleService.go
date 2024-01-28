package NexsoftRoleService

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
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
)

type getListNexsoftRoleService struct {
	FileName string
	service.GetListData
}

var GetListNexsoftRoleService = getListNexsoftRoleService{}.New()

func (input getListNexsoftRoleService) New() (output getListNexsoftRoleService) {
	output.FileName = "GetListNexsoftRoleService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"role_id",
		"description",
		"created_at",
		"created_name",
	}
	output.ValidSearchBy = []string{
		"role_id",
		"description",
	}
	return
}

func (input getListNexsoftRoleService) InitiateNexsoftRole(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		countData     int
		searchByParam []in.SearchByParam
		db            = serverconfig.ServerAttribute.DBConnection
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListRoleValidOperator)
	if err.Error != nil {
		return
	}

	userID, isOnlyHaveOwnPermission := service.CheckIsOnlyHaveOwnPermission(*contextModel)
	if !isOnlyHaveOwnPermission {
		userID = 0
	}

	countData, err = dao.NexsoftRoleDAO.GetCountNexsoftRole(db, searchByParam, userID)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListRoleValidOperator,
		CountData:     countData,
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_GET_LIST_ROLE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input getListNexsoftRoleService) GetListNexsoftRole(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct     in.GetListDataDTO
		searchByParam   []in.SearchByParam
		nexsoftRoleList []interface{}
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListRoleValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	err = input.SetDefaultOrder(request, constanta.CreatedAtDesc, &inputStruct, input.ValidOrderBy)
	if err.Error != nil {
		return
	}

	nexsoftRoleList, err = input.getListNexsoftRoles(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.convertListToDTOOut(nexsoftRoleList)
	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_ROLE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input getListNexsoftRoleService) getListNexsoftRoles(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output []interface{}, err errorModel.ErrorModel) {
	var (
		userID                  int64
		isOnlyHaveOwnPermission bool
	)

	userID, isOnlyHaveOwnPermission = service.CheckIsOnlyHaveOwnPermission(*contextModel)
	if !isOnlyHaveOwnPermission {
		userID = 0
	}

	output, err = dao.NexsoftRoleDAO.GetListNexsoftRole(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, userID)
	return
}

func (input getListNexsoftRoleService) convertListToDTOOut(dbResult []interface{}) (result []out.ViewListRoleDTOOut) {
	for _, item := range dbResult {
		itemOnRepo := item.(repository.RoleModel)
		result = append(result, out.ViewListRoleDTOOut{
			ID:          itemOnRepo.ID.Int64,
			RoleID:      itemOnRepo.RoleID.String,
			Description: itemOnRepo.Description.String,
			CreatedBy:   itemOnRepo.CreatedBy.Int64,
			CreatedAt:   itemOnRepo.CreatedAt.Time,
			CreatedName: itemOnRepo.CreatedName.String,
			UpdatedAt:   itemOnRepo.UpdatedAt.Time,
		})
	}
	return
}
