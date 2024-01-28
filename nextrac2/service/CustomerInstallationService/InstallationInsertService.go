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
)

func (input customerInstallationService) insertInstallationService(tx *sql.Tx, customerInstallationModel repository.CustomerInstallationModel, idxSite int, dataAudit *[]repository.AuditSystemModel,
	queue *int64, contextModel *applicationModel.ContextModel, trackUnique map[int][]repository.CustomerInstallationTracking, scope map[string]interface{}) (err errorModel.ErrorModel) {

	for idxInst := range customerInstallationModel.CustomerInstallationData[idxSite].Installation {
		err = input.doInsertInstallationService(tx, customerInstallationModel, idxSite, idxInst, dataAudit, queue, contextModel, trackUnique, scope)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) doInsertInstallationService(tx *sql.Tx, headerModel repository.CustomerInstallationModel, idxSite int, idxInst int, dataAudit *[]repository.AuditSystemModel, queue *int64, contextModel *applicationModel.ContextModel,
	trackUnique map[int][]repository.CustomerInstallationTracking, scope map[string]interface{}) (err errorModel.ErrorModel) {

	var (
		installationID       = input.generateConstanta(constanta.NewInstallationID, contextModel)
		customer             = input.generateConstanta(constanta.CustomerID, contextModel)
		product              = input.generateConstanta(constanta.Product, contextModel)
		productID            = input.generateConstanta(constanta.ProductID, contextModel)
		uniqueID1            = input.generateConstanta(constanta.UniqueID1, contextModel)
		uniqueID2            = input.generateConstanta(constanta.UniqueID2, contextModel)
		prepErrBasic         string
		prepErrProduct       string
		prepErrUniqueKey     string
		prepErrParentNoExist string
		//productCode          string
		idReturnInstallation int64
		parentClientType     int64
		isCheckParentInst    bool
		tempTrackInst        = make(map[repository.CustomerInstallationTracking]bool)
		tempRepo             []repository.CustomerInstallationTracking
		lv1                  = headerModel
		lv2                  = headerModel.CustomerInstallationData[idxSite]
		lv3                  = headerModel.CustomerInstallationData[idxSite].Installation[idxInst]
		//fileName             = "InstallationInsertService.go"
		//funcName             = "doInsertInstallationService"
		db                   = serverconfig.ServerAttribute.DBConnection
		isClientExist        bool
		clientMappingID      int64
	)

	//--- Prepare Error Basic
	prepErrBasic = fmt.Sprintf(`[%s Site %s: %d]`, installationID, customer, lv2.CustomerID.Int64)
	prepErrProduct = fmt.Sprintf(`%s %s`, product, prepErrBasic)
	//prepErrCodeUnique := fmt.Sprintf(`%s %s`, uniqueID1, prepErrBasic)
	prepErrUniqueKey = fmt.Sprintf(`%s & %s %s`, uniqueID1, uniqueID2, prepErrBasic)
	prepErrParentNoExist = fmt.Sprintf(`[%s Site %s: %d, %s: %d, Unique ID 1: %s, Unique ID 2: %s]`,
		installationID, customer, lv2.CustomerID.Int64, productID, lv3.ProductID.Int64, lv3.UniqueID1.String, lv3.UniqueID2.String)

	//--- Check Product And Parent Client Type
	_, isCheckParentInst, parentClientType, err = input.checkProductAndParent(lv3, scope, prepErrProduct)
	if err.Error != nil {
		return
	}

	//lv3.UniqueID1.String = strings.ToUpper(lv3.UniqueID1.String)
	//uniqueID1Code := lv3.UniqueID1.String[0:2]
	//productCodeSplit := productCode[0:2]
	//if uniqueID1Code != productCodeSplit {
	//	err = errorModel.GenerateFieldFormatWithRuleError(fileName, funcName, constanta.InstallationUniqueRegex, prepErrCodeUnique, "")
	//	return
	//}

	//--- Check Duplication Unique Key
	err = input.checkDuplicateUnique12(lv1, lv3, idxSite, idxInst, prepErrUniqueKey, trackUnique, tempTrackInst)
	if err.Error != nil {
		return
	}

	//--- Check Parent Customer Installation
	if isCheckParentInst {
		err = input.checkParentCustomerInstallation(lv1, lv2, lv3, idxInst, parentClientType, contextModel, prepErrParentNoExist)
		if err.Error != nil {
			return
		}
	}

	//--- Installation Number
	err = input.getInstallationNumber(queue, lv1, idxSite)
	if err.Error != nil {
		return
	}
	
	//--- Get Client Mapping ID If Exist On DB
	isClientExist, clientMappingID, err = dao.CustomerInstallationDAO.GetClientMappingIDCustomerInstallationByUniqueID(db, lv1.CustomerInstallationData[idxSite].Installation[idxInst])
	if err.Error != nil {
		return
	}
	if isClientExist {
		lv1.CustomerInstallationData[idxSite].Installation[idxInst].ClientMappingID.Int64 = clientMappingID
	}

	//--- Insert Customer Installation
	idReturnInstallation, err = dao.CustomerInstallationDAO.InsertCustomerInstallation(tx, lv1, idxSite, idxInst, *queue)
	if err.Error != nil {
		return
	}

	//--- Record Audit
	*dataAudit = append(*dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.CustomerInstallationDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: idReturnInstallation},
		Action:     sql.NullInt32{Int32: constanta.ActionAuditInsertConstanta},
	})

	//--- Record Unique
	tempRepo = input.prepareTrackingUniqueID(tempTrackInst)
	trackUnique[idxSite] = tempRepo

	//--- Finish
	err = errorModel.GenerateNonErrorModel()
	return
}
