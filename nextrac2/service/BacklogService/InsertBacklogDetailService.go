package BacklogService

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"os"
	"time"
)

func (input backlogService) InsertDetailBacklog(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		inputStruct []*in.BacklogRequest
		funcName    = "InsertDetailBacklog"
	)

	inputStruct, _, _, err = input.ReadRequestMultipartForm(request, contextModel, input.validateInsertDetailBacklog)
	if err.Error != nil {
		return
	}

	_, err = input.InsertServiceWithAudit(funcName, inputStruct, contextModel, input.doInsertDetailBacklog, nil)
	if err.Error != nil {
		return
	}

	output.Status = input.GetResponseMessage("SUCCESS_INSERT_MESSAGE", contextModel)
	return
}

func (input backlogService) doInsertDetailBacklog(tx *sql.Tx, inputStructInterface interface{}, contextModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	var (
		funcName     = "doInsertDetailBacklog"
		inputStruct  = inputStructInterface.([]*in.BacklogRequest)
		listBacklog  []in.BacklogRequest
		insertedID   int64
		employeeOnDB repository.EmployeeModel
		scope        map[string]interface{}
		db           = serverconfig.ServerAttribute.DBConnection
	)

	mappingScopeDB := make(map[string]applicationModel.MappingScopeDB)
	mappingScopeDB[constanta.EmployeeDataScope] = applicationModel.MappingScopeDB{
		View:  "e.id",
		Count: "e.id",
	}

	createdBy := contextModel.LimitedByCreatedBy       //--- Add userID when have own permission
	scope, err = input.validateDataScope(contextModel) //--- Get scope
	if err.Error != nil {
		return
	}

	// handle Form Perubahan
	_, err = input.HandleFile(tx, inputStruct, contextModel)
	if err.Error != nil {
		return
	}

	// cara bodoh untuk menghilangkan pointer
	for _, itemBacklog := range inputStruct {
		listBacklog = append(listBacklog, *itemBacklog)
	}

	inputModel := input.convertDTOToModel(listBacklog, *contextModel, timeNow)
	for _, inputModelItem := range inputModel {
		employeeOnDB, err = dao.EmployeeDAO.ViewEmployee(db, repository.EmployeeModel{ID: inputModelItem.EmployeeId}, createdBy, scope, mappingScopeDB)
		if err.Error != nil {
			return
		}

		if employeeOnDB.ID.Int64 < 1 {
			err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.EmployeeID)
			return
		}

		// validasi form developer sesuai department name pic
		if inputModelItem.DepartmentCode.String == "developer" {
			if employeeOnDB.DepartmentId.Int64 != 1 {
				err = errorModel.GenerateSimpleErrorModel(400, "Form tidak sesuai dengan department pic")
				return
			}
		}

		// validasi form qa sesuai department name pic
		if inputModelItem.DepartmentCode.String == "qaqc" {
			if employeeOnDB.DepartmentId.Int64 != 2 {
				err = errorModel.GenerateSimpleErrorModel(400, "Form tidak sesuai dengan department pic")
				return
			}
		}

		insertedID, err = dao.BacklogDAO.InsertBacklog(tx, inputModelItem)
		if err.Error != nil {
			err = input.checkDuplicateError(err)
			return
		}

		dataAudit = append(dataAudit, repository.AuditSystemModel{
			TableName:  sql.NullString{String: dao.EmployeeDAO.TableName},
			PrimaryKey: sql.NullInt64{Int64: insertedID},
		})
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input backlogService) HandleFile(tx *sql.Tx, listBacklog []*in.BacklogRequest, contextModel *applicationModel.ContextModel) (newUpdatedFileUpload int64, err errorModel.ErrorModel) {
	for _, backlog := range listBacklog {
		if !util.IsStringEmpty(backlog.File.Base64) {
			if util.IsStringEmpty(backlog.File.FileName) {
				err = errorModel.GenerateSimpleErrorModel(400, "data file tidak lengkap")
				return
			}

			if util.IsStringEmpty(backlog.File.Type) {
				err = errorModel.GenerateSimpleErrorModel(400, "data file tidak lengkap")
				return
			}

			var (
				timeUnixStr    = fmt.Sprintf(`%d`, time.Now().Unix())
				fileName       = backlog.File.FileName + timeUnixStr + backlog.File.Type
				fileUploadRepo repository.FileUpload
				fileByte       []byte
				errs           error
				fileSize       int64
				idInserted     int64
			)

			// convert base64 to []byte
			fileByte, errs = base64.StdEncoding.DecodeString(backlog.File.Base64)
			if errs != nil {
				err = errorModel.GenerateSimpleErrorModel(400, "error decode base64")
				return
			}

			// Local Upload
			fileUploadRepo, err = uploadLocalDirectory(fileName, fileByte, contextModel)
			if err.Error != nil {
				return
			}

			// CDN Upload
			fileRequestCDN := []in.MultipartFileDTO{
				{
					FileContent: fileByte,
					Filename:    fileName,
					Size:        fileSize,
					Host:        config.ApplicationConfiguration.GetCDN().Host,
					Path:        config.ApplicationConfiguration.GetCDN().RootPath,
					FileID:      contextModel.AuthAccessTokenModel.ResourceUserID,
				},
			}

			fmt.Println("Local File Location :", fileUploadRepo.Host.String+fileUploadRepo.Path.String+fileUploadRepo.FileName.String)
			fmt.Println("Done Upload to Local Server")
			fmt.Println("File Name : ", fileName)

			containerName := constanta.ContainerBacklog + service.GetAzureDateContainer()
			err = service.UploadFileToLocalCDN(containerName, &fileRequestCDN, contextModel.AuthAccessTokenModel.ResourceUserID)
			if err.Error != nil {
				return
			}

			fmt.Println("Success Upload to CDN Server")

			// update item file_upload
			fileUploadRepo.Host.String = config.ApplicationConfiguration.GetCDN().RootPath
			fileUploadRepo.Path.String = config.ApplicationConfiguration.GetCDN().Suffix + containerName
			fileUploadRepo.FileName.String = fileRequestCDN[0].Filename

			fmt.Println("Save Table File Upload -> " + fileUploadRepo.Host.String + fileUploadRepo.Path.String + fileName)

			// Azure Upload
			fileRequestAzure := []in.MultipartFileDTO{
				{
					FileContent: fileByte,
					Filename:    fileName,
					Size:        fileSize,
					Host:        config.ApplicationConfiguration.GetAzure().Host,
					Path:        config.ApplicationConfiguration.GetAzure().Suffix,
					FileID:      contextModel.AuthAccessTokenModel.ResourceUserID,
				},
			}

			err = service.UploadFileToAzure(&fileRequestAzure)
			if err.Error != nil {
				fmt.Println("Upload Azure Failed")
			} else {
				// update item file_upload
				fileUploadRepo.Host.String = config.ApplicationConfiguration.GetAzure().Host
				fileUploadRepo.Path.String = config.ApplicationConfiguration.GetAzure().Suffix + service.GetAzureDateContainer()
				fileUploadRepo.FileName.String = fileRequestAzure[0].Filename
			}

			fmt.Println("Success Upload to Azure Server")
			fmt.Println(" Last File Upload ->", fileUploadRepo.Host.String+fileUploadRepo.Path.String+
				service.GetAzureDateContainer()+fileUploadRepo.FileName.String)

			// save to table file_upload
			idInserted, err = dao.FileUploadDAO.InsertFileUploadInfoForBacklog(tx, fileUploadRepo)
			if err.Error != nil {
				return
			}

			backlog.FileUploadId = idInserted
			newUpdatedFileUpload = idInserted
		}
	}

	return
}

func uploadLocalDirectory(fileName string, fileByte []byte, contextModel *applicationModel.ContextModel) (fileUploadLocal repository.FileUpload, err errorModel.ErrorModel) {
	var (
		errs error
	)

	dataDirectoryLocal := config.ApplicationConfiguration.GetDataDirectory().BaseDirectoryPath +
		config.ApplicationConfiguration.GetDataDirectory().ImportPath + config.ApplicationConfiguration.GetDataDirectory().Backlog +
		string(os.PathSeparator)

	_ = os.MkdirAll(dataDirectoryLocal, 0770)
	//_ = os.MkdirAll(dataDirectoryLocal, 777)

	errs = ioutil.WriteFile(dataDirectoryLocal+fileName, fileByte, 0660)
	//errs = ioutil.WriteFile(dataDirectoryLocal+fileName, fileByte, 777)
	if errs != nil {
		err = errorModel.GenerateSimpleErrorModel(500, "error upload local server")
		return
	}

	// Read File For Get Size
	//fileInfo, errs := os.Stat(dataDirectoryLocal + fileName)
	//if errs != nil {
	//	err = errorModel.GenerateSimpleErrorModel(500, "error get local file")
	//	return
	//}

	fileUploadLocal = repository.FileUpload{
		FileName: sql.NullString{String: fileName},
		FileSize: sql.NullInt64{Int64: 1},
		Konektor: sql.NullString{String: dao.BacklogDAO.TableName},
		Host:     sql.NullString{String: config.ApplicationConfiguration.GetDataDirectory().BaseDirectoryPath},
		Path: sql.NullString{String: config.ApplicationConfiguration.GetDataDirectory().ImportPath +
			config.ApplicationConfiguration.GetDataDirectory().Backlog + "/"},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		CreatedAt:     sql.NullTime{Time: time.Now()},
		UpdatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		UpdatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: time.Now()},
	}

	return
}

func reqUpload(byteData []byte, fileName string, pathFile string, pathUpload string) (string, string, error) {
	// byteData, _ := ioutil.ReadFile(pathFile + fileName)
	// var uuid = util.GetUUID()

	pathDirectory := config.ApplicationConfiguration.GetCDN().RootPath + pathUpload
	_ = os.MkdirAll(pathDirectory, 0770)
	errs := ioutil.WriteFile(pathDirectory+fileName, []byte(byteData), 0660)
	if errs != nil {
		fmt.Println("error upload ", errs.Error())
		// return
	}
	fmt.Println("upload success")
	fmt.Println("save into local CDN success")

	hostSuffix := config.ApplicationConfiguration.GetCDN().Host
	fmt.Println("Host ===> ", config.ApplicationConfiguration.GetCDN().Host)
	//fmt.Println("Suffix ===> ", config.ApplicationConfiguration.GetCDN().Suffix)

	return pathDirectory, hostSuffix, errs
}

func (input backlogService) convertDTOToModel(listInputStruct []in.BacklogRequest, contextModel applicationModel.ContextModel, timeNow time.Time) (listBacklog []repository.BacklogModel) {
	for _, inputStruct := range listInputStruct {
		backlogItem := repository.BacklogModel{
			ID:              sql.NullInt64{Int64: inputStruct.ID},
			Layer1:          sql.NullString{String: inputStruct.Layer1},
			Layer2:          sql.NullString{String: inputStruct.Layer2},
			Layer3:          sql.NullString{String: inputStruct.Layer3},
			Layer4:          sql.NullString{String: inputStruct.Layer4},
			Layer5:          sql.NullString{String: inputStruct.Layer5},
			RedmineNumber:   sql.NullInt64{Int64: inputStruct.RedmineNumber},
			Sprint:          sql.NullString{String: inputStruct.Sprint},
			SprintName:      sql.NullString{String: inputStruct.SprintName},
			EmployeeId:      sql.NullInt64{Int64: inputStruct.PicId},
			Status:          sql.NullString{String: inputStruct.Status},
			Mandays:         sql.NullFloat64{Float64: inputStruct.Mandays},
			EstimateTime:    sql.NullFloat64{Float64: inputStruct.MandaysDone},
			FlowChanged:     sql.NullString{String: inputStruct.FlowChanged},
			AdditionalData:  sql.NullString{String: inputStruct.AdditionalData},
			Note:            sql.NullString{String: inputStruct.Note},
			Url:             sql.NullString{String: inputStruct.Url},
			Page:            sql.NullString{String: inputStruct.Page},
			Tracker:         sql.NullString{String: inputStruct.Tracker},
			Feature:         sql.NullInt64{Int64: inputStruct.Feature},
			ReferenceTicket: sql.NullInt64{Int64: inputStruct.ReferenceTicket},
			FileUploadId:    sql.NullInt64{Int64: inputStruct.FileUploadId},
			DepartmentCode:  sql.NullString{String: inputStruct.DepartmentCode},
			Subject:         sql.NullString{String: inputStruct.Subject},
			UpdatedBy:       sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			UpdatedAt:       sql.NullTime{Time: timeNow},
			UpdatedClient:   sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
			CreatedBy:       sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
			CreatedAt:       sql.NullTime{Time: timeNow},
			CreatedClient:   sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		}

		listBacklog = append(listBacklog, backlogItem)
	}

	return
}

func (input backlogService) validateInsertDetailBacklog(inputStruct []*in.BacklogRequest) (err errorModel.ErrorModel) {
	for _, item := range inputStruct {
		err = item.ValidateInsert()
		if err.Error != nil {
			return
		}
	}

	return
}
