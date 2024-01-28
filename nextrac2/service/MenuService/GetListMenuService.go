package MenuService

import (
	"github.com/gorilla/mux"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	util2 "nexsoft.co.id/nextrac2/util"
	"strconv"
)

func (input menuService) GetListParentMenu(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	output.Data.Content, err = input.doViewMenuFull(false)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_PARENT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input menuService) GetListParentMenuSysadmin(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	output.Data.Content, err = input.doViewMenuFull(true)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_PARENT_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input menuService) removeIndexMenu(menuList []out.MenuList, index int) []out.MenuList {
	return append(menuList[:index], menuList[index+1:]...)
}

func (input menuService) doViewParentAndServe(isSysadmin bool) (output interface{}, err errorModel.ErrorModel) {
	var menuDTOOut []out.MenuList
	var menuResult []out.MenuList
	var globalMenu out.MenuList

	menuResult, err = dao.MenuDAO.ViewParentMenuList(serverconfig.ServerAttribute.DBConnection, isSysadmin)
	if err.Error != nil {
		return
	}

	if !isSysadmin {
		for i := 0; i < len(menuResult); i++ {
			if menuResult[i].Name == "Admin" {
				menuResult = input.removeIndexMenu(menuResult, i)
				break
			}
		}
	}

	globalMenu.Name = "Global"
	globalMenu.EnName = "Global"
	globalMenu.AvailableAction = input.generateGlobalPermission()

	menuDTOOut = append(menuDTOOut, globalMenu)
	menuDTOOut = append(menuDTOOut, menuResult...)

	output = menuDTOOut
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input menuService) GetListServiceMenu(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "GetListServiceMenu"

	var MenuDTOOut []out.MenuList

	menuID, _ := strconv.Atoi(mux.Vars(request)["ID"])

	if menuID < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.ServiceMenu)
		return
	}

	MenuDTOOut, err = dao.MenuDAO.ViewServiceMenuList(serverconfig.ServerAttribute.DBConnection, menuID)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_SERVICE_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	output.Data.Content = MenuDTOOut
	return
}

func (input menuService) GetListMenuItem(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	funcName := "GetListMenuItem"

	var MenuDTOOut []out.NewMenuItemList

	menuID, _ := strconv.Atoi(mux.Vars(request)["ID"])

	if menuID < 1 {
		err = errorModel.GenerateUnknownDataError(input.FileName, funcName, constanta.MenuItem)
		return
	}

	MenuDTOOut, err = dao.MenuDAO.ViewMenuItemList(serverconfig.ServerAttribute.DBConnection, menuID)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: GenerateI18NMessage("SUCCESS_GET_ITEM_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	output.Data.Content = MenuDTOOut
	return
}

func (input menuService) doViewMenuFull(isSysadmin bool) (output interface{}, err errorModel.ErrorModel) {
	var (
		menuDTOOut []out.ParentMenuList
		menuResult []out.ParentMenuList
		globalMenu out.ParentMenuList
	)

	if !isSysadmin {
		menuResult, err = dao.MenuDAO.ViewMenuListForPermission(serverconfig.ServerAttribute.DBConnection)
		if err.Error != nil {
			return
		}

		for i := 0; i < len(menuResult); i++ {
			if menuResult[i].Name == "Admin" {
				menuResult = append(menuResult[:i], menuResult[i+1:]...)
				break
			}
		}
	} else {
		menuResult, err = dao.MenuDAO.ViewMenuListAdminForPermission(serverconfig.ServerAttribute.DBConnection)
		if err.Error != nil {
			return
		}
	}

	globalMenu.Name = "Global"
	globalMenu.EnName = "Global"
	globalMenu.AvailableAction = input.generateGlobalPermission()

	menuDTOOut = append(menuDTOOut, globalMenu)
	menuDTOOut = append(menuDTOOut, menuResult...)

	output = menuDTOOut
	err = errorModel.GenerateNonErrorModel()
	return
}

func (input menuService) generateGlobalPermission() string {
	result :=
		InsertDataGlobalPermission + ", " + ViewDataGlobalPermission + ", " + UpdateDataGlobalPermission + ", " +
			DeleteDataPermission + ", " + ChangePaswwordDataGlobalPermission
	return result
}
