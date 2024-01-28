package service

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	util2 "nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
	"os"
	"strings"
	"time"
)

type ImportDataService struct {
	AbstractService
	RowStartAt		int
}

type ImportValidatorData struct {
	FieldName	string
	HeaderName	string
	MinLength	int
	MaxLength	int
	Validator	func(string, string, string, string, int, int) (interface{}, errorModel.ErrorModel)
}

type ImportDataType struct {
	ConvertToDTO	func([]string, *map[int]map[string]int64, *applicationModel.ContextModel) (interface{}, errorModel.ErrorModel)
	DoInsert		func(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel)
	Truncate		func(tx *sql.Tx, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel)
	Delimiter		rune
	Message			string
	IsActive		bool
	ImportValidator	[]ImportValidatorData
}

type fileNamePath struct {
	pathInbound 		string
	pathDone			string
	pathFailed			string
	pathProcessNextrac	string
}

func varFilePath() fileNamePath {
	//------ Folder import
	importPath := config.ApplicationConfiguration.GetDataDirectory().ImportPath

	//------ Folder customer
	customerPath := config.ApplicationConfiguration.GetDataDirectory().CustomerPath

	//------ Full path
	path := config.ApplicationConfiguration.GetDataDirectory().BaseDirectoryPath + importPath + customerPath

	return fileNamePath {
		pathInbound:	 	path + config.ApplicationConfiguration.GetDataDirectory().InboundPath + "/",
		pathDone: 			path + config.ApplicationConfiguration.GetDataDirectory().DonePath + "/",
		pathFailed: 		path + config.ApplicationConfiguration.GetDataDirectory().FailedPath + "/",
		pathProcessNextrac: path + config.ApplicationConfiguration.GetDataDirectory().ProcessPath + "/",
	}
}

func (input ImportDataService) ReadConfirmImportBody(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.ImportConfirmRequest) errorModel.ErrorModel) (inputStruct in.ImportConfirmRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	_ = json.Unmarshal([]byte(stringBody), &inputStruct)
	err = validation(&inputStruct)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input ImportDataService) ReadRequestMultipartForm(request *http.Request, contextModel *applicationModel.ContextModel,
	validation func(input *in.ImportRequest) errorModel.ErrorModel) (inputStruct in.ImportRequest, buffer *bytes.Buffer, extension string, err errorModel.ErrorModel) {

	fileName := "ImportDataService.go"
	funcName := "readRequestMultipartForm"

	var errorS error
	var file multipart.File
	var handler *multipart.FileHeader

	//------- Maximum 10 MB
	errorS = request.ParseMultipartForm(constanta.SizeMaximumRequireFileImport)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	//------- Get request from multipart form
	file, handler, errorS = request.FormFile(config.ApplicationConfiguration.GetDataDirectory().KeyFile)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	//------- Get request content
	content := request.MultipartForm.Value[config.ApplicationConfiguration.GetDataDirectory().KeyContent][0]
	byteContent := []byte(content)
	_ = json.Unmarshal(byteContent, &inputStruct)

	//------- Validation input request content
	err = validation(&inputStruct)
	if err.Error != nil {
		return
	}

	defer func() {
		_ = file.Close()
	}()

	//------- Write to contextModel
	contextModel.LoggerModel.ByteIn = int(handler.Size)

	//------- Get extension with split
	extensionTemp := strings.Split(handler.Filename, ".")
	if len(extensionTemp) == 1 {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, "file")
		return
	}

	//------- Error if extension is not "csv"
	if extensionTemp[len(extensionTemp) - 1] != "csv" {
		err = errorModel.GenerateFormatFieldError(fileName, funcName, "file")
		return
	}

	extension = extensionTemp[len(extensionTemp) - 1]

	//------- Buffer read
	buffer = bytes.NewBuffer(nil)
	if _, errorS = io.Copy(buffer, file); errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input ImportDataService) ValidateImportData(fileName string, funcName string, buff *bytes.Buffer, extension string,
	delimiter rune, validator []ImportValidatorData, typeData string, locale string) (fileDataName string, result [][]string, totalData int,
	err errorModel.ErrorModel, multipleError []out.MultipleErrorResponse) {

	var isValid bool

	if buff == nil {
		err = errorModel.GenerateEmptyFieldError(fileName, funcName, "file")
		return
	}

	r := bytes.NewReader(buff.Bytes())
	if strings.ToLower(extension) == "csv" {
		isValid = true
		result, totalData, err, multipleError = input.readCSVData(fileName, funcName, r, delimiter, validator, locale, typeData)
	}

	if err.Error != nil || len(multipleError) != 0 {
		return
	}

	if isValid {
		fileDataName, err = input.WriteFileImportData(buff.Bytes(), typeData, extension)
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input ImportDataService) readCSVData(fileName string, funcName string, reader *bytes.Reader, delimiter rune,
	validator []ImportValidatorData, locale string, typeData string) (result [][]string, totalData int, err errorModel.ErrorModel, multipleError []out.MultipleErrorResponse) {

	//------ Preparing csv
	csvData := csv.NewReader(reader)
	csvData.Comma = delimiter
	csvData.LazyQuotes = true
	counter := 0
	indexData := 0

	//------ Loop
	for {
		//------ Read csv
		record, errorS := csvData.Read()
		if errorS == io.EOF {
			break
		}

		//------ Read header
		if counter == input.RowStartAt - 1 {
			for i := 0; i < len(validator); i++ {

				if i == 10 && typeData == "customer" {
					for j := 10; j <= 17; j++ {
						if !strings.Contains(record[j], validator[i].HeaderName) {
							err = errorModel.GenerateFormatFieldError(fileName, funcName, "Header : " + record[j])
							return
						}
					}
				} else if record[i] != validator[i].HeaderName {
					err = errorModel.GenerateFormatFieldError(fileName, funcName, "Header : " + record[i])
					return
				}

			}
		}

		//------ Header is not read
		if counter < input.RowStartAt {
			counter++
			continue
		}

		if errorS != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
			return
		}

		indexData++
		//------ Validate row data
		resultTemp, err := input.validateRowData(fileName, funcName, record, validator, typeData)
		if err.Error != nil {

			//------ Add if any multiple error
			//todo -----multiple error import
			if len(multipleError) < 10 {
				multipleError = append(multipleError, out.MultipleErrorResponse{
					ID: 		int64(indexData),
					CausedBy: 	util.GenerateI18NErrorMessage(err, locale),
				})

				err = errorModel.GenerateNonErrorModel()
			} else if len(multipleError) == 10 {
				break
			}

			continue
		}

		//---------- Limit 10 for sample
		if len(result) < 10 {
			result = append(result, resultTemp[0])
		}
	}

	totalData = indexData
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input ImportDataService) validateRowData(fileName string, funcName string, record []string, validator []ImportValidatorData,
	typeData string) (result [][]string, err errorModel.ErrorModel) {

	for i := 0; i < len(validator); i++ {
		//------- Special case
		var expDataFillColumn int
		if i == 10 && typeData == "customer" {
			for j := 17; j >= 10; j-- {
				if record[j] != "" {
					expDataFillColumn = j
					break
				}
			}
			_, err = validator[i].Validator(fileName, funcName, record[expDataFillColumn], validator[i].FieldName, validator[i].MinLength, validator[i].MaxLength)
			if err.Error != nil {
				return
			}
		} else {
			//------- Normal case
			_, err = validator[i].Validator(fileName, funcName, record[i], validator[i].FieldName, validator[i].MinLength, validator[i].MaxLength)
			if err.Error != nil {
				return
			}
		}
	}

	result = append(result, record)
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input ImportDataService) WriteFileImportData(file []byte, typeData string, extension string) (fileDataName string, err errorModel.ErrorModel) {
	fileName := "ImportDataService.go"
	funcName := "WriteFileImportData"

	//------ Folder path inbound
	inboundPath := varFilePath().pathInbound

	//------ Generate UUID
	uuid := util2.GetUUID()

	//------ FileName preparation
	fileDataName = typeData + "_" + uuid + "." + extension

	//------ Create folder
	_ = os.MkdirAll(inboundPath, 0770)

	//------ Write file
	errorS := ioutil.WriteFile(inboundPath + fileDataName, file, 0660)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	return fileDataName, errorModel.GenerateNonErrorModel()
}

func (input ImportDataService) ConfirmImportData(inputStruct in.ImportConfirmRequest, contextModel *applicationModel.ContextModel) (confirm bool, pathNew string, err errorModel.ErrorModel) {
	if inputStruct.Confirm {
		pathNew, err = input.MoveFileImportData(inputStruct.Filename, "", contextModel, false)
		confirm = true
		if err.Error != nil {
			return
		}
	} else {
		err = input.DeleteFileImportData(varFilePath().pathInbound, inputStruct.Filename)
		confirm = false
		if err.Error != nil {
			return
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input ImportDataService) DeleteFileImportData(path string, fileDataName string) (err errorModel.ErrorModel) {
	fileName := "ImportDataService.go"
	funcName := "DeleteFileImportData"
	pathFileName := path + fileDataName

	errorS := os.Remove(pathFileName)
	if errorS != nil {
		return errorModel.GenerateUnknownError(fileName, funcName, errorS)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input ImportDataService) MoveFileImportData(fileDataName string, finalFileName string, contextModel *applicationModel.ContextModel, isMoveToDone bool) (path string, err errorModel.ErrorModel) {
	fileName := "ImportDataService.go"
	funcName := "MoveFileImportData"
	var logModel applicationModel.LoggerModel

	//------ Path preparation
	inboundPathTemp := varFilePath().pathInbound
	processPathTemp := varFilePath().pathProcessNextrac
	donePathTemp := varFilePath().pathDone

	if isMoveToDone {
		//------ File moving to done
		var fileDataNameDone string
		if strings.Contains(fileDataName, ".csv") {
			fileDataName = strings.ReplaceAll(fileDataName, ".csv", ".process")
		}
		processPath := processPathTemp + fileDataName
		if finalFileName != "" {
			fileDataNameDone = finalFileName
		} else {
			fileDataNameDone = strings.ReplaceAll(fileDataName, ".process", ".csv")
		}

		arrPath := strings.Split(finalFileName, "/")

		var donePath string

		if len(arrPath) > 0 {
			year := arrPath[0]
			month := arrPath[1]
			date := arrPath[2]
			file := arrPath[3]

			_ = os.MkdirAll(donePathTemp + "/" + year + "/" + month + "/" + date, 0770)
			donePath = donePathTemp + year + "/" + month + "/" + date + "/" + file
		} else {
			_ = os.MkdirAll(donePathTemp, 0770)
			donePath = donePathTemp + fileDataNameDone
		}

		errorS := os.Rename(processPath, donePath)
		if errorS != nil {
			logModel = contextModel.LoggerModel
			logModel.Status = 500
			logModel.Message = errorS.Error()
			return "", errorModel.GenerateUnknownError(fileName, funcName, errorS)
		}
		path = donePath
	} else {
		//------ File moving to process
		inboundPath := inboundPathTemp + fileDataName
		fileDataNameProcess := strings.ReplaceAll(fileDataName, ".csv", ".process")
		_ = os.MkdirAll(processPathTemp, 0770)
		processPath := processPathTemp + fileDataNameProcess
		errorS := os.Rename(inboundPath, processPath)
		if errorS != nil {
			logModel = contextModel.LoggerModel
			logModel.Status = 500
			logModel.Message = errorS.Error()
			return "", errorModel.GenerateUnknownError(fileName, funcName, errorS)
		}
		path = processPath
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input ImportDataService) DoImportDataToDB(jobProcessModel repository.JobProcessModel, data [][]string,
	convertToDTO func([]string, *map[int]map[string]int64, *applicationModel.ContextModel) (interface{}, errorModel.ErrorModel),
	wrapAudit func(action int32, funcName string, inputStruct interface{}, contextModel *applicationModel.ContextModel, serve func(*sql.Tx, interface{}, *applicationModel.ContextModel, time.Time) (interface{}, []repository.AuditSystemModel, errorModel.ErrorModel), additionalAfterCommit func(interface{}, applicationModel.ContextModel)) (output interface{}, err errorModel.ErrorModel),
	doInsert func(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel),
	contextModel *applicationModel.ContextModel, inputStruct in.ImportConfirmRequest,
	truncate func(tx *sql.Tx, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel)) () {

	fileName := "ImportDataService.go"
	funcName := "DoImportDataToDB"
	var output []out.MultipleErrorResponse
	var err errorModel.ErrorModel
	counter := 0

	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())

	defer func() {
		if r := recover(); r != nil {
			util.InputLog(errorModel.GenerateRecoverError(), contextModel.LoggerModel)
			//todo update error
		}
	}()

	var listID map[int]map[string]int64
	listID = make(map[int]map[string]int64)
	if inputStruct.Truncate && inputStruct.TypeData == "customer" {
		var errorResponse out.MultipleErrorResponse
		var txTruncate *sql.Tx
		var errorS error

		func() {
			defer func() {
				if err.Error != nil {
					_ = txTruncate.Rollback()

					errorResponse = out.MultipleErrorResponse {
						ID: 		0,
						CausedBy: 	util.GenerateI18NErrorMessage(err, contextModel.AuthAccessTokenModel.Locale),
					}

					jobProcessModel.Status.String = constanta.JobProcessErrorStatus
					jobProcessModel.ContentDataOut.String = util2.StructToJSON(errorResponse)
					go input.DeleteFileImportData(varFilePath().pathProcessNextrac, inputStruct.Filename)
					err = dao.JobProcessDAO.UpdateErrorJobProcess(serverconfig.ServerAttribute.DBConnection, jobProcessModel)
					if err.Error != nil {
						logModel.Message = err.Error.Error()
						logModel.Status = err.Code
						if err.CausedBy != nil {
							logModel.Message = err.CausedBy.Error()
							logModel.Status = 500
						}
						util2.LogError(logModel.ToLoggerObject())
					}
				} else {
					_ = txTruncate.Commit()
				}
			}()

			txTruncate, errorS = serverconfig.ServerAttribute.DBConnection.Begin()
			if errorS != nil {
				err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
				return
			}
			err = truncate(txTruncate, contextModel)
			if err.Error != nil {
				return
			}

			err = errorModel.GenerateNonErrorModel()
			return
		}()
	}

	for i := 1; i < len(data); i++ {
		func() {
			defer func() {
				if err.Error != nil {
					output = append(output, out.MultipleErrorResponse{
						ID:       int64(i),
						CausedBy: util.GenerateI18NErrorMessage(err, contextModel.AuthAccessTokenModel.Locale),
					})
				}
			}()
			var dtoIn interface{}
			dtoIn, err = convertToDTO(data[i], &listID, contextModel)
			if err.Error != nil {
				return
			}

			_, err = wrapAudit(int32(constanta.ActionAuditInsertConstanta), funcName, dtoIn, contextModel, doInsert, func(_ interface{}, _ applicationModel.ContextModel) {})
			if err.Error != nil {
				return
			}
		}()

		updateDBEvery := CountUpdateDBJobProcessCounter(jobProcessModel.Total.Int32)
		if counter == updateDBEvery || i == len(data) - 1 {

			//-- todo must create duration
			jobProcessModel.UpdatedAt.Time = time.Now()

			if len(output) > 0 {
				jobProcessModel.Status.String = constanta.JobProcessOnProgressErrorStatus
			}

			if len(output) > 0 {
				jobProcessModel.ContentDataOut.String = util2.StructToJSON(output)
			}

			if i == len(data) - 1 {
				jobProcessModel.Counter.Int32 += int32(counter + 1)
				path, _ := input.MoveFileImportData(inputStruct.Filename, jobProcessModel.FilenameIn.String, contextModel, true)
				jobProcessModel.FilenameIn.String = path
				err = dao.JobProcessDAO.UpdateJobProcessCounter(serverconfig.ServerAttribute.DBConnection, jobProcessModel)
				if err.Error != nil {
					logModel.Message = err.Error.Error()
					logModel.Status = err.Code
					if err.CausedBy != nil {
						logModel.Message = err.CausedBy.Error()
						logModel.Status = 500
					}
					util2.LogError(logModel.ToLoggerObject())
				}
			} else {
				jobProcessModel.Counter.Int32 += int32(counter)
			}

			err = dao.JobProcessDAO.UpdateJobProcessCounter(serverconfig.ServerAttribute.DBConnection, jobProcessModel)
			if err.Error != nil {
				logModel.Message = err.Error.Error()
				logModel.Status = err.Code
				if err.CausedBy != nil {
					logModel.Message = err.CausedBy.Error()
					logModel.Status = 500
				}
				util2.LogError(logModel.ToLoggerObject())
			}
			counter = 0
		}
		counter++
	}
}