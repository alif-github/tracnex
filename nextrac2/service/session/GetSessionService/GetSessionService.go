package GetSessionService

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	model2 "nexsoft.co.id/nextrac2/resource_common_service/model"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/session"
	"nexsoft.co.id/nextrac2/util"
	"time"
)

type getSessionService struct {
	service.AbstractService
}

var GetSessionService = getSessionService{}.New()

func (input getSessionService) New() (output getSessionService) {
	output.FileName = "GetSessionService.go"
	return
}

func (input getSessionService) StartService(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		authentication model2.AuthenticationModel
		menuDTOOut     []out.ParentMenu
		personProfile  repository.UserModel
	)

	_ = json.Unmarshal([]byte(contextModel.AuthAccessTokenModel.Authentication), &authentication)
	fmt.Println("Authentication -> ", authentication)
	fmt.Println("ContextModel -> ", contextModel)
	personProfile, err = dao.UserDAO.GetUserByClientID(serverconfig.ServerAttribute.DBConnection, repository.UserModel{
		ClientID: sql.NullString{
			String: contextModel.AuthAccessTokenModel.ClientID,
		},
	})

	if err.Error != nil {
		return
	}

	if personProfile.IsSystemAdmin.Bool {
		var menuDTOOutTemp []out.ParentMenu
		authentication.Role.Permission["nexsoft.home"] = []string{"view"}
		menuDTOOut = []out.ParentMenu{{
			ID:              0,
			Name:            "Home",
			EnName:          "Home",
			Sequence:        0,
			IconName:        "HomeIcon",
			Background:      "",
			AvailableAction: "view",
			MenuCode:        "nexsoft.home",
			ServiceMenu:     nil,
		}}
		menuDTOOutTemp, err = dao.MenuDAO.ViewMenuListAdmin(serverconfig.ServerAttribute.DBConnection)
		menuDTOOut = append(menuDTOOut, menuDTOOutTemp...)
	} else {
		menuDTOOut, err = dao.MenuDAO.ViewMenuList(serverconfig.ServerAttribute.DBConnection)
	}

	if err.Error != nil {
		return
	}

	output.Data.Content = out.GetSessionDTOOut{
		FirstName:      personProfile.FirstName.String,
		LastName:       personProfile.LastName.String,
		Username:       personProfile.Username.String,
		Role:           authentication.Role.Role,
		Locale:         contextModel.AuthAccessTokenModel.Locale,
		UserID:         contextModel.AuthAccessTokenModel.ResourceUserID,
		IdCard:         personProfile.IdCard.String,
		Position:       personProfile.Position.String,
		Department:     personProfile.Department.String,
		IsHaveMember:   personProfile.IsHaveMember.Bool,
		Currency:       personProfile.Currency.String,
		PlatformDevice: personProfile.PlatformDevice.String,
		Scope:          contextModel.AuthAccessTokenModel.Scope,
		CurrentTime:    time.Now().Format(constanta.DefaultTimeFormat),
		Permission:     authentication.Role.Permission,
		Menu:           menuDTOOut,
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: session.GenerateLoginI18NMessage("GET_SESSION_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}

func (input getSessionService) GetCurrentDateTimeService(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {

	output.Data.Content = out.GetSessionDateTimeDTOOut{
		CurrentTime: time.Now().Format(constanta.DefaultTimeFormat),
	}

	output.Status = out.StatusResponse{
		Code:    util.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: session.GenerateLoginI18NMessage("GET_CURRENT_DATETIME_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	return
}
