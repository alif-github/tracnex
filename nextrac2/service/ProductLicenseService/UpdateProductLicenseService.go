package ProductLicenseService

import (
	"database/sql"
	"github.com/Azure/go-autorest/autorest/date"
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
	"time"
)

func (input productLicenseService) UpdateProductLicense(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateProductLicense"
		inputStruct in.ProductLicenseRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateProductLicense)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doUpdateProductLicense, func(interface{}, applicationModel.ContextModel) {
		//--- Function Additional
	})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_UPDATE_MESSAGE", contextModel)
	return
}

func (input productLicenseService) validateUpdateProductLicense(inputStruct *in.ProductLicenseRequest) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateProductLicense()
}

func (input productLicenseService) doUpdateProductLicense(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName           = "doUpdateProductLicense"
		fileName           = "UpdateProductLicenseService.go"
		inputStruct        = inputStructInterface.(in.ProductLicenseRequest)
		db                 = serverconfig.ServerAttribute.DBConnection
		productLicenseOnDB repository.ProductLicenseModel
		scope              map[string]interface{}
	)

	productLicenseModel := repository.ProductLicenseModel{
		ID:                     sql.NullInt64{Int64: inputStruct.ID},
		LicenseStatus:          sql.NullInt32{Int32: inputStruct.LicenseStatus},
		TerminationDescription: sql.NullString{String: inputStruct.TerminationDescription},
		UpdatedBy:              sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient:          sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:              sql.NullTime{Time: timeNow},
		CreatedBy:              sql.NullInt64{Int64: contextModel.LimitedByCreatedBy},
	}

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	productLicenseOnDB, err = dao.ProductLicenseDAO.GetProductLicenseForUpdate(db, productLicenseModel, scope, input.MappingScopeDB, input.ListScope)
	if err.Error != nil {
		return
	}

	if productLicenseOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ProductLicenseID)
		return
	}

	if inputStruct.LicenseStatus == constanta.ProductLicenseStatusActive {
		err = input.validateReActiveProduct(productLicenseOnDB, timeNow)
		if err.Error != nil {
			return
		}
	}

	//--- Check For Own
	if err = input.CheckUserLimitedByOwnAccess(contextModel, productLicenseOnDB.CreatedBy.Int64); err.Error != nil {
		return
	}

	if productLicenseOnDB.UpdatedAt.Time.Sub(inputStruct.UpdatedAt) != time.Duration(0) {
		err = errorModel.GenerateDataLockedError(fileName, funcName, constanta.UpdatedAt)
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.ProductLicenseDAO.TableName, productLicenseOnDB.ID.Int64, 0)...)
	err = dao.ProductLicenseDAO.UpdateProductLicense(tx, productLicenseModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input productLicenseService) validateReActiveProduct(productOnDB repository.ProductLicenseModel, timeNow time.Time) (err errorModel.ErrorModel) {
	comparedTime := date.Date{Time: timeNow}
	expiredProduct := date.Date{Time: productOnDB.ProductValidThru.Time}

	if expiredProduct.Before(comparedTime.Time) {
		err = errorModel.GenerateRequestError(input.FileName, "validateReActiveProduct", "EXPIRED_LICENSE")
		return
	}

	return
}
