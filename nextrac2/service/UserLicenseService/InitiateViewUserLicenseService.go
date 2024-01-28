package UserLicenseService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input userLicenseService) InitiateViewUserLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	var inputStruct in.ViewUserLicenseRequest
	var countData interface{}
	var ValidOrderByUserRegDetail = []string{
		"user_id",
		"salesman_id",
		"salesman_id",
		"email",
		"no_telp",
		"salesman_category",
		"reg_date",
		"android_id",
		"status",
	}

	_, inputStruct, searchByParam, err = input.readCountDataUserRegistrationDetail(request, input.ValidSearchBy, nil)
	if err.Error != nil {
		return
	}

	countData, err = dao.UserRegistrationDetailDAO.GetCountUserRegistrationDetail(serverconfig.ServerAttribute.DBConnection, searchByParam, inputStruct, contextModel.LimitedByCreatedBy)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  ValidOrderByUserRegDetail,
		ValidSearchBy: nil,
		ValidLimit:    input.ValidLimit,
		ValidOperator: nil,
		CountData:     countData.(int),
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_VIEW_USER_LICENSE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
