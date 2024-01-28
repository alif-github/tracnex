package CompanyProfileService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_request"
	"nexsoft.co.id/nextrac2/resource_master_data/dto/master_data_response"
	"nexsoft.co.id/nextrac2/resource_master_data/master_data_dao"
	"nexsoft.co.id/nextrac2/serverconfig"
)

func (input companyProfileService) ViewCompanyProfile(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct master_data_request.CompanyProfileGetListRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateViewCompanyProfile)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.DoViewCompanyProfile(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input companyProfileService) DoViewCompanyProfile(inputStruct master_data_request.CompanyProfileGetListRequest, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		funcName                      = "DoViewCompanyProfile"
		customerParentOnDB            repository.CustomerModel
		idProvince, idDistrict        int64
		idSubDistrict, idUrbanVillage int64
		idPostalCode                  int64
		companyProfileOnMDB           master_data_response.ViewCompanyProfileResponse
	)

	companyProfileOnMDB, err = master_data_dao.GetViewCompanyProfileFromMasterData(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	//--- Get region data from DB
	idProvince, err = dao.ProvinceDAO.GetProvinceIDByMdbID(serverconfig.ServerAttribute.DBConnection, repository.ProvinceModel{
		MDBProvinceID: sql.NullInt64{Int64: companyProfileOnMDB.ProvinceID},
	}, false)
	if err.Error != nil {
		return
	}

	if idProvince < 1 {
		companyProfileOnMDB.ProvinceName = ""
	}

	if idProvince > 0 {
		idDistrict, err = dao.DistrictDAO.GetDistrictIDByMdbID(serverconfig.ServerAttribute.DBConnection, repository.DistrictModel{
			MdbDistrictID: sql.NullInt64{Int64: companyProfileOnMDB.DistrictID},
			ProvinceID:    sql.NullInt64{Int64: companyProfileOnMDB.ProvinceID},
		})

		if err.Error != nil {
			return
		}

		if idDistrict < 1 {
			companyProfileOnMDB.DistrictName = ""
		}
	}

	if idDistrict > 0 {
		idSubDistrict, err = dao.SubDistrictDAO.GetSubDistrictIDByMdbID(serverconfig.ServerAttribute.DBConnection, repository.SubDistrictModel{
			MDBSubDistrictID: sql.NullInt64{Int64: companyProfileOnMDB.SubDistrictID},
			DistrictID:       sql.NullInt64{Int64: companyProfileOnMDB.DistrictID},
		})

		if err.Error != nil {
			return
		}

		if idSubDistrict < 1 {
			companyProfileOnMDB.SubDistrictName = ""
		}
	}

	if idSubDistrict > 0 {
		idUrbanVillage, err = dao.UrbanVillageDAO.GetUrbanVillageIDByMdbID(serverconfig.ServerAttribute.DBConnection, repository.UrbanVillageModel{
			MDBUrbanVillageID: sql.NullInt64{Int64: companyProfileOnMDB.UrbanVillageID},
			SubDistrictID:     sql.NullInt64{Int64: companyProfileOnMDB.SubDistrictID},
		})

		if err.Error != nil {
			return
		}

		if idUrbanVillage < 1 {
			companyProfileOnMDB.UrbanVillageName = ""
		}
	}

	if idUrbanVillage > 0 {
		idPostalCode, err = dao.PostalCodeDAO.GetPostalCodeIDByMdbID(serverconfig.ServerAttribute.DBConnection, repository.PostalCodeModel{
			MDBPostalCodeID: sql.NullInt64{Int64: companyProfileOnMDB.PostalCodeID},
			UrbanVillageID:  sql.NullInt64{Int64: companyProfileOnMDB.UrbanVillageID},
		})

		if err.Error != nil {
			return
		}

		if idPostalCode < 1 {
			companyProfileOnMDB.PostalCode = ""
		}
	}

	companyProfileOnMDB.ProvinceID = idProvince
	companyProfileOnMDB.DistrictID = idDistrict
	companyProfileOnMDB.SubDistrictID = idSubDistrict
	companyProfileOnMDB.UrbanVillageID = idUrbanVillage
	companyProfileOnMDB.PostalCodeID = idPostalCode

	finalResponse := input.convertMDBResponseToDTOOut(companyProfileOnMDB)
	if finalResponse.ID < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.CompanyProfile)
		return
	}

	if finalResponse.CompanyParent > 0 {
		customerParentOnDB, err = dao.CustomerDAO.GetCustomerByMdbCompanyProfile(serverconfig.ServerAttribute.DBConnection, repository.CustomerModel{
			MDBCompanyProfileID: sql.NullInt64{Int64: finalResponse.CompanyParent},
		})

		if err.Error != nil {
			return
		}

		if customerParentOnDB.ID.Int64 > 0 {
			finalResponse.IsParentCompany = false
			finalResponse.CompanyParentName = customerParentOnDB.CustomerName.String
			finalResponse.CustomerParentID = customerParentOnDB.ID.Int64
		}
	}

	output = finalResponse
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input companyProfileService) convertMDBResponseToDTOOut(inputStruct master_data_response.ViewCompanyProfileResponse) out.ViewCompanyProfileResponse {
	return out.ViewCompanyProfileResponse{
		ID:                inputStruct.ID,
		Npwp:              inputStruct.NPWP,
		Status:            inputStruct.Status,
		CompanyTitleID:    inputStruct.CompanyTitleID,
		CompanyTitle:      inputStruct.CompanyTitle,
		Name:              inputStruct.Name,
		CompanyParent:     inputStruct.CompanyParent,
		CompanyParentName: inputStruct.CompanyParentName,
		Address:           inputStruct.Address1,
		Address2:          inputStruct.Address2,
		Address3:          inputStruct.Address3,
		Hamlet:            inputStruct.Hamlet,
		Neighbourhood:     inputStruct.Neighbourhood,
		ProvinceID:        inputStruct.ProvinceID,
		ProvinceName:      inputStruct.ProvinceName,
		DistrictID:        inputStruct.DistrictID,
		DistrictName:      inputStruct.DistrictName,
		SubDistrictID:     inputStruct.SubDistrictID,
		SubDistrictName:   inputStruct.SubDistrictName,
		UrbanVillageID:    inputStruct.UrbanVillageID,
		UrbanVillageName:  inputStruct.UrbanVillageName,
		PostalCodeID:      inputStruct.PostalCodeID,
		PostalCode:        inputStruct.PostalCode,
		Latitude:          inputStruct.Latitude,
		Longitude:         inputStruct.Longitude,
		Phone:             inputStruct.Phone,
		Fax:               inputStruct.Fax,
		Email:             inputStruct.Email,
		AlternativeEmail:  inputStruct.AlternativeEmail,
		UpdatedAt:         inputStruct.UpdatedAt,
		CreatedBy:         inputStruct.CreatedBy,
	}
}

func (input companyProfileService) ValidateViewCompanyProfile(inputStruct *master_data_request.CompanyProfileGetListRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}
