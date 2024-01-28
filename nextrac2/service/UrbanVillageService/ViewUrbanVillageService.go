package UrbanVillageService

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

func (input urbanVillageService) ViewUrbanVillage(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.UrbanVillageRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewUrbanVillage)
	if err.Error != nil {
		return
	}

	output.Data.Content, err =input.doViewUrbanVillage(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status =input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input urbanVillageService) doViewUrbanVillage(inputStruct in.UrbanVillageRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	funcName := "doViewUrbanVillage"

	urbanVillageOnDB, err := dao.UrbanVillageDAO.ViewUrbanVillage(serverconfig.ServerAttribute.DBConnection, repository.UrbanVillageModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	})
	if err.Error != nil {
		return
	}

	if urbanVillageOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.UrbanVillage)
		return
	}

	result = input.convertModelToResponseDetail(urbanVillageOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input urbanVillageService) convertModelToResponseDetail(inputModel repository.UrbanVillageModel) out.UrbanVillageDetailResponse {
	return out.UrbanVillageDetailResponse{
		ID:              inputModel.ID.Int64,
		SubDistrictID:   inputModel.SubDistrictID.Int64,
		SubDistrictName: inputModel.SubDistrictName.String,
		Code:            inputModel.Code.String,
		Name:            inputModel.Name.String,
		Status:          inputModel.Status.String,
		CreatedBy:       inputModel.CreatedBy.Int64,
		CreatedAt:       inputModel.CreatedAt.Time,
		UpdatedBy:       inputModel.UpdatedBy.Int64,
		UpdatedAt:       inputModel.UpdatedAt.Time,
	}
}

func (input urbanVillageService) validateViewUrbanVillage(inputStruct *in.UrbanVillageRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}
