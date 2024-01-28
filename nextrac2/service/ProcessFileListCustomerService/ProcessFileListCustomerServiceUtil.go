package ProcessFileListCustomerService

import (
	"encoding/csv"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/util"
	"os"
	"strconv"
	"strings"
	"time"
)

func (input processFileListCustomerService) readRequestMultipartForm(request *http.Request, contextModel *applicationModel.ContextModel) (fileInfo os.FileInfo, isTruncate bool, err errorModel.ErrorModel) {
	fileName := "ProcessFileListCustomerServiceUtil.go"
	funcName := "readRequestMultipartForm"

	var errorS error
	var file multipart.File
	var handler *multipart.FileHeader

	//------- Maximum 10 MB
	errorS = request.ParseMultipartForm(constanta.SizeMaximumRequireFileImport)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	//------- Get request from multipart form
	file, handler, errorS = request.FormFile(config.ApplicationConfiguration.GetDataDirectory().KeyFile)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	statusTruncate := request.MultipartForm.Value["truncate"][0]
	if statusTruncate != "" {
		isTruncate, errorS = strconv.ParseBool(statusTruncate)
		if errorS != nil {
			err = errorModel.GenerateInvalidJSONRequestError(fileName, funcName, errorS)
			return
		}
	}

	defer func() {
		_ = file.Close()
	}()

	//------- Do validate basic before inbound
	var fileNameCustomerList string
	fileNameCustomerList, err = input.validateBasicBeforeInboundCustomerList(handler, contextModel)
	if err.Error != nil {
		return
	}

	//------- Create dir if not exist
	input.makeNewFolderIsNotExist(varPath().pathInbound)

	//------- Upload file to inbound process 2
	fileBytes, errorS := ioutil.ReadAll(file)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	//------- Name of file import
	fileNameImport := fileNameCustomerList + ".csv"

	//------- Upload file to inbound process 3
	errorS = ioutil.WriteFile(varPath().pathInbound + fileNameImport, fileBytes, 0660)
	//_, errorS = tempFile.Write(fileBytes)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	//------- Get info file (date, size) in inbound
	fileInfo, _ = os.Stat(varPath().pathInbound + fileNameImport)

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input processFileListCustomerService) validateBasicBeforeInboundCustomerList(handler *multipart.FileHeader, contextModel *applicationModel.ContextModel) (fileCustomerListName string, err errorModel.ErrorModel) {

	//------- Get info from handler
	nameImportFile := handler.Filename
	sizeImportFile := handler.Size
	typeImportFile := handler.Header.Get(constanta.GetTypeImportFile)

	//------- Copy to struct
	fileInfo := in.FileInfo{
		FileName: nameImportFile,
		SizeFile: sizeImportFile,
		TypeFile: typeImportFile,
	}

	//------- Do validate basic
	fileCustomerListName, err = input.doValidateBasicCustomerList(&fileInfo, contextModel)
	if err.Error != nil {
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input processFileListCustomerService) doValidateBasicCustomerList(inputStruct *in.FileInfo, contextModel *applicationModel.ContextModel) (fileCustomerListName string, err errorModel.ErrorModel) {
	funcName := "doValidateBasicCustomerList"
	var tempFileName []string
	var message string

	//------- Check extension
	if strings.Contains(inputStruct.FileName, ".") {
		tempFileName = strings.Split(inputStruct.FileName, ".")
		fileCustomerListName = tempFileName[0]
	} else {
		message = input.generateI18NMessage("FAILED_UNKNOWN_EXTENSION_MESSAGE", contextModel.AuthAccessTokenModel.Locale)
		err = errorModel.GenerateInvalidJSONRequestError(input.FileName, funcName, errors.New(message))
		return
	}

	//------- Check type file, must txt or csv
	if inputStruct.TypeFile != constanta.DefineTypeImportFile {
		message = input.generateI18NMessage("FAILED_FILE_REQUIRE_MESSAGE", contextModel.AuthAccessTokenModel.Locale)
		err = errorModel.GenerateInvalidJSONRequestError(input.FileName, funcName, errors.New(message))
		return
	}

	//------- Check data corrupt
	if inputStruct.SizeFile == 0 {
		message = input.generateI18NMessage("FAILED_CORRUPT_MESSAGE", contextModel.AuthAccessTokenModel.Locale)
		err = errorModel.GenerateInvalidJSONRequestError(input.FileName, funcName, errors.New(message))
		return
	}

	//------- Check title name file
	if fileCustomerListName != "trac_list_customer" {
		message = input.generateI18NMessage("FAILED_UNKNOWN_TITLE_MESSAGE", contextModel.AuthAccessTokenModel.Locale)
		err = errorModel.GenerateInvalidJSONRequestError(input.FileName, funcName, errors.New(message))
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input processFileListCustomerService) readFileCustomerList(fileNameCustomerList string) (contentFile [][]string, err errorModel.ErrorModel) {
	funcName := "readFileCustomerList"

	//------- Open file csv
	csvFile, errorS := os.Open(varPath().pathInbound + fileNameCustomerList)
	if errorS != nil {
		err = errorModel.GenerateInvalidJSONRequestError(input.FileName, funcName, errorS)
		return
	}

	defer func() {
		_ = csvFile.Close()
	}()

	//------- Read file csv
	contentFile, errorS = csv.NewReader(csvFile).ReadAll()
	if errorS != nil {
		err = errorModel.GenerateInvalidJSONRequestError(input.FileName, funcName, errorS)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input processFileListCustomerService) dateConvert(date string) (output time.Time) {
	const (
		layoutISO = "2006-01-02"
		layoutUS  = "2006-01-02T15:04:05Z"
	)

	t, _ := time.Parse(layoutISO, date)
	timeString := t.Format(layoutUS)
	output, _ = time.Parse(layoutUS, timeString)
	return
}

func (input processFileListCustomerService) generateI18NMessage(messageID string, language string) (output string) {
	return util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.ImportFileCustomerServiceBundle, messageID, language, nil)
}

func (input processFileListCustomerService) checkFolderIsExist(folderPath string) (fileInfo os.FileInfo , isExist bool) {
	var errorS error

	fileInfo, errorS = os.Stat(folderPath)
	if os.IsNotExist(errorS) {
		return nil, false
	}

	return fileInfo, true
}

func (input processFileListCustomerService) makeNewFolderIsNotExist(pathFolder string) {
	_ = os.MkdirAll(pathFolder, 0770)
}