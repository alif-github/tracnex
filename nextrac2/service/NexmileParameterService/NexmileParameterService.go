package NexmileParameterService

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type nexmileParameterService struct {
	service.AbstractService
}

var NexmileParameterService = nexmileParameterService{}.New()

func (input nexmileParameterService) New() (output nexmileParameterService) {
	output.FileName = "NexmileParameterService.go"
	output.ServiceName = "NEXMILE_PARAMETER"
	return
}

func (input nexmileParameterService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.NexmileParameterRequest) errorModel.ErrorModel) (inputStruct in.NexmileParameterRequest, err errorModel.ErrorModel) {
	var (
		funcName   = "readBodyAndValidate"
		stringBody string
	)

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	err = validation(&inputStruct)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input nexmileParameterService) createModelNexmileParameter(inputStruct in.NexmileParameterRequest, contextModel *applicationModel.ContextModel, timeNow time.Time) (nexmileParameterModel repository.NexmileParameterModel, nexmileParameterModelMap repository.NexmileParameterModelMap) {
	var (
		parameterValueModelCol    []repository.ParameterValueModel
		parameterValueModelColMap = make(map[string]repository.ParameterValueModel)
	)

	for _, itemParameterData := range inputStruct.ParameterData {
		model := repository.ParameterValueModel{
			ParameterID:    sql.NullString{String: itemParameterData.ParameterID},
			ParameterValue: sql.NullString{String: itemParameterData.Value},
		}

		parameterValueModelColMap[itemParameterData.ParameterID] = model
		parameterValueModelCol = append(parameterValueModelCol, model)
	}

	nexmileParameterModelMap = repository.NexmileParameterModelMap{
		CLientID:      sql.NullString{String: inputStruct.ClientID},
		ClientTypeID:  sql.NullInt64{Int64: inputStruct.ClientTypeID},
		UniqueID1:     sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:     sql.NullString{String: inputStruct.UniqueID2},
		ParameterData: parameterValueModelColMap,
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	nexmileParameterModel = repository.NexmileParameterModel{
		ClientID:      sql.NullString{String: inputStruct.ClientID},
		ClientTypeID:  sql.NullInt64{Int64: inputStruct.ClientTypeID},
		UniqueID1:     sql.NullString{String: inputStruct.UniqueID1},
		UniqueID2:     sql.NullString{String: inputStruct.UniqueID2},
		ParameterData: parameterValueModelCol,
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}

	return
}

func (input nexmileParameterService) checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_nexmileparameter_parameterid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.ParameterID)
		}
	}
	return err
}

func (input nexmileParameterService) readBodyAndValidateForView(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.NexmileParameterRequestForView) errorModel.ErrorModel) (inputStruct in.NexmileParameterRequestForView, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidateForView"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	err = validation(&inputStruct)
	return
}
