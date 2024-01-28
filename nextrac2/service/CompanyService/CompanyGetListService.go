package CompanyService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
)

type companyGetListService struct {
	service.GetListData
	FileName string
}

var CompanyGetListService = companyGetListService{}.New()

func (input companyGetListService) New() (output companyGetListService) {
	output.FileName = "CompanyGetListService.go"
	output.ValidSearchBy = []string{"company_name"}
	output.ValidOrderBy = []string{"id"}
	output.ValidLimit = service.DefaultLimit
	return
}

func (input companyGetListService) GetCompanyList(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var inputStruct in.GetListDataDTO
	var searchByParam []in.SearchByParam
	var createdBy int64 = 0

	inputStruct, searchByParam, err = input.ReadAndValidateGetListData(request, input.ValidSearchBy, input.ValidOrderBy, applicationModel.GetListInternalCompanyValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	createdBy, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*context)
	if !isOnlyHaveOwnAccess {
		createdBy = 0
	}

	companies, err := dao.CompanyDAO.GetCompanyList(serverconfig.ServerAttribute.DBConnection, inputStruct, searchByParam, createdBy)
	if err.Error != nil {
		return
	}

	output.Data.Content = input.convertRepoToDTO(companies)
	output.Status = service.GetResponseMessages("SUCCESS_GET_MESSAGE", context)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input companyGetListService) InitiateCompany(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var searchByParam []in.SearchByParam
	_, searchByParam, err = input.ReadAndValidateGetCountData(request, input.ValidSearchBy, applicationModel.GetListInternalCompanyValidOperator)
	if err.Error != nil {
		return
	}

	countData, err := dao.CompanyDAO.GetCountCompany(serverconfig.ServerAttribute.DBConnection, searchByParam)
	if err.Error != nil {
		return
	}

	output.Data.Content = out.InitiateGetListDataDTOOut{
		ValidOrderBy:  input.ValidOrderBy,
		ValidSearchBy: input.ValidSearchBy,
		ValidLimit:    input.ValidLimit,
		CountData:     int(countData),
		ValidOperator: applicationModel.GetListInternalCompanyValidOperator,
	}
	output.Status = service.GetResponseMessages("SUCCESS_INITIATE_MESSAGE", context)

	return
}

func (input companyGetListService) convertRepoToDTO(data []interface{}) (companies []out.CompanyForViewResponse) {
	for _, item := range data {
		company := item.(repository.CompanyModel)
		file, _ := dao.FileUploadDAO.GetFileByParentIDAndCategory(serverconfig.ServerAttribute.DBConnection, company.ID.Int64, dao.CompanyDAO.TableName)
		companies = append(companies, out.CompanyForViewResponse{
			ID:                company.ID.Int64,
			CompanyTitle:      company.CompanyTitle.String,
			CompanyName:       company.CompanyName.String,
			PhotoIcon:         file.Host.String + file.Path.String,
			Address:           company.Address.String,
			Telephone:         company.Telephone.String,
			UpdatedAt:         company.UpdatedAt.Time,
		})
	}
	return
}

