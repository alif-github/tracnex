package SalesmanService

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
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input salesmanService) GetListAdminSalesman(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListSalesmanValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListSalesman(inputStruct, searchByParam, contextModel, true)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_SALESMAN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanService) GetListSalesman(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListSalesmanValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListSalesman(inputStruct, searchByParam, contextModel, false)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_SALESMAN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanService) InitiateGetListSalesman(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		countData      interface{}
		searchByParam  []in.SearchByParam
		genderSalesman []out.GenderSalesman
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListSalesmanValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateListSalesman(searchByParam, contextModel, false)
	if err.Error != nil {
		return
	}

	genderSalesman = append(genderSalesman, out.GenderSalesman{
		Code:       "L",
		GenderName: "Pria",
	}, out.GenderSalesman{
		Code:       "P",
		GenderName: "Wanita",
	}, out.GenderSalesman{
		Code:       "N",
		GenderName: "Tidak Terdefinisi",
	})

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListSalesmanValidOperator,
		EnumData:      genderSalesman,
		CountData:     countData.(int),
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_GET_LIST_SALESMAN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanService) doGetListSalesman(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel, isAdmin bool) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult        []interface{}
		scope           map[string]interface{}
		regionalDeleted bool
		db              = serverconfig.ServerAttribute.DBConnection
	)

	if isAdmin {
		scope = make(map[string]interface{})
		scope[constanta.SalesmanDataScope] = []interface{}{"all"}
		scope[constanta.ProvinceDataScope] = []interface{}{"all"}
		scope[constanta.DistrictDataScope] = []interface{}{"all"}
		contextModel.LimitedByCreatedBy = 0
		regionalDeleted = true
	} else {
		scope, err = input.validateDataScopeSalesman(contextModel)
		if err.Error != nil {
			return
		}
	}

	//--- Get List Salesman
	dbResult, err = dao.SalesmanDAO.GetListSalesman(db, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB, regionalDeleted)
	if err.Error != nil {
		return
	}

	output = input.convertToListDTOOut(dbResult)
	return
}

func (input salesmanService) doInitiateListSalesman(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel, isAdmin bool) (output interface{}, err errorModel.ErrorModel) {
	var scope map[string]interface{}

	scope, err = input.validateDataScopeSalesman(contextModel)
	if err.Error != nil {
		return
	}

	output, err = dao.SalesmanDAO.GetCountSalesman(serverconfig.ServerAttribute.DBConnection, searchByParam, contextModel.LimitedByCreatedBy, scope, input.MappingScopeDB)
	return
}

func (input salesmanService) convertToListDTOOut(dbResult []interface{}) (result []out.ListSalesman) {
	for _, dbResultItem := range dbResult {
		repo := dbResultItem.(repository.ListSalesmanModel)
		result = append(result, out.ListSalesman{
			ID:        repo.ID.Int64,
			FirstName: repo.FirstName.String,
			LastName:  repo.LastName.String,
			Status:    repo.Status.String,
			Address:   repo.Address.String,
			District:  repo.District.String,
			Province:  repo.Province.String,
			Phone:     repo.Phone.String,
			Email:     repo.Email.String,
			UpdatedAt: repo.UpdatedAt.Time,
		})
	}

	return result
}
