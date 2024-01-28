package PersonProfileService

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

func (input personProfileService) ViewPersonProfile(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct master_data_request.PersonProfileGetListRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateViewPersonProfile)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.DoViewPersonProfile(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_GET_LIST_MESSAGE", contextModel)
	return
}

func (input personProfileService) DoViewPersonProfile(inputStruct master_data_request.PersonProfileGetListRequest, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	funcName := "DoViewPersonProfile"
	personProfileOnMDB, err := master_data_dao.GetViewPersonProfileFromMasterData(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	idProvince, err := dao.ProvinceDAO.GetProvinceIDByMdbID(serverconfig.ServerAttribute.DBConnection, repository.ProvinceModel{
		MDBProvinceID: sql.NullInt64{Int64: personProfileOnMDB.ProvinceID},
	}, false)
	if err.Error != nil {
		return
	}

	idDistrict, err := dao.DistrictDAO.GetDistrictIDByMdbID(serverconfig.ServerAttribute.DBConnection, repository.DistrictModel{
		MdbDistrictID: sql.NullInt64{Int64: personProfileOnMDB.DistrictID},
		ProvinceID:    sql.NullInt64{Int64: personProfileOnMDB.ProvinceID},
	})
	if err.Error != nil {
		return
	}

	personProfileOnMDB.ProvinceID = idProvince
	personProfileOnMDB.DistrictID = idDistrict

	finalResponse := input.convertMDBResponseToDTOOut(personProfileOnMDB)
	if finalResponse.ID < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.PersonProfile)
		return
	}

	output = finalResponse
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input personProfileService) convertMDBResponseToDTOOut(inputStruct master_data_response.ViewPersonProfileResponse) out.ViewDetailPersonProfileResponse {
	return out.ViewDetailPersonProfileResponse{
		ID:              inputStruct.ID,
		Nik:             inputStruct.Nik,
		PersonTitleID:   inputStruct.PersonTitleID,
		PersonTitleName: inputStruct.PersonTitleName,
		FirstName:       inputStruct.FirstName,
		LastName:        inputStruct.LastName,
		Sex:             inputStruct.Sex,
		Address:         inputStruct.Address1,
		Address2:        inputStruct.Address2,
		Address3:        inputStruct.Address3,
		Hamlet:          inputStruct.Hamlet,
		Neighbourhood:   inputStruct.Neighbourhood,
		ProvinceID:      inputStruct.ProvinceID,
		ProvinceName:    inputStruct.ProvinceName,
		DistrictID:      inputStruct.DistrictID,
		DistrictName:    inputStruct.DistrictName,
		Phone:           inputStruct.Phone,
		Email:           inputStruct.Email,
	}
}

func (input personProfileService) ValidateViewPersonProfile(inpustStruct *master_data_request.PersonProfileGetListRequest) errorModel.ErrorModel {
	return inpustStruct.ValidateView()
}
