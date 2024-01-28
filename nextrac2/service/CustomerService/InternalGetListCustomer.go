package CustomerService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input customerService) InternalGetListCustomer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	inputStruct, err := input.readBodyAndValidate(request, contextModel, input.validateInternalListProvince)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doInternalGetListCustomer(inputStruct)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input customerService) doInternalGetListCustomer(inputStruct in.CustomerRequest) (output interface{}, err errorModel.ErrorModel) {
	getListData, searchByParam, err := input.convertCustomerDTOToGetListDTO(inputStruct)
	if err.Error != nil {
		return
	}

	err = getListData.ValidateUpdatedAtRange()
	if err.Error != nil {
		return
	}

	dbResult, err := dao.CustomerDAO.InternalGetListCustomer(serverconfig.ServerAttribute.DBConnection, getListData, searchByParam)
	if err.Error != nil {
		return
	}

	output = input.convertToInternalCustomerGetListResponse(dbResult)

	return
}

func (input customerService) convertCustomerDTOToGetListDTO(inputStruct in.CustomerRequest) (getListData in.GetListDataDTO, searchByParam []in.SearchByParam, err errorModel.ErrorModel) {
	getListData.Page = inputStruct.Page
	getListData.Limit = inputStruct.Limit
	getListData.OrderBy = inputStruct.OrderBy
	getListData.UpdatedAtStartString = inputStruct.UpdatedAtStart
	getListData.UpdatedAtEndString = inputStruct.UpdatedAtEnd

	//Input valid search by
	return
}

func (input customerService) validateInternalListProvince(inputStruct *in.CustomerRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInputPageLimitAndOrderBy(input.ValidLimit, input.ValidOrderBy)
}

func (input customerService) convertToInternalCustomerGetListResponse(dataOnDB []interface{}) (output []out.InternalGetListCustomerResponse) {
	for _, data := range dataOnDB {
		customerItem := data.(repository.CustomerModel)
		output = append(output, out.InternalGetListCustomerResponse{
			ID:                  customerItem.ID.Int64,
			MDBCompanyProfileId: customerItem.MDBCompanyProfileID.Int64,
			NPWP:                customerItem.Npwp.String,
			IsPrincipal:         customerItem.IsPrincipal.Bool,
			IsParent:            customerItem.IsParent.Bool,
			CompanyTitle:        customerItem.CompanyTitle.String,
			CustomerName:        customerItem.CustomerName.String,
			Address:             customerItem.Address.String,
			Phone:               customerItem.Phone.String,
			CompanyEmail:        customerItem.CompanyEmail.String,
		})
	}
	return
}

func (input customerService) InternalCountCustomer(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	inputStruct, err := input.readBodyAndValidate(request, contextModel, input.validateInternalCountCustomer)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doInternalCountCustomer(inputStruct)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INITIATE_MESSAGE", contextModel)
	return
}

func (input customerService) doInternalCountCustomer(inputStruct in.CustomerRequest) (output int, err errorModel.ErrorModel) {
	getListData, searchByParam, err := input.convertCustomerDTOToGetListDTO(inputStruct)
	if err.Error != nil {
		return
	}

	err = getListData.ValidateUpdatedAtRange()
	if err.Error != nil {
		return
	}

	output, err = dao.CustomerDAO.InternalGetCountCustomer(serverconfig.ServerAttribute.DBConnection, searchByParam, getListData)
	if err.Error != nil {
		return
	}
	return
}

func (input customerService) validateInternalCountCustomer(_ *in.CustomerRequest) (err errorModel.ErrorModel) {
	return errorModel.GenerateNonErrorModel()
}
