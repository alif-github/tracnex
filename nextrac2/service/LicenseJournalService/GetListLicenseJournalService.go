package LicenseJournalService

import (
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
)

func (input licenseJournalService) InitiateGetListLicenseJournal(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		db            = serverconfig.ServerAttribute.DBConnection
		searchByParam []in.SearchByParam
		count         int
		inputStruct   in.LicenseJournalRequest
	)

	//--- Read Body And Validate
	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateCountDataLicenseJournal)
	if err.Error != nil {
		return
	}

	if !inputStruct.LastSync.IsZero() {
		searchByParam = []in.SearchByParam{{
			SearchKey:   "pl.updated_at >=",
			SearchValue: inputStruct.LastSync.Format(constanta.DefaultTimeFormat),
			SearchType:  constanta.Filter,
		}}
	}

	//--- Query Get Count
	count, err = dao.ProductLicenseDAO.GetCountProductLicenseSalesJournal(db, searchByParam)
	if err.Error != nil {
		return
	}

	output.Data.Content = count
	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_GET_LIST_LICENSE_JOURNAL_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseJournalService) GetListLicenseJournal(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		db            = serverconfig.ServerAttribute.DBConnection
		searchByParam []in.SearchByParam
		result        []interface{}
		inputStruct   in.LicenseJournalRequest
		userParam     in.GetListDataDTO
	)

	//--- Read Body And Validate
	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.validateGetListDataLicenseJournal)
	if err.Error != nil {
		return
	}

	//--- Last Sync If Exist
	if !inputStruct.LastSync.IsZero() {
		searchByParam = []in.SearchByParam{{
			SearchKey:   "pl.updated_at >=",
			SearchValue: inputStruct.LastSync.Format(constanta.DefaultTimeFormat),
			SearchType:  constanta.Filter,
		}}
	}

	//--- Convert To In Get List Data DTO
	if inputStruct.Limit > 0 {
		userParam = in.GetListDataDTO{
			AbstractDTO: in.AbstractDTO{
				Page:  inputStruct.Page,
				Limit: inputStruct.Limit,
			},
		}
	}

	//--- Query Get Count
	result, err = dao.ProductLicenseDAO.GetProductLicenseSalesJournal(db, userParam, searchByParam)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.convertToListDTOOut(result)
	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_LICENSE_JOURNAL_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseJournalService) convertToListDTOOut(dbResult []interface{}) (result []out.ListLicenseJournal) {
	for _, dbResultItem := range dbResult {
		var repo repository.LicenseSalesJournal
		repo = dbResultItem.(repository.LicenseSalesJournal)
		result = append(result, out.ListLicenseJournal{
			ID:                repo.ID.Int64,
			ClientID:          repo.ClientID.String,
			LicenseStatus:     repo.LicenseStatusID.Int64,
			StatusDescription: repo.LicenseStatus.String,
			UniqueID1:         repo.UniqueID1.String,
			UniqueID2:         repo.UniqueID2.String,
			ProductName:       repo.ProductName.String,
			ClientType:        repo.ClientType.String,
			AllowActivation:   repo.AllowActivation.String,
			NoOfUser:          repo.NoOfUser.Int64,
			ProductValidFrom:  repo.ProductValidFrom.Time.Format(constanta.DefaultInstallationTimeFormat),
			ProductValidThru:  repo.ProductValidThru.Time.Format(constanta.DefaultInstallationTimeFormat),
			IsUserConcurrent:  repo.IsUserConcurrent.String,
			TotalLicense:      repo.TotalLicense.Int64,
			TotalActivated:    repo.TotalActivated.Int64,
		})
	}

	return result
}

func (input licenseJournalService) validateCountDataLicenseJournal(inputStruct *in.LicenseJournalRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInitiateLicenseJournal()
}

func (input licenseJournalService) validateGetListDataLicenseJournal(inputStruct *in.LicenseJournalRequest) errorModel.ErrorModel {
	return inputStruct.ValidateGetListLicenseJournal()
}
