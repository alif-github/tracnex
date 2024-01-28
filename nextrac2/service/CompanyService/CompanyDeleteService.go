package CompanyService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type companyDeleteService struct {
	FileName string
	service.AbstractService
}

var CompanyDeleteService = companyDeleteService{FileName: "CompanyDeleteService.go"}

func (input companyDeleteService) DeleteCompany(request *http.Request, context *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	inputStruct, err := input.readParamAndBody(request, context)
	if err.Error != nil {
		return
	}

	_, err = input.ServiceWithDataAuditPreparedByService("DeleteCompany", inputStruct, context, input.deleteCompany, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", context)

	return
}

func (input companyDeleteService) deleteCompany(tx *sql.Tx, body interface{}, context *applicationModel.ContextModel, now time.Time) (_ interface{}, auditData []repository.AuditSystemModel, err errorModel.ErrorModel) {
	funcName := "deleteCompany"
	authAccessToken := context.AuthAccessTokenModel

	companyBody := body.(in.CompanyRequest)

	companyRepository := input.getCompanyRepository(companyBody.ID, authAccessToken, now)
	if err.Error != nil {
		return
	}

	userID, isOnlyHaveOwnAccess := service.CheckIsOnlyHaveOwnPermission(*context)
	if isOnlyHaveOwnAccess {
		companyRepository.CreatedBy.Int64 = userID
	}

	companyOnDB, err := dao.CompanyDAO.GetDetailCompanyForUpdate(tx, companyRepository.ID.Int64)
	if err.Error != nil {
		return
	}

	if companyOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, "company")
		return
	}

	err = input.validation(companyOnDB, companyBody)
	if err.Error != nil {
		return
	}

	auditData = append(auditData, service.GetAuditData(tx, constanta.ActionAuditDeleteConstanta, *context, now, dao.CompanyDAO.TableName, companyOnDB.ID.Int64, userID)...)

	err = dao.CompanyDAO.DeleteCompany(tx, companyRepository)
	if err.Error != nil {
		return
	}

	return
}

func (input companyDeleteService) readParamAndBody(request *http.Request, contextModel *applicationModel.ContextModel) (inputStruct in.CompanyRequest, err errorModel.ErrorModel) {
	id, err := service.ReadPathParamID(request)
	if err.Error != nil {
		return
	}

	inputStruct, bodySize, err := getCompanyBody(request, input.FileName)
	inputStruct.ID = id
	contextModel.LoggerModel.ByteIn = bodySize
	return
}

func (input companyDeleteService) getCompanyRepository(id int64, authAccessToken model2.AuthAccessTokenModel, now time.Time) repository.CompanyModel {
	return repository.CompanyModel{
		ID:            sql.NullInt64{Int64: id},
		UpdatedClient: sql.NullString{String: authAccessToken.ClientID},
		UpdatedAt:     sql.NullTime{Time: now},
		UpdatedBy:     sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		Deleted:       sql.NullBool{Bool: true},
	}
}

func (input companyDeleteService) validation(companyOnDB repository.CompanyModel, companyBody in.CompanyRequest) (err errorModel.ErrorModel) {
	result, err := TimeStrToTime(companyBody.UpdatedAtStr, "updated_at")
	companyBody.UpdatedAt = result
	if err.Error != nil {
		return
	}

	err = service.OptimisticLock(companyOnDB.UpdatedAt.Time, companyBody.UpdatedAt, input.FileName, "company")
	return
}

func TimeStrToTime(timeStr string, fieldName string) (output time.Time, error errorModel.ErrorModel) {
	output, err := time.Parse(constanta.DefaultTimeFormat, timeStr)

	if err != nil {
		error = errorModel.GenerateFormatFieldError("AbstractDTO.go", "TimeStrToTime", fieldName)
		return
	}
	return output, errorModel.GenerateNonErrorModel()
}

