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
	"time"
)

func (input customerInstallationService) deleteAndUpdateInstallationService(tx *sql.Tx, customerInstallationModel repository.CustomerInstallationModel, idxSite, idxInst int, dataAudit *[]repository.AuditSystemModel, contextModel *applicationModel.ContextModel, timeNow time.Time, isUpdate bool,
	trackUnique map[int][]repository.CustomerInstallationTracking, scope map[string]interface{}) (err errorModel.ErrorModel) {
	var (
		fileName             = "InstallationUpdateDeleteService.go"
		funcName             = "deleteAndUpdateInstallationService"
		installationID       = input.generateConstanta(constanta.InstallationID, contextModel)
		customerID           = input.generateConstanta(constanta.CustomerID, contextModel)
		installation         = input.generateConstanta(constanta.Installation, contextModel)
		updatedAt            = input.generateConstanta(constanta.UpdatedAt, contextModel)
		productID            = input.generateConstanta(constanta.Product, contextModel)
		uniqueID1            = input.generateConstanta(constanta.UniqueID1, contextModel)
		uniqueID2            = input.generateConstanta(constanta.UniqueID2, contextModel)
		tempTrackInst        = make(map[repository.CustomerInstallationTracking]bool)
		db                   = serverconfig.ServerAttribute.DBConnection
		lvl2                 = customerInstallationModel.CustomerInstallationData[idxSite]
		lvl3                 = customerInstallationModel.CustomerInstallationData[idxSite].Installation[idxInst]
		prepErrBasic         string
		prepErrProduct       string
		prepErrUniqueKey     string
		prepErrParentNoExist string
		//productCode          string
		parentClientType  int64
		isCheckParentInst bool
		resultDB          repository.CustomerInstallationDetail
	)

	//--- Prepare Error
	prepErrBasic = fmt.Sprintf(`[%s: %d Site %s: %d]`, installationID, lvl3.InstallationID.Int64, customerID, lvl2.CustomerID.Int64)
	prepErrProduct = fmt.Sprintf(`%s %s`, productID, prepErrBasic)
	//prepErrCodeUnique := fmt.Sprintf(`%s %s`, uniqueID1, prepErrBasic)
	prepErrUniqueKey = fmt.Sprintf(`%s & %s %s`, uniqueID1, uniqueID2, prepErrBasic)
	prepErrParentNoExist = fmt.Sprintf(`[%s: %d, Site %s: %d, %s: %d, Unique ID 1: %s, Unique ID 2: %s]`, installationID, lvl3.InstallationID.Int64, customerID,
		lvl2.CustomerID.Int64, productID, lvl3.ProductID.Int64, lvl3.UniqueID1.String, lvl3.UniqueID2.String)

	//--- Check Updated At (Must Fill)
	if lvl3.UpdatedAt.Time.IsZero() {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, fmt.Sprintf(`%s %s`, updatedAt, prepErrBasic))
		return
	}

	//--- Check Customer Installation
	resultDB, err = dao.CustomerInstallationDAO.GetCustomerInstallationForUpdate(db, customerInstallationModel, idxSite, idxInst)
	if resultDB.InstallationID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, fmt.Sprintf(`%s %s`, installation, prepErrBasic))
		return
	}

	if resultDB.IsUsed.Bool {
		err = errorModel.GenerateDataUsedError(fileName, funcName, fmt.Sprintf(`%s %s`, installation, prepErrBasic))
		return
	}

	if resultDB.UpdatedAt.Time != lvl3.UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(fileName, funcName, fmt.Sprintf(`%s %s`, updatedAt, prepErrBasic))
		return
	}

	//--- Check Product
	_, isCheckParentInst, parentClientType, err = input.checkProductAndParent(lvl3, scope, prepErrProduct)
	if err.Error != nil {
		return
	}

	if isUpdate {
		//lvl3.UniqueID1.String = strings.ToUpper(lvl3.UniqueID1.String)
		//uniqueID1Code := lvl3.UniqueID1.String[0:2]
		//productCodeSplit := productCode[0:2]
		//if uniqueID1Code != productCodeSplit {
		//	err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.InstallationUniqueRegex, prepErrCodeUnique, "")
		//	return
		//}

		//--- Update Customer Installation
		err = input.doUpdateInstallationService(tx, customerInstallationModel, lvl2, lvl3, idxSite, idxInst, prepErrUniqueKey, prepErrParentNoExist, resultDB, trackUnique, tempTrackInst, isCheckParentInst, parentClientType, contextModel, timeNow, dataAudit)
		if err.Error != nil {
			return
		}
	} else {
		//--- Delete Customer Installation
		err = input.doDeleteCustomerInstallation(tx, customerInstallationModel, lvl2, lvl3, idxSite, idxInst, prepErrParentNoExist, resultDB, contextModel, timeNow, dataAudit)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) doGetProductChildAndInstallationChild(db *sql.DB, parentInfo repository.ProductModel, resultDB repository.CustomerInstallationDetail, lvl2 repository.CustomerInstallationData,
	lvl3 repository.CustomerInstallationDetail) (childIDProduct map[int64]bool, childInstallation map[int64]repository.CustomerInstallationDetail, parentInstModel repository.CustomerInstallationForConfig, err errorModel.ErrorModel) {

	//--- Get All Product ID Child
	childIDProduct, err = dao.ProductDAO.GetProductChild(db, repository.ProductModel{ClientTypeID: sql.NullInt64{Int64: parentInfo.ParentClientTypeID.Int64}})
	if err.Error != nil {
		return
	}

	//--- Check Child
	parentInstModel = repository.CustomerInstallationForConfig{
		ID:        sql.NullInt64{Int64: lvl3.InstallationID.Int64},
		UniqueID1: sql.NullString{String: resultDB.UniqueID1.String},
		UniqueID2: sql.NullString{String: resultDB.UniqueID2.String},
		SiteID:    sql.NullInt64{Int64: lvl2.SiteID.Int64},
		ProductID: sql.NullInt64{Int64: resultDB.ProductID.Int64},
	}

	childInstallation, err = dao.CustomerInstallationDAO.GetInstallationChild(db, parentInstModel)
	if err.Error != nil {
		return
	}

	return
}
