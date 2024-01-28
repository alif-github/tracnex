package CompanyService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
)

type companyDetailService struct {
	FileName string
	service.AbstractService
}

var CompanyDetailService = companyDetailService{FileName: "CompanyDetailService.go"}

func (input companyDetailService) DetailCompany(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "DetailCompany"

	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	detail, err := dao.CompanyDAO.GetDetailCompany(serverconfig.ServerAttribute.DBConnection, id)
	if err.Error != nil {
		return
	}

	if detail.ID.Int64 == 0 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "company")
		return
	}

	output.Data.Content = input.getResponse(detail)
	output.Status = service.GetResponseMessages("SUCCESS_GET_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input companyDetailService) getResponse(model repository.CompanyModel) out.CompanyDetailResponse {
	file, _ := dao.FileUploadDAO.GetFileByParentIDAndCategory(serverconfig.ServerAttribute.DBConnection, model.ID.Int64, dao.CompanyDAO.TableName)
	return out.CompanyDetailResponse{
		ID:            model.ID.Int64,
		CompanyTitle:  model.CompanyTitle.String,
		CompanyName:   model.CompanyName.String,
		PhotoIcon:     file.Host.String + file.Path.String,
        Address:       model.Address.String,
        Address2:      model.Address2.String,
        Neighbourhood: model.Neighbourhood.String,
        Hamlet:        model.Hamlet.String,
        ProvinceId:    model.ProvinceID.Int64,
        ProvinceName:  model.ProvinceName.String,
        DistrictId:    model.DistrictID.Int64,
        DistrictName:  model.DistrictName.String,
        SubDistrictId: model.SubDistrictID.Int64,
        SubDistrictName: model.SubDistrictName.String,
        Village:       model.UrbanVillageName.String,
        VillageId:     model.UrbanVillageID.Int64,
        PostalCodeId:  model.PostalCodeID.Int64,
        PostalCode:    model.PostalCode.String,
        Longitude:     model.Longitude.String,
        Latitude:      model.Latitude.String,
        Telephone:     model.Telephone.String,
        TelephoneAlternate: model.AlternateTelephone.String,
        Fax    :       model.Fax.String,
        Email :        model.CompanyEmail.String,
        AlternateEmail: model.AlternativeCompanyEmail.String,
        Npwp:          model.Npwp.String,
        TaxName:       model.TaxName.String,
        TaxAddress:    model.TaxAddress.String,
        UpdatedAt:     model.UpdatedAt.Time,
	}
}
