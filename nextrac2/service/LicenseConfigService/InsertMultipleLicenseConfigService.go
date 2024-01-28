package LicenseConfigService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
	"strings"
	"time"
)

func (input licenseConfigService) InsertMultipleLicenseConfig(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName    = "InsertMultipleLicenseConfig"
		inputStruct in.LicenseConfigMultipleRequest
	)

	inputStruct, err = input.readBodyAndValidateMultipleLicenseConfig(request, contextModel, input.validateInsertMultiple)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertMultipleLicenseConfig, func(_ interface{}, _ applicationModel.ContextModel) {
		//--- Additional function
	})

	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_INSERT_LICENSE_CONFIG", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) doInsertMultipleLicenseConfig(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (output interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		inputStruct          = inputStructInterface.(in.LicenseConfigMultipleRequest)
		mapLicenseConfigID   = make(map[int64]bool)
		db                   = serverconfig.ServerAttribute.DBConnection
		licenseConfigModel   repository.LicenseConfigModel
		licenseConfigModelDB []repository.LicenseConfigModel
		scopeLimit           map[string]interface{}
	)

	//--- Validate data scope
	scopeLimit, err = input.validateDataScope(contextModel)
	if err.Error != nil {
		return
	}

	//--- Create Model
	err = input.createModelForInsertMultiple(inputStruct, &licenseConfigModel, mapLicenseConfigID, contextModel, timeNow)
	if err.Error != nil {
		return
	}

	//--- Check License Config ID
	licenseConfigModelDB, err = input.validateIDLicenseConfig(db, licenseConfigModel, mapLicenseConfigID, scopeLimit, timeNow)
	if err.Error != nil {
		return
	}

	//--- Insert License Config
	dataAudit, err = input.insertLicenseConfig(tx, licenseConfigModelDB)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) validateIDLicenseConfig(db *sql.DB, licenseConfigModel repository.LicenseConfigModel, mapLicenseConfigID map[int64]bool, scopeLimit map[string]interface{}, timeNow time.Time) (licenseConfigModelOnDB []repository.LicenseConfigModel, err errorModel.ErrorModel) {
	var (
		fileName               = "InsertMultipleLicenseConfigService.go"
		funcName               = "validateIDLicenseConfig"
		excludeLicenseConfigID []string
	)

	licenseConfigModelOnDB, err = dao.LicenseConfigDAO.GetLicenseConfigDataForDuplicate(db, licenseConfigModel, scopeLimit, input.MappingScopeDB, timeNow)
	if err.Error != nil {
		return
	}

	if len(licenseConfigModelOnDB) < 1 {
		for _, itemOnInput := range licenseConfigModel.LicenseConfigIDs {
			excludeLicenseConfigID = append(excludeLicenseConfigID, strconv.Itoa(int(itemOnInput.ID.Int64)))
		}
	} else {
		for _, itemOnDB := range licenseConfigModelOnDB {
			_, ok := mapLicenseConfigID[itemOnDB.ID.Int64]
			if !ok {
				excludeLicenseConfigID = append(excludeLicenseConfigID, strconv.Itoa(int(itemOnDB.ID.Int64)))
			}
		}
	}

	if len(excludeLicenseConfigID) > 0 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, fmt.Sprintf(`ID: [%s]`, strings.Join(excludeLicenseConfigID, ", ")))
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) insertLicenseConfig(tx *sql.Tx, licenseConfigModel []repository.LicenseConfigModel) (dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var idReturning []int64
	idReturning, err = dao.LicenseConfigDAO.InsertMultipleExtendedLicenseConfig(tx, licenseConfigModel)
	if err.Error != nil {
		return
	}

	for _, itemIdReturning := range idReturning {
		dataAudit = append(dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.LicenseConfigDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: itemIdReturning},
		})
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input licenseConfigService) validateInsertMultiple(inputStruct *in.LicenseConfigMultipleRequest) errorModel.ErrorModel {
	return inputStruct.ValidateInsertMultipleLicenseConfig()
}

func (input licenseConfigService) readBodyAndValidateMultipleLicenseConfig(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.LicenseConfigMultipleRequest) errorModel.ErrorModel) (inputStruct in.LicenseConfigMultipleRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidateMultipleLicenseConfig"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	if request.Method != "GET" {
		errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
		if errorS != nil {
			err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
			return
		}
	}

	err = validation(&inputStruct)
	return
}
