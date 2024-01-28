package UserLicenseService

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

func (input userLicenseService) ViewTransferUserLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.ViewUserLicenseRequest

	inputStruct, err = input.readBodyAndValidateForViewTransferKeyUserLicense(request, input.validateViewTransferKeyUserLicense)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewTransferKeyUserLicense(inputStruct)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_USER_LICENSE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseService) validateViewTransferKeyUserLicense(inputStruct *in.ViewUserLicenseRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewTransferKeyUserLicense()
}

func (input userLicenseService) readBodyAndValidateForViewTransferKeyUserLicense(request *http.Request, validation func(inputStruct *in.ViewUserLicenseRequest) errorModel.ErrorModel) (inputStruct in.ViewUserLicenseRequest, err errorModel.ErrorModel) {
	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.UserLicenseId == 0 {
		inputStruct.UserLicenseId = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func (input userLicenseService) doViewTransferKeyUserLicense(inputStruct in.ViewUserLicenseRequest) (output out.UserLicenseTransferKeyResponse, err errorModel.ErrorModel) {
	fileName := "ViewTransferUserLicenseService.go"
	funcName := "doViewTransferKeyUserLicense"

	var userLicenseOnDB repository.UserLicenseModel

	userLicenseOnDB, err = dao.UserLicenseDAO.ViewDetailUserLicense(serverconfig.ServerAttribute.DBConnection, repository.UserLicenseModel{ID: sql.NullInt64{Int64: inputStruct.UserLicenseId}})
	if err.Error != nil {
		return
	}

	if userLicenseOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.UserLicense)
		return
	}

	//clientTypeId, err := dao.UserLicenseDAO.GetClientTypeById(serverconfig.ServerAttribute.DBConnection, userLicenseOnDB)
	//if err.Error != nil {
	//	return
	//}
	//
	//userLicenseOnDB.ClientType.Int64 = clientTypeId

	if userLicenseOnDB.ClientTypeId.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ClientTypeID)
		return
	}

	output = input.convertUserLicenseTransferKeyToDTOOut(userLicenseOnDB)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseService) convertUserLicenseTransferKeyToDTOOut(userLicenseOnDB repository.UserLicenseModel) out.UserLicenseTransferKeyResponse {
	return out.UserLicenseTransferKeyResponse{
		ID:               userLicenseOnDB.ID.Int64,
		LicenseConfigId:  userLicenseOnDB.LicenseConfigId.Int64,
		InstallationId:   userLicenseOnDB.InstallationId.Int64,
		ParentCustomerId: userLicenseOnDB.ParentCustomerId.Int64,
		ParentCustomer:   userLicenseOnDB.ParentCustomerName.String,
		CustomerId:       userLicenseOnDB.CustomerId.Int64,
		SiteId:           userLicenseOnDB.SiteId.Int64,
		CustomerName:     userLicenseOnDB.CustomerName.String,
		UniqueId1:        userLicenseOnDB.UniqueId1.String,
		UniqueId2:        userLicenseOnDB.UniqueId2.String,
		ProductName:      userLicenseOnDB.ProductName.String,
		TotalActivated:   userLicenseOnDB.TotalActivated.Int64,
		TotalLicense:     userLicenseOnDB.TotalLicense.Int64,
		LicenseValidFrom: userLicenseOnDB.ProductValidFrom.Time,
		LicenseValidThru: userLicenseOnDB.ProductValidThru.Time,
		UpdatedAt:        userLicenseOnDB.UpdatedAt.Time,
		ClientTypeId:     userLicenseOnDB.ClientTypeId.Int64,
	}
}
