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

func (input userLicenseService) ViewDetailUserLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	parameterGetList, userLicenseViewStruct := input.readBodyCustomUserLicenseView(request)

	output.Data.Content, err = input.doViewDetailUserLicense(request, userLicenseViewStruct, contextModel, parameterGetList)
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

func (input userLicenseService) readBodyAndValidateForView(request *http.Request, validation func(inputStruct *in.UserLicenseRequest) errorModel.ErrorModel) (inputStruct in.UserLicenseRequest, err errorModel.ErrorModel) {
	id, _ := strconv.Atoi(mux.Vars(request)["id"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func (input userLicenseService) validateView(inputStruct *in.UserLicenseRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewDetailUserLicense()
}

func (input userLicenseService) doViewDetailUserLicense(request *http.Request, inputStruct in.ViewUserLicenseRequest, contextModel *applicationModel.ContextModel, parameterGetList in.GetListDataDTO) (output out.UserLicenseDetailResponse, err errorModel.ErrorModel) {
	fileName := "ViewDetailUserLicense.go"
	funcName := "doViewDetailUserLicense"
	var userLicenseOnDB repository.UserLicenseModel
	var listUserRegistrationDetailOnDB []repository.UserRegistrationDetailModel

	userLicenseParam := repository.UserLicenseModel{
		ID: sql.NullInt64{Int64: inputStruct.UserLicenseId},
	}

	userLicenseParam.CreatedBy.Int64 = 0

	userLicenseOnDB, err = dao.UserLicenseDAO.ViewDetailUserLicense(serverconfig.ServerAttribute.DBConnection, userLicenseParam)
	if err.Error != nil {
		return
	}

	if userLicenseOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.UserLicense)
		return
	}

	// GET LIST USER REGISTRATION DETAIL BY USER_LICENSE_ID
	listUserRegistrationDetailOnDB, err = input.getListUserRegistrationDetail(request, contextModel, parameterGetList)
	if err.Error != nil {
		return
	}

	output = input.convertUserLicenseOnBDToDTOOut(userLicenseOnDB, listUserRegistrationDetailOnDB)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userLicenseService) getListUserRegistrationDetail(request *http.Request, contextModel *applicationModel.ContextModel, parameterGetList in.GetListDataDTO) (listUserRegistrationDetailOnDB []repository.UserRegistrationDetailModel, err errorModel.ErrorModel) {
	var viewUserLicenseRequest in.ViewUserLicenseRequest

	var validOrderBy = []string{
		"user_id",
		"salesman_id",
		"email",
		"no_telp",
		"salesman_category",
		"reg_date",
		"android_id",
		"status",
	}

	parameterGetList, viewUserLicenseRequest, _, err = input.readGetListUserRegistrationDetail(request, nil, validOrderBy, nil, input.ValidLimit)
	if err.Error != nil {
		return
	}

	listUserRegistrationDetailOnDB, err = input.doGetListUserRegistrationDetail(parameterGetList, viewUserLicenseRequest, contextModel)
	if err.Error != nil {
		return
	}

	return
}

func (input userLicenseService) doGetListUserRegistrationDetail(parameterGetList in.GetListDataDTO, viewUserLicenseRequest in.ViewUserLicenseRequest, contextModel *applicationModel.ContextModel) (output []repository.UserRegistrationDetailModel, err errorModel.ErrorModel) {

	dbResult, err := dao.UserRegistrationDetailDAO.GetListUserRegistrationDetail(serverconfig.ServerAttribute.DBConnection, parameterGetList, contextModel.LimitedByCreatedBy, viewUserLicenseRequest)
	if err.Error != nil {
		return
	}

	output = input.convertUserRegistrationDetailToModel(dbResult)

	return
}

func (input userLicenseService) convertUserRegistrationDetailToModel(dbResult []interface{}) (listUserRegistrationDetailOnDB []repository.UserRegistrationDetailModel) {
	for _, resultItem := range dbResult {
		item := resultItem.(repository.UserRegistrationDetailModel)
		listUserRegistrationDetailOnDB = append(listUserRegistrationDetailOnDB, repository.UserRegistrationDetailModel{
			UserRegDetailID:  item.UserRegDetailID,
			UserID:           item.UserID,
			SalesmanId:       item.SalesmanId,
			Email:            item.Email,
			NoTelp:           item.NoTelp,
			SalesmanCategory: item.SalesmanCategory,
			RegDate:          item.RegDate,
			AndroidID:        item.AndroidID,
			Status:           item.Status,
		})
	}

	return listUserRegistrationDetailOnDB
}

func (input userLicenseService) convertUserLicenseOnBDToDTOOut(userLicenseOnDB repository.UserLicenseModel, userRegistrationDetailOnDB []repository.UserRegistrationDetailModel) out.UserLicenseDetailResponse {
	var listUserRegistrationDetailResponse []out.UserRegistrationDetailResponse

	for _, itemUserRegDetail := range userRegistrationDetailOnDB {
		listUserRegistrationDetailResponse = append(listUserRegistrationDetailResponse, out.UserRegistrationDetailResponse{
			UserRegistrationDetailID: itemUserRegDetail.UserRegDetailID.Int64,
			UserId:                   itemUserRegDetail.UserID.String,
			SalesmanId:               itemUserRegDetail.SalesmanId.String,
			Email:                    itemUserRegDetail.Email.String,
			NoTelp:                   itemUserRegDetail.NoTelp.String,
			SalesmanCategory:         itemUserRegDetail.SalesmanCategory.String,
			RegDate:                  itemUserRegDetail.RegDate.Time,
			AndroidId:                itemUserRegDetail.AndroidID.String,
			Status:                   itemUserRegDetail.Status.String,
		})
	}

	return out.UserLicenseDetailResponse{
		ID:                      userLicenseOnDB.ID.Int64,
		LicenseConfigId:         userLicenseOnDB.LicenseConfigId.Int64,
		InstallationId:          userLicenseOnDB.InstallationId.Int64,
		ParentCustomerId:        userLicenseOnDB.ParentCustomerId.Int64,
		ParentCustomer:          userLicenseOnDB.ParentCustomerName.String,
		CustomerId:              userLicenseOnDB.CustomerId.Int64,
		SiteId:                  userLicenseOnDB.SiteId.Int64,
		CustomerName:            userLicenseOnDB.CustomerName.String,
		UniqueId1:               userLicenseOnDB.UniqueId1.String,
		UniqueId2:               userLicenseOnDB.UniqueId2.String,
		ProductName:             userLicenseOnDB.ProductName.String,
		TotalActivated:          userLicenseOnDB.TotalActivated.Int64,
		TotalLicense:            userLicenseOnDB.TotalLicense.Int64,
		LicenseValidFrom:        userLicenseOnDB.ProductValidFrom.Time,
		LicenseValidThru:        userLicenseOnDB.ProductValidThru.Time,
		UpdatedAt:               userLicenseOnDB.UpdatedAt.Time,
		UserRegistrationDetails: listUserRegistrationDetailResponse,
	}
}
