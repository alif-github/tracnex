package DataGroupService

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/ClientTypeService"
	"nexsoft.co.id/nextrac2/service/CustomerCategoryService"
	"nexsoft.co.id/nextrac2/service/CustomerGroupService"
	"nexsoft.co.id/nextrac2/service/EmployeeService"
	"nexsoft.co.id/nextrac2/service/MasterDataService/DistrictService"
	"nexsoft.co.id/nextrac2/service/MasterDataService/ProvinceService"
	"nexsoft.co.id/nextrac2/service/ProductGroupService"
	"nexsoft.co.id/nextrac2/service/SalesmanService"
	"sort"
)

func (input dataGroupService) ViewDataGroup(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.DataGroupRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewDataGroup(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input dataGroupService) doViewDataGroup(inputStruct in.DataGroupRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	funcName := "doViewDataGroup"
	dataGroup := repository.DataGroupModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	}

	dataGroup.CreatedBy.Int64 = contextModel.LimitedByCreatedBy

	dataGroup, err = dao.DataGroupDAO.ViewDataGroup(serverconfig.ServerAttribute.DBConnection, dataGroup)
	if err.Error != nil {
		return
	}

	if dataGroup.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.DataGroup)
		return
	}

	result, err = reformatDAOtoDTO(dataGroup)
	return
}

func reformatDAOtoDTO(dataGroup repository.DataGroupModel) (result out.ViewDetailDataGroupResponse, err errorModel.ErrorModel) {
	var groupsScope map[string][]string
	var detailScopes []out.DetailScopes

	_ = json.Unmarshal([]byte(dataGroup.Scope.String), &groupsScope)
	//scope := service.GenerateInitiateDataGroupDTOOut(groupsScope, groupsScope, true)

	// todo get detail scope
	for key := range groupsScope {
		var newGroupScope = make(map[string][]string)
		var detailScope out.DetailScopes
		detailScope.Key = key + ":" + service.GeneratePermissionKey(groupsScope[key])
		detailScope.Menu = key
		sort.Strings(groupsScope[key])
		if common.ValidateStringContainInStringArray(groupsScope[key], "all") {
			detailScope.Scope = append(detailScope.Scope, out.DetailScope{
				Label: key + ":" + groupsScope[key][sort.SearchStrings(groupsScope[key], "all")],
				Value: out.ScopeValue{
					Name: groupsScope[key][sort.SearchStrings(groupsScope[key], "all")],
				},
			})
		} else {
			newGroupScope[key] = groupsScope[key]
			switch key {
			case constanta.CustomerGroupDataScope:
				detailScope.Scope, err = CustomerGroupService.CustomerGroupService.GetListScopeCustomerGroup(newGroupScope)
				if err.Error != nil {
					return
				}
				break
			case constanta.CustomerCategoryDataScope:
				detailScope.Scope, err = CustomerCategoryService.CustomerCategoryService.GetListScopeCustomerCategory(newGroupScope)
				if err.Error != nil {
					return
				}
				break
			case constanta.ProductGroupDataScope:
				detailScope.Scope, err = ProductGroupService.ProductGroupService.GetListScopeProductGroup(newGroupScope)
				if err.Error != nil {
					return
				}
				break
			case constanta.ProvinceDataScope:
				detailScope.Scope, err = ProvinceService.ProvinceService.GetListScopeProvince(newGroupScope)
				if err.Error != nil {
					return
				}
				break
			case constanta.DistrictDataScope:
				detailScope.Scope, err = DistrictService.DistrictService.GetListScopeDistrict(newGroupScope)
				if err.Error != nil {
					return
				}
				break
			case constanta.ClientTypeDataScope:
				detailScope.Scope, err = ClientTypeService.ClientTypeService.GetListScopeClientType(newGroupScope)
				if err.Error != nil {
					return
				}
				break
			case constanta.SalesmanDataScope:
				detailScope.Scope, err = SalesmanService.SalesmanService.GetListScopeSalesman(newGroupScope)
				if err.Error != nil {
					return
				}
				break
			case constanta.EmployeeDataScope:
				detailScope.Scope, err = EmployeeService.EmployeeService.GetListScopeEmployee(newGroupScope)
				if err.Error != nil {
					return
				}
				break
			}
		}
		detailScopes = append(detailScopes, detailScope)
	}

	result = out.ViewDetailDataGroupResponse{
		ID:          dataGroup.ID.Int64,
		GroupID:     dataGroup.GroupID.String,
		Description: dataGroup.Description.String,
		Scope:       detailScopes,
		CreatedBy:   dataGroup.CreatedBy.Int64,
		UpdatedAt:   dataGroup.UpdatedAt.Time,
		CreatedAt:   dataGroup.CreatedAt.Time,
		UpdatedBy:   dataGroup.UpdatedBy.Int64,
		UpdatedName: dataGroup.UpdatedName.String,
	}

	return
}

func (input dataGroupService) validateView(inputStruct *in.DataGroupRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewDataGroup()
}
