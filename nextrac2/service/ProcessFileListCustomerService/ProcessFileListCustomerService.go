package ProcessFileListCustomerService

import (
	"database/sql"
	"errors"
	"io/ioutil"
	"net/http"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"os"
	"strings"
	"time"
)

type processFileListCustomerService struct {
	in.AbstractDTO
	service.AbstractService
}

var ProcessFileListCustomerService = processFileListCustomerService{}.New()

func (input processFileListCustomerService) New() (output processFileListCustomerService) {
	output.FileName = "ProcessFileListCustomerService.go"
	return
}

type fileCustomer struct {
	pathInbound 		string
	pathDone			string
	pathFailed			string
	pathProcessNextrac	string
}

func varPath() fileCustomer {
	//------ Folder import
	importPath := config.ApplicationConfiguration.GetDataDirectory().ImportPath

	//------ Folder customer
	customerPath := config.ApplicationConfiguration.GetDataDirectory().CustomerPath

	//------ Full path
	path := config.ApplicationConfiguration.GetDataDirectory().BaseDirectoryPath + importPath + customerPath

	return fileCustomer {
		pathInbound:	 	path + config.ApplicationConfiguration.GetDataDirectory().InboundPath + "/",
		pathDone: 			path + config.ApplicationConfiguration.GetDataDirectory().DonePath + "/",
		pathFailed: 		path + config.ApplicationConfiguration.GetDataDirectory().FailedPath + "/",
		pathProcessNextrac: path + config.ApplicationConfiguration.GetDataDirectory().ProcessPath + "/",
	}
}

func (input processFileListCustomerService) ImportFileCustomerList(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	//------- Time begin here
	timeNow := time.Now()

	//------- Read body and validate request file
	fileInfo, isTruncate, err := input.readRequestMultipartForm(request, contextModel)
	if err.Error != nil {
		return
	}

	//------- Validate 2 in inbound and generate file
	var fileNameCustomerList string
	output.Data.Content, fileNameCustomerList, err = input.DoGenerateFileCustomerList(fileInfo, contextModel, timeNow, isTruncate)
	if err.Error != nil {
		return
	}

	//------- Move file from process to done
	input.makeNewFolderIsNotExist(varPath().pathDone)
	folderInProcess := strings.ReplaceAll(varPath().pathProcessNextrac + fileNameCustomerList, ".csv", ".process")
	go func() {
		_ = os.Rename(folderInProcess, strings.ReplaceAll(varPath().pathDone + fileNameCustomerList, ".process", ".csv"))
	}()

	output.Status = out.StatusResponse{
		Code: 		util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message:	input.generateI18NMessage("SUCCESS_IMPORT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input processFileListCustomerService) DoGenerateFileCustomerList(fileInfo os.FileInfo, contextModel *applicationModel.ContextModel,
	timeNow time.Time, isTruncate bool) (output interface{}, fileNameCustomerList string, err errorModel.ErrorModel) {

	funcName := "DoGenerateFileCustomerList"
	var contentFile [][]string
	var errorS error
	var tx *sql.Tx

	defer func() {
		if errorS != nil || err.Error != nil {
			errorS = tx.Rollback()
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
			}
		} else {
			errorS = tx.Commit()
			if errorS != nil {
				err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
			}
		}
	}()

	tx, errorS = serverconfig.ServerAttribute.DBConnection.Begin()
	if errorS != nil {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errorS)
		return
	}

	//------- Validate file customer list in folder inbound
	_, fileNameCustomerList, err = input.validateFileCustomerList(fileInfo, contextModel)
	if err.Error != nil {
		return
	}

	//------- Read file customer list in folder inbound
	contentFile, err = input.readFileCustomerList(fileNameCustomerList)
	if err.Error != nil {
		return
	}

	//------- Move file to process
	folderInInbound := varPath().pathInbound + fileNameCustomerList
	input.makeNewFolderIsNotExist(varPath().pathProcessNextrac)
	go func() {
		_ = os.Rename(folderInInbound, strings.ReplaceAll(varPath().pathProcessNextrac + fileNameCustomerList, ".csv", ".process"))
	}()

	//------- Truncate the list table and reset sequence
	if isTruncate {
		err = input.truncateTableCustomerList(tx)
		if err.Error != nil {
			return
		}
	}

	//------- Save list customer in DB
	var failedIndex int
	var importFileStruct []out.ImportFile
	var statFileStruct	out.StatFile
	totalRow := 0
	totalSuccessRow := 0
	totalFailedRow := 0

	//------- Count the total row
	go func() {
		for indexTemp := range contentFile {
			if strings.Contains(contentFile[indexTemp][0], "|ND6") || strings.Contains(contentFile[indexTemp][0], "|NF6") {
				totalRow++
			}
		}
	}()

	for index := range contentFile {
		if index < 1 {continue}

		if strings.Contains(contentFile[index][0], "|ND6") || strings.Contains(contentFile[index][0], "|NF6") {
			failedIndex, err = input.SaveListCustomerInDB(contentFile[index], index, contextModel, timeNow, tx)
			if failedIndex > 0 && err.Error != nil {
				importFileStruct = append(importFileStruct, out.ImportFile{
					NumberData: 	failedIndex,
					ErrorMessage: 	err,
				})
				totalFailedRow++
			} else {
				totalSuccessRow++
			}
		} else {
			continue
		}
	}

	statFileStruct = out.StatFile{
		FileSuccess: 	totalSuccessRow,
		FileFailed: 	totalFailedRow,
		FileAmount: 	totalRow,
		FileDetail: 	importFileStruct,
	}

	output = statFileStruct
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input processFileListCustomerService) truncateTableCustomerList(tx *sql.Tx) (err errorModel.ErrorModel) {
	funcName := "truncateTableCustomerList"

	//------- Step 1 truncate table customer
	err = dao.CustomerListDAO.TruncateTableCustomer(tx)
	if err.Error != nil {
		return
	}

	//------- Step 2 count row
	var counts int64
	counts, err = dao.CustomerListDAO.CountRowCustomerList(tx)
	if err.Error != nil {
		return
	}

	//------- Step 3 reset sequence ID
	if counts == 0 {
		err = dao.CustomerListDAO.SetValSequenceCustomerListTable(tx)
		if err.Error != nil {
			return
		}
	} else {
		err = errorModel.GenerateInternalDBServerError(input.FileName, funcName, errors.New("counts is not zero"))
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input processFileListCustomerService) validateFileCustomerList(_ os.FileInfo, contextModel *applicationModel.ContextModel) (fileNameCustomerList string, fileNameWithExtension string, err errorModel.ErrorModel) {
	funcName := "validateFileCustomerList"
	var message string
	countFileProcess := 0

	//------- Read directory inbound
	files, errorS := ioutil.ReadDir(varPath().pathInbound)
	if errorS != nil {
		err = errorModel.GenerateUnknownError(input.FileName, funcName, errorS)
		return
	}

	//------- Get and validate current file
	for _, fileItem := range files {
		//------- Get file which want to process
		//if !fileItem.ModTime().Equal(fileInfo.ModTime()) {
		//	go input.moveFileToFailed(fileItem.Name(), fileItem.Name())
		//	continue
		//}

		//------- Check extension
		var tempFileName []string
		if strings.Contains(fileItem.Name(), ".") {
			tempFileName = strings.Split(fileItem.Name(), ".")
		} else {
			go input.moveFileToFailed(fileItem.Name(), fileItem.Name())
			continue
		}

		fileNameCustomerList = tempFileName[0]

		var resourceTargetFile string
		var listName string
		var object string

		//------- Check name file
		if strings.Contains(fileNameCustomerList, "_") {
			temp := strings.Split(fileNameCustomerList, "_")
			resourceTargetFile = temp[0]
			listName = temp[1]
			object = temp[2]
		} else {
			go input.moveFileToFailed(fileItem.Name(), fileItem.Name())
			continue
		}

		//------- Check name file level 1
		if resourceTargetFile != config.ApplicationConfiguration.GetServerResourceID() {
			go input.moveFileToFailed(fileItem.Name(), fileItem.Name())
			continue
		}

		//------- Check name file level 2
		if listName != constanta.NameFileLevel2 {
			go input.moveFileToFailed(fileItem.Name(), fileItem.Name())
			continue
		}

		//------- Check name file level 3
		if object != constanta.NameFileLevel3 {
			go input.moveFileToFailed(fileItem.Name(), fileItem.Name())
			continue
		}

		countFileProcess++
		fileNameWithExtension = fileItem.Name()
	}

	//------- File only 1 who have processing
	if countFileProcess != 1 {
		message = input.generateI18NMessage("FAILED_DATA_PROCESS_MESSAGE", contextModel.AuthAccessTokenModel.Locale)
		err = errorModel.GenerateInvalidJSONRequestError(input.FileName, funcName, errors.New(message))
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input processFileListCustomerService) moveFileToFailed(oldFileName string, newFileName string) {

	//------- Make Mkdir for path failed
	input.makeNewFolderIsNotExist(varPath().pathFailed)

	//------- Move file to failed
	_ = os.Rename(varPath().pathInbound + oldFileName, varPath().pathFailed + newFileName)
}