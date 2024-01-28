package PostalCodeService

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

func (input postalCodeService) ViewPostalCode(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.PostalCodeRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateViewPostalCode)
	if err.Error != nil {
		return
	}

	output.Data.Content, err =input.doViewPostalCode(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status =input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	return
}

func (input postalCodeService) doViewPostalCode(inputStruct in.PostalCodeRequest, contextModel *applicationModel.ContextModel) (result interface{}, err errorModel.ErrorModel) {
	funcName := "doViewUrbanVillage"

	urbanVillageOnDB, err := dao.PostalCodeDAO.ViewPostalCode(serverconfig.ServerAttribute.DBConnection, repository.PostalCodeModel{
		ID: sql.NullInt64{Int64: inputStruct.ID},
	})
	if err.Error != nil {
		return
	}

	if urbanVillageOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.PostalCode)
		return
	}

	result = input.convertModelToResponseDetail(urbanVillageOnDB)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input postalCodeService) convertModelToResponseDetail(inputModel repository.PostalCodeModel) out.PostalCodeDetailResponse {
	return out.PostalCodeDetailResponse{
		ID:               inputModel.ID.Int64,
		UrbanVillageID:   inputModel.UrbanVillageID.Int64,
		UrbanVillageName: inputModel.UrbanVillageName.String,
		Code:             inputModel.Code.String,
		Status:           inputModel.Status.String,
		CreatedBy:        inputModel.CreatedBy.Int64,
		CreatedAt:        inputModel.CreatedAt.Time,
		UpdatedBy:        inputModel.UpdatedBy.Int64,
		UpdatedAt:        inputModel.UpdatedAt.Time,
	}
}

func (input postalCodeService) validateViewPostalCode(inputStruct *in.PostalCodeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateView()
}