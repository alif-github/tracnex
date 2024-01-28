package out

import "time"

type ReportResponse struct {
	NIK           int64   `json:"nik"`
	Name          string  `json:"name"`
	Department    string  `json:"department"`
	BacklogManday float64 `json:"backlog_manday"`
	ActualManday  float64 `json:"actual_manday"`
	Tracker       string  `json:"tracker"`
	MandayRate    float64 `json:"manday_rate"`
	Manday        float64 `json:"manday"`
	BacklogTicket string  `json:"ticket_backlog"`
	ActualTicket  string  `json:"actual_ticket"`
}

type ReportHistoryResponse struct {
	ID                int64     `json:"id"`
	Department        string    `json:"department"`
	SuccessTicket     string    `json:"success_ticket"`
	PaymentDate       time.Time `json:"payment_date"`
	PersonResponsible string    `json:"person_responsible"`
}

type ViewReportHistoryResponse struct {
	ID                int64                 `json:"id"`
	Department        string                `json:"department"`
	SuccessTicket     string                `json:"success_ticket"`
	PaymentDate       time.Time             `json:"payment_date"`
	PersonResponsible string                `json:"person_responsible"`
	Data              ResultsReportResponse `json:"data"`
}

type ResultsReportResponse struct {
	Results []PersonReportResponse `json:"results"`
}

type PersonReportResponse struct {
	NIK         int64                   `json:"nik"`
	Name        string                  `json:"name"`
	Department  string                  `json:"department"`
	TotalManday float64                 `json:"total_manday"`
	Detail      []DetailMandayByTracker `json:"detail"`
}

type DetailMandayByTracker struct {
	Tracker       string  `json:"tracker"`
	MandayRate    float64 `json:"manday_rate"`
	BacklogManday float64 `json:"backlog_manday"`
	TicketBacklog string  `json:"backlog_ticket"`
	ActualManday  float64 `json:"actual_manday"`
	TicketActual  string  `json:"actual_ticket"`
	Manday        float64 `json:"manday"`
}

type TrackerDeveloper struct {
	Task float64 `json:"task"`
}

type TrackerQA struct {
	Automation float64 `json:"automation"`
	Manual     float64 `json:"manual"`
}

type TrackerInfraDevOps struct {
	Task float64 `json:"task"`
}
