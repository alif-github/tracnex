package SubDistrictService

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

func (input subDistrictService) InitiateSubDistrict(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData interface{}

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListSubDistrictValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiateSubDistrict(searchByParam, contextModel)

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListSubDistrictValidOperator,
		CountData:     countData.(int),
	}
	return
}

func (input subDistrictService) doInitiateSubDistrict(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	output, err = dao.SubDistrictDAO.GetCountSubDistrict(serverconfig.ServerAttribute.DBConnection, searchByParam, repository.SubDistrictModel{})
	if err.Error != nil {
		return 0, err
	}

	return
}

func (input subDistrictService) GetListSubDistrict(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListSubDistrictValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListSubDistrict(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input subDistrictService) doGetListSubDistrict(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	funcName := "doGetListSubDistrict"
	var dbResult []interface{}

	for _, param := range searchByParam {
		if param.SearchKey == "district_id" {
			var districtOnDB repository.DistrictModel
			idDistrict, _ := strconv.Atoi(param.SearchValue)
			districtOnDB, err = dao.DistrictDAO.GetDistrictByID(serverconfig.ServerAttribute.DBConnection, repository.DistrictModel{ID: sql.NullInt64{Int64: int64(idDistrict)}})
			if err.Error != nil {
				return
			}

			if districtOnDB.ID.Int64 < 1 {
				err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.District)
				return
			}
		}
	}

	dbResult, err = dao.SubDistrictDAO.GetListSubDistrict(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, repository.SubDistrictModel{})
	if err.Error != nil {
		return
	}

	output = input.convertModelToResponseList(dbResult)

	return
}

func (input subDistrictService) convertModelToResponseList(dbResult []interface{}) (result []out.SubDistrictResponse) {
	for _, resultItem := range dbResult {
		item := resultItem.(repository.SubDistrictModel)
		result = append(result, out.SubDistrictResponse{
			ID:         item.ID.Int64,
			DistrictID: item.DistrictID.Int64,
			Code:       item.Code.String,
			Name:       item.Name.String,
			Status:     item.Status.String,
			CreatedBy:  item.CreatedBy.Int64,
			UpdatedAt:  item.UpdatedAt.Time,
		})
	}
	return result
}
