package endpoint

import (
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/resource_common_service/common"
	"nexsoft.co.id/nextrac2/resource_common_service/model"
	"strconv"
	"strings"
	"time"
)

func ValidatePermissionWithRole(mustHavePermission string, roleModel model.AuthenticationRoleModel) (string, errorModel.ErrorModel) {
	funcName := "ValidatePermissionWithRole"
	var isValid = false
	var permissionAllowed string

	validationResult, errField := util.IsNexsoftPermissionStandardValid(mustHavePermission)
	if !validationResult {
		return permissionAllowed, errorModel.GenerateFieldFormatWithRuleError("ServiceUtil.go", funcName, errField, "Permission", "")
	}

	splitMustHavePermission := strings.Split(mustHavePermission, ":")
	menu := splitMustHavePermission[0]
	permission := splitMustHavePermission[1]
	isValid, permissionAllowed = roleChecker(permission, roleModel.Permission[permission])
	if !isValid {
		return permissionAllowed, errorModel.GenerateForbiddenAccessClientError("ServiceUtil.go", funcName)
	}

	isValid = false
	splitDotMenu := strings.Split(menu, ".")
	size := len(splitDotMenu)
	for size > 0 {
		menu = ""
		for i := 0; i < size; i++ {
			menu += splitDotMenu[i]
			if i < size-1 {
				menu += "."
			}
		}
		if roleModel.Permission[menu] != nil {
			isValid, permissionAllowed = roleChecker(permission, roleModel.Permission[menu])
			if isValid {
				break
			}
		}
		size--
	}

	if !isValid {
		return permissionAllowed, errorModel.GenerateForbiddenAccessClientError("ServiceUtil.go", funcName)
	}

	return permissionAllowed, errorModel.GenerateNonErrorModel()
}

func roleChecker(permissionNeed string, listPermission []string) (bool, string) {
	for i := 0; i < len(listPermission); i++ {
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
		if listPermission[i] == permissionNeed+"-all" {
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

func ReadError(err model.ResourceCommonErrorModel) errorModel.ErrorModel {
	fileName := "EndpointUtil.go"
	funcName := "ReadError"
	if err.Error != nil {
		if err.Error.Error() == common.GenerateUnauthorizedTokenError().Error.Error() {
			return errorModel.GenerateUnauthorizedClientError(fileName, funcName)
		} else if err.Error.Error() == common.GenerateForbiddenByResourceID().Error.Error() {
			return errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		} else if err.Error.Error() == common.GenerateForbiddenByScope().Error.Error() {
			return errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		} else if err.Error.Error() == common.GenerateExpiredToken().Error.Error() {
			return errorModel.GenerateExpiredTokenError(fileName, funcName)
		} else if err.Error.Error() == common.GenerateInvalidMethode().Error.Error() {
			return errorModel.GenerateUnauthorizedClientError(fileName, funcName)
		}
		return errorModel.ErrorModel{
			Code:     err.Code,
			Error:    err.Error,
			FileName: "AbstractEndpoint.go",
			FuncName: "ServeJWTTokenValidationEndpoint",
		}
	} else {
		return errorModel.GenerateNonErrorModel()
	}
}

func GenerateSignature(message string, key string, request *http.Request) map[string]string {
	output := make(map[string]string)
	internalToken := request.Header.Get(constanta.TokenHeaderNameConstanta)
	digest := util.GenerateMessageDigest(message)
	timestamp := time.Now().Format(constanta.DefaultTimeFormat)
	signature := util.GenerateSignature(request.Method, request.RequestURI, internalToken, digest, timestamp, key)

	output[constanta.SignatureHeaderNameConstanta] = signature
	output[constanta.TimestampSignatureHeaderNameConstanta] = timestamp
	return output
}

func generateGetMainVersionApp() int {
	version := strings.Split(config.ApplicationConfiguration.GetServerVersion(), ".")
	versionNum, _ := strconv.Atoi(version[0])
	return versionNum
}

func GetMenuCode(redesignCode, normalCode string) string {
	var version = generateGetMainVersionApp()
	if version == constanta.VersionRedesign {
		return redesignCode
	}
	return normalCode
}
