package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx"
	inc "nexsoft.co.id/nexcommon/includes"
	"sync"
)

// this structure needs to be defined properly.
// still investigating.
type InsertData struct {
	Id int64
}

type InsertRc struct {
	Rc   inc.Rc
	Data InsertData
}

type DbRowMetaData struct {
	Created_by      int64  `json:"created_by"`
	Created_at      string `json:"created_at"`
	Last_updated_by int64  `json:"last_updated_by"`
	Last_updated_at string `json:"last_updated_at"`
	Record_status   string `json:"record_status"`
}

// =============================================================================
// =============================================================================
// =============================================================================
// =============================================================================

type nex_dbInfo struct {
	instance      *sql.DB
	driver        string
	connectionStr string
	userid        string
	password      string
	DBName        string
	schema        string
	setParams     []string
}

var instance *sql.DB
var once sync.Once

func DB_returnedOk(_rc inc.Rc) bool {
	if _rc.Code == "OK" || _rc.Code == "" {
		return true
	} else {
		return false
	}
}

// =============================================================================
//
// =============================================================================
func GetDbConnection(target string) (*sql.DB, error) {
	fmt.Println("GetDbConnection: entry")
	defer fmt.Println("GetDbConnection: exit")

	_params := []string{"set search_path='nex_api_interface'"}
	_dbInfo := nex_dbInfo{nil, "postgres",
		"user=postgres password=Pa55w0rd dbname=postgres sslmode=disable",
		"", "", "", "", _params}
	_db, _err := getInstance(_dbInfo)
	return _db, _err
}

// =============================================================================
//
// =============================================================================
func getInstance(connInfo nex_dbInfo) (*sql.DB, error) {
	fmt.Println("getInstance: entry")
	var _errOpen error
	defer fmt.Println("getInstance: exit")
	once.Do(func() {
		instance, _errOpen = sql.Open(connInfo.driver, connInfo.connectionStr)

		if _errOpen == nil {
			for _, _param := range connInfo.setParams {
				instance.Exec(_param)
			}
		} else {
			fmt.Println(_errOpen)
			instance = nil
		}
	})
	return instance, _errOpen
}
