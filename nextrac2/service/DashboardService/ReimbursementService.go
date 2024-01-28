package DashboardService

import (
	"net/http"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"time"
)

type reimbursementService struct {
	service.AbstractService
	service.GetListData
}

var ReimbursementService = reimbursementService{}.New()

func (input reimbursementService) New() (output reimbursementService) {
	output.FileName = "ReimbursementService.go"
	output.ServiceName = "Reimbursement Dashboard"
	return
}

func (input reimbursementService) ViewDashboardCountReimbursement(_ *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		db         = serverconfig.ServerAttribute.DBConnection
		timeNow    = time.Now()
		outputTemp = make(map[string]interface{})
		count      int64
	)

	count, err = dao.DashboardDAO.GetReimbursementInMonth(db, timeNow)
	if err.Error != nil {
		return
	}

	//--- Output Content
	outputTemp["reimbursement"] = count
	output.Data.Content = outputTemp

	//--- Output Message
	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}
