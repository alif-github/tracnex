package Task

import "fmt"

type syncElasticBank struct {
	AbstractScheduledTask
}

var SyncElasticBank = syncElasticBank{}.New()
var SyncElasticBank2 = syncElasticBank{AbstractScheduledTask{RunType: "scheduler.sync_data_bank2"}}

func (input syncElasticBank) New() (output syncElasticBank) {
	output.RunType = "scheduler.sync_data_bank"
	return
}

func (input syncElasticBank) Start() {
	input.StartTask(input.RunType, input.syncElasticBank)
}

func (input syncElasticBank) syncElasticBank() {
	//ElasticSearchSyncService.ElasticSearchSyncService.DoSyncElasticBank(applicationModel.ContextModel{
	//	AuthAccessTokenModel: model.AuthAccessTokenModel{
	//		RedisAuthAccessTokenModel:  model.RedisAuthAccessTokenModel{ResourceUserID: 1},
	//		AuthenticationServerUserID: 1,
	//	},
	//	DBSchema: config.ApplicationConfiguration.GetPostgreSQLDefaultSchema(),
	//})
	//fmt.Println("Start Scheduler")
	//timeNow := time.Now()
	//fmt.Println(timeNow.String())
	fmt.Println("Test Scheduller")
}
