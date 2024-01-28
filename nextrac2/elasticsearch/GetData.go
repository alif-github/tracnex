package elasticsearch

import (
	"context"
	"gopkg.in/olivere/elastic.v7"
)

func GetRowByID(client *elastic.Client, indexName string, id string) (result string, err error) {
	ctx := context.Background()
	get, err := client.Get().Index(indexName).Id(id).Do(ctx)
	if err != nil {
		return
	}
	var byteData []byte
	byteData, err = get.Source.MarshalJSON()
	if err != nil {
		return
	}
	return string(byteData), nil
}
