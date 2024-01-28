package elasticsearch

import (
	"context"
	"encoding/json"
	"gopkg.in/olivere/elastic.v7"
	"nexsoft.co.id/nexcommon/util"
	"strconv"
)

type SearchingConfiguration struct {
	DescendingSort bool
	Limit          int
	Offset         int
	Pretty         bool
}

func GetListDataEqualsToValuesQuery(client *elastic.Client, indexName string, fieldName string, fieldValuesSearch []interface{}, searchingConfiguration SearchingConfiguration) (result string, resultFoundSize int, err error) {
	ctx := context.Background()
	termQuery := elastic.NewTermsQuery(fieldName, fieldValuesSearch...)
	return startSearch(ctx, client, indexName, termQuery, searchingConfiguration)
}

func GetListDataLikeValueSearchQuery(client *elastic.Client, indexName string, fieldName string, fieldValueSearch string, searchingConfiguration SearchingConfiguration) (result string, resultFoundSize int, err error) {
	ctx := context.Background()
	termQuery := elastic.NewWildcardQuery(fieldName, fieldValueSearch)
	return startSearch(ctx, client, indexName, termQuery, searchingConfiguration)
}

func GetListDataInMultiFieldsSearchQuery(client *elastic.Client, indexName string, fieldSearch map[string]string, searchingConfiguration SearchingConfiguration) (result string, resultFoundSize int, err error) {
	ctx := context.Background()
	var QuerySearch []elastic.Query
	for key := range fieldSearch {
		var termQuery elastic.Query
		var dataInt int

		dataInt, err = strconv.Atoi(fieldSearch[key])
		if err != nil {
			termQuery = elastic.NewMatchPhraseQuery(key, fieldSearch[key])
		} else {
			termQuery = elastic.NewBoolQuery().Filter(elastic.NewRangeQuery(key).
				From(dataInt).
				To(dataInt))
		}

		QuerySearch = append(QuerySearch, termQuery)
	}
	termQuery := elastic.NewBoolQuery().Must(QuerySearch...)
	return startSearch(ctx, client, indexName, termQuery, searchingConfiguration)
}

func startSearch(ctx context.Context, client *elastic.Client, indexName string, query elastic.Query, searchingConfiguration SearchingConfiguration) (result string, resultFoundSize int, err error) {
	search := client.Search().Index(indexName).Query(query).From(searchingConfiguration.Offset)

	if searchingConfiguration.Limit > 0 {
		search = search.Size(searchingConfiguration.Limit)
	}
	if searchingConfiguration.Pretty {
		search = search.Pretty(true)
	}

	searchResult, err := search.Do(ctx)
	if err != nil {
		return
	}

	resultFoundSize = int(searchResult.Hits.TotalHits.Value)
	if searchResult.Hits.TotalHits.Value > 0 {
		var arrayResult []interface{}
		for _, hit := range searchResult.Hits.Hits {
			var byteData []byte
			var temp interface{}
			byteData, err = hit.Source.MarshalJSON()
			if err != nil {
				return
			}
			_ = json.Unmarshal(byteData, &temp)

			arrayResult = append(arrayResult, temp)
		}
		result = util.StructToJSON(arrayResult)
	}

	return
}
