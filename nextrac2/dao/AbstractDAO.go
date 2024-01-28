package dao

import (
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/elasticsearch"
	"nexsoft.co.id/nextrac2/serverconfig"
	"strconv"
)

type AbstractDAO struct {
	FileName           string
	TableName          string
	ElasticSearchIndex string
}

type FieldStatus struct {
	IsCheck   bool
	FieldName string
	Value     interface{}
}

type DefaultFieldMustCheck struct {
	ID        FieldStatus
	Deleted   FieldStatus
	Status    FieldStatus
	CreatedBy FieldStatus
	UpdatedAt FieldStatus
}

func (DefaultFieldMustCheck) GetDefaultField(isCheckStatus bool, createdBy int64) DefaultFieldMustCheck {
	return DefaultFieldMustCheck{
		ID:        FieldStatus{FieldName: "id"},
		Deleted:   FieldStatus{FieldName: "deleted"},
		Status:    FieldStatus{FieldName: "status", IsCheck: isCheckStatus},
		CreatedBy: FieldStatus{FieldName: "created_by", Value: createdBy},
	}
}

func (input AbstractDAO) doUpdateAtElasticSearch(id string, userParam interface{}) {
	go elasticsearch.AsyncUpdateRowByID(serverconfig.ServerAttribute.ElasticClient, config.ApplicationConfiguration.GetServerResourceID()+"."+input.ElasticSearchIndex, id, userParam)
}

func (input AbstractDAO) DoInsertAtElasticSearch(id string, userParam interface{}) {
	go elasticsearch.AsyncInsertRow(serverconfig.ServerAttribute.ElasticClient, config.ApplicationConfiguration.GetServerResourceID()+"."+input.ElasticSearchIndex, id, userParam)
}

func (input AbstractDAO) doDeleteAtElasticSearch(id string) {
	go elasticsearch.AsyncDeleteRowByID(serverconfig.ServerAttribute.ElasticClient, config.ApplicationConfiguration.GetServerResourceID()+"."+input.ElasticSearchIndex, id)
}

func (input AbstractDAO) doGetListDataFromElastic(userParam in.GetListDataDTO, indexName string, searchBy []in.SearchByParam) (resultStr string, resultFoundSize int, err error) {
	mapSearch := SearchByParamToMapSearch(searchBy)

	configSearch := elasticsearch.SearchingConfiguration{
		Offset: CountOffset(userParam.Page, userParam.Limit),
		Limit:  userParam.Limit,
	}

	resultStr, resultFoundSize, err = elasticsearch.GetListDataInMultiFieldsSearchQuery(serverconfig.ServerAttribute.ElasticClient, config.ApplicationConfiguration.GetServerResourceID()+"."+indexName, mapSearch, configSearch)
	if err != nil {
		return
	}

	return
}

func (input AbstractDAO) doGetCountDataFromElastic(searchBy []in.SearchByParam, isCheckStatus bool, createdBy int64) (resultStr string, resultFoundSize int, err error) {
	mapSearch := SearchByParamToMapSearch(searchBy)

	if isCheckStatus {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:   "status",
			SearchValue: "A",
		})
	}
	if createdBy > 0 {
		searchBy = append(searchBy, in.SearchByParam{
			SearchKey:   "created_by",
			SearchValue: strconv.Itoa(int(createdBy)),
		})
	}

	resultStr, resultFoundSize, err = elasticsearch.GetListDataInMultiFieldsSearchQuery(serverconfig.ServerAttribute.ElasticClient, config.ApplicationConfiguration.GetServerResourceID()+"."+input.ElasticSearchIndex, mapSearch, elasticsearch.SearchingConfiguration{})
	if err != nil {
		return
	}

	return
}
