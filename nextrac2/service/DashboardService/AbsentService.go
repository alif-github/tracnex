package DashboardService

import (
	"fmt"
	"net/http"
	"nexsoft.co.id/nextrac2/constanta"
	"nexsoft.co.id/nextrac2/dao"
	"nexsoft.co.id/nextrac2/dto/in"
	"nexsoft.co.id/nextrac2/dto/out"
	"nexsoft.co.id/nextrac2/model/applicationModel"
	"nexsoft.co.id/nextrac2/model/errorModel"
	"nexsoft.co.id/nextrac2/serverconfig"
	"nexsoft.co.id/nextrac2/service"
	"nexsoft.co.id/nextrac2/service/AbsentService"
	"time"
)

type absentDashboardService struct {
	service.AbstractService
	service.GetListData
}

var AbsentDashboardService = absentDashboardService{}.New()

func (input absentDashboardService) New() (output absentDashboardService) {
	output.FileName = "AbsentDashboardService.go"
	output.ServiceName = "Absent Dashboard"
	return
}

func (input absentDashboardService) ViewAbsentAverage(request *http.Request, contextModel *applicationModel.ContextModel) (output out.Payload, header map[string]string, err errorModel.ErrorModel) {
	var (
		db             = serverconfig.ServerAttribute.DBConnection
		searchBy       = []string{"period"}
		absentValidOpt map[string]applicationModel.DefaultOperator
		outputTemp     map[string]interface{}
		searchByParam  []in.SearchByParam
		average        float64
		inputStruct    *in.GetListDataDTO
	)

	absentValidOpt = make(map[string]applicationModel.DefaultOperator)
	absentValidOpt["period"] = applicationModel.DefaultOperator{
		DataType: "char",
		Operator: []string{"eq", "like"},
	}

	filter := service.GenerateQueryValue(request.URL.Query()["filter"])
	inputStruct = &in.GetListDataDTO{
		Filter: filter,
	}

	searchByParam, err = inputStruct.ValidateGetCountData(searchBy, absentValidOpt)
	if err.Error != nil {
		return
	}

	err = AbsentService.AbsentService.PeriodCheck(&searchByParam)
	if err.Error != nil {
		return
	}

	average, err = dao.DashboardDAO.GetAbsentInLastPeriod(db, searchByParam)
	if err.Error != nil {
		return
	}

	var (
		s, s1, e, e1   string
		st, et         time.Time
		dateView, date string
	)

	for _, item := range searchByParam {
		if item.SearchKey == "period_start" {
			st, _ = time.Parse(constanta.DefaultTimeFormat, item.SearchValue)
			s = st.Format(`02/01/2006`)
			s1 = st.Format(constanta.DefaultTimeFormatForFile)
		} else if item.SearchKey == "period_end" {
			et, _ = time.Parse(constanta.DefaultTimeFormat, item.SearchValue)
			e = et.Format(`02/01/2006`)
			e1 = et.Format(constanta.DefaultTimeFormatForFile)
		}
	}

	//--- Date View
	dateView = fmt.Sprintf(`%s-%s`, s, e)
	date = fmt.Sprintf(`%s-%s`, s1, e1)

	//--- Output Content
	outputTemp = make(map[string]interface{})
	outputTemp["absent"] = average
	outputTemp["period"] = date
	outputTemp["view"] = dateView
	output.Data.Content = outputTemp

	//--- Output Message
	output.Status = input.GetResponseMessage("SUCCESS_VIEW_MESSAGE", contextModel)
	err = errorModel.GenerateNonErrorModel()
	return
}
