package CustomerInstallationService

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input customerInstallationService) validateSiteService(tx *sql.Tx, customerInstallationModel repository.CustomerInstallationModel, dataAudit *[]repository.AuditSystemModel, indexSite int, contextModel *applicationModel.ContextModel, timeNow time.Time,
	trackUniqueIDMap map[int][]repository.CustomerInstallationTracking, scope map[string]interface{}) (err errorModel.ErrorModel) {

	var (
		fileName          = "SiteNoActionService.go"
		funcName          = "validateSiteService"
		custIns           = "Customer Installation"
		db                = serverconfig.ServerAttribute.DBConnection
		siteID            = input.generateConstanta(constanta.SiteID, contextModel)
		customerID        = input.generateConstanta(constanta.CustomerID, contextModel)
		existPrepareError = fmt.Sprintf(`%s [%s : %d]`, siteID, customerID, customerInstallationModel.CustomerInstallationData[indexSite].CustomerID.Int64)
		isExist           bool
		queue             int64
	)

	isExist, err = dao.CustomerSiteDAO.CheckCustomerSiteOnly(db, customerInstallationModel, indexSite, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if !isExist {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, existPrepareError)
		return
	}

	//--- Data Audit
	*dataAudit = append(*dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.CustomerSiteDAO.TableName, customerInstallationModel.CustomerInstallationData[indexSite].SiteID.Int64, 0)...)

	for i := 0; i < len(customerInstallationModel.CustomerInstallationData[indexSite].Installation); i++ {

		if customerInstallationModel.CustomerInstallationData[indexSite].Installation[i].InstallationID.Int64 == 0 && customerInstallationModel.CustomerInstallationData[indexSite].Installation[i].Action.Int32 == int32(constanta.ActionInsertCode) {

			//--- Installation id equal 0 or empty and action equal 1
			err = input.doInsertInstallationService(tx, customerInstallationModel, indexSite, i, dataAudit, &queue, contextModel, trackUniqueIDMap, scope)
			if err.Error != nil {
				return
			}

		} else if customerInstallationModel.CustomerInstallationData[indexSite].Installation[i].InstallationID.Int64 > 0 && customerInstallationModel.CustomerInstallationData[indexSite].Installation[i].Action.Int32 == int32(constanta.ActionDeleteCode) {

			//--- Installation id more than 0 and action equal 3
			err = input.deleteAndUpdateInstallationService(tx, customerInstallationModel, indexSite, i, dataAudit, contextModel, timeNow, false, trackUniqueIDMap, scope)
			if err.Error != nil {
				return
			}

		} else if customerInstallationModel.CustomerInstallationData[indexSite].Installation[i].InstallationID.Int64 > 0 && customerInstallationModel.CustomerInstallationData[indexSite].Installation[i].Action.Int32 == int32(constanta.ActionUpdateCode) {

			//--- Installation id more than 0 and action equal 2
			err = input.deleteAndUpdateInstallationService(tx, customerInstallationModel, indexSite, i, dataAudit, contextModel, timeNow, true, trackUniqueIDMap, scope)
			if err.Error != nil {
				return
			}

		} else {

			//--- Error if no one action for installation
			err = errorModel.GenerateWrongActionInstallation(fileName, funcName, custIns, constanta.CustomerInstallationRegex, indexSite+1, i+1)
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) checkParentCustomerInstallation(lvl1 repository.CustomerInstallationModel, lvl2 repository.CustomerInstallationData, lvl3 repository.CustomerInstallationDetail,
	idxInst int, parentClientType int64, contextModel *applicationModel.ContextModel, prepErrBasic string) (err errorModel.ErrorModel) {

	var (
		db                = serverconfig.ServerAttribute.DBConnection
		customer          repository.CustomerInstallationForConfig
		custInstallation  repository.CustomerInstallationDetail
		resultParent      map[int64]repository.CustomerInstallationDetail
		productParent     map[int64]bool
		isCheckAllRequest bool
	)

	customer = repository.CustomerInstallationForConfig{
		ParentCustomerID: sql.NullInt64{Int64: lvl1.ParentCustomerID.Int64},
		CustomerID:       sql.NullInt64{Int64: lvl2.CustomerID.Int64},
		SiteID:           sql.NullInt64{Int64: lvl2.SiteID.Int64},
	}

	custInstallation = repository.CustomerInstallationDetail{
		UniqueID1:          sql.NullString{String: lvl3.UniqueID1.String},
		UniqueID2:          sql.NullString{String: lvl3.UniqueID2.String},
		ParentClientTypeID: sql.NullInt64{Int64: parentClientType},
	}

	//--- Get All Parent Customer Installation
	resultParent, err = dao.CustomerInstallationDAO.GetAllParentCustomerInstallation(db, customer, custInstallation)
	if err.Error != nil {
		return
	}

	//--- Get All Product Parent
	productParent, err = dao.ProductDAO.GetProductParentByClientTypeID(db, repository.ProductModel{ClientTypeID: sql.NullInt64{Int64: parentClientType}})
	if err.Error != nil {
		return
	}

	if len(resultParent) < 1 {
		//--- Check All Request
		isCheckAllRequest = true
	} else {
		//--- Request 1 Directly Return Function
		if len(lvl2.Installation) == 1 {
			_, ok := resultParent[lvl3.InstallationID.Int64]
			if !ok {
				return
			}
		}
	}

	return input.checkNewInsertedAndUpdatedParent(lvl2, lvl3, idxInst, parentClientType, isCheckAllRequest, resultParent, productParent, contextModel, prepErrBasic)
}
