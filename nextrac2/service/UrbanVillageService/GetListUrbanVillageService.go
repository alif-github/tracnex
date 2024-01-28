package UrbanVillageService

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
	"strconv"
)

func (input urbanVillageService) InitiateUrbanVillage(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData interface{}

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListUrbanVillageValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateUrbanVillage(searchByParam, contextModel)

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListUrbanVillageValidOperator,
		CountData:     countData.(int),
	}
	return
}

func (input urbanVillageService) doInitiateUrbanVillage(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	output, err = dao.UrbanVillageDAO.GetCountUrbanVillage(serverconfig.ServerAttribute.DBConnection, searchByParam, repository.UrbanVillageModel{
		CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	})

	if err.Error != nil {
		return 0, err
	}

	return
}

func (input urbanVillageService) GetListUrbanVillage(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListUrbanVillageValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListUrbanVillage(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input urbanVillageService) doGetListUrbanVillage(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	funcName := "doGetListUrbanVillage"
	var dbResult []interface{}

	for _, param := range searchByParam {
		if param.SearchKey == "sub_district_id" {
			var subDistrictOnDB repository.SubDistrictModel
			idSubDistrict, _ := strconv.Atoi(param.SearchValue)
			subDistrictOnDB, err = dao.SubDistrictDAO.GetSubDistrictByIDForGetList(serverconfig.ServerAttribute.DBConnection, int64(idSubDistrict), 0, false)
			if err.Error != nil {
				return
			}

			if subDistrictOnDB.ID.Int64 < 1 {
				err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.SubDistrict)
				return
			}
		}
	}

	dbResult, err = dao.UrbanVillageDAO.GetListUrbanVillage(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, repository.UrbanVillageModel{
		CreatedBy: sql.NullInt64{Int64: 0},
	})

	output = input.convertModelToResponseList(dbResult)

	return
}

func (input urbanVillageService) convertModelToResponseList(dbResult []interface{}) (result []out.UrbanVillageResponse) {
	for _, resultItem := range dbResult {
		item := resultItem.(repository.UrbanVillageModel)
		result = append(result, out.UrbanVillageResponse{
			ID:            item.ID.Int64,
			SubDistrictID: item.SubDistrictID.Int64,
			Code:          item.Code.String,
			Name:          item.Name.String,
			Status:        item.Status.String,
			CreatedBy:     item.CreatedBy.Int64,
			UpdatedAt:     item.UpdatedAt.Time,
		})
	}

	return result
}
