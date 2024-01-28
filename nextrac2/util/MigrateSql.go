package util

import (
	"database/sql"
	"errors"
	"github.com/gobuffalo/packr/v2"
	migrate "github.com/rubenv/sql-migrate"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"strconv"
)

func RollbackSchema(db *sql.DB) (err errorModel.ErrorModel) {
	fileName := "MigrateSql.go"
	funcName := "RollbackSchema"

	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations", "../../sql_migrations"),
	}

	if db != nil {
		n, errs := migrate.Exec(db, "postgres", migrations, migrate.Down)
		if errs != nil {
			logModel := applicationModel.GenerateLogModel("-", config.ApplicationConfiguration.GetServerResourceID())
			logModel.Message = errs.Error()
			logModel.Status = 500
			util.LogError(logModel.ToLoggerObject())
			err = errorModel.GenerateInternalDBServerError(fileName, funcName, errs)
		} else {
			logModel := applicationModel.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion(), config.ApplicationConfiguration.GetServerResourceID())
			logModel.Status = 200
			logModel.Message = "Dropped " + strconv.Itoa(n) + " migrations!"
			util.LogInfo(logModel.ToLoggerObject())
		}
	} else {
		logModel := applicationModel.GenerateLogModel("-", config.ApplicationConfiguration.GetServerResourceID())
		logModel.Message = "null database"
		logModel.Status = 500
		util.LogError(logModel.ToLoggerObject())
		err = errorModel.GenerateInternalDBServerError(fileName, funcName, errors.New("null database"))
	}
	return
}