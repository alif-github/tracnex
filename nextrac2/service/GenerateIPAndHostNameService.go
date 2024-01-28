package service

import (
	"database/sql"
	"fmt"
	"nexsoft.co.id/nexcommon/util"
	"nexsoft.co.id/nextrac2/config"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/repository"
	"os"
	"strconv"
	"time"
)

var GenerateIPAndHostNameService = generateIPAndHostNameService{}.New()

func (input generateIPAndHostNameService) New() (output generateIPAndHostNameService) {
	output.FileName = "GenerateIPAndHostNameService.go"
	return
}

type generateIPAndHostNameService struct {
	AbstractService
	FileName string
}

func (input generateIPAndHostNameService) GenerateIPAndServerID() {
	_, err := input.ServiceWithDataAuditPreparedByService("UpdateDataGroupService", nil, &applicationModel.ContextModel{}, input.generateIPAndServerID, func(_ interface{}, _ applicationModel.ContextModel) {})
	if err.Error != nil {
		return
	}
}

func (input generateIPAndHostNameService) generateIPAndServerID(tx *sql.Tx, _ interface{}, loggerModel *applicationModel.ContextModel, timeNow time.Time) (_ interface{}, dataAudit []repository.AuditSystemModel, err errorModel.ErrorModel) {
	ip, errs := util.GenerateIPAddress(config.ApplicationConfiguration.GetServerEthernet())
	if errs != nil {
		fmt.Println(errs)
		os.Exit(3)
	}

	hostname, errs := util.GenerateHostname()
	if errs != nil {
		fmt.Println(errs)
		os.Exit(3)
	}

	hostServerModel, err := dao.HostServerDAO.GetHostByHostName(tx, repository.HostServerModel{
		HostName: sql.NullString{String: hostname},
	})
	if err.Error != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	protocol := config.ApplicationConfiguration.GetServerProtocol()
	port := config.ApplicationConfiguration.GetServerPort()

	url := protocol + "://" + ip + ":" + strconv.Itoa(port)
	hostServerModel.HostName.String = hostname
	hostServerModel.HostURL.String = url
	hostServerModel.CreatedAt.Time = timeNow
	hostServerModel.UpdatedAt.Time = timeNow

	fmt.Println(fmt.Sprintf(`[HostName] -> %s, [HostURL] -> %s, [HostID] -> %d`, hostname, url, hostServerModel.ID.Int64))
	if hostServerModel.ID.Int64 == 0 {
		var id int64
		id, err = dao.HostServerDAO.InsertHostnameAndIP(tx, hostServerModel)
		if err.Error != nil {
			fmt.Println(err)
			os.Exit(3)
		}
		dataAudit = append(dataAudit, GetAuditData(tx, constanta.ActionAuditInsertConstanta, *loggerModel, timeNow, dao.HostServerDAO.TableName, id, loggerModel.LimitedByCreatedBy)...)
	} else {
		dataAudit = append(dataAudit, GetAuditData(tx, constanta.ActionAuditUpdateConstanta, *loggerModel, timeNow, dao.HostServerDAO.TableName, hostServerModel.ID.Int64, loggerModel.LimitedByCreatedBy)...)
		err = dao.HostServerDAO.UpdateHostnameAndIP(tx, hostServerModel, timeNow)
		if err.Error != nil {
			fmt.Println(err)
			os.Exit(3)
		}
	}

	err = errorModel.GenerateNonErrorModel()
	return
}
