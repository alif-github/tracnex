package CustomerInstallationService

import (
	"database/sql"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input customerInstallationService) UpdateCustomerInstallation(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "UpdateCustomerInstallation"
		inputStruct in.CustomerSiteInstallationRequest
	)

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateUpdateCustomerSiteInstallation)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAuditCustom(funcName, inputStruct, contextModel, input.doUpdateCustomerInstallation, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_UPDATE_CUSTOMER_INSTALLATION_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) doUpdateCustomerInstallation(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, isServiceUpdate bool, err errorModel.ErrorModel) {
	var (
		fileName                  = "UpdateCustomerInstallationService.go"
		funcName                  = "doUpdateCustomerInstallation"
		db                        = serverconfig.ServerAttribute.DBConnection
		inputStruct               = inputStructInterface.(in.CustomerSiteInstallationRequest)
		isExist                   bool
		customerInstallationModel repository.CustomerInstallationModel
		scope                     map[string]interface{}
		mappingScopeDB            map[string]applicationModel.MappingScopeDB
	)

	//-- New Mapping Scope DB
	mappingScopeDB = input.newMappingScopeDB()

	scope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	//--- Create Model For Process
	input.createModel(inputStruct, &customerInstallationModel, contextModel, timeNow)

	//--- Check Parent Customer ID
	isExist, err = dao.CustomerDAO.GetCustomerParent(db, repository.CustomerModel{ID: sql.NullInt64{Int64: customerInstallationModel.ParentCustomerID.Int64}}, scope, mappingScopeDB)
	if err.Error != nil {
		return
	}

	if !isExist {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, constanta.ParentCustomerID)
		return
	}

	//--- Core Process
	dataAudit, err = input.siteService(customerInstallationModel, tx, contextModel, timeNow, &isServiceUpdate, scope)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) siteService(customerInstallationModel repository.CustomerInstallationModel, tx *sql.Tx, contextModel *applicationModel.ContextModel, timeNow time.Time, isServiceUpdate *bool, scope map[string]interface{}) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		fileName         = "UpdateCustomerInstallationService.go"
		funcName         = "siteService"
		cusSite          = "Customer Site"
		site             = customerInstallationModel.CustomerInstallationData
		trackUniqueIDMap = make(map[int][]repository.CustomerInstallationTracking)
	)

	for i := 0; i < len(site); i++ {
		if site[i].SiteID.Int64 > 0 && site[i].Action.Int32 == int32(constanta.ActionDeleteCode) {

			//--- Delete if site id more than 0 and action equal 3
			if err = input.deleteSiteService(tx, customerInstallationModel, &dataAudit, i, contextModel, timeNow, scope); err.Error != nil {
				return
			}

			*isServiceUpdate = true
		} else if site[i].SiteID.Int64 == 0 && site[i].Action.Int32 == int32(constanta.ActionInsertCode) {

			//--- Insert if site id empty or equal 0 and action equal 1
			if err = input.insertSiteService(tx, customerInstallationModel, &dataAudit, i, contextModel, trackUniqueIDMap, scope); err.Error != nil {
				return
			}

		} else if site[i].SiteID.Int64 > 0 && site[i].Action.Int32 == int32(constanta.ActionNoActionCode) {

			//--- No one action in site id more than 0 and action equal 4
			if err = input.validateSiteService(tx, customerInstallationModel, &dataAudit, i, contextModel, timeNow, trackUniqueIDMap, scope); err.Error != nil {
				return
			}

			*isServiceUpdate = true
		} else {

			//--- Error no one corresponding site service suitable
			err = errorModel.GenerateWrongAction(fileName, funcName, cusSite, constanta.CustomerSiteRegex, i+1)
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) validateUpdateCustomerSiteInstallation(inputStruct *in.CustomerSiteInstallationRequest, contextModel *applicationModel.ContextModel) errorModel.ErrorModel {
	return inputStruct.ValidateUpdateCustomerInstallation(contextModel)
}

func (input customerInstallationService) checkNewInsertedAndUpdatedParent(lvl2 repository.CustomerInstallationData, lvl3 repository.CustomerInstallationDetail, idxInst int, parentClientType int64,
	isCheckAllRequest bool, parentInstOnDB map[int64]repository.CustomerInstallationDetail, productParent map[int64]bool, contextModel *applicationModel.ContextModel, prepErrBasic string) (err errorModel.ErrorModel) {

	var (
		fileName     = "UpdateCustomerInstallationService.go"
		funcName     = "checkNewInsertedAndUpdatedParent"
		db           = serverconfig.ServerAttribute.DBConnection
		installation = input.generateConstanta(constanta.Installation, contextModel)
		instTemp     = make([]repository.CustomerInstallationDetail, len(lvl2.Installation))
		productModel []repository.ProductModel
		isValid      bool
	)

	copy(instTemp, lvl2.Installation)
	instTemp = append(instTemp[:idxInst], instTemp[idxInst+1:]...)
	if len(instTemp) < 1 {
		err = errorModel.GenerateErrorParentInstallationNotFound(fileName, funcName, fmt.Sprintf(`%s %s`, installation, prepErrBasic))
		return
	}

	for _, itemInst := range instTemp {
		if itemInst.UniqueID1.String == lvl3.UniqueID1.String && itemInst.UniqueID2.String == lvl3.UniqueID2.String {
			if !isCheckAllRequest {
				v, ok := parentInstOnDB[itemInst.InstallationID.Int64]
				if !ok {
					if itemInst.Action.Int32 == constanta.ActionInsertCode || itemInst.InstallationID.Int64 < 1 {
						_, okProd := productParent[itemInst.ProductID.Int64]
						if okProd {
							return
						}
					}
					continue
				}

				if itemInst.Action.Int32 == constanta.ActionDeleteCode {
					delete(parentInstOnDB, itemInst.InstallationID.Int64)
					continue
				}

				if v.ProductID.Int64 == itemInst.ProductID.Int64 {
					return
				}
			}

			productModel = append(productModel, repository.ProductModel{
				ID: sql.NullInt64{Int64: itemInst.ProductID.Int64},
			})
		}
	}

	//--- Product Model Empty Then Error
	if len(productModel) < 1 {
		if len(parentInstOnDB) > 0 {
			return
		}

		err = errorModel.GenerateErrorParentInstallationNotFound(fileName, funcName, fmt.Sprintf(`%s %s`, installation, prepErrBasic))
		return
	}

	//--- Check Valid Parent Product
	isValid, err = dao.ProductDAO.CheckValidParentProductByIDAndClientType(db, productModel, parentClientType)
	if err.Error != nil {
		return
	}

	if !isValid {
		err = errorModel.GenerateErrorParentInstallationNotFound(fileName, funcName, fmt.Sprintf(`%s %s`, installation, prepErrBasic))
		return
	}

	return
}

func (input customerInstallationService) checkProductAndParent(lv3 repository.CustomerInstallationDetail, scope map[string]interface{}, prepError string) (productCode string, isCheckParentInst bool, parentClientType int64, err errorModel.ErrorModel) {
	var (
		fileName     = "UpdateCustomerInstallation.go"
		funcName     = "checkProductAndParent"
		db           = serverconfig.ServerAttribute.DBConnection
		productOnDB  repository.ProductModel
		productModel repository.ProductModel
	)

	productModel = repository.ProductModel{ID: sql.NullInt64{Int64: lv3.ProductID.Int64}}
	productOnDB, err = dao.ProductDAO.CheckProductAndParentClientTypeIsExist(db, productModel, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if productOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, prepError)
		return
	}

	if productOnDB.ParentClientTypeID.Int64 > 0 {
		isCheckParentInst = true
		parentClientType = productOnDB.ParentClientTypeID.Int64
	}

	productCode = productOnDB.ProductID.String
	return
}

func (input customerInstallationService) checkDuplicateUnique12(lv1 repository.CustomerInstallationModel, lv3 repository.CustomerInstallationDetail, idxSite, idxInst int,
	prepErrUniqueKey string, trackUnique map[int][]repository.CustomerInstallationTracking, tempTrackInst map[repository.CustomerInstallationTracking]bool) (err errorModel.ErrorModel) {

	var (
		fileName      = "UpdateCustomerInstallation.go"
		funcName      = "checkDuplicateUnique12"
		db            = serverconfig.ServerAttribute.DBConnection
		isExistUnique bool
		tracking      repository.CustomerInstallationTracking
	)

	//--- Check Different Site
	isExistUnique, err = dao.CustomerInstallationDAO.CheckCustomerInstallationInDifferentSite(db, lv1, idxSite, idxInst)
	if err.Error != nil {
		return
	}

	if isExistUnique {
		err = errorModel.GenerateDataUsedError(fileName, funcName, prepErrUniqueKey)
		return
	}

	//--- Model Mapping
	tracking = repository.CustomerInstallationTracking{
		KeyID:     sql.NullInt64{Int64: int64(idxSite)},
		UniqueID1: sql.NullString{String: lv3.UniqueID1.String},
		UniqueID2: sql.NullString{String: lv3.UniqueID2.String},
	}

	//--- Map tracking
	tempTrackInst[tracking] = true

	//--- Check track global
	err = input.trackingUniqueID(trackUnique, idxSite, lv3.UniqueID1.String, lv3.UniqueID2.String, prepErrUniqueKey)
	if err.Error != nil {
		return
	}

	return
}
