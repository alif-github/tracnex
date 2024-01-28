package applicationModel

import model2 "nexsoft.co.id/nextrac2/resource_common_service/model"

type ContextModel struct {
	LoggerModel          LoggerModel
	PermissionHave       string
	IsSignatureCheck     bool
	IsInternal           bool
	LimitedByCreatedBy   int64
	AuthAccessTokenModel model2.AuthAccessTokenModel
	DBSchema             string
	IsAdmin              bool
}

type MappingScopeDB struct {
	View      string
	Count     string
	TableName string
}
