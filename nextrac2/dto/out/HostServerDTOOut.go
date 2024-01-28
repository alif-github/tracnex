package out

import "time"

type HostServerResponse struct {
	ID        int64     `json:"id"`
	HostName  string    `json:"host_name"`
	CreatedBy int64     `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ViewHostServerResponse struct {
	Host          Host          `json:"host"`
	ListScheduler []ListScheduler `json:"list_scheduler"`
}

type Host struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type ListScheduler struct {
	Name      string `json:"name"`
	Cron      string `json:"cron"`
	RunStatus bool   `json:"run_status"`
}
