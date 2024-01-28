package AuditMonitoringService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type auditMonitoringService struct {
	service.AbstractService
	service.GetListData
}

var AuditMonitoringService = auditMonitoringService{}.New()

func (input auditMonitoringService) New() (output auditMonitoringService) {
	output.FileName = "AuditMonitoringService.go"
	output.ServiceName = "AUDIT_SYSTEM"
	output.ValidLimit = []int{5}
	output.ValidLimit = append(output.ValidLimit, service.DefaultLimit...)
	output.ValidOrderBy = []string{"created_at", "id", "table_name", "created_name"}
	output.ValidSearchBy = []string{"table_name", "primary_key", "created_by", "created_client", "menu_code"}

	output.MappingScopeDB = make(map[string]applicationModel.MappingScopeDB)
	output.MappingScopeDB[constanta.CustomerGroupDataScope] = applicationModel.MappingScopeDB{
		View: dao.CustomerGroupDAO.TableName,
	}

	output.MappingScopeDB[constanta.CustomerCategoryDataScope] = applicationModel.MappingScopeDB{
		View: dao.CustomerCategoryDAO.TableName,
	}

	output.MappingScopeDB[constanta.ProvinceDataScope] = applicationModel.MappingScopeDB{
		View: dao.ProvinceDAO.TableName,
	}

	output.MappingScopeDB[constanta.DistrictDataScope] = applicationModel.MappingScopeDB{
		View: dao.DistrictDAO.TableName,
	}

	output.MappingScopeDB[constanta.ProductGroupDataScope] = applicationModel.MappingScopeDB{
		View: dao.ProductGroupDAO.TableName,
	}

	output.MappingScopeDB[constanta.ClientTypeDataScope] = applicationModel.MappingScopeDB{
		View: dao.ClientTypeDAO.TableName,
	}

	output.MappingScopeDB[constanta.SalesmanDataScope] = applicationModel.MappingScopeDB{
		View: dao.SalesmanDAO.TableName,
	}
	return
}

func (input auditMonitoringService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.AuditMonitoringRequest) errorModel.ErrorModel) (inputStruct in.AuditMonitoringRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)

	return
}

func (input auditMonitoringService) validateDataScope(contextModel *applicationModel.ContextModel) (output map[string]interface{}, err errorModel.ErrorModel) {
	output = service.ValidateScope(contextModel, []string{
		constanta.CustomerGroupDataScope,
		constanta.CustomerCategoryDataScope,
		constanta.ProvinceDataScope,
		constanta.DistrictDataScope,
		constanta.ProductGroupDataScope,
		constanta.ClientTypeDataScope,
		constanta.SalesmanDataScope,
	})
	if output == nil {
		err = errorModel.GenerateDataScopeNotDefinedYet(input.FileName, "validateDataScope")
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
