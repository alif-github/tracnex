package GetSessionService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service/session"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input getSessionService) GetDashboardView(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		timeNow                    = time.Now()
		dataActiveCustomerSiteOnDB int64
		dataActiveLicenseOnDB      []repository.ProductLicenseModel
	)

	dataActiveCustomerSiteOnDB, err = dao.CustomerSiteDAO.GetActiveCustomerSite(serverconfig.ServerAttribute.DBConnection)
	if err.Error != nil {
		return
	}

	dataActiveLicenseOnDB, err = dao.ProductLicenseDAO.GetActiveLicenseForDashboard(serverconfig.ServerAttribute.DBConnection)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.convertToDashboardView(dataActiveCustomerSiteOnDB, dataActiveLicenseOnDB, timeNow)

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: session.GenerateLoginI18NMessage("GET_SESSION_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}
	return
}

func (input getSessionService) convertToDashboardView(dataActiveCustomerSiteOnDB int64, dataActiveLicenseOnDB []repository.ProductLicenseModel, timeNow time.Time) (output out.ViewDashboardResponse) {
	// Set Total Active Customer
	output.TotalCustomerActive = dataActiveCustomerSiteOnDB
	timeExpired1 := time.Date(timeNow.Year(), timeNow.Month()+1, 1, 0, 0, 0, 0, timeNow.Location())
	timeExpired2 := time.Date(timeNow.Year(), timeNow.Month()+2, 1, 0, 0, 0, 0, timeNow.Location())

	// Set Total Active License
	output.TotalLicense = out.ViewLicenseDashboard{
		Total: int64(len(dataActiveLicenseOnDB)),
		Detail: []out.DetailViewLicenseDashboard{
			{
				Month: int64(timeExpired1.Month()),
				Year:  int64(timeExpired1.Year()),
			},
			{
				Month: int64(timeExpired2.Month()),
				Year:  int64(timeExpired2.Year()),
			},
		},
	}

	for _, item := range dataActiveLicenseOnDB {
		if item.ProductValidThru.Time.Month() == timeExpired1.Month() {
			output.TotalLicense.Detail[0].Total++
		} else if item.ProductValidThru.Time.Month() == timeExpired2.Month() {
			output.TotalLicense.Detail[1].Total++
		}
	}
	return
}
