package ProductService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input productService) InsertProduct(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertProduct"
		inputStruct in.ProductRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateInsert)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertProduct, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- additional function
	})

	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_PRODUCT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) doInsertProduct(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		productID    int64
		inputStruct  = inputStructInterface.(in.ProductRequest)
		productModel repository.ProductModel
		scopeLimit   map[string]interface{}
	)

	scopeLimit, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	err = input.validateRelation(inputStruct, contextModel, scopeLimit)
	if err.Error != nil {
		return
	}

	productModel = input.createModelProduct(inputStruct, contextModel, timeNow)
	productID, err = dao.ProductDAO.InsertProduct(tx, productModel)
	if err.Error != nil {
		err = CheckDuplicateError(err)
		return
	}

	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.ProductDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: productID},
	})

	productModel.ID.Int64 = productID
	err = input.doInsertToProductComponent(tx, productModel, &dataAudit)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) validateInsert(inputStruct *in.ProductRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertProduct()
}

func (input productService) doInsertToProductComponent(tx *sql.Tx, productModel repository.ProductModel, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {
	if len(productModel.ProductComponentModel) > 0 {
		var productComponentID []int64
		productComponentID, err = dao.ProductComponentDAO.InsertProductComponent(tx, productModel)
		if err.Error != nil {
			return
		}

		for _, valueProductComponentID := range productComponentID {
			*dataAudit = append(*dataAudit, repository.AuditSystemModel{
				TableName:  sql.NullString{String: dao.ProductComponentDAO.TableName},
				PrimaryKey: sql.NullInt64{Int64: valueProductComponentID},
			})
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
