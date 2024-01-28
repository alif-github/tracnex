package PostalCodeService

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

func (input postalCodeService) InitiatePostalCode(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var countData interface{}

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListPostalCodeValidOperator)
	if err.Error != nil {
		return
	}

	countData, err = input.doInitiatePostalCode(searchByParam, contextModel)

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListPostalCodeValidOperator,
		CountData:     countData.(int),
	}
	return
}

func (input postalCodeService) doInitiatePostalCode(searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	output, err = dao.PostalCodeDAO.GetCountPostalCode(serverconfig.ServerAttribute.DBConnection, searchByParam, repository.PostalCodeModel{})

	if err.Error != nil {
		return 0, err
	}

	return
}

func (input postalCodeService) GetListPostalCode(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListPostalCodeValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListPostalCode(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input postalCodeService) doGetListPostalCode(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	funcName := "doGetListPostalCode"
	var dbResult []interface{}

	for _, param := range searchByParam {
		if param.SearchKey == "urban_village_id" {
			var urbanVillageOnDB repository.UrbanVillageModel
			idUrbanVillage, _ := strconv.Atoi(param.SearchValue)
			urbanVillageOnDB, err = dao.UrbanVillageDAO.GetUrbanVillageByIDForGetList(serverconfig.ServerAttribute.DBConnection, int64(idUrbanVillage), false)
			if err.Error != nil {
				return
			}

			if urbanVillageOnDB.ID.Int64 < 1 {
				err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.UrbanVillage)
				return
			}
		}
	}

	dbResult, err = dao.PostalCodeDAO.GetListPostalCode(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, repository.PostalCodeModel{
		CreatedBy: sql.NullInt64{Int64: 0},
	})

	output = input.convertModelToResponseList(dbResult)

	return
}

func (input postalCodeService) convertModelToResponseList(dbResult []interface{}) (result []out.PostalCodeResponse) {
	for _, resultItem := range dbResult {
		item := resultItem.(repository.PostalCodeModel)
		result = append(result, out.PostalCodeResponse{
			ID:             item.ID.Int64,
			UrbanVillageID: item.UrbanVillageID.Int64,
			Code:           item.Code.String,
			Status:         item.Status.String,
			CreatedBy:      item.CreatedBy.Int64,
			UpdatedAt:      item.UpdatedAt.Time,
		})
	}

	return result
}
