package CustomerCategoryService

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
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"time"
)

func (input customerCategoryService) InsertCustomerCategory(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "InsertCustomerCategory"
	var inputStruct in.CustomerCategoryRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return 
	}
	
	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertCustomerCategory, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input customerCategoryService) doInsertCustomerCategory(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	inputStruct := inputStructInterface.(in.CustomerCategoryRequest)

	customerCategoryModel := input.convertStructToModelInsert(inputStruct, contextModel.AuthAccessTokenModel, timeNow)

	idCustomerCatg, err := dao.CustomerCategoryDAO.InsertCustomerCategory(tx, customerCategoryModel)
	if err.Error != nil {
		err = input.CheckDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.CustomerCategoryDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idCustomerCatg},
	})

	var dataAuditTemp repository.AuditSystemModel
	dataAuditTemp, err =input.GenerateDataScope(tx, idCustomerCatg, dao.CustomerCategoryDAO.TableName, constanta.CustomerCategoryDataScope, contextModel.AuthAccessTokenModel.ResourceUserID, contextModel.AuthAccessTokenModel.ClientID, timeNow)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, dataAuditTemp)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerCategoryService) convertStructToModelInsert(inputStruct in.CustomerCategoryRequest, authAccessModel model.AuthAccessTokenModel, timeNow time.Time) (output repository.CustomerCategoryModel) {
	return repository.CustomerCategoryModel{
		CustomerCategoryID:   sql.NullString{String: inputStruct.CustomerCategoryID},
		CustomerCategoryName: sql.NullString{String: inputStruct.CustomerCategoryName},
		CreatedBy:            sql.NullInt64{Int64: authAccessModel.ResourceUserID},
		CreatedAt:            sql.NullTime{Time: timeNow},
		CreatedClient:        sql.NullString{String: authAccessModel.ClientID},
		UpdatedBy:            sql.NullInt64{Int64: authAccessModel.ResourceUserID},
		UpdatedAt:            sql.NullTime{Time: timeNow},
		UpdatedClient:        sql.NullString{String: authAccessModel.ClientID},
	}
}

func (input customerCategoryService) validateInsert(inputStruct *in.CustomerCategoryRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsert()
}