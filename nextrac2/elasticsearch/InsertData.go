package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/olivere/elastic.v7"
	util2 "nexsoft.co.id/nexcommon/util"
)

func AsyncInsertRow(client *elastic.Client, indexName string, id string, newValue interface{}) {
	err := InsertRow(client, indexName, id, util2.StructToJSON(newValue))
	if err != nil {
		fmt.Print(err)
		//todo Async insert to db
	}
}

func InsertRow(client *elastic.Client, indexName string, id string, body string) error {
	ctx := context.Background()
	ind, err := client.Index().Index(indexName).Id(id).BodyJson(body).Do(ctx)
	if err != nil {
		return err
	}

	if ind.Status != 0 {
		return errors.New(ind.Result)
	}

	return nil
}
