package ProcessFileListCustomerService

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"io/ioutil"
	"net/http"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/backgroundJobModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"strings"
	"time"
)

func (input importService) ConfirmImport(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var confirm bool
	var pathNew string

	//------ Read request body from JSON
	inputStruct, err := input.ReadConfirmImportBody(request, contextModel, input.validateConfirm)
	if err.Error != nil {
		return
	}

	//------ Confirm decision if true, then move to process, else, delete the file
	confirm, pathNew, err = input.ConfirmImportData(inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	//------ Read data
	if confirm {
		var jobProcessModel repository.JobProcessModel

		jobProcessModel, err = input.readDataForProcess(inputStruct, pathNew, contextModel)
		if err.Error != nil {
			return
		}

		output.Data.Content = out.ViewConfirmImportJobProcessResponse {
			JobID: jobProcessModel.JobID.String,
		}
	}

	output.Status = out.StatusResponse {
		Code: 		util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: 	GenerateI18NMessage("SUCCESS_CONFIRM_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input importService) validateConfirm(inputStruct *in.ImportConfirmRequest) errorModel.ErrorModel {
	return inputStruct.ValidateConfirm()
}

func (input importService) readDataForProcess(inputStruct in.ImportConfirmRequest, path string, contextModel *applicationModel.ContextModel) (jobProcessModel repository.JobProcessModel, err errorModel.ErrorModel) {
	fileName := "ConfirmImport.go"
	funcName := "readDataForProcess"

	//------ If type data import validator is nil, error (basic validation must exist)
	if input.AvailableTypeImport[inputStruct.TypeData].ImportValidator == nil {
		err = errorModel.GenerateUnsupportedRequestParam(fileName, funcName)
		return
	}

	//------ Read file path
	byteData, errorS := ioutil.ReadFile(path)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	return input.doImport(inputStruct, byteData, contextModel)
}

func (input importService) doImport(inputStruct in.ImportConfirmRequest, byteData []byte, contextModel *applicationModel.ContextModel) (jobProcessModel repository.JobProcessModel, err errorModel.ErrorModel) {
	fileName := "ConfirmImport.go"
	funcName := "doImport"
	timeNow := time.Now()
	var extension string
	var finalFileName string

	jobProcessModel = service.GetJobProcess(backgroundJobModel.ChildTask{
		Group: 	constanta.Import,
		Type: 	constanta.File,
		Name: 	"Import " + inputStruct.TypeData + " Data",
	}, *contextModel, timeNow)

	importJob := input.AvailableTypeImport[inputStruct.TypeData]

	buffer := bytes.NewReader(byteData)
	csvData := csv.NewReader(buffer)
	csvData.Comma = importJob.Delimiter
	csvData.LazyQuotes = true
	records, errorS := csvData.ReadAll()
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	jobProcessModel.Total.Int32 = int32(len(records) - 1)
	err = dao.JobProcessDAO.InsertJobProcess(serverconfig.ServerAttribute.DBConnection, jobProcessModel)
	if err.Error != nil {
		return
	}

	fileNameSplit := strings.Split(inputStruct.Filename, ".")
	if len(fileNameSplit) == 2 {
		extension = fileNameSplit[1]
		finalFileName = timeNow.Format(constanta.DefaultTimeFormatForFile) + "/" + jobProcessModel.JobID.String + "." + extension
	}

	hostname, _ := util2.GenerateHostname()
	url, _ := dao.HostServerDAO.GetHostUrlDbByHostname(serverconfig.ServerAttribute.DBConnection, repository.HostServerModel{HostName: sql.NullString{String: hostname}})

	jobProcessModel.UrlIn.String = url.HostURL.String
	jobProcessModel.FilenameIn.String = finalFileName

	go input.DoImportDataToDB(jobProcessModel, records, importJob.ConvertToDTO,
		input.ServiceWithDataAuditGetByAuditService, importJob.DoInsert, contextModel,
		inputStruct, importJob.Truncate)

	err = errorModel.GenerateNonErrorModel()
	return
}