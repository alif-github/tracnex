package ClientTypeService

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
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

func (input clientTypeService) InsertClientType(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "InsertClientType"
	var inputStruct in.ClientTypeRequest

	inputStruct, err = input.readBodyAndValidate(request, contextModel, input.ValidateInsertClientType)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertClientType, nil)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_CLIENT_TYPE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input clientTypeService) doInsertClientType(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct   = inputStructInterface.(in.ClientTypeRequest)
		inputModel    = input.convertDTOToModel(inputStruct, contextModel.AuthAccessTokenModel, timeNow)
		dataAuditTemp repository.AuditSystemModel
	)

	err = input.validateRemark(inputStruct, 3)
	if err.Error != nil {
		return
	}

	insertedID, err := dao.ClientTypeDAO.InsertClientType(tx, inputModel)
	if err.Error != nil {
		return
	}

	dataAuditTemp, err = input.GenerateDataScope(tx, insertedID, dao.ClientTypeDAO.TableName, constanta.ClientTypeDataScope, contextModel.AuthAccessTokenModel.ResourceUserID, contextModel.AuthAccessTokenModel.ClientID, timeNow)
	if err.Error != nil {
		return
	}

	dataAudit = append(dataAudit, dataAuditTemp)
	dataAudit = append(dataAudit, repository.AuditSystemModel{
		TableName:  sql.NullString{String: dao.ClientTypeDAO.TableName},
		PrimaryKey: sql.NullInt64{Int64: insertedID},
	})

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input clientTypeService) validateRemark(inputStruct in.ClientTypeRequest, levelRemark int64) (err errorModel.ErrorModel) {
	funcName := "validateRemark"
	for _, remark := range inputStruct.Remarks {
		var (
			remarksOnDB []repository.RemarkModel
			remarkOnDB  repository.RemarkModel
		)

		remarkOnDB, err = dao.RemarkDAO.GetRemarkByID(serverconfig.ServerAttribute.DBConnection, repository.RemarkModel{
			ID: sql.NullInt64{Int64: remark},
		})

		if err.Error != nil {
			return
		}

		if remarkOnDB.ID.Int64 == 0 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.Remark)
			return
		}

		if remarkOnDB.Level.Int64 != levelRemark {
			err = errorModel.GenerateCannotBeUsedError(input.FileName, funcName, constanta.Remark)
			return
		}

		remarksOnDB, err = dao.RemarkDAO.GetParentFromChild(serverconfig.ServerAttribute.DBConnection, repository.RemarkModel{
			ID: sql.NullInt64{Int64: remark},
		})

		if err.Error != nil {
			return
		}

		for _, remarkModel := range remarksOnDB {
			if remarkModel.ParentID.Int64 == 0 {
				remarkOnDB = remarkModel
				break
			}
		}

		if remarkOnDB.Name.String != constanta.ActionClientTypeRemarkName {
			err = errorModel.GenerateCannotBeUsedError(input.FileName, funcName, constanta.Remark)
			return
		}
	}

	return
}

func (input clientTypeService) convertDTOToModel(inputStruct in.ClientTypeRequest, authAccessToken model.AuthAccessTokenModel, timeNow time.Time) repository.ClientTypeModel {
	return repository.ClientTypeModel{
		ID:                 sql.NullInt64{Int64: inputStruct.ID},
		ClientType:         sql.NullString{String: inputStruct.ClientType},
		Description:        sql.NullString{String: inputStruct.Description},
		ParentClientTypeID: sql.NullInt64{Int64: inputStruct.ParentClientTypeID},
		CreatedBy:          sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		CreatedAt:          sql.NullTime{Time: timeNow},
		CreatedClient:      sql.NullString{String: authAccessToken.ClientID},
		UpdatedBy:          sql.NullInt64{Int64: authAccessToken.ResourceUserID},
		UpdatedAt:          sql.NullTime{Time: timeNow},
		UpdatedClient:      sql.NullString{String: authAccessToken.ClientID},
	}
}

func (input clientTypeService) convertDTOToModelForView(repo repository.ClientTypeModel) out.ListClientTypeResponse {
	return out.ListClientTypeResponse{
		ID:          repo.ID.Int64,
		ClientType:  repo.ClientType.String,
		Description: repo.Description.String,
		CreatedAt:   repo.CreatedAt.Time,
		UpdatedName: repo.UpdatedName.String,
		UpdatedAt:   repo.UpdatedAt.Time,
	}
}

func (input clientTypeService) ValidateInsertClientType(inputStruct *in.ClientTypeRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsert()
}
