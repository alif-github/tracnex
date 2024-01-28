package CustomerInstallationService

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
)

func GenerateI18NMessage(messageID string, language string) (output string) {
	return util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.CustomerInstallationBundle, messageID, language, nil)
}

func (input customerInstallationService) trackingUniqueID(trackUnique map[int][]repository.CustomerInstallationTracking, idxSite int, uniqueID1, uniqueID2, message string) (err errorModel.ErrorModel) {
	var (
		fileName   = "CustomerInstallationServiceUtil.go"
		funcName   = "trackingUniqueID"
		errMessage = fmt.Sprintf(`error because duplicate in insert site`)
	)

	for key, value := range trackUnique {
		if idxSite != key {
			for _, valueItem := range value {
				if uniqueID1 == valueItem.UniqueID1.String && uniqueID2 == valueItem.UniqueID2.String {
					service.LogMessage(errMessage, 200)
					err = errorModel.GenerateDataUsedError(fileName, funcName, message)
					return
				}
			}
		}

		if err.Error != nil {
			return
		}
	}

	return
}

func (input customerInstallationService) prepareTrackingUniqueID(tempTrackInst map[repository.CustomerInstallationTracking]bool) (tempRepo []repository.CustomerInstallationTracking) {
	for k := range tempTrackInst {
		tempRepo = append(tempRepo, repository.CustomerInstallationTracking{
			KeyID:     sql.NullInt64{Int64: k.KeyID.Int64},
			UniqueID1: sql.NullString{String: k.UniqueID1.String},
			UniqueID2: sql.NullString{String: k.UniqueID2.String},
		})
	}
	return
}

func (input customerInstallationService) generateConstanta(constanta string, contextModel *applicationModel.ContextModel) (result string) {
	return util.GenerateConstantaI18n(constanta, contextModel.AuthAccessTokenModel.Locale, nil)
}

func (input customerInstallationService) getInstallationNumber(queue *int64, lv1 repository.CustomerInstallationModel, idxSite int) (err errorModel.ErrorModel) {
	if *queue < 1 {
		*queue, err = dao.CustomerInstallationDAO.GetInstallationNumberLastInstallation(serverconfig.ServerAttribute.DBConnection, lv1, idxSite)
		if err.Error != nil {
			return
		}
	} else {
		*queue += 1
	}
	return
}
