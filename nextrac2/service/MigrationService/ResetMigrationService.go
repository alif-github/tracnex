package MigrationService

import (
	"database/sql"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"time"
)

func (input migrationService) ResetMigration(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName = "ResetMigration"
	)

	inputStruct, err := input.readBodyAndValidate(request, contextModel, input.validateResetMigration)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.ServiceWithDataAuditPreparedByService(funcName, inputStruct, contextModel, input.doResetMigration, func(interface{}, applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_DELETE_MESSAGE", contextModel)
	return
}

func (input migrationService) doResetMigration(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (result interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct = inputStructInterface.(in.ResetMigrationRequest)
	)

	if !inputStruct.Reset {
		// Created Input Model
		for _, idMigration := range inputStruct.ID {
			model := repository.MigrationModel{
				ID: sql.NullString{String: idMigration},
			}

			// Reset Migration Selected
			err = dao.MigrationDAO.ResetMigration(tx, model)
			if err.Error != nil {
				return
			}
		}
	} else {
		// Reset All Migration
		err = dao.MigrationDAO.ResetMigration(tx, repository.MigrationModel{})
		if err.Error != nil {
			return
		}
	}

	return
}

func (input migrationService) validateResetMigration(inputStruct *in.ResetMigrationRequest) errorModel.ErrorModel {
	return inputStruct.ValidateReset()
}
