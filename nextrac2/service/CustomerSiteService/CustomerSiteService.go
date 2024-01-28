package CustomerSiteService

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type customerSiteService struct {
	service.AbstractService
	service.GetListData
}

var CustomerSiteService = customerSiteService{}.New()

func (input customerSiteService) New() (output customerSiteService) {
	output.FileName = "CustomerSiteService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
	}
	output.ValidSearchBy = []string{"id"}
	return
}

func (input customerSiteService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.CustomerSiteRequest) errorModel.ErrorModel) (inputStruct in.CustomerSiteRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if !util.IsStringEmpty(stringBody) {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func (input customerSiteService) validateRelation(customerID []int64, customerSiteModel *repository.CustomerSiteModel) (err errorModel.ErrorModel) {
	var (
		fileName = input.FileName
		funcName = "validateRelation"
	)

	for index, valueCustomerID := range customerID {
		var isExist bool
		isExist, err = dao.CustomerDAO.CheckCustomerIsExist(serverconfig.ServerAttribute.DBConnection, repository.CustomerModel{ID: sql.NullInt64{Int64: valueCustomerID}}, nil, nil)
		if err.Error != nil {
			return
		}

		if !isExist {
			if index == 0 {
				err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ParentCustomerID)
				return
			} else {
				err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.CustomerID)
				return
			}
		} else {
			if index == 0 {
				customerSiteModel.ParentCustomerID.Int64 = valueCustomerID
			} else {
				customerSiteModel.CustomerID.Int64 = valueCustomerID
			}
		}
	}

	return errorModel.GenerateNonErrorModel()
}
