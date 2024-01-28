package CustomerService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"time"
)

func (input customerService) InternalGetListDistributor(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	input.ValidOrderBy = []string{"id"}
	input.ValidSearchBy = []string{
		"license_variant",
		"client_type",
		"updated_at",
		"updated_at_start",
		"updated_at_end",
	}

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListMasterDistributorValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	//--- Validate Time Updated
	input.validateTimeUpdated(&inputStruct, &searchByParam)
	output.Data.Content, err = input.doGetListDistributor(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input customerService) doGetListDistributor(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, _ *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var dbResult []interface{}
	dbResult, err = dao.CustomerDAO.InternalGetListDistributor(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam)
	if err.Error != nil {
		return
	}

	output = input.convertToInternalGetLisDistributorResponse(dbResult)
	return
}

func (input customerService) validateTimeUpdated(inputStruct *in.GetListDataDTO, searchByParam *[]in.SearchByParam) errorModel.ErrorModel {
	var (
		fileName = "InternalGetListDistributor.go"
		funcName = "validateTimeUpdated"
		errS     error
	)

	for i := 0; i < len(*searchByParam); i++ {
		v := (*searchByParam)[i]
		if v.SearchKey == "updated_at" {
			_, errS = time.Parse(constanta.DefaultInstallationTimeFormat, v.SearchValue)
			if errS != nil {
				return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.UpdatedAt)
			}
		}

		if v.SearchKey == "updated_at_start" {
			inputStruct.UpdatedAtStartString = v.SearchValue
			inputStruct.UpdatedAtStart, errS = time.Parse(constanta.DefaultInstallationTimeFormat, v.SearchValue)
			if errS != nil {
				return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.UpdatedAtStart)
			}
			*searchByParam = append((*searchByParam)[:i], (*searchByParam)[i+1:]...)
			i = -1
			continue
		}

		if v.SearchKey == "updated_at_end" {
			inputStruct.UpdatedAtEndString = v.SearchValue
			inputStruct.UpdatedAtEnd, errS = time.Parse(constanta.DefaultInstallationTimeFormat, v.SearchValue)
			inputStruct.UpdatedAtEnd = inputStruct.UpdatedAtEnd.Add(24 * time.Hour).Add(-1 * time.Second)
			if errS != nil {
				return errorModel.GenerateFormatFieldError(fileName, funcName, constanta.UpdatedAtEnd)
			}
			*searchByParam = append((*searchByParam)[:i], (*searchByParam)[i+1:]...)
			i = -1
			continue
		}
	}

	return errorModel.GenerateNonErrorModel()
}

func (input customerService) convertToInternalGetLisDistributorResponse(dataOnDB []interface{}) (output []out.InternalGetListDistributorResponse) {
	for _, data := range dataOnDB {
		customerItem := data.(repository.CustomerModel)
		output = append(output, out.InternalGetListDistributorResponse{
			ClientID:            customerItem.ClientID.String,
			AuthUserID:          customerItem.AuthUserID.Int64,
			MDBCompanyProfileId: customerItem.MDBCompanyProfileID.Int64,
			ClientTypeName:      customerItem.ClientType.String,
			DistributorID:       customerItem.ID.Int64,
			PrincipalID:         customerItem.PrincipalID.Int64,
			DistTitle:           customerItem.CompanyTitle.String,
			DistName:            customerItem.CustomerName.String,
			DistNPWP:            customerItem.Npwp.String,
			DistAddress:         customerItem.Address.String,
			DistHamlet:          customerItem.Hamlet.String,
			DistNeighbourhood:   customerItem.Neighbourhood.String,
			DistCountry:         customerItem.CountryID.Int64,
			DistProvince:        customerItem.ProvinceID.Int64,
			DistDistrict:        customerItem.DistrictID.Int64,
			DistSubdistrict:     customerItem.SubDistrictID.Int64,
			DistUrbanvillage:    customerItem.UrbanVillageID.Int64,
			DistPostalcode:      customerItem.PostalCodeID.Int64,
			DistLicenseVariant:  customerItem.LicenseVariantName.String,
			Longitude:           customerItem.Longitude.Float64,
			Latitude:            customerItem.Latitude.Float64,
			DistPhone:           customerItem.Phone.String,
			DistFax:             customerItem.Fax.String,
			DistEmail:           customerItem.CompanyEmail.String,
			DistJoindate:        customerItem.InstallationDate.Time,
			DistFromdate:        customerItem.ProductValidFrom.Time,
			DistExpirydate:      customerItem.ProductValidThru.Time,
			CompanyID:           customerItem.UniqueID1.String,
			BranchID:            customerItem.UniqueID2.String,
			ActivationDate:      customerItem.ActivationDate.Time,
			UpdatedAt:           customerItem.UpdatedAt.Time,
		})
	}

	return
}
