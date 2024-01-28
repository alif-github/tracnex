package elasticsearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/olivere/elastic.v7"
	"nexsoft.co.id/nexcommon/util"
)

func AsyncUpdateRowByID(client *elastic.Client, indexName string, id string, newValue interface{}) {
	err := UpdateRowByID(client, indexName, id, util.StructToJSON(newValue))
	if err != nil {
		fmt.Print(err)
		//todo Async insert to db
	}
}

func UpdateRowByID(client *elastic.Client, indexName string, id string, newValue string) error {
	ctx := context.Background()

	var upsert map[string]interface{}
	_ = json.Unmarshal([]byte(newValue), &upsert)

	update, err := client.Update().Index(indexName).Id(id).Doc(upsert).Do(ctx)
	if err != nil {
		return err
	}

	if update.Status != 0 {
		return errors.New(update.Result)
	}

	return nil
}
