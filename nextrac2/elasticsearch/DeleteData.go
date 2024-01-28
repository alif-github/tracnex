package elasticsearch

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/olivere/elastic.v7"
)

func AsyncDeleteRowByID(client *elastic.Client, indexName string, id string) {
	err := DeleteRowByID(client, indexName, id)
	if err != nil {
		fmt.Print(err)
		//todo Async insert to db
	}
}

func DeleteRowByID(client *elastic.Client, indexName string, id string) error {
	ctx := context.Background()
	update, err := client.Delete().Index(indexName).Id(id).Do(ctx)
	if err != nil {
		return err
	}

	if update.Status != 0 {
		return errors.New(update.Result)
	}

	return nil
}
