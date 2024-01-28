package CRUDUserService

import (
	"database/sql"
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"nexsoft.co.id/nextrac2/resource_common_service"
	"nexsoft.co.id/nextrac2/serverconfig"
	userService2 "nexsoft.co.id/nextrac2/service/UserService"
	"nexsoft.co.id/nextrac2/token"
	util2 "nexsoft.co.id/nextrac2/util"
	"strings"
)

func (input userService) getListUserActiveFromRedis() (payload []token.PayloadJWTToken, err errorModel.ErrorModel) {
	var (
		funcName  = "getListUserActiveFromRedis"
		keyMatch  = constanta.SessionUser + "*"
		jwtToken  string
		listToken []string
	)

	// get sys-user token from redis
	iter := serverconfig.ServerAttribute.RedisClientSession.Scan(0, keyMatch, 0).Iterator()
	for iter.Next() {
		plainTokenFromRedis := iter.Val()
		jwtToken = strings.Replace(plainTokenFromRedis, constanta.SessionUser, " ", 1)
		listToken = append(listToken, jwtToken)
	}

	if errs := iter.Err(); errs != nil {
		err = errorModel.GenerateRedisError(input.FileName, funcName)
		return
	}

	// unmarshal payload
	for _, tokenIndex := range listToken {
		payloadJWT, errs := resource_common_service.ConvertJWTToPayload(tokenIndex)
		if errs.Error != nil {
			return
		}

		fmt.Println("Token JWT : ", payloadJWT)
		fmt.Println("ClientID Token JWT : ", payloadJWT.ClientID)
		payload = append(payload, payloadJWT)
	}
	fmt.Println("Get List User Active From Redis : ", payload)

	return
}

func (input userService) GetListUserSysUserActive(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		searchByParam []in.SearchByParam
		validOrderBy  = []string{"username"}
		validSearchBy = []string{"username"}
	)

	_, searchByParam, err = input.ReadAndValidateGetListData(request, validSearchBy, validOrderBy, applicationModel.GetListUserActiveValidOperator, input.ValidLimit)
	if err.Error != nil {
		return
	}

	output.Data.Content, err = input.doGetListUserActive(searchByParam)
	if err.Error != nil {
		return
	}

	output.Status = out.StatusResponse{
		Code:    util2.GenerateConstantaI18n("SUCCESS", contextModel.AuthAccessTokenModel.Locale, nil),
		Message: userService2.GenerateI18NMessage("SUCCESS_LIST_MESSAGE", contextModel.AuthAccessTokenModel.Locale),
	}

	err = errorModel.GenerateNonErrorModel()
	return
}

func (input userService) doGetListUserActive(searchByParam []in.SearchByParam) (output []out.SessionLoginResponse, err errorModel.ErrorModel) {
	var (
		listUserActive     []token.PayloadJWTToken
		listResponse       []out.SessionLoginResponse
		listResponseUnique []out.SessionLoginResponse
	)

	listUserActive, err = input.getListUserActiveFromRedis()
	if err.Error != nil {
		return
	}

	// get all detail user from userActive
	for _, userActive := range listUserActive {
		userOnDB, errs := dao.UserDAO.GetDetailUserForCheckSessionSysuser(serverconfig.ServerAttribute.DBConnection, repository.UserModel{ClientID: sql.NullString{String: userActive.ClientID}})
		if errs.Error != nil {
			return
		}

		if userOnDB.Username.String != "" {
			listResponse = append(listResponse, out.SessionLoginResponse{
				Username: userOnDB.Username.String,
				Phone:    userOnDB.Phone.String,
				Email:    userOnDB.Email.String,
			})
		}
	}

	if searchByParam == nil {
		for _, responseUnique := range listResponse {
			if !contains(listResponseUnique, responseUnique) {
				listResponseUnique = append(listResponseUnique, responseUnique)
			}
		}

		output = listResponseUnique
		return
	}

	output = listResponse
	return
}

func contains(elements []out.SessionLoginResponse, value out.SessionLoginResponse) bool {
	for _, s := range elements {
		if value == s {
			return true
		}
	}

	return false
}
