package NexsoftRoleService

import (
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/util"
	"strconv"
	"strings"
)

func GenerateI18NMessage(messageID string, language string) (output string) {
	return util.GenerateI18NServiceMessage(serverconfig.ServerAttribute.NexsoftRoleServiceBundle, messageID, language, nil)
}

func readPathParam(request *http.Request) (id int64, err errorModel.ErrorModel) {
	funcName := "readPathParam"

	strId, ok := mux.Vars(request)["ID"]
	idParam, errConvert := strconv.Atoi(strId)
	id = int64(idParam)

	if !ok || errConvert != nil {
		err = errorModel.GenerateUnsupportedRequestParam("ComplaintSubServiceUtil.go", funcName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func checkDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_nexsoft_role_roleid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.Role)
		}
	}

	return err
}

func getDefaultPermission(inputStruct in.RoleRequest) (output []string) {
	var (
		isGetDefaultPermission bool
		defaultPermission      []string
	)

	defaultPermission = []string{
		"nexsoft.province:view",
		"nexsoft.district:view",
		"nexsoft.customer-group:view",
		"nexsoft.customer-category:view",
		"nexsoft.product-group:view",
		"nexsoft.client-type:view",
		"nexsoft.profile-setting:view-own",
		"nexsoft.profile-setting:update-own",
		"nexsoft.audit-monitoring:view",
		"nexsoft.change-password:update-own",
	}

	for _, permission := range inputStruct.Permission {
		splitPermission := strings.Split(permission, ":")
		if splitPermission[0] == "nexsoft.data-group" {
			if splitPermission[1] == "insert" || splitPermission[1] == "update" {
				isGetDefaultPermission = true
				break
			}
		}
	}

	if isGetDefaultPermission {
		for _, permission := range inputStruct.Permission {
			for j, defPermission := range defaultPermission {
				if permission == defPermission {
					defaultPermission = append(defaultPermission[:j], defaultPermission[j+1:]...)
					break
				}
			}
		}
	} else {
		defaultPermission = nil
	}

	return defaultPermission
}
