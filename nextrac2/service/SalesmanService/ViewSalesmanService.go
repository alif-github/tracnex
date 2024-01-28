package SalesmanService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
)

func (input salesmanService) ViewSalesman(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.SalesmanRequest
	inputStruct, err = input.readBodyAndValidateForView(request, input.validateView)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doViewSalesman(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_VIEW_SALESMAN_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanService) doViewSalesman(inputStruct in.SalesmanRequest, contextModel *applicationModel.ContextModel) (output out.ViewSalesman, err errorModel.ErrorModel) {
	var (
		fileName           = "ViewSalesmanService.go"
		funcName           = "doViewSalesman"
		db                 = serverconfig.ServerAttribute.DBConnection
		resultSalesmanView repository.ViewSalesmanModel
		salesmanModel      repository.SalesmanModel
	)

	salesmanModel = repository.SalesmanModel{
		ID:        sql.NullInt64{Int64: inputStruct.ID},
		CreatedBy: sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	}

	resultSalesmanView, err = dao.SalesmanDAO.ViewSalesman(db, salesmanModel)
	if err.Error != nil {
		return
	}

	if resultSalesmanView.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, "Salesman")
		return
	}

	output = input.convertSalesmanModelToDTOOut(resultSalesmanView)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input salesmanService) convertSalesmanModelToDTOOut(salesmanModel repository.ViewSalesmanModel) out.ViewSalesman {
	return out.ViewSalesman{
		ID:            salesmanModel.ID.Int64,
		PersonTitleID: salesmanModel.PersonTitleID.Int64,
		PersonTitle:   salesmanModel.PersonTitle.String,
		Sex:           salesmanModel.Sex.String,
		Nik:           salesmanModel.Nik.String,
		Address:       salesmanModel.Address.String,
		Hamlet:        salesmanModel.Hamlet.String,
		Neighbourhood: salesmanModel.Neighbourhood.String,
		FirstName:     salesmanModel.FirstName.String,
		LastName:      salesmanModel.LastName.String,
		ProvinceID:    salesmanModel.MdbProvinceID.Int64,
		Province:      salesmanModel.Province.String,
		DistrictID:    salesmanModel.MdbDistrictID.Int64,
		District:      salesmanModel.District.String,
		Phone:         salesmanModel.Phone.String,
		Email:         salesmanModel.Email.String,
		Status:        salesmanModel.Status.String,
		CreatedAt:     salesmanModel.CreatedAt.Time,
		UpdatedAt:     salesmanModel.UpdatedAt.Time,
		UpdatedBy:     salesmanModel.UpdatedBy.Int64,
		UpdatedName:   salesmanModel.UpdatedName.String,
	}
}

func (input salesmanService) validateView(inputStruct *in.SalesmanRequest) errorModel.ErrorModel {
	return inputStruct.ValidateViewSalesman()
}
