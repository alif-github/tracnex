package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"hash/fnv"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/backgroundJobModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_request"
	"nexsoft.co.id/nextrac2/resource_common_service/dto/authentication_response"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var source = rand.NewSource(time.Now().UnixNano())

func GenerateQueryValue(queryValues []string) string {
	if queryValues == nil {
		return ""
	} else {
		return queryValues[0]
	}
}

func CheckIsOnlyHaveOwnPermission(contextModel applicationModel.ContextModel) (int64, bool) {
	if strings.Contains(contextModel.PermissionHave, "own") {
		return contextModel.AuthAccessTokenModel.ResourceUserID, true
	}
	return 0, false
}

func IsHaveAllPermission(contextModel *applicationModel.ContextModel) bool {
	return strings.Contains(contextModel.PermissionHave, "all")
}

func GenerateInitiateRoleDTOOut(listPermission map[string][]string, rolePermissionFromDB map[string][]string, needIsChecked bool) (result out.InitiateInsertUpdateRoleDTOOut) {
	for key := range listPermission {
		var permissionsDTOOut out.Permissions
		for i := 0; i < len(listPermission[key]); i++ {
			var permissionDTOOut out.Permission
			permissionDTOOut.Label = listPermission[key][i]
			permissionDTOOut.Value = key + ":" + listPermission[key][i]
			if needIsChecked {
				permissionDTOOut.IsChecked = common.ValidateStringContainInStringArray(rolePermissionFromDB[key], listPermission[key][i])
			}
			permissionsDTOOut.Permission = append(permissionsDTOOut.Permission, permissionDTOOut)
		}
		permissionsDTOOut.Key = key + ":" + GeneratePermissionKey(listPermission[key])
		permissionsDTOOut.Menu = key
		result.Permissions = append(result.Permissions, permissionsDTOOut)
	}
	return
}

func GenerateInitiateDataGroupDTOOut(listScope map[string][]string, roleScopesFromDB map[string][]string, needIsChecked bool) (result out.InitiateInsertUpdateDataGroupDTOOut) {
	for key := range listScope {
		var permissionsDTOOut out.Scopes
		for i := 0; i < len(listScope[key]); i++ {
			var permissionDTOOut out.Scope
			permissionDTOOut.Label = listScope[key][i]
			permissionDTOOut.Value = key + ":" + listScope[key][i]
			if needIsChecked {
				permissionDTOOut.IsChecked = common.ValidateStringContainInStringArray(roleScopesFromDB[key], listScope[key][i])
			}
			permissionsDTOOut.Scope = append(permissionsDTOOut.Scope, permissionDTOOut)
		}
		permissionsDTOOut.Key = key + ":" + GeneratePermissionKey(listScope[key])
		permissionsDTOOut.Menu = key
		result.Scopes = append(result.Scopes, permissionsDTOOut)
	}
	return
}

func GeneratePermissionKey(arrayOfPermission []string) (result string) {
	for i := 0; i < len(arrayOfPermission); i++ {
		result += arrayOfPermission[i]
		if i != len(arrayOfPermission)-1 {
			result += ", "
		}
	}
	return
}

func GenerateHashMapPermissionAndDataScope(listPermission []string, isRemoveDuplicate bool, isDataGroup bool) (result map[string][]string) {
	sort.Slice(listPermission, func(i, j int) bool {
		if strings.Contains(listPermission[i], ".") && strings.Contains(listPermission[j], ".") {
			idxI := strings.Split(listPermission[i], ".")
			idxJ := strings.Split(listPermission[j], ".")
			if idxI[0] == idxJ[0] {
				return len(idxI) < len(idxJ)
			}
		} else if strings.Contains(listPermission[i], ".") {
			return false
		} else if strings.Contains(listPermission[j], ".") {
			return true
		}
		return listPermission[i] < listPermission[j]
	})

	result = make(map[string][]string)
	for i := 0; i < len(listPermission); i++ {
		var menu string
		permission := listPermission[i]
		var validationResult bool
		if !isDataGroup {
			validationResult, _ = util.IsNexsoftPermissionStandardValid(permission)
		} else {
			validationResult = true
		}
		if !validationResult {
			continue
		}
		splitPermission := strings.Split(permission, ":")
		splitDotMenu := strings.Split(splitPermission[0], ".")
		sizeWithDot := len(splitDotMenu)
		if isRemoveDuplicate {
			var isAvailable = false
			for sizeWithDot > 0 {
				menu = ""
				for i := 0; i < sizeWithDot; i++ {
					menu += splitDotMenu[i]
					if i < sizeWithDot-1 {
						menu += "."
					}
				}
				if common.ValidateStringContainInStringArray(result[menu], splitPermission[1]) {
					isAvailable = true
					break
				}
				sizeWithDot--
			}

			if !isAvailable {
				result[splitPermission[0]] = append(result[splitPermission[0]], splitPermission[1])
			}
		} else {
			result[splitPermission[0]] = append(result[splitPermission[0]], splitPermission[1])
		}
	}

	return
}

func ValidateRole(permission map[string][]string, fileName string, funcName string) {
	for key := range permission {
		permission[key] = getRoles(permission[key])
	}

}

func getRoles(arrRoles []string) (result []string) {
	var isValid bool
	var role string

	isValid, role = RoleChecker(common.ViewData, arrRoles)
	if isValid {
		result = append(result, role)
	}

	isValid, role = RoleChecker(common.UpdateData, arrRoles)
	if isValid {
		result = append(result, role)
	}

	isValid, role = RoleChecker(common.InsertData, arrRoles)
	if isValid {
		result = append(result, role)
	}

	isValid, role = RoleChecker(common.DeleteData, arrRoles)
	if isValid {
		result = append(result, role)
	}

	isValid, role = RoleChecker(common.ChangePassword, arrRoles)
	if isValid {
		result = append(result, role)
	}

	return
}

func RoleChecker(permissionNeed string, listPermission []string) (bool, string) {
	for i := 0; i < len(listPermission); i++ {
		if listPermission[i] == permissionNeed+"-all" {
			newListPermission := listPermission[(i + 1):]
			for j := 0; j < len(newListPermission); j++ {
				if newListPermission[j] == permissionNeed {
					return true, newListPermission[j]
				}
			}
			return true, listPermission[i]
		}
		if listPermission[i] == permissionNeed {
			return true, listPermission[i]
		}
		if listPermission[i] == permissionNeed+"-own" {
			newListPermission := listPermission[(i + 1):]
			for j := 0; j < len(newListPermission); j++ {
				if newListPermission[j] == permissionNeed {
					return true, newListPermission[j]
				}
			}
			return true, listPermission[i]
		}
	}
	return false, ""
}

func CheckScopeAndGetTeamIDAndAssetCategoryID(result map[string][]string) (team []int64, assetCategory []int64) {
	for key := range result {
		id := splitScopeToGetID(key)
		if id > 0 {
			if strings.Contains(key, "team") {
				team = append(team, id)
			} else {
				assetCategory = append(assetCategory, id)
			}
		}
	}
	return
}

func splitScopeToGetID(scope string) int64 {
	scopeSplit := strings.Split(scope, ".")
	if len(scopeSplit) > 1 {
		id, _ := strconv.Atoi(scopeSplit[1])
		return int64(id)
	} else {
		return 0
	}
}

func CheckDBError(err errorModel.ErrorModel, constraint string) bool {
	if err.Error != nil && err.CausedBy != nil {
		return strings.Contains(err.CausedBy.Error(), "\""+constraint+"\"")
	}
	return false
}

func GetAzureDateContainer() (result string) {
	timeNow := time.Now()

	year, month, day := timeNow.Date()

	result += strconv.Itoa(year) + "/"
	if int(month) < 10 {
		result += "0"
	}
	result += strconv.Itoa(int(month)) + "/"

	if day < 10 {
		result += "0"
	}
	result += strconv.Itoa(day) + "/"

	return
}

func ReformatFileName(filename string, userID int64) (result string) {
	fileNameSplit := strings.Split(filename, ".")
	fileNameWithoutExtension := strings.Join(fileNameSplit[0:len(fileNameSplit) - 1], ".")
	extension := fileNameSplit[len(fileNameSplit) - 1]
	fileNameWithoutExtension = strings.ToLower(fileNameWithoutExtension)
	regex := regexp.MustCompile("[a-z0-9-_]")
	for i := 0; i < len(fileNameWithoutExtension); i++ {
		if regex.MatchString(string(fileNameWithoutExtension[i])) {
			result += string(fileNameWithoutExtension[i])
		}
	}

	randomNumber := GenerateRandom("0123456789", 6)
	result += "_" + strconv.Itoa(int(userID))
	result += "_" + randomNumber

	return result + "." + extension
}

func GenerateRandom(possibleChar string, length int) string {
	random := rand.New(source)
	b := make([]rune, length)
	letterRunes := []rune(possibleChar)

	for i := range b {
		b[i] = letterRunes[random.Intn(len(letterRunes))]
	}

	return string(b)
}

func ReadFileWithMultipart(request *http.Request, numOfFile int, validator func(in.MultipartFileDTO) errorModel.ErrorModel) (result []in.MultipartFileDTO, totalSize int64, err errorModel.ErrorModel) {
	funcName := "readFileWithMultipart"
	var temp in.MultipartFileDTO
	var errs error

	for i := 1; i <= numOfFile; i++ {
		var file multipart.File
		file, temp.Filename, temp.Size, errs = util.ReadMultipartFile(request, "photo"+strconv.Itoa(i))
		if errs != nil {
			if errs.Error() == "http: no such file" {
				err = errorModel.GenerateNonErrorModel()
				return
			}
			return
		}
		if temp.Size > 0 {
			totalSize += temp.Size
			buf := bytes.NewBuffer(nil)
			if _, errs := io.Copy(buf, file); errs != nil {
				err = errorModel.GenerateUnknownError("ServiceUtil.go", funcName, errs)
				return
			}
			temp.FileContent = buf.Bytes()
			err = validator(temp)
			if err.Error != nil {
				return
			}
			result = append(result, temp)
		}
	}

	return
}

func UploadFileToLocalCDN(containerName string, file *[]in.MultipartFileDTO, id int64) (err errorModel.ErrorModel) {
	funcFileName := "service/Util.go"
	funcName := "UploadFileToLocalCDN"
	if string(containerName[len(containerName)-1]) == "/" {
		containerName = containerName[0 : len(containerName)-1]
	}

	if file != nil {
		temp := *file
		for i := 0; i < len(temp); i++ {
			fileName := ReformatFileName(temp[i].Filename, id)

			directory := config.ApplicationConfiguration.GetCDN().RootPath + config.ApplicationConfiguration.GetCDN().Suffix +
				containerName + config.ApplicationConfiguration.GetCDN().Suffix
			_ = os.MkdirAll(directory, 0770)

			fmt.Println("Upload to CDN")
			fmt.Println("Directory File Name Old : ", directory+temp[i].Filename)
			fmt.Println("Directory File Name Reformat : ", directory+fileName)
			//fmt.Println("File Content  : ", temp[i].FileContent)
			//fmt.Println("File Content convert : ", []byte(temp[i].FileContent))

			errs := ioutil.WriteFile(directory+fileName, temp[i].FileContent, 0660)
			if errs != nil {
				fmt.Println("Error Write File : ", errs.Error())
				return errorModel.GenerateUnknownError(funcFileName, funcName, errs)
			}

			fmt.Println("Success Write File, File Path : " + directory + fileName)

			temp[i].Host = config.ApplicationConfiguration.GetCDN().Host + config.ApplicationConfiguration.GetCDN().Suffix
			temp[i].Host = temp[i].Host[0 : len(temp[i].Host)-1]
			temp[i].Path = strings.Replace(directory+fileName, config.ApplicationConfiguration.GetCDN().RootPath+
				config.ApplicationConfiguration.GetCDN().Suffix, "", -1)
			temp[i].Path = "/" + temp[i].Path

			fmt.Println("Log : " + temp[i].Host + temp[i].Path)
		}
		file = &temp
	}

	return errorModel.GenerateNonErrorModel()
}

func UploadFileToAzure(file *[]in.MultipartFileDTO) (err errorModel.ErrorModel) {
	funcFileName := "ServiceUtil.go"
	funcName := "UploadFileToAzure"

	if file != nil {
		temp := *file
		for i := 0; i < len(temp); i++ {
			credential, errCred := azblob.NewSharedKeyCredential(config.ApplicationConfiguration.GetAzure().AccountName,
				config.ApplicationConfiguration.GetAzure().AccountKey)

			if errCred != nil {
				fmt.Println("Error Credential ", errCred.Error())
				return errorModel.GenerateUnknownError(funcFileName, funcName, errCred)
			}
			pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

			host := config.ApplicationConfiguration.GetAzure().Host
			suffix := config.ApplicationConfiguration.GetAzure().Suffix

			var hostSuffix = host + suffix + GetAzureDateContainer()
			URL, _ := url.Parse(hostSuffix)

			containerURL := azblob.NewContainerURL(*URL, pipeline)
			fmt.Printf("Creating a container named %s\n", URL)

			//containerURL := GetContainer()
			ctx := context.Background()
			//file, errOpen := ioutil.ReadFile(temp[i].Path + temp[i].Filename)
			//if errOpen != nil {
			//	fmt.Println("Error Open => ", errOpen.Error())
			//} else {
			//	fmt.Println("Open file success")
			//}

			blobURL := containerURL.NewBlockBlobURL(temp[i].Filename)
			_, errs := azblob.UploadBufferToBlockBlob(ctx, temp[i].FileContent, blobURL, azblob.UploadToBlockBlobOptions{BlockSize: 4 * 1024 * 1024, Parallelism: 16})
			if errs != nil {
				fmt.Println("Error blob ===>", errs.Error())
				return errorModel.GenerateUnknownError(funcFileName, funcName, errs)
			}

			// p := filepath.FromSlash("path/to/file")
			newPath := URL.String() + temp[i].Filename
			fmt.Println("new path =====>" + newPath)
		}
		file = &temp
	}
	return errorModel.GenerateNonErrorModel()
}

func GetContainer() azblob.ContainerURL {
	host := config.ApplicationConfiguration.GetAzure().Host
	suffix := config.ApplicationConfiguration.GetAzure().Suffix
	URL, _ := url.Parse(host + suffix)
	containerURL := azblob.NewContainerURL(*URL, serverconfig.ServerAttribute.AzurePipeline)

	return containerURL
}

func DeleteFileFromCDN(file []in.MultipartFileDTO) {
	if file == nil {
		return
	}

	ctx := context.Background()
	for i := 0; i < len(file); i++ {
		URL, _ := url.Parse(file[i].Host)
		if file[i].Host != "" && file[i].Path != "" {
			if file[i].Host != config.ApplicationConfiguration.GetCDN().RootPath {
				containerURL := azblob.NewContainerURL(*URL, serverconfig.ServerAttribute.AzurePipeline)
				path := file[i].Path[1:]
				blobURL := containerURL.NewBlockBlobURL(path)

				_, _ = blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
				//if err != nil {
				//	return errorModel.GenerateUnknownError("", "DeleteFileFromCDN", err)
				//}
			}
			_ = os.Remove(config.ApplicationConfiguration.GetCDN().RootPath + "/" + file[i].Path)
		}
	}

	return
}

func UploadListFileToAzure(data interface{}, contextModel applicationModel.ContextModel, saveToDB func(*sql.Tx, in.MultipartFileDTO) errorModel.ErrorModel) {
	if data == nil {
		return
	}

	dataFile := data.([]in.MultipartFileDTO)
	var err errorModel.ErrorModel

	err = UploadFileToAzure(&dataFile)
	if err.Error != nil {
		contextModel.LoggerModel.Status = 500
		contextModel.LoggerModel.Message = err.CausedBy.Error()
		util.LogError(contextModel.LoggerModel.ToLoggerObject())
		return
	}

	tx, errs := serverconfig.ServerAttribute.DBConnection.Begin()
	if errs != nil {
		contextModel.LoggerModel.Status = 500
		contextModel.LoggerModel.Message = errs.Error()
		util.LogError(contextModel.LoggerModel.ToLoggerObject())
		return
	}

	defer func() {
		if errs != nil && err.Error != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	for i := 0; i < len(dataFile); i++ {
		err = saveToDB(tx, dataFile[i])
		if err.Error != nil {
			contextModel.LoggerModel.Status = 500
			contextModel.LoggerModel.Message = err.CausedBy.Error()
			util.LogError(contextModel.LoggerModel.ToLoggerObject())
			return
		}
	}
}

func AppendString(current *string, added string, delimiter string) {
	temp := *current
	if temp == "" {
		temp += added
	} else {
		temp += delimiter + "" + added
	}

}

func GetJobProcess(task backgroundJobModel.ChildTask, contextModel applicationModel.ContextModel, timeNow time.Time) repository.JobProcessModel {
	return repository.JobProcessModel{
		Level:         sql.NullInt32{},
		JobID:         sql.NullString{String: util.GetUUID()},
		Group:         sql.NullString{String: task.Group},
		Type:          sql.NullString{String: task.Type},
		Name:          sql.NullString{String: task.Name},
		Status:        sql.NullString{String: constanta.JobProcessOnProgressStatus},
		CreatedBy:     sql.NullInt64{Int64: contextModel.AuthAccessTokenModel.ResourceUserID},
		CreatedAt:     sql.NullTime{Time: timeNow},
		CreatedClient: sql.NullString{String: contextModel.AuthAccessTokenModel.ClientID},
		UpdatedAt:     sql.NullTime{Time: timeNow},
	}
}

func ValidateSignature(messageDigest string, key string, request *http.Request) bool {
	signature := request.Header.Get(constanta.SignatureHeaderNameConstanta)
	timestamp := request.Header.Get(constanta.TimestampSignatureHeaderNameConstanta)
	internalToken := request.Header.Get(constanta.TokenHeaderNameConstanta)

	if signature == "" || timestamp == "" || internalToken == "" {
		return false
	}
	return util.ValidateSignature(request.Method, request.RequestURI, internalToken, messageDigest, timestamp, key, signature)
}

func DeleteTokenFromRedis(listToken []string) {
	for i := 0; i < len(listToken); i++ {
		serverconfig.ServerAttribute.RedisClient.Del(listToken[i])
	}
}

func SortArrayOfDataScope(listPermission []out.Scopes) []out.Scopes {
	sort.Slice(listPermission, func(i, j int) bool {
		key1 := strings.Replace(listPermission[i].Key, ".", "", -1)
		key2 := strings.Replace(listPermission[j].Key, ".", "", -1)
		return key1 < key2
	})
	return listPermission
}

func CheckAdditionalInformation(currentAdditionalInfo []model.AdditionalInformation, additionalInfoFromDTO []model.AdditionalInformation) (output []model.AdditionalInformation, err errorModel.ErrorModel) {
	funcName := "CheckAdditionalInformation"
	for i := 0; i < len(additionalInfoFromDTO); i++ {
		var found = false
		for j := 0; j < len(currentAdditionalInfo); j++ {
			if currentAdditionalInfo[j].Key == additionalInfoFromDTO[i].Key {
				found = true
				if currentAdditionalInfo[j].Value != additionalInfoFromDTO[i].Value {
					err = errorModel.GenerateDifferentAdditionalInformationValueError("ServiceUtil.go", funcName)
					return
				}
				break
			}
		}
		if !found {
			output = append(output, additionalInfoFromDTO[i])
		}
	}
	err = errorModel.GenerateNonErrorModel()
	return
}

func appendQuery(data map[string]interface{}, field string, value interface{}) map[string]interface{} {
	temp := data[field]
	switch value.(type) {
	case []map[string][]string:
		data = appendListScope(data, field, value.([]map[string][]string))
		break
	case []interface{}:
		if temp == nil {
			if len(data) == 0 {
				data[field] = value
			} else {
				isFound := false
				for key := range data {
					switch data[key].(type) {
					case []map[string][]string:
						temp := data[key].([]map[string][]string)
						for i := 0; i < len(temp); i++ {
							for keyOnHash := range temp[i] {
								if keyOnHash == field {
									isFound = true
									break
								}
							}
						}
					}
				}
				if !isFound {
					data[field] = value
				}
			}
		}
		break
	}
	return data
}

func appendListScope(data map[string]interface{}, field string, value []map[string][]string) map[string]interface{} {
	var fieldValue []map[string][]string

	if data[field] == nil {
		fieldValue = append(fieldValue, value...)
	} else {
		switch data[field].(type) {
		case []string:
			fieldValue = append(fieldValue, value...)
		case []map[string][]string:
			fieldValue = data[field].([]map[string][]string)
			fieldValue = append(fieldValue, value...)
		}
	}
	data[field] = fieldValue

	for i := 0; i < len(fieldValue); i++ {
		dataSlice := fieldValue[i]
		if dataSlice != nil {
			for key := range dataSlice {
				if data[key] != nil {
					delete(data, key)
				}
			}
		}
	}
	return data
}

func ReadParameterByPermissionAndName(permission string, name string, userID int64) (result string, err errorModel.ErrorModel) {
	isFound := false
	splitDotMenu := strings.Split(permission, ".")
	menu := permission
	size := len(splitDotMenu)
	for size > 0 {
		var parameterModel repository.ParameterModel
		menu = ""
		for i := 0; i < size; i++ {
			menu += splitDotMenu[i]
			if i < size-1 {
				menu += "."
			}
		}

		parameterModel, err = dao.ParameterDAO.GetParameterByNameAndCode(serverconfig.ServerAttribute.DBConnection, repository.ParameterModel{
			Permission: sql.NullString{String: permission},
			Name:       sql.NullString{String: name},
		})
		if err.Error != nil {
			return
		}

		if parameterModel.ID.Int64 != 0 {
			isFound = true
			result = parameterModel.Value.String

			var userParameterModel repository.UserParameterModel
			userParameterModel, err = dao.ParameterDAO.GetUserParameter(serverconfig.ServerAttribute.DBConnection, repository.UserParameterModel{
				UserID: sql.NullInt64{Int64: userID},
			})
			if userParameterModel.ID.Int64 != 0 {
				var listParameterUser map[string]string
				_ = json.Unmarshal([]byte(userParameterModel.ParameterValue.String), &listParameterUser)
				if listParameterUser["p"+strconv.Itoa(int(parameterModel.ID.Int64))] != "" {
					result = listParameterUser["p"+strconv.Itoa(int(parameterModel.ID.Int64))]
				}
			}
			break
		}
		size--
	}

	if isFound {
		err = errorModel.GenerateNonErrorModel()
	} else {
		err = errorModel.GenerateParameterNotFoundError("ServiceUtil.go", "ReadParameterByPermissionAndName")
	}

	return
}

func AddResourceNexcloudToNexcloud(addResourceStruct in.AddResourceNexcloud, contextModel *applicationModel.ContextModel) (err errorModel.ErrorModel) {
	funcName := "addResourceNexcloudToNexcloud"
	var payloadMessage model.PayloadResponse

	internalToken := resource_common_service.GenerateInternalToken(constanta.NexCloudResourceID, 0, contextModel.AuthAccessTokenModel.ClientID, config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
	nexcloudAPIServer := config.ApplicationConfiguration.GetNexcloudAPI()
	addResourceNexcloudUrl := nexcloudAPIServer.Host + nexcloudAPIServer.PathRedirect.AddResourceClient

	header := make(map[string][]string)
	header[common.AuthorizationHeaderConstanta] = []string{internalToken}

	statusCode, _, bodyResult, errorS := common.HitAPI(addResourceNexcloudUrl, header, util.StructToJSON(addResourceStruct), "POST", *contextModel)

	if errorS != nil {
		err = errorModel.GenerateUnknownError("ServiceUtil.go", funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &payloadMessage)

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(payloadMessage.Status.Message)
		//todo if client id has been used, then handle
		err = errorModel.GenerateAuthenticationServerError("ServiceUtil.go", funcName, statusCode, payloadMessage.Status.Code, causedBy)
		return
	}
	return
}

func CheckClientOrUserInAuth(checkClientUserStruct authentication_request.CheckClientOrUser, contextModel *applicationModel.ContextModel) (checkClientUserResp authentication_response.CheckClientOrUserResponse, err errorModel.ErrorModel) {
	funcName := "CheckClientOrUserInAuth"

	internalToken := resource_common_service.GenerateInternalToken("auth", 0, "", config.ApplicationConfiguration.GetServerResourceID(), constanta.IndonesianLanguage)
	authenticationServer := config.ApplicationConfiguration.GetAuthenticationServer()
	checkClientUserUrl := authenticationServer.Host + authenticationServer.PathRedirect.InternalClient.CheckClientUser

	statusCode, bodyResult, errorS := common.HitCheckClientUserAuthenticationServer(internalToken, checkClientUserUrl, checkClientUserStruct, contextModel)
	if errorS != nil {
		err = errorModel.GenerateUnknownError("ServiceUtil.go", funcName, errorS)
		return
	}

	_ = json.Unmarshal([]byte(bodyResult), &checkClientUserResp)

	if statusCode == 200 {
		err = errorModel.GenerateNonErrorModel()
	} else {
		causedBy := errors.New(checkClientUserResp.Nexsoft.Payload.Status.Message)
		err = errorModel.GenerateAuthenticationServerError("ServiceUtil.go", funcName, statusCode, checkClientUserResp.Nexsoft.Payload.Status.Code, causedBy)
		return
	}

	return
}

func CustomFailedResponsePayload(outputStruct interface{}, erorS errorModel.ErrorModel,
	contextModel *applicationModel.ContextModel) (outputPayload interface{}, message string, err errorModel.ErrorModel) {

	var customPayloadOutput out.Payload
	message = util2.GenerateI18NErrorMessage(erorS, contextModel.AuthAccessTokenModel.Locale)

	customPayloadOutput.Data = out.PayloadData{
		Content: &outputStruct,
	}

	//true will change in func writeSuccessResponse be false, default bool is false
	customPayloadOutput.Status = out.StatusResponse{
		Success: true,
		Code:    erorS.Error.Error(),
		Message: message,
		Detail:  erorS.AdditionalInformation,
	}

	outputPayload = customPayloadOutput
	err = errorModel.GenerateNonErrorModel()
	return
}

func CustomSuccessResponsePayload(outputStruct interface{}, message string,
	contextModel *applicationModel.ContextModel) (code string, outputPayload interface{}) {

	var customSuccessPayloadOutput out.Payload
	code = util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil)

	customSuccessPayloadOutput.Data = out.PayloadData{
		Content: &outputStruct,
	}

	customSuccessPayloadOutput.Status = out.StatusResponse{
		Code:    code,
		Message: message,
	}

	outputPayload = customSuccessPayloadOutput
	return
}

func NewErrorAddResource(bundle *i18n.Bundle, contextModel *applicationModel.ContextModel, messageID string, fileName string, funcName string) (detail string, err errorModel.ErrorModel) {

	detail = util2.GenerateI18NServiceMessage(bundle, messageID, contextModel.AuthAccessTokenModel.Locale, nil)
	err = errorModel.GenerateAuthenticationServerAddResourceError(fileName, funcName, []string{detail})

	return
}

func ReadPathParamID(request *http.Request) (id int64, err errorModel.ErrorModel) {
	funcName := "ReadPathParamID"

	strId, ok := mux.Vars(request)["ID"]
	idParam, errConvert := strconv.Atoi(strId)
	id = int64(idParam)

	if !ok || errConvert != nil {
		err = errorModel.GenerateUnsupportedRequestParam("ServiceUtil.go", funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func CountUpdateDBJobProcessCounter(total int32) int {
	updateDBEvery := int(float32(total) * 1.0 / 100.0)
	if updateDBEvery < 10 {
		updateDBEvery = 10
	}

	return updateDBEvery
}

func ValidateScope(contextModel *applicationModel.ContextModel, checkedScope []string) map[string]interface{} {
	var authenticationModel model.AuthenticationModel
	_ = json.Unmarshal([]byte(contextModel.AuthAccessTokenModel.RedisAuthAccessTokenModel.Authentication), &authenticationModel)
	return CheckScope(authenticationModel.Data.Scope, checkedScope)
}

func CheckScope(listScope map[string]interface{}, checkedScope []string) map[string]interface{} {
	result := make(map[string]interface{})
	numberOfAddedQuery := 0
	for key := range listScope {
		switch listScope[key].(type) {
		case []interface{}:
			for i := 0; i < len(checkedScope); i++ {
				scopeQuery, isExist := checkIsScopeContains(key, checkedScope[i])
				if isExist {
					numberOfAddedQuery++
					result = appendQuery(result, checkedScope[i], listScope[scopeQuery])
				}
			}
		case []map[string][]string:
			data := listScope[key].([]map[string][]string)
			isExist := isScopeExistInListHashmap(data, checkedScope)
			if isExist {
				numberOfAddedQuery++
				keySplit := strings.Split(key, ".")
				result = appendQuery(result, keySplit[len(keySplit)-1], listScope[key].([]map[string][]string))
			}
		}
	}

	if numberOfAddedQuery >= len(checkedScope) {
		return result
	} else {
		return nil
	}
}

func checkIsScopeContains(listScopeKey string, scope string) (string, bool) {
	splitMustHavePermission := strings.Split(scope, ":")
	menu := splitMustHavePermission[0]
	splitDotMenu := strings.Split(scope, ".")
	size := len(splitDotMenu)
	for size > 0 {
		menu = ""
		for i := 0; i < size; i++ {
			menu += splitDotMenu[i]
			if i < size-1 {
				menu += "."
			}
		}
		if menu == listScopeKey {
			return menu, true
		}
		size--
	}
	return "", false
}

func isScopeExistInListHashmap(listScope []map[string][]string, listNeedScope []string) (isAllExist bool) {
	for i := 0; i < len(listScope); i++ {
		for key := range listScope[i] {
			isAllExist = false
			for j := 0; j < len(listNeedScope); j++ {
				_, isExist := checkIsScopeContains(key, listNeedScope[j])
				isAllExist = isExist || isAllExist
			}
			if !isAllExist {
				return false
			}
		}
	}
	return true
}

func GetArrayInterfaceFromStringCollection(inputStr []string) (output []interface{}) {
	for _, item := range inputStr {
		output = append(output, item)
	}
	return output
}

func LogMessage(message string, statusCode int) {
	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
	logModel.Status = statusCode
	logModel.Message = message

	if statusCode == 200 {
		util.LogInfo(logModel.ToLoggerObject())
	} else {
		util.LogError(logModel.ToLoggerObject())
	}
}

func GetErrorMessage(err errorModel.ErrorModel, contextModel applicationModel.ContextModel) string {
	errCode := err.Error.Error()
	errMessage := util2.GenerateI18NErrorMessage(err, contextModel.AuthAccessTokenModel.Locale)
	if errMessage == errCode {
		if err.CausedBy != nil {
			errMessage = err.CausedBy.Error()
		}
	}

	if err.Code == 500 {
		errMessage = err.CausedBy.Error()
	}
	return errMessage
}

func LogMessageWithErrorModel(err errorModel.ErrorModel) {
	logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
	logModel.Code = err.Error.Error()
	logModel.Status = err.Code
	logModel.Class = "[" + err.FileName + "," + err.FuncName + "]"
	errMessage := util2.GenerateI18NErrorMessage(err, constanta.DefaultApplicationsLanguage)
	if errMessage == logModel.Code {
		if err.CausedBy != nil {
			errMessage = err.CausedBy.Error()
		}
	}
	logModel.Message = errMessage
	if logModel.Status == 200 {
		util.LogInfo(logModel.ToLoggerObject())
	} else {
		util.LogError(logModel.ToLoggerObject())
	}
}

func RandToken(n int) (string, error) {
	bytesTemp := make([]byte, n)
	if _, err := rand.Read(bytesTemp); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytesTemp)[:n], nil
}

func RandTimeToken(n int, salt string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(time.Now().String()))
	_, _ = h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))[:n]
}

func RefactorArrayAggInt(numberStr string) (numberResult []int, err errorModel.ErrorModel) {
	var (
		fileName      = "ServiceUtil.go"
		funcName      = "RefactorArrayAggInt"
		numberStrTemp []string
	)

	if !util.IsStringEmpty(numberStr) {
		numberStr = strings.ReplaceAll(numberStr, "{", "")
		numberStr = strings.ReplaceAll(numberStr, "}", "")
		isContains := strings.Contains(numberStr, ",")

		if isContains {
			numberStrTemp = strings.Split(numberStr, ",")
		} else {
			numberStrTemp = append(numberStrTemp, numberStr)
		}
	}

	for _, valueStrNumber := range numberStrTemp {
		var (
			numberInt int
			errorS    error
		)

		numberInt, errorS = strconv.Atoi(valueStrNumber)
		if errorS != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
			return
		}

		numberResult = append(numberResult, numberInt)
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func GetFileForDownload(fileName string, extension string, isZip bool, zipFileName string) (file *os.File, header map[string]string, err errorModel.ErrorModel) {
	var (
		funcName     = "getFileLog"
		fileNameUsed = fileName
		errs         error
	)

	if isZip {
		if !strings.Contains(zipFileName, "home") && !strings.Contains(zipFileName, ":") {
			zipFileName = "." + zipFileName
		}
		errs := util.ZipFiles(zipFileName, []string{fileName})
		if errs != nil {
			err = errorModel.GenerateUnknownError(fileName, funcName, errs)
			return
		}
		fileNameUsed = zipFileName
	}

	file, errs = os.Open(fileNameUsed)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	fileHeader := make([]byte, 512)
	_, errs = file.Read(fileHeader)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	fileContentType := http.DetectContentType(fileHeader)
	fileStat, _ := file.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	contentDispositionName := fileStat.Name()
	if len(strings.Split(fileStat.Name(), ".")) == 1 {
		contentDispositionName += "." + extension
	}

	header = make(map[string]string)
	header["Content-Disposition"] = "attachment; filename=" + contentDispositionName
	header["Content-Type"] = fileContentType
	header["Content-Length"] = fileSize

	_, errs = file.Seek(0, 0)
	if errs != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errs)
		return
	}

	return
}

func GeneratorLicense(action string, args map[string]string) (stdout []byte, err errorModel.ErrorModel) {
	var (
		fileName = "ServiceUtil.go"
		funcName = "GeneratorLicense"
		errorS   error
		logModel = applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
	)

	defer func() {
		if err.Error != nil || stdout == nil {
			logModel.Status = 500
			logModel.Message = "[FAILED] Generator License Failed Build Encrypt Key -----> " + err.CausedBy.Error()
			util.LogError(logModel.ToLoggerObject())
		} else {
			logModel.Status = 200
			logModel.Message = "[SUCCESS] Generator License Success Build Encrypt Key"
			util.LogInfo(logModel.ToLoggerObject())
		}
	}()

	switch runtime.GOOS {
	case "windows":
		var (
			goPath, finalPath, fullPath string
			paths                       []string
		)

		goPath = os.Getenv("GOPATH")
		paths = strings.Split(goPath, ";")
		finalPath = strings.Replace(paths[0], "\\", "/", -1)
		fullPath = fmt.Sprintf(finalPath + config.ApplicationConfiguration.GetGenerator().Path + "/generator.exe")

		cmd := exec.Command(fullPath, action, args["args1"], args["args2"])
		stdout, errorS = cmd.Output()

		break
	case "linux":
		var rootPath, path, fullPath string

		rootPath = config.ApplicationConfiguration.GetGenerator().RootPath
		path = config.ApplicationConfiguration.GetGenerator().Path
		fullPath = rootPath + path
		
		//--- Cek
		fmt.Println(fmt.Sprintf(`Generator %s -> %s %s %s %s`, action, fullPath+"/generator", action, args["args1"], args["args2"]))
		
		cmd := exec.Command(fullPath+"/generator", action, args["args1"], args["args2"])
		stdout, errorS = cmd.Output()
		
		//--- Cek
		fmt.Println(fmt.Sprintf(`Hasil %s -> %s`, action, string(stdout)))

		break
	default:
		LogMessage("No runtime detected", 200)
	}

	if errorS != nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errorS)
		return
	}

	if stdout == nil {
		err = errorModel.GenerateUnknownError(fileName, funcName, errors.New("generator error, stdout nil"))
		return
	}

	return
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GenerateRandomNumberToString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var letterRunes = []rune("0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}


func OptimisticLock(updatedAtFromDB time.Time, updatedAtFromBody time.Time, fileName string, fieldName string) (err errorModel.ErrorModel) {
	funcName := "OptimisticLock"

	if updatedAtFromDB != updatedAtFromBody {
		fmt.Println("Updated At From DB OptimisticLock : ", updatedAtFromDB)
		fmt.Println("Updated At From Body OptimisticLock : ", updatedAtFromBody)
		err = errorModel.GenerateDataLockedError(fileName, funcName, fieldName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func GetFileBytes(file multipart.File) (result []byte, errModel errorModel.ErrorModel) {
	funcName := "GetFileBytes"

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		errModel = errorModel.GenerateUnknownError("ServiceUtil.go", funcName, err)
		return
	}

	return buf.Bytes(), errorModel.GenerateNonErrorModel()

}
