package LogFileService

import (
	"bufio"
	"fmt"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/service"
	"os"
	"strconv"
	"strings"
	"time"
)

type downloadLogService struct {
	service.AbstractService
}

var DownloadLogService = downloadLogService{}.New()

func (input downloadLogService) New() (output downloadLogService) {
	output.FileName = "DownloadLogService.go"
	return
}

func (input downloadLogService) StartService(r *http.Request, _ *applicationModel.ContextModel) (file *os.File, header map[string]string, err errorModel.ErrorModel) {
	var (
		listFile      = config.ApplicationConfiguration.GetLogFile()
		fileName      string
		zipLocation   string
		zipFileName   string
		fileNameSplit []string
		isFull        bool
		startTime     time.Time
		endTime       time.Time
	)

	//--- Validate
	startTime, endTime, isFull, err = input.validateRequestInput(r)
	if err.Error != nil {
		return
	}

	for i := 0; i < len(listFile); i++ {
		if listFile[i] != "stdout" {
			fileName = listFile[i]
			break
		}
	}

	if !isFull {
		//--- Create New File Log
		err = input.createNewFileLog(&fileName, startTime, endTime)
		if err.Error != nil {
			return
		}
	}

	fileNameSplit = strings.Split(fileName, "/")
	zipLocation = strings.Join(fileNameSplit[0:len(fileNameSplit)-1], "/")
	zipFileName = "log.zip"

	return service.GetFileForDownload(fileName, "", true, zipLocation+"/"+zipFileName)
}

func (input downloadLogService) createNewFileLog(fileName *string, inputStartTime, inputEndTime time.Time) (err errorModel.ErrorModel) {
	var (
		serviceName        = "DownloadLogFileService.go"
		funcName           = "createNewFileLog"
		newFileLog         = "new_nextrac.log"
		timeNow            = time.Now()
		startTime, endTime time.Time
		errS               error
		oldFile            *os.File
		outputFile         *os.File
		isStart            bool
	)

	startTime = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), inputStartTime.Hour(), inputStartTime.Minute(), inputStartTime.Second(), 0, timeNow.Location())
	if !inputEndTime.IsZero() {
		endTime = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), inputEndTime.Hour(), inputEndTime.Minute(), inputEndTime.Second(), 0, timeNow.Location())
	}

	oldFile, errS = os.Open(*fileName)
	if errS != nil {
		err = errorModel.GenerateUnknownError(serviceName, funcName, errS)
		return
	}

	defer func() {
		_ = oldFile.Close()
	}()

	//--- Create File
	outputFile, errS = os.Create(newFileLog)
	if errS != nil {
		err = errorModel.GenerateUnknownError(serviceName, funcName, errS)
		return
	}

	defer func() {
		_ = outputFile.Close()
	}()

	writer := bufio.NewWriter(outputFile)
	defer func() {
		_ = writer.Flush()
	}()

	scanner := bufio.NewScanner(oldFile)
	for scanner.Scan() {
		var (
			customLine, line string
			customLineCol    []string
			timeStamp        string
			timeLog          time.Time
			errorS           error
		)

		line = scanner.Text()
		if (!strings.Contains(line, fmt.Sprintf(`{"level"`)) || endTime.IsZero()) && isStart {
			_, errorS = writer.WriteString(line + "\n")
			if errorS != nil {
				err = errorModel.GenerateUnknownError(serviceName, funcName, errS)
				return
			}
			continue
		}

		customLine = strings.ReplaceAll(line, fmt.Sprintf(`"`), "")
		customLine = strings.ReplaceAll(customLine, fmt.Sprintf(`{`), "")
		customLine = strings.ReplaceAll(customLine, fmt.Sprintf(`}`), "")
		customLineCol = strings.Split(customLine, ",")
		for _, itemCustomLineCol := range customLineCol {
			var isTimeStamp bool
			s := strings.ReplaceAll(itemCustomLineCol, ":", ",")
			s = strings.Replace(s, ",", ":", 1)
			sCol := strings.Split(s, ":")
			for j, itemS := range sCol {
				if j == 0 {
					if itemS == "timestamp" {
						continue
					}
					break
				} else {
					timeStamp = strings.ReplaceAll(itemS, ",", ":")
					isTimeStamp = true
					break
				}
			}
			if isTimeStamp {
				break
			}
		}

		timeLog, errS = time.Parse(fmt.Sprintf(`2006-01-02T15:04:05.9999999-07:00`), timeStamp)
		if errS != nil {
			err = errorModel.GenerateUnknownError(serviceName, funcName, errS)
			return
		}

		if timeLog.Equal(startTime) || timeLog.After(startTime) {
			if timeLog.Before(endTime) || endTime.IsZero() {
				_, errorS = writer.WriteString(line + "\n")
				if errorS != nil {
					err = errorModel.GenerateUnknownError(serviceName, funcName, errS)
					return
				}
				isStart = true
				continue
			}
		}

		if !timeLog.IsZero() {
			if timeLog.After(endTime) && !endTime.IsZero() {
				break
			}
		}

		isStart = false
	}

	errS = scanner.Err()
	if errS != nil {
		err = errorModel.GenerateUnknownError(serviceName, funcName, errS)
		return
	}

	*fileName = newFileLog
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input downloadLogService) validateRequestInput(r *http.Request) (startTime, endTime time.Time, isFull bool, err errorModel.ErrorModel) {
	var (
		serviceName  = "DownloadLogFileService.go"
		funcName     = "validateRequestInput"
		defaultTime  = "15:04:05"
		isFullStr    string
		startTimeStr string
		endTimeStr   string
		errorS       error
	)

	isFullStr = service.GenerateQueryValue(r.URL.Query()["is_full"])
	startTimeStr = service.GenerateQueryValue(r.URL.Query()["time_start"])
	endTimeStr = service.GenerateQueryValue(r.URL.Query()["time_end"])

	if !util.IsStringEmpty(isFullStr) {
		isFull, errorS = strconv.ParseBool(isFullStr)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(serviceName, funcName, errorS)
			return
		}
	}

	if !isFull {
		if util.IsStringEmpty(startTimeStr) {
			err = errorModel.GenerateUnknownDataError(serviceName, funcName, "Start Time")
			return
		}

		startTime, errorS = time.Parse(defaultTime, startTimeStr)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(serviceName, funcName, errorS)
			return
		}

		if !util.IsStringEmpty(endTimeStr) {
			endTime, errorS = time.Parse(defaultTime, endTimeStr)
			if errorS != nil {
				err = errorModel.GenerateUnknownError(serviceName, funcName, errorS)
				return
			}
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
