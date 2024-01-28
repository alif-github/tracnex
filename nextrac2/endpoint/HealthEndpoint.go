package endpoint

import (
	"fmt"
	"net/http"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/model"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/serverconfig"
)

type healthEndpoint struct {
	AbstractEndpoint
}

var HealthEndpoint healthEndpoint


func (hs healthEndpoint) GetHealthStatus(responseWriter http.ResponseWriter, _ *http.Request) {
	//funcName := "GetHealthStatus"
	HealthEndpoint.FileName = "HealthEndpoint.go"
	responseWriter.WriteHeader(200)
	status := "UP"
	statusRedis := "UP"
	statusDB := "UP"
	//statusElastic := "UP"

	errRedis := serverconfig.ServerAttribute.RedisClient.Ping()
	if errRedis.Err() != nil {
		statusRedis = "DOWN"
	}

	//health := serverconfig.ServerAttribute.ElasticClient.ClusterHealth()
	//stat, errElastic := health.Do(context.Background())
	//if errElastic != nil {
	//	return
	//}
	//statusElastic = stat.Status

	errDb := serverconfig.ServerAttribute.DBConnection.Ping()
	if errDb != nil {
		statusDB = "DOWN"
	}

	serverStatus := applicationModel.ServerStatus{
		Status:        status,
		Redis:         statusRedis,
		Database:      statusDB,
		//ElasticSearch: statusElastic,
	}

	_, err := responseWriter.Write([]byte(util.StructToJSON(serverStatus)))
	if err != nil {
		logModel := model.GenerateLogModel(config.ApplicationConfiguration.GetServerVersion())
		logModel.Status = 500
		logModel.Message = fmt.Sprintf(err.Error())
		util.LogInfo(logModel.ToLoggerObject())
	}
}