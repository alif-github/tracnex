package Task

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"os"
	"runtime"
	"strings"
	"time"
)

type houseKeepingFileUpload struct {
	AbstractScheduledTask
}

var HouseKeepingFileUpload = houseKeepingFileUpload{}.New()

func (input houseKeepingFileUpload) New() (output houseKeepingFileUpload) {
	output.RunType = "scheduler.house_keeping_file_upload"
	return
}

func (input houseKeepingFileUpload) Start() {
	if config.ApplicationConfiguration.GetSchedulerStatus().IsActive {
		input.StartTask(input.RunType, input.mainHouseKeepingFileUpload)
	}
}

func (input houseKeepingFileUpload) mainHouseKeepingFileUpload() {
	input.doHouseKeepingFileUpload()
}

type mappingRemoveFile struct {
	FileName    string `json:"file_name"`
	IsDirectory bool   `json:"is_directory"`
	PathName    string `json:"path_name"`
}

func (input houseKeepingFileUpload) doHouseKeepingFileUpload() {
	//-- Cron : 0 2 1 */1 *
	//-- At 02:00 on day-of-month 1 in every month.

	var (
		startTime     = time.Now()
		baseDir       = config.ApplicationConfiguration.GetDataDirectory().BaseDirectoryPath
		importPath    = config.ApplicationConfiguration.GetDataDirectory().ImportPath
		group         = "House Keeping"
		typeJob       = "File Upload"
		nameSch       = "Scheduler House Keeping File"
		jobProcess    repository.JobProcessModel
		err           errorModel.ErrorModel
		mapRemoveFile []mappingRemoveFile
		logModel      = applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
		jobID         = util.GetUUID()
		counting      = 0
	)

	if runtime.GOOS == "windows" {
		goPaths := os.Getenv("GOPATH")
		paths := strings.Split(goPaths, ";")
		baseDir = strings.Replace(paths[0], "\\", "/", -1)
	}

	rootPath := baseDir + importPath

	defer func() {
		if err.Error != nil {
			var (
				errUpdate  errorModel.ErrorModel
				errMessage string
				timeDone   = time.Now()
			)

			jsonMapRemoveFileButErr, _ := json.Marshal(mapRemoveFile)

			jobProcess.Status.String = constanta.JobProcessErrorStatus
			jobProcess.UpdatedAt.Time = timeDone
			jobProcess.ContentDataOut.String = string(jsonMapRemoveFileButErr)
			jobProcess.Counter.Int32 = int32(counting)
			jobProcess.Total.Int32 = int32(counting)
			jobProcess.MessageAlert.String = err.Error.Error()

			errUpdate = dao.JobProcessDAO.UpdateFullJobProcess(serverconfig.ServerAttribute.DBConnection, jobProcess)
			if errUpdate.Error != nil {
				errMessage = errUpdate.Error.Error()
			} else {
				errMessage = err.Error.Error()
			}

			logModel.Message = fmt.Sprintf(`----------> [FINISH] [ERROR: %s] Scheduler process house keeping file upload is failed`, errMessage)
			logModel.Time = int64(timeDone.Sub(startTime))
			logModel.RequestID = jobID
			logModel.Status = 500
			util.LogError(logModel.ToLoggerObject())
		}
	}()

	logModel.Message = "----------> [START] Scheduler process house keeping file upload"
	logModel.Status = 200
	util.LogInfo(logModel.ToLoggerObject())

	jobProcess = repository.JobProcessModel{
		Level:     sql.NullInt32{Int32: 0},
		JobID:     sql.NullString{String: jobID},
		Group:     sql.NullString{String: group},
		Type:      sql.NullString{String: typeJob},
		Name:      sql.NullString{String: nameSch},
		Status:    sql.NullString{String: constanta.JobProcessOnProgressStatus},
		CreatedBy: sql.NullInt64{Int64: constanta.SystemID},
		CreatedAt: sql.NullTime{Time: startTime},
		UpdatedAt: sql.NullTime{Time: startTime},
	}

	if err = dao.JobProcessDAO.InsertJobProcess(serverconfig.ServerAttribute.DBConnection, jobProcess); err.Error != nil {
		return
	}

	if err = input.processHouseKeeping(startTime, rootPath, &mapRemoveFile, &counting, jobID); err.Error != nil {
		return
	}

	timeFinish := time.Now()

	jsonMapRemoveFile, _ := json.Marshal(mapRemoveFile)

	jobProcess.Status.String = constanta.JobProcessDoneStatus
	jobProcess.UpdatedAt.Time = timeFinish
	jobProcess.Counter.Int32 = int32(counting)
	jobProcess.Total.Int32 = int32(counting)
	jobProcess.ContentDataOut.String = string(jsonMapRemoveFile)

	if err = dao.JobProcessDAO.UpdateFullJobProcess(serverconfig.ServerAttribute.DBConnection, jobProcess); err.Error != nil {
		return
	}

	duration := timeFinish.Sub(startTime)

	logModel.Message = "----------> [FINISH] Scheduler process house keeping file upload"
	logModel.RequestID = jobID
	logModel.Time = int64(duration.Seconds())
	logModel.Status = 200

	util.LogInfo(logModel.ToLoggerObject())
}

func (input houseKeepingFileUpload) processHouseKeeping(timeNow time.Time, path string, mapRemoveFile *[]mappingRemoveFile, counting *int, jobID string) (err errorModel.ErrorModel) {
	var last6Month = timeNow.AddDate(0, -6, 0)
	if err = input.openDirectory(path, last6Month, mapRemoveFile, counting, jobID); err.Error != nil {
		return
	}

	return
}

func (input houseKeepingFileUpload) openDirectory(path string, timeLimit time.Time, mapRemoveFile *[]mappingRemoveFile, counting *int, jobID string) (err errorModel.ErrorModel) {
	var (
		fileName = "HouseKeepingFileUpload.go"
		funcName = "openDirectory"
	)

	files, errS := ioutil.ReadDir(path)
	if errS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errS)
		return
	}

	for _, file := range files {

		filePathPoint := path + "/" + file.Name()

		//-- Is file or folder and file modified before last6Month or file modified on last6Month
		if file.ModTime().Before(timeLimit) || file.ModTime().Equal(timeLimit) {

			errorS := os.RemoveAll(filePathPoint)
			if errorS != nil {

				errLogModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())

				errLogModel.Message = "----------> [ERROR] : " + errorS.Error()
				errLogModel.RequestID = jobID
				errLogModel.Status = 500

				util.LogError(errLogModel.ToLoggerObject())
			}

			*mapRemoveFile = append(*mapRemoveFile, mappingRemoveFile{
				FileName:    file.Name(),
				IsDirectory: file.IsDir(),
				PathName:    filePathPoint,
			})

			*counting++

		} else {
			if file.IsDir() {
				if err = input.openDirectory(filePathPoint, timeLimit, mapRemoveFile, counting, jobID); err.Error != nil {
					return
				}
			}
		}
	}

	return
}
