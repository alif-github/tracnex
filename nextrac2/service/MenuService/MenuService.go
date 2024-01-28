package MenuService

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"strconv"
)

type menuService struct {
	service.GetListData
	service.AbstractService
}

var MenuService = menuService{}.New()

func (input menuService) New() (output menuService) {
	output.FileName = "MenuService.go"
	return
}

func (input menuService) readBodyAndValidate(request *http.Request, contextModel *applicationModel.ContextModel, validation func(input *in.MenuRequest) errorModel.ErrorModel) (inputStruct in.MenuRequest, err errorModel.ErrorModel) {
	var stringBody string

	stringBody, err = input.ReadBody(request, contextModel)
	if err.Error != nil {
		return
	}

	errS := json.Unmarshal([]byte(stringBody), &inputStruct)
	if errS != nil {
		err = errorModel.GenerateInvalidRequestError(input.FileName, "readBodyAndValidate", errS)
		return
	}

	id, _ := strconv.Atoi(mux.Vars(request)["ID"])
	if inputStruct.ID == 0 {
		inputStruct.ID = int64(id)
	}

	err = validation(&inputStruct)
	return
}

func (input menuService) doCheckParentMenu(parentMenuID int64, tableName string) (err errorModel.ErrorModel) {
	fileName := "UpdateMenuService.go"
	funcName := "doCheckParentMenu"
	var menuParentOnDB repository.MenuModel

	menuParentOnDB, err  = dao.MenuDAO.GetMenuForUpdate(serverconfig.ServerAttribute.DBConnection, repository.MenuModel{
		ID: 		sql.NullInt64{Int64: parentMenuID},
		CreatedBy:	sql.NullInt64{Int64: 0},
	}, tableName)
	if err.Error != nil {return}

	if menuParentOnDB.ID.Int64 < 1 {
		err = errorModel.GenerateUnknownDataError(fileName, funcName, "id " + tableName)
		return
	}

	err = errorModel.GenerateNonErrorModel()
	return
}