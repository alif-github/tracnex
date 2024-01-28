package SubDistrictService

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
)

func (input subDistrictService) ViewSubDistrictService(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.SubDistrictRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewSubDistrict)
	if err.Error != nil {
		return
	}

	output.Data.Content, err =input.doViewSubDistrictService(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status =input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input subDistrictService) doViewSubDistrictService(inputStruct in.SubDistrictRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	funcName := "doViewSubDistrictService"

	subDistrictOnDB, err := dao.SubDistrictDAO.ViewSubDistrict(serverconfig.ServerAttribute.DBConnection, repository.SubDistrictModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	})
	if err.Error != nil {
		return
	}

	if subDistrictOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.UrbanVillage)
		return
	}

	result = input.convertModelToResponseDetail(subDistrictOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input subDistrictService) convertModelToResponseDetail(inputStruct repository.SubDistrictModel) out.SubDistrictDetailResponse {
	return out.SubDistrictDetailResponse{
		ID:           inputStruct.ID.Int64,
		DistrictID:   inputStruct.DistrictID.Int64,
		DistrictName: inputStruct.DistrictName.String,
		Code:         inputStruct.Code.String,
		Name:         inputStruct.Name.String,
		Status:       inputStruct.Status.String,
		CreatedBy:    inputStruct.CreatedBy.Int64,
		CreatedAt:    inputStruct.CreatedAt.Time,
		UpdatedBy:    inputStruct.UpdatedBy.Int64,
		UpdatedAt:    inputStruct.UpdatedAt.Time,
	}
}

func (input subDistrictService) validateViewSubDistrict(inputStruct *in.SubDistrictRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}
