package ProductService

import (
	"database/sql"
	"errors"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"strconv"
	"time"
)

func (input productService) UpdateProduct(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateProduct"
		inputStruct in.ProductRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateProduct)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAuditCustom(funcName, inputStruct, contextModel, input.doUpdateProduct, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("OK", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_UPDATE_PRODUCT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	return
}

func (input productService) doUpdateProduct(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, isServiceUpdate bool, err errorModel.ErrorModel) {
	var (
		inputStruct                    = inputStructInterface.(in.ProductRequest)
		productOnDB                    []repository.GetForUpdateProduct
		productComponentModelForUpdate []repository.ProductComponentModel
		productComponentModelForDelete []repository.ProductComponentModel
		productComponentModelForInsert []repository.ProductComponentModel
		scopeLimit                     map[string]interface{}
	)

	isServiceUpdate = true
	scopeLimit, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	//--- Create model
	parameterProductModel := input.createModelProduct(inputStruct, contextModel, timeNow)

	//--- Lock product get for update
	productOnDB, err = input.getProductForUpdate(parameterProductModel, inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	//-- Check and validation db relation
	err = input.validateRelation(inputStruct, contextModel, scopeLimit)
	if err.Error != nil {
		return
	}

	//-- Data audit before update
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ProductDAO.TableName, productOnDB[0].ID.Int64, 0)...)

	//-- Update product
	err = dao.ProductDAO.UpdateProduct(tx, parameterProductModel)
	if err.Error != nil {
		err = CheckDuplicateError(err)
		return
	}

	//-- Check request product component ID and product component ID db comparison
	err = input.checkProductComponentCompareByDB(parameterProductModel, productOnDB, contextModel)
	if err.Error != nil {
		return
	}

	//-- Check format request in product component
	err = input.checkFormatRequestProduct(&productComponentModelForUpdate, &productComponentModelForDelete, &productComponentModelForInsert, parameterProductModel)
	if err.Error != nil {
		return
	}

	//-- Append to update model product component
	err = input.updateHelpServiceProductComponent(tx, contextModel, timeNow, parameterProductModel, productComponentModelForUpdate, &dataAudit)
	if err.Error != nil {
		return
	}

	//-- Append to delete model product component
	err = input.deleteHelpServiceProductComponent(tx, contextModel, timeNow, parameterProductModel, productComponentModelForDelete, &dataAudit)
	if err.Error != nil {
		return
	}

	//-- Append to insert model product component
	err = input.insertHelpServiceProductComponent(tx, contextModel, timeNow, parameterProductModel, &dataAudit, productComponentModelForInsert)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) validateUpdateProduct(inputStruct *in.ProductRequest) errorModel.ErrorModel {
	return inputStruct.ValidationUpdateProduct()
}

func (input productService) getProductForUpdate(parameterProductModel repository.ProductModel, inputStruct in.ProductRequest, contextModel *applicationModel.ContextModel) (productOnDB []repository.GetForUpdateProduct, err errorModel.ErrorModel) {
	fileName := "UpdateProductService.go"
	funcName := "getProductForUpdate"

	productOnDB, err = dao.ProductDAO.GetProductForUpdate(serverconfig.ServerAttribute.DBConnection, parameterProductModel)
	if err.Error != nil {
		return
	}

	_, err = dao.ProductDAO.LockProductForUpdate(serverconfig.ServerAttribute.DBConnection, parameterProductModel)
	if err.Error != nil {
		return
	}

	if len(productOnDB) < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Product)
		return
	}

	if productOnDB[0].ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.Product)
		return
	}

	if productOnDB[0].IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(fileName, funcName, constanta.Product)
		return
	}

	err = input.CheckUserLimitedByOwnAccess(contextModel, productOnDB[0].CreatedBy.Int64)
	if err.Error != nil {
		return
	}

	if inputStruct.UpdatedAt != productOnDB[0].UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.UpdatedAt)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) checkProductComponentCompareByDB(parameterProductModel repository.ProductModel, productOnDB []repository.GetForUpdateProduct, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	fileName := "UpdateProductService.go"
	funcName := "checkProductComponentCompareByDB"

	for _, itemProductComponentRequest := range parameterProductModel.ProductComponentModel {
		for idx, itemProductComponentOnDB := range productOnDB {

			//---- If new product component
			if itemProductComponentRequest.ID.Int64 == 0 {
				break
			}

			//---- If product component on request equal product component on db
			if itemProductComponentRequest.ID.Int64 == itemProductComponentOnDB.ProductComponentID.Int64 {
				break
			}

			//---- If this is last lap product on db looping
			if len(productOnDB)-(idx+1) == 0 {
				err = errorModel.GenerateUnknownDataError(fileName, funcName, util.GenerateConstantaI18n(constanta.ProductComponentID, contextModel.AuthAccessTokenModel.Locale, nil)+" "+strconv.Itoa(int(itemProductComponentRequest.ID.Int64)))
				return
			}

		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) updateHelpServiceProductComponent(tx *sql.Tx, contextModel *applicationModel.ContextModel, timeNow time.Time, parameterProductModel repository.ProductModel,
	productComponentModel []repository.ProductComponentModel, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {

	for i := 0; i < len(productComponentModel) && len(productComponentModel) > 0; i++ {
		//-------------------- Input Audit System
		*dataAudit = append(*dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ProductComponentDAO.TableName, productComponentModel[i].ID.Int64, 0)...)

		//-------------------- Update Product Component
		err = dao.ProductComponentDAO.UpdateProductComponent(tx, productComponentModel[i], parameterProductModel.ID.Int64)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) deleteHelpServiceProductComponent(tx *sql.Tx, contextModel *applicationModel.ContextModel, timeNow time.Time, parameterProductModel repository.ProductModel,
	productComponentModel []repository.ProductComponentModel, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {

	for i := 0; i < len(productComponentModel) && len(productComponentModel) > 0; i++ {
		//-------------------- Input Audit System
		*dataAudit = append(*dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.ProductComponentDAO.TableName, productComponentModel[i].ID.Int64, 0)...)
	}

	if len(productComponentModel) > 0 {
		//-------------------- Multi Delete Product Component
		err = dao.ProductComponentDAO.DeleteProductComponentMultiple(tx, productComponentModel, parameterProductModel)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) checkFormatRequestProduct(productComponentModelForUpdate *[]repository.ProductComponentModel, productComponentModelForDelete *[]repository.ProductComponentModel,
	productComponentModelForInsert *[]repository.ProductComponentModel, parameterProductModel repository.ProductModel) (err errorModel.ErrorModel) {

	fileName := "UpdateProductService.go"
	funcName := "checkFormatRequestProduct"
	prepareError := "product component no. "

	for idx, itemProductComponent := range parameterProductModel.ProductComponentModel {
		if itemProductComponent.ID.Int64 > 0 && !itemProductComponent.Deleted.Bool {
			*productComponentModelForUpdate = append(*productComponentModelForUpdate, itemProductComponent)
		} else if itemProductComponent.ID.Int64 > 0 && itemProductComponent.Deleted.Bool {
			*productComponentModelForDelete = append(*productComponentModelForDelete, itemProductComponent)
		} else if itemProductComponent.ID.Int64 == 0 && !itemProductComponent.Deleted.Bool {
			*productComponentModelForInsert = append(*productComponentModelForInsert, itemProductComponent)
		} else {
			err = errorModel.GenerateInvalidRequestError(fileName, funcName, errors.New(prepareError+strconv.Itoa(idx)))
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) insertHelpServiceProductComponent(tx *sql.Tx, contextModel *applicationModel.ContextModel, timeNow time.Time, parameterProductModel repository.ProductModel,
	dataAudit *[]repository.AuditSystemModel, productComponentModelForInsert []repository.ProductComponentModel) (err errorModel.ErrorModel) {

	parameterProductModel.ProductComponentModel = productComponentModelForInsert
	parameterProductModel.CreatedBy.Int64 = contextModel.AuthAccessTokenModel.ResourceUserID
	parameterProductModel.CreatedClient.String = contextModel.AuthAccessTokenModel.ClientID
	parameterProductModel.CreatedAt.Time = timeNow

	if len(parameterProductModel.ProductComponentModel) > 0 {
		var resultIdInsertComponent []int64
		resultIdInsertComponent, err = dao.ProductComponentDAO.InsertProductComponent(tx, parameterProductModel)
		if err.Error != nil {
			return
		}

		for _, valueResultId := range resultIdInsertComponent {
			*dataAudit = append(*dataAudit, repository.AuditSystemModel{
				TableName:  sql.NullString{String: dao.ProductComponentDAO.TableName},
				PrimaryKey: sql.NullInt64{Int64: valueResultId},
				Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
			})
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
