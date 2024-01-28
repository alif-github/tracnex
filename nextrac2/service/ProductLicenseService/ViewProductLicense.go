package ProductLicenseService

import (
	"database/sql"
	"github.com/gorilla/mux"
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
	"strconv"
)

func (input productLicenseService) ViewProductLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ProductLicenseRequest

	inputStruct, err = input.readBodyAndValidateForViewProductLicense(request, input.validateViewProductLicense)
	if err.Error != nil {
		return
	}

	err = input.checkExistsProductLicense(inputStruct.ID)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewProductLicenseById(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_PRODUCT_LICENSE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) readBodyAndValidateForViewProductLicense(request *http.Request, validation func(input *in.ProductLicenseRequest) errorModel.ErrorModel) (inputStruct in.ProductLicenseRequest, err errorModel.ErrorModel) {
	ProductLicenseId, _ := strconv.Atoi(mux.Vars(request)["id"])

	inputStruct.ID = int64(ProductLicenseId)

	err = validation(&inputStruct)

	return
}

func (input productLicenseService) validateViewProductLicense(inputStruct *in.ProductLicenseRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewProductLicense()
}

func (input productLicenseService) checkExistsProductLicense(idProductLicense int64) (err errorModel.ErrorModel) {
	fileName := "ViewProductLicense.go"
	funcName := "checkExistsProductLicense"

	productLicenseExists, err := dao.ProductLicenseDAO.CheckExistsProductLicenseById(serverconfig.ServerAttribute.DBConnection, idProductLicense)

	if productLicenseExists == false {
		return errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ProductLicenseID)
	}

	return errorModel.GenerateNonErrorModel()
}

func (input productLicenseService) doViewProductLicenseById(inputStruct in.ProductLicenseRequest, contextModel *applicationModel.ContextModel) (output out.ProductLicenseDetailResponse, err errorModel.ErrorModel) {

	var (
		listComponentDB       []repository.ProductComponentModel
		DetailProductLicense  repository.DetailProductLicense
		listComponentResponse []out.ListProductComponentProductLicense
		scope                 map[string]interface{}
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	DetailProductLicense, err = dao.ProductLicenseDAO.ViewDetailProductLicenseById(serverconfig.ServerAttribute.DBConnection, repository.DetailProductLicense{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	}, scope, input.MappingScopeDB, input.ListScope)
	if err.Error != nil {
		return
	}

	if DetailProductLicense.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, "doViewProductLicenseById", constanta.ID)
		return
	}

	for _, component := range DetailProductLicense.Components {
		listComponentResponse = append(listComponentResponse, out.ListProductComponentProductLicense{
			ComponentID:    component.ComponentID.Int64,
			ComponentName:  component.ComponentName.String,
			ComponentValue: component.ComponentValue.String,
		})
	}

	DetailProductLicense.Components = listComponentDB

	result := out.ProductLicenseDetailResponse{
		LicenseID:              DetailProductLicense.ID.Int64,
		ProductKey:             DetailProductLicense.ProductKey.String,
		ActivationDate:         DetailProductLicense.ActivationDate.Time,
		LicenseStatus:          DetailProductLicense.LicenseStatus.Int32,
		TerminationDescription: DetailProductLicense.TerminationDescription.String,
		LicenseConfigId:        DetailProductLicense.LicenseConfigId.Int64,
		InstallationId:         DetailProductLicense.InstallationId.Int64,
		ParentCustomerId:       DetailProductLicense.ParentCustomerId.Int64,
		ParentCustomer:         DetailProductLicense.ParentCustomer.String,
		CustomerId:             DetailProductLicense.CustomerId.Int64,
		SiteId:                 DetailProductLicense.SiteId.Int64,
		Customer:               DetailProductLicense.Customer.String,
		ClientId:               DetailProductLicense.ClientId.String,
		Product:                DetailProductLicense.Product.String,
		Client:                 DetailProductLicense.Client.String,
		LicenseVariant:         DetailProductLicense.LicenseVariant.String,
		LicenseType:            DetailProductLicense.LicenseType.String,
		DeploymentMethod:       DetailProductLicense.DeploymentMethod.String,
		NumberOfUser:           DetailProductLicense.NumberOfUser.Int64,
		ConcurentUser:          DetailProductLicense.ConcurentUser.String,
		UniqueId1:              DetailProductLicense.UniqueId1.String,
		UniqueId2:              DetailProductLicense.UniqueId2.String,
		LicenseValidFrom:       DetailProductLicense.LicenseValidFrom.Time,
		LicenseValidThru:       DetailProductLicense.LicenseValidThru.Time,
		Created:                DetailProductLicense.CreatedAt.Time,
		Modified:               DetailProductLicense.UpdatedAt.Time,
		ModifiedBy:             DetailProductLicense.AliasName.String,
		Module1:                DetailProductLicense.Module1.String,
		Module2:                DetailProductLicense.Module2.String,
		Module3:                DetailProductLicense.Module3.String,
		Module4:                DetailProductLicense.Module4.String,
		Module5:                DetailProductLicense.Module5.String,
		Module6:                DetailProductLicense.Module6.String,
		Module7:                DetailProductLicense.Module7.String,
		Module8:                DetailProductLicense.Module8.String,
		Module9:                DetailProductLicense.Module9.String,
		Module10:               DetailProductLicense.Module10.String,
		Components:             listComponentResponse,
	}

	output = result
	err = errorModel.GenerateNonErrorModel()
	return
}
