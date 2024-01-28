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
	util2 "nexsoft.co.id/nextrac2/util"
	"time"
)

func (input customerInstallationService) deleteSiteService(tx *sql.Tx, customerModel repository.CustomerInstallationModel, dataAudit *[]repository.AuditSystemModel, indexSite int, contextModel *applicationModel.ContextModel, timeNow time.Time, scope map[string]interface{}) (err errorModel.ErrorModel) {
	var (
		fileName          = "SiteDeleteService.go"
		funcName          = "deleteSiteService"
		custID            = util2.GenerateConstantaI18n(constanta.CustomerID, contextModel.AuthAccessTokenModel.Locale, nil)
		db                = serverconfig.ServerAttribute.DBConnection
		preparingError    string
		isUsedSite        bool
		idInstallation    string
		idInstallationInt []int
		resultOnDB        repository.CustomerSiteModel
	)

	preparingError = fmt.Sprintf(`%s %d`, custID, customerModel.CustomerInstallationData[indexSite].CustomerID.Int64)

	//--- Check Updated At (Must Fill)
	if customerModel.CustomerInstallationData[indexSite].UpdatedAt.Time.IsZero() {
		updatedAt := util2.GenerateConstantaI18n(constanta.UpdatedAt, contextModel.AuthAccessTokenModel.Locale, nil)
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, fmt.Sprintf(`%s %s`, updatedAt, preparingError))
		return
	}

	//--- Check Customer Site Exist
	_, isUsedSite, idInstallation, err = dao.CustomerSiteDAO.CheckCustomerSiteIsExist(db, customerModel, indexSite, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	//--- Is Used ?
	if isUsedSite {
		err = errorModel.GenerateDataUsedError(fileName, funcName, fmt.Sprintf(`Site %s`, preparingError))
		return
	}

	//--- Customer Installation Refactor
	err = input.customerInstallationWhenExist(idInstallation, &idInstallationInt)
	if err.Error != nil {
		return
	}

	//--- Get Customer By Site ID
	resultOnDB, err = dao.CustomerSiteDAO.GetCustomerSiteForUpdate(db, customerModel, indexSite)
	if err.Error != nil {
		return
	}

	//--- Check Site ID
	if resultOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, fmt.Sprintf(`Site %s`, preparingError))
		return
	}

	//--- Lock
	if resultOnDB.UpdatedAt.Time != customerModel.CustomerInstallationData[indexSite].UpdatedAt.Time {
		err = errorModel.GenerateDataLockedError(fileName, funcName, preparingError)
		return
	}

	for _, valueIDInstallation := range idInstallationInt {
		var (
			isUsed          bool
			installationStr = util2.GenerateConstantaI18n(constanta.InstallationID, contextModel.AuthAccessTokenModel.Locale, nil)
		)

		isUsed, err = dao.LicenseConfigDAO.CheckInstallationIsUsed(db, repository.LicenseConfigModel{InstallationID: sql.NullInt64{Int64: int64(valueIDInstallation)}})
		if err.Error != nil {
			return
		}

		if isUsed {
			err = errorModel.GenerateDataUsedError(fileName, funcName, fmt.Sprintf(`%s: %d Site %s`, installationStr, valueIDInstallation, preparingError))
			return
		}

		*dataAudit = append(*dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.CustomerInstallationDAO.TableName, int64(valueIDInstallation), 0)...)
	}

	if len(idInstallationInt) > 0 {
		err = dao.CustomerInstallationDAO.DeleteCustomerInstallationBySiteID(tx, customerModel, indexSite)
		if err.Error != nil {
			return
		}
	}

	*dataAudit = append(*dataAudit, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *contextModel, timeNow, dao.CustomerSiteDAO.TableName, customerModel.CustomerInstallationData[indexSite].SiteID.Int64, 0)...)
	err = dao.CustomerSiteDAO.DeleteCustomerSite(tx, customerModel, indexSite)
	return
}

func (input customerInstallationService) customerInstallationWhenExist(idInstallation string, idInstallationInt *[]int) (err errorModel.ErrorModel) {
	var idInstallationIntTemp []int

	idInstallationIntTemp, err = service.RefactorArrayAggInt(idInstallation)
	if err.Error != nil {
		return
	}

	for _, valueIdInstallationIntTemp := range idInstallationIntTemp {
		*idInstallationInt = append(*idInstallationInt, valueIdInstallationIntTemp)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
