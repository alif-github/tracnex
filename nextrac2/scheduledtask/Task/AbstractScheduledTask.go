package Task

import (
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/serverconfig"
	"os"
)

type AbstractScheduledTask struct {
	RunType string
}

func (input AbstractScheduledTask) StartTask(runType string, cmd func()) {
	hostName, _ := os.Hostname()
	serverRunModel, err := dao.SchedulerDAO.GetServerRunByHostNameAndRunType(serverconfig.ServerAttribute.DBConnection, hostName, runType)
	if err.Error != nil {
		//todo log error
		return
	}

	if serverRunModel.ID.Int64 != 0 {
		cmd()
	}
}
