package CustomerInstallationService

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

func (input customerInstallationService) ViewCustomerSiteInInstallationService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.CustomerSiteInstallationRequest

	inputStruct, err = input.readBodyAndValidateForViewSite(request, input.validateViewSite)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewCustomerSite(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_CUSTOMER_INSTALLATION_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) ViewCustomerInstallationInInstallationService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.CustomerInstallationDetailRequest

	inputStruct, err = input.readBodyAndValidateForViewInstallation(request, input.validateViewInstallation)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewCustomerInstallation(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_CUSTOMER_INSTALLATION_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) doViewCustomerSite(inputStruct in.CustomerSiteInstallationRequest, contextModel *applicationModel.ContextModel) (output out.CustomerSiteInstallation, err errorModel.ErrorModel) {
	var (
		resultColCustomerSite []out.CustomerSite
		resultDBCustomer      repository.CustomerModel
		outputTemp            out.CustomerSiteInstallation
		db                    = serverconfig.ServerAttribute.DBConnection
		scope                 map[string]interface{}
		fileName              = "ViewCustomerInstallationService.go"
		funcName              = "doViewCustomerSite"
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	resultColCustomerSite, err = dao.CustomerSiteDAO.ViewDetailCustomerSite(db, repository.CustomerInstallationModel{ParentCustomerID: sql.NullInt64{Int64: inputStruct.ParentCustomerID}}, inputStruct.Page, inputStruct.Limit, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	resultDBCustomer, err = dao.CustomerDAO.GetNameCustomer(db, repository.CustomerModel{ID: sql.NullInt64{Int64: inputStruct.ParentCustomerID}}, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if resultDBCustomer.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Customer)
		return
	}

	outputTemp.ParentCustomerID = resultDBCustomer.ID.Int64
	outputTemp.ParentCustomerName = resultDBCustomer.CustomerName.String
	outputTemp.CustomerSite = resultColCustomerSite

	for index := range outputTemp.CustomerSite {
		outputTemp.CustomerSite[index].UpdatedAt += "Z"
	}

	output = outputTemp
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) doViewCustomerInstallation(inputStruct in.CustomerInstallationDetailRequest, contextModel *applicationModel.ContextModel) (output []out.CustomerInstallationDetailList, err errorModel.ErrorModel) {
	var (
		resultOnDB []repository.CustomerInstallationDetail
		scope      map[string]interface{}
		db         = serverconfig.ServerAttribute.DBConnection
		ciModel    repository.CustomerInstallationData
		cidModel   repository.CustomerInstallationDetail
	)

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	ciModel, cidModel = input.inputModelView(inputStruct)
	resultOnDB, err = dao.CustomerInstallationDAO.ViewCustomerInstallation(db, ciModel, cidModel, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	for _, valueResultOnDb := range resultOnDB {
		output = append(output, out.CustomerInstallationDetailList{
			InstallationID:        valueResultOnDb.InstallationID.Int64,
			ProductID:             valueResultOnDb.ProductID.Int64,
			ClientTypeID:          valueResultOnDb.ClientTypeID.Int64,
			ClientTypeDependantID: valueResultOnDb.ParentClientTypeID.Int64,
			ProductGroupID:        valueResultOnDb.ProductGroupID.Int64,
			ProductCode:           valueResultOnDb.ProductCode.String,
			ProductName:           valueResultOnDb.ProductName.String,
			ProductDescription:    valueResultOnDb.ProductDescription.String,
			Remark:                valueResultOnDb.Remark.String,
			UniqueID1:             valueResultOnDb.UniqueID1.String,
			UniqueID2:             valueResultOnDb.UniqueID2.String,
			InstallationDate:      valueResultOnDb.InstallationDate.Time,
			InstallationStatus:    valueResultOnDb.InstallationStatus.String,
			ProductValidFrom:      valueResultOnDb.ProductValidFrom.Time,
			ProductValidThru:      valueResultOnDb.ProductValidThru.Time,
			DayRange:              valueResultOnDb.DayRange.Int64,
			UpdatedAt:             valueResultOnDb.UpdatedAt.Time,
		})
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) validateViewSite(inputStruct *in.CustomerSiteInstallationRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewCustomerSite()
}

func (input customerInstallationService) validateViewInstallation(inputStruct *in.CustomerInstallationDetailRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewCustomerInstallation()
}

func (input customerInstallationService) inputModelView(inputStruct in.CustomerInstallationDetailRequest) (ciModel repository.CustomerInstallationData, cidModel repository.CustomerInstallationDetail) {
	ciModel = repository.CustomerInstallationData{
		CustomerID: sql.NullInt64{Int64: inputStruct.ParentCustomerID},
		SiteID:     sql.NullInt64{Int64: inputStruct.SiteID},
	}

	cidModel = repository.CustomerInstallationDetail{
		ClientTypeID: sql.NullInt64{Int64: inputStruct.ClientTypeID},
		IsLicense:    sql.NullBool{Bool: inputStruct.IsLicense},
	}

	return
}
