package out

import "time"

type CronSchedulerResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	RunType   string    `json:"run_type"`
	CreatedBy int64     `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    bool      `json:"status"`
}

type ViewCronSchedulerResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	RunType   string    `json:"run_type"`
	Cron      string    `json:"cron"`
	CreatedBy int64     `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    bool      `json:"status"`
}
