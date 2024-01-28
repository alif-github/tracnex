package service

import (
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"strconv"
)

type GetListData struct {
	ValidSearchBy []string
	ValidOrderBy  []string
	ValidLimit    []int
}

var DefaultLimit = []int{
	10,
	20,
	50,
	100,
}

func (input GetListData) readGetListData(request *http.Request) (inputStruct in.GetListDataDTO) {
	inputStruct.Page, _ = strconv.Atoi(GenerateQueryValue(request.URL.Query()["page"]))
	inputStruct.Limit, _ = strconv.Atoi(GenerateQueryValue(request.URL.Query()["limit"]))
	inputStruct.Filter = GenerateQueryValue(request.URL.Query()["filter"])
	inputStruct.Search = GenerateQueryValue(request.URL.Query()["search"])
	inputStruct.OrderBy = GenerateQueryValue(request.URL.Query()["order"])

	return
}

func (input GetListData) ReadAndValidateGetListData(request *http.Request, validSearchKey []string, validOrderBy []string, validOperator map[string]applicationModel.DefaultOperator, validLimit []int) (inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, err errorModel.ErrorModel) {
	inputStruct = input.readGetListData(request)
	searchByParam, err = inputStruct.ValidateGetListData(validSearchKey, validOrderBy, validOperator, validLimit)
	return
}

func (input GetListData) ReadAndValidateGetCountData(request *http.Request, validSearchBy []string, validOperator map[string]applicationModel.DefaultOperator) (inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, err errorModel.ErrorModel) {
	inputStruct = input.readGetListData(request)
	searchByParam, err = inputStruct.ValidateGetCountData(validSearchBy, validOperator)
	return
}

func (input GetListData) ReadAndValidateGetListDataWithID(request *http.Request, validSearchKey []string, validOrderBy []string, validOperator map[string]applicationModel.DefaultOperator, validLimit []int) (inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, err errorModel.ErrorModel) {
	inputStruct = input.readGetListData(request)
	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}
	searchByParam, err = inputStruct.ValidateGetListDataWithID(inputStruct.ID, validSearchKey, validOrderBy, validOperator, validLimit)
	if err.Error != nil {
		return
	}

	return
}

func (input GetListData) ReadAndValidateGetCountDataWithID(request *http.Request, validSearchBy []string, validOperator map[string]applicationModel.DefaultOperator) (inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, err errorModel.ErrorModel) {
	inputStruct = input.readGetListData(request)
	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	searchByParam, err = inputStruct.ValidateGetCountDataWithID(inputStruct.ID, validSearchBy, validOperator)
	return
}

func (input GetListData) SetDefaultOrder(request *http.Request, orderBy string, inputStruct *in.GetListDataDTO, validOrderBy []string) (err errorModel.ErrorModel) {
	if util.IsStringEmpty(GenerateQueryValue(request.URL.Query()["order"])) {
		inputStruct.OrderBy = orderBy
		inputStruct.OrderBy, err = in.ValidateOrderBy(validOrderBy, inputStruct.OrderBy)
	}

	return
}
