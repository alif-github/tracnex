package RoleService

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type roleService struct {
	service.AbstractService
	service.GetListData
	service.MultiDeleteData
}

var RoleService = roleService{}.New()

func (input roleService) New() (output roleService) {
	output.FileName = "RoleService.go"
	output.ValidLimit = service.DefaultLimit
	output.ValidOrderBy = []string{
		"id",
		"role_id",
		"description",
		"created_at",
		"created_name",
	}
	output.ValidSearchBy = []string{
		"role_id",
		"description",
	}
	return
}

func (input roleService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.RoleRequest) errorModel.ErrorModel) (inputStruct in.RoleRequest, err errorModel.ErrorModel) {
	funcName := "readBodyAndValidate"
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errorS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errorS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, funcName, errorS)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	inputStruct.ID = int64(id)

	err = validation(&inputStruct)
	return
}

func (input roleService) readParamAndValidate(request *http.Request, validation func(input *in.RoleRequest) errorModel.ErrorModel) (inputStruct in.RoleRequest, err errorModel.ErrorModel) {
	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	inputStruct.ID = int64(id)

	err = validation(&inputStruct)
	return
}

func (input roleService) CheckDuplicateError(err errorModel.ErrorModel) errorModel.ErrorModel {
	if err.CausedBy != nil {
		if service.CheckDBError(err, "uq_role_roleid") {
			return errorModel.GenerateDataUsedError(err.FileName, err.FuncName, constanta.Role)
		}
	}

	return err
}

func (input roleService) removeIndexPerm(permission []string, index int) []string {
	return append(permission[:index], permission[index+1:]...)
}

func (input roleService) checkRoleLimitedByLimitedCreatedBy(contextModel *applicationModel.ContextModel, resultGetOnDB repository.RoleModel) (err errorModel.ErrorModel) {
	fileName := "RoleService.go"
	funcName := "checkRoleLimitedByLimitedCreatedBy"

	// ---------- Check Created By Limited ----------
	if contextModel.LimitedByCreatedBy > 0 && (resultGetOnDB.CreatedBy.Int64 != contextModel.LimitedByCreatedBy) {
		err = errorModel.GenerateForbiddenAccessClientError(fileName, funcName)
		return
	}
	// -----------------------------------------------

	return errorModel.GenerateNonErrorModel()
}
