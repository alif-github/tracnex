package CustomerListService

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
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input customerListService) ViewCustomer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.CustomerListRequest

	inputStruct, err = input.readBodyAndValidateForView(request, contextModel, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewCustomer(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_CUSTOMER_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListService) doViewCustomer(inputStruct in.CustomerListRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	fileName := "ViewCustomerListService.go"
	funcName := "doViewCustomer"

	customerListModel := repository.CustomerListModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	}

	customerListModel.CreatedBy.Int64 = contextModel.LimitedByCreatedBy

	customerListOnDB, err := dao.CustomerListDAO.ViewCustomer(serverconfig.ServerAttribute.DBConnection, customerListModel)
	if err.Error != nil {
		return
	}

	if customerListOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ID)
		return
	}

	result = input.convertDAOToDTO(customerListOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerListService) convertDAOToDTO(resultModel repository.CustomerListModel) (output interface{}) {
	return out.ViewCustomerResponse{
		CompanyID: 		resultModel.CompanyID.String,
		BranchID: 		resultModel.BranchID.String,
		CompanyName: 	resultModel.CompanyName.String,
		City: 			resultModel.City.String,
		Implementer: 	resultModel.Implementer.String,
		Implementation: resultModel.Implementation.Time,
		Product: 		resultModel.Product.String,
		Version: 		resultModel.Version.String,
		LicenseType: 	resultModel.LicenseType.String,
		UserAmount: 	resultModel.UserAmount.Int64,
		ExpDate: 		resultModel.ExpDate.Time,
		CreatedBy: 		resultModel.CreatedBy.Int64,
		CreatedAt: 		resultModel.CreatedAt.Time,
		CreatedClient: 	resultModel.CreatedClient.String,
		UpdatedBy: 		resultModel.UpdatedBy.Int64,
		UpdatedAt: 		resultModel.UpdatedAt.Time,
		UpdatedClient: 	resultModel.UpdatedClient.String,
	}
}

func (input customerListService) validateView(inputStruct *in.CustomerListRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewCustomer()
}