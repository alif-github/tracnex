package ProductService

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
	"nexsoft.co.id/nextrac2/service"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input productService) DeleteProduct(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "DeleteProduct"
		inputStruct in.ProductRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateDelete)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doDeleteProduct, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_DELETE_PRODUCT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) doDeleteProduct(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName     = "DeleteProductService.go"
		funcName     = "doDeleteProduct"
		inputStruct  = inputStructInterface.(in.ProductRequest)
		productModel repository.ProductModel
		productOnDB  []repository.GetForUpdateProduct
	)

	productModel = repository.ProductModel{
		ID:            sql.NullInt64{Int64: inputStruct.ID},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	productOnDB, err = dao.ProductDAO.GetProductForUpdate(serverconfig.ServerAttribute.DBConnection, productModel)
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

	if productOnDB[0].UpdatedAt.Time != inputStruct.UpdatedAt {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.Product)
		return
	}

	// ----------- Write Audit Product
	dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.ProductDAO.TableName, productOnDB[0].ID.Int64, 0)...)

	// ----------- Write Audit Product Component
	for _, valueProductOnDB := range productOnDB {
		dataAudit = append(dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.ProductComponentDAO.TableName, valueProductOnDB.ProductComponentID.Int64, 0)...)
	}

	// ----------- Update for delete
	encodedStr, errorS := service.RandToken(8)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	productModel.ProductName.String = productOnDB[0].ProductName.String + encodedStr
	productModel.ProductID.String = productOnDB[0].ProductID.String + encodedStr
	err = dao.ProductDAO.DeleteProduct(tx, productModel)
	if err.Error != nil {
		return
	}

	err = dao.ProductComponentDAO.DeleteProductComponent(tx, productModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productService) validateDelete(inputStruct *in.ProductRequest) errorModel.ErrorModel {
	return inputStruct.ValidationDeleteProduct()
}
