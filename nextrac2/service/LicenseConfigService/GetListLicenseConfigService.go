package LicenseConfigService

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

func (input licenseConfigService) InitiateGetListLicenseConfig(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		count         int
		newScope      map[string]interface{}
		db            = serverconfig.ServerAttribute.DBConnection
	)

	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListLicenseConfigValidOperator)
	if err.Error != nil {
		return
	}

	//--- Validate Data Scope
	newScope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	count, err = dao.LicenseConfigDAO.GetCountLicenseConfig(db, searchByParam, contextModel.LimitedByCreatedBy, newScope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	enum := make(map[string][]string)
	enum["allow_activation"] = []string{"Y", "N"}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		ValidOperator: applicationModel.GetListLicenseConfigValidOperator,
		EnumData:      enum,
		CountData:     count,
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INITIATE_GET_LIST_LICENSE_CONFIG_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) GetListLicenseConfig(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		fileName      = "GetListLicenseConfigService.go"
		funcName      = "GetListLicenseConfig"
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListLicenseConfigValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	if inputStruct.OrderBy == input.ValidOrderBy[0] {
		inputStruct.OrderBy = input.ValidOrderBy[0]
	}

	//--- Check Enum Allow Activation
	for _, itemSearchByParam := range searchByParam {
		if itemSearchByParam.SearchKey == "allow_activation" {
			if itemSearchByParam.SearchValue != "Y" && itemSearchByParam.SearchValue != "N" {
				err = errorModel.GenerateFormatFieldError(fileName, funcName, constanta.AllowActivation)
				return
			}
		}
	}

	output.Data.Content, err = input.doGetListLicenseConfig(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_LICENSE_CONFIG_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) SelectAllLicenseConfig(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct   in.GetListDataDTO
		searchByParam []in.SearchByParam
	)

	inputStruct, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListLicenseConfigValidOperator)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doSelectAllLicenseConfig(inputStruct, searchByParam, contextModel)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_LIST_LICENSE_CONFIG_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) doGetListLicenseConfig(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult []interface{}
		newScope map[string]interface{}
		db       = serverconfig.ServerAttribute.DBConnection
	)

	newScope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	if inputStruct.OrderBy == input.ValidOrderBy[0] {
		inputStruct.OrderBy = input.ValidOrderBy[0]
	}
	dbResult, err = dao.LicenseConfigDAO.GetListLicenseConfig(db, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, newScope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	output = input.convertToListDTOOut(dbResult)
	return
}

func (input licenseConfigService) doSelectAllLicenseConfig(inputStruct in.GetListDataDTO, searchByParam []in.SearchByParam, contextModel *applicationModel.ContextModel) (output interface{}, err errorModel.ErrorModel) {
	var (
		dbResult         []interface{}
		newScope         map[string]interface{}
		LicenseConfigIDs []int64
		db               = serverconfig.ServerAttribute.DBConnection
	)

	newScope, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	inputStruct.OrderBy = input.ValidOrderBy[0]
	dbResult, err = dao.LicenseConfigDAO.SelectAllLicenseConfigGetID(db, inputStruct, searchByParam, contextModel.LimitedByCreatedBy, newScope, input.MappingScopeDB)
	if err.Error != nil {
		return
	}

	for _, dbResultItem := range dbResult {
		id := dbResultItem.(int64)
		LicenseConfigIDs = append(LicenseConfigIDs, id)
	}

	output = out.ListLicenseConfigIDs{
		TotalLicenseConfigID: int64(len(LicenseConfigIDs)),
		LicenseConfigID:      LicenseConfigIDs,
	}
	return
}

func (input licenseConfigService) convertToListDTOOut(dbResult []interface{}) (result []out.ListLicenseConfigModel) {
	for _, dbResultItem := range dbResult {
		var repo repository.LicenseConfigModel
		repo = dbResultItem.(repository.LicenseConfigModel)
		result = append(result, out.ListLicenseConfigModel{
			LicenseConfigID:    repo.ID.Int64,
			CustomerName:       repo.Customer.String,
			UniqueID1:          repo.UniqueID1.String,
			UniqueID2:          repo.UniqueID2.String,
			InstallationID:     repo.InstallationID.Int64,
			ProductName:        repo.ProductName.String,
			ClientTypeID:       repo.ClientTypeID.Int64,
			LicenseVariantName: repo.LicenseVariantName.String,
			LicenseTypeName:    repo.LicenseTypeName.String,
			ProductValidFrom:   repo.ProductValidFrom.Time.Format(constanta.DefaultInstallationTimeFormat),
			ProductValidThru:   repo.ProductValidThru.Time.Format(constanta.DefaultInstallationTimeFormat),
			AllowActivation:    repo.AllowActivation.String,
			IsExtendChecklist:  repo.IsExtendChecklist.Bool,
			PaymentStatus:      repo.PaymentStatus.String,
			UpdatedAt:          repo.UpdatedAt.Time,
		})
	}

	return result
}
