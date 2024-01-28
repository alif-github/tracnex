package backgroundJobModel

import (
	"database/sql"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
)

type BackgroundServiceModel struct {
	SearchByParam []in.SearchByParam
	IsCheckStatus bool
	CreatedBy     int64
	Data          interface{}
}

type ChildTask struct {
	Group        string
	Type         string
	Name         string
	Data         BackgroundServiceModel
	GetCountData func(*sql.DB, []in.SearchByParam, bool, int64) (int, errorModel.ErrorModel)
	DoJob        func(*sql.DB, interface{}, *repository.JobProcessModel) errorModel.ErrorModel
	DoJobWithCtx func(*sql.DB, interface{}, *repository.JobProcessModel, applicationModel.ContextModel) errorModel.ErrorModel
}
