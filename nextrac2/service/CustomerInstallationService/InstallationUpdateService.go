package CustomerInstallationService

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

func (input customerInstallationService) doUpdateInstallationService(tx *sql.Tx, customerInstallationModel repository.CustomerInstallationModel, lvl2 repository.CustomerInstallationData, lvl3 repository.CustomerInstallationDetail, idxSite, idxInst int, prepErrUniqueKey, prepErrParentNoExist string, resultDB repository.CustomerInstallationDetail,
	trackUnique map[int][]repository.CustomerInstallationTracking, tempTrackInst map[repository.CustomerInstallationTracking]bool, isCheckParentInst bool, parentClientType int64, contextModel *applicationModel.ContextModel, timeNow time.Time, dataAudit *[]repository.AuditSystemModel) (err errorModel.ErrorModel) {

	//--- Check Duplication Unique Key
	err = input.checkDuplicateUnique12(customerInstallationModel, lvl3, idxSite, idxInst, prepErrUniqueKey, trackUnique, tempTrackInst)
	if err.Error != nil {
		return
	}

	//--- Check Parent Customer Installation
	if isCheckParentInst {
		err = input.checkParentCustomerInstallation(customerInstallationModel, lvl2, lvl3, idxInst, parentClientType, contextModel, prepErrParentNoExist)
		if err.Error != nil {
			return
		}
	} else {
		if !(lvl3.UniqueID1.String == resultDB.UniqueID1.String && lvl3.UniqueID2.String == resultDB.UniqueID2.String && lvl3.ProductID.Int64 == resultDB.ProductID.Int64) {
			//--- Validate Update Parent Customer Installation
			err = input.validateUpdateCustomerInstallation(resultDB, idxInst, lvl2, lvl3, prepErrParentNoExist)
			if err.Error != nil {
				return
			}
		}
	}

	//--- Update Customer Installation
	*dataAudit = append(*dataAudit, service.GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *contextModel, timeNow, dao.CustomerInstallationDAO.TableName, resultDB.InstallationID.Int64, 0)...)
	err = dao.CustomerInstallationDAO.UpdateCustomerInstallationByInstallationID(tx, customerInstallationModel, idxSite, idxInst)
	if err.Error != nil {
		return
	}

	tempRepo := input.prepareTrackingUniqueID(tempTrackInst)
	trackUnique[idxSite] = tempRepo
	return
}

func (input customerInstallationService) validateUpdateCustomerInstallation(resultDB repository.CustomerInstallationDetail, idxInst int, lvl2 repository.CustomerInstallationData, lvl3 repository.CustomerInstallationDetail, prepErrParentNoExist string) (err errorModel.ErrorModel) {
	var (
		fileName          = "InstallationUpdateService.go"
		funcName          = "validateUpdateCustomerInstallation"
		db                = serverconfig.ServerAttribute.DBConnection
		installationTemp  []repository.CustomerInstallationDetail
		isParent          bool
		isValid           bool
		childIDProduct    map[int64]bool
		parentIDProduct   map[int64]bool
		parentInfo        repository.ProductModel
		childInstallation map[int64]repository.CustomerInstallationDetail
		otherParent       map[int64]repository.CustomerInstallationDetail
		parentInstModel   repository.CustomerInstallationForConfig
	)

	//--- Check Parent App Valid
	isParent, parentInfo, err = dao.ProductDAO.CheckValidParentProductByProductID(db, repository.ProductModel{ID: sql.NullInt64{Int64: resultDB.ProductID.Int64}})
	if err.Error != nil {
		return
	}

	if isParent {

		//--- Parent ID Product
		parentIDProduct, err = dao.ProductDAO.GetProductParentByClientTypeID(db, repository.ProductModel{ClientTypeID: sql.NullInt64{Int64: parentInfo.ParentClientTypeID.Int64}})
		if err.Error != nil {
			return
		}

		//--- Get Child ID Product and ID Installation
		childIDProduct, childInstallation, parentInstModel, err = input.doGetProductChildAndInstallationChild(db, parentInfo, resultDB, lvl2, lvl3)
		if err.Error != nil {
			return
		}

		//--- Validate All Request
		isValid, installationTemp, err = input.doValidateAllRequestWithDataOnDBForUpdate(childIDProduct, childInstallation, lvl2, idxInst, resultDB, parentIDProduct)
		if err.Error != nil {
			return
		}

		//--- Valid Then Return
		if isValid {
			return
		}

		if len(childInstallation) > 0 {

			//--- Check Any Parent
			otherParent, err = dao.CustomerInstallationDAO.CheckOtherParentInstallation(db, parentInstModel, parentInfo)
			if err.Error != nil {
				return
			}

			for _, itemInst := range installationTemp {
				v, ok := otherParent[itemInst.InstallationID.Int64]
				if ok {

					//--- If Delete Then Delete On DB Result
					if itemInst.Action.Int32 == constanta.ActionDeleteCode {
						delete(otherParent, itemInst.InstallationID.Int64)
						continue
					}

					//--- If Unique ID Key Ok And Product Ok
					if itemInst.UniqueID1.String == v.UniqueID1.String && itemInst.UniqueID2.String == v.UniqueID2.String && itemInst.ProductID.Int64 == v.ProductID.Int64 {
						return
					} else {
						delete(childInstallation, itemInst.InstallationID.Int64)
						continue
					}
				}

				if itemInst.Action.Int32 == constanta.ActionInsertCode || itemInst.Action.Int32 == constanta.ActionUpdateCode {
					if itemInst.UniqueID1.String == lvl3.UniqueID1.String && itemInst.UniqueID2.String == lvl3.UniqueID2.String {
						_, okParent := parentIDProduct[itemInst.ProductID.Int64]
						if okParent {
							return
						}
					}
				}
			}

			//--- On Update Service
			if lvl3.Action.Int32 == constanta.ActionUpdateCode && resultDB.UniqueID1.String == lvl3.UniqueID1.String && resultDB.UniqueID2.String == lvl3.UniqueID2.String {
				_, okParent := parentIDProduct[lvl3.ProductID.Int64]
				if okParent {
					return
				}
			}

			if len(otherParent) < 1 {
				err = errorModel.GenerateErrorParentAppUpdatedDeleted(fileName, funcName, prepErrParentNoExist)
				return
			}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input customerInstallationService) doValidateAllRequestWithDataOnDBForUpdate(childIDProduct map[int64]bool, childInstallation map[int64]repository.CustomerInstallationDetail, lvl2 repository.CustomerInstallationData, idxInst int,
	resultDB repository.CustomerInstallationDetail, parentIDProduct map[int64]bool) (isValid bool, installationTemp []repository.CustomerInstallationDetail, err errorModel.ErrorModel) {
	installationTemp = make([]repository.CustomerInstallationDetail, len(lvl2.Installation))

	if len(childInstallation) > 0 {
		//--- Prepare Request Data
		copy(installationTemp, lvl2.Installation)
		installationTemp = append(installationTemp[:idxInst], installationTemp[idxInst+1:]...)
		if len(installationTemp) < 1 {
			//--- If Insert New Installation Product Parent
			if (lvl2.Installation[0].Action.Int32 == constanta.ActionInsertCode || lvl2.Installation[0].Action.Int32 == constanta.ActionUpdateCode) && lvl2.Installation[0].UniqueID1.String == resultDB.UniqueID1.String && lvl2.Installation[0].UniqueID2.String == resultDB.UniqueID2.String {
				_, okParent := parentIDProduct[lvl2.Installation[0].ProductID.Int64]
				if okParent {
					isValid = true
					return
				}
			}
		}

		//--- Validate Request Data
		for _, itemInst := range installationTemp {
			//--- If Insert New Installation Product Parent
			if itemInst.InstallationID.Int64 < 1 {
				if itemInst.Action.Int32 == constanta.ActionInsertCode && itemInst.UniqueID1.String == resultDB.UniqueID1.String && itemInst.UniqueID2.String == resultDB.UniqueID2.String {
					_, okParent := parentIDProduct[itemInst.ProductID.Int64]
					if okParent {
						isValid = true
						return
					}
				}
				continue
			}

			//--- If Data Exist Then Validate (Update or Delete)
			v, ok := childInstallation[itemInst.InstallationID.Int64]
			if ok {
				//--- If Delete Then Delete On DB Result
				if itemInst.Action.Int32 == constanta.ActionDeleteCode {
					delete(childInstallation, itemInst.InstallationID.Int64)
					continue
				}

				//--- If Unique ID Key Ok Then Check Product
				if itemInst.UniqueID1.String == v.UniqueID1.String && itemInst.UniqueID2.String == v.UniqueID2.String {
					if itemInst.ProductID.Int64 != v.ProductID.Int64 {
						_, okProduct := childIDProduct[itemInst.ProductID.Int64]
						if !okProduct {
							delete(childInstallation, itemInst.InstallationID.Int64)
							continue
						}
					}
				} else {
					delete(childInstallation, itemInst.InstallationID.Int64)
					continue
				}
			} else {
				if itemInst.UniqueID1.String == resultDB.UniqueID1.String && itemInst.UniqueID2.String == resultDB.UniqueID2.String {
					_, okParent := parentIDProduct[itemInst.ProductID.Int64]
					if okParent {
						isValid = true
						return
					}
				}
			}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
