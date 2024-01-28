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
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
)

func (input customerInstallationService) insertSiteService(tx *sql.Tx, customerInstallationModel repository.CustomerInstallationModel, dataAudit *[]repository.AuditSystemModel, i int, contextModel *applicationModel.ContextModel,
	trackUniqueIDMap map[int][]repository.CustomerInstallationTracking, scope map[string]interface{}) (err errorModel.ErrorModel) {

	var (
		idSiteDB           int64
		idInstallationColl []int64
		db                 = serverconfig.ServerAttribute.DBConnection
	)

	//--- Validate Site
	err = input.validateForInsertSiteService(customerInstallationModel, i, contextModel, scope)
	if err.Error != nil {
		return
	}

	//--- Validate Installation
	err = input.validateInsertInstallation(customerInstallationModel, i, contextModel, trackUniqueIDMap, scope)
	if err.Error != nil {
		return
	}

	//--- Insert Customer Site
	idSiteDB, err = dao.CustomerSiteDAO.InsertCustomerSite(tx, customerInstallationModel, i)
	if err.Error != nil {
		return
	}

	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.CustomerSiteDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idSiteDB},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	//--- Set Site ID
	customerInstallationModel.CustomerInstallationData[i].SiteID.Int64 = idSiteDB
	if len(customerInstallationModel.CustomerInstallationData[i].Installation) > 0 {

		//--- Get Client Mapping ID If Exist
		for idxInstallation, itemInstallation := range customerInstallationModel.CustomerInstallationData[i].Installation {
			var (
				isClientExist   bool
				clientMappingID int64
			)
			//--- Get Client Mapping ID If Exist On DB
			isClientExist, clientMappingID, err = dao.CustomerInstallationDAO.GetClientMappingIDCustomerInstallationByUniqueID(db, itemInstallation)
			if err.Error != nil {
				return
			}
			if isClientExist {
				customerInstallationModel.CustomerInstallationData[i].Installation[idxInstallation].ClientMappingID.Int64 = clientMappingID
			}
		}
		
		//--- Insert Customer Installation Multiple
		idInstallationColl, err = dao.CustomerInstallationDAO.InsertMultiCustomerInstallation(tx, customerInstallationModel, i)
		if err.Error != nil {
			return
		}

		for _, valueIdInstallationColl := range idInstallationColl {
			*dataAudit = append(*dataAudit, repository.AuditSystemModel{
				TableName:  sql.NullString{String: dao.CustomerInstallationDAO.TableName},
				PrimaryKey: sql.NullInt64{Int64: valueIdInstallationColl},
				Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
			})
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) validateForInsertSiteService(customerInstallationModel repository.CustomerInstallationModel, i int, contextModel *applicationModel.ContextModel, scope map[string]interface{}) (err errorModel.ErrorModel) {
	var (
		fileName          = "SiteInsertService.go"
		funcName          = "validateForInsertSiteService"
		prepareError      = util2.GenerateConstantaI18n(constanta.CustomerID, contextModel.AuthAccessTokenModel.Locale, nil) + " " + strconv.Itoa(int(customerInstallationModel.CustomerInstallationData[i].CustomerID.Int64))
		isExistCustomerID bool
	)

	custModel := repository.CustomerModel{ID: sql.NullInt64{Int64: customerInstallationModel.CustomerInstallationData[i].CustomerID.Int64}}
	isExistCustomerID, err = dao.CustomerDAO.CheckCustomerIsExist(serverconfig.ServerAttribute.DBConnection, custModel, scope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	if !isExistCustomerID {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, prepareError)
		return
	}

	isExistCustomerID = false
	custInstModel := repository.CustomerInstallationData{CustomerID: sql.NullInt64{Int64: customerInstallationModel.CustomerInstallationData[i].CustomerID.Int64}}
	isExistCustomerID, err = dao.CustomerSiteDAO.CheckCustomerSiteIsExistByParentIDAndCustomerID(serverconfig.ServerAttribute.DBConnection, custInstModel, customerInstallationModel.ParentCustomerID.Int64)
	if err.Error != nil {
		return
	}

	if isExistCustomerID {
		err = errorModel.GenerateDuplicateErrorWithParam(fileName, funcName, prepareError)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) validateInsertInstallation(customerInstallation repository.CustomerInstallationModel, idxSite int, contextModel *applicationModel.ContextModel,
	trackUnique map[int][]repository.CustomerInstallationTracking, scope map[string]interface{}) (err errorModel.ErrorModel) {

	var (
		fileName      = "CustomerInstallationService.go"
		funcName      = "validateInsertInstallation"
		itemSite      = customerInstallation.CustomerInstallationData[idxSite]
		tempTrackInst = make(map[repository.CustomerInstallationTracking]bool)
		productID     = input.generateConstanta(constanta.ProductID, contextModel)
		customerID    = input.generateConstanta(constanta.CustomerID, contextModel)
		uniqueID1     = input.generateConstanta(constanta.UniqueID1, contextModel)
		uniqueID2     = input.generateConstanta(constanta.UniqueID2, contextModel)
		tempRepo      []repository.CustomerInstallationTracking
	)

	for idxInst, itemInst := range itemSite.Installation {
		var (
			isCheckParentInst    bool
			parentClientType     int64
			prepErrProduct       string
			prepErrUniqueKey     string
			prepErrParentNoExist string
			//productCode          string
		)

		//--- Prepare Error
		prepErrProduct = fmt.Sprintf(`%s: %d [Site %s: %d]`, productID, itemInst.ProductID.Int64, customerID, itemSite.CustomerID.Int64)
		prepErrUniqueKey = fmt.Sprintf(`%s & %s [Site %s: %d]`, uniqueID1, uniqueID2, customerID, itemSite.CustomerID.Int64)
		prepErrParentNoExist = fmt.Sprintf(`[Site %s: %d, %s: %d, Unique ID 1: %s, Unique ID 2: %s]`, customerID, itemSite.CustomerID.Int64, productID, itemInst.ProductID.Int64, itemInst.UniqueID1.String, itemInst.UniqueID2.String)

		//--- Installation ID Must Less Than 1
		if itemInst.InstallationID.Int64 > 0 {
			err = errorModel.GenerateFieldMustEmptyInstallation(fileName, funcName, constanta.InstallationID)
			return
		}

		//--- Updated At Must Empty
		if !itemInst.UpdatedAt.Time.IsZero() {
			err = errorModel.GenerateFieldMustEmptyInstallation(fileName, funcName, constanta.UpdatedAt)
			return
		}

		//--- Action Must 1 (Insert)
		if itemInst.Action.Int32 != constanta.ActionInsertCode {
			err = errorModel.GenerateActionMustInsertInstallation(fileName, funcName, constanta.ActionCode)
			return
		}

		//--- Check Product And Parent Client Type
		_, isCheckParentInst, parentClientType, err = input.checkProductAndParent(itemInst, scope, prepErrProduct)
		if err.Error != nil {
			return
		}

		//itemInst.UniqueID1.String = strings.ToUpper(itemInst.UniqueID1.String)
		//uniqueID1Code := itemInst.UniqueID1.String[0:2]
		//productCodeSplit := productCode[0:2]
		//if uniqueID1Code != productCodeSplit {
		//	err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.InstallationUniqueRegex, constanta.UniqueID1, "")
		//	return
		//}

		//--- Check Duplication Unique Key
		err = input.checkDuplicateUnique12(customerInstallation, itemInst, idxSite, idxInst, prepErrUniqueKey, trackUnique, tempTrackInst)
		if err.Error != nil {
			return
		}

		if isCheckParentInst {
			err = input.checkNewInsertedAndUpdatedParent(itemSite, itemInst, idxInst, parentClientType, true, nil, nil, contextModel, prepErrParentNoExist)
			if err.Error != nil {
				return
			}
		}
	}

	//--- Record Unique
	tempRepo = input.prepareTrackingUniqueID(tempTrackInst)
	trackUnique[idxSite] = tempRepo

	//--- Finish
	err = errorModel.GenerateNonErrorModel()
	return
}
